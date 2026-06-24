import re, sys

def r(path, old, new, lbl):
    c = open(path, 'rb').read().decode('utf-8')
    if old in c:
        c = c.replace(old, new, 1)
        open(path, 'wb').write(c.encode('utf-8'))
        sys.stderr.write('OK: ' + lbl + '\n')
        return True
    sys.stderr.write('MISS: ' + lbl + ' in ' + path + '\n')
    return False

# P0-3 evm_engine.go
r('x/humanity/keeper/evm_engine.go',
  'e.syncBalancesFromDB(sdb)',
  '// P0-3: Go-State is authoritative. Removed syncBalancesFromDB — it overwrote\n// correct Go-state with stale EVM-memory values causing balance divergence.',
  'P0-3 syncBalancesFromDB')

# P0-4 register.go
c = open('x/humanity/keeper/register.go', 'rb').read().decode('utf-8')
old4 = ('if regErr := a.state.RegisterHuman(wallet); regErr != nil {\n'
        'fmt.Printf("[REGISTER] Warning: native balance grant failed (contract registration still succeeded): %v\\n", regErr)')
new4 = ('if regErr := a.state.RegisterHuman(wallet); regErr != nil {\n'
        '// P0-4: EVM succeeded but Go-State failed — retry to prevent permanent divergence\n'
        'registered := false\n'
        'for retry := 1; retry <= 3; retry++ {\n'
        'time.Sleep(time.Duration(retry) * 500 * time.Millisecond)\n'
        'if err2 := a.state.RegisterHuman(wallet); err2 == nil { registered = true; break }\n'
        '}\n'
        'if !registered {\n'
        'Log.Error("CRITICAL: RegisterHuman failed 3x after EVM success", "wallet", wallet, "error", regErr)\n'
        '}')
if old4 in c:
    c = c.replace(old4, new4, 1)
    open('x/humanity/keeper/register.go', 'wb').write(c.encode('utf-8'))
    sys.stderr.write('OK: P0-4 RegisterHuman retry\n')
else:
    sys.stderr.write('MISS: P0-4\n')
    idx = c.find('RegisterHuman')
    sys.stderr.write('  RegisterHuman at: ' + str(idx) + '\n')

# P0-5 UBI zero-share
r('x/humanity/keeper/state.go',
  'total := poolAcc.Balance.Float()\nshare := total / float64(len(humanAddrs))',
  'total := poolAcc.Balance.Float()\nshare := total / float64(len(humanAddrs))\n// P0-5: prevent funds vanishing via rounding\nif round6(share) == 0 { fmt.Printf("[UBI] Share %.10f too small — pool left intact\\n", share); return }',
  'P0-5 UBI zero-share')

# P0-5 validators
r('x/humanity/keeper/state.go',
  'total := poolAcc.Balance.Float()\n// P0-2: credit recipients BEFORE zeroing pool — crash-safe',
  'total := poolAcc.Balance.Float()\n// P0-5: prevent zero-share destroying funds\nif len(nodeShares) > 0 && round6(total/float64(len(nodeShares))) == 0 {\nfmt.Printf("[VALIDATORS] Share rounds to 0 — pool left intact\\n")\nreturn\n}\n// P0-2: credit recipients BEFORE zeroing pool — crash-safe',
  'P0-5 validators zero-share')

# P1-2 api.go
r('x/humanity/keeper/api.go',
  'if secretOK || keyAuthorizedEarly {\nGlobalPeerRegistry.Register(req.URL)',
  '// P1-2: URL registration requires PEER_SECRET or challenge-response sig. Knowing\n// a validator address from /api/blocks is NOT enough — prevents impersonation.\nif secretOK || sigOK {\nGlobalPeerRegistry.Register(req.URL)',
  'P1-2 URL auth sigOK')

# P1-5 demurrage synchronous flag
r('x/humanity/keeper/state.go',
  '// P3-6: async write so GET request is not blocked by a DB write.\ngo func(addr string) {\ncs.mu.Lock(); defer cs.mu.Unlock()\nif a, ok := cs.accounts[addr]; ok && !a.Demurrage14DayWarningShown {\na.Demurrage14DayWarningShown = true\ncs.saveAccountToDB(a)\n}\n}(address)',
  '// P1-5: set flag synchronously in memory to prevent duplicate notices.\nacc.Demurrage14DayWarningShown = true\ngo func(addr string) {\ncs.mu.Lock(); defer cs.mu.Unlock()\nif a, ok := cs.accounts[addr]; ok { cs.saveAccountToDB(a) }\n}(address)',
  'P1-5 demurrage sync')

# P2-1 evm_storage.go
c = open('x/humanity/keeper/evm_storage.go', 'rb').read().decode('utf-8')
# Find the bad truncation
old21 = 'big.NewInt(int64(acc.Balance.Float()))'
if old21 in c:
    c = c.replace(old21, 'big.NewInt(int64(acc.Balance))', -1)
    open('x/humanity/keeper/evm_storage.go', 'wb').write(c.encode('utf-8'))
    sys.stderr.write('OK: P2-1 Balance.Float() truncation\n')
else:
    sys.stderr.write('MISS: P2-1\n')

# P2-2 api.go handleWealthCap
c = open('x/humanity/keeper/api.go', 'rb').read().decode('utf-8')
if 'GetWealthCapInfo' not in c:
    # Try to find and replace the old implementation
    old22 = ('accs := a.state.GetAllAccounts()\nvar total float64\nn := 0\n'
             'for _, acc := range accs {\nif acc.IsHuman { total += acc.Balance.Float(); n++ }\n}\n'
             'avg := 0.0\nif n > 0 { avg = total / float64(n) }\n'
             'mult := 5.0\nif n > 5 { mult = float64(n) }\nif mult > 25 { mult = 25 }\n'
             'capAEQ := mult * avg')
    new22 = ('// P2-2: use canonical GetWealthCapInfo = bootstrapMultiplierLocked * getAverageBalanceLocked\ncapAEQ, mult, avg, n := a.state.GetWealthCapInfo()')
    if old22 in c:
        c = c.replace(old22, new22, 1)
        # Also fix the json encode
        c = c.replace(
            '"cap_aeq": capAEQ, "multiplier": mult, "average_aeq": avg,\n"humans": n, "total_supply": total,',
            '"cap_aeq": capAEQ, "multiplier": mult, "average_aeq": avg, "humans": n,',
            1
        )
        open('x/humanity/keeper/api.go', 'wb').write(c.encode('utf-8'))
        sys.stderr.write('OK: P2-2 handleWealthCap\n')
    else:
        sys.stderr.write('MISS: P2-2\n')
else:
    sys.stderr.write('SKIP P2-2: GetWealthCapInfo already used\n')

# P2-11 genesis block timestamp
r('x/humanity/keeper/block.go',
  'Timestamp:    time.Date(2026, 6, 13, 0, 0, 0, 0, time.UTC).Unix(),',
  'Timestamp:    genesisTimestamp(), // P2-11: reads from genesis.json when available',
  'P2-11 genesis timestamp')

# P3-2 warnedUnknownProposers cap
r('x/humanity/keeper/block.go',
  'if !dag.warnedUnknownProposers[proposer] {\ndag.warnedUnknownProposers[proposer] = true',
  'if len(dag.warnedUnknownProposers) > 500 { dag.warnedUnknownProposers = make(map[string]bool) }\nif !dag.warnedUnknownProposers[proposer] {\ndag.warnedUnknownProposers[proposer] = true',
  'P3-2 warnedUnknownProposers cap')

# P3-6 p2p.go
r('x/humanity/keeper/p2p.go',
  'buf := make([]byte, 1024)\nn, err := s.Read(buf)\nif err != nil {\nreturn\n}\ndata := buf[:n]',
  'data, err := io.ReadAll(io.LimitReader(s, 64*1024)) // P3-6: 64KB max\nif err != nil || len(data) == 0 {\nreturn\n}',
  'P3-6 p2p LimitReader')

# P3-10 decimal.go overflow guard
r('x/humanity/keeper/decimal.go',
  'func (d Decimal) MulFloat(f float64) Decimal { return Decimal(math.Round(float64(d) * f)) }',
  'func (d Decimal) MulFloat(f float64) Decimal {\n// P3-10: overflow guard\nresult := math.Round(float64(d) * f)\nconst maxD = float64(math.MaxInt64)\nif result > maxD { return Decimal(math.MaxInt64) }\nif result < -maxD { return Decimal(math.MinInt64) }\nreturn Decimal(result)\n}',
  'P3-10 MulFloat overflow guard')

sys.stderr.write('Done.\n')
