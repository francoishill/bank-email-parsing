package fnb

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

const bulkLines = `
    * FNB :-) R300.00 paid from cheq a/c..123456 @ Online Banking. Avail R15000. Ref.Doorstep Chef. 29Oct 15:07
    * FNB :-) R350.00 paid from cheq a/c..123456 @ Smartapp. Avail R15000. Ref.Electricity 00112233. 29Oct 13:50
    * FNB :-) R4431.70 paid from cheq a/c..123456 @ Eft. Avail R15000. Ref.Wesbank_fi0123456789. 29Oct 00:00
    * FNB :-) R300.00 paid from cheq a/c..123456 @ Online Banking. Avail R15000. Ref.Rates and Taxes. 28Oct 15:40
    * FNB :-) R1065.70 paid from cheq a/c..123456 @ Online Banking. Avail R15000. Ref.Apartment levy. 28Oct 15:40 * FNB :-) R4784.09 paid from cheq a/c..123456 @ Eft. Avail R15000. Ref.Fnb H Loan 0000001. 28Oct 00:00
    * FNB :-) R125.80 reserved for purchase @ Cafe shop from cheq a/c..123456 using card..1234
    * FNB :-) R1056.07 reserved for purchase @ Airline ticket from cheq a/c..123456 using card..1234
    * FNB :-) R173.11 reserved for purchase @ McDonalds from cheq a/c..123456 using card..1234
    * FNB :-) R1081.96 withdrawn from cheq a/c..123456 using card..1234 @ E12345678912312
    * FNB :-) R194.60 reserved for purchase @ Mcdonalds from cheq a/c..123456 using card..1234
    * FNB :-) R186.37 reserved for purchase @ Cafe corner shop from cheq a/c..123456 using card..1234
    * FNB :-) R141.64 reserved for purchase @ Paypal *dropboxirel from cheq a/c..123456 using card..1234 * FNB :-) R585.11 reserved for purchase @ The magazine shop from cheq a/c..123456 using card..1234
    * FNB :-) R2000.00 paid from cheq a/c..123456 @ Eft. Ref.Jill. 26Oct 00:00
    * FNB :-) R866.83 withdrawn from cheq a/c..123456 using card..1234 @ N123456. Avail R15000. 25Oct 21:24
    * FNB :-) R1305.00 paid from cheq a/c..123456 @ Online Banking. Avail R15000. Ref.Donation. 6Oct 15:30 * FNB :-) R200.00 t/fer from cheq a/c..123456 to card a/c..234567 @ Online Banking. Avail R15000. 5Oct 20:45
`

func getLines() []string {
	splittedLines := strings.Split(bulkLines, "\n")

	cleanedLines := []string{}
	for _, s := range splittedLines {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			cleanedLines = append(cleanedLines, s)
		}
	}

	return cleanedLines
}

func TestParsingLines(t *testing.T) {
	Convey("Testing parsing FNB lines", t, func() {
		lines := getLines()
		expenses := ParseLinesAsExpenses(lines)

		Convey("Check type counts", func() {
			So(len(expenses), ShouldEqual, len(lines)+3)
			So(expenses.CountPaidFromExpenseTypes(), ShouldEqual, 8)
			So(expenses.CountReservedForPurchaseExpenseTypes(), ShouldEqual, 7)
			So(expenses.CountWithdrawnFromExpenseTypes(), ShouldEqual, 2)
		})

		expected := []*Expense{
			{"Doorstep Chef", 300.00, "123456", &PaidFromExpenseType{}},
			{"Electricity 00112233", 350.00, "123456", &PaidFromExpenseType{}},
			{"Wesbank_fi0123456789", 4431.70, "123456", &PaidFromExpenseType{}},
			{"Rates and Taxes", 300.00, "123456", &PaidFromExpenseType{}},
			{"Apartment levy", 1065.70, "123456", &PaidFromExpenseType{}},
			{"Fnb H Loan 0000001", 4784.09, "123456", &PaidFromExpenseType{}},
			{"Cafe shop", 125.80, "123456", &ReservedForPurchaseExpenseType{}},
			{"Airline ticket", 1056.07, "123456", &ReservedForPurchaseExpenseType{}},
			{"McDonalds", 173.11, "123456", &ReservedForPurchaseExpenseType{}},
			{"Withdrawal", 1081.96, "123456", &WithdrawnFromExpenseType{}},
			{"Mcdonalds", 194.60, "123456", &ReservedForPurchaseExpenseType{}},
			{"Cafe corner shop", 186.37, "123456", &ReservedForPurchaseExpenseType{}},
			{"Paypal *dropboxirel", 141.64, "123456", &ReservedForPurchaseExpenseType{}},
			{"The magazine shop", 585.11, "123456", &ReservedForPurchaseExpenseType{}},
			{"Jill", 2000.00, "123456", &PaidFromExpenseType{}},
			{"Withdrawal", 866.83, "123456", &WithdrawnFromExpenseType{}},
			{"Donation", 1305.00, "123456", &PaidFromExpenseType{}},
			{"card a/c..234567", 200.00, "123456", &TransferredExpenseType{}},
		}

		Convey("Check variable details (amount, description, type)", func() {
			for i, e := range expected {
				actual := expenses[i]
				So(e.Description, ShouldEqual, actual.Description)
				So(e.Amount, ShouldAlmostEqual, actual.Amount, 0.001)
				So(e.AccountNumberShort, ShouldEqual, actual.AccountNumberShort)
				So(e.Type, ShouldHaveSameTypeAs, actual.Type)
			}
		})
	})
}
