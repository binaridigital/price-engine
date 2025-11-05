// path: pkg/grpcapi/server.go
package grpcapi

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pricev1 "github.com/binaridigital/price-engine/proto/price/v1"
)

type Hub struct {
	mu   sync.RWMutex
	subs map[string]map[chan *pricev1.Candle]struct{}
}

func NewHub() *Hub {
	return &Hub{subs: make(map[string]map[chan *pricev1.Candle]struct{})}
}

func (h *Hub) Publish(c *pricev1.Candle) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	group := h.subs[c.Symbol]
	for ch := range group {
		select {
		case ch <- c:
		default:
		}
	}
}

func (h *Hub) Subscribe(symbol string) (chan *pricev1.Candle, func()) {
	ch := make(chan *pricev1.Candle, 1024)
	h.mu.Lock()
	if _, ok := h.subs[symbol]; !ok {
		h.subs[symbol] = make(map[chan *pricev1.Candle]struct{})
	}
	h.subs[symbol][ch] = struct{}{}
	h.mu.Unlock()
	unsub := func() {
		h.mu.Lock()
		if group, ok := h.subs[symbol]; ok {
			delete(group, ch)
			close(ch)
			if len(group) == 0 {
				delete(h.subs, symbol)
			}
		}
		h.mu.Unlock()
	}
	return ch, unsub
}

type Server struct {
	pricev1.UnimplementedPriceStreamServer
	hub              *Hub
	engineIntervalMs int64
}

func NewServer(hub *Hub, engineInterval time.Duration) *Server {
	return &Server{hub: hub, engineIntervalMs: engineInterval.Milliseconds()}
}

func (s *Server) StreamAggregates(req *pricev1.SubscribeRequest, stream pricev1.PriceStream_StreamAggregatesServer) error {
	if req.GetSymbol() == "" {
		return errors.New("symbol required")
	}
	if req.GetIntervalMs() != 0 && req.GetIntervalMs() != s.engineIntervalMs {
		return errors.New("requested interval not supported by current engine instance")
	}
	ch, unsub := s.hub.Subscribe(req.GetSymbol())
	defer unsub()

	for {
		select {
		case <-stream.Context().Done():
			return context.Canceled
		case c := <-ch:
			if err := stream.Send(c); err != nil {
				return err
			}
		}
	}
}

func Serve(addr string, s *Server) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pricev1.RegisterPriceStreamServer(grpcServer, s)
	reflection.Register(grpcServer)
	return grpcServer.Serve(lis)
}
