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

type predicateClass_ struct {
	// This class does not define any constants.
}

// Private Class Namespace Reference

var predicateClass = &predicateClass_{
	// This class does not initialize any constants.
}

// Public Class Namespace Access

func PredicateClass() PredicateClassLike {
	return predicateClass
}

// Public Class Constructors

func (c *predicateClass_) FromAssertion(
	assertion AssertionLike,
	isInverted bool,
) PredicateLike {
	var predicate = &predicate_{
		// This class does not initialize any attributes.
	}
	predicate.SetAssertion(assertion)
	predicate.SetInverted(isInverted)
	return predicate
}

// CLASS INSTANCES

// Private Class Type Definition

type predicate_ struct {
	assertion  AssertionLike
	isInverted bool
}

// Public Interface

func (v *predicate_) GetAssertion() AssertionLike {
	return v.assertion
}

func (v *predicate_) IsInverted() bool {
	return v.isInverted
}

func (v *predicate_) SetAssertion(assertion AssertionLike) {
	if assertion == nil {
		panic("An assertion must not be nil.")
	}
	v.assertion = assertion
}

func (v *predicate_) SetInverted(isInverted bool) {
	v.isInverted = isInverted
}
