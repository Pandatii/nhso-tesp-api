// api/index.go
package api

import (
	"bytes"
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var userAuthenData = map[string]interface{}{}

// ฝังไฟล์ CSV ลงในโค้ด
//
//go:embed realPerson.csv
var realPersonData []byte

// Response โครงสร้างข้อมูลสำหรับ API response
type Response struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data,omitempty"`
}

// printLog ฟังก์ชันสำหรับแสดง log
func printLog(level string, message string, data interface{}) {
	fmt.Printf("[%s] %s: %v\n", level, message, data)
}

// findProductByPID ค้นหาข้อมูลจาก CSV โดยใช้ PID
func findProductByPID(pid string) (map[string]string, error) {
	// สร้าง reader สำหรับอ่านข้อมูล CSV จาก bytes
	reader := csv.NewReader(bytes.NewReader(realPersonData))

	// อ่านข้อมูล CSV ทั้งหมด
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	// ตรวจสอบว่ามีข้อมูลหรือไม่
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV is empty or only has headers")
	}

	// แถวแรกเป็น headers
	headers := records[0]

	// ค้นหาข้อมูลจาก column แรก (PID)
	for i := 1; i < len(records); i++ {
		if len(records[i]) > 0 && records[i][0] == pid {
			// พบข้อมูล - สร้าง map ของข้อมูล
			result := make(map[string]string)

			for j := 0; j < len(headers) && j < len(records[i]); j++ {
				result[headers[j]] = records[i][j]
			}

			return result, nil
		}
	}

	// ไม่พบข้อมูล
	return nil, fmt.Errorf("product with PID %s not found", pid)
}

// Handler function สำหรับ Vercel serverless
/*
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

		fmt.Println("API Authen")
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
*/

/*
// Handler function สำหรับ Vercel serverless
func Handler(w http.ResponseWriter, r *http.Request) {
	// ตั้งค่า response header
	w.Header().Set("Content-Type", "application/json")

	// สร้าง reader สำหรับอ่านข้อมูล CSV จาก bytes
	reader := csv.NewReader(bytes.NewReader(csvData))

	// อ่านข้อมูล CSV ทั้งหมด
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		http.Error(w, "Error reading CSV data", http.StatusInternalServerError)
		return
	}

	// ตรวจสอบว่ามีข้อมูลหรือไม่
	if len(records) == 0 {
		fmt.Println("CSV is empty")
		http.Error(w, "CSV data is empty", http.StatusInternalServerError)
		return
	}

	// สมมติว่าแถวแรกเป็น headers
	headers := records[0]

	// แปลงข้อมูลให้อยู่ในรูปแบบที่ใช้งานง่าย (array ของ objects)
	var data []map[string]string

	for i := 1; i < len(records); i++ {
		row := make(map[string]string)

		for j := 0; j < len(headers) && j < len(records[i]); j++ {
			row[headers[j]] = records[i][j]
		}

		data = append(data, row)
	}

	// สร้าง response
	response := Response{
		Message: "CSV data processed successfully",
		Data:    data,
		Total:   len(data),
	}

	// แปลงเป็น JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error creating JSON:", err)
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	// ส่ง response
	w.Write(jsonResponse)
}
*/

// readCSVFromFS อ่านไฟล์ CSV จาก embed.FS
func readCSVFromFS(fs embed.FS, filename string) ([][]string, error) {
	// อ่านไฟล์จาก embed.FS
	fileContent, err := fs.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading embedded file %s: %v", filename, err)
	}

	// อ่าน CSV จาก bytes
	reader := csv.NewReader(bytes.NewReader(fileContent))

	// อ่านทุกแถวจาก CSV
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error parsing CSV from %s: %v", filename, err)
	}

	return records, nil
}

// csvToMap แปลงข้อมูล CSV เป็น array ของ map
func csvToMap(records [][]string) ([]map[string]string, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("CSV data is empty")
	}

	// แถวแรกเป็น headers
	headers := records[0]

	// แปลงข้อมูล
	result := make([]map[string]string, 0, len(records)-1)

	for i := 1; i < len(records); i++ {
		row := make(map[string]string)

		for j := 0; j < len(headers) && j < len(records[i]); j++ {
			row[headers[j]] = records[i][j]
		}

		result = append(result, row)
	}

	return result, nil
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
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed. Only GET is supported.",
		})
		return
	}

	// รับค่า pid จาก query parameter
	pid := r.URL.Query().Get("pid")

	// ตรวจสอบว่ามีการระบุ pid หรือไม่
	if pid == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Missing required parameter: pid",
		})
		return
	}

	printLog("INFO", "Searching for product", pid)

	// ค้นหาข้อมูล
	product, err := findProductByPID(pid)
	if err != nil {
		// ตรวจสอบว่าเป็น error "not found" หรือไม่
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: fmt.Sprintf("Product with PID %s not found", pid),
			})
		} else {
			// กรณีเกิด error อื่นๆ
			printLog("ERROR", "Error finding product", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Error processing request: " + err.Error(),
			})
		}
		return
	}

	// ส่งข้อมูลกลับ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Product found",
		Data:    product,
	})
}

/*
// Handler function สำหรับ Vercel serverless
func Handler(w http.ResponseWriter, r *http.Request) {
	// ตั้งค่า response header
	w.Header().Set("Content-Type", "application/json")

	// ดึงชื่อไฟล์จาก query parameter (ถ้ามี)
	requestedFile := r.URL.Query().Get("file")

	// สร้าง response
	response := Response{
		Message:   "CSV data processed successfully",
		FileData:  make(map[string][]map[string]string),
		TotalRows: make(map[string]int),
		Files:     []string{},
	}

	// ลิสต์ไฟล์ทั้งหมดที่ฝังมา
	entries, err := csvFiles.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading embedded directory:", err)
		http.Error(w, "Error listing CSV files", http.StatusInternalServerError)
		return
	}

	// จำนวนไฟล์ที่พบ
	fileCount := 0

	// ทำงานกับไฟล์แต่ละไฟล์
	for _, entry := range entries {
		// ข้ามหากไม่ใช่ไฟล์หรือไม่ใช่ไฟล์ .csv
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".csv") {
			continue
		}

		// เพิ่มในรายการไฟล์
		response.Files = append(response.Files, entry.Name())
		fileCount++

		// ถ้ามีการระบุชื่อไฟล์และไม่ตรงกับไฟล์ปัจจุบัน ให้ข้ามไป
		if requestedFile != "" && requestedFile != entry.Name() {
			continue
		}

		// อ่านและแปลงข้อมูล CSV
		records, err := readCSVFromFS(csvFiles, entry.Name())
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", entry.Name(), err)
			continue
		}

		// แปลงเป็น array ของ map
		data, err := csvToMap(records)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", entry.Name(), err)
			continue
		}

		// เพิ่มข้อมูลลงใน response
		response.FileData[entry.Name()] = data
		response.TotalRows[entry.Name()] = len(data)
	}

	// อัปเดตจำนวนไฟล์
	response.FileCount = fileCount

	// กรณีระบุชื่อไฟล์แต่ไม่พบ
	if requestedFile != "" && len(response.FileData) == 0 {
		http.Error(w, fmt.Sprintf("Requested file '%s' not found", requestedFile), http.StatusNotFound)
		return
	}

	// แปลงเป็น JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error creating JSON:", err)
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	// ส่ง response
	w.Write(jsonResponse)
}*/

func makeAuthenData() {
	userAuthenData := make(map[string]interface{})

	userAuthenData["2182893981963"] = []map[string]interface{}{
		{"age": "32 ปี 5 เดือน 21 วัน", "birthDate": "1992-09-06", "firstName": "ทดสอบ10", "fullName": "ทดสอบ10 ทดสอบ", "lastName": "ทดสอบ", "mainInscl": "UCS", "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ", "middleName": "", "nationCode": "099", "nationDescription": "ไทย", "paidModel": "1", "personalId": "1443852933786", "provinceCode": "17", "provinceName": "สิงห์บุรี", "serviceDate": "2025-02-02 7:54:00", "sex": "หญิง", "statusAuthen": "TRUE", "statusMessage": "พบข้อมูลการ authen", "subInscl": "", "subInsclName": ""},
	}
}

/*
// apiHandler จัดการ API endpoint
func apiHandlerAuthen(w http.ResponseWriter, r *http.Request) {
	// รับค่า parameter จาก URL
	query := r.URL.Query()

	// ตรวจสอบค่า PID และ serviceDate
	pid := query.Get("personalId")
	serviceDate := query.Get("serviceDate")

	// ถ้ามีค่า PID ให้ค้นหาข้อมูลใน Excel
	if pid != "" {

		jsonData, err := findJSONByPIDFromExcel(pid, serviceDate)
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
*/

/*
func findJSONByPIDFromExcel(pid string, serviceDate string) (interface{}, error) {

	excelFileName := "Authen.xlsx"

	// เปิดไฟล์ Excel
	xlFile, err := xlsx.OpenFile(excelFileName)
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

*/

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
