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
	byt "bytes"
	sts "strings"
)

type Comment string

const TokenComment TokenType = "Comment"

// This private method appends a formatted comment to the result.
func (v *formatter) formatComment(comment Comment) {
	v.appendString(string(comment))
}

// This method adds a comment token with the current scanner information
// to the token channel. It returns true if a comment token was found.
func (v *scanner) foundComment() bool {
	var s = v.source[v.nextByte:]
	var matches = scanComment(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.line += sts.Count(matches[0], EOL)
		v.emitToken(TokenComment)
		return true
	}
	return false
}

// This function returns for the specified string an array of the matching
// subgroups for a commment token. The first string in the array is the entire
// matched string. Since a comment can be recursive a regular expression is not
// used in this implementation.
func scanComment(v []byte) []string {
	var result []string
	var space = []byte(" ")
	var eol = []byte(EOL)
	var bangAngle = []byte("!>" + EOL)
	var angleBang = []byte("<!")
	if !byt.HasPrefix(v, bangAngle) {
		return result
	}
	var angleBangAllowed = false
	var current = 3 // Skip the leading '!>\n' characters.
	var last = 0
	var level = 1
	for level > 0 {
		var s = v[current:]
		switch {
		case len(s) == 0:
			return result
		case byt.HasPrefix(s, eol):
			angleBangAllowed = true
			last = current
			current++
		case byt.HasPrefix(s, bangAngle):
			current += 3 // Skip the '!>\n' characters.
			level++      // Start a nested narrative.
		case byt.HasPrefix(s, angleBang):
			if !angleBangAllowed {
				return result
			}
			current += 2 // Skip the '<!' characters.
			level--      // Terminate the current narrative.
		default:
			if angleBangAllowed && !byt.HasPrefix(s, space) {
				angleBangAllowed = false
			}
			current++ // Accept the next character.
		}
	}
	result = append(result, string(v[:current])) // Includes bang literals.
	result = append(result, string(v[3:last]))   // Excludes bang literals.
	return result
}
