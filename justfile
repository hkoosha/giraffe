set shell := ['bash', '-c']
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
  just build test lint

@path:
    echo "$PATH"

