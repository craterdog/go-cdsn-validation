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
	col "github.com/craterdog/go-collection-framework/v3"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type expressionClass_ struct {
	// This class does not define any constants.
}

// Private Class Namespace Reference

var expressionClass = &expressionClass_{
	// This class does not initialize any constants.
}

// Public Class Namespace Access

func ExpressionClass() ExpressionClassLike {
	return expressionClass
}

// Public Class Constructors

func (c *expressionClass_) FromAlternatives(alternatives col.Sequential[AlternativeLike]) ExpressionLike {
	var expression = &expression_{
		// This class does not initialize any attributes.
	}
	expression.SetAlternatives(alternatives)
	return expression
}

// CLASS INSTANCES

// Private Class Type Definition

type expression_ struct {
	alternatives col.Sequential[AlternativeLike]
	isMultilined bool
}

// Public Interface

func (v *expression_) GetAlternatives() col.Sequential[AlternativeLike] {
	return v.alternatives
}

func (v *expression_) IsMultilined() bool {
	return v.isMultilined
}

func (v *expression_) SetAlternatives(alternatives col.Sequential[AlternativeLike]) {
	if alternatives == nil || alternatives.IsEmpty() {
		panic("An expression must have at least one alternative.")
	}
	v.alternatives = alternatives
}

func (v *expression_) SetMultilined(isMultilined bool) {
	v.isMultilined = isMultilined
}
