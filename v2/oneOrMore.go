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

// ONE OR MORE GROUPING INTERFACE

// This interface defines the methods supported by all one-or-more-like
// components.
type OneOrMoreLike interface {
	GetExpression() ExpressionLike
	SetExpression(expression ExpressionLike)
}

// This constructor creates a new one or more grouping.
func OneOrMore(expression ExpressionLike) OneOrMoreLike {
	var v = &oneOrMore{}
	v.SetExpression(expression)
	return v
}

// ONE OR MORE GROUPING IMPLEMENTATION

// This type defines the structure and methods associated with an one or more
// grouping.
type oneOrMore struct {
	expression ExpressionLike
}

// This method returns the expression for this one or more grouping.
func (v *oneOrMore) GetExpression() ExpressionLike {
	return v.expression
}

// This method sets the expression for this one or more grouping.
func (v *oneOrMore) SetExpression(expression ExpressionLike) {
	if expression == nil {
		panic("An one or more grouping requires an expression.")
	}
	v.expression = expression
}

// This method attempts to parse an one or more grouping. It returns the one or
// more grouping and whether or not the one or more grouping was successfully parsed.
func (v *parser) parseOneOrMore() (OneOrMoreLike, *Token, bool) {
	var ok bool
	var token *Token
	var expression ExpressionLike
	var oneOrMore OneOrMoreLike
	_, token, ok = v.parseDelimiter("<")
	if !ok {
		// This is not an one or more grouping.
		return oneOrMore, token, false
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
	_, token, ok = v.parseDelimiter(">")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(">",
			"$factor",
			"$expression")
		panic(message)
	}
	oneOrMore = OneOrMore(expression)
	return oneOrMore, token, true
}

// This private method appends a formatted one or more group to the result.
func (v *formatter) formatOneOrMore(group OneOrMoreLike) {
	var expression = group.GetExpression()
	v.appendString("<")
	v.formatExpression(expression)
	v.appendString(">")
}
