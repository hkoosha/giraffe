set shell := ['bash', '-c', '-e', '-u', '-o', 'pipefail']
set dotenv-filename := './local/env'
set allow-duplicate-variables

tools_dir := justfile_directory() / 'build/tools' / shell('uname -s')
export PATH := tools_dir + ':' + env_var('PATH')
export GOPRIVATE := 'github.com/hkoosha'

import './xscripts/just/build.just'
import './xscripts/just/containers.just'
import './xscripts/just/fix.just'
import './xscripts/just/git.just'
import './xscripts/just/req.just'

import? './local/justfile'

import? 'local/justfile'

default:
  just build test fix

@path:
    echo "$PATH"


fixall:
   for i in $(cat go.work | grep -E $'^\t+(core|contrib)' | tr -d $'\t'); do echo "$i"; \
       ( cd "./$i" && golangci-lint -c ./.golangci.yml run --fix; ); done
   just fix

