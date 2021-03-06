package main

import (
	log "github.com/hashicorp/go-hclog"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
	vsphere "github.com/martezr/vault-plugin-auth-vsphere"
)

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: vsphere.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		log.L().Error("plugin shutting down", "error", err)
		os.Exit(1)
	}
}
