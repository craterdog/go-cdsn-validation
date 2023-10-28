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
	GetElement() ElementLike
	SetElement(element ElementLike)
	GetPrecedence() PrecedenceLike
	SetPrecedence(precedence PrecedenceLike)
	GetInversion() InversionLike
	SetInversion(inversion InversionLike)
}

// This constructor creates a new predicate.
func Predicate(glyph GlyphLike, element ElementLike, precedence PrecedenceLike, inversion InversionLike) PredicateLike {
	if glyph == nil && element == nil && precedence == nil && inversion == nil {
		panic("A predicate requires at least one of its attributes to be set.")
	}
	var v = &predicate{}
	v.SetGlyph(glyph)
	v.SetElement(element)
	v.SetPrecedence(precedence)
	v.SetInversion(inversion)
	return v
}

// PREDICATE IMPLEMENTATION

// This type defines the structure and methods associated with a predicate.
type predicate struct {
	glyph      GlyphLike
	element    ElementLike
	precedence PrecedenceLike
	inversion  InversionLike
}

// This method returns the glyph for this predicate.
func (v *predicate) GetGlyph() GlyphLike {
	return v.glyph
}

// This method sets the glyph for this predicate.
func (v *predicate) SetGlyph(glyph GlyphLike) {
	if glyph != nil {
		v.element = nil
		v.glyph = glyph
		v.inversion = nil
		v.precedence = nil
	}
}

// This method returns the element for this predicate.
func (v *predicate) GetElement() ElementLike {
	return v.element
}

// This method sets the element for this predicate.
func (v *predicate) SetElement(element ElementLike) {
	if element != nil {
		v.element = element
		v.glyph = nil
		v.inversion = nil
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
		v.element = nil
		v.glyph = nil
		v.inversion = nil
		v.precedence = precedence
	}
}

// This method returns the inversion for this predicate.
func (v *predicate) GetInversion() InversionLike {
	return v.inversion
}

// This method sets the inversion for this predicate.
func (v *predicate) SetInversion(inversion InversionLike) {
	if inversion != nil {
		v.element = nil
		v.glyph = nil
		v.inversion = inversion
		v.precedence = nil
	}
}
