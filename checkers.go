// Package fancycheck provides fancy gocheck-compatible checkers.
package fancycheck

import (
	"bytes"
	"fmt"

	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"

	"gopkg.in/check.v1"
)

var StringEquals check.Checker = &stringEqualsChecker{
	check.CheckerInfo{Name: "StringEquals", Params: []string{"obtained", "expected"}},
}

type stringEqualsChecker struct {
	check.CheckerInfo
}

func (c *stringEqualsChecker) Check(args []interface{}, names []string) (result bool, err string) {
	defer func() {
		if v := recover(); v != nil {
			result, err = false, fmt.Sprint(v)
		}
	}()
	obtained := fmt.Sprint(args[0])
	expected := fmt.Sprint(args[1])
	if obtained == expected {
		return true, ""
	}
	var buf bytes.Buffer
	d := diffmatchpatch.New()
	diffs := d.DiffMain(obtained, expected, true)
	deleteSprint := color.New(color.CrossedOut, color.BgRed).SprintFunc()
	insertSprint := color.New(color.Bold, color.BgGreen).SprintFunc()
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			buf.WriteString(deleteSprint(diff.Text))
		case diffmatchpatch.DiffEqual:
			buf.WriteString(diff.Text)
		case diffmatchpatch.DiffInsert:
			buf.WriteString(insertSprint(diff.Text))
		}
	}
	return false, buf.String()
}
