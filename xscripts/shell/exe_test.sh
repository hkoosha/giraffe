#!/usr/bin/env bash
set -euo pipefail
MY_GIRRAFE_ROOT=$(realpath "$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)/../..")
export MY_GIRRAFE_ROOT
cd "$MY_GIRRAFE_ROOT"
source "$MY_GIRRAFE_ROOT/xscripts/shell/internal/test.sh"


giraffe_test "$@"

echo 'all good'
