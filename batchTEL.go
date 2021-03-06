// Copyright 2018 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import "fmt"

// BatchTEL is a batch that handles SEC payment type Telephone-Initiated Entries (TEL)
// Telephone-Initiated Entries (TEL) are consumer debit transactions. The NACHA Operating Rules permit TEL entries when the Originator obtains the Receiver’s authorization for the debit entry orally via the telephone.
// An entry based upon a Receiver’s oral authorization must utilize the TEL (Telephone-Initiated Entry) Standard Entry Class (SEC) Code.
type BatchTEL struct {
	batch
}

// NewBatchTEL returns a *BatchTEL
func NewBatchTEL(bh *BatchHeader) *BatchTEL {
	batch := new(BatchTEL)
	batch.SetControl(NewBatchControl())
	batch.SetHeader(bh)
	return batch
}

// Validate ensures the batch meets NACHA rules specific to the SEC type TEL
func (batch *BatchTEL) Validate() error {
	// basic verification of the batch before we validate specific rules.
	if err := batch.verify(); err != nil {
		return err
	}
	// Add configuration based validation for this type.
	// TEL can not have an addenda
	if err := batch.isAddendaCount(0); err != nil {
		return err
	}

	// Add type specific validation.
	if batch.Header.StandardEntryClassCode != "TEL" {
		msg := fmt.Sprintf(msgBatchSECType, batch.Header.StandardEntryClassCode, "TEL")
		return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "StandardEntryClassCode", Msg: msg}
	}
	// can not have credits in TEL batches
	for _, entry := range batch.Entries {
		if entry.CreditOrDebit() != "D" {
			msg := fmt.Sprintf(msgBatchTransactionCodeCredit, entry.IndividualName)
			return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "TransactionCode", Msg: msg}
		}
	}

	return batch.isPaymentTypeCode()
}

// Create builds the batch sequence numbers and batch control. Additional creation
func (batch *BatchTEL) Create() error {
	// generates sequence numbers and batch control
	if err := batch.build(); err != nil {
		return err
	}

	return batch.Validate()
}
