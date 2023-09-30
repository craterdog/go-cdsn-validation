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

// DEFINITION INTERFACE

// This interface defines the methods supported by all definition-like
// components.
type DefinitionLike interface {
	GetSYMBOL() SYMBOL
	SetSYMBOL(symbol SYMBOL)
	GetExpression() ExpressionLike
	SetExpression(expression ExpressionLike)
}

// This constructor creates a new definition.
func Definition(symbol SYMBOL, expression ExpressionLike) DefinitionLike {
	var v = &definition{}
	v.SetSYMBOL(symbol)
	v.SetExpression(expression)
	return v
}

// DEFINITION IMPLEMENTATION

// This type defines the structure and methods associated with a definition.
type definition struct {
	symbol     SYMBOL
	expression ExpressionLike
}

// This method returns the symbol for this definition.
func (v *definition) GetSYMBOL() SYMBOL {
	return v.symbol
}

// This method sets the symbol for this definition.
func (v *definition) SetSYMBOL(symbol SYMBOL) {
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
