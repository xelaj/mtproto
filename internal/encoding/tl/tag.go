// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
	"github.com/pkg/errors"
)

const tagName = "tl"

type fieldTag struct {
	index            int  // flags:<N>
	encodedInBitflag bool // encoded_in_bitflags
	ignore           bool // -
	optional         bool // omitempty
}

func parseTag(s reflect.StructTag) (*fieldTag, error) {
	tags, err := structtag.Parse(string(s))
	if err != nil {
		return nil, errors.Wrap(err, "parsing field tags")
	}

	tag, err := tags.Get(tagName)
	if err != nil {
		// ну не нашли и не нашли, че бубнить то
		return nil, nil
	}

	info := &fieldTag{}

	if tag.Name == "-" {
		info.ignore = true
		return info, nil
	}

	var flagIndexSet bool
	if strings.HasPrefix(tag.Name, "flag:") {
		num := strings.TrimPrefix(tag.Name, "flag:")
		info.index, err = strconv.Atoi(num)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing index number '%s'", num)
		}

		// поля внутри битфлагов всегда optional
		info.optional = true

		flagIndexSet = true
	}

	if haveInSlice("encoded_in_bitflags", tag.Options) {
		if !flagIndexSet {
			return nil, errors.New("have 'encoded_in_bitflag' option without flag index")
		}

		info.encodedInBitflag = true
	}

	if haveInSlice("omitempty", tag.Options) {
		info.optional = true
	}

	return info, nil
}

//! slicetricks
func haveInSlice(s string, slice []string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}

	return false
}
