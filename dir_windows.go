// Copyright 2014 The LevelDB-Go and Pebble Authors. All rights reserved. Use
// of this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

//go:build windows

package vfs

import (
	"os"
	"syscall"

	"github.com/cockroachdb/errors"
)

type windowsDir struct {
	File
}

func (windowsDir) Sync() error {
	// Silently ignore Sync() on Windows. This is the same behavior as
	// RocksDB. See port/win/io_win.cc:WinDirectory::Fsync().
	return nil
}

type wrapOsFile struct {
	*os.File
}

func (w wrapOsFile) Preallocate(offset, length int64) error {
	return nil
}

func (w wrapOsFile) SyncTo(length int64) (fullSync bool, err error) {
	return true, w.File.Sync()
}

func (w wrapOsFile) SyncData() error {
	return w.File.Sync()
}

func (w wrapOsFile) Prefetch(offset int64, length int64) error {
	return nil
}

func (defaultFS) OpenDir(name string) (File, error) {
	f, err := os.OpenFile(name, syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return windowsDir{File: wrapOsFile{File: f}}, nil
}
