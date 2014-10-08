package btcjam

import (
  "github.com/gorilla/mux"
)

const (
  ListingsRoute = "Listings"

  ProfileRoute  = "Profile"
)

func NewAPIRouter() *mux.Router {
  m := mux.NewRouter()

  m.StrictSlash(true)

  m.Path("/listings").Methods("GET").Name(ListingsRoute)

  m.Path("/me").Methods("GET").Name(ProfileRoute)

  return m
}
