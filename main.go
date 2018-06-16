package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/mholt/archiver"
	"github.com/sbaildon/wow-addon-downloader/providers"
	_ "github.com/sbaildon/wow-addon-downloader/providers/curseforge"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
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

const (
	downloadError = 100
	mkdirError    = 101
	saveError     = 102
	unzipError    = 103
	missingError  = 104
)

func errorBar(bar *mpb.Bar, errorCode int64) {
	bar.SetTotal(errorCode, true)
	bar.IncrBy(int(bar.Total() - bar.Current()))
}

func download(provider providers.Provider, u url.URL, config config, bar *mpb.Bar, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := provider.GetName(u)
	if err != nil {
		errorBar(bar, missingError)
		return
	}

	_, err = provider.GetVersion(u)
	if err != nil {
		errorBar(bar, missingError)
		return
	}
	bar.Increment()

	downloadURL, err := provider.DownloadURL(u)
	if err != nil {
		errorBar(bar, downloadError)
		return
	}

	resp, err := http.Get(downloadURL)
	if err != nil {
		errorBar(bar, downloadError)
		return
	}
	defer resp.Body.Close()
	bar.Increment()

	if resp.StatusCode != http.StatusOK {
		errorBar(bar, downloadError)
		return
	}

	/* Create a temporary directory for saving files */
	tempDir, err := ioutil.TempDir("", "wow-addon-downloader")
	if err != nil {
		errorBar(bar, mkdirError)
		return
	}
	defer os.RemoveAll(tempDir)

	/* Save zip to temporary directory */
	out, err := os.Create(path.Join(tempDir, path.Base(resp.Request.URL.String())))
	if err != nil {
		errorBar(bar, saveError)
		return
	}
	defer out.Close()
	io.Copy(out, resp.Body)
	bar.Increment()

	/* Unzip files */
	err = archiver.Zip.Open(out.Name(), config.System.AddonDir)
	if err != nil {
		errorBar(bar, saveError)
		return
	}

	bar.Increment()
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
		log.Println("Problem reading config file")
	}

	var config config
	err = yaml.Unmarshal(configSource, &config)
	if err != nil {
		log.Println("Can't understand config file. Is it malformed?")

	}

	var wg sync.WaitGroup
	pool := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(10))

	var steps = []string{"searching", "downloading", "saving", "unzipping", "done"}

	for _, url := range config.AddOns {
		provider, err := providers.GetProvider(url.URL.Hostname())

		if err != nil {
			bar := pool.AddBar(1,
				mpb.BarClearOnComplete(),
				mpb.PrependDecorators(
					decor.StaticName(fmt.Sprintf("%s:", url.String()), 0, decor.DSyncSpace+decor.DidentRight),
					decor.StaticName(fmt.Sprint("unsupported"), 16, 1),
				),
			)
			bar.Increment()
			continue
		}

		bar := pool.AddBar(int64(len(steps)-1),
			mpb.BarClearOnComplete(),
			mpb.AppendDecorators(
				decor.OnComplete(decor.Percentage(5, 0), "", 0, decor.DwidthSync),
			),
			mpb.PrependDecorators(
				decor.StaticName(fmt.Sprintf("%s:", url.String()), 0, decor.DSyncSpace+decor.DidentRight),
				decor.DynamicName(func(s *decor.Statistics) string {
					switch s.Total {
					case downloadError:
						return "failed to download"
					case mkdirError:
						return "failed to create addon directory"
					case saveError:
						return "unable to save addon"
					case unzipError:
						return "unable to unzip addon"
					case missingError:
						return "unable to locate"
					default:
						return fmt.Sprintf("%s", steps[s.Current])
					}
				}, 16, 1),
				decor.OnComplete(
					decor.Elapsed(3, decor.DSyncSpace), "", 0, decor.DwidthSync,
				),
			),
		)

		wg.Add(1)
		go download(provider, *url.URL, config, bar, &wg)
	}

	pool.Wait()
	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
	}()
}
