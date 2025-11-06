# Price Engine

A real-time price aggregation engine that streams trade data from multiple exchanges, aggregates them into candles (OHLCV), and serves them via gRPC. Supports crypto and forex markets with ISO 4217 currency code compliance and regulatory-aligned data structures.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Service](#running-the-service)
- [Command-Line Flags](#command-line-flags)
- [Sample Commands](#sample-commands)
- [Testing](#testing)
- [Docker](#docker)
- [Regulatory Compliance & ISO Standards](#regulatory-compliance--iso-standards)

## Prerequisites

- **Go 1.24+** - [Install Go](https://go.dev/doc/install)
- **Protocol Buffers Compiler (protoc)** - [Install protoc](https://grpc.io/docs/protoc-installation/)
- **protoc-gen-go** and **protoc-gen-go-grpc** plugins:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```
  Make sure `$GOPATH/bin` or `$HOME/go/bin` is in your `PATH`.

- **Optional**: `grpcurl` for testing gRPC endpoints:
  ```bash
  brew install grpcurl  # macOS
  # or download from https://github.com/fullstorydev/grpcurl/releases
  ```

## Installation

1. **Clone the repository** (if not already done):
   ```bash
   git clone <repository-url>
   cd price-engine
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Generate protocol buffer files**:
   ```bash
   make proto
   ```
   Or manually:
   ```bash
   protoc -I proto --go_out=. --go-grpc_out=. proto/price/v1/price.proto
   ```

4. **Build the application** (optional):
   ```bash
   go build -o bin/price-engine ./cmd/aggregator
   ```

## Configuration

### Environment Variables

Some exchange connectors require API keys:

- **TraderMade** (for forex data):
  ```bash
  export TRADERMADE_API_KEY=your_api_key_here
  ```

- **TwelveData** (for forex/stocks data):
  ```bash
  export TWELVEDATA_API_KEY=your_api_key_here
  ```

- **Binance**: No API key required (uses public WebSocket)

## Running the Service

### Quick Start

Run with default settings (Binance, BTCUSDT, 1-second candles):
```bash
make run
```

Or using `go run`:
```bash
go run ./cmd/aggregator
```

### Using the Binary

If you built the binary:
```bash
./bin/price-engine
```

## Command-Line Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--grpc-addr` | string | `:8080` | gRPC server listen address (e.g., `:8080`, `localhost:9090`) |
| `--symbols` | string | `BTCUSDT` | Comma-separated list of symbols to track (e.g., `BTCUSDT,EURUSD,ETHUSDT`) |
| `--exchanges` | string | `binance` | Comma-separated list of exchange connectors: `binance`, `tradermade`, `twelvedata` |
| `--interval` | duration | `1s` | Aggregation window for candles (e.g., `1s`, `5s`, `1m`, `5m`, `1h`) |
| `--kafka-enable` | bool | `false` | Enable publishing aggregated candles to Kafka |
| `--kafka-brokers` | string | `localhost:9092` | Comma-separated list of Kafka broker addresses |
| `--kafka-topic` | string | `agg.candles.v1` | Kafka topic name for publishing candles |

### Flag Details

#### `--grpc-addr`
The address and port where the gRPC server will listen for incoming connections.
- Format: `host:port` or `:port`
- Examples:
  - `:8080` - Listen on all interfaces, port 8080
  - `localhost:9090` - Listen only on localhost, port 9090
  - `0.0.0.0:8080` - Listen on all interfaces, port 8080

#### `--symbols`
Symbols to track from the exchanges. Format depends on the exchange:
- **Crypto (Binance)**: Use exchange format (e.g., `BTCUSDT`, `ETHUSDT`, `BNBUSDT`)
- **Forex (TraderMade/TwelveData)**: Use standard pairs (e.g., `EURUSD`, `GBPUSD`, `USDJPY`)

Multiple symbols can be specified comma-separated. The engine will subscribe to all symbols on all specified exchanges.

#### `--exchanges`
Exchange connectors to use. Available options:
- `binance` - Binance cryptocurrency exchange (no API key required)
- `tradermade` - TraderMade forex data (requires `TRADERMADE_API_KEY`)
- `twelvedata` - TwelveData market data (requires `TWELVEDATA_API_KEY`)

You can specify multiple exchanges comma-separated. The engine will aggregate data from all specified exchanges.

#### `--interval`
Time window for candle aggregation. Accepts Go duration format:
- `1s` - 1 second
- `5s` - 5 seconds
- `30s` - 30 seconds
- `1m` - 1 minute
- `5m` - 5 minutes
- `1h` - 1 hour

The engine creates OHLCV (Open, High, Low, Close, Volume) candles for each interval.

#### `--kafka-enable`
Enable publishing aggregated candles to Kafka. When enabled, all candles are published to the specified Kafka topic.

#### `--kafka-brokers`
Kafka broker addresses. For multiple brokers, use comma-separated list:
- Single broker: `localhost:9092`
- Multiple brokers: `kafka1:9092,kafka2:9092,kafka3:9092`

#### `--kafka-topic`
Kafka topic name where candles will be published. Default is `agg.candles.v1`.

## Sample Commands

### Basic Usage

**1. Run with default settings (Binance, BTCUSDT, 1s candles):**
```bash
go run ./cmd/aggregator
```

**2. Track multiple crypto symbols:**
```bash
go run ./cmd/aggregator \
  --symbols=BTCUSDT,ETHUSDT,BNBUSDT \
  --exchanges=binance \
  --interval=5s
```

**3. Track forex pairs with TraderMade:**
```bash
export TRADERMADE_API_KEY=your_key_here
go run ./cmd/aggregator \
  --symbols=EURUSD,GBPUSD,USDJPY \
  --exchanges=tradermade \
  --interval=1s \
  --grpc-addr=:8080
```

**4. Track multiple symbols from multiple exchanges:**
```bash
export TRADERMADE_API_KEY=your_key_here
go run ./cmd/aggregator \
  --symbols=BTCUSDT,EURUSD \
  --exchanges=binance,tradermade \
  --interval=1s
```

**5. Use custom gRPC port:**
```bash
go run ./cmd/aggregator \
  --grpc-addr=:9090 \
  --symbols=BTCUSDT \
  --exchanges=binance
```

**6. Generate 5-minute candles:**
```bash
go run ./cmd/aggregator \
  --symbols=BTCUSDT \
  --exchanges=binance \
  --interval=5m
```

**7. Enable Kafka publishing:**
```bash
go run ./cmd/aggregator \
  --symbols=BTCUSDT \
  --exchanges=binance \
  --interval=1s \
  --kafka-enable=true \
  --kafka-brokers=localhost:9092 \
  --kafka-topic=price-candles
```

**8. Use custom Kafka brokers:**
```bash
go run ./cmd/aggregator \
  --symbols=BTCUSDT \
  --exchanges=binance \
  --kafka-enable=true \
  --kafka-brokers=kafka1:9092,kafka2:9092 \
  --kafka-topic=agg.candles.v1
```

**9. Track forex with TwelveData:**
```bash
export TWELVEDATA_API_KEY=your_key_here
go run ./cmd/aggregator \
  --symbols=EURUSD,GBPUSD \
  --exchanges=twelvedata \
  --interval=1s
```

**10. Production-like setup with multiple symbols and Kafka:**
```bash
export TRADERMADE_API_KEY=your_key_here
go run ./cmd/aggregator \
  --grpc-addr=:8080 \
  --symbols=BTCUSDT,ETHUSDT,EURUSD,GBPUSD \
  --exchanges=binance,tradermade \
  --interval=1s \
  --kafka-enable=true \
  --kafka-brokers=kafka:9092 \
  --kafka-topic=agg.candles.v1
```

### Using Make Commands

**Run with defaults:**
```bash
make run
```

**Regenerate proto files:**
```bash
make proto
```

**Build Docker image:**
```bash
make docker
```

## Testing

### Test gRPC Endpoint

**List available services:**
```bash
grpcurl -plaintext localhost:8080 list
```

**List methods:**
```bash
grpcurl -plaintext localhost:8080 list price.v1.PriceStream
```

**Stream aggregated candles:**
```bash
grpcurl -plaintext \
  -d '{"symbol":"BTCUSDT","interval_ms":1000}' \
  localhost:8080 \
  price.v1.PriceStream/StreamAggregates
```

**Stream with custom interval:**
```bash
grpcurl -plaintext \
  -d '{"symbol":"EURUSD","interval_ms":5000}' \
  localhost:8080 \
  price.v1.PriceStream/StreamAggregates
```

**Using the Makefile test command:**
```bash
make test-grpc
```

### Expected Output

When streaming, you'll see JSON candle data like:
```json
{
  "symbol": "BTCUSDT",
  "windowStartMs": "1762415854000",
  "windowEndMs": "1762415855000",
  "open": 103204,
  "high": 103214.85,
  "low": 103204,
  "close": 103214.84,
  "volume": 0.02838,
  "vwap": 103210.32,
  "isFinal": true,
  "exchange": "agg",
  "lastTradeTs": "1762415854451",
  "tradeCount": "226",
  "instrumentType": "IT_CRYPTO_SPOT",
  "priceType": "PT_UNSPECIFIED",
  "baseCcy": "",
  "quoteCcy": ""
}
```

## Docker

### Build Docker Image

```bash
make docker
```

Or manually:
```bash
docker build -t rt-price-engine:dev .
```

### Run with Docker

**Basic run:**
```bash
docker run -p 8080:8080 rt-price-engine:dev
```

**With custom flags:**
```bash
docker run -p 8080:8080 rt-price-engine:dev \
  --symbols=BTCUSDT,ETHUSDT \
  --exchanges=binance \
  --interval=5s
```

**With environment variables:**
```bash
docker run -p 8080:8080 \
  -e TRADERMADE_API_KEY=your_key_here \
  rt-price-engine:dev \
  --symbols=EURUSD \
  --exchanges=tradermade
```

**With Kafka:**
```bash
docker run -p 8080:8080 \
  --network kafka-network \
  rt-price-engine:dev \
  --symbols=BTCUSDT \
  --exchanges=binance \
  --kafka-enable=true \
  --kafka-brokers=kafka:9092
```

## Troubleshooting

### Common Issues

1. **"protoc-gen-go: program not found"**
   - Install the protoc plugins and ensure they're in your PATH:
     ```bash
     go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
     export PATH=$PATH:$(go env GOPATH)/bin
     ```

2. **"TRADERMADE_API_KEY not set" or "TWELVEDATA_API_KEY not set"**
   - Make sure you've exported the required environment variables before running

3. **"no connectors configured"**
   - Ensure at least one valid exchange is specified in `--exchanges`

4. **Port already in use**
   - Change the `--grpc-addr` to use a different port
   - Or stop the process using the port: `lsof -ti:8080 | xargs kill`

5. **Connection errors to exchanges**
   - Check your internet connection
   - Verify API keys are correct (for TraderMade/TwelveData)
   - Check if the exchange is accessible from your network

## Architecture

The price engine:
1. **Ingests** trade data from multiple exchange connectors via WebSocket/REST
2. **Merges** trades from all sources into a single stream
3. **Aggregates** trades into OHLCV candles based on the specified interval
4. **Publishes** candles via gRPC streaming API
5. **Optionally** publishes candles to Kafka

## Regulatory Compliance & ISO Standards

### Streaming Data
✅ **Yes, the engine streams data in real-time** via gRPC server-side streaming. Clients subscribe to symbol streams and receive continuous candle updates as they're generated.

### ISO Compliance

The price engine includes several ISO-compliant and regulatory-aligned features:

#### ✅ Implemented Standards

1. **ISO 4217 Currency Codes**
   - Base and quote currency fields (`base_ccy`, `quote_ccy`) use ISO 4217 3-letter codes
   - Validation function ensures only valid ISO codes are accepted
   - Currently supports 24+ major currencies (USD, EUR, GBP, JPY, CHF, AUD, etc.)
   - Automatically splits FX pairs (e.g., "EURUSD" → base: "EUR", quote: "USD")

2. **Instrument Classification**
   - `InstrumentType` enum distinguishes:
     - `IT_CRYPTO_SPOT` - Cryptocurrency spot trades
     - `IT_FX_SPOT` - Foreign exchange spot trades
   - Automatically inferred from symbol format

3. **Price Type Classification**
   - `PriceType` enum supports:
     - `PT_TRADE` - Executed trade prices
     - `PT_BID` - Bid prices
     - `PT_ASK` - Ask prices
     - `PT_MID` - Mid prices (bid+ask)/2

4. **Timestamp Standards**
   - All timestamps in Unix milliseconds (ISO 8601 compatible)
   - Window boundaries clearly defined (`window_start_ms`, `window_end_ms`)
   - Last trade timestamp included for auditability

5. **Data Structure**
   - OHLCV (Open, High, Low, Close, Volume) standard candle format
   - VWAP (Volume-Weighted Average Price) calculation
   - Trade count for transparency
   - Final flag to indicate completed candles

#### ⚠️ Considerations for Full Regulatory Compliance

For production use in regulated financial institutions, consider the following enhancements:

1. **ISO 4217 Coverage**
   - Current implementation includes ~24 major currencies
   - Full ISO 4217 standard includes 170+ currency codes
   - **Recommendation**: Extend `pkg/common/iso4217.go` with complete ISO 4217 list

2. **Additional Regulatory Fields** (not currently included):
   - **ISIN** (International Securities Identification Number) for securities
   - **CFI** (Classification of Financial Instruments) codes
   - **MIC** (Market Identifier Code) for exchange identification
   - **Venue** and **Market Segment** identifiers
   - **Data Source Attribution** (more detailed than current "exchange" field)
   - **Data Quality Flags** (e.g., indicative vs. firm prices)

3. **Timestamp Precision**
   - Currently uses millisecond precision
   - Some regulators may require microsecond or nanosecond precision
   - **Recommendation**: Ensure all timestamps are explicitly UTC

4. **Price Type Annotation**
   - Currently defaults to `PT_UNSPECIFIED` for crypto trades
   - **Recommendation**: Set `PT_TRADE` for executed trades from exchanges

5. **Audit & Lineage**
   - Consider adding:
     - Original trade IDs for traceability
     - Data source timestamps
     - Processing timestamps
     - Version/sequence numbers

6. **Data Validation**
   - Add validation for:
     - Price reasonableness checks
     - Volume validation
     - Timestamp ordering
     - Missing data handling

### Current Compliance Status

| Standard/Requirement | Status | Notes |
|---------------------|--------|-------|
| ISO 4217 Currency Codes | ✅ Partial | 24+ major currencies, can be extended |
| Instrument Classification | ✅ Yes | Crypto and FX spot types defined |
| Price Type Classification | ✅ Yes | Trade/Bid/Ask/Mid types defined |
| Timestamp Format | ✅ Yes | Unix milliseconds (ISO 8601 compatible) |
| Streaming Protocol | ✅ Yes | gRPC server-side streaming |
| OHLCV Standard | ✅ Yes | Standard candle format |
| ISIN/CFI Codes | ❌ No | Not implemented |
| MIC Codes | ❌ No | Not implemented |
| Full ISO 4217 Coverage | ⚠️ Partial | ~24/170+ currencies |
| UTC Timestamp Guarantee | ⚠️ Implicit | Should be explicitly validated |
| Audit Trail | ⚠️ Basic | Exchange name only |

### Recommendations for Financial Institutions

1. **Extend ISO 4217 list** to include all required currencies
2. **Add explicit UTC validation** for all timestamps
3. **Implement MIC codes** for exchange identification
4. **Add data quality flags** and validation rules
5. **Enhance audit trail** with source timestamps and trade IDs
6. **Set appropriate PriceType** based on data source (trade vs. quote)
7. **Consider adding ISIN/CFI** if dealing with securities
8. **Implement data retention policies** for regulatory reporting
9. **Add monitoring/alerting** for data quality issues
10. **Document data lineage** and transformation rules

### Example: Extending ISO 4217

To add more currencies, update `pkg/common/iso4217.go`:

```go
var iso4217 = map[string]struct{}{
  // Existing currencies...
  "USD": {}, "EUR": {}, "GBP": {},
  // Add more as needed:
  "XAU": {}, // Gold (if needed)
  "XAG": {}, // Silver (if needed)
  // ... full ISO 4217 list
}
```

## License

[Add your license information here]

