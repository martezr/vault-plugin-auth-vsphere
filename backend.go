package vsphere

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// backend implements logical.Backend
type backend struct {
	*framework.Backend

	RolesMap *framework.PolicyMap
}

// Factory returns a new backend as logical.Backend.
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend

	b.RolesMap = &framework.PolicyMap{
		PathMap: framework.PathMap{
			Name: "roles",
		},
		DefaultKey: "default",
	}

	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
			SealWrapStorage: []string{
				"config",
			},
		},

		Paths: []*framework.Path{
			pathConfig(&b),
			pathLogin(&b),
			pathListRoles(&b),
			pathRoles(&b),
		},

		BackendType: logical.TypeCredential,
	}

	return &b
}

const backendHelp = `
The vSphere Auth Backend allows authentication for vSphere virtual machines.
`
