// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	"strings"
	"regexp"
	"bytes"
)

func shortForm(typ string) string {
	return strings.ToLower(string(typ[0]))
}

func cmdShortForm(s string) string {
	re := regexp.MustCompile(`[A-Z]`)
	var b bytes.Buffer
	for _, el := range re.FindAllString(s, -1) {
		b.WriteString(strings.ToLower(el))
	}
	return b.String()
}
