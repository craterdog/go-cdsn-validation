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

// PREDICATE INTERFACE

// This interface defines the methods supported by all predicate-like
// components.
type PredicateLike interface {
	GetRange() RangeLike
	SetRange(range_ RangeLike)
	GetRepetition() RepetitionLike
	SetRepetition(repetition_ RepetitionLike)
	GetFactor() FactorLike
	SetFactor(factor_ FactorLike)
}

// This constructor creates a new predicate.
func Predicate(range_ RangeLike, repetition_ RepetitionLike, factor_ FactorLike) PredicateLike {
	var v = &predicate_{}
	v.SetRange(range_)
	v.SetRepetition(repetition_)
	v.SetFactor(factor_)
	return v
}

// PREDICATE IMPLEMENTATION

// This type defines the structure and methods associated with a predicate.
type predicate_ struct {
	range_      RangeLike
	repetition_ RepetitionLike
	factor_     FactorLike
}

// This method returns the range for this predicate.
func (v *predicate_) GetRange() RangeLike {
	return v.range_
}

// This method sets the range for this predicate.
func (v *predicate_) SetRange(range_ RangeLike) {
	if range_ != nil {
		v.range_ = range_
		v.repetition_ = nil
		v.factor_ = nil
	}
}

// This method returns the repetition for this predicate.
func (v *predicate_) GetRepetition() RepetitionLike {
	return v.repetition_
}

// This method sets the repetition for this predicate.
func (v *predicate_) SetRepetition(repetition_ RepetitionLike) {
	if repetition_ != nil {
		v.range_ = nil
		v.repetition_ = repetition_
		v.factor_ = nil
	}
}

// This method returns the factor for this predicate.
func (v *predicate_) GetFactor() FactorLike {
	return v.factor_
}

// This method sets the factor for this predicate.
func (v *predicate_) SetFactor(factor_ FactorLike) {
	if factor_ != nil {
		v.range_ = nil
		v.repetition_ = nil
		v.factor_ = factor_
	}
}
