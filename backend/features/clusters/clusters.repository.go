package clusters

import (
	"context"

	"truck-request-portal/pkg/database"
)

func GetClusters(ctx context.Context) ([]Cluster, error) {
	query := `
		SELECT id, cluster_name, hub_name, region, dock_number, backlogs, backlogs_ts
		FROM clusters
		ORDER BY cluster_name ASC
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clusters []Cluster
	for rows.Next() {
		var cluster Cluster
		if err := rows.Scan(
			&cluster.ID,
			&cluster.Name,
			&cluster.HubName,
			&cluster.Region,
			&cluster.DockNumber,
			&cluster.Backlogs,
			&cluster.BacklogsTS,
		); err != nil {
			return nil, err
		}
		clusters = append(clusters, cluster)
	}

	return clusters, rows.Err()
}
