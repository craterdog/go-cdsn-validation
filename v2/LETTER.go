/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
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
	reg "regexp"
	uni "unicode"
)

type Letter string

const TokenLetter TokenType = "Letter"
const letter = lowercase + `|` + uppercase

// This private method appends a formatted letter to the result.
func (v *formatter) formatLetter(letter Letter) {
	v.appendString(string(letter))
}

// This method adds a letter token with the current scanner information
// to the token channel. It returns true if a letter token was found.
func (v *scanner) foundLetter() bool {
	var s = v.source[v.nextByte:]
	var matches = scanLetter(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenLetter)
		return true
	}
	return false
}

// This method attempts to parse a letter token. It returns the token and
// whether or not a letter token was found.
func (v *parser) parseLetter() (Letter, *Token, bool) {
	var letter Letter
	var token = v.nextToken()
	if token.Type != TokenLetter {
		v.backupOne()
		return letter, token, false
	}
	if v.isToken && uni.IsLower(rune(token.Value[0])) {
		panic(fmt.Sprintf("A token definition contains a rulename: %v\n", token.Value))
	}
	letter = Letter(token.Value)
	var symbol = Symbol("$" + token.Value)
	var definition = v.symbols.GetValue(symbol)
	v.symbols.SetValue(symbol, definition)
	return letter, token, true
}

// This scanner is used for matching letter tokens.
var letterScanner = reg.MustCompile(`^(?:` + letter + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a letter token. The first string in the array is the
// entire matched string.
func scanLetter(v []byte) []string {
	return bytesToStrings(letterScanner.FindSubmatch(v))
}
