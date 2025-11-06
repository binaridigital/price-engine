/**
 * React hook for real-time price streaming
 */

import { useEffect, useState, useCallback } from 'react';
import { subscribePriceStream } from '@/lib/grpc-client';
import { Candle } from '@/types/price_pb';

interface UsePriceStreamOptions {
  symbol: string;
  intervalMs?: number;
  enabled?: boolean;
}

interface UsePriceStreamReturn {
  data: Candle | null;
  error: Error | null;
  isLoading: boolean;
  isConnected: boolean;
}

export function usePriceStream({
  symbol,
  intervalMs = 1000,
  enabled = true,
}: UsePriceStreamOptions): UsePriceStreamReturn {
  const [data, setData] = useState<Candle | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    if (!enabled || !symbol) {
      return;
    }

    setIsLoading(true);
    setError(null);

    const cleanup = subscribePriceStream(
      symbol,
      intervalMs,
      (candle) => {
        setData(candle);
        setIsLoading(false);
        setIsConnected(true);
        setError(null);
      },
      (err) => {
        setError(err);
        setIsConnected(false);
        setIsLoading(false);
      }
    );

    return () => {
      cleanup();
      setIsConnected(false);
    };
  }, [symbol, intervalMs, enabled]);

  return {
    data,
    error,
    isLoading,
    isConnected,
  };
}

