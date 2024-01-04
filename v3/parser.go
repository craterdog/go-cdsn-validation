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
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v3"
	sts "strings"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type parserClass_ struct {
	channelSize int
	stackSize   int
}

// Private Class Namespace Reference

var parserClass = &parserClass_{
	channelSize: 128,
	stackSize:   4,
}

// Public Class Namespace Access

func ParserClass() ParserClassLike {
	return parserClass
}

// Public Class Constructors

func (c *parserClass_) Default() ParserLike {
	var parser = &parser_{
		names: col.CatalogClass[string, string]().Empty(),
		next:  col.StackClass[*token_]().WithCapacity(c.stackSize),
	}
	return parser
}

// CLASS INSTANCES

// Private Class Type Definition

type parser_ struct {
	document string
	names    col.CatalogLike[string, string]
	next     col.StackLike[*token_] // A stack of unprocessed retrieved tokens.
	tokens   chan *token_           // A queue of unread tokens from the scanner.
}

// Public Interface

func (v *parser_) ParseDocument(document string) DocumentLike {
	// Start a scanner running in a separate Go routine.
	v.document = document
	v.tokens = make(chan *token_, parserClass.channelSize)
	ScannerClass().FromDocument(v.document, v.tokens)

	// Parse the tokens from the scanner.
	var grammar, token, ok = v.parseGrammar()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("grammar",
			"$document",
			"$grammar",
		)
		panic(message)
	}
	_, token, ok = v.parseEOF()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("EOF",
			"$document",
			"$grammar",
		)
		panic(message)
	}

	// Make sure all names have associated definitions.
	var iterator = v.names.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var name = association.GetKey()
		var definition = association.GetValue()
		if len(definition) == 0 {
			var message = fmt.Sprintf(
				"The grammar is missing a definition for name: %v\n",
				name,
			)
			panic(message)
		}
	}

	return DocumentClass().FromGrammar(grammar)
}

// Private Interface

// This private class method returns an error message containing the context for
// a parsing error.
func (v *parser_) formatError(token *token_) string {
	var message = fmt.Sprintf(
		"An unexpected token was received by the parser: %v\n",
		token,
	)
	var line = token.GetLine()
	var lines = sts.Split(v.document, "\n")

	message += "\033[36m"
	if line > 1 {
		message += fmt.Sprintf("%04d: ", line-1) + string(lines[line-2]) + "\n"
	}
	message += fmt.Sprintf("%04d: ", line) + string(lines[line-1]) + "\n"

	message += " \033[32m>>>─"
	var count = 0
	for count < token.GetPosition() {
		message += "─"
		count++
	}
	message += "⌃\033[36m\n"

	if line < len(lines) {
		message += fmt.Sprintf("%04d: ", line+1) + string(lines[line]) + "\n"
	}
	message += "\033[0m\n"

	return message
}

// This private class method is useful when creating scanner and parser error
// messages that include the required grammatical rules.
func (v *parser_) generateGrammar(expected string, symbols ...string) string {
	var message = "Was expecting '" + expected + "' from:\n"
	for _, symbol := range symbols {
		message += fmt.Sprintf(
			"  \033[32m%v: \033[33m%v\033[0m\n\n",
			symbol,
			grammar[symbol],
		)
	}
	return message
}

// This private class method attempts to read the next token from the token
// stream and return it.
func (v *parser_) getNextToken() *token_ {
	var next *token_
	if v.next.IsEmpty() {
		var token, ok = <-v.tokens
		if !ok {
			panic("The token channel terminated without an EOF token.")
		}
		next = token
		if next.GetType() == TokenClass().GetError() {
			var message = v.formatError(next)
			panic(message)
		}
	} else {
		next = v.next.RemoveTop()
	}
	return next
}

func (v *parser_) parseAlternative() (AlternativeLike, *token_, bool) {
	var ok bool
	var token *token_
	var alternative AlternativeLike
	var note string
	var factor FactorLike
	var factors = col.ListClass[FactorLike]().Empty()
	factor, token, ok = v.parseFactor()
	if !ok {
		// This is not an alternative.
		return alternative, token, false
	}
	for {
		factors.AppendValue(factor)
		factor, _, ok = v.parseFactor()
		if !ok {
			// There are no more factors.
			alternative = AlternativeClass().FromFactors(factors)
			note, _, ok = v.parseNote()
			if ok {
				alternative.SetNote(note)
			}
			return alternative, token, true
		}
	}
}

func (v *parser_) parseAssertion() (AssertionLike, *token_, bool) {
	var ok bool
	var token *token_
	var element ElementLike
	var glyph GlyphLike
	var precedence PrecedenceLike
	var assertion AssertionLike
	element, token, ok = v.parseElement()
	if ok {
		assertion = AssertionClass().FromElement(element)
		return assertion, token, true
	}
	glyph, token, ok = v.parseGlyph()
	if ok {
		assertion = AssertionClass().FromGlyph(glyph)
		return assertion, token, true
	}
	precedence, token, ok = v.parsePrecedence()
	if ok {
		assertion = AssertionClass().FromPrecedence(precedence)
		return assertion, token, true
	}
	return assertion, token, false
}

func (v *parser_) parseCardinality() (CardinalityLike, *token_, bool) {
	var ok bool
	var token *token_
	var cardinality CardinalityLike
	var constraint ConstraintLike
	_, token, ok = v.parseDelimiter("?")
	if ok {
		constraint = ConstraintClass().FromRange("0", "1")
		cardinality = CardinalityClass().FromConstraint(constraint)
		return cardinality, token, true
	}
	_, token, ok = v.parseDelimiter("*")
	if ok {
		constraint = ConstraintClass().FromRange("0", "")
		cardinality = CardinalityClass().FromConstraint(constraint)
		return cardinality, token, true
	}
	_, token, ok = v.parseDelimiter("+")
	if ok {
		constraint = ConstraintClass().FromRange("1", "")
		cardinality = CardinalityClass().FromConstraint(constraint)
		return cardinality, token, true
	}
	_, token, ok = v.parseDelimiter("{")
	if ok {
		constraint, token, ok = v.parseConstraint()
		if !ok {
			var message = v.formatError(token)
			message += v.generateGrammar("constraint",
				"$cardinality",
				"$constraint",
			)
			panic(message)
		}
		_, token, ok = v.parseDelimiter("}")
		if !ok {
			var message = v.formatError(token)
			message += v.generateGrammar("}",
				"$cardinality",
				"$constraint",
			)
			panic(message)
		}
		cardinality = CardinalityClass().FromConstraint(constraint)
		return cardinality, token, true
	}
	return cardinality, token, false
}

func (v *parser_) parseCharacter() (string, *token_, bool) {
	var character string
	var token = v.getNextToken()
	if token.GetType() != TokenClass().GetCharacter() {
		v.putBack(token)
		return character, token, false
	}
	character = token.GetValue()
	return character, token, true
}

func (v *parser_) parseComment() (string, *token_, bool) {
	var comment string
	var token = v.getNextToken()
	if token.GetType() != TokenClass().GetComment() {
		v.putBack(token)
		return comment, token, false
	}
	comment = token.GetValue()
	return comment, token, true
}

func (v *parser_) parseConstraint() (ConstraintLike, *token_, bool) {
	var ok bool
	var token *token_
	var constraint ConstraintLike
	var first, last string
	first, token, ok = v.parseNumber()
	if !ok {
		// This is not a constraint.
		return constraint, token, false
	}
	_, _, ok = v.parseDelimiter("..")
	if ok {
		last, token, _ = v.parseNumber() // The last number is optional.
		constraint = ConstraintClass().FromRange(first, last)
		return constraint, token, true
	}
	constraint = ConstraintClass().FromNumber(first)
	return constraint, token, true
}

func (v *parser_) parseDefinition() (DefinitionLike, *token_, bool) {
	var ok bool
	var token *token_
	var symbol string
	var expression ExpressionLike
	var definition DefinitionLike
	symbol, token, ok = v.parseSymbol()
	if !ok {
		// This is not a definition.
		return definition, token, false
	}
	var name = symbol[1:]
	var existing = v.names.GetValue(name)
	if len(existing) > 0 {
		var message = v.formatError(token)
		message += "This symbol has already been defined in this grammar:\n"
		message += "    " + existing + "\n"
		panic(message)
	}
	_, token, ok = v.parseDelimiter(":")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar(":",
			"$definition",
			"$expression",
		)
		panic(message)
	}
	expression, token, ok = v.parseExpression()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("expression",
			"$definition",
			"$expression",
		)
		panic(message)
	}
	definition = DefinitionClass().FromSymbolAndExpression(symbol, expression)
	var formatter = FormatterClass().Default()
	v.names.SetValue(name, formatter.FormatDefinition(definition))
	return definition, token, true
}

func (v *parser_) parseDelimiter(delimiter string) (string, *token_, bool) {
	var token = v.getNextToken()
	if token.GetType() == TokenClass().GetEOF() || token.GetValue() != delimiter {
		v.putBack(token)
		return delimiter, token, false
	}
	return delimiter, token, true
}

func (v *parser_) parseElement() (ElementLike, *token_, bool) {
	var ok bool
	var token *token_
	var element ElementLike
	var intrinsic string
	var name string
	var literal string
	intrinsic, token, ok = v.parseIntrinsic()
	if ok {
		element = ElementClass().FromIntrinsic(intrinsic)
		return element, token, true
	}
	literal, token, ok = v.parseLiteral()
	if ok {
		element = ElementClass().FromLiteral(literal)
		return element, token, true
	}
	name, token, ok = v.parseName()
	if ok {
		element = ElementClass().FromName(name)
		return element, token, true
	}
	return element, token, false
}

func (v *parser_) parseEOF() (*token_, *token_, bool) {
	var token = v.getNextToken()
	if token.GetType() != TokenClass().GetEOF() {
		v.putBack(token)
		return token, token, false
	}
	return token, token, true
}

func (v *parser_) parseEOL() (*token_, *token_, bool) {
	var token = v.getNextToken()
	if token.GetType() != TokenClass().GetEOL() {
		v.putBack(token)
		return token, token, false
	}
	return token, token, true
}

func (v *parser_) parseExpression() (ExpressionLike, *token_, bool) {
	var ok bool
	var token *token_
	var expression ExpressionLike
	var alternative AlternativeLike
	var alternatives = col.ListClass[AlternativeLike]().Empty()

	// Handle in-line case.
	alternative, token, ok = v.parseAlternative()
	if ok {
		for {
			alternatives.AppendValue(alternative)
			_, token, ok = v.parseDelimiter("|")
			if !ok {
				// No more alternatives.
				expression = ExpressionClass().FromAlternatives(alternatives)
				expression.SetMultilined(false)
				return expression, token, true
			}
			alternative, token, ok = v.parseAlternative()
			if !ok {
				var message = v.formatError(token)
				message += v.generateGrammar("alternative",
					"$expression",
					"$alternative",
				)
				panic(message)
			}
		}
	}

	// Handle multi-line case.
	_, _, ok = v.parseEOL()
	if !ok {
		// This is not an expression.
		return expression, token, false
	}
	alternative, token, ok = v.parseAlternative()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("alternative",
			"$expression",
			"$alternative",
		)
		panic(message)
	}
	for {
		alternatives.AppendValue(alternative)
		var _, token, ok = v.parseEOL()
		if !ok {
			var message = v.formatError(token)
			message += v.generateGrammar("EOL",
				"$expression",
				"$alternative",
			)
			panic(message)
		}
		alternative, token, ok = v.parseAlternative()
		if !ok {
			// No more alternatives.
			expression = ExpressionClass().FromAlternatives(alternatives)
			expression.SetMultilined(true)
			return expression, token, true
		}
	}
}

func (v *parser_) parseFactor() (FactorLike, *token_, bool) {
	var ok bool
	var token *token_
	var factor FactorLike
	var predicate PredicateLike
	var cardinality CardinalityLike
	predicate, token, ok = v.parsePredicate()
	if !ok {
		// This is not a factor.
		return factor, token, false
	}
	factor = FactorClass().FromPredicate(predicate)
	cardinality, token, ok = v.parseCardinality()
	if ok {
		// The cardinality is optional.
		factor.SetCardinality(cardinality)
	}
	return factor, token, true
}

func (v *parser_) parseGlyph() (GlyphLike, *token_, bool) {
	var ok bool
	var token *token_
	var glyph GlyphLike
	var first, last string
	first, token, ok = v.parseCharacter()
	if !ok {
		// This is not a glyph.
		return glyph, token, false
	}
	_, _, ok = v.parseDelimiter("..")
	if !ok {
		// The range of characters is optional.
		glyph = GlyphClass().FromCharacter(first)
		return glyph, token, true
	}
	last, token, ok = v.parseCharacter()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("glyph",
			"$glyph",
		)
		panic(message)
	}
	glyph = GlyphClass().FromRange(first, last)
	return glyph, token, true
}

func (v *parser_) parseGrammar() (GrammarLike, *token_, bool) {
	var ok bool
	var token *token_
	var grammar GrammarLike
	var statement StatementLike
	var statements = col.ListClass[StatementLike]().Empty()
	for {
		statement, token, ok = v.parseStatement()
		if !ok {
			// There are no more statements.
			grammar = GrammarClass().FromStatements(statements)
			return grammar, token, true
		}
		statements.AppendValue(statement)
		_, token, ok = v.parseEOL()
		if !ok {
			var message = v.formatError(token)
			message += v.generateGrammar("EOL",
				"$grammar",
				"$statement",
			)
			panic(message)
		}
		for ok {
			// Absorb any blank lines.
			_, _, ok = v.parseEOL()
		}
	}
}

func (v *parser_) parseIntrinsic() (string, *token_, bool) {
	var intrinsic string
	var token = v.getNextToken()
	if token.GetType() != TokenClass().GetIntrinsic() {
		v.putBack(token)
		return intrinsic, token, false
	}
	intrinsic = token.GetValue()
	return intrinsic, token, true
}

func (v *parser_) parseLiteral() (string, *token_, bool) {
	var literal string
	var token = v.getNextToken()
	if token.GetType() != TokenClass().GetLiteral() {
		v.putBack(token)
		return literal, token, false
	}
	literal = token.GetValue()
	return literal, token, true
}

func (v *parser_) parseName() (string, *token_, bool) {
	var name string
	var token = v.getNextToken()
	if token.GetType() != TokenClass().GetName() {
		v.putBack(token)
		return name, token, false
	}
	name = token.GetValue()
	var definition = v.names.GetValue(name) // Returns "" if not found.
	v.names.SetValue(name, definition)
	return name, token, true
}

func (v *parser_) parseNote() (string, *token_, bool) {
	var note string
	var token = v.getNextToken()
	if token.GetType() != TokenClass().GetNote() {
		v.putBack(token)
		return note, token, false
	}
	note = token.GetValue()
	return note, token, true
}

func (v *parser_) parseNumber() (string, *token_, bool) {
	var number string
	var token = v.getNextToken()
	if token.GetType() != TokenClass().GetNumber() {
		v.putBack(token)
		return number, token, false
	}
	number = token.GetValue()
	return number, token, true
}

func (v *parser_) parsePrecedence() (PrecedenceLike, *token_, bool) {
	var ok bool
	var token *token_
	var precedence PrecedenceLike
	var expression ExpressionLike
	_, token, ok = v.parseDelimiter("(")
	if !ok {
		// This is not a precedence.
		return precedence, token, false
	}
	expression, token, ok = v.parseExpression()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("expression",
			"$precedence",
			"$expression",
		)
		panic(message)
	}
	_, token, ok = v.parseDelimiter(")")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar(")",
			"$precedence",
			"$expression",
		)
		panic(message)
	}
	precedence = PrecedenceClass().FromExpression(expression)
	return precedence, token, true
}

func (v *parser_) parsePredicate() (PredicateLike, *token_, bool) {
	var ok bool
	var token *token_
	var isInverted bool
	var assertion AssertionLike
	var predicate PredicateLike
	_, _, isInverted = v.parseDelimiter("~")
	if isInverted {
		assertion, token, ok = v.parseAssertion()
		if ok {
			predicate = PredicateClass().FromAssertion(assertion, isInverted)
			return predicate, token, true
		}
		var message = v.formatError(token)
		message += v.generateGrammar("assertion",
			"$predicate",
			"$assertion",
		)
		panic(message)
	}
	assertion, token, ok = v.parseAssertion()
	if ok {
		predicate = PredicateClass().FromAssertion(assertion, isInverted)
		return predicate, token, true
	}
	return predicate, token, false
}

func (v *parser_) parseStatement() (StatementLike, *token_, bool) {
	var ok bool
	var token *token_
	var statement StatementLike
	var comment string
	var definition DefinitionLike
	comment, token, ok = v.parseComment()
	if ok {
		statement = StatementClass().FromComment(comment)
		return statement, token, true
	}
	definition, token, ok = v.parseDefinition()
	if ok {
		statement = StatementClass().FromDefinition(definition)
		return statement, token, true
	}
	return statement, token, false
}

func (v *parser_) parseSymbol() (string, *token_, bool) {
	var symbol string
	var token = v.getNextToken()
	if token.GetType() != TokenClass().GetSymbol() {
		v.putBack(token)
		return symbol, token, false
	}
	symbol = token.GetValue()
	return symbol, token, true
}

func (v *parser_) putBack(token *token_) {
	v.next.AddValue(token)
}

var grammar = map[string]string{
	"$document":    `grammar EOF  ! Terminated with an end-of-file marker.`,
	"$grammar":     `statement+`,
	"$statement":   `(COMMENT | definition) EOL+`,
	"$definition":  `SYMBOL ":" EOL? expression  ! This works for tokens and rules.`,
	"$expression":  `alternative ("|" alternative)*`,
	"$alternative": `factor+ (NOTE EOL)?`,
	"$factor":      `predicate cardinality?  ! The default cardinality is one.`,
	"$predicate":   `"~"? assertion`,
	"$assertion":   `element | glyph | precedence`,
	"$element":     `INTRINSIC | LITERAL | NAME`,
	"$glyph":       `CHARACTER (".." CHARACTER)?  ! The range of characters is inclusive.`,
	"$precedence":  `"(" expression ")"`,
	"$cardinality": `
      "?"  ! Zero or one instance of a predicate.
    | "*"  ! Zero or more instances of a predicate.
    | "+"  ! One or more instances of a predicate.
    | "{" constraint "}"  ! Constrains the number of instances of a predicate.`,
	"$constraint": `NUMBER (".." NUMBER?)?  ! The range of numbers is inclusive.`,
}
