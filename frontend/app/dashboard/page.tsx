import { DashboardMetrics } from '@/components/DashboardMetrics';

export default function DashboardPage() {
  return (
    <section className="grid" style={{ gap: '32px' }}>
      <div>
        <h1 className="section-title">Cash flow control center</h1>
        <p className="text-muted">
          Track reminders, invoices, and outstanding balances in one place.
        </p>
      </div>
      <DashboardMetrics />
    </section>
  );
}
