// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package command

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

	"github.com/xoe-labs/ddd-gen/pkg/gen_app/directive"
)

// StructTag Key
var (
	tagKey = "command"
)

// A simple regexp pattern to match tag values
var (
	newTagPattern           = regexp.MustCompile(`new(,non-identifiable)?`)
	delTagPattern           = regexp.MustCompile(`del`)
	topicTagPattern         = regexp.MustCompile(`topic,([^;]+)`)
	withoutPolicyTagPattern = regexp.MustCompile(`w/o policy`)
	adaptersTagPattern      = regexp.MustCompile(`adapters(?:,([^;]+:[^;]+))+`) // adapters,a1:github.com/foo/bar.Adapter1,a2:github.com/foo/bar.Adapter2
)

func generateDoc(docFile string) {
	if !fileExists(docFile) {
		df := directive.GenDoc(docFile)
		df.Save(docFile)
	}
}

func generate(genPath, sourceTypeName string, struuct *types.Struct, conf directive.ParsedConfig) error {
	log.Printf("Generating code in: %s\n", genPath)
	log.Println("  using arguments ...")
	log.Printf("\taggregate:    %s %s\n", conf.AggEntityStruct.Qual, conf.AggEntityStruct.Id)
	log.Printf("\tidentifiable: %s %s\n", conf.IdentifiableInterface.Qual, conf.IdentifiableInterface.Id)
	log.Printf("\tpoliceable:   %s %s\n", conf.PoliceableInterface.Qual, conf.PoliceableInterface.Id)
	log.Printf("\tpolicer:      %s %s\n", conf.PolicerInterface.Qual, conf.PolicerInterface.Id)
	log.Printf("\trepository:   %s %s\n", conf.RepositoryInterface.Qual, conf.RepositoryInterface.Id)
	log.Println("  infered types ...")
	log.Printf("\tidentiferTyp:        %s %s\n", conf.IdentifierTyp.Qual, conf.IdentifierTyp.Id)
	log.Printf("\tuserTyp:             %s %s\n", conf.UserTyp.Qual, conf.UserTyp.Id)
	log.Printf("\televationTokenTyp:   %s %s\n", conf.ElevationTokenTyp.Qual, conf.ElevationTokenTyp.Id)
	log.Println("  error constructors ...")
	log.Printf("\tauthorizationErrorNew:     %s %s\n", conf.AuthorizationErrorNew.Qual, conf.AuthorizationErrorNew.Id)
	log.Printf("\tidentificationErrorNew:    %s %s\n", conf.IdentificationErrorNew.Qual, conf.IdentificationErrorNew.Id)
	log.Printf("\trepositoryErrorNew:        %s %s\n", conf.RepositoryErrorNew.Qual, conf.RepositoryErrorNew.Id)
	log.Printf("\tdomainErrorNew:            %s %s\n", conf.DomainErrorNew.Qual, conf.DomainErrorNew.Id)

	// 2. iterate over  fields
	for i := 0; i < struuct.NumFields(); i++ {
		field := struuct.Field(i)
		tag := reflect.StructTag(struuct.Tag(i))

		var (
			cmd                 string
			topic               string
			withPolicy          bool
			withCommandStub     bool
			newWithIdentifiable bool
			allAdapters         []directive.NamedQualId
			extraAdapters       []directive.NamedQualId
			genNew              bool
			genDel              bool
		)
		cmd = field.Name()
		withPolicy = true
		withCommandStub = true
		newWithIdentifiable = true

		// match and classify fields according to tags
		if tagKeyV, ok := tag.Lookup(tagKey); ok {
			if matches := topicTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				topic = matches[1]
			}
			if matches := withoutPolicyTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				withPolicy = false
			}
			if matches := newTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				genNew = true
				if len(matches) > 1 && matches[1] != "" {
					newWithIdentifiable = false
				}
			}
			if matches := delTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				genDel = true
			}
			if matches := adaptersTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				for _, m := range matches[1:] {
					ss := strings.Split(m, ":")
					if !isValidQualId(ss[1]) {
						return fmt.Errorf("'adapters' tag value %s:%s does not contain a valid full qualifier", ss[0], ss[1])
					}
					allAdapters = append(allAdapters, directive.NamedQualId{Name: ss[0], QualId: splitQual(ss[1])})
					extraAdapters = append(extraAdapters, directive.NamedQualId{Name: ss[0], QualId: splitQual(ss[1])})
				}
			}
		}
		if genNew && genDel {
			return fmt.Errorf("only one tag value of new, del can be set")
		}
		if topic == "" {
			topic = getLastTitledWord(cmd)
		}
		if withPolicy == true {
			allAdapters = append(allAdapters, directive.NamedQualId{Name: "pol", QualId: conf.PolicerInterface})
		}
		allAdapters = append(allAdapters, directive.NamedQualId{Name: "agg", QualId: conf.RepositoryInterface})

		topic = strings.Title(topic)
		log.Printf("topic %s -> %s: generating handler and stub\n", topic, cmd)

		fileBaseName := toSnakeCase(cmd)
		if getLastTitledWord(cmd) != topic {
			fileBaseName = fileBaseName + "_" + strings.ToLower(topic)

		}
		genFile := path.Join(genPath, fileBaseName+"_gen.go")
		stubFile := path.Join(genPath, fileBaseName+".go")

		// Remove existing generated file
		if fileExists(genFile) {
			if err := os.Remove(genFile); err != nil {
				return err
			}
		}

		genTyp := directive.UpdTyp
		if genNew {
			genTyp = directive.AddTyp
		} else if genDel {
			genTyp = directive.RemTyp
		}

		gf := directive.GenCommand(cmd, topic, withPolicy, newWithIdentifiable, allAdapters, extraAdapters, genTyp, conf)
		if err := gf.Save(genFile); err != nil {
			return err
		}

		if !fileExists(stubFile) {
			sf := directive.StubCommand(cmd, topic, withPolicy, withCommandStub, newWithIdentifiable, extraAdapters, genTyp, conf)
			if err := sf.Save(stubFile); err != nil {
				return err
			}
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
