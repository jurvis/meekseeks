package meeseeks

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"

	"github.com/tucnak/telebot"
)

const cwMarketsSummaryURL = "https://api.cryptowat.ch/markets/%s/%s/summary"

type MarketSummaryResponse struct {
	Result struct {
		Price struct {
			Last   float64 `json:"last"`
			High   float64 `json:"high"`
			Low    float64 `json:"low"`
			Change struct {
				Percentage float64 `json:"percentage"`
				Absolute   float64 `json:"absolute"`
			} `json:"change"`
		} `json:"price"`
		Volume float64 `json:"volume"`
	} `json:"result"`
}

func (m *Meeseeks) Ticker(msg *message) {
	if len(msg.Args) == 0 {
		so := &telebot.SendOptions{ReplyTo: *msg.Message, ReplyMarkup: telebot.ReplyMarkup{ForceReply: true, Selective: true}}
		m.SendMessage(msg.Chat, "/ticker: Do a crypto exchange rate converstion\nHere are some commands to try: \n* btcusd \n* btcusd --market coinbase", so)
	}

	pair, market := parseArgs(msg.Args)
	if pair == "" || market == "" {
		// get coinbase btcusd
		pair = "btcusd"
		market = "coinbase"
	}

	urlString := fmt.Sprintf(cwMarketsSummaryURL, market, pair)
	resp, err := http.Get(urlString)
	
	if err != nil {
		m.log.Printf("failure retrieving exchange information from crypto.wat.ch for trading pair %s in %s market: %s", pair, market, err)
		return
	}

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m.log.Printf("failure reading json results from cryptowat.ch for trading pair %s in %s market: %s", pair, market, err)
		return
	}

	var res MarketSummaryResponse
	err = json.Unmarshal(jsonBody, &res)

	if err != nil {
		m.log.Printf("failure unmarshalling json  from cryptowat.ch for trading pair %s in %s market: %s", pair, market, err)
		return
	}

	resMsg := "%s %s ticker | Best bid: %s, 24 hour volume: %s, 24 hour low: %s, 24 hour high: %s"
	rm := fmt.Sprintf(resMsg, 
			strings.Title(market), 
			strings.ToUpper(pair), 
			strconv.FormatFloat(res.Result.Price.Last, 'f', -1, 64),
			strconv.FormatFloat(res.Result.Volume, 'f', -1, 64),
			strconv.FormatFloat(res.Result.Price.Low, 'f', -1, 64),
			strconv.FormatFloat(res.Result.Price.High, 'f', -1, 64))
	m.SendMessage(msg.Chat, rm, nil)
}


// Helper functions
func parseArgs(args []string) (pair string, market string) {
	pair = ""
	market = ""

	for i, a := range args {
		pair = args[0]

		if (a == "--market") {
			market = args[i + 1]
		}
	}

	return pair, market
}