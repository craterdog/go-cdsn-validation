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
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v2"
	sts "strings"
	uni "unicode"
)

// PARSER INTERFACE

// This function parses the specified document source retrieved from a POSIX
// compliant file and returns the corresponding CDSN grammar that was used
// to generate the document using the CDSN formatting capabilities.
// A POSIX compliant file must end with an EOF marker.
func ParseDocument(source []byte) GrammarLike {
	var ok bool
	var token *Token
	var grammar GrammarLike
	var tokens = make(chan Token, 256)
	Scanner(source, tokens) // Starts scanning in a separate go routine.
	var p = &parser{
		symbols: col.Catalog[Symbol, DefinitionLike](),
		source:  source,
		next:    col.StackWithCapacity[*Token](4),
		tokens:  tokens,
	}
	grammar, token, ok = p.parseGrammar()
	if !ok {
		var message = p.formatError(token)
		message += generateGrammar("statement",
			"$grammar",
			"$statement")
		panic(message)
	}
	var iterator = col.Iterator[col.Binding[Symbol, DefinitionLike]](p.symbols)
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var symbol = association.GetKey()
		var definition = association.GetValue()
		if definition == nil {
			panic(fmt.Sprintf("Missing a definition for symbol: %v\n", symbol))
		}
	}
	return grammar
}

// PARSER IMPLEMENTATION

// This type defines the structure and methods for the parser agent.
type parser struct {
	symbols        col.CatalogLike[Symbol, DefinitionLike]
	source         []byte
	next           col.StackLike[*Token] // The stack of the retrieved tokens that have been put back.
	tokens         chan Token            // The queue of unread tokens coming from the scanner.
	p1, p2, p3, p4 *Token                // The previous four tokens that have been retrieved.
	isToken        bool                  // Whether or not the current definition is a token definition.
}

// This method puts back the current token onto the token stream so that it can
// be retrieved by another parsing method.
func (v *parser) backupOne() {
	v.next.AddValue(v.p1)
	v.p1, v.p2, v.p3, v.p4 = v.p2, v.p3, v.p4, nil
}

// This method returns an error message containing the context for a parsing
// error.
func (v *parser) formatError(token *Token) string {
	var message = fmt.Sprintf("An unexpected token was received by the parser: %v\n", token)
	var line = token.Line
	var lines = sts.Split(string(v.source), EOL)

	message += "\033[36m"
	if line > 1 {
		message += fmt.Sprintf("%04d: ", line-1) + string(lines[line-2]) + EOL
	}
	message += fmt.Sprintf("%04d: ", line) + string(lines[line-1]) + EOL

	message += " \033[32m>>>─"
	var count = 0
	for count < token.Position {
		message += "─"
		count++
	}
	message += "⌃\033[36m\n"

	if line < len(lines) {
		message += fmt.Sprintf("%04d: ", line+1) + string(lines[line]) + EOL
	}
	message += "\033[0m\n"

	return message
}

// This method attempts to read the next token from the token stream and return
// it.
func (v *parser) nextToken() *Token {
	var next *Token
	if v.next.IsEmpty() {
		var token, ok = <-v.tokens
		if !ok {
			panic("The token channel terminated without an EOF or error token.")
		}
		next = &token
		if next.Type == TokenError {
			var message = v.formatError(next)
			panic(message)
		}
	} else {
		next = v.next.RemoveTop()
	}
	v.p4, v.p3, v.p2, v.p1 = v.p3, v.p2, v.p1, next
	return next
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

// This method attempts to parse a comment. It returns the comment and whether
// or not a comment was successfully parsed.
func (v *parser) parseComment() (Comment, *Token, bool) {
	var comment Comment
	var token = v.nextToken()
	if token.Type != TokenComment {
		v.backupOne()
		return comment, token, false
	}
	comment = Comment(token.Value)
	return comment, token, true
}

// This method attempts to parse the specified delimiter. It returns
// the token and whether or not the delimiter was found.
func (v *parser) parseDelimiter(delimiter string) (string, *Token, bool) {
	var token = v.nextToken()
	if token.Type == TokenEOF || token.Value != delimiter {
		v.backupOne()
		return delimiter, token, false
	}
	return delimiter, token, true
}

// This method attempts to parse the end-of-file (EOF) marker. It returns
// the token and whether or not an EOF marker was found. Note that the POSIX
// standard requires that the last byte in a file be an end-of-line (EOL)
// character.
func (v *parser) parseEOF() (*Token, *Token, bool) {
	var token = v.nextToken()
	if token.Type != TokenEOF {
		v.backupOne()
		return token, token, false
	}
	return token, token, true
}

// This method attempts to parse a factor. It returns the factor and whether or
// not the factor was successfully parsed.
func (v *parser) parseFactor() (Factor, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	factor, token, ok = v.parseElement()
	if !ok {
		factor, token, ok = v.parseRange()
	}
	if !ok {
		factor, token, ok = v.parseInverse()
	}
	if !ok {
		factor, token, ok = v.parseExactlyN()
	}
	if !ok {
		factor, token, ok = v.parseZeroOrOne()
	}
	if !ok {
		factor, token, ok = v.parseZeroOrMore()
	}
	if !ok {
		factor, token, ok = v.parseOneOrMore()
	}
	return factor, token, ok
}

// This method attempts to parse an intrinsic. It returns the intrinsic and
// whether or not the intrinsic was successfully parsed.
func (v *parser) parseIntrinsic() (Intrinsic, *Token, bool) {
	var intrinsic Intrinsic
	var token = v.nextToken()
	if token.Type != TokenIntrinsic {
		v.backupOne()
		return intrinsic, token, false
	}
	intrinsic = Intrinsic(token.Value)
	return intrinsic, token, true
}

// This method attempts to parse a name token. It returns the token and
// whether or not a name token was found.
func (v *parser) parseName() (Name, *Token, bool) {
	var name Name
	var token = v.nextToken()
	if token.Type != TokenName {
		v.backupOne()
		return name, token, false
	}
	if v.isToken && uni.IsLower(rune(token.Value[0])) {
		panic(fmt.Sprintf("A token definition contains a rulename: %v\n", token.Value))
	}
	name = Name(token.Value)
	var symbol = Symbol("$" + token.Value)
	var definition = v.symbols.GetValue(symbol)
	v.symbols.SetValue(symbol, definition)
	return name, token, true
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

// This method attempts to parse a number. It returns the number and
// whether or not a number was successfully parsed.
func (v *parser) parseNumber() (Number, *Token, bool) {
	var number Number
	var token = v.nextToken()
	if token.Type != TokenNumber {
		v.backupOne()
		return number, token, false
	}
	number = Number(token.Value)
	return number, token, true
}

// This method attempts to parse a string. It returns the string and whether
// or not the string was successfully parsed.
func (v *parser) parseString() (String, *Token, bool) {
	var string_ String
	var token = v.nextToken()
	if token.Type != TokenString {
		v.backupOne()
		return string_, token, false
	}
	string_ = String(token.Value)
	return string_, token, true
}

// This method attempts to parse a symbol. It returns the symbol and
// whether or not the symbol was successfully parsed.
func (v *parser) parseSymbol() (Symbol, *Token, bool) {
	var symbol Symbol
	var token = v.nextToken()
	if token.Type != TokenSymbol {
		v.backupOne()
		return symbol, token, false
	}
	symbol = Symbol(token.Value)
	return symbol, token, true
}

// This method attempts to parse an element. It returns the element and whether
// or not the element was successfully parsed.
func (v *parser) parseElement() (Factor, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	factor, token, ok = v.parseIntrinsic()
	if !ok {
		factor, token, ok = v.parseString()
	}
	if !ok {
		factor, token, ok = v.parseNumber()
	}
	if !ok {
		factor, token, ok = v.parseName()
	}
	return factor, token, ok
}

// GRAMMAR UTILITIES

// This map captures the syntax expressions for Crater Dog Syntax Notation.
// It is useful when creating scanner and parser error messages.
var grammar_ = map[string]string{
	"$grammar":     `<statement> EOF  ! Terminated with an end-of-file marker.`,
	"$statement":   `COMMENT | definition`,
	"$definition":  `symbol ":" expression  ! This works for both tokens and rules.`,
	"$symbol":      `RULESYMBOL | TOKENSYMBOL`,
	"$expression":  `alternative {"|" alternative}`,
	"$alternative": `<factor> [NOTE]`,
	"$factor": `
      inverse
    | exactlyN
    | zeroOrOne
    | zeroOrMore
    | oneOrMore
    | range  ! Ranges must be parsed before RUNEs.
   	| token`,
	"$token":      `INTRINSIC | RUNE | STRING | NUMBER | NAME`,
	"$inverse":    `"~" factor`,
	"$exactlyN":   `"(" expression ")" [NUMBER]  ! The default is exactly one.`,
	"$zeroOrOne":  `"[" expression "]"`,
	"$zeroOrMore": `"{" expression "}"`,
	"$oneOrMore":  `"<" expression ">"`,
	"$range":      `RUNE ".." RUNE  ! A range includes the first and last RUNEs listed.`,
}

const header = `!>
    A formal expression of Crater Dog Syntax Notation™ (CDSN) using Crater Dog
    Syntax Notation™ itself.  This language grammar consists of rule expressions
    and token expressions.

    Each rule name begins with a lowercase letter.  The rules are applied in the
    order listed. So—for example—within a factor, a range of RUNEs takes
    precedence over an individual RUNE.

    Each token name begins with an uppercase letter.  The INTRINSIC tokens are
    environment and language specific, and are therefore left undefined. The
    tokens are also scanned in the order listed.  So an INTRINSIC token takes
    precedence over an IDENTIFIER token.
<!
`

func FormatGrammar() string {
	var builder sts.Builder
	builder.WriteString(header)
	var unsorted = make([]string, len(grammar_))
	var index = 0
	for key := range grammar_ {
		unsorted[index] = key
		index++
	}
	var keys = col.ListFromArray(unsorted)
	keys.SortValues()
	var iterator = col.Iterator[string](keys)
	for iterator.HasNext() {
		var key = iterator.GetNext()
		var value = grammar_[key]
		builder.WriteString(fmt.Sprintf("%s: %s\n\n", key, value))
	}
	return builder.String()
}

// PRIVATE FUNCTIONS

func generateGrammar(expected string, symbols ...string) string {
	var message = "Was expecting '" + expected + "' from:\n"
	for _, symbol := range symbols {
		message += fmt.Sprintf("  \033[32m%v: \033[33m%v\033[0m\n\n", symbol, grammar_[symbol])
	}
	return message
}
