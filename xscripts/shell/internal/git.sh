function giraffe_git_dl() {
  [[ $# -eq 4 ]] || _giraffe_die 'bad args'

  local -r my_repo="$1"
  local -r my_ref="$2"
  local -r my_file="$3"
  local -r my_write_to="$4"

  local -r my_url="https://api.github.com/repos/hkoosha/$my_repo/contents/$my_file?ref=${my_ref}"
  echo "downloading: $my_url"
  gh api "$my_url" -H "Accept: application/vnd.github.raw" > "$my_write_to"
}

