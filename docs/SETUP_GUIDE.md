# 🚛 Truck Request Portal - Zero-Admin Setup Guide

Welcome! This guide is designed specifically for environments where you **do not have administrator rights** on your computer. 

We will use web-based tools (which require no installation) and "portable" software that installs only to your personal user folder.

---

## 📋 Table of Contents
1. [The Recommended Path: GitHub Codespaces (Zero Install)](#1-the-recommended-path-github-codespaces-zero-install)
2. [The Local Path: Portable Tools (No Admin Rights)](#2-the-local-path-portable-tools-no-admin-rights)
3. [Step-by-Step Service Setup (Free Tier)](#3-step-by-step-service-setup-free-tier)
4. [Project Structure Rules](#4-project-structure-rules)
5. [Environment Variables Template](#5-environment-variables-template)
6. [Pre-Flight Checklist](#6-pre-flight-checklist)

---

## 1. The Recommended Path: GitHub Codespaces (Zero Install)
**Highly Recommended.** This gives you a free, pre-configured Linux computer right inside your web browser. Node.js, Go, Git, and Docker are *already installed*.

**How to start:**
1. Create a free GitHub account: [github.com/signup](https://github.com/signup)
2. Create a new repository named `truck-request-portal`.
3. Click the green **"Code"** button.
4. Select the **"Codespaces"** tab.
5. Click **"Create codespace on main"**.
6. A new browser tab will open with a full VS Code editor. Open the terminal at the bottom and you are ready to code!
*Note: GitHub provides 60 hours of free Codespaces usage per month.*

---

## 2. The Local Path: Portable Tools (No Admin Rights)
If you prefer to work on your own computer, use these specific "User-Level" or "Portable" versions.

- **VS Code:** Download the **"User Installer"** from [code.visualstudio.com](https://code.visualstudio.com/).
- **Git:** Download **PortableGit** (`.7z.exe` or `.zip`) from [Git for Windows Releases](https://github.com/git-for-windows/git/releases). Extract to `C:\Users\YourName\PortableGit` and add the `bin` folder to your User "Path" environment variable.
- **Node.js:** Download the **Windows Binary (.zip)** from [nodejs.org](https://nodejs.org/en/download/). Extract to `C:\Users\YourName\nodejs` and add to User "Path".
- **Golang:** Download the **Windows ZIP archive** (NOT the `.msi`) from [go.dev/dl](https://go.dev/dl/). Extract to `C:\Users\YourName\go` and add `bin` to User "Path".
- **Docker:** ⚠️ **SKIP locally.** Docker requires admin rights. We will run the backend directly via `go run cmd/server/main.go` locally, and only use Docker when deploying to the cloud.

---

## 3. Step-by-Step Service Setup (Free Tier)
*All of these are 100% cloud-based. You only need a web browser.*

### A. Supabase (Database)
1. Go to [supabase.com](https://supabase.com) and sign up.
2. Click **"New Project"**. Name it `truck-request-portal`, create a strong database password, and choose a region close to you.
3. Wait 2-3 minutes for it to provision.
4. Go to **Project Settings (gear icon) > API**.
5. Copy the **Project URL** and the **`anon` public key**.
6. Go to the **SQL Editor** (sidebar), click "New Query", paste our `init.sql` script, and click **"Run"**.

### B. Upstash (Redis Caching & Rate Limiting)
1. Go to [upstash.com](https://upstash.com) and sign up.
2. Click **"Create Database"**.
3. Name it `truck-cache`, choose the same region as Supabase, and ensure **"TLS (SSL)"** is enabled.
4. Scroll down to **"REST API"**.
5. Copy the **URL** and the **Token**.

### C. Clerk (Authentication)
1. Go to [clerk.com](https://clerk.com) and sign up.
2. Click **"Create Application"**. Name it `Truck Portal`.
3. Choose **"Email & Username"** for sign-in options (Username will be used for Backroom Ops ID).
4. Go to **API Keys** in the left sidebar.
5. Copy the **Publishable Key** (for Frontend) and **Secret Key** (for Backend).
6. *(Later Step)*: Go to **Webhooks**, add an endpoint `https://your-domain.com/api/v1/users/webhook`, and select `user.created` and `user.updated` events.

### D. Resend (Email Service)
1. Go to [resend.com](https://resend.com) and sign up.
2. Go to **API Keys** in the left sidebar.
3. Click **"Create API Key"**, name it `Truck Portal`, and set permissions to "Sending".
4. Copy the generated **API Key**. *(Note: On the free tier, you can only send emails to your own verified email address until you verify a custom domain).*

### E. Vercel (Frontend Deployment)
1. Go to [vercel.com](https://vercel.com) and sign up **using your GitHub account**.
2. Click **"Add New Project"**.
3. Select your `truck-request-portal` GitHub repository and click **"Import"**.
4. In the "Environment Variables" section, paste all your `VITE_...` variables from your `.env` file.
5. Click **"Deploy"**. Vercel will give you a live `.vercel.app` URL.

### F. Cloudflare (DNS & CDN)
1. Go to [cloudflare.com](https://cloudflare.com) and sign up.
2. Click **"Add a Site"** and enter the domain you bought from Namecheap.
3. Cloudflare will scan your DNS records. Click **"Continue"**.
4. Cloudflare will give you two **Nameservers** (e.g., `bob.ns.cloudflare.com`).
5. Go to your **Namecheap** dashboard > Domain List > Manage > Nameservers > Select "Custom DNS" and paste the two Cloudflare nameservers.
6. Back in Cloudflare, go to **DNS > Records** and add the `CNAME` and `A` records that Vercel provided you in Step E.

### G. PostHog (Analytics)
1. Go to [posthog.com](https://posthog.com) and sign up.
2. Create a new project.
3. Go to **Project Settings**.
4. Copy the **Project API Key**.

### H. Sentry (Error Tracking)
1. Go to [sentry.io](https://sentry.io) and sign up.
2. Create a new project. Select **React** for the frontend, and later create another project for **Go**.
3. Copy the **DSN** (Data Source Name) URL provided.

### I. Better Stack (Uptime Monitoring)
1. Go to [betterstack.com](https://betterstack.com) and sign up.
2. Create a new **Monitor**.
3. Enter your Vercel deployment URL (e.g., `https://truck-portal.vercel.app`).
4. Set the check frequency to 3 minutes and add your email for alerts.

### J. Pinecone (Vector DB - For Advanced Search in Phase 3)
1. Go to [pinecone.io](https://pinecone.io) and sign up.
2. Create a new **Index**.
3. Name it `truck-requests`, set dimensions to `1536` (standard for OpenAI embeddings), and select the **Starter (Free)** tier.
4. Copy the **API Key** and the **Index Host URL**.

---

## 4. Project Structure Rules
To prevent "vibe coding" and ensure maintainability, we strictly follow the **3-Layer Setup** grouped by **Feature**, not file type.

**Do NOT do this:** `/controllers`, `/services`, `/repositories`
**DO this:** Group everything related to a feature together.

```text
truck-request-portal/
├── backend/
│   ├── cmd/server/main.go          # App entry point
│   ├── features/
│   │   ├── requests/               # Everything about requests
│   │   │   ├── requests.controller.go  # Handles HTTP only
│   │   │   ├── requests.service.go     # Business logic only
│   │   │   ├── requests.repository.go  # Database queries only
│   │   │   └── requests.model.go       # Data structures
│   │   ├── users/
│   │   └── clusters/
│   └── pkg/                        # Shared utilities (DB, Cache, Middleware)
├── frontend/
│   ├── src/
│   │   ├── features/               # Feature-specific UI (Auth, Requests, Clusters)
│   │   ├── store/                  # Zustand client state
│   │   └── App.tsx
└── database/
    └── init.sql                    # Supabase schema