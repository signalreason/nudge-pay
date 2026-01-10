import { InvoicesPanel } from '@/components/InvoicesPanel';

export default function InvoicesPage() {
  return (
    <section className="grid" style={{ gap: '32px' }}>
      <div>
        <h1 className="section-title">Invoice runway</h1>
        <p className="text-muted">Schedule reminders right when invoices go out.</p>
      </div>
      <InvoicesPanel />
    </section>
  );
}
