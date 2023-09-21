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

// GROUPING IMPLEMENTATION

// This constructor creates a new group.
func Group(expression ExpressionLike, type_ GroupType, number Number) GroupLike {
	var v = &group{}
	v.SetExpression(expression)
	v.SetType(type_)
	v.SetNumber(number)
	return v
}

// This type defines the structure and methods associated with a group.
type group struct {
	expression ExpressionLike
	type_      GroupType
	number     Number
}

// This method returns the expression for this group.
func (v *group) GetExpression() ExpressionLike {
	return v.expression
}

// This method sets the expression for this group.
func (v *group) SetExpression(expression ExpressionLike) {
	if expression == nil {
		panic("A group requires an expression.")
	}
	v.expression = expression
}

// This method returns the group type for this group.
func (v *group) GetType() GroupType {
	return v.type_
}

// This method sets the group type for this group.
func (v *group) SetType(type_ GroupType) {
	v.type_ = type_
}

// This method returns the number for this group.
func (v *group) GetNumber() Number {
	return v.number
}

// This method sets the number for this group.
func (v *group) SetNumber(number Number) {
	v.number = number
}

// This method attempts to parse a zero or more group. It returns the zero or
// more group and whether or not the zero or more group was successfully parsed.
func (v *parser) parseZeroOrMore() (GroupLike, *Token, bool) {
	var ok bool
	var token *Token
	var expression ExpressionLike
	var group GroupLike
	_, token, ok = v.parseDelimiter("{")
	if !ok {
		// This is not a zero or more group.
		return group, token, false
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
	_, token, ok = v.parseDelimiter("}")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("}",
			"$factor",
			"$expression")
		panic(message)
	}
	var number, _, _ = v.parseNumber() // The number is optional.
	group = Group(expression, ZeroOrMore, number)
	return group, token, true
}

// This method attempts to parse a zero or one group. It returns the zero or
// one group and whether or not the zero or one group was successfully parsed.
func (v *parser) parseZeroOrOne() (GroupLike, *Token, bool) {
	var ok bool
	var token *Token
	var expression ExpressionLike
	var group GroupLike
	_, token, ok = v.parseDelimiter("[")
	if !ok {
		// This is not a zero or one group.
		return group, token, false
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
	_, token, ok = v.parseDelimiter("]")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("]",
			"$factor",
			"$expression")
		panic(message)
	}
	var number, _, _ = v.parseNumber() // The number is optional.
	group = Group(expression, ZeroOrOne, number)
	return group, token, true
}

// This method attempts to parse an exact count group. It returns the
// exact count group and whether or not the exact count group was
// successfully parsed.
func (v *parser) parseExactlyN() (GroupLike, *Token, bool) {
	var ok bool
	var token *Token
	var expression ExpressionLike
	var group GroupLike
	_, token, ok = v.parseDelimiter("(")
	if !ok {
		// This is not a precedence group.
		return group, token, false
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
	_, token, ok = v.parseDelimiter(")")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar(")",
			"$factor",
			"$expression")
		panic(message)
	}
	var number, _, _ = v.parseNumber() // The number is optional.
	group = Group(expression, ExactlyN, number)
	return group, token, true
}

// This method attempts to parse a one-or-more group. It returns the one or
// more group and whether or not the one or more group was successfully
// parsed.
func (v *parser) parseOneOrMore() (GroupLike, *Token, bool) {
	var ok bool
	var token *Token
	var expression ExpressionLike
	var group GroupLike
	_, token, ok = v.parseDelimiter("<")
	if !ok {
		// This is not a one-or-more group.
		return group, token, false
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
	var number, _, _ = v.parseNumber() // The number is optional.
	group = Group(expression, OneOrMore, number)
	return group, token, true
}
