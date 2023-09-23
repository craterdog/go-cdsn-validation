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

import (
	fmt "fmt"
	sts "strings"
)

// FORMATTER INTERFACE

// This function returns the bytes containing the canonical format for the
// specified grammar including the POSIX standard EOF marker.
func FormatDocument(grammar GrammarLike) []byte {
	var v = &formatter{}
	v.formatGrammar(grammar)
	var string_ = v.getResult()
	return []byte(string_)
}

// FORMATTER IMPLEMENTATION

// This type defines the structure and methods for a canonical formatter agent.
type formatter struct {
	indentation int
	depth       int
	result      sts.Builder
}

// This method returns the canonically formatted string result.
func (v *formatter) getResult() string {
	var result = v.result.String()
	v.result.Reset()
	return result
}

// This method appends the specified string to the result.
func (v *formatter) appendString(s string) {
	v.result.WriteString(s)
}

// This method appends a properly indented newline to the result.
func (v *formatter) appendNewline() {
	var separator = "\n"
	var levels = v.depth + v.indentation
	for level := 0; level < levels; level++ {
		separator += "    "
	}
	v.result.WriteString(separator)
}

// This private method appends a formatted element to the result.
func (v *formatter) formatElement(element Element) {
	switch e := element.(type) {
	case Intrinsic:
		v.formatIntrinsic(e)
	case String:
		v.formatString(e)
	case Number:
		v.formatNumber(e)
	case Name:
		v.formatName(e)
	default:
		panic(fmt.Sprintf("Attempted to format:\n    element: %v\n    type: %t\n", e, element))
	}
}

// This private method appends a formatted factor to the result.
func (v *formatter) formatFactor(factor Factor) {
	switch f := factor.(type) {
	case *range_:
		v.formatRange(f)
	case *inverse:
		v.formatInverse(f)
	case *exactlyN:
		v.formatExactlyN(f)
	case *zeroOrOne:
		v.formatZeroOrOne(f)
	case *zeroOrMore:
		v.formatZeroOrMore(f)
	case *oneOrMore:
		v.formatOneOrMore(f)
	default:
		v.formatElement(f)
	}
}
