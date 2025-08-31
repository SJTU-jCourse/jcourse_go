package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"jcourse_go/internal/interface/dto"
)

func TestOptCourseReviewHandler(t *testing.T) {
	t.Run("BasicBindInBody", func(t *testing.T) {
		body := dto.OptCourseReviewRequest{
			CourseName:    "CourseName",
			ReviewContent: "ReviewContent",
		}
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		bodyBytes, err := json.Marshal(body)
		assert.Nil(t, err)
		req := httptest.NewRequest("POST", "/review/opt", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.Handle("POST", "/review/opt", OptCourseReviewHandler)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusBadRequest, w.Code)
	})

	t.Run("SuccessWriteResponse", func(t *testing.T) {
		body := dto.OptCourseReviewRequest{
			CourseName:    "CourseName",
			ReviewContent: "ReviewContent",
		}
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		bodyBytes, err := json.Marshal(body)
		assert.Nil(t, err)
		req := httptest.NewRequest("POST", "/review/opt", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler := func(c *gin.Context) {
			responseJson := `{"suggestion":"suggestion","result":"result"}`
			var response dto.OptCourseReviewResponse
			err := json.Unmarshal([]byte(responseJson), &response)
			if err != nil {
				c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
				return
			}
			c.JSON(http.StatusOK, response)
		}
		r.Handle("POST", "/review/opt", handler)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("response: %v", w.Body.String())
	})
}

func TestGetMatchCoursesHandler(t *testing.T) {
	t.Run("BasicBindInBody", func(t *testing.T) {
		body := dto.GetMatchCourseRequest{
			Description: "test",
		}
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		bodyBytes, err := json.Marshal(body)
		assert.Nil(t, err)
		req := httptest.NewRequest("POST", "/course/match", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req2 := httptest.NewRequest("POST", "/course/match?description=test", nil)

		r.Handle("POST", "/course/match", GetMatchCoursesHandler)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusBadRequest, w.Code)

		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusBadRequest, w2.Code)
	})
}
