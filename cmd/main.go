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

package main

import (
	"flag"
	"os"
	"os/exec"

	"github.com/pterm/pterm"

	"github.com/cheina97/gowatch/pkg/watch"
)

var version = "development"

func main() {
	cfg := watch.NewConfig()
	flag.StringVar(&cfg.PatternFilePath, "f", "", "Path to the json file which contains the patterns to check")
	flag.BoolVar(&cfg.Quiet, "q", false, "Disable output")
	flag.IntVar(&cfg.Concurrency, "c", 4, "Number of concurrent workers")
	flag.Var(&cfg.Signal, "s", "Signal to send to the command's process when its output matches a pattern")
	v := flag.Bool("v", false, "Show version")
	debug := flag.Bool("d", false, "Enable debug mode")
	flag.Parse()
	if *v {
		pterm.Printfln("gowatch version: %s", version)
	}
	if *debug {
		pterm.EnableDebugMessages()
	}

	if cfg.PatternFilePath == "" {
		pterm.Error.Println("No pattern file specified with -f flag")
		flag.Usage()
		os.Exit(1)
	}
	if err := cfg.ReadPatterns(); err != nil {
		pterm.Error.Printf("Error reading pattern file: %s\n", err)
		os.Exit(1)
	}

	if len(flag.Args()) == 0 {
		pterm.Error.Println("No command specified")
		flag.Usage()
		os.Exit(1)
	}

	//nolint:gosec // This is a command specified by the user
	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)

	w := watch.NewWatcher(cfg, cmd)

	if err := w.RunMonitorAndWait(); err != nil {
		pterm.Error.Println(err.Error())
		os.Exit(1)
	}
}
