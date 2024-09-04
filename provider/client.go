package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func (c *Client) DeleteProject(projectId string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/projects/%s", projectId), nil)
	return err
}

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

func (c *Client) DeleteBranch(projectId, branchId string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/projects/%s/branches/%s", projectId, branchId), nil)
	return err
}

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

func (c *Client) GetEndpoint(projectId, endpointId string) (*EndpointState, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/projects/%s/endpoints/%s", projectId, endpointId), nil)
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

func (c *Client) UpdateEndpoint(projectId, endpointId string, branchId string, endpointType string) (*EndpointState, error) {
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

	resp, err := c.doRequest("PATCH", fmt.Sprintf("/projects/%s/endpoints/%s", projectId, endpointId), body)
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

func (c *Client) DeleteEndpoint(projectId, endpointId string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/projects/%s/endpoints/%s", projectId, endpointId), nil)
	return err
}

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
			OwnerName: "default",
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

func (c *Client) GetDatabase(projectId, branchId, databaseName string) (*DatabaseState, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/projects/%s/branches/%s/databases/%s", projectId, branchId, databaseName), nil)
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

func (c *Client) UpdateDatabase(projectId, branchId, databaseName, newName string) (*DatabaseState, error) {
	body := struct {
		Database struct {
			Name string `json:"name"`
		} `json:"database"`
	}{
		Database: struct {
			Name string `json:"name"`
		}{
			Name: newName,
		},
	}

	resp, err := c.doRequest("PATCH", fmt.Sprintf("/projects/%s/branches/%s/databases/%s", projectId, branchId, databaseName), body)
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

func (c *Client) DeleteDatabase(projectId, branchId, databaseName string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/projects/%s/branches/%s/databases/%s", projectId, branchId, databaseName), nil)
	return err
}

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
		Id:        result.Role.Name,
		CreatedAt: result.Role.CreatedAt,
	}, nil
}

func (c *Client) GetRole(projectId, branchId, roleName string) (*RoleState, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/projects/%s/branches/%s/roles/%s", projectId, branchId, roleName), nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Role struct {
			Name      string `json:"name"`
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
		Id:        result.Role.Name,
		CreatedAt: result.Role.CreatedAt,
	}, nil
}

func (c *Client) UpdateRole(projectId, branchId, roleName, newName string) (*RoleState, error) {
	body := struct {
		Role struct {
			Name string `json:"name"`
		} `json:"role"`
	}{
		Role: struct {
			Name string `json:"name"`
		}{
			Name: newName,
		},
	}

	resp, err := c.doRequest("PATCH", fmt.Sprintf("/projects/%s/branches/%s/roles/%s", projectId, branchId, roleName), body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Role struct {
			Name      string `json:"name"`
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
		Id:        result.Role.Name,
		CreatedAt: result.Role.CreatedAt,
	}, nil
}

func (c *Client) DeleteRole(projectId, branchId, roleName string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/projects/%s/branches/%s/roles/%s", projectId, branchId, roleName), nil)
	return err
}

func IsNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "404 Not Found")
}
