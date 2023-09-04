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

import (
	col "github.com/craterdog/go-collection-framework/v2"
)

// OPTION IMPLEMENTATION

// This constructor creates a new option.
func Option(factors col.Sequential[Factor], note Note) OptionLike {
	var v = &option{}
	v.SetFactors(factors)
	v.SetNote(note)
	return v
}

// This type defines the structure and methods associated with an option.
type option struct {
	factors col.Sequential[Factor]
	note Note
}

// This method returns the factors for this option.
func (v *option) GetFactors() col.Sequential[Factor] {
	return v.factors
}

// This method sets the factors for this option.
func (v *option) SetFactors(factors col.Sequential[Factor]) {
	if factors == nil || factors.IsEmpty() {
		panic("An option requires at least one factor.")
	}
	v.factors = factors
}

// This method returns the note for this option.
func (v *option) GetNote() Note {
	return v.note
}

// This method sets the note for this option.
func (v *option) SetNote(note Note) {
	v.note = note
}
