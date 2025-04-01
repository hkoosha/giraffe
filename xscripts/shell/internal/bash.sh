function _giraffe_die() {
    echo "[error]: $*" >&2
    exit 1
}

function _giraffe_debug() {
    echo "[debug]: $*" > /dev/null
}

function _giraffe_kernel() {
  local my_kernel

  if [[ $# -eq 0 ]]; then
    my_kernel="$(uname -s)"
  elif [[ $# -eq 1 ]]; then
    my_kernel="$1"
  else
    _giraffe_die 'too many args to _giraffe_kernel'
  fi

  my_kernel="$(echo "$my_kernel" | tr '[:upper:]' '[:lower:]')"

  case "$my_kernel" in
    'linux' | 'darwin')
      echo "$my_kernel"
      ;;

    *)
      _giraffe_die "unknown kernel: $my_kernel"
      ;;
  esac
}

function _giraffe_machine() {
  local my_machine

  if [[ $# -eq 0 ]]; then
    my_machine="$(uname -m)"
  elif [[ $# -eq 1 ]]; then
    my_machine="$1"
  else
    _giraffe_die 'too many args to _giraffe_machine'
  fi

  my_machine="$(echo "$my_machine" | tr '[:upper:]' '[:lower:]')"

  case "$my_machine" in
    'x86_64' | 'amd64')
      echo 'amd64'
      ;;

    'aarch64' | 'arm64')
      echo 'arm64'
      ;;

    *)
      _giraffe_die "unknown machine: $my_machine"
      ;;
  esac
}

# Source in order of dependency:

source "$MY_GIRRAFE_ROOT/xscripts/shell/internal/git.sh"

source "$MY_GIRRAFE_ROOT/xscripts/shell/internal/tools.sh"
source "$MY_GIRRAFE_ROOT/xscripts/shell/internal/build.sh"
source "$MY_GIRRAFE_ROOT/xscripts/shell/internal/lint.sh"
source "$MY_GIRRAFE_ROOT/xscripts/shell/internal/test.sh"

