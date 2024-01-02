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

type precedenceClass_ struct {
	// This class does not define any constants.
}

// Private Class Namespace Reference

var precedenceClass = &precedenceClass_{
	// This class does not initialize any constants.
}

// Public Class Namespace Access

func PrecedenceClass() PrecedenceClassLike {
	return precedenceClass
}

// Public Class Constructors

func (c *precedenceClass_) FromExpression(
	expression ExpressionLike,
) PrecedenceLike {
	var precedence = &precedence_{
		// This class does not initialize any attributes.
	}
	precedence.SetExpression(expression)
	return precedence
}

// CLASS INSTANCES

// Private Class Type Definition

type precedence_ struct {
	expression ExpressionLike
}

// Public Interface

func (v *precedence_) GetExpression() ExpressionLike {
	return v.expression
}

func (v *precedence_) SetExpression(expression ExpressionLike) {
	if expression == nil {
		panic("The expression within a precedence cannot be nil.")
	}
	v.expression = expression
}
