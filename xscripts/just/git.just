trigger-ci:
    git commit --allow-empty -m "trigger CI"

[private]
get-file-branch branch file:
  ./xscripts/shell/git_dl_file.sh \
    giraffe \
    '{{branch}}' \
    '{{file}}' \
    '{{file}}'

[private]
get-file-curr-branch file:
  just get-file-branch \
    "$(git rev-parse --abbrev-ref HEAD)" \
    '{{file}}'

get-file branch-or-file file='':
  if [[ -z '{{ file }}' ]]; then just get-file-curr-branch '{{branch-or-file}}'            ; \
  else                           just get-file-branch      '{{branch-or-file}}' '{{file}}' ; fi
    
