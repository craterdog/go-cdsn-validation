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

// INVERSION INTERFACE

// This interface defines the methods supported by all inverse-like
// components.
type InverseLike interface {
	GetFactor() Factor
	SetFactor(factor Factor)
}

// This constructor creates a new inverse.
func Inverse(factor Factor) InverseLike {
	var v = &inverse{}
	v.SetFactor(factor)
	return v
}

// INVERSION IMPLEMENTATION

// This type defines the structure and methods associated with an inverse.
type inverse struct {
	factor Factor
}

// This method returns the factor for this inverse.
func (v *inverse) GetFactor() Factor {
	return v.factor
}

// This method sets the factor for this inverse.
func (v *inverse) SetFactor(factor Factor) {
	if factor == nil {
		panic("An inverse requires a factor.")
	}
	v.factor = factor
}

// This method attempts to parse an inverse. It returns the inverse and
// whether or not the inverse was successfully parsed.
func (v *parser) parseInverse() (InverseLike, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	var inverse InverseLike
	_, token, ok = v.parseDelimiter("~")
	if !ok {
		// This is not an inverse.
		return inverse, token, false
	}
	factor, token, ok = v.parseFactor()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("factor",
			"$inverse",
			"$factor")
		panic(message)
	}
	inverse = Inverse(factor)
	return inverse, token, true
}

// This private method appends a formatted inverse to the result.
func (v *formatter) formatInverse(inverse InverseLike) {
	v.appendString("~")
	var factor = inverse.GetFactor()
	v.formatFactor(factor)
}
