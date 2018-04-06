package providers

import (
	"net/url"
	"fmt"
)

type Provider interface {
	DownloadURL(url.URL) string
	GetName() string
	GetVersion() string
}

var Providers map[string]Provider

func init() {
	Providers = make(map[string]Provider)
}

func AddProvider(name string, provider Provider) {
	Providers[name] = provider
}

func GetProvider(provider string) (Provider, error) {
	if val, ok := Providers[provider]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("Provider not found")
}
