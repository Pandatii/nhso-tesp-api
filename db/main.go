package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tealeg/xlsx"
)

func main() {
	// กำหนดชื่อไฟล์ Excel และ SQLite
	excelFileName := "Authen.xlsx"
	dbFileName := "Authen.db"

	// ลบไฟล์ฐานข้อมูลเดิมหากมีอยู่
	os.Remove(dbFileName)

	// เปิดหรือสร้างฐานข้อมูล SQLite
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		fmt.Printf("เปิดฐานข้อมูลไม่สำเร็จ: %v\n", err)
		return
	}
	defer db.Close()

	// สร้างตาราง
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS user_data (
		pid TEXT PRIMARY KEY,
		json_data TEXT NOT NULL,
		service_date TEXT
	)
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Printf("สร้างตารางไม่สำเร็จ: %v\n", err)
		return
	}

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
	pidColumnIndex := 2
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

	// เตรียม statement สำหรับ insert
	insertStmt, err := db.Prepare("INSERT INTO user_data (pid, json_data, service_date) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Printf("เตรียม statement ไม่สำเร็จ: %v\n", err)
		return
	}
	defer insertStmt.Close()

	// อ่านข้อมูลแต่ละแถว (ข้ามแถวแรกซึ่งเป็นหัวข้อ)
	recordCount := 0
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

		var jsonData []byte
		var serviceDate string

		// ดึงข้อมูล serviceDate ถ้ามี
		if serviceDateColumnIndex != -1 && serviceDateColumnIndex < len(row.Cells) {
			serviceDate = row.Cells[serviceDateColumnIndex].String()
		}

		if jsonColumnIndex != -1 && jsonColumnIndex < len(row.Cells) {
			// ถ้ามีคอลัมน์ json ให้ใช้ข้อมูลจากคอลัมน์นั้น
			jsonStr := row.Cells[jsonColumnIndex].String()
			if jsonStr != "" {
				// ตรวจสอบว่าเป็น JSON ที่ถูกต้องหรือไม่
				var jsonObj interface{}
				if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err == nil {
					jsonData = []byte(jsonStr)
				} else {
					fmt.Printf("แปลง JSON สำหรับ PID %s ไม่สำเร็จ: %v\n", pid, err)
					// สร้างข้อมูลจากคอลัมน์อื่นๆ แทน
					rowData := createRowData(row, headers)
					jsonData, _ = json.Marshal(rowData)
				}
			} else {
				// ถ้าคอลัมน์ json ว่างเปล่า ให้สร้างข้อมูลจากคอลัมน์อื่นๆ
				rowData := createRowData(row, headers)
				jsonData, _ = json.Marshal(rowData)
			}
		} else {
			// ถ้าไม่มีคอลัมน์ json ให้สร้างข้อมูลจากคอลัมน์ทั้งหมด
			rowData := createRowData(row, headers)
			jsonData, _ = json.Marshal(rowData)
		}

		// เพิ่มข้อมูลลงในฐานข้อมูล
		_, err = insertStmt.Exec(pid, string(jsonData), serviceDate)
		if err != nil {
			fmt.Printf("เพิ่มข้อมูลสำหรับ PID %s ไม่สำเร็จ: %v\n", pid, err)
			continue
		}
		recordCount++
	}

	// สร้าง index เพื่อให้ค้นหาข้อมูลได้เร็วขึ้น
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_pid ON user_data(pid)")
	if err != nil {
		fmt.Printf("สร้าง index ไม่สำเร็จ: %v\n", err)
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_service_date ON user_data(service_date)")
	if err != nil {
		fmt.Printf("สร้าง index ไม่สำเร็จ: %v\n", err)
	}

	fmt.Printf("แปลงข้อมูลสำเร็จแล้ว นำเข้า %d records ลงในฐานข้อมูล %s\n", recordCount, dbFileName)
}

// createRowData สร้าง map ข้อมูลจากแถวใน Excel
func createRowData(row *xlsx.Row, headers []string) map[string]interface{} {
	rowData := make(map[string]interface{})

	for j, cell := range row.Cells {
		if j < len(headers) && headers[j] != "PID" {
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
