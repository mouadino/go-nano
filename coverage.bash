#!/usr/bin/env bash
# Taken from github.com/go-kit/kit

set -e

function go_files { find . -name '*_test.go' ; }
function filter { grep -v '/_' ; }
function remove_relative_prefix { sed -e 's/^\.\///g' ; }

function directories {
  go_files | filter | remove_relative_prefix | while read f
  do
    dirname $f
  done
}

function unique_directories { directories | sort | uniq ; }

PATHS=${1:-$(unique_directories)}

function report {
  for path in $PATHS
  do
    go test -coverprofile=$path/cover.coverprofile ./$path
  done
}

function combine {
  gover
}

function clean {
  find . -name cover.coverprofile | xargs rm
}

report
combine
clean

if [ -n "${WEB+x}" ]
then
  go tool cover -html=gover.coverprofile
fi

