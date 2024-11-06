/*
 * Copyright 2024 Arsene Tochemey Gandote
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package discoverygo

import (
	"bytes"
	"io"
	"regexp"

	"github.com/groupcache/discovery-go/logger"
)

// logWriter is a wrapper for membership logging
type logWriter struct {
	logger logger.Logger
	info   *regexp.Regexp
	debug  *regexp.Regexp
	warn   *regexp.Regexp
	error  *regexp.Regexp
}

// make sure that the logWriter implements the io.Writer interface fully
var _ io.Writer = (*logWriter)(nil)

// newLogWriter create an instance of logWriter
func newLogWriter(logger logger.Logger) *logWriter {
	return &logWriter{
		logger: logger,
		info:   regexp.MustCompile(`\[INFO\] (.+)`),
		debug:  regexp.MustCompile(`\[DEBUG\] (.+)`),
		warn:   regexp.MustCompile(`\[WARN\] (.+)`),
		error:  regexp.MustCompile(`\[ERROR\] (.+)`),
	}
}

// Write writes len(p) bytes from p to the underlying data stream.
func (l *logWriter) Write(message []byte) (n int, err error) {
	// trim all spaces
	trimmed := bytes.TrimSpace(message)
	// get the text value of the log message
	text := string(trimmed)

	// parse info message
	matches := l.info.FindStringSubmatch(text)
	if len(matches) > 1 {
		// info message found
		infoText := matches[1]
		l.logger.Info(infoText)
		return len(message), nil
	}

	// parse debug message
	matches = l.debug.FindStringSubmatch(text)
	if len(matches) > 1 {
		// debug message found
		debugText := matches[1]
		l.logger.Debug(debugText)
		return len(message), nil
	}

	// parse warn messages
	matches = l.warn.FindStringSubmatch(text)
	if len(matches) > 1 {
		// debug message found
		warnText := matches[1]
		l.logger.Warn(warnText)
		return len(message), nil
	}

	// parse error messages
	matches = l.error.FindStringSubmatch(text)
	if len(matches) > 1 {
		// error message found
		errorText := matches[1]
		l.logger.Error(errorText)
		return len(message), nil
	}

	return len(message), nil
}
