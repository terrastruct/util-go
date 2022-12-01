#!/bin/sh
set -eu
if [ ! -d "$(dirname "$0")/ci/sub/.git" ]; then
  git submodule update --init
fi
. "$(dirname "$0")/ci/sub/lib.sh"
cd "$(dirname "$0")"

sh_c detect_changed_files

job_parseflags "$@"
runjob fmt ci_go_fmt &
runjob lint ci_go_lint &
runjob build ci_go_build &
runjob test ci_go_test &
ci_waitjobs
