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

type Intrinsic string

const TokenIntrinsic TokenType = "Intrinsic"
const intrinsic = `LOWERCASE|UPPERCASE|DIGIT|SEPARATOR|WHITESPACE|ESCAPE|EOL|EOF`

// This method adds an intrinsic token with the current scanner information
// to the token channel. It returns true if an intrinsic token was found.
func (v *scanner) foundIntrinsic() bool {
	var s = v.source[v.nextByte:]
	var matches = scanIntrinsic(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenIntrinsic)
		return true
	}
	return false
}

// This scanner is used for matching intrinsic tokens.
var intrinsicScanner = reg.MustCompile(`^(?:` + intrinsic + `)`)

// This function returns for the specified string an array of the matching
// subgroups for an intrinsic token. The first string in the array is the
// entire matched string.
func scanIntrinsic(v []byte) []string {
	return bytesToStrings(intrinsicScanner.FindSubmatch(v))
}
