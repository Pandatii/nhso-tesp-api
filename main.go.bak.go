package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/tealeg/xlsx"
)

// Message โครงสร้างข้อมูลสำหรับการตอบกลับ
type Message struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// ExcelData โครงสร้างข้อมูลสำหรับข้อมูลที่ได้จาก Excel
type ExcelData struct {
	PID        string            `json:"pid"`
	Fields     map[string]string `json:"fields"`
	JsonString string            `json:"json,omitempty"`
}

// findJSONByPID ค้นหาข้อมูล JSON จาก Excel โดยใช้ PID
func findJSONByPIDCol3(pid string, excelPath string, serviceDate string) (json.RawMessage, error) {
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

// findJSONByPID ค้นหาข้อมูล JSON จาก Excel โดยใช้ PID
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

	// ค้นหาตำแหน่งคอลัมน์ json
	for i, cell := range sheet.Rows[0].Cells {
		headerText := cell.String()
		if headerText == "json" {
			jsonColumnIndex = i
			break
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

// handleAPI จัดการ endpoint หลักของ API
func handleAPIAuthen(w http.ResponseWriter, r *http.Request) {
	// กำหนด header
	w.Header().Set("Content-Type", "application/json")

	// รับค่า parameter จาก URL
	query := r.URL.Query()

	// ตรวจสอบค่า PID
	pid := query.Get("personalId")
	serviceDate := query.Get("serviceDate")

	// ถ้ามีค่า PID ให้ค้นหาข้อมูลใน Excel
	if pid != "" {
		// ตำแหน่งไฟล์ Excel (อาจจะต้องปรับเปลี่ยนตามการเก็บไฟล์ใน Vercel)
		excelPath := "data/Mockup API Authen&realPerson_Authen.xlsx"

		data, err := findJSONByPIDCol3(pid, excelPath, serviceDate)
		if err == nil {
			// ส่งข้อมูลที่พบกลับไปโดยตรง
			json.NewEncoder(w).Encode(data)
			return
		} else {
			// ถ้าไม่พบข้อมูลหรือมีข้อผิดพลาด
			errorResponse := map[string]string{
				"error": fmt.Sprintf("ไม่พบข้อมูลสำหรับ PID: %s", pid),
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	// ถ้าไม่ได้ระบุ PID ให้ส่งข้อความแจ้งเตือน
	errorResponse := map[string]string{
		"error": "กรุณาระบุค่า PID (เช่น: /api?pid=P001)",
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errorResponse)
}

// handleAPI จัดการ endpoint หลักของ API
func handleAPI(w http.ResponseWriter, r *http.Request) {
	// กำหนด header
	w.Header().Set("Content-Type", "application/json")

	// รับค่า parameter จาก URL
	query := r.URL.Query()

	// ตรวจสอบค่า PID
	pid := query.Get("PID")

	// ถ้ามีค่า PID ให้ค้นหาข้อมูลใน Excel
	if pid != "" {
		// ตำแหน่งไฟล์ Excel (อาจจะต้องปรับเปลี่ยนตามการเก็บไฟล์ใน Vercel)
		excelPath := "data/Mockup API Authen&realPerson_RealPerson.xlsx"

		data, err := findJSONByPID(pid, excelPath)
		if err == nil {
			// ส่งข้อมูลที่พบกลับไปโดยตรง
			json.NewEncoder(w).Encode(data)
			return
		} else {
			// ถ้าไม่พบข้อมูลหรือมีข้อผิดพลาด
			errorResponse := map[string]string{
				"error": fmt.Sprintf("ไม่พบข้อมูลสำหรับ PID: %s", pid),
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	// ถ้าไม่ได้ระบุ PID ให้ส่งข้อความแจ้งเตือน
	errorResponse := map[string]string{
		"error": "กรุณาระบุค่า PID (เช่น: /api?pid=P001)",
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errorResponse)
}

// handleRoot จัดการเส้นทางหลัก
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	msg := Message{
		Status:    "success",
		Message:   "ยินดีต้อนรับสู่ API ลองใช้ endpoint /api เพื่อดูข้อมูล",
		Timestamp: time.Now(),
	}

	json.NewEncoder(w).Encode(msg)
}

func main() {
	// รับค่า port จากตัวแปรสภาพแวดล้อมหรือใช้ค่าเริ่มต้น
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ตั้งค่าเส้นทาง
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/api/RealPerson", handleAPI)
	http.HandleFunc("/authencodeapi/CheckAuthenStatus", handleAPIAuthen)

	fmt.Println("เริ่มต้นเซิร์ฟเวอร์ที่พอร์ต:", port)
	// เริ่มเซิร์ฟเวอร์
	http.ListenAndServe(":"+port, nil)
}
