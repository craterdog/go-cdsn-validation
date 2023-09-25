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

// This string type is used as a type identifier for each token.
type TokenType string

// This enumeration defines all possible token types including the error token.
const (
	TokenError   TokenType = "Error"
	TokenLiteral TokenType = "Literal"
	TokenEOF     TokenType = "EOF"
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
	case v.Type == TokenEOF:
		s = "<EOF>"
	case len(v.Value) > 60:
		s = fmt.Sprintf("%.60q...", v.Value)
	default:
		s = fmt.Sprintf("%q", v.Value)
	}
	return fmt.Sprintf("Token [type: %s, line: %d, position: %d]: %s", v.Type, v.Line, v.Position, s)
}

// SCANNER INTERFACE

// The POSIX standard end-of-line character.
const EOL = "\n"

// This function creates a new scanner initialized with the specified array
// of bytes. The scanner will automatically generating tokens that match the
// corresponding regular expressions.
func ScanTokens(source []byte, tokens chan Token) *scanner {
	var v = &scanner{source: source, line: 1, position: 1, tokens: tokens}
	go v.generateTokens() // Start scanning in the background.
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

// This method adds a token of the specified type with the current scanner
// information to the token channel. It then resets the first byte index to the
// next byte index position. It returns the token type of the type added to the
// channel.
func (v *scanner) emitToken(tType TokenType) {
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
		// The last byte in a POSIX standard file must be an EOL character.
		if byt.HasPrefix(v.source[v.nextByte-1:], []byte(EOL)) {
			v.emitToken(TokenEOF)
			return true
		}
	}
	return false
}

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

// This method adds a literal token with the current scanner information to
// the token channel. It returns true if a literal token was found.
func (v *scanner) foundLiteral() bool {
	var s = v.source[v.nextByte:]
	var matches = scanLiteral(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenLiteral)
		return true
	}
	return false
}

// This method adds a name token with the current scanner information
// to the token channel. It returns true if a name token was found.
func (v *scanner) foundName() bool {
	var s = v.source[v.nextByte:]
	var matches = scanName(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenName)
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

// This method adds a string token with the current scanner information to the
// token channel. It returns true if a string token was found.
func (v *scanner) foundString() bool {
	var s = v.source[v.nextByte:]
	var matches = scanString(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenString)
		return true
	}
	return false
}

// This method adds a symbol token with the current scanner information to
// the token channel. It returns true if a symbol token was found.
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

// This method tells the scanner to ignore any whitespace.  It returns true if
// whitespace was found.
func (v *scanner) foundWhitespace() bool {
	var s = v.source[v.nextByte:]
	var matches = scanWhitespace(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.firstByte = v.nextByte
		v.line += sts.Count(matches[0], EOL)
		v.position += len(matches[0]) - sts.LastIndex(matches[0], EOL) - 1
		return true
	}
	return false
}

// This method continues scanning tokens from the source bytes until an error
// occurs or the end of file is reached. It then closes the token channel.
func (v *scanner) generateTokens() {
	for v.processToken() {
	}
	close(v.tokens)
}

// This method attempts to scan any token starting with the next rune in the
// source bytes. It checks for each type of token as the cases for the switch
// statement. If that token type is found, this method returns true and skips
// the rest of the cases.  If no valid token is found, or a TokenEOF is found
// this method returns false.
func (v *scanner) processToken() bool {
	switch {
	case v.foundWhitespace():
	case v.foundIntrinsic():
	case v.foundNote():
	case v.foundComment():
	case v.foundCharacter():
	case v.foundString():
	case v.foundNumber():
	case v.foundName():
	case v.foundSymbol():
	case v.foundLiteral():
	case v.foundEOF():
		// We are at the end of the source bytes.
		return false
	default:
		// No valid token was found.
		v.foundError()
		return false
	}
	return true
}

// PRIVATE DEFINITIONS

const literal = `[~:|()[\]{}<>]|\.\.`

// This scanner is used for matching literal tokens.
var literalScanner = reg.MustCompile(`^(?:` + literal + `)`)

// This scanner is used for matching whitespace.
var whitespaceScanner = reg.MustCompile(`^(?:` + whitespace + `)`)

// PRIVATE FUNCTIONS

func bytesToStrings(bytes [][]byte) []string {
	var strings = make([]string, len(bytes))
	for index, array := range bytes {
		strings[index] = string(array)
	}
	return strings
}

// This function returns for the specified string an array of the matching
// subgroups for a literal. The first string in the array is the entire
// matched string.
func scanLiteral(v []byte) []string {
	return bytesToStrings(literalScanner.FindSubmatch(v))
}

// This function returns for the specified string an array of the matching
// subgroups for any whitespace. The first string in the array is the
// entire matched string.
func scanWhitespace(v []byte) []string {
	return bytesToStrings(whitespaceScanner.FindSubmatch(v))
}
