//+build mage

package main

import (
	"runtime"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/build"
)

var cb build.BeforeBuildCallback = func(cfg build.Config) (build.Config, error) {
	cfg.EnableCGo = true
	return cfg, nil
}

func BuildPlugin() {
	if runtime.GOOS == "windows" {
		var err = build.SetBeforeBuildCallback(cb)

		if err != nil {
			log.DefaultLogger.Error(err.Error())
		}
	}

	build.BuildAll()
}

var Default = BuildPlugin
