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

type cardinalityClass_ struct {
	// This class does not define any constants.
}

// Private Class Namespace Reference

var cardinalityClass = &cardinalityClass_{
	// This class does not initialize any constants.
}

// Public Class Namespace Access

func CardinalityClass() CardinalityClassLike {
	return cardinalityClass
}

// Public Class Constructors

func (c *cardinalityClass_) FromConstraint(constraint ConstraintLike) CardinalityLike {
	var cardinality = &cardinality_{
		// This class does not initialize any attributes.
	}
	cardinality.SetConstraint(constraint)
	return cardinality
}

// CLASS INSTANCES

// Private Class Type Definition

type cardinality_ struct {
	constraint ConstraintLike
}

// Public Interface

func (v *cardinality_) GetConstraint() ConstraintLike {
	return v.constraint
}

func (v *cardinality_) SetConstraint(constraint ConstraintLike) {
	if constraint == nil {
		panic("A constraint must not be nil.")
	}
	v.constraint = constraint
}
