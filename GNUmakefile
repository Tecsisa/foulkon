PACKAGES = $(shell go list ./... | grep -v '/vendor/')
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods \
         -nilfunc -printf -rangeloops -shift -structtags -unsafeptr
EXTERNAL_TOOLS=\
	golang.org/x/tools/cmd/cover \
	github.com/axw/gocov/gocov \
	gopkg.in/matm/v1/gocov-html \

GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: test vet

dev: format generate
	@sh -c "'$(PWD)/scripts/build.sh'"

bin: generate
	@sh -c "'$(PWD)/scripts/build.sh'"

release:
	@$(MAKE) bin

cov:
	gocov test ./... | gocov-html > /tmp/coverage.html
	open /tmp/coverage.html

test: generate
	@echo "--> Running go fmt" ;
	@if [ -n "`go fmt ${PACKAGES}`" ]; then \
		echo "[ERR] go fmt updated formatting. Please commit formatted code first."; \
		exit 1; \
	fi
	@sh -c "'$(PWD)/scripts/test.sh'"

cover:
	go list ./... | xargs -n1 go test --cover

format:
	@echo "--> Running go fmt"
	@go fmt $(PACKAGES)

generate:
	@echo "--> Running go generate"
	@go generate $(PACKAGES)

vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@echo "--> Running go tool vet"
	@go tool vet $(VETARGS) ${GOFILES_NOVENDOR} ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "[LINT] Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		echo ""; \
	fi

	@git grep -n `echo "log"".Print"` | grep -v 'vendor/' ; if [ $$? -eq 0 ]; then \
		echo "[LINT] Found "log"".Printf" calls. These should use foulkon's logger instead."; \
		echo ""; \
	fi

# bootstrap the build by downloading additional tools
bootstrap:
	@for tool in $(EXTERNAL_TOOLS) ; do \
		echo "Installing $$tool" ; \
    go get $$tool; \
	done

travis:
	@sh -c "'$(PWD)/scripts/travis.sh'"

.PHONY: all bin cov integ test vet web web-push test-nodep