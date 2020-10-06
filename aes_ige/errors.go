package ige

import "github.com/pkg/errors"

var (
	ErrDataTooSmall     = errors.New("AES256IGE: data too small")
	ErrDataNotDivisible = errors.New("AES256IGE: data not divisible by block size")
)
