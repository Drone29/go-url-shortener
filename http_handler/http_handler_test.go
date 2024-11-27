package http_handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testHTTP(method string, url string, handler func(http.ResponseWriter, *http.Request)) http.ResponseWriter {
	w := httptest.NewRecorder()
	// mock request
	req := httptest.NewRequest(method, url, strings.NewReader(""))

	handler(w, req)

	return w
}

func TestPOSTHandler(t *testing.T) {

}
