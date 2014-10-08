package btcjam

import (
  "net/http"
  "testing"
  "reflect"
)

func TestListingsService_List(t *testing.T) {
  setup()
  defer teardown()

  client.SetAuthentication("application_id", "application_secret")

  want := []*Listing{&Listing{Data: ListingData{Id: 42}}}

  var called bool
  test_mux.HandleFunc(urlPath(t, ListingsRoute), func(w http.ResponseWriter, r *http.Request) {
    called = true
    testMethod(t, r, "GET")

    testQueryString(t, r, "appid", "application_id")
    testQueryString(t, r, "secret", "application_secret")

    writeJSON(w, want)
  })

  listings, err := client.Listings.List()
  if err != nil {
    t.Errorf("Listings.List returned error: %v", err)
  }

  if !called {
    t.Fatal("!called")
  }

  if !reflect.DeepEqual(listings, want) {
    t.Errorf("Listings.List returned %+v, want %+v", listings, want, "\n")
  }
}
