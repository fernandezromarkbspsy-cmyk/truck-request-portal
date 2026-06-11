package clusters

import "time"

type Cluster struct {
	ID         string     `json:"id"`
	Name       string     `json:"cluster_name"`
	HubName    string     `json:"hub_name"`
	Region     string     `json:"region"`
	DockNumber string     `json:"dock_number"`
	Backlogs   *string    `json:"backlogs"`
	BacklogsTS *time.Time `json:"backlogs_ts"`
}
