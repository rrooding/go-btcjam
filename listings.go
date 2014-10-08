package btcjam

// This holds all listing data
type Listing struct {
  Data ListingData `json:"listing"`
}

// All the data belonging to a single listing
type ListingData struct {
  Id                    int
  Title                 string
  Term                  int     `json:"term_days"`
  Description           string
  Amount                string
  AmountFunded          string
  NumberOfPayments      int
  PaymentCycle          int     `json:"payment_cycle_days"`
  StartDate             string
  EndDate               string
  Rate                  string
  PeriodicPaymentAmount float64
  DenominatedIn         string
  ListingStatus         string
  IsSecured             bool
  LoanPurpose           string
  PercentSecured        int
  PercentFunded         int
  FundingThreshold      float64
  ListingScore          int
  User
}

// All information about a user belonging to a listing
type User struct {
  Country               string
  PositiveRepCount      int     `json:"positive_count_reputation"`
  NegativeRepCount      int     `json:"negative_count_reputation"`
  CanBorrow             bool
  CanTrade              bool
  BTCTalkVerified       bool    `json:"bitcointalk_accound_verified"`
  BTCJamScore           string  `json:"btcjam_score_numeric"`
  BTCJamCode            string  `json:"btcjam_score"`
  AddressVerified       bool
  IdentityVerified      bool
  PhoneVerified         bool
  FacebookConnected     bool
  FacebookFriendCount   int
  LinkedinConnected     bool
  EbayConnected         bool
  EbayAccountDate       string
  EbayFeedbackScore     string
  PaypalConnected       bool    `json:"paypal_verified_account_connected"`
  PaypalAccountDate     string
  RepaidLoansCount      int
  RepaidLoansAmount     float64
  LateLoansCount        int
  LateLoansAmount       float64
  OpenCreditLinesCount  int
  OpenCreditLinesAmount float64
  MadeLatePayments      int
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

  return listings, nil
}
