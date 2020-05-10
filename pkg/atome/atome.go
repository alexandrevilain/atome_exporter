package atome

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"
)

const (
	baseURL           = "https://esoftlink.esoftthings.com"
	sessionCookieName = "PHPSESSID"

	storageCookieKey         = "cookie"
	storageUserIDKey         = "user_id"
	storageSubscriptionIDKey = "subscription_id"

	maxRetry = 3
)

// A Client is an Atome client API
type Client struct {
	logger *logrus.Logger
	client *http.Client
	debug  bool

	username string
	password string

	cookie *http.Cookie
	user   *authenticateResponse
}

// NewClient creates a new instance of an atome API client
func NewClient(logger *logrus.Logger, username, password string, debug bool) *Client {
	return &Client{
		logger:   logger,
		client:   &http.Client{},
		debug:    debug,
		username: username,
		password: password,
	}
}

// Authenticate is ...
func (c *Client) Authenticate() error {
	body, err := json.Marshal(authenticateRequest{
		Email:         c.username,
		PlainPassword: c.password,
	})
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/api/user/login.json", baseURL)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := c.client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if c.debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Fatal(err)
		}

		c.logger.Debugf("%q \n", dump)
	}

	var sessionCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == sessionCookieName {
			sessionCookie = cookie
		}
	}

	if sessionCookie == nil {
		return errors.New("Authentication cookie not found")
	}

	c.logger.Printf("Cookie expires at: %s", sessionCookie.Expires)
	c.logger.Printf("Cookie retrieved value: %s", sessionCookie.Value)

	var response authenticateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	c.cookie = sessionCookie
	c.user = &response

	c.logger.Printf("Authenticated as %s %s", response.Firstname, response.Lastname)

	return nil
}

// RetriveDayConsumption returns the current consumption
func (c *Client) RetriveDayConsumption(retry int) (*Consumption, error) {

	url := fmt.Sprintf("%s/api/subscription/%d/%s/consumption.json?period=sod", baseURL, c.user.ID, c.user.Subscriptions[0].Reference)

	c.logger.Printf("Making request to: %s \n", url)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.AddCookie(c.cookie)
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if c.debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Fatal(err)
		}

		c.logger.Debugf("%q \n", dump)
	}

	if resp.StatusCode == 403 {
		if retry <= maxRetry {
			c.Authenticate()
			return c.RetriveDayConsumption(retry + 1)
		}
		return nil, errors.New("unable to get data: authentication error")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("can't get data, got status: %d", resp.StatusCode)
	}

	var response lastConsumptionResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &Consumption{
		Total:     response.Total,
		Price:     response.Price,
		Co2Impact: response.ImpactCo2,
	}, nil
}
