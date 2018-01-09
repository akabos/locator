SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

GOPKG_TOML := $(SOURCEDIR)/Gopkg.toml
GOPKG_LOCK := $(SOURCEDIR)/Gopkg.lock
VENDORDIR ?= $(SOURCEDIR)/vendor

GO ?= $(shell which go)
GOOS ?= $(shell ${GO} env GOOS)
GOARCH ?= $(shell ${GO} env GOARCH)

GOPKG ?= $(shell which dep)

BINSUFFIX ?= $(GOOS)_$(GOARCH)
BINPREFIX ?= locator
BIN := $(SOURCEDIR)/$(BINPREFIX).$(BINSUFFIX)

.PHONY: all
all: $(BIN)

$(GOPKG_LOCK): $(GOPKG_TOML)
	@echo Ensure deps
	@$(GOPKG) ensure -v
	@touch $@

$(BIN): $(GOPKG_LOCK) $(SOURCES)
	@echo Build $@
	@$(GO) build -o $@

.PHONY: dep
dep: $(GOPKG_LOCK)

.PHONY: clean
clean:
	@$(GO) clean
	@rm -f $(BINPREFIX).*

.PHONY: distclean
distclean: clean
	@rm -f $(GLIDELOCK)
	@rm -rf $(VENDORDIR)
