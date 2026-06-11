package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HandleCreateRequest(w http.ResponseWriter, r *http.Request) {
	var input CreateRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value("clerkUserID").(string)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := CreateRequestService(r.Context(), userID, &input)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Request created successfully"})
}

func HandleGetPendingRequests(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	resp, err := GetPendingRequestsService(r.Context(), page, limit)
	if err != nil {
		http.Error(w, "Failed to fetch requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func HandleApproveRequest(w http.ResponseWriter, r *http.Request) {
	requestID := chi.URLParam(r, "id")
	if requestID == "" {
		http.Error(w, "Missing request ID", http.StatusBadRequest)
		return
	}

	err := ApproveRequestService(r.Context(), requestID)
	if err != nil {
		http.Error(w, "Failed to approve request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Request approved"})
}

func HandleGetApprovedRequests(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	resp, err := GetApprovedRequestsService(r.Context(), page, limit)
	if err != nil {
		http.Error(w, "Failed to fetch approved requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func HandleAssignTruck(w http.ResponseWriter, r *http.Request) {
	requestID := chi.URLParam(r, "id")
	if requestID == "" {
		http.Error(w, "Missing request ID", http.StatusBadRequest)
		return
	}

	var input AssignRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err := AssignTruckService(r.Context(), requestID, &input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to assign truck: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Truck assigned successfully"})
}

func HandleRejectRequest(w http.ResponseWriter, r *http.Request) {
	requestID := chi.URLParam(r, "id")
	if requestID == "" {
		http.Error(w, "Missing request ID", http.StatusBadRequest)
		return
	}

	var input RejectRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err := RejectRequestService(r.Context(), requestID, &input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to reject request: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Request rejected and routed back to FTE Ops"})
}
