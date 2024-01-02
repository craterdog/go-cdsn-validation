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

// CLASS NAMESPACE

// Private Class Namespace Type

type glyphClass_ struct {
	// This class does not define any constants.
}

// Private Class Namespace Reference

var glyphClass = &glyphClass_{
	// This class does not initialize any constants.
}

// Public Class Namespace Access

func GlyphClass() GlyphClassLike {
	return glyphClass
}

// Public Class Constructors

func (c *glyphClass_) FromCharacter(character string) GlyphLike {
	var glyph = &glyph_{
		// This class does not initialize any attributes.
	}
	glyph.SetFirst(character)
	return glyph
}

func (c *glyphClass_) FromRange(first, last string) GlyphLike {
	var glyph = &glyph_{
		// This class does not initialize any attributes.
	}
	glyph.SetFirst(first)
	glyph.SetLast(last)
	return glyph
}

// CLASS INSTANCES

// Private Class Type Definition

type glyph_ struct {
	first string
	last  string
}

// Public Interface

func (v *glyph_) GetFirst() string {
	return v.first
}

func (v *glyph_) GetLast() string {
	return v.last
}

func (v *glyph_) SetFirst(first string) {
	if len(first) < 1 {
		panic("A glyph requires a first character.")
	}
	v.first = first
}

func (v *glyph_) SetLast(last string) {
	v.last = last
}
