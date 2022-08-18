package main

import (
	"strings"
	"testing"

	"github.com/bearmini/sexp"
	"github.com/tj/assert"
)

func TestNextSexp(t *testing.T) {
	testData := []struct {
		Name     string
		Pattern  string
		Expected *sexp.Sexp
	}{
		{
			Name:     "pattern 1",
			Pattern:  ";; only a comment line",
			Expected: nil,
		},
		{
			Name:    "pattern 2 - only one assertion",
			Pattern: `(assert_return (invoke "add" (i32.const 1) (i32.const 1)) (i32.const 2))`,
			Expected: &sexp.Sexp{
				Children: []*sexp.Sexp{
					{Atom: &sexp.Token{Type: sexp.TokenTypeSymbol, Value: "assert_return"}},
					{Children: []*sexp.Sexp{
						{Atom: &sexp.Token{Type: sexp.TokenTypeSymbol, Value: "invoke"}},
						{Atom: &sexp.Token{Type: sexp.TokenTypeString, Value: `"add"`}},
						{Children: []*sexp.Sexp{
							{Atom: &sexp.Token{Type: sexp.TokenTypeSymbol, Value: "i32.const"}},
							{Atom: &sexp.Token{Type: sexp.TokenTypeNumber, Value: "1"}},
						}},
						{Children: []*sexp.Sexp{
							{Atom: &sexp.Token{Type: sexp.TokenTypeSymbol, Value: "i32.const"}},
							{Atom: &sexp.Token{Type: sexp.TokenTypeNumber, Value: "1"}},
						}},
					}},
					{Children: []*sexp.Sexp{
						{Atom: &sexp.Token{Type: sexp.TokenTypeSymbol, Value: "i32.const"}},
						{Atom: &sexp.Token{Type: sexp.TokenTypeNumber, Value: "2"}},
					}},
				},
			},
		},
	}

	for _, data := range testData {
		data := data // capture
		t.Run(data.Name, func(t *testing.T) {
			//t.Parallel()

			r := NewWastReader(strings.NewReader(data.Pattern))
			a, err := r.NextSexp()
			assert.NoError(t, err)
			assert.Equal(t, data.Expected, a)
		})
	}
}
