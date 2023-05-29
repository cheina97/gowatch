// Copyright 2023-2023 cheina97
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package watch contains the Watcher struct and the Config struct.
package watch

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
	"syscall"

	"github.com/pterm/pterm"
)

// Config contains the flags data and the patterns to check.
type Config struct {
	// Concurrency is the number of concurrent workers.
	Concurrency int
	// Quiet disables the output of the program.
	Quiet bool
	// PatternFilePath is the path to the json file which contains the patterns to check.
	PatternFilePath string
	// ConfigPatterns is the struct which contains the patterns to check.
	Patterns []regexp.Regexp
	// Signal is the signal to send to the command's process when its output matches a pattern.
	Signal Signal
}

// NewConfig returns a new Config.
func NewConfig() *Config {
	return &Config{
		Patterns: []regexp.Regexp{},
		Signal: Signal{
			syscall.SIGKILL,
		},
	}
}

// ReadPatterns reads the patterns from the file specified in PatternFilePath.
func (cfg *Config) ReadPatterns() error {
	jsonFile, err := os.Open(cfg.PatternFilePath)
	if err != nil {
		return err
	}

	defer func() {
		if err := jsonFile.Close(); err != nil {
			pterm.Error.Printf("Error closing file: %s\n", err)
		}
	}()

	b, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	patterns := []string{}
	err = json.Unmarshal(b, &patterns)
	if err != nil {
		return err
	}

	for _, p := range patterns {
		r, err := regexp.Compile(p)
		if err != nil {
			return err
		}
		cfg.Patterns = append(cfg.Patterns, *r)
	}

	return nil
}
