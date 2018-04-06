package curseforge

import (
	"net/url"
	"github.com/sbaildon/wow-addon-downloader/providers"
)

func init() {
	providers.AddProvider("wow.curseforge.com", &CurseForge{})
}

type CurseForge struct {
}

func (p *CurseForge) DownloadURL(u url.URL) string {
	return "https://media.forgecdn.net/files/2531/838/BadBoy-v7.3.92.zip"
}

func (p *CurseForge) GetName() string {
	return "BadBoy: Spam Blocker & Reporter"
}

func (p *CurseForge) GetVersion() string {
	return "v7.3.92"
}
