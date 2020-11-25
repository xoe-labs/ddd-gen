// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package command

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/xoe-labs/ddd-gen/pkg/gen_app/generator"
)

type Config struct {
	Adapters generator.Adapters
	Errors   generator.Errors
	Objects  generator.Objects
}

func NewConfig(
	entity,
	authorizationErrorNew,
	targetIdentificationErrorNew,
	storageLoadingErrorNew,
	storageSavingErrorNew,
	domainErrorNew string,
) (*Config, error) {
	if !isValidQualId(entity) {
		return nil, fmt.Errorf("'%s' is not a valid full qualifier entity", entity)
	}
	if !isValidQualId(authorizationErrorNew) {
		return nil, fmt.Errorf("'%s' is not a valid full qualifier authorizationErrorNew", authorizationErrorNew)
	}
	if !isValidQualId(targetIdentificationErrorNew) {
		return nil, fmt.Errorf("'%s' is not a valid full qualifier targetIdentificationErrorNew", targetIdentificationErrorNew)
	}
	if !isValidQualId(storageLoadingErrorNew) {
		return nil, fmt.Errorf("'%s' is not a valid full qualifier storageLoadingErrorNew", storageLoadingErrorNew)
	}
	if !isValidQualId(storageSavingErrorNew) {
		return nil, fmt.Errorf("'%s' is not a valid full qualifier storageSavingErrorNew", storageSavingErrorNew)
	}
	if !isValidQualId(domainErrorNew) {
		return nil, fmt.Errorf("'%s' is not a valid full qualifier domainErrorNew", domainErrorNew)
	}
	return &Config{
		Adapters: generator.Adapters{},
		Objects: generator.Objects{
			Entity:     splitQual(entity),
		},
		Errors: generator.Errors{
			AuthorizationErrorNew:        splitQual(authorizationErrorNew),
			TargetIdentificationErrorNew: splitQual(targetIdentificationErrorNew),
			StorageLoadingErrorNew:       splitQual(storageLoadingErrorNew),
			StorageSavingErrorNew:        splitQual(storageSavingErrorNew),
			DomainErrorNew:               splitQual(domainErrorNew),
		},
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
