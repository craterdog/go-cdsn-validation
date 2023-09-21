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
)

// COMPILER INTERFACE

// This function compiles the specified grammar into its corresponding parser.
func CompileGrammar(grammar GrammarLike) {
	var v = &compiler{}
	v.compileGrammar(grammar)
	// add #package#.go file if one does not yet exist
	// generate scanner.go file
	// generate parser.go file
}

// COMPILER IMPLEMENTATION

// This type defines the structure and methods for a compiler agent.
type compiler struct {
	packageName string
}

// This private method compiles an alternative.
func (v *compiler) compileAlternative(alternative AlternativeLike) {
	var factor Factor
	var factors = alternative.GetFactors()
	var iterator = col.Iterator(factors)
	for iterator.HasNext() {
		factor = iterator.GetNext()
		v.compileFactor(factor)
	}
}

// This private method compiles a definition.
func (v *compiler) compileDefinition(definition DefinitionLike) {
	var symbol = definition.GetSymbol()
	v.compileSymbol(symbol)
	var expression = definition.GetExpression()
	v.compileExpression(expression)
}

// This private method compiles an expression.
func (v *compiler) compileExpression(expression ExpressionLike) {
	var alternatives = expression.GetAlternatives()
	var iterator = col.Iterator(alternatives)
	for iterator.HasNext() {
		var alternative = iterator.GetNext()
		v.compileAlternative(alternative)
	}
}

// This private method compiles an exactly N expression.
func (v *compiler) compileExactlyN(expression ExpressionLike, number Number) {
}

// This private method compiles a zero or one expression.
func (v *compiler) compileZeroOrOne(expression ExpressionLike) {
}

// This private method compiles a zero or more expression.
func (v *compiler) compileZeroOrMore(expression ExpressionLike) {
}

// This private method compiles a one or more expression.
func (v *compiler) compileOneOrMore(expression ExpressionLike) {
}

// This private method compiles a factor.
func (v *compiler) compileFactor(factor Factor) {
	switch f := factor.(type) {
	case Intrinsic:
		v.compileIntrinsic(f)
	case String:
		v.compileString(f)
	case Number:
		v.compileNumber(f)
	case Name:
		v.compileName(f)
	case RangeLike:
		v.compileRange(f)
	case InverseLike:
		v.compileInverse(f)
	case GroupLike:
		v.compileGroup(f)
	default:
		panic(fmt.Sprintf("Attempted to compile:\n    factor: %v\n    type: %t\n", f, factor))
	}
}

// This private method compiles a grammar.
func (v *compiler) compileGrammar(grammar GrammarLike) {
	var statements = grammar.GetStatements()
	var iterator = col.Iterator(statements)
	for iterator.HasNext() {
		var statement = iterator.GetNext()
		v.compileStatement(statement)
	}
}

// This private method compiles a group.
func (v *compiler) compileGroup(group GroupLike) {
	var expression = group.GetExpression()
	var type_ = group.GetType()
	var number = group.GetNumber()
	switch type_ {
	case ExactlyN:
		v.compileExactlyN(expression, number)
	case ZeroOrOne:
		v.compileZeroOrOne(expression)
	case ZeroOrMore:
		v.compileZeroOrMore(expression)
	case OneOrMore:
		v.compileOneOrMore(expression)
	default:
		panic(fmt.Sprintf("Attempted to compile an invalid group type: %v\n", type_))
	}
}

// This private method compiles an intrinsic.
func (v *compiler) compileIntrinsic(intrinsic Intrinsic) {
}

// This private method compiles an inverse.
func (v *compiler) compileInverse(inverse InverseLike) {
	var factor = inverse.GetFactor()
	v.compileFactor(factor)
}

// This private method compiles a name.
func (v *compiler) compileName(name Name) {
}

// This private method compiles a number.
func (v *compiler) compileNumber(number Number) {
}

// This private method compiles a range.
func (v *compiler) compileRange(range_ RangeLike) {
	//var first = range_.GetFirstCharacter()
	//var last = range_.GetLastCharacter()
}

// This private method compiles a statement.
func (v *compiler) compileStatement(statement StatementLike) {
	var definition = statement.GetDefinition()
	if definition != nil {
		v.compileDefinition(definition)
	}
}

// This private method compiles a string.
func (v *compiler) compileString(string_ String) {
}

// This private method compiles a symbol.
func (v *compiler) compileSymbol(symbol Symbol) {
	if len(v.packageName) < 1 {
		// The symbol for the first definition defines the package name.
		v.packageName = string(symbol)[1:]
	}
}
