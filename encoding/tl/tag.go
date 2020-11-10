package tl

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	tagName              = "tl"
	optEncodedInBitflags = "encoded_in_bitflags"
	optFlagPrefix        = "flag:"
	optIgnore            = "-"
)

type tagInfo struct {
	index             int
	encodedInBitflags bool
	ignore            bool
}

func parseTag(s string) (info tagInfo, err error) {
	vals := strings.Split(s, ",")
	if len(vals) == 0 {
		err = fmt.Errorf("bad tag: %s", s)
		return
	}

	if haveInSlice(optIgnore, vals) {
		if len(vals) != 1 {
			err = fmt.Errorf("got '%s' with multiple options", optIgnore)
			return
		}

		info.ignore = true
		return
	}

	flag, haveFlag := haveStartsWith(optFlagPrefix, vals)
	if haveFlag {
		num := flag[len(optFlagPrefix):] // get index
		info.index, err = strconv.Atoi(num)
		if err != nil {
			err = fmt.Errorf("parse flag index '%s': %w", num, err)
			return
		}
	}

	if haveInSlice(optEncodedInBitflags, vals) {
		if !haveFlag {
			err = fmt.Errorf("have '%s' option without flag index", optEncodedInBitflags)
			return
		}

		info.encodedInBitflags = true
	}

	return
}

func haveInSlice(s string, slice []string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}

	return false
}

func haveStartsWith(s string, slice []string) (string, bool) {
	for _, item := range slice {
		if strings.HasPrefix(item, s) {
			return item, true
		}
	}

	return "", false
}
