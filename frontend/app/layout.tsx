import './globals.css';
import Link from 'next/link';
import { Ticket, LogIn, UserPlus } from 'lucide-react';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className="bg-slate-50 text-slate-800 font-sans antialiased min-h-screen flex flex-col justify-between">
        <header className="sticky top-0 z-40 backdrop-blur-md bg-white/80 border-b border-slate-200">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
            <Link href="/" className="flex items-center gap-2 font-bold text-xl text-slate-900">
              <div className="bg-indigo-600 text-white p-2 rounded-xl shadow-md shadow-indigo-500/20">
                <Ticket className="w-5 h-5" />
              </div>
              <span>Hermes<span className="text-indigo-600">.</span></span>
            </Link>
            <nav className="hidden md:flex items-center gap-2">
              <Link href="/" className="px-4 py-2 rounded-lg text-sm font-medium text-slate-600 hover:text-slate-900 transition">Explore Events</Link>
              <Link href="/dashboard" className="px-4 py-2 rounded-lg text-sm font-medium text-slate-600 hover:text-slate-900 transition">My Tickets</Link>
              <Link href="/admin" className="px-4 py-2 rounded-lg text-sm font-medium text-indigo-600 hover:text-indigo-700 transition">Admin Center</Link>
            </nav>
            <div className="flex items-center gap-3">
              <Link href="/login" className="px-4 py-2 text-sm font-medium text-slate-700 hover:text-slate-900 transition">Sign In</Link>
              <Link href="/register" className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-xl shadow-md shadow-indigo-500/20 transition">Get Started</Link>
            </div>
          </div>
        </header>

        <main className="flex-grow max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 w-full">
          {children}
        </main>

        <footer className="border-t border-slate-200 bg-white py-6 text-center text-xs text-slate-400">
          Project Hermes Event Ticketing Platform &copy; 2026. Secured with Go & PostgreSQL.
        </footer>
      </body>
    </html>
  );
}