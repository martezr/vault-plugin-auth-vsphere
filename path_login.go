package vsphere

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const sourceHeader string = "vault-plugin-auth-vsphere"

type vm struct {
	Name       string `json:"name"`
	Datacenter string `json:"datacenter"`
	SecretKey  string `json:"secretkey"`
	Folder     string `json:"folder"`
}

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields: map[string]*framework.FieldSchema{
			"vmname": {
				Type:        framework.TypeString,
				Description: "The name of the computer account.",
			},
			"datacenter": {
				Type:        framework.TypeString,
				Description: "The name of the vSphere datacenter.",
			},
			"secretkey": {
				Type:        framework.TypeString,
				Description: "The secret key associated with the virtual machine.",
			},
			"folder": {
				Type:        framework.TypeString,
				Description: "The folder associated with the virtual machine.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLogin,
			logical.AliasLookaheadOperation: b.pathLoginAliasLookahead,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLoginAliasLookahead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	vmname := strings.ToLower(d.Get("vmname").(string))
	if vmname == "" {
		return nil, fmt.Errorf("missing vmname")
	}

	datacenter := strings.ToLower(d.Get("datacenter").(string))
	if datacenter == "" {
		return nil, fmt.Errorf("missing datacenter")
	}

	secretkey := strings.ToLower(d.Get("secretkey").(string))
	if secretkey == "" {
		return nil, fmt.Errorf("missing secretkey")
	}

	folder := strings.ToLower(d.Get("folder").(string))
	if folder == "" {
		return nil, fmt.Errorf("missing folder")
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: vmname,
			},
		},
	}, nil
}

func (b *backend) pathLogin(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Load and validate auth method configuration
	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse("could not load configuration"), nil
	}

	// Validate vmname argument
	vmname := strings.ToLower(d.Get("vmname").(string))

	if vmname == "" {
		return logical.ErrorResponse("missing vmname"), nil
	}

	folder := d.Get("folder").(string)

	if folder == "" {
		return logical.ErrorResponse("missing folder"), nil
	}

	datacenter := d.Get("datacenter").(string)
	if datacenter == "" {
		return nil, fmt.Errorf("missing datacenter")
	}

	secretkey := d.Get("secretkey").(string)
	if secretkey == "" {
		return nil, fmt.Errorf("missing secretkey")
	}

	vAuthServer := config.vAuthServer
	vAuthServer = "localhost"

	url := "http://" + vAuthServer + ":8090/vm/" + vmname
	b.Logger().Info(url)
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}

	var vmoutput vm

	json.NewDecoder(res.Body).Decode(&vmoutput)

	b.Logger().Info("Testing Stuff")

	if vmname != vmoutput.Name {
		return logical.ErrorResponse("Invalid VM name"), nil
	}
	b.Logger().Info(folder)
	b.Logger().Info(vmoutput.Folder)
	if folder != vmoutput.Folder {
		return logical.ErrorResponse("Invalid Folder"), nil
	}

	if datacenter != vmoutput.Datacenter {
		return logical.ErrorResponse("Invalid datacenter"), nil
	}

	if secretkey != vmoutput.SecretKey {
		return logical.ErrorResponse("Invalid secret key"), nil
	}

	//	vmdata, _ := b.vm(ctx, req.Storage, vmname)

	//	policies := vmdata.Policies

	resp := &logical.Response{
		Auth: &logical.Auth{
			Metadata: map[string]string{
				"vmname": vmname,
				//				"policies": strings.Join(policies, ","),
			},
			DisplayName: vmname,
			LeaseOptions: logical.LeaseOptions{
				TTL:       30,
				Renewable: false,
			},
			Alias: &logical.Alias{
				Name: vmname,
			},
		},
	}
	//	resp.Auth.Policies = append(resp.Auth.Policies, policies...)

	return resp, nil
}

const pathLoginSyn = `
Log in thVM name, vSphere datacenter, vSphere VM folder and a generated secret key.
`

const pathLoginDesc = `
This endpoint authenticates using the VM name, vSphere datacenter, vSphere VM folder and a generated secret key.
`
