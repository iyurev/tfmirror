package localfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iyurev/tfmirror/pkg/config"
	tfmerr "github.com/iyurev/tfmirror/pkg/errors"
	"github.com/iyurev/tfmirror/pkg/tools"
	"github.com/iyurev/tfmirror/pkg/types"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	progbar "github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// v1/providers/hashicorp/random/versions
const (
	RegistryHost  = "registry.terraform.io"
	ProvidersPath = "v1/providers"
)

func init() {
	zlog.WithLevel(zerolog.InfoLevel)
}

type Client struct {
	httpClient   http.Client
	Host         string
	ProvidersUrl string
	WorkDir      string
	Conf         *config.Conf
}

func NewHttpClient(conf *config.Conf) (*Client, error) {
	headers := http.Header{}
	headers.Add("Accept", "application/json")
	t := &http.Transport{}
	c := http.Client{
		Timeout:   time.Second * time.Duration(conf.Client.TimeOut),
		Transport: t,
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var workDir = fmt.Sprintf("%s/workdir", cwd)
	if conf.Client.WorkDir != "" {
		workDir = conf.Client.WorkDir
	}
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, err
	}

	zlog.Info()
	client := &Client{
		httpClient:   c,
		Host:         RegistryHost,
		ProvidersUrl: fmt.Sprintf("https://%s/%s", RegistryHost, ProvidersPath),
		WorkDir:      workDir,
		Conf:         conf,
	}
	return client, nil
}

func (c *Client) ListVersions(providerSource string) (*types.AvailableVersionsResponse, error) {
	url := fmt.Sprintf("%s/%s/versions", c.ProvidersUrl, providerSource)
	body, err := c.DoRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	availableVersions := &types.AvailableVersionsResponse{}
	err = json.Unmarshal(body, availableVersions)
	if err != nil {
		return nil, err
	}
	return availableVersions, nil
}

func (c *Client) DoRequest(method, url string, requestBody io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, err
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if tfmerr.IsWrongStatusCode(response.StatusCode) {
		return nil, tfmerr.StatusCodeErr(response)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) GetPackage(providerSource, version string, platform types.Platform) (*types.PackageMetadata, error) {
	// https://registry.terraform.io/v1/providers/hashicorp/random/2.0.0/download/linux/amd64'
	url := fmt.Sprintf("%s/%s/%s/download/%s/%s", c.ProvidersUrl, providerSource, version, platform.Os, platform.Arch)
	respBody, err := c.DoRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	packageMeta := &types.PackageMetadata{}
	if err := json.Unmarshal(respBody, packageMeta); err != nil {
		return nil, err
	}
	return packageMeta, nil

}

func (c *Client) pkgLocalDirPath(providerSource string) string {
	return fmt.Sprintf("%s/%s/%s", c.WorkDir, c.Host, providerSource)
}

func (c *Client) MakePkgDir(providerSource string) error {
	if err := os.MkdirAll(c.pkgLocalDirPath(providerSource), 0750); err != nil {
		return err
	}
	return nil
}

func (c *Client) LocalArchivePath(provSource, archiveName string) string {
	return fmt.Sprintf("%s/%s/%s", c.WorkDir, provSource, archiveName)
}

func (c *Client) DownloadPackage(pkgMeta *types.PackageMetadata, destDir string) error {
	request, err := http.NewRequest(http.MethodGet, pkgMeta.DownloadUrl, nil)
	if err != nil {
		return err
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	pb := progbar.DefaultBytes(response.ContentLength, fmt.Sprintf("Downloading provider - %s", pkgMeta.Filename))
	destFilePath := fmt.Sprintf("%s/%s", destDir, pkgMeta.Filename)
	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()
	if _, err := io.Copy(io.MultiWriter(destFile, pb), response.Body); err != nil {
		return err
	}
	return nil
}

func (c *Client) DownloadProvider(provSource, version string, downloadList []*types.PackageMetadata) error {
	var localProviderDirPath = fmt.Sprintf("%s/%s", c.WorkDir, provSource)
	var localVerIndexPath = fmt.Sprintf("%s/%s.json", localProviderDirPath, version)
	var localIndexPath = fmt.Sprintf("%s/index.json", localProviderDirPath)
	var localIndex = types.NewLocalIndex()

	err := localIndex.Unmarshal(localIndexPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Local meta index file doesn't exists, we'll create it and go ahead.")
		} else {
			return err
		}
	}
	var localVerIndex = types.NewProviderLocalIndex()
	log.Printf("Add version %s to index.json", version)
	localIndex.AddProviderIndex(version)
	err = localVerIndex.Unmarshal(localVerIndexPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Local provider meta index file doesn't exists, we create it and go ahead.")
		} else {
			return err
		}
	}
	err = os.MkdirAll(localProviderDirPath, 0755)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup

	for _, pkgMeta := range downloadList {
		wg.Add(1)
		go func(pkgMeta *types.PackageMetadata) {
			archPath := c.LocalArchivePath(provSource, pkgMeta.Filename)
			needToDownload, err := NeedToDownload(archPath, pkgMeta, localVerIndex)
			if err != nil {
				log.Fatal(err)
			}
			if needToDownload {
				err = c.DownloadPackage(pkgMeta, localProviderDirPath)
				if err != nil {
					log.Fatal(err)
				}
				archMeta, err := types.NewArchiveMeta(archPath)
				if err != nil {
					log.Fatal(err)
				}
				err = localVerIndex.AddMeta(archMeta, pkgMeta.GetPlatform())
				if err != nil {
					log.Fatal(err)
				}
			}
			wg.Done()
		}(pkgMeta)
	}
	wg.Wait()
	if err := localVerIndex.Save(fmt.Sprintf("%s/%s.json", localProviderDirPath, version)); err != nil {
		return err
	}
	if err := localIndex.Save(localIndexPath); err != nil {
		return err
	}

	return nil
}

func (c *Client) DownloadMain() error {
	for _, provConf := range c.Conf.Providers {
		availableVersions, err := c.ListVersions(provConf.Source)
		if err != nil {
			return err
		}
		for _, version := range availableVersions.Versions {
			if provConf.HasVersion(version.Version) {
				var toDownloadList = make([]*types.PackageMetadata, 0)
				for _, platform := range provConf.Platforms {
					if version.HasPlatform(platform) {
						// Add to download list
						pkgMeta, err := c.GetPackage(provConf.Source, version.Version, platform)
						if err != nil {
							return err
						}
						toDownloadList = append(toDownloadList, pkgMeta)
					}
				}
				if err := c.DownloadProvider(provConf.Source, version.Version, toDownloadList); err != nil {
					return err
				}

			}
		}
	}

	return nil
}

func NeedToDownload(filePath string, pkgMeta *types.PackageMetadata, localProvIndex *types.ProviderLocalVersionMetadata) (bool, error) {
	fileExists, err := tools.IsExists(filePath)
	if err != nil {
		return true, err
	}
	if fileExists {
		fileHash, err := tools.Hash1(filePath)
		if err != nil {
			return true, err
		}
		if hashMatched := localProvIndex.HasHash(pkgMeta.GetPlatform().Name(), fileHash); hashMatched {
			return false, err
		}
		return true, nil

	}
	return true, nil
}
