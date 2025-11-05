# path: Makefile
APP=rt-price-engine
PKG=github.com/binaridigital/price-engine

.PHONY: proto
proto:
	protoc -I proto --go_out=. --go-grpc_out=. proto/price/v1/price.proto

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run:
	go run ./cmd/aggregator \
	  --grpc-addr=:8080 \
	  --symbols=BTCUSDT \
	  --exchanges=binance \
	  --interval=1s

.PHONY: docker
docker:
	docker build -t $(APP):dev .

.PHONY: test-grpc
test-grpc:
	grpcurl -plaintext localhost:8080 list price.v1.PriceStream || true
	grpcurl -plaintext -d '{"symbol":"BTCUSDT","interval_ms":1000}' localhost:8080 price.v1.PriceStream/StreamAggregates
