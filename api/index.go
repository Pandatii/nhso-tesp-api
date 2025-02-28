package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// userData เก็บข้อมูลแบบ hard-coded จากไฟล์ Excel
var userAuthenData = map[string]interface{}{
	
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// กำหนด CORS headers
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
		// หน้าหลัก - รายงานผลการค้นหาไฟล์
		response := map[string]interface{}{
			"status":        "success",
			"message":       "ยินดีต้อนรับสู่ API",
			"usage":         "ลองใช้ endpoint /api?pid=P001 หรือ /api?pid=P001&serviceDate=2025-02-15"
		}
		writeJSON(w, response)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api") {
		// API endpoint
		apiHandlerRealPerson(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/authencodeapi") {
		// API endpoint
		apiHandlerAuthen(w, r)
		return
	}

	// กรณีไม่พบเส้นทาง
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "เส้นทางไม่ถูกต้อง",
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

// homeHandler จัดการเส้นทางหลัก
func homeHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "ยินดีต้อนรับสู่ API ลองใช้ endpoint /api?pid=P001 หรือ /api?pid=P001&serviceDate=2025-02-15 เพื่อดูข้อมูล JSON",
	})
}

// apiHandler จัดการ API endpoint
func apiHandlerRealPerson(w http.ResponseWriter, r *http.Request) {
	// รับค่า parameter จาก URL
	query := r.URL.Query()

	// ตรวจสอบค่า PID และ serviceDate
	pid := query.Get("PID")

	// ถ้ามีค่า PID ให้ค้นหาข้อมูลใน Excel
	if pid != "" {
		// ตำแหน่งไฟล์ Excel (ต้องอยู่ใน folder api)
		excelPath := "./Mockup API Authen&realPerson_RealPerson.xlsx"

		jsonData, err := findJSONByPID(pid, excelPath)
		if err == nil {
			// ส่งข้อมูล JSON กลับไปโดยตรง
			w.Write(jsonData)
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
			w.Write(jsonData)
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

// findJSONByPID ค้นหาข้อมูล JSON จาก Excel โดยใช้ PID และ serviceDate
func findJSONByPID(pid string, excelPath string) (json.RawMessage, error) {
	// เปิดไฟล์ Excel
	xlFile, err := xlsx.OpenFile(excelPath)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถเปิดไฟล์ Excel ได้: %v", err)
	}

	// สมมติว่า sheet แรกคือที่เราต้องการ
	if len(xlFile.Sheets) == 0 {
		return nil, errors.New("ไม่พบ sheet ในไฟล์ Excel")
	}
	sheet := xlFile.Sheets[0]

	// อ่านหัวข้อคอลัมน์ (สมมติว่าอยู่ในแถวแรก)
	if len(sheet.Rows) == 0 {
		return nil, errors.New("Excel ไม่มีข้อมูล")
	}

	jsonColumnIndex := -1 // ตำแหน่งคอลัมน์ json

	// ค้นหาตำแหน่งคอลัมน์ที่ต้องการ
	for i, cell := range sheet.Rows[0].Cells {
		headerText := cell.String()
		if headerText == "json" {
			jsonColumnIndex = i
		}
	}

	if jsonColumnIndex == -1 {
		return nil, errors.New("ไม่พบคอลัมน์ 'json' ในไฟล์ Excel")
	}

	// ค้นหาแถวที่มีค่า PID ตรงกัน
	for i := 1; i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]
		if len(row.Cells) == 0 {
			continue
		}

		// ตรวจสอบว่า PID ตรงกับที่ต้องการหรือไม่
		rowPID := row.Cells[0].String()
		if rowPID == pid {

			// ตรวจสอบว่ามีข้อมูลในคอลัมน์ json หรือไม่
			if jsonColumnIndex < len(row.Cells) {
				jsonStr := row.Cells[jsonColumnIndex].String()
				if jsonStr == "" {
					return json.RawMessage("{}"), nil
				}

				// ตรวจสอบว่าข้อมูลเป็น JSON ที่ถูกต้องหรือไม่
				var jsonData json.RawMessage
				if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
					return nil, fmt.Errorf("ข้อมูลในคอลัมน์ 'json' ไม่ใช่ JSON ที่ถูกต้อง: %v", err)
				}

				return jsonData, nil
			}
			return nil, errors.New("ไม่พบข้อมูลในคอลัมน์ 'json'")
		}
	}

	return nil, fmt.Errorf("not found PID: %s", pid)
}

// findJSONByPID ค้นหาข้อมูล JSON จาก Excel โดยใช้ PID และ serviceDate
func findJSONByPIDCol2(pid string, serviceDate string, excelPath string) (json.RawMessage, error) {
	// เปิดไฟล์ Excel
	xlFile, err := xlsx.OpenFile(excelPath)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถเปิดไฟล์ Excel ได้: %v", err)
	}

	// สมมติว่า sheet แรกคือที่เราต้องการ
	if len(xlFile.Sheets) == 0 {
		return nil, errors.New("ไม่พบ sheet ในไฟล์ Excel")
	}
	sheet := xlFile.Sheets[0]

	// อ่านหัวข้อคอลัมน์ (สมมติว่าอยู่ในแถวแรก)
	if len(sheet.Rows) == 0 {
		return nil, errors.New("Excel ไม่มีข้อมูล")
	}

	jsonColumnIndex := -1        // ตำแหน่งคอลัมน์ json
	serviceDateColumnIndex := -1 // ตำแหน่งคอลัมน์ serviceDate

	// ค้นหาตำแหน่งคอลัมน์ที่ต้องการ
	for i, cell := range sheet.Rows[0].Cells {
		headerText := cell.String()
		if headerText == "json" {
			jsonColumnIndex = i
		} else if headerText == "serviceDate serviceDate" {
			serviceDateColumnIndex = i
		}
	}

	if jsonColumnIndex == -1 {
		return nil, errors.New("ไม่พบคอลัมน์ 'json' ในไฟล์ Excel")
	}

	if serviceDateColumnIndex == -1 && serviceDate != "" {
		return nil, errors.New("ไม่พบคอลัมน์ 'serviceDate' ในไฟล์ Excel")
	}

	// ค้นหาแถวที่มีค่า PID ตรงกัน
	for i := 1; i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]
		if len(row.Cells) == 0 {
			continue
		}

		// ตรวจสอบว่า PID ตรงกับที่ต้องการหรือไม่
		rowPID := row.Cells[2].String()
		if rowPID == pid {
			// ถ้ามีการระบุ serviceDate ให้ตรวจสอบว่าตรงกันหรือไม่
			if serviceDate != "" && serviceDateColumnIndex < len(row.Cells) {
				rowServiceDate := row.Cells[serviceDateColumnIndex].String()
				// ตัดเอาเฉพาะ 10 ตัวแรก (YYYY-MM-DD)
				rowServiceDatePrefix := ""
				if len(rowServiceDate) >= 10 {
					rowServiceDatePrefix = rowServiceDate[:10]
				} else {
					rowServiceDatePrefix = rowServiceDate
				}

				serviceDatePrefix := ""
				if len(serviceDate) >= 10 {
					serviceDatePrefix = serviceDate[:10]
				} else {
					serviceDatePrefix = serviceDate
				}

				// ถ้า serviceDate ไม่ตรงกัน ให้ข้ามแถวนี้
				if rowServiceDatePrefix != serviceDatePrefix {
					continue
				}
			}

			// ตรวจสอบว่ามีข้อมูลในคอลัมน์ json หรือไม่
			if jsonColumnIndex < len(row.Cells) {
				jsonStr := row.Cells[jsonColumnIndex].String()
				if jsonStr == "" {
					return json.RawMessage("{}"), nil
				}

				// ตรวจสอบว่าข้อมูลเป็น JSON ที่ถูกต้องหรือไม่
				var jsonData json.RawMessage
				if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
					return nil, fmt.Errorf("ข้อมูลในคอลัมน์ 'json' ไม่ใช่ JSON ที่ถูกต้อง: %v", err)
				}

				return jsonData, nil
			}
			return nil, errors.New("ไม่พบข้อมูลในคอลัมน์ 'json'")
		}
	}

	// ถ้า serviceDate ถูกระบุแต่ไม่พบข้อมูลที่ตรงกัน
	if serviceDate != "" {
		return nil, fmt.Errorf("ไม่พบข้อมูลสำหรับ PID: %s และ serviceDate: %s", pid, serviceDate)
	}

	return nil, fmt.Errorf("ไม่พบข้อมูลสำหรับ PID: %s", pid)
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
