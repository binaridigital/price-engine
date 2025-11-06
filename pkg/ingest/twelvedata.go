// path: pkg/ingest/twelvedata.go
package ingest

import (
  "context"
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "os"
  "time"

  "github.com/binaridigital/price-engine/pkg/common"
)

// Twelve Data REST price endpoint example: https://api.twelvedata.com/price\?symbol\=EUR/USD\&apikey\=...
// We poll ~250ms for MVP. Replace with WS when you enable it on your plan.
type TwelveData struct {
  apiKey string
  httpc  *http.Client
}

func NewTwelveData() *TwelveData {
  return &TwelveData{
    apiKey: os.Getenv("TWELVEDATA_API_KEY"),
    httpc:  &http.Client{ Timeout: 3 * time.Second },
  }
}

func (t *TwelveData) Name() string { return "twelvedata" }

func (t *TwelveData) Start(ctx context.Context, symbol string) (<-chan common.Trade, <-chan error) {
  out := make(chan common.Trade, 2048)
  errc := make(chan error, 1)
  if t.apiKey == "" {
    close(out); close(errc)
    go func(){ errc <- fmt.Errorf("TWELVEDATA_API_KEY not set") }()
    return out, errc
  }
  symslash := symbol
  if len(symbol) == 6 { symslash = symbol[:3] + "/" + symbol[3:] }

  go func() {
    defer close(out); defer close(errc)
    ticker := time.NewTicker(250 * time.Millisecond)
    defer ticker.Stop()

    for {
      select {
      case <-ctx.Done():
        return
      case <-ticker.C:
        url := fmt.Sprintf("https://api.twelvedata.com/price?symbol=%s&apikey=%s", symslash, t.apiKey)
        req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
        resp, err := t.httpc.Do(req)
        if err != nil { continue }
        body, _ := io.ReadAll(resp.Body)
        _ = resp.Body.Close()

        // Response can be {"price":"1.12345"} or {"price":1.12345}
        var m map[string]json.RawMessage
        if err := json.Unmarshal(body, &m); err != nil { continue }
        var f float64
        if p, ok := m["price"]; ok {
          if err := json.Unmarshal(p, &f); err != nil { continue }
        } else { continue }

        tr := common.Trade{
          Symbol:   normFXSymbol(symbol),
          Price:    f,
          Qty:      1,
          Exchange: "twelvedata",
          TS:       time.Now(),
        }
        select {
        case out <- tr:
        case <-ctx.Done():
          return
        }
      }
    }
  }()
  return out, errc
}
