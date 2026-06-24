package keeper

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
)

// Decimal stores exact monetary amounts as int64 micro-units.
// 1 AEQ = 1_000_000 micro-AEQ. This eliminates float64 rounding errors
// that accumulate across many ledger operations.
//
// JSON serializes as a decimal number (e.g. 1000.000000) for API
// backward-compatibility — existing clients see no change.
type Decimal int64

const DecimalPrecision = int64(1_000_000)

// NewDecimal converts a float64 AEQ amount to Decimal, rounding to 6dp.
func NewDecimal(aeq float64) Decimal {
	return Decimal(math.Round(aeq * float64(DecimalPrecision)))
}

// NewDecimalFromMicro creates a Decimal directly from micro-units.
func NewDecimalFromMicro(micro int64) Decimal { return Decimal(micro) }

// Float returns the AEQ value as float64 (for display/legacy code).
func (d Decimal) Float() float64 { return float64(d) / float64(DecimalPrecision) }

// Micro returns the underlying int64 micro-unit value.
func (d Decimal) Micro() int64 { return int64(d) }

// Add returns d + other. Exact.
func (d Decimal) Add(other Decimal) Decimal { return d + other }

// Sub returns d - other. Exact.
func (d Decimal) Sub(other Decimal) Decimal { return d - other }

// MulFloat multiplies by a float64 (e.g. for rate/fee calculations), rounding.
func (d Decimal) MulFloat(f float64) Decimal {
	// P3-10: guard against int64 overflow when f > 1 and d is large
	result := math.Round(float64(d) * f)
	const maxD = float64(math.MaxInt64)
	if result > maxD {
		return Decimal(math.MaxInt64)
	}
	if result < -maxD {
		return Decimal(math.MinInt64)
	}
	return Decimal(result)
}

// DivDecimal divides two Decimals and returns a Decimal result (not a ratio).
// Computes (d * precision) / other to maintain scale.
func (d Decimal) DivDecimal(other Decimal) Decimal {
	if other == 0 {
		return 0
	}
	return Decimal(new(big.Int).Div(
		new(big.Int).Mul(big.NewInt(int64(d)), big.NewInt(DecimalPrecision)),
		big.NewInt(int64(other)),
	).Int64())
}

// IsZero returns true when d == 0.
func (d Decimal) IsZero() bool { return d == 0 }

// IsPositive returns true when d > 0.
func (d Decimal) IsPositive() bool { return d > 0 }

// IsNegative returns true when d < 0.
func (d Decimal) IsNegative() bool { return d < 0 }

// Neg returns -d.
func (d Decimal) Neg() Decimal { return -d }

// String formats as "X.YYYYYY AEQ" for debugging.
func (d Decimal) String() string { return fmt.Sprintf("%.6f", d.Float()) }

// MarshalJSON serializes as a plain float number for API compatibility.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Float())
}

// UnmarshalJSON accepts a float number and converts to micro-units.
func (d *Decimal) UnmarshalJSON(b []byte) error {
	var f float64
	if err := json.Unmarshal(b, &f); err != nil {
		return err
	}
	*d = NewDecimal(f)
	return nil
}

// AMMSwapOut computes the constant-product AMM output amount:
// amountOut = (reserveOut * amountInAfterFee) / (reserveIn + amountInAfterFee)
// Uses big.Int arithmetic to avoid int64 overflow for large reserves.
func AMMSwapOut(reserveIn, reserveOut, amountInAfterFee Decimal) Decimal {
	ri := new(big.Int).SetInt64(int64(reserveIn))
	ro := new(big.Int).SetInt64(int64(reserveOut))
	ai := new(big.Int).SetInt64(int64(amountInAfterFee))
	// numerator = reserveOut * amountIn
	num := new(big.Int).Mul(ro, ai)
	// denominator = reserveIn + amountIn
	den := new(big.Int).Add(ri, ai)
	if den.Sign() == 0 {
		return 0
	}
	result := new(big.Int).Div(num, den)
	// P0-4: guard against overflow — result must fit in int64.
	// If it doesn't, the swap is too large for the current pool size.
	if result.BitLen() > 63 || result.Sign() < 0 {
		return 0 // caller will detect zero output and reject the swap
	}
	return Decimal(result.Int64())
}
