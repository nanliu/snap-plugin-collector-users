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
	"errors"
	"math"
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

type mcMock struct {
	mock.Mock
}

func (mc *mcMock) Execute() (uint64, error) {
	args := mc.Called()

	return args.Get(0).(uint64), args.Error(1)
}

var mockMts = []plugin.PluginMetricType{
	plugin.PluginMetricType{Namespace_: []string{"intel", "utmp", "users", "logged"}},
	plugin.PluginMetricType{Namespace_: []string{"intel", "utmp", "users", "logged_avg"}},
	plugin.PluginMetricType{Namespace_: []string{"intel", "utmp", "users", "logged_min"}},
	plugin.PluginMetricType{Namespace_: []string{"intel", "utmp", "users", "logged_max"}},
}

func TestGetConfigPolicy(t *testing.T) {
	usersPlugin := New()

	Convey("getting config policy", t, func() {
		So(func() { usersPlugin.GetConfigPolicy() }, ShouldNotPanic)
		_, err := usersPlugin.GetConfigPolicy()
		So(err, ShouldBeNil)
	})
}

func TestGetMetricTypes(t *testing.T) {
	var cfg plugin.PluginConfigType
	usersPlugin := New()

	Convey("getting exposed metric types", t, func() {
		So(func() { usersPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)

		results, err := usersPlugin.GetMetricTypes(cfg)

		So(err, ShouldBeNil)
		So(len(results), ShouldEqual, 4) // this plugin exposed 4 metrics (fixed numbers of metrics)
	})
}

func TestCollectMetrics(t *testing.T) {
	mockData := []uint64{6, 10, 2}
	usersPlugin := New()

	Convey("successful execution of getting users stat", t, func() {

		Convey("the first data collecting", func() {
			mc := &mcMock{}
			mc.On("Execute").Return(mockData[0], nil)
			usersPlugin.exec = mc

			results, err := usersPlugin.CollectMetrics(mockMts)

			So(err, ShouldBeNil)
			So(len(results), ShouldEqual, len(mockMts))

			for _, metric := range results {
				// min, max and avg number of logged users are equal
				So(metric.Data(), ShouldEqual, mockData[0])
			}
		})

		Convey("the second data collecting", func() {
			mc := &mcMock{}
			mc.On("Execute").Return(mockData[1], nil)
			usersPlugin.exec = mc

			results, err := usersPlugin.CollectMetrics(mockMts)

			So(err, ShouldBeNil)
			So(len(results), ShouldEqual, len(mockMts))

			for _, metric := range results {
				// min, max and avg number of logged users are calculated based on two measurements
				So(checkValues(metric, mockData[:2]), ShouldEqual, true)
			}
		})

		Convey("the third data collecting", func() {
			mc := &mcMock{}
			mc.On("Execute").Return(mockData[2], nil)
			usersPlugin.exec = mc

			results, err := usersPlugin.CollectMetrics(mockMts)

			So(err, ShouldBeNil)
			So(len(results), ShouldEqual, len(mockMts))

			for _, metric := range results {
				// min, max and avg number of logged users are calculated based on three measurements
				So(checkValues(metric, mockData[:3]), ShouldEqual, true)
			}
		})

	})

	Convey("failure execution of getting users stat ", t, func() {
		mc := &mcMock{}
		mc.On("Execute").Return(uint64(0), errors.New("x"))
		usersPlugin.exec = mc

		results, err := usersPlugin.CollectMetrics(mockMts)

		So(err, ShouldNotBeNil)
		So(results, ShouldBeNil)
	})

}

func checkValues(metric plugin.PluginMetricType, mockData []uint64) bool {
	result := false

	// get last namespace's item
	last := len(metric.Namespace()) - 1
	// approximation error, acceptance boundary is (+/-) 0.5
	approxErr := 0.5
	switch metric.Namespace()[last] {
	case nLoggedUsers:
		// take the last set mock data
		if metric.Data() == mockData[len(mockData)-1] {
			result = true
		}
	case nLoggedUsersMin:
		if metric.Data() == min(mockData) {
			result = true
		}

	case nLoggedUsersMax:
		if metric.Data() == max(mockData) {
			result = true
		}
	case nLoggedUsersAvg:
		if math.Abs(metric.Data().(float64)-avg(mockData)) <= approxErr {
			result = true
		}
	}

	return result
}

func min(values []uint64) uint64 {
	min := values[0]
	for _, val := range values[1:] {
		if min > val {
			min = val
		}
	}

	return min
}

func max(values []uint64) uint64 {
	max := values[0]
	for _, val := range values[1:] {
		if max < val {
			max = val
		}
	}

	return max
}

func avg(values []uint64) float64 {
	items := len(values)
	if items == 0 {
		return 0
	}

	sum := values[0]
	for _, val := range values[1:] {
		sum += val
	}

	return float64(sum) / float64(items)
}
