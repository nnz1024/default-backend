// +build !template

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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func errorHandler(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var file string
		var err error
		var f *os.File
		start := time.Now()
		count := true

		if os.Getenv("DEBUG") != "" {
			for _, header := range Headers {
				w.Header().Set(header, r.Header.Get(header))
			}
		}

		format, cext := getFormatInfo(r.Header.Get(Headers["FormatHeader"]))
		code := getErrorCode(r.Header.Get(Headers["CodeHeader"]))

		for _, ext := range cext {
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			file = fmt.Sprintf("%v/%v%v", path, code, ext)
			f, err = os.Open(file)
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Printf("Error opening file for code %v and extensions %v: %v", code, cext, err)
			// scode := strconv.Itoa(code)
			scode := code / 100
			for _, ext := range cext {
				if !strings.HasPrefix(ext, ".") {
					ext = "." + ext
				}
				file = fmt.Sprintf("%v/%dxx%v", path, scode, ext)
				f, err = os.Open(file)
				if err == nil {
					break
				}
			}
			if err != nil {
				log.Printf("Error opening file for code %dxx and extensions %v: %v", scode, cext, err)
				http.NotFound(w, r)
				return
			}
			// Do not count Nxx requests
			count = false
		}
		defer f.Close()
		log.Printf("Serving custom error response for code %v and format %v from file %v", code, format, file)

		w.WriteHeader(code)
		w.Header().Set(Headers["ContentType"], format)
		io.Copy(w, f)

		if count {
			duration := time.Since(start).Seconds()

			proto := fmt.Sprintf("%s.%s", strconv.Itoa(r.ProtoMajor), strconv.Itoa(r.ProtoMinor))

			requestCount.WithLabelValues(proto).Inc()
			requestDuration.WithLabelValues(proto).Observe(duration)
		}
	}
}
