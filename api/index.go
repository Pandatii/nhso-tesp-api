// api/index.go
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var userAuthenData = map[string]interface{}{}

// Response โครงสร้างข้อมูลสำหรับ API response
type Response struct {
	Message   string    `json:"message"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

// Handler function สำหรับ Vercel serverless
func Handler(w http.ResponseWriter, r *http.Request) {
	// ตั้งค่า response header
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// จัดการกับ OPTIONS request (สำหรับ CORS preflight)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	makeAuthenData()

	writeJSON(w, userAuthenData)

	// สร้าง response object
	response := Response{
		Message:   "Hello from Go API on Vercel!",
		Version:   "1.0.0",
		Timestamp: time.Now(),
	}

	if strings.HasPrefix(r.URL.Path, "/authencodeapi") {
		// API endpoint
		apiHandlerAuthen(w, r)
		return
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

func makeAuthenData() {
	userAuthenData := make(map[string]interface{})

	userAuthenData["2182893981963"] = []map[string]interface{}{
		{"age": "32 ปี 5 เดือน 21 วัน", "birthDate": "1992-09-06", "firstName": "ทดสอบ10", "fullName": "ทดสอบ10 ทดสอบ", "lastName": "ทดสอบ", "mainInscl": "UCS", "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ", "middleName": "", "nationCode": "099", "nationDescription": "ไทย", "paidModel": "1", "personalId": "1443852933786", "provinceCode": "17", "provinceName": "สิงห์บุรี", "serviceDate": "2025-02-02 7:54:00", "sex": "หญิง", "statusAuthen": "TRUE", "statusMessage": "พบข้อมูลการ authen", "subInscl": "", "subInsclName": ""},
	}
}

// apiHandler จัดการ API endpoint
func apiHandlerAuthen(w http.ResponseWriter, r *http.Request) {
	// รับค่า parameter จาก URL
	query := r.URL.Query()

	// ตรวจสอบค่า PID และ serviceDate
	pid := query.Get("personalId")
	serviceDate := query.Get("serviceDate")

	// ถ้ามีค่า PID ให้ค้นหาข้อมูลใน Excel
	if pid != "" {

		jsonData, err := findJSONByPIDFromAuthenData(pid, serviceDate)
		if err == nil {
			// ส่งข้อมูล JSON กลับไปโดยตรง
			writeJSON(w, jsonData)
			return
		} else {
			// ถ้าไม่พบข้อมูลหรือมีข้อผิดพลาด
			errorResponse := map[string]string{
				"error": fmt.Sprintf("เกิดข้อผิดพลาด: %v", err),
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	// ถ้าไม่ได้ระบุ PID ให้ส่งข้อความแจ้งเตือน
	errorResponse := map[string]string{
		"error": "กรุณาระบุค่า PID (เช่น: /api?pid=P001)",
		"info":  "สามารถระบุ serviceDate เพิ่มเติมได้ (เช่น: /api?pid=P001&serviceDate=2025-02-15)",
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errorResponse)
}

// findDataByPID ค้นหาข้อมูลตาม PID และ serviceDate
func findJSONByPIDFromAuthenData(pid string, serviceDate string) (interface{}, error) {
	// ค้นหาข้อมูลตาม PID
	data, exists := userAuthenData[pid]
	if !exists {
		return nil, fmt.Errorf("ไม่พบข้อมูลสำหรับ PID: %s", pid)
	}

	// ถ้ามีการระบุ serviceDate ให้ตรวจสอบว่าตรงกันหรือไม่
	if serviceDate != "" {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("ข้อมูลไม่อยู่ในรูปแบบที่ถูกต้อง")
		}

		// ตรวจสอบว่ามี serviceDate หรือไม่
		dateValue, hasDate := dataMap["serviceDate"]
		if !hasDate {
			return nil, fmt.Errorf("ไม่พบข้อมูล serviceDate สำหรับ PID: %s", pid)
		}

		// ตรวจสอบว่า serviceDate ตรงกันหรือไม่ (เปรียบเทียบเฉพาะ 10 ตัวแรก)
		storedDate, ok := dateValue.(string)
		if !ok {
			return nil, fmt.Errorf("ข้อมูล serviceDate ไม่อยู่ในรูปแบบที่ถูกต้อง")
		}

		// ตัดเอาเฉพาะ 10 ตัวแรก (YYYY-MM-DD)
		storedDatePrefix := ""
		if len(storedDate) >= 10 {
			storedDatePrefix = storedDate[:10]
		} else {
			storedDatePrefix = storedDate
		}

		serviceDatePrefix := ""
		if len(serviceDate) >= 10 {
			serviceDatePrefix = serviceDate[:10]
		} else {
			serviceDatePrefix = serviceDate
		}

		// ถ้า serviceDate ไม่ตรงกัน
		if storedDatePrefix != serviceDatePrefix {
			return nil, fmt.Errorf("ไม่พบข้อมูลสำหรับ PID: %s และ serviceDate: %s", pid, serviceDate)
		}
	}

	return data, nil
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
