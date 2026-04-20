package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"drumkit-take-home/internal/load"
	"drumkit-take-home/internal/turvo"
)

type Server struct {
	turvoClient *turvo.Client
}

func NewServer(turvoClient *turvo.Client) *Server {
	return &Server{turvoClient: turvoClient}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/v1/loads", s.handleListLoads)
	mux.HandleFunc("/v1/integrations/webhooks/loads", s.handleCreateLoad)
	return s.withMiddleware(mux)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok": true,
	})
}

func (s *Server) handleListLoads(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params, err := listParamsFromRequest(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	result, err := s.turvoClient.ListLoads(ctx, params)
	if err != nil {
		log.Printf("list loads failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleCreateLoad(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input load.Load
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": fmt.Sprintf("invalid request body: %v", err),
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	result, err := s.turvoClient.CreateLoad(ctx, input)
	if err != nil {
		status := http.StatusInternalServerError
		if isCreateLoadValidationError(err) {
			status = http.StatusBadRequest
		}
		log.Printf("create load failed: %v", err)
		writeJSON(w, status, map[string]any{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (s *Server) withMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		started := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(started))
	})
}

func listParamsFromRequest(r *http.Request) (load.ListParams, error) {
	page, err := positiveIntQueryParam(r, "page", 1)
	if err != nil {
		return load.ListParams{}, err
	}

	limit, err := positiveIntQueryParam(r, "limit", 20)
	if err != nil {
		return load.ListParams{}, err
	}
	if limit > 100 {
		return load.ListParams{}, fmt.Errorf("limit must be less than or equal to 100")
	}

	return load.ListParams{
		Status:               r.URL.Query().Get("status"),
		CustomerID:           r.URL.Query().Get("customerId"),
		PickupDateSearchFrom: r.URL.Query().Get("pickupDateSearchFrom"),
		PickupDateSearchTo:   r.URL.Query().Get("pickupDateSearchTo"),
		Page:                 page,
		Limit:                limit,
	}, nil
}

func positiveIntQueryParam(r *http.Request, key string, fallback int) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}

	return parsed, nil
}

func isCreateLoadValidationError(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "is required") || strings.Contains(message, "must be a positive integer") || strings.Contains(message, "must be a valid rfc3339 timestamp")
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write json failed: %v", err)
	}
}
