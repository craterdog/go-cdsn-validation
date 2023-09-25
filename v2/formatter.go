/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
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
	ref "reflect"
	sts "strings"
)

// FORMATTER INTERFACE

// This function returns the bytes containing the canonical format for the
// specified grammar including the POSIX standard EOF marker.
func FormatDocument(grammar GrammarLike) []byte {
	var v = &formatter{}
	v.formatGrammar(grammar)
	var string_ = v.getResult()
	return []byte(string_)
}

// FORMATTER IMPLEMENTATION

// This type defines the structure and methods for a canonical formatter agent.
type formatter struct {
	indentation int
	depth       int
	result      sts.Builder
}

// This method appends a properly indented newline to the result.
func (v *formatter) appendNewline() {
	var separator = "\n"
	var levels = v.depth + v.indentation
	for level := 0; level < levels; level++ {
		separator += "    "
	}
	v.result.WriteString(separator)
}

// This method appends the specified string to the result.
func (v *formatter) appendString(s string) {
	v.result.WriteString(s)
}

// This private method appends a formatted alternative to the result.
func (v *formatter) formatAlternative(alternative AlternativeLike) {
	var factor Factor
	var factors = alternative.GetFactors()
	var iterator = col.Iterator(factors)
	if iterator.HasNext() {
		factor = iterator.GetNext()
		v.formatFactor(factor)
	}
	for iterator.HasNext() {
		v.appendString(" ")
		factor = iterator.GetNext()
		v.formatFactor(factor)
	}
	var note = alternative.GetNote()
	if len(note) > 0 {
		v.appendString("  ")
		v.formatNote(note)
	}
}

// This private method appends a formatted character to the result.
func (v *formatter) formatCharacter(character Character) {
	v.appendString(string(character))
}

// This private method appends a formatted comment to the result.
func (v *formatter) formatComment(comment Comment) {
	v.appendString(string(comment))
}

// This private method appends a formatted definition to the result.
func (v *formatter) formatDefinition(definition DefinitionLike) {
	var symbol = definition.GetSymbol()
	v.formatSymbol(symbol)
	v.appendString(":")
	v.depth++
	var expression = definition.GetExpression()
	if expression.IsMultilined() {
		v.appendNewline()
		v.appendString("  ")
	} else {
		v.appendString(" ")
	}
	v.formatExpression(expression)
	v.depth--
}

// This private method appends a formatted element to the result.
func (v *formatter) formatElement(element Element) {
	switch e := element.(type) {
	case Intrinsic:
		v.formatIntrinsic(e)
	case String:
		v.formatString(e)
	case Number:
		v.formatNumber(e)
	case Name:
		v.formatName(e)
	default:
		panic(fmt.Sprintf("Attempted to format:\n    element: %v\n    type: %t\n", e, element))
	}
}

// This private method appends a formatted exactly N group to the result.
func (v *formatter) formatExactlyN(group ExactlyNLike) {
	var expression = group.GetExpression()
	v.appendString("(")
	v.formatExpression(expression)
	v.appendString(")")
	var n = group.GetN()
	if len(n) > 0 {
		v.formatNumber(n)
	}
}

// This private method appends a formatted expression to the result.
func (v *formatter) formatExpression(expression ExpressionLike) {
	var alternatives = expression.GetAlternatives()
	var iterator = col.Iterator(alternatives)
	var alternative = iterator.GetNext()
	v.formatAlternative(alternative)
	for iterator.HasNext() {
		alternative = iterator.GetNext()
		if expression.IsMultilined() {
			v.appendNewline()
		} else {
			v.appendString(" ")
		}
		v.appendString("| ")
		v.formatAlternative(alternative)
	}
}

// This private method appends a formatted factor to the result.
func (v *formatter) formatFactor(factor Factor) {
	if ref.ValueOf(factor).Kind() == ref.String {
		v.formatElement(factor)
		return
	}
	switch f := factor.(type) {
	case *range_:
		v.formatRange(f)
	case *inverse:
		v.formatInverse(f)
	default:
		v.formatGrouping(Grouping(f))
	}
}

// This private method appends a formatted grammar to the result.
func (v *formatter) formatGrammar(grammar GrammarLike) {
	var statements = grammar.GetStatements()
	var iterator = col.Iterator(statements)
	for iterator.HasNext() {
		var statement = iterator.GetNext()
		v.formatStatement(statement)
	}
}

// This private method appends a formatted grouping to the result.
func (v *formatter) formatGrouping(grouping Grouping) {
	switch g := grouping.(type) {
	case *exactlyN:
		v.formatExactlyN(g)
	case *zeroOrOne:
		v.formatZeroOrOne(g)
	case *zeroOrMore:
		v.formatZeroOrMore(g)
	case *oneOrMore:
		v.formatOneOrMore(g)
	default:
		panic(fmt.Sprintf("Attempted to format:\n    grouping: %v\n    type: %t\n", g, grouping))
	}
}

// This private method appends a formatted intrinsic to the result.
func (v *formatter) formatIntrinsic(intrinsic Intrinsic) {
	v.appendString(string(intrinsic))
}

// This private method appends a formatted inverse to the result.
func (v *formatter) formatInverse(inverse InverseLike) {
	v.appendString("~")
	var factor = inverse.GetFactor()
	v.formatFactor(factor)
}

// This private method appends a formatted name to the result.
func (v *formatter) formatName(name Name) {
	v.appendString(string(name))
}

// This private method appends a formatted note to the result.
func (v *formatter) formatNote(note Note) {
	v.appendString(string(note))
}

// This private method appends a formatted number to the result.
func (v *formatter) formatNumber(number Number) {
	v.appendString(string(number))
}

// This private method appends a formatted one or more group to the result.
func (v *formatter) formatOneOrMore(group OneOrMoreLike) {
	var expression = group.GetExpression()
	v.appendString("<")
	v.formatExpression(expression)
	v.appendString(">")
}

// This private method appends a formatted range to the result.
func (v *formatter) formatRange(range_ RangeLike) {
	var first = range_.GetFirstCharacter()
	var last = range_.GetLastCharacter()
	v.formatCharacter(first)
	if len(last) > 0 {
		v.appendString("..")
		v.formatCharacter(last)
	}
}

// This private method appends a formatted statement to the result.
func (v *formatter) formatStatement(statement StatementLike) {
	var comment = statement.GetComment()
	if len(comment) > 0 {
		v.appendNewline()
		v.formatComment(comment)
	} else {
		var definition = statement.GetDefinition()
		v.formatDefinition(definition)
	}
	v.appendNewline()
	v.appendNewline()
}

// This private method appends a formatted string to the result.
func (v *formatter) formatString(string_ String) {
	v.appendString(string(string_))
}

// This private method appends a formatted symbol to the result.
func (v *formatter) formatSymbol(symbol Symbol) {
	v.appendString(string(symbol))
}

// This private method appends a formatted zero or more group to the result.
func (v *formatter) formatZeroOrMore(group ZeroOrMoreLike) {
	var expression = group.GetExpression()
	v.appendString("{")
	v.formatExpression(expression)
	v.appendString("}")
}

// This private method appends a formatted zero or one group to the result.
func (v *formatter) formatZeroOrOne(group ZeroOrOneLike) {
	var expression = group.GetExpression()
	v.appendString("[")
	v.formatExpression(expression)
	v.appendString("]")
}

// This method returns the canonically formatted string result.
func (v *formatter) getResult() string {
	var result = v.result.String()
	v.result.Reset()
	return result
}
