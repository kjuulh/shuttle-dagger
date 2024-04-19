// A generated module for ShuttleDagger functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/shuttle-dagger/internal/dagger"
	"fmt"
	"runtime"
)

type Shuttle struct{}

//https://github.com/lunarway/shuttle/releases/download/v0.23.0/shuttle-linux-amd64

const version = "v0.23.0"

func (m *Shuttle) ShuttleContainer() *dagger.File {
	os := runtime.GOOS
	arch := runtime.GOARCH

	return dag.HTTP(
		fmt.Sprintf(
			"https://github.com/lunarway/shuttle/releases/download/%s/shuttle-%s-%s",
			version, os, arch,
		),
	)
}

// Returns lines that match a pattern in the files of the provided Directory
func (m *Shuttle) GrepDir(ctx context.Context, directoryArg *Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}
