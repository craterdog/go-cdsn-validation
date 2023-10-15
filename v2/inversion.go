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

// INVERSION INTERFACE

// This interface defines the methods supported by all inversion-like
// components.
type InversionLike interface {
	GetPredicate() PredicateLike
	SetPredicate(predicate PredicateLike)
}

// This constructor creates a new inversion.
func Inversion(predicate PredicateLike) InversionLike {
	var v = &inversion{}
	v.SetPredicate(predicate)
	return v
}

// INVERSION IMPLEMENTATION

// This type defines the structure and methods associated with a inversion.
type inversion struct {
	predicate PredicateLike
}

// This method returns the predicate for this inversion.
func (v *inversion) GetPredicate() PredicateLike {
	return v.predicate
}

// This method sets the predicate for this inversion.
func (v *inversion) SetPredicate(predicate PredicateLike) {
	if predicate == nil {
		panic("A inversion requires a predicate.")
	}
	v.predicate = predicate
}
