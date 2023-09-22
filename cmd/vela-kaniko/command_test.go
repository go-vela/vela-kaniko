// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os/exec"
	"testing"
)

func TestDocker_execCmd(t *testing.T) {
	// setup types
	e := exec.Command("echo", "hello")

	err := execCmd(e)
	if err != nil {
		t.Errorf("execCmd returned err: %v", err)
	}
}
