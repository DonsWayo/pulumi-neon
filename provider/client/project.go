package client

import (
	"github.com/kislerdm/neon-sdk-go"
)

type ProjectState struct {
	ProjectArgs
	Id        string
	CreatedAt string
}

type ProjectArgs struct {
	Name     string
	RegionId string
}

func (c *Client) CreateProject(name, regionId string) (*ProjectState, error) {
	ctx := c.GetContext()
	project, err := c.sdk.Project.Create(ctx, neon.ProjectCreateRequestV2{
		Project: neon.ProjectCreateRequestV2Project{
			Name:     name,
			RegionId: regionId,
		},
	})
	if err != nil {
		return nil, err
	}

	return &ProjectState{
		ProjectArgs: ProjectArgs{
			Name:     project.Project.Name,
			RegionId: project.Project.RegionId,
		},
		Id:        project.Project.ID,
		CreatedAt: project.Project.CreatedAt,
	}, nil
}

func (c *Client) GetProject(projectId string) (*ProjectState, error) {
	ctx := c.GetContext()
	project, err := c.sdk.Project.Get(ctx, projectId)
	if err != nil {
		return nil, err
	}

	return &ProjectState{
		ProjectArgs: ProjectArgs{
			Name:     project.Project.Name,
			RegionId: project.Project.RegionId,
		},
		Id:        project.Project.ID,
		CreatedAt: project.Project.CreatedAt,
	}, nil
}

func (c *Client) UpdateProject(projectId string, name string) (*ProjectState, error) {
	ctx := c.GetContext()
	project, err := c.sdk.Project.Update(ctx, projectId, neon.ProjectUpdateRequest{
		Project: neon.ProjectUpdateRequestProject{
			Name: name,
		},
	})
	if err != nil {
		return nil, err
	}

	return &ProjectState{
		ProjectArgs: ProjectArgs{
			Name:     project.Project.Name,
			RegionId: project.Project.RegionId,
		},
		Id:        project.Project.ID,
		CreatedAt: project.Project.CreatedAt,
	}, nil
}

func (c *Client) DeleteProject(projectId string) error {
	ctx := c.GetContext()
	_, err := c.sdk.Project.Delete(ctx, projectId)
	return err
}