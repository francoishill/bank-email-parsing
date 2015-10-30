package fnb

type ExpenseTypeVisitor interface {
	VisitPaidFrom(paidFrom *PaidFromExpenseType)
	VisitReservedForPurchase(reservedForPurchase *ReservedForPurchaseExpenseType)
	VisitWithdrawnFrom(withdrawnFrom *WithdrawnFromExpenseType)
	VisitTransferred(transferred *TransferredExpenseType)
}
