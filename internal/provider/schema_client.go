// Package provider.go
package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
	"net/http"
	"sort"
	"strings"
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
	SchemaType string `json:"schemaType"`
}

type SchemaModel struct {
	Subject    types.String `tfsdk:"subject"`
	Schema     types.String `tfsdk:"schema"`
	SchemaType types.String `tfsdk:"schema_type"`
	Version    types.Int64  `tfsdk:"version"`
	Id         types.Int64  `tfsdk:"id"`
}

func (c *ClientSchema) CreateSchema(schema *SchemaModel) error {

	schemaInfo := SchemaInfo{
		Schema:     schema.Schema.ValueString(),
		SchemaType: schema.SchemaType.ValueString(),
	}

	reqBody, err := json.Marshal(schemaInfo)
	if err != nil {
		return err
	}

	subject := strings.Trim(schema.Subject.ValueString(), "\"")

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/subjects/%s/versions", c.APIURL, subject), bytes.NewBuffer(reqBody))
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
		Id int64 `json:"id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&schemaResponse)
	if err != nil {
		return err
	}

	reqVersion, err := http.NewRequest("GET", fmt.Sprintf("%s/subjects/%s/versions", c.APIURL, subject), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")

	respVersion, err := http.DefaultClient.Do(reqVersion)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(respVersion.Body)

	if respVersion.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(respVersion.Body)
		return errors.New(string(bodyBytes))
	}

	var schemaVersionResponse []int

	err = json.NewDecoder(respVersion.Body).Decode(&schemaVersionResponse)
	if err != nil {
		return err
	}

	schemaVersionResponseLen := len(schemaVersionResponse)
	sort.Ints(schemaVersionResponse)
	version := schemaVersionResponse[schemaVersionResponseLen-1]

	schema.Id = types.Int64Value(schemaResponse.Id)
	schema.Version = types.Int64Value(int64(version))

	return nil
}

func (c *ClientSchema) GetSchema(subject string, version string) (*SchemaModel, error) {
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
		Subject    types.String `tfsdk:"subject"`
		Version    types.Int64  `tfsdk:"version"`
		Id         int64        `tfsdk:"id"`
		Schema     types.String `tfsdk:"schema"`
		SchemaType types.String `tfsdk:"schema_type"`
	}

	err = json.NewDecoder(resp.Body).Decode(&schemaResponse)
	if err != nil {
		return nil, err
	}

	schema := &SchemaModel{
		Subject: schemaResponse.Subject,
		Version: schemaResponse.Version,
		//Id:         schemaResponse.Id,
		Schema:     schemaResponse.Schema,
		SchemaType: schemaResponse.SchemaType,
	}

	return schema, nil
}

func (c *ClientSchema) UpdateSchema(schema *SchemaModel) error {
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
