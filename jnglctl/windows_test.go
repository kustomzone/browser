// Copyright 2016 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestFindUpspinBinaries(t *testing.T) {
	testFiles := []struct {
		n string
		b bool
	}{
		{"some/path/upspin-foo.exe", true},
		{"some/path/upspin-bar.txt", false},
		{"thirst/path/upspin-baz.bat", true},
		{"fourth/path/upspin-qux.com", true},
		{"yet/another/upspin-foo.cmd", true},
	}
	tmpDir, err := ioutil.TempDir("", "upspin-binary-test-")
	if err != nil {
		t.Fatalf("could not create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	paths := map[string]bool{}
	for _, tf := range testFiles {
		d := filepath.Dir(filepath.Join(tmpDir, tf.n))
		paths[d] = true
		if err = os.MkdirAll(d, 0700); err != nil {
			t.Fatalf("could not create %s: %v", d, err)
			continue
		}
		f, err := os.Create(filepath.Join(tmpDir, tf.n))
		if err != nil {
			t.Fatalf("could not create temporary file %s: %v", tf.n, err)
			continue
		}
		f.Close()
	}

	defer os.Setenv("PATH", os.Getenv("PATH"))
	var newPath string
	for k, _ := range paths {
		newPath += k + string(filepath.ListSeparator)
	}
	err = os.Setenv("PATH", newPath)
	if err != nil {
		t.Fatalf("could not set PATH: %v", err)
	}

	binaries := findUpspinBinaries()
	sort.Strings(binaries)
	wanted := "baz;foo;foo;qux"
	got := strings.Join(binaries, ";")
	if wanted != got {
		t.Fatalf("expected %q, got %q", wanted, got)
	}
}
