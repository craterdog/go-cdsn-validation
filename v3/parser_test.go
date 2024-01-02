/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package cdsn_test

import (
	fmt "fmt"
	cds "github.com/craterdog/go-cdsn-validation/v2"
	ass "github.com/stretchr/testify/assert"
	osx "os"
	sts "strings"
	tes "testing"
)

const grammarsDirectory = "./grammars/"

func TestParsingRoundtrips(t *tes.T) {
	var files, err = osx.ReadDir(grammarsDirectory)
	if err != nil {
		panic("Could not find the " + grammarsDirectory + " directory.")
	}

	for _, file := range files {
		var parser = cds.ParserClass().Default()
		var validator = cds.ValidatorClass().Default()
		var formatter = cds.FormatterClass().Default()
		var filename = grammarsDirectory + file.Name()
		if sts.HasSuffix(filename, ".cdsn") {
			fmt.Println(filename)
			var bytes, _ = osx.ReadFile(filename)
			var expected = string(bytes)
			var document = parser.ParseDocument(expected)
			validator.ValidateDocument(document)
			var actual = formatter.FormatDocument(document)
			ass.Equal(t, expected, actual)
		}
	}
}

func TestRuleInTokenDefinition(t *tes.T) {
	var parser = cds.ParserClass().Default()
	var validator = cds.ValidatorClass().Default()
	var document = `$BAD: rule
$rule: "bad"
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"The definition for $BAD is invalid:\nA token definition cannot contain a rule name.\n",
				e,
			)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()

	validator.ValidateDocument(parser.ParseDocument(document))
}

func TestDoubleInversion(t *tes.T) {
	var parser = cds.ParserClass().Default()
	var validator = cds.ValidatorClass().Default()
	var document = `$BAD: ~~CONTROL
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"An unexpected token was received by the parser: Token [type: Delimiter, line: 1, position: 8]: \"~\"\n\x1b[36m0001: $BAD: ~~CONTROL\n \x1b[32m>>>─────────⌃\x1b[36m\n0002: \n\x1b[0m\nWas expecting 'assertion' from:\n  \x1b[32m$predicate: \x1b[33m\"~\"? assertion\x1b[0m\n\n  \x1b[32m$assertion: \x1b[33melement | glyph | precedence\x1b[0m\n\n",
				e,
			)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()

	validator.ValidateDocument(parser.ParseDocument(document))
}

func TestInvertedString(t *tes.T) {
	var parser = cds.ParserClass().Default()
	var validator = cds.ValidatorClass().Default()
	var document = `$BAD: ~"ow"
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"The definition for $BAD is invalid:\nA multi-character literal is not allowed in an inversion.\n",
				e,
			)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()

	validator.ValidateDocument(parser.ParseDocument(document))
}

func TestInvertedRule(t *tes.T) {
	var parser = cds.ParserClass().Default()
	var validator = cds.ValidatorClass().Default()
	var document = `$bad: ~rule
$rule: "rule"
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"The definition for $bad is invalid:\nAn inverted assertion cannot contain a rule name.\n",
				e,
			)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()

	validator.ValidateDocument(parser.ParseDocument(document))
}

func TestMissingRule(t *tes.T) {
	var parser = cds.ParserClass().Default()
	var validator = cds.ValidatorClass().Default()
	var document = `$bad: rule
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"The grammar is missing a definition for name: rule\n",
				e,
			)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()

	validator.ValidateDocument(parser.ParseDocument(document))
}

func TestDuplicateRule(t *tes.T) {
	var parser = cds.ParserClass().Default()
	var validator = cds.ValidatorClass().Default()
	var document = `$bad: "bad"
$bad: "worse"
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"An unexpected token was received by the parser: Token [type: Symbol, line: 2, position: 1]: \"$bad\"\n\x1b[36m0001: $bad: \"bad\"\n0002: $bad: \"worse\"\n \x1b[32m>>>──⌃\x1b[36m\n0003: \n\x1b[0m\nThis symbol has already been defined in this grammar.\n",
				e,
			)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()

	validator.ValidateDocument(parser.ParseDocument(document))
}

func TestNestedInversions(t *tes.T) {
	var parser = cds.ParserClass().Default()
	var validator = cds.ValidatorClass().Default()
	var document = `$BAD: ~(WORSE | ~BAD)
$WORSE: CONTROL
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"The definition for $BAD is invalid:\nInverted assertions cannot be nested.\n",
				e,
			)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()

	validator.ValidateDocument(parser.ParseDocument(document))
}
