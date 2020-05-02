package atome

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/alexandrevilain/atome_exporter/pkg/storage"
	"github.com/sirupsen/logrus"
)

const (
	baseURL           = "https://esoftlink.esoftthings.com"
	sessionCookieName = "PHPSESSID"

	storageCookieKey         = "cookie"
	storageUserIDKey         = "user_id"
	storageSubscriptionIDKey = "subscription_id"
)

// A Client is an Atome client API
type Client struct {
	logger   *logrus.Logger
	username string
	password string
	client   *http.Client
	storage  *storage.Storage
}

// NewClient creates a new instance of an atome API client
func NewClient(logger *logrus.Logger, username, password string, storage *storage.Storage) *Client {
	return &Client{
		logger:   logger,
		username: username,
		password: password,
		client:   &http.Client{},
		storage:  storage,
	}
}

// Authenticate is ..
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

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%q \n", dump)

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

	var response authenticateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	err = c.storage.Put(storageCookieKey, sessionCookie.Value)
	if err != nil {
		return err
	}

	err = c.storage.Put(storageUserIDKey, response.ID)
	if err != nil {
		return err
	}

	err = c.storage.Put(storageSubscriptionIDKey, response.Subscriptions[0].Reference)
	if err != nil {
		return err
	}

	c.logger.Printf("Authenticated as %s %s", response.Firstname, response.Lastname)

	return nil
}

// RetriveDayConsumption is ...
func (c *Client) RetriveDayConsumption() (*Consumption, error) {
	var sessionCookieValue string
	err := c.storage.Get(sessionCookieName, &sessionCookieValue)
	if err != nil {
		return nil, err
	}

	var userID int
	err = c.storage.Get(storageUserIDKey, &userID)
	if err != nil {
		return nil, err
	}

	var subscriptionID string
	err = c.storage.Get(storageSubscriptionIDKey, &subscriptionID)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/subscription/%d/%s/consumption.json?period=sod", baseURL, userID, subscriptionID)

	c.logger.Printf("Making request to: %s \n", url)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: sessionCookieValue})
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%q \n", dump)

	if resp.StatusCode == 403 {
		c.Authenticate()
		return c.RetriveDayConsumption()
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
