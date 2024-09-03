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
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string

const Name string = "neon"

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
	if preview {
		return name, ProjectState{ProjectArgs: input}, nil
	}

	client := NewClient(ctx.GetConfig().(*Config).ApiKey)
	project, err := client.CreateProject(input.Name, input.RegionId)
	if err != nil {
		return "", ProjectState{}, err
	}

	return name, *project, nil
}

func (p Project) Read(ctx p.Context, id string, inputs ProjectArgs, state ProjectState) (string, ProjectArgs, ProjectState, error) {
	client := NewClient(ctx.GetConfig().(*Config).ApiKey)
	project, err := client.GetProject(state.Id)
	if err != nil {
		return "", ProjectArgs{}, ProjectState{}, err
	}

	return id, project.ProjectArgs, *project, nil
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
	if preview {
		return name, BranchState{BranchArgs: input}, nil
	}

	client := NewClient(ctx.GetConfig().(*Config).ApiKey)
	branch, err := client.CreateBranch(input.ProjectId, input.Name)
	if err != nil {
		return "", BranchState{}, err
	}

	return name, *branch, nil
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
	if preview {
		return name, EndpointState{EndpointArgs: input}, nil
	}

	client := NewClient(ctx.GetConfig().(*Config).ApiKey)
	endpoint, err := client.CreateEndpoint(input.ProjectId, input.BranchId, input.Type)
	if err != nil {
		return "", EndpointState{}, err
	}

	return name, *endpoint, nil
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
	if preview {
		return name, DatabaseState{DatabaseArgs: input}, nil
	}

	client := NewClient(ctx.GetConfig().(*Config).ApiKey)
	database, err := client.CreateDatabase(input.ProjectId, input.BranchId, input.Name)
	if err != nil {
		return "", DatabaseState{}, err
	}

	return name, *database, nil
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
	if preview {
		return name, RoleState{RoleArgs: input}, nil
	}

	client := NewClient(ctx.GetConfig().(*Config).ApiKey)
	role, err := client.CreateRole(input.ProjectId, input.BranchId, input.Name)
	if err != nil {
		return "", RoleState{}, err
	}

	return name, *role, nil
}
