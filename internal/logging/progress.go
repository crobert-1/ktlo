// Copyright Splunk, Inc.
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

package logging

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const LevelDefault Level = "default"
const LevelDetailed Level = "detailed"
const LevelNone Level = "none"
const PrefixDependencyPRStatus = "dependency-pr-status"
const PrefixFlakyTestAnalysis = "flaky-test-analysis"
const PrefixGenerateReport = "generate-report"
const PrefixVulnerability = "vulnerability"
const progressFormat = "[%s] %s %s\n"
const progressTimestampFormat = "2006-01-02 15:04:05.000 -0700 MST"

type Level string

type Logger struct {
	level  Level
	prefix string
	writer io.Writer
}

func New(writer io.Writer, prefix string, level Level) *Logger {
	return &Logger{level: level, prefix: prefix, writer: writer}
}

func (level Level) IsValid() bool {
	return level == LevelDefault || level == LevelDetailed || level == LevelNone
}

func (level *Level) UnmarshalText(text []byte) error {
	*level = Level(strings.ToLower(strings.TrimSpace(string(text))))
	return nil
}

func (logger *Logger) Progressf(format string, arguments ...any) {
	logger.write(LevelDefault, format, arguments...)
}

func (logger *Logger) Detailedf(format string, arguments ...any) {
	logger.write(LevelDetailed, format, arguments...)
}

func (logger *Logger) write(messageLevel Level, format string, arguments ...any) {
	if logger == nil || logger.writer == nil || logger.level == LevelNone {
		return
	}
	if messageLevel == LevelDetailed && logger.level != LevelDetailed {
		return
	}
	timestamp := time.Now().Format(progressTimestampFormat)
	_, _ = fmt.Fprintf(logger.writer, progressFormat, logger.prefix, timestamp, fmt.Sprintf(format, arguments...))
}
