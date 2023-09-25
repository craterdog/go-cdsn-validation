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
const intrinsic = `LOWERCASE|UPPERCASE|DIGIT|SEPARATOR|ESCAPE|EOL|EOF`

// This scanner is used for matching intrinsic tokens.
var intrinsicScanner = reg.MustCompile(`^(?:` + intrinsic + `)`)

// This function returns for the specified string an array of the matching
// subgroups for an intrinsic token. The first string in the array is the
// entire matched string.
func scanIntrinsic(v []byte) []string {
	return bytesToStrings(intrinsicScanner.FindSubmatch(v))
}
