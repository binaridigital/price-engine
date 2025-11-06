# Price Engine Frontend

React/Next.js frontend for real-time price data visualization and analytics.

## 🏗️ Architecture

```
frontend/
├── src/
│   ├── components/          # React components
│   │   ├── prices/         # Price-related components
│   │   ├── exchanges/      # Exchange comparison & management
│   │   ├── order-flow/     # Order flow visualization
│   │   ├── transactions/   # Transaction tracking
│   │   ├── wallets/        # Wallet monitoring
│   │   └── metrics/        # Metrics & monitoring
│   ├── pages/              # Next.js pages
│   ├── lib/                # Utilities & helpers
│   ├── hooks/              # Custom React hooks
│   ├── types/              # TypeScript types
│   └── styles/             # Global styles
├── public/                 # Static assets
└── package.json
```

## 🚀 Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

```bash
cd frontend
npm install
# or
yarn install
```

### Development

```bash
npm run dev
# or
yarn dev
```

Open [http://localhost:3000](http://localhost:3000)

### Build

```bash
npm run build
npm run start
```

## 📦 Tech Stack

- **Next.js 14+** - React framework
- **TypeScript** - Type safety
- **gRPC-Web** - Real-time streaming
- **Recharts** - Charting library
- **Tailwind CSS** - Styling
- **Zustand/Redux** - State management
- **React Query** - Data fetching

## 🔌 gRPC-Web Integration

The frontend connects to the price engine via gRPC-Web for real-time streaming.

### Setup

1. Generate gRPC-Web client from proto files
2. Configure connection to price engine backend
3. Use streaming hooks for real-time data

## 📱 Pages

- `/` - Dashboard
- `/prices` - Price charts & data
- `/exchanges` - Exchange comparison
- `/order-flow` - Order flow visualization
- `/transactions` - Transaction tracking
- `/wallets` - Wallet monitoring
- `/metrics` - System metrics

## 🎨 Components

### Price Components
- `PriceChart` - Real-time price chart
- `PriceTable` - Price data table
- `PriceComparison` - Multi-exchange comparison

### Exchange Components
- `ExchangeList` - List of exchanges
- `ExchangeComparison` - Side-by-side comparison
- `ExchangeHealth` - Exchange status monitoring

### Order Flow Components
- `OrderBook` - Order book visualization
- `FootprintChart` - Footprint chart
- `VolumeProfile` - Volume profile chart
- `CVDChart` - Cumulative Volume Delta

## 🔄 Real-time Streaming

Uses gRPC-Web for server-side streaming:

```typescript
// Example hook usage
const { data, error } = usePriceStream('BTCUSDT', {
  interval: 1000,
  exchanges: ['binance', 'coinbase']
});
```

## 📝 Development

See [DEV.md](../DEV.md) for development workflow and branch strategy.

## 🧪 Testing

```bash
npm run test
```

## 📚 Documentation

- [Component Documentation](./docs/components.md)
- [API Integration](./docs/api.md)
- [Real-time Streaming](./docs/streaming.md)

