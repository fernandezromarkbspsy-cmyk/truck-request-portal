package requests

import (
	"context"
	"fmt"
	"time"
	"truck-request-portal/pkg/database"
)

func CreateRequestService(ctx context.Context, userID string, input *CreateRequestInput) error {
	// 1. Validate & Fetch Cluster Data for Auto-Mapping
	var region, dockDefault, backlogs string
	var backlogsTs *time.Time
	clusterQuery := `SELECT region, dock_number, backlogs, backlogs_ts FROM clusters WHERE cluster_name = $1`
	err := database.DB.QueryRow(ctx, clusterQuery, input.Cluster).Scan(&region, &dockDefault, &backlogs, &backlogsTs)
	if err != nil {
		return fmt.Errorf("invalid cluster: %v", err)
	}

	// 2. Allow editable dock_no override
	finalDock := dockDefault
	if input.DockNo != "" {
		finalDock = input.DockNo
	}

	// 3. Prepare Record
	req := &Request{
		Cluster:           input.Cluster,
		Region:            region,
		DockNo:            finalDock,
		Backlogs:          backlogs,
		BacklogsTimestamp: backlogsTs,
		ObOpsPic:          userID,
		TruckSize:         input.TruckSize,
		TruckType:         input.TruckType,
		Status:            "PENDING",
		UpdatedAt:         time.Now(),
	}

	// 4. Save to DB
	return CreateRequest(ctx, req)
}

func GetPendingRequestsService(ctx context.Context, page, limit int) (PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	} // Enforce max limit
	return GetPendingRequests(ctx, page, limit)
}

func ApproveRequestService(ctx context.Context, requestID string) error {
	return UpdateStatusToApproved(ctx, requestID)
}

func GetApprovedRequestsService(ctx context.Context, page, limit int) (PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	} // Enforce max limit
	return GetApprovedRequests(ctx, page, limit)
}

func AssignTruckService(ctx context.Context, requestID string, input *AssignRequestInput) error {
	if input.PlateNumber == "" {
		return fmt.Errorf("plate number is required")
	}
	return AssignTruck(ctx, requestID, input.PlateNumber, input.ProvideTime)
}

func RejectRequestService(ctx context.Context, requestID string, input *RejectRequestInput) error {
	if input.RejectionRemarks == "" {
		return fmt.Errorf("rejection remarks are required")
	}
	return RejectRequest(ctx, requestID, input.RejectionRemarks)
}
