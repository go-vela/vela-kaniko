package main

import (
	"fmt"
	"regexp"
	"strings"
)

// helper function to validate fields being passed by user to plugin
func validate(e *env, p *plugin) error {

	// if the auto tag flag is set auto tag with commit or tag
	if p.AutoTag {
		switch e.BuildEvent {
		case "tag":
			p.Tags = append(p.Tags, e.BuildTag)
		default:
			p.Tags = append(p.Tags, e.BuildCommit)
		}
	}

	if len(p.Registry) == 0 {
		return fmt.Errorf("Plugin field registry is mandatory")
	}
	if len(p.Repo) == 0 {
		return fmt.Errorf("Plugin field repo is mandatory")
	}
	if len(p.Tags) == 0 {
		re := regexp.MustCompile(`^[A-Za-z0-9\-\.\_]*$`)

		// check each tag value for valid docker tag syntax
		// See docker docs for examples: https://docs.docker.com/engine/reference/commandline/tag/#Extended%20description
		for _, tag := range p.Tags {
			if !re.MatchString(tag) {
				return fmt.Errorf("Plugin tag %s has unaccepted rune. (Valid char set: abcdefghijklmnopqrstuvwxyz0123456789_-.ABCDEFGHIJKLMNOPQRSTUVWXY", tag)
			}
		}
		return fmt.Errorf("Plugin field tags or auto tag is mandatory")
	}
	if !strings.Contains("panic|fatal|error|warn|info|debug", p.LogLevel) {
		return fmt.Errorf("Plugin field debug accepted values are: panic|fatal|error|warn|info|debug")
	}

	if len(p.CacheRepo) != 0 {
		if !p.Cache {
			return fmt.Errorf("Plugin field cache is mandatory when using cache repo")
		}
	}

	if !p.DryRun {
		if len(p.Username) == 0 {
			return fmt.Errorf("Plugin field username or dry_run must be set")
		}

		if len(p.Password) == 0 {
			return fmt.Errorf("Plugin field password or dry_run must be set")
		}
	}

	return nil
}
