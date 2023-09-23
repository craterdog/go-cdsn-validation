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

type Name string

const TokenName TokenType = "Name"
const name = `(?:` + lowercase + `|` + uppercase + `)(?:` + lowercase + `|` + uppercase + `|` + digit + `)*`

// This private method appends a formatted name to the result.
func (v *formatter) formatName(name Name) {
	v.appendString(string(name))
}

// This method adds a name token with the current scanner information
// to the token channel. It returns true if a name token was found.
func (v *scanner) foundName() bool {
	var s = v.source[v.nextByte:]
	var matches = scanName(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenName)
		return true
	}
	return false
}

// This method attempts to parse a name token. It returns the token and
// whether or not a name token was found.
func (v *parser) parseName() (Name, *Token, bool) {
	var name Name
	var token = v.nextToken()
	if token.Type != TokenName {
		v.backupOne()
		return name, token, false
	}
	if v.isToken && uni.IsLower(rune(token.Value[0])) {
		panic(fmt.Sprintf("A token definition contains a rulename: %v\n", token.Value))
	}
	name = Name(token.Value)
	var symbol = Symbol("$" + token.Value)
	var definition = v.symbols.GetValue(symbol)
	v.symbols.SetValue(symbol, definition)
	return name, token, true
}

// This scanner is used for matching name tokens.
var nameScanner = reg.MustCompile(`^(?:` + name + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a name token. The first string in the array is the
// entire matched string.
func scanName(v []byte) []string {
	return bytesToStrings(nameScanner.FindSubmatch(v))
}
