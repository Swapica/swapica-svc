package runner

import (
	"encoding/json"
	"fmt"
	"github.com/Swapica/order-aggregator-svc/resources"
	"github.com/Swapica/swapica-svc/internal/config"
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/internal/proxy"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
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
	socket     *websocket.Conn
	aggregator *jsonapi.Connector
	log        *logan.Entry
	proxyRepo  proxy.ProxyRepo
	chains     data.ChainsQ
	tokens     data.TokensQ
}

func NewRunner(cfg config.Config, proxyRepo proxy.ProxyRepo, chains data.ChainsQ, tokens data.TokensQ) (*Runner, error) {
	socket, _, err := websocket.DefaultDialer.Dial(cfg.RunnerCfg().Ws, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial ws url")
	}

	u, _ := url.Parse(cfg.RunnerCfg().Url)
	conn := jsonapi.NewConnector(signed.NewClient(&http.Client{Timeout: cfg.RunnerCfg().Timeout}, u))

	runner := &Runner{
		socket:     socket,
		aggregator: conn,
		log:        cfg.Log(),
		proxyRepo:  proxyRepo,
		chains:     chains,
		tokens:     tokens,
	}

	return runner, nil
}

func (r *Runner) fetchOrders() (orders []resources.Order, err error) {
	var limit = 100
	var response resources.OrderListResponse
	u, _ := url.Parse(fmt.Sprintf("/integrations/order-aggregator/orders?filter[auto_execute]=true&filter[state]=1&page[number]=0&page[limit]=%d", limit))
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
	u, _ := url.Parse(fmt.Sprintf("/integrations/order-aggregator/match_orders?filter[auto_execute]=true&filter[state]=2&page[number]=0&page[limit]=%d", limit))
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

func (r *Runner) ExecuteOrders() {
	_, err := r.fetchOrders()
	if err != nil {
		r.log.WithError(err).Error("failed to get order list from aggregator")
	}

	// TODO: process orders

	_, err = r.fetchMatches()
	if err != nil {
		r.log.WithError(err).Error("failed to get match list from aggregator")
	}

	// TODO: process match orders
}

func (r *Runner) getChains(srcId, destId string) (src *data.Chain, dest *data.Chain, err error) {
	src, err = r.chains.FilterByID(srcId).Get()
	if err != nil {
		return
	}

	if src == nil {
		return nil, nil, errors.New("src chain not found")
	}

	dest, err = r.chains.FilterByID(destId).Get()
	if dest != nil {
		return nil, nil, errors.New("dest chain not found")
	}
	return
}

func (r *Runner) executeOrder(data []byte) error {
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

	match, err := r.proxyRepo.Get(destChainId).GetMatch(big.NewInt(*order.Attributes.MatchId))
	if err != nil {
		return errors.Wrap(err, "failed to get swapica match")
	}

	params := types.ExecuteOrderParams{
		SrcChain:  *srcChain,
		DestChain: *destChain,
		Order:     swapicaOrder,
		Match:     match,
		Receiver:  match.Creator.String(),
		Sender:    swapicaOrder.Creator.String(),
	}

	tx, err := r.proxyRepo.Get(srcChainId).ExecuteOrder(params)
	if err != nil {
		return errors.Wrap(err, "failed to execute order")
	}
	if tx == nil {
		return errors.New("failed to build transaction")
	}
	return nil
}

func (r *Runner) executeMatch(data []byte) error {
	var match resources.Match
	err := json.Unmarshal(data, &match)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message data as order")
	}

	srcChainId, destChainId := match.Relationships.SrcChain.Data.ID, match.Relationships.OriginChain.Data.ID
	srcChain, destChain, err := r.getChains(srcChainId, destChainId)
	if err != nil {
		return errors.Wrap(err, "failed to get chains")
	}

	order, err := r.proxyRepo.Get(srcChainId).GetOrder(big.NewInt(match.Attributes.OriginOrderId))
	if err != nil {
		return errors.Wrap(err, "failed to get swapica order")
	}

	swapicaMatch, err := r.proxyRepo.Get(destChainId).GetMatch(big.NewInt(match.Attributes.MatchId))
	if err != nil {
		return errors.Wrap(err, "failed to get swapica match")
	}

	params := types.ExecuteMatchParams{
		SrcChain:  *srcChain,
		DestChain: *destChain,
		Order:     order,
		Match:     swapicaMatch,
		Receiver:  swapicaMatch.Creator.String(),
		Sender:    order.Creator.String(),
	}

	tx, err := r.proxyRepo.Get(srcChainId).ExecuteMatch(params)
	if err != nil {
		return errors.Wrap(err, "failed to execute match")
	}
	if tx == nil {
		return errors.New("failed to build transaction")
	}
	return nil
}

func (r *Runner) handleNewOrder(message []byte) error {
	var recvMessage = struct {
		Action string `json:"action,required"`
		Data   []byte `json:"data,required"`
	}{}
	err := json.Unmarshal(message, &recvMessage)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message")
	}

	switch recvMessage.Action {
	case "add-order":
		err = r.executeOrder(recvMessage.Data)
	case "add-match":
		err = r.executeMatch(recvMessage.Data)
	}
	return err
}

func (r *Runner) ListenNewOrders() {
	for {
		_, message, err := r.socket.ReadMessage()
		if err != nil {
			r.log.Debug("read error:", err)
		}

		err = r.handleNewOrder(message)
		if err != nil {
			r.log.WithError(err).Error("failed to handle new order/match")
		}
	}
}
