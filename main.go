package main

import (
	"io/ioutil"
	"io"
	"log"
	"net/url"
	"net/http"
	"sync"
	"os"
	"path"
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/mholt/archiver"
	"github.com/gosuri/uiprogress"

	_ "github.com/shibukawa/configdir"
	"github.com/sbaildon/wow-addon-downloader/providers"
	_ "github.com/sbaildon/wow-addon-downloader/providers/curseforge"
)

type yamlurl struct {
	*url.URL
}

func (j *yamlurl) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	url, err := url.Parse(s)
	j.URL = url
	return err
}

func download(provider providers.Provider, u url.URL, config config, bar *uiprogress.Bar, wg *sync.WaitGroup) {
	defer wg.Done()

	bar.Incr()
	_, err := provider.GetName(u)
	if err != nil {
		log.Println(err)
		return
	}

	bar.Incr()
	_, err = provider.GetVersion(u)
	if err != nil {
		log.Println(err)
		return
	}

	var downloadURL = provider.DownloadURL(u)
	resp, err := http.Get(downloadURL)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Request failed, status not ok")
		return
	}

	/* Create a temporary directory for saving files */
	dir, err := ioutil.TempDir("", "wow-addon-downloader")
	if err != nil {
		log.Println(err)
		return
	}
	defer os.RemoveAll(dir)

	bar.Incr()
	/* Save zip to temporary directory */
	out, err := os.Create(path.Join(dir, path.Base(resp.Request.URL.String())))
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	bar.Incr()
	/* Unzip files */
	err = archiver.Zip.Open(out.Name(), config.System.AddonDir)
	if err != nil {
		log.Println(err)
		return
	}

	bar.Incr()

}

type config struct {
	System struct {
		AddonDir string `yaml:"addon_dir"`
	} `yaml:"system"`
	AddOns []yamlurl `yaml:"addons"`
}

func main() {
	configSource, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatal("Problem reading config file")
	}

	var config config
	err = yaml.Unmarshal(configSource, &config)
	if err != nil {
		log.Fatal("Can't understand config file. Is it malformed?")
	}

	var wg sync.WaitGroup

	uiprogress.Start()

	var steps = []string{"starting", "fetching name", "fetching version", "downloading", "writing to disk", "unzipping", "done"}

	for _, url := range config.AddOns {
		provider, err :=  providers.GetProvider(url.URL.Hostname())

		bar := uiprogress.AddBar(len(steps))
		bar.PrependElapsed()
		if err != nil {
			bar.PrependFunc(func(b *uiprogress.Bar) string{
				return fmt.Sprintf("%s: %-15s", url.String(), "unsupported")
			})
			bar.Set(0)
			continue
		}

		bar.PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("%s: %-15s", url.String(), steps[b.Current()-1])
		})

		wg.Add(1)
		go download(provider, *url.URL, config, bar, &wg)
	}

	wg.Wait()
	fmt.Println("done")
}
