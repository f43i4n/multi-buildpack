package main

import (
	"fmt"
	"os"
	"path"

	"code.cloudfoundry.org/buildpackapplifecycle/buildpackrunner"
	"github.com/cloudfoundry/libbuildpack"
)

func writeProcfile(mainCommand string, additionalCommands []string) error {
	buildDir := os.Args[1]
	procfile := path.Join(buildDir, "Procfile")

	f, err := os.Create(procfile)
	if err != nil {
		return err
	}
	defer f.Close()

	for i, cmd := range additionalCommands {
		_, err = f.WriteString(fmt.Sprintf("proc_%v: %v\n", i+1, cmd))
		if err != nil {
			return err
		}
	}

	_, err = f.WriteString(fmt.Sprintf("main: %v\n", mainCommand))
	if err != nil {
		return err
	}

	return nil
}

func WriteStartCommand(stagingInfoFile string, outputFile string, additionalCommands []string) error {
	var stagingInfo buildpackrunner.DeaStagingInfo
	var webStartCommand map[string]string

	err := libbuildpack.NewYAML().Load(stagingInfoFile, &stagingInfo)
	if err != nil {
		return err
	}

	if len(additionalCommands) != 0 {
		if err := writeProcfile(stagingInfo.StartCommand, additionalCommands); err != nil {
			return err
		}

		webStartCommand = map[string]string{
			"web": "./bin/forego start",
		}
	} else {

		webStartCommand = map[string]string{
			"web": stagingInfo.StartCommand,
		}
	}

	release := buildpackrunner.Release{
		DefaultProcessTypes: webStartCommand,
	}

	return libbuildpack.NewYAML().Write(outputFile, &release)
}
