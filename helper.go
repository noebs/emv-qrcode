package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/noebs/emv-qrcode/emv/mpm"
)

//QR response to be used in Cashq and other products
type QR struct {
	MerchantID     string  `json:"merchant_id,omitempty"`
	MerchantName   string  `json:"merchant_name,omitempty"`
	MerchantBankID string  `json:"merchant_bank_id,omitempty"`
	Amount         float64 `json:"amount,omitempty"`
	AcquirerID     string  `json:"acquirer_id,omitempty"`
	RawQR          *mpm.EMVQR
}

func (qr *QR) parseAccount(e *mpm.EMVQR) error {
	// We have to iterate through all *possible* tags. Previously it was 26
	// Merchant ids are from 2...50

	if id, ok := e.MerchantAccountInformation["26"]; ok { // for the common case
		p := id.Value.PaymentNetworkSpecific
		// 0 = acquirer id, 1 = merchant account
		//Acquirer ID and merchant ID are mandatory per EMV Book 4 QR specs
		if len(p) < 2 {
			return errors.New("specs are wrong")
		}
		qr.AcquirerID = p[0].Value
		qr.MerchantBankID = p[1].Value
		qr.MerchantID = p[1].Value
	} else {
		// Iterate through EVERY variable!
		for i := 2; i <= 51; i++ {
			parsedIndex := fmt.Sprintf("%02d", i)
			if id, ok := e.MerchantAccountInformation[mpm.ID(parsedIndex)]; ok { // for the common case
				p := id.Value.PaymentNetworkSpecific
				// 0 = acquirer id, 1 = merchant account
				//Acquirer ID and merchant ID are mandatory per EMV Book 4 QR specs
				if len(p) < 2 {
					return errors.New("specs are wrong")
				}
				qr.AcquirerID = p[0].Value
				qr.MerchantBankID = p[1].Value
				qr.MerchantID = p[1].Value

			}
		}
	}
	return nil
}

func (qr *QR) init(e *mpm.EMVQR) error {
	var err error
	qr.RawQR = e

	qr.MerchantName = e.MerchantName.Value
	qr.Amount, err = strconv.ParseFloat(e.TransactionAmount.Value, 32)
	if err != nil {
		return err
	}
	if err = qr.parseAccount(e); err != nil {
		return err
	}
	return nil

}
