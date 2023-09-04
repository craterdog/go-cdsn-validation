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
		source: source,
		next:   col.StackWithCapacity[*Token](4),
		tokens: tokens,
	}
	grammar, token, ok = p.parseGrammar()
	if !ok {
		var message = p.formatError(token)
		message += generateGrammar("statement",
			"$grammar",
			"$statement")
		panic(message)
	}
	return grammar
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

// This method attempts to parse a count. It returns the count and whether or
// not the count was successfully parsed.
func (v *parser) parseCount() (CountLike, *Token, bool) {
	var ok bool
	var token *Token
	var digit Digit
	var digits = col.List[Digit]()
	var count CountLike
	for {
		digits.AddValue(digit)
		digit, token, ok = v.parseDigit()
		if !ok {
			// No more digits.
			break
		}
	}
	count = Count(digits)
	return count, token, true
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

// This method attempts to parse a digit. It returns the digit and
// whether or not a digit was successfully parsed.
func (v *parser) parseDigit() (Digit, *Token, bool) {
	var digit Digit
	var token = v.nextToken()
	if token.Type != TokenDigit {
		v.backupOne()
		return digit, token, false
	}
	digit = Digit(token.Value)
	return digit, token, true
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
		factor, token, ok = v.parseInversion()
	}
	if !ok {
		factor, token, ok = v.parseExactCount()
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
	if !ok {
		factor, token, ok = v.parseCharacter()
	}
	if !ok {
		factor, token, ok = v.parseLiteral()
	}
	if !ok {
		factor, token, ok = v.parseIntrinsic()
	}
	if !ok {
		factor, token, ok = v.parseIdentifier()
	}
	return factor, token, ok
}

// This method attempts to parse a grammar. It returns the grammar and whether
// or not the grammar was successfully parsed.
func (v *parser) parseGrammar() (GrammarLike, *Token, bool) {
	var ok bool
	var token *Token
	var statement StatementLike
	var statements = col.List[StatementLike]()
	var grammar GrammarLike
	statement, token, ok = v.parseStatement()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("statement",
			"$grammar",
			"$statement")
		panic(message)
	}
	for {
		statements.AddValue(statement)
		statement, _, ok = v.parseStatement()
		if !ok {
			// No more statements.
			break
		}
	}
	_, token, ok = v.parseEOF()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("EOF",
			"$grammar",
			"$statement")
		panic(message)
	}
	grammar = Grammar(statements)
	return grammar, token, true
}

// This method attempts to parse an identifier token. It returns
// the token and whether or not an identifier token was found.
func (v *parser) parseIdentifier() (Identifier, *Token, bool) {
	var identifier Identifier
	var token = v.nextToken()
	if token.Type != TokenIdentifier {
		v.backupOne()
		return identifier, token, false
	}
	identifier = Identifier(token.Value)
	return identifier, token, true
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
	var count, _, _ = v.parseCount()  // The count is optional.
	grouping = Grouping(rule, OneOrMore, count)
	return grouping, token, true
}

// This method attempts to parse an option. It returns the option and whether or
// not the option was successfully parsed.
func (v *parser) parseOption() (OptionLike, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	var factors = col.List[Factor]()
	var note Note
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
	note, _, _ = v.parseNote() // The note is optional.
	option = Option(factors, note)
	return option, token, true
}

// This method attempts to parse an exact count grouping. It returns the
// exact count grouping and whether or not the exact count grouping was
// successfully parsed.
func (v *parser) parseExactCount() (GroupingLike, *Token, bool) {
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
	var count, _, _ = v.parseCount()  // The count is optional.
	grouping = Grouping(rule, ExactCount, count)
	return grouping, token, true
}

// This method attempts to parse a production. It returns the production
// and whether or not the production was successfully parsed.
func (v *parser) parseProduction() (ProductionLike, *Token, bool) {
	var ok bool
	var token *Token
	var symbol Symbol
	var rule RuleLike
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
			"$rule")
		panic(message)
	}
	rule, token, ok = v.parseRule()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("rule",
			"$production",
			"$SYMBOL",
			"$rule")
		panic(message)
	}
	production = Production(symbol, rule)
	return production, token, true
}

// This method attempts to parse a range. It returns the range
// and whether or not the range was successfully parsed.
func (v *parser) parseRange() (RangeLike, *Token, bool) {
	var ok bool
	var token *Token
	var first Character
	var last Character
	var range_ RangeLike
	first, token, ok = v.parseCharacter()
	if !ok {
		// This is not a range.
		return range_, token, false
	}
	_, token, ok = v.parseDelimiter("..")
	if !ok {
		// This is not a range.
		v.backupOne() // Put back the character.
		return range_, token, false
	}
	last, token, ok = v.parseCharacter()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("CHARACTER",
			"$range",
			"$CHARACTER")
		panic(message)
	}
	range_ = Range(first, last)
	return range_, token, true
}

// This method attempts to parse a rule. It returns the rule and whether or not
// the rule was successfully parsed.
func (v *parser) parseRule() (RuleLike, *Token, bool) {
	var ok bool
	var token *Token
	var option OptionLike
	var options = col.List[OptionLike]()
	var rule RuleLike
	v.parseEOL() // The EOL is optional.
	option, token, ok = v.parseOption()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("option",
			"$rule",
			"$option")
		panic(message)
	}
	for {
		options.AddValue(option)
		v.parseEOL() // The EOL is optional.
		_, _, ok = v.parseDelimiter("|")
		if !ok {
			// No more options.
			break
		}
		option, token, ok = v.parseOption()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("option",
				"$rule",
				"$option")
			panic(message)
		}
	}
	rule = Rule(options)
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
	comment, _, ok = v.parseComment()
	if !ok {
		production, token, ok = v.parseProduction()
		if !ok {
			// This is not a statement.
			return statement, token, false
		}
	}
	for {
		_, _, ok = v.parseEOL()
		if !ok {
			// No more blank lines.
			break
		}
	}
	statement = Statement(comment, production)
	return statement, token, true
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
	var count, _, _ = v.parseCount()  // The count is optional.
	grouping = Grouping(rule, ZeroOrMore, count)
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
	var count, _, _ = v.parseCount()  // The count is optional.
	grouping = Grouping(rule, ZeroOrOne, count)
	return grouping, token, true
}

// GRAMMAR UTILITIES

// This map captures the syntax rules for Crater Dog Syntax Notation.
// It is useful when creating scanner and parser error messages.
var grammar_ = map[string]string{
	"$NOTE":        `"! " {~EOL}`,
	"$COMMENT":     `"!>" EOL  {COMMENT | ~"<!"} EOL "<!"`,
	"$CHARACTER":   `"'" ~"'" "'"`,
	"$LITERAL":     `'"' <~'"'> '"'`,
	"$INTRINSIC":   `"LETTER" | "DIGIT" | "EOL" | "EOF"`,
	"$IDENTIFIER":  `LETTER {LETTER | DIGIT}`,
	"$SYMBOL":      `"$" IDENTIFIER`,
	"$grammar":      `<statement> EOF`,
	"$statement":   `(COMMENT | production) <EOL>`,
	"$production":  `SYMBOL ":" rule [NOTE]`,
	"$rule":        `option {"|" alternative}`,
	"$option":      `<factor>`,
	"$alternative": `[[NOTE] EOL] option`,
	"$range":       `CHARACTER ".." CHARACTER`,
	"$factor": `
    range        |
    "~" factor   |
    "(" rule ")" |
    "[" rule "]" |
    "{" rule "}" |
    "<" rule ">" |
    CHARACTER    |
    LITERAL      |
    INTRINSIC    |
    IDENTIFIER`,
}

const header = `!>
    A formal definition of Crater Dog Syntax Notation™ (CDSN) using Crater Dog
    Syntax Notation™ itself. Token names are identified by all CAPITAL
    letters and rule names are identified by lowerCamelCase letters.

    The INTRINSIC tokens are environment dependent and therefore left undefined.
    The tokens are scanned in the order listed so an INTRINSIC token takes
    precedence over an IDENTIFIER token.

    The rules are applied in the order listed as well, so within a factor a
    range takes precedence over an individual CHARACTER.  The starting rule is
    the "$grammar" rule.
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
