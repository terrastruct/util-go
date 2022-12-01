#!/bin/sh
set -eu
if [ ! -e "$(dirname "$0")/ci/.git" ]; then
  set -x
  git submodule update --init
  set +x
fi
. "$(dirname "$0")/ci/lib.sh"
cd "$(dirname "$0")"

sh_c detect_changed_files

fmt() {
  if is_changed README.md; then
    sh_c tocsubst --skip 2 README.md
  fi
  ci_go_fmt
}

job_parseflags "$@"
runjob fmt fmt &
runjob lint ci_go_lint &
runjob build ci_go_build &
runjob test ci_go_test &
ci_waitjobs
