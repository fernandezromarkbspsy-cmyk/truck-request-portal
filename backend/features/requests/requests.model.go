package requests

import "time"

type Request struct {
	ID                string     `json:"id"`
	RequestTimestamp  time.Time  `json:"request_timestamp"`
	Cluster           string     `json:"cluster"`
	Region            string     `json:"region"`
	DockNo            string     `json:"dock_no"`
	Backlogs          string     `json:"backlogs"`
	BacklogsTimestamp *time.Time `json:"backlogs_timestamp"`
	ObOpsPic          string     `json:"ob_ops_pic"`
	TruckSize         string     `json:"truck_size"`
	TruckType         string     `json:"truck_type"`
	Status            string     `json:"status"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type CreateRequestInput struct {
	Cluster   string `json:"cluster"`
	DockNo    string `json:"dock_no"`
	TruckSize string `json:"truck_size"`
	TruckType string `json:"truck_type"`
}

type PaginatedResponse struct {
	Data       []Request `json:"data"`
	TotalCount int64     `json:"total_count"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
}

type AssignRequestInput struct {
	PlateNumber string    `json:"plate_number"`
	ProvideTime time.Time `json:"provide_time"`
}

type RejectRequestInput struct {
	RejectionRemarks string `json:"rejection_remarks"`
}
