'use client';

import { useEffect, useState } from 'react';
import { Metrics, getMetrics } from '@/lib/api';
import { getToken } from '@/lib/auth';

export function DashboardMetrics() {
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token = getToken();
    if (!token) {
      setError('Please log in to view metrics.');
      return;
    }

    getMetrics(token)
      .then((data) => setMetrics(data))
      .catch((err) => setError(err instanceof Error ? err.message : 'Failed to load metrics'));
  }, []);

  if (error) {
    return <div className="card">{error}</div>;
  }

  if (!metrics) {
    return <div className="card">Loading metrics...</div>;
  }

  const outstanding = (metrics.outstanding_cents / 100).toFixed(2);

  return (
    <div className="grid two">
      <div className="card">
        <div className="kpi">${outstanding}</div>
        <div className="text-muted">Outstanding receivables</div>
      </div>
      <div className="card">
        <div className="kpi">{metrics.overdue}</div>
        <div className="text-muted">Overdue invoices</div>
      </div>
      <div className="card">
        <div className="kpi">{metrics.upcoming_reminders}</div>
        <div className="text-muted">Reminders scheduled (7d)</div>
      </div>
      <div className="card">
        <div className="kpi">{metrics.clients}</div>
        <div className="text-muted">Active clients</div>
      </div>
    </div>
  );
}
