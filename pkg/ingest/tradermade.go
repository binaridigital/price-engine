// path: pkg/ingest/tradermade.go
package ingest

import (
  "context"
  "encoding/json"
  "fmt"
  "log"
  "math"
  "os"
  "strings"
  "time"

  "github.com/binaridigital/price-engine/pkg/common"
  "nhooyr.io/websocket"
)

type TraderMade struct {
  apiKey string
}

func NewTraderMade() *TraderMade {
  return &TraderMade{apiKey: os.Getenv("TRADERMADE_API_KEY")}
}

func (t *TraderMade) Name() string { return "tradermade" }

// Start connects to TraderMade WS and streams FX ticks as Trades.
// Requires: export TRADERMADE_API_KEY=xxxx
// Docs (WS streaming & examples): tradermade.com/docs/streaming-data-api
func (t *TraderMade) Start(ctx context.Context, symbol string) (<-chan common.Trade, <-chan error) {
  out := make(chan common.Trade, 2048)
  errc := make(chan error, 1)

  sym := normFXSymbol(symbol)
  if t.apiKey == "" {
    close(out); close(errc)
    go func(){ errc <- fmt.Errorf("TRADERMADE_API_KEY not set") }()
    return out, errc
  }

  // Endpoint per docs. Auth via query param; subscription message sent after connect.
  // If your plan requires a different host/path or header auth, adjust here.
  url := fmt.Sprintf("wss://marketdata.tradermade.com/feedadv?api_key=%s", t.apiKey)
  backoff := 500 * time.Millisecond

  go func() {
    defer close(out); defer close(errc)
    for {
      c, _, err := websocket.Dial(ctx, url, nil)
      if err != nil {
        errc <- fmt.Errorf("tradermade dial: %w", err)
        time.Sleep(backoff)
        backoff = time.Duration(math.Min(float64(backoff*2), float64(30*time.Second)))
        continue
      }
      backoff = 500 * time.Millisecond

      // Subscribe: see docs; pairs are compact (EURUSD, GBPUSD…)
      sub := map[string]interface{}{
        "subscribe": []string{strings.ToUpper(sym)},
      }
      b, _ := json.Marshal(sub)
      _ = c.Write(ctx, websocket.MessageText, b)

      readCtx, cancel := context.WithCancel(ctx)
      readerErr := make(chan error, 1)

      go func() {
        defer cancel()
        for {
          _, data, e := c.Read(readCtx)
          if e != nil {
            readerErr <- fmt.Errorf("tradermade read: %w", e)
            return
          }
          // Message example (field names can vary by plan):
          // {"symbol":"EURUSD","bid":1.12345,"ask":1.12358,"mid":1.123515,"ts":1730869995123}
          var m map[string]interface{}
          if err := json.Unmarshal(data, &m); err != nil {
            continue
          }
          symAny, ok := m["symbol"]
          if !ok { continue }
          ps := normFXSymbol(fmt.Sprint(symAny))

          // Price: use mid if present; else avg(bid,ask); else skip
          var price float64
          if v, ok := m["mid"]; ok {
            price, _ = asFloat(v)
          } else if b, bok := m["bid"]; bok {
            if a, aok := m["ask"]; aok {
              bb, _ := asFloat(b); aa, _ := asFloat(a)
              if bb > 0 && aa > 0 { price = (bb + aa) / 2 }
            }
          }
          if price == 0 { continue }

          // Timestamp
          ts := time.Now()
          if v, ok := m["ts"]; ok {
            if tsms, ok := asFloat(v); ok {
              // ts likely in ms
              ts = time.UnixMilli(int64(tsms))
            }
          }

          tmsg := common.Trade{
            Symbol:   ps,
            Price:    price,
            Qty:      1,            // quote ticks have no volume; use 1 to compute "VWAP" as time-avg
            Exchange: "tradermade",
            TS:       ts,
          }
          select {
          case out <- tmsg:
          case <-readCtx.Done():
            return
          }
        }
      }()

      // Wait for stop or read error -> reconnect
      select {
      case <-ctx.Done():
        _ = c.Close(websocket.StatusNormalClosure, "context done")
        return
      case re := <-readerErr:
        _ = c.Close(websocket.StatusAbnormalClosure, "reconnect")
        log.Printf("tradermade reconnect: %v", re)
        time.Sleep(backoff)
        backoff = time.Duration(math.Min(float64(backoff*2), float64(30*time.Second)))
      }
    }
  }()

  return out, errc
}

func asFloat(v interface{}) (float64, bool) {
  switch x := v.(type) {
  case float64: return x, true
  case json.Number:
    f, err := x.Float64(); if err == nil { return f, true }
  case string:
    var f float64
    if err := json.Unmarshal([]byte(`"`+x+`"`), &f); err == nil { return f, true }
  }
  return 0, false
}
