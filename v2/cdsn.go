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

// These characters include all lowercase unicode letters.
const lowercase = `\p{Ll}`

// These characters include all uppercase unicode letters.
const uppercase = `\p{Lu}`

// These characters include all unicode digits.
const digit = `\p{Nd}`

// These characters can be used to separate words in names.
const separator = `_`

// This string contains the actual characters `\` and `n`, not EOL.
const eol = `\n`

// These characters are interpreted as escape characters by the scanner.
const (
	escape  = `\\(?:(?:` + unicode + `)|[abfnrtv'"\\])`
	unicode = `u` + base16 + `{4}|U` + base16 + `{8}`
	base16  = `[0-9a-f]`
)

// These characters are treated as whitespace by the scanner and ignored.
const (
	whitespace = `(?:` + space + `|` + eol + `)+`
	space      = ` `
)