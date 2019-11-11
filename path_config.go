package vsphere

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"vauth_server": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "vAuth Server",
			},
		},

		ExistenceCheck: b.pathConfigExistCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigCreateOrUpdate,
			logical.CreateOperation: b.pathConfigCreateOrUpdate,
			logical.ReadOperation:   b.pathConfigRead,
		},

		HelpSynopsis: pathConfigSyn,
    HelpDescription: pathConfigDesc,
	}
}

func (b *backend) pathConfigExistCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return false, err
	}

	if config == nil {
		return false, nil
	}

	return true, nil
}

func (b *backend) pathConfigCreateOrUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		cfg = &config{}
	}

	val, ok := data.GetOk("vauth_server")
	if ok {
		cfg.vAuthServer = val.(string)
	} else if req.Operation == logical.CreateOperation {
		cfg.vAuthServer = data.Get("vauth_server").(string)
	}
	if cfg.vAuthServer == "" {
		return logical.ErrorResponse("config parameter `vauth_server` cannot be empty"), nil
	}

	entry, err := logical.StorageEntryJSON("config", cfg)

	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"vauth_server": config.vAuthServer,
		},
	}
	return resp, nil
}

// Config returns the configuration for this backend.
func (b *backend) Config(ctx context.Context, s logical.Storage) (*config, error) {
	entry, err := s.Get(ctx, "config")

	if err != nil {
		return nil, err
	}

	var result config
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, fmt.Errorf("error reading configuration: %s", err)
		}
		return &result, nil
	}

	return nil, nil
}

type config struct {
	vAuthServer string `json:"vauth_server"`
}

const pathConfigSyn = `
VMware vSphere auth configuration
`

const pathConfigDesc = `
Use this endpoint to set vAuth endpoint settings.
`
