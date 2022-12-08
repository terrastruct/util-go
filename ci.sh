#!/bin/sh
set -eu
if [ ! -e "$(dirname "$0")/ci/.git" ]; then
  set -x
  git submodule update --init
  set +x
fi
. "$(dirname "$0")/ci/lib.sh"
cd "$(dirname "$0")"

job_parseflags "$@"
runjob fmt ./ci/bin/fmt.sh &
runjob lint ci_go_lint &
runjob build 'go build ./...' &
runjob test 'go test ./...' &
ci_waitjobs
