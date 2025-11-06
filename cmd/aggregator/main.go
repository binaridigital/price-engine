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

  "github.com/binaridigital/price-engine/pkg/aggregate"
  "github.com/binaridigital/price-engine/pkg/common"
  "github.com/binaridigital/price-engine/pkg/grpcapi"
  "github.com/binaridigital/price-engine/pkg/ingest"
  pkafka "github.com/binaridigital/price-engine/pkg/kafka"
)

func main() {
  grpcAddr   := flag.String("grpc-addr", ":8080", "gRPC listen address")
  symbolsCSV := flag.String("symbols", "BTCUSDT", "comma-separated symbols (e.g., BTCUSDT,EURUSD)")
  exchanges  := flag.String("exchanges", "binance", "comma-separated connectors: binance,tradermade,twelvedata")
  interval   := flag.Duration("interval", time.Second, "aggregation window (e.g., 1s)")
  // Kafka (optional)
  kafkaEnable  := flag.Bool("kafka-enable", false, "publish to Kafka")
  kafkaBrokers := flag.String("kafka-brokers", "localhost:9092", "kafka brokers (comma)")
  kafkaTopic   := flag.String("kafka-topic", "agg.candles.v1", "kafka topic")
  flag.Parse()

  ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
  defer cancel()

  // Build connectors
  var conns []ingest.Connector
  for _, name := range strings.Split(*exchanges, ",") {
    name = strings.TrimSpace(strings.ToLower(name))
    switch name {
    case "binance":
      conns = append(conns, ingest.NewBinance())
    case "tradermade":
      conns = append(conns, ingest.NewTraderMade())
    case "twelvedata":
      conns = append(conns, ingest.NewTwelveData())
    case "", "none":
      // skip
    default:
      log.Printf("unknown connector: %s (skipped)", name)
    }
  }
  if len(conns) == 0 {
    log.Fatal("no connectors configured")
  }

  // Start ingestion
  var tradeChans []<-chan common.Trade
  var errChans []<-chan error
  syms := strings.Split(*symbolsCSV, ",")
  for _, c := range conns {
    for _, s := range syms {
      tc, ec := c.Start(ctx, strings.TrimSpace(s))
      tradeChans = append(tradeChans, tc)
      errChans = append(errChans, ec)
    }
  }
  // Error logger
  for _, ec := range errChans {
    go func(ch <-chan error) {
      for e := range ch {
        if e != nil { log.Printf("ingest error: %v", e) }
      }
    }(ec)
  }

  // Merge + aggregate
  merged := ingest.MergeTrades(ctx, tradeChans...)
  candles := aggregate.Run(ctx, merged, *interval)

  // Hub + gRPC
  hub := grpcapi.NewHub()
  go func() {
    log.Printf("gRPC listening on %s", *grpcAddr)
    if err := grpcapi.Serve(*grpcAddr, grpcapi.NewServer(hub, *interval)); err != nil {
      log.Fatalf("grpc serve: %v", err)
    }
  }()

  // Kafka (optional)
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
      if !ok { return }
      hub.Publish(c)
      if producer != nil {
        _ = producer.Publish(ctx, c)
      }
    }
  }
}
