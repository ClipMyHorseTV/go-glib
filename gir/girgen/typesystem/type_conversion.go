package typesystem

import "fmt"

type ConvertibleType interface {
	Type
	CanTransferToGlib(transfer TransferOwnership) bool
	CanTransferFromGlib(transfer TransferOwnership) bool

	GetTransferToGlibFunction(transfer TransferOwnership) string
	GetTransferFromGlibFunction(transfer TransferOwnership) string

	GoUnsafeFromGlibFullFunction() string
	GoUnsafeFromGlibNoneFunction() string
	GoUnsafeFromGlibBorrowFunction() string

	GoUnsafeToGlibFullFunction() string
	GoUnsafeToGlibNoneFunction() string
}

type BaseConversions struct {
	FromGlibBorrowFunction string
	FromGlibFullFunction   string
	FromGlibNoneFunction   string
	ToGlibNoneFunction     string
	ToGlibFullFunction     string
}

// CanTransferFromGlib implements ConvertibleType.
func (b BaseConversions) CanTransferFromGlib(transfer TransferOwnership) bool {
	switch transfer {
	case TransferNone:
		return b.FromGlibNoneFunction != ""
	case TransferFull:
		return b.FromGlibFullFunction != ""
	case TransferBorrow:
		return b.FromGlibBorrowFunction != ""
	default:
		return false
	}
}

// CanTransferToGlib implements ConvertibleType.
func (b BaseConversions) CanTransferToGlib(transfer TransferOwnership) bool {
	switch transfer {
	case TransferNone:
		return b.ToGlibNoneFunction != ""
	case TransferFull:
		return b.ToGlibFullFunction != ""
	default:
		return false
	}
}

// GetTransferFromGlibFunction implements ConvertibleType.
func (b BaseConversions) GetTransferFromGlibFunction(transfer TransferOwnership) string {
	switch transfer {
	case TransferNone:
		return b.FromGlibNoneFunction
	case TransferFull:
		return b.FromGlibFullFunction
	case TransferBorrow:
		return b.FromGlibBorrowFunction
	default:
		return ""
	}
}

// GetTransferToGlibFunction implements ConvertibleType.
func (b BaseConversions) GetTransferToGlibFunction(transfer TransferOwnership) string {
	switch transfer {
	case TransferNone:
		return b.ToGlibNoneFunction
	case TransferFull:
		return b.ToGlibFullFunction
	default:
		return ""
	}
}

func newDefaultBaseConversions(name string) BaseConversions {
	return BaseConversions{
		FromGlibBorrowFunction: fmt.Sprintf("Unsafe%sFromGlibBorrow", name),
		FromGlibNoneFunction:   fmt.Sprintf("Unsafe%sFromGlibNone", name),
		FromGlibFullFunction:   fmt.Sprintf("Unsafe%sFromGlibFull", name),

		ToGlibNoneFunction: fmt.Sprintf("Unsafe%sToGlibNone", name),
		ToGlibFullFunction: fmt.Sprintf("Unsafe%sToGlibFull", name),
	}
}

// GoUnsafeFromGlibBorrowFunction implements ConvertibleType.
func (b BaseConversions) GoUnsafeFromGlibBorrowFunction() string {
	return b.FromGlibBorrowFunction
}

// GoUnsafeFromGlibFullFunction implements ConvertibleType.
func (b BaseConversions) GoUnsafeFromGlibFullFunction() string {
	return b.FromGlibFullFunction
}

// GoUnsafeFromGlibNoneFunction implements ConvertibleType.
func (b BaseConversions) GoUnsafeFromGlibNoneFunction() string {
	return b.FromGlibNoneFunction
}

// GoUnsafeToGlibFullFunction implements ConvertibleType.
func (b BaseConversions) GoUnsafeToGlibFullFunction() string {
	return b.ToGlibFullFunction
}

// GoUnsafeToGlibNoneFunction implements ConvertibleType.
func (b BaseConversions) GoUnsafeToGlibNoneFunction() string {
	return b.ToGlibNoneFunction
}
