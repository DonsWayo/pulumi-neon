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

type Database struct{}

type DatabaseArgs struct {
	ProjectId string `pulumi:"projectId"`
	BranchId  string `pulumi:"branchId"`
	Name      string `pulumi:"name"`
}

type DatabaseState struct {
	DatabaseArgs
	Id        string `pulumi:"id"`
	CreatedAt string `pulumi:"createdAt"`
}

func (d Database) Create(ctx context.Context, name string, input DatabaseArgs, preview bool) (string, DatabaseState, error) {
	if preview {
		return name, DatabaseState{DatabaseArgs: input}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", DatabaseState{}, fmt.Errorf("missing configuration")
	}

	databaseData := struct {
		Name      string `json:"name"`
		OwnerName string `json:"owner_name"`
	}{
		Name:      input.Name,
		OwnerName: "default",
	}

	jsonData, err := json.Marshal(map[string]interface{}{"database": databaseData})
	if err != nil {
		return "", DatabaseState{}, fmt.Errorf("failed to marshal database data: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/projects/%s/branches/%s/databases", baseURL, input.ProjectId, input.BranchId), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", DatabaseState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", DatabaseState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", DatabaseState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", DatabaseState{}, fmt.Errorf("failed to create database: %s", string(body))
	}

	var result struct {
		Database struct {
			Id        int64  `json:"id"`
			Name      string `json:"name"`
			OwnerName string `json:"owner_name"`
			ProjectId string `json:"project_id"`
			BranchId  string `json:"branch_id"`
			CreatedAt string `json:"created_at"`
		} `json:"database"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", DatabaseState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return name, DatabaseState{
		DatabaseArgs: DatabaseArgs{
			ProjectId: result.Database.ProjectId,
			BranchId:  result.Database.BranchId,
			Name:      result.Database.Name,
		},
		Id:        fmt.Sprintf("%d", result.Database.Id),
		CreatedAt: result.Database.CreatedAt,
	}, nil
}

func (d Database) Read(ctx context.Context, id string, inputs DatabaseArgs, state DatabaseState) (string, DatabaseArgs, DatabaseState, error) {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", DatabaseArgs{}, DatabaseState{}, fmt.Errorf("missing configuration")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s/branches/%s/databases/%s", baseURL, state.ProjectId, state.BranchId, state.Name), nil)
	if err != nil {
		return "", DatabaseArgs{}, DatabaseState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", DatabaseArgs{}, DatabaseState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", DatabaseArgs{}, DatabaseState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		if IsNotFoundError(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))) {
			return "", DatabaseArgs{}, DatabaseState{}, nil
		}
		return "", DatabaseArgs{}, DatabaseState{}, fmt.Errorf("failed to read database: %s", string(body))
	}

	var result struct {
		Database struct {
			Id        int64  `json:"id"`
			Name      string `json:"name"`
			OwnerName string `json:"owner_name"`
			ProjectId string `json:"project_id"`
			BranchId  string `json:"branch_id"`
			CreatedAt string `json:"created_at"`
		} `json:"database"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", DatabaseArgs{}, DatabaseState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return id, DatabaseArgs{
		ProjectId: result.Database.ProjectId,
		BranchId:  result.Database.BranchId,
		Name:      result.Database.Name,
	}, DatabaseState{
		DatabaseArgs: DatabaseArgs{
			ProjectId: result.Database.ProjectId,
			BranchId:  result.Database.BranchId,
			Name:      result.Database.Name,
		},
		Id:        fmt.Sprintf("%d", result.Database.Id),
		CreatedAt: result.Database.CreatedAt,
	}, nil
}

func (d Database) Update(ctx context.Context, id string, olds DatabaseState, news DatabaseArgs, preview bool) (DatabaseState, error) {
	if preview {
		return DatabaseState{
			DatabaseArgs: news,
			Id:           olds.Id,
			CreatedAt:    olds.CreatedAt,
		}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return DatabaseState{}, fmt.Errorf("missing configuration")
	}

	databaseData := struct {
		Name string `json:"name"`
	}{
		Name: news.Name,
	}

	jsonData, err := json.Marshal(map[string]interface{}{"database": databaseData})
	if err != nil {
		return DatabaseState{}, fmt.Errorf("failed to marshal database data: %v", err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/projects/%s/branches/%s/databases/%s", baseURL, news.ProjectId, news.BranchId, olds.Name), bytes.NewBuffer(jsonData))
	if err != nil {
		return DatabaseState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return DatabaseState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DatabaseState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return DatabaseState{}, fmt.Errorf("failed to update database: %s", string(body))
	}

	var result struct {
		Database struct {
			Id        int64  `json:"id"`
			Name      string `json:"name"`
			OwnerName string `json:"owner_name"`
			ProjectId string `json:"project_id"`
			BranchId  string `json:"branch_id"`
			CreatedAt string `json:"created_at"`
		} `json:"database"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return DatabaseState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return DatabaseState{
		DatabaseArgs: DatabaseArgs{
			ProjectId: result.Database.ProjectId,
			BranchId:  result.Database.BranchId,
			Name:      result.Database.Name,
		},
		Id:        fmt.Sprintf("%d", result.Database.Id),
		CreatedAt: result.Database.CreatedAt,
	}, nil
}

func (d Database) Delete(ctx context.Context, id string, state DatabaseState) error {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return fmt.Errorf("missing configuration")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s/branches/%s/databases/%s", baseURL, state.ProjectId, state.BranchId, state.Name), nil)
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
		return fmt.Errorf("failed to delete database: %s", string(body))
	}

	return nil
}