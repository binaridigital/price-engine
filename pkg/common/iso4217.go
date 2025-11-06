// path: pkg/common/iso4217.go
package common

// Minimal ISO 4217 allow-list for fast validation.
// Extend as needed; keep uppercase 3-letter codes.
// Sources: ISO 4217 overview & tables. (See docs)
// This is intentionally compact for MVP.
var iso4217 = map[string]struct{}{
  "USD": {}, "EUR": {}, "GBP": {}, "JPY": {}, "CHF": {}, "AUD": {}, "NZD": {}, "CAD": {},
  "SEK": {}, "NOK": {}, "DKK": {}, "PLN": {}, "CZK": {}, "HUF": {}, "TRY": {}, "ILS": {},
  "ZAR": {}, "MXN": {}, "BRL": {}, "RUB": {}, "HKD": {}, "SGD": {}, "CNH": {}, "CNY": {},
  "INR": {}, "KRW": {}, "TWD": {}, "THB": {}, "MYR": {}, "IDR": {}, "PHP": {}, "AED": {},
}

func IsISO4217(s string) bool {
  _, ok := iso4217[s]
  return ok
}

// SplitFX tries to split a compact pair like "EURUSD" -> "EUR","USD", true.
// Returns empty strings/false if not a known ISO pair.
func SplitFX(sym string) (string, string, bool) {
  if len(sym) != 6 { return "", "", false }
  b := sym[:3]
  q := sym[3:]
  if IsISO4217(b) && IsISO4217(q) { return b, q, true }
  return "", "", false
}
