/*
Copyright 2017 The Kubernetes Authors All rights reserved.

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
package util

import (
	"fmt"
	"github.com/golang/glog"
	"math/rand"
	"regexp"
	"runtime"
	"time"

	"k8s.io/node-problem-detector/pkg/types"
)

var PodOOMRegex = regexp.MustCompile(`node:(\S+)\s+pod:(\S+)\s+namespace:(\S+)\s+uuid:(\S+)`)

// GenerateConditionChangeEvent generates an event for condition change.
func GenerateConditionChangeEvent(t string, status types.ConditionStatus, reason string, timestamp time.Time, msg string) types.Event {
	return types.Event{
		Severity:  types.Info,
		Timestamp: timestamp,
		Reason:    reason,
		Message:   fmt.Sprintf("Node condition %s is now: %s, reason: %s, msg: %s", t, status, reason, msg),
	}
}

func GetStartTime(now time.Time, uptimeDuration time.Duration, lookbackStr string, delayStr string) (time.Time, error) {
	startTime := now.Add(-uptimeDuration)

	// Delay startTime if delay duration is set, so that the log watcher can skip
	// the logs in delay duration and wait until the node is stable.
	if delayStr != "" {
		delay, err := time.ParseDuration(delayStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to parse delay duration %q: %v", delayStr, err)
		}
		// Notice that when delay > uptime, startTime is actually after now, which is fine.
		startTime = startTime.Add(delay)
	}

	// Addjust startTime according to lookback duration
	lookbackStartTime := now
	if lookbackStr != "" {
		lookback, err := time.ParseDuration(lookbackStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to parse lookback duration %q: %v", lookbackStr, err)
		}
		lookbackStartTime = now.Add(-lookback)
	}
	if startTime.Before(lookbackStartTime) {
		startTime = lookbackStartTime
	}

	return startTime, nil
}

func Recovery() error {
	if err := recover(); err != nil {
		var e error
		switch r := err.(type) {
		case error:
			e = r
		default:
			e = fmt.Errorf("%v", r)
		}
		stack := make([]byte, 2048)
		length := runtime.Stack(stack, true)
		glog.Errorf("[%s] %s %s\n", "PANIC RECOVER", e, stack[:length])
		return e
	}
	return nil
}

func RandomDurationMinute(seedMinutes int64) time.Duration {
	rand.Seed(time.Now().Unix())
	return time.Duration(rand.Int63n(seedMinutes) * int64(time.Minute))
}
