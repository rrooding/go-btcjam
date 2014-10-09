package btcjam

import (
  "net/http"
  "net/http/httptest"
  "net/url"
  "testing"
  "encoding/json"
)

var (
  // mux is the HTTP request multiplexer used with the test server.
  test_mux *http.ServeMux

  // client is the BTCJam API client being tested.
  client *Client

  // Server is a test HTTP server to mock responses.
  server *httptest.Server
)

func setup() {
  test_mux = http.NewServeMux()
  server = httptest.NewServer(test_mux)

  client = NewClient(nil)
  url, _ := url.Parse(server.URL)
  client.BaseURL = url
}

func teardown() {
  server.Close()
}

func urlPath(t *testing.T, routeName string) string {
  url, err := client.url(routeName, nil)
  if err != nil {
    t.Fatalf("Error constructing URL path for route %q: %s", routeName, err)
  }
  return "/" + url.Path
}

func writeJSON(w http.ResponseWriter, v interface{}) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  err := json.NewEncoder(w).Encode(v)
  if err != nil {
    panic("writeJSON: " + err.Error())
  }
}

func testMethod(t *testing.T, r *http.Request, want string) {
  if want != r.Method {
    t.Errorf("Request method = %v, want %v", r.Method, want)
  }
}

func testQueryString(t *testing.T, r *http.Request, key string, wantValue string) {
  value := r.URL.Query().Get(key)
  if value != wantValue {
    t.Errorf("Expected query param %v to have value '%v', got '%v'", key, wantValue, value)
  }
}

func Test_decodeRawJsonFloatDecodesIntegers(t *testing.T) {
  wantValue := 1.0
  message := []byte(`1`)

  value, err := decodeRawJsonFloat(message)
  if err != nil {
    t.Fatalf("Error decoding %#v", message)
  }
  if value != wantValue {
    t.Errorf("Expected value to be %f, got %f instead", wantValue, value)
  }
}

func Test_decodeRawJsonFloatDecodesString(t *testing.T) {
  wantValue := 1.234
  message := []byte(`"1.234"`)

  value, err := decodeRawJsonFloat(message)
  if err != nil {
    t.Fatalf("Error decoding %#v", message)
  }
  if value != wantValue {
    t.Errorf("Expected value to be %f, got %f instead", wantValue, value)
  }
}

func Test_decodeRawJsonFloatDecodesFloat(t *testing.T) {
  wantValue := 1.234
  message := []byte(`1.234`)

  value, err := decodeRawJsonFloat(message)
  if err != nil {
    t.Fatalf("Error decoding %#v", message)
  }
  if value != wantValue {
    t.Errorf("Expected value to be %f, got %f instead", wantValue, value)
  }
}
