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
	GetGlyph() GlyphLike
	SetGlyph(glyph GlyphLike)
	GetRepetition() RepetitionLike
	SetRepetition(repetition RepetitionLike)
	GetFactor() FactorLike
	SetFactor(factor FactorLike)
}

// This constructor creates a new predicate.
func Predicate(glyph GlyphLike, repetition RepetitionLike, factor FactorLike) PredicateLike {
	var v = &predicate{}
	v.SetGlyph(glyph)
	v.SetRepetition(repetition)
	v.SetFactor(factor)
	return v
}

// PREDICATE IMPLEMENTATION

// This type defines the structure and methods associated with a predicate.
type predicate struct {
	glyph      GlyphLike
	repetition RepetitionLike
	factor     FactorLike
}

// This method returns the glyph for this predicate.
func (v *predicate) GetGlyph() GlyphLike {
	return v.glyph
}

// This method sets the glyph for this predicate.
func (v *predicate) SetGlyph(glyph GlyphLike) {
	if glyph != nil {
		v.glyph = glyph
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
		v.glyph = nil
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
		v.glyph = nil
		v.repetition = nil
		v.factor = factor
	}
}
