GO ?= go
GOEXPERIMENT_JSONV2 ?= jsonv2
PKGS ?= ./...
GOFILES := $(shell find . -name '*.go' -not -path './.git/*' -not -path './research/*')

.PHONY: fmt fmt-check vet build test test-nojsonv2 bench fuzz ci clean

fmt:
	@$(GO) fmt $(PKGS)
	@gofmt -w $(GOFILES)

fmt-check:
	@out=$$(gofmt -l $(GOFILES)); \
	if [ -n "$$out" ]; then \
		echo "Unformatted files:"; \
		echo "$$out"; \
		exit 1; \
	fi

vet:
	@GOEXPERIMENT=$(GOEXPERIMENT_JSONV2) $(GO) vet $(PKGS)

build:
	@GOEXPERIMENT=$(GOEXPERIMENT_JSONV2) $(GO) build $(PKGS)

test:
	@GOEXPERIMENT=$(GOEXPERIMENT_JSONV2) $(GO) test $(PKGS)

test-nojsonv2:
	@GOEXPERIMENT= $(GO) test $(PKGS)

bench:
	@GOEXPERIMENT=$(GOEXPERIMENT_JSONV2) $(GO) test -run '^$$' -bench BenchmarkUnmarshalADF -benchmem $(PKGS)

fuzz:
	@GOEXPERIMENT=$(GOEXPERIMENT_JSONV2) $(GO) test -fuzz=FuzzUnmarshalADF -run='^$$' $(PKGS)

ci: fmt-check vet build test test-nojsonv2

clean:
	@$(GO) clean -cache -testcache
