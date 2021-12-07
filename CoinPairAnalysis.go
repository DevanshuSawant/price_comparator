package main

import (
	//"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	"sort"
	"strconv"
	"strings"
	"time"
	//"github.com/adshao/go-binance/v2"
	//"github.com/joho/godotenv"
	//"reflect"
)

type ResponseWazirx struct {		//It follows the structure of the response from the GET request to wazirx API .It converts it into strings/floats etc(using variables)
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

type Account struct {
	MakerCommission  int64     `json:"makerCommission"`
	TakerCommission  int64     `json:"takerCommission"`
	BuyerCommission  int64     `json:"buyerCommission"`
	SellerCommission int64     `json:"sellerCommission"`
	CanTrade         bool      `json:"canTrade"`
	CanWithdraw      bool      `json:"canWithdraw"`
	CanDeposit       bool      `json:"canDeposit"`
	Balances         []struct {
		Asset  string `json:"asset"`
		Free   string `json:"free"`
		Locked string `json:"locked"`
	} `json:"balances"`
}



func main() {
	start := time.Now()

	wazirxResponse := time.Now()
	respWazirx, err := http.Get("https://api.wazirx.com/api/v2/market-status")		//GET request to wazirx for ticker prices and withdrawal/deposit enabled lists
	if err != nil {
		log.Fatal(err)
	}
	wazirxResponseTime :=time.Since(wazirxResponse)
	defer respWazirx.Body.Close()
	bodyWazirx, err := ioutil.ReadAll(respWazirx.Body) // response body is []byte
	if err != nil {
		fmt.Println("wrong here")
	}
	var resultWazirx ResponseWazirx
	if err := json.Unmarshal(bodyWazirx, &resultWazirx); err != nil { // Parse []byte to the go struct pointer
		fmt.Println(err)
	}

	binanceResponse := time.Now()
	respBinance, err := http.Get("https://api.binance.com/api/v3/ticker/price")		//GET request to binance for ticker prices
	if err != nil {
		log.Fatal(err)
	}
	binanceResponseTime := time.Since(binanceResponse)
	defer respBinance.Body.Close()
	bodyBinance, err := ioutil.ReadAll(respBinance.Body) // response body is []byte
	if err != nil {
		fmt.Println("wrong here")
	}
	var resultBinance []map[string]string
	if err := json.Unmarshal(bodyBinance, &resultBinance); err != nil { // Parse []byte to the go struct pointer
		fmt.Println(err)
	}

	withdrawalprices := make(map[string]string)
	depositprices := make(map[string]string)

	for i := 0; i < len(resultWazirx.Assets); i ++ {		//converts json array into comparable map eg key:ADAUSDT value:1.4 
		pairSymbol := strings.ToUpper(resultWazirx.Assets[i].Type) + "USDT"
		if resultWazirx.Assets[i].Withdrawal == "enabled" {
			for j := 0; j < len(resultWazirx.Markets); j ++ {
				if resultWazirx.Assets[i].Type == resultWazirx.Markets[j].BaseMarket && resultWazirx.Markets[j].QuoteMarket == "usdt"  {
					withdrawalprices[pairSymbol] = resultWazirx.Markets[j].Buy
				}
			}
		}
		if resultWazirx.Assets[i].Deposit =="enabled" {
			for j := 0; j <  len(resultBinance); j ++ {
				if resultBinance[j]["symbol"] == pairSymbol {
					depositprices[pairSymbol] = resultBinance[j]["price"]
				}
			}
		}
	}

	binanceCoin := make(map[string]float64)
	wazirxCoin := make(map[string]float64)

	for kw, vw := range withdrawalprices {		//compares both to get percentage diffrences between the two
		for kd, vd := range depositprices {
			if kw == kd {
				switch  {
					case vw == vd:
						fmt.Println(kw, "both equal")		

					case vw > vd:
						s, err1 := strconv.ParseFloat(vw, 64);
						g, err2 := strconv.ParseFloat(vd, 64);
						if err1 == nil && err2 == nil {
							x := ((s-g)/g)*100
							//fmt.Println("wazirx  higher%", x, "for", kw)
							binanceCoin[kw] = x
						}
					
					case vw < vd:
						s, err1 := strconv.ParseFloat(vw, 64);
						g, err2 := strconv.ParseFloat(vd, 64);
						if err1 == nil && err2 == nil {
							x := ((g-s)/s)*100
							//fmt.Println("binance higher%", x, "for", kw)
							wazirxCoin[kw] = x
						}
				}
			}
		}									
	}

    keysBinance := make([]string, 0, len(binanceCoin))
    for key1 := range binanceCoin {
        keysBinance = append(keysBinance, key1)
    }
    sort.Slice(keysBinance, func(i, j int) bool { return binanceCoin[keysBinance[i]] > binanceCoin[keysBinance[j]] })
    for _, key1 := range keysBinance {
    	fmt.Printf("%s, %f\n", key1, binanceCoin[key1])

    }

	keysWazirx := make([]string, 0, len(wazirxCoin))
    for key2 := range wazirxCoin {
        keysWazirx = append(keysWazirx, key2)
    }
    sort.Slice(keysWazirx, func(i, j int) bool { return wazirxCoin[keysWazirx[i]] > wazirxCoin[keysWazirx[j]] })
    for _, key2 := range keysWazirx {
    	fmt.Printf("%s, %f\n", key2, wazirxCoin[key2])

    }
/*
	dotEnv := godotenv.Load()
	if dotEnv != nil {
	  log.Fatal("Error loading .env file")
	}

	var (
		apiKey = os.Getenv("apiKey")
		secretKey = os.Getenv("secretKey")
	)
	binance.UseTestnet = true
	client := binance.NewClient(apiKey, secretKey)

	binanceBot := time.Now()
	Account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
   		fmt.Println(err)
    	return
	}
	fmt.Println(Account.Balances[0].Asset)

	order, err := client.NewCreateOrderService().Symbol("BUSDETH").
        Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
        TimeInForce(binance.TimeInForceTypeGTC).Quantity("1").
        Price("4000").Do(context.Background())
		if err != nil {
   		 fmt.Println(err)
   		 return
		}
	fmt.Println(order)

	fmt.Println(res1)
*/

	//binanceBotTime := time.Since(binanceBot)

/*
	fmt.Println(keysWazirx)
	fmt.Println(keysBinance)
	fmt.Println(wazirxCoin)
	fmt.Println(binanceCoin)
*/
	fmt.Println(time.Since(start))		//for measuring the processing time
	time := time.Since(start) - wazirxResponseTime - binanceResponseTime //-binanceBotTime
	fmt.Println("Time For Code To Run (Excluding API Response Time):", time)
	
}