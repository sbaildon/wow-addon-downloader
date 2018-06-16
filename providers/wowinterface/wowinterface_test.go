package wowinterface_test

import (
	"net/url"
	"testing"

	"github.com/sbaildon/wow-addon-downloader/providers"
	"github.com/stretchr/testify/assert"
)

const (
	testAddOn       = "http://www.wowinterface.com/downloads/info8736-BadBoySpamBlockerReporter.html"
	expectedVersion = "v7.3.126"
	expectedName    = "BadBoy: Spam Blocker & Reporter"
)

/* TestSum is slick */
func TestVersion(t *testing.T) {
	assert := assert.New(t)

	endpoint, err := url.Parse(testAddOn)
	assert.NoError(err)

	provider, err := providers.GetProvider(endpoint.Hostname())
	assert.NoError(err)

	version, err := provider.GetVersion(*endpoint)
	assert.NoError(err)
	assert.Equal(expectedVersion, version, "the correct version should be found")
}

/* TestGetName does stuff */
func TestGetName(t *testing.T) {
	assert := assert.New(t)

	endpoint, err := url.Parse(testAddOn)
	assert.NoError(err)

	provider, err := providers.GetProvider(endpoint.Hostname())
	assert.NoError(err)

	name, err := provider.GetName(*endpoint)
	assert.NoError(err)
	assert.Equal(expectedName, name, "the correct name should be found")
}
