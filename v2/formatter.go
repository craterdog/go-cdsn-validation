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

// This method returns the canonically formatted string result.
func (v *formatter) getResult() string {
	var result = v.result.String()
	v.result.Reset()
	return result
}

// This method appends the specified string to the result.
func (v *formatter) appendString(s string) {
	v.result.WriteString(s)
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

// This private method appends a formatted count to the result.
func (v *formatter) formatCount(count CountLike) {
	var digits = count.GetDigits()
	var iterator = col.Iterator(digits)
	for iterator.HasNext() {
		var digit = iterator.GetNext()
		v.appendString(string(digit))
	}
}

// This private method appends a formatted factor to the result.
func (v *formatter) formatFactor(factor Factor) {
	switch f := factor.(type) {
	case Character:
		v.formatCharacter(f)
	case Literal:
		v.formatLiteral(f)
	case Intrinsic:
		v.formatIntrinsic(f)
	case Identifier:
		v.formatIdentifier(f)
	case InversionLike:
		v.formatInversion(f)
	case GroupingLike:
		v.formatGrouping(f)
	case RangeLike:
		v.formatRange(f)
	default:
		panic(fmt.Sprintf("Attempted to format:\n    factor: %v\n    type: %t\n", f, factor))
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
func (v *formatter) formatGrouping(grouping GroupingLike) {
	var rule = grouping.GetRule()
	var type_ = grouping.GetType()
	var count = grouping.GetCount()
	switch type_ {
	case Optional:
		v.appendString("[")
		v.formatRule(rule)
		v.appendString("]")
	case ExactCount:
		v.appendString("(")
		v.formatRule(rule)
		v.appendString(")")
	case MinimumCount:
		v.appendString("<")
		v.formatRule(rule)
		v.appendString(">")
	case MaximumCount:
		v.appendString("{")
		v.formatRule(rule)
		v.appendString("}")
	default:
		panic(fmt.Sprintf("Attempted to format an invalid grouping type: %v\n", type_))
	}
	v.formatCount(count)
}

// This private method appends a formatted identifier to the result.
func (v *formatter) formatIdentifier(identifier Identifier) {
	v.appendString(string(identifier))
}

// This private method appends a formatted intrinsic to the result.
func (v *formatter) formatIntrinsic(intrinsic Intrinsic) {
	v.appendString(string(intrinsic))
}

// This private method appends a formatted inversion to the result.
func (v *formatter) formatInversion(inversion InversionLike) {
	v.appendString("~")
	var factor = inversion.GetFactor()
	v.formatFactor(factor)
}

// This private method appends a formatted literal to the result.
func (v *formatter) formatLiteral(literal Literal) {
	v.appendString(string(literal))
}

// This private method appends a formatted note to the result.
func (v *formatter) formatNote(note Note) {
	v.appendString(string(note))
}

// This private method appends a formatted production to the result.
func (v *formatter) formatProduction(production ProductionLike) {
	var symbol = production.GetSymbol()
	v.formatSymbol(symbol)
	v.appendString(":")
	v.depth++
	var rule = production.GetRule()
	if rule.IsMultilined() {
		v.appendNewline()
		v.appendString("  ")
	} else {
		v.appendString(" ")
	}
	v.formatRule(rule)
	v.depth--
}

// This private method appends a formatted range to the result.
func (v *formatter) formatRange(range_ RangeLike) {
	var first = range_.GetFirstCharacter()
	v.formatCharacter(first)
	v.appendString("..")
	var last = range_.GetLastCharacter()
	v.formatCharacter(last)
}

// This private method appends a formatted rule to the result.
func (v *formatter) formatRule(rule RuleLike) {
	var alternatives = rule.GetAlternatives()
	var iterator = col.Iterator(alternatives)
	var alternative = iterator.GetNext()
	v.formatAlternative(alternative)
	for iterator.HasNext() {
		alternative = iterator.GetNext()
		if rule.IsMultilined() {
			v.appendNewline()
		} else {
			v.appendString(" ")
		}
		v.appendString("| ")
		v.formatAlternative(alternative)
	}
}

// This private method appends a formatted statement to the result.
func (v *formatter) formatStatement(statement StatementLike) {
	var comment = statement.GetComment()
	if len(comment) > 0 {
		v.formatComment(comment)
	} else {
		var production = statement.GetProduction()
		v.formatProduction(production)
	}
	v.appendNewline()
	v.appendNewline()
}

// This private method appends a formatted symbol to the result.
func (v *formatter) formatSymbol(symbol Symbol) {
	v.appendString(string(symbol))
}
