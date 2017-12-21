package config

import (
	"net/http"
	"fmt"
	"github.com/pkg/errors"
	"encoding/json"
	"encoding/base64"
)

type Payload struct {
	Name            string            `json:"name"`
	Profiles        []string          `json:"profiles"`
	Label           string            `json:"label"`
	PropertySources []PropertySources `json:"propertySources"`
}

type PropertySources struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

type SpringConfig struct {
	Name     string
	URI      string
	Profile  string
	Username string
	Password string
}

var client = &http.Client{}

func NewSpringConfig(name string, uri string, profile string, username string, password string) *SpringConfig {
	if name == "" || uri == "" || username == "" || password == "" {
		panic(errors.New("One required parameter is empty"))
	}

	if profile == "" {
		profile = "master"
	}

	return &SpringConfig{name, uri, profile, username, password}
}

func (c *SpringConfig) Read() (map[string]string, error) {

	url := fmt.Sprintf("%s/%s/%s/%s", c.URI, c.Name, c.Profile, c.Profile)

	credentialForB64 := fmt.Sprintf("%s:%s", c.Username, c.Password)
	authHeader := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(credentialForB64)))

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", authHeader)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Incorrect response with StatusCode=%d", resp.StatusCode))
	}

	d := json.NewDecoder(resp.Body)

	var payload Payload
	err = d.Decode(&payload)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for k, v := range payload.PropertySources[0].Source {
		result[k] = fmt.Sprintf("%v", v)
	}

	return result, nil
}
