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

// GetProject retrieves a Neon project by ID
func (c *Client) GetProject(projectId string) (*ProjectState, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/projects/%s", projectId), nil)
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

// UpdateProject updates an existing Neon project
func (c *Client) UpdateProject(projectId string, name string) (*ProjectState, error) {
	body := struct {
		Project struct {
			Name string `json:"name"`
		} `json:"project"`
	}{
		Project: struct {
			Name string `json:"name"`
		}{
			Name: name,
		},
	}

	resp, err := c.doRequest("PATCH", fmt.Sprintf("/projects/%s", projectId), body)
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

// DeleteProject deletes an existing Neon project
func (c *Client) DeleteProject(projectId string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/projects/%s", projectId), nil)
	return err
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

// GetBranch retrieves a Neon branch by ID
func (c *Client) GetBranch(projectId, branchId string) (*BranchState, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/projects/%s/branches/%s", projectId, branchId), nil)
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

// UpdateBranch updates an existing Neon branch
func (c *Client) UpdateBranch(projectId, branchId, name string) (*BranchState, error) {
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

	resp, err := c.doRequest("PATCH", fmt.Sprintf("/projects/%s/branches/%s", projectId, branchId), body)
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

// DeleteBranch deletes an existing Neon branch
func (c *Client) DeleteBranch(projectId, branchId string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/projects/%s/branches/%s", projectId, branchId), nil)
	return err
}

// CreateEndpoint creates a new endpoint in a Neon project
func (c *Client) CreateEndpoint(projectId, branchId, endpointType string) (*EndpointState, error) {
	body := struct {
		Endpoint struct {
			BranchId string `json:"branch_id"`
			Type     string `json:"type"`
		} `json:"endpoint"`
	}{
		Endpoint: struct {
			BranchId string `json:"branch_id"`
			Type     string `json:"type"`
		}{
			BranchId: branchId,
			Type:     endpointType,
		},
	}

	resp, err := c.doRequest("POST", fmt.Sprintf("/projects/%s/endpoints", projectId), body)
	if err != nil {
		return nil, err
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

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}

	return &EndpointState{
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

// CreateDatabase creates a new database in a Neon project
func (c *Client) CreateDatabase(projectId, branchId, name string) (*DatabaseState, error) {
	body := struct {
		Database struct {
			Name      string `json:"name"`
			OwnerName string `json:"owner_name"`
		} `json:"database"`
	}{
		Database: struct {
			Name      string `json:"name"`
			OwnerName string `json:"owner_name"`
		}{
			Name:      name,
			OwnerName: "default", // We're using a default owner here. You might want to make this configurable.
		},
	}

	resp, err := c.doRequest("POST", fmt.Sprintf("/projects/%s/branches/%s/databases", projectId, branchId), body)
	if err != nil {
		return nil, err
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

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}

	return &DatabaseState{
		DatabaseArgs: DatabaseArgs{
			ProjectId: result.Database.ProjectId,
			BranchId:  result.Database.BranchId,
			Name:      result.Database.Name,
		},
		Id:        fmt.Sprintf("%d", result.Database.Id),
		CreatedAt: result.Database.CreatedAt,
	}, nil
}

// CreateRole creates a new role in a Neon project
func (c *Client) CreateRole(projectId, branchId, name string) (*RoleState, error) {
	body := struct {
		Role struct {
			Name string `json:"name"`
		} `json:"role"`
	}{
		Role: struct {
			Name string `json:"name"`
		}{
			Name: name,
		},
	}

	resp, err := c.doRequest("POST", fmt.Sprintf("/projects/%s/branches/%s/roles", projectId, branchId), body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Role struct {
			Name      string `json:"name"`
			Password  string `json:"password"`
			Protected bool   `json:"protected"`
			CreatedAt string `json:"created_at"`
		} `json:"role"`
	}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}

	return &RoleState{
		RoleArgs: RoleArgs{
			ProjectId: projectId,
			BranchId:  branchId,
			Name:      result.Role.Name,
		},
		Id:        result.Role.Name, // Using the name as the ID since the API doesn't return a separate ID
		CreatedAt: result.Role.CreatedAt,
	}, nil
}