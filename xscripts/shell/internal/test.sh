function giraffe_test() {
  if [[ $# -eq 1 ]] && [[ "$1" == '-v' ]]; then
    go test -v -short ./...
  elif [[ $# -gt 0 ]]; then
    _giraffe_die 'bad args for testing'
  else
    go test -short ./...
  fi
}

function giraffe_test_dev() {
  local tmp
  tmp="$(mktemp)"
  giraffe_test "$@" 2>&1 | tee "$tmp" || true

  echo
  echo
  echo

  grep -Ev '^ok|\?' < "$tmp" || true
}

