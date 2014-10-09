package btcjam

import (
  "fmt"
  "strings"
  "net/url"
  "net/http"
  "reflect"
  "strconv"
  "encoding/json"

  "github.com/google/go-querystring/query"
)

type Client struct {
  // Services used to communicate with different parts of BTCJam.
  Listings ListingsService

  // Base URL for API requests, which should have a trailing slash.
  BaseURL *url.URL

  // Contains a struct with the authentication options for the API
  AuthenticationOptions *AuthOptions

  // HTTP client used to communicate with the BTCJam API.
  httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
  if httpClient == nil {
    cloned := *http.DefaultClient
    httpClient = &cloned
  }

  c := new(Client)
  c.httpClient = httpClient

  c.Listings = &listingsService{c}

  c.BaseURL = &url.URL{Scheme: "https", Host: "btcjam.com", Path: "/api/v1/"}

  return c
}

// SetAuthentication allows the user to set the id and secret
// used to authenticate against the API
func (c *Client) SetAuthentication(applicationID string, applicationSecret string) {
  c.AuthenticationOptions = &AuthOptions{applicationID, applicationSecret}
}

// apiRouter is used to generate URLs for the BTCJam API
var apiRouter = NewAPIRouter()

// url generates the URL to the named BTCJam endpoint, using the specified
// route variables and query options.
func (c *Client) url(apiRouteName string, opt interface{}) (*url.URL, error) {
  route := apiRouter.Get(apiRouteName)
  if route == nil {
    return nil, fmt.Errorf("no API route named %q", apiRouteName)
  }

  url, err := route.URL()
  if err != nil {
    return nil, err
  }

  // Make the route URL path relative to BaseURL by trimming the lead "/"
  url.Path = strings.TrimPrefix(url.Path, "/")

  if opt != nil {
    err = addOptions(url, opt)
    if err != nil {
      return nil, err
    }
  }

  return url, nil
}

func (c *Client) NewRequest(method, urlStr string) (*http.Request, error) {
  rel, err := url.Parse(urlStr)
  if err != nil {
    return nil, err
  }

  u := c.BaseURL.ResolveReference(rel)

  req, err := http.NewRequest(method, u.String(), nil)
  if err != nil {
    return nil, err
  }

  return req, nil
}

// Do sends an API request and returns the results.
func (c *Client) Do(req *http.Request, v interface{}) error {
  resp, err := c.httpClient.Do(req)
  if err != nil {
    return err
  }

  defer resp.Body.Close()

  err = CheckResponse(resp)
  if err != nil {
    return err
  }

  if v != nil {
    err = json.NewDecoder(resp.Body).Decode(v)
    if err != nil {
      return fmt.Errorf("error reading response from %s %s: %s", req.Method, req.URL.RequestURI(), err)
    }
  }

  return nil
}

// addOptions adds the parameters in opt as URL query parameters to u. opt
// must be a struct whose fields may contain "url" tags
func addOptions(u *url.URL, opt interface{}) error {
  v := reflect.ValueOf(opt)
  if v.Kind() == reflect.Ptr && v.IsNil() {
    return nil
  }

  qs, err := query.Values(opt)
  if err != nil {
    return err
  }

  u.RawQuery = qs.Encode()

  return nil
}

// AuthOptions specifies key and secret used to authenticate to the API
type AuthOptions struct {
  ApplicationID string `url:"appid,omitempty"`
  ApplicationSecret string `url:"secret,omitempty"`
}

// CheckResponse checks the API response for errors. A response is considered
// an error if the status code is outside the 200 range
func CheckResponse(r *http.Response) error {
  if c := r.StatusCode; 200 <= c && c <= 299 {
    return nil
  }

  return fmt.Errorf(
    "%v %v: %d",
    r.Request.Method,
    r.Request.URL,
    r.StatusCode)
}

// decodeRawJsonFloat decodes a json.RawMessage field to a float by trying to convert
// it from different formats. This is needed because for example the amount_funded
// field in the BTCJam API can return a float (1.0000), a string ("1.0000") or an
// integer (0). The standard Decode cannot handle this.
func decodeRawJsonFloat(field json.RawMessage) (val float64, err error) {
    var n int
    if err = json.Unmarshal(field, &n); err == nil {
      val = float64(n)
      return val, nil
    }

    var s string
    if err = json.Unmarshal(field, &s); err == nil {
      val, err = strconv.ParseFloat(s, 64)
      if err != nil {
        return 0.0, err
      }
      return val, nil
    }

    if err = json.Unmarshal(field, &val); err == nil {
      return val, nil
    }

    return 0.0, fmt.Errorf("Could not convert some fields: %#v", field)
}
