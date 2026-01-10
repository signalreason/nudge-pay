import { RemindersPanel } from '@/components/RemindersPanel';

export default function RemindersPage() {
  return (
    <section className="grid" style={{ gap: '32px' }}>
      <div>
        <h1 className="section-title">Reminder timeline</h1>
        <p className="text-muted">See what is queued and what was already sent.</p>
      </div>
      <RemindersPanel />
    </section>
  );
}
