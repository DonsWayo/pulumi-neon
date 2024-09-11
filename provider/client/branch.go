package client

import (
	"log"

	"github.com/kislerdm/neon-sdk-go"
)

type BranchState struct {
	BranchArgs
	Id        string
	CreatedAt string
}

type BranchArgs struct {
	ProjectId string
	Name      string
}

func (c *Client) CreateBranch(projectId, name string) (*BranchState, error) {
	ctx := c.GetContext()
	log.Printf("CreateBranch: Starting with projectId=%s, name=%s", projectId, name)

	branch, err := c.sdk.Branch.Create(ctx, projectId, neon.BranchCreateRequest{
		Branch: neon.BranchCreateRequestBranch{
			Name: name,
		},
		Endpoints: []neon.BranchCreateRequestEndpoint{
			{Type: "read_only"},
		},
	})
	if err != nil {
		log.Printf("CreateBranch: Error occurred: %v", err)
		return nil, err
	}

	log.Printf("CreateBranch: Branch created successfully: id=%s", branch.Branch.ID)
	return &BranchState{
		BranchArgs: BranchArgs{
			ProjectId: projectId,
			Name:      branch.Branch.Name,
		},
		Id:        branch.Branch.ID,
		CreatedAt: branch.Branch.CreatedAt,
	}, nil
}

func (c *Client) GetBranch(projectId, branchId string) (*BranchState, error) {
	ctx := c.GetContext()
	branch, err := c.sdk.Branch.Get(ctx, projectId, branchId)
	if err != nil {
		return nil, err
	}

	return &BranchState{
		BranchArgs: BranchArgs{
			ProjectId: projectId,
			Name:      branch.Branch.Name,
		},
		Id:        branch.Branch.ID,
		CreatedAt: branch.Branch.CreatedAt,
	}, nil
}

func (c *Client) UpdateBranch(projectId, branchId, name string) (*BranchState, error) {
	ctx := c.GetContext()
	branch, err := c.sdk.Branch.Update(ctx, projectId, branchId, neon.BranchUpdateRequest{
		Branch: neon.BranchUpdateRequestBranch{
			Name: name,
		},
	})
	if err != nil {
		return nil, err
	}

	return &BranchState{
		BranchArgs: BranchArgs{
			ProjectId: projectId,
			Name:      branch.Branch.Name,
		},
		Id:        branch.Branch.ID,
		CreatedAt: branch.Branch.CreatedAt,
	}, nil
}

func (c *Client) DeleteBranch(projectId, branchId string) error {
	ctx := c.GetContext()
	_, err := c.sdk.Branch.Delete(ctx, projectId, branchId)
	return err
}