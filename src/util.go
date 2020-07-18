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

import (
	"log"
	"mime"
	"strconv"
)

func getFormatInfo(formatHeader string) (format string, cext []string) {
	format = DefaultFormat
	cext = defaultCext()

	if formatHeader == "" {
		log.Printf("Format not specified. Using %v", format)
		return
	}

	newCext, err := mime.ExtensionsByType(formatHeader)
	if err != nil {
		log.Printf("Error reading media type extension: %v. Using %v", err, cext)
		return
	}
	if len(newCext) == 0 {
		log.Printf("Couldn't get media type extension. Using %v", cext)
		return
	}

	format = formatHeader
	cext = newCext
	return
}

func getErrorCode(codeHeader string) (code int) {
	code = 404

	if codeHeader == "" {
		log.Printf("No HTTP code provided. Using %v", code)
		return
	}

	newCode, err := strconv.Atoi(codeHeader)
	if err != nil {
		log.Printf("Error reading HTTP code: %v. Using %v", err, code)
		return
	}
	if (newCode >= 600) || (newCode < 100) {
		log.Printf("Invalid HTTP code: %v. Using %v", newCode, code)
		return
	}

	code = newCode
	return
}
