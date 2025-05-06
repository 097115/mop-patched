// Copyright (c) 2013-2023 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	"net/http"
)

const marketURL = `https://query1.finance.yahoo.com/v7/finance/quote?crumb=%s&symbols=%s`
const marketURLQueryParts = `&range=1d&interval=5m&indicators=close&includeTimestamps=false&includePrePost=false&corsDomain=finance.yahoo.com&.tsrc=finance`

// Market stores current market information displayed in the top three lines of
// the screen. The market data is fetched and parsed from the HTML page above.
type Market struct {
	IsClosed  bool              // True when U.S. markets are closed.
	Dow       map[string]string // Hash of Dow Jones indicators.
	Nasdaq    map[string]string // Hash of NASDAQ indicators.
	Sp500     map[string]string // Hash of S&P 500 indicators.
	Btc       map[string]string
	Tokyo     map[string]string
	HongKong  map[string]string
	London    map[string]string
	Frankfurt map[string]string
	Yield     map[string]string
	Silver    map[string]string
	Yen       map[string]string
	Rub       map[string]string
	Gbp       map[string]string
	Euro      map[string]string
	Gold      map[string]string
	errors    string // Error(s), if any.
	url       string // URL with symbols to fetch data
	cookies   string // cookies for auth
	crumb     string // crumb for the cookies, to be applied as a query param
}

// Returns new initialized Market struct.
func NewMarket() *Market {
	market := &Market{}
	market.IsClosed = false
	market.Dow = make(map[string]string)
	market.Nasdaq = make(map[string]string)
	market.Sp500 = make(map[string]string)
	market.Btc = make(map[string]string)

	market.Tokyo = make(map[string]string)
	market.HongKong = make(map[string]string)
	market.London = make(map[string]string)
	market.Frankfurt = make(map[string]string)

	market.Yield = make(map[string]string)
	market.Silver = make(map[string]string)
	market.Yen = make(map[string]string)
	market.Euro = make(map[string]string)
	market.Gold = make(map[string]string)

	market.cookies = fetchCookies()
	market.crumb = fetchCrumb(market.cookies)
	market.url = fmt.Sprintf(marketURL, market.crumb, `^DJI,^IXIC,^GSPC,BTC-USD,^N225,^HSI,^FTSE,^GDAXI,JPY=X,RUB=X,GBP=X,EUR=X,^TNX,SI=F,GC=F`) + marketURLQueryParts

	market.errors = ``

	return market
}

// Fetch downloads HTML page from the 'marketURL', parses it, and stores resulting data
// in internal hashes. If download or data parsing fails Fetch populates 'market.errors'.
func (market *Market) Fetch() (self *Market) {
	self = market // <-- This ensures we return correct market after recover() from panic().
	defer func() {
		if err := recover(); err != nil {
			market.errors = "" // Don't pollute the screen
		}
	}()

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	request, err := http.NewRequest("GET", market.url, nil)
	if err != nil {
		panic(err)
	}

	request.Header = http.Header{
		"Accept":          {"*/*"},
		"Accept-Language": {"en-US,en;q=0.5"},
		"Connection":      {"keep-alive"},
		"Content-Type":    {"application/json"},
		"Cookie":          {market.cookies},
		"Host":            {"query1.finance.yahoo.com"},
		"Origin":          {"https://finance.yahoo.com"},
		"Referer":         {"https://finance.yahoo.com"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-site"},
		"TE":              {"trailers"},
		"User-Agent":      {userAgent},
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	body = market.isMarketOpen(body)
	return market.extract(body)
}

// Ok returns two values: 1) boolean indicating whether the error has occurred,
// and 2) the error text itself.
func (market *Market) Ok() (bool, string) {
	return market.errors == ``, market.errors
}

// -----------------------------------------------------------------------------
func (market *Market) isMarketOpen(body []byte) []byte {
	// TBD -- CNN page doesn't seem to have market open/close indicator.
	return body
}

// -----------------------------------------------------------------------------
func assign(results []map[string]interface{}, position int, changeAsPercent bool) map[string]string {
	out := make(map[string]string)
	out[`change`] = float2Str(results[position]["regularMarketChange"].(float64))
	out[`latest`] = float2Str(results[position]["regularMarketPrice"].(float64))
	if changeAsPercent {
		out[`change`] = float2Str(results[position]["regularMarketChangePercent"].(float64)) + `%`
	} else {
		out[`percent`] = float2Str(results[position]["regularMarketChangePercent"].(float64)) + `%`
	}
	return out
}

// -----------------------------------------------------------------------------
func (market *Market) extract(body []byte) *Market {
	d := map[string]map[string][]map[string]interface{}{}
	err := json.Unmarshal(body, &d)
	if err != nil {
		panic(err)
	}
	results := d["quoteResponse"]["result"]
	market.Dow = assign(results, 0, false)
	market.Nasdaq = assign(results, 1, false)
	market.Sp500 = assign(results, 2, false)
	market.Btc = assign(results, 3, false)
	market.Tokyo = assign(results, 4, false)
	market.HongKong = assign(results, 5, false)
	market.London = assign(results, 6, false)
	market.Frankfurt = assign(results, 7, false)

	market.Yen = assign(results, 8, false)
	market.Rub = assign(results, 9, false)
	market.Gbp = assign(results, 10, false)
	market.Euro = assign(results, 11, false)
	market.Yield[`name`] = `10-year Yield`
	market.Yield = assign(results, 12, false)
	market.Silver = assign(results, 13, false)
	market.Gold = assign(results, 14, false)

	return market
}
