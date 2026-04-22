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

// Package version exposes build-time metadata stamped via -ldflags.
package version

// These vars are populated at build time via -ldflags by GoReleaser /
// the Taskfile build recipe. Defaults are set so `go run` works locally.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
