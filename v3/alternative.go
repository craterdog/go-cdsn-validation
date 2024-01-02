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
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v3"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type alternativeClass_ struct {
	// This class does not define any constants.
}

// Private Class Namespace Reference

var alternativeClass = &alternativeClass_{
	// This class does not initialize any constants.
}

// Public Class Namespace Access

func AlternativeClass() AlternativeClassLike {
	return alternativeClass
}

// Public Class Constructors

func (c *alternativeClass_) FromFactors(factors col.Sequential[FactorLike]) AlternativeLike {
	var alternative = &alternative_{
		// This class does not initialize any attributes.
	}
	alternative.SetFactors(factors)
	return alternative
}

// CLASS INSTANCES

// Private Class Type Definition

type alternative_ struct {
	factors col.Sequential[FactorLike]
	note    string
}

// Public Interface

func (v *alternative_) GetFactors() col.Sequential[FactorLike] {
	return v.factors
}

func (v *alternative_) GetNote() string {
	return v.note
}

func (v *alternative_) SetFactors(factors col.Sequential[FactorLike]) {
	if factors == nil || factors.IsEmpty() {
		panic("An alternative must have at least one factor.")
	}
	v.factors = factors
}

func (v *alternative_) SetNote(note string) {
	if len(note) < 2 {
		var message = fmt.Sprintf("An invalid note was found:\n    %v\n", note)
		panic(message)
	}
	v.note = note
}
