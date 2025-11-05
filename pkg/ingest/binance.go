// path: pkg/ingest/binance.go
package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"nhooyr.io/websocket"

	"github.com/binari-digital/price-engine/pkg/common"
)

type binanceConnector struct{}

func NewBinance() Connector { return &binanceConnector{} }
func (b *binanceConnector) Name() string { return "binance" }

type binanceTradeMsg struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	TradeID   int64  `json:"t"`
	Price     string `json:"p"`
	Qty       string `json:"q"`
	BuyerID   int64  `json:"b"`
	SellerID  int64  `json:"a"`
	TradeTime int64  `json:"T"`
	IsMaker   bool   `json:"m"`
}

func (b *binanceConnector) Start(ctx context.Context, symbol string) (<-chan common.Trade, <-chan error) {
	trades := make(chan common.Trade, 2048)
	errc := make(chan error, 1)

	go func() {
		defer close(trades)
		defer close(errc)

		sym := strings.ToLower(symbol)
		url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@trade", sym)

		backoff := time.Millisecond * 500
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			c, _, err := websocket.Dial(ctx, url, nil)
			if err != nil {
				errc <- fmt.Errorf("binance dial: %w", err)
				time.Sleep(backoff)
				backoff = time.Duration(math.Min(float64(backoff*2), float64(30*time.Second)))
				continue
			}
			backoff = time.Millisecond * 500

			readCtx, cancel := context.WithCancel(ctx)
			readerErr := make(chan error, 1)
			go func() {
				defer cancel()
				for {
					typ, data, rerr := c.Read(readCtx)
					if rerr != nil {
						readerErr <- rerr
						return
					}
					if typ != websocket.MessageText {
						continue
					}
					var m binanceTradeMsg
					if err := json.Unmarshal(data, &m); err != nil {
						errc <- fmt.Errorf("binance unmarshal: %w", err)
						continue
					}
					price, _ := strconv.ParseFloat(m.Price, 64)
					qty, _ := strconv.ParseFloat(m.Qty, 64)
					t := common.Trade{
						Symbol:   strings.ToUpper(m.Symbol),
						Price:    price,
						Qty:      qty,
						Exchange: "binance",
						TS:       time.UnixMilli(m.TradeTime),
					}
					select {
					case trades <- t:
					case <-readCtx.Done():
						return
					}
				}
			}()

			select {
			case <-ctx.Done():
				_ = c.Close(websocket.StatusNormalClosure, "context done")
				return
			case re := <-readerErr:
				_ = c.Close(websocket.StatusAbnormalClosure, "reconnect")
				log.Printf("binance reconnect: %v", re)
				time.Sleep(backoff)
				backoff = time.Duration(math.Min(float64(backoff*2), float64(30*time.Second)))
			}
		}
	}()

	return trades, errc
}
