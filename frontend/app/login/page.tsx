'use client';
import { useState } from 'react';
import { api } from '@/lib/api';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const router = useRouter();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      const res = await api.post('/login', { email, password });
      if (res.status === 200) {
        // Extract user data and role from the backend response
        const data = res.data;
        const role = data.role || data.user?.role;

        // Optionally save user session data to localStorage
        localStorage.setItem('user', JSON.stringify(data.user || data));

        // Redirect based on role
        if (role === 'admin') {
          router.push('/admin');
        } else {
          router.push('/dashboard');
        }
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Invalid credentials or network error.');
    }
  };

  return (
    <div className="max-w-md mx-auto py-12">
      <div className="bg-white rounded-3xl border border-slate-200 shadow-xl p-8 space-y-6">
        <div className="text-center space-y-2">
          <h2 className="text-2xl font-bold text-slate-900">Welcome Back</h2>
          <p className="text-sm text-slate-500">Sign in to access your tickets and event dashboard</p>
        </div>

        {error && (
          <div className="p-3 bg-red-50 border border-red-200 text-red-700 text-sm rounded-xl font-medium">
            {error}
          </div>
        )}

        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <label className="block text-xs font-semibold text-slate-900 uppercase mb-1">Email Address</label>
            <input 
              type="email" 
              value={email} 
              onChange={(e) => setEmail(e.target.value)} 
              required 
              className="w-full px-4 py-3 rounded-xl border border-slate-200 text-slate-900 bg-white focus:ring-2 focus:ring-indigo-500 focus:outline-none text-sm" 
            />
          </div>
          <div>
            <label className="block text-xs font-semibold text-slate-900 uppercase mb-1">Password</label>
            <input 
              type="password" 
              value={password} 
              onChange={(e) => setPassword(e.target.value)} 
              required 
              className="w-full px-4 py-3 rounded-xl border border-slate-200 text-slate-900 bg-white focus:ring-2 focus:ring-indigo-500 focus:outline-none text-sm" 
            />
          </div>
          <button type="submit" className="w-full py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-xl shadow-lg shadow-indigo-500/20 transition text-sm">
            Sign In
          </button>
        </form>

        <div className="text-center text-sm text-slate-500">
          <p>Don't have an account? <Link href="/register" className="text-indigo-600 font-semibold hover:underline">Sign up</Link></p>
        </div>
      </div>
    </div>
  );
}