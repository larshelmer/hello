package v1api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessageHandlerWrongMethod(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/v1/message", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("PUT /v1/message; status == %v, want %v", status, http.StatusMethodNotAllowed)
	}
}
