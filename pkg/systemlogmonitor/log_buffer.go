/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package systemlogmonitor

import (
	"regexp"
	"strings"
	"time"

	"k8s.io/node-problem-detector/pkg/systemlogmonitor/types"
)

// LogBuffer buffers the logs and supports match in the log buffer with regular expression.
type LogBuffer interface {
	// Push pushes log into the log buffer. It will panic if the buffer is full.
	Push(*types.Log)
	// Poll poll log from the log buffer. It returns nil if no log is available.
	Poll() *types.Log
	// Match with regular expression in the log buffer.
	Match(string) []*types.Log

	Clean()

	SetLookback(lookback *time.Duration)
}

type logBuffer struct {
	// buffer is a simple ring buffer.
	buffer   []*types.Log
	msg      []string
	max      int
	write    int
	read     int
	size     int
	lookback *time.Duration
}

// NewLogBuffer creates log buffer with max line number limit. Because we only match logs
// in the log buffer, the max buffer line number is also the max pattern line number we
// support. Smaller buffer line number means less memory and cpu usage, but also means less
// lines of patterns we support.
func NewLogBuffer(maxLines int) *logBuffer {
	return &logBuffer{
		buffer: make([]*types.Log, maxLines, maxLines),
		msg:    make([]string, maxLines, maxLines),
		max:    maxLines,
		write:  0,
		read:   0,
		size:   0,
	}
}

func (b *logBuffer) SetLookback(lookback *time.Duration) {
	b.lookback = lookback
}

func (b *logBuffer) Push(log *types.Log) {
	if b.size == b.max {
		b.read = (b.read + 1) % b.max
		b.size--
	}
	b.buffer[b.write] = log
	b.msg[b.write] = log.Message
	b.write = (b.write + 1) % b.max
	b.size++
}

func (b *logBuffer) Poll() *types.Log {
	if b.size == 0 {
		return nil
	}
	item := b.buffer[b.read]
	b.read = (b.read + 1) % b.max
	b.size--
	return item
}

func (b *logBuffer) Clean() {
	b.read = b.write
	b.size = 0
}

func (b *logBuffer) IsEmpty() bool {
	return b.size == 0
}

func (b *logBuffer) IsFull() bool {
	return b.size == b.max
}

// TODO(random-liu): Cache regexp if garbage collection becomes a problem someday.
func (b *logBuffer) Match(expr string) []*types.Log {
	// The expression should be checked outside
	reg := regexp.MustCompile(expr + `\z`)
	var matched []*types.Log

	for i := b.read; i != b.write; i = (i + 1) % b.max {
		if b.lookback != nil && b.buffer[i] != nil && b.buffer[i].Timestamp.Before(time.Now().Add(-*b.lookback)) {
			// not ontime log
			continue
		}
		if b.buffer[i] != nil && reg.MatchString(b.buffer[i].Message) {
			matched = append(matched, b.buffer[i])
		}
	}
	return matched
}

// concatLogs concatenates multiple lines of logs into one string.
func concatLogs(logs []string) string {
	return strings.Join(logs, "\n")
}
