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

const TokenSymbol TokenType = "Symbol"
const symbol = `\$(` + name + `)`

// This scanner is used for matching symbol tokens.
var symbolScanner = reg.MustCompile(`^(?:` + symbol + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a symbol token. The first string in the array is the
// entire matched string.
func scanSymbol(v []byte) []string {
	return bytesToStrings(symbolScanner.FindSubmatch(v))
}
