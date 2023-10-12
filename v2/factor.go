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
	SetPrecedence(precedence_ PrecedenceLike)
	GetElement() ElementLike
	SetElement(element_ ElementLike)
}

// This constructor creates a new factor.
func Factor(precedence_ PrecedenceLike, element_ ElementLike) FactorLike {
	var v = &factor_{}
	v.SetPrecedence(precedence_)
	v.SetElement(element_)
	return v
}

// FACTOR IMPLEMENTATION

// This type defines the structure and methods associated with a factor.
type factor_ struct {
	precedence_ PrecedenceLike
	element_    ElementLike
}

// This method returns the precedence for this factor.
func (v *factor_) GetPrecedence() PrecedenceLike {
	return v.precedence_
}

// This method sets the precedence for this factor.
func (v *factor_) SetPrecedence(precedence_ PrecedenceLike) {
	if precedence_ != nil {
		v.precedence_ = precedence_
		v.element_ = nil
	}
}

// This method returns the element for this factor.
func (v *factor_) GetElement() ElementLike {
	return v.element_
}

// This method sets the element for this factor.
func (v *factor_) SetElement(element_ ElementLike) {
	if element_ != nil {
		v.precedence_ = nil
		v.element_ = element_
	}
}
