package main

import (
	"os"

	"github.com/spf13/afero"

	config "github.com/tsingson/discovery/conf"
)

var ( // global variable
	// cacheSize                int
	// cacheTimeOut             int64
	path, logPath string
	cfg           *config.Config
)

func init() {
	var err error
	afs := afero.NewOsFs()
	{ // setup path for storage of log / configuration / cache
		// path = "/Users/qinshen/git/linksmart/bin"  // for test
		path, err = GetCurrentExecDir()
		if err != nil {
			// fmt.Println("无法读取可执行程序的存储路径")
			panic("无法读取可执行程序的存储路径")
			os.Exit(-1)
		}

	}
	{ // load config for discovery daemon
		configToml := path + "/discoveryd-config.toml"

		cfg, err = config.LoadConfig(configToml)
		if err != nil {
			// fmt.Println("无法读取可执行程序的存储路径")
			panic("无法读取可执行程序的存储路径")
			os.Exit(-1)
		}

	}
	{
		logPath = path + "/log"
		check, _ := afero.DirExists(afs, logPath)
		if !check {
			err = afs.MkdirAll(logPath, 0755)
			if err != nil {
				panic("mkdir log path fail")
				os.Exit(-1)
			}
		}
	}
}
