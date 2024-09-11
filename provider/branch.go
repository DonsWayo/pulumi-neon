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

type Branch struct{}

type BranchArgs struct {
	ProjectId string `pulumi:"projectId"`
	Name      string `pulumi:"name"`
}

type BranchState struct {
	BranchArgs
	Id        string `pulumi:"id"`
	CreatedAt string `pulumi:"createdAt"`
}

func (b Branch) Create(ctx context.Context, name string, input BranchArgs, preview bool) (string, BranchState, error) {
	if preview {
		return name, BranchState{BranchArgs: input}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", BranchState{}, fmt.Errorf("missing configuration")
	}

	branchData := struct {
		Name string `json:"name"`
	}{
		Name: input.Name,
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"branch": branchData,
		"endpoints": []map[string]string{
			{"type": "read_only"},
		},
	})
	if err != nil {
		return "", BranchState{}, fmt.Errorf("failed to marshal branch data: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/projects/%s/branches", baseURL, input.ProjectId), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", BranchState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", BranchState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", BranchState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", BranchState{}, fmt.Errorf("failed to create branch: %s", string(body))
	}

	var result struct {
		Branch struct {
			Id        string `json:"id"`
			Name      string `json:"name"`
			ProjectId string `json:"project_id"`
			CreatedAt string `json:"created_at"`
		} `json:"branch"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", BranchState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return name, BranchState{
		BranchArgs: BranchArgs{
			ProjectId: result.Branch.ProjectId,
			Name:      result.Branch.Name,
		},
		Id:        result.Branch.Id,
		CreatedAt: result.Branch.CreatedAt,
	}, nil
}

func (b Branch) Read(ctx context.Context, id string, inputs BranchArgs, state BranchState) (string, BranchArgs, BranchState, error) {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", BranchArgs{}, BranchState{}, fmt.Errorf("missing configuration")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s/branches/%s", baseURL, state.ProjectId, state.Id), nil)
	if err != nil {
		return "", BranchArgs{}, BranchState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", BranchArgs{}, BranchState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", BranchArgs{}, BranchState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		if IsNotFoundError(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))) {
			return "", BranchArgs{}, BranchState{}, nil
		}
		return "", BranchArgs{}, BranchState{}, fmt.Errorf("failed to read branch: %s", string(body))
	}

	var result struct {
		Branch struct {
			Id        string `json:"id"`
			Name      string `json:"name"`
			ProjectId string `json:"project_id"`
			CreatedAt string `json:"created_at"`
		} `json:"branch"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", BranchArgs{}, BranchState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return id, BranchArgs{
		ProjectId: result.Branch.ProjectId,
		Name:      result.Branch.Name,
	}, BranchState{
		BranchArgs: BranchArgs{
			ProjectId: result.Branch.ProjectId,
			Name:      result.Branch.Name,
		},
		Id:        result.Branch.Id,
		CreatedAt: result.Branch.CreatedAt,
	}, nil
}

func (b Branch) Update(ctx context.Context, id string, olds BranchState, news BranchArgs, preview bool) (BranchState, error) {
	if preview {
		return BranchState{
			BranchArgs: news,
			Id:         olds.Id,
			CreatedAt:  olds.CreatedAt,
		}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return BranchState{}, fmt.Errorf("missing configuration")
	}

	branchData := struct {
		Name string `json:"name"`
	}{
		Name: news.Name,
	}

	jsonData, err := json.Marshal(map[string]interface{}{"branch": branchData})
	if err != nil {
		return BranchState{}, fmt.Errorf("failed to marshal branch data: %v", err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/projects/%s/branches/%s", baseURL, news.ProjectId, olds.Id), bytes.NewBuffer(jsonData))
	if err != nil {
		return BranchState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return BranchState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return BranchState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return BranchState{}, fmt.Errorf("failed to update branch: %s", string(body))
	}

	var result struct {
		Branch struct {
			Id        string `json:"id"`
			Name      string `json:"name"`
			ProjectId string `json:"project_id"`
			CreatedAt string `json:"created_at"`
		} `json:"branch"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return BranchState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return BranchState{
		BranchArgs: BranchArgs{
			ProjectId: result.Branch.ProjectId,
			Name:      result.Branch.Name,
		},
		Id:        result.Branch.Id,
		CreatedAt: result.Branch.CreatedAt,
	}, nil
}

func (b Branch) Delete(ctx context.Context, id string, state BranchState) error {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return fmt.Errorf("missing configuration")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s/branches/%s", baseURL, state.ProjectId, state.Id), nil)
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
		return fmt.Errorf("failed to delete branch: %s", string(body))
	}

	return nil
}