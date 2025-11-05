// path: cmd/aggregator/main.go
package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/binari-digital/price-engine/pkg/aggregate"
	"github.com/binari-digital/price-engine/pkg/common"
	"github.com/binari-digital/price-engine/pkg/grpcapi"
	"github.com/binari-digital/price-engine/pkg/ingest"
	pkafka "github.com/binari-digital/price-engine/pkg/kafka"
)

func main() {
	var (
		grpcAddr     = flag.String("grpc-addr", ":8080", "gRPC listen addr")
		symbolsCSV   = flag.String("symbols", "BTCUSDT", "comma-separated symbols")
		exchangesCSV = flag.String("exchanges", "binance", "comma-separated exchanges (binance)")
		intervalStr  = flag.String("interval", "1s", "aggregation interval (e.g., 1s, 100ms)")
		kafkaEnable  = flag.Bool("kafka-enable", false, "enable Kafka publishing")
		kafkaBrokers = flag.String("kafka-brokers", "localhost:9092", "comma-separated broker list")
		kafkaTopic   = flag.String("kafka-topic", "agg.candles.v1", "Kafka topic")
	)
	flag.Parse()

	interval, err := time.ParseDuration(*intervalStr)
	if err != nil || interval <= 0 {
		log.Fatalf("invalid interval: %v", *intervalStr)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var tradeChans []<-chan common.Trade
	var errChans []<-chan error

	exchanges := strings.Split(*exchangesCSV, ",")
	symbols := strings.Split(*symbolsCSV, ",")

	for _, ex := range exchanges {
		switch strings.ToLower(strings.TrimSpace(ex)) {
		case "binance":
			conn := ingest.NewBinance()
			for _, sym := range symbols {
				tch, ech := conn.Start(ctx, strings.TrimSpace(sym))
				tradeChans = append(tradeChans, tch)
				errChans = append(errChans, ech)
			}
		default:
			log.Fatalf("unsupported exchange: %s", ex)
		}
	}

	for _, ec := range errChans {
		go func(c <-chan error) {
			for {
				select {
				case <-ctx.Done():
					return
				case e, ok := <-c:
					if !ok {
						return
					}
					log.Printf("ingest error: %v", e)
				}
			}
		}(ec)
	}

	merged := ingest.MergeTrades(ctx, tradeChans...)
	candles := aggregate.Run(ctx, merged, interval)

	hub := grpcapi.NewHub()
	go func() {
		log.Printf("gRPC listening on %s", *grpcAddr)
		if err := grpcapi.Serve(*grpcAddr, grpcapi.NewServer(hub, interval)); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	var producer *pkafka.Producer
	if *kafkaEnable {
		producer = pkafka.NewProducer(strings.Split(*kafkaBrokers, ","), *kafkaTopic)
		defer producer.Close()
		log.Printf("Kafka enabled -> topic=%s brokers=%s", *kafkaTopic, *kafkaBrokers)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case c, ok := <-candles:
			if !ok {
				return
			}
			hub.Publish(c)
			if producer != nil {
				_ = producer.Publish(ctx, c)
			}
		}
	}
}
