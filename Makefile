# Copyright (c) 2013-2024 by Michael Dvorkin and contributors. All Rights Reserved.
# Use of this source code is governed by a MIT-style license that can
# be found in the LICENSE file.

VERSION = 1.0.0-patched

run:
	go run ./main.go

build:
	go build -x -o ./mop-$(shell go env GOOS)-$(shell go env GOARCH)

install:
	go install -x

buildall:
	GOOS=linux   GOARCH=386   go build $(GOFLAGS) -o ./mop-$(VERSION)-linux-32
	GOOS=linux   GOARCH=amd64 go build $(GOFLAGS) -o ./mop-$(VERSION)-linux-64
