//+build mage

package main

import (
	"runtime"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/build"
)

// Callback for setting CGO_ENABLED for building on windows.
// See https://github.com/grafana/grafana-plugin-sdk-go/blob/002de0b2b925ac930952dc380682100dee27b439/build/common.go#L108.
var cb build.BeforeBuildCallback = func(cfg build.Config) (build.Config, error) {
	cfg.EnableCGo = true
	return cfg, nil
}

// Main entry-point for building plugin.
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
