package client

import (
	"github.com/kislerdm/neon-sdk-go"
)

type EndpointState struct {
	EndpointArgs
	Id        string
	Host      string
	CreatedAt string
}

type EndpointArgs struct {
	ProjectId string
	BranchId  string
	Type      string
}

func (c *Client) CreateEndpoint(projectId, branchId, endpointType string) (*EndpointState, error) {
	ctx := c.GetContext()
	endpoint, err := c.sdk.Endpoint.Create(ctx, projectId, neon.EndpointCreateRequest{
		Endpoint: neon.EndpointCreateRequestEndpoint{
			BranchID: branchId,
			Type:     endpointType,
		},
	})
	if err != nil {
		return nil, err
	}

	return &EndpointState{
		EndpointArgs: EndpointArgs{
			ProjectId: projectId,
			BranchId:  endpoint.Endpoint.BranchID,
			Type:      endpoint.Endpoint.Type,
		},
		Id:        endpoint.Endpoint.ID,
		Host:      endpoint.Endpoint.Host,
		CreatedAt: endpoint.Endpoint.CreatedAt,
	}, nil
}

func (c *Client) GetEndpoint(projectId, endpointId string) (*EndpointState, error) {
	ctx := c.GetContext()
	endpoint, err := c.sdk.Endpoint.Get(ctx, projectId, endpointId)
	if err != nil {
		return nil, err
	}

	return &EndpointState{
		EndpointArgs: EndpointArgs{
			ProjectId: projectId,
			BranchId:  endpoint.Endpoint.BranchID,
			Type:      endpoint.Endpoint.Type,
		},
		Id:        endpoint.Endpoint.ID,
		Host:      endpoint.Endpoint.Host,
		CreatedAt: endpoint.Endpoint.CreatedAt,
	}, nil
}

func (c *Client) UpdateEndpoint(projectId, endpointId, branchId, endpointType string) (*EndpointState, error) {
	ctx := c.GetContext()
	endpoint, err := c.sdk.Endpoint.Update(ctx, projectId, endpointId, neon.EndpointUpdateRequest{
		Endpoint: neon.EndpointUpdateRequestEndpoint{
			BranchID: branchId,
			Type:     endpointType,
		},
	})
	if err != nil {
		return nil, err
	}

	return &EndpointState{
		EndpointArgs: EndpointArgs{
			ProjectId: projectId,
			BranchId:  endpoint.Endpoint.BranchID,
			Type:      endpoint.Endpoint.Type,
		},
		Id:        endpoint.Endpoint.ID,
		Host:      endpoint.Endpoint.Host,
		CreatedAt: endpoint.Endpoint.CreatedAt,
	}, nil
}

func (c *Client) DeleteEndpoint(projectId, endpointId string) error {
	ctx := c.GetContext()
	_, err := c.sdk.Endpoint.Delete(ctx, projectId, endpointId)
	return err
}