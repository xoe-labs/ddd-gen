// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package gen_domain

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/xoe-labs/ddd-gen/pkg/gen_domain/generator"
)

type Config struct {
	Typ    string
	Entity generator.QualId
}

func NewConfig(
	typ,
	entity string,
) (*Config, error) {
	if !isValidQualId(entity) {
		return nil, fmt.Errorf("'%s' is not a valid full qualifier entity", entity)
	}
	if typ == "" {
		return nil, fmt.Errorf("'typ' is empty")
	}
	return &Config{
		Typ:    typ,
		Entity: splitQual(entity),
	}, nil
}

func isValidQualId(s string) bool {
	idx := strings.LastIndex(s, ".")
	if idx != -1 {
		id := s[strings.LastIndex(s, ".")+1:] // suggested identifier
		if strings.Index(id, "/") == -1 {     // no '/' in suggested identifier
			return unicode.IsUpper(rune(id[0])) // starts with upper case (is exported)
		}
	}
	return false
}

func splitQual(s string) generator.QualId {
	imp := s[:strings.LastIndex(s, ".")]
	id := s[strings.LastIndex(s, ".")+1:]
	return generator.QualId{
		Qual: imp,
		Id:   id,
	}
}
