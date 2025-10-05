package gtesting

import (
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	TestWriteDebug = "GIRAFFE_TEST_WRITE_DEBUG"

	TestWriteDir = "/tmp/giraffe/test_debug"
)

func IsExtraTestDebug() bool {
	return true
	// en := strings.ToLower(strings.TrimSpace(os.Getenv(TestWriteDebug)))
	// return en == "1" || en == "true"
}

func EnsureDir() error {
	return os.MkdirAll(TestWriteDir, 0o755)
}

func write(
	t *testing.T,
	out string,
	content string,
) error {
	t.Helper()

	if !IsExtraTestDebug() {
		return nil
	}

	t.Logf(out, content)

	_, err := os.Stat(out)
	if !os.IsNotExist(err) || err == nil {
		// return errors.New("output exists: " + out)
		t.Log("output exists", out)
	}

	file, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func Write(
	t *testing.T,
	name string,
	content string,
) {
	t.Helper()

	tName := regexp.
		MustCompile("[^a-zA-Z0-9_]").
		ReplaceAllString(t.Name(), "_")
	for strings.Contains(tName, "__") {
		tName = strings.ReplaceAll(tName, "__", "_")
	}
	name = tName + "__" + name

	NoError(t, EnsureDir())

	out := path.Join(TestWriteDir, name)

	if strings.Contains(out, "..") ||
		!strings.HasPrefix(out, TestWriteDir) ||
		!strings.HasPrefix(out, "/") {
		require.Fail(t, "invalid output name: "+name)

		// Do NOT remove this return.
		return
	}

	NoError(t, write(t, out, content))
}
