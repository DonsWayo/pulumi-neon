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

type Project struct{}

type ProjectArgs struct {
	Name     string `pulumi:"name"`
	RegionId string `pulumi:"regionId"`
}

type ProjectState struct {
	ProjectArgs
	Id        string `pulumi:"id"`
	CreatedAt string `pulumi:"createdAt"`
}

func (p Project) Create(ctx context.Context, name string, input ProjectArgs, preview bool) (string, ProjectState, error) {
	if preview {
		return name, ProjectState{ProjectArgs: input}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", ProjectState{}, fmt.Errorf("missing configuration")
	}

	projectData := struct {
		Name     string `json:"name"`
		RegionId string `json:"region_id"`
	}{
		Name:     input.Name,
		RegionId: input.RegionId,
	}

	jsonData, err := json.Marshal(map[string]interface{}{"project": projectData})
	if err != nil {
		return "", ProjectState{}, fmt.Errorf("failed to marshal project data: %v", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/projects", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", ProjectState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", ProjectState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ProjectState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", ProjectState{}, fmt.Errorf("failed to create project: %s", string(body))
	}

	var result struct {
		Project struct {
			Id        string `json:"id"`
			Name      string `json:"name"`
			RegionId  string `json:"region_id"`
			CreatedAt string `json:"created_at"`
		} `json:"project"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", ProjectState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return name, ProjectState{
		ProjectArgs: ProjectArgs{
			Name:     result.Project.Name,
			RegionId: result.Project.RegionId,
		},
		Id:        result.Project.Id,
		CreatedAt: result.Project.CreatedAt,
	}, nil
}

func (p Project) Read(ctx context.Context, id string, inputs ProjectArgs, state ProjectState) (string, ProjectArgs, ProjectState, error) {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", ProjectArgs{}, ProjectState{}, fmt.Errorf("missing configuration")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s", baseURL, state.Id), nil)
	if err != nil {
		return "", ProjectArgs{}, ProjectState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", ProjectArgs{}, ProjectState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ProjectArgs{}, ProjectState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		if IsNotFoundError(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))) {
			return "", ProjectArgs{}, ProjectState{}, nil
		}
		return "", ProjectArgs{}, ProjectState{}, fmt.Errorf("failed to read project: %s", string(body))
	}

	var result struct {
		Project struct {
			Id        string `json:"id"`
			Name      string `json:"name"`
			RegionId  string `json:"region_id"`
			CreatedAt string `json:"created_at"`
		} `json:"project"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", ProjectArgs{}, ProjectState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return id, ProjectArgs{
		Name:     result.Project.Name,
		RegionId: result.Project.RegionId,
	}, ProjectState{
		ProjectArgs: ProjectArgs{
			Name:     result.Project.Name,
			RegionId: result.Project.RegionId,
		},
		Id:        result.Project.Id,
		CreatedAt: result.Project.CreatedAt,
	}, nil
}

func (p Project) Update(ctx context.Context, id string, olds ProjectState, news ProjectArgs, preview bool) (ProjectState, error) {
	if preview {
		return ProjectState{
			ProjectArgs: news,
			Id:          olds.Id,
			CreatedAt:   olds.CreatedAt,
		}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return ProjectState{}, fmt.Errorf("missing configuration")
	}

	projectData := struct {
		Name string `json:"name"`
	}{
		Name: news.Name,
	}

	jsonData, err := json.Marshal(map[string]interface{}{"project": projectData})
	if err != nil {
		return ProjectState{}, fmt.Errorf("failed to marshal project data: %v", err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/projects/%s", baseURL, olds.Id), bytes.NewBuffer(jsonData))
	if err != nil {
		return ProjectState{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ProjectState{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ProjectState{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return ProjectState{}, fmt.Errorf("failed to update project: %s", string(body))
	}

	var result struct {
		Project struct {
			Id        string `json:"id"`
			Name      string `json:"name"`
			RegionId  string `json:"region_id"`
			CreatedAt string `json:"created_at"`
		} `json:"project"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return ProjectState{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return ProjectState{
		ProjectArgs: ProjectArgs{
			Name:     result.Project.Name,
			RegionId: result.Project.RegionId,
		},
		Id:        result.Project.Id,
		CreatedAt: result.Project.CreatedAt,
	}, nil
}

func (p Project) Delete(ctx context.Context, id string, state ProjectState) error {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return fmt.Errorf("missing configuration")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s", baseURL, state.Id), nil)
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
		return fmt.Errorf("failed to delete project: %s", string(body))
	}

	return nil
}