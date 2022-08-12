//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
	"github.com/samber/lo"
)

const (
	versionGolangcilint = "v1.48.0"
)

func Prepare() {
	lo.Must0(sh.RunV("lefthook", "install"))

	lo.Must0(sh.RunV("go", "install", fmt.Sprintf("github.com/golangci/golangci-lint/cmd/golangci-lint@%s", versionGolangcilint)))

	sh.RunV("asdf", "reshim", "golang") // nolint:errcheck
	sh.RunV("goenv", "rehash")          // nolint:errcheck
}
