package v1api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockDB struct {
	addCount     int
	readCount    int
	readMultiple bool
}

func (mdb *mockDB) InitData(path string) error {
	return nil
}
func (mdb *mockDB) Read() (*[]string, error) {
	mdb.readCount++
	if mdb.readMultiple {
		return &[]string{"motd", "motd 2"}, nil
	}
	return &[]string{"motd"}, nil
}
func (mdb *mockDB) Add(message string) error {
	mdb.addCount++
	return nil
}

type failDB struct {
	addCount  int
	readCount int
	readEmpty bool
}

func (mdb *failDB) InitData(path string) error {
	return errors.New("something is wrong")
}
func (mdb *failDB) Read() (*[]string, error) {
	mdb.readCount++
	if mdb.readEmpty {
		return new([]string), nil
	}
	return nil, errors.New("something is wrong")
}
func (mdb *failDB) Add(message string) error {
	mdb.addCount++
	return errors.New("something is wrong")
}

func TestMessageHandlerWrongMethod(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/v1/message/", nil)
	rr := httptest.NewRecorder()
	e := env{db: &mockDB{}}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("PUT /v1/message; status == %v, want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestMessageHandlerFailedRead(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/message/", nil)
	rr := httptest.NewRecorder()
	db := failDB{}
	e := env{db: &db}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("GET /v1/message; status == %v, want %v", status, http.StatusInternalServerError)
	}
	if db.readCount != 1 {
		t.Errorf("GET /v1/message, Read called %v times, want %v", db.readCount, 1)
	}
}

func TestMessageHandlerRead(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/message/", nil)
	rr := httptest.NewRecorder()
	db := mockDB{}
	e := env{db: &db}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /v1/message; status == %v, want %v", status, http.StatusOK)
	}
	want := "[\"motd\"]"
	if rr.Body.String() != want {
		t.Errorf("GET /v1/message == %v, want %v", rr.Body.String(), want)
	}
	if db.readCount != 1 {
		t.Errorf("GET /v1/message, Read called %v times, want %v", db.readCount, 1)
	}
}

func TestMessageHandlerReadOne(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/message/0", nil)
	rr := httptest.NewRecorder()
	db := mockDB{}
	e := env{db: &db}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /v1/message/0; status == %v, want %v", status, http.StatusOK)
	}
	want := "\"motd\""
	if rr.Body.String() != want {
		t.Errorf("GET /v1/message/0 == %v, want %v", rr.Body.String(), want)
	}
	if db.readCount != 1 {
		t.Errorf("GET /v1/message/0, Read called %v times, want %v", db.readCount, 1)
	}
}

func TestMessageHandlerReadSecond(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/message/1", nil)
	rr := httptest.NewRecorder()
	db := mockDB{}
	db.readMultiple = true
	e := env{db: &db}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /v1/message/1; status == %v, want %v", status, http.StatusOK)
	}
	want := "\"motd 2\""
	if rr.Body.String() != want {
		t.Errorf("GET /v1/message/1 == %v, want %v", rr.Body.String(), want)
	}
	if db.readCount != 1 {
		t.Errorf("GET /v1/message/1, Read called %v times, want %v", db.readCount, 1)
	}
}

func TestMessageHandlerReadNonNumericID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/message/abc", nil)
	rr := httptest.NewRecorder()
	db := mockDB{}
	e := env{db: &db}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("GET /v1/message/abc; status == %v, want %v", status, http.StatusNotFound)
	}
	if db.readCount != 0 {
		t.Errorf("GET /v1/message/abc, Read called %v times, want %v", db.readCount, 0)
	}
}

func TestMessageHandlerReadNegativeID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/message/-1", nil)
	rr := httptest.NewRecorder()
	db := mockDB{}
	e := env{db: &db}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("GET /v1/message/-1; status == %v, want %v", status, http.StatusNotFound)
	}
	if db.readCount != 0 {
		t.Errorf("GET /v1/message/-1, Read called %v times, want %v", db.readCount, 0)
	}
}

func TestMessageHandlerReadToHighID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/message/2", nil)
	rr := httptest.NewRecorder()
	db := mockDB{}
	e := env{db: &db}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("GET /v1/message/2; status == %v, want %v", status, http.StatusNotFound)
	}
	if db.readCount != 1 {
		t.Errorf("GET /v1/message/2, Read called %v times, want %v", db.readCount, 1)
	}
}

func TestMessageHandlerReadOneFailedRead(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/message/0", nil)
	rr := httptest.NewRecorder()
	db := failDB{}
	e := env{db: &db}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("GET /v1/message/0; status == %v, want %v", status, http.StatusInternalServerError)
	}
	if db.readCount != 1 {
		t.Errorf("GET /v1/message/0, Read called %v times, want %v", db.readCount, 1)
	}
}

func TestMessageHandlerAdd(t *testing.T) {
	req, _ := http.NewRequest("POST", "/v1/message/", strings.NewReader("\"a new motd\""))
	rr := httptest.NewRecorder()
	db := mockDB{}
	e := env{db: &db}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("POST /v1/message; status == %v, want %v", status, http.StatusCreated)
		t.Logf("Body of response: %v", rr.Body.String())
	}
	if db.addCount != 1 {
		t.Errorf("POST /v1/message, Add called %v times, want %v", db.addCount, 1)
	}
}

func TestMessageHandlerFailedAdd(t *testing.T) {
	req, _ := http.NewRequest("POST", "/v1/message/", strings.NewReader("\"a new motd\""))
	rr := httptest.NewRecorder()
	db := failDB{}
	e := env{db: &db}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("POST /v1/message; status == %v, want %v", status, http.StatusInternalServerError)
		t.Logf("Body of response: %v", rr.Body.String())
	}
	if db.addCount != 1 {
		t.Errorf("POST /v1/message, Add called %v times, want %v", db.addCount, 1)
	}
}

func TestMessageHandlerBadInput(t *testing.T) {
	req, _ := http.NewRequest("POST", "/v1/message/", strings.NewReader("a new motd"))
	rr := httptest.NewRecorder()
	e := env{db: &mockDB{}}
	handler := http.HandlerFunc(e.messageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("POST /v1/message; status == %v, want %v", status, http.StatusBadRequest)
		t.Logf("Body of response: %v", rr.Body.String())
	}
}

func TestGetRandomMessageHandlerWrongMethod(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/v1/message/random", nil)
	rr := httptest.NewRecorder()
	e := env{db: &mockDB{}}
	handler := http.HandlerFunc(e.getRandomMessageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("PUT /v1/message; status == %v, want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestGetRandomMessageHandlerRead(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/message/random", nil)
	rr := httptest.NewRecorder()
	db := mockDB{}
	e := env{db: &db}
	handler := http.HandlerFunc(e.getRandomMessageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /v1/message; status == %v, want %v", status, http.StatusOK)
	}
	want := "\"motd\""
	if rr.Body.String() != want {
		t.Errorf("GET /v1/message == %v, want %v", rr.Body.String(), want)
	}
	if db.readCount != 1 {
		t.Errorf("GET /v1/message, Read called %v times, want %v", db.readCount, 1)
	}
}

func TestGetRandomMessageHandlerFailedRead(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/message/random", nil)
	rr := httptest.NewRecorder()
	db := failDB{}
	e := env{db: &db}
	handler := http.HandlerFunc(e.getRandomMessageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("GET /v1/message; status == %v, want %v", status, http.StatusInternalServerError)
	}
	if db.readCount != 1 {
		t.Errorf("GET /v1/message, Read called %v times, want %v", db.readCount, 1)
	}
}

func TestGetRandomMessageHandlerDbEmpty(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/message/random", nil)
	rr := httptest.NewRecorder()
	db := failDB{}
	db.readEmpty = true
	e := env{db: &db}
	handler := http.HandlerFunc(e.getRandomMessageHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("GET /v1/message; status == %v, want %v", status, http.StatusNoContent)
	}
	if db.readCount != 1 {
		t.Errorf("GET /v1/message, Read called %v times, want %v", db.readCount, 1)
	}
}
