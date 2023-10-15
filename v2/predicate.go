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
	GetInversion() InversionLike
	SetInversion(inversion InversionLike)
	GetRepetition() RepetitionLike
	SetRepetition(repetition RepetitionLike)
	GetFactor() FactorLike
	SetFactor(factor FactorLike)
}

// This constructor creates a new predicate.
func Predicate(inversion InversionLike, repetition RepetitionLike, factor FactorLike) PredicateLike {
	if inversion == nil && repetition == nil && factor == nil {
		panic("A predicate requires at least one of its attributes to be set.")
	}
	var v = &predicate{}
	v.SetInversion(inversion)
	v.SetRepetition(repetition)
	v.SetFactor(factor)
	return v
}

// PREDICATE IMPLEMENTATION

// This type defines the structure and methods associated with a predicate.
type predicate struct {
	inversion  InversionLike
	repetition RepetitionLike
	factor     FactorLike
}

// This method returns the inversion for this predicate.
func (v *predicate) GetInversion() InversionLike {
	return v.inversion
}

// This method sets the inversion for this predicate.
func (v *predicate) SetInversion(inversion InversionLike) {
	if inversion != nil {
		v.inversion = inversion
		v.repetition = nil
		v.factor = nil
	}
}

// This method returns the repetition for this predicate.
func (v *predicate) GetRepetition() RepetitionLike {
	return v.repetition
}

// This method sets the repetition for this predicate.
func (v *predicate) SetRepetition(repetition RepetitionLike) {
	if repetition != nil {
		v.inversion = nil
		v.repetition = repetition
		v.factor = nil
	}
}

// This method returns the factor for this predicate.
func (v *predicate) GetFactor() FactorLike {
	return v.factor
}

// This method sets the factor for this predicate.
func (v *predicate) SetFactor(factor FactorLike) {
	if factor != nil {
		v.inversion = nil
		v.repetition = nil
		v.factor = factor
	}
}
