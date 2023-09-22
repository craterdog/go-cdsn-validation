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
	GetN() Number
	SetN(n Number)
}

// This constructor creates a new exactly N grouping.
func ExactlyN(expression ExpressionLike, n Number) ExactlyNLike {
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
	n          Number
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
func (v *exactlyN) GetN() Number {
	return v.n
}

// This method sets the number for this exactly N grouping.
func (v *exactlyN) SetN(n Number) {
	v.n = n
}

// This method attempts to parse an exactly N grouping. It returns the exactly
// N grouping and whether or not the exactly N grouping was successfully parsed.
func (v *parser) parseExactlyN() (ExactlyNLike, *Token, bool) {
	var ok bool
	var token *Token
	var expression ExpressionLike
	var exactlyN ExactlyNLike
	_, token, ok = v.parseLiteral("(")
	if !ok {
		// This is not an exactly N grouping.
		return exactlyN, token, false
	}
	expression, token, ok = v.parseExpression()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("expression",
			"$factor",
			"$expression")
		panic(message)
	}
	expression.SetMultilined(false)
	_, token, ok = v.parseLiteral(")")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(")",
			"$factor",
			"$expression")
		panic(message)
	}
	var n, _, _ = v.parseNumber() // The number is optional.
	exactlyN = ExactlyN(expression, n)
	return exactlyN, token, true
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
