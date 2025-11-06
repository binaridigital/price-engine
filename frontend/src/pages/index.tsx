/**
 * Main dashboard page
 */

import { PriceChart } from '@/components/prices/PriceChart';

export default function Dashboard() {
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto py-6 px-4">
          <h1 className="text-3xl font-bold text-gray-900">
            Price Engine Dashboard
          </h1>
        </div>
      </header>

      <main className="max-w-7xl mx-auto py-6 px-4">
        <div className="grid grid-cols-1 gap-6">
          <div className="bg-white rounded-lg shadow p-6">
            <PriceChart symbol="BTCUSDT" intervalMs={1000} />
          </div>
        </div>
      </main>
    </div>
  );
}

