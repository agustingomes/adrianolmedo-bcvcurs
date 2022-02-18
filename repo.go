package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const URL = "http://www.bcv.org.ve/"

var curs = Currencies{
	{
		ID:     "euro",
		Iso:    "EUR",
		Symbol: "€",
	},
	{
		ID:     "yuan",
		Iso:    "CNY",
		Symbol: "¥",
	},
	{
		ID:     "lira",
		Iso:    "TRY",
		Symbol: "₺",
	},
	{
		ID:     "rublo",
		Iso:    "RUB",
		Symbol: "₽",
	},
	{
		ID:     "dolar",
		Iso:    "USD",
		Symbol: "$",
	},
}

func getAll() (Currencies, error) {
	body, err := bodyFromURL(URL)
	if err != nil {
		cfg.Logger.Log("level", "error", "msg", err.Error(), "caller", logCaller(1))
		return nil, err
	}
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		cfg.Logger.Log("level", "error", "msg", err.Error(), "caller", logCaller(1))
		return nil, err
	}

	for _, cur := range curs {
		value, err := findValueByID(cur.ID, doc)
		if err != nil {
			cfg.Logger.Log("level", "error", "msg", err.Error(), "caller", logCaller(1))
			return nil, err
		}
		cur.Value = value
	}
	return curs, nil
}

func getUnique(key int) (Currency, error) {
	body, err := bodyFromURL(URL)
	if err != nil {
		cfg.Logger.Log("level", "error", "msg", err.Error(), "caller", logCaller(1))
		return Currency{}, err
	}
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		cfg.Logger.Log("level", "error", "msg", err.Error(), "caller", logCaller(1))
		return Currency{}, err
	}

	value, err := findValueByID(curs[key].ID, doc)
	if err != nil {
		cfg.Logger.Log("level", "error", "msg", err.Error(), "caller", logCaller(1))
		return Currency{}, err
	}

	return Currency{
		Value:  value,
		Iso:    curs[key].Iso,
		Symbol: curs[key].Symbol,
	}, nil
}

func bodyFromURL(url string) (body io.ReadCloser, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code from source %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func findValueByID(id string, doc *goquery.Document) (float64, error) {
	s := doc.Find("div[id='" + id + "']").Find("strong").Text()
	if s == "" {
		return 0, ErrGettingData
	}

	s = strings.TrimSpace(s)
	s = strings.Replace(s, ",", ".", -1)
	return strconv.ParseFloat(s, 64)
}
