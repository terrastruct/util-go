#!/bin/sh
set -eu
if [ ! -e "$(dirname "$0")/ci/sub/.git" ]; then
  set -x
  git submodule update --init
  set +x
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
