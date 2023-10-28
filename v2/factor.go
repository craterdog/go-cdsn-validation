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
	GetPredicate() PredicateLike
	SetPredicate(predicate PredicateLike)
	GetCardinality() CardinalityLike
	SetCardinality(cardinality CardinalityLike)
}

// This constructor creates a new factor.
func Factor(predicate PredicateLike, cardinality CardinalityLike) FactorLike {
	if predicate == nil {
		panic("A factor requires a predicate to be set.")
	}
	var v = &factor{}
	v.SetPredicate(predicate)
	v.SetCardinality(cardinality)
	return v
}

// FACTOR IMPLEMENTATION

// This type defines the structure and methods associated with a factor.
type factor struct {
	predicate   PredicateLike
	cardinality CardinalityLike
}

// This method returns the predicate for this factor.
func (v *factor) GetPredicate() PredicateLike {
	return v.predicate
}

// This method sets the predicate for this factor.
func (v *factor) SetPredicate(predicate PredicateLike) {
	if predicate != nil {
		v.predicate = predicate
	}
}

// This method returns the cardinality for this factor.
func (v *factor) GetCardinality() CardinalityLike {
	return v.cardinality
}

// This method sets the cardinality for this factor.
func (v *factor) SetCardinality(cardinality CardinalityLike) {
	v.cardinality = cardinality
}
