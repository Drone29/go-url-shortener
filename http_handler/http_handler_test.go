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

// helpers

func testHTTP(method, url, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	// mock request
	req := httptest.NewRequest(method, url, strings.NewReader(body))

	shorten(w, req)

	return w
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
	mock_db := dbCollectionMock{}
	db = &mock_db

	w := testHTTP("POST", "/shorten", `{"url": "http://someurl"}`)
	if w.Code != http.StatusCreated {
		t.Errorf("invalid response code %v", w.Code)
	}
	if len(mock_db.data) != 1 {
		t.Errorf("nothing was inserted into db")
	}
	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Errorf("io error %v", err)
	}
	url_data := URLData{}
	if err := json.Unmarshal(body, &url_data); err != nil {
		t.Errorf("json error %v", err)
	}
	if url_data.String() != mock_db.data[0].String() {
		t.Errorf("invalid response: %v", url_data)
	}

	// check if requested again
	w = testHTTP("POST", "/shorten", `{"url": "http://someurl"}`)
	if w.Code != http.StatusOK {
		t.Errorf("invalid response code %v", w.Code)
	}
	if len(mock_db.data) != 1 {
		t.Errorf("shouldn't have inserted into db")
	}
	body, err = io.ReadAll(w.Body)
	if err != nil {
		t.Errorf("io error %v", err)
	}
	if err := json.Unmarshal(body, &url_data); err != nil {
		t.Errorf("json error %v", err)
	}
	if url_data.String() != mock_db.data[0].String() {
		t.Errorf("invalid response: %v", url_data)
	}
}
