// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os/exec"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestDocker_Plugin_Exec_BadWrite(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "push",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    false,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:     true,
			CacheName: "index.docker.io/target/vela-kaniko",
			Name:      "index.docker.io/target/vela-kaniko",
			Tags:      []string{"latest"},
			AutoTag:   true,
		},
	}

	err := p.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestDocker_Plugin_Exec_BadExec(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "push",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    false,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:     true,
			CacheName: "index.docker.io/target/vela-kaniko",
			Name:      "index.docker.io/target/vela-kaniko",
			Tags:      []string{"latest"},
			AutoTag:   true,
		},
	}

	err := p.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestDocker_Plugin_Command(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
			IgnorePath:   []string{"/tmp", ".git"},
			LogTimestamp: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "foo",
		},
		Registry: &Registry{
			Name:               "index.docker.io",
			Username:           "octocat",
			Password:           "superSecretPassword",
			DryRun:             true,
			PushRetry:          1,
			InsecureRegistries: []string{"insecure.docker.local", "docker.local"},
			InsecurePull:       true,
			InsecurePush:       true,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--ignore-path=/tmp",
		"--ignore-path=.git",
		"--log-timestamp",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--insecure-registry=insecure.docker.local",
		"--insecure-registry=docker.local",
		"--insecure-pull",
		"--insecure",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_AutoTag_TagBuild(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "foo",
		},
		Registry: &Registry{
			Name:               "index.docker.io",
			Username:           "octocat",
			Password:           "superSecretPassword",
			DryRun:             true,
			PushRetry:          1,
			InsecureRegistries: []string{"insecure.docker.local", "docker.local"},
			InsecurePull:       true,
			InsecurePush:       true,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	// configure repo tags using auto_tag and build info
	p.Repo.ConfigureAutoTagBuildTags(p.Build)

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--destination=index.docker.io/target/vela-kaniko:v0.0.0",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--insecure-registry=insecure.docker.local",
		"--insecure-registry=docker.local",
		"--insecure-pull",
		"--insecure",
		"--verbosity=info",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_AutoTag_PushBuild(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "push",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "foo",
		},
		Registry: &Registry{
			Name:               "index.docker.io",
			Username:           "octocat",
			Password:           "superSecretPassword",
			DryRun:             true,
			PushRetry:          1,
			InsecureRegistries: []string{"insecure.docker.local", "docker.local"},
			InsecurePull:       true,
			InsecurePush:       true,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	// configure repo tags using auto_tag and build info
	p.Repo.ConfigureAutoTagBuildTags(p.Build)

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--destination=index.docker.io/target/vela-kaniko:7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--insecure-registry=insecure.docker.local",
		"--insecure-registry=docker.local",
		"--insecure-pull",
		"--insecure",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_Labels(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			Labels:            []string{"key1=tag1"},
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=key1=tag1",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_MultipleTopics(t *testing.T) {
	// setup types
	label := testLabel()
	label.Topics = []string{"foo", "bar"}

	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             label,
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=foo,bar",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_MultipleTopicsWithFilter(t *testing.T) {
	// setup types
	label := testLabel()
	label.Topics = []string{"id123", "bar"}

	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			TopicsFilter:      "^id",
			Label:             label,
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_MultipleTopicsNoTopics(t *testing.T) {
	// setup types
	label := testLabel()
	label.Topics = []string{}

	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             label,
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_CustomLabels(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
			IgnorePath:   []string{"/tmp", ".git"},
			LogTimestamp: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "foo",
		},
		Registry: &Registry{
			Name:               "index.docker.io",
			Username:           "octocat",
			Password:           "superSecretPassword",
			DryRun:             true,
			PushRetry:          1,
			InsecureRegistries: []string{"insecure.docker.local", "docker.local"},
			InsecurePull:       true,
			InsecurePush:       true,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	p.Repo.Label.CustomSet = []string{"label1=foo", "label2=bar"}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--ignore-path=/tmp",
		"--ignore-path=.git",
		"--log-timestamp",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--label=label1=foo",
		"--label=label2=bar",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--insecure-registry=insecure.docker.local",
		"--insecure-registry=docker.local",
		"--insecure-pull",
		"--insecure",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_SnapshotMode(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			SnapshotMode: "redo",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "foo",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--snapshot-mode=redo",
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_UseNewRun(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			UseNewRun:    true,
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "foo",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--use-new-run",
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_TarPath(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			TarPath:      "build",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "foo",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--tar-path=build",
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_UseSingleSnapshot(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:          "tag",
			Sha:            "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:            "v0.0.0",
			SingleSnapshot: true,
			IgnoreVarRun:   true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "foo",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--single-snapshot",
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_UseIgnoreVarRunFalse(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: false,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "foo",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=false",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_UseIgnoreVarRunTrue(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "foo",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_ForceBuildMetaData(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:               []string{"foo=bar"},
			Context:            ".",
			Dockerfile:         "Dockerfile",
			Target:             "foo",
			ForceBuildMetadata: true,
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--force-build-metadata",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_Mirror(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "foo",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Mirror:    "company.mirror.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--registry-mirror=company.mirror.io",
		"--push-retry=1",
		"--target=foo",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_With_Compression(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "foo",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Compression:       "zstd",
			CompressionLevel:  3,
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--compression=zstd",
		"--compression-level=3",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_NoCacheRepo(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "push",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    true,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_NoDryRun(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "push",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "",
		},
		Registry: &Registry{
			Name:      "index.docker.io",
			Username:  "octocat",
			Password:  "superSecretPassword",
			DryRun:    false,
			PushRetry: 1,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--push-retry=1",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_CustomPlatform(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:           []string{"foo=bar"},
			Context:        ".",
			Dockerfile:     "Dockerfile",
			Target:         "foo",
			CustomPlatform: "linux/arm64/v8",
		},
		Registry: &Registry{
			Name:               "index.docker.io",
			Username:           "octocat",
			Password:           "superSecretPassword",
			DryRun:             true,
			PushRetry:          1,
			InsecureRegistries: []string{"insecure.docker.local", "docker.local"},
			InsecurePull:       true,
			InsecurePush:       true,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             testLabel(),
			CompressedCaching: true,
		},
	}

	want := exec.Command(
		kanikoBin,
		"--ignore-var-run=true",
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-kaniko",
		"--context=.",
		"--destination=index.docker.io/target/vela-kaniko:latest",
		"--label=org.opencontainers.image.created=now",
		"--label=org.opencontainers.image.url=git.example.com",
		"--label=org.opencontainers.image.revision=deadbeef",
		"--label=io.vela.build.author=octocat@example.com",
		"--label=io.vela.build.number=1",
		"--label=io.vela.build.repo=octocat/scripts",
		"--label=io.vela.build.commit=deadbeef",
		"--label=io.vela.build.url=git.example.com",
		"--label=io.vela.build.link=https://vela.example.com/velaOrg/velaRepo/1",
		"--label=io.vela.build.host=vela-worker",
		"--label=io.vela.build.topics=id123",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--push-retry=1",
		"--target=foo",
		"--custom-platform=linux/arm64/v8",
		"--insecure-registry=insecure.docker.local",
		"--insecure-registry=docker.local",
		"--insecure-pull",
		"--insecure",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !strings.EqualFold(sortCmdArgs(got).String(), sortCmdArgs(want).String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Validate(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "push",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "",
		},
		Registry: &Registry{
			Name:     "index.docker.io",
			Username: "octocat",
			Password: "superSecretPassword",
			DryRun:   false,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             &Label{},
			CompressedCaching: true,
		},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDocker_Plugin_Validate_NoBuild(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "",
		},
		Registry: &Registry{
			Name:     "index.docker.io",
			Username: "octocat",
			Password: "superSecretPassword",
			DryRun:   false,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			CompressedCaching: true,
		},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Plugin_Validate_AutoTag_InvalidBuildTag(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "tag",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "-v0.0.0",
			SnapshotMode: "full",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "",
		},
		Registry: &Registry{
			Name:     "index.docker.io",
			Username: "octocat",
			Password: "superSecretPassword",
			DryRun:   false,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			Label:             &Label{},
			CompressedCaching: true,
		},
	}

	// configure auto_tag using the invalid build tag
	p.Repo.ConfigureAutoTagBuildTags(p.Build)

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Plugin_Validate_NoImage(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "push",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{},
		Registry: &Registry{
			Name:     "index.docker.io",
			Username: "octocat",
			Password: "superSecretPassword",
			DryRun:   false,
		},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			CompressedCaching: true,
		},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Plugin_Validate_NoRegistry(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "push",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "",
		},
		Registry: &Registry{},
		Repo: &Repo{
			Cache:             true,
			CacheName:         "index.docker.io/target/vela-kaniko",
			Name:              "index.docker.io/target/vela-kaniko",
			Tags:              []string{"latest"},
			AutoTag:           true,
			CompressedCaching: true,
		},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Plugin_Validate_NoRepo(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event:        "push",
			Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:          "v0.0.0",
			IgnoreVarRun: true,
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
			Target:     "",
		},
		Registry: &Registry{
			Name:     "index.docker.io",
			Username: "octocat",
			Password: "superSecretPassword",
			DryRun:   false,
		},
		Repo: &Repo{},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func testLabel() *Label {
	return &Label{
		AuthorEmail: "octocat@example.com",
		Commit:      "deadbeef",
		Created:     "now",
		FullName:    "octocat/scripts",
		Number:      1,
		Topics:      []string{"id123"},
		URL:         "git.example.com",
		BuildURL:    "https://vela.example.com/velaOrg/velaRepo/1",
		Host:        "vela-worker",
	}
}

func sortCmdArgs(cmd *exec.Cmd) *exec.Cmd {
	labels := []string{}
	otherArgs := []string{}

	for _, arg := range cmd.Args {
		if strings.HasPrefix(arg, "--label") {
			labels = append(labels, arg)
		} else {
			otherArgs = append(otherArgs, arg)
		}
	}

	sort.Strings(labels)

	cmd.Args = append(otherArgs, labels...)

	return cmd
}
