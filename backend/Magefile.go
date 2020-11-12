//+build mage

package main

import (
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
)

// Callback for setting CGO_ENABLED for building on windows.
// See https://github.com/grafana/grafana-plugin-sdk-go/blob/002de0b2b925ac930952dc380682100dee27b439/build/common.go#L108.
var cb build.BeforeBuildCallback = func(cfg build.Config) (build.Config, error) {
	cfg.EnableCGo = true
	return cfg, nil
}

var err = build.SetBeforeBuildCallback(cb)

// Default configures the default target.
var Default = build.BuildAll
