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

// REPETITION INTERFACE

// This interface defines the methods supported by all repetition-like
// components.
type RepetitionLike interface {
	GetCONSTRAINT() CONSTRAINT
	SetCONSTRAINT(constraint CONSTRAINT)
	GetFactor() FactorLike
	SetFactor(factor FactorLike)
}

// This constructor creates a new repetition.
func Repetition(constraint CONSTRAINT, factor FactorLike) RepetitionLike {
	var v = &repetition{}
	v.SetCONSTRAINT(constraint)
	v.SetFactor(factor)
	return v
}

// REPETITION IMPLEMENTATION

// This type defines the structure and methods associated with a repetition.
type repetition struct {
	constraint CONSTRAINT
	factor     FactorLike
}

// This method returns the number for this repetition.
func (v *repetition) GetCONSTRAINT() CONSTRAINT {
	return v.constraint
}

// This method sets the number for this repetition.
func (v *repetition) SetCONSTRAINT(constraint CONSTRAINT) {
	if len(constraint) == 0 {
		panic("A repetition requires a constraint.")
	}
	v.constraint = constraint
}

// This method returns the factor for this repetition.
func (v *repetition) GetFactor() FactorLike {
	return v.factor
}

// This method sets the factor for this repetition.
func (v *repetition) SetFactor(factor FactorLike) {
	if factor == nil {
		panic("A repetition requires a factor.")
	}
	v.factor = factor
}
