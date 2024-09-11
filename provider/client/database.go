package client

import (
	"fmt"
	"log"

	"github.com/kislerdm/neon-sdk-go"
)

type DatabaseState struct {
	DatabaseArgs
	Id        string
	CreatedAt string
}

type DatabaseArgs struct {
	ProjectId string
	BranchId  string
	Name      string
}

func (c *Client) CreateDatabase(projectId, branchId, name string) (*DatabaseState, error) {
	ctx := c.GetContext()
	log.Printf("Creating database: projectId=%s, branchId=%s, name=%s", projectId, branchId, name)

	database, err := c.sdk.Database.Create(ctx, projectId, branchId, neon.DatabaseCreateRequest{
		Database: neon.DatabaseCreateRequestDatabase{
			Name:      name,
			OwnerName: "default",
		},
	})
	if err != nil {
		log.Printf("Error creating database: %v", err)
		return nil, err
	}

	log.Printf("Database created successfully: id=%s", database.Database.ID)
	return &DatabaseState{
		DatabaseArgs: DatabaseArgs{
			ProjectId: projectId,
			BranchId:  branchId,
			Name:      database.Database.Name,
		},
		Id:        database.Database.ID,
		CreatedAt: database.Database.CreatedAt,
	}, nil
}

func (c *Client) GetDatabase(projectId, branchId, databaseName string) (*DatabaseState, error) {
	ctx := c.GetContext()
	database, err := c.sdk.Database.Get(ctx, projectId, branchId, databaseName)
	if err != nil {
		return nil, err
	}

	return &DatabaseState{
		DatabaseArgs: DatabaseArgs{
			ProjectId: projectId,
			BranchId:  branchId,
			Name:      database.Database.Name,
		},
		Id:        database.Database.ID,
		CreatedAt: database.Database.CreatedAt,
	}, nil
}

func (c *Client) UpdateDatabase(projectId, branchId, databaseName, newName string) (*DatabaseState, error) {
	ctx := c.GetContext()
	database, err := c.sdk.Database.Update(ctx, projectId, branchId, databaseName, neon.DatabaseUpdateRequest{
		Database: neon.DatabaseUpdateRequestDatabase{
			Name: newName,
		},
	})
	if err != nil {
		return nil, err
	}

	return &DatabaseState{
		DatabaseArgs: DatabaseArgs{
			ProjectId: projectId,
			BranchId:  branchId,
			Name:      database.Database.Name,
		},
		Id:        database.Database.ID,
		CreatedAt: database.Database.CreatedAt,
	}, nil
}

func (c *Client) DeleteDatabase(projectId, branchId, databaseName string) error {
	ctx := c.GetContext()
	_, err := c.sdk.Database.Delete(ctx, projectId, branchId, databaseName)
	return err
}