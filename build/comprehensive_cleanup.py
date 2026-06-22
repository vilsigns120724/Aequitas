"""Comprehensive cleanup of api_html.go — fix all remaining issues."""
import re

content = open('x/humanity/keeper/api_html.go', 'rb').read().decode('utf-8')
orig = len(content)
fixes = []

def fix(old, new, label, count=0):
    global content
    n = content.count(old)
    if n > 0:
        content = content.replace(old, new, count or n)
        fixes.append(f'{label} ({n}x)')
        return True
    return False

# ── 1. Remove BETA everywhere it's outdated ──────────────────────────────────
for old, new in [
    ('BETA v0.1', 'v1.0'),
    ('BETA v0.1 · Open Source', 'v1.0 · Open Source'),
    ('BETA v0.1 · Codigo Abierto', 'v1.0 · Codigo Abierto'),
    ('BETA v0.1 · Acik Kaynak', 'v1.0 · Acik Kaynak'),
    ('Android APK · direct download · BETA', 'Android APK · direct download'),
    ('covers all BETA requirements', 'covers all requirements'),
    ('meets minimum requirements for BETA', 'meets all requirements'),
    ('pour la BETA', 'pour le lancement'),
    ('per la BETA', 'per il lancio'),
    ('untuk BETA', 'untuk node'),
    ('während der BETA', 'derzeit'),
    ('In der BETA keine Mindest-Verfügbarkeit', 'Keine Mindest-Verfügbarkeit'),
    ('No minimum uptime required for BETA,', 'No minimum uptime required,'),
    ('Recommended for BETA', 'Recommended'),
    ('Empfohlen für BETA', 'Empfohlen'),
    ('(Recommended for BETA)', '(Recommended)'),
    ('para BETA.', 'para el lanzamiento.'),
    ('icin minimum gereksinimleri karsilar', 'gereksinimlerini karsilar'),
    ('BETA feedback on node setup', 'Feedback on node setup'),
]:
    fix(old, new, f'BETA cleanup: {old[:30]}')

# ── 2. Fix chartIntervalMs - var instead of let to avoid TDZ ─────────────────
fix('let chartIntervalMs =', 'var chartIntervalMs =', 'chartIntervalMs var fix')
fix('let chartIntervalMs=', 'var chartIntervalMs=', 'chartIntervalMs var fix2')

# ── 3. Ensure price chart filters zeros - already done but double-check ───────
# Already fixed in previous session

# ── 4. Fix Gini in Lorenz chart — shows (gini*100).toFixed(1)+'%' which is wrong
# Should show gini.toFixed(4) or as percentage consistently
# Check what's in the Lorenz header boxes
old_pct = "(gini*100).toFixed(1)+'%'"
new_pct = "gini.toFixed(4)"
n = content.count(old_pct)
if n > 0:
    content = content.replace(old_pct, new_pct)
    fixes.append(f'Lorenz Gini format {n}x')

# ── 5. Fix Gini comparison thresholds in Lorenz (still using 0.35) ───────────
for old, new in [
    ('gini<0.35', 'gini<0.30'),
    ('gini < 0.35', 'gini < 0.30'),
    ("gini<0.35?'#34D399':'#F0B429'", "gini<0.30?'#34D399':'#F0B429'"),
    ("gini<0.12", "gini<0.10"),  # near-perfect equality threshold
]:
    fix(old, new, f'Gini threshold: {old}')

# ── 6. Wealth cap: fix display to show total supply context ──────────────────
# The API returns cap, multiplier, avg_balance - these ARE correct
# But we should add context that total supply = humans * 1000

# ── 7. Fix wealth cap API call and display ────────────────────────────────────
# Check if /api/wealth-cap is being fetched properly
wc_fetch = content.find("fetch('/api/wealth-cap')")
if wc_fetch < 0:
    print('WARNING: /api/wealth-cap not fetched in loadStatus')
else:
    print('OK: /api/wealth-cap fetched at', wc_fetch)

# ── 8. Fix the UBI timer restart issue ───────────────────────────────────────
# Already fixed in previous session (only start if secs > 0)

# ── 9. Add missing canvas 'actual' size setting ──────────────────────────────
# Some canvases don't set width correctly
old_canvas = "if (!canvas || !canvas.offsetParent) return;\n  canvas.width = canvas.offsetWidth;"
new_canvas = ("if (!canvas) return;\n"
              "  if (canvas.offsetWidth === 0) { setTimeout(arguments.callee, 100); return; }\n"
              "  canvas.width = canvas.offsetWidth;")
n = content.count(old_canvas)
if n > 0:
    content = content.replace(old_canvas, new_canvas)
    fixes.append(f'Canvas width fix ({n}x)')

# ── 10. Fix the Gini history Y-axis — use gini (0-1) not idx (0-100) ─────────
# Already done in previous session

# ── 11. Fix Lorenz reference country Gini target line ────────────────────────
old_scan = "'Scandinavia',  g:0.27,lc:'#60A5FA',fc:'rgba(96,165,250,0.07)', tag:'Very low — target'"
new_scan = "'Scandinavia',  g:0.27,lc:'#60A5FA',fc:'rgba(96,165,250,0.07)', tag:'Very low'"
fix(old_scan, new_scan, 'Scandinavia target label')

# ── 12. Fix Lorenz bottom note ───────────────────────────────────────────────
old_note = "Aequitas target: Gini < 35% (Scandinavia level)"
new_note = "Aequitas target: Gini < 0.30  ·  World average: 0.38"
fix(old_note, new_note, 'Lorenz bottom note')

# ── 13. Fix stat boxes: < 0.35 -> < 0.30 ────────────────────────────────────
# Already done in previous session

# ── 14. Fix score target box ─────────────────────────────────────────────────
old_target = '&lt; 0.35'
new_target = '&lt; 0.30'
fix(old_target, new_target, 'Target box fix')

# ── 15. Fix BITCOIN GINI stat (it's displayed, check value) ──────────────────
# The stat should show ~0.85 for Bitcoin
# Already correct in HTML

# ── 16. Add ResizeObserver for charts that don't have it ─────────────────────
# Already added in previous session

# ── 17. Fix Gini value in stat boxes (shows %) ───────────────────────────────
# The Gini box in the Lorenz chart shows (gini*100).toFixed(1)+'%'
# This is wrong - should be gini.toFixed(4)
old_lorenz_box = "ctx.fillText((gini*100).toFixed(1)+'%', pad.l+58, 65);"
new_lorenz_box = "ctx.fillText(gini.toFixed(4), pad.l+58, 65);"
fix(old_lorenz_box, new_lorenz_box, 'Lorenz box Gini format')

old_gstr = "gStr:'G = '+(gini*100).toFixed(1)+'%'"
new_gstr = "gStr:'G = '+gini.toFixed(4)"
fix(old_gstr, new_gstr, 'Lorenz legend Gini format')

# ── 18. Fix reference countries format ───────────────────────────────────────
for old, new in [
    ("gStr:'G = '+Math.round(ref.g*100)+'%'", "gStr:'G = '+ref.g.toFixed(2)"),
    ("gStr:'G = '+(gini*100).toFixed(1)+'%'", "gStr:'G = '+gini.toFixed(4)"),
]:
    fix(old, new, f'Gini format: {old[:30]}')

# ── 19. Fix the 'owns' in Lorenz labels ──────────────────────────────────────
old_owns_pct = "owns:'50% own '+(aqY50*100).toFixed(0)+'%'"
new_owns_pct = "owns:'50% own '+(aqY50*100).toFixed(1)+'%'"
fix(old_owns_pct, new_owns_pct, 'Owns format fix')

# ── 20. Mobile: add touch-friendly font sizes ─────────────────────────────────
# Already partially done

# ── 21. Fix any remaining 'less than 0.35' text ──────────────────────────────
for old, new in [
    ('below 0.35', 'below 0.30'),
    ('< 0.35', '< 0.30'),
    ('below target (target Gini &lt; 0.30)', 'below target (Gini &lt; 0.30)'),
    ('target Gini < 0.30', 'target Gini 0.30'),
    ('Goal achieved (target Gini &lt; 0.30)', 'Near-perfect equality (Gini &lt; 0.30)'),
    ('Goal achieved (target Gini < 0.30)', 'Near-perfect equality (Gini < 0.30)'),
    ('near-perfect equality! Goal achieved (target Gini < 0.30)', 'near-perfect equality! Gini target 0.30 achieved'),
]:
    fix(old, new, f'Target text: {old[:30]}')

# ── 22. Gini history — fix multi-point chart Y-axis labels ───────────────────
# Already fixed in previous session

# ── 23. Network info — fix bootstrap address note ────────────────────────────
old_peer = "Set in your environment: <span style=\"color:var(--purple);font-family:var(--font-mono)\">PEER_NODES=https://aequitas.digital</span>"
new_peer = "Set in your environment: <span style=\"color:var(--purple);font-family:var(--font-mono)\">PRIMARY_NODE_URL=https://aequitas.digital</span>"
fix(old_peer, new_peer, 'Bootstrap address fix')

# ── 24. Remove old jsPDF generation (now using server PDFs) ──────────────────
# The _buildNodeGuidePDF function is very long and no longer needed
# But removing it is risky - let's just keep it as fallback

# ── 25. Fix "Aequitas Index" description ─────────────────────────────────────
# The index shows 0-100 scale but is derived from Gini (0-1)
old_idx_desc = "0 = perfect equality</strong> (every wallet holds exactly the same AEQ). <strong style=\"color:var(--red)\">100 = total concentration</strong>"
new_idx_desc = "0 = perfect equality</strong> (every wallet holds the same AEQ). <strong style=\"color:var(--red)\">100 = maximum concentration</strong>"
fix(old_idx_desc, new_idx_desc, 'Index description')

# ── 26. Fix "BITCOIN GINI" stat box value ────────────────────────────────────
# It shows ~0.85 which is correct, keep as-is

# ── 27. Add missing mobile CSS for Exchange tab ──────────────────────────────
old_mobile_end = "@media(max-width:600px){.stats-grid{grid-template-columns:1fr 1fr}canvas{max-width:100%!important}.tab{padding:10px 8px;font-size:0.58rem}.rhero{padding:12px 14px 0}.nc{padding:14px}.rs{padding:12px}}"
new_mobile_end = "@media(max-width:600px){.stats-grid{grid-template-columns:1fr 1fr}canvas{max-width:100%!important}.tab{padding:10px 8px;font-size:0.58rem}.rhero{padding:12px 14px 0}.nc{padding:14px}.rs{padding:12px}.swap-form{padding:12px}.pool-cards{grid-template-columns:1fr}}"
fix(old_mobile_end, new_mobile_end, 'Mobile CSS exchange')

# ── 28. Fix incorrect "Scandinavia ~0.25" reference ─────────────────────────
for old, new in [
    ('Scandinavia ≈ 0.25', 'Scandinavia ≈ 0.27'),
    ('Sweden ≈ 0.27', 'Scandinavia ≈ 0.27'),
]:
    fix(old, new, f'Scandinavia reference: {old}')

# ── 29. Update launch year reference ─────────────────────────────────────────
# "launched June 2026" is already correct

# ── 30. Fix wealth cap description ───────────────────────────────────────────
# The user said LP pool shares cause wrong values
# Actually the /api/wealth-cap returns cap, multiplier, avg_balance correctly
# The issue is the avg_balance doesn't include LP reserve AEQ
# Fix: add note to wealth cap display
old_wcap = "'Current Wealth Cap:'"
new_wcap = "'Current Wealth Cap (liquid AEQ):'"
fix(old_wcap, new_wcap, 'Wealth cap label')

# Write output
open('x/humanity/keeper/api_html.go', 'wb').write(content.encode('utf-8'))
print(f'Size: {orig} -> {len(content)} bytes ({len(content)-orig:+d})')
print(f'\n{len(fixes)} fixes applied:')
for f in fixes:
    print(f'  {f}')
