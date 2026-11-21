'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { Ticket as TicketIcon } from 'lucide-react';

export default function DashboardPage() {
  const [tickets, setTickets] = useState<any[]>([]);

  useEffect(() => {
    api.get('/my-tickets')
      .then(res => setTickets(res.data.tickets || []))
      .catch(err => console.error('Failed to load tickets', err));
  }, []);

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-900">My Booked Tickets</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {tickets.map((ticket) => (
          <div key={ticket.ID} className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 space-y-4">
            <div className="flex justify-between items-start">
              <span className="text-xs font-bold uppercase tracking-wider text-indigo-600 bg-indigo-50 px-2.5 py-1 rounded-lg">Booking ID: #{ticket.ID}</span>
              <span className="text-sm font-semibold text-slate-700">{ticket.Ticket} Tickets</span>
            </div>
            <h3 className="text-lg font-bold text-slate-900">{ticket.Event ? ticket.Event.title : 'Event Booking'}</h3>
            <p className="text-xs text-slate-400">QR confirmation code has been dispatched via SendGrid.</p>
          </div>
        ))}
        {tickets.length === 0 && (
          <p className="text-slate-500">No tickets found. Explore events to make your first booking.</p>
        )}
      </div>
    </div>
  );
}