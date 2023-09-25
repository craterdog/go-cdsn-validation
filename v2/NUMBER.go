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

type Number string

const TokenNumber TokenType = "Number"
const number = digit + `+`

// This private method appends a formatted number to the result.
func (v *formatter) formatNumber(number Number) {
	v.appendString(string(number))
}

// This method adds a number token with the current scanner information to the
// token channel. It returns true if a number token was found.
func (v *scanner) foundNumber() bool {
	var s = v.source[v.nextByte:]
	var matches = scanNumber(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenNumber)
		return true
	}
	return false
}

// This scanner is used for matching number tokens.
var numberScanner = reg.MustCompile(`^(?:` + number + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a number token. The first string in the array is the
// entire matched string.
func scanNumber(v []byte) []string {
	return bytesToStrings(numberScanner.FindSubmatch(v))
}
