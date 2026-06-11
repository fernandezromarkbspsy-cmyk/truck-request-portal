import { SignIn } from '@clerk/clerk-react';
import { useUser } from '@clerk/clerk-react';
import { useNavigate } from 'react-router-dom';
import { useEffect } from 'react';

export default function LoginPage() {
  const { isSignedIn, user } = useUser();
  const navigate = useNavigate();

  useEffect(() => {
    if (isSignedIn && user) {
      // After successful login, redirect based on role (fetched from backend or Clerk metadata)
      // For now, we redirect to a generic dashboard, the backend middleware will enforce actual access
      navigate('/dashboard');
    }
  }, [isSignedIn, user, navigate]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <div className="w-full max-w-md p-8 bg-white rounded-lg shadow-md">
        <h1 className="text-2xl font-bold text-center text-gray-800 mb-6">
          Truck Request Portal
        </h1>
        {/* Clerk's pre-built component handles Email or Username (ops_id) login securely */}
        <SignIn 
          routing="path" 
          path="/login" 
          appearance={{
            elements: {
              formButtonPrimary: "bg-blue-600 hover:bg-blue-700 text-sm normal-case",
              card: "shadow-none",
            }
          }}
        />
      </div>
    </div>
    </div>
  );
}