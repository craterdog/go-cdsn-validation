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

// PRECEDENCE INTERFACE

// This interface defines the methods supported by all precedence-like
// components.
type PrecedenceLike interface {
	GetExpression() ExpressionLike
	SetExpression(expression ExpressionLike)
}

// This constructor creates a new precedence.
func Precedence(expression ExpressionLike) PrecedenceLike {
	var v = &precedence{}
	v.SetExpression(expression)
	return v
}

// PRECEDENCE IMPLEMENTATION

// This type defines the structure and methods associated with a precedence.
type precedence struct {
	expression ExpressionLike
}

// This method returns the expression for this precedence.
func (v *precedence) GetExpression() ExpressionLike {
	return v.expression
}

// This method sets the expression for this precedence.
func (v *precedence) SetExpression(expression ExpressionLike) {
	if expression == nil {
		panic("A precedence requires an expression.")
	}
	v.expression = expression
}
