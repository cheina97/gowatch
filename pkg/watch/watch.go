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

package watch

import (
	"bufio"
	"os/exec"
	"sync"
	"syscall"

	"github.com/pterm/pterm"
)

// Watcher is the struct which contains the flags data and the patterns to check.
type Watcher struct {
	cfg         *Config
	cmd         *exec.Cmd
	chList      []chan string
	m           sync.RWMutex
	w           sync.WaitGroup
	stop        bool
	stopline    string
	stoppattern string
}

// NewWatcher returns a new Watcher.
func NewWatcher(cfg *Config, cmd *exec.Cmd) *Watcher {
	w := &Watcher{
		cfg:  cfg,
		cmd:  cmd,
		stop: false,
	}
	w.chList = make([]chan string, cfg.Concurrency)
	for i := range w.chList {
		w.chList[i] = make(chan string)
	}
	return w
}

// RunMonitorAndWait runs the monitor routines and waits for the command to finish.
func (w *Watcher) RunMonitorAndWait() error {
	w.initMonitorRoutine()

	r, err := w.cmd.StdoutPipe()
	if err != nil {
		pterm.Error.Printf("Error getting stdout pipe: %s\n", err)
		return err
	}

	if err := w.cmd.Start(); err != nil {
		pterm.Error.Printf("Error starting command: %s\n", err)
		return err
	}

	scanner := bufio.NewScanner(r)
	selectedChannel := 0
	for scanner.Scan() {
		line := scanner.Text()
		if !w.cfg.Quiet {
			pterm.Println(line)
		}
		w.chList[selectedChannel] <- line
		selectedChannel = (selectedChannel + 1) % w.cfg.Concurrency
	}

	for i := range w.chList {
		close(w.chList[i])
	}

	err = w.cmd.Wait()
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			pterm.Warning.Printf("Command exited with error: %s\n", err)
		default:
			pterm.Error.Printf("Error waiting for command: %s\n", err)
		}
	}
	w.w.Wait()
	if w.stop {
		pterm.Warning.Printfln("Pattern %s matched: %s",
			pterm.FgCyan.Sprint(w.stoppattern),
			pterm.FgCyan.Sprint(w.stopline),
		)
		pterm.Warning.Println("Command stopped")
		return nil
	}
	pterm.Success.Println("Command finished")
	return nil
}

// initMonitorRoutine initializes the monitor routines.
func (w *Watcher) initMonitorRoutine() {
	for i := range w.chList {
		w.w.Add(1)
		go monitor(w, i)
	}
}

// monitor is the routine which monitors the output.
func monitor(w *Watcher, index int) {
	pterm.Debug.Printf("Starting monitor routine %d\n", index)
	c := w.chList[index]
	for line := range c {
		w.m.RLock()
		stop := w.stop
		w.m.RUnlock()
		if stop {
			continue
		}
		for i := range w.cfg.Patterns {
			if !w.cfg.Patterns[i].MatchString(line) {
				continue
			}
			w.m.Lock()
			stop = w.stop
			if !stop {
				err := syscall.Kill(w.cmd.Process.Pid, w.cfg.Signal.Signal)
				if err != nil {
					pterm.Error.Printf("Error sending signal to process: %s\n", err)
				}
				w.stopline = line
				w.stoppattern = w.cfg.Patterns[i].String()
				w.stop = true
			}
			w.m.Unlock()
		}
	}
	w.w.Done()
}
