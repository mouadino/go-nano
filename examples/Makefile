.PHONY: build

# Prepend our _vendor directory to the system GOPATH
# so that import path resolution will prioritize
# our third party snapshots.
GOPATH := ${PWD}/_vendor:${GOPATH}
export GOPATH

default: build

build:
		go build -v -o _build/upper ./upper
		go build -v -o _build/demo ./demo

