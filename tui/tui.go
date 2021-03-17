package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/tview"
	"reality.rehab/pastrino/news"
	"reality.rehab/pastrino/portfolio"
)

func New(p *portfolio.Portfolio, ns *news.NewsSentiment) {
	app := tview.NewApplication()

	newsBox := tview.NewTextView().SetRegions(true).SetDynamicColors(true)

	var output = []string{}
	for i, n := range ns.NewsEntries {
		sString := "[green]Good"
		if n.Sentiment == 0 {
			sString = "[red]Bad"
		}
		title := fmt.Sprintf("[white](%d) %s", i, n.Title)
		content := fmt.Sprintf("[yellow]\tMentions: %s (%s)\n\tSentiment: %s\n\tLink: %s", n.Mentions.Name, n.Mentions.Ticker, sString, n.Link)
		output = append(output, title)
		output = append(output, content)
	}

	newsBox.SetText(strings.Join(output, "\n")).Highlight(strconv.Itoa(1))
	newsBox.Box.SetTitle("News").SetBorder(true)

	stonkBox := tview.NewTable().SetBorders(true)
	stonkBox.Box.SetTitle("Stonks").SetBorder(true)

	stonkBox.SetCell(0, 0, tview.NewTableCell("[yellow]Ticker"))
	stonkBox.SetCell(0, 1, tview.NewTableCell("[yellow]Name"))
	stonkBox.SetCell(0, 2, tview.NewTableCell("[yellow]Quantity"))
	stonkBox.SetCell(0, 3, tview.NewTableCell("[yellow]Value"))

	c := 1
	for _, s := range p.Stonks {
		tickerText := fmt.Sprintf("[white]%s", s.Ticker)
		nameText := fmt.Sprintf("[white]%s", s.Name)
		quantityText := fmt.Sprintf("[white]%f", s.Quantity)

		valueText := "[white]Unknown"
		if s.GotValue {
			valueText = fmt.Sprintf("[white]%f %s", s.Value, s.Currency)
		}

		stonkBox.SetCell(c, 0, tview.NewTableCell(tickerText))
		stonkBox.SetCell(c, 1, tview.NewTableCell(nameText))
		stonkBox.SetCell(c, 2, tview.NewTableCell(quantityText))
		stonkBox.SetCell(c, 3, tview.NewTableCell(valueText))

		c++
	}

	/*for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			stonkBox.SetCell(x, y,
				tview.NewTableCell("hi").
					SetAlign(tview.AlignCenter))
		}
	}*/

	flex := tview.NewFlex().
		AddItem(stonkBox, 0, 2, false).
		AddItem(newsBox, 0, 2, false)

	err := app.SetRoot(flex, true).Run()
	if err != nil {
		panic(err)
	}
}
