/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package cdsn

type Element any

// This method attempts to parse an element. It returns the element and whether
// or not the element was successfully parsed.
func (v *parser) parseElement() (Factor, *Token, bool) {
	var ok bool
	var token *Token
	var factor Factor
	factor, token, ok = v.parseIntrinsic()
	if !ok {
		factor, token, ok = v.parseString()
	}
	if !ok {
		factor, token, ok = v.parseNumber()
	}
	if !ok {
		factor, token, ok = v.parseName()
	}
	if !ok {
		factor, token, ok = v.parseLetter()
	}
	return factor, token, ok
}
