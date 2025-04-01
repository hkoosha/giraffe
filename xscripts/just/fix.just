fix:
  go mod tidy
  @./xscripts/shell/exe_fix.sh

@lint:
  ./xscripts/shell/exe_lint.sh


six *args:
  go mod tidy
  ./xscripts/shell/exe_fix.sh {{args}}
