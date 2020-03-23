package vsphere

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathConfig returns the path configuration for CRUD operations on the backend
// configuration.
func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"vauth_url": {
				Type:        framework.TypeString,
				Description: "vAuth URL address (https://vauth.grt.local)",
			},
		},
		ExistenceCheck: b.pathConfigExistCheck,
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigCreateOrUpdate,
			logical.CreateOperation: b.pathConfigCreateOrUpdate,
			logical.ReadOperation:   b.pathConfigRead,
		},

		HelpSynopsis: pathConfigSyn,
	}
}

// config contains the URL of the vAuth server used for VM validation
type config struct {
	VAuthURL string `json:"vauth_url"`
}

// pathConfigExistCheck checks for the existence of a configuration
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

// pathConfigCreateOrUpdate handles create and update commands to the config
func (b *backend) pathConfigCreateOrUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		cfg = &config{}
	}

	val, ok := data.GetOk("vauth_url")
	if ok {
		cfg.VAuthURL = val.(string)
	} else if req.Operation == logical.CreateOperation {
		cfg.VAuthURL = data.Get("vauth_url").(string)
	}
	if cfg.VAuthURL == "" {
		return logical.ErrorResponse("config parameter `vauth_url` cannot be empty"), nil
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

// pathConfigWrite handles create and update commands to the config
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
			"vauth_url": config.VAuthURL,
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

const pathConfigSyn = `
This path allows you to configure the VMware vSphere auth provider to interact with the vAuth Identity Platform
for authenticating virtual machines.`
