/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
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
	fmt "fmt"
	reg "regexp"
	sts "strings"
	utf "unicode/utf8"
)

// TOKENS

// This integer type is used as a type identifier for each token.
type TokenType int

// This enumeration defines all possible token types including the error token.
const (
	TokenError TokenType = iota
	TokenCatagory
	TokenComment
	TokenDelimiter
	TokenEOF
	TokenEOL
	TokenIdentifier
	TokenNote
	TokenSymbol
)

// This method returns the string representation for each token type.
func (v TokenType) String() string {
	return [...]string{
		"Error",
		"Catagory",
		"Comment",
		"Delimiter",
		"EOF",
		"EOL",
		"Identifier",
		"Note",
		"Symbol",
	}[v]
}

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
	case v.Type == TokenEOF:
		s = "<EOF>"
	case v.Type == TokenEOL:
		s = "<EOL>"
	case len(v.Value) > 60:
		s = fmt.Sprintf("%.60q...", v.Value)
	default:
		s = fmt.Sprintf("%q", v.Value)
	}
	return fmt.Sprintf("Token [type: %s, line: %d, position: %d]: %s", v.Type, v.Line, v.Position, s)
}

// SCANNER

// This constructor creates a new scanner initialized with the specified array
// of bytes. The scanner will scan in tokens matching the corresponding regular
// expressions.
func Scanner(source []byte, tokens chan Token) *scanner {
	var v = &scanner{source: source, line: 1, position: 1, tokens: tokens}
	go v.scanTokens() // Start scanning in the background.
	return v
}

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
	line      int // The line number in the source string of the next rune.
	position  int // The position in the current line of the first rune in the next token.
	tokens    chan Token
}

// This method continues scanning tokens from the source array until an error
// occurs or the end of file is reached. It then closes the token channel.
func (v *scanner) scanTokens() {
	for v.scanToken() {
	}
	close(v.tokens)
}

// This method attempts to scan any token starting with the next rune in the
// source array. It checks for each type of token as the cases for the switch
// statement. If that token type is found, this method returns true and skips
// the rest of the cases.  If no valid token is found, or a TokenEOF is found
// this method returns false.
func (v *scanner) scanToken() bool {
	v.skipSpaces()
	switch {
	case v.foundEOL():
	case v.foundCatagory():
	case v.foundComment():
	case v.foundDelimiter():
	case v.foundIdentifier():
	case v.foundNote():
	case v.foundSymbol():
	case v.foundEOF():      // Must be last.
		// We are at the end of the source array.
		return false
	default:
		// No valid token was found.
		v.foundError()
		return false
	}
	return true
}

// This method scans through any spaces in the source array and sets the next
// byte index to the next non-white-space rune.
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
		case "\t":
			tValue = "<TAB>"
		case "\f":
			tValue = "<FF>"
		case "\r":
			tValue = "<CR>"
		case "\v":
			tValue = "<VTAB>"
		}
	}
	var token = Token{tType, tValue, v.line, v.position}
	//fmt.Println(token)
	v.tokens <- token
	v.firstByte = v.nextByte
	v.position += sts.Count(tValue, "") - 1 // Add the number of runes in the token.
	return tType
}

// This method adds a catagory token with the current scanner information to the
// token channel. It returns true if a catagory token was found.
func (v *scanner) foundCatagory() bool {
	var s = v.source[v.nextByte:]
	var matches = scanCatagory(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenCatagory)
		return true
	}
	return false
}

// This method adds a comment token with the current scanner information
// to the token channel. It returns true if a comment token was found.
func (v *scanner) foundComment() bool {
	var s = v.source[v.nextByte:]
	var matches = scanComment(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenComment)
		return true
	}
	return false
}

// This method adds a delimiter token with the current scanner information to
// the token channel. It returns true if a delimiter token was found.
func (v *scanner) foundDelimiter() bool {
	var s = v.source[v.nextByte:]
	var matches = scanDelimiter(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenDelimiter)
		return true
	}
	return false
}

// This method adds an error token with the current scanner information to the
// token channel.
func (v *scanner) foundError() {
	var bytes = v.source[v.nextByte:]
	var _, width = utf.DecodeRune(bytes)
	v.nextByte += width
	v.emitToken(TokenError)
}

// This method adds an EOF token with the current scanner information to the
// token channel. It returns true if an EOF token was found.
func (v *scanner) foundEOF() bool {
	if v.nextByte == len(v.source) {
		v.emitToken(TokenEOF)
		return true
	}
	return false
}

// This method adds an EOL token with the current scanner information to the
// token channel. It returns true if an EOL token was found.
func (v *scanner) foundEOL() bool {
	var s = v.source[v.nextByte:]
	if byt.HasPrefix(s, []byte(EOL)) {
		v.nextByte++
		v.emitToken(TokenEOL)
		v.line++
		v.position = 1
		return true
	}
	return false
}

// This method adds an identifier token with the current scanner information
// to the token channel. It returns true if an identifier token was found.
func (v *scanner) foundIdentifier() bool {
	var s = v.source[v.nextByte:]
	var matches = scanIdentifier(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenIdentifier)
		return true
	}
	return false
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

// This method adds a symbol token with the current scanner information to the
// token channel. It returns true if a symbol token was found.
func (v *scanner) foundSymbol() bool {
	var s = v.source[v.nextByte:]
	var matches = scanSymbol(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenSymbol)
		return true
	}
	return false
}

// This scanner is used for matching catagory tokens.
var catagoryScanner = reg.MustCompile(`^(?:` + catagory + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a catagory token. The first string in the array is the
// entire matched string.
func scanCatagory(v []byte) []string {
	return bytesToStrings(catagoryScanner.FindSubmatch(v))
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
	result = append(result, string(v[:current])) // Includes bang delimeters.
	result = append(result, string(v[3:last]))   // Excludes bang delimeters.
	return result
}

// This function returns for the specified string an array of the matching
// subgroups for a delimiter. The first string in the array is the entire
// matched string.
func scanDelimiter(v []byte) []string {
	var result []string
	for _, delimiter := range delimiters {
		if byt.HasPrefix(v, delimiter) {
			result = append(result, string(delimiter))
		}
	}
	return result
}

// This scanner is used for matching identifier tokens.
var identifierScanner = reg.MustCompile(`^(?:` + identifier + `)`)

// This function returns for the specified string an array of the matching
// subgroups for an identifier token. The first string in the array is the
// entire matched string.
func scanIdentifier(v []byte) []string {
	return bytesToStrings(identifierScanner.FindSubmatch(v))
}

// This scanner is used for matching note tokens.
var noteScanner = reg.MustCompile(`^(?:` + note + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a note token. The first string in the array is the
// entire matched string.
func scanNote(v []byte) []string {
	return bytesToStrings(noteScanner.FindSubmatch(v))
}

// This scanner is used for matching symbol tokens.
var symbolScanner = reg.MustCompile(`^(?:` + symbol + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a symbol token. The first string in the array is the
// entire matched string.
func scanSymbol(v []byte) []string {
	return bytesToStrings(symbolScanner.FindSubmatch(v))
}

// CONSTANT DEFINITIONS

// This is the POSIX standard end-of-line character constant.
const EOL = "\n"

// These constant definitions capture regular expression subpatterns.
const (
	catagory   = `LETTER|DIGIT`
	digit      = `\pN` // All unicode digits.
	eol        = `\n`
	identifier = letter + `(?:` + letter + `|` + digit + `)*`
	letter     = `\pL` // All unicode letters and connectors like underscores.
	note       = `! [^` + eol + `]*`
	symbol     = `\$` + identifier
)

// This array contains the set of delimiters that may be used to separate the
// other tokens.
var delimiters = [][]byte{
	[]byte(`..`),
	[]byte(`! `),
	[]byte(`~`),
	[]byte(`:`),
	[]byte(`|`),
	[]byte(`'`),
	[]byte(`"`),
	[]byte(`(`),
	[]byte(`)`),
	[]byte(`[`),
	[]byte(`]`),
	[]byte(`{`),
	[]byte(`}`),
	[]byte(`<`),
	[]byte(`>`),
}

// PRIVATE FUNCTIONS

func bytesToStrings(bytes [][]byte) []string {
	var strings = make([]string, len(bytes))
	for index, array := range bytes {
		strings[index] = string(array)
	}
	return strings
}
