# Price Engine - Roadmap & TODO

This document outlines planned features, enhancements, and improvements for the Price Engine project.

## 🎯 Priority Levels
- **P0** - Critical / High Priority
- **P1** - Important / Medium Priority  
- **P2** - Nice to Have / Low Priority

## 📋 Feature Categories
- **Core Features** - Essential functionality
- **Data Sources** - Exchange and provider integrations
- **Data Quality** - Monitoring, validation, metadata
- **Advanced Features** - Order books, order flow, analytics
- **Infrastructure** - Performance, scalability, observability

---

## 🚀 Core Features

### 1. Multi-Exchange Support with Comparison & Aggregation [P0]

**Goal**: Support top cryptocurrency exchanges with per-exchange views, comparison, and aggregated data.

**Requirements**:
- [ ] Add support for top 10+ cryptocurrency exchanges:
  - [ ] Coinbase Pro / Advanced Trade
  - [ ] Kraken
  - [ ] Bitfinex
  - [ ] OKX (formerly OKEx)
  - [ ] Bybit
  - [ ] KuCoin
  - [ ] Gate.io
  - [ ] Huobi Global
  - [ ] Bitstamp
  - [ ] Gemini

**Features**:
- [ ] **Per-Exchange Data View**
  - [ ] Separate candle streams per exchange
  - [ ] Exchange-specific metadata (fees, limits, etc.)
  - [ ] Exchange health/status indicators
  - [ ] API endpoint: `/exchanges/{exchange}/symbols/{symbol}`

- [ ] **Exchange Comparison**
  - [ ] Side-by-side price comparison across exchanges
  - [ ] Price spread analysis (highest vs lowest)
  - [ ] Volume comparison across exchanges
  - [ ] Arbitrage opportunity detection
  - [ ] API endpoint: `/compare?symbol=BTCUSDT&exchanges=binance,coinbase,kraken`

- [ ] **Aggregated View**
  - [ ] Volume-weighted average price (VWAP) across all exchanges
  - [ ] Consensus price (median/average of all exchanges)
  - [ ] Total volume across all exchanges
  - [ ] Price confidence score based on exchange agreement
  - [ ] API endpoint: `/aggregate?symbol=BTCUSDT`

**Technical Considerations**:
- Extend `Connector` interface to support exchange metadata
- Add exchange-specific configuration (API keys, rate limits)
- Implement exchange health monitoring
- Add exchange ranking/weighting for aggregation
- Handle exchange-specific symbol formats (normalization)

**Proto Updates**:
```protobuf
message ExchangeMetadata {
  string exchange_id = 1;
  string exchange_name = 2;
  bool is_active = 3;
  int64 last_update_ms = 4;
  double fee_rate = 5;
  // ... more metadata
}

message Candle {
  // ... existing fields
  ExchangeMetadata exchange_metadata = 18;
  repeated string source_exchanges = 19; // for aggregated candles
  double consensus_price = 20;
  double price_spread = 21;
}
```

---

## 🔗 Data Sources

### 2. Decentralized Exchange (DEX) Integration [P1]

**Goal**: Integrate top DEX protocols to provide decentralized exchange price data.

**Target DEXs**:
- [ ] Uniswap (V2, V3, V4)
- [ ] PancakeSwap
- [ ] SushiSwap
- [ ] Curve Finance
- [ ] Balancer
- [ ] 1inch
- [ ] dYdX
- [ ] Orca (Solana)

**Features**:
- [ ] Real-time price feeds from DEX pools
- [ ] Liquidity pool data (TVL, depth)
- [ ] Slippage calculations
- [ ] Multi-chain support (Ethereum, BSC, Polygon, Arbitrum, Optimism, Solana, etc.)
- [ ] DEX-specific metadata (pool addresses, AMM type, fees)

**Technical Considerations**:
- Web3 integration (Ethereum JSON-RPC, Solana RPC)
- Smart contract interaction for pool data
- Handle multiple chains/networks
- Gas price tracking for DEX transactions
- MEV (Maximal Extractable Value) data

**Implementation**:
- New package: `pkg/ingest/dex/`
- DEX connector interface similar to CEX connectors
- Support for both on-chain queries and subgraph APIs (The Graph)

---

### 3. Staking Data Integration [P1]

**Goal**: Provide staking data from top staking providers and protocols.

**Target Providers**:
- [ ] Lido (ETH staking)
- [ ] Rocket Pool
- [ ] Coinbase Staking
- [ ] Kraken Staking
- [ ] Binance Staking
- [ ] Staked.us
- [ ] Figment Networks
- [ ] Alluvial (Liquid Collective)

**Features**:
- [ ] Staking APY/APR rates
- [ ] Staking pool sizes (total staked)
- [ ] Validator performance metrics
- [ ] Unstaking/unbonding periods
- [ ] Staking rewards history
- [ ] Liquid staking token prices (stETH, rETH, etc.)

**Data Points**:
- [ ] Current staking rate
- [ ] Historical staking rates
- [ ] Total value locked (TVL)
- [ ] Validator count
- [ ] Slashing events
- [ ] Reward distribution

**Technical Considerations**:
- API integrations with staking providers
- On-chain data for decentralized staking protocols
- Handle different staking mechanisms (PoS, liquid staking, etc.)
- Support multiple assets (ETH, SOL, ATOM, DOT, etc.)

**Implementation**:
- New package: `pkg/ingest/staking/`
- Staking provider connector interface
- Staking-specific data models

---

## 📊 Data Quality & Monitoring

### 4. Latency Tracking & Monitoring [P0]

**Goal**: Track and monitor data latency to ensure real-time performance.

**Features**:
- [ ] **Latency Metrics**
  - [ ] Source latency (exchange → engine)
  - [ ] Processing latency (ingestion → aggregation)
  - [ ] Delivery latency (engine → client)
  - [ ] End-to-end latency

- [ ] **Latency Monitoring**
  - [ ] Real-time latency dashboard
  - [ ] Latency percentiles (p50, p95, p99)
  - [ ] Latency alerts (threshold-based)
  - [ ] Historical latency trends

- [ ] **Latency Metadata**
  - [ ] Include latency in candle metadata
  - [ ] Source timestamp vs processing timestamp
  - [ ] Network round-trip time (RTT)

**Implementation**:
- Add latency tracking to `pkg/ingest/` connectors
- Add latency fields to proto messages
- Implement latency monitoring service
- Add Prometheus metrics for latency
- Create latency dashboard (Grafana)

**Proto Updates**:
```protobuf
message LatencyMetrics {
  int64 source_latency_ms = 1;      // Exchange → Engine
  int64 processing_latency_ms = 2;  // Ingestion → Aggregation
  int64 delivery_latency_ms = 3;    // Engine → Client
  int64 end_to_end_latency_ms = 4;  // Total latency
  int64 source_timestamp_ms = 5;    // Original timestamp from exchange
  int64 received_timestamp_ms = 6;  // When engine received it
  int64 processed_timestamp_ms = 7; // When aggregation completed
}

message Candle {
  // ... existing fields
  LatencyMetrics latency = 22;
}
```

---

### 5. Enhanced Metadata [P0]

**Goal**: Add comprehensive metadata to all data points for regulatory compliance and analytics.

**Metadata Categories**:

- [ ] **Source Metadata**
  - [ ] Exchange/provider name and ID
  - [ ] API endpoint used
  - [ ] Data source type (REST, WebSocket, on-chain)
  - [ ] Source timestamp (original)
  - [ ] Source sequence number/ID

- [ ] **Quality Metadata**
  - [ ] Data quality score
  - [ ] Confidence level
  - [ ] Validation status
  - [ ] Missing data flags
  - [ ] Outlier detection flags

- [ ] **Regulatory Metadata**
  - [ ] MIC (Market Identifier Code) for exchanges
  - [ ] Venue identifier
  - [ ] Market segment
  - [ ] Regulatory jurisdiction
  - [ ] Data classification (public, proprietary, etc.)

- [ ] **Technical Metadata**
  - [ ] Protocol version
  - [ ] Schema version
  - [ ] Compression applied
  - [ ] Encryption status
  - [ ] Data lineage (transformation steps)

- [ ] **Business Metadata**
  - [ ] Trading hours
  - [ ] Market status (open, closed, pre-market, after-hours)
  - [ ] Holiday calendar
  - [ ] Fee structure
  - [ ] Minimum/maximum trade sizes

**Implementation**:
- Extend proto messages with metadata fields
- Create metadata validation layer
- Add metadata enrichment pipeline
- Implement metadata versioning

**Proto Updates**:
```protobuf
message DataMetadata {
  SourceMetadata source = 1;
  QualityMetadata quality = 2;
  RegulatoryMetadata regulatory = 3;
  TechnicalMetadata technical = 4;
  BusinessMetadata business = 5;
}

message Candle {
  // ... existing fields
  DataMetadata metadata = 23;
}
```

---

## 🔍 Advanced Features

### 6. Exchange Wallet Tracking [P2]

**Goal**: Track transactions in/out of top exchanges by monitoring their hot/cold wallets.

**Note**: This may be implemented as a separate microservice/project due to performance and latency considerations.

**Features**:
- [ ] **Wallet Monitoring**
  - [ ] Track hot wallets (operational wallets)
  - [ ] Track cold wallets (storage wallets)
  - [ ] Monitor deposit/withdrawal addresses
  - [ ] Track wallet balances

- [ ] **Transaction Analysis**
  - [ ] Inflow/outflow tracking
  - [ ] Large transaction detection (whale alerts)
  - [ ] Transaction patterns
  - [ ] Exchange reserve tracking

- [ ] **Exchange-Specific Tracking**
  - [ ] Binance wallets
  - [ ] Coinbase wallets
  - [ ] Kraken wallets
  - [ ] Other major exchanges

**Technical Considerations**:
- Blockchain node connections (Bitcoin, Ethereum, etc.)
- Handle multiple blockchains
- High-volume transaction processing
- Real-time vs batch processing trade-offs
- Storage requirements for historical data

**Architecture Options**:
1. **Separate Microservice** (Recommended)
   - Independent service: `wallet-tracker`
   - Communicates with price-engine via events/messaging
   - Own database for wallet/transaction data
   - Can scale independently

2. **Integrated Module**
   - Add to price-engine if performance allows
   - Shared infrastructure
   - Simpler deployment

**Decision Criteria**:
- If latency impact < 10ms → Integrate
- If storage requirements > 100GB → Separate service
- If processing requires different infrastructure → Separate service

**If Separate Project**:
- Project name: `wallet-tracker` or `exchange-monitor`
- Communication: Kafka events, gRPC, or REST API
- Data sharing: Common data format/proto definitions

---

### 7. Level 1, 2, and 3 Order Book Data [P1]

**Goal**: Provide comprehensive order book data (Level 1, 2, 3) and trade data.

**Order Book Levels**:

- [ ] **Level 1 (Top of Book)**
  - [ ] Best bid price
  - [ ] Best ask price
  - [ ] Bid size
  - [ ] Ask size
  - [ ] Last trade price
  - [ ] Last trade size

- [ ] **Level 2 (Market Depth)**
  - [ ] Full order book (all price levels)
  - [ ] Aggregated depth (price levels with total size)
  - [ ] Order book snapshots
  - [ ] Order book updates (deltas)
  - [ ] Depth visualization data

- [ ] **Level 3 (Full Order Book)**
  - [ ] Individual orders (not just price levels)
  - [ ] Order IDs
  - [ ] Order timestamps
  - [ ] Order types (limit, market, stop, etc.)
  - [ ] Order status (active, filled, cancelled)

**Additional Trade Data**:
- [ ] Trade-by-trade feed
- [ ] Trade IDs
- [ ] Trade side (buy/sell)
- [ ] Trade type (market, limit, etc.)
- [ ] Maker/taker information
- [ ] Trade fees

**Features**:
- [ ] Real-time order book streaming
- [ ] Order book snapshots on demand
- [ ] Historical order book data
- [ ] Order book depth analysis
- [ ] Spread calculations
- [ ] Liquidity metrics

**Technical Considerations**:
- High-frequency updates (1000+ updates/second)
- Efficient delta compression
- Memory management for large order books
- WebSocket connections for real-time updates
- Order book reconstruction from deltas

**Implementation**:
- New package: `pkg/orderbook/`
- Order book data structures (red-black tree, sorted maps)
- Delta compression algorithms
- Order book snapshot/update protocol

**Proto Updates**:
```protobuf
message OrderBookLevel {
  double price = 1;
  double size = 2;
  int32 order_count = 3; // For Level 2 aggregated
}

message OrderBook {
  string symbol = 1;
  int64 timestamp_ms = 2;
  repeated OrderBookLevel bids = 3;
  repeated OrderBookLevel asks = 4;
  int64 sequence = 5;
  bool is_snapshot = 6;
}

message Order {
  string order_id = 1;
  double price = 2;
  double size = 3;
  OrderSide side = 4;
  OrderType type = 5;
  int64 timestamp_ms = 6;
  OrderStatus status = 7;
}

service OrderBookStream {
  rpc StreamOrderBook(SubscribeRequest) returns (stream OrderBook);
  rpc StreamOrders(SubscribeRequest) returns (stream Order);
}
```

---

### 8. Order Flow Data & Diagrams [P1]

**Goal**: Provide data for order flow analysis and visualization.

**Order Flow Metrics**:
- [ ] **Volume Analysis**
  - [ ] Buy volume vs sell volume
  - [ ] Volume at price levels
  - [ ] Cumulative volume delta (CVD)
  - [ ] Volume profile

- [ ] **Trade Flow**
  - [ ] Aggressive buy vs aggressive sell
  - [ ] Market buy vs market sell
  - [ ] Limit order fills
  - [ ] Stop order triggers

- [ ] **Order Flow Imbalance**
  - [ ] Bid/ask imbalance
  - [ ] Order flow delta
  - [ ] Liquidity imbalance
  - [ ] Pressure indicators

**Data for Diagrams**:
- [ ] **Footprint Charts**
  - [ ] Volume at each price level
  - [ ] Buy/sell volume split
  - [ ] Delta (buy - sell) per level

- [ ] **Volume Profile**
  - [ ] Volume distribution by price
  - [ ] Value area (VA)
  - [ ] Point of control (POC)

- [ ] **Order Flow Imbalance Charts**
  - [ ] Real-time imbalance visualization
  - [ ] Historical imbalance trends

- [ ] **Cumulative Volume Delta (CVD)**
  - [ ] Running total of buy - sell volume
  - [ ] CVD trends and divergences

**Features**:
- [ ] Real-time order flow streaming
- [ ] Historical order flow data
- [ ] Order flow indicators/calculations
- [ ] Pre-computed diagrams data
- [ ] Custom timeframes for analysis

**Technical Considerations**:
- High-frequency data processing
- Efficient aggregation algorithms
- Real-time calculations
- Historical data storage
- Visualization-friendly data format

**Implementation**:
- New package: `pkg/orderflow/`
- Order flow calculation engine
- Diagram data generators
- Real-time streaming for order flow metrics

**Proto Updates**:
```protobuf
message OrderFlowMetrics {
  string symbol = 1;
  int64 timestamp_ms = 2;
  double buy_volume = 3;
  double sell_volume = 4;
  double volume_delta = 5; // buy - sell
  double cumulative_volume_delta = 6;
  double bid_ask_imbalance = 7;
  map<double, PriceLevelFlow> price_level_flow = 8; // For footprint charts
}

message PriceLevelFlow {
  double price = 1;
  double buy_volume = 2;
  double sell_volume = 3;
  double total_volume = 4;
  double delta = 5;
}

service OrderFlowStream {
  rpc StreamOrderFlow(SubscribeRequest) returns (stream OrderFlowMetrics);
}
```

---

## 🏗️ Infrastructure & Performance

### Additional Improvements

- [ ] **Performance Optimizations**
  - [ ] Connection pooling for exchange APIs
  - [ ] Rate limit management
  - [ ] Caching layer for frequently accessed data
  - [ ] Batch processing optimizations
  - [ ] Memory optimization for high-frequency data

- [ ] **Scalability**
  - [ ] Horizontal scaling support
  - [ ] Load balancing for gRPC services
  - [ ] Distributed aggregation (if needed)
  - [ ] Database sharding strategies

- [ ] **Observability**
  - [ ] Comprehensive logging
  - [ ] Distributed tracing (OpenTelemetry)
  - [ ] Metrics dashboard (Prometheus + Grafana)
  - [ ] Alerting system
  - [ ] Health check endpoints

- [ ] **Data Storage**
  - [ ] Historical data storage (TimescaleDB, InfluxDB)
  - [ ] Data retention policies
  - [ ] Backup and recovery
  - [ ] Data archival strategies

- [ ] **Security**
  - [ ] API authentication/authorization
  - [ ] Rate limiting per client
  - [ ] Encryption in transit and at rest
  - [ ] API key management
  - [ ] Audit logging

---

## 📅 Implementation Phases

### Phase 1: Foundation (Weeks 1-4)
- Multi-exchange support (3-5 top exchanges)
- Per-exchange data views
- Basic aggregation
- Latency tracking
- Enhanced metadata

### Phase 2: Advanced Data Sources (Weeks 5-8)
- DEX integration (top 3-5 DEXs)
- Staking data integration (top 3-5 providers)
- Exchange comparison features
- Advanced aggregation

### Phase 3: Order Books & Flow (Weeks 9-12)
- Level 1 order book
- Level 2 order book
- Order flow metrics
- Order flow diagrams data

### Phase 4: Advanced Features (Weeks 13-16)
- Level 3 order book
- Wallet tracking (separate service if needed)
- Advanced analytics
- Performance optimizations

---

## 📝 Notes

- **Exchange Wallet Tracking**: Decision on separate service vs. integration should be made after performance testing
- **Order Book Levels**: Start with Level 1 and 2, add Level 3 based on demand
- **DEX Integration**: Prioritize based on TVL and trading volume
- **Staking Providers**: Focus on liquid staking protocols first (higher demand)
- **Metadata**: Implement incrementally, starting with most critical fields

---

## 🔗 Related Projects

- **wallet-tracker** (if separate): Exchange wallet monitoring service
- **order-flow-visualizer**: Frontend for order flow diagrams (future)
- **price-engine-client**: Client SDKs for different languages (future)

---

*Last Updated: [Current Date]*
*Version: 1.0*

