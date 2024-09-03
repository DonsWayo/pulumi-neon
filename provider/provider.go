// Copyright 2016-2023, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"fmt"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string

const Name string = "neon"

var providerLogger = log.New(log.Writer(), "[neon-provider] ", log.LstdFlags|log.Lshortfile)

func Provider() p.Provider {
	// We tell the provider what resources it needs to support.
	return infer.Provider(infer.Options{
		Resources: []infer.InferredResource{
			infer.Resource[Project, ProjectArgs, ProjectState](),
			infer.Resource[Branch, BranchArgs, BranchState](),
			infer.Resource[Endpoint, EndpointArgs, EndpointState](),
			infer.Resource[Database, DatabaseArgs, DatabaseState](),
			infer.Resource[Role, RoleArgs, RoleState](),
		},
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"provider": "index",
		},
		Config: infer.Config[*Config](),
	})
}

type Config struct {
	ApiKey string `pulumi:"apiKey"`
}

func (c *Config) Validate() error {
	if c.ApiKey == "" {
		return fmt.Errorf("apiKey is required")
	}
	return nil
}

// Project resource
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

func (p Project) Create(ctx p.Context, name string, input ProjectArgs, preview bool) (string, ProjectState, error) {
	providerLogger.Printf("Creating Project: %s", name)
	if preview {
		return name, ProjectState{ProjectArgs: input}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", ProjectState{}, fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	project, err := client.CreateProject(input.Name, input.RegionId)
	if err != nil {
		providerLogger.Printf("Failed to create Project: %v", err)
		return "", ProjectState{}, fmt.Errorf("failed to create project: %v", err)
	}

	providerLogger.Printf("Project created successfully: %s", project.Id)
	return name, *project, nil
}

func (p Project) Read(ctx p.Context, id string, inputs ProjectArgs, state ProjectState) (string, ProjectArgs, ProjectState, error) {
	providerLogger.Printf("Reading Project: %s", id)
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", ProjectArgs{}, ProjectState{}, fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	project, err := client.GetProject(state.Id)
	if err != nil {
		if IsNotFoundError(err) {
			providerLogger.Printf("Project not found: %s", id)
			return "", ProjectArgs{}, ProjectState{}, nil
		}
		providerLogger.Printf("Failed to read Project: %v", err)
		return "", ProjectArgs{}, ProjectState{}, fmt.Errorf("failed to read project: %v", err)
	}

	providerLogger.Printf("Project read successfully: %s", project.Id)
	return id, project.ProjectArgs, *project, nil
}

func (p Project) Update(ctx p.Context, id string, olds ProjectState, news ProjectArgs, preview bool) (ProjectState, error) {
	providerLogger.Printf("Updating Project: %s", id)
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

	client := NewClient(config.ApiKey)
	project, err := client.UpdateProject(olds.Id, news.Name)
	if err != nil {
		providerLogger.Printf("Failed to update Project: %v", err)
		return ProjectState{}, fmt.Errorf("failed to update project: %v", err)
	}

	providerLogger.Printf("Project updated successfully: %s", project.Id)
	return *project, nil
}

func (p Project) Delete(ctx p.Context, id string, state ProjectState) error {
	providerLogger.Printf("Deleting Project: %s", id)
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	err := client.DeleteProject(state.Id)
	if err != nil {
		if IsNotFoundError(err) {
			providerLogger.Printf("Project not found during deletion: %s", id)
			return nil
		}
		providerLogger.Printf("Failed to delete Project: %v", err)
		return fmt.Errorf("failed to delete project: %v", err)
	}
	providerLogger.Printf("Project deleted successfully: %s", id)
	return nil
}

// Branch resource
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

func (b Branch) Create(ctx p.Context, name string, input BranchArgs, preview bool) (string, BranchState, error) {
	providerLogger.Printf("Creating Branch: %s", name)
	if preview {
		return name, BranchState{BranchArgs: input}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", BranchState{}, fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	branch, err := client.CreateBranch(input.ProjectId, input.Name)
	if err != nil {
		providerLogger.Printf("Failed to create Branch: %v", err)
		return "", BranchState{}, fmt.Errorf("failed to create branch: %v", err)
	}

	providerLogger.Printf("Branch created successfully: %s", branch.Id)
	return name, *branch, nil
}

func (b Branch) Read(ctx p.Context, id string, inputs BranchArgs, state BranchState) (string, BranchArgs, BranchState, error) {
	providerLogger.Printf("Reading Branch: %s", id)
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", BranchArgs{}, BranchState{}, fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	branch, err := client.GetBranch(state.ProjectId, state.Id)
	if err != nil {
		if IsNotFoundError(err) {
			providerLogger.Printf("Branch not found: %s", id)
			return "", BranchArgs{}, BranchState{}, nil
		}
		providerLogger.Printf("Failed to read Branch: %v", err)
		return "", BranchArgs{}, BranchState{}, fmt.Errorf("failed to read branch: %v", err)
	}

	providerLogger.Printf("Branch read successfully: %s", branch.Id)
	return id, branch.BranchArgs, *branch, nil
}

func (b Branch) Update(ctx p.Context, id string, olds BranchState, news BranchArgs, preview bool) (BranchState, error) {
	providerLogger.Printf("Updating Branch: %s", id)
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

	client := NewClient(config.ApiKey)
	branch, err := client.UpdateBranch(news.ProjectId, olds.Id, news.Name)
	if err != nil {
		providerLogger.Printf("Failed to update Branch: %v", err)
		return BranchState{}, fmt.Errorf("failed to update branch: %v", err)
	}

	providerLogger.Printf("Branch updated successfully: %s", branch.Id)
	return *branch, nil
}

func (b Branch) Delete(ctx p.Context, id string, state BranchState) error {
	providerLogger.Printf("Deleting Branch: %s", id)
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	err := client.DeleteBranch(state.ProjectId, state.Id)
	if err != nil {
		if IsNotFoundError(err) {
			providerLogger.Printf("Branch not found during deletion: %s", id)
			return nil
		}
		providerLogger.Printf("Failed to delete Branch: %v", err)
		return fmt.Errorf("failed to delete branch: %v", err)
	}
	providerLogger.Printf("Branch deleted successfully: %s", id)
	return nil
}

// Endpoint resource
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

func (e Endpoint) Create(ctx p.Context, name string, input EndpointArgs, preview bool) (string, EndpointState, error) {
	providerLogger.Printf("Creating Endpoint: %s", name)
	if preview {
		return name, EndpointState{EndpointArgs: input}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", EndpointState{}, fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	endpoint, err := client.CreateEndpoint(input.ProjectId, input.BranchId, input.Type)
	if err != nil {
		providerLogger.Printf("Failed to create Endpoint: %v", err)
		return "", EndpointState{}, fmt.Errorf("failed to create endpoint: %v", err)
	}

	providerLogger.Printf("Endpoint created successfully: %s", endpoint.Id)
	return name, *endpoint, nil
}

func (e Endpoint) Read(ctx p.Context, id string, inputs EndpointArgs, state EndpointState) (string, EndpointArgs, EndpointState, error) {
	providerLogger.Printf("Reading Endpoint: %s", id)
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", EndpointArgs{}, EndpointState{}, fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	endpoint, err := client.GetEndpoint(state.ProjectId, state.Id)
	if err != nil {
		if IsNotFoundError(err) {
			providerLogger.Printf("Endpoint not found: %s", id)
			return "", EndpointArgs{}, EndpointState{}, nil
		}
		providerLogger.Printf("Failed to read Endpoint: %v", err)
		return "", EndpointArgs{}, EndpointState{}, fmt.Errorf("failed to read endpoint: %v", err)
	}

	providerLogger.Printf("Endpoint read successfully: %s", endpoint.Id)
	return id, endpoint.EndpointArgs, *endpoint, nil
}

func (e Endpoint) Update(ctx p.Context, id string, olds EndpointState, news EndpointArgs, preview bool) (EndpointState, error) {
	providerLogger.Printf("Updating Endpoint: %s", id)
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

	client := NewClient(config.ApiKey)
	endpoint, err := client.UpdateEndpoint(news.ProjectId, olds.Id, news.BranchId, news.Type)
	if err != nil {
		providerLogger.Printf("Failed to update Endpoint: %v", err)
		return EndpointState{}, fmt.Errorf("failed to update endpoint: %v", err)
	}

	providerLogger.Printf("Endpoint updated successfully: %s", endpoint.Id)
	return *endpoint, nil
}

func (e Endpoint) Delete(ctx p.Context, id string, state EndpointState) error {
	providerLogger.Printf("Deleting Endpoint: %s", id)
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	err := client.DeleteEndpoint(state.ProjectId, state.Id)
	if err != nil {
		if IsNotFoundError(err) {
			providerLogger.Printf("Endpoint not found during deletion: %s", id)
			return nil
		}
		providerLogger.Printf("Failed to delete Endpoint: %v", err)
		return fmt.Errorf("failed to delete endpoint: %v", err)
	}
	providerLogger.Printf("Endpoint deleted successfully: %s", id)
	return nil
}

// Database resource
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

func (d Database) Create(ctx p.Context, name string, input DatabaseArgs, preview bool) (string, DatabaseState, error) {
	providerLogger.Printf("Creating Database: %s", name)
	if preview {
		return name, DatabaseState{DatabaseArgs: input}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", DatabaseState{}, fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	database, err := client.CreateDatabase(input.ProjectId, input.BranchId, input.Name)
	if err != nil {
		providerLogger.Printf("Failed to create Database: %v", err)
		return "", DatabaseState{}, fmt.Errorf("failed to create database: %v", err)
	}

	providerLogger.Printf("Database created successfully: %s", database.Id)
	return name, *database, nil
}

func (d Database) Read(ctx p.Context, id string, inputs DatabaseArgs, state DatabaseState) (string, DatabaseArgs, DatabaseState, error) {
	providerLogger.Printf("Reading Database: %s", id)
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", DatabaseArgs{}, DatabaseState{}, fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	database, err := client.GetDatabase(state.ProjectId, state.BranchId, state.Name)
	if err != nil {
		if IsNotFoundError(err) {
			providerLogger.Printf("Database not found: %s", id)
			return "", DatabaseArgs{}, DatabaseState{}, nil
		}
		providerLogger.Printf("Failed to read Database: %v", err)
		return "", DatabaseArgs{}, DatabaseState{}, fmt.Errorf("failed to read database: %v", err)
	}

	providerLogger.Printf("Database read successfully: %s", database.Id)
	return id, database.DatabaseArgs, *database, nil
}

func (d Database) Update(ctx p.Context, id string, olds DatabaseState, news DatabaseArgs, preview bool) (DatabaseState, error) {
	providerLogger.Printf("Updating Database: %s", id)
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

	client := NewClient(config.ApiKey)
	database, err := client.UpdateDatabase(news.ProjectId, news.BranchId, olds.Name, news.Name)
	if err != nil {
		providerLogger.Printf("Failed to update Database: %v", err)
		return DatabaseState{}, fmt.Errorf("failed to update database: %v", err)
	}

	providerLogger.Printf("Database updated successfully: %s", database.Id)
	return *database, nil
}

func (d Database) Delete(ctx p.Context, id string, state DatabaseState) error {
	providerLogger.Printf("Deleting Database: %s", id)
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	err := client.DeleteDatabase(state.ProjectId, state.BranchId, state.Name)
	if err != nil {
		if IsNotFoundError(err) {
			providerLogger.Printf("Database not found during deletion: %s", id)
			return nil
		}
		providerLogger.Printf("Failed to delete Database: %v", err)
		return fmt.Errorf("failed to delete database: %v", err)
	}
	providerLogger.Printf("Database deleted successfully: %s", id)
	return nil
}

// Role resource
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

func (r Role) Create(ctx p.Context, name string, input RoleArgs, preview bool) (string, RoleState, error) {
	providerLogger.Printf("Creating Role: %s", name)
	if preview {
		return name, RoleState{RoleArgs: input}, nil
	}

	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", RoleState{}, fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	role, err := client.CreateRole(input.ProjectId, input.BranchId, input.Name)
	if err != nil {
		providerLogger.Printf("Failed to create Role: %v", err)
		return "", RoleState{}, fmt.Errorf("failed to create role: %v", err)
	}

	providerLogger.Printf("Role created successfully: %s", role.Id)
	return name, *role, nil
}

func (r Role) Read(ctx p.Context, id string, inputs RoleArgs, state RoleState) (string, RoleArgs, RoleState, error) {
	providerLogger.Printf("Reading Role: %s", id)
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", RoleArgs{}, RoleState{}, fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	role, err := client.GetRole(state.ProjectId, state.BranchId, state.Name)
	if err != nil {
		if IsNotFoundError(err) {
			providerLogger.Printf("Role not found: %s", id)
			return "", RoleArgs{}, RoleState{}, nil
		}
		providerLogger.Printf("Failed to read Role: %v", err)
		return "", RoleArgs{}, RoleState{}, fmt.Errorf("failed to read role: %v", err)
	}

	providerLogger.Printf("Role read successfully: %s", role.Id)
	return id, role.RoleArgs, *role, nil
}

func (r Role) Update(ctx p.Context, id string, olds RoleState, news RoleArgs, preview bool) (RoleState, error) {
	providerLogger.Printf("Updating Role: %s", id)
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

	client := NewClient(config.ApiKey)
	role, err := client.UpdateRole(news.ProjectId, news.BranchId, olds.Name, news.Name)
	if err != nil {
		providerLogger.Printf("Failed to update Role: %v", err)
		return RoleState{}, fmt.Errorf("failed to update role: %v", err)
	}

	providerLogger.Printf("Role updated successfully: %s", role.Id)
	return *role, nil
}

func (r Role) Delete(ctx p.Context, id string, state RoleState) error {
	providerLogger.Printf("Deleting Role: %s", id)
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return fmt.Errorf("missing configuration")
	}

	client := NewClient(config.ApiKey)
	err := client.DeleteRole(state.ProjectId, state.BranchId, state.Name)
	if err != nil {
		if IsNotFoundError(err) {
			providerLogger.Printf("Role not found during deletion: %s", id)
			return nil
		}
		providerLogger.Printf("Failed to delete Role: %v", err)
		return fmt.Errorf("failed to delete role: %v", err)
	}
	providerLogger.Printf("Role deleted successfully: %s", id)
	return nil
}

// Helper function to check if an error is a "not found" error
func IsNotFoundError(err error) bool {
	// Implement the logic to determine if the error is a "not found" error
	// This will depend on how the Neon API indicates resource not found errors
	// For example:
	// return strings.Contains(err.Error(), "not found")
	return false
}
