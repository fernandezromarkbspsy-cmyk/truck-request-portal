import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

interface RequestData {
  id: string;
  request_timestamp: string;
  cluster: string;
  region: string;
  dock_no: string;
  truck_size: string;
  truck_type: string;
  ob_ops_pic: string;
}

export default function FteOpsDashboard() {
  const [page, setPage] = useState(1);
  const limit = 10;
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    queryKey: ['pendingRequests', page, limit],
    queryFn: async () => {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/requests?page=${page}&limit=${limit}`);
      if (!res.ok) throw new Error('Failed to fetch');
      return res.json();
    },
  });

  const approveMutation = useMutation({
    mutationFn: async (id: string) => {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/requests/${id}/approve`, { method: 'PUT' });
      if (!res.ok) throw new Error('Approval failed');
      return res.json();
    },
    onSuccess: () => {
      // Optimistic UI update: refetch pending list immediately
      queryClient.invalidateQueries({ queryKey: ['pendingRequests'] });
    },
  });

  if (isLoading) return <div className="p-6 text-blue-600">Loading dashboard...</div>;
  if (error) return <div className="p-6 text-red-600">Error loading requests. Check connection.</div>;

  return (
    <div className="max-w-6xl mx-auto p-6 bg-gray-50 min-h-screen">
      <h1 className="text-2xl font-bold text-gray-800 mb-6">FTE Ops Dashboard - Pending Requests</h1>
      
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead className="bg-gray-100 border-b">
            <tr>
              <th className="p-4 font-semibold text-gray-700">Cluster</th>
              <th className="p-4 font-semibold text-gray-700">Region / Dock</th>
              <th className="p-4 font-semibold text-gray-700">Truck Details</th>
              <th className="p-4 font-semibold text-gray-700">Requested By</th>
              <th className="p-4 font-semibold text-gray-700">Time</th>
              <th className="p-4 font-semibold text-gray-700">Action</th>
            </tr>
          </thead>
          <tbody>
            {data?.data?.length === 0 ? (
              <tr><td colSpan={6} className="p-6 text-center text-gray-500">No pending requests. Good job!</td></tr>
            ) : (
              data?.data?.map((req: RequestData) => (
                <tr key={req.id} className="border-b hover:bg-gray-50 transition">
                  <td className="p-4">{req.cluster}</td>
                  <td className="p-4">{req.region} / {req.dock_no}</td>
                  <td className="p-4 font-mono text-sm">{req.truck_size} • {req.truck_type}</td>
                  <td className="p-4">{req.ob_ops_pic}</td>
                  <td className="p-4 text-sm text-gray-600">{new Date(req.request_timestamp).toLocaleTimeString()}</td>
                  <td className="p-4">
                    <button
                      onClick={() => approveMutation.mutate(req.id)}
                      disabled={approveMutation.isPending}
                      className="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed transition"
                    >
                      {approveMutation.isPending ? 'Processing...' : 'Approve'}
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination Controls */}
      <div className="flex justify-between items-center mt-4 text-sm text-gray-600">
        <button 
          disabled={page === 1} 
          onClick={() => setPage(p => p - 1)}
          className="px-3 py-2 bg-white border rounded-md hover:bg-gray-100 disabled:opacity-50"
        >Previous</button>
        <span>Page {page} of {Math.ceil((data?.total_count || 0) / limit)}</span>
        <button 
          disabled={page * limit >= (data?.total_count || 0)} 
          onClick={() => setPage(p => p + 1)}
          className="px-3 py-2 bg-white border rounded-md hover:bg-gray-100 disabled:opacity-50"
        >Next</button>
      </div>
    </div>
  );
}