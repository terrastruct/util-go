#!/bin/sh
set -eu
if [ ! -e "$(dirname "$0")/ci/.git" ]; then
  set -x
  git submodule update --init
  set +x
fi
. "$(dirname "$0")/ci/lib.sh"
cd "$(dirname "$0")"

if [ -n "${CI-}" ]; then
  go install golang.org/x/tools/cmd/goimports@v0.4.0
  npm install -g prettier@2.8.1
fi

job_parseflags "$@"
runjob fmt ./ci/bin/fmt.sh &
runjob lint ci_go_lint &
runjob build 'go build ./...' &
runjob test 'go test ./...' &
ci_waitjobs
