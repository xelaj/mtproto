package tl

import (
	"fmt"
	"strconv"
	"strings"
)

type tagInfo struct {
	index            int
	encodedInBitflag bool
	ignore           bool
}

func parseTag(s string) (info tagInfo, err error) {
	vals := strings.Split(s, ",")
	if len(vals) == 0 {
		err = fmt.Errorf("bad tag: %s", s)
		return
	}

	if haveInSlice("-", vals) {
		if len(vals) != 1 {
			err = fmt.Errorf("got '-' with multiple options")
			return
		}

		info.ignore = true
		return
	}

	flag, haveFlag := haveStartsWith("flag:", vals)
	if haveFlag {
		num := flag[len("flag:"):] // get index
		info.index, err = strconv.Atoi(num)
		if err != nil {
			err = fmt.Errorf("parse flag index '%s': %w", num, err)
			return
		}
	}

	if haveInSlice("encoded_in_bitflags", vals) {
		if !haveFlag {
			err = fmt.Errorf("have 'encoded_in_bitflag' option without flag index")
			return
		}

		info.encodedInBitflag = true
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
