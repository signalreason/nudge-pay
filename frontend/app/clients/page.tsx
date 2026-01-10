import { ClientsPanel } from '@/components/ClientsPanel';

export default function ClientsPage() {
  return (
    <section className="grid" style={{ gap: '32px' }}>
      <div>
        <h1 className="section-title">Client roster</h1>
        <p className="text-muted">Keep billing contacts tidy and ready to nudge.</p>
      </div>
      <ClientsPanel />
    </section>
  );
}
