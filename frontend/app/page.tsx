import Link from 'next/link';

export default function HomePage() {
  return (
    <section className="hero">
      <div>
        <h1>Chase invoices without chasing clients.</h1>
        <p>
          NudgePay schedules tasteful reminders, logs every touch, and keeps your
          cash flow predictable. Built for freelancers and small agencies who
          bill by the project.
        </p>
        <div style={{ display: 'flex', gap: '12px', marginTop: '20px' }}>
          <Link className="button" href="/signup">
            Start free trial
          </Link>
          <Link className="button secondary" href="/login">
            View dashboard
          </Link>
        </div>
      </div>
      <div className="card">
        <h2 className="section-title">Today at a glance</h2>
        <div className="grid two">
          <div>
            <div className="kpi">$18,450</div>
            <div className="text-muted">Outstanding receivables</div>
          </div>
          <div>
            <div className="kpi">12</div>
            <div className="text-muted">Reminders scheduled</div>
          </div>
          <div>
            <div className="kpi">4</div>
            <div className="text-muted">Overdue invoices</div>
          </div>
          <div>
            <div className="kpi">2.4d</div>
            <div className="text-muted">Average response time</div>
          </div>
        </div>
      </div>
    </section>
  );
}
