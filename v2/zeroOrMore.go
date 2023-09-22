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

// ZERO OR MORE GROUPING INTERFACE

// This interface defines the methods supported by all zero-or-more-like
// components.
type ZeroOrMoreLike interface {
	GetExpression() ExpressionLike
	SetExpression(expression ExpressionLike)
}

// This constructor creates a new zero or more grouping.
func ZeroOrMore(expression ExpressionLike) ZeroOrMoreLike {
	var v = &zeroOrMore{}
	v.SetExpression(expression)
	return v
}

// ZERO OR MORE GROUPING IMPLEMENTATION

// This type defines the structure and methods associated with an zero or more
// grouping.
type zeroOrMore struct {
	expression ExpressionLike
}

// This method returns the expression for this zero or more grouping.
func (v *zeroOrMore) GetExpression() ExpressionLike {
	return v.expression
}

// This method sets the expression for this zero or more grouping.
func (v *zeroOrMore) SetExpression(expression ExpressionLike) {
	if expression == nil {
		panic("An zero or more grouping requires an expression.")
	}
	v.expression = expression
}

// This method attempts to parse an zero or more grouping. It returns the zero or
// more grouping and whether or not the zero or more grouping was successfully parsed.
func (v *parser) parseZeroOrMore() (ZeroOrMoreLike, *Token, bool) {
	var ok bool
	var token *Token
	var expression ExpressionLike
	var zeroOrMore ZeroOrMoreLike
	_, token, ok = v.parseLiteral("{")
	if !ok {
		// This is not an zero or more grouping.
		return zeroOrMore, token, false
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
	_, token, ok = v.parseLiteral("}")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("}",
			"$factor",
			"$expression")
		panic(message)
	}
	zeroOrMore = ZeroOrMore(expression)
	return zeroOrMore, token, true
}

// This private method appends a formatted zero or more group to the result.
func (v *formatter) formatZeroOrMore(group ZeroOrMoreLike) {
	var expression = group.GetExpression()
	v.appendString("{")
	v.formatExpression(expression)
	v.appendString("}")
}
