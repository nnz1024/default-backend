/*
Copyright 2017 The Kubernetes Authors.
Copyright 2020 Sergey Ptashnik

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

package main

const (
	// ErrFilesPathVar is the name of the environment variable indicating
	// the location on disk of files served by the handler
	ErrFilesPathVar = "ERROR_FILES_PATH"

	// DefaultFormat is default response format
	DefaultFormat = "text/html"

	// ListenPort is the TCP port which default-backend will listen
	ListenPort = 8080
)

// The simplest way to create immutable array
func defaultCext() []string {
	// Must correspond to DefaultFormat
	return []string{".html", ".htm"}
}

// Headers is a map which will be passed to template engine
// Also used for references in the internal logic and debug purposes
var Headers = map[string]string{
	"FormatHeader": "X-Format",
	"CodeHeader":   "X-Code",
	"ContentType":  "Content-Type",
	"OriginalURI":  "X-Original-URI",
	"Namespace":    "X-Namespace",
	"IngressName":  "X-Ingress-Name",
	"ServiceName":  "X-Service-Name",
	"ServicePort":  "X-Service-Port",
	"RequestID":    "X-Request-ID",
}
