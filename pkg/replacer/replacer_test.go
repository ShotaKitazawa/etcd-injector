package replacer

import (
	"os"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// TODO
		})
	}
}

func TestMain(m *testing.M) {
	// test
	status := m.Run()

	// exit
	os.Exit(status)
}
