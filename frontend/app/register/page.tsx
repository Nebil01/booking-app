'use client';
import { useState } from 'react';
import { api } from '@/lib/api';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function RegisterPage() {
  const [form, setForm] = useState({ first_name: '', last_name: '', email: '', password: '' });
  const [error, setError] = useState('');
  const router = useRouter();

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      const res = await api.post('/register', form);
      if (res.status === 200 || res.status === 201) {
        router.push('/login');
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Registration failed.');
    }
  };

  return (
    <div className="max-w-md mx-auto py-12">
      <div className="bg-white rounded-3xl border border-slate-200 shadow-xl p-8 space-y-6">
        <div className="text-center space-y-2">
          <h2 className="text-2xl font-bold text-slate-900">Create an Account</h2>
          <p className="text-sm text-slate-900">Join Hermes to book exclusive events</p>
        </div>

        {error && (
          <div className="p-3 bg-red-50 border border-red-200 text-red-700 text-sm rounded-xl font-medium">
            {error}
          </div>
        )}

        <form onSubmit={handleRegister} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-semibold text-slate-900 uppercase mb-1">First Name</label>
              <input type="text" value={form.first_name} onChange={(e) => setForm({...form, first_name: e.target.value})} required className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none text-sm" />
            </div>
            <div>
              <label className="block text-xs font-semibold text-slate-900 uppercase mb-1">Last Name</label>
              <input type="text" value={form.last_name} onChange={(e) => setForm({...form, last_name: e.target.value})} required className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none text-sm" />
            </div>
          </div>
          <div>
            <label className="block text-xs font-semibold text-slate-900 uppercase mb-1">Email Address</label>
            <input type="email" value={form.email} onChange={(e) => setForm({...form, email: e.target.value})} required className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none text-sm" />
          </div>
          <div>
            <label className="block text-xs font-semibold text-slate-900 uppercase mb-1">Password</label>
            <input type="password" value={form.password} onChange={(e) => setForm({...form, password: e.target.value})} required className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none text-sm" />
          </div>
          <button type="submit" className="w-full py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-xl shadow-lg shadow-indigo-500/20 transition text-sm">
            Create Account
          </button>
        </form>

        <div className="text-center text-sm text-slate-500">
          <p>Already have an account? <Link href="/login" className="text-indigo-600 font-semibold hover:underline">Sign in</Link></p>
        </div>
      </div>
    </div>
  );
}