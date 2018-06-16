package providers

import (
	"fmt"
	"net/url"
)

// Provider does something
type Provider interface {
	DownloadURL(url.URL) (string, error)
	GetName(url.URL) (string, error)
	GetVersion(url.URL) (string, error)
}

var providers map[string]Provider

func init() {
	providers = make(map[string]Provider)
}

// AddProvider lets you add a provider
func AddProvider(name string, provider Provider) {
	providers[name] = provider
}

// GetProvider gets you a provider
func GetProvider(provider string) (Provider, error) {
	if val, ok := providers[provider]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("%s is unsupported", provider)
}
