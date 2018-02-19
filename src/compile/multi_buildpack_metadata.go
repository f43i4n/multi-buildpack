package main

import (
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
)

// MultiBuildpackMetadata is a struct to parse multi-buildpack.yml
type MultiBuildpackMetadata struct {
	Buildpacks []string `yaml:"buildpacks"`

	// additional Commands for start command
	AdditionalCommands []string `yaml:"additionalCommands"`
}

// returns list of buildpacks and additional commands
func GetBuildpacks(dir string, logger *libbuildpack.Logger) ([]string, []string, error) {
	metadata := &MultiBuildpackMetadata{}

	err := libbuildpack.NewYAML().Load(filepath.Join(dir, "multi-buildpack.yml"), metadata)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Error("A multi-buildpack.yml file must be provided at your app root to use this buildpack.")
		} else {
			logger.Error("The multi-buildpack.yml file is malformed.")
		}
		return nil, nil, err
	}

	return metadata.Buildpacks, metadata.AdditionalCommands, nil
}
