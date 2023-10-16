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

// PREDICATE INTERFACE

// This interface defines the methods supported by all predicate-like
// components.
type PredicateLike interface {
	GetFactor() FactorLike
	SetFactor(factor FactorLike)
	GetRepetition() RepetitionLike
	SetRepetition(repetition RepetitionLike)
}

// This constructor creates a new predicate.
func Predicate(factor FactorLike, repetition RepetitionLike) PredicateLike {
	if factor == nil {
		panic("A predicate requires a factor to be set.")
	}
	var v = &predicate{}
	v.SetFactor(factor)
	v.SetRepetition(repetition)
	return v
}

// PREDICATE IMPLEMENTATION

// This type defines the structure and methods associated with a predicate.
type predicate struct {
	factor     FactorLike
	repetition RepetitionLike
}

// This method returns the factor for this predicate.
func (v *predicate) GetFactor() FactorLike {
	return v.factor
}

// This method sets the factor for this predicate.
func (v *predicate) SetFactor(factor FactorLike) {
	if factor != nil {
		v.factor = factor
	}
}

// This method returns the repetition for this predicate.
func (v *predicate) GetRepetition() RepetitionLike {
	return v.repetition
}

// This method sets the repetition for this predicate.
func (v *predicate) SetRepetition(repetition RepetitionLike) {
	v.repetition = repetition
}
