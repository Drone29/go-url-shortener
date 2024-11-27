package http_handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/db_interface"
)

// mock interface

type dbCollectionMock struct {
	data   []URLData
	id_cnt int
}

var mock_db = dbCollectionMock{}

func (collection *dbCollectionMock) InsertOne(doc any) (id string, err error) {
	t, ok := doc.(URLData)
	if ok {
		t.ID = fmt.Sprintf("%d", collection.id_cnt)
		collection.data = append(collection.data, t)
		collection.id_cnt++
		return t.ID, nil
	}
	return "", fmt.Errorf("invalid doc type %T", t)
}

func (collection *dbCollectionMock) FindOne(filter any, result any) error {
	f, ok := filter.(URLData)
	if !ok {
		return fmt.Errorf("invalid filter type %T", f)
	}
	r, ok := result.(*URLData)
	if !ok {
		return fmt.Errorf("invalid result type %T", r)
	}
	for _, data := range collection.data {
		if f.URL == data.URL || f.ShortCode == data.ShortCode {
			*r = data
			return nil
		}
	}
	return db_interface.ErrNoDocuments
}

// update doc
func (collection *dbCollectionMock) UpdateOne(filter any, update_with any) error {
	f, ok := filter.(URLData)
	if !ok {
		return fmt.Errorf("invalid filter type %T", f)
	}
	r, ok := update_with.(URLData)
	if !ok {
		return fmt.Errorf("invalid result type %T", r)
	}
	for i, data := range collection.data {
		if f.URL == data.URL || f.ShortCode == data.ShortCode {
			collection.data[i] = r
			return nil
		}
	}
	return db_interface.ErrNoDocuments
}

// delete doc
func (collection *dbCollectionMock) DeleteOne(filter any) error {
	f, ok := filter.(URLData)
	if !ok {
		return fmt.Errorf("invalid filter type %T", f)
	}
	for i, data := range collection.data {
		if f.URL == data.URL || f.ShortCode == data.ShortCode {
			temp := collection.data[i+1:]                      // save everything after i
			collection.data = collection.data[:i]              // truncate until i
			collection.data = append(collection.data, temp...) // concatenate
			return nil
		}
	}
	return db_interface.ErrNoDocuments
}

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
		ID:        "1",
		URL:       "http://someurl.com",
		ShortCode: "abc123",
	})

	w := testHTTP("GET", "/shorten/abc123", "")
	if w.Code != http.StatusOK {
		t.Errorf("invalid response code %v", w.Code)
	}
	if _, err := testResult(w, mock_db.data[0]); err != nil {
		t.Errorf("%v", err)
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
		t.Errorf("invalid response: %v", url_data)
	}
}
