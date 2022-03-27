package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// This package is used to generate the JSON used to run the crosscompile job

type matrix struct {
	Goos      string `json:"goos"`
	Goarch    string `json:"goarch"`
	GoVersion string `json:"go_version"`
}

func main() {
	matrixItems := make([]matrix, 0)

	goversions := []string{"1.17.x", "1.18.x"}

	add := func(version, os, arch string) {
		matrixItems = append(matrixItems, matrix{
			Goos:      os,
			Goarch:    arch,
			GoVersion: version,
		})
	}

	for _, version := range goversions {

		// Add all the archs for linux
		add(version, "linux", "amd64")
		add(version, "linux", "386")
		add(version, "linux", "arm")
		add(version, "linux", "arm64")
		add(version, "linux", "ppc64")
		add(version, "linux", "ppc64le")
		add(version, "linux", "s390x")
		add(version, "linux", "mips")
		add(version, "linux", "mipsle")
		add(version, "linux", "mips64")
		add(version, "linux", "mips64le")

		// Add basic x86 and arm architectures for everything else
		oss := []string{"darwin", "freebsd", "netbsd", "openbsd", "dragonfly", "solaris", "windows"}
		archs := []string{"amd64", "386", "arm", "arm64"}

		for _, arch := range archs {
			for _, os := range oss {
				if (arch == "386" || arch == "arm") && os == "darwin" {
					// Golang dropped support for darwin 32bits since go1.15

				} else if (os == "dragonfly" || os == "solaris") && arch != "amd64" {
					// dragonfly and solaris only build for amd64
					continue
				} else {
					add(version, os, arch)
				}
			}
		}
	}

	// Make sure darwin 32bit still builds for 1.14
	add("1.14.x", "darwin", "386")
	add("1.14.x", "darwin", "arm")

	// Run a single test for an old go version.
	add("1.6.x", "linux", "amd64")

	json, err := json.Marshal(matrixItems)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	n, err := os.Stdout.Write(json)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if n != len(json) {
		fmt.Println("Didn't fully write JSON data")
		os.Exit(1)
	}
}
