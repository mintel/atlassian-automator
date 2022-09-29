//go:build tools
// +build tools

package main

import (
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/google/go-cmp/cmp"
	_ "github.com/stretchr/testify/assert"
)
