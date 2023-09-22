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

const TokenNote TokenType = 7

// This private method appends a formatted note to the result.
func (v *formatter) formatNote(note Note) {
	v.appendString(string(note))
}

// This method adds a note token with the current scanner information to the
// token channel. It returns true if a note token was found.
func (v *scanner) foundNote() bool {
	var s = v.source[v.nextByte:]
	var matches = scanNote(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenNote)
		return true
	}
	return false
}

// This method attempts to parse a note. It returns the note and whether or not
// the note was successfully parsed.
func (v *parser) parseNote() (Note, *Token, bool) {
	var note Note
	var token = v.nextToken()
	if token.Type != TokenNote {
		v.backupOne()
		return note, token, false
	}
	note = Note(token.Value)
	return note, token, true
}

// This scanner is used for matching note tokens.
var noteScanner = reg.MustCompile(`^(?:` + note + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a note token. The first string in the array is the
// entire matched string.
func scanNote(v []byte) []string {
	return bytesToStrings(noteScanner.FindSubmatch(v))
}
