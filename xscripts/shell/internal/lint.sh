function giraffe_lint() {
  (set -x;
  golangci-lint \
    -c "${MY_GIRRAFE_ROOT}/.golangci.yml" \
    run \
    "$@"
  )
}

function giraffe_fix() {
  giraffe_lint --fix "$@"
}
