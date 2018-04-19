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
)
