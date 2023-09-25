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

type Note string

const TokenNote TokenType = "Note"
const note = `! [^` + eol + `]*`

// This scanner is used for matching note tokens.
var noteScanner = reg.MustCompile(`^(?:` + note + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a note token. The first string in the array is the
// entire matched string.
func scanNote(v []byte) []string {
	return bytesToStrings(noteScanner.FindSubmatch(v))
}
