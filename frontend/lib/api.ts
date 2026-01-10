export const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {})
    },
    ...options
  });

  if (!response.ok) {
    const payload = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(payload.error || 'Request failed');
  }

  return response.json() as Promise<T>;
}

export function register(payload: {
  email: string;
  password: string;
  org_name: string;
}): Promise<{ token: string }> {
  return request('/api/auth/register', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export function login(payload: {
  email: string;
  password: string;
}): Promise<{ token: string }> {
  return request('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export function getMetrics(token: string) {
  return request('/api/metrics', {
    headers: { Authorization: `Bearer ${token}` }
  });
}

export function listClients(token: string) {
  return request('/api/clients', {
    headers: { Authorization: `Bearer ${token}` }
  });
}

export function createClient(token: string, payload: {
  name: string;
  email: string;
  company: string;
  phone?: string;
  notes?: string;
}) {
  return request('/api/clients', {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: JSON.stringify(payload)
  });
}

export function listInvoices(token: string) {
  return request('/api/invoices', {
    headers: { Authorization: `Bearer ${token}` }
  });
}

export function createInvoice(token: string, payload: {
  client_id: string;
  number: string;
  amount_cents: number;
  currency: string;
  due_date: string;
  reminder_offsets?: number[];
}) {
  return request('/api/invoices', {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: JSON.stringify(payload)
  });
}

export function listReminders(token: string) {
  return request('/api/reminders', {
    headers: { Authorization: `Bearer ${token}` }
  });
}

export function listOutbox(token: string) {
  return request('/api/outbox', {
    headers: { Authorization: `Bearer ${token}` }
  });
}
