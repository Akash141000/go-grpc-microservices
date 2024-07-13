package main

import (
	"context"
	"encoding/json"
	"fmt"
	"grpc-go/types"
	"math/rand"
	"net/http"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type JSONAPIServer struct {
	listenAddr string
	svc        PriceFetcher
}

func NewJSONAPIServer(listenAddr string, svc PriceFetcher) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc:        svc,
	}
}

func (s *JSONAPIServer) Run() {
	http.HandleFunc("/", makeHttpHandlerFunc(s.handlePriceFetched))
	http.ListenAndServe(s.listenAddr, nil)
}

type RequestId string

type Key string

func makeHttpHandlerFunc(apiFn APIFunc) http.HandlerFunc {
	ctx := context.Background()
	ctx = context.WithValue(ctx, Key("requestId"), rand.Intn(10000000))
	return func(w http.ResponseWriter, r *http.Request) {
		if err := apiFn(ctx, w, r); err != nil {
			// w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
	}
}

func (s *JSONAPIServer) handlePriceFetched(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ticker := r.URL.Query().Get("ticker")

	price, err := s.svc.FetchPrice(ctx, ticker)
	if err != nil {
		return err
	}
	fmt.Println("Price fetched in api", price)
	priceResp := &types.PriceResponse{
		Ticker: ticker,
		Price:  price,
	}
	return writeJSON(w, http.StatusOK, &priceResp)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
