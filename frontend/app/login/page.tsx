import { AuthForm } from '@/components/AuthForm';

export default function LoginPage() {
  return (
    <div className="grid two">
      <div className="card">
        <h1 className="section-title">Welcome back</h1>
        <p className="text-muted">Pick up where your collections left off.</p>
      </div>
      <div className="card">
        <AuthForm mode="login" />
      </div>
    </div>
  );
}
