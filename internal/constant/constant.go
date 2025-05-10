package constant

type LoanPeriodUnit string

const (
	PeriodUnitWeek  = "WEEK"
	PeriodUnitMonth = "MONTH"
)

type LoanPaymentStatus string

const (
	LoanPaymentStatusUnpaid = "UNPAID"
	LoanPaymentStatusPaid   = "PAID"
)
