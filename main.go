package main

import (
	"fmt"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	gecko "github.com/superoo7/go-gecko/v3"
)

var geckoClient = gecko.NewClient(nil)
var coinsList = getCoinsList(geckoClient)
var currencyList = getCurrencyList(geckoClient)
var coinSuggestions = getCoinSuggestions(coinsList)
var currencySuggestions = getCurrencySuggestions(currencyList)

func executor(input string) {
	input = strings.TrimSpace(input)
	args := strings.Split(input, " ")

	coinId := args[0]

	currency := "USD"
	if len(args) > 1 {
		currency = strings.ToUpper(args[1])
	}

	coin, err := getCoin(coinId)
	if err != nil {
		coin, err := getCoinBySymbol(coinId)
		printPrices(currency, coin)
		if err != nil || coin == nil {
			fmt.Println("Error getting coin with ID of ", input)
		}
	} else {
		printPrices(currency, coin)
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	w := in.CurrentLineBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(w, " ")
	argsLen := len(args)
	lastArg := in.GetWordBeforeCursor()

	var suggestions []prompt.Suggest
	switch argsLen {
	case 1:
		suggestions = coinSuggestions
	case 2:
		suggestions = currencySuggestions
	default:
		return []prompt.Suggest{}
	}

	return prompt.FilterHasPrefix(suggestions, lastArg, true)
}

func main() {
	p := prompt.New(executor, completer)
	p.Run()
}
