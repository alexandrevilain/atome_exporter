package atome

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/alexandrevilain/atome_exporter/pkg/storage"
)

const (
	baseURL           = "https://esoftlink.esoftthings.com"
	sessionCookieName = "PHPSESSID"
	databaseFilename  = "atome_exporter.db"

	storageCookieKey         = "cookie"
	storageUserIDKey         = "user_id"
	storageSubscriptionIDKey = "subscription_id"
)

// A Client is an Atome client API
type Client struct {
	username string
	password string
	client   *http.Client
	storage  *storage.Storage
}

// NewClient creates a new instance of an atome API client
func NewClient(username, password string) (*Client, error) {
	storage, err := storage.New(databaseFilename, "atome")
	if err != nil {
		return nil, err
	}

	return &Client{
		username: username,
		password: password,
		client:   &http.Client{},
		storage:  storage,
	}, nil
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

	var sessionCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == sessionCookieName {
			sessionCookie = cookie
		}
	}

	if sessionCookie == nil {
		return errors.New("Authentication cookie not found")
	}

	log.Printf("Cookie expires at: %s", sessionCookie.Expires)

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

	log.Printf("Authenticated as %s %s", response.Firstname, response.Lastname)

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
