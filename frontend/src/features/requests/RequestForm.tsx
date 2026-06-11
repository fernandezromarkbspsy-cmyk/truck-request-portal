// Replace the old handleSubmit in RequestForm.tsx with this:
const submitMutation = useMutation({
  mutationFn: async (formData: any) => {
    const res = await fetch(`${import.meta.env.VITE_API_URL}/requests`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(formData),
    });
    if (!res.ok) throw new Error('Submission failed');
    return res.json();
  },
  onSuccess: () => {
    alert('Request submitted successfully!');
    // Reset form or redirect
  }
});

const handleSubmit = (e: React.FormEvent) => {
  e.preventDefault();
  submitMutation.mutate({
    cluster: selectedCluster,
    dock_no: dockNo,
    truck_size: truckSize,
    truck_type: truckType
  });
};