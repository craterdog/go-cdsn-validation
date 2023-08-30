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

// This map captures the syntax rules for Bali Wirth Syntax Notation.
// It is useful when creating scanner and parser error messages.
var grammar = map[string]string{
	"$INTRINSIC": `"LETTER" | "DIGIT | EOF"  ! Language specific definitions.`,
	"$COMMENT": `"!>" EOL  {COMMENT | ~"<!"} EOL "<!"  ! Allows recursion.`,
	"$EOL": `"\n"  ! Standard POSIX definition.`,
	"$IDENTIFIER": `LETTER {LETTER | DIGIT}`,
	"$NOTE": `"! " {~EOL}`,
	"$SYMBOL": `"$" IDENTIFIER`,
	"$factor": `
    INTRINSIC              |
    IDENTIFIER             |
    literal [".." literal] |
    "~" factor             |  ! Indicates NOT the factor.
    "(" rule ")"           |  ! Indicates that the rule is evaluated first.
    "[' rule "]"           |  ! Indicates zero or one repetitions of the rule.
    "{" rule "}"           |  ! Indicates zero or more repetitions of the rule.
    "<' rule ">"              ! Indicates one or more repetitions of the rule.`,
	"$grammar": `<production> EOF  ! EOF is the end-of-file marker.`,
	"$literal": `"'" <~"'"> "'" | '"' <~'"'> '"'`,
	"$production": `(COMMENT | SYMBOL ":" rule [NOTE]) <EOL>`,
	"$rule": `<factor> {"|" [[NOTE] EOL] <factor>}`,
}

const header = `!>
    A formal definition of Bali Wirth Syntax Notation™ (BWSN) using Bali Wirth
    Syntax Notation™ itself. The token names are identified by all CAPITAL
    characters and the rule names are identified by lowerCamelCase characters.
    The starting rule is "$grammar".
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
