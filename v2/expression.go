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

// RULE IMPLEMENTATION

// This constructor creates a new expression.
func Expression(alternatives col.Sequential[AlternativeLike]) ExpressionLike {
	var v = &expression{}
	v.SetAlternatives(alternatives)
	return v
}

// This type defines the structure and methods associated with an expression.
type expression struct {
	multilined   bool
	alternatives col.Sequential[AlternativeLike]
}

// This method determines whether or not this expression is multlined.
func (v *expression) IsMultilined() bool {
	return v.multilined
}

// This method sets whether or not this expression is multlined.
func (v *expression) SetMultilined(multilined bool) {
	v.multilined = multilined
}

// This method returns the alternatives for this expression.
func (v *expression) GetAlternatives() col.Sequential[AlternativeLike] {
	return v.alternatives
}

// This method sets the alternatives for this expression.
func (v *expression) SetAlternatives(alternatives col.Sequential[AlternativeLike]) {
	if alternatives == nil || alternatives.IsEmpty() {
		panic("A expression requires at least one alternative.")
	}
	var iterator = col.Iterator(alternatives)
	for iterator.HasNext() {
		var alternative = iterator.GetNext()
		if alternatives.GetSize() > 1 && (alternative.GetFactors().GetSize() > 2 || len(alternative.GetNote()) > 0) {
			v.multilined = true
			break
		}
	}
	v.alternatives = alternatives
}
