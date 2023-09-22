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
	reg "regexp"
)

type Symbol string

const TokenSymbol TokenType = 10

// This private method appends a formatted symbol to the result.
func (v *formatter) formatSymbol(symbol Symbol) {
	v.appendString(string(symbol))
}

// This method adds a symbol token with the current scanner information to
// the token channel. It returns true if a symbol token was found.
func (v *scanner) foundSymbol() bool {
	var s = v.source[v.nextByte:]
	var matches = scanSymbol(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenSymbol)
		return true
	}
	return false
}

// This method attempts to parse a symbol. It returns the symbol and
// whether or not the symbol was successfully parsed.
func (v *parser) parseSymbol() (Symbol, *Token, bool) {
	var symbol Symbol
	var token = v.nextToken()
	if token.Type != TokenSymbol {
		v.backupOne()
		return symbol, token, false
	}
	symbol = Symbol(token.Value)
	return symbol, token, true
}

// This scanner is used for matching symbol tokens.
var symbolScanner = reg.MustCompile(`^(?:` + symbol + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a symbol token. The first string in the array is the
// entire matched string.
func scanSymbol(v []byte) []string {
	return bytesToStrings(symbolScanner.FindSubmatch(v))
}
