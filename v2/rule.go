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

// This constructor creates a new rule.
func Rule(options col.Sequential[OptionLike]) RuleLike {
	var v = &rule{}
	v.SetOptions(options)
	return v
}

// This type defines the structure and methods associated with a rule.
type rule struct {
	multilined bool
	options col.Sequential[OptionLike]
}

// This method determines whether or not this rule is multlined.
func (v *rule) IsMultilined() bool {
	return v.multilined
}

// This method returns the options for this rule.
func (v *rule) GetOptions() col.Sequential[OptionLike] {
	return v.options
}

// This method sets the options for this rule.
func (v *rule) SetOptions(options col.Sequential[OptionLike]) {
	if options == nil || options.IsEmpty() {
		panic("A rule requires at least one option.")
	}
	var iterator = col.Iterator(options)
	for iterator.HasNext() {
		var option = iterator.GetNext()
		if options.GetSize() > 1 && (option.GetFactors().GetSize() > 2 || len(option.GetNote()) > 0) {
			v.multilined = true
			break
		}
	}
	v.options = options
}
