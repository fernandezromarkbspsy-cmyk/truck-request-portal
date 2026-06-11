import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

interface ApprovedRequest {
  id: string;
  cluster: string;
  region: string;
  dock_no: string;
  truck_size: string;
  truck_type: string;
  ob_ops_pic: string;
  request_timestamp: string;
}

export default function FteMMDashboard() {
  const [page, setPage] = useState(1);
  const limit = 10;
  const queryClient = useQueryClient();

  // Modal State
  const [selectedRequest, setSelectedRequest] = useState<ApprovedRequest | null>(null);
  const [modalType, setModalType] = useState<'assign' | 'reject' | null>(null);
  
  // Form State
  const [plateNumber, setPlateNumber] = useState('');
  const [provideTime, setProvideTime] = useState('');
  const [rejectionRemarks, setRejectionRemarks] = useState('');

  const { data, isLoading, error } = useQuery({
    queryKey: ['approvedRequests', page, limit],
    queryFn: async () => {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/requests/approved?page=${page}&limit=${limit}`);
      if (!res.ok) throw new Error('Failed to fetch approved requests');
      return res.json();
    },
  });

  const assignMutation = useMutation({
    mutationFn: async (payload: { id: string; plate_number: string; provide_time: string }) => {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/requests/${payload.id}/assign`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ plate_number: payload.plate_number, provide_time: payload.provide_time }),
      });
      if (!res.ok) throw new Error('Assignment failed');
      return res.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['approvedRequests'] });
      closeModal();
    },
  });

  const rejectMutation = useMutation({
    mutationFn: async (payload: { id: string; rejection_remarks: string }) => {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/requests/${payload.id}/reject`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ rejection_remarks: payload.rejection_remarks }),
      });
      if (!res.ok) throw new Error('Rejection failed');
      return res.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['approvedRequests'] });
      closeModal();
    },
  });

  const closeModal = () => {
    setSelectedRequest(null);
    setModalType(null);
    setPlateNumber('');
    setProvideTime('');
    setRejectionRemarks('');
  };

  if (isLoading) return <div className="p-6 text-blue-600">Loading approved requests...</div>;
  if (error) return <div className="p-6 text-red-600">Error loading requests. Check connection.</div>;

  return (
    <div className="max-w-6xl mx-auto p-6 bg-gray-50 min-h-screen">
      <h1 className="text-2xl font-bold text-gray-800 mb-6">FTE MM Dashboard - Assign Trucks</h1>
      
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead className="bg-gray-100 border-b">
            <tr>
              <th className="p-4 font-semibold text-gray-700">Cluster</th>
              <th className="p-4 font-semibold text-gray-700">Region / Dock</th>
              <th className="p-4 font-semibold text-gray-700">Truck Details</th>
              <th className="p-4 font-semibold text-gray-700">Approved By</th>
              <th className="p-4 font-semibold text-gray-700">Time</th>
              <th className="p-4 font-semibold text-gray-700 text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            {data?.data?.length === 0 ? (
              <tr><td colSpan={6} className="p-6 text-center text-gray-500">No approved requests awaiting assignment.</td></tr>
            ) : (
              data?.data?.map((req: ApprovedRequest) => (
                <tr key={req.id} className="border-b hover:bg-gray-50 transition">
                  <td className="p-4 font-medium">{req.cluster}</td>
                  <td className="p-4">{req.region} / {req.dock_no}</td>
                  <td className="p-4 font-mono text-sm">{req.truck_size} • {req.truck_type}</td>
                  <td className="p-4">{req.ob_ops_pic}</td>
                  <td className="p-4 text-sm text-gray-600">{new Date(req.request_timestamp).toLocaleString()}</td>
                  <td className="p-4 text-right space-x-2">
                    <button
                      onClick={() => { setSelectedRequest(req); setModalType('assign'); }}
                      className="bg-blue-600 hover:bg-blue-700 text-white px-3 py-2 rounded-md text-sm font-medium transition"
                    >
                      Assign Truck
                    </button>
                    <button
                      onClick={() => { setSelectedRequest(req); setModalType('reject'); }}
                      className="bg-red-600 hover:bg-red-700 text-white px-3 py-2 rounded-md text-sm font-medium transition"
                    >
                      Reject
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
        <button disabled={page === 1} onClick={() => setPage(p => p - 1)} className="px-3 py-2 bg-white border rounded-md hover:bg-gray-100 disabled:opacity-50">Previous</button>
        <span>Page {page} of {Math.ceil((data?.total_count || 0) / limit)}</span>
        <button disabled={page * limit >= (data?.total_count || 0)} onClick={() => setPage(p => p + 1)} className="px-3 py-2 bg-white border rounded-md hover:bg-gray-100 disabled:opacity-50">Next</button>
      </div>

      {/* Modals */}
      {modalType && selectedRequest && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg shadow-xl max-w-md w-full p-6">
            <h2 className="text-xl font-bold mb-4">
              {modalType === 'assign' ? 'Assign Truck' : 'Reject Request'}
            </h2>
            <p className="text-sm text-gray-600 mb-4">Request ID: {selectedRequest.id} ({selectedRequest.cluster})</p>
            
            {modalType === 'assign' ? (
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Plate Number *</label>
                  <input type="text" value={plateNumber} onChange={(e) => setPlateNumber(e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md p-2" required />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Provide Time *</label>
                  <input type="datetime-local" value={provideTime} onChange={(e) => setProvideTime(e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md p-2" required />
                </div>
                <div className="flex justify-end space-x-3 mt-6">
                  <button onClick={closeModal} className="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md">Cancel</button>
                  <button 
                    onClick={() => assignMutation.mutate({ id: selectedRequest.id, plate_number: plateNumber, provide_time: provideTime })}
                    disabled={assignMutation.isPending || !plateNumber || !provideTime}
                    className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md disabled:opacity-50"
                  >
                    {assignMutation.isPending ? 'Assigning...' : 'Confirm Assignment'}
                  </button>
                </div>
              </div>
            ) : (
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Rejection Remarks *</label>
                  <textarea value={rejectionRemarks} onChange={(e) => setRejectionRemarks(e.target.value)} rows={4} className="mt-1 block w-full border border-gray-300 rounded-md p-2" placeholder="e.g., No 10W trucks available at this time." required />
                </div>
                <div className="flex justify-end space-x-3 mt-6">
                  <button onClick={closeModal} className="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md">Cancel</button>
                  <button 
                    onClick={() => rejectMutation.mutate({ id: selectedRequest.id, rejection_remarks: rejectionRemarks })}
                    disabled={rejectMutation.isPending || !rejectionRemarks.trim()}
                    className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md disabled:opacity-50"
                  >
                    {rejectMutation.isPending ? 'Rejecting...' : 'Confirm Rejection'}
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}