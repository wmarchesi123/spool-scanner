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
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// Home page - just redirects to instructions
func (h *Handler) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Spool Scanner</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body style="font-family: sans-serif; padding: 20px; max-width: 600px; margin: 0 auto;">
    <h1>Spool Scanner</h1>
    <p>To use this system:</p>
    <ol>
        <li>Scan the QR code or NFC tag on your filament spool</li>
        <li>Select which printer you're loading it into</li>
        <li>Confirm the assignment</li>
    </ol>
    <p>The QR/NFC tags should contain URLs like:</p>
    <code>http://spool-scanner/select/SPOOL_ID</code>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

// Spool selection page
func (h *Handler) handleSpoolSelect(w http.ResponseWriter, r *http.Request) {
	// Extract spool ID from URL
	// URL format: /select/SPOOL_123
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	spoolID := parts[2]
	if spoolID == "" {
		http.Error(w, "Spool ID required", http.StatusBadRequest)
		return
	}

	// We'll implement the template in the next step
	tmplStr := `
<!DOCTYPE html>
<html>
<head>
	<title>Select Printer - Spool {{.SpoolID}}</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="/static/style.css">
	<script src="//unpkg.com/alpinejs" defer></script>
</head>
<body>
	<div class="container" x-data="spoolSelector" x-init="init()">
		<!-- Loading State -->
		<div x-show="step === 'loading'" class="loading-container">
			<div class="loading-spinner"></div>
			<p>Loading spool information...</p>
		</div>

		<!-- Error State -->
		<div x-show="step === 'error'" class="message message-error">
			<strong>Error:</strong> <span x-text="error"></span>
		</div>

		<!-- Spool Info -->
		<div x-show="spool && (step === 'select-printer' || step === 'confirming')" class="spool-info">
			<div class="spool-header">
				<div class="spool-color" :style="{ backgroundColor: spool.color_hex || spool.color || '#888' }"></div>
				<div class="spool-title">
					<h2 x-text="spool.name || (spool.material + ' - ' + spool.vendor)"></h2>
					<p x-text="spool.vendor"></p>
				</div>
			</div>
			<div class="spool-details">
				<div class="spool-detail">
					<span class="spool-detail-label">Total Weight</span>
					<span class="spool-detail-value" x-text="formatWeight(spool.weight || 0)"></span>
				</div>
				<div class="spool-detail">
					<span class="spool-detail-label">Used</span>
					<span class="spool-detail-value" x-text="formatWeight(spool.used || 0)"></span>
				</div>
				<div class="spool-detail">
					<span class="spool-detail-label">Remaining</span>
					<span class="spool-detail-value" x-text="formatWeight(spool.remaining || 0)"></span>
				</div>
			</div>
		</div>

		<!-- Printer Selection -->
		<div x-show="step === 'select-printer'">
			<h3>Select Printer:</h3>
			<div class="printer-grid">
				<template x-for="printer in printers" :key="printer.id">
					<div class="printer-card"
						:class="{
							'selected': selectedPrinter && selectedPrinter.id === printer.id,
							'disabled': printer.status === 'Printing'
						}"
						@click="printer.status !== 'Printing' && selectPrinter(printer)"
					>
						<div class="printer-name" x-text="printer.name"></div>
						<div class="printer-status"
							:class="{
								'ready': printer.status === 'Ready',
								'printing': printer.status === 'Printing',
								'error': printer.status === 'Error'
							}"
						>
							<span x-text="getPrinterStatusEmoji(printer.status)"></span>
							<span x-text="printer.status"></span>
						</div>
						
						<!-- Current Spool Info -->
						<div x-show="printer.current_spool" class="current-spool">
							<div class="current-spool-label">Current:</div>
							<div class="current-spool-info">
								<span class="spool-color-dot" 
									:style="{ backgroundColor: printer.current_spool ? printer.current_spool.color : '#888' }">
								</span>
								<span x-text="printer.current_spool ? printer.current_spool.name : ''"></span>
							</div>
						</div>
					</div>
				</template>
			</div>

			<div x-show="error" class="message message-error" x-text="error"></div>

			<div class="button-group">
				<button 
					class="btn btn-primary" 
					:disabled="!selectedPrinter"
					@click="confirmAssignment()"
				>
					Confirm Selection
				</button>
				<a href="/" class="btn btn-secondary">Cancel</a>
			</div>
		</div>

		<!-- Confirming State -->
		<div x-show="step === 'confirming'" class="loading-container">
			<div class="loading-spinner"></div>
			<p>Assigning spool to <span x-text="selectedPrinter.name"></span>...</p>
		</div>

		<!-- Success State -->
		<div x-show="step === 'success'" class="success-container">
			<h2>Success!</h2>
			<p>Spool has been assigned to <strong x-text="selectedPrinter.name"></strong></p>
			<p class="message message-success">
				<span x-show="selectedPrinter.url">Redirecting to OctoPrint...</span>
				<span x-show="!selectedPrinter.url">Redirecting...</span>
			</p>
			<div x-show="selectedPrinter.url" style="margin-top: 20px;">
				<a :href="selectedPrinter.url" class="btn btn-primary">
					Go to OctoPrint Now
				</a>
			</div>
		</div>
	</div>

	<script>
		const SPOOL_ID = "{{.SpoolID}}";
	</script>
	<script src="/static/app.js"></script>
</body>
</html>
`

	tmpl, err := template.New("select").Parse(tmplStr)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := struct {
		SpoolID string
	}{
		SpoolID: spoolID,
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

// API: Get list of printers
func (h *Handler) handleGetPrinters(w http.ResponseWriter, r *http.Request) {
	type CurrentSpool struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Material string `json:"material"`
		Color    string `json:"color"`
	}

	type PrinterResponse struct {
		ID           string        `json:"id"`
		Name         string        `json:"name"`
		Status       string        `json:"status"`
		URL          string        `json:"url"`
		CurrentSpool *CurrentSpool `json:"current_spool,omitempty"`
	}

	printers := make([]PrinterResponse, len(h.config.Printers))
	for i, p := range h.config.Printers {
		status := "Unknown"
		var currentSpool *CurrentSpool

		if client, ok := h.octoprintClients[p.ID]; ok {
			// Get printer state
			if state, err := client.GetPrinterState(); err == nil {
				if state.Flags.Printing {
					status = "Printing"
				} else if state.Flags.Ready {
					status = "Ready"
				} else if state.Flags.Error {
					status = "Error"
				} else {
					status = state.Text
				}
			}

			// Get current spool (tool 0 for now)
			if spoolID, err := client.GetCurrentSpool(0); err == nil && spoolID != "" {
				// Fetch spool details from Spoolman
				if spool, err := h.spoolmanClient.GetSpool(spoolID); err == nil {
					color := "#888888"
					if spool.Filament.ColorHex != "" {
						color = spool.Filament.ColorHex
						if !strings.HasPrefix(color, "#") {
							color = "#" + color
						}
					}

					displayName := spool.Filament.Name
					if spool.Filament.Material != "" {
						displayName = fmt.Sprintf("%s | %s", displayName, spool.Filament.Material)
					}

					currentSpool = &CurrentSpool{
						ID:       spoolID,
						Name:     displayName,
						Material: spool.Filament.Material,
						Color:    color,
					}
				}
			}
		}

		printers[i] = PrinterResponse{
			ID:           p.ID,
			Name:         p.Name,
			Status:       status,
			URL:          p.OctoPrintURL,
			CurrentSpool: currentSpool,
		}
	}

	h.respondJSON(w, map[string]interface{}{
		"printers": printers,
	})
}

// API: Get spool details
func (h *Handler) handleGetSpool(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		h.respondError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	spoolID := parts[3]

	spool, err := h.spoolmanClient.GetSpool(spoolID)
	if err != nil {
		h.respondError(w, "Spool not found", http.StatusNotFound)
		return
	}

	// Debug logging
	log.Printf("Spool data from Spoolman: %+v", spool)
	log.Printf("Filament data: %+v", spool.Filament)

	// Extract color from ColorHex or use name as fallback
	color := "#888888" // Default gray
	if spool.Filament.ColorHex != "" {
		color = spool.Filament.ColorHex
		if !strings.HasPrefix(color, "#") {
			color = "#" + color
		}
	}

	displayName := spool.Filament.Name
	if spool.Filament.Material != "" {
		displayName = fmt.Sprintf("%s | %s", displayName, spool.Filament.Material)
	}

	log.Printf("Display Name: %+v", displayName)
	log.Printf("Material: %+v", spool.Filament.Material)

	// In handleGetSpool, update the response building:
	response := map[string]interface{}{
		"id":        spoolID,
		"name":      displayName,
		"material":  spool.Filament.Material,
		"color":     color,
		"color_hex": color,
		"vendor":    spool.Filament.Vendor.Name,
		// Use spool weight if available, otherwise use filament weight
		"weight":    spool.InitialWeight,
		"used":      spool.UsedWeight,
		"remaining": spool.InitialWeight - spool.UsedWeight,
		// Include filament weight as fallback
		"filament_weight": spool.Filament.Weight,
	}

	log.Printf("Sending response: %+v", response)

	h.respondJSON(w, response)
}

// API: Assign spool to printer
func (h *Handler) handleAssignSpool(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SpoolID   string `json:"spool_id"`
		PrinterID string `json:"printer_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the OctoPrint client for this printer
	client, ok := h.octoprintClients[req.PrinterID]
	if !ok {
		h.respondError(w, "Invalid printer ID", http.StatusBadRequest)
		return
	}

	// Set the active spool (tool 0 for now, can be extended)
	if err := client.SetActiveSpool(req.SpoolID, 0); err != nil {
		log.Printf("Failed to set spool: %v", err)
		h.respondError(w, "Failed to set spool", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Spool %s assigned to printer %s", req.SpoolID, req.PrinterID),
	})
}
