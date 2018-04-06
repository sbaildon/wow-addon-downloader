package curseforge

import (
	"path"
	"net/url"
	"net/http"
	"golang.org/x/net/html"
	"log"
	"github.com/sbaildon/wow-addon-downloader/providers"
)

func init() {
	providers.AddProvider("wow.curseforge.com", &CurseForge{})
}

type CurseForge struct {
}

func (p *CurseForge) DownloadURL(u url.URL) string {
	u.Path = path.Join(u.Path, "/files/latest")
	return u.String()
}

func (p *CurseForge) GetName() string {
	return "BadBoy: Spam Blocker & Reporter"
}

func (p *CurseForge) GetVersion(u url.URL) string {
	log.Print("Starting version fetch")
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}

	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		if tt != html.StartTagToken {
			continue
		}

		t := z.Token()
		if t.Data != "a" {
			continue
		}

		for _, a := range t.Attr {
			if a.Key == "data-name" {
				return a.Val
			}
		}
	}
}
