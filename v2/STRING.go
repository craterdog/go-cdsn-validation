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

type String string

const TokenString TokenType = "String"
const string_ = `["](?:` + escape + `|[^"` + eol + `])+["]`

// This method adds a string token with the current scanner information to the
// token channel. It returns true if a string token was found.
func (v *scanner) foundString() bool {
	var s = v.source[v.nextByte:]
	var matches = scanString(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenString)
		return true
	}
	return false
}

// This scanner is used for matching string tokens.
var stringScanner = reg.MustCompile(`^(?:` + string_ + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a string token. The first string in the array is the
// entire matched string.
func scanString(v []byte) []string {
	return bytesToStrings(stringScanner.FindSubmatch(v))
}
