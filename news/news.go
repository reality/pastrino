package news

import (
	"fmt"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"reality.rehab/pastrino/portfolio"
)

type NewsSentiment struct {
	newsentries []*NewsEntry
}

type NewsEntry struct {
	id        string
	link      string
	text      string
	mentions  *portfolio.Stonk
	sentiment float64
}

func New() *NewsSentiment {
	ns := &NewsSentiment{}
	ns.newsentries = []*NewsEntry{}
	return ns
}

func (ns *NewsSentiment) ParseNews(p *portfolio.Portfolio, links []string) {
	fp := gofeed.NewParser()
	for _, l := range links {
		feed, err := fp.ParseURL(l)
		if err != nil {
			fmt.Printf("Failed to read news feed: %s", l)
		}

		for _, item := range feed.Items {
			desc := bluemonday.StrictPolicy().Sanitize(item.Description)
			text := bluemonday.StrictPolicy().Sanitize(item.Content)
			for _, s := range p.Stonks {
				// TODO only copmile once
				nPattern, _ := regexp.Compile(`(?i)\b` + s.Name + `\b`)
				tPattern, _ := regexp.Compile(`(?i)\b` + s.Ticker + `\b`)
				descMatches := nPattern.MatchString(text) || nPattern.MatchString(desc)
				textMatches := tPattern.MatchString(text) || tPattern.MatchString(desc)
				if descMatches || textMatches {
					ns.newsentries = append(ns.newsentries, &NewsEntry{
						id:       item.Title,
						link:     item.Link,
						text:     text,
						mentions: s,
					})

					fmt.Printf("News story mentioning a stonk you own: \n")
					fmt.Printf("\t%s\n", item.Title)
					fmt.Printf("\tMentions: %s\n", s.Name)
					fmt.Printf("\tLink: %s\n", item.Link)
					fmt.Println()
				}
			}
		}
	}
}
