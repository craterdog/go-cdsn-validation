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
	GetGlyph() GlyphLike
	SetGlyph(glyph GlyphLike)
	GetElement() ElementLike
	SetElement(element ElementLike)
	GetPrecedence() PrecedenceLike
	SetPrecedence(precedence PrecedenceLike)
	GetInversion() InversionLike
	SetInversion(inversion InversionLike)
}

// This constructor creates a new factor.
func Factor(glyph GlyphLike, element ElementLike, precedence PrecedenceLike, inversion InversionLike) FactorLike {
	if glyph == nil && element == nil && precedence == nil && inversion == nil {
		panic("A factor requires at least one of its attributes to be set.")
	}
	var v = &factor{}
	v.SetGlyph(glyph)
	v.SetElement(element)
	v.SetPrecedence(precedence)
	v.SetInversion(inversion)
	return v
}

// FACTOR IMPLEMENTATION

// This type defines the structure and methods associated with a factor.
type factor struct {
	glyph      GlyphLike
	element    ElementLike
	precedence PrecedenceLike
	inversion  InversionLike
}

// This method returns the glyph for this factor.
func (v *factor) GetGlyph() GlyphLike {
	return v.glyph
}

// This method sets the glyph for this factor.
func (v *factor) SetGlyph(glyph GlyphLike) {
	if glyph != nil {
		v.element = nil
		v.glyph = glyph
		v.inversion = nil
		v.precedence = nil
	}
}

// This method returns the element for this factor.
func (v *factor) GetElement() ElementLike {
	return v.element
}

// This method sets the element for this factor.
func (v *factor) SetElement(element ElementLike) {
	if element != nil {
		v.element = element
		v.glyph = nil
		v.inversion = nil
		v.precedence = nil
	}
}

// This method returns the precedence for this factor.
func (v *factor) GetPrecedence() PrecedenceLike {
	return v.precedence
}

// This method sets the precedence for this factor.
func (v *factor) SetPrecedence(precedence PrecedenceLike) {
	if precedence != nil {
		v.element = nil
		v.glyph = nil
		v.inversion = nil
		v.precedence = precedence
	}
}

// This method returns the inversion for this factor.
func (v *factor) GetInversion() InversionLike {
	return v.inversion
}

// This method sets the inversion for this factor.
func (v *factor) SetInversion(inversion InversionLike) {
	if inversion != nil {
		v.element = nil
		v.glyph = nil
		v.inversion = inversion
		v.precedence = nil
	}
}
