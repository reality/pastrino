package main

import (
	"fmt"
	"os"

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

	FEEDS := []string{
		"https://feeds.theguardian.com/theguardian/world/rss",
		"http://www.washingtonmonthly.com/rss2full_author.xml",
		"https://www.businessinsider.com/rss?op=1&r=US&IR=T",
		"https://www.globenewswire.com/RssFeed/orgclass/1/feedTitle/GlobeNewswire%20-%20News%20about%20Public%20Companies",
		"https://www.globenewswire.com/RssFeed/",
		"https://www.globenewswire.com/RssFeed/country/United%20States/feedTitle/GlobeNewswire%20-%20News%20from%20United%20Statesa",
		"https://www.cbsnews.com/latest/rss/moneywatch",
		"http://feeds.marketwatch.com/marketwatch/topstories/",
		"https://stockstotrade.com/blog/feed/",
		"https://www.reddit.com/r/stocks/.rss",
		"https://www.reddit.com/r/StockMarket/.rss",
		"https://www.reddit.com/r/UKInvesting/.rss",
		"http://economictimes.indiatimes.com/markets/stocks/rssfeeds/2146842.cms",
		"https://www.cnbc.com/id/20409666/device/rss/rss.html?x=1",
	}

	fmt.Println("Loading portfolio...")

	historyFile := os.Args[1]
	p := portfolio.New(historyFile)
	p.String()

	fmt.Println("\nExamining recent news...\n")

	ns := news.New()
	ns.ParseNews(p, FEEDS)

	ns.String()
}
