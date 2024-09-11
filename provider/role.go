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

type Role struct{}

type RoleArgs struct {
	ProjectId string `pulumi:"projectId"`
	BranchId  string `pulumi:"branchId"`
	Name      string `pulumi:"name"`
}

type RoleState struct {
	RoleArgs
	Id        string `pulumi:"id"`
	CreatedAt string `pulumi:"createdAt"`
}

func (r Role) Create(ctx context.Context, name string, input RoleArgs, preview bool) (string, RoleState, error) {
	if preview {
		return name, RoleState{RoleArgs: input}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", RoleState{}, fmt.Errorf("missing configuration")
	}

	roleData := struct {
		Name string `json:"name"`
	}{
		Name: input.Name,
	}

	jsonData, err := json.Marshal(map[string]interface{}{"role": roleData})
	if err != nil {
		return "", RoleState{}, fmt.Errorf("failed to marshal role data: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/projects/%s/branches/%s/roles", baseURL, input.ProjectId, input.BranchId), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", RoleState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", RoleState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", RoleState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", RoleState{}, fmt.Errorf("failed to create role: %s", string(body))
	}

	var result struct {
		Role struct {
			Name      string `json:"name"`
			Password  string `json:"password"`
			Protected bool   `json:"protected"`
			CreatedAt string `json:"created_at"`
		} `json:"role"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", RoleState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return name, RoleState{
		RoleArgs: RoleArgs{
			ProjectId: input.ProjectId,
			BranchId:  input.BranchId,
			Name:      result.Role.Name,
		},
		Id:        result.Role.Name,
		CreatedAt: result.Role.CreatedAt,
	}, nil
}

func (r Role) Read(ctx context.Context, id string, inputs RoleArgs, state RoleState) (string, RoleArgs, RoleState, error) {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", RoleArgs{}, RoleState{}, fmt.Errorf("missing configuration")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s/branches/%s/roles/%s", baseURL, state.ProjectId, state.BranchId, state.Name), nil)
	if err != nil {
		return "", RoleArgs{}, RoleState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", RoleArgs{}, RoleState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", RoleArgs{}, RoleState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		if IsNotFoundError(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))) {
			return "", RoleArgs{}, RoleState{}, nil
		}
		return "", RoleArgs{}, RoleState{}, fmt.Errorf("failed to read role: %s", string(body))
	}

	var result struct {
		Role struct {
			Name      string `json:"name"`
			Protected bool   `json:"protected"`
			CreatedAt string `json:"created_at"`
		} `json:"role"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", RoleArgs{}, RoleState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return id, RoleArgs{
		ProjectId: state.ProjectId,
		BranchId:  state.BranchId,
		Name:      result.Role.Name,
	}, RoleState{
		RoleArgs: RoleArgs{
			ProjectId: state.ProjectId,
			BranchId:  state.BranchId,
			Name:      result.Role.Name,
		},
		Id:        result.Role.Name,
		CreatedAt: result.Role.CreatedAt,
	}, nil
}

func (r Role) Update(ctx context.Context, id string, olds RoleState, news RoleArgs, preview bool) (RoleState, error) {
	if preview {
		return RoleState{
			RoleArgs:  news,
			Id:        olds.Id,
			CreatedAt: olds.CreatedAt,
		}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return RoleState{}, fmt.Errorf("missing configuration")
	}

	roleData := struct {
		Name string `json:"name"`
	}{
		Name: news.Name,
	}

	jsonData, err := json.Marshal(map[string]interface{}{"role": roleData})
	if err != nil {
		return RoleState{}, fmt.Errorf("failed to marshal role data: %v", err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/projects/%s/branches/%s/roles/%s", baseURL, news.ProjectId, news.BranchId, olds.Name), bytes.NewBuffer(jsonData))
	if err != nil {
		return RoleState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return RoleState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RoleState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return RoleState{}, fmt.Errorf("failed to update role: %s", string(body))
	}

	var result struct {
		Role struct {
			Name      string `json:"name"`
			Protected bool   `json:"protected"`
			CreatedAt string `json:"created_at"`
		} `json:"role"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return RoleState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return RoleState{
		RoleArgs: RoleArgs{
			ProjectId: news.ProjectId,
			BranchId:  news.BranchId,
			Name:      result.Role.Name,
		},
		Id:        result.Role.Name,
		CreatedAt: result.Role.CreatedAt,
	}, nil
}

func (r Role) Delete(ctx context.Context, id string, state RoleState) error {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return fmt.Errorf("missing configuration")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s/branches/%s/roles/%s", baseURL, state.ProjectId, state.BranchId, state.Name), nil)
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
		return fmt.Errorf("failed to delete role: %s", string(body))
	}

	return nil
}