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
	tokens         chan Token            // The queue of unread tokens coming from the scanner.
	p1, p2, p3, p4 *Token                // The previous four tokens that have been retrieved.
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

// This method attempts to parse an alternative. It returns the alternative
// and whether or not the alternative was successfully parsed.
func (v *parser) parseAlternative() (AlternativeLike, *Token, bool) {
	var ok bool
	var token *Token
	var note Note
	var option OptionLike
	var alternative AlternativeLike
	note, token, ok = v.parseNote()
	if ok {
		_, token, ok = v.parseEOL()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("EOL",
				"$alternative",
				"$NOTE",
				"$option")
			panic(message)
		}
	}
	option, token, ok = v.parseOption()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("option",
			"$alternative",
			"$NOTE",
			"$option")
		panic(message)
	}
	alternative = Alternative(note, option)
	return alternative, token, true
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
	factor, token, ok = v.parseRange() // This must be first.
	if !ok {
		factor, token, ok = v.parseIntrinsic() // This must be second.
	}
	if !ok {
		factor, token, ok = v.parseIdentifier()
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

// This method attempts to parse an identifier token. It returns
// the token and whether or not an identifier token was found.
func (v *parser) parseIdentifier() (*Token, *Token, bool) {
	var token = v.nextToken()
	if token.Type != TokenIdentifier {
		v.backupOne()
		return token, token, false
	}
	return token, token, true
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

// This method attempts to parse an inversion. It returns the inversion and
// whether or not the inversion was successfully parsed.
func (v *parser) parseInversion() (InversionLike, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	var inversion InversionLike
	_, token, ok = v.parseDelimiter("~")
	if !ok {
		// This is not an inversion.
		return inversion, token, false
	}
	factor, token, ok = v.parseFactor()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("factor",
			"$factor")
		panic(message)
	}
	inversion = Inversion(factor)
	return inversion, token, true
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

// This method attempts to parse a note. It returns the note and whether
// or not the note was successfully parsed.
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

// This method attempts to parse a one or more grouping. It returns the
// one or more grouping and whether or not the one or more grouping was
// successfully parsed.
func (v *parser) parseOneOrMore() (GroupingLike, *Token, bool) {
	var ok bool
	var token *Token
	var rule RuleLike
	var grouping GroupingLike
	_, token, ok = v.parseDelimiter("<")
	if !ok {
		// This is not a one or more grouping.
		return grouping, token, false
	}
	rule, token, ok = v.parseRule()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("rule",
			"$factor",
			"$rule")
		panic(message)
	}
	_, token, ok = v.parseDelimiter(">")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(">",
			"$factor",
			"$rule")
		panic(message)
	}
	grouping = Grouping(rule, OneOrMore)
	return grouping, token, true
}

// This method attempts to parse an option. It returns the option and whether or
// not the option was successfully parsed.
func (v *parser) parseOption() (OptionLike, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	var factors = col.List[Factor]()
	var option OptionLike
	factor, token, ok = v.parseFactor()
	if !ok {
		// An option must have at least one factor.
		return option, token, false
	}
	for {
		factors.AddValue(factor)
		factor, token, ok = v.parseFactor()
		if !ok {
			// No more factors.
			break
		}
	}
	option = Option(factors)
	return option, token, true
}

// This method attempts to parse a precedence grouping. It returns the
// precedence grouping and whether or not the precedence grouping was
// successfully parsed.
func (v *parser) parsePrecedence() (GroupingLike, *Token, bool) {
	var ok bool
	var token *Token
	var rule RuleLike
	var grouping GroupingLike
	_, token, ok = v.parseDelimiter("(")
	if !ok {
		// This is not a precedence grouping.
		return grouping, token, false
	}
	rule, token, ok = v.parseRule()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("rule",
			"$factor",
			"$rule")
		panic(message)
	}
	_, token, ok = v.parseDelimiter(")")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(")",
			"$factor",
			"$rule")
		panic(message)
	}
	grouping = Grouping(rule, Precedence)
	return grouping, token, true
}

// This method attempts to parse a production. It returns the production
// and whether or not the production was successfully parsed.
func (v *parser) parseProduction() (ProductionLike, *Token, bool) {
	var ok bool
	var token *Token
	var symbol Symbol
	var rule RuleLike
	var note Note
	var production ProductionLike
	symbol, token, ok = v.parseSymbol()
	if !ok {
		// This is not a production.
		return production, token, false
	}
	_, token, ok = v.parseDelimiter(":")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(":",
			"$production",
			"$SYMBOL",
			"$rule",
			"$NOTE")
		panic(message)
	}
	rule, token, ok = v.parseRule()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("rule",
			"$production",
			"$SYMBOL",
			"$rule",
			"$NOTE")
		panic(message)
	}
	note, token, ok = v.parseNote()
	production = Production(symbol, rule, note)
	return production, token, true
}

// This method attempts to parse a range. It returns the range and whether
// or not the range was successfully parsed.
func (v *parser) parseRange() (Range, *Token, bool) {
	var range_ Range
	var token = v.nextToken()
	if token.Type != TokenRange {
		v.backupOne()
		return range_, token, false
	}
	range_ = Range(token.Value)
	return range_, token, true
}

// This method attempts to parse a rule. It returns the rule and whether or not
// the rule was successfully parsed.
func (v *parser) parseRule() (RuleLike, *Token, bool) {
	var ok bool
	var token *Token
	var option OptionLike
	var alternative AlternativeLike
	var alternatives = col.List[AlternativeLike]()
	var rule RuleLike
	option, token, ok = v.parseOption()
	if !ok {
		// This is not a rule.
		return rule, token, false
	}
	for {
		_, token, ok = v.parseDelimiter("|")
		if !ok {
			// No more alternatives.
			break
		}
		alternative, token, ok = v.parseAlternative()
		alternatives.AddValue(alternative)
	}
	rule = Rule(option, alternatives)
	return rule, token, true
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

// This method attempts to parse a symbol. It returns the symbol and whether
// or not the symbol was successfully parsed.
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

// This method attempts to parse a zero or more grouping. It returns the
// zero or more grouping and whether or not the zero or more grouping was
// successfully parsed.
func (v *parser) parseZeroOrMore() (GroupingLike, *Token, bool) {
	var ok bool
	var token *Token
	var rule RuleLike
	var grouping GroupingLike
	_, token, ok = v.parseDelimiter("{")
	if !ok {
		// This is not a zero or more grouping.
		return grouping, token, false
	}
	rule, token, ok = v.parseRule()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("rule",
			"$factor",
			"$rule")
		panic(message)
	}
	_, token, ok = v.parseDelimiter("}")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("}",
			"$factor",
			"$rule")
		panic(message)
	}
	grouping = Grouping(rule, ZeroOrMore)
	return grouping, token, true
}

// This method attempts to parse a zero or one grouping. It returns the
// zero or one grouping and whether or not the zero or one grouping was
// successfully parsed.
func (v *parser) parseZeroOrOne() (GroupingLike, *Token, bool) {
	var ok bool
	var token *Token
	var rule RuleLike
	var grouping GroupingLike
	_, token, ok = v.parseDelimiter("[")
	if !ok {
		// This is not a zero or one grouping.
		return grouping, token, false
	}
	rule, token, ok = v.parseRule()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("rule",
			"$factor",
			"$rule")
		panic(message)
	}
	_, token, ok = v.parseDelimiter("]")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("]",
			"$factor",
			"$rule")
		panic(message)
	}
	grouping = Grouping(rule, ZeroOrOne)
	return grouping, token, true
}
