// path: pkg/common/types.go
package common

import "time"

type Trade struct {
	Symbol   string
	Price    float64
	Qty      float64
	Exchange string
	TS       time.Time
}
