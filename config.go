package main

import (
	"crypto/tls"
	"net/http"

	"github.com/atlassian/go-vtm"
)

// Config is the configuration structure used to instantiate the Stingray
// provider.
type Config struct {
	URL       string
	Username  string
	Password  string
	VerifySSL bool
}

//Client returns back a new client for the provider to use
func (c *Config) Client() (*stingray.Client, error) {
	client := newClient(c)

	return client, nil
}

func newClient(c *Config) *stingray.Client {
	if c.VerifySSL {
		return stingray.NewClient(nil, c.URL, c.Username, c.Password)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	return stingray.NewClient(httpClient, c.URL, c.Username, c.Password)
}
