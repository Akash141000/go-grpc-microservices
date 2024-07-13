package main

import (
	"context"
	"fmt"
)

var mockPrices = map[string]float64{
	"ETH": 20,
	"BTC": 10,
}

type PriceFetcher interface {
	FetchPrice(context.Context, string) (float64, error)
}

type priceFetcher struct{}

func (p *priceFetcher) FetchPrice(ctx context.Context, ticker string) (float64, error) {
	return MockFetchPrice(ctx, ticker)
}

func MockFetchPrice(ctx context.Context, ticker string) (float64, error) {
	price, ok := mockPrices[ticker]
	if !ok {
		return price, fmt.Errorf("the given ticket %v is not supported", ticker)
	}

	fmt.Println("Mock price fetched", price)
	return price, nil
}
