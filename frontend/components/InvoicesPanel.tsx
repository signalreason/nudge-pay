'use client';

import { FormEvent, useEffect, useState } from 'react';
import { createInvoice, listClients, listInvoices } from '@/lib/api';
import { getToken } from '@/lib/auth';

interface Client {
  id: string;
  name: string;
}

interface Invoice {
  id: string;
  number: string;
  amount_cents: number;
  currency: string;
  due_date: string;
  status: string;
}

export function InvoicesPanel() {
  const [clients, setClients] = useState<Client[]>([]);
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [clientID, setClientID] = useState('');
  const [number, setNumber] = useState('');
  const [amount, setAmount] = useState('');
  const [dueDate, setDueDate] = useState('');
  const [error, setError] = useState<string | null>(null);

  const token = typeof window !== 'undefined' ? getToken() : null;

  useEffect(() => {
    if (!token) {
      setError('Log in to manage invoices.');
      return;
    }

    Promise.all([listClients(token), listInvoices(token)])
      .then(([clientsData, invoicesData]) => {
        setClients(clientsData.clients as Client[]);
        setInvoices(invoicesData.invoices as Invoice[]);
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Failed to load'));
  }, [token]);

  async function handleCreate(event: FormEvent) {
    event.preventDefault();
    if (!token) {
      return;
    }
    setError(null);

    const amountCents = Math.round(parseFloat(amount) * 100);
    if (!clientID || !number || !amountCents || !dueDate) {
      setError('Fill out all invoice fields.');
      return;
    }

    try {
      await createInvoice(token, {
        client_id: clientID,
        number,
        amount_cents: amountCents,
        currency: 'USD',
        due_date: dueDate,
        reminder_offsets: [-3, 0, 7]
      });
      setNumber('');
      setAmount('');
      setDueDate('');
      const invoicesData = await listInvoices(token);
      setInvoices(invoicesData.invoices as Invoice[]);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create');
    }
  }

  return (
    <div className="grid" style={{ gap: '24px' }}>
      <div className="card">
        <h2 className="section-title">New invoice</h2>
        <form className="form" onSubmit={handleCreate}>
          <select
            className="input"
            value={clientID}
            onChange={(event) => setClientID(event.target.value)}
            required
          >
            <option value="">Select client</option>
            {clients.map((client) => (
              <option key={client.id} value={client.id}>
                {client.name}
              </option>
            ))}
          </select>
          <input
            className="input"
            placeholder="Invoice number"
            value={number}
            onChange={(event) => setNumber(event.target.value)}
            required
          />
          <input
            className="input"
            placeholder="Amount (USD)"
            value={amount}
            onChange={(event) => setAmount(event.target.value)}
            required
          />
          <input
            className="input"
            type="date"
            value={dueDate}
            onChange={(event) => setDueDate(event.target.value)}
            required
          />
          <button className="button" type="submit">
            Schedule invoice
          </button>
        </form>
        {error ? <div className="text-muted">{error}</div> : null}
      </div>
      <div className="card">
        <h2 className="section-title">Invoices</h2>
        <table className="table">
          <thead>
            <tr>
              <th>Number</th>
              <th>Amount</th>
              <th>Due</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {invoices.map((invoice) => (
              <tr key={invoice.id}>
                <td>{invoice.number}</td>
                <td>
                  {(invoice.amount_cents / 100).toFixed(2)} {invoice.currency}
                </td>
                <td>{invoice.due_date.slice(0, 10)}</td>
                <td>
                  <span className="badge">{invoice.status}</span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
