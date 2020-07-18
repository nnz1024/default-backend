// +build template

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

// Enums in Go looks strange at the first time
const (
	// No variable escaping at template substitution
	escapeNone = iota
	// Use html/template with semantical escaping
	escapeHTML
	// Use text/template, but escape vars with JSEscapeString()
	escapeJSON
)

// TemplatingSettings struct holding templating settings
// for given MIME type
type TemplatingSettings struct {
	// Name of the envirement variable which enables
	// templating for this MIME type
	enableVar string
	// Type of template escaping (see escapeXXX consts definition)
	escapeType int
}

// Templating is a list of template formats and their settings
// Key is the MIME type, and value is its templating settings
var Templating = map[string]TemplatingSettings{
	"text/plain": {
		enableVar:  "TEMPLATE_TEXT",
		escapeType: escapeNone,
	},
	"text/html": {
		enableVar:  "TEMPLATE_HTML",
		escapeType: escapeHTML,
	},
	"application/json": {
		enableVar:  "TEMPLATE_JSON",
		escapeType: escapeJSON,
	},
}
