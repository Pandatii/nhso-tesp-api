package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// userData เก็บข้อมูลแบบ hard-coded จากไฟล์ Excel
var userData = map[string]interface{}{
	    "1443852933786": {
	        "age": "32 ปี 5 เดือน 21 วัน",
	        "birthDate": "1992-09-06",
	        "firstName": "ทดสอบ10",
	        "fullName": "ทดสอบ10 ทดสอบ",
	        "hospMain": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospMainOp": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospSub": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "lastName": "ทดสอบ",
	        "mainInscl": "UCS",
	        "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ",
	        "middleName": "",
	        "nationCode": "099",
	        "nationDescription": "ไทย",
	        "paidModel": "1",
	        "personalId": "1443852933786",
	        "provinceCode": "17",
	        "provinceName": "สิงห์บุรี",
	        "serviceDate": "2025-02-02 7:54:00",
	        "serviceHistories": {
	            "claimCode": "PP1321912142",
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว",
	            "serviceCode": "",
	            "serviceName": "",
	            "sourceChannel": "API",
	            "tel": "06384759394"
	        },
	        "sex": "หญิง",
	        "statusAuthen": "TRUE",
	        "statusMessage": "พบข้อมูลการ authen",
	        "subInscl": "",
	        "subInsclName": ""
	    },
	    "1464459697208": {
	        "age": "32 ปี 5 เดือน 29 วัน\r",
	        "birthDate": "1992-08-29",
	        "firstName": "ทดสอบ2",
	        "fullName": "ทดสอบ2 ทดสอบ",
	        "hospMain": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospMainOp": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospSub": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "lastName": "ทดสอบ",
	        "mainInscl": "UCS",
	        "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ",
	        "middleName": "",
	        "nationCode": "099",
	        "nationDescription": "ไทย",
	        "paidModel": "1",
	        "personalId": "1464459697208",
	        "provinceCode": "17",
	        "provinceName": "สิงห์บุรี",
	        "serviceDate": "2025-02-02 7:54:00",
	        "serviceHistories": {
	            "claimCode": "PP1321912134",
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว",
	            "serviceCode": "",
	            "serviceName": "",
	            "sourceChannel": "API",
	            "tel": "06384759386"
	        },
	        "sex": "หญิง",
	        "statusAuthen": "TRUE",
	        "statusMessage": "พบข้อมูลการ authen",
	        "subInscl": "",
	        "subInsclName": ""
	    },
	    "1498157950656": {
	        "age": "32 ปี 0 เดือน 0 วัน",
	        "birthDate": "1992-08-28",
	        "firstName": "ทดสอบ1",
	        "fullName": "ทดสอบ1 ทดสอบ",
	        "hospMain": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospMainOp": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospSub": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "lastName": "ทดสอบ",
	        "mainInscl": "UCS",
	        "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ",
	        "middleName": "",
	        "nationCode": "099",
	        "nationDescription": "ไทย",
	        "paidModel": "1",
	        "personalId": "1498157950656",
	        "provinceCode": "17",
	        "provinceName": "สิงห์บุรี",
	        "serviceDate": "2025-02-02 7:54:00",
	        "serviceHistories": {
	            "claimCode": "PP1321912133",
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว",
	            "serviceCode": "",
	            "serviceName": "",
	            "sourceChannel": "API",
	            "tel": "06384759385"
	        },
	        "sex": "ชาย",
	        "statusAuthen": "TRUE",
	        "statusMessage": "พบข้อมูลการ authen",
	        "subInscl": "",
	        "subInsclName": ""
	    },
	    "2182893981963": {
	        "age": "32 ปี 0 เดือน 0 วัน",
	        "birthDate": "1992-08-28",
	        "firstName": "นิวัด",
	        "fullName": "นิวัด แพน",
	        "hospMain": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospMainOp": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospSub": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "lastName": "แพน",
	        "mainInscl": "UCS",
	        "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ",
	        "middleName": "",
	        "nationCode": "099",
	        "nationDescription": "ไทย",
	        "paidModel": "1",
	        "personalId": "2182893981963",
	        "provinceCode": "17",
	        "provinceName": "สิงห์บุรี",
	        "serviceDate": "2025-02-02 7:54:00",
	        "serviceHistories": {
	            "claimCode": "PP1321912132",
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว",
	            "serviceCode": "",
	            "serviceName": "",
	            "sourceChannel": "API",
	            "tel": "06384759384"
	        },
	        "sex": "หญิง",
	        "statusAuthen": "TRUE",
	        "statusMessage": "พบข้อมูลการ authen",
	        "subInscl": "",
	        "subInsclName": ""
	    },
	    "2860764678719": {
	        "age": "34 ปี 3 เดือน 26 วัน",
	        "birthDate": "1990-01-11",
	        "firstName": "ทดสอบ3",
	        "fullName": "ทดสอบ3 ทดสอบ",
	        "hospMain": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospMainOp": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospSub": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "lastName": "ทดสอบ",
	        "mainInscl": "UCS",
	        "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ",
	        "middleName": "",
	        "nationCode": "099",
	        "nationDescription": "ไทย",
	        "paidModel": "1",
	        "personalId": "2860764678719",
	        "provinceCode": "17",
	        "provinceName": "สิงห์บุรี",
	        "serviceDate": "2025-02-02 7:54:00",
	        "serviceHistories": {
	            "claimCode": "PP1321912135",
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว",
	            "serviceCode": "",
	            "serviceName": "",
	            "sourceChannel": "API",
	            "tel": "06384759387"
	        },
	        "sex": "หญิง",
	        "statusAuthen": "TRUE",
	        "statusMessage": "พบข้อมูลการ authen",
	        "subInscl": "",
	        "subInsclName": ""
	    },
	    "3900618707359": {
	        "age": "31 ปี 1 เดือน 21 วัน",
	        "birthDate": "1994-01-04",
	        "firstName": "ทดสอบ8",
	        "fullName": "ทดสอบ8 ทดสอบ",
	        "hospMain": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospMainOp": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospSub": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "lastName": "ทดสอบ",
	        "mainInscl": "UCS",
	        "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ",
	        "middleName": "",
	        "nationCode": "099",
	        "nationDescription": "ไทย",
	        "paidModel": "1",
	        "personalId": "3900618707359",
	        "provinceCode": "17",
	        "provinceName": "สิงห์บุรี",
	        "serviceDate": "2025-02-02 7:54:00",
	        "serviceHistories": {
	            "claimCode": "PP1321912140",
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว",
	            "serviceCode": "",
	            "serviceName": "",
	            "sourceChannel": "API",
	            "tel": "06384759392"
	        },
	        "sex": "ชาย",
	        "statusAuthen": "TRUE",
	        "statusMessage": "พบข้อมูลการ authen",
	        "subInscl": "",
	        "subInsclName": ""
	    },
	    "4069865132985": {
	        "age": "40 ปี 11 เดือน 6 วัน",
	        "birthDate": "1982-07-22",
	        "firstName": "ทดสอบ6",
	        "fullName": "ทดสอบ6 ทดสอบ",
	        "hospMain": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospMainOp": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospSub": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "lastName": "ทดสอบ",
	        "mainInscl": "UCS",
	        "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ",
	        "middleName": "",
	        "nationCode": "099",
	        "nationDescription": "ไทย",
	        "paidModel": "1",
	        "personalId": "4069865132985",
	        "provinceCode": "17",
	        "provinceName": "สิงห์บุรี",
	        "serviceDate": "2025-02-02 7:54:00",
	        "serviceHistories": {
	            "claimCode": "PP1321912138",
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว",
	            "serviceCode": "",
	            "serviceName": "",
	            "sourceChannel": "API",
	            "tel": "06384759390"
	        },
	        "sex": "ชาย",
	        "statusAuthen": "TRUE",
	        "statusMessage": "พบข้อมูลการ authen",
	        "subInscl": "",
	        "subInsclName": ""
	    },
	    "5468367033903": {
	        "age": "32 ปี 5 เดือน 24 วัน",
	        "birthDate": "1992-09-03",
	        "firstName": "ทดสอบ7",
	        "fullName": "ทดสอบ7 ทดสอบ",
	        "hospMain": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospMainOp": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospSub": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "lastName": "ทดสอบ",
	        "mainInscl": "UCS",
	        "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ",
	        "middleName": "",
	        "nationCode": "099",
	        "nationDescription": "ไทย",
	        "paidModel": "1",
	        "personalId": "5468367033903",
	        "provinceCode": "17",
	        "provinceName": "สิงห์บุรี",
	        "serviceDate": "2025-02-02 7:54:00",
	        "serviceHistories": {
	            "claimCode": "PP1321912139",
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว",
	            "serviceCode": "",
	            "serviceName": "",
	            "sourceChannel": "API",
	            "tel": "06384759391"
	        },
	        "sex": "ชาย",
	        "statusAuthen": "TRUE",
	        "statusMessage": "พบข้อมูลการ authen",
	        "subInscl": "",
	        "subInsclName": ""
	    },
	    "7124746737416": {
	        "age": "27 ปี 9 เดือน 29 วัน",
	        "birthDate": "1997-04-29",
	        "firstName": "ทดสอบ9",
	        "fullName": "ทดสอบ9 ทดสอบ",
	        "hospMain": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospMainOp": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospSub": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "lastName": "ทดสอบ",
	        "mainInscl": "UCS",
	        "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ",
	        "middleName": "",
	        "nationCode": "099",
	        "nationDescription": "ไทย",
	        "paidModel": "1",
	        "personalId": "7124746737416",
	        "provinceCode": "17",
	        "provinceName": "สิงห์บุรี",
	        "serviceDate": "2025-02-02 7:54:00",
	        "serviceHistories": {
	            "claimCode": "PP1321912141",
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว",
	            "serviceCode": "",
	            "serviceName": "",
	            "sourceChannel": "API",
	            "tel": "06384759393"
	        },
	        "sex": "หญิง",
	        "statusAuthen": "TRUE",
	        "statusMessage": "พบข้อมูลการ authen",
	        "subInscl": "",
	        "subInsclName": ""
	    },
	    "9676671309548": {
	        "age": "32 ปี 5 เดือน 26 วัน",
	        "birthDate": "1992-09-01",
	        "firstName": "ทดสอบ5",
	        "fullName": "ทดสอบ5 ทดสอบ",
	        "hospMain": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospMainOp": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospSub": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "lastName": "ทดสอบ",
	        "mainInscl": "UCS",
	        "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ",
	        "middleName": "",
	        "nationCode": "099",
	        "nationDescription": "ไทย",
	        "paidModel": "1",
	        "personalId": "9676671309548",
	        "provinceCode": "17",
	        "provinceName": "สิงห์บุรี",
	        "serviceDate": "2025-02-02 7:54:00",
	        "serviceHistories": {
	            "claimCode": "PP1321912137",
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว",
	            "serviceCode": "",
	            "serviceName": "",
	            "sourceChannel": "API",
	            "tel": "06384759389"
	        },
	        "sex": "หญิง",
	        "statusAuthen": "TRUE",
	        "statusMessage": "พบข้อมูลการ authen",
	        "subInscl": "",
	        "subInsclName": ""
	    },
	    "9976922247707": {
	        "age": "32 ปี 5 เดือน 27 วัน",
	        "birthDate": "1992-08-31",
	        "firstName": "ทดสอบ4",
	        "fullName": "ทดสอบ4 ทดสอบ",
	        "hospMain": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospMainOp": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "hospSub": {
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว"
	        },
	        "lastName": "ทดสอบ",
	        "mainInscl": "UCS",
	        "mainInsclName": "สิทธิหลักประกันสุขภาพแห่งชาติ",
	        "middleName": "",
	        "nationCode": "099",
	        "nationDescription": "ไทย",
	        "paidModel": "1",
	        "personalId": "9976922247707",
	        "provinceCode": "17",
	        "provinceName": "สิงห์บุรี",
	        "serviceDate": "2025-02-02 7:54:00",
	        "serviceHistories": {
	            "claimCode": "PP1321912136",
	            "hcode": "11305",
	            "hname": "รพ.บ้านแพ้ว",
	            "serviceCode": "",
	            "serviceName": "",
	            "sourceChannel": "API",
	            "tel": "06384759388"
	        },
	        "sex": "ชาย",
	        "statusAuthen": "TRUE",
	        "statusMessage": "พบข้อมูลการ authen",
	        "subInscl": "",
	        "subInsclName": ""
	    }
}

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
		w.Write([]byte(`{"status":"error","message":"เกิดข้อผิดพลาดในการสร้าง JSON"}`))
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
