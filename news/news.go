package news

import (
	"fmt"
	"regexp"

	"github.com/cdipaolo/sentiment"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"reality.rehab/pastrino/config"
	"reality.rehab/pastrino/portfolio"
)

type NewsSentiment struct {
	newsentries []*NewsEntry
}

func (ns *NewsSentiment) String() {
	for _, n := range ns.newsentries {
		n.String()
	}
}

type NewsEntry struct {
	id        string
	link      string
	text      string
	mentions  *portfolio.Stonk
	published string
	sentiment uint8
}

func (n *NewsEntry) String() {
	fmt.Printf("News story mentioning a stonk you own: \n")
	fmt.Printf("\t%s\n", n.id)
	fmt.Printf("\tMentions: %s\n", n.mentions.Name)
	if n.published != "" {
		fmt.Printf("\tPublished: %s\n", n.published)
	}
	fmt.Printf("\tSentiment: %d\n", n.sentiment)
	fmt.Printf("\tLink: %s\n", n.link)
	fmt.Println()
}

func New() *NewsSentiment {
	ns := &NewsSentiment{}
	ns.newsentries = []*NewsEntry{}
	return ns
}

func (ns *NewsSentiment) ParseNews(p *portfolio.Portfolio, config *config.Config) {
	patterns := buildRegexen(p)
	fp := gofeed.NewParser()

	model, err := sentiment.Restore()
	if err != nil {
		panic(fmt.Sprintf("Could not restore sentiment analysis model!\n\t%v\n", err))
	}

	for _, l := range config.Feeds {
		feed, err := fp.ParseURL(l)
		if err != nil {
			fmt.Printf("Failed to read news feed: %s", l)
		}

		for _, item := range feed.Items {
			allText := item.Description + " . " + item.Content
			allText = bluemonday.StrictPolicy().Sanitize(allText)
			for _, s := range p.Stonks {
				match := false
				for _, rg := range patterns[s.Ticker] {
					if !match {
						match = rg.MatchString(allText)
					}
				}

				if match {
					analysis := model.SentimentAnalysis(allText, sentiment.English)

					entry := &NewsEntry{
						id:        item.Title,
						link:      item.Link,
						text:      allText,
						mentions:  s,
						sentiment: analysis.Score,
					}

					if item.Published != "" {
						entry.published = item.Published
					} else if item.Updated != "" { // For some reason, Reddit seems to put the publish date in updated field
						entry.published = item.Updated
					}

					ns.newsentries = append(ns.newsentries, entry)
				}
			}
		}
	}
}

func buildRegexen(p *portfolio.Portfolio) map[string][]*regexp.Regexp {
	m := make(map[string][]*regexp.Regexp)

	for _, s := range p.Stonks {
		rs := []*regexp.Regexp{}

		tPattern, err := regexp.Compile(`(?i)\b` + s.Ticker + `\b`)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("Failed to build match pattern from %s\n", s.Ticker)
		}
		nPattern, err := regexp.Compile(`(?i)\b` + s.Name + `\b`)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("Failed to build match pattern from %s\n", s.Name)
		}
		rs = append(rs, tPattern)
		rs = append(rs, nPattern)

		for _, k := range s.Keywords {
			kPattern, err := regexp.Compile(`(?i)\b` + k + `\b`)
			if err != nil {
				fmt.Println(err)
				fmt.Printf("Failed to build match pattern from keyword %s\n", k)
			}
			rs = append(rs, kPattern)
		}

		m[s.Ticker] = rs
	}

	return m
}
