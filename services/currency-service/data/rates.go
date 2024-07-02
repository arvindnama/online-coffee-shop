package data

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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
	e.logger.Debug("Fetching Rates from Central Bank Exchange")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(centralBankExchangeRateUrl)

	if err != nil {
		e.logger.Error(fmt.Sprintf("Error Fetching Rates from Central Bank Exchange:: %#v\n", err))
		return err
	}

	if resp.StatusCode != http.StatusOK {
		e.logger.Error(fmt.Sprintf("Error Fetching Rates from Central Bank Exchange:: %#v\n", resp))
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
	e.rates["EUR"] = 1

	e.logger.Debug(fmt.Sprintf("Rate Cache:: %#v\n", e.rates))

	return nil
}

func (e *ExchangeRates) GetRate(base string, dest string) (float64, error) {
	br, err := e.getRateEuroAsBase(base)

	if err != nil {
		return 0, err
	}

	dr, err := e.getRateEuroAsBase(dest)

	if err != nil {
		return 0, err
	}

	return br * dr, nil
}

func (e *ExchangeRates) getRateEuroAsBase(dest string) (float64, error) {
	dr, ok := e.rates[dest]

	if !ok {
		return 0, fmt.Errorf("rate not found for currency %v", dest)
	}

	return dr, nil
}

// MonitorRates check the exchange server and updates the cache
func (e *ExchangeRates) MonitorRates(interval time.Duration) chan struct{} {
	res := make(chan struct{})

	//[learning]: NewTicket is like setInterval
	// it returns a channel and notifies the channel every intervalset
	timer := time.NewTicker(interval)

	go func() {
		for {
			//[learning]: Select is used to wait on channel (it does close the channel once done)
			select {
			case <-timer.C:
				// randomize the rates.
				for k, v := range e.rates {
					// Change only 10% of value
					change := (rand.Float64() / 10)

					// + / - change

					direction := rand.Intn(1)

					if direction == 0 {
						// reduce
						change = 1 - change
					} else {
						// increase
						change = 1 + change
					}
					e.rates[k] = v * change
				}
				res <- struct{}{}
			}
		}
	}()
	return res
}
