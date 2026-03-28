/*
 * Copyright (c) The Kowabunga Project
 * Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
 * SPDX-License-Identifier: Apache-2.0
 */

package kaktus

import (
	"os"
	"testing"

	"github.com/kowabunga-cloud/common/klog"
)

func TestMain(m *testing.M) {
	klog.Init("kaktus-test", []klog.LoggerConfiguration{
		{Type: "console", Enabled: true, Level: "ERROR"},
	})
	os.Exit(m.Run())
}
