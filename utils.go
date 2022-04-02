package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	gecko "github.com/superoo7/go-gecko/v3"
	"github.com/superoo7/go-gecko/v3/types"
)

func getCoinsList(client *gecko.Client) types.CoinList {
	coinsList, err := client.CoinsList()
	if err != nil {
		log.Fatal(err)
	}
	return *coinsList
}

func getCurrencyList(client *gecko.Client) types.SimpleSupportedVSCurrencies {
	currencyList, err := client.SimpleSupportedVSCurrencies()
	if err != nil {
		log.Fatal(err)
	}
	return *currencyList
}

func getCoinSuggestions(coinsList types.CoinList) []prompt.Suggest {
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

func getCurrencySuggestions(currencyList types.SimpleSupportedVSCurrencies) []prompt.Suggest {
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

func getCoin(id string) (*types.CoinsID, error) {
	cg := gecko.NewClient(nil)
	coin, err := cg.CoinsID(id, false, true, true, false, false, false)
	return coin, err
}

func getCoinBySymbol(symbol string) (*types.CoinsID, error) {
	symbol = strings.ToLower(symbol)
	client := gecko.NewClient(nil)
	coinsList := getCoinsList(client)
	var coin *types.CoinsID = nil
	var id string = ""

	switch symbol {
	case "dot":
		id = "polkadot"
	case "eth":
		id = "ethereum"
	default:
		for _, val := range coinsList {
			if val.Symbol == symbol {
				id = val.ID
			}
		}
	}

	if id != "" {
		coin, err := getCoin(id)
		if err == nil {
			return coin, nil
		}
	}

	errorMsg := fmt.Sprint("Could not find a coin with symbol of ", symbol)
	return coin, errors.New(errorMsg)
}

func printPrices(currency string, coin *types.CoinsID) {
	if coin != nil && *coin.Tickers != nil {
		fmt.Println(coin.ID, currency)
		for _, val := range *coin.Tickers {
			if strings.ToUpper(val.Target) == currency || val.Base == currency {
				displayString := fmt.Sprintf("  %s: %f", val.Market.Name, val.Last)
				fmt.Println(displayString)
			}
		}
	}
}
