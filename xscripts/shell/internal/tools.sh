function _giraffe_go_tools_export() {
  GOBIN="${MY_GIRRAFE_ROOT}/build/$(uname -s)"
  PATH="$GOBIN:$PATH"

  export GOBIN
  export PATH
}

function _giraffe_go_tools_install() {
  [[ $# -lt 1 ]] && _giraffe_die 'bad args'

  local -r repo="$1"
  shift

  if [[ $# -eq 0 ]]; then
    local -r gobin="${MY_GIRRAFE_ROOT}/build/$(uname -s)"
  elif [[ $# -eq 1 ]]; then
    local -r gobin="$1"
  else
    _giraffe_die 'bad args'
  fi

  if ! grep -q "$repo" "${MY_GIRRAFE_ROOT}/tools.go"; then
    _giraffe_die "add the tool to the tools.go file first: $repo"
  fi
  
  mkdir -p "$gobin"
  GOBIN="$gobin" go get "$repo"
}

