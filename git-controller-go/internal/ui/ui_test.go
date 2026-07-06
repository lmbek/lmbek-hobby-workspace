package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func capture(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()
	fn()
	_ = w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestUI_Primitives_PrintSomething(t *testing.T) {
	out := capture(t, func() { Header("hello") })
	if !strings.Contains(strings.ToUpper(out), "HELLO") {
		t.Fatalf("header did not print expected text: %q", out)
	}

	out = capture(t, func() { Info("info %d", 1) })
	if !strings.Contains(out, "info 1") {
		t.Fatalf("info missing: %q", out)
	}

	out = capture(t, func() { Success("ok") })
	if !strings.Contains(out, "ok") {
		t.Fatalf("success missing: %q", out)
	}

	out = capture(t, func() { Warn("careful") })
	if !strings.Contains(out, "careful") {
		t.Fatalf("warn missing: %q", out)
	}

	out = capture(t, func() { Error("bad") })
	if !strings.Contains(out, "bad") {
		t.Fatalf("error missing: %q", out)
	}

	out = capture(t, func() { Step(2, "do") })
	if !strings.Contains(out, "[2] do") {
		t.Fatalf("step missing: %q", out)
	}

	out = capture(t, func() { Note("note") })
	if !strings.Contains(out, "Note:") {
		t.Fatalf("note missing: %q", out)
	}
}
