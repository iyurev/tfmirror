package localfs

import (
	"github.com/iyurev/tfmirror/pkg/config"
	"github.com/iyurev/tfmirror/pkg/types"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

var (
	testConf         *config.Conf
	testMirrorClient *Client
	testPkg          *types.PackageMetadata
)

func TestMain(t *testing.M) {
	var err error
	testConf, err = config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	testMirrorClient, err = NewHttpClient(testConf)
	if err != nil {
		log.Fatal(err)
	}
	testPkg, err = testMirrorClient.GetPackage("hashicorp/kubernetes", "2.20.0", *types.NewPlatform(types.OsLinux, types.ArchAmd64))
	if err != nil {
		log.Fatal(err)
	}
	t.Run()
}

func TestClient_DownloadProvider(t *testing.T) {
	err := testMirrorClient.DownloadProvider("hashicorp/kubernetes", "2.20.0", []*types.PackageMetadata{testPkg})
	_ = time.Now()
	require.NoError(t, err)

}

// func TestClient_DownloadMain(t *testing.T) {
// 	err := testMirrorClient.DownloadMain()
// 	require.NoError(t, err)
// }
