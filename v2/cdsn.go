/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

/*
This package defines a parser and formatter for documents written using Crater
Dog Syntax Notation™ (CDSN).  The parser performs validation on the resulting
parse tree.  The formatter takes a validated parse tree and generates the
corresponding CDSN document using the canonical format.
*/
package cdsn

// CONFIGURATION PARAMETERS

// These constant definitions capture regular expression subpatterns for the
// intrinsic token types.
const (
	eol       = `\n`     // This contains the actual characters `\` and `n`, not EOL.
	ignored   = ` |\n`   // These characters are treated as whitespace by the parser.
	lowercase = `\p{Ll}` // This includes all unicode lowercase letters.
	uppercase = `\p{Lu}` // This includes all unicode uppercase letters.
	digit     = `\p{Nd}` // This includes all unicode digits.
	separator = `_`      // This can be used to separate words in names.
)
