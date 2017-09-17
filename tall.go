package meeseeks

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"
	"strconv"
)

func (m *Meeseeks) Tall(msg *message) {
	responseMessage := ""
	perExchangeFormat := "%s %s last: %s, vol: %s | "

	for k, v := range exchangeCodes {
		urlString := fmt.Sprintf(cwMarketsSummaryURL, k, "btcusd")
		resp, err := http.Get(urlString)

		if err != nil {
			m.log.Printf("failure retrieving exchange information from crypto.wat.ch for market %s: %s", v, err)
			return
		}

		jsonBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			m.log.Printf("failure reading json results from cryptowat.ch for market %s: %s", v, err)
			return
		}

		var res MarketSummaryResponse
		err = json.Unmarshal(jsonBody, &res)

		if err != nil {
			m.log.Printf("failure unmarshalling json  from cryptowat.ch for market %s: %s", v, err)
			return
		}

		exchangeSummary := fmt.Sprintf(perExchangeFormat, v, "BTCUSD", 
			strconv.FormatFloat(res.Result.Price.Last, 'f', -1, 64),
			strconv.FormatFloat(res.Result.Volume, 'f', -1, 64))

		responseMessage = responseMessage + exchangeSummary
	}

	m.SendMessage(msg.Chat, responseMessage, nil)
}

var exchangeCodes = map[string]string {
	"bitfinex": "Bitfinex",
	"gdax": "GDAX",
	"bitstamp": "Bitstamp",
	"kraken": "Kraken",
	"btce": "BTC-E",
	"cryptsy": "Cryptsy",
	"okcoin": "OKCoin",
	"gemini": "Gemini",
}