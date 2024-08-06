package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jcourse_go/dal"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func printTables(db *gorm.DB) {
	var tables []string
	derr := db.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tables).Error
	if derr != nil {
		panic("failed to get tables")
	}

	// 打印所有表名
	fmt.Println("Tables:")
	for _, table := range tables {
		fmt.Println(table)
	}
}
func prettyJsonLog(w *httptest.ResponseRecorder) {
	data := w.Body.Bytes()
	// 使用 json.Indent 函数格式化 JSON 字符串
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, data, "", "    ")

	// 打印格式化的 JSON 字符串
	prettyJSON.WriteTo(os.Stdout)
}
func baseConfig() (*httptest.ResponseRecorder, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	err := dal.InitSqliteDBTest("../gorm.db")
	if err != nil {
		panic(err)
	}
	w := httptest.NewRecorder()
	r := gin.Default()
	return w, r
}
func TestGetTrainingPlanHandler(t *testing.T) {
	w, r := baseConfig()

	r.GET("/api/trainingplan/:trainingPlanID", GetTrainingPlanHandler)
	req, _ := http.NewRequest("GET", "/api/trainingplan/1", nil)
	r.ServeHTTP(w, req)
	// log

	prettyJsonLog(w)
}

func TestSearchTrainingPlanHandler(t *testing.T) {
	querys := map[string]struct {
		q      string
		status int
	}{
		"valid-multi-req":   {"entry_year=2019&page=1&page_size=3", http.StatusOK},
		"valid-single-req":  {"entry_year=2019&major_name=数学与应用数学&page=1&page_size=3", http.StatusOK},
		"invalid-value-req": {"entry_year=不存在", http.StatusBadRequest},
		"invalid-page-req":  {"page=-1&page_size=10086", http.StatusOK},
		"overflow-page-req": {"page=10086&page_size=1", http.StatusOK},
		"invalid-key-req":   {"hello=1", http.StatusBadRequest},
	}

	for name, s := range querys {
		t.Run(name, func(t *testing.T) {
			w, r := baseConfig()
			r.GET("/api/trainingplan/query", SearchTrainingPlanHandler)
			req, _ := http.NewRequest("GET", "/api/trainingplan/query", nil)
			req.URL.RawQuery = s.q
			r.ServeHTTP(w, req)
			// log
			if s.status != w.Code {
				t.Errorf("Expected status code %d, but got %d", s.status, w.Code)
			}
			prettyJsonLog(w)
		})
	}
}
