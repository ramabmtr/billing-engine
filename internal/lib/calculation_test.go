package lib

import (
	"testing"

	"github.com/ramabmtr/billing-engine/internal/constant"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCalculateTotalRepayment(t *testing.T) {
	tests := []struct {
		name               string
		principal          decimal.Decimal
		annualInterestRate decimal.Decimal
		period             int
		periodUnit         constant.LoanPeriodUnit
		expected           decimal.Decimal
	}{
		{
			name:               "Weekly Period - 1 year",
			principal:          decimal.NewFromInt(1_000_000),
			annualInterestRate: decimal.NewFromInt(10),
			period:             52,
			periodUnit:         constant.PeriodUnitWeek,
			expected:           decimal.NewFromInt(1_100_000),
		},
		{
			name:               "Weekly Period - 6 months",
			principal:          decimal.NewFromInt(1_000_000),
			annualInterestRate: decimal.NewFromInt(10),
			period:             26,
			periodUnit:         constant.PeriodUnitWeek,
			expected:           decimal.NewFromInt(1_050_000),
		},
		{
			name:               "Monthly Period - 1 year",
			principal:          decimal.NewFromInt(1_000_000),
			annualInterestRate: decimal.NewFromInt(10),
			period:             12,
			periodUnit:         constant.PeriodUnitMonth,
			expected:           decimal.NewFromInt(1_100_000),
		},
		{
			name:               "Monthly Period - 6 months",
			principal:          decimal.NewFromInt(1_000_000),
			annualInterestRate: decimal.NewFromInt(10),
			period:             6,
			periodUnit:         constant.PeriodUnitMonth,
			expected:           decimal.NewFromInt(1_050_000),
		},
		{
			name:               "Zero Interest Rate",
			principal:          decimal.NewFromInt(1_000_000),
			annualInterestRate: decimal.NewFromInt(0),
			period:             12,
			periodUnit:         constant.PeriodUnitMonth,
			expected:           decimal.NewFromInt(1_000_000),
		},
		{
			name:               "Decimal Interest Rate",
			principal:          decimal.NewFromInt(1_000_000),
			annualInterestRate: decimal.NewFromFloat(5.5),
			period:             12,
			periodUnit:         constant.PeriodUnitMonth,
			expected:           decimal.NewFromInt(1_055_000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTotalRepayment(
				tt.principal,
				tt.annualInterestRate,
				tt.period,
				tt.periodUnit,
			)
			assert.True(t, tt.expected.Equal(result),
				"Expected %s but got %s", tt.expected.String(), result.String())
		})
	}
}
