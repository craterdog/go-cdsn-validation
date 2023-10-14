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

// FACTOR INTERFACE

// This interface defines the methods supported by all factor-like
// components.
type FactorLike interface {
	GetPrecedence() PrecedenceLike
	SetPrecedence(precedence PrecedenceLike)
	GetElement() ElementLike
	SetElement(element ElementLike)
}

// This constructor creates a new factor.
func Factor(precedence PrecedenceLike, element ElementLike) FactorLike {
	var v = &factor{}
	v.SetPrecedence(precedence)
	v.SetElement(element)
	return v
}

// FACTOR IMPLEMENTATION

// This type defines the structure and methods associated with a factor.
type factor struct {
	precedence PrecedenceLike
	element    ElementLike
}

// This method returns the precedence for this factor.
func (v *factor) GetPrecedence() PrecedenceLike {
	return v.precedence
}

// This method sets the precedence for this factor.
func (v *factor) SetPrecedence(precedence PrecedenceLike) {
	if precedence != nil {
		v.precedence = precedence
		v.element = nil
	}
}

// This method returns the element for this factor.
func (v *factor) GetElement() ElementLike {
	return v.element
}

// This method sets the element for this factor.
func (v *factor) SetElement(element ElementLike) {
	if element != nil {
		v.precedence = nil
		v.element = element
	}
}
