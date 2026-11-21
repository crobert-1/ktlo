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
	"time"
)

const PrefixDependencyPRStatus = "dependency-pr-status"
const PrefixFlakyTestAnalysis = "flaky-test-analysis"
const PrefixGenerateReport = "generate-report"
const PrefixVulnerability = "vulnerability"
const progressFormat = "[%s] %s %s\n"
const progressTimestampFormat = "2006-01-02 15:04:05.000 -0700 MST"

func Progressf(writer io.Writer, prefix string, format string, arguments ...any) {
	if writer == nil {
		return
	}
	timestamp := time.Now().Format(progressTimestampFormat)
	_, _ = fmt.Fprintf(writer, progressFormat, prefix, timestamp, fmt.Sprintf(format, arguments...))
}
