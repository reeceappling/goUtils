package s3

import (
	"github.com/reeceappling/goUtils/v2/utils/test"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	settings := test.RegisterTests(true, false)
	if !settings.RunTests {
		os.Exit(0)
	}

	code := m.Run()
	os.Exit(code)
}
