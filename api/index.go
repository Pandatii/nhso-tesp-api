package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/tealeg/xlsx"
)

// Handler เป็นฟังก์ชันหลักที่จะถูกเรียกโดย Vercel
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
	if r.URL.Path == "/" {
		// หน้าหลัก
		homeHandler(w, r)
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
		// ตำแหน่งไฟล์ Excel (ต้องอยู่ใน folder api)
		excelPath := "../data/Mockup API Authen&realPerson_Authen.xlsx"

		jsonData, err := findJSONByPIDCol2(pid, serviceDate, excelPath)
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

	return nil, fmt.Errorf("ไม่พบข้อมูลสำหรับ PID: %s", pid)
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
		} else if headerText == "serviceDate" {
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
