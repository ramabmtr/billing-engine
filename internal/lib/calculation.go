package lib

import (
	"github.com/ramabmtr/billing-engine/internal/constant"
	"github.com/shopspring/decimal"
)

var periodToYears = map[constant.LoanPeriodUnit]decimal.Decimal{
	constant.PeriodUnitWeek:  decimal.NewFromInt(52),
	constant.PeriodUnitMonth: decimal.NewFromInt(12),
}

func CalculateTotalRepayment(
	principal decimal.Decimal,
	annualInterestRate decimal.Decimal,
	period int,
	periodUnit constant.LoanPeriodUnit,
) decimal.Decimal {
	// Convert period to years
	py := periodToYears[periodUnit]
	years := decimal.NewFromInt(int64(period)).Div(py)

	rate := annualInterestRate.Div(decimal.NewFromInt(100))
	interest := principal.Mul(rate).Mul(years)

	return principal.Add(interest)
}
