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
	"strings"
)

type ClientTopic struct {
	BOOTSTRAPSERVERS string
}

func NewClientTopic(bootstrapServers string) *ClientTopic {
	return &ClientTopic{
		BOOTSTRAPSERVERS: bootstrapServers,
	}
}

type TopicInfo struct {
	Topic string `json:"topic"`
}

type TopicModel struct {
	Topic             types.String `tfsdk:"topic"`
	ReplicationFactor types.Int64  `tfsdk:"replication_factor"`
	Partitions        types.Int64  `tfsdk:"partitions"`
	CleanupPolicy     types.String `tfsdk:"cleanup_policy"`
}

func (c *ClientTopic) CreateTopic(schema *TopicModel) error {

	topicInfo := TopicInfo{
		Topic: schema.Topic.ValueString(),
	}

	reqBody, err := json.Marshal(topicInfo)
	if err != nil {
		return err
	}

	topic := strings.Trim(schema.Topic.ValueString(), "\"")

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/subjects/%s/versions", c.BOOTSTRAPSERVERS, topic), bytes.NewBuffer(reqBody))
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

	var topicResponse struct {
		Id int64 `json:"id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&topicResponse)
	if err != nil {
		return err
	}

	//schema.ReplicationFactor = types.Int64Value(topicResponse.Id)
	//version, err := c.GetLatestVersion(subject)
	//if err != nil {
	//	return err
	//}
	//
	//schema.ReplicationFactor = types.Int64Value(version)
	return nil
}

//func (c *ClientTopic) GetLatestVersion(topic string) (int64, error) {
//
//	req, err := http.NewRequest("GET", fmt.Sprintf("%s/subjects/%s/versions", c.APIURL, topic), nil)
//	if err != nil {
//		return 0, err
//	}
//	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return 0, err
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//		}
//	}(resp.Body)
//
//	if resp.StatusCode != http.StatusOK {
//		bodyBytes, _ := io.ReadAll(resp.Body)
//		return 0, errors.New(string(bodyBytes))
//	}
//
//	var topicVersionResponse []int
//
//	err = json.NewDecoder(resp.Body).Decode(&topicVersionResponse)
//	if err != nil {
//		return 0, err
//	}
//
//	schemaVersionResponseLen := len(topicVersionResponse)
//	sort.Ints(topicVersionResponse)
//	version := topicVersionResponse[schemaVersionResponseLen-1]
//
//	return int64(version), nil
//
//}

//func (c *ClientTopic) GetTopic(topic string, version int64) (*TopicModel, error) {
//	req, err := http.NewRequest("GET", fmt.Sprintf("%s/subjects/%s/versions/%d", c.APIURL, topic, version), nil)
//	if err != nil {
//		return nil, err
//	}
//	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//
//		}
//	}(resp.Body)
//
//	if resp.StatusCode != http.StatusOK {
//		bodyBytes, _ := io.ReadAll(resp.Body)
//		return nil, errors.New(string(bodyBytes))
//	}
//
//	var topicResponse struct {
//		Topic             string `json:"topic"`
//		ReplicationFactor int64  `json:"replication_factor"`
//		Partitions        int64  `json:"partitions"`
//		CleanupPolicy     string `json:"cleanup_policy"`
//	}
//
//	err = json.NewDecoder(resp.Body).Decode(&topicResponse)
//	if err != nil {
//		return nil, err
//	}
//
//	//var schemaType string
//	//if topicResponse.SchemaType != "" {
//	//	schemaType = topicResponse.SchemaType
//	//} else {
//	//	schemaType = "AVRO"
//	//}
//
//	schema := &TopicModel{
//		Topic:             types.StringValue(topicResponse.Topic),
//		Partitions:        types.Int64Value(topicResponse.Partitions),
//		ReplicationFactor: types.Int64Value(topicResponse.ReplicationFactor),
//		CleanupPolicy:     types.StringValue(topicResponse.CleanupPolicy),
//	}
//
//	return schema, nil
//}

//func (c *ClientTopic) UpdateTopic(schema *TopicModel) error {
//	// Alias to CreateSchema as schemas are created as new versions
//	return c.CreateTopic(schema)
//}

//func (c *ClientTopic) DeleteTopic(topic string) error {
//	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/subjects/%s", c.topic, topic), nil)
//	if err != nil {
//		return err
//	}
//	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return err
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//
//		}
//	}(resp.Body)
//
//	if resp.StatusCode != http.StatusOK {
//		bodyBytes, _ := io.ReadAll(resp.Body)
//		return errors.New(string(bodyBytes))
//	}
//
//	return nil
//}
