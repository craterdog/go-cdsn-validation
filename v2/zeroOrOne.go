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

// ZERO OR ONE GROUPING INTERFACE

// This interface defines the methods supported by all zero-or-one-like
// components.
type ZeroOrOneLike interface {
	GetExpression() ExpressionLike
	SetExpression(expression ExpressionLike)
}

// This constructor creates a new zero or more grouping.
func ZeroOrOne(expression ExpressionLike) ZeroOrOneLike {
	var v = &zeroOrOne{}
	v.SetExpression(expression)
	return v
}

// ZERO OR ONE GROUPING IMPLEMENTATION

// This type defines the structure and methods associated with an zero or more
// grouping.
type zeroOrOne struct {
	expression ExpressionLike
}

// This method returns the expression for this zero or more grouping.
func (v *zeroOrOne) GetExpression() ExpressionLike {
	return v.expression
}

// This method sets the expression for this zero or more grouping.
func (v *zeroOrOne) SetExpression(expression ExpressionLike) {
	if expression == nil {
		panic("An zero or more grouping requires an expression.")
	}
	v.expression = expression
}

// This private method appends a formatted zero or one group to the result.
func (v *formatter) formatZeroOrOne(group ZeroOrOneLike) {
	var expression = group.GetExpression()
	v.appendString("[")
	v.formatExpression(expression)
	v.appendString("]")
}
