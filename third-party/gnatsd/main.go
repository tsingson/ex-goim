// Copyright 2012-2018 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"runtime"
	"time"

	gnatsd "github.com/nats-io/gnatsd/server"
)

func main() {
	runtime.MemProfileRate = 0
	runtime.GOMAXPROCS(128)
	signal := make(chan struct{})
	s:= RunDefaultServer()
	s.Start()
	defer s.Shutdown()
	<-signal
}

// So we can pass tests and benchmarks..
type tLogger interface {
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// DefaultTestOptions are default options for the unit tests.
var DefaultTestOptions = &gnatsd.Options{
	Host:           "127.0.0.1",
	Port:           4222,
	NoLog:          true,
	NoSigs:         true,
	MaxControlLine: 2048,
}

// RunDefaultServer starts a new Go routine based gnatsd using the default options
func RunDefaultServer() *gnatsd.Server {
	return RunServer(DefaultTestOptions)
}

// To turn on gnatsd tracing and debugging and logging which are
// normally suppressed.
var (
	doLog   = false
	doTrace = false
	doDebug = false
)

// RunServer starts a new Go routine based gnatsd
func RunServer(opts *gnatsd.Options) *gnatsd.Server {
	// if opts == nil {
	// 	opts = &DefaultTestOptions
	// }
	// Optionally override for individual debugging of tests
	// opts.NoLog = !doLog
	// opts.Trace = doTrace
	opts.Debug = true

	s := gnatsd.New(opts)
	// if err != nil || s == nil {
	// 	panic(fmt.Sprintf("No NATS Server object returned: %v", err))
	// }

	if doLog {
		s.ConfigureLogger()
	}

	// Run gnatsd in Go routine.
	go s.Start()

	// Wait for accept loop(s) to be started
	if !s.ReadyForConnections(10 * time.Second) {
		panic("Unable to start NATS Server in Go Routine")
	}
	return s
}

// LoadConfig loads a configuration from a filename
func LoadConfig(configFile string) *gnatsd.Options {
	opts, err := gnatsd.ProcessConfigFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("Error processing configuration file: %v", err))
	}
	return opts
}

// RunServerWithConfig starts a new Go routine based gnatsd with a configuration file.
func RunServerWithConfig(configFile string) (srv *gnatsd.Server, opts *gnatsd.Options) {
	opts = LoadConfig(configFile)
	srv = RunServer(opts)
	return
}
