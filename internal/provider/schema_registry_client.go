// Package provider
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

//type ClientSchemaRegistry struct {
//	APIURL string
//}

//func NewClientSchemaRegistry(apiURL string) *ClientSchemaRegistry {
//	return &ClientSchemaRegistry{
//		APIURL: apiURL,
//	}
//}

func NewClientSchemaRegistry(apiURL types.String) *RedPandaProvider {
	return &RedPandaProvider{
		SchemaRegistryApiUrl: apiURL,
	}
}

type SchemaRegistryInfo struct {
	Schema     string `json:"schema"`
	SchemaType string `json:"schemaType"`
}

type SchemaRegistryModel struct {
	Subject    types.String `tfsdk:"subject"`
	Schema     types.String `tfsdk:"schema"`
	SchemaType types.String `tfsdk:"schema_type"`
	Version    types.Int64  `tfsdk:"version"`
	Id         types.Int64  `tfsdk:"id"`
}

func (p *RedPandaProvider) CreateSchema(schema *SchemaRegistryModel) error {

	SchemaRegistryInfo := SchemaRegistryInfo{
		Schema:     schema.Schema.ValueString(),
		SchemaType: schema.SchemaType.ValueString(),
	}

	reqBody, err := json.Marshal(SchemaRegistryInfo)
	if err != nil {
		return err
	}

	subject := strings.Trim(schema.Subject.ValueString(), "\"")
	schemaRegistryUrl := p.SchemaRegistryApiUrl.ValueString()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/subjects/%s/versions", schemaRegistryUrl, subject), bytes.NewBuffer(reqBody))
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

	schema.Id = types.Int64Value(schemaResponse.Id)
	version, err := p.GetLatestVersion(subject)
	if err != nil {
		return err
	}

	schema.Id = types.Int64Value(version)
	return nil
}

func (p *RedPandaProvider) GetLatestVersion(subject string) (int64, error) {

	schemaRegistryUrl := p.SchemaRegistryApiUrl.String()
	fmt.Println("url: " + schemaRegistryUrl)
	fmt.Println("subject: " + subject)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/subjects/%s/versions", schemaRegistryUrl, subject), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, errors.New(string(bodyBytes))
	}

	var schemaVersionResponse []int

	err = json.NewDecoder(resp.Body).Decode(&schemaVersionResponse)
	if err != nil {
		return 0, err
	}

	schemaVersionResponseLen := len(schemaVersionResponse)
	sort.Ints(schemaVersionResponse)
	version := schemaVersionResponse[schemaVersionResponseLen-1]

	return int64(version), nil

}

func (p *RedPandaProvider) GetSchema(subject string, version int64) (*SchemaRegistryModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/subjects/%s/versions/%d", p.SchemaRegistryApiUrl, subject, version), nil)
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
		Version    int64  `json:"version"`
		Id         int64  `json:"id"`
		Schema     string `json:"schema"`
		SchemaType string `json:"schemaType"`
	}

	err = json.NewDecoder(resp.Body).Decode(&schemaResponse)
	if err != nil {
		return nil, err
	}

	var schemaType string
	if schemaResponse.SchemaType != "" {
		schemaType = schemaResponse.SchemaType
	} else {
		schemaType = "AVRO"
	}

	schema := &SchemaRegistryModel{
		Subject:    types.StringValue(schemaResponse.Subject),
		Version:    types.Int64Value(schemaResponse.Version),
		Id:         types.Int64Value(schemaResponse.Id),
		Schema:     types.StringValue(schemaResponse.Schema),
		SchemaType: types.StringValue(schemaType),
	}

	return schema, nil
}

func (p *RedPandaProvider) UpdateSchema(schema *SchemaRegistryModel) error {
	// Alias to CreateSchema as schemas are created as new versions
	return p.CreateSchema(schema)
}

func (p *RedPandaProvider) DeleteSchema(subject string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/subjects/%s", p.SchemaRegistryApiUrl, subject), nil)
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
