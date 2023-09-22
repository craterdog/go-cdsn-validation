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
	uni "unicode"
)

// DEFINITION INTERFACE

// This interface defines the methods supported by all definition-like
// components.
type DefinitionLike interface {
	GetSymbol() Symbol
	SetSymbol(symbol Symbol)
	GetExpression() ExpressionLike
	SetExpression(expression ExpressionLike)
}

// This constructor creates a new definition.
func Definition(symbol Symbol, expression ExpressionLike) DefinitionLike {
	var v = &definition{}
	v.SetSymbol(symbol)
	v.SetExpression(expression)
	return v
}

// DEFINITION IMPLEMENTATION

// This type defines the structure and methods associated with a definition.
type definition struct {
	symbol     Symbol
	expression ExpressionLike
}

// This method returns the symbol for this definition.
func (v *definition) GetSymbol() Symbol {
	return v.symbol
}

// This method sets the symbol for this definition.
func (v *definition) SetSymbol(symbol Symbol) {
	if len(symbol) == 0 {
		panic("A definition requires a symbol.")
	}
	v.symbol = symbol
}

// This method returns the expression for this definition.
func (v *definition) GetExpression() ExpressionLike {
	return v.expression
}

// This method sets the expression for this definition.
func (v *definition) SetExpression(expression ExpressionLike) {
	if expression == nil {
		panic("A definition requires an expression.")
	}
	v.expression = expression
}

// This method attempts to parse a definition. It returns the definition and
// whether or not the definition was successfully parsed.
func (v *parser) parseDefinition() (DefinitionLike, *Token, bool) {
	var ok bool
	var token *Token
	var symbol Symbol
	var expression ExpressionLike
	var definition DefinitionLike
	symbol, token, ok = v.parseSymbol()
	if !ok {
		// This is not a definition.
		return definition, token, false
	}
	v.isToken = uni.IsUpper(rune(symbol[1]))
	_, token, ok = v.parseLiteral(":")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(":",
			"$definition",
			"$SYMBOL",
			"$expression")
		panic(message)
	}
	expression, token, ok = v.parseExpression()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("expression",
			"$definition",
			"$SYMBOL",
			"$expression")
		panic(message)
	}
	definition = Definition(symbol, expression)
	v.symbols.SetValue(symbol, definition)
	return definition, token, true
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
