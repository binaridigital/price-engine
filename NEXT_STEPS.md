# Next Steps - Development Roadmap

## 🎯 Immediate Actions (This Week)

### 1. Start Multi-Exchange Support [P0] ⭐ **PRIORITY #1**

**Branch:** `feature/multi-exchange-support`

**Tasks:**
```bash
# Create feature branch
git checkout main
git pull origin main
git checkout -b feature/multi-exchange-support
```

**Implementation Steps:**
1. Add Coinbase connector (`pkg/ingest/coinbase.go`)
   - Follow Binance connector pattern
   - Implement WebSocket connection
   - Handle reconnection logic
   - Add unit tests

2. Add Kraken connector (`pkg/ingest/kraken.go`)
   - Similar pattern to Coinbase
   - Handle Kraken-specific API quirks

3. Extend Connector interface
   - Add metadata methods
   - Add health check methods
   - Add exchange-specific config

4. Update aggregation logic
   - Support per-exchange streams
   - Add exchange metadata to candles
   - Implement basic multi-exchange aggregation

**Estimated Time:** 2-3 weeks

---

### 2. Latency Tracking [P0] ⭐ **PRIORITY #2** (Can be done in parallel)

**Branch:** `feature/latency-tracking`

**Tasks:**
```bash
git checkout -b feature/latency-tracking
```

**Implementation Steps:**
1. Add latency tracking to connectors
   - Record source timestamp
   - Record received timestamp
   - Calculate source latency

2. Update proto file
   - Add `LatencyMetrics` message
   - Add latency field to `Candle`

3. Update aggregator
   - Track processing latency
   - Include latency in output candles

4. Regenerate proto files
   ```bash
   make proto
   ```

**Estimated Time:** 1 week

---

### 3. Enhanced Metadata [P0] ⭐ **PRIORITY #3** (Can be done in parallel)

**Branch:** `feature/enhanced-metadata`

**Tasks:**
```bash
git checkout -b feature/enhanced-metadata
```

**Implementation Steps:**
1. Update proto file
   - Add metadata message types
   - Add metadata field to Candle

2. Implement metadata enrichment
   - Source metadata
   - Quality metadata
   - Regulatory metadata (MIC codes)

3. Add metadata validation
   - Validate ISO codes
   - Validate timestamps
   - Validate data quality

**Estimated Time:** 1-2 weeks

---

## 📋 Development Checklist

### Before Starting Work

- [ ] Read [DEV.md](./DEV.md) for workflow
- [ ] Create feature branch from `main`
- [ ] Update local `main` branch
- [ ] Review related code/patterns
- [ ] Check TODO.md for requirements

### During Development

- [ ] Write code following Go conventions
- [ ] Add unit tests for new code
- [ ] Update documentation
- [ ] Commit frequently with good messages
- [ ] Keep branch updated with main

### Before PR

- [ ] All tests pass
- [ ] Code builds successfully
- [ ] Linting passes
- [ ] Documentation updated
- [ ] Self-review completed
- [ ] Create PR with description

### After PR Approval

- [ ] Merge to main
- [ ] Delete feature branch
- [ ] Update TODO.md (check off items)
- [ ] Tag release if needed

---

## 🏗️ Architecture Decisions Needed

### 1. Multi-Exchange Aggregation Strategy

**Question:** How should we aggregate prices from multiple exchanges?

**Options:**
- A) Volume-weighted average (VWAP)
- B) Median price
- C) Simple average
- D) Configurable (user chooses)

**Recommendation:** Start with VWAP, add options later

---

### 2. Exchange Health Monitoring

**Question:** How do we track exchange health/status?

**Options:**
- A) Simple up/down status
- B) Detailed metrics (latency, error rate, etc.)
- C) Health score (0-100)

**Recommendation:** Start with A, evolve to B

---

### 3. Frontend gRPC-Web Setup

**Question:** How to set up gRPC-Web for Next.js?

**Tasks:**
1. Install gRPC-Web dependencies
2. Set up proxy for gRPC-Web (envoy or similar)
3. Generate TypeScript types from proto
4. Create React hooks for streaming

**Recommendation:** Use Envoy proxy for gRPC-Web translation

---

## 🔧 Setup Tasks

### Backend Setup

1. **Add new exchange connectors:**
   ```bash
   # Create new connector file
   touch pkg/ingest/coinbase.go
   ```

2. **Update main.go:**
   - Add new exchange case
   - Add configuration options

3. **Add tests:**
   ```bash
   touch pkg/ingest/coinbase_test.go
   ```

### Frontend Setup

1. **Initialize Next.js project:**
   ```bash
   cd frontend
   npm install
   ```

2. **Set up gRPC-Web:**
   - Install dependencies
   - Configure proxy
   - Generate TypeScript types

3. **Create initial components:**
   - Price chart
   - Exchange list
   - Dashboard

---

## 📊 Progress Tracking

### Week 1 Goals
- [ ] Multi-exchange branch created
- [ ] Coinbase connector implemented
- [ ] Basic tests added
- [ ] PR created for review

### Week 2 Goals
- [ ] Kraken connector implemented
- [ ] Exchange metadata added
- [ ] Per-exchange streaming working
- [ ] Basic aggregation implemented

### Week 3 Goals
- [ ] Exchange comparison API
- [ ] Frontend foundation
- [ ] Real-time price chart
- [ ] Documentation updated

---

## 🚨 Blockers & Risks

### Potential Blockers

1. **Exchange API Rate Limits**
   - Risk: Hitting rate limits
   - Mitigation: Implement rate limiting, use WebSocket where possible

2. **gRPC-Web Setup Complexity**
   - Risk: Complex setup for Next.js
   - Mitigation: Use Envoy proxy, follow tutorials

3. **Performance with Multiple Exchanges**
   - Risk: High memory/CPU usage
   - Mitigation: Profile early, optimize as needed

### Dependencies

- Exchange API access (may need API keys)
- gRPC-Web proxy setup
- Frontend dependencies

---

## 📚 Resources

### Documentation
- [DEV.md](./DEV.md) - Development workflow
- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [TODO.md](./TODO.md) - Feature roadmap
- [README.md](./README.md) - User documentation

### External Resources
- [gRPC-Web Documentation](https://grpc.io/docs/platforms/web/)
- [Next.js Documentation](https://nextjs.org/docs)
- [Go Best Practices](https://go.dev/doc/effective_go)

---

*Last Updated: [Current Date]*
*Next Review: [Weekly]*

