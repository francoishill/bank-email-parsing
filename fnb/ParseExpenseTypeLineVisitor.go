package fnb

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"regexp"
	"strconv"
)

type ParseExpenseTypeLineVisitor struct {
	stringLine    string
	FoundExpenses []*Expense
}

var paidFromPattern = regexp.MustCompile(`R([0-9]+\.[0-9]{2}) paid from cheq a/c\.\.([0-9]+).* Ref\.([^\.]*)\.`)
var reservedForPattern = regexp.MustCompile(`R([0-9]+\.[0-9]{2}) reserved for purchase @ (.+) from cheq a/c\.\.([0-9]+) using card\.\.[0-9]+`)
var withdrawnFromPattern = regexp.MustCompile(`R([0-9]+\.[0-9]{2}) withdrawn from cheq a/c\.\.([0-9]+) using card\.\.[0-9]+ @`)
var transferredPattern = regexp.MustCompile(`R([0-9]+\.[0-9]{2}) t/fer from .* a/c\.\.([0-9]+) to (.* a/c\.\.[0-9]+) @`)

func extractExpenseFromPattern(expenseType ExpenseType, pattern *regexp.Regexp, line string, randValIndex, descriptionIndex, accountNumIndex int) []*Expense {
	if !pattern.MatchString(line) {
		return nil
	}

	submatches := pattern.FindAllStringSubmatch(line, -1)

	expenses := []*Expense{}
	for _, m := range submatches {
		randVal := mustStringToFloat32(m[randValIndex])
		refVal := m[descriptionIndex]
		accountNum := m[accountNumIndex]

		expenses = append(expenses, &Expense{
			refVal,
			randVal,
			accountNum,
			expenseType,
		})
	}

	return expenses
}

func (p *ParseExpenseTypeLineVisitor) VisitPaidFrom(paidFrom *PaidFromExpenseType) {
	if expenses := extractExpenseFromPattern(paidFrom, paidFromPattern, p.stringLine, 1, 3, 2); len(expenses) > 0 {
		p.FoundExpenses = append(p.FoundExpenses, expenses...)
	}
}

func (p *ParseExpenseTypeLineVisitor) VisitReservedForPurchase(reservedForPurchase *ReservedForPurchaseExpenseType) {
	if expenses := extractExpenseFromPattern(reservedForPurchase, reservedForPattern, p.stringLine, 1, 2, 3); len(expenses) > 0 {
		p.FoundExpenses = append(p.FoundExpenses, expenses...)
	}
}

func (p *ParseExpenseTypeLineVisitor) VisitWithdrawnFrom(withdrawnFrom *WithdrawnFromExpenseType) {
	if expenses := extractExpenseFromPattern(withdrawnFrom, withdrawnFromPattern, p.stringLine, 1, 1, 2); len(expenses) > 0 {
		for _, e := range expenses {
			e.Description = "Withdrawal"
		}
		p.FoundExpenses = append(p.FoundExpenses, expenses...)
	}
}

func (p *ParseExpenseTypeLineVisitor) VisitTransferred(transferred *TransferredExpenseType) {
	if expenses := extractExpenseFromPattern(transferred, transferredPattern, p.stringLine, 1, 3, 2); len(expenses) > 0 {
		p.FoundExpenses = append(p.FoundExpenses, expenses...)
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

			if len(visitor.FoundExpenses) > 0 {
				expenses = append(expenses, visitor.FoundExpenses...)
			}
		}
	}

	return expenses
}
