package handler

import (
	"net/http"
	"testing"
)

func TestGetTeacherDetailHandler(t *testing.T) {
	w, r := baseConfig()

	r.GET("/api/teacher/:teacherID", GetTeacherDetailHandler)
	req, _ := http.NewRequest("GET", "/api/teacher/1", nil)
	r.ServeHTTP(w, req)
	// log

	prettyJsonLog(w)
}

func TestGetTeacherListHandler(t *testing.T) {
	w, r := baseConfig()

	r.GET("/api/teacher", GetTeacherListHandler)
	req, _ := http.NewRequest("GET", "/api/teacher", nil)
	r.ServeHTTP(w, req)
	// log

	prettyJsonLog(w)
}

func TestSearchTeacherListHandler(t *testing.T) {
	querys := map[string]struct {
		q      string
		status int
	}{
		"valid-multi-req":   {"name=古金宇&page=1&page_size=3", http.StatusOK},
		"valid-single-req":  {"name=古金宇&department=电子信息与电气工程学院&page=1&page_size=3", http.StatusOK},
		"invalid-value-req": {"name=你谁", http.StatusOK},
		"invalid-page-req":  {"page=-1&page_size=10086", http.StatusOK},
		"overflow-page-req": {"page=10086&page_size=1", http.StatusOK},
		"invalid-key-req":   {"hello=1", http.StatusBadRequest},
	}

	for name, s := range querys {
		t.Run(name, func(t *testing.T) {
			w, r := baseConfig()
			r.GET("/api/teacher/query", SearchTeacherListHandler)
			req, _ := http.NewRequest("GET", "/api/teacher/query", nil)
			req.URL.RawQuery = s.q
			r.ServeHTTP(w, req)
			// log
			if w.Code != s.status {
				t.Errorf("Expected status code %d, but got %d", s.status, w.Code)
			}
			prettyJsonLog(w)
		})
	}
}
