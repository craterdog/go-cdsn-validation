/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package cdsn_test

import (
	fmt "fmt"
	val "github.com/craterdog/go-cdsn-validation/v2"
	osx "os"
	tes "testing"
)

const cdsn = "./test/cdsn.cdsn"

func TestGenerateGrammar(t *tes.T) {
	var err = osx.WriteFile(cdsn, []byte(val.FormatGrammar()), 0644)
	if err != nil {
		var message = fmt.Sprintf("Could not create the cdsn file: %v.", err)
		panic(message)
	}
}
