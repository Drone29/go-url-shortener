package http_handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var mock_db = dbCollectionMock{}

// helpers

func testHTTP(method, url, body string) *httptest.ResponseRecorder {
	db = &mock_db

	w := httptest.NewRecorder()
	// mock request
	req := httptest.NewRequest(method, url, strings.NewReader(body))

	shorten(w, req)

	return w
}

func testResult(w *httptest.ResponseRecorder, ref URLData) (*URLData, error) {
	body, err := io.ReadAll(w.Body)
	if err != nil {
		return nil, fmt.Errorf("io error %v", err)
	}
	url_data := URLData{}
	if err := json.Unmarshal(body, &url_data); err != nil {
		return nil, fmt.Errorf("json error %v", err)
	}
	if url_data.String() != ref.String() {
		return nil, fmt.Errorf("invalid response: %v", url_data)
	}
	return &url_data, nil
}

// tests

// POST

func TestPOSTInvalidURL(t *testing.T) {
	if w := testHTTP("POST", "/shorten/abc", ""); w.Code != http.StatusNotFound {
		t.Errorf("invalid response code %v", w.Code)
	}
}

func TestPOSTEmptyBody(t *testing.T) {
	if w := testHTTP("POST", "/shorten", ""); w.Code != http.StatusBadRequest {
		t.Errorf("invalid response code %v", w.Code)
	}
}

func TestPOSTInvalidBody(t *testing.T) {
	//invalid url type
	if w := testHTTP("POST", "/shorten", `{"url": 123}`); w.Code != http.StatusBadRequest {
		t.Errorf("invalid response code %v", w.Code)
	}
	// incomplete json
	if w := testHTTP("POST", "/shorten", `{"url": "http://someurl"`); w.Code != http.StatusBadRequest {
		t.Errorf("invalid response code %v", w.Code)
	}
	//absent url
	if w := testHTTP("POST", "/shorten", `{"notaurl": "http://someurl"}`); w.Code != http.StatusBadRequest {
		t.Errorf("invalid response code %v", w.Code)
	}
}

func TestPOSTInsertion(t *testing.T) {

	w := testHTTP("POST", "/shorten", `{"url": "http://someurl"}`)
	if w.Code != http.StatusCreated {
		t.Errorf("invalid response code %v", w.Code)
	}
	if len(mock_db.data) != 1 {
		t.Errorf("nothing was inserted into db")
	}
	if _, err := testResult(w, mock_db.data[0]); err != nil {
		t.Errorf("%v", err)
	}

	// check if requested again
	w = testHTTP("POST", "/shorten", `{"url": "http://someurl"}`)
	if w.Code != http.StatusOK {
		t.Errorf("invalid response code %v", w.Code)
	}
	if len(mock_db.data) != 1 {
		t.Errorf("shouldn't have inserted into db")
	}
	if _, err := testResult(w, mock_db.data[0]); err != nil {
		t.Errorf("%v", err)
	}
}

// GET
func TestGETInvalidURL(t *testing.T) {
	if w := testHTTP("GET", "/shorten/", ""); w.Code != http.StatusNotFound {
		t.Errorf("invalid response code %v", w.Code)
	}
}

func TestGETNoData(t *testing.T) {
	mock_db.data = mock_db.data[:0] //clear data
	if w := testHTTP("GET", "/shorten/abc123", ""); w.Code != http.StatusNotFound {
		t.Errorf("invalid response code %v", w.Code)
	}
}

func TestGETRetrieve(t *testing.T) {
	mock_db.data = mock_db.data[:0] //clear data
	// add record to db
	mock_db.data = append(mock_db.data, URLData{
		ID:          "1",
		URL:         "http://someurl.com",
		ShortCode:   "abc123",
		AccessCount: 3,
	})

	w := testHTTP("GET", "/shorten/abc123", "")
	if w.Code != http.StatusOK {
		t.Errorf("invalid response code %v", w.Code)
	}
	url_data, err := testResult(w, mock_db.data[0])
	if err != nil {
		t.Errorf("%v", err)
	}
	if mock_db.data[0].AccessCount != 4 {
		t.Error("should increment access counter")
	}
	if url_data.AccessCount != 0 {
		t.Error("shouldn't return access count for simple get")
	}
}

func TestGETStats(t *testing.T) {
	mock_db.data = mock_db.data[:0] //clear data
	// add record to db
	mock_db.data = append(mock_db.data, URLData{
		ID:          "1",
		URL:         "http://someurl.com",
		ShortCode:   "abc123",
		AccessCount: 3,
	})

	w := testHTTP("GET", "/shorten/abc123/stats", "")
	if w.Code != http.StatusOK {
		t.Errorf("invalid response code %v", w.Code)
	}
	url_data, err := testResult(w, mock_db.data[0])
	if err != nil {
		t.Errorf("%v", err)
	}
	if url_data.AccessCount != 3 {
		t.Errorf("should return access count for stats")
	}
}

// PUT
func TestPUTInvalidURL(t *testing.T) {
	if w := testHTTP("PUT", "/shorten/", ""); w.Code != http.StatusNotFound {
		t.Errorf("invalid response code %v", w.Code)
	}
}

func TestPUTInvalidBody(t *testing.T) {
	//invalid url type
	if w := testHTTP("PUT", "/shorten/abc123", `{"url": 123}`); w.Code != http.StatusBadRequest {
		t.Errorf("invalid response code %v", w.Code)
	}
	// incomplete json
	if w := testHTTP("PUT", "/shorten/abc123", `{"url": "http://someurl"`); w.Code != http.StatusBadRequest {
		t.Errorf("invalid response code %v", w.Code)
	}
	//absent url
	if w := testHTTP("PUT", "/shorten/abc123", `{"notaurl": "http://someurl"}`); w.Code != http.StatusBadRequest {
		t.Errorf("invalid response code %v", w.Code)
	}
}

func TestPUTEmptyBody(t *testing.T) {
	mock_db.data = mock_db.data[:0] //clear data
	if w := testHTTP("PUT", "/shorten/abc123", ""); w.Code != http.StatusBadRequest {
		t.Errorf("invalid response code %v", w.Code)
	}
}

func TestPUTNoData(t *testing.T) {
	mock_db.data = mock_db.data[:0] //clear data
	if w := testHTTP("PUT", "/shorten/abc123", `{"url": "http://somenewurl"}`); w.Code != http.StatusNotFound {
		t.Errorf("invalid response code %v", w.Code)
	}
}

func TestPUTChangeData(t *testing.T) {
	mock_db.data = mock_db.data[:0] //clear data
	// add record to db
	mock_db.data = append(mock_db.data, URLData{
		ID:          "1",
		URL:         "http://someurl.com",
		ShortCode:   "abc123",
		AccessCount: 3,
	})

	w := testHTTP("PUT", "/shorten/abc123", `{"url": "http://somenewurl"}`)
	if w.Code != http.StatusOK {
		t.Errorf("invalid response code %v", w.Code)
	}
	url_data, err := testResult(w, mock_db.data[0])
	if err != nil {
		t.Errorf("%v", err)
	}
	if mock_db.data[0].URL != "http://somenewurl" {
		t.Errorf("data didn't change")
	}
	if mock_db.data[0].AccessCount != 3 {
		t.Error("shouldn't increment access counter")
	}
	if url_data.AccessCount != 0 {
		t.Error("shouln't return access counter")
	}
}

// DELETE
func TestDELETEInvalidURL(t *testing.T) {
	if w := testHTTP("DELETE", "/shorten/", ""); w.Code != http.StatusNotFound {
		t.Errorf("invalid response code %v", w.Code)
	}
}

func TestDELETENoData(t *testing.T) {
	mock_db.data = mock_db.data[:0] //clear data
	if w := testHTTP("DELETE", "/shorten/abc123", ""); w.Code != http.StatusNotFound {
		t.Errorf("invalid response code %v", w.Code)
	}
}

func TestDELETERemoveRecord(t *testing.T) {
	mock_db.data = mock_db.data[:0] //clear data
	// add record to db
	mock_db.data = append(mock_db.data, URLData{
		ID:          "1",
		URL:         "http://someurl.com",
		ShortCode:   "abc123",
		AccessCount: 3,
	})

	w := testHTTP("DELETE", "/shorten/abc123", "")
	if w.Code != http.StatusNoContent {
		t.Errorf("invalid response code %v", w.Code)
	}
	if len(mock_db.data) > 0 {
		t.Errorf("no records were deleted")
	}
}
