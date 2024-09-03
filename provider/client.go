package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://console.neon.tech/api/v2"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", baseURL, path)

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// CreateProject creates a new Neon project
func (c *Client) CreateProject(name, regionId string) (*ProjectState, error) {
	body := struct {
		Project struct {
			Name     string `json:"name"`
			RegionId string `json:"region_id"`
		} `json:"project"`
	}{
		Project: struct {
			Name     string `json:"name"`
			RegionId string `json:"region_id"`
		}{
			Name:     name,
			RegionId: regionId,
		},
	}

	resp, err := c.doRequest("POST", "/projects", body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Project struct {
			Id        string `json:"id"`
			Name      string `json:"name"`
			RegionId  string `json:"region_id"`
			CreatedAt string `json:"created_at"`
		} `json:"project"`
	}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}

	return &ProjectState{
		ProjectArgs: ProjectArgs{
			Name:     result.Project.Name,
			RegionId: result.Project.RegionId,
		},
		Id:        result.Project.Id,
		CreatedAt: result.Project.CreatedAt,
	}, nil
}

// CreateBranch creates a new branch in a Neon project
func (c *Client) CreateBranch(projectId, name string) (*BranchState, error) {
	body := struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	}{
		Branch: struct {
			Name string `json:"name"`
		}{
			Name: name,
		},
	}

	resp, err := c.doRequest("POST", fmt.Sprintf("/projects/%s/branches", projectId), body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Branch struct {
			Id        string `json:"id"`
			Name      string `json:"name"`
			ProjectId string `json:"project_id"`
			CreatedAt string `json:"created_at"`
		} `json:"branch"`
	}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}

	return &BranchState{
		BranchArgs: BranchArgs{
			ProjectId: result.Branch.ProjectId,
			Name:      result.Branch.Name,
		},
		Id:        result.Branch.Id,
		CreatedAt: result.Branch.CreatedAt,
	}, nil
}

// Add more methods for other resources (Endpoint, Database, Role) here