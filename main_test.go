package main

import "testing"

func TestVersionLine(t *testing.T) {
	oldVersion, oldCommit, oldDate := Version, Commit, Date
	t.Cleanup(func() {
		Version, Commit, Date = oldVersion, oldCommit, oldDate
	})

	Version = "dev"
	if got := versionLine(); got != "copyright dev (commit none, date unknown)" {
		t.Fatalf("versionLine() = %q", got)
	}

	Version = "1.2.3"
	Commit = "abc123"
	Date = "2026-11-08"
	if got := versionLine(); got != "copyright 1.2.3 (commit abc123, date 2026-11-08)" {
		t.Fatalf("versionLine() = %q", got)
	}
}
