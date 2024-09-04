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
	mockClient := &MockClient{}
	oldNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient
	}
	defer func() { NewClient = oldNewClient }()

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
	mockClient := &MockClient{}
	oldNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient
	}
	defer func() { NewClient = oldNewClient }()

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
	mockClient := &MockClient{}
	oldNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient
	}
	defer func() { NewClient = oldNewClient }()

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
	mockClient := &MockClient{}
	oldNewClient := NewClient
	NewClient = func(apiKey string) *Client {
		return mockClient
	}
	defer func() { NewClient = oldNewClient }()

	// Set expectations
	mockClient.On("DeleteProject", state.Id).Return(nil)

	// Call the Delete method
	err := p.Delete(ctx, id, state)

	// Assert the results
	assert.NoError(t, err)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

// Mock client for testing
type MockClient struct {
	mock.Mock
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

// Add more test methods for other resources and operations
