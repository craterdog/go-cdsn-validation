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

// This private method appends a formatted annotation to the result.
func (v *formatter) formatAnnotation(annotation Annotation) {
	v.appendString(string(annotation))
}

// This private method appends a formatted definition to the result.
func (v *formatter) formatDefinition(definition DefinitionLike) {
	var alternatives = definition.GetAlternatives()
	var iterator = col.Iterator(alternatives)
	var alternative = iterator.GetNext()
	v.formatAlternative(alternative)
	for iterator.HasNext() {
		alternative = iterator.GetNext()
		if definition.IsMultilined() {
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
	switch f := factor.(type) {
	case InverseLike:
		v.formatInverse(f)
	case GroupLike:
		v.formatGroup(f)
	case RangeLike:
		v.formatRange(f)
	case Intrinsic:
		v.formatIntrinsic(f)
	case Rune:
		v.formatRune(f)
	case String:
		v.formatString(f)
	case Number:
		v.formatNumber(f)
	case Name:
		v.formatName(f)
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

// This private method appends a formatted group to the result.
func (v *formatter) formatGroup(group GroupLike) {
	var definition = group.GetDefinition()
	var type_ = group.GetType()
	var number = group.GetNumber()
	switch type_ {
	case ExactlyN:
		v.appendString("(")
		v.formatDefinition(definition)
		v.appendString(")")
	case ZeroOrOne:
		v.appendString("[")
		v.formatDefinition(definition)
		v.appendString("]")
	case ZeroOrMore:
		v.appendString("{")
		v.formatDefinition(definition)
		v.appendString("}")
	case OneOrMore:
		v.appendString("<")
		v.formatDefinition(definition)
		v.appendString(">")
	default:
		panic(fmt.Sprintf("Attempted to format an invalid group type: %v\n", type_))
	}
	v.formatNumber(number)
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

// This private method appends a formatted production to the result.
func (v *formatter) formatProduction(production ProductionLike) {
	var symbol = production.GetSymbol()
	v.formatSymbol(symbol)
	v.appendString(":")
	v.depth++
	var definition = production.GetDefinition()
	if definition.IsMultilined() {
		v.appendNewline()
		v.appendString("  ")
	} else {
		v.appendString(" ")
	}
	v.formatDefinition(definition)
	v.depth--
}

// This private method appends a formatted range to the result.
func (v *formatter) formatRange(range_ RangeLike) {
	var first = range_.GetFirstRune()
	v.formatRune(first)
	v.appendString("..")
	var last = range_.GetLastRune()
	v.formatRune(last)
}

// This private method appends a formatted rune to the result.
func (v *formatter) formatRune(rune_ Rune) {
	v.appendString(string(rune_))
}

// This private method appends a formatted statement to the result.
func (v *formatter) formatStatement(statement StatementLike) {
	var annotation = statement.GetAnnotation()
	if len(annotation) > 0 {
		v.appendNewline()
		v.formatAnnotation(annotation)
	} else {
		var production = statement.GetProduction()
		v.formatProduction(production)
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
