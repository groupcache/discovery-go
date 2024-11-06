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

package logger

import "io"

// Logger represents an active logging object that generates lines of
// output to an io.Writer.
type Logger interface {
	// Info starts a new message with info level.
	Info(...any)
	// Infof starts a new message with info level.
	Infof(string, ...any)
	// Warn starts a new message with warn level.
	Warn(...any)
	// Warnf starts a new message with warn level.
	Warnf(string, ...any)
	// Error starts a new message with error level.
	Error(...any)
	// Errorf starts a new message with error level.
	Errorf(string, ...any)
	// Fatal starts a new message with fatal level. The os.Exit(1) function
	// is called which terminates the program immediately.
	Fatal(...any)
	// Fatalf starts a new message with fatal level. The os.Exit(1) function
	// is called which terminates the program immediately.
	Fatalf(string, ...any)
	// Panic starts a new message with panic level. The panic() function
	// is called which stops the ordinary flow of a goroutine.
	Panic(...any)
	// Panicf starts a new message with panic level. The panic() function
	// is called which stops the ordinary flow of a goroutine.
	Panicf(string, ...any)
	// Debug starts a new message with debug level.
	Debug(...any)
	// Debugf starts a new message with debug level.
	Debugf(string, ...any)
	// LogOutput returns the log output that is set
	LogOutput() []io.Writer
}
