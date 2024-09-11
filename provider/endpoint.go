package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Endpoint struct{}

type EndpointArgs struct {
	ProjectId string `pulumi:"projectId"`
	BranchId  string `pulumi:"branchId"`
	Type      string `pulumi:"type"`
}

type EndpointState struct {
	EndpointArgs
	Id        string `pulumi:"id"`
	Host      string `pulumi:"host"`
	CreatedAt string `pulumi:"createdAt"`
}

func (e Endpoint) Create(ctx context.Context, name string, input EndpointArgs, preview bool) (string, EndpointState, error) {
	if preview {
		return name, EndpointState{EndpointArgs: input}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", EndpointState{}, fmt.Errorf("missing configuration")
	}

	endpointData := struct {
		BranchId string `json:"branch_id"`
		Type     string `json:"type"`
	}{
		BranchId: input.BranchId,
		Type:     input.Type,
	}

	jsonData, err := json.Marshal(map[string]interface{}{"endpoint": endpointData})
	if err != nil {
		return "", EndpointState{}, fmt.Errorf("failed to marshal endpoint data: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/projects/%s/endpoints", baseURL, input.ProjectId), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", EndpointState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", EndpointState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", EndpointState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", EndpointState{}, fmt.Errorf("failed to create endpoint: %s", string(body))
	}

	var result struct {
		Endpoint struct {
			Id        string `json:"id"`
			Host      string `json:"host"`
			ProjectId string `json:"project_id"`
			BranchId  string `json:"branch_id"`
			Type      string `json:"type"`
			CreatedAt string `json:"created_at"`
		} `json:"endpoint"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", EndpointState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return name, EndpointState{
		EndpointArgs: EndpointArgs{
			ProjectId: result.Endpoint.ProjectId,
			BranchId:  result.Endpoint.BranchId,
			Type:      result.Endpoint.Type,
		},
		Id:        result.Endpoint.Id,
		Host:      result.Endpoint.Host,
		CreatedAt: result.Endpoint.CreatedAt,
	}, nil
}

func (e Endpoint) Read(ctx context.Context, id string, inputs EndpointArgs, state EndpointState) (string, EndpointArgs, EndpointState, error) {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", EndpointArgs{}, EndpointState{}, fmt.Errorf("missing configuration")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s/endpoints/%s", baseURL, state.ProjectId, state.Id), nil)
	if err != nil {
		return "", EndpointArgs{}, EndpointState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", EndpointArgs{}, EndpointState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", EndpointArgs{}, EndpointState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		if IsNotFoundError(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))) {
			return "", EndpointArgs{}, EndpointState{}, nil
		}
		return "", EndpointArgs{}, EndpointState{}, fmt.Errorf("failed to read endpoint: %s", string(body))
	}

	var result struct {
		Endpoint struct {
			Id        string `json:"id"`
			Host      string `json:"host"`
			ProjectId string `json:"project_id"`
			BranchId  string `json:"branch_id"`
			Type      string `json:"type"`
			CreatedAt string `json:"created_at"`
		} `json:"endpoint"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", EndpointArgs{}, EndpointState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return id, EndpointArgs{
		ProjectId: result.Endpoint.ProjectId,
		BranchId:  result.Endpoint.BranchId,
		Type:      result.Endpoint.Type,
	}, EndpointState{
		EndpointArgs: EndpointArgs{
			ProjectId: result.Endpoint.ProjectId,
			BranchId:  result.Endpoint.BranchId,
			Type:      result.Endpoint.Type,
		},
		Id:        result.Endpoint.Id,
		Host:      result.Endpoint.Host,
		CreatedAt: result.Endpoint.CreatedAt,
	}, nil
}

func (e Endpoint) Update(ctx context.Context, id string, olds EndpointState, news EndpointArgs, preview bool) (EndpointState, error) {
	if preview {
		return EndpointState{
			EndpointArgs: news,
			Id:           olds.Id,
			Host:         olds.Host,
			CreatedAt:    olds.CreatedAt,
		}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return EndpointState{}, fmt.Errorf("missing configuration")
	}

	endpointData := struct {
		BranchId string `json:"branch_id"`
		Type     string `json:"type"`
	}{
		BranchId: news.BranchId,
		Type:     news.Type,
	}

	jsonData, err := json.Marshal(map[string]interface{}{"endpoint": endpointData})
	if err != nil {
		return EndpointState{}, fmt.Errorf("failed to marshal endpoint data: %v", err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/projects/%s/endpoints/%s", baseURL, news.ProjectId, olds.Id), bytes.NewBuffer(jsonData))
	if err != nil {
		return EndpointState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return EndpointState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return EndpointState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return EndpointState{}, fmt.Errorf("failed to update endpoint: %s", string(body))
	}

	var result struct {
		Endpoint struct {
			Id        string `json:"id"`
			Host      string `json:"host"`
			ProjectId string `json:"project_id"`
			BranchId  string `json:"branch_id"`
			Type      string `json:"type"`
			CreatedAt string `json:"created_at"`
		} `json:"endpoint"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return EndpointState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return EndpointState{
		EndpointArgs: EndpointArgs{
			ProjectId: result.Endpoint.ProjectId,
			BranchId:  result.Endpoint.BranchId,
			Type:      result.Endpoint.Type,
		},
		Id:        result.Endpoint.Id,
		Host:      result.Endpoint.Host,
		CreatedAt: result.Endpoint.CreatedAt,
	}, nil
}

func (e Endpoint) Delete(ctx context.Context, id string, state EndpointState) error {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return fmt.Errorf("missing configuration")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s/endpoints/%s", baseURL, state.ProjectId, state.Id), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete endpoint: %s", string(body))
	}

	return nil
}