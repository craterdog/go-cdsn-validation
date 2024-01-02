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
)

// CLASS NAMESPACE

// Private Class Namespace Type

type definitionClass_ struct {
	// This class does not define any constants.
}

// Private Class Namespace Reference

var definitionClass = &definitionClass_{
	// This class does not initialize any constants.
}

// Public Class Namespace Access

func DefinitionClass() DefinitionClassLike {
	return definitionClass
}

// Public Class Constructors

func (c *definitionClass_) FromSymbolAndExpression(
	symbol string,
	expression ExpressionLike,
) DefinitionLike {
	var definition = &definition_{
		// This class does not initialize any attributes.
	}
	definition.SetSymbol(symbol)
	definition.SetExpression(expression)
	return definition
}

// CLASS INSTANCES

// Private Class Type Definition

type definition_ struct {
	expression ExpressionLike
	symbol     string
}

// Public Interface

func (v *definition_) GetExpression() ExpressionLike {
	return v.expression
}

func (v *definition_) GetSymbol() string {
	return v.symbol
}

func (v *definition_) SetExpression(expression ExpressionLike) {
	if expression == nil {
		panic("An expression cannot be nil.")
	}
	v.expression = expression
}

func (v *definition_) SetSymbol(symbol string) {
	if len(symbol) < 2 {
		var message = fmt.Sprintf("An invalid symbol was found:\n    %v\n", symbol)
		panic(message)
	}
	v.symbol = symbol
}
