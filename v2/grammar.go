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
	"$CHARACTER":   `LETTER | DIGIT`,
	"$COMMENT":     `"!>" EOL  {COMMENT | ~"<!"} EOL "<!"  ! Supports nested comments.`,
	"$EOL":         `"\n"  ! Standard POSIX definition.`,
	"$IDENTIFIER":  `LETTER {CHARACTER}`,
	"$INTRINSIC":   `"LETTER" | "DIGIT" | "EOF"`,
	"$LITERAL":     `"'" <~"'"> "'" | '"' <~'"'> '"'`,
	"$NOTE":        `"! " {~EOL}`,
	"$RANGE":       `CHARACTER ".." CHARACTER`,
	"$SYMBOL":      `"$" IDENTIFIER`,
	"$alternative": `[[NOTE] EOL] option`,
	"$factor": `
    RANGE        |  ! A range takes precedence over other factor types.
    INTRINSIC    |  ! These tokens have character set specific definitions.
    IDENTIFIER   |  ! These tokens map to symbols defined in other statements.
    LITERAL      |  ! Represents a literal string in quotes.
    "~" factor   |  ! Indicates the inverse of the factor.
    "(" rule ")" |  ! Indicates that the rule is evaluated first.
    "[" rule "]" |  ! Indicates zero or one repetitions of the rule.
    "{" rule "}" |  ! Indicates zero or more repetitions of the rule.
    "<" rule ">"    ! Indicates one or more repetitions of the rule.`,
	"$option":     `<factor>`,
	"$production": `SYMBOL ":" rule [NOTE]`,
	"$rule":       `option {"|" alternative}`,
	"$source":     `<statement> EOF  ! EOF is the end-of-file marker.`,
	"$statement":  `(COMMENT | production) <EOL>  ! EOL is the end-of-line character.`,
}

const header = `!>
    A formal definition of Crater Dog Syntax Notation™ (CDSN) using Crater Dog
    Syntax Notation™ itself. The token names are identified by all CAPITAL
    characters and the rule names are identified by lowerCamelCase characters.
    The starting rule is "$source".
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
