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

// EXACTLY N GROUPING INTERFACE

// This interface defines the methods supported by all exactly-n-like
// components.
type ExactlyNLike interface {
	GetExpression() ExpressionLike
	SetExpression(expression ExpressionLike)
	GetN() NUMBER
	SetN(n NUMBER)
}

// This constructor creates a new exactly N grouping.
func ExactlyN(expression ExpressionLike, n NUMBER) ExactlyNLike {
	var v = &exactlyN{}
	v.SetExpression(expression)
	v.SetN(n)
	return v
}

// EXACTLY N GROUPING IMPLEMENTATION

// This type defines the structure and methods associated with an exactly N
// grouping.
type exactlyN struct {
	expression ExpressionLike
	n          NUMBER
}

// This method returns the expression for this exactly N grouping.
func (v *exactlyN) GetExpression() ExpressionLike {
	return v.expression
}

// This method sets the expression for this exactly N grouping.
func (v *exactlyN) SetExpression(expression ExpressionLike) {
	if expression == nil {
		panic("An exactly N grouping requires an expression.")
	}
	v.expression = expression
}

// This method returns the number for this exactly N grouping.
func (v *exactlyN) GetN() NUMBER {
	return v.n
}

// This method sets the number for this exactly N grouping.
func (v *exactlyN) SetN(n NUMBER) {
	v.n = n
}
