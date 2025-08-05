package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func splitLines(s string) []string {
	// Supports both \n and \r\n line endings
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, "\r")
	}
	// Optionally, filter out empty lines
	var result []string
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			result = append(result, l)
		}
	}
	return result
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	format := r.URL.Query().Get("f")
	if format == "" {
		format = "json"
	}

	if format != "json" && format != "csv" {
		http.Error(w, "format must be 'json' or 'csv'", http.StatusBadRequest)
		return
	}

	if url == "" {
		http.Error(w, "missing url parameter", http.StatusBadRequest)
		return
	}

	// Fetch URL contents
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "failed to fetch url", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read response", http.StatusInternalServerError)
		return
	}
	lines := splitLines(string(body))

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lines)
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		cw := csv.NewWriter(w)
		for _, line := range lines {
			_ = cw.Write([]string{line})
		}
		cw.Flush()
	}
}

func main() {
	http.HandleFunc("/tojson", handleGet)
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
