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
	byt "bytes"
	fmt "fmt"
	reg "regexp"
	sts "strings"
	utf "unicode/utf8"
)

// TOKENS

// This string type is used as a type identifier for each type of token.
type TokenType string

// This enumeration defines all possible token types including the error token.
const (
	TokenEOL TokenType = "EOL"
	TokenEOF TokenType = "EOF"
	TokenError TokenType = "Error"
)

// This type defines the structure and methods for each token returned by the
// scanner.
type Token struct {
	Type     TokenType
	Value    string
	Line     int // The line number of the token in the input string.
	Position int // The position in the line of the first rune of the token.
}

// This method returns the canonical string version of this token.
func (v Token) String() string {
	var s string
	switch {
	case len(v.Value) > 60:
		s = fmt.Sprintf("%.60q...", v.Value)
	default:
		s = fmt.Sprintf("%q", v.Value)
	}
	return fmt.Sprintf("Token [type: %s, line: %d, position: %d]: %s", v.Type, v.Line, v.Position, s)
}

// SCANNER INTERFACE

// This constructor creates a new scanner initialized with the specified array
// of bytes. The scanner will scan in tokens matching the corresponding regular
// expressions.
func Scanner(source []byte, tokens chan Token) *scanner {
	var v = &scanner{source: source, line: 1, position: 1, tokens: tokens}
	go v.scanTokens() // Start scanning in the background.
	return v
}

// SCANNER IMPLEMENTATION

// This type defines the structure and methods for the scanner agent. The source
// bytes can be viewed like this:
//
//   | byte 0 | byte 1 | byte 2 | byte 3 | byte 4 | byte 5 | ... | byte N-1 |
//   | rune 0 |      rune 1     |      rune 2     | rune 3 | ... | rune R-1 |
//
// Runes can be one to eight bytes long.

type scanner struct {
	source    []byte
	firstByte int // The zero based index of the first possible byte in the next token.
	nextByte  int // The zero based index of the next possible byte in the next token.
	line      int // The line number in the source bytes of the next rune.
	position  int // The position in the current line of the first rune in the next token.
	tokens    chan Token
}

// This method continues scanning tokens from the source bytes until an error
// occurs or the end of file is reached. It then closes the token channel.
func (v *scanner) scanTokens() {
	for v.scanToken() {
	}
	close(v.tokens)
}

// This method attempts to scan any token starting with the next rune in the
// source bytes. It checks for each type of token as the cases for the switch
// statement. If that token type is found, this method returns true and skips
// the rest of the cases.  If no valid token is found, or a TokenEOF is found
// this method returns false.
func (v *scanner) scanToken() bool {
	v.skipSpaces()
	// Look for grammar specific tokens.
	for _, function := range foundFunctions {
		if function(v) {
			found = true
			break
		}
	}
	if found {
		return true
	}
	// Look for a POSIX end-of-file marker.
	found = foundEOF(v)
	if found {
		// We are at the end of the source bytes so signal to end scanning.
		return false
	}
	foundError(v)
	// We found an illegal character so signal to end scanning.
	return false
}

// This method scans through any spaces in the source bytes and sets the next
// byte index to the next non-space rune.
func (v *scanner) skipSpaces() {
	if v.nextByte < len(v.source) {
		for {
			if v.source[v.nextByte] != ' ' {
				break
			}
			v.nextByte++
			v.position++
		}
		v.firstByte = v.nextByte
	}
}

// This method adds a token of the specified type with the current scanner
// information to the token channel. It then resets the first byte index to the
// next byte index position. It returns the token type of the type added to the
// channel.
func (v *scanner) emitToken(tType TokenType) TokenType {
	var tValue = string(v.source[v.firstByte:v.nextByte])
	if tType == TokenEOF {
		tValue = "<EOF>"
	}
	if tType == TokenError {
		switch tValue {
		case "\a":
			tValue = "<BELL>"
		case "\b":
			tValue = "<BKSP>"
		case "\f":
			tValue = "<FF>"
		case "\n":
			tValue = "<EOL>"
		case "\t":
			tValue = "<TAB>"
		case "\r":
			tValue = "<CR>"
		case "\v":
			tValue = "<VTAB>"
		}
	}
	var token = Token{tType, tValue, v.line, v.position}
	//fmt.Println(token) // Uncomment this line for debugging purposes.
	v.tokens <- token
	v.firstByte = v.nextByte
	v.position += sts.Count(tValue, "") - 1 // Add the number of runes in the token.
	return tType
}

// PRIVATE FUNCTIONS

// These constant definitions capture regular expression subpatterns.
const (
	intrinsic  = `LOWERCASE|UPPERCASE|DIGIT|EOL|EOF`
	lowercase  = `\p{Ll}` // All unicode lowercase letters.
	uppercase  = `\p{Lu}` // All unicode upppercase letters.
	digit      = `\p{Nd}` // All unicode digits.
	eol        = `\n` // POSIX standard end-of-line character.
)

// This type defines the signature for all found<Token>() functions.
type foundFunction func(v *scanner) bool

// This array lists the found functions in their precedence order.
var foundFunctions = [...]foundFunction{
	foundIntrinsic,
	#foundFunctions#
	foundEOL,
}

// This function adds an error token with the current scanner information to the
// token channel.
func foundError(v *scanner) bool {
	var bytes = v.source[v.nextByte:]
	var _, width = utf.DecodeRune(bytes)
	v.nextByte += width
	v.emitToken(TokenError)
	return true
}

// This function adds an EOF token with the current scanner information to the
// token channel. It returns true if an EOF token was found.
func foundEOF(v *scanner) bool {
	var s = v.source[v.nextByte:]
	// The last byte in a POSIX standard file must be an EOL character.
	if byt.HasPrefix(s, []byte(EOL)) && v.nextByte+1 == len(v.source) {
		v.nextByte++
		v.emitToken(TokenEOF)
		return true
	}
	return false
}

// This function adds an EOL token with the current scanner information to the
// token channel. It returns true if an EOL token was found.
func foundEOL(v *scanner) bool {
	var s = v.source[v.nextByte:]
	// A normal EOL character cannot be the last byte in the file.
	if byt.HasPrefix(s, []byte(EOL)) && v.nextByte+1 < len(v.source) {
		v.nextByte++
		v.emitToken(TokenEOL)
		v.line++
		v.position = 1
		return true
	}
	return false
}

// to the token channel. It returns true if an intrinsic token was found.
func foundIntrinsic(v *scanner) bool {
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

// This function converts an array of byte arrays into an array of strings.
func bytesToStrings(bytes [][]byte) []string {
	var strings = make([]string, len(bytes))
	for index, array := range bytes {
		strings[index] = string(array)
	}
	return strings
}
