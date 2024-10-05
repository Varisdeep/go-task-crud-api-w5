package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Device represents a mobile device
type Device struct {
	ID    int    `json:"id"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	Color string `json:"color"`
	RAM   string `json:"ram"`
	ROM   string `json:"rom"`
}

// In-memory storage for devices
var devices []Device
var mu sync.Mutex
var nextID = 1

// Create a new device
func createDevice(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var device Device
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	device.ID = nextID
	nextID++
	devices = append(devices, device)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(device)
}

// Get all devices
func getDevices(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// Get a device by ID
func getDeviceByID(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)
	id := vars["id"]

	for _, device := range devices {
		if device.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(device)
			return
		}
	}
	http.NotFound(w, r)
}

// Update an existing device
func updateDevice(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)
	id := vars["id"]

	for i, device := range devices {
		if device.ID == id {
			var updatedDevice Device
			if err := json.NewDecoder(r.Body).Decode(&updatedDevice); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			updatedDevice.ID = device.ID
			devices[i] = updatedDevice
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedDevice)
			return
		}
	}
	http.NotFound(w, r)
}

// Delete a device
func deleteDevice(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)
	id := vars["id"]

	for i, device := range devices {
		if device.ID == id {
			devices = append(devices[:i], devices[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}

// Main function
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/devices", createDevice).Methods("POST")
	r.HandleFunc("/devices", getDevices).Methods("GET")
	r.HandleFunc("/devices/{id}", getDeviceByID).Methods("GET")
	r.HandleFunc("/devices/{id}", updateDevice).Methods("PUT")
	r.HandleFunc("/devices/{id}", deleteDevice).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}
