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
	GetPredicates() col.Sequential[PredicateLike]
	SetPredicates(predicates col.Sequential[PredicateLike])
	GetNOTE() NOTE
	SetNOTE(note NOTE)
}

// This constructor creates a new alternative.
func Alternative(predicates col.Sequential[PredicateLike], note NOTE) AlternativeLike {
	var v = &alternative{}
	v.SetPredicates(predicates)
	v.SetNOTE(note)
	return v
}

// ALTERNATIVE IMPLEMENTATION

// This type defines the structure and methods associated with an alternative.
type alternative struct {
	predicates col.Sequential[PredicateLike]
	note       NOTE
}

// This method returns the predicates for this alternative.
func (v *alternative) GetPredicates() col.Sequential[PredicateLike] {
	return v.predicates
}

// This method sets the predicates for this alternative.
func (v *alternative) SetPredicates(predicates col.Sequential[PredicateLike]) {
	if predicates == nil || predicates.IsEmpty() {
		panic("An alternative requires at least one predicate.")
	}
	v.predicates = predicates
}

// This method returns the note for this alternative.
func (v *alternative) GetNOTE() NOTE {
	return v.note
}

// This method sets the note for this alternative.
func (v *alternative) SetNOTE(note NOTE) {
	v.note = note
}
