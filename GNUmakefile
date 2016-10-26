PACKAGES = $(shell go list ./... | grep -v '/vendor/')
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods \
         -nilfunc -printf -rangeloops -shift -structtags -unsafeptr
EXTERNAL_TOOLS=\
	golang.org/x/tools/cmd/cover \
	golang.org/x/tools/cmd/goimports \

GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: test vet

dev: generate format imports
	@FOULKON_DEV=1 sh -c "'$(PWD)/scripts/build.sh'"

deps:
	curl https://glide.sh/get | sh
	glide install

bin: generate format imports
	@sh -c "'$(PWD)/scripts/build.sh'"

demo:
		@sh -c "'$(PWD)/demo/build-demo.sh'"

release:
	@$(MAKE) bin

generate:
	@echo "--> Running go generate"
	@go generate $(PACKAGES)

format:
	@echo "--> Running go fmt" ; \
  if [ -n "`go fmt ${PACKAGES}`" ]; then \
      echo "[FMT] go fmt updated formatting. Please commit formatted code first."; \
      exit 1; \
  fi

imports:
	@echo "--> Running goimports" ; \
  if [ -n "`goimports -l ${GOFILES_NOVENDOR}`" ]; then \
      echo "[IMPORTS] go imports found incorrect imports format. Please check these files:"; \
			echo "`goimports -l ${GOFILES_NOVENDOR}`"; \
      exit 1; \
  fi

test: generate format imports
	@sh -c "'$(PWD)/scripts/test.sh'"

vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@echo "--> Running go tool vet"
	@go tool vet $(VETARGS) ${GOFILES_NOVENDOR} ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "[VET] Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		echo ""; \
	fi

	@git grep -n `echo "log"".Print"` | grep -v 'vendor/' ; if [ $$? -eq 0 ]; then \
		echo "[LINT] Found "log"".Printf" calls. These should use foulkon's logger instead."; \
		echo ""; \
	fi

changelog:
	github_changelog_generator

# bootstrap the build by downloading additional tools
bootstrap: deps
	@for tool in $(EXTERNAL_TOOLS) ; do \
		echo "Installing $$tool" ; \
    go get $$tool; \
	done

travis:
	@sh -c "'$(PWD)/scripts/travis.sh'"

.PHONY: all dev deps bin demo release generate format imports test vet bootstrap travis
