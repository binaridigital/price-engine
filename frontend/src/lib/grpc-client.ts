/**
 * gRPC-Web client for connecting to price engine backend
 */

import { PriceStreamClient } from '@/types/price_grpc_web_pb';
import { SubscribeRequest, Candle } from '@/types/price_pb';

// gRPC-Web client instance
let client: PriceStreamClient | null = null;

/**
 * Initialize gRPC-Web client
 */
export function initGrpcClient(url?: string): PriceStreamClient {
  const grpcUrl = url || process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  
  if (!client) {
    client = new PriceStreamClient(grpcUrl);
  }
  
  return client;
}

/**
 * Get gRPC client instance
 */
export function getGrpcClient(): PriceStreamClient {
  if (!client) {
    return initGrpcClient();
  }
  return client;
}

/**
 * Subscribe to price stream
 */
export function subscribePriceStream(
  symbol: string,
  intervalMs: number,
  onData: (candle: Candle) => void,
  onError: (error: Error) => void
): () => void {
  const grpcClient = getGrpcClient();
  const request = new SubscribeRequest();
  request.setSymbol(symbol);
  request.setIntervalMs(intervalMs);

  const stream = grpcClient.streamAggregates(request, {});
  
  stream.on('data', (candle: Candle) => {
    onData(candle);
  });
  
  stream.on('error', (error: Error) => {
    onError(error);
  });
  
  // Return cleanup function
  return () => {
    stream.cancel();
  };
}

