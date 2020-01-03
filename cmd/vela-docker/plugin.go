package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
)

// helper function to write a docker conf to kaniko user dir
func dockerConf(p plugin) error {

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", p.Username, p.Password)))

	err := ioutil.WriteFile(fmt.Sprint("/kaniko/.docker/config.json"), []byte(fmt.Sprintf(conf, p.Registry, auth)), 0644)
	if err != nil {
		return err
	}

	return nil
}

// helper function to convert the plugin configuration into kaniko CLI flags
func buildCommand(p plugin) []string {

	flags := []string{}

	// get the working directory for context
	dir, err := os.Getwd()
	if err != nil {
		return nil
	}

	// Add required and default fields to kaniko command
	logrus.Debug("...add required kaniko flags")
	for _, tag := range p.Tags {
		flags = append(flags, fmt.Sprintf("--destination=%s:%s", p.Repo, tag))
	}

	flags = append(flags, fmt.Sprintf("--verbosity=%s", p.LogLevel))
	if len(p.Context) != 0 {
		flags = append(flags, fmt.Sprintf("--context=%s/%s", dir, p.Context))
	} else if dir != "/" {
		flags = append(flags, fmt.Sprintf("--context=%s", dir))
	}

	// handle adding optional fields to kaniko commands
	logrus.Debug("...add optional kaniko flags")
	if len(p.Dockerfile) != 0 {
		flags = append(flags, fmt.Sprintf("--dockerfile=%s", p.Dockerfile))
	}
	if p.DryRun {
		flags = append(flags, fmt.Sprint("--no-push"))
	}
	if len(p.BuildArgs) != 0 {
		for _, arg := range p.BuildArgs {
			flags = append(flags, fmt.Sprintf("--build-arg=%s", arg))
		}
	}
	if p.Cache {
		flags = append(flags, fmt.Sprint("--cache"))
		if len(p.CacheRepo) != 0 {
			flags = append(flags, fmt.Sprintf("--cache-repo=%s", p.CacheRepo))
		} else {
			flags = append(flags, fmt.Sprintf("--cache-repo=%s", p.Repo))
		}
	}

	logrus.Debugf("...flags added %+v", flags)

	return flags
}

// helper function to run the kaniko binary against provided plugin configuration
func run(flags []string) error {

	cmd := exec.Command("/kaniko/executor", flags...)
	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}

	if errStdout != nil || errStderr != nil {
		return fmt.Errorf("Error: %s", err)
	}

	return nil
}
