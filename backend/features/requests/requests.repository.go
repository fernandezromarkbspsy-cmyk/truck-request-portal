package requests

import (
	"context"
	"time"
	"truck-request-portal/pkg/database"
)

func CreateRequest(ctx context.Context, req *Request) error {
	// Rule: Explicit columns only. No SELECT *.
	query := `
		INSERT INTO requests (
			cluster, region, dock_no, backlogs, backlogs_timestamp,
			ob_ops_pic, truck_size, truck_type, status, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := database.DB.Exec(ctx, query,
		req.Cluster, req.Region, req.DockNo, req.Backlogs, req.BacklogsTimestamp,
		req.ObOpsPic, req.TruckSize, req.TruckType, req.Status, req.UpdatedAt,
	)
	return err
}

func GetPendingRequests(ctx context.Context, page, limit int) (PaginatedResponse, error) {
	var resp PaginatedResponse
	offset := (page - 1) * limit

	// 1. Get total count for pagination UI
	countQuery := `SELECT COUNT(*) FROM requests WHERE status = 'PENDING'`
	err := database.DB.QueryRow(ctx, countQuery).Scan(&resp.TotalCount)
	if err != nil {
		return resp, err
	}

	// 2. Fetch data with explicit LIMIT/OFFSET (Rule: Never load 50k rows)
	dataQuery := `
		SELECT id, request_timestamp, cluster, region, dock_no, backlogs, backlogs_timestamp,
		       ob_ops_pic, truck_size, truck_type, status, updated_at
		FROM requests
		WHERE status = 'PENDING'
		ORDER BY request_timestamp DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := database.DB.Query(ctx, dataQuery, limit, offset)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Request
		err := rows.Scan(&r.ID, &r.RequestTimestamp, &r.Cluster, &r.Region, &r.DockNo, &r.Backlogs, &r.BacklogsTimestamp,
			&r.ObOpsPic, &r.TruckSize, &r.TruckType, &r.Status, &r.UpdatedAt)
		if err != nil {
			return resp, err
		}
		resp.Data = append(resp.Data, r)
	}
	resp.Page = page
	resp.Limit = limit
	return resp, rows.Err()
}

func UpdateStatusToApproved(ctx context.Context, requestID string) error {
	query := `
		UPDATE requests 
		SET status = 'APPROVED', approved_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND status = 'PENDING'
	`
	_, err := database.DB.Exec(ctx, query, requestID)
	return err
}

func GetApprovedRequests(ctx context.Context, page, limit int) (PaginatedResponse, error) {
	var resp PaginatedResponse
	offset := (page - 1) * limit

	// Rule: Explicit columns only. No SELECT *.
	countQuery := `SELECT COUNT(*) FROM requests WHERE status = 'APPROVED'`
	err := database.DB.QueryRow(ctx, countQuery).Scan(&resp.TotalCount)
	if err != nil {
		return resp, err
	}

	dataQuery := `
		SELECT id, cluster, region, dock_no, truck_size, truck_type, ob_ops_pic, request_timestamp
		FROM requests
		WHERE status = 'APPROVED'
		ORDER BY request_timestamp DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := database.DB.Query(ctx, dataQuery, limit, offset)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Request
		err := rows.Scan(&r.ID, &r.Cluster, &r.Region, &r.DockNo, &r.TruckSize, &r.TruckType, &r.ObOpsPic, &r.RequestTimestamp)
		if err != nil {
			return resp, err
		}
		resp.Data = append(resp.Data, r)
	}
	resp.Page = page
	resp.Limit = limit
	return resp, rows.Err()
}

func AssignTruck(ctx context.Context, requestID string, plateNumber string, provideTime time.Time) error {
	// Only update if status is currently APPROVED (prevents race conditions)
	query := `
		UPDATE requests 
		SET status = 'ASSIGNED', plate_number = $1, provide_time = $2, updated_at = NOW()
		WHERE id = $1 AND status = 'APPROVED'
	`
	// Note: Using $1 for ID twice is a PostgreSQL quirk, better to use $1, $2, $3
	query = `
		UPDATE requests 
		SET status = 'ASSIGNED', plate_number = $1, provide_time = $2, updated_at = NOW()
		WHERE id = $3 AND status = 'APPROVED'
	`
	_, err := database.DB.Exec(ctx, query, plateNumber, provideTime, requestID)
	return err
}

func RejectRequest(ctx context.Context, requestID string, remarks string) error {
	query := `
		UPDATE requests 
		SET status = 'REJECTED_BY_MM', rejection_remarks = $1, rejected_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND status = 'APPROVED'
	`
	_, err := database.DB.Exec(ctx, query, remarks, requestID)
	return err
}
