package client

import (
	"github.com/kislerdm/neon-sdk-go"
)

type RoleState struct {
	RoleArgs
	Id        string
	CreatedAt string
}

type RoleArgs struct {
	ProjectId string
	BranchId  string
	Name      string
}

func (c *Client) CreateRole(projectId, branchId, name string) (*RoleState, error) {
	ctx := c.GetContext()
	role, err := c.sdk.Role.Create(ctx, projectId, branchId, neon.RoleCreateRequest{
		Role: neon.RoleCreateRequestRole{
			Name: name,
		},
	})
	if err != nil {
		return nil, err
	}

	return &RoleState{
		RoleArgs: RoleArgs{
			ProjectId: projectId,
			BranchId:  branchId,
			Name:      role.Role.Name,
		},
		Id:        role.Role.Name,
		CreatedAt: role.Role.CreatedAt,
	}, nil
}

func (c *Client) GetRole(projectId, branchId, roleName string) (*RoleState, error) {
	ctx := c.GetContext()
	role, err := c.sdk.Role.Get(ctx, projectId, branchId, roleName)
	if err != nil {
		return nil, err
	}

	return &RoleState{
		RoleArgs: RoleArgs{
			ProjectId: projectId,
			BranchId:  branchId,
			Name:      role.Role.Name,
		},
		Id:        role.Role.Name,
		CreatedAt: role.Role.CreatedAt,
	}, nil
}

func (c *Client) UpdateRole(projectId, branchId, roleName, newName string) (*RoleState, error) {
	ctx := c.GetContext()
	role, err := c.sdk.Role.Update(ctx, projectId, branchId, roleName, neon.RoleUpdateRequest{
		Role: neon.RoleUpdateRequestRole{
			Name: newName,
		},
	})
	if err != nil {
		return nil, err
	}

	return &RoleState{
		RoleArgs: RoleArgs{
			ProjectId: projectId,
			BranchId:  branchId,
			Name:      role.Role.Name,
		},
		Id:        role.Role.Name,
		CreatedAt: role.Role.CreatedAt,
	}, nil
}

func (c *Client) DeleteRole(projectId, branchId, roleName string) error {
	ctx := c.GetContext()
	_, err := c.sdk.Role.Delete(ctx, projectId, branchId, roleName)
	return err
}