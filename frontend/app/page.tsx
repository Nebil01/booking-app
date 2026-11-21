'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { Calendar, Users, ArrowRight, Sparkles } from 'lucide-react';

interface Event {
  ID: number;
  title?: string;
  event_title?: string;
  date: string;
  total_capacity: number;
  remaining_tickets: number;
}

export default function EventsPage() {
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedEvent, setSelectedEvent] = useState<Event | null>(null);
  const [ticketCount, setTicketCount] = useState(1);
  const [message, setMessage] = useState('');

  useEffect(() => {
    fetchEvents();
  }, []);

  const fetchEvents = async () => {
    try {
      const res = await api.get('/events');
      setEvents(res.data.events || res.data);
    } catch (err) {
      console.error('Failed to load events', err);
    } finally {
      setLoading(false);
    }
  };

  const handleBook = async () => {
    if (!selectedEvent) return;
    try {
      await api.post('/book', { user_ticket: ticketCount });
      setMessage('Tickets booked successfully! QR code sent to your email.');
      setSelectedEvent(null);
      fetchEvents();
    } catch (err: any) {
      setMessage(err.response?.data?.error || 'Booking failed');
    }
  };

  return (
    <div className="space-y-8">
      {/* Hero Section */}
      <div className="relative overflow-hidden rounded-3xl bg-slate-900 text-white p-8 md:p-12 shadow-2xl">
        <div className="absolute inset-0 opacity-20 bg-[radial-gradient(#6366f1_1px,transparent_1px)] [background-size:16px_16px]"></div>
        <div className="relative z-10 max-w-2xl space-y-4">
          <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-indigo-500/20 text-indigo-300 border border-indigo-500/30">
            <Sparkles className="w-3.5 h-3.5" /> Secure Live Booking Platform
          </span>
          <h1 className="text-4xl sm:text-5xl font-extrabold tracking-tight">Find your next unforgettable event.</h1>
          <p className="text-slate-300 text-base">Book verified experiences instantly with cryptographic QR verification.</p>
        </div>
      </div>

      {message && (
        <div className="p-4 bg-indigo-50 border border-indigo-200 text-indigo-700 rounded-xl text-sm font-medium">
          {message}
        </div>
      )}

      {/* Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {events.map((event) => (
          <div key={event.ID} className="bg-white rounded-2xl border border-slate-200 shadow-sm hover:shadow-xl transition-all duration-300 flex flex-col justify-between overflow-hidden">
            <div className="p-6 space-y-4">
              <div className="flex items-center justify-between">
                <span className="px-3 py-1 bg-slate-100 text-slate-700 font-medium text-xs rounded-lg flex items-center gap-1.5">
                  <Calendar className="w-3.5 h-3.5" /> {event.date || 'TBA'}
                </span>
                <span className={`px-3 py-1 text-xs font-bold rounded-full ${event.remaining_tickets > 10 ? 'bg-emerald-50 text-emerald-700' : 'bg-red-50 text-red-700'}`}>
                  {event.remaining_tickets} left
                </span>
              </div>
              <h3 className="text-xl font-bold text-slate-900">{event.event_title || event.title}</h3>
              <div className="text-sm text-slate-500 flex items-center gap-2">
                <Users className="w-4 h-4" /> Capacity: <span className="font-semibold text-slate-700">{event.total_capacity}</span>
              </div>
            </div>
            <div className="p-6 pt-0">
              <button onClick={() => setSelectedEvent(event)} className="w-full py-3 px-4 bg-slate-900 hover:bg-indigo-600 text-white font-semibold rounded-xl shadow-md transition flex items-center justify-center gap-2">
                <span>Book Tickets</span>
                <ArrowRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Booking Modal */}
      {selectedEvent && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4">
          <div className="bg-white rounded-3xl max-w-md w-full p-8 shadow-2xl space-y-6">
            <h3 className="text-xl font-bold text-slate-900">Book: {selectedEvent.event_title || selectedEvent.title}</h3>
            <div>
              <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">Ticket Quantity</label>
              <input type="number" min="1" value={ticketCount} onChange={(e) => setTicketCount(Number(e.target.value))} className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none text-sm" />
            </div>
            <div className="flex gap-3">
              <button onClick={() => setSelectedEvent(null)} className="w-1/2 py-3 bg-slate-100 font-semibold text-slate-600 rounded-xl">Cancel</button>
              <button onClick={handleBook} className="w-1/2 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-xl shadow-lg">Confirm</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}