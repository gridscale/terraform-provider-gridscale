TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PLATFORMS=windows linux darwin
ARCHES=amd64
BUILDDIR=build
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=gridscale
VERSION=1.0.0
EXECUTABLE_NAME=terraform-provider-$(PKG_NAME)_v$(VERSION)

default: build

build: fmtcheck
	go install

buildallplatforms:
	$(foreach platform,$(PLATFORMS), \
	    $(foreach arch,$(ARCHES), \
	        GOOS=$(platform) GOARCH=$(arch) mkdir -p $(BUILDDIR)/$(platform)-$(arch); go build -o $(BUILDDIR)/$(platform)-$(arch)/$(EXECUTABLE_NAME);))
	@echo "Renaming Windows file"
	@if [ -f $(BUILDDIR)/windows_amd64/$(EXECUTABLE_NAME) ]; then mv $(BUILDDIR)/windows_amd64/$(EXECUTABLE_NAME) $(BUILDDIR)/windows_amd64/$(EXECUTABLE_NAME).exe; fi

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
	echo "Creating required symlinks in $(GOPATH)/src/$(WEBSITE_REPO)"
	ln -s ../../../../ext/providers/$(PKG_NAME)/website/docs $(GOPATH)/src/$(WEBSITE_REPO)/content/source/docs/providers/$(PKG_NAME)
	ln -s ../../../ext/providers/$(PKG_NAME)/website/$(PKG_NAME).erb $(GOPATH)/src/$(WEBSITE_REPO)/content/source/layouts/$(PKG_NAME).erb
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
	echo "Creating required symlinks in $(GOPATH)/src/$(WEBSITE_REPO)"
	ln -s ../../../../ext/providers/$(PKG_NAME)/website/docs $(GOPATH)/src/$(WEBSITE_REPO)/content/source/docs/providers/$(PKG_NAME)
	ln -s ../../../ext/providers/$(PKG_NAME)/website/$(PKG_NAME).erb $(GOPATH)/src/$(WEBSITE_REPO)/content/source/layouts/$(PKG_NAME).erb
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile website website-test

