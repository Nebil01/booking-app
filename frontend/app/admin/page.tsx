'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { Plus, Trash2, Calendar, Users } from 'lucide-react';

export default function AdminPage() {
  const [events, setEvents] = useState<any[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [form, setForm] = useState({ title: '', date: '', total_capacity: 100 });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    fetchEvents();
  }, []);

  const fetchEvents = async () => {
    try {
      const res = await api.get('/events');
      setEvents(res.data.events || res.data);
    } catch (err) {
      console.error('Failed to load events');
    }
  };

  const handleCreateEvent = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    try {
      await api.post('/admin/events', {
        title: form.title,
        date: form.date,
        total_capacity: Number(form.total_capacity)
      });
      setSuccess('Event created successfully!');
      setShowModal(false);
      setForm({ title: '', date: '', total_capacity: 100 });
      fetchEvents();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Access denied or invalid payload.');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this event?')) return;
    try {
      await api.delete(`/admin/events/${id}`);
      fetchEvents();
    } catch (err: any) {
      alert(err.response?.data?.error || 'Failed to delete event.');
    }
  };

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-slate-900">Admin Control Center</h2>
          <p className="text-sm text-slate-500">Manage platform events and venue capacities securely.</p>
        </div>
        <button onClick={() => setShowModal(true)} className="px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-xl text-sm shadow-md transition flex items-center gap-2">
          <Plus className="w-4 h-4" /> Create New Event
        </button>
      </div>

      {error && <div className="p-4 bg-red-50 border border-red-200 text-red-700 rounded-xl text-sm font-medium">{error}</div>}
      {success && <div className="p-4 bg-emerald-50 border border-emerald-200 text-emerald-700 rounded-xl text-sm font-medium">{success}</div>}

      <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-500 uppercase tracking-wider">
              <th className="p-4">ID</th>
              <th className="p-4">Title</th>
              <th className="p-4">Date</th>
              <th className="p-4">Capacity</th>
              <th className="p-4">Remaining</th>
              <th className="p-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 text-sm">
            {events.map((event) => (
              <tr key={event.ID}>
                <td className="p-4 font-mono text-slate-500">{event.ID}</td>
                <td className="p-4 font-semibold text-slate-900">{event.event_title || event.title}</td>
                <td className="p-4 text-slate-600">{event.date || 'TBA'}</td>
                <td className="p-4 text-slate-600">{event.total_capacity}</td>
                <td className="p-4 font-bold text-indigo-600">{event.remaining_tickets}</td>
                <td className="p-4 text-right">
                  <button onClick={() => handleDelete(event.ID)} className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition" title="Delete Event">
                    <Trash2 className="w-4 h-4" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Create Event Modal */}
      {showModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4">
          <div className="bg-white rounded-3xl max-w-md w-full p-8 shadow-2xl space-y-6">
            <h3 className="text-xl font-bold text-slate-900">Create New Event</h3>
            <form onSubmit={handleCreateEvent} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">Event Title</label>
                <input type="text" value={form.title} onChange={(e) => setForm({...form, title: e.target.value})} required className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none text-sm" />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">Date</label>
                <input type="text" placeholder="YYYY-MM-DD" value={form.date} onChange={(e) => setForm({...form, date: e.target.value})} required className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none text-sm" />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">Total Capacity</label>
                <input type="number" min="1" value={form.total_capacity} onChange={(e) => setForm({...form, total_capacity: Number(e.target.value)})} required className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none text-sm" />
              </div>
              <div className="flex gap-3 pt-2">
                <button type="button" onClick={() => setShowModal(false)} className="w-1/2 py-3 bg-slate-100 font-semibold text-slate-600 rounded-xl">Cancel</button>
                <button type="submit" className="w-1/2 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-xl shadow-lg">Publish</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}