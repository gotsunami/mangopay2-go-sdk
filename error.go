// Copyright 2015 GoTsunami. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

// Documents error code from https://docs.mangopay.com/api-references/error-codes/

const (
	// PayIn web errors
	ErrUserNotRedirected                        = "001031"
	ErrUserCancelledPayment                     = "001031"
	ErrUserFillingPaymentCardDetails            = "001032"
	ErrUserNotRedirectedPaymentSessionExpired   = "001033"
	ErrUserLetPaymentSessionExpireWithoutPaying = "001034"

	// Generic transaction errors
	ErrUserNotCompleteTransaction = "101001"
	ErrTransactionCancelledByUser = "101002"
	ErrTransactionAmountTooHigh   = "001011"

	// 3DSecure errors
	Err3DSNotAvailable         = "101399"
	Err3DSSessionExpired       = "101304"
	Err3DSCardNotCompatible    = "101303"
	Err3DSCardNotEnrolled      = "101302"
	Err3DSAuthenticationFailed = "101301"
)
