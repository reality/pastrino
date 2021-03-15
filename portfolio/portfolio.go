package portfolio

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/piquette/finance-go/equity"
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
	quantity float64
	value    float64
	currency string
	gotValue bool
}

func (s *Stonk) getValue() {
	q, err := equity.Get(s.Ticker)
	if err != nil || q == nil {
		fmt.Printf("Unable to get value for %s\n", s.Ticker)
	} else {
		if q.RegularMarketPrice != 0 {
			s.value = s.quantity * q.RegularMarketPrice
			s.currency = q.CurrencyID
			s.gotValue = true
		}
	}
}

func (s *Stonk) String() {
	fmt.Printf("%s (%s): \n", s.Ticker, s.Name)
	fmt.Printf("\tQuantity: %f\n", s.quantity)

	if s.gotValue {
		fmt.Printf("\tValue:    %f%s\n", s.value, s.currency)
	} else {
		fmt.Printf("\tValue:    No info found >:(\n")
	}
}

// Here we load in the
// TODO: To extend this to other sites, we will simply have to make different adapters
func New(historyFile string) *Portfolio {
	p := &Portfolio{}
	p.Stonks = make(map[string]*Stonk)

	fillPortfoliofromT212(p, historyFile)

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
				quantity: 0,
				gotValue: false,
			}
			p.Stonks[record[3]] = stonk
		}

		quantity, err := strconv.ParseFloat(record[5], 64)

		// by Action
		switch record[0] {
		case "Market buy":
			stonk.quantity += quantity
		case "Market sell":
			stonk.quantity -= quantity
		}
	}

	// Purge the ones that were all sold
	for _, v := range p.Stonks {
		if v.quantity == 0 {
			delete(p.Stonks, v.Ticker)
		}
	}
}
