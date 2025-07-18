// Copyright 2025 William Marchesi

// Author: William Marchesi
// Email: will@marchesi.io
// Website: https://marchesi.io/

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wmarchesi123/go-3dprint-client/config"
	"github.com/wmarchesi123/go-3dprint-client/octoprint"
	"github.com/wmarchesi123/go-3dprint-client/spoolman"
)

type Handler struct {
	config           *config.Config
	mux              *http.ServeMux
	octoprintClients map[string]*octoprint.Client
	spoolmanClient   *spoolman.Client
}

func NewHandler() *Handler {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	h := &Handler{
		config:           cfg,
		mux:              http.NewServeMux(),
		octoprintClients: make(map[string]*octoprint.Client),
		spoolmanClient:   spoolman.NewClient(cfg.SpoolmanURL),
	}

	// Initialize OctoPrint clients for each printer
	for _, printer := range cfg.Printers {
		h.octoprintClients[printer.ID] = octoprint.NewClient(printer.OctoPrintURL, printer.APIKey)
	}

	h.setupRoutes()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	h.mux.ServeHTTP(w, r)
}

func (h *Handler) setupRoutes() {
	h.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	h.mux.HandleFunc("/", h.handleHome)
	h.mux.HandleFunc("/select/", h.handleSpoolSelect)
	h.mux.HandleFunc("/api/printers", h.handleGetPrinters)
	h.mux.HandleFunc("/api/spool/", h.handleGetSpool)
	h.mux.HandleFunc("/api/assign", h.handleAssignSpool)
}

func (h *Handler) respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) respondError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
