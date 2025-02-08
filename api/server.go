package api

import (
	"encoding/json"
	"net/http"
)

type PressureCheckFunc func() (bool, error)

func StartServer(checkFunc PressureCheckFunc) {
	http.HandleFunc("/pressure-check", func(w http.ResponseWriter, r *http.Request) {
		pressureCheckHandler(w, checkFunc)
	})
	http.ListenAndServe(":8080", nil)
}

func pressureCheckHandler(w http.ResponseWriter, checkFunc PressureCheckFunc) {
	significantChange, err := checkFunc()

	if err != nil {
		http.Error(w, "Error retrieving data", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"significant_change": significantChange,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
