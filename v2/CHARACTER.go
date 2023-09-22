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

const TokenCharacter TokenType = "Character"
const character = `['][^'][']`

type Character string

// This private method appends a formatted character to the result.
func (v *formatter) formatCharacter(character Character) {
	v.appendString(string(character))
}

// This method adds a character token with the current scanner information
// to the token channel. It returns true if a character token was found.
func (v *scanner) foundCharacter() bool {
	var s = v.source[v.nextByte:]
	var matches = scanCharacter(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenCharacter)
		return true
	}
	return false
}

// This method attempts to parse a character. It returns the character and
// whether or not a character was successfully parsed.
func (v *parser) parseCharacter() (Character, *Token, bool) {
	var character Character
	var token = v.nextToken()
	if token.Type != TokenCharacter {
		v.backupOne()
		return character, token, false
	}
	character = Character(token.Value)
	return character, token, true
}

// This scanner is used for matching character tokens.
var characterScanner = reg.MustCompile(`^(?:` + character + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a character token. The first string in the array is the
// entire matched string.
func scanCharacter(v []byte) []string {
	return bytesToStrings(characterScanner.FindSubmatch(v))
}
