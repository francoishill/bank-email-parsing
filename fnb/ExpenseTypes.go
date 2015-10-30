package fnb

var AllExpenseTypes []ExpenseType = []ExpenseType{
	&PaidFromExpenseType{},
	&ReservedForPurchaseExpenseType{},
	&WithdrawnFromExpenseType{},
	&TransferredExpenseType{},
}

type PaidFromExpenseType struct{}

func (p *PaidFromExpenseType) Accept(visitor ExpenseTypeVisitor) {
	visitor.VisitPaidFrom(p)
}

type ReservedForPurchaseExpenseType struct{}

func (r *ReservedForPurchaseExpenseType) Accept(visitor ExpenseTypeVisitor) {
	visitor.VisitReservedForPurchase(r)
}

type WithdrawnFromExpenseType struct{}

func (w *WithdrawnFromExpenseType) Accept(visitor ExpenseTypeVisitor) {
	visitor.VisitWithdrawnFrom(w)
}

type TransferredExpenseType struct{}

func (t *TransferredExpenseType) Accept(visitor ExpenseTypeVisitor) {
	visitor.VisitTransferred(t)
}
