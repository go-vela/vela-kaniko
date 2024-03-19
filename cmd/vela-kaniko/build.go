// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// Build represents the plugin configuration for build information.
type Build struct {
	// event generated for build
	Event string
	// SHA-1 hash generated for commit
	Sha string
	// control how to snapshot the filesystem. - options (full|redo|time)
	SnapshotMode string
	// tag generated for build
	Tag string
	// https://github.com/GoogleContainerTools/kaniko#flag---use-new-run
	UseNewRun bool
	// https://github.com/GoogleContainerTools/kaniko#flag---tar-path
	TarPath string
	// https://github.com/GoogleContainerTools/kaniko#flag---single-snapshot
	SingleSnapshot bool
	// https://github.com/GoogleContainerTools/kaniko#flag---ignore-var-run
	IncludeVarRun bool
}

// SnapshotModeValues represents the available options for setting a snapshot mode.
//
// https://github.com/GoogleContainerTools/kaniko#--snapshotmode
var SnapshotModeValues = []string{"full", "redo", "time"}

// Validate verifies the Build is properly configured.
func (b *Build) Validate() error {
	logrus.Trace("validating build plugin configuration")

	// verify event is provided
	if len(b.Event) == 0 {
		return fmt.Errorf("no build event provided")
	}

	// verify sha is provided
	if len(b.Sha) == 0 {
		return fmt.Errorf("no build sha provided")
	}

	// verify the snapshot mode is a valid value
	if len(b.SnapshotMode) != 0 {
		// check if the value is a valid option
		if !isSnapshotModeValid(b.SnapshotMode) {
			return fmt.Errorf("snapshot mode was not a valid value - valid options (full|redo|time)")
		}
	}

	return nil
}

// isSnapshotModeValid checks if a value is within the list of accepted values.
func isSnapshotModeValid(value string) bool {
	// loop through snapshot values checking the value against the list
	for _, mode := range SnapshotModeValues {
		// when the value equals the mode return
		if strings.EqualFold(value, mode) {
			return true
		}
	}

	return false
}
