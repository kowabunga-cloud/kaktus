/*
 * Copyright (c) The Kowabunga Project
 * Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
 * SPDX-License-Identifier: Apache-2.0
 */

package kaktus

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestNewNfsExport_ValidID(t *testing.T) {
	export := NewNfsExport("42", "myexport", "cephfs", "/data", "rw", []int32{3, 4}, []string{"10.0.0.1"})
	if export.ID != 42 {
		t.Errorf("ID: expected 42, got %d", export.ID)
	}
}

func TestNewNfsExport_NamePrefixed(t *testing.T) {
	export := NewNfsExport("1", "myexport", "cephfs", "/data", "rw", []int32{3}, []string{})
	if export.Name != "/myexport" {
		t.Errorf("Name: expected '/myexport', got %q", export.Name)
	}
}

func TestNewNfsExport_InvalidID(t *testing.T) {
	export := NewNfsExport("not-a-number", "myexport", "cephfs", "/data", "rw", []int32{3}, []string{})
	if export.ID != 0 || export.Name != "" {
		t.Errorf("expected empty NfsExport for invalid ID, got %+v", export)
	}
}

func TestNewNfsExport_Fields(t *testing.T) {
	protocols := []int32{3, 4}
	clients := []string{"10.0.0.1", "10.0.0.2"}
	export := NewNfsExport("7", "share", "cephfs", "/mnt/share", "ro", protocols, clients)

	if export.FS != "cephfs" {
		t.Errorf("FS: expected 'cephfs', got %q", export.FS)
	}
	if export.Path != "/mnt/share" {
		t.Errorf("Path: expected '/mnt/share', got %q", export.Path)
	}
	if export.Access != "ro" {
		t.Errorf("Access: expected 'ro', got %q", export.Access)
	}
	if len(export.Protocols) != 2 || export.Protocols[0] != 3 || export.Protocols[1] != 4 {
		t.Errorf("Protocols: expected [3 4], got %v", export.Protocols)
	}
	if len(export.Clients) != 2 || export.Clients[0] != "10.0.0.1" || export.Clients[1] != "10.0.0.2" {
		t.Errorf("Clients: expected [10.0.0.1 10.0.0.2], got %v", export.Clients)
	}
}

func nfsTestServer(t *testing.T, handler http.HandlerFunc) (host string, port int) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("failed to parse test server URL: %v", err)
	}
	p, err := strconv.Atoi(u.Port())
	if err != nil {
		t.Fatalf("failed to parse test server port: %v", err)
	}
	return u.Hostname(), p
}

func TestNfsExport_CreateBackend(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody NfsExport

	host, port := nfsTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &gotBody)
		w.WriteHeader(http.StatusOK)
	})

	export := NewNfsExport("5", "myshare", "cephfs", "/data", "rw", []int32{3}, []string{"10.0.0.1"})
	if err := export.CreateBackend(host, port); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method: expected POST, got %s", gotMethod)
	}
	if gotPath != "/api/v1/export" {
		t.Errorf("path: expected '/api/v1/export', got %s", gotPath)
	}
	if gotBody.ID != 5 {
		t.Errorf("body ID: expected 5, got %d", gotBody.ID)
	}
	if gotBody.Name != "/myshare" {
		t.Errorf("body Name: expected '/myshare', got %q", gotBody.Name)
	}
}

func TestNfsExport_UpdateBackend(t *testing.T) {
	var gotMethod, gotPath string

	host, port := nfsTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	export := NewNfsExport("3", "share", "cephfs", "/data", "rw", []int32{3}, []string{})
	if err := export.UpdateBackend(host, port); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method: expected PUT, got %s", gotMethod)
	}
	if gotPath != "/api/v1/export/3" {
		t.Errorf("path: expected '/api/v1/export/3', got %s", gotPath)
	}
}

func TestNfsExport_DeleteBackend(t *testing.T) {
	var gotMethod, gotPath string

	host, port := nfsTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	export := NewNfsExport("9", "share", "cephfs", "/data", "rw", []int32{3}, []string{})
	if err := export.DeleteBackend(host, port); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method: expected DELETE, got %s", gotMethod)
	}
	if gotPath != "/api/v1/export/9" {
		t.Errorf("path: expected '/api/v1/export/9', got %s", gotPath)
	}
}

func TestNfsExport_CreateBackend_ServerError(t *testing.T) {
	host, port := nfsTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	export := NewNfsExport("1", "share", "cephfs", "/data", "rw", []int32{3}, []string{})
	// resty does not treat HTTP error status codes as errors; only network/transport errors fail.
	// This test confirms the method completes without a transport-level error.
	err := export.CreateBackend(host, port)
	if err != nil {
		t.Errorf("unexpected transport error: %v", err)
	}
}

func TestNewNfsConnectionSettings(t *testing.T) {
	ncs, err := NewNfsConnectionSettings()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ncs == nil {
		t.Error("expected non-nil NfsConnectionSettings")
	}
}
