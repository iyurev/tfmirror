package types

import (
	"encoding/json"
	"fmt"
	"github.com/iyurev/tfmirror/pkg/tools"
	"os"
	"path"
)

const (
	OsLinux  = "linux"
	OsDarwin = "darwin"

	ArchAmd64 = "amd64"
	ArchArm64 = "arm64"
)

// GpgPublicKey GPG public key representation
type GpgPublicKey struct {
	KeyId          string `json:"key_id"`
	AsciiArmor     string `json:"ascii_armor"`
	TrustSignature string `json:"trust_signature"`
	Source         string `json:"source"`
	SourceUrl      string `json:"source_url"`
}

// SigningKeys  Represent list of GPG keys
type SigningKeys struct {
	GpgPublicKeys []GpgPublicKey `json:"gpg_public_keys"`
}

// Platform contains operating system and cpu architecture name.
type Platform struct {
	Os   string `json:"os"`
	Arch string `json:"arch"`
}

// NewPlatform return a new instance of Platform
func NewPlatform(os, arch string) *Platform {
	return &Platform{
		Os:   os,
		Arch: arch,
	}
}

// Name return string in format <OS_NAME>_<ARCH_NAME>, e.g. linux_amd64.
func (p *Platform) Name() string {
	return fmt.Sprintf("%s_%s", p.Os, p.Arch)
}

// PackageMetadata contains metadata about provider package.
type PackageMetadata struct {
	Protocols           []string    `json:"protocols"`
	Os                  string      `json:"os"`
	Arch                string      `json:"arch"`
	Filename            string      `json:"filename"`
	DownloadUrl         string      `json:"download_url"`
	ShasumsUrl          string      `json:"shasums_url"`
	ShasumsSignatureUrl string      `json:"shasums_signature_url"`
	Shasum              string      `json:"shasum"`
	SigningKeys         SigningKeys `json:"signing_keys"`
}

func (p *PackageMetadata) GetPlatform() *Platform {
	return &Platform{
		p.Os,
		p.Arch,
	}
}

// VersionObject is a full description of  provider version for set of available platforms.
type VersionObject struct {
	Version   string     `json:"version"`
	Protocols []string   `json:"protocols"`
	Platforms []Platform `json:"platforms"`
}

// HasPlatform
func (v *VersionObject) HasPlatform(platform Platform) bool {
	for _, p := range v.Platforms {
		if p.Os == platform.Os && p.Arch == platform.Arch {
			return true
		}
	}
	return false
}

type AvailableVersionsResponse struct {
	Id       string          `json:"id"`
	Versions []VersionObject `json:"versions"`
	Warnings string          `json:"warnings"`
}

// ProviderLocalVersionMetadata represents metadata for local provider packages.
type ProviderLocalVersionMetadata struct {
	Archives map[string]ProviderPlatformLocalMeta `json:"archives"`
}

// HasHash return true if ProviderLocalVersionMetadata already has provided H1 hash.
func (p *ProviderLocalVersionMetadata) HasHash(platformName, hash string) bool {
	if archive, ok := p.Archives[platformName]; ok {
		for _, h := range archive.Hashes {
			if h == hash {
				return true
			}
		}
	}
	return false
}

func NewProviderLocalIndex() *ProviderLocalVersionMetadata {
	return &ProviderLocalVersionMetadata{
		Archives: make(map[string]ProviderPlatformLocalMeta),
	}
}

func (p *ProviderLocalVersionMetadata) AddMeta(meta *ProviderPlatformLocalMeta, platform *Platform) error {
	name := platform.Name()
	p.Archives[name] = *meta
	return nil
}

func (p *ProviderLocalVersionMetadata) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *ProviderLocalVersionMetadata) Unmarshal(localFilePath string) error {
	f, err := os.ReadFile(localFilePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(f, p); err != nil {
		return err
	}
	return nil
}

func (p *ProviderLocalVersionMetadata) Save(localFilePath string) error {
	f, err := os.Create(localFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := p.Marshal()
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

// LocalProviderIndex represent index.json for local provider directory.
type LocalProviderIndex struct {
	Versions map[string]interface{} `json:"versions"`
}

// NewLocalIndex returns a new instance of LocalProviderIndex
func NewLocalIndex() *LocalProviderIndex {
	return &LocalProviderIndex{
		Versions: make(map[string]interface{}),
	}
}

// AddProviderIndex adds new provider index to index.json
func (l *LocalProviderIndex) AddProviderIndex(version string) {
	l.Versions[version] = make(map[string]string)
}

func (l *LocalProviderIndex) Marshal() ([]byte, error) {
	return json.Marshal(l)
}

func (l *LocalProviderIndex) Unmarshal(filePath string) error {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(f, l); err != nil {
		return err
	}
	return nil
}

func (l *LocalProviderIndex) Save(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := l.Marshal()
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

type ProviderPlatformLocalMeta struct {
	Hashes []string `json:"hashes"`
	Url    string   `json:"url"`
}

func NewArchiveMeta(localArchiveLoc string) (*ProviderPlatformLocalMeta, error) {
	archHash, err := tools.Hash1(localArchiveLoc)
	if err != nil {
		return nil, err
	}
	localUrl := path.Base(localArchiveLoc)
	return &ProviderPlatformLocalMeta{
		Url:    localUrl,
		Hashes: []string{archHash},
	}, nil
}
