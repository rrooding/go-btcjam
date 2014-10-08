package btcjam

// This holds all profile data
type Profile struct {
}

// ProfileService communicates with the profile endpoint of the BTCJam API.
type ProfileService interface {
  // Get the profile
  Get() (*Profile, error)
}

// profileService implements ProfileService
type profileService struct {
  client *Client
}

// Get the profile
func (s *profileService) Get() (*Profile, error) {
  url, err := s.client.url(ProfileRoute, nil)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest("GET", url.String())
  if err != nil {
    return nil, err
  }

  var profile *Profile
  err = s.client.Do(req, &profile)
  if err != nil {
    return nil, err
  }

  return profile, nil
}
