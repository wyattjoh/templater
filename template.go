package main

import (
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// Code here is sourced from the dockerize repository: https://github.com/jwilder/dockerize

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func contains(item map[string]string, key string) bool {
	if _, ok := item[key]; ok {
		return true
	}
	return false
}

func defaultValue(args ...interface{}) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("default called with no values")
	}

	if len(args) > 0 {
		if args[0] != nil {
			return args[0].(string), nil
		}
	}

	if len(args) > 1 {
		if args[1] == nil {
			return "", fmt.Errorf("default called with nil default value")
		}

		if _, ok := args[1].(string); !ok {
			return "", fmt.Errorf("default is not a string value. hint: surround it w/ double quotes")
		}

		return args[1].(string), nil
	}

	return "", fmt.Errorf("default called with no default value")
}

func parseURL(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Fatalf("unable to parse url %s: %s", rawurl, err)
	}
	return u
}

func add(arg1, arg2 int) int {
	return arg1 + arg2
}

func isTrue(s string) bool {
	b, err := strconv.ParseBool(strings.ToLower(s))
	if err == nil {
		return b
	}
	return false
}

func generateFile(templatePath, destPath string, ctx *Context) bool {
	tmpl := template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"contains": contains,
		"exists":   exists,
		"split":    strings.Split,
		"replace":  strings.Replace,
		"default":  defaultValue,
		"parseUrl": parseURL,
		"atoi":     strconv.Atoi,
		"add":      add,
		"isTrue":   isTrue,
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
	})

	if len(delims) > 0 {
		tmpl = tmpl.Delims(delims[0], delims[1])
	}
	tmpl, err := tmpl.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("unable to parse template: %s", err)
	}

	dest := os.Stdout
	if destPath != "" {
		dest, err = os.Create(destPath)
		if err != nil {
			log.Fatalf("unable to create %s", err)
		}
		defer dest.Close()
	}

	err = tmpl.ExecuteTemplate(dest, filepath.Base(templatePath), ctx)
	if err != nil {
		log.Fatalf("template error: %s\n", err)
	}

	if fi, err := os.Stat(destPath); err == nil {
		if err := dest.Chmod(fi.Mode()); err != nil {
			log.Fatalf("unable to chmod temp file: %s\n", err)
		}
		if err := dest.Chown(int(fi.Sys().(*syscall.Stat_t).Uid), int(fi.Sys().(*syscall.Stat_t).Gid)); err != nil {
			log.Fatalf("unable to chown temp file: %s\n", err)
		}
	}

	return true
}
