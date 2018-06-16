package wowinterface

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

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
	providers.AddProvider("www.wowinterface.com", &WoWInterface{})
	pageCache = cache.New(cache.NoExpiration, cache.NoExpiration)
}

// WoWInterface is a provider for wowi
type WoWInterface struct{}

// DownloadURL does something
func (p WoWInterface) DownloadURL(u url.URL) (string, error) {
	reg, _ := regexp.Compile("info")
	stringer := reg.ReplaceAllString(u.String(), "download")
	new, _ := url.Parse(stringer)

	resp, err := cacheFetch(*new)

	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return "", err
	}

	val, ok := doc.Find(".manuallink").Find("a").Attr("href")
	if !ok {
		return "", errors.New("Expected download href")
	}

	return val, nil
}

// GetName does something
func (p WoWInterface) GetName(u url.URL) (string, error) {
	resp, err := cacheFetch(u)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return "", err
	}

	var title string
	doc.Find("meta").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if prop, found := s.Attr("property"); found {
			if prop == "og:title" {
				title, _ = s.Attr("content")
				return false
			}
		}
		return true
	})

	if len(title) > 0 {
		return title, nil
	}

	return "", fmt.Errorf("Unable to find title")
}

// GetVersion does something
func (p WoWInterface) GetVersion(u url.URL) (string, error) {
	resp, err := cacheFetch(u)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return "", err
	}

	element := doc.Find("#version")
	ver := element.Text()

	reg, _ := regexp.Compile("(Version: )")
	return reg.ReplaceAllString(ver, ""), nil
}
