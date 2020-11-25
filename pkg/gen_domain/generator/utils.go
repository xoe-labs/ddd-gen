// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	"strings"
)

func shortForm(typ string) string {
	return strings.ToLower(string(typ[0]))
}
