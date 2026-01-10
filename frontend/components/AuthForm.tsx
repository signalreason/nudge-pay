'use client';

import { useRouter } from 'next/navigation';
import { FormEvent, useState } from 'react';
import { login, register } from '@/lib/api';
import { setToken } from '@/lib/auth';

interface AuthFormProps {
  mode: 'login' | 'signup';
}

export function AuthForm({ mode }: AuthFormProps) {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [orgName, setOrgName] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const payload =
        mode === 'signup'
          ? await register({ email, password, org_name: orgName })
          : await login({ email, password });
      setToken(payload.token);
      router.push('/dashboard');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong');
    } finally {
      setLoading(false);
    }
  }

  return (
    <form className="form" onSubmit={handleSubmit}>
      {mode === 'signup' ? (
        <label>
          Studio or agency name
          <input
            className="input"
            value={orgName}
            onChange={(event) => setOrgName(event.target.value)}
            placeholder="Signal Studio"
            required
          />
        </label>
      ) : null}
      <label>
        Email
        <input
          className="input"
          type="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          placeholder="you@studio.com"
          required
        />
      </label>
      <label>
        Password
        <input
          className="input"
          type="password"
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          placeholder="••••••••"
          required
        />
      </label>
      {error ? <div className="text-muted">{error}</div> : null}
      <button className="button" type="submit" disabled={loading}>
        {loading ? 'Working...' : mode === 'signup' ? 'Create account' : 'Log in'}
      </button>
    </form>
  );
}
