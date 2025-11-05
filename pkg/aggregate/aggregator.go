// path: pkg/aggregate/aggregator.go
package aggregate

import (
	"context"
	"sync"
	"time"

	"github.com/binaridigital/price-engine/pkg/common"
	pricev1 "github.com/binaridigital/price-engine/proto/price/v1"
)

type window struct {
	startMs int64
	endMs   int64
	open    float64
	high    float64
	low     float64
	close   float64
	vol     float64
	sumPV   float64
	sumV    float64
	count   uint64
	lastTs  int64
	init    bool
}

func Run(ctx context.Context, trades <-chan common.Trade, interval time.Duration) <-chan *pricev1.Candle {
	out := make(chan *pricev1.Candle, 2048)
	windows := make(map[string]*window)
	var mu sync.Mutex

	flush := func(sym string, w *window, final bool) {
		if w == nil || !w.init {
			return
		}
		vwap := 0.0
		if w.sumV > 0 {
			vwap = w.sumPV / w.sumV
		}
		c := &pricev1.Candle{
			Symbol:        sym,
			WindowStartMs: w.startMs,
			WindowEndMs:   w.endMs,
			Open:          w.open,
			High:          w.high,
			Low:           w.low,
			Close:         w.close,
			Volume:        w.vol,
			Vwap:          vwap,
			IsFinal:       final,
			Exchange:      "agg",
			LastTradeTs:   w.lastTs,
			TradeCount:    w.count,
		}
		select {
		case out <- c:
		default:
			select {
			case out <- c:
			case <-ctx.Done():
			}
		}
	}

	ticker := time.NewTicker(interval / 2)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-ticker.C:
				cutoff := now.Add(-10 * time.Millisecond).UnixMilli()
				mu.Lock()
				for sym, w := range windows {
					if w != nil && w.init && w.endMs <= cutoff {
						flush(sym, w, true)
						delete(windows, sym)
					}
				}
				mu.Unlock()
			}
		}
	}()

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-trades:
				if !ok {
					return
				}
				winStart := t.TS.Truncate(interval).UnixMilli()
				winEnd := t.TS.Truncate(interval).Add(interval).UnixMilli()

				mu.Lock()
				w := windows[t.Symbol]
				if w == nil || w.startMs != winStart {
					if w != nil && w.init {
						flush(t.Symbol, w, true)
					}
					w = &window{startMs: winStart, endMs: winEnd}
					windows[t.Symbol] = w
				}
				if !w.init {
					w.open = t.Price
					w.high = t.Price
					w.low = t.Price
					w.init = true
				}
				if t.Price > w.high {
					w.high = t.Price
				}
				if t.Price < w.low {
					w.low = t.Price
				}
				w.close = t.Price
				w.vol += t.Qty
				w.sumPV += t.Price * t.Qty
				w.sumV += t.Qty
				w.count++
				w.lastTs = t.TS.UnixMilli()
				flush(t.Symbol, w, false)
				mu.Unlock()
			}
		}
	}()

	return out
}
