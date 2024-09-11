package provider

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string

const Name string = "neon"

func Provider() provider.Provider {
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
	ApiKey  string  `pulumi:"apiKey"`
	Version *string `pulumi:"version,optional"`
}

func (c *Config) Validate() error {
	if c.ApiKey == "" {
		return fmt.Errorf("apiKey is required")
	}
	return nil
}

// IsNotFoundError checks if the error is a "not found" error
func IsNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "404 Not Found")
}
