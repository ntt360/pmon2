package utils

import (
	"regexp"
	"strings"
)

type Args struct {
	params map[string][]string
}

func ParseArgs(args []string) *Args {
	a := Args{params: map[string][]string{}}
	var currentKey string

	for _, arg := range args {
		reg := regexp.MustCompile(`^--(\w*)$`)
		if !reg.MatchString(arg) {
			if len(currentKey) <= 0 {
				continue
			}
			vals, ok := a.params[currentKey]
			if ok {
				vals = append(vals, arg)
			} else {
				vals = []string{arg}
			}
			a.params[currentKey] = vals
		} else {
			// 截取字符非 -- 符号作为key
			key := arg[2:]
			if len(key) <= 0 {
				key = "def_params"
			}
			currentKey = key
		}
	}

	return &a
}

func (a *Args) Get(key string) string {
	return strings.Join(a.params[key], " ")
}
