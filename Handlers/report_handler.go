package handlers

import (
	"net/http"
	"os"
	"strings"
	"encoding/json"
	"path/filepath"
)

type ReportHandler struct{}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{}
}

func (h *ReportHandler) GetReports(w http.ResponseWriter, r *http.Request){
	files, err := os.ReadDir("reports")

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to read reports directory")
		return
	}

	var reports []string

	for _, file := range files {
		if strings.HasSuffix(file.Name(),".json"){
			reports = append(reports, file.Name())
		}
	}


	resp, err := json.Marshal(reports)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to serialize reports list")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)



}


func (h *ReportHandler) GetReportByDate(w http.ResponseWriter, r *http.Request) {
	date := strings.TrimPrefix(r.URL.Path, "/reports/")
	if date == "" {
		WriteError(w, http.StatusBadRequest, "missing report date")
		return
	}

	filename := "sales_report_" + date + ".json"
	path := filepath.Join("reports", filename)

	raw, err := os.ReadFile(path)
	if err != nil {
		WriteError(w, http.StatusNotFound, "report not found")
		return
	}

	// Unmarshal to validate JSON structure
	var report map[string]interface{}
	err = json.Unmarshal(raw, &report)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "invalid report format")
		return
	}

	// Marshal again for response
	data, err := json.Marshal(report)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to serialize report")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
