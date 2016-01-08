/*
http://www.apache.org/licenses/LICENSE-2.0.txt
Copyright 2016 Intel Corporation
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package users

import (
	"os/exec"
	"strings"
)

const who = "/usr/bin/who"

// Execution ia an interface and has a single function Execution which returns the number of logged-in users
type Execution interface {
	Execute() (uint64, error)
}

// Executor can execute command `who -u` on local machine
type Executor struct {
}

// Execute executes command `who -u` on local machine (which gives a list of logged-in users) and returns the number of logged-in users
func (e *Executor) Execute() (uint64, error) {

	// execute `who -u`
	outBytes, err := exec.Command(who, "-u").Output()
	if err != nil {
		return 0, err
	}

	// trim white space, specially the last enter in output
	out := strings.TrimSpace(string(outBytes))

	// the number of logged users equals the number of lines in output
	users := len(strings.Split(out, "\n"))

	return uint64(users), nil
}
