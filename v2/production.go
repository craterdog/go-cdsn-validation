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

// PRODUCTION IMPLEMENTATION

// This constructor creates a new production.
func Production(symbol Symbol, rule RuleLike, note Note) ProductionLike {
	var v = &production{}
	v.SetSymbol(symbol)
	v.SetRule(rule)
	v.SetNote(note)
	return v
}

// This type defines the structure and methods associated with a production.
type production struct {
	symbol  Symbol
	rule    RuleLike
	note    Note
}

// This method returns the symbol for this production.
func (v *production) GetSymbol() Symbol {
	return v.symbol
}

// This method sets the symbol for this production.
func (v *production) SetSymbol(symbol Symbol) {
	if symbol == nil {
		panic("A production requires a symbol.")
	}
	v.symbol = symbol
}

// This method returns the rule for this production.
func (v *production) GetRule() RuleLike {
	return v.rule
}

// This method sets the rule for this production.
func (v *production) SetRule(rule RuleLike) {
	if rule == nil {
		panic("A production requires a rule.")
	}
	v.rule = rule
}

// This method returns the note for this production.
func (v *production) GetNote() Note {
	return v.note
}

// This method sets the note for this production.
func (v *production) SetNote(note Note) {
	v.note = note
}
