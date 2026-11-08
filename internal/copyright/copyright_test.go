package copyright

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/simplycubed/copyright/internal/config"
)

func TestChecker_FindsKeyword(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "main.go")
	content := "package main\n\n// Copyright 2025 SimplyCubed, GK\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	checker := NewChecker(&CheckerOptions{
		LinesToScan: 20,
		Keywords:    []string{"copyright"},
	})
	result, err := checker.Check(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if !result.HasCopyright {
		t.Fatal("expected copyright to be found")
	}
	if result.LineFound != 3 {
		t.Fatalf("LineFound = %d, want 3", result.LineFound)
	}
}

func TestChecker_RequireAll(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "main.go")
	if err := os.WriteFile(path, []byte("// Copyright 2025 SimplyCubed, GK\npackage main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	checker := NewChecker(&CheckerOptions{
		LinesToScan: 20,
		Keywords:    []string{"Copyright", "SPDX-License-Identifier"},
		RequireAll:  true,
	})
	result, err := checker.Check(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if result.HasCopyright {
		t.Fatal("expected missing copyright when one required keyword is absent")
	}
}

func TestChecker_RespectsLineLimit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "main.go")
	if err := os.WriteFile(path, []byte("package main\n// Copyright 2025 SimplyCubed, GK\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	checker := NewChecker(&CheckerOptions{
		LinesToScan: 1,
		Keywords:    []string{"copyright"},
	})
	result, err := checker.Check(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if result.HasCopyright {
		t.Fatal("keyword past LinesToScan must not count")
	}
}

func TestHeaderGenerator_RendersBothSPDXNames(t *testing.T) {
	cfg := config.NewDefault()
	cfg.Copyright.Template = "Copyright {{.Year}} {{.Company}}\nSPDX-License-Identifier: {{.SPDX}}\n{{.SPDXIdentifier}}"
	cfg.Copyright.Year = "2025"
	cfg.Copyright.Company = "SimplyCubed, GK"
	cfg.Copyright.SPDXIdentifier = "Apache-2.0"

	got, err := NewHeaderGenerator(cfg).Generate()
	if err != nil {
		t.Fatal(err)
	}
	want := "Copyright 2025 SimplyCubed, GK\nSPDX-License-Identifier: Apache-2.0\nApache-2.0"
	if got != want {
		t.Fatalf("Generate() = %q, want %q", got, want)
	}
}

func TestHeaderGenerator_NilConfig(t *testing.T) {
	if _, err := NewHeaderGenerator(nil).Generate(); err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestHeaderGenerator_InvalidTemplate(t *testing.T) {
	cfg := config.NewDefault()
	cfg.Copyright.Template = "{{.Year"
	_, err := NewHeaderGenerator(cfg).Generate()
	if err == nil || !strings.Contains(err.Error(), "parse") {
		t.Fatalf("expected parse error, got %v", err)
	}
}
