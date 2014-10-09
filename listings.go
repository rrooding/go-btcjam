package btcjam

import (
  "fmt"
  "encoding/json"
)

// This holds all listing data
type Listing struct {
  Data                  ListingData     `json:"listing"`
}

// All the data belonging to a single listing
type ListingData struct {
  Id                    int
  Title                 string
  Term                  int
  Description           string
  Amount                float64         `json:",string"`
  AmountFundedRaw       json.RawMessage `json:"amount_funded,omitempty"`
  AmountFunded          float64
  NumberOfPayments      int             `json:"number_of_payments"`
  PaymentCycle          int             `json:"payment_cycle_days"`
  StartDate             string          `json:"start_date"`
  EndDate               string          `json:"end_date"`
  Rate                  float64         `json:"rate,string"`
  PeriodicPaymentAmount float64         `json:"periodic_payment_amount"`
  DenominatedIn         string          `json:"denominated_in"`
  ListingStatus         string          `json:"listing_status"`
  LoanPurpose           string          `json:"loan_purpose"`
  PercentFunded         int             `json:"percent_funded"`
  ListingScore          int             `json:"listing_score"`
  User                  User            `json:"user"`
}

// All information about a user belonging to a listing
type User struct {
  Country               string
  PositiveRepCount      int             `json:"positive_count_reputation"`
  NegativeRepCount      int             `json:"negative_count_reputation"`
  CanBorrow             bool            `json:"can_borrow"`
  CanTrade              bool            `json:"can_trade"`
  BTCTalkVerified       bool            `json:"bitcointalk_accound_verified"`
  BTCJamScore           float64         `json:"btcjam_score_numeric,string"`
  BTCJamCode            string          `json:"btcjam_score"`
  AddressVerified       bool            `json:"address_verified"`
  IdentityVerified      bool            `json:"identity_verified"`
  PhoneVerified         bool            `json:"phone_verified"`
  FacebookConnected     bool            `json:"facebook_connected"`
  FacebookFriendCount   int             `json:"facebook_friend_count"`
  LinkedinConnected     bool            `json:"linkedin_connected"`
  PaypalConnected       bool            `json:"paypal_verified_account_connected"`
  PaypalAccountDate     string          `json:"paypal_account_date"`
  RepaidLoansCount      int             `json:"repaid_loans_count"`
  RepaidLoansAmount     float64         `json:"repaid_loans_amount"`
  LateLoansCount        int             `json:"late_loans_count"`
  LateLoansAmount       float64         `json:"late_loans_amount"`
  OpenCreditLinesCount  int             `json:"open_credit_lines_count"`
  OpenCreditLinesAmount float64         `json:"open_credit_lines_amount"`
  MadeLatePaymentsCount int             `json:"made_late_payments_count"`
}

// ListingsService communicates with the listings endpoint of the BTCJam API.
type ListingsService interface {
  // List all listings
  List() ([]*Listing, error)
}

// listingsService implements ListingsService
type listingsService struct {
  client *Client
}

// Fetches all listings from BTCJam
func (s *listingsService) List() ([]*Listing, error) {
  opt := s.client.AuthenticationOptions
  url, err := s.client.url(ListingsRoute, opt)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest("GET", url.String())
  if err != nil {
    return nil, err
  }

  var listings []*Listing
  err = s.client.Do(req, &listings)
  if err != nil {
    return nil, err
  }

  for _, listing := range listings {
    if listing.Data.AmountFundedRaw == nil { continue }

    listing.Data.AmountFunded, err = decodeRawJsonFloat(listing.Data.AmountFundedRaw)
    if err != nil {
      fmt.Printf("[btcjam] error decoding AmountFunded from %#v for listing %d", listing.Data.AmountFundedRaw, listing.Data.Id)
      listing.Data.AmountFunded = 0.0
    }
  }

  return listings, nil
}
