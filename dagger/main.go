// Shuttle is a code-sharing and templating tool to easily inherit code actions across an org, whether they be in the same repository or spread over many
//
// Shuttle is built and operated by Lunar, and is a widely used tool to share out deployment plans and code actions across the org.

package main

import (
	"context"
	"dagger/shuttle-dagger/internal/dagger"
	"fmt"
	"log"
	"runtime"
)

type Shuttle struct{}

//https://github.com/lunarway/shuttle/releases/download/v0.23.0/shuttle-linux-amd64

const version = "v0.23.0"

// ShuttleBin Provides the raw binary for shuttle as a file
func (m *Shuttle) ShuttleBin() *dagger.File {
	os := runtime.GOOS
	arch := runtime.GOARCH

	log.Printf("os: %s, arch: %s", os, arch)

	return dag.HTTP(
		fmt.Sprintf(
			"https://github.com/lunarway/shuttle/releases/download/%s/shuttle-%s-%s",
			version, os, arch,
		),
	)
}

// Exec allows calling into a user generated query for the shuttle binary
func (m *Shuttle) Exec(
	directory *Directory,
	args ...string,
) *dagger.Container {
	return dag.Container().
		From("alpine:latest").
		WithFile("/usr/local/bin/shuttle", m.ShuttleBin(), dagger.ContainerWithFileOpts{
			Permissions: 755,
		}).
		WithDirectory("/mnt", directory).
		WithWorkdir("/mnt").
		WithExec(append([]string{"shuttle"}, args...))
}

// Version prints the version of the shuttle binary
func (m *Shuttle) Version(
	ctx context.Context,
	directory *Directory,
) (string, error) {
	return m.Exec(directory, "version").Stdout(ctx)
}

// Prepare pulls the latest plan from shuttle if found and exports it to the host directory at .shuttle
func (m *Shuttle) Prepare(
	ctx context.Context,
	directory *Directory,
	// +default=false
	skipPull bool,
) *dagger.Container {
	shuttle := m.Exec(directory, "prepare", "--skip-pull", fmt.Sprint(skipPull))

	return shuttle
}

// Run executes a user defined query on a prepared shuttle repository
func (m *Shuttle) Run(
	ctx context.Context,
	directory *Directory,
	args ...string,
) *dagger.Container {
	shuttle := m.Prepare(ctx, directory, true)

	return shuttle.WithExec(append([]string{"shuttle", "run"}, args...))
}
