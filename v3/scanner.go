/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package cdsn

import (
	//fmt "fmt"
	reg "regexp"
	sts "strings"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type scannerClass_ struct {
	characterMatcher *reg.Regexp
	commentMatcher   *reg.Regexp
	delimiterMatcher *reg.Regexp
	eolMatcher       *reg.Regexp
	intrinsicMatcher *reg.Regexp
	literalMatcher   *reg.Regexp
	nameMatcher      *reg.Regexp
	noteMatcher      *reg.Regexp
	numberMatcher    *reg.Regexp
	spaceMatcher     *reg.Regexp
	symbolMatcher    *reg.Regexp
}

// Private Class Namespace Reference

var scannerClass = &scannerClass_{
	characterMatcher: reg.MustCompile(`^(?:` + character_ + `)`),
	commentMatcher:   reg.MustCompile(`^(?:` + comment_ + `)`),
	delimiterMatcher: reg.MustCompile(`^(?:` + delimiter_ + `)`),
	eolMatcher:       reg.MustCompile(`^(?:` + eol_ + `)`),
	intrinsicMatcher: reg.MustCompile(`^(?:` + intrinsic_ + `)`),
	literalMatcher:   reg.MustCompile(`^(?:` + literal_ + `)`),
	nameMatcher:      reg.MustCompile(`^(?:` + name_ + `)`),
	noteMatcher:      reg.MustCompile(`^(?:` + note_ + `)`),
	numberMatcher:    reg.MustCompile(`^(?:` + number_ + `)`),
	spaceMatcher:     reg.MustCompile(`^(?:` + space_ + `)`),
	symbolMatcher:    reg.MustCompile(`^(?:` + symbol_ + `)`),
}

// Public Class Namespace Access

func ScannerClass() *scannerClass_ {
	return scannerClass
}

// Public Class Constants

func (c *scannerClass_) GetCharacterMatcher() *reg.Regexp {
	return c.characterMatcher
}

func (c *scannerClass_) GetCommentMatcher() *reg.Regexp {
	return c.commentMatcher
}

func (c *scannerClass_) GetDelimiterMatcher() *reg.Regexp {
	return c.delimiterMatcher
}

func (c *scannerClass_) GetEOLMatcher() *reg.Regexp {
	return c.eolMatcher
}

func (c *scannerClass_) GetIntrinsicMatcher() *reg.Regexp {
	return c.intrinsicMatcher
}

func (c *scannerClass_) GetLiteralMatcher() *reg.Regexp {
	return c.literalMatcher
}

func (c *scannerClass_) GetNameMatcher() *reg.Regexp {
	return c.nameMatcher
}

func (c *scannerClass_) GetNoteMatcher() *reg.Regexp {
	return c.noteMatcher
}

func (c *scannerClass_) GetNumberMatcher() *reg.Regexp {
	return c.numberMatcher
}

func (c *scannerClass_) GetSpaceMatcher() *reg.Regexp {
	return c.spaceMatcher
}

func (c *scannerClass_) GetSymbolMatcher() *reg.Regexp {
	return c.symbolMatcher
}

// Public Class Constructors

func (c *scannerClass_) FromDocument(
	document string,
	tokens chan *token_,
) *scanner_ {
	var scanner = &scanner_{
		line:     1,
		position: 1,
		runes:    []rune(document),
		tokens:   tokens,
	}
	go scanner.scanTokens() // Start scanning tokens in the background.
	return scanner
}

// Public Class Functions

func (c *scannerClass_) MatchCharacter(text string) []string {
	return c.characterMatcher.FindStringSubmatch(text)
}

func (c *scannerClass_) MatchComment(text string) []string {
	return c.commentMatcher.FindStringSubmatch(text)
}

func (c *scannerClass_) MatchDelimiter(text string) []string {
	return c.delimiterMatcher.FindStringSubmatch(text)
}

func (c *scannerClass_) MatchEOL(text string) []string {
	return c.eolMatcher.FindStringSubmatch(text)
}

func (c *scannerClass_) MatchIntrinsic(text string) []string {
	return c.intrinsicMatcher.FindStringSubmatch(text)
}

func (c *scannerClass_) MatchLiteral(text string) []string {
	return c.literalMatcher.FindStringSubmatch(text)
}

func (c *scannerClass_) MatchName(text string) []string {
	return c.nameMatcher.FindStringSubmatch(text)
}

func (c *scannerClass_) MatchNote(text string) []string {
	return c.noteMatcher.FindStringSubmatch(text)
}

func (c *scannerClass_) MatchNumber(text string) []string {
	return c.numberMatcher.FindStringSubmatch(text)
}

func (c *scannerClass_) MatchSpace(text string) []string {
	return c.spaceMatcher.FindStringSubmatch(text)
}

func (c *scannerClass_) MatchSymbol(text string) []string {
	return c.symbolMatcher.FindStringSubmatch(text)
}

// CLASS INSTANCES

// Private Class Type Definition

type scanner_ struct {
	first    int // A zero based index of the first possible rune in the next token.
	line     int // The line number in the document of the next rune.
	next     int // A zero based index of the next possible rune in the next token.
	position int // The position in the current line of the next rune.
	runes    []rune
	tokens   chan *token_
}

// Private Interface

// This private class method adds a token of the specified type with the current
// scanner information to the token channel. It then resets the first rune index
// to the next rune index position. It returns the token type of the type added
// to the channel.
func (v *scanner_) emitToken(tokenType string) string {
	var tokenValue = string(v.runes[v.first:v.next])
	switch tokenValue {
	case "\a":
		tokenValue = "<BELL>"
	case "\b":
		tokenValue = "<BKSP>"
	case "\t":
		tokenValue = "<TAB>"
	case "\f":
		tokenValue = "<FF>"
	case "\n":
		tokenValue = "<EOL>"
	case "\r":
		tokenValue = "<CR>"
	case "\v":
		tokenValue = "<VTAB>"
	}
	var token = TokenClass().FromContext(v.line, v.position, tokenType, tokenValue)
	//fmt.Println(token) // Uncomment when debugging.
	v.tokens <- token
	v.position += v.next - v.first
	v.first = v.next
	return tokenType
}

func (v *scanner_) foundCharacter() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.characterMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetCharacter())
		return true
	}
	return false
}

func (v *scanner_) foundComment() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.commentMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetComment())
		v.line += sts.Count(matches[0], "\n")
		v.position = 1
		return true
	}
	return false
}

func (v *scanner_) foundDelimiter() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.delimiterMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetDelimiter())
		return true
	}
	return false
}

func (v *scanner_) foundEOL() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.eolMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetEOL())
		v.line++
		v.position = 1
		return true
	}
	return false
}

func (v *scanner_) foundError() {
	v.next++
	v.emitToken(TokenClass().GetError())
}

func (v *scanner_) foundIntrinsic() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.intrinsicMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetIntrinsic())
		return true
	}
	return false
}

func (v *scanner_) foundLiteral() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.literalMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetLiteral())
		return true
	}
	return false
}

func (v *scanner_) foundName() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.nameMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetName())
		return true
	}
	return false
}

func (v *scanner_) foundNote() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.noteMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetNote())
		return true
	}
	return false
}

func (v *scanner_) foundNumber() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.numberMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetNumber())
		return true
	}
	return false
}

func (v *scanner_) foundSpace() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.spaceMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.position += v.next - v.first
		v.first = v.next
		// Don't pass spaces along to the parser.
		return true
	}
	return false
}

func (v *scanner_) foundSymbol() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.symbolMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetSymbol())
		return true
	}
	return false
}

func (v *scanner_) scanTokens() {
loop:
	for v.next < len(v.runes) {
		switch {
		case v.foundCharacter():
		case v.foundComment():
		case v.foundDelimiter():
		case v.foundEOL():
		case v.foundIntrinsic():
		case v.foundLiteral():
		case v.foundName():
		case v.foundNote():
		case v.foundNumber():
		case v.foundSpace():
		case v.foundSymbol():
		default:
			v.foundError()
			break loop
		}
	}
	v.emitToken(TokenClass().GetEOF())
	close(v.tokens)
}

// These private constants define the regular expression sub-patterns that make
// up all token types.  Unfortunately there is no way to make them private to
// the scanner class namespace since they must be TRUE Go constants to be
// initialized in this way.  We add an underscore to lessen the chance of a name
// collision with other private Go class constants.
const (
	any_       = `.|` + eol_
	base16_    = `[0-9a-f]`
	character_ = `['][^` + control_ + `][']`
	comment_   = `!>(?:` + any_ + `)*?<!`
	control_   = `\p{Cc}`
	delimiter_ = `[~?*+:|(){}]|\.\.`
	digit_     = `\p{Nd}`
	eol_       = `\n`
	escape_    = `\\(?:(?:` + unicode_ + `)|[abfnrtv'"\\])`
	intrinsic_ = `ANY|LOWER|UPPER|DIGIT|ESCAPE|CONTROL|EOL|EOF`
	letter_    = lower_ + `|` + upper_
	literal_   = `["](?:` + escape_ + `|[^"` + eol_ + `])+["]`
	lower_     = `\p{Ll}`
	name_      = `(?:` + letter_ + `)(?:_?(?:` + letter_ + `|` + digit_ + `))*`
	note_      = `! [^` + eol_ + `]*`
	number_    = `(?:` + digit_ + `)+`
	space_     = `[ ]+`
	symbol_    = `\$(` + name_ + `)`
	unicode_   = `x` + base16_ + `{2}|u` + base16_ + `{4}|U` + base16_ + `{8}`
	upper_     = `\p{Lu}`
)
