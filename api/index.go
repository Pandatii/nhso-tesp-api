// api/index.go
package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// Response โครงสร้างข้อมูลสำหรับ API response
type Response struct {
	Message   string    `json:"message"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

// Handler function สำหรับ Vercel serverless
func Handler(w http.ResponseWriter, r *http.Request) {
	// ตั้งค่า response header
	w.Header().Set("Content-Type", "application/json")
	
	// สร้าง response object
	response := Response{
		Message:   "Hello from Go API on Vercel!",
		Version:   "1.0.0",
		Timestamp: time.Now(),
	}
	
	// แปลง response เป็น JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	
	// เขียน response
	w.Write(jsonResponse)
}
