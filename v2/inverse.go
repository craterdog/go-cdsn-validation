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

// This private method appends a formatted inverse to the result.
func (v *formatter) formatInverse(inverse InverseLike) {
	v.appendString("~")
	var factor = inverse.GetFactor()
	v.formatFactor(factor)
}
