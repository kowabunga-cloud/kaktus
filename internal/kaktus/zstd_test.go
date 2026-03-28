/*
 * Copyright (c) The Kowabunga Project
 * Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
 * SPDX-License-Identifier: Apache-2.0
 */

package kaktus

import (
	"bytes"
	"io"
	"testing"

	"github.com/klauspost/compress/zstd"
)

func compressZstd(t *testing.T, data []byte) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	enc, err := zstd.NewWriter(&buf)
	if err != nil {
		t.Fatalf("failed to create zstd encoder: %v", err)
	}
	if _, err := enc.Write(data); err != nil {
		t.Fatalf("failed to write zstd data: %v", err)
	}
	if err := enc.Close(); err != nil {
		t.Fatalf("failed to close zstd encoder: %v", err)
	}
	return &buf
}

func TestNewZstdDecompressor_Create(t *testing.T) {
	buf := compressZstd(t, []byte("hello world"))
	dec, err := NewZstdDecompressor(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec == nil {
		t.Error("expected non-nil decompressor")
	}
	_ = dec.Close()
}

func TestNewZstdDecompressor_ReadsData(t *testing.T) {
	original := []byte("the quick brown fox jumps over the lazy dog")
	buf := compressZstd(t, original)

	dec, err := NewZstdDecompressor(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		_ = dec.Close()
	}()

	result, err := io.ReadAll(dec)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if !bytes.Equal(result, original) {
		t.Errorf("decompressed data mismatch: expected %q, got %q", original, result)
	}
}

func TestNewZstdDecompressor_EmptyInput(t *testing.T) {
	buf := compressZstd(t, []byte{})
	dec, err := NewZstdDecompressor(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		_ = dec.Close()
	}()

	result, err := io.ReadAll(dec)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty output, got %q", result)
	}
}

func TestNewZstdDecompressor_InvalidData(t *testing.T) {
	// zstd.NewReader does not validate data until first Read.
	buf := bytes.NewReader([]byte("this is not valid zstd compressed data"))
	dec, err := NewZstdDecompressor(buf)
	if err != nil {
		// Some implementations may reject on creation; that is also acceptable.
		return
	}
	defer func() {
		_ = dec.Close()
	}()
	_, err = io.ReadAll(dec)
	if err == nil {
		t.Error("expected error reading invalid zstd data, got nil")
	}
}

func TestZstdDecompressor_Close(t *testing.T) {
	buf := compressZstd(t, []byte("test"))
	dec, err := NewZstdDecompressor(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dec.Close(); err != nil {
		t.Errorf("Close returned unexpected error: %v", err)
	}
}

func TestNewZstdDecompressor_LargeData(t *testing.T) {
	original := bytes.Repeat([]byte("kowabunga-kaktus-"), 10_000)
	buf := compressZstd(t, original)

	dec, err := NewZstdDecompressor(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		_ = dec.Close()
	}()

	result, err := io.ReadAll(dec)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if !bytes.Equal(result, original) {
		t.Errorf("large data decompression mismatch (len got=%d, want=%d)", len(result), len(original))
	}
}
