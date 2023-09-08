/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package #package#

import (
	reg "regexp"
)

// This token type is used to refer to all #token#s.
const Token#Token# TokenType = "#Token#"

// This regular expression matches all #token#s.
const #token#  = `['][^'][']`

// This function attempts to parse a #token#. It returns the #token# and
// whether or not a #token# was successfully parsed.
func parse#Token#(v *parser) (#Token#, *Token, bool) {
	var #token# #Token#
	var token = v.nextToken()
	if token.Type != Token#Token# {
		v.backupOne()
		return #token#, token, false
	}
	#token# = #Token#(token.Value)
	return #token#, token, true
}

// This function adds a #token# token with the current scanner information
// to the token channel. It returns true if a #token# token was found.
func found#Token#(v *scanner) bool {
	var s = v.source[v.nextByte:]
	var matches = scan#Token#(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token#Token#)
		return true
	}
	return false
}

// This scanner is used for matching #token# tokens.
var #token#Scanner = reg.MustCompile(`^(?:` + #token# + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a #token# token. The first string in the array is the
// entire matched string.
func scan#Token#(v []byte) []string {
	return bytesToStrings(#token#Scanner.FindSubmatch(v))
}
