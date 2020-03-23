package vsphere

import (
	"context"
	"encoding/json"

	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const sourceHeader string = "vault-plugin-auth-vsphere"

type vm struct {
	Name       string `json:"name"`
	Datacenter string `json:"datacenter"`
	SecretKey  string `json:"secretkey"`
	Role       string `json:"role"`
}

// pathLogin returns the path configurations for login endpoints
func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields: map[string]*framework.FieldSchema{
			"vmname": {
				Type:        framework.TypeString,
				Description: "The name of the virtual machine.",
			},
			"datacenter": {
				Type:        framework.TypeString,
				Description: "The name of the vSphere datacenter.",
			},
			"secretkey": {
				Type:        framework.TypeString,
				Description: "The secret key associated with the virtual machine.",
			},
			"role": {
				Type:        framework.TypeString,
				Description: "The role associated with the virtual machine.",
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

	role := strings.ToLower(d.Get("role").(string))
	if role == "" {
		return nil, fmt.Errorf("missing role")
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: vmname,
			},
		},
	}, nil
}

// pathLogin is used to authenticate to this backend
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

	role := d.Get("role").(string)

	if role == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	datacenter := d.Get("datacenter").(string)
	if datacenter == "" {
		return nil, fmt.Errorf("missing datacenter")
	}

	secretkey := d.Get("secretkey").(string)
	if secretkey == "" {
		return nil, fmt.Errorf("missing secretkey")
	}

	VAuthURL := config.VAuthURL
	url := VAuthURL + "/vm/" + vmname
	b.Logger().Info(url)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	//	var netClient = &http.Client{
	//		Timeout: time.Second * 10,
	//	}

	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}

	var vmoutput vm

	json.NewDecoder(res.Body).Decode(&vmoutput)

	if vmname != vmoutput.Name {
		return logical.ErrorResponse("Invalid VM name"), nil
	}

	if role != vmoutput.Role {
		return logical.ErrorResponse("Invalid role"), nil
	}

	if datacenter != vmoutput.Datacenter {
		return logical.ErrorResponse("Invalid datacenter"), nil
	}

	if secretkey != vmoutput.SecretKey {
		return logical.ErrorResponse("Invalid secret key"), nil
	}

	vmdata, _ := b.role(ctx, req.Storage, role)

	policies := vmdata.Policies

	resp := &logical.Response{
		Auth: &logical.Auth{
			Metadata: map[string]string{
				"vmname":   vmname,
				"policies": strings.Join(policies, ","),
			},
			DisplayName: vmname,
			LeaseOptions: logical.LeaseOptions{
				TTL:       30 * time.Minute,
				MaxTTL:    60 * time.Minute,
				Renewable: true,
			},
			Alias: &logical.Alias{
				Name: vmname,
			},
		},
	}
	resp.Auth.Policies = append(resp.Auth.Policies, policies...)

	return resp, nil
}

const pathLoginSyn = `
Log in with the VM name, vSphere datacenter, vSphere VM role and generated secret key.
`
const pathLoginDesc = `
This endpoint authenticates using the VM name, vSphere datacenter, vSphere VM role and a generated secret key.
`
