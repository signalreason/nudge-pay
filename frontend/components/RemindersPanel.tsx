'use client';

import { useEffect, useState } from 'react';
import { OutboxEmail, Reminder, listOutbox, listReminders } from '@/lib/api';
import { getToken } from '@/lib/auth';

export function RemindersPanel() {
  const [reminders, setReminders] = useState<Reminder[]>([]);
  const [outbox, setOutbox] = useState<OutboxEmail[]>([]);
  const [error, setError] = useState<string | null>(null);

  const token = typeof window !== 'undefined' ? getToken() : null;

  useEffect(() => {
    if (!token) {
      setError('Log in to view reminders.');
      return;
    }

    Promise.all([listReminders(token), listOutbox(token)])
      .then(([remindersData, outboxData]) => {
        setReminders(remindersData.reminders);
        setOutbox(outboxData.outbox);
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Failed to load'));
  }, [token]);

  if (error) {
    return <div className="card">{error}</div>;
  }

  return (
    <div className="grid" style={{ gap: '24px' }}>
      <div className="card">
        <h2 className="section-title">Upcoming reminders</h2>
        <table className="table">
          <thead>
            <tr>
              <th>Invoice</th>
              <th>Scheduled</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {reminders.map((reminder) => (
              <tr key={reminder.id}>
                <td>{reminder.invoice_number}</td>
                <td>{reminder.scheduled_for.slice(0, 10)}</td>
                <td>
                  <span className="badge">{reminder.status}</span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="card">
        <h2 className="section-title">Recent outbox</h2>
        <table className="table">
          <thead>
            <tr>
              <th>Recipient</th>
              <th>Subject</th>
              <th>Sent</th>
            </tr>
          </thead>
          <tbody>
            {outbox.map((item) => (
              <tr key={item.id}>
                <td>{item.to_email}</td>
                <td>{item.subject}</td>
                <td>{item.created_at.slice(0, 10)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
