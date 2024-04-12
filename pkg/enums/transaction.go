package enums

type TxType string

const (
	CREDIT TxType = "Credit"
	DEBIT  TxType = "Debit"
)

type TxStatus string

const (
	RESERVED TxStatus = "Reserved"
	SUCCESS  TxStatus = "Success"
	FAILED   TxStatus = "Failed"
)
