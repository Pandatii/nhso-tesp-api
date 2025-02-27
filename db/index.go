package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// dbPath คือเส้นทางของไฟล์ฐานข้อมูล
var dbPath = filepath.Join("db", "Authen.db")

// Handler เป็นฟังก์ชันหลักที่จะถูกเรียกโดย Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// กำหนด CORS headers และ Content-Type เป็น UTF-8 อย่างชัดเจน
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// จัดการกับ OPTIONS request (สำหรับ CORS preflight)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// ตรวจสอบเส้นทาง
	if r.URL.Path == "/" || r.URL.Path == "" {
		// หน้าหลัก - ตรวจสอบว่ามีไฟล์ฐานข้อมูลหรือไม่
		dbExists := false
		if _, err := os.Stat(dbPath); err == nil {
			dbExists = true
		}

		// แสดงข้อมูลเกี่ยวกับฐานข้อมูล
		info := map[string]interface{}{
			"status":      "success",
			"message":     "ยินดีต้อนรับสู่ API",
			"usage":       "ลองใช้ endpoint /api?pid=P001 หรือ /api?pid=P001&serviceDate=2025-02-15",
			"db_exists":   dbExists,
			"db_path":     dbPath,
			"current_dir": mustGetwd(),
		}

		// ถ้ามีฐานข้อมูล ลองนับจำนวนแถว
		if dbExists {
			count, err := getRecordCount()
			if err == nil {
				info["record_count"] = count
			} else {
				info["db_error"] = err.Error()
			}
		}

		writeJSON(w, info)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/db") {
		// API endpoint
		apiHandler(w, r)
		return
	}

	// กรณีไม่พบเส้นทาง
	w.WriteHeader(http.StatusNotFound)
	writeJSON(w, map[string]interface{}{
		"status": "error",
		"error":  "เส้นทางไม่ถูกต้อง DB",
	})
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

// mustGetwd ดึงค่า working directory ปัจจุบัน
func mustGetwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return dir
}

// apiHandler จัดการ API endpoint
func apiHandler(w http.ResponseWriter, r *http.Request) {
	// รับค่า parameter จาก URL
	query := r.URL.Query()

	// ตรวจสอบค่า PID และ serviceDate
	pid := query.Get("pid")
	serviceDate := query.Get("serviceDate")

	// ถ้ามีค่า PID ให้ค้นหาข้อมูล
	if pid != "" {
		// ตรวจสอบว่ามีไฟล์ฐานข้อมูลหรือไม่
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			// ไม่พบไฟล์ฐานข้อมูล
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("ไม่พบไฟล์ฐานข้อมูล: %s", dbPath),
			})
			return
		}

		// ค้นหาข้อมูลจากฐานข้อมูล
		data, err := findDataByPID(pid, serviceDate)
		if err == nil {
			// พบข้อมูล ส่งกลับไป
			var jsonResponse interface{}
			if err := json.Unmarshal([]byte(data), &jsonResponse); err == nil {
				writeJSON(w, jsonResponse)
			} else {
				// ถ้าแปลง JSON ไม่ได้ ให้ส่งเป็นข้อความ
				w.Write([]byte(data))
			}
			return
		} else {
			// ไม่พบข้อมูล
			w.WriteHeader(http.StatusNotFound)
			writeJSON(w, map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("เกิดข้อผิดพลาด: %v", err),
			})
			return
		}
	}

	// ถ้าไม่ได้ระบุ PID ให้ส่งข้อความแจ้งเตือน
	w.WriteHeader(http.StatusBadRequest)
	writeJSON(w, map[string]interface{}{
		"status": "error",
		"error":  "กรุณาระบุค่า PID (เช่น: /api?pid=P001)",
		"info":   "สามารถระบุ serviceDate เพิ่มเติมได้ (เช่น: /api?pid=P001&serviceDate=2025-02-15)",
	})
}

// findDataByPID ค้นหาข้อมูลจากฐานข้อมูล SQLite ตาม PID และ serviceDate
func findDataByPID(pid string, serviceDate string) (string, error) {
	// เปิดฐานข้อมูล
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", fmt.Errorf("เปิดฐานข้อมูลไม่สำเร็จ: %v", err)
	}
	defer db.Close()

	// สร้าง query ตามเงื่อนไข
	var query string
	var args []interface{}

	if serviceDate == "" {
		// ค้นหาเฉพาะ PID
		query = "SELECT json_data FROM user_data WHERE pid = ?"
		args = []interface{}{pid}
	} else {
		// ค้นหาทั้ง PID และ serviceDate
		// ตัดเอาเฉพาะ 10 ตัวแรกของ serviceDate
		if len(serviceDate) > 10 {
			serviceDate = serviceDate[:10]
		}
		query = "SELECT json_data FROM user_data WHERE pid = ? AND substr(service_date, 1, 10) = ?"
		args = []interface{}{pid, serviceDate}
	}

	// ค้นหาข้อมูล
	var jsonData string
	err = db.QueryRow(query, args...).Scan(&jsonData)
	if err != nil {
		if err == sql.ErrNoRows {
			if serviceDate == "" {
				return "", fmt.Errorf("ไม่พบข้อมูลสำหรับ PID: %s", pid)
			} else {
				return "", fmt.Errorf("ไม่พบข้อมูลสำหรับ PID: %s และ serviceDate: %s", pid, serviceDate)
			}
		}
		return "", fmt.Errorf("ค้นหาข้อมูลไม่สำเร็จ: %v", err)
	}

	return jsonData, nil
}

// getRecordCount นับจำนวนแถวในฐานข้อมูล
func getRecordCount() (int, error) {
	// เปิดฐานข้อมูล
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return 0, fmt.Errorf("เปิดฐานข้อมูลไม่สำเร็จ: %v", err)
	}
	defer db.Close()

	// นับจำนวนแถว
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM user_data").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("นับจำนวนแถวไม่สำเร็จ: %v", err)
	}

	return count, nil
}
