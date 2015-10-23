# -*- Makefile -*-

VERSION = 0.1.0dev

PREFIX ?= /usr/local

BUILD_NAME = goverview-$(GOOS)-$(GOARCH)
BUILD_DIR = build/$(BUILD_NAME)

DIST_DIR = dist

TAR ?= tar

GO ?= go

GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)

GOPATH = $(realpath $(BUILD_DIR)):$(realpath vendor)
export GOPATH

GOLINT ?= golint
GOLINT_ARGS =

JSHINT ?= jshint
JSHINT_ARGS =

UGLIFYJS ?= uglifyjs
UGLIFYJS_ARGS = --comments --compress --mangle --screw-ie8

UGLIFYCSS ?= uglifycss
UGLIFYCSS_ARGS =

MYTH ?= myth
MYTH_ARGS =

PANDOC ?= pandoc
PANDOC_ARGS = --standalone --to man

mesg_start = echo "$(shell tty -s && tput setaf 4)$(1):$(shell tty -s && tput sgr0) $(2)"
mesg_step = echo "$(1)"
mesg_ok = echo "result: $(shell tty -s && tput setaf 2)ok$(shell tty -s && tput sgr0)"
mesg_fail = (echo "result: $(shell tty -s && tput setaf 1)fail$(shell tty -s && tput sgr0)" && false)

SCRIPT_OUTPUT = static/js/goverview.min.js
STYLE_OUTPUT = static/css/style.min.css
MAN_OUTPUT = $(BUILD_DIR)/man/man1/goverview.1

all: build

clean:
	@$(call mesg_start,clean,Cleaning build directory...)
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(SCRIPT_OUTPUT) $(STYLE_OUTPUT) cmd/goverview/data.go && \
		$(call mesg_ok) || $(call mesg_fail)

build: build-script build-style build-doc build-bin

build-script: $(SCRIPT_OUTPUT)

build-style: $(STYLE_OUTPUT)

build-doc: $(MAN_OUTPUT)

build-bin: $(BUILD_DIR)/bin/goverview

install: install-doc install-bin

install-doc: build-doc
	@$(call mesg_start,install,Installing manuals files...)
	@install -d -m 0755 $(PREFIX)/share/man && cp -Rp $(BUILD_DIR)/man/man1 $(PREFIX)/share/man && \
		$(call mesg_ok) || $(call mesg_fail)
	@$(call mesg_start,install,Installing examples files...)
	@install -d -m 0755 $(PREFIX)/share/goverview/examples && cp -Rp docs/examples $(PREFIX)/share/goverview && \
		$(call mesg_ok) || $(call mesg_fail)

install-bin: build-bin
	@$(call mesg_start,install,Installing binaries...)
	@install -d -m 0755 $(PREFIX)/bin && cp $(BUILD_DIR)/bin/* $(PREFIX)/bin && \
		$(call mesg_ok) || $(call mesg_fail)

lint: clean lint-script lint-bin

lint-script:
	@$(call mesg_start,lint,Checking scripts with JSHint...)
	-@$(JSHINT) $(JSHINT_ARGS) src/js/goverview.js && \
		$(call mesg_ok) || $(call mesg_fail)

lint-bin:
	@$(call mesg_start,lint,Checking sources with Golint...)
	@$(GOLINT) $(GOLINT_ARGS) cmd/... && $(GOLINT) $(GOLINT_ARGS) pkg/... && \
		$(call mesg_ok) || $(call mesg_fail)

data:
	@$(call mesg_start,build,Generating embeddable data...)
	@go-bindata -nomemcopy -prefix=static/ -o cmd/goverview/data.go static/* && \
		$(call mesg_ok) || $(call mesg_fail)

$(SCRIPT_OUTPUT): src/js/goverview.js
	@$(call mesg_start,static,Packing script file...)
	@install -d -m 0755 $(dir $(SCRIPT_OUTPUT)) && \
		$(UGLIFYJS) $(UGLIFYJS_ARGS) --output $(SCRIPT_OUTPUT) src/js/goverview.js && \
		$(call mesg_ok) || $(call mesg_fail)

$(STYLE_OUTPUT): src/css/style.css
	@$(call mesg_start,static,Applying vendor-specific rules to style file...)
	@install -d -m 0755 $(BUILD_DIR)/tmp && install -d -m 0755 $(dir $(STYLE_OUTPUT)) && \
		$(MYTH) $(MYTH_ARGS) src/css/style.css >$(BUILD_DIR)/tmp/style.css && \
		$(call mesg_ok) || $(call mesg_fail)
	@$(call mesg_start,static,Packing style file...)
	@$(UGLIFYCSS) $(UGLIFYCSS_ARGS) $(BUILD_DIR)/tmp/style.css >$(STYLE_OUTPUT) && \
		$(call mesg_ok) || $(call mesg_fail)

$(MAN_OUTPUT): src/man/goverview.1.md
	@$(call mesg_start,docs,Generating manual page...)
	@install -d -m 0755 $(dir $@) && $(PANDOC) $(PANDOC_ARGS) src/man/goverview.1.md >$@ && \
		$(call mesg_ok) || $(call mesg_fail)

$(BUILD_DIR)/bin/%: $(BUILD_DIR)/src/github.com/vbatoufflet/goverview data
	@$(call mesg_start,build,Building $(notdir $@)...)
	@install -d -m 0755 $(BUILD_DIR)/bin && $(GO) build \
			-ldflags "-X main.version=$(VERSION)" \
			-o $@ cmd/$(notdir $@)/*.go && \
		$(call mesg_ok) || $(call mesg_fail)

$(BUILD_DIR)/src/github.com/vbatoufflet/goverview:
	@$(call mesg_start,build,Creating source symlink...)
	@install -m 0755 -d $(BUILD_DIR)/src/github.com/vbatoufflet && \
		ln -s ../../../../.. $(BUILD_DIR)/src/github.com/vbatoufflet/goverview && \
		$(call mesg_ok) || $(call mesg_fail)

.PHONY: dist
dist: build
	@$(call mesg_start,dist,Creating distribution directory...)
	@install -d -m 0755 $(DIST_DIR)/$(BUILD_NAME) && \
		$(call mesg_ok) || $(call mesg_fail)
	@$(call mesg_start,dist,Installing files into distribution directory...)
	@cp $(BUILD_DIR)/bin/* $(BUILD_DIR)/man/*/* README.md LICENSE $(DIST_DIR)/$(BUILD_NAME)/ && \
		$(call mesg_ok) || $(call mesg_fail)
	@$(call mesg_start,dist,Build distribution tarball...)
	@$(TAR) -C $(DIST_DIR) -czf $(DIST_DIR)/$(BUILD_NAME:goverview-%=goverview-$(VERSION)-%).tar.gz $(BUILD_NAME) && \
		$(call mesg_ok) || $(call mesg_fail)
