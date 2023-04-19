// client.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
)

type Client struct {
	APIURL string
}

func NewClient(apiURL string) *Client {
	return &Client{
		APIURL: apiURL,
	}
}

type Topic struct {
	Name             string `json:"name"`
	Partitions       int    `json:"partitions"`
	ReplicationFactor int    `json:"replication_factor"`
}

func (c *Client) CreateTopic(topic *Topic) error {
	reqBody, err := json.Marshal(topic)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/topics", c.APIURL), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(bodyBytes))
	}

	return nil
}

func (c *Client) GetTopic(topicName string) (*Topic, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/topics/%s", c.APIURL, topicName), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(bodyBytes))
	}

	var topic Topic
	err = json.NewDecoder(resp.Body).Decode(&topic)
	if err != nil {
		return nil, err
	}

	return &topic, nil
}


func (c *Client) UpdateTopic(topic *Topic) error {
	reqBody, err := json.Marshal(topic)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/topics/%s", c.APIURL, topic.Name), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(bodyBytes))
	}

	return nil
}

func (c *Client) DeleteTopic(topicName string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/topics/%s", c.APIURL, topicName), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(bodyBytes))
	}

	return nil
}

type Schema struct {
	Subject     string `json:"subject"`
	Schema      string `json:"schema"`
	SchemaType  string `json:"schema_type,omitempty"`
	Version     int    `json:"version,omitempty"`
	Id          int    `json:"id,omitempty"`
}

func (c *Client) CreateSchema(schema *Schema) error {
	reqBody, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/subjects/%s/versions", c.APIURL, schema.Subject), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(bodyBytes))
	}

	var schemaResponse struct {
		Id int `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&schemaResponse)
	if err != nil {
		return err
	}

	schema.Id = schemaResponse.Id

	return nil
}


func (c *Client) GetSchema(subject string, version string) (*Schema, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/subjects/%s/versions/%s", c.APIURL, subject, version), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(bodyBytes))
	}

	var schemaResponse struct {
		Subject    string `json:"subject"`
		Version    int    `json:"version"`
		Id         int    `json:"id"`
		Schema     string `json:"schema"`
		SchemaType string `json:"schemaType,omitempty"`
	}

	err = json.NewDecoder(resp.Body).Decode(&schemaResponse)
	if err != nil {
		return nil, err
	}

	schema := &Schema{
		Subject:    schemaResponse.Subject,
		Version:    schemaResponse.Version,
		Id:         schemaResponse.Id,
		Schema:     schemaResponse.Schema,
		SchemaType: schemaResponse.SchemaType,
	}

	return schema, nil
}


func (c *Client) UpdateSchema(schema *Schema) error {
    // Alias to CreateSchema as schemas are created as new versions
	return c.CreateSchema(schema)
}


func (c *Client) DeleteSchema(subject string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/subjects/%s", c.APIURL, subject), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(bodyBytes))
	}

	return nil
}
