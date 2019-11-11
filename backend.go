package vsphere

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend

	b.VmsMap = &framework.PolicyMap{
		PathMap: framework.PathMap{
			Name: "vms",
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
			pathListVms(&b),
			pathVms(&b),
		},

		BackendType: logical.TypeCredential,
	}

	return &b
}

type backend struct {
	*framework.Backend

	VmsMap *framework.PolicyMap
}

const backendHelp = `
The vSphere Auth Backend allows authentication for vSphere virtual machines.
`
