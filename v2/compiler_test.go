/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
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
	cds "github.com/craterdog/go-cdsn-validation/v2"
	osx "os"
	tes "testing"
)

const testDirectory = "./test/"

func TestCompiler(t *tes.T) {

	var filename = testDirectory + "test.cdsn"
	fmt.Println(filename)
	var source, err = osx.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var document = cds.ParseDocument(source)
	cds.CompileDocument(testDirectory, "test", document)
}
