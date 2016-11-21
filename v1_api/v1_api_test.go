package v1api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDB struct{}

func (mdb *mockDB) InitData(path string) error {
	return nil
}
func (mdb *mockDB) Read() (*[]string, error) {
	return nil, nil
}
func (mdb *mockDB) Add(message string) error {
	return nil
}

func TestMessageHandlerWrongMethod(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/v1/message", nil)
	rr := httptest.NewRecorder()
	e := env{db: &mockDB{}}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("PUT /v1/message; status == %v, want %v", status, http.StatusMethodNotAllowed)
	}
}
