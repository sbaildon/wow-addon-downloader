package curseforge

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/PuerkitoBio/goquery"
	"github.com/patrickmn/go-cache"
	"github.com/sbaildon/wow-addon-downloader/providers"
	"golang.org/x/net/html"
)

func contains(s []html.Attribute, key string, value string) bool {
	for _, attribute := range s {
		if attribute.Key == key {
			if attribute.Val == value {
				return true
			}
			return false
		}
	}
	return false
}

func cacheFetch(u url.URL) ([]byte, error) {
	body, cacheHit := pageCache.Get(u.String())
	if !cacheHit {
		resp, err := http.Get(u.String())
		if err != nil {
			return nil, fmt.Errorf("Can't connect to %s", u.String())
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, fmt.Errorf("Unable to read response body")
		}

		pageCache.Add(u.String(), body, cache.NoExpiration)
		return body, nil
	}

	return body.([]byte), nil
}

var pageCache *cache.Cache

func init() {
	providers.AddProvider("wow.curseforge.com", &CurseForge{})
	providers.AddProvider("www.wowace.com", &CurseForge{})
	pageCache = cache.New(cache.NoExpiration, cache.NoExpiration)
}

// CurseForge is a provider for curse
type CurseForge struct{}

// DownloadURL does something
func (p CurseForge) DownloadURL(u url.URL) (string, error) {
	u.Path = path.Join(u.Path, "/files/latest")
	return u.String(), nil
}

// GetName does something
func (p CurseForge) GetName(u url.URL) (string, error) {
	resp, err := cacheFetch(u)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return "", err
	}

	var name string
	doc.Find("meta").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if prop, found := s.Attr("property"); found {
			if prop == "og:title" {
				name, _ = s.Attr("content")
				return false
			}
		}
		return true
	})

	if len(name) > 0 {
		return name, nil
	}

	return "", fmt.Errorf("Unable to find name")
}

// GetVersion does something
func (p CurseForge) GetVersion(u url.URL) (string, error) {
	files, _ := url.Parse(u.String() + "/files")
	resp, err := cacheFetch(*files)

	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return "", err
	}

	var version string
	doc.Find(".listing-project-file").Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if attr, found := s.Attr("data-name"); found {
			version = attr
			return false
		}
		return true
	})

	if len(version) > 0 {
		return version, nil
	}

	return "", fmt.Errorf("Unable to find version")
}
