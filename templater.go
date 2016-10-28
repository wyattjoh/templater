package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

type sliceVar []string

func (s *sliceVar) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (s *sliceVar) String() string {
	return strings.Join(*s, ",")
}

var (
	delims        []string
	templatesFlag sliceVar
	contextFlag   string
)

func main() {
	flag.Var(&templatesFlag, "template", "Template (/template:/dest). Can be passed multiple times")
	flag.StringVar(&contextFlag, "context", "", "Filename for json context file")

	flag.Parse()

	if flag.NArg() == 0 && flag.NFlag() == 0 {
		os.Exit(1)
	}

	ctx := &Context{}

	var err error

	if contextFlag != "" {

		// load in the context
		ctx, err = loadContext(contextFlag)
		if err != nil {
			log.Fatalf("can't load context file: %s", err.Error())
		}
	}

	for _, t := range templatesFlag {
		template, dest := t, ""
		if strings.Contains(t, ":") {
			parts := strings.Split(t, ":")
			if len(parts) != 2 {
				log.Fatalf("bad template argument: %s. expected \"/template:/dest\"", t)
			}
			template, dest = parts[0], parts[1]
		}
		generateFile(template, dest, ctx)
	}
}
