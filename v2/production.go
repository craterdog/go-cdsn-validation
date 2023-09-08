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

// PRODUCTION IMPLEMENTATION

// This constructor creates a new production.
func Production(symbol Symbol, definition DefinitionLike) ProductionLike {
	var v = &production{}
	v.SetSymbol(symbol)
	v.SetDefinition(definition)
	return v
}

// This type defines the structure and methods associated with a production.
type production struct {
	symbol     Symbol
	definition DefinitionLike
}

// This method returns the symbol for this production.
func (v *production) GetSymbol() Symbol {
	return v.symbol
}

// This method sets the symbol for this production.
func (v *production) SetSymbol(symbol Symbol) {
	if len(symbol) == 0 {
		panic("A production requires a symbol.")
	}
	v.symbol = symbol
}

// This method returns the definition for this production.
func (v *production) GetDefinition() DefinitionLike {
	return v.definition
}

// This method sets the definition for this production.
func (v *production) SetDefinition(definition DefinitionLike) {
	if definition == nil {
		panic("A production requires a definition.")
	}
	v.definition = definition
}
