'use client';

import { FormEvent, useEffect, useState } from 'react';
import { Client, createClient, listClients } from '@/lib/api';
import { getToken } from '@/lib/auth';

export function ClientsPanel() {
  const [clients, setClients] = useState<Client[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [company, setCompany] = useState('');

  const token = typeof window !== 'undefined' ? getToken() : null;

  useEffect(() => {
    if (!token) {
      setError('Log in to manage clients.');
      return;
    }
    listClients(token)
      .then((data) => setClients(data.clients))
      .catch((err) => setError(err instanceof Error ? err.message : 'Failed to load'));
  }, [token]);

  async function handleCreate(event: FormEvent) {
    event.preventDefault();
    if (!token) {
      return;
    }
    setError(null);
    try {
      await createClient(token, {
        name,
        email,
        company: company || '-'
      });
      setName('');
      setEmail('');
      setCompany('');
      const data = await listClients(token);
      setClients(data.clients);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create');
    }
  }

  return (
    <div className="grid" style={{ gap: '24px' }}>
      <div className="card">
        <h2 className="section-title">Add client</h2>
        <form className="form" onSubmit={handleCreate}>
          <input
            className="input"
            placeholder="Client name"
            value={name}
            onChange={(event) => setName(event.target.value)}
            required
          />
          <input
            className="input"
            placeholder="Client email"
            type="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            required
          />
          <input
            className="input"
            placeholder="Company"
            value={company}
            onChange={(event) => setCompany(event.target.value)}
          />
          <button className="button" type="submit">
            Save client
          </button>
        </form>
        {error ? <div className="text-muted">{error}</div> : null}
      </div>
      <div className="card">
        <h2 className="section-title">Clients</h2>
        <table className="table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Company</th>
            </tr>
          </thead>
          <tbody>
            {clients.map((client) => (
              <tr key={client.id}>
                <td>{client.name}</td>
                <td>{client.email}</td>
                <td>{client.company}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
