package portfolio

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/piquette/finance-go/equity"
	"reality.rehab/pastrino/config"
)

type Portfolio struct {
	Stonks map[string]*Stonk
}

func (p *Portfolio) String() {
	for _, s := range p.Stonks {
		s.String()
	}
}

type Stonk struct {
	Ticker   string
	Name     string
	Quantity float64
	Value    float64
	Currency string
	Keywords []string
	GotValue bool
}

func (s *Stonk) getValue() {
	q, err := equity.Get(s.Ticker)
	if err != nil || q == nil {
		fmt.Printf("Unable to get Value for %s\n", s.Ticker)
	} else {
		if q.RegularMarketPrice != 0 {
			s.Value = s.Quantity * q.RegularMarketPrice
			s.Currency = q.CurrencyID
			s.GotValue = true
		}
	}
}

func (s *Stonk) String() {
	fmt.Printf("%s (%s): \n", s.Ticker, s.Name)
	fmt.Printf("\tQuantity: %f\n", s.Quantity)

	if s.GotValue {
		fmt.Printf("\tValue:    %f%s\n", s.Value, s.Currency)
	} else {
		fmt.Printf("\tValue:    No info found >:(\n")
	}
}

// Here we load in the
// TODO: To extend this to other sites, we will simply have to make different adapters
func New(config *config.Config) *Portfolio {
	p := &Portfolio{}
	p.Stonks = make(map[string]*Stonk)

	if config.T212File != "" {
		fillPortfoliofromT212(p, config.T212File)
	}
	if config.WatchListFile != "" {
		fillPortfoliofromWatchlist(p, config.WatchListFile)
	}

	for _, v := range p.Stonks {
		v.getValue()
	}

	return p
}

// Essentially, we just construct the current
// TODO: In the future, we will probably want to instill the indidivual Stonks with the trading history, etc
// TODO: for code clarity, might help to make it read the headers
// TODO: not sure if there are non market buys or sells (e.g. limits etc), since I don't have any in my file
func fillPortfoliofromT212(p *Portfolio, historyFile string) {
	file, err := os.Open(historyFile)
	if err != nil {
		fmt.Println("Could not load file", err)
	}

	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Could not parse CSV entry", err)
		}
		if record[0] == "Action" || record[3] == "" {
			continue // header
		}

		stonk, ok := p.Stonks[record[3]]
		if !ok {
			stonk = &Stonk{
				Ticker:   record[3],
				Name:     record[4],
				Quantity: 0,
				GotValue: false,
			}
			p.Stonks[record[3]] = stonk
		}

		Quantity, err := strconv.ParseFloat(record[5], 64)

		// by Action
		switch record[0] {
		case "Market buy":
			stonk.Quantity += Quantity
		case "Market sell":
			stonk.Quantity -= Quantity
		}
	}

	// Purge the ones that were all sold
	for _, v := range p.Stonks {
		if v.Quantity == 0 {
			delete(p.Stonks, v.Ticker)
		}
	}
}

// TODO could probably abstract out the field-based parsing boilerplate
func fillPortfoliofromWatchlist(p *Portfolio, watchFile string) {
	file, err := os.Open(watchFile)
	if err != nil {
		fmt.Println("Could not load file", err)
	}

	r := csv.NewReader(file)
	r.Comma = '\t'
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Could not parse TSV entry", err)
		}
		if record[0] == "Ticker" {
			continue // header
		}

		stonk, ok := p.Stonks[record[0]]
		if !ok {
			stonk = &Stonk{
				Ticker:   record[0],
				Name:     record[1],
				Quantity: 0,
				GotValue: false,
			}
			p.Stonks[record[0]] = stonk
		}

		if record[2] != "" {
			stonk.Keywords = strings.Split(record[2], ",")
		}

	}
}
