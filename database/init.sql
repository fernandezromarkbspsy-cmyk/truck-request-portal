-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Users Table
CREATE TYPE user_role AS ENUM ('ops_pic', 'fte_ops', 'fte_mm', 'dock_officer', 'doc_officer');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    clerk_id TEXT UNIQUE NOT NULL, -- Links to Clerk Auth
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    ops_id TEXT UNIQUE,
    role user_role NOT NULL,
    is_fte BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for fast role/ops_id lookups
CREATE INDEX idx_users_ops_id ON users(ops_id);
CREATE INDEX idx_users_clerk_id ON users(clerk_id);

-- 2. Clusters Table (Lookup Reference)
CREATE TABLE clusters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cluster_name TEXT UNIQUE NOT NULL,
    hub_name TEXT NOT NULL,
    region TEXT NOT NULL,
    dock_number TEXT NOT NULL,
    backlogs TEXT,
    backlogs_ts TIMESTAMP WITH TIME ZONE
);

-- Index for fast cluster name lookups
CREATE INDEX idx_clusters_name ON clusters(cluster_name);

-- 3. Requests Table (Core Transaction)
CREATE TYPE request_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED_BY_MM', 'ASSIGNED', 'DOCKED', 'CANCELLED');
CREATE TYPE truck_size_enum AS ENUM ('4W', '6W', '10W', '6WF');
CREATE TYPE truck_type_enum AS ENUM ('WETLEASE', 'DRYLEASE');

CREATE TABLE requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Auto-mapped fields (from clusters table)
    cluster TEXT REFERENCES clusters(cluster_name) ON DELETE RESTRICT,
    region TEXT,
    dock_no TEXT,
    backlogs TEXT,
    backlogs_timestamp TIMESTAMP WITH TIME ZONE,
    
    -- Personnel
    ob_fte TEXT,
    ob_ops_pic TEXT,
    midmile_fte TEXT,
    
    -- Truck Details
    truck_size truck_size_enum,
    truck_type truck_type_enum,
    plate_number TEXT,
    provide_time TIMESTAMP WITH TIME ZONE,
    
    -- Docking Details
    linehaul_trip_no TEXT,
    docked_time TIMESTAMP WITH TIME ZONE,
    driver_id TEXT,
    
    -- Workflow & Status
    status request_status DEFAULT 'PENDING',
    rejection_remarks TEXT,
    approved_at TIMESTAMP WITH TIME ZONE,
    rejected_at TIMESTAMP WITH TIME ZONE,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for fast dashboard filtering (Prevents full table scans)
CREATE INDEX idx_requests_status ON requests(status);
CREATE INDEX idx_requests_cluster ON requests(cluster);
CREATE INDEX idx_requests_ob_ops_pic ON requests(ob_ops_pic);

-- 4. Audit Logs Table (For security and tracking)
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID REFERENCES requests(id) ON DELETE CASCADE,
    action TEXT NOT NULL,
    performed_by UUID REFERENCES users(id),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_request_id ON audit_logs(request_id);