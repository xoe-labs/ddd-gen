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
	// "golang.org/x/tools/go/packages"
)

// StructTag Key
var (
	tagKey = "command"
)

// A simple regexp pattern to match tag values
var (
	newTagPattern           = regexp.MustCompile(`new`)
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

func generate(genPath, sourceTypeName string, struuct *types.Struct, conf Config) error {
	log.Printf("Generating code in: %s\n", genPath)
	log.Println("  using arguments ...")
	log.Printf("\taggregate:    %s\n", conf.aggEntityStruct)
	log.Printf("\tidentifiable: %s\n", conf.identifiableInterface)
	log.Printf("\tpoliceable:   %s\n", conf.policeableInterface)
	log.Printf("\tpolicer:      %s\n", conf.policerInterface)
	log.Printf("\trepository:   %s\n", conf.repositoryInterface)
	var (
		aggEntity    string = conf.aggEntityStruct
		identifiable string = conf.identifiableInterface
		policeable   string = conf.policeableInterface
	)

	// 2. iterate over  fields
	for i := 0; i < struuct.NumFields(); i++ {
		field := struuct.Field(i)
		tag := reflect.StructTag(struuct.Tag(i))

		var (
			cmd             string
			topic           string
			withPolicy      bool
			withCommandStub bool
			adapters        []struct{ Id, Qual string }
			genNew          bool
			genDel          bool
		)
		cmd = field.Name()
		withPolicy = true
		withCommandStub = true
		adapters = []struct{ Id, Qual string }{}

		// 2.2 match and classify fields according to tags
		if tagKeyV, ok := tag.Lookup(tagKey); ok {
			if matches := topicTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				topic = matches[1]
			}
			if matches := withoutPolicyTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				withPolicy = false
			}
			if matches := newTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				genNew = true
			}
			if matches := delTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				genDel = true
			}
			if matches := adaptersTagPattern.FindStringSubmatch(tagKeyV); matches != nil {
				for _, m := range matches[1:] {
					ss := strings.Split(m, ":")
					adapters = append(adapters, struct{ Id, Qual string }{Id: ss[0], Qual: ss[1]})
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
			adapters = append(adapters, struct{ Id, Qual string }{Id: "pol", Qual: conf.policerInterface})
		}
		adapters = append(adapters, struct{ Id, Qual string }{Id: "agg", Qual: conf.repositoryInterface})

		topic = strings.Title(topic)
		log.Printf("topic %s -> %s: generating handler and stub\n", topic, cmd)

		genFile := path.Join(genPath, cmd+"_gen.go")
		stubFile := path.Join(genPath, cmd+".go")

		// Remove existing generated file
		if fileExists(genFile) {
			if err := os.Remove(genFile); err != nil {
				return err
			}
		}
		gf := directive.GenCommand(cmd, topic, withPolicy, adapters, aggEntity)
		if err := gf.Save(genFile); err != nil {
			return err
		}

		if !fileExists(stubFile) {
			sf := directive.StubCommand(cmd, topic, withPolicy, withCommandStub, aggEntity, identifiable, policeable)
			if err := sf.Save(stubFile); err != nil {
				return err
			}
		}
	}

	return nil
}

func getLastTitledWord(s string) string {
	b := []byte(s)
	re := regexp.MustCompile(`[A-Z][a-z]+`)
	slsl := re.FindAllIndex(b, -1)
	lastsl := slsl[len(slsl)-1]
	return string(b[lastsl[0]:lastsl[1]])

}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
