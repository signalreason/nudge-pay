import '@/styles/globals.css';
import Link from 'next/link';
import { BackgroundShapes } from '@/components/BackgroundShapes';

export const metadata = {
  title: 'NudgePay',
  description: 'Automated invoice reminders for cash-flow sanity.'
};

export default function RootLayout({
  children
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <div className="app-shell">
          <BackgroundShapes />
          <header className="top-bar">
            <div>
              <span className="logo">NudgePay</span>
              <span className="tag" style={{ marginLeft: '8px' }}>
                v0.1
              </span>
            </div>
            <nav className="nav">
              <Link href="/dashboard">Dashboard</Link>
              <Link href="/clients">Clients</Link>
              <Link href="/invoices">Invoices</Link>
              <Link href="/reminders">Reminders</Link>
            </nav>
            <div style={{ display: 'flex', gap: '12px' }}>
              <Link className="button secondary" href="/login">
                Log in
              </Link>
              <Link className="button" href="/signup">
                Start free
              </Link>
            </div>
          </header>
          <main>{children}</main>
        </div>
      </body>
    </html>
  );
}
