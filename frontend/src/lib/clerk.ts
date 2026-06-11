import { ClerkProvider } from '@clerk/clerk-react';

// VITE_CLERK_PUBLISHABLE_KEY must be in your .env file
const clerkPubKey = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;

if (!clerkPubKey) {
  throw new Error("Missing Clerk Publishable Key");
}

export { ClerkProvider, clerkPubKey };