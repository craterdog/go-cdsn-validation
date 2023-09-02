/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package cdsn

import (
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v2"
	sts "strings"
)

// This map captures the syntax rules for Crater Dog Syntax Notation.
// It is useful when creating scanner and parser error messages.
var grammar = map[string]string{
	"$NOTE":        `"! " {~EOL}`,
	"$COMMENT":     `"!>" EOL  {COMMENT | ~"<!"} EOL "<!"`,
	"$CHARACTER":   `"'" ~"'" "'"`,
	"$LITERAL":     `'"' <~'"'> '"'`,
	"$INTRINSIC":   `"LETTER" | "DIGIT" | "EOL" | "EOF"`,
	"$IDENTIFIER":  `LETTER {LETTER | DIGIT}`,
	"$SYMBOL":      `"$" IDENTIFIER`,
	"$source":      `<statement> EOF`,
	"$statement":   `(COMMENT | production) <EOL>`,
	"$production":  `SYMBOL ":" rule [NOTE]`,
	"$rule":        `option {"|" alternative}`,
	"$option":      `<factor>`,
	"$alternative": `[[NOTE] EOL] option`,
	"$range":       `CHARACTER ".." CHARACTER`,
	"$factor": `
    range        |
    "~" factor   |
    "(" rule ")" |
    "<" rule ">" |
    "[" rule "]" |
    "{" rule "}" |
    CHARACTER    |
    LITERAL      |
    INTRINSIC    |
    IDENTIFIER`,
}

const header = `!>
    A formal definition of Crater Dog Syntax Notation™ (CDSN) using Crater Dog
    Syntax Notation™ itself. Token names are identified by all CAPITAL
    letters and rule names are identified by lowerCamelCase letters.

    The INTRINSIC tokens are environment dependent and therefore left undefined.
    The tokens are scanned in the order listed so an INTRINSIC token takes
    precedence over an IDENTIFIER token.

    The rules are applied in the order listed as well, so within a factor a
    range takes precedence over an individual CHARACTER.  The starting rule is
    the "$source" rule.
<!

`

func FormatGrammar() string {
	var builder sts.Builder
	builder.WriteString(header)
	var unsorted = make([]string, len(grammar))
	var index = 0
	for key := range grammar {
		unsorted[index] = key
		index++
	}
	var keys = col.ListFromArray(unsorted)
	keys.SortValues()
	var iterator = col.Iterator[string](keys)
	for iterator.HasNext() {
		var key = iterator.GetNext()
		var value = grammar[key]
		builder.WriteString(fmt.Sprintf("%s: %s\n\n", key, value))
	}
	return builder.String()
}

// PRIVATE FUNCTIONS

func generateGrammar(expected string, symbols ...string) string {
	var message = "Was expecting '" + expected + "' from:\n"
	for _, symbol := range symbols {
		message += fmt.Sprintf("  \033[32m%v: \033[33m%v\033[0m\n\n", symbol, grammar[symbol])
	}
	return message
}
