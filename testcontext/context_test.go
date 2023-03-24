// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package testcontext_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"common/testcontext"
)

func TestBasic(t *testing.T) {
	ctx := testcontext.New(t)

	ctx.Go(func() error {
		time.Sleep(time.Millisecond)
		return nil
	})

	t.Log(ctx.Dir("a", "b", "c"))
	t.Log(ctx.File("a", "w", "c.txt"))
}

func TestSubDir(t *testing.T) {
	ctx := testcontext.New(t)
	path := ctx.File("a", "w/c.txt")
	err := os.WriteFile(path, []byte{1}, 0644)
	assert.NoError(t, err)
}

func TestMessage(t *testing.T) {
	var subtest test

	ctx := testcontext.NewWithTimeout(&subtest, 50*time.Millisecond)
	ctx.Go(func() error {
		time.Sleep(time.Second)
		return nil
	})
	ctx.Cleanup()

	t.Log(subtest.errors[0])

	assert.Contains(t, subtest.errors[0], "Test exceeded timeout")
	assert.Contains(t, subtest.errors[0], "some goroutines are still running")

	assert.Contains(t, subtest.errors[1], "TestMessage")
}

type test struct {
	logs   []string
	errors []string
	fatals []string
}

func (t *test) Name() string      { return "Example" }
func (t *test) Helper()           {}
func (t *test) Cleanup(fn func()) {}

func (t *test) Log(args ...interface{})   { t.logs = append(t.logs, fmt.Sprint(args...)) }
func (t *test) Error(args ...interface{}) { t.errors = append(t.errors, fmt.Sprint(args...)) }
func (t *test) Fatal(args ...interface{}) { t.fatals = append(t.fatals, fmt.Sprint(args...)) }
