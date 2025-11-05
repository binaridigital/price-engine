// path: pkg/ingest/exchanges.go
package ingest

import (
	"context"
	"sync"

	"github.com/binaridigital/price-engine/pkg/common"
)

type Connector interface {
	Name() string
	// Start returns a trade channel and an error channel. Both close on ctx.Done().
	Start(ctx context.Context, symbol string) (<-chan common.Trade, <-chan error)
}

// MergeTrades fans in multiple trade channels into one output channel.
func MergeTrades(ctx context.Context, inputs ...<-chan common.Trade) <-chan common.Trade {
	out := make(chan common.Trade, 2048)
	var wg sync.WaitGroup
	wg.Add(len(inputs))
	for _, ch := range inputs {
		go func(c <-chan common.Trade) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case t, ok := <-c:
					if !ok {
						return
					}
					select {
					case out <- t:
					case <-ctx.Done():
						return
					}
				}
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
