package gtesting

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/hkoosha/giraffe/core/t11y"
)

func NoError(
	t *testing.T,
	err any,
) {
	t.Helper()

	if err == nil {
		return
	}

	content := []struct{ label, content string }{
		{"G11y Trace", t11y.FmtStacktraceOf(err)},
		{"Error", fmt.Sprintf("Received unexpected error:\n%+v", err)},
		{"Test", t.Name()},
	}

	longest := len("Error Trace")

	var sb0 strings.Builder
	for _, v := range content {
		sb1 := strings.Builder{}
		for i, scanner := 0, bufio.NewScanner(strings.NewReader(v.content)); scanner.Scan(); i++ {
			if i != 0 {
				sb1.WriteString("\n\t" + strings.Repeat(" ", longest+1) + "\t")
			}
			sb1.WriteString(scanner.Text())
		}

		sb0.WriteByte('\t')
		sb0.WriteString(v.label)
		sb0.WriteByte(':')
		sb0.WriteString(strings.Repeat(" ", longest-len(v.label)))
		sb0.WriteByte('\t')
		sb0.WriteString(sb1.String())
		sb0.WriteByte('\n')
	}

	t.Errorf("\n%s", ""+sb0.String())
	t.FailNow()
}
