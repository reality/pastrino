package main

import (
	"flag"
	"fmt"

	"reality.rehab/pastrino/config"
	"reality.rehab/pastrino/news"
	"reality.rehab/pastrino/portfolio"
	"reality.rehab/pastrino/tui"
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

	var fancyInterface bool

	flag.StringVar(&config.T212File, "t212", "", "Specify a Trading212 transaction history file.")
	flag.StringVar(&config.WatchListFile, "wl", "", "Specify a watchlist transaction history file (see README for format).")
	flag.BoolVar(&fancyInterface, "fi", false, "Use fancy terminal interface")

	flag.Parse()

	fmt.Println("Loading portfolio...")
	p := portfolio.New(config)

	fmt.Println("\nExamining recent news...")
	ns := news.New()
	ns.ParseNews(p, config)

	if fancyInterface {
		tui.New(p, ns)
	} else {
		p.String()
		ns.String()
	}
}
