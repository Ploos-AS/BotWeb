package registry

import (
	"os"
	"path/filepath"
	"testing"
)

func writeRegistry(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "bots.json")
	if err := os.WriteFile(path, []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadDefaultsPBMPProfileToBotM0(t *testing.T) {
	r, err := Load(writeRegistry(t, `{"bots":[{"id":"legacy","name":"Legacy","transport":"unix","endpoint":"/tmp/legacy.sock"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if got := r.Bots[0].Profile; got != "bot-m0" {
		t.Fatalf("profile=%q, want bot-m0", got)
	}
}

func TestLoadPreservesEndpointM0Profile(t *testing.T) {
	r, err := Load(writeRegistry(t, `{"bots":[{"id":"ambnc","name":"AmBNC","profile":"endpoint-m0","transport":"unix","endpoint":"/tmp/ambnc.sock"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if got := r.Bots[0].Profile; got != "endpoint-m0" {
		t.Fatalf("profile=%q, want endpoint-m0", got)
	}
}
