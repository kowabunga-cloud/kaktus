/*
 * Copyright (c) The Kowabunga Project
 * Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
 * SPDX-License-Identifier: Apache-2.0
 */

package kaktus

import (
	"fmt"
	"math"
)

func toUint8(i int) (uint8, error) {
	if i < 0 || i > math.MaxUint8 {
		return 0, fmt.Errorf("value %d out of uint8 range", i)
	}
	return uint8(i), nil
}

func toUint16(i int) (uint16, error) {
	if i < 0 || i > math.MaxUint16 {
		return 0, fmt.Errorf("value %d out of uint16 range", i)
	}
	return uint16(i), nil
}

func toUint64(i int64) (uint64, error) {
	if i < 0 {
		return 0, fmt.Errorf("value %d is negative, cannot convert to uint64", i)
	}
	return uint64(i), nil
}
