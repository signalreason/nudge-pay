export const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface Client {
  id: string;
  name: string;
  email: string;
  company: string;
  phone: string;
  notes: string;
}

export interface Invoice {
  id: string;
  number: string;
  amount_cents: number;
  currency: string;
  due_date: string;
  status: string;
}

export interface Reminder {
  id: string;
  invoice_number: string;
  scheduled_for: string;
  sent_at?: string | null;
  status: string;
}

export interface OutboxEmail {
  id: string;
  to_email: string;
  subject: string;
  created_at: string;
}

export interface Metrics {
  clients: number;
  invoices: number;
  overdue: number;
  upcoming_reminders: number;
  outstanding_cents: number;
}

interface ClientsResponse {
  clients: Client[];
}

interface InvoicesResponse {
  invoices: Invoice[];
}

interface RemindersResponse {
  reminders: Reminder[];
}

interface OutboxResponse {
  outbox: OutboxEmail[];
}

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

export function getMetrics(token: string): Promise<Metrics> {
  return request<Metrics>('/api/metrics', {
    headers: { Authorization: `Bearer ${token}` }
  });
}

export function listClients(token: string): Promise<ClientsResponse> {
  return request<ClientsResponse>('/api/clients', {
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

export function listInvoices(token: string): Promise<InvoicesResponse> {
  return request<InvoicesResponse>('/api/invoices', {
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

export function listReminders(token: string): Promise<RemindersResponse> {
  return request<RemindersResponse>('/api/reminders', {
    headers: { Authorization: `Bearer ${token}` }
  });
}

export function listOutbox(token: string): Promise<OutboxResponse> {
  return request<OutboxResponse>('/api/outbox', {
    headers: { Authorization: `Bearer ${token}` }
  });
}
