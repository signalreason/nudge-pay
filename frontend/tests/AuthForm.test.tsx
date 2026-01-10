import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { AuthForm } from '@/components/AuthForm';

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn() })
}));

vi.mock('@/lib/api', () => ({
  login: vi.fn(() => Promise.resolve({ token: 'token' })),
  register: vi.fn(() => Promise.resolve({ token: 'token' }))
}));

vi.mock('@/lib/auth', () => ({
  setToken: vi.fn()
}));

describe('AuthForm', () => {
  it('renders signup fields', () => {
    render(<AuthForm mode="signup" />);
    expect(screen.getByText('Studio or agency name')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('you@studio.com')).toBeInTheDocument();
  });

  it('renders login fields', () => {
    render(<AuthForm mode="login" />);
    expect(screen.queryByText('Studio or agency name')).toBeNull();
    expect(screen.getByPlaceholderText('you@studio.com')).toBeInTheDocument();
  });
});
