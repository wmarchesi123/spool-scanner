package config

import (
	"fmt"
	"os"
)

type Config struct {
	SpoolmanURL string    `json:"spoolman_url"`
	Printers    []Printer `json:"printers"`
}

type Printer struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	OctoPrintURL string `json:"octoprint_url"`
	APIKey       string `json:"api_key"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		SpoolmanURL: os.Getenv("SPOOLMAN_URL"),
		Printers:    []Printer{},
	}

	// Load printers from environment
	// Format: PRINTER_1_NAME=Kitchen,PRINTER_1_URL=http://...,PRINTER_1_KEY=...
	for i := 1; i <= 10; i++ { // Support up to 10 printers
		name := os.Getenv(fmt.Sprintf("PRINTER_%d_NAME", i))
		if name == "" {
			continue
		}

		printer := Printer{
			ID:           fmt.Sprintf("printer-%d", i),
			Name:         name,
			OctoPrintURL: os.Getenv(fmt.Sprintf("PRINTER_%d_URL", i)),
			APIKey:       os.Getenv(fmt.Sprintf("PRINTER_%d_KEY", i)),
		}

		if printer.OctoPrintURL == "" || printer.APIKey == "" {
			return nil, fmt.Errorf("incomplete configuration for printer %d", i)
		}

		cfg.Printers = append(cfg.Printers, printer)
	}

	if len(cfg.Printers) == 0 {
		return nil, fmt.Errorf("no printers configured")
	}

	if cfg.SpoolmanURL == "" {
		return nil, fmt.Errorf("SPOOLMAN_URL not set")
	}

	return cfg, nil
}
