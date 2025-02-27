package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

func main() {
	// กำหนดชื่อไฟล์ Excel และไฟล์ output
	excelFileName := "Authen.xlsx"
	outputFileName := "authen_code.go"

	// เปิดไฟล์ Excel
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf("เปิดไฟล์ Excel ไม่สำเร็จ: %v\n", err)
		return
	}

	// สมมติว่า sheet แรกคือข้อมูลที่ต้องการ
	if len(xlFile.Sheets) == 0 {
		fmt.Println("ไม่พบ sheet ในไฟล์ Excel")
		return
	}
	sheet := xlFile.Sheets[0]

	// อ่านหัวข้อคอลัมน์ (แถวแรก)
	if len(sheet.Rows) == 0 {
		fmt.Println("ไม่พบข้อมูลในไฟล์ Excel")
		return
	}

	// หาตำแหน่งคอลัมน์ที่สำคัญ
	pidColumnIndex := 2 // สมมติว่า PID อยู่ในคอลัมน์แรก
	jsonColumnIndex := -1
	serviceDateColumnIndex := -1
	headers := []string{}

	// สร้าง slice เก็บหัวข้อคอลัมน์
	for i, cell := range sheet.Rows[0].Cells {
		headerText := cell.String()
		headers = append(headers, headerText)

		if headerText == "json" {
			jsonColumnIndex = i
		} else if headerText == "serviceDate serviceDate" {
			serviceDateColumnIndex = i
		}
	}

	// สร้าง map เก็บข้อมูลทั้งหมด โดยใช้ PID เป็น key
	data := make(map[string]interface{})

	// อ่านข้อมูลแต่ละแถว (ข้ามแถวแรกซึ่งเป็นหัวข้อ)
	for i := 1; i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]
		if len(row.Cells) == 0 {
			continue
		}

		// อ่าน PID จากคอลัมน์แรก
		pid := row.Cells[pidColumnIndex].String()
		if pid == "" {
			continue
		}

		rowData := make(map[string]interface{})

		if jsonColumnIndex != -1 && jsonColumnIndex < len(row.Cells) {
			// ถ้ามีคอลัมน์ json ให้ใช้ข้อมูลจากคอลัมน์นั้น
			jsonStr := row.Cells[jsonColumnIndex].String()
			if jsonStr != "" {
				// แปลง JSON string เป็น map
				var jsonData interface{}
				if err := json.Unmarshal([]byte(jsonStr), &jsonData); err == nil {
					if jsonMap, ok := jsonData.(map[string]interface{}); ok {
						rowData = jsonMap
					}
				} else {
					fmt.Printf("แปลง JSON สำหรับ PID %s ไม่สำเร็จ: %v\n", pid, err)
					// สร้างข้อมูลจากคอลัมน์อื่นๆ แทน
					rowData = createRowData(row, headers)
				}
			} else {
				// ถ้าคอลัมน์ json ว่างเปล่า ให้สร้างข้อมูลจากคอลัมน์อื่นๆ
				rowData = createRowData(row, headers)
			}
		} else {
			// ถ้าไม่มีคอลัมน์ json ให้สร้างข้อมูลจากคอลัมน์ทั้งหมด
			rowData = createRowData(row, headers)
		}

		// เพิ่มค่า serviceDate ถ้ามี
		if serviceDateColumnIndex != -1 && serviceDateColumnIndex < len(row.Cells) {
			serviceDate := row.Cells[serviceDateColumnIndex].String()
			if serviceDate != "" {
				rowData["serviceDate"] = serviceDate
			}
		}

		data[pid] = rowData
	}

	// สร้างโค้ด Go แบบ hard-coded
	code := generateGoCode(data)

	// บันทึกไฟล์โค้ด
	err = os.WriteFile(outputFileName, []byte(code), 0644)
	if err != nil {
		fmt.Printf("บันทึกไฟล์ไม่สำเร็จ: %v\n", err)
		return
	}

	fmt.Printf("สร้างโค้ด Go แล้ว บันทึกไว้ที่: %s\n", outputFileName)
	fmt.Printf("จำนวน records ทั้งหมด: %d\n", len(data))
}

// createRowData สร้าง map ข้อมูลจากแถวใน Excel
func createRowData(row *xlsx.Row, headers []string) map[string]interface{} {
	rowData := make(map[string]interface{})

	for j, cell := range row.Cells {
		if j < len(headers) && headers[j] != "PID" && headers[j] != "json" && headers[j] != "serviceDate" {
			// ลองแปลงเป็นตัวเลข
			if value, err := cell.Float(); err == nil {
				rowData[headers[j]] = value
			} else {
				// ถ้าแปลงเป็นตัวเลขไม่ได้ให้เก็บเป็นข้อความ
				rowData[headers[j]] = cell.String()
			}
		}
	}

	return rowData
}

// generateGoCode สร้างโค้ด Go แบบ hard-coded จากข้อมูล
func generateGoCode(data map[string]interface{}) string {
	// แปลงข้อมูลเป็น JSON แบบสวยงาม
	jsonData, _ := json.MarshalIndent(data, "", "    ")

	// สร้างโค้ด Go
	code := `package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// userData เก็บข้อมูลแบบ hard-coded จากไฟล์ Excel
var userData = map[string]interface{}{
`

	// แปลง JSON เป็นโค้ด Go
	jsonStr := string(jsonData)
	lines := strings.Split(jsonStr, "\n")

	// เพิ่มด้วย tab สำหรับทุกบรรทัดยกเว้นบรรทัดแรกและสุดท้าย
	for i, line := range lines {
		if i == 0 {
			// บรรทัดแรก (เปลี่ยน { เป็น empty ในโค้ด Go เพราะเราเพิ่ม { ไปแล้ว)
			continue
		} else if i == len(lines)-1 {
			// บรรทัดสุดท้าย
			code += "}\n\n"
		} else {
			// บรรทัดระหว่าง
			code += "\t" + line + "\n"
		}
	}

	// เพิ่มฟังก์ชันสำหรับ API
	code += `// Handler เป็นฟังก์ชันหลักที่จะถูกเรียกโดย Vercel
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
		// หน้าหลัก
		response := map[string]interface{}{
			"status":  "success",
			"message": "ยินดีต้อนรับสู่ API",
			"usage":   "ลองใช้ endpoint /api?pid=P001 หรือ /api?pid=P001&serviceDate=2025-02-15",
		}
		writeJSON(w, response)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api") {
		// API endpoint
		apiHandler(w, r)
		return
	}

	// กรณีไม่พบเส้นทาง
	w.WriteHeader(http.StatusNotFound)
	writeJSON(w, map[string]interface{}{
		"status": "error",
		"error":  "เส้นทางไม่ถูกต้อง",
	})
}

// writeJSON ช่วยเขียน JSON response ให้รองรับภาษาไทย
func writeJSON(w http.ResponseWriter, data interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(` + "`" + `{"status":"error","message":"เกิดข้อผิดพลาดในการสร้าง JSON"}` + "`" + `))
		return
	}
	w.Write(jsonBytes)
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
		// ค้นหาข้อมูลจาก userData
		data, err := findDataByPID(pid, serviceDate)
		if err == nil {
			// พบข้อมูล ส่งกลับไป
			writeJSON(w, data)
			return
		} else {
			// ไม่พบข้อมูล
			w.WriteHeader(http.StatusNotFound)
			writeJSON(w, map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("%v", err),
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

// findDataByPID ค้นหาข้อมูลตาม PID และ serviceDate
func findDataByPID(pid string, serviceDate string) (interface{}, error) {
	// ค้นหาข้อมูลตาม PID
	data, exists := userData[pid]
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
`

	return code
}
