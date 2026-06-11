package clusters

import (
	"encoding/json"
	"net/http"
)

func HandleGetClusters(w http.ResponseWriter, r *http.Request) {
	clusters, err := GetClusters(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch clusters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clusters)
}
