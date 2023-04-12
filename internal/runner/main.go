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
	"github.com/gorilla/websocket"
	jsonapi "gitlab.com/distributed_lab/json-api-connector"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/connectors/signed"
	"math/big"
	"net/http"
	"net/url"
)

type Runner struct {
	socketUrl       string
	aggregator      *jsonapi.Connector
	relayer         *jsonapi.Connector
	log             *logan.Entry
	proxyRepo       proxy.ProxyRepo
	chains          data.ChainsQ
	tokens          data.TokensQ
	sendAutoExecute bool
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
		socketUrl:       cfg.RunnerCfg().Aggregator.Ws,
		aggregator:      aggregator,
		relayer:         relayer,
		log:             cfg.Log(),
		proxyRepo:       proxyRepo,
		chains:          chains,
		tokens:          tokens,
		sendAutoExecute: cfg.RunnerCfg().SendAutoExecute,
	}

	return runner, nil
}

func (r *Runner) fetchOrders() (orders []resources.Order, err error) {
	var limit = 100
	var response resources.OrderListResponse
	u, _ := url.Parse(fmt.Sprintf("/orders?filter[auto_execute]=true&filter[state]=1&page[number]=0&page[limit]=%d", limit))
	for {
		err = r.aggregator.Get(u, &response)
		if err != nil {
			return nil, err
		}
		u, _ = url.Parse(response.Links.Next)
		orders = append(orders, response.Data...)
		if len(response.Data) < limit {
			break
		}
	}
	return orders, nil
}

func (r *Runner) fetchMatches() (matches []resources.Match, err error) {
	var limit = 100
	var response resources.MatchListResponse
	u, _ := url.Parse(fmt.Sprintf("/match_orders?filter[auto_execute]=true&filter[state]=2&page[number]=0&page[limit]=%d", limit))
	for {
		err = r.aggregator.Get(u, &response)
		if err != nil {
			return nil, errors.Wrap(err, "failed to send get request")
		}

		u, err = url.Parse(response.Links.Next)
		if err != nil {
			return nil, errors.Wrap(err, "failed parse url")
		}

		matches = append(matches, response.Data...)
		if len(response.Data) < limit {
			break
		}
	}
	return matches, nil
}

func (r *Runner) ExecuteOrders() error {
	_, err := r.fetchOrders()
	if err != nil {
		return errors.Wrap(err, "failed to get order list from aggregator")
	}

	// TODO: process orders

	_, err = r.fetchMatches()
	if err != nil {
		return errors.Wrap(err, "failed to get match list from aggregator")
	}

	// TODO: process match orders
	return nil
}

func (r *Runner) getChains(srcId, destId string) (src *data.Chain, dest *data.Chain, err error) {
	src, err = r.chains.New().FilterByID(srcId).Get()
	if err != nil {
		return
	}

	if src == nil {
		return nil, nil, errors.New("src chain not found")
	}

	dest, err = r.chains.New().FilterByID(destId).Get()
	if dest == nil {
		return nil, nil, errors.New("dest chain not found")
	}
	return
}

func (r *Runner) sendTxToRelayer(chainId string, data interface{}) error {
	u, _ := url.Parse("/transaction")
	evmTx, ok := data.(resources2.EvmTransaction)
	if !ok {
		r.log.Debug(fmt.Sprintf("expected type EvmTransaction, got %T", data))
		return errors.New(fmt.Sprintf("invalid tx data format"))
	}

	evmTxReq := resources3.EvmTransactionRequest{
		Data: resources3.EvmTransaction{
			Key: resources3.Key{Type: "evm_transaction"},
			Attributes: resources3.EvmTransactionAttributes{
				Chain: chainId,
				Data:  evmTx.Attributes.TxBody.Data,
			},
		},
	}

	err := r.relayer.PostJSON(u, evmTxReq, context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "relayer return error")
	}
	return nil
}

func (r *Runner) executeOrder(data []byte) error {
	var match resources.Match
	err := json.Unmarshal(data, &match)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message data as match")
	}

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

	err = r.sendTxToRelayer(origChainId, tx)
	if err != nil {
		return errors.Wrap(err, "failed to send tx to relayer")
	}

	return nil
}

func (r *Runner) executeMatch(data []byte) error {
	var order resources.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message data as order")
	}

	srcChainId, destChainId := order.Relationships.SrcChain.Data.ID, order.Relationships.DestinationChain.Data.ID
	srcChain, destChain, err := r.getChains(srcChainId, destChainId)
	if err != nil {
		return errors.Wrap(err, "failed to get chains")
	}

	swapicaOrder, err := r.proxyRepo.Get(srcChainId).GetOrder(big.NewInt(order.Attributes.OrderId))
	if err != nil {
		return errors.Wrap(err, "failed to get swapica order")
	}

	if enums.State(swapicaOrder.Status.State) != enums.Executed {
		return errors.New("cannot execute a match if order is not executed")
	}

	match, err := r.proxyRepo.Get(destChainId).GetMatch(swapicaOrder.Status.MatchId)
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

	tx, err := r.proxyRepo.Get(destChainId).ExecuteMatch(params)
	if err != nil {
		return errors.Wrap(err, "failed to get execute match tx")
	}
	if tx == nil {
		return errors.New("failed to build transaction")
	}

	err = r.sendTxToRelayer(destChainId, tx)
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
		err = r.executeOrder(recvMessage.Data)
	case "update-order":
		// when the order is updated and its state is Executed we can execute the match
		err = r.executeMatch(recvMessage.Data)
	}
	return err
}

func (r *Runner) ListenNewOrders() error {
	socket, _, err := websocket.DefaultDialer.Dial(r.socketUrl, nil)
	if err != nil {
		return errors.Wrap(err, "failed to dial ws url")
	}

	defer socket.Close()

	r.log.Info("Try to start listen new orders")
	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			r.log.Debug("read error:", err)
		}

		err = r.handleOrderUpdates(message)
		if err != nil {
			r.log.WithError(err).Debug("failed to handle order/match updates")
		}
	}
}

func (r *Runner) Run(_ context.Context) error {
	if r.sendAutoExecute {
		if err := r.ExecuteOrders(); err != nil {
			return err
		}
	}
	err := r.ListenNewOrders()
	return err
}
