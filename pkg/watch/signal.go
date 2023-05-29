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
	"fmt"
	"strconv"
	"syscall"

	"github.com/pterm/pterm"
)

// Signal is a wrapper which implement the flag.Value interface.
type Signal struct {
	syscall.Signal
}

// Set implements the flag.Value interface.
func (s *Signal) Set(value string) error {
	pterm.Debug.Printf("Signal.Set(%s)\n", value)
	sn, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	if sn < 0 || sn >= 64 {
		return fmt.Errorf("invalid signal number %d", sn)
	}

	s.Signal = syscall.Signal(sn)
	return nil
}

func (s *Signal) String() string {
	return s.Signal.String()
}
