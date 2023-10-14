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

// RANGE INTERFACE

// This interface defines the methods supported by all glyph-like components.
type GlyphLike interface {
	GetFirstCHARACTER() CHARACTER
	SetFirstCHARACTER(first CHARACTER)
	GetLastCHARACTER() CHARACTER
	SetLastCHARACTER(last CHARACTER)
}

// This constructor creates a new glyph.
func Glyph(first CHARACTER, last CHARACTER) GlyphLike {
	var v = &glyph{}
	v.SetFirstCHARACTER(first)
	v.SetLastCHARACTER(last)
	return v
}

// RANGE IMPLEMENTATION

// This type defines the structure and methods associated with a glyph.
type glyph struct {
	first CHARACTER
	last  CHARACTER
}

// This method returns the first character in the range for this glyph.
func (v *glyph) GetFirstCHARACTER() CHARACTER {
	return v.first
}

// This method sets the first character in the range for this glyph.
func (v *glyph) SetFirstCHARACTER(first CHARACTER) {
	if len(first) == 0 {
		panic("A glyph requires at least one character.")
	}
	v.first = first
}

// This method returns the last character in the range for this glyph.
func (v *glyph) GetLastCHARACTER() CHARACTER {
	return v.last
}

// This method sets the last character in the range for this glyph.
func (v *glyph) SetLastCHARACTER(last CHARACTER) {
	v.last = last
}
