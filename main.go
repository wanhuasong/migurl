package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/juju/errors"
	"github.com/wanhuasong/migurl/migration"
)

var (
	apiBaseURL            string
	qiniuPublicDomain     string
	genericStorageBaseURL string
	cfgFile               string
)

func main() {
	initFlags()

	m, err := migration.NewMigration(cfgFile, apiBaseURL, qiniuPublicDomain, genericStorageBaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", errors.Trace(err))
		os.Exit(1)
	}
	if err = m.Do(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", errors.Trace(err))
		os.Exit(1)
	}
}

func initFlags() {
	flag.StringVar(&apiBaseURL, "api-base-url", "", "specify api_base_url while using local storage")
	flag.StringVar(&qiniuPublicDomain, "qiniu-public-domain", "", "specify qiniu_public_domain while using qiniu storage")
	flag.StringVar(&genericStorageBaseURL, "generic-storage-base-url", "", "specify generic_storage_base_url while using generic storage")
	flag.StringVar(&cfgFile, "config", "./config.json", "config file")
	flag.Parse()
}
