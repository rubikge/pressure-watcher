package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPressureCheck_SignificantChange(t *testing.T) {
	mockCheckFunc := func() (bool, error) {
		return true, nil
	}

	w := httptest.NewRecorder()

	pressureCheckHandler(w, mockCheckFunc)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %v", res.Status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["significant_change"] != true {
		t.Errorf("Expected significant_change to be true, got %v", response["significant_change"])
	}
}

func TestPressureCheck_NoSignificantChange(t *testing.T) {
	mockCheckFunc := func() (bool, error) {
		return false, nil
	}

	w := httptest.NewRecorder()

	pressureCheckHandler(w, mockCheckFunc)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %v", res.Status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["significant_change"] != false {
		t.Errorf("Expected significant_change to be false, got %v", response["significant_change"])
	}
}

func TestPressureCheck_Error(t *testing.T) {
	mockCheckFunc := func() (bool, error) {
		return false, fmt.Errorf("mock error")
	}

	w := httptest.NewRecorder()

	pressureCheckHandler(w, mockCheckFunc)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status 500 Internal Server Error, got %v", res.Status)
	}
}
