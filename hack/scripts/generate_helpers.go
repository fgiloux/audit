// Copyright 2021 The Audit Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Deprecated
// This script is only a helper
// todo: remove after 4.9-GA
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/operator-framework/audit/hack"
	"github.com/operator-framework/audit/pkg"
)

// nolint:gocyclo,dupl,lll,staticcheck
func main() {

	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fullReportsPath := filepath.Join(currentPath, hack.ReportsPath)

	dirs := map[string]string{
		"redhat_certified_operator_index": "registry.redhat.io/redhat/certified-operator-index",
		"redhat_redhat_marketplace_index": "registry.redhat.io/redhat/redhat-marketplace-index",
		"redhat_redhat_operator_index":    "registry.redhat.io/redhat/redhat-operator-index",
	}

	dirGreen := filepath.Join(fullReportsPath, "deprecate-green")
	command := exec.Command("rm", "-rf", dirGreen)
	_, err = pkg.RunCommand(command)
	if err != nil {
		log.Errorf("running command :%s", err)
	}

	command = exec.Command("mkdir", dirGreen)
	_, err = pkg.RunCommand(command)
	if err != nil {
		log.Errorf("running command :%s", err)
	}

	// nolint:scopelint
	for dir := range dirs {
		pathToWalk := filepath.Join(fullReportsPath, dir)
		if _, err := os.Stat(pathToWalk); err != nil && os.IsNotExist(err) {
			continue
		}

		// Walk in the testdata dir and generates the deprecated-api custom dashboard for
		// all bundles JSON reports available there
		err = filepath.Walk(pathToWalk, func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() && strings.HasSuffix(info.Name(), "json") &&
				strings.HasPrefix(info.Name(), "bundles") {

				// Ignore the tag images 4.6 and 4.7
				if strings.Contains(info.Name(), "v4.7") ||
					strings.Contains(info.Name(), "v4.6") ||
					strings.Contains(info.Name(), "v4.8") {
					return nil
				}

				command := exec.Command("go",
					"run",
					"./hack/scripts/deprecated-bundles-repo/deprecate-green/deprecated_index.go",
					fmt.Sprintf("--file=%s", filepath.Join(pathToWalk, info.Name())),
					fmt.Sprintf("--output=%s", dirGreen))
				_, err = pkg.RunCommand(command)
				if err != nil {
					log.Warnf("running command :%s", err)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	dirAll := filepath.Join(fullReportsPath, "deprecate-all")
	command = exec.Command("rm", "-rf", dirAll)
	_, err = pkg.RunCommand(command)
	if err != nil {
		log.Errorf("running command :%s", err)
	}

	command = exec.Command("mkdir", dirAll)
	_, err = pkg.RunCommand(command)
	if err != nil {
		log.Errorf("running command :%s", err)
	}

	// nolint:scopelint
	for dir := range dirs {
		pathToWalk := filepath.Join(fullReportsPath, dir)
		if _, err := os.Stat(pathToWalk); err != nil && os.IsNotExist(err) {
			continue
		}

		// Walk in the testdata dir and generates the deprecated-api custom dashboard for
		// all bundles JSON reports available there
		err = filepath.Walk(pathToWalk, func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() && strings.HasSuffix(info.Name(), "json") &&
				strings.HasPrefix(info.Name(), "bundles") {

				if strings.Contains(info.Name(), "v4.7") ||
					strings.Contains(info.Name(), "v4.6") ||
					strings.Contains(info.Name(), "v4.8") {
					return nil
				}

				command := exec.Command("go",
					"run",
					"./hack/scripts/deprecated-bundles-repo/deprecate-all/deprecated_index.go",
					fmt.Sprintf("--file=%s", filepath.Join(pathToWalk, info.Name())),
					fmt.Sprintf("--output=%s", dirAll))
				_, err = pkg.RunCommand(command)
				if err != nil {
					log.Warnf("running command :%s", err)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	command = exec.Command("rm", "-rf", filepath.Join(fullReportsPath, "contacts"))
	_, err = pkg.RunCommand(command)
	if err != nil {
		log.Errorf("running command :%s", err)
	}

	command = exec.Command("rm", "-rf", filepath.Join(fullReportsPath, "packages"))
	_, err = pkg.RunCommand(command)
	if err != nil {
		log.Errorf("running command :%s", err)
	}

	// nolint:scopelint
	for dir := range dirs {
		pathToWalk := filepath.Join(fullReportsPath, dir)
		if _, err := os.Stat(pathToWalk); err != nil && os.IsNotExist(err) {
			continue
		}

		// Walk in the testdata dir and generates the deprecated-api custom dashboard for
		// all bundles JSON reports available there
		err = filepath.Walk(pathToWalk, func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() && strings.HasSuffix(info.Name(), "json") &&
				strings.HasPrefix(info.Name(), "bundles") {

				// Ignore the tag images 4.6 and 4.7
				if strings.Contains(info.Name(), "v4.7") ||
					strings.Contains(info.Name(), "v4.6") {
					return nil
				}

				command := exec.Command("go",
					"run",
					"./hack/scripts/packages/generate.go",
					fmt.Sprintf("--file=%s", filepath.Join(pathToWalk, info.Name())))
				_, err = pkg.RunCommand(command)
				if err != nil {
					log.Warnf("running command :%s", err)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	ivsDir := filepath.Join(fullReportsPath, "ivs")
	command = exec.Command("rm", "-rf", ivsDir)
	_, err = pkg.RunCommand(command)
	if err != nil {
		log.Errorf("running command :%s", err)
	}

	command = exec.Command("mkdir", ivsDir)
	_, err = pkg.RunCommand(command)
	if err != nil {
		log.Errorf("running command :%s", err)
	}

	// nolint:scopelint
	for dir := range dirs {
		pathToWalk := filepath.Join(fullReportsPath, dir)
		if _, err := os.Stat(pathToWalk); err != nil && os.IsNotExist(err) {
			continue
		}

		// Walk in the testdata dir and generates the deprecated-api custom dashboard for
		// all bundles JSON reports available there
		err = filepath.Walk(pathToWalk, func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() && strings.HasSuffix(info.Name(), "json") &&
				strings.HasPrefix(info.Name(), "bundles") {

				if strings.Contains(info.Name(), "v4.7") ||
					strings.Contains(info.Name(), "v4.6") ||
					strings.Contains(info.Name(), "v4.8") ||
					strings.Contains(info.Name(), "redhat_community_operator_index") ||
					strings.Contains(info.Name(), "bundles_registry.redhat.io_redhat_redhat_operator_index") {
					return nil
				}

				command := exec.Command("go",
					"run",
					"./hack/scripts/ivs-emails/generate.go",
					fmt.Sprintf("--file=%s", filepath.Join(pathToWalk, info.Name())),
					"--mongo=hack/scripts/ivs-emails/mongo-query-join-results-prod.json",
					fmt.Sprintf("--output=%s", ivsDir))
				_, err = pkg.RunCommand(command)
				if err != nil {
					log.Warnf("running command :%s", err)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}

}
