package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"os"

	prompt "github.com/c-bata/go-prompt"
	gecko "github.com/superoo7/go-gecko/v3"
	geckoTypes "github.com/superoo7/go-gecko/v3/types"
	table "github.com/jedib0t/go-pretty/v6/table"
)

func getCoinsList(client *gecko.Client) geckoTypes.CoinList {
	coinsList, err := client.CoinsList()
	if err != nil {
		log.Fatal(err)
	}
	return *coinsList
}

func getCurrencyList(client *gecko.Client) geckoTypes.SimpleSupportedVSCurrencies {
	currencyList, err := client.SimpleSupportedVSCurrencies()
	if err != nil {
		log.Fatal(err)
	}
	return *currencyList
}

func getCoinSuggestions(coinsList geckoTypes.CoinList) []prompt.Suggest {
	suggestionsById := make([]prompt.Suggest, len(coinsList))
	for i, val := range coinsList {
		item := prompt.Suggest{
			Text:        val.ID,
			Description: val.Symbol,
		}
		suggestionsById[i] = item
	}

	suggestionsBySymbol := make([]prompt.Suggest, len(coinsList))
	for i, val := range coinsList {
		item := prompt.Suggest{
			Text:        val.Symbol,
			Description: val.ID,
		}
		suggestionsBySymbol[i] = item
	}

	return append(suggestionsById, suggestionsBySymbol...)
}

func getCurrencySuggestions(currencyList geckoTypes.SimpleSupportedVSCurrencies) []prompt.Suggest {
	suggestions := make([]prompt.Suggest, len(currencyList))
	for i, val := range currencyList {
		item := prompt.Suggest{
			Text:        val,
			Description: val,
		}
		suggestions[i] = item
	}
	return suggestions
}

func getCoin(id string) (*geckoTypes.CoinsID, error) {
	cg := gecko.NewClient(nil)
	coin, err := cg.CoinsID(id, false, true, true, false, false, false)
	return coin, err
}

func getCoinBySymbol(symbol string) (*geckoTypes.CoinsID, error) {
	symbol = strings.ToLower(symbol)
	client := gecko.NewClient(nil)
	coinsList := getCoinsList(client)
	var coin *geckoTypes.CoinsID = nil

	for _, val := range coinsList {
		if val.Symbol == symbol && !strings.HasPrefix(val.ID, "binance-peg") {
			// Found coin; get and return data
			coin, err := getCoin(val.ID)
			if err == nil {
				return coin, nil
			}
		}
	}

	errorMsg := fmt.Sprint("Could not find a coin with symbol of ", symbol)
	return coin, errors.New(errorMsg)
}

func printPrices(currency string, coin *geckoTypes.CoinsID) {
	if coin != nil && *coin.Tickers != nil {
		t := table.NewWriter()
		t.SetStyle(table.StyleColoredGreenWhiteOnBlack)
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Exchange", fmt.Sprintf("Price (%s)", currency)})
		fmt.Println(coin.ID, currency)
		for _, val := range *coin.Tickers {
			if strings.ToUpper(val.Target) == currency || val.Base == currency {
				t.AppendRow(table.Row{val.Market.Name, val.Last})
			}
		}
		t.Render()
	}
}
