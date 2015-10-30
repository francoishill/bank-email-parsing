package fnb

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"regexp"
	"strconv"
)

type ParseExpenseTypeLineVisitor struct {
	stringLine   string
	FoundExpense *Expense
}

var paidFromPattern = regexp.MustCompile(`R([0-9]+\.[0-9]{2}) paid from cheq a/c\.\.([0-9]+).* Ref\.([^\.]*)\.`)
var reservedForPattern = regexp.MustCompile(`R([0-9]+\.[0-9]{2}) reserved for purchase @ (.+) from cheq a/c\.\.([0-9]+) using card\.\.[0-9]+`)
var withdrawnFromPattern = regexp.MustCompile(`R([0-9]+\.[0-9]{2}) withdrawn from cheq a/c\.\.([0-9]+) using card\.\.[0-9]+ @`)
var transferredPattern = regexp.MustCompile(`R([0-9]+\.[0-9]{2}) t/fer from .* a/c\.\.([0-9]+) to (.* a/c\.\.[0-9]+) @`)

func extractExpenseFromPattern(expenseType ExpenseType, pattern *regexp.Regexp, line string, randValIndex, descriptionIndex, accountNumIndex int) *Expense {
	if !pattern.MatchString(line) {
		return nil
	}

	submatches := pattern.FindAllStringSubmatch(line, -1)
	firstMatch := submatches[0]

	randVal := mustStringToFloat32(firstMatch[randValIndex])
	refVal := firstMatch[descriptionIndex]
	accountNum := firstMatch[accountNumIndex]

	return &Expense{
		refVal,
		randVal,
		accountNum,
		expenseType,
	}
}

func (p *ParseExpenseTypeLineVisitor) VisitPaidFrom(paidFrom *PaidFromExpenseType) {
	if expense := extractExpenseFromPattern(paidFrom, paidFromPattern, p.stringLine, 1, 3, 2); expense != nil {
		p.FoundExpense = expense
	}
}

func (p *ParseExpenseTypeLineVisitor) VisitReservedForPurchase(reservedForPurchase *ReservedForPurchaseExpenseType) {
	if expense := extractExpenseFromPattern(reservedForPurchase, reservedForPattern, p.stringLine, 1, 2, 3); expense != nil {
		p.FoundExpense = expense
	}
}

func (p *ParseExpenseTypeLineVisitor) VisitWithdrawnFrom(withdrawnFrom *WithdrawnFromExpenseType) {
	if expense := extractExpenseFromPattern(withdrawnFrom, withdrawnFromPattern, p.stringLine, 1, 1, 2); expense != nil {
		expense.Description = "Withdrawal"
		p.FoundExpense = expense
	}
}

func (p *ParseExpenseTypeLineVisitor) VisitTransferred(transferred *TransferredExpenseType) {
	if expense := extractExpenseFromPattern(transferred, transferredPattern, p.stringLine, 1, 3, 2); expense != nil {
		p.FoundExpense = expense
	}
}

func mustStringToFloat32(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	CheckError(err)
	return float32(f)
}

func ParseLinesAsExpenses(cleanedLines []string) ExpenseSlice {
	expenses := []*Expense{}
	for _, l := range cleanedLines {
		for _, et := range AllExpenseTypes {
			visitor := &ParseExpenseTypeLineVisitor{}
			visitor.stringLine = l
			et.Accept(visitor)

			if visitor.FoundExpense != nil {
				expenses = append(expenses, visitor.FoundExpense)
				break
			}
		}
	}

	return expenses
}
