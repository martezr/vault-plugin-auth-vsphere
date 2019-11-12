package vsphere

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathListVms(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "vms/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathVMList,
		},

		HelpSynopsis:    pathUserHelpSyn,
		HelpDescription: pathUserHelpDesc,
	}
}

func pathVms(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "vms/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the virtual machine.",
			},

			"policies": {
				Type:        framework.TypeCommaStringSlice,
				Description: "Comma-separated list of policies associated to the vm.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathVMDelete,
			logical.ReadOperation:   b.pathVMRead,
			logical.UpdateOperation: b.pathVMWrite,
			logical.CreateOperation: b.pathVMWrite,
		},

		HelpSynopsis:    pathUserHelpSyn,
		HelpDescription: pathUserHelpDesc,
	}
}

func (b *backend) pathVMDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "vm/"+d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) vm(ctx context.Context, s logical.Storage, vmname string) (*VMEntry, error) {
	if vmname == "" {
		return nil, fmt.Errorf("missing vmname")
	}

	entry, err := s.Get(ctx, "vm/"+strings.ToLower(vmname))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result VMEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathVMRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	vm, err := b.vm(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if vm == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"policies": vm.Policies,
		},
	}, nil
}

func (b *backend) pathVMWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := strings.ToLower(d.Get("name").(string))

	var policies = policyutil.ParsePolicies(d.Get("policies"))
	for _, policy := range policies {
		if policy == "root" {
			return logical.ErrorResponse("root policy cannot be granted by an auth method"), nil
		}
	}

	// Store it
	entry, err := logical.StorageEntryJSON("vm/"+d.Get("name").(string), &VMEntry{
		Name:     name,
		Policies: policies,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathVMList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	vms, err := req.Storage.List(ctx, "vm/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vms), nil
}

// VMEntry stores all the options that are set on a VM
type VMEntry struct {
	Name     string
	Policies []string
}

const pathUserHelpSyn = `
Manage users allowed to authenticate.
`

const pathUserHelpDesc = `
This endpoint allows you to create, read, update, and delete users
that are allowed to authenticate.
Deleting a user will not revoke auth for prior authenticated users
with that name. To do this, do a revoke on "login/<username>" for
the username you want revoked. If you don't need to revoke login immediately,
then the next renew will cause the lease to expire.
`
