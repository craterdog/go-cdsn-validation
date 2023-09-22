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

// ALTERNATIVE INTERFACE

// This interface defines the methods supported by all alternative-like
// components.
type AlternativeLike interface {
	GetFactors() col.Sequential[Factor]
	SetFactors(factors col.Sequential[Factor])
	GetNote() Note
	SetNote(note Note)
}

// This constructor creates a new alternative.
func Alternative(factors col.Sequential[Factor], note Note) AlternativeLike {
	var v = &alternative{}
	v.SetFactors(factors)
	v.SetNote(note)
	return v
}

// ALTERNATIVE IMPLEMENTATION

// This type defines the structure and methods associated with an alternative.
type alternative struct {
	factors col.Sequential[Factor]
	note    Note
}

// This method returns the factors for this alternative.
func (v *alternative) GetFactors() col.Sequential[Factor] {
	return v.factors
}

// This method sets the factors for this alternative.
func (v *alternative) SetFactors(factors col.Sequential[Factor]) {
	if factors == nil || factors.IsEmpty() {
		panic("An alternative requires at least one factor.")
	}
	v.factors = factors
}

// This method returns the note for this alternative.
func (v *alternative) GetNote() Note {
	return v.note
}

// This method sets the note for this alternative.
func (v *alternative) SetNote(note Note) {
	v.note = note
}

// This method attempts to parse an alternative. It returns the alternative and
// whether or not the alternative was successfully parsed.
func (v *parser) parseAlternative() (AlternativeLike, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	var factors = col.List[Factor]()
	var note Note
	var alternative AlternativeLike
	factor, token, ok = v.parseFactor()
	if !ok {
		// An alternative must have at least one factor.
		return alternative, token, false
	}
	for {
		factors.AddValue(factor)
		factor, token, ok = v.parseFactor()
		if !ok {
			// No more factors.
			break
		}
	}
	note, _, _ = v.parseNote() // The note is optional.
	alternative = Alternative(factors, note)
	return alternative, token, true
}

// This private method appends a formatted alternative to the result.
func (v *formatter) formatAlternative(alternative AlternativeLike) {
	var factor Factor
	var factors = alternative.GetFactors()
	var iterator = col.Iterator(factors)
	if iterator.HasNext() {
		factor = iterator.GetNext()
		v.formatFactor(factor)
	}
	for iterator.HasNext() {
		v.appendString(" ")
		factor = iterator.GetNext()
		v.formatFactor(factor)
	}
	var note = alternative.GetNote()
	if len(note) > 0 {
		v.appendString("  ")
		v.formatNote(note)
	}
}
