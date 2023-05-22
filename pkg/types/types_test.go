package types

import (
	"encoding/json"
	"log"
	"testing"
)

var (
	testProviderArch = "sample/mirror/registry.terraform.io/hashicorp/kubernetes/terraform-provider-kubernetes_2.16.0_darwin_amd64.zip"
)

func TestHash1(t *testing.T) {
	h, err := Hash1("sample/mirror/registry.terraform.io/hashicorp/kubernetes/terraform-provider-kubernetes_2.16.0_darwin_amd64.zip")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(h)
}

func TestNewArchiveMeta(t *testing.T) {
	// testPlatform := NewPlatform("linux", "amd64")
	archMeta, err := NewArchiveMeta(testProviderArch)
	if err != nil {
		t.Fatal(err)
	}
	archMetaJson, err := json.Marshal(archMeta)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Archive Metadata: %s\n", archMetaJson)
}

func TestNewProviderLocalIndex(t *testing.T) {
	localIndex := NewProviderLocalIndex()
	testPlatform := NewPlatform("linux", "amd64")
	archMeta, err := NewArchiveMeta(testProviderArch)
	if err != nil {
		t.Fatal(err)
	}
	if err := localIndex.AddMeta(archMeta, testPlatform); err != nil {
		t.Fatal(err)
	}
	j, err := localIndex.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Provider local index json: %s\n", j)
}

func TestNewLocalIndex(t *testing.T) {
	li := NewLocalIndex()
	li.AddProviderIndex("2.16.1")
	j, err := li.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", j)
}
