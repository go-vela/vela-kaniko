// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

const kanikoBin = "/kaniko/executor"

// execCmd is a helper function to
// run the provided command.
func execCmd(e *exec.Cmd) error {
	logrus.Tracef("executing cmd %s", strings.Join(e.Args, " "))

	// set command stdout to OS stdout
	e.Stdout = os.Stdout
	// set command stderr to OS stderr
	e.Stderr = os.Stderr

	// output "trace" string for command
	fmt.Println("$", strings.Join(e.Args, " "))

	return e.Run()
}

// versionCmd is a helper function to output
// the client version information.
func versionCmd(ctx context.Context) *exec.Cmd {
	logrus.Trace("creating kaniko version command")

	// variable to store flags for command
	var flags []string

	// add flag for version kaniko command
	flags = append(flags, "version")

	return exec.CommandContext(ctx, kanikoBin, flags...)
}
