// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package gen_app

import (
	"fmt"
	// "go/structes"
	"go/types"
	"log"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/xoe-labs/ddd-gen/pkg/gen_app/generator"
)

// StructTag Key
var (
	tagKey = "command"
)

// A simple regexp pattern to match tag values
var (
	topicTagPattern         = regexp.MustCompile(`topic,([^;]+)`)
	withoutPolicyTagPattern = regexp.MustCompile(`w/o policy`)
	adaptersTagPattern      = regexp.MustCompile(`adapters(?:,([^;]+:[^;]+))+`) // adapters,a1:github.com/foo/bar.Adapter1,a2:github.com/foo/bar.Adapter2
)

func generateDoc(docFile string) {
	if !fileExists(docFile) {
		df := generator.GenCommandDoc(docFile)
		df.Save(docFile)
	}
}

func analyzeStructAndGenerateCommandWrappers(genPath, sourceTypeName string, useFactStorage bool, struuct *types.Struct, adapters generator.Adapters, objects generator.Objects, errors generator.Errors) error {
	// determin the fully qualified package path
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName}, genPath)
	if err != nil {
		return err
	}
	pkgPath := pkgs[0].PkgPath
	log.Printf("Generating package: %s\n", pkgPath)
	log.Println("  using object interfaces ...")
	log.Printf("\t%s\n", objects.Entity)
	// log.Printf("\t%s\n", objects.TargetIdAssertable)
	log.Printf("\t%s\n", objects.Target)
	log.Printf("\t%s\n", objects.Entity)
	log.Printf("\t%s\n", objects.Actor)
	log.Printf("\t%s\n", objects.DomainCommandHandler)
	log.Println("  using adapter interfaces ...")
	// log.Printf("\t%s\n", adapters.StorageR)
	log.Printf("\t%s\n", adapters.StorageRW)
	log.Printf("\t%s\n", adapters.Policer)
	// log.Printf("\t%s\n", adapters.DomServiceAdapters)
	log.Println("  using error constructors ...")
	log.Printf("\t%s\n", errors.AuthorizationErrorNew)
	log.Printf("\t%s\n", errors.TargetIdentificationErrorNew)
	log.Printf("\t%s\n", errors.StorageLoadingErrorNew)
	log.Printf("\t%s\n", errors.StorageSavingErrorNew)
	log.Printf("\t%s\n", errors.DomainErrorNew)

	// 2. iterate over  fields
	for i := 0; i < struuct.NumFields(); i++ {
		field := struuct.Field(i)
		tag := reflect.StructTag(struuct.Tag(i))

		var (
			cmd                   string
			topic                 string
			withPolicyEnforcement bool
		)
		cmd = field.Name()
		withPolicyEnforcement = true

		// match and classify fields according to tags
		if tagKeyV, ok := tag.Lookup(tagKey); ok {
			if matches := topicTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				topic = matches[1]
			}
			if matches := withoutPolicyTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				withPolicyEnforcement = false
			}
			if matches := adaptersTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				for _, m := range matches[1:] {
					ss := strings.Split(m, ":")
					if !isValidQualId(ss[1]) {
						return fmt.Errorf("'adapters' tag value %s:%s does not contain a valid full qualifier", ss[0], ss[1])
					}
					adapters.DomServiceAdapters = append(adapters.DomServiceAdapters, generator.NamedQualId{Name: ss[0], QualId: splitQual(ss[1])})
				}
			}
		}
		if topic == "" {
			topic = getLastTitledWord(cmd)
		}

		topic = strings.Title(topic)
		log.Printf("topic %s -> %s: generating handler and stub\n", topic, cmd)

		fileBaseName := toSnakeCase(cmd)
		if getLastTitledWord(cmd) != topic {
			fileBaseName = fileBaseName + "_" + strings.ToLower(topic)

		}
		genFile := path.Join(genPath, fileBaseName+"_gen.go")

		// Remove existing generated file
		if fileExists(genFile) {
			if err := os.Remove(genFile); err != nil {
				return err
			}
		}
		gf := generator.GenCommandHandlerWrapper(cmd, topic, useFactStorage, withPolicyEnforcement, adapters, objects, errors)
		if err := gf.Save(genFile); err != nil {
			return err
		}

	}

	return nil
}

var (
	matchFirstLetterFollowedByCapWord = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllLowCapTransition          = regexp.MustCompile("([a-z0-9])([A-Z])")
	matchAllCapWord                   = regexp.MustCompile("[A-Z][a-z]+")
)

func getLastTitledWord(s string) string {
	b := []byte(s)
	slsl := matchAllCapWord.FindAllIndex(b, -1)
	lastsl := slsl[len(slsl)-1]
	return string(b[lastsl[0]:lastsl[1]])

}

func toSnakeCase(str string) string {
	snake := matchFirstLetterFollowedByCapWord.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllLowCapTransition.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
