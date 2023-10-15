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
	GetInversion() InversionLike
	SetInversion(inversion InversionLike)
	GetPrecedence() PrecedenceLike
	SetPrecedence(precedence PrecedenceLike)
}

// This constructor creates a new predicate.
func Predicate(factor FactorLike, inversion InversionLike, precedence PrecedenceLike) PredicateLike {
	if factor == nil && inversion == nil && precedence == nil {
		panic("A predicate requires at least one of its attributes to be set.")
	}
	var v = &predicate{}
	v.SetFactor(factor)
	v.SetInversion(inversion)
	v.SetPrecedence(precedence)
	return v
}

// PREDICATE IMPLEMENTATION

// This type defines the structure and methods associated with a predicate.
type predicate struct {
	factor     FactorLike
	inversion  InversionLike
	precedence PrecedenceLike
}

// This method returns the factor for this predicate.
func (v *predicate) GetFactor() FactorLike {
	return v.factor
}

// This method sets the factor for this predicate.
func (v *predicate) SetFactor(factor FactorLike) {
	if factor != nil {
		v.factor = factor
		v.inversion = nil
		v.precedence = nil
	}
}

// This method returns the inversion for this predicate.
func (v *predicate) GetInversion() InversionLike {
	return v.inversion
}

// This method sets the inversion for this predicate.
func (v *predicate) SetInversion(inversion InversionLike) {
	if inversion != nil {
		v.factor = nil
		v.inversion = inversion
		v.precedence = nil
	}
}

// This method returns the precedence for this predicate.
func (v *predicate) GetPrecedence() PrecedenceLike {
	return v.precedence
}

// This method sets the precedence for this predicate.
func (v *predicate) SetPrecedence(precedence PrecedenceLike) {
	if precedence != nil {
		v.factor = nil
		v.inversion = nil
		v.precedence = precedence
	}
}
