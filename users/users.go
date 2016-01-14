// +build linux

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
	"math"
	"os"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
)

const (
	// Name of plugin
	Name = "users"
	// Version of plugin
	Version = 1
	// Type of plugin
	Type = plugin.CollectorPluginType
)

var nsPrefix = []string{"intel", "utmp", "users"}

const (
	// name of available metrics
	nLoggedUsers    = "logged"
	nLoggedUsersMin = "logged_min"
	nLoggedUsersMax = "logged_max"
	nLoggedUsersAvg = "logged_avg"
)

// Users is the main structure, keeps obtained users statistic
type Users struct {
	data map[string]interface{}
	exec Execution
	avg  average
}

// average keeps items needed to calculate an average
type average struct {
	start time.Time // the start time of average calculation (start collecting)
	now   time.Time // current time
}

// CollectMetrics returns values of desired metrics defined in mts
func (users *Users) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	metrics := make([]plugin.PluginMetricType, len(mts))

	err := users.getUsersStats()
	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()
	for i, m := range mts {
		if v, ok := users.data[parseNamespace(m.Namespace())]; ok {
			metrics[i] = plugin.PluginMetricType{
				Namespace_: m.Namespace(),
				Data_:      v,
				Source_:    hostname,
				Timestamp_: time.Now(),
			}
		}
	}
	return metrics, nil
}

// GetConfigPolicy returns config policy
func (users *Users) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

// GetMetricTypes returns the metric types exposed by snap-plugin-collector-users
func (users *Users) GetMetricTypes(_ plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	mts := []plugin.PluginMetricType{}

	for m := range users.data {
		metric := plugin.PluginMetricType{Namespace_: createNamespace(m)}
		mts = append(mts, metric)
	}
	return mts, nil
}

// New returns snap-plugin-collector-users instance
func New() *Users {
	users := &Users{data: map[string]interface{}{}, avg: average{}, exec: &Executor{}}
	users.init()
	return users
}

// createNamespace returns metric namespace as a slice of strings composed from prefix and metric name
func createNamespace(name string) []string {
	return append(nsPrefix, name)
}

// init initializes users plugin
func (users *Users) init() {
	users.data[nLoggedUsers] = uint64(0)
	users.data[nLoggedUsersMin] = uint64(0)
	users.data[nLoggedUsersMax] = uint64(0)

	// avg logged users metrics initialize as float64
	users.data[nLoggedUsersAvg] = float64(0)
}

// getUsersStats extracts users stats and put them into Users structure
func (users *Users) getUsersStats() error {
	lastLogged, err := users.exec.Execute()
	if err != nil {
		return err
	}

	// set logged users value
	users.data[nLoggedUsers] = lastLogged

	if users.avg.start.Second() == 0 {
		// first getting users stats, initialize values of metrics
		users.data[nLoggedUsersMin] = lastLogged
		users.data[nLoggedUsersMax] = lastLogged
		users.data[nLoggedUsersAvg] = float64(lastLogged)

		// put average structure items
		users.avg.start = time.Now()
		users.avg.now = users.avg.start

		return nil
	}

	// min logged users
	if users.data[nLoggedUsersMin].(uint64) > lastLogged {
		users.data[nLoggedUsersMin] = lastLogged
	}
	// max logged users
	if users.data[nLoggedUsersMax].(uint64) < lastLogged {
		users.data[nLoggedUsersMax] = lastLogged
	}

	// avg logged users, calculating an average value
	last := users.avg.now
	avg := users.data[nLoggedUsersAvg].(float64)
	users.avg.now = time.Now()
	duration := users.avg.now.Sub(users.avg.start).Seconds()
	interval := users.avg.now.Sub(last).Seconds()

	// formula: `avg = (cnt-1)*avg/cnt+val/cnt`, where cnt is the number of measurements
	users.data[nLoggedUsersAvg] = roundToPlaces(avg+(float64(lastLogged)-avg)/(round(duration/interval)+1), 2)

	return nil
}

// parseNamespace performs reverse operation to createNamespace, extracts metric name from namespace
func parseNamespace(ns []string) (name string) {
	if len(ns) > len(nsPrefix) {
		name = ns[len(nsPrefix)]
	}
	return name
}

// round  returns the integral value that is nearest to x
func round(x float64) float64 {
	return math.Floor(x + 0.5)
}

// roundToPlaces returns the value that is nearest to x rounded to the certain decimal places
func roundToPlaces(x float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return round(x*shift) / shift
}
