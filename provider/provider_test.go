package provider

import (
	"testing"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockContext struct {
	p.Context
	config *Config
}

func (m *mockContext) GetConfig() (interface{}, bool) {
	return m.config, true
}

// Mock client for testing
type MockClient struct {
	mock.Mock
	*Client
}

func (m *MockClient) CreateProject(name, regionId string) (*ProjectState, error) {
	args := m.Called(name, regionId)
	return args.Get(0).(*ProjectState), args.Error(1)
}

func (m *MockClient) GetProject(projectId string) (*ProjectState, error) {
	args := m.Called(projectId)
	return args.Get(0).(*ProjectState), args.Error(1)
}

func (m *MockClient) UpdateProject(projectId, name string) (*ProjectState, error) {
	args := m.Called(projectId, name)
	return args.Get(0).(*ProjectState), args.Error(1)
}

func (m *MockClient) DeleteProject(projectId string) error {
	args := m.Called(projectId)
	return args.Error(0)
}

func (m *MockClient) CreateBranch(projectId, name string) (*BranchState, error) {
	args := m.Called(projectId, name)
	return args.Get(0).(*BranchState), args.Error(1)
}

func (m *MockClient) GetBranch(projectId, branchId string) (*BranchState, error) {
	args := m.Called(projectId, branchId)
	return args.Get(0).(*BranchState), args.Error(1)
}

func (m *MockClient) UpdateBranch(projectId, branchId, name string) (*BranchState, error) {
	args := m.Called(projectId, branchId, name)
	return args.Get(0).(*BranchState), args.Error(1)
}

func (m *MockClient) DeleteBranch(projectId, branchId string) error {
	args := m.Called(projectId, branchId)
	return args.Error(0)
}

func (m *MockClient) CreateEndpoint(projectId, branchId, endpointType string) (*EndpointState, error) {
	args := m.Called(projectId, branchId, endpointType)
	return args.Get(0).(*EndpointState), args.Error(1)
}

func (m *MockClient) GetEndpoint(projectId, endpointId string) (*EndpointState, error) {
	args := m.Called(projectId, endpointId)
	return args.Get(0).(*EndpointState), args.Error(1)
}

func (m *MockClient) UpdateEndpoint(projectId, endpointId, branchId, endpointType string) (*EndpointState, error) {
	args := m.Called(projectId, endpointId, branchId, endpointType)
	return args.Get(0).(*EndpointState), args.Error(1)
}

func (m *MockClient) DeleteEndpoint(projectId, endpointId string) error {
	args := m.Called(projectId, endpointId)
	return args.Error(0)
}

func (m *MockClient) CreateDatabase(projectId, branchId, name string) (*DatabaseState, error) {
	args := m.Called(projectId, branchId, name)
	return args.Get(0).(*DatabaseState), args.Error(1)
}

func (m *MockClient) GetDatabase(projectId, branchId, name string) (*DatabaseState, error) {
	args := m.Called(projectId, branchId, name)
	return args.Get(0).(*DatabaseState), args.Error(1)
}

func (m *MockClient) UpdateDatabase(projectId, branchId, oldName, newName string) (*DatabaseState, error) {
	args := m.Called(projectId, branchId, oldName, newName)
	return args.Get(0).(*DatabaseState), args.Error(1)
}

func (m *MockClient) DeleteDatabase(projectId, branchId, name string) error {
	args := m.Called(projectId, branchId, name)
	return args.Error(0)
}

func TestProjectCreate(t *testing.T) {
	p := Project{}
	name := "test-project"
	input := ProjectArgs{
		Name:     "Test Project",
		RegionId: "us-east-1",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("CreateProject", input.Name, input.RegionId).Return(&ProjectState{
		ProjectArgs: input,
		Id:          "test-project-id",
		CreatedAt:   "2023-05-01T00:00:00Z",
	}, nil)

	// Call the Create method
	id, state, err := p.Create(ctx, name, input, false)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, name, id)
	assert.Equal(t, input.Name, state.Name)
	assert.Equal(t, input.RegionId, state.RegionId)
	assert.Equal(t, "test-project-id", state.Id)
	assert.Equal(t, "2023-05-01T00:00:00Z", state.CreatedAt)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestProjectRead(t *testing.T) {
	p := Project{}
	id := "test-project"
	input := ProjectArgs{
		Name:     "Test Project",
		RegionId: "us-east-1",
	}
	state := ProjectState{
		ProjectArgs: input,
		Id:          "test-project-id",
		CreatedAt:   "2023-05-01T00:00:00Z",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("GetProject", state.Id).Return(&state, nil)

	// Call the Read method
	readId, readInput, readState, err := p.Read(ctx, id, input, state)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, id, readId)
	assert.Equal(t, input, readInput)
	assert.Equal(t, state, readState)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestProjectUpdate(t *testing.T) {
	p := Project{}
	id := "test-project"
	olds := ProjectState{
		ProjectArgs: ProjectArgs{
			Name:     "Old Project",
			RegionId: "us-east-1",
		},
		Id:        "test-project-id",
		CreatedAt: "2023-05-01T00:00:00Z",
	}
	news := ProjectArgs{
		Name:     "New Project",
		RegionId: "us-east-1",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("UpdateProject", olds.Id, news.Name).Return(&ProjectState{
		ProjectArgs: news,
		Id:          olds.Id,
		CreatedAt:   olds.CreatedAt,
	}, nil)

	// Call the Update method
	updatedState, err := p.Update(ctx, id, olds, news, false)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, news.Name, updatedState.Name)
	assert.Equal(t, news.RegionId, updatedState.RegionId)
	assert.Equal(t, olds.Id, updatedState.Id)
	assert.Equal(t, olds.CreatedAt, updatedState.CreatedAt)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestProjectDelete(t *testing.T) {
	p := Project{}
	id := "test-project"
	state := ProjectState{
		ProjectArgs: ProjectArgs{
			Name:     "Test Project",
			RegionId: "us-east-1",
		},
		Id:        "test-project-id",
		CreatedAt: "2023-05-01T00:00:00Z",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("DeleteProject", state.Id).Return(nil)

	// Call the Delete method
	err := p.Delete(ctx, id, state)

	// Assert the results
	assert.NoError(t, err)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestBranchCreate(t *testing.T) {
	b := Branch{}
	name := "test-branch"
	input := BranchArgs{
		ProjectId: "test-project-id",
		Name:      "Test Branch",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("CreateBranch", input.ProjectId, input.Name).Return(&BranchState{
		BranchArgs: input,
		Id:         "test-branch-id",
		CreatedAt:  "2023-05-01T00:00:00Z",
	}, nil)

	// Call the Create method
	id, state, err := b.Create(ctx, name, input, false)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, name, id)
	assert.Equal(t, input.ProjectId, state.ProjectId)
	assert.Equal(t, input.Name, state.Name)
	assert.Equal(t, "test-branch-id", state.Id)
	assert.Equal(t, "2023-05-01T00:00:00Z", state.CreatedAt)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestBranchRead(t *testing.T) {
	b := Branch{}
	id := "test-branch"
	input := BranchArgs{
		ProjectId: "test-project-id",
		Name:      "Test Branch",
	}
	state := BranchState{
		BranchArgs: input,
		Id:         "test-branch-id",
		CreatedAt:  "2023-05-01T00:00:00Z",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("GetBranch", state.ProjectId, state.Id).Return(&state, nil)

	// Call the Read method
	readId, readInput, readState, err := b.Read(ctx, id, input, state)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, id, readId)
	assert.Equal(t, input, readInput)
	assert.Equal(t, state, readState)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestBranchUpdate(t *testing.T) {
	b := Branch{}
	id := "test-branch"
	olds := BranchState{
		BranchArgs: BranchArgs{
			ProjectId: "test-project-id",
			Name:      "Old Branch",
		},
		Id:        "test-branch-id",
		CreatedAt: "2023-05-01T00:00:00Z",
	}
	news := BranchArgs{
		ProjectId: "test-project-id",
		Name:      "New Branch",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("UpdateBranch", news.ProjectId, olds.Id, news.Name).Return(&BranchState{
		BranchArgs: news,
		Id:         olds.Id,
		CreatedAt:  olds.CreatedAt,
	}, nil)

	// Call the Update method
	updatedState, err := b.Update(ctx, id, olds, news, false)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, news.ProjectId, updatedState.ProjectId)
	assert.Equal(t, news.Name, updatedState.Name)
	assert.Equal(t, olds.Id, updatedState.Id)
	assert.Equal(t, olds.CreatedAt, updatedState.CreatedAt)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestBranchDelete(t *testing.T) {
	b := Branch{}
	id := "test-branch"
	state := BranchState{
		BranchArgs: BranchArgs{
			ProjectId: "test-project-id",
			Name:      "Test Branch",
		},
		Id:        "test-branch-id",
		CreatedAt: "2023-05-01T00:00:00Z",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("DeleteBranch", state.ProjectId, state.Id).Return(nil)

	// Call the Delete method
	err := b.Delete(ctx, id, state)

	// Assert the results
	assert.NoError(t, err)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestEndpointCreate(t *testing.T) {
	e := Endpoint{}
	name := "test-endpoint"
	input := EndpointArgs{
		ProjectId: "test-project-id",
		BranchId:  "test-branch-id",
		Type:      "read_write",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("CreateEndpoint", input.ProjectId, input.BranchId, input.Type).Return(&EndpointState{
		EndpointArgs: input,
		Id:           "test-endpoint-id",
		Host:         "test-endpoint-host",
		CreatedAt:    "2023-05-01T00:00:00Z",
	}, nil)

	// Call the Create method
	id, state, err := e.Create(ctx, name, input, false)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, name, id)
	assert.Equal(t, input.ProjectId, state.ProjectId)
	assert.Equal(t, input.BranchId, state.BranchId)
	assert.Equal(t, input.Type, state.Type)
	assert.Equal(t, "test-endpoint-id", state.Id)
	assert.Equal(t, "test-endpoint-host", state.Host)
	assert.Equal(t, "2023-05-01T00:00:00Z", state.CreatedAt)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestEndpointRead(t *testing.T) {
	e := Endpoint{}
	id := "test-endpoint"
	input := EndpointArgs{
		ProjectId: "test-project-id",
		BranchId:  "test-branch-id",
		Type:      "read_write",
	}
	state := EndpointState{
		EndpointArgs: input,
		Id:           "test-endpoint-id",
		Host:         "test-endpoint-host",
		CreatedAt:    "2023-05-01T00:00:00Z",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("GetEndpoint", state.ProjectId, state.Id).Return(&state, nil)

	// Call the Read method
	readId, readInput, readState, err := e.Read(ctx, id, input, state)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, id, readId)
	assert.Equal(t, input, readInput)
	assert.Equal(t, state, readState)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestEndpointUpdate(t *testing.T) {
	e := Endpoint{}
	id := "test-endpoint"
	olds := EndpointState{
		EndpointArgs: EndpointArgs{
			ProjectId: "test-project-id",
			BranchId:  "old-branch-id",
			Type:      "read_only",
		},
		Id:        "test-endpoint-id",
		Host:      "test-endpoint-host",
		CreatedAt: "2023-05-01T00:00:00Z",
	}
	news := EndpointArgs{
		ProjectId: "test-project-id",
		BranchId:  "new-branch-id",
		Type:      "read_write",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("UpdateEndpoint", news.ProjectId, olds.Id, news.BranchId, news.Type).Return(&EndpointState{
		EndpointArgs: news,
		Id:           olds.Id,
		Host:         olds.Host,
		CreatedAt:    olds.CreatedAt,
	}, nil)

	// Call the Update method
	updatedState, err := e.Update(ctx, id, olds, news, false)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, news.ProjectId, updatedState.ProjectId)
	assert.Equal(t, news.BranchId, updatedState.BranchId)
	assert.Equal(t, news.Type, updatedState.Type)
	assert.Equal(t, olds.Id, updatedState.Id)
	assert.Equal(t, olds.Host, updatedState.Host)
	assert.Equal(t, olds.CreatedAt, updatedState.CreatedAt)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestEndpointDelete(t *testing.T) {
	e := Endpoint{}
	id := "test-endpoint"
	state := EndpointState{
		EndpointArgs: EndpointArgs{
			ProjectId: "test-project-id",
			BranchId:  "test-branch-id",
			Type:      "read_write",
		},
		Id:        "test-endpoint-id",
		Host:      "test-endpoint-host",
		CreatedAt: "2023-05-01T00:00:00Z",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("DeleteEndpoint", state.ProjectId, state.Id).Return(nil)

	// Call the Delete method
	err := e.Delete(ctx, id, state)

	// Assert the results
	assert.NoError(t, err)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestDatabaseCreate(t *testing.T) {
	d := Database{}
	name := "test-database"
	input := DatabaseArgs{
		ProjectId: "test-project-id",
		BranchId:  "test-branch-id",
		Name:      "TestDatabase",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("CreateDatabase", input.ProjectId, input.BranchId, input.Name).Return(&DatabaseState{
		DatabaseArgs: input,
		Id:           "test-database-id",
		CreatedAt:    "2023-05-01T00:00:00Z",
	}, nil)

	// Call the Create method
	id, state, err := d.Create(ctx, name, input, false)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, name, id)
	assert.Equal(t, input.ProjectId, state.ProjectId)
	assert.Equal(t, input.BranchId, state.BranchId)
	assert.Equal(t, input.Name, state.Name)
	assert.Equal(t, "test-database-id", state.Id)
	assert.Equal(t, "2023-05-01T00:00:00Z", state.CreatedAt)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestDatabaseRead(t *testing.T) {
	d := Database{}
	id := "test-database"
	input := DatabaseArgs{
		ProjectId: "test-project-id",
		BranchId:  "test-branch-id",
		Name:      "TestDatabase",
	}
	state := DatabaseState{
		DatabaseArgs: input,
		Id:           "test-database-id",
		CreatedAt:    "2023-05-01T00:00:00Z",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("GetDatabase", state.ProjectId, state.BranchId, state.Name).Return(&state, nil)

	// Call the Read method
	readId, readInput, readState, err := d.Read(ctx, id, input, state)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, id, readId)
	assert.Equal(t, input, readInput)
	assert.Equal(t, state, readState)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestDatabaseUpdate(t *testing.T) {
	d := Database{}
	id := "test-database"
	olds := DatabaseState{
		DatabaseArgs: DatabaseArgs{
			ProjectId: "test-project-id",
			BranchId:  "test-branch-id",
			Name:      "OldDatabase",
		},
		Id:        "test-database-id",
		CreatedAt: "2023-05-01T00:00:00Z",
	}
	news := DatabaseArgs{
		ProjectId: "test-project-id",
		BranchId:  "test-branch-id",
		Name:      "NewDatabase",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("UpdateDatabase", news.ProjectId, news.BranchId, olds.Name, news.Name).Return(&DatabaseState{
		DatabaseArgs: news,
		Id:           olds.Id,
		CreatedAt:    olds.CreatedAt,
	}, nil)

	// Call the Update method
	updatedState, err := d.Update(ctx, id, olds, news, false)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, news.ProjectId, updatedState.ProjectId)
	assert.Equal(t, news.BranchId, updatedState.BranchId)
	assert.Equal(t, news.Name, updatedState.Name)
	assert.Equal(t, olds.Id, updatedState.Id)
	assert.Equal(t, olds.CreatedAt, updatedState.CreatedAt)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestDatabaseDelete(t *testing.T) {
	d := Database{}
	id := "test-database"
	state := DatabaseState{
		DatabaseArgs: DatabaseArgs{
			ProjectId: "test-project-id",
			BranchId:  "test-branch-id",
			Name:      "TestDatabase",
		},
		Id:        "test-database-id",
		CreatedAt: "2023-05-01T00:00:00Z",
	}

	// Mock the context
	ctx := &mockContext{
		config: &Config{ApiKey: "test-api-key"},
	}

	// Create a mock client
	mockClient := new(MockClient)
	originalNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient.Client
	}
	defer func() { NewClient = originalNewClient }()

	// Set expectations
	mockClient.On("DeleteDatabase", state.ProjectId, state.BranchId, state.Name).Return(nil)

	// Call the Delete method
	err := d.Delete(ctx, id, state)

	// Assert the results
	assert.NoError(t, err)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

// Add more test methods for other resources and operations
