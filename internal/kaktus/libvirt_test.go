/*
 * Copyright (c) The Kowabunga Project
 * Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
 * SPDX-License-Identifier: Apache-2.0
 */

package kaktus

import (
	"encoding/xml"
	"testing"

	virtxml "libvirt.org/go/libvirtxml"
)

// getGuestMachine tests

func TestGetGuestMachine_FoundWithCanonical(t *testing.T) {
	machines := []virtxml.CapsGuestMachine{
		{Name: "pc-i440fx-2.9", Canonical: "pc"},
		{Name: "q35"},
	}
	result := getGuestMachine(machines, "pc-i440fx-2.9")
	if result != "pc" {
		t.Errorf("expected canonical 'pc', got %q", result)
	}
}

func TestGetGuestMachine_FoundNoCanonical(t *testing.T) {
	machines := []virtxml.CapsGuestMachine{
		{Name: "q35"},
	}
	result := getGuestMachine(machines, "q35")
	if result != "q35" {
		t.Errorf("expected name 'q35', got %q", result)
	}
}

func TestGetGuestMachine_NotFound(t *testing.T) {
	machines := []virtxml.CapsGuestMachine{
		{Name: "q35"},
	}
	result := getGuestMachine(machines, "pc")
	if result != "" {
		t.Errorf("expected empty string for missing machine, got %q", result)
	}
}

func TestGetGuestMachine_EmptyList(t *testing.T) {
	result := getGuestMachine([]virtxml.CapsGuestMachine{}, "pc")
	if result != "" {
		t.Errorf("expected empty string for empty list, got %q", result)
	}
}

func TestGetGuestMachine_CanonicalTakesPrecedence(t *testing.T) {
	// When canonical is set it should be returned, not the original name.
	machines := []virtxml.CapsGuestMachine{
		{Name: "pc-i440fx-6.2", Canonical: "pc-i440fx"},
	}
	result := getGuestMachine(machines, "pc-i440fx-6.2")
	if result != "pc-i440fx" {
		t.Errorf("expected canonical 'pc-i440fx', got %q", result)
	}
}

// getGuestMachineName tests

func TestGetGuestMachineName_ReturnsCanonical(t *testing.T) {
	guest := &virtxml.CapsGuest{
		Arch: virtxml.CapsGuestArch{
			Machines: []virtxml.CapsGuestMachine{
				{Name: "pc-i440fx-2.9", Canonical: "pc"},
			},
		},
	}
	name, err := getGuestMachineName(guest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "pc" {
		t.Errorf("expected canonical 'pc', got %q", name)
	}
}

func TestGetGuestMachineName_ReturnsNameWhenNoCanonical(t *testing.T) {
	guest := &virtxml.CapsGuest{
		Arch: virtxml.CapsGuestArch{
			Machines: []virtxml.CapsGuestMachine{
				{Name: "q35"},
			},
		},
	}
	name, err := getGuestMachineName(guest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "q35" {
		t.Errorf("expected 'q35', got %q", name)
	}
}

func TestGetGuestMachineName_MultipleMachines(t *testing.T) {
	// The first machine's name is used as the lookup target.
	guest := &virtxml.CapsGuest{
		Arch: virtxml.CapsGuestArch{
			Machines: []virtxml.CapsGuestMachine{
				{Name: "pc-i440fx-7.0", Canonical: "pc"},
				{Name: "q35"},
			},
		},
	}
	name, err := getGuestMachineName(guest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "pc" {
		t.Errorf("expected canonical of first machine 'pc', got %q", name)
	}
}

// xmlUnmarshal tests

func TestXmlUnmarshal_Valid(t *testing.T) {
	type Simple struct {
		XMLName xml.Name `xml:"root"`
		Value   string   `xml:"value"`
	}
	var result Simple
	err := xmlUnmarshal(`<root><value>hello</value></root>`, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Value != "hello" {
		t.Errorf("expected 'hello', got %q", result.Value)
	}
}

func TestXmlUnmarshal_MultipleFields(t *testing.T) {
	type Config struct {
		XMLName xml.Name `xml:"config"`
		Name    string   `xml:"name"`
		Port    int      `xml:"port"`
	}
	var result Config
	err := xmlUnmarshal(`<config><name>kaktus</name><port>8080</port></config>`, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "kaktus" {
		t.Errorf("Name: expected 'kaktus', got %q", result.Name)
	}
	if result.Port != 8080 {
		t.Errorf("Port: expected 8080, got %d", result.Port)
	}
}

func TestXmlUnmarshal_InvalidXML(t *testing.T) {
	type Simple struct {
		XMLName xml.Name `xml:"root"`
	}
	var result Simple
	// Mismatched tags are invalid XML and must produce an error.
	err := xmlUnmarshal(`<root><child></other></root>`, &result)
	if err == nil {
		t.Error("expected error for mismatched XML tags, got nil")
	}
}
