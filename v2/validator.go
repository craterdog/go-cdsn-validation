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
	uni "unicode"
)

// VALIDATOR INTERFACE

// This function validates the specified grammar and generates an exception if
// the grammar is not valid.
func ValidateGrammar(grammar GrammarLike) {
	var v = &validator{}
	v.validateGrammar(grammar)
}

// VALIDATOR IMPLEMENTATION

// This type defines the structure and methods for a validator agent.
type validator struct {
	isToken bool
}

// This private method validates an alternative.
func (v *validator) validateAlternative(alternative AlternativeLike) {
	var factor Factor
	var factors = alternative.GetFactors()
	var iterator = col.Iterator(factors)
	for iterator.HasNext() {
		factor = iterator.GetNext()
		v.validateFactor(factor)
	}
}

// This private method validates a definition.
func (v *validator) validateDefinition(definition DefinitionLike) {
	var symbol = definition.GetSymbol()
	v.isToken = uni.IsUpper(rune(symbol[1]))
	var expression = definition.GetExpression()
	v.validateExpression(expression)
}

// This private method validates an expression.
func (v *validator) validateExpression(expression ExpressionLike) {
	var alternatives = expression.GetAlternatives()
	var iterator = col.Iterator(alternatives)
	for iterator.HasNext() {
		var alternative = iterator.GetNext()
		v.validateAlternative(alternative)
	}
}

// This private method validates a factor.
func (v *validator) validateFactor(factor Factor) {
	switch f := factor.(type) {
	case Name:
		v.validateName(f)
	case InverseLike:
		v.validateInverse(f)
	case GroupLike:
		v.validateGroup(f)
	}
}

// This private method validates a grammar.
func (v *validator) validateGrammar(grammar GrammarLike) {
	var statements = grammar.GetStatements()
	var iterator = col.Iterator(statements)
	for iterator.HasNext() {
		var statement = iterator.GetNext()
		v.validateStatement(statement)
	}
}

// This private method validates a group.
func (v *validator) validateGroup(group GroupLike) {
	var expression = group.GetExpression()
	v.validateExpression(expression)
}

// This private method validates an inverse.
func (v *validator) validateInverse(inverse InverseLike) {
	var factor = inverse.GetFactor()
	v.validateFactor(factor)
}

// This private method validates a name.
func (v *validator) validateName(name Name) {
	if v.isToken && uni.IsLower(rune(name[0])) {
		panic(fmt.Sprintf("A token definition contains a rulename: %v\n", name))
	}
}

// This private method validates a statement.
func (v *validator) validateStatement(statement StatementLike) {
	var definition = statement.GetDefinition()
	if definition != nil {
		v.validateDefinition(definition)
	}
}
