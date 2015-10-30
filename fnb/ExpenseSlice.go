package fnb

type ExpenseSlice []*Expense

func (e ExpenseSlice) CountPaidFromExpenseTypes() int {
	cnt := 0
	for _, ex := range e {
		if _, ok := ex.Type.(*PaidFromExpenseType); ok {
			cnt++
		}
	}
	return cnt
}

func (e ExpenseSlice) CountReservedForPurchaseExpenseTypes() int {
	cnt := 0
	for _, ex := range e {
		if _, ok := ex.Type.(*ReservedForPurchaseExpenseType); ok {
			cnt++
		}
	}
	return cnt
}

func (e ExpenseSlice) CountWithdrawnFromExpenseTypes() int {
	cnt := 0
	for _, ex := range e {
		if _, ok := ex.Type.(*WithdrawnFromExpenseType); ok {
			cnt++
		}
	}
	return cnt
}
