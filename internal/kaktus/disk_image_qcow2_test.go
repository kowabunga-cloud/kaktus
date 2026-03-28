/*
 * Copyright (c) The Kowabunga Project
 * Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
 * SPDX-License-Identifier: Apache-2.0
 */

package kaktus

import (
	"testing"

	"github.com/lima-vm/go-qcow2reader/image/qcow2"
)

func TestQcowEncryptionMethod(t *testing.T) {
	tests := []struct {
		method   qcow2.CryptMethod
		expected string
	}{
		{qcow2.CryptMethodNone, "unencrypted"},
		{qcow2.CryptMethodAES, "AES-encrypted"},
		{qcow2.CryptMethodLUKS, "LUKS-encrypted"},
		{qcow2.CryptMethod(99), ""},
	}
	for _, tc := range tests {
		result := qcowEncryptionMethod(tc.method)
		if result != tc.expected {
			t.Errorf("qcowEncryptionMethod(%d): expected %q, got %q", tc.method, tc.expected, result)
		}
	}
}

func TestQcowCompressionType(t *testing.T) {
	tests := []struct {
		ct       qcow2.CompressionType
		expected string
	}{
		{qcow2.CompressionTypeZlib, "zlib-compressed"},
		{qcow2.CompressionTypeZstd, "zstd-compressed"},
		{qcow2.CompressionType(99), ""},
	}
	for _, tc := range tests {
		result := qcowCompressionType(tc.ct)
		if result != tc.expected {
			t.Errorf("qcowCompressionType(%d): expected %q, got %q", tc.ct, tc.expected, result)
		}
	}
}
