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

type documentClass_ struct {
	// This class does not define any constants.
}

// Private Class Namespace Reference

var documentClass = &documentClass_{
	// This class does not initialize any constants.
}

// Public Class Namespace Access

func DocumentClass() DocumentClassLike {
	return documentClass
}

// Public Class Constructors

func (c *documentClass_) FromGrammar(
	grammar GrammarLike,
) DocumentLike {
	var document = &document_{
		// This class does not initialize any attributes.
	}
	document.SetGrammar(grammar)
	return document
}

// CLASS INSTANCES

// Private Class Type Definition

type document_ struct {
	grammar GrammarLike
}

// Public Interface

func (v *document_) GetGrammar() GrammarLike {
	return v.grammar
}

func (v *document_) SetGrammar(grammar GrammarLike) {
	if grammar == nil {
		panic("A grammar within a document cannot be nil.")
	}
	v.grammar = grammar
}
