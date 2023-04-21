// client_schema.go
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ClientSchema struct {
	APIURL string
}

func NewClientSchema(apiURL string) *ClientSchema {
	return &ClientSchema{
		APIURL: apiURL,
	}
}

type SchemaInfo struct {
	Schema     string `json:"schema"`
	SchemaType string `json:"schemaType,omitempty"`
}

type ResourceRedPandaSchemaModel struct {
	Subject    string `json:"subject"`
	Schema     string `json:"schema"`
	SchemaType string `json:"schemaType,omitempty"`
	Version    int    `json:"version,omitempty"`
	Id         int    `json:"id,omitempty"`
}

func (c *ClientSchema) CreateSchema(schema *ResourceRedPandaSchemaModel) error {

	schemaInfo := SchemaInfo{
		Schema:     schema.Schema,
		SchemaType: schema.SchemaType,
	}

	reqBody, err := json.Marshal(schemaInfo)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/subjects/%s/versions", c.APIURL, schema.Subject), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
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

func (c *ClientSchema) GetSchema(subject string, version string) (*ResourceRedPandaSchemaModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/subjects/%s/versions/%s", c.APIURL, subject, version), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
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

	schema := &ResourceRedPandaSchemaModel{
		Subject:    schemaResponse.Subject,
		Version:    schemaResponse.Version,
		Id:         schemaResponse.Id,
		Schema:     schemaResponse.Schema,
		SchemaType: schemaResponse.SchemaType,
	}

	return schema, nil
}

func (c *ClientSchema) UpdateSchema(schema *ResourceRedPandaSchemaModel) error {
	// Alias to CreateSchema as schemas are created as new versions
	return c.CreateSchema(schema)
}

func (c *ClientSchema) DeleteSchema(subject string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/subjects/%s", c.APIURL, subject), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return errors.New(string(bodyBytes))
	}

	return nil
}
