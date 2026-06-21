package keeper

import (
	"fmt"
	"math"
	"testing"
)

// newTestState creates an in-memory ChainState with no file/DB backend.
func newTestState() *ChainState {
	return &ChainState{
		accounts: make(map[string]*AccountState),
		pool:     &PoolState{},
		useDB:    false,
	}
}

// addHuman inserts a registered human account directly into state for testing.
// LastActivityAt=0 means effectiveBalance returns Balance with no demurrage.
func addHuman(cs *ChainState, addr string, balance float64) {
	cs.accounts[addr] = &AccountState{
		Address: addr,
		Balance: NewDecimal(balance),
		IsHuman: true,
	}
}

// --- CalcGini ---

func TestCalcGini_TwoEqualHumans(t *testing.T) {
	cs := newTestState()
	addHuman(cs, "0x01", 1000)
	addHuman(cs, "0x02", 1000)
	g := cs.CalcGini()
	if g != 0.0 {
		t.Errorf("equal balances: want Gini=0, got %v", g)
	}
}

func TestCalcGini_TwoHumansMaxConcentration(t *testing.T) {
	// One human has nearly everything. With n=2, biased Gini would cap at 0.5;
	// the unbiased estimator (×n/(n-1)) must push it close to 1.
	cs := newTestState()
	addHuman(cs, "0x01", 1)
	addHuman(cs, "0x02", 9999)
	g := cs.CalcGini()
	// biased = 9998/(2*10000) ≈ 0.4999; unbiased ≈ 0.9998 — must be > 0.99
	if g < 0.99 {
		t.Errorf("near-total concentration: want Gini>0.99, got %v", g)
	}
}

func TestCalcGini_ThreeEqualHumans(t *testing.T) {
	cs := newTestState()
	addHuman(cs, "0x01", 500)
	addHuman(cs, "0x02", 500)
	addHuman(cs, "0x03", 500)
	g := cs.CalcGini()
	if g != 0.0 {
		t.Errorf("equal balances (3): want Gini=0, got %v", g)
	}
}

func TestCalcGini_ThreeHumansKnownValue(t *testing.T) {
	// balances=[1,2,6], n=3, sum=9
	// biased Gini: numerator = (-2)*1 + 0*2 + 2*6 = 10; gini = 10/27 ≈ 0.3704
	// unbiased = 0.3704 * (3/2) ≈ 0.5556
	cs := newTestState()
	addHuman(cs, "0x01", 1)
	addHuman(cs, "0x02", 2)
	addHuman(cs, "0x03", 6)
	g := cs.CalcGini()
	want := 10.0 / 27.0 * 3.0 / 2.0 // exact unbiased value
	if math.Abs(g-want) > 1e-9 {
		t.Errorf("three humans [1,2,6]: want %v, got %v", want, g)
	}
}

func TestCalcGini_SingleHuman_ReturnsZero(t *testing.T) {
	cs := newTestState()
	addHuman(cs, "0x01", 1000)
	g := cs.CalcGini()
	if g != 0.0 {
		t.Errorf("single human: want 0, got %v", g)
	}
}

func TestCalcGini_NonHumanBalancesIgnored(t *testing.T) {
	cs := newTestState()
	addHuman(cs, "0x01", 100)
	addHuman(cs, "0x02", 100)
	// Non-human with huge balance — must not skew Gini
	cs.accounts["0x99"] = &AccountState{Address: "0x99", Balance: NewDecimal(1_000_000), IsHuman: false}
	g := cs.CalcGini()
	if g != 0.0 {
		t.Errorf("non-human balance must be ignored: want Gini=0, got %v", g)
	}
}

// --- CalcAequitasIndex ---

func TestCalcAequitasIndex_EqualHumans(t *testing.T) {
	cs := newTestState()
	addHuman(cs, "0x01", 1000)
	addHuman(cs, "0x02", 1000)
	idx := cs.CalcAequitasIndex()
	if idx != 0.0 {
		t.Errorf("equal balances: want Index=0, got %v", idx)
	}
}

func TestCalcAequitasIndex_IsGiniTimes100Truncated(t *testing.T) {
	cs := newTestState()
	addHuman(cs, "0x01", 1)
	addHuman(cs, "0x02", 2)
	addHuman(cs, "0x03", 6)
	gini := cs.CalcGini()
	idx := cs.CalcAequitasIndex()
	// Index = floor(gini*1000) / 10  (truncates to 1 decimal place)
	want := math.Trunc(gini*1000) / 10
	if idx != want {
		t.Errorf("index=%v, want gini*100 truncated to 1dp=%v (gini=%v)", idx, want, gini)
	}
}

// --- bootstrapMultiplierLocked ---

func TestBootstrapMultiplier_ZeroHumans(t *testing.T) {
	cs := newTestState()
	cs.mu.RLock()
	m := cs.bootstrapMultiplierLocked()
	cs.mu.RUnlock()
	if m != 5.0 {
		t.Errorf("0 humans: want 5.0, got %v", m)
	}
}

func TestBootstrapMultiplier_OneHuman(t *testing.T) {
	cs := newTestState()
	addHuman(cs, "0x01", 0)
	cs.mu.RLock()
	m := cs.bootstrapMultiplierLocked()
	cs.mu.RUnlock()
	if m != 5.0 {
		t.Errorf("1 human: want 5.0, got %v", m)
	}
}

func TestBootstrapMultiplier_FiveHumans(t *testing.T) {
	cs := newTestState()
	for i := 0; i < 5; i++ {
		addHuman(cs, fmt.Sprintf("0xF%d", i), 0)
	}
	cs.mu.RLock()
	m := cs.bootstrapMultiplierLocked()
	cs.mu.RUnlock()
	if m != 5.0 {
		t.Errorf("5 humans: want 5.0, got %v", m)
	}
}

func TestBootstrapMultiplier_TenHumans(t *testing.T) {
	cs := newTestState()
	for i := 0; i < 10; i++ {
		addHuman(cs, fmt.Sprintf("0xA%d", i), 0)
	}
	cs.mu.RLock()
	m := cs.bootstrapMultiplierLocked()
	cs.mu.RUnlock()
	if m != 10.0 {
		t.Errorf("10 humans: want 10.0, got %v", m)
	}
}

func TestBootstrapMultiplier_TwentyFourHumans(t *testing.T) {
	cs := newTestState()
	for i := 0; i < 24; i++ {
		addHuman(cs, fmt.Sprintf("0xB%d", i), 0)
	}
	cs.mu.RLock()
	m := cs.bootstrapMultiplierLocked()
	cs.mu.RUnlock()
	if m != 24.0 {
		t.Errorf("24 humans: want 24.0, got %v", m)
	}
}

func TestBootstrapMultiplier_TwentyFiveHumans(t *testing.T) {
	cs := newTestState()
	for i := 0; i < 25; i++ {
		addHuman(cs, fmt.Sprintf("0xC%d", i), 0)
	}
	cs.mu.RLock()
	m := cs.bootstrapMultiplierLocked()
	cs.mu.RUnlock()
	if m != 25.0 {
		t.Errorf("25 humans: want 25.0 (full cap), got %v", m)
	}
}

func TestBootstrapMultiplier_ThirtyHumans(t *testing.T) {
	cs := newTestState()
	for i := 0; i < 30; i++ {
		addHuman(cs, fmt.Sprintf("0xD%d", i), 0)
	}
	cs.mu.RLock()
	m := cs.bootstrapMultiplierLocked()
	cs.mu.RUnlock()
	if m != 25.0 {
		t.Errorf("30 humans: want 25.0 (capped at wealthCapMultiplier), got %v", m)
	}
}

// --- enforceWealthCapLocked ---

func TestEnforceWealthCap_BelowCap_NoChange(t *testing.T) {
	cs := newTestState()
	// 2 humans at 500 each: avg=500, multiplier=5 (floor), cap=2500 > 500 → no cap
	addHuman(cs, "0x01", 500)
	addHuman(cs, "0x02", 500)
	acc := cs.accounts["0x01"]
	cs.mu.Lock()
	cs.enforceWealthCapLocked(acc)
	cs.mu.Unlock()
	if acc.Balance.Float() != 500 {
		t.Errorf("under cap: balance should stay 500, got %v", acc.Balance.Float())
	}
}

func TestEnforceWealthCap_AboveCap_ExcessRedistributed(t *testing.T) {
	// The cap bites only when N > multiplier (otherwise the rich account
	// pulls the average up and is never capped by its own contribution).
	// With 26 humans (multiplier=25) and 25 poor at 10 AEQ + 1 rich at 100_000:
	//   avg = (100_000 + 25*10) / 26  ≈ 3855.77
	//   cap = 25 * avg                ≈ 96394.23
	//   excess = 100_000 - cap        ≈ 3605.77
	cs := newTestState()
	for i := 0; i < 25; i++ {
		addHuman(cs, fmt.Sprintf("0xpoor%02d", i), 10)
	}
	addHuman(cs, "0xrich", 100_000)
	acc := cs.accounts["0xrich"]
	avg := (100_000.0 + 25*10.0) / 26.0
	expectedCap := 25.0 * avg
	cs.mu.Lock()
	cs.enforceWealthCapLocked(acc)
	cs.mu.Unlock()
	if math.Abs(acc.Balance.Float()-expectedCap) > 1e-6 {
		t.Errorf("after cap: want balance=%.6f, got %.6f", expectedCap, acc.Balance.Float())
	}
	poolTotal := cs.accounts[validatorsPoolAddr].Balance.Float() +
		cs.accounts[lpPoolAddr].Balance.Float() +
		cs.accounts[ubiPoolAddr].Balance.Float() +
		cs.accounts[treasuryPoolAddr].Balance.Float()
	expectedExcess := 100_000.0 - expectedCap
	if math.Abs(poolTotal-expectedExcess) > 1e-6 {
		t.Errorf("pool total: want %.6f excess redistributed, got %.6f", expectedExcess, poolTotal)
	}
}

func TestEnforceWealthCap_PoolAddresses_Exempt(t *testing.T) {
	cs := newTestState()
	addHuman(cs, "0x01", 100)
	// Pool address with huge balance — must be exempt from cap
	cs.accounts[validatorsPoolAddr] = &AccountState{Address: validatorsPoolAddr, Balance: NewDecimal(1_000_000)}
	acc := cs.accounts[validatorsPoolAddr]
	cs.mu.Lock()
	cs.enforceWealthCapLocked(acc)
	cs.mu.Unlock()
	if acc.Balance.Float() != 1_000_000 {
		t.Errorf("pool address must be exempt from cap, balance changed to %v", acc.Balance.Float())
	}
}
