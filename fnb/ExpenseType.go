package fnb

type ExpenseType interface {
	Accept(ExpenseTypeVisitor)
}
