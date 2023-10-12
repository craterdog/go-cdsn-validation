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
	byt "bytes"
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v2"
	osx "os"
	sts "strings"
	uni "unicode"
)

// COMPILER INTERFACE

// This function compiles the specified document into its corresponding parser.
func CompileDocument(directory, packageName string, document DocumentLike) {
	var v = &compiler{directory: directory, packageName: packageName}
	v.initializeConfiguration()
	v.initializeScanner()
	v.initializeParser()
	v.compileDocument(document)
	v.finalizeScanner()
	v.finalizeParser()
}

// COMPILER IMPLEMENTATION

// This constant defines the maximum depth allowed to prevent overly complex
// grammars and infinite recursion in rule definitions.
const maximumDepth = 8

// This private function determines whether or not the specified name is a token
// name.
func isTokenName(name NAME) bool {
	return uni.IsUpper(rune(name[1]))
}

// This private function replaces all occurences of the target string with the
// specified name.
func replaceName(template string, target string, name string) string {
	var nameLower, nameUpper string
	var targetLower = "#" + target + "#"
	var targetUpper = "#" + sts.ToUpper(target[0:1]) + target[1:] + "#"
	if isTokenName(NAME(name)) {
		nameLower = sts.ToLower(name)
		nameUpper = name
	} else {
		nameLower = name
		nameUpper = sts.ToUpper(name[0:1]) + name[1:]
	}
	template = sts.ReplaceAll(template, targetLower, nameLower)
	template = sts.ReplaceAll(template, targetUpper, nameUpper)
	return template
}

// This type defines the structure and methods for a compiler agent.
type compiler struct {
	depth         int
	directory     string
	packageName   string
	scannerBuffer byt.Buffer
	parserBuffer  byt.Buffer
}

// This method creates a new configuration (package.go) file if one
// does not already exist.
func (v *compiler) initializeConfiguration() {
	var err error
	var template []byte
	var configuration = v.directory + "/package.go"
	_, err = osx.Open(configuration)
	if err != nil {
		template, err = osx.ReadFile("./templates/package.tp")
		if err != nil {
			panic(err)
		}
		template = byt.ReplaceAll(template, []byte("#package#"), []byte(v.packageName))
		err = osx.WriteFile(configuration, template, 0666)
		if err != nil {
			panic(err)
		}
	}
}

// This method creates the byte buffer for the generated scanner code.
func (v *compiler) initializeScanner() {
	var template, err = osx.ReadFile("./templates/scanner.tp")
	if err != nil {
		panic(err)
	}
	template = byt.ReplaceAll(template, []byte("#package#"), []byte(v.packageName))
	v.scannerBuffer.Write(template[0 : len(template)-1])
}

// This method creates the byte buffer for the generated parser code.
func (v *compiler) initializeParser() {
	var template, err = osx.ReadFile("./templates/parser.tp")
	if err != nil {
		panic(err)
	}
	template = byt.ReplaceAll(template, []byte("#package#"), []byte(v.packageName))
	v.parserBuffer.Write(template[0 : len(template)-1])
}

// This method appends the scan token template for the specified name to
// the scanner byte buffer.
func (v *compiler) appendScanToken(name NAME, re string) {
	const template = `

type #Token# string
const Token#Token# TokenType = "#Token#"
const #token# = #tokenRE#
var #token#Scanner = reg.MustCompile(` + "`^(?:` + #token# + `)`)" + `

// This method adds a new #token# token with the current scanner information
// to the token channel. It returns true if a new #token# token was found.
func (v *scanner) scan#Token#() bool {
	var s = v.source[v.nextByte:]
	var matches = bytesToStrings(#token#Scanner.FindSubmatch(s))
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token#Token#)
		v.line += sts.Count(matches[0], EOL)
		return true
	}
	return false
}`
	var snippet = replaceName(template, "token", string(name))
	snippet = replaceName(snippet, "tokenRE", re)
	v.scannerBuffer.WriteString(snippet)
}

// This method appends the parse token template for the specified name
// to the parser byte buffer.
func (v *compiler) appendParseToken(name NAME) {
	const template = `

// This method attempts to parse a new #token# token. It returns the token
// and whether or not the token was successfully parsed.
func (v *parser) parse#Token#() (#Token#, *Token, bool) {
	var #token#_ #Token#
	var token = v.nextToken()
	if token.Type != Token#Token# {
		v.backupOne(token)
		return #token#_, token, false
	}
	#token#_ = #Token#(token.Value)
	return #token#_, token, true
}`
	var snippet = replaceName(template, "token", string(name))
	v.parserBuffer.WriteString(snippet)
}

// This method appends the parse rule start template for the specified
// name to the parser byte buffer.
func (v *compiler) appendParseRuleStart(name NAME) {
	const template = `

// This method attempts to parse a new #rule#. It returns the #rule#
// and whether or not the #rule# was successfully parsed.
func (v *parser) parse#Rule#() (#Rule#Like, *Token, bool) {
	var ok bool
	var token *Token
	var #rule#_ #Rule#Like`
	var snippet = replaceName(template, "rule", string(name))
	v.parserBuffer.WriteString(snippet)
}

// This method appends the parse rule end template for the specified
// name and arguments to the parser byte buffer.
func (v *compiler) appendParseRuleEnd(name NAME, arguments string) {
	const template = `
	#rule#_ = #Rule#(#arguments#)
	return #rule#_, token, true
}`
	var snippet = replaceName(template, "rule", string(name))
	snippet = replaceName(snippet, "arguments", arguments)
	v.parserBuffer.WriteString(snippet)
}

// This method writes the byte buffer for the generated scanner code into
// a file.
func (v *compiler) finalizeScanner() {
	var filename = v.directory + "scanner.go"
	var err = osx.WriteFile(filename, v.scannerBuffer.Bytes(), 0666)
	if err != nil {
		panic(err)
	}
}

// This method writes the byte buffer for the generated parser code into
// a file.
func (v *compiler) finalizeParser() {
	var filename = v.directory + "parser.go"
	var err = osx.WriteFile(filename, v.parserBuffer.Bytes(), 0666)
	if err != nil {
		panic(err)
	}
}

// This method increments the depth of the compilation by one and checks
// for run-away recursion.
func (v *compiler) incrementDepth() {
	v.depth++
	if v.depth > maximumDepth {
		var message = fmt.Sprintf("Exceeded the maximum depth:\n%s\n", v.parserBuffer.String())
		panic(message)
	}
}

// This method decrements the depth of the compilation by one and checks
// for run-away recursion.
func (v *compiler) decrementDepth() {
	v.depth--
}

// This method compiles the specified definition.
func (v *compiler) compileDefinition(definition DefinitionLike) {
	var symbol = definition.GetSYMBOL()
	var expression = definition.GetExpression()
	var name = symbol.GetNAME()
	switch {
	case string(name) == "INTRINSIC":
		// Intrinsics are automatically part of every parser.
	case isTokenName(name):
		var re sts.Builder
		v.compileTokenExpression(expression, &re)
		v.appendScanToken(name, re.String())
		v.appendParseToken(name)
	default:
		var arguments sts.Builder
		v.appendParseRuleStart(name)
		v.compileRuleExpression(expression, &arguments)
		v.appendParseRuleEnd(name, arguments.String())
	}
}

// This method compiles the specified document.
func (v *compiler) compileDocument(document DocumentLike) {
	var statements = document.GetStatements()
	var iterator = col.Iterator(statements)
	for iterator.HasNext() {
		var statement = iterator.GetNext()
		v.compileStatement(statement)
	}
}

// This method compiles the specified statement.
func (v *compiler) compileStatement(statement Statement) {
	switch actual := statement.(type) {
	case *definition:
		v.compileDefinition(actual)
	case COMMENT:
	}
}

// This method compiles the specified token alternative.
func (v *compiler) compileTokenAlternative(alternative AlternativeLike, re *sts.Builder) {
	var predicates = alternative.GetPredicates()
	var iterator = col.Iterator(predicates)
	var predicate = iterator.GetNext()
	v.compileTokenPredicate(predicate, re)
	for iterator.HasNext() {
		re.WriteString("|")
		predicate = iterator.GetNext()
		v.compileTokenPredicate(predicate, re)
	}
}

// This method compiles the specified token constraint.
func (v *compiler) compileTokenCONSTRAINT(constraint CONSTRAINT, re *sts.Builder) {
	re.WriteString(string(constraint))
}

// This method compiles the specified token element.
func (v *compiler) compileTokenElement(element ElementLike, re *sts.Builder) {
	var intrinsic = element.GetINTRINSIC()
	var name = element.GetNAME()
	var string_ = element.GetSTRING()
	switch {
	case len(intrinsic) > 0:
		v.compileTokenINTRINSIC(intrinsic, re)
	case len(name) > 0:
		v.compileTokenNAME(name, re)
	case len(string_) > 0:
		v.compileTokenSTRING(string_, re)
	}
}

// This method compiles the specified rule expression.
func (v *compiler) compileRuleExpression(expression ExpressionLike, arguments *sts.Builder) {
	v.incrementDepth()
	arguments.WriteString("foo")
	v.decrementDepth()
}

// This method compiles the specified token expression.
func (v *compiler) compileTokenExpression(expression ExpressionLike, re *sts.Builder) {
	var alternatives = expression.GetAlternatives()
	var iterator = col.Iterator(alternatives)
	var alternative = iterator.GetNext()
	v.compileTokenAlternative(alternative, re)
	for iterator.HasNext() {
		re.WriteString("|")
		alternative = iterator.GetNext()
		v.compileTokenAlternative(alternative, re)
	}
}

// This method compiles the specified token factor.
func (v *compiler) compileTokenFactor(factor FactorLike, re *sts.Builder) {
	var precedence = factor.GetPrecedence()
	var element = factor.GetElement()
	switch {
	case precedence != nil:
		v.compileTokenPrecedence(precedence, re)
	case element != nil:
		v.compileTokenElement(element, re)
	}
}

// This method compiles the specified token intrinsic.
func (v *compiler) compileTokenINTRINSIC(intrinsic INTRINSIC, re *sts.Builder) {
	switch string(intrinsic) {
	case "ANY":
		re.WriteString(any_)
	case "LOWER_CASE":
		re.WriteString(lowerCase)
	case "UPPER_CASE":
		re.WriteString(upperCase)
	case "DIGIT":
		re.WriteString(digit)
	case "SEPARATOR":
		re.WriteString(separator)
	case "ESCAPE":
		re.WriteString(escape)
	case "EOL":
		re.WriteString(eol)
	}
}

// This method compiles the specified token name.
func (v *compiler) compileTokenNAME(name NAME, re *sts.Builder) {
	re.WriteString(string(name))
}

// This method compiles the specified token precedence.
func (v *compiler) compileTokenPrecedence(precedence PrecedenceLike, re *sts.Builder) {
	var expression = precedence.GetExpression()
	re.WriteString("(?:")
	v.compileTokenExpression(expression, re)
	re.WriteString(")")
}

// This method compiles the specified token predicate.
func (v *compiler) compileTokenPredicate(predicate PredicateLike, re *sts.Builder) {
	var range_ = predicate.GetRange()
	var repetition = predicate.GetRepetition()
	var factor = predicate.GetFactor()
	switch {
	case range_ != nil:
		v.compileTokenRange(range_, re)
	case repetition != nil:
		v.compileTokenRepetition(repetition, re)
	case factor != nil:
		v.compileTokenFactor(factor, re)
	}
}

// This method compiles the specified token range.
func (v *compiler) compileTokenRange(range_ RangeLike, re *sts.Builder) {
	var first = range_.GetFirstCHARACTER()
	re.WriteString(string(first))
	var last = range_.GetLastCHARACTER()
	if len(last) > 0 {
		re.WriteString("-")
		re.WriteString(string(last))
	}
}

// This method compiles the specified token repetition.
func (v *compiler) compileTokenRepetition(repetition RepetitionLike, re *sts.Builder) {
	var constraint = repetition.GetCONSTRAINT()
	var factor = repetition.GetFactor()
	switch string(constraint) {
	case "~":
		re.WriteString("[^")
		v.compileTokenFactor(factor, re)
		re.WriteString("]")
	case "?", "*", "+":
		v.compileTokenFactor(factor, re)
		v.compileTokenCONSTRAINT(constraint, re)
	default:
		v.compileTokenFactor(factor, re)
		re.WriteString("{")
		v.compileTokenCONSTRAINT(constraint, re)
		re.WriteString("}")
	}
}

// This method compiles the specified token string.
func (v *compiler) compileTokenSTRING(string_ STRING, re *sts.Builder) {
	var s = string(string_)
	re.WriteString(s[1 : len(s)-1])
}
