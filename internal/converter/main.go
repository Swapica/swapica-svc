package converter

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	lru "github.com/golang/groupcache/lru"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const STATS_TOKEN = "USDT"

type Converter interface {
	ConvertToStat(amount *big.Float, source string) (*big.Float, error)
	Convert(amount *big.Float, source string, dest string) (*big.Float, error)
}

type converter struct {
	cache priceCache
	log   *logan.Entry
}

func NewConverter() Converter {
	return &converter{
		cache: priceCache{
			store:        lru.New(defaultCacheSize),
			itemLifespan: defaultItemLifespan,
		},
		log: logan.New(),
	}
}

func (c *converter) Get(key string) (*big.Float, bool) {
	raw, isExist := c.cache.Get(key)
	if !isExist {
		return nil, isExist
	}
	value, ok := raw.(priceCacheValue)
	if !ok {
		c.log.WithField("raw", raw).Error("failed to cast cached value")
		return nil, ok
	}
	if value.expirationTime.After(time.Now()) {
		return value.price, true
	}
	return nil, false
}

func (c *converter) ConvertToStat(amount *big.Float, source string) (result *big.Float, err error) {
	return c.Convert(amount, source, STATS_TOKEN)
}

func (c *converter) Convert(amount *big.Float, source string, dest string) (result *big.Float, err error) {
	if strings.EqualFold(source, dest) {
		return amount, nil
	}

	res := new(big.Float)

	cacheKey := strings.ToLower(source + dest)
	price, ok := c.Get(cacheKey)

	if ok {
		return res.Mul(amount, price), nil
	}

	cacheKey = strings.ToLower(dest + source)
	price, ok = c.Get(cacheKey)
	if ok {
		return res.Mul(amount, revertPrice(price)), nil
	}

	price = c.getPrice(source, dest)
	if price.Cmp(big.NewFloat(0)) == 0 {
		return nil, errors.New("failed to convert")
	}

	c.cache.Add(strings.ToLower(source+dest), price)
	return res.Mul(amount, price), nil
}

func (c *converter) getPrice(source, dest string) *big.Float {
	for _, f := range priceFeeds {
		price, err := f(source, dest)

		if err != nil || price.Cmp(big.NewFloat(0)) == 0 {
			price, err = f(dest, source)
			if err != nil || price.Cmp(big.NewFloat(0)) == 0 {
				continue
			}
			price = revertPrice(price)
		}

		return price
	}
	return big.NewFloat(0)
}

func revertPrice(price *big.Float) *big.Float {
	res := big.NewFloat(1)
	return res.Quo(res, price)
}

var priceFeeds = []func(string, string) (*big.Float, error){
	getBinancePrice,
	getGateIoPrice,
	getKuCoinPrice,
	getMEXCPrice,
	getHuobiPrice,
}

func getGateIoPrice(source string, dest string) (*big.Float, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.gate.io/api2/1/ticker/%s_%s", source, dest))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get price")
	}

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	var price gateIoPriceResponse
	err = json.Unmarshal(raw, &price)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	if price.Last == nil || price.Last.Cmp(big.NewFloat(0)) == 0 {
		return nil, errors.New("the price is zero")
	}

	return price.Last, nil
}

type gateIoPriceResponse struct {
	Last *big.Float `json:"last"`
}

func getBinancePrice(source string, dest string) (*big.Float, error) {
	source = strings.ToUpper(source)
	dest = strings.ToUpper(dest)
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=" + source + dest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get price")
	}

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	var price binancePriceResponse
	err = json.Unmarshal(raw, &price)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	if price.Price == nil || price.Price.Cmp(big.NewFloat(0)) == 0 {
		return nil, errors.New("the price is zero")
	}

	return price.Price, nil
}

type binancePriceResponse struct {
	Price *big.Float `json:"price"`
}

func getKuCoinPrice(source, dest string) (*big.Float, error) {
	source = strings.ToUpper(source)
	dest = strings.ToUpper(dest)

	resp, err := http.Get(fmt.Sprintf("https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=%s-%s", source, dest))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get price")
	}

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	var price kuCoinPriceResponse
	err = json.Unmarshal(raw, &price)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	if price.Data.Price == nil || price.Data.Price.Cmp(big.NewFloat(0)) == 0 {
		return nil, errors.New("the price is zero")
	}
	return price.Data.Price, nil
}

type kuCoinPriceResponse struct {
	Data struct {
		Price *big.Float `json:"price"`
	} `json:"data"`
}

func getMEXCPrice(source, dest string) (*big.Float, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.mexc.com/open/api/v2/market/ticker?symbol=%s_%s", source, dest))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get price")
	}

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	var price mexcPriceResponse
	err = json.Unmarshal(raw, &price)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	if len(price.Data) == 0 {
		return nil, errors.New("failed to get price")
	}

	return price.Data[0].Last, nil
}

type mexcPriceResponse struct {
	Data []mexcData `json:"data"`
}

type mexcData struct {
	Last *big.Float `json:"last"`
}

func getHuobiPrice(source, dest string) (*big.Float, error) {
	source = strings.ToLower(source)
	dest = strings.ToLower(dest)

	resp, err := http.Get("https://api.huobi.pro/market/trade?symbol=" + source + dest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get price")
	}

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	var priceResponse huobiPriceResponse
	err = json.Unmarshal(raw, &priceResponse)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	if len(priceResponse.Tick.Data) == 0 {
		return nil, errors.New("failed to get price")
	}

	amount, err := priceResponse.Tick.Data[0].Price.Float64()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get price")
	}

	return big.NewFloat(amount), nil
}

type huobiPriceResponse struct {
	Tick struct {
		Data []huobiData `json:"data"`
	} `json:"tick"`
}

type huobiData struct {
	Price json.Number `json:"price"`
}
