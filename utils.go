package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	gecko "github.com/superoo7/go-gecko/v3"
	geckoTypes "github.com/superoo7/go-gecko/v3/types"
	table "github.com/jedib0t/go-pretty/v6/table"
	humanize "github.com/dustin/go-humanize"
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

func printCoinStats(coin *geckoTypes.CoinsID) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendRows([]table.Row{
		{
			fmt.Sprintf("ID: %v", coin.ID),
			fmt.Sprintf("Name: %v", coin.Name),
			fmt.Sprintf("Symbol: %v", coin.Symbol),
			fmt.Sprintf("MarketCapRank: %v", coin.MarketCapRank),
		},
	})

	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateColumns = true
	t.Style().Options.SeparateRows = true

	t.Render()
}

func printPrices(currencies []string, coin *geckoTypes.CoinsID) {
	if coin != nil && *coin.Tickers != nil {
		fmt.Println()
		t := table.NewWriter()
		t.SetStyle(table.StyleColoredGreenWhiteOnBlack)
		t.SetOutputMirror(os.Stdout)

		var priceData = map[string]map[string]string{}
		for _, ticker := range *coin.Tickers {
			for _, currency := range currencies {
				key := ""
				if strings.ToUpper(ticker.Base) == currency {
					key = strings.ToUpper(ticker.Base)
				} else if strings.ToUpper(ticker.Target) == currency {
					key = strings.ToUpper(ticker.Target)
				}
				if key != "" {
					if priceData[ticker.Market.Name] == nil {
						priceData[ticker.Market.Name] = make(map[string]string)
					}
					priceData[ticker.Market.Name][currency] = getPriceDisplayString(ticker.Last, currency)
				}
			}
		}

		for k, v := range priceData {
			rowData := make(table.Row, len(currencies) + 1)
			rowData[0] = k
			for i, currency := range currencies {
				var val string
				if price, ok := v[currency]; ok {
					val = price
				} else {
					val = "-"
				}
				rowData[i + 1] = val
			}
			t.AppendRow(rowData)
		}

		t.SortBy([]table.SortBy{
			{Name: "Exchange", Mode: table.Asc},
		})

		// Add headers
		headers := make(table.Row, len(currencies) + 1)
		headers[0] = "Exchange"
		for i, currency := range currencies {
			headers[i + 1] = fmt.Sprintf("%s", currency)
		}
		t.AppendHeader(headers)

		t.Render()
	}
}

func getPriceDisplayString(price float64, currency string) string {
	displayString := humanize.FormatFloat(getCurrencyFormatString(currency), price)
	switch currency {
	case "USD":
		displayString = "$" + displayString
	case "EUR":
		displayString = "€" + displayString
	}
	return displayString
}

func getCurrencyFormatString(currency string) string{
	switch currency {
	case "BTC":
		return "#,###.########"
	default:
		return "#,###.##"
	}
}
