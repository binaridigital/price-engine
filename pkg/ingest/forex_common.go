// path: pkg/ingest/forex_common.go
package ingest

import (
  "strings"
)

func normFXSymbol(s string) string {
  // Accept "EURUSD" or "EUR/USD" -> "EURUSD"
  s = strings.ToUpper(strings.TrimSpace(s))
  return strings.ReplaceAll(s, "/", "")
}
