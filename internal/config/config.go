package config

import (
	"flag"
	"log"
	"strings"
	"time"

	env "github.com/caarlos0/env/v6"
	// "go.starlark.net/lib/time"
)

type Config struct {
	RequestsAddresses []string
	LogFileName       string
	ReceiveTimeout    time.Duration
}

type configEnv struct {
	RequestsAddresses string        `env:"REQUESTS_ADDRESSES"`
	LogFileName       string        `env:"LOG_FILENAME"`
	ReceiveTimeout    time.Duration `env:"RECEIVE_TIMEOUT"`
}

func Init() (Config, error) {

	var ServerConfig Config
	// tmpFile := filepath.Join(os.TempDir(), "short-url-db.json")
	var RequestsAddresses string
	flag.StringVar(&RequestsAddresses, "ra", "google.com, mail.ru, xcap-portal.slb.ru", "request addresses")
	flag.StringVar(&ServerConfig.LogFileName, "l", "log.json", "log file name")
	// flag.DurationVar(&ServerConfig.ReceiveTimeout, "rt", time.Second*2, "receive timeout")
	flag.DurationVar(&ServerConfig.ReceiveTimeout, "rt", time.Millisecond*400, "receive timeout")
	flag.Parse()

	var cfg configEnv

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.RequestsAddresses != "" {
		arr := strings.Split(cfg.RequestsAddresses, ",")
		ServerConfig.RequestsAddresses = arr
	} else {
		arr := strings.Split(RequestsAddresses, ",")
		ServerConfig.RequestsAddresses = arr
	}

	if cfg.LogFileName != "" {
		ServerConfig.LogFileName = cfg.LogFileName
	}

	if cfg.ReceiveTimeout != 0 {
		ServerConfig.ReceiveTimeout = cfg.ReceiveTimeout
	}

	return ServerConfig, nil
}
