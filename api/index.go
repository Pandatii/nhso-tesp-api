// api/index.go
package api

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ฝังไฟล์ CSV ลงในโค้ด
//
//go:embed realPerson.csv
var realPersonData []byte

//go:embed authen.csv
var authenData []byte

// ErrorResponse โครงสร้างข้อมูลสำหรับ error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// printLog ฟังก์ชันสำหรับแสดง log
func printLog(level string, message string, data interface{}) {
	fmt.Printf("[%s] %s: %v\n", level, message, data)
}

// writeJSON ช่วยเขียน JSON response ให้รองรับภาษาไทย
func writeJSON(w http.ResponseWriter, data interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","message":"เกิดข้อผิดพลาดในการสร้าง JSON"}`))
		return
	}
	w.Write(jsonBytes)
}

func apiHandlerAuthen(w http.ResponseWriter, r *http.Request) {

	// รับค่า pid จาก query parameter
	pid := r.URL.Query().Get("pid")

	serviceDate := r.URL.Query().Get("serviceDate")

	// ตรวจสอบว่ามีการระบุ pid หรือไม่
	if pid == "" || serviceDate == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Missing required parameter: pid or serviceDate",
		})
		return
	}

	printLog("INFO", "Searching for message with PID", pid)

	// สร้าง reader สำหรับอ่านข้อมูล CSV
	reader := csv.NewReader(bytes.NewReader(authenData))

	// อ่านข้อมูล CSV ทั้งหมด
	records, err := reader.ReadAll()
	if err != nil {
		printLog("ERROR", "Failed to read CSV", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Error reading data source: " + err.Error(),
		})
		return
	}

	// ตรวจสอบว่ามีข้อมูลหรือไม่
	if len(records) < 1 {
		printLog("ERROR", "CSV is empty", nil)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Data source is empty",
		})
		return
	}

	// ค้นหาข้อมูลจาก column แรก (PID)
	for i := 0; i < len(records); i++ {
		// ตรวจสอบว่ามีข้อมูลอย่างน้อย 2 columns
		if len(records[i]) < 2 {
			continue
		}

		// ตรวจสอบว่า PID ตรงกับที่ต้องการหรือไม่
		if records[i][0] == pid && strings.HasPrefix(records[i][1], serviceDate) {
			// พบข้อมูล - ดึง JSON message จาก column ที่สอง
			jsonMessage := records[i][2]

			// ตรวจสอบว่า JSON message ถูกต้องหรือไม่
			var jsonData interface{}
			err := json.Unmarshal([]byte(jsonMessage), &jsonData)
			if err != nil {
				printLog("ERROR", "Invalid JSON in data source", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{
					Success: false,
					Message: "Invalid JSON format in data source: " + err.Error(),
				})
				return
			}

			// ส่ง JSON message กลับไปโดยตรง
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// เขียน JSON โดยตรงโดยไม่ผ่าน struct
			w.Write([]byte(jsonMessage))
			return
		}
	}

	// ไม่พบข้อมูล
	printLog("INFO", "No data found for PID", pid)
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Message: fmt.Sprintf("No data found for PID: %s", pid),
	})
}

func apiHandlerRealPerson(w http.ResponseWriter, r *http.Request) {
	// รับค่า pid จาก query parameter
	pid := r.URL.Query().Get("PID")

	// ตรวจสอบว่ามีการระบุ pid หรือไม่
	if pid == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Missing required parameter: pid",
		})
		return
	}

	printLog("INFO", "Searching for message with PID", pid)

	// สร้าง reader สำหรับอ่านข้อมูล CSV
	reader := csv.NewReader(bytes.NewReader(realPersonData))

	// อ่านข้อมูล CSV ทั้งหมด
	records, err := reader.ReadAll()
	if err != nil {
		printLog("ERROR", "Failed to read CSV", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Error reading data source: " + err.Error(),
		})
		return
	}

	// ตรวจสอบว่ามีข้อมูลหรือไม่
	if len(records) < 1 {
		printLog("ERROR", "CSV is empty", nil)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Data source is empty",
		})
		return
	}

	// ค้นหาข้อมูลจาก column แรก (PID)
	for i := 0; i < len(records); i++ {
		// ตรวจสอบว่ามีข้อมูลอย่างน้อย 2 columns
		if len(records[i]) < 2 {
			continue
		}

		// ตรวจสอบว่า PID ตรงกับที่ต้องการหรือไม่
		if records[i][0] == pid {
			// พบข้อมูล - ดึง JSON message จาก column ที่สอง
			jsonMessage := records[i][1]

			// ตรวจสอบว่า JSON message ถูกต้องหรือไม่
			var jsonData interface{}
			err := json.Unmarshal([]byte(jsonMessage), &jsonData)
			if err != nil {
				printLog("ERROR", "Invalid JSON in data source", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{
					Success: false,
					Message: "Invalid JSON format in data source: " + err.Error(),
				})
				return
			}

			// ส่ง JSON message กลับไปโดยตรง
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// เขียน JSON โดยตรงโดยไม่ผ่าน struct
			w.Write([]byte(jsonMessage))
			return
		}
	}

	// ไม่พบข้อมูล
	printLog("INFO", "No data found for PID", pid)
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Message: fmt.Sprintf("No data found for PID: %s", pid),
	})
}

// Handler function สำหรับ Vercel serverless
func Handler(w http.ResponseWriter, r *http.Request) {
	// ตั้งค่า response header
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	// จัดการกับ CORS preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// ตรวจสอบว่าเป็น GET request
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Method not allowed. Only GET is supported.",
		})
		return
	}

	if strings.HasPrefix(r.URL.Path, "/authencodeapi") {

		fmt.Println("API Authen")
		// API endpoint
		apiHandlerAuthen(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api") {

		fmt.Println("API Authen")
		// API endpoint
		apiHandlerRealPerson(w, r)
		return
	}

	// กรณีไม่พบเส้นทาง
	w.WriteHeader(http.StatusNotFound)
	writeJSON(w, map[string]interface{}{
		"status": "error",
		"error":  "เส้นทางไม่ถูกต้อง",
	})
}
