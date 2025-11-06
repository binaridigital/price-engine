/**
 * Real-time price chart component
 */

import { usePriceStream } from '@/hooks/usePriceStream';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

interface PriceChartProps {
  symbol: string;
  intervalMs?: number;
}

export function PriceChart({ symbol, intervalMs = 1000 }: PriceChartProps) {
  const { data, error, isLoading, isConnected } = usePriceStream({
    symbol,
    intervalMs,
  });

  // Transform data for chart
  const chartData = data ? [{
    time: new Date(data.getWindowStartMs()).toLocaleTimeString(),
    price: data.getClose(),
    open: data.getOpen(),
    high: data.getHigh(),
    low: data.getLow(),
  }] : [];

  if (error) {
    return (
      <div className="p-4 bg-red-50 border border-red-200 rounded">
        <p className="text-red-800">Error: {error.message}</p>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="p-4 bg-gray-50 border border-gray-200 rounded">
        <p>Connecting to price stream...</p>
      </div>
    );
  }

  return (
    <div className="w-full h-96 p-4">
      <div className="mb-4">
        <h2 className="text-xl font-bold">{symbol}</h2>
        <div className="flex items-center gap-2">
          <span className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
          <span className="text-sm text-gray-600">
            {isConnected ? 'Connected' : 'Disconnected'}
          </span>
        </div>
      </div>
      
      {data && (
        <div className="mb-4">
          <div className="text-3xl font-bold">${data.getClose().toFixed(2)}</div>
          <div className="text-sm text-gray-600">
            Volume: {data.getVolume().toFixed(4)} | 
            VWAP: ${data.getVwap().toFixed(2)}
          </div>
        </div>
      )}

      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={chartData}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="time" />
          <YAxis domain={['auto', 'auto']} />
          <Tooltip />
          <Line 
            type="monotone" 
            dataKey="price" 
            stroke="#8884d8" 
            strokeWidth={2}
            dot={false}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}

