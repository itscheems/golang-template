// Copyright 2026 itscheems
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

// Package log is a thin wrapper around the stdlib slog package that
// stamps every record with the calling service's name. It exists to
// demonstrate cross-module imports inside the workspace and gives
// downstream services a single, opinionated logger entry point.
package log

import (
	"log/slog"
	"os"
)

// New returns a structured logger tagged with the given service name.
func New(service string) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", service)
}
