package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Markets []struct {
		BaseMarket     string  `json:"baseMarket"`
		QuoteMarket    string  `json:"quoteMarket"`
		MinBuyAmount   float64 `json:"minBuyAmount,omitempty"`
		MinSellAmount  float64 `json:"minSellAmount,omitempty"`
		BasePrecision  int     `json:"basePrecision,omitempty"`
		QuotePrecision int     `json:"quotePrecision,omitempty"`
		Status         string  `json:"status"`
		Fee            struct {
			Bid struct {
				Maker float64 `json:"maker"`
				Taker float64 `json:"taker"`
			} `json:"bid"`
			Ask struct {
				Maker float64 `json:"maker"`
				Taker float64 `json:"taker"`
			} `json:"ask"`
		} `json:"fees"`
		Low                string      `json:"low,omitempty"`
		High               string      `json:"high,omitempty"`
		Last               string      `json:"last,omitempty"`
		Type               string      `json:"type"`
		Open               json.Number `json:"open,omitempty"`
		Volume             string      `json:"volume,omitempty"`
		Sell               string      `json:"sell,omitempty"`
		Buy                string      `json:"buy,omitempty"`
		At                 int         `json:"at,omitempty"`
		MaxBuyAmount       int         `json:"maxBuyAmount,omitempty"`
		MinBuyVolume       int         `json:"minBuyVolume,omitempty"`
		MaxBuyVolume       int         `json:"maxBuyVolume,omitempty"`
		FeePercentOnProfit float64     `json:"feePercentOnProfit,omitempty"`
	} `json:"markets"`
	Assets []struct {
		Type              string  `json:"type"`
		Name              string  `json:"name"`
		Deposit           string  `json:"deposit"`
		Withdrawal        string  `json:"withdrawal"`
		ListingType       string  `json:"listingType"`
		Category          string  `json:"category"`
		WithdrawFee       float64 `json:"withdrawFee,omitempty"`
		MinWithdrawAmount string  `json:"minWithdrawAmount,omitempty"`
		MaxWithdrawAmount float64 `json:"maxWithdrawAmount,omitempty"`
		MinDepositAmount  float64 `json:"minDepositAmount,omitempty"`
		Confirmations     int     `json:"confirmations,omitempty"`
	} `json:"assets"`
}

func main() {

	resp, err := http.Get("https://api.wazirx.com/api/v2/market-status")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	if err != nil {
		fmt.Println("wrong here")
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println(err)
	}

	fmt.Println(result.Assets[0].Type)

	withdrawalprices := make(map[string]string)
	withdrawalprices

}
