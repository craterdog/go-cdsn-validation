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

type factorClass_ struct {
	// This class does not define any constants.
}

// Private Class Namespace Reference

var factorClass = &factorClass_{
	// This class does not initialize any constants.
}

// Public Class Namespace Access

func FactorClass() FactorClassLike {
	return factorClass
}

// Public Class Constructors

func (c *factorClass_) FromPredicate(
	predicate PredicateLike,
) FactorLike {
	var factor = &factor_{
		// This class does not initialize any attributes.
	}
	factor.SetPredicate(predicate)
	return factor
}

// CLASS INSTANCES

// Private Class Type Definition

type factor_ struct {
	cardinality CardinalityLike
	predicate   PredicateLike
}

// Public Interface

func (v *factor_) GetCardinality() CardinalityLike {
	return v.cardinality
}

func (v *factor_) GetPredicate() PredicateLike {
	return v.predicate
}

func (v *factor_) SetCardinality(cardinality CardinalityLike) {
	if cardinality == nil {
		panic("An cardinality cannot be nil.")
	}
	v.cardinality = cardinality
}

func (v *factor_) SetPredicate(predicate PredicateLike) {
	if predicate == nil {
		panic("A predicate within a factor cannot be nil.")
	}
	v.predicate = predicate
}
