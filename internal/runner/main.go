package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Swapica/order-aggregator-svc/resources"
	resources3 "github.com/Swapica/relayer-svc/resources"
	"github.com/Swapica/swapica-svc/internal/config"
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/internal/proxy"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/enums"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	resources2 "github.com/Swapica/swapica-svc/resources"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gorilla/websocket"
	jsonapi "gitlab.com/distributed_lab/json-api-connector"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/connectors/signed"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
)

type Runner struct {
	socketUrl  string
	aggregator *jsonapi.Connector
	relayer    *jsonapi.Connector
	log        *logan.Entry
	proxyRepo  proxy.ProxyRepo
	chains     data.ChainsQ
	tokens     data.TokensQ
	useRelayer bool
}

func NewRunner(cfg config.Config, proxyRepo proxy.ProxyRepo, chains data.ChainsQ, tokens data.TokensQ) (*Runner, error) {

	aggregator := jsonapi.NewConnector(
		signed.NewClient(
			&http.Client{Timeout: cfg.RunnerCfg().Timeout},
			cfg.RunnerCfg().Aggregator.Endpoint,
		),
	)
	relayer := jsonapi.NewConnector(
		signed.NewClient(
			&http.Client{Timeout: cfg.RunnerCfg().Timeout},
			cfg.RunnerCfg().RelayerEndpoint,
		),
	)

	runner := &Runner{
		socketUrl:  cfg.RunnerCfg().Aggregator.Ws,
		aggregator: aggregator,
		relayer:    relayer,
		log:        cfg.Log(),
		proxyRepo:  proxyRepo,
		chains:     chains,
		tokens:     tokens,
		useRelayer: cfg.RunnerCfg().UseRelayer,
	}

	return runner, nil
}

func (r *Runner) fetchOrders() (orders []resources.Order, err error) {
	var limit = 100
	var response resources.OrderListResponse
	val, _ := url.ParseQuery(fmt.Sprintf("filter[use_relayer]=true&filter[state]=1&page[number]=0&page[limit]=%d", limit))
	for {
		u, _ := url.Parse("/orders?" + val.Encode())
		err = r.aggregator.Get(u, &response)
		if err != nil {
			return nil, err
		}

		orders = append(orders, response.Data...)
		if len(response.Data) < limit || response.Links == nil {
			break
		}

		nextUrl, _ := url.Parse(response.Links.Next)
		val = nextUrl.Query()
	}
	return orders, nil
}

func (r *Runner) fetchMatches() (matches []resources.Match, err error) {
	var limit = 100
	var response resources.MatchListResponse
	val, _ := url.ParseQuery(fmt.Sprintf("filter[use_relayer]=true&filter[state]=2&page[number]=0&page[limit]=%d", limit))
	for {
		u, _ := url.Parse("/match_orders?" + val.Encode())
		err = r.aggregator.Get(u, &response)
		if err != nil {
			return nil, errors.Wrap(err, "failed to send get request")
		}

		matches = append(matches, response.Data...)
		if len(response.Data) < limit || response.Links == nil {
			break
		}

		nextUrl, _ := url.Parse(response.Links.Next)
		val = nextUrl.Query()
	}
	return matches, nil
}

func (r *Runner) ExecuteOrders() error {
	orders, err := r.fetchOrders()
	if err != nil {
		return errors.Wrap(err, "failed to get order list from aggregator")
	}

	for _, order := range orders {
		_ = r.executeMatch(order)
	}

	matches, err := r.fetchMatches()
	if err != nil {
		return errors.Wrap(err, "failed to get match order list from aggregator")
	}

	for _, match := range matches {
		err = r.executeOrder(match)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Runner) getChains(srcId, destId string) (src *data.Chain, dest *data.Chain, err error) {
	src, err = r.chains.New().FilterByID(srcId).Get()
	if err != nil {
		return
	}

	if src == nil {
		src, err = r.getChainByParams(srcId)
		if err != nil {
			return
		}

		if src == nil {
			return nil, nil, errors.New("src chain not found")
		}
	}

	dest, err = r.chains.New().FilterByID(destId).Get()
	if dest == nil {
		dest, err = r.getChainByParams(destId)
		if err != nil {
			return
		}

		if dest == nil {
			return nil, nil, errors.New("dest chain not found")
		}
	}
	return
}

type ChainParams struct {
	ChainId uint64 `json:"chain_id"`
}

func (r *Runner) getChainByParams(chainId string) (ch *data.Chain, err error) {
	allChains, err := r.chains.Select()
	if err != nil {
		return
	}

	for _, chain := range allChains {
		var cainParams ChainParams
		err = json.Unmarshal(chain.ChainParams, &cainParams)
		if err != nil {
			return
		}

		if chainId == strconv.FormatUint(cainParams.ChainId, 10) {
			ch, err = r.chains.New().FilterByID(chain.ID).Get()
			if err != nil {
				return
			}
		}
	}
	return
}

func (r *Runner) sendTxToRelayer(chainId string, data interface{}, token common.Address,
	receiver common.Address, contractAddress common.Address, amount string, rpc string) error {
	u, _ := url.Parse("/transaction")
	evmTx, ok := data.(resources2.EvmTransaction)
	if !ok {
		r.log.Debug(fmt.Sprintf("expected type EvmTransaction, got %T", data))
		return errors.New(fmt.Sprintf("invalid tx data format"))
	}

	executeTxData, err := hexutil.Decode(evmTx.Attributes.TxBody.Data)
	if err != nil {
		r.log.Debug(fmt.Sprintf("err in decode tx data"))
		return errors.New(fmt.Sprintf("invalid tx data"))
	}

	//commission, err := CommissionEstimate(executeTxData, contractAddress, token, amount, rpc)

	relayerTx, err := EncodeExecuteParams(executeCalldata{
		Token:      token,
		Commission: big.NewInt(10),
		Receiver:   receiver,
		CoreData:   executeTxData,
	})

	evmTxReq := resources3.EvmTransactionRequest{
		Data: resources3.EvmTransaction{
			Key: resources3.Key{Type: "evm_transaction"},
			Attributes: resources3.EvmTransactionAttributes{
				Chain: chainId,
				Data:  relayerTx,
			},
		},
	}

	err = r.relayer.PostJSON(u, evmTxReq, context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "relayer returns error")
	}
	return nil
}

func (r *Runner) executeOrder(match resources.Match) error {
	srcChainId, origChainId := match.Relationships.SrcChain.Data.ID, match.Relationships.OriginChain.Data.ID
	srcChain, origChain, err := r.getChains(srcChainId, origChainId)
	if err != nil {
		return errors.Wrap(err, "failed to get chains")
	}

	order, err := r.proxyRepo.Get(origChainId).GetOrder(big.NewInt(match.Attributes.OriginOrderId))
	if err != nil {
		return errors.Wrap(err, "failed to get swapica order")
	}

	swapicaMatch, err := r.proxyRepo.Get(srcChainId).GetMatch(big.NewInt(match.Attributes.MatchId))
	if err != nil {
		return errors.Wrap(err, "failed to get swapica match")
	}

	params := types.ExecuteOrderParams{
		SrcChain:  *origChain,
		DestChain: *srcChain,
		Order:     order,
		Match:     swapicaMatch,
		Receiver:  swapicaMatch.Creator.String(),
		Sender:    order.Creator.String(),
	}

	tx, err := r.proxyRepo.Get(origChainId).ExecuteOrder(params)
	if err != nil {
		return errors.Wrap(err, "failed to get execute order tx")
	}
	if tx == nil {
		return errors.New("failed to build transaction")
	}

	err = r.sendTxToRelayer(origChainId, tx, order.TokenToBuy, swapicaMatch.Creator, common.HexToAddress(origChain.SwapContract), order.AmountToBuy.String(), origChain.RpcEndpoint)
	if err != nil {
		return errors.Wrap(err, "failed to send tx to relayer")
	}
	return nil
}

func (r *Runner) executeMatch(order resources.Order) error {
	srcChainId, destChainId := order.Relationships.SrcChain.Data.ID, order.Relationships.DestinationChain.Data.ID
	srcChain, destChain, err := r.getChains(srcChainId, destChainId)
	if err != nil {
		return errors.Wrap(err, "failed to get chains")
	}

	swapicaOrder, err := r.proxyRepo.Get(srcChain.ID).GetOrder(big.NewInt(order.Attributes.OrderId))
	if err != nil {
		return errors.Wrap(err, "failed to get swapica order")
	}

	if enums.State(swapicaOrder.Status.State) != enums.Executed {
		return errors.New("cannot execute a match if order is not executed")
	}

	match, err := r.proxyRepo.Get(destChain.ID).GetMatch(swapicaOrder.Status.MatchId)
	if err != nil {
		return errors.Wrap(err, "failed to get swapica match")
	}

	params := types.ExecuteMatchParams{
		SrcChain:  *destChain,
		DestChain: *srcChain,
		Order:     swapicaOrder,
		Match:     match,
		Receiver:  match.Creator.String(),
		Sender:    swapicaOrder.Creator.String(),
	}

	tx, err := r.proxyRepo.Get(destChain.ID).ExecuteMatch(params)
	if err != nil {
		return errors.Wrap(err, "failed to get execute match tx")
	}
	if tx == nil {
		return errors.New("failed to build transaction")
	}

	err = r.sendTxToRelayer(destChain.ID, tx, match.TokenToSell, match.Creator, common.HexToAddress(destChain.SwapContract), match.AmountToSell.String(), destChain.RpcEndpoint)
	if err != nil {
		return errors.Wrap(err, "failed to send tx to relayer")
	}
	return nil
}

func (r *Runner) handleOrderUpdates(message []byte) error {
	var recvMessage = struct {
		Action string `json:"action,required"`
		Data   []byte `json:"data,required"`
	}{}

	err := json.Unmarshal(message, &recvMessage)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message")
	}

	switch recvMessage.Action {
	case "add-match":
		// when a match is added we can execute the order
		var match resources.Match
		err = json.Unmarshal(recvMessage.Data, &match)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal message data as match")
		}
		if match.Attributes.UseRelayer {
			err = r.executeOrder(match)
		}
	case "update-order":
		// when the order is updated and its state is Executed we can execute the match
		var order resources.Order
		err = json.Unmarshal(recvMessage.Data, &order)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal message data as order")
		}
		if order.Attributes.UseRelayer {
			err = r.executeMatch(order)
		}
	}
	return err
}

func (r *Runner) ListenNewOrders() error {
	socket, _, err := websocket.DefaultDialer.Dial(r.socketUrl, nil)
	if err != nil {
		return errors.Wrap(err, "failed to dial ws url")
	}

	defer func(socket *websocket.Conn) {
		err := socket.Close()
		if err != nil {
			r.log.WithError(err).Debug("failed to close socket")
		}
	}(socket)

	r.log.Info("Try to start listen new orders")
	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			r.log.WithError(err).Debug("websocket read error")
			return err
		}

		err = r.handleOrderUpdates(message)
		if err != nil {
			r.log.WithError(err).Debug("failed to handle order/match updates")
		}
	}
}

func (r *Runner) Run(_ context.Context) error {
	if r.useRelayer {
		if err := r.ExecuteOrders(); err != nil {
			return err
		}
	}
	err := r.ListenNewOrders()
	return err
}
