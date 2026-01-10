import { AuthForm } from '@/components/AuthForm';

export default function SignupPage() {
  return (
    <div className="grid two">
      <div className="card">
        <h1 className="section-title">Start your reminder studio</h1>
        <p className="text-muted">
          Spin up a new workspace and start nudging invoices today.
        </p>
      </div>
      <div className="card">
        <AuthForm mode="signup" />
      </div>
    </div>
  );
}
