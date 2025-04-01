function giraffe_build() {
  if [[ $# -eq 2 ]]; then
    local -r my_bin="${1}/${2}"
  elif [[ $# -eq 1 ]]; then
    local -r my_bin="$MY_GIRRAFE_ROOT/ybin/${1}"
  elif [[ $# -eq 0 ]]; then
    local -r my_bin="$MY_GIRRAFE_ROOT/ybin/giraffe"
  else
    _giraffe_die "too many args: $*"
  fi

  # local f=(
  #   "-X flag0=value1"
  #   "-X flag1=value1"
  # )
  # local -r flags="${f[*]}"

  # echo "go build -ldflags $flags -o $my_bin ."
  # go build -ldflags "$flags" -o "$my_bin" .
  #

  echo "go build -o $my_bin ."
  go build -o "$my_bin" ./...
}
