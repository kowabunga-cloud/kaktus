/*
 * Copyright (c) The Kowabunga Project
 * Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
 * SPDX-License-Identifier: Apache-2.0
 */

package kaktus

import (
	"os"
	"testing"
)

func writeConfigFile(t *testing.T, content string) *os.File {
	t.Helper()
	f, err := os.CreateTemp("", "kaktus-cfg-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(f.Name())
	})
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatalf("failed to seek: %v", err)
	}
	return f
}

func TestKaktusConfigParser_LibvirtTCP(t *testing.T) {
	f := writeConfigFile(t, `
libvirt:
  protocol: tcp
  address: 192.168.1.10
  port: 16509
`)
	cfg, err := KaktusConfigParser(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Libvirt.Protocol != "tcp" {
		t.Errorf("protocol: expected 'tcp', got %q", cfg.Libvirt.Protocol)
	}
	if cfg.Libvirt.Address != "192.168.1.10" {
		t.Errorf("address: expected '192.168.1.10', got %q", cfg.Libvirt.Address)
	}
	if cfg.Libvirt.Port != 16509 {
		t.Errorf("port: expected 16509, got %d", cfg.Libvirt.Port)
	}
}

func TestKaktusConfigParser_LibvirtTLS(t *testing.T) {
	f := writeConfigFile(t, `
libvirt:
  protocol: tls
  address: 10.0.0.1
  port: 16514
  tls:
    key: /etc/ssl/client.key
    cert: /etc/ssl/client.crt
    ca: /etc/ssl/ca.crt
`)
	cfg, err := KaktusConfigParser(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Libvirt.TLS.PrivateKey != "/etc/ssl/client.key" {
		t.Errorf("TLS key: expected '/etc/ssl/client.key', got %q", cfg.Libvirt.TLS.PrivateKey)
	}
	if cfg.Libvirt.TLS.PublicCert != "/etc/ssl/client.crt" {
		t.Errorf("TLS cert: expected '/etc/ssl/client.crt', got %q", cfg.Libvirt.TLS.PublicCert)
	}
	if cfg.Libvirt.TLS.CA != "/etc/ssl/ca.crt" {
		t.Errorf("TLS CA: expected '/etc/ssl/ca.crt', got %q", cfg.Libvirt.TLS.CA)
	}
}

func TestKaktusConfigParser_Ceph(t *testing.T) {
	f := writeConfigFile(t, `
ceph:
  plugin: /usr/lib/kaktus/ceph.so
  monitor:
    name: ceph-mon-0
    address: 10.0.0.5
    port: 3300
`)
	cfg, err := KaktusConfigParser(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Ceph.PluginLib != "/usr/lib/kaktus/ceph.so" {
		t.Errorf("plugin: expected '/usr/lib/kaktus/ceph.so', got %q", cfg.Ceph.PluginLib)
	}
	if cfg.Ceph.Monitor.Name != "ceph-mon-0" {
		t.Errorf("monitor name: expected 'ceph-mon-0', got %q", cfg.Ceph.Monitor.Name)
	}
	if cfg.Ceph.Monitor.Address != "10.0.0.5" {
		t.Errorf("monitor address: expected '10.0.0.5', got %q", cfg.Ceph.Monitor.Address)
	}
	if cfg.Ceph.Monitor.Port != 3300 {
		t.Errorf("monitor port: expected 3300, got %d", cfg.Ceph.Monitor.Port)
	}
}

func TestKaktusConfigParser_InvalidYAML(t *testing.T) {
	f := writeConfigFile(t, "libvirt: [unclosed")
	_, err := KaktusConfigParser(f)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestKaktusConfigParser_EmptyFile(t *testing.T) {
	f := writeConfigFile(t, "")
	cfg, err := KaktusConfigParser(f)
	if err != nil {
		t.Fatalf("unexpected error for empty file: %v", err)
	}
	if cfg.Libvirt.Address != "" {
		t.Errorf("expected empty libvirt address for empty config, got %q", cfg.Libvirt.Address)
	}
	if cfg.Ceph.PluginLib != "" {
		t.Errorf("expected empty ceph plugin for empty config, got %q", cfg.Ceph.PluginLib)
	}
}
