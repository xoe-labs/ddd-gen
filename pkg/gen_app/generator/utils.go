// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	"bytes"
	"regexp"
	"strings"
)

// Utils ...

func cmdShortForm(s string) string {
	re := regexp.MustCompile(`[A-Z]`)
	var b bytes.Buffer
	for _, el := range re.FindAllString(s, -1) {
		b.WriteString(strings.ToLower(el))
	}
	return b.String()
}

func splitQual(s string) (string, string) {
	imp := s[:strings.LastIndex(s, ".")]
	id := s[strings.LastIndex(s, ".")+1:]
	return imp, id
}
