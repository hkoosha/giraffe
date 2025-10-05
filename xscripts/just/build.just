build:
    go build .

test:
  GIRRAFE_TESTS=unit \
    ./xscripts/shell/exe_test_dev.sh

btest:
  GIRRAFE_TESTS=unit \
    ./xscripts/shell/exe_test_dev.sh \
    | grep -Ev '^\?' \
    | grep -Ev 'ok'

test-integ:
  GIRRAFE_TESTS=integ \
    ./xscripts/shell/exe_test_dev.sh

@summary:
  GIRRAFE_TESTS=unit \
    ./xscripts/shell/exe_test_dev.sh 2>&1 | grep -Ev '^(\?|ok|#) '

