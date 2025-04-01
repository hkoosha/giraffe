build:
    go build .

test:
  GIRRAFE_TESTS=unit \
    ./xscripts/shell/exe_test_dev.sh

test-integ:
  GIRRAFE_TESTS=integ \
    ./xscripts/shell/exe_test_dev.sh

