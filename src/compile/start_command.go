package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"code.cloudfoundry.org/buildpackapplifecycle/buildpackrunner"
	"github.com/cloudfoundry/libbuildpack"
)

func writeProfileLine(file *os.File, name string, command string) error {
	// escape "
	command = strings.Replace(command, "\"", "\\\"", -1)
	_, err := file.WriteString(fmt.Sprintf("%v: bash -c \"%v\"\n", name, command))
	return err
}

func writeProcfile(mainCommand string, additionalCommands []string) error {
	buildDir := os.Args[1]
	procfile := path.Join(buildDir, "multi_procfile")

	f, err := os.Create(procfile)
	if err != nil {
		return err
	}
	defer f.Close()

	for i, cmd := range additionalCommands {
		err = writeProfileLine(f, fmt.Sprintf("proc_%v", i+1), cmd)
		if err != nil {
			return err
		}
	}

	err = writeProfileLine(f, "main", mainCommand)
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
			"web": "./bin/forego start -f multi_procfile",
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
