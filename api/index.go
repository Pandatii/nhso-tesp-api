package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

// Message เน�เธ�เธฃเธ�เธชเธฃเน�เธฒเธ�เธ�เน�เธญเธกเธนเธฅเธชเธณเธซเธฃเธฑเธ�เธ�เธฒเธฃเธ•เธญเธ�เธ�เธฅเธฑเธ�
type Message struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// ExcelData เน�เธ�เธฃเธ�เธชเธฃเน�เธฒเธ�เธ�เน�เธญเธกเธนเธฅเธชเธณเธซเธฃเธฑเธ�เธ�เน�เธญเธกเธนเธฅเธ—เธตเน�เน�เธ”เน�เธ�เธฒเธ� Excel
type ExcelData struct {
	PID        string            `json:"pid"`
	Fields     map[string]string `json:"fields"`
	JsonString string            `json:"json,omitempty"`
}

// findJSONByPID เธ�เน�เธ�เธซเธฒเธ�เน�เธญเธกเธนเธฅ JSON เธ�เธฒเธ� Excel เน�เธ”เธขเน�เธ�เน� PID
func findJSONByPIDCol3(pid string, excelPath string, serviceDate string) (json.RawMessage, error) {
	// เน€เธ�เธดเธ”เน�เธ�เธฅเน� Excel
	xlFile, err := xlsx.OpenFile(excelPath)
	if err != nil {
		return nil, fmt.Errorf("เน�เธกเน�เธชเธฒเธกเธฒเธฃเธ–เน€เธ�เธดเธ”เน�เธ�เธฅเน� Excel เน�เธ”เน�: %v", err)
	}

	// เธชเธกเธกเธ•เธดเธงเน�เธฒ sheet เน�เธฃเธ�เธ�เธทเธญเธ—เธตเน�เน€เธฃเธฒเธ•เน�เธญเธ�เธ�เธฒเธฃ
	if len(xlFile.Sheets) == 0 {
		return nil, errors.New("เน�เธกเน�เธ�เธ� sheet เน�เธ�เน�เธ�เธฅเน� Excel")
	}
	sheet := xlFile.Sheets[0]

	// เธญเน�เธฒเธ�เธซเธฑเธงเธ�เน�เธญเธ�เธญเธฅเธฑเธกเธ�เน� (เธชเธกเธกเธ•เธดเธงเน�เธฒเธญเธขเธนเน�เน�เธ�เน�เธ–เธงเน�เธฃเธ�)
	if len(sheet.Rows) == 0 {
		return nil, errors.New("Excel เน�เธกเน�เธกเธตเธ�เน�เธญเธกเธนเธฅ")
	}

	jsonColumnIndex := -1        // เธ•เธณเน�เธซเธ�เน�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� json
	serviceDateColumnIndex := -1 // เธ•เธณเน�เธซเธ�เน�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� serviceDate

	// เธ�เน�เธ�เธซเธฒเธ•เธณเน�เธซเธ�เน�เธ�เธ�เธญเธฅเธฑเธกเธ�เน�เธ—เธตเน�เธ•เน�เธญเธ�เธ�เธฒเธฃ
	for i, cell := range sheet.Rows[0].Cells {
		headerText := cell.String()
		if headerText == "json" {
			jsonColumnIndex = i
		} else if headerText == "serviceDate serviceDate" {
			serviceDateColumnIndex = i
		}
	}

	if jsonColumnIndex == -1 {
		return nil, errors.New("เน�เธกเน�เธ�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� 'json' เน�เธ�เน�เธ�เธฅเน� Excel")
	}

	if serviceDateColumnIndex == -1 && serviceDate != "" {
		return nil, errors.New("เน�เธกเน�เธ�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� 'serviceDate' เน�เธ�เน�เธ�เธฅเน� Excel")
	}

	// เธ�เน�เธ�เธซเธฒเน�เธ–เธงเธ—เธตเน�เธกเธตเธ�เน�เธฒ PID เธ•เธฃเธ�เธ�เธฑเธ�
	for i := 1; i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]
		if len(row.Cells) == 0 {
			continue
		}

		// เธ•เธฃเธงเธ�เธชเธญเธ�เธงเน�เธฒ PID เธ•เธฃเธ�เธ�เธฑเธ�เธ—เธตเน�เธ•เน�เธญเธ�เธ�เธฒเธฃเธซเธฃเธทเธญเน�เธกเน�
		rowPID := row.Cells[2].String()
		if rowPID == pid {
			// เธ–เน�เธฒเธกเธตเธ�เธฒเธฃเธฃเธฐเธ�เธธ serviceDate เน�เธซเน�เธ•เธฃเธงเธ�เธชเธญเธ�เธงเน�เธฒเธ•เธฃเธ�เธ�เธฑเธ�เธซเธฃเธทเธญเน�เธกเน�
			if serviceDate != "" && serviceDateColumnIndex < len(row.Cells) {
				rowServiceDate := row.Cells[serviceDateColumnIndex].String()
				// เธ•เธฑเธ”เน€เธญเธฒเน€เธ�เธ�เธฒเธฐ 10 เธ•เธฑเธงเน�เธฃเธ� (YYYY-MM-DD)
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

				// เธ–เน�เธฒ serviceDate เน�เธกเน�เธ•เธฃเธ�เธ�เธฑเธ� เน�เธซเน�เธ�เน�เธฒเธกเน�เธ–เธงเธ�เธตเน�
				if rowServiceDatePrefix != serviceDatePrefix {
					continue
				}
			}

			// เธ•เธฃเธงเธ�เธชเธญเธ�เธงเน�เธฒเธกเธตเธ�เน�เธญเธกเธนเธฅเน�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� json เธซเธฃเธทเธญเน�เธกเน�
			if jsonColumnIndex < len(row.Cells) {
				jsonStr := row.Cells[jsonColumnIndex].String()
				if jsonStr == "" {
					return json.RawMessage("{}"), nil
				}

				// เธ•เธฃเธงเธ�เธชเธญเธ�เธงเน�เธฒเธ�เน�เธญเธกเธนเธฅเน€เธ�เน�เธ� JSON เธ—เธตเน�เธ–เธนเธ�เธ•เน�เธญเธ�เธซเธฃเธทเธญเน�เธกเน�
				var jsonData json.RawMessage
				if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
					return nil, fmt.Errorf("เธ�เน�เธญเธกเธนเธฅเน�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� 'json' เน�เธกเน�เน�เธ�เน� JSON เธ—เธตเน�เธ–เธนเธ�เธ•เน�เธญเธ�: %v", err)
				}

				return jsonData, nil
			}
			return nil, errors.New("เน�เธกเน�เธ�เธ�เธ�เน�เธญเธกเธนเธฅเน�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� 'json'")
		}
	}

	// เธ–เน�เธฒ serviceDate เธ–เธนเธ�เธฃเธฐเธ�เธธเน�เธ•เน�เน�เธกเน�เธ�เธ�เธ�เน�เธญเธกเธนเธฅเธ—เธตเน�เธ•เธฃเธ�เธ�เธฑเธ�
	if serviceDate != "" {
		return nil, fmt.Errorf("เน�เธกเน�เธ�เธ�เธ�เน�เธญเธกเธนเธฅเธชเธณเธซเธฃเธฑเธ� PID: %s เน�เธฅเธฐ serviceDate: %s", pid, serviceDate)
	}

	return nil, fmt.Errorf("เน�เธกเน�เธ�เธ�เธ�เน�เธญเธกเธนเธฅเธชเธณเธซเธฃเธฑเธ� PID: %s", pid)
}

// findJSONByPID เธ�เน�เธ�เธซเธฒเธ�เน�เธญเธกเธนเธฅ JSON เธ�เธฒเธ� Excel เน�เธ”เธขเน�เธ�เน� PID
func findJSONByPID(pid string, excelPath string) (json.RawMessage, error) {
	// เน€เธ�เธดเธ”เน�เธ�เธฅเน� Excel
	xlFile, err := xlsx.OpenFile(excelPath)
	if err != nil {
		return nil, fmt.Errorf("เน�เธกเน�เธชเธฒเธกเธฒเธฃเธ–เน€เธ�เธดเธ”เน�เธ�เธฅเน� Excel เน�เธ”เน�: %v", err)
	}

	// เธชเธกเธกเธ•เธดเธงเน�เธฒ sheet เน�เธฃเธ�เธ�เธทเธญเธ—เธตเน�เน€เธฃเธฒเธ•เน�เธญเธ�เธ�เธฒเธฃ
	if len(xlFile.Sheets) == 0 {
		return nil, errors.New("เน�เธกเน�เธ�เธ� sheet เน�เธ�เน�เธ�เธฅเน� Excel")
	}
	sheet := xlFile.Sheets[0]

	// เธญเน�เธฒเธ�เธซเธฑเธงเธ�เน�เธญเธ�เธญเธฅเธฑเธกเธ�เน� (เธชเธกเธกเธ•เธดเธงเน�เธฒเธญเธขเธนเน�เน�เธ�เน�เธ–เธงเน�เธฃเธ�)
	if len(sheet.Rows) == 0 {
		return nil, errors.New("Excel เน�เธกเน�เธกเธตเธ�เน�เธญเธกเธนเธฅ")
	}

	jsonColumnIndex := -1 // เธ•เธณเน�เธซเธ�เน�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� json

	// เธ�เน�เธ�เธซเธฒเธ•เธณเน�เธซเธ�เน�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� json
	for i, cell := range sheet.Rows[0].Cells {
		headerText := cell.String()
		if headerText == "json" {
			jsonColumnIndex = i
			break
		}
	}

	if jsonColumnIndex == -1 {
		return nil, errors.New("เน�เธกเน�เธ�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� 'json' เน�เธ�เน�เธ�เธฅเน� Excel")
	}

	// เธ�เน�เธ�เธซเธฒเน�เธ–เธงเธ—เธตเน�เธกเธตเธ�เน�เธฒ PID เธ•เธฃเธ�เธ�เธฑเธ�
	for i := 1; i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]
		if len(row.Cells) == 0 {
			continue
		}

		// เธ•เธฃเธงเธ�เธชเธญเธ�เธงเน�เธฒ PID เธ•เธฃเธ�เธ�เธฑเธ�เธ—เธตเน�เธ•เน�เธญเธ�เธ�เธฒเธฃเธซเธฃเธทเธญเน�เธกเน�
		rowPID := row.Cells[0].String()
		if rowPID == pid {
			// เธ•เธฃเธงเธ�เธชเธญเธ�เธงเน�เธฒเธกเธตเธ�เน�เธญเธกเธนเธฅเน�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� json เธซเธฃเธทเธญเน�เธกเน�
			if jsonColumnIndex < len(row.Cells) {
				jsonStr := row.Cells[jsonColumnIndex].String()
				if jsonStr == "" {
					return json.RawMessage("{}"), nil
				}

				// เธ•เธฃเธงเธ�เธชเธญเธ�เธงเน�เธฒเธ�เน�เธญเธกเธนเธฅเน€เธ�เน�เธ� JSON เธ—เธตเน�เธ–เธนเธ�เธ•เน�เธญเธ�เธซเธฃเธทเธญเน�เธกเน�
				var jsonData json.RawMessage
				if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
					return nil, fmt.Errorf("เธ�เน�เธญเธกเธนเธฅเน�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� 'json' เน�เธกเน�เน�เธ�เน� JSON เธ—เธตเน�เธ–เธนเธ�เธ•เน�เธญเธ�: %v", err)
				}

				return jsonData, nil
			}
			return nil, errors.New("เน�เธกเน�เธ�เธ�เธ�เน�เธญเธกเธนเธฅเน�เธ�เธ�เธญเธฅเธฑเธกเธ�เน� 'json'")
		}
	}

	return nil, fmt.Errorf("เน�เธกเน�เธ�เธ�เธ�เน�เธญเธกเธนเธฅเธชเธณเธซเธฃเธฑเธ� PID: %s", pid)
}

// handleAPI เธ�เธฑเธ”เธ�เธฒเธฃ endpoint เธซเธฅเธฑเธ�เธ�เธญเธ� API
func handleAPIAuthen(w http.ResponseWriter, r *http.Request) {
	// เธ�เธณเธซเธ�เธ” header
	w.Header().Set("Content-Type", "application/json")

	// เธฃเธฑเธ�เธ�เน�เธฒ parameter เธ�เธฒเธ� URL
	query := r.URL.Query()

	// เธ•เธฃเธงเธ�เธชเธญเธ�เธ�เน�เธฒ PID
	pid := query.Get("personalId")
	serviceDate := query.Get("serviceDate")

	// เธ–เน�เธฒเธกเธตเธ�เน�เธฒ PID เน�เธซเน�เธ�เน�เธ�เธซเธฒเธ�เน�เธญเธกเธนเธฅเน�เธ� Excel
	if pid != "" {
		// เธ•เธณเน�เธซเธ�เน�เธ�เน�เธ�เธฅเน� Excel (เธญเธฒเธ�เธ�เธฐเธ•เน�เธญเธ�เธ�เธฃเธฑเธ�เน€เธ�เธฅเธตเน�เธขเธ�เธ•เธฒเธกเธ�เธฒเธฃเน€เธ�เน�เธ�เน�เธ�เธฅเน�เน�เธ� Vercel)
		excelPath := "../data/Mockup API Authen&realPerson_Authen.xlsx"

		data, err := findJSONByPIDCol3(pid, excelPath, serviceDate)
		if err == nil {
			// เธชเน�เธ�เธ�เน�เธญเธกเธนเธฅเธ—เธตเน�เธ�เธ�เธ�เธฅเธฑเธ�เน�เธ�เน�เธ”เธขเธ•เธฃเธ�
			json.NewEncoder(w).Encode(data)
			return
		} else {
			// เธ–เน�เธฒเน�เธกเน�เธ�เธ�เธ�เน�เธญเธกเธนเธฅเธซเธฃเธทเธญเธกเธตเธ�เน�เธญเธ�เธดเธ”เธ�เธฅเธฒเธ”
			errorResponse := map[string]string{
				"error": fmt.Sprintf("เน�เธกเน�เธ�เธ�เธ�เน�เธญเธกเธนเธฅเธชเธณเธซเธฃเธฑเธ� PID: %s", pid),
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	// เธ–เน�เธฒเน�เธกเน�เน�เธ”เน�เธฃเธฐเธ�เธธ PID เน�เธซเน�เธชเน�เธ�เธ�เน�เธญเธ�เธงเธฒเธกเน�เธ�เน�เธ�เน€เธ•เธทเธญเธ�
	errorResponse := map[string]string{
		"error": "เธ�เธฃเธธเธ“เธฒเธฃเธฐเธ�เธธเธ�เน�เธฒ PID (เน€เธ�เน�เธ�: /api?pid=P001)",
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errorResponse)
}

// handleAPI เธ�เธฑเธ”เธ�เธฒเธฃ endpoint เธซเธฅเธฑเธ�เธ�เธญเธ� API
func handleAPI(w http.ResponseWriter, r *http.Request) {
	// เธ�เธณเธซเธ�เธ” header
	w.Header().Set("Content-Type", "application/json")

	// เธฃเธฑเธ�เธ�เน�เธฒ parameter เธ�เธฒเธ� URL
	query := r.URL.Query()

	// เธ•เธฃเธงเธ�เธชเธญเธ�เธ�เน�เธฒ PID
	pid := query.Get("PID")

	// เธ–เน�เธฒเธกเธตเธ�เน�เธฒ PID เน�เธซเน�เธ�เน�เธ�เธซเธฒเธ�เน�เธญเธกเธนเธฅเน�เธ� Excel
	if pid != "" {
		// เธ•เธณเน�เธซเธ�เน�เธ�เน�เธ�เธฅเน� Excel (เธญเธฒเธ�เธ�เธฐเธ•เน�เธญเธ�เธ�เธฃเธฑเธ�เน€เธ�เธฅเธตเน�เธขเธ�เธ•เธฒเธกเธ�เธฒเธฃเน€เธ�เน�เธ�เน�เธ�เธฅเน�เน�เธ� Vercel)
		excelPath := "../data/Mockup API Authen&realPerson_RealPerson.xlsx"

		data, err := findJSONByPID(pid, excelPath)
		if err == nil {
			// เธชเน�เธ�เธ�เน�เธญเธกเธนเธฅเธ—เธตเน�เธ�เธ�เธ�เธฅเธฑเธ�เน�เธ�เน�เธ”เธขเธ•เธฃเธ�
			json.NewEncoder(w).Encode(data)
			return
		} else {
			// เธ–เน�เธฒเน�เธกเน�เธ�เธ�เธ�เน�เธญเธกเธนเธฅเธซเธฃเธทเธญเธกเธตเธ�เน�เธญเธ�เธดเธ”เธ�เธฅเธฒเธ”
			errorResponse := map[string]string{
				"error": fmt.Sprintf("เน�เธกเน�เธ�เธ�เธ�เน�เธญเธกเธนเธฅเธชเธณเธซเธฃเธฑเธ� PID: %s", pid),
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	// เธ–เน�เธฒเน�เธกเน�เน�เธ”เน�เธฃเธฐเธ�เธธ PID เน�เธซเน�เธชเน�เธ�เธ�เน�เธญเธ�เธงเธฒเธกเน�เธ�เน�เธ�เน€เธ•เธทเธญเธ�
	errorResponse := map[string]string{
		"error": "เธ�เธฃเธธเธ“เธฒเธฃเธฐเธ�เธธเธ�เน�เธฒ PID (เน€เธ�เน�เธ�: /api?pid=P001)",
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errorResponse)
}

// handleRoot เธ�เธฑเธ”เธ�เธฒเธฃเน€เธชเน�เธ�เธ—เธฒเธ�เธซเธฅเธฑเธ�
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	msg := Message{
		Status:    "success",
		Message:   "เธขเธดเธ�เธ”เธตเธ•เน�เธญเธ�เธฃเธฑเธ�เธชเธนเน� API เธฅเธญเธ�เน�เธ�เน� endpoint /api เน€เธ�เธทเน�เธญเธ”เธนเธ�เน�เธญเธกเธนเธฅ",
		Timestamp: time.Now(),
	}

	json.NewEncoder(w).Encode(msg)
}

// Handler เน€เธ�เน�เธ�เธ�เธฑเธ�เธ�เน�เธ�เธฑเธ�เธซเธฅเธฑเธ�เธ—เธตเน�เธ�เธฐเธ–เธนเธ�เน€เธฃเธตเธขเธ�เน�เธ”เธข Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// เธ�เธณเธซเธ�เธ” CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// เธ�เธฑเธ”เธ�เธฒเธฃเธ�เธฑเธ� OPTIONS request (เธชเธณเธซเธฃเธฑเธ� CORS preflight)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// เธ•เธฃเธงเธ�เธชเธญเธ�เน€เธชเน�เธ�เธ—เธฒเธ�
	if r.URL.Path == "/" {
		// เธซเธ�เน�เธฒเธซเธฅเธฑเธ�
		homeHandler(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api") {
		// API endpoint
		handleAPI(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/authencodeapi") {
		// API endpoint
		handleAPIAuthen(w, r)
		return
	}

	// เธ�เธฃเธ“เธตเน�เธกเน�เธ�เธ�เน€เธชเน�เธ�เธ—เธฒเธ�
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "เน€เธชเน�เธ�เธ—เธฒเธ�เน�เธกเน�เธ–เธนเธ�เธ•เน�เธญเธ�",
	})
}

// homeHandler เธ�เธฑเธ”เธ�เธฒเธฃเน€เธชเน�เธ�เธ—เธฒเธ�เธซเธฅเธฑเธ�
func homeHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "เธขเธดเธ�เธ”เธตเธ•เน�เธญเธ�เธฃเธฑเธ�เธชเธนเน� API เธฅเธญเธ�เน�เธ�เน� endpoint /api?pid=P001 เธซเธฃเธทเธญ /api?pid=P001&serviceDate=2025-02-15 เน€เธ�เธทเน�เธญเธ”เธนเธ�เน�เธญเธกเธนเธฅ JSON",
	})
}
