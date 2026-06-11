import React from 'react';
import ReactDOM from 'react-dom/client';
import { ClerkProvider } from './lib/clerk';
import App from './App';
import './index.css'; // Tailwind CSS imports

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ClerkProvider publishableKey={import.meta.env.VITE_CLERK_PUBLISHABLE_KEY}>
      <App />
    </ClerkProvider>
  </React.StrictMode>
);