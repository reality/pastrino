package main

import (
	"fmt"
	"os"

	"reality.rehab/pastrino/config"
	"reality.rehab/pastrino/news"
	"reality.rehab/pastrino/portfolio"
)

func main() {
	// two operations: 1. examine existing stocks and suggest actions
	//								 2. suggest new stocks

	// 1. First we'll look at our balances, and associated:
	//   1. danger scores
	//	 2. hold onto it iveness re recent news

	// 2. Then we will look at recent news, sentiment analysis on it

	// 3. Modify scores and take a look at offsets, suggest movements

	// So I think the modules we'll need are something like:
	// 1. A ledger or current account built from the trading212 export system
	// TODO manually or automatically associate additional keywords (such as e.g. relevant sectors) to (maybe yahoo would be able to help us there)

	config := config.New("./config.json")

	fmt.Println("Loading portfolio...")

	historyFile := os.Args[1]
	p := portfolio.New(historyFile)
	p.String()

	fmt.Println("\nExamining recent news...\n")

	ns := news.New()
	ns.ParseNews(p, config)

	ns.String()
}
