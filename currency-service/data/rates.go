package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
)

const centralBankExchangeRateUrl = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

type ExchangeRates struct {
	logger hclog.Logger
	rates  map[string]float64
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

func NewExchangeRates(logger hclog.Logger) (*ExchangeRates, error) {
	e := &ExchangeRates{logger, map[string]float64{}}
	e.getRates()
	return e, nil
}

func (e *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get(centralBankExchangeRateUrl)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected Status code 200 got %v", resp.StatusCode)

	}
	defer resp.Body.Close()
	md := &Cubes{}

	xml.NewDecoder(resp.Body).Decode(&md)

	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}

		e.rates[c.Currency] = r
	}

	return nil
}
