package btcjam

import (
  "testing"
  "net/http"
  "net/url"

  "github.com/gorilla/mux"
)

func TestMatch(t *testing.T) {
  router := NewAPIRouter()

  tests := []struct {
    path          string
    wantRouteName string
  }{
    // Listings
    {
      path:          "/listings",
      wantRouteName: ListingsRoute,
    },

    // Profile
    {
      path:          "/me",
      wantRouteName: ProfileRoute,
    },
  }

  for _, test := range tests {
    var routeMatch mux.RouteMatch
    match := router.Match(&http.Request{Method: "GET", URL: &url.URL{Path: test.path}}, &routeMatch)

    if !match {
      t.Errorf("%s: Wanted match, did not match", test.path)
    }

    if routeName := routeMatch.Route.GetName(); routeName != test.wantRouteName {
      t.Errorf("%s: got matched route %q, want %q", test.path, routeName, test.wantRouteName)
    }
  }
}
