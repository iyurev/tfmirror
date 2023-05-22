package localfs

import (
	"github.com/iyurev/tfmirror/pkg/config"
	"github.com/iyurev/tfmirror/pkg/mirror/localfs"
	"log"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	client, err := localfs.NewHttpClient(conf)
	if err != nil {
		log.Fatal(err)
	}
	err = client.DownloadMain(conf)
	if err != nil {
		log.Fatal(err)
	}

}
