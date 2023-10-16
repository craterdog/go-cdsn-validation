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
	//fmt "fmt"
	sts "strings"
	utf "unicode/utf8"
)

// SCANNER INTERFACE

// The POSIX standard end-of-line character.
const EOL = "\n"

// These constants define the token types that don't have regular expression
// patterns.
const (
	TokenEOF   TokenType = "EOF"
	TokenERROR TokenType = "ERROR"
)

// This function creates a new scanner initialized with the specified array
// of bytes. The scanner will automatically generating tokens that match the
// corresponding regular expressions.
func ScanTokens(source []byte, tokens chan Token) *scanner {
	var v = &scanner{source: source, line: 1, position: 1, tokens: tokens}
	go v.generateTokens() // Start scanning in the background.
	return v
}

// SCANNER IMPLEMENTATION

// This private function converts an array of byte arrays into an array of
// strings.
func bytesToStrings(bytes [][]byte) []string {
	var strings = make([]string, len(bytes))
	for index, array := range bytes {
		strings[index] = string(array)
	}
	return strings
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
	line      int // The line number in the source bytes of the next rune.
	position  int // The position in the current line of the first rune in the next token.
	tokens    chan Token
}

// This method adds a new error token with the current scanner information
// to the token channel.
func (v *scanner) atError() {
	var bytes = v.source[v.nextByte:]
	var _, width = utf.DecodeRune(bytes)
	v.nextByte += width
	v.emitToken(TokenERROR)
}

// This method determines whether or not the scanner is at the end of the source
// bytes and adds an EOF token with the current scanner information to the token
// channel if it is at the end.
func (v *scanner) atEOF() bool {
	if v.nextByte == len(v.source) {
		// The last byte in a POSIX standard file must be an EOL character.
		if byt.HasPrefix(v.source[v.nextByte-1:], []byte(EOL)) {
			v.emitToken(TokenEOF)
			return true
		}
	}
	return false
}

// This method adds a token of the specified type with the current scanner
// information to the token channel. It then resets the first byte index to the
// next byte index position. It returns the token type of the type added to the
// channel.
func (v *scanner) emitToken(tType TokenType) {
	var tValue = string(v.source[v.firstByte:v.nextByte])
	if tType == TokenEOF {
		tValue = "<EOFL>"
	}
	if tType == TokenERROR {
		switch tValue {
		case "\a":
			tValue = "<BELL>"
		case "\b":
			tValue = "<BKSP>"
		case "\t":
			tValue = "<HTAB>"
		case "\f":
			tValue = "<FMFD>"
		case "\n":
			tValue = "<EOLN>"
		case "\r":
			tValue = "<CRTN>"
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
	case v.atEOF():
		// We are at the end of the source bytes.
		return false
	case v.scanWHITESPACE():
	case v.scanDELIMITER():
	case v.scanINTRINSIC():
	case v.scanNOTE():
	case v.scanCOMMENT():
	case v.scanNUMBER():
	case v.scanCHARACTER():
	case v.scanLITERAL():
	case v.scanNAME():
	case v.scanSYMBOL():
	case v.scanLIMIT():
	default:
		// No valid token was found.
		v.atError()
		return false
	}
	// Successfully processed a token.
	return true
}

// This method tells the scanner to ignore any whitespace.  It returns true if
// whitespace was found.
func (v *scanner) scanWHITESPACE() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(whitespaceScanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.firstByte = v.nextByte
		var length = len(matches[0])
		var eolCount = sts.Count(matches[0], EOL)
		if eolCount > 0 {
			v.line += eolCount
			var index = sts.LastIndex(matches[0], EOL)
			v.position = length - index
		} else {
			v.position += length
		}
		return true
	}
	return false
}

// This method adds a new delimiter token with the current scanner information
// to the token channel. It returns true if a new delimiter token was found.
func (v *scanner) scanDELIMITER() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(delimiterScanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenDELIMITER)
		v.line += sts.Count(matches[0], EOL)
		return true
	}
	return false
}

// This method adds a new intrinsic token with the current scanner information
// to the token channel. It returns true if a new intrinsic token was found.
func (v *scanner) scanINTRINSIC() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(intrinsicScanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenINTRINSIC)
		v.line += sts.Count(matches[0], EOL)
		return true
	}
	return false
}

// This method adds a new note token with the current scanner information
// to the token channel. It returns true if a new note token was found.
func (v *scanner) scanNOTE() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(noteScanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenNOTE)
		v.line += sts.Count(matches[0], EOL)
		return true
	}
	return false
}

// This method adds a new comment token with the current scanner information
// to the token channel. It returns true if a new comment token was found.
func (v *scanner) scanCOMMENT() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(commentScanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenCOMMENT)
		v.line += sts.Count(matches[0], EOL)
		return true
	}
	return false
}

// This method adds a new number token with the current scanner information
// to the token channel. It returns true if a new number token was found.
func (v *scanner) scanNUMBER() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(numberScanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenNUMBER)
		v.line += sts.Count(matches[0], EOL)
		return true
	}
	return false
}

// This method adds a new character token with the current scanner information
// to the token channel. It returns true if a new character token was found.
func (v *scanner) scanCHARACTER() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(characterScanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenCHARACTER)
		v.line += sts.Count(matches[0], EOL)
		return true
	}
	return false
}

// This method adds a new string token with the current scanner information
// to the token channel. It returns true if a new string token was found.
func (v *scanner) scanLITERAL() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(literalScanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenLITERAL)
		v.line += sts.Count(matches[0], EOL)
		return true
	}
	return false
}

// This method adds a new name token with the current scanner information
// to the token channel. It returns true if a new name token was found.
func (v *scanner) scanNAME() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(nameScanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenNAME)
		v.line += sts.Count(matches[0], EOL)
		return true
	}
	return false
}

// This method adds a new symbol token with the current scanner information
// to the token channel. It returns true if a new symbol token was found.
func (v *scanner) scanSYMBOL() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(symbolScanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenSYMBOL)
		v.line += sts.Count(matches[0], EOL)
		return true
	}
	return false
}

// This method adds a new limit token with the current scanner information
// to the token channel. It returns true if a new limit token was found.
func (v *scanner) scanLIMIT() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(limitScanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenLIMIT)
		v.line += sts.Count(matches[0], EOL)
		return true
	}
	return false
}
