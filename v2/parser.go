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
	stc "strconv"
	sts "strings"
	utf "unicode/utf8"
)

// PARSER INTERFACE

// This function parses the specified document source retrieved from a POSIX
// compliant file and returns the corresponding CDSN statements that were used
// to generate the document using the CDSN formatting capabilities.
// A POSIX compliant file must end with an EOF marker.
func ParseDocument(source []byte) col.Sequential[StatementLike] {
	var ok bool
	var token *Token
	var statements col.Sequential[StatementLike]
	var tokens = make(chan Token, 256)
	Scanner(source, tokens) // Starts scanning in a separate go routine.
	var p = &parser{
		source: source,
		next:   col.StackWithCapacity[*Token](4),
		tokens: tokens,
	}
	statements, token, ok = p.parseStatements()
	if !ok {
		var message = p.formatError(token)
		message += generateGrammar("statement",
			"$source",
			"$statement")
		panic(message)
	}
	_, token, ok = p.parseEOF()
	if !ok {
		var message = p.formatError(token)
		message += generateGrammar("EOF",
			"$source",
			"$statement")
		panic(message)
	}
	return statements
}

// PARSER IMPLEMENTATION

// This type defines the structure and methods for the parser agent.
type parser struct {
	source         []byte
	next           col.StackLike[*Token] // The stack of the retrieved tokens that have been put back.
	tokens         chan Token        // The queue of unread tokens coming from the scanner.
	p1, p2, p3, p4 *Token            // The previous four tokens that have been retrieved.
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

// This method attempts to parse a sequence of statements. It returns the
// sequence of statements and whether or not the sequence of statements was
// successfully parsed.
func (v *parser) parseStatements() (col.Sequential[StatementLike], *Token, bool) {
	var ok bool
	var token *Token
	var statement StatementLike
	var statements = col.List[StatementLike]()
	statement, token, ok = v.parseStatement()
	if !ok {
		// A grammar must have at least one statement.
		return statements, token, false
	}
	for {
		statements.AddValue(statement)
		statement, token, ok = v.parseStatement()
		if !ok {
			// No more statements.
			break
		}
	}
	return statements, token, true
}

// This method attempts to parse a sequence of statements. It returns the
// sequence of statements and whether or not the sequence of statements was
// successfully parsed.
func (v *parser) parseStatement() (StatementLike, *Token, bool) {
	var ok bool
	var token *Token
	var comment Comment
	var production ProductionLike
	var statement StatementLike
	comment, token, ok = v.parseComment()
	if !ok {
		production, token, ok = v.parseProduction()
		if !ok {
			// This is not a statement.
			return statement, token, false
		}
	}
	_, token, ok = v.parseEOL()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("EOL",
			"$statement",
			"$COMMENT",
			"$production")
		panic(message)
	}
	statement = Statement(comment, production)
	return statement, token, true
}

// This method attempts to parse an association between a key and value. It
// returns the association and whether or not the association was successfully
// parsed.
func (v *parser) parseAssociation() (AssociationLike[Key, Value], *Token, bool) {
	var ok bool
	var token *Token
	var key Key
	var value Value
	var association AssociationLike[Key, Value]
	key, token, ok = v.parsePrimitive()
	if !ok {
		// This is not an association.
		return association, token, false
	}
	_, token, ok = v.parseDelimiter(":")
	if !ok {
		// This is not an association.
		v.backupOne() // Put back the primitive key token.
		return association, token, false
	}
	// This must be an association.
	value, token, ok = v.parseValue()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("value",
			"$association",
			"$key",
			"$value")
		panic(message)
	}
	association = Association[Key, Value](key, value)
	return association, token, true
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
	if token.Type != TokenEOL {
		v.backupOne()
		return token, token, false
	}
	token = v.nextToken()
	if token.Type != TokenEOF {
		v.backupOne() // Put back the EOL character.
		v.backupOne()
		return token, token, false
	}
	return token, token, true
}

// This method attempts to parse the end-of-line (EOL) token. It returns
// the token and whether or not an EOF token was found.
func (v *parser) parseEOL() (*Token, *Token, bool) {
	var token = v.nextToken()
	if token.Type != TokenEOL {
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
	factor, token, ok = v.parseIntrinsic()
	if !ok {
		factor, token, ok = v.parseIdentifier()
	}
	if !ok {
		factor, token, ok = v.parseRange()
	}
	if !ok {
		factor, token, ok = v.parseInversion()
	}
	if !ok {
		factor, token, ok = v.parsePrecedence()
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

// This method attempts to parse a intrinsic. It returns the intrinsic and
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

// This method attempts to parse a literal. It returns the literal and whether
// or not the literal was successfully parsed.
func (v *parser) parseLiteral() (Literal, *Token, bool) {
	var literal Literal
	var token = v.nextToken()
	if token.Type != TokenLiteral {
		v.backupOne()
		return literal, token, false
	}
	literal = Literal(token.Value)
	return literal, token, true
}
