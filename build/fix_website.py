import re, sys

content = open('x/humanity/keeper/api_html.go', 'rb').read().decode('utf-8')
orig_len = len(content)
fixes = []

# 1. Gini target 0.35 -> 0.30
for old, new in [
    ('&lt;0.35', '&lt;0.30'),
    ('below 0.35', 'below 0.30'),
    ('< 0.35', '< 0.30'),
    ('Gini below 0.35', 'Gini below 0.30'),
    ('gini < 0.35', 'gini < 0.30'),
    ('gini<0.35', 'gini<0.30'),
    ('Gini target &lt;0.35', 'Gini target &lt;0.30'),
    ('Gini target <0.35', 'Gini target <0.30'),
    ('Gini below 0.30 — comparable to the most equal developed economies', 'Gini below 0.30 — comparable to Scandinavian countries'),
    ('Gini below 0.30 — comparable to the most equal developed nations', 'Gini below 0.30 — comparable to Scandinavian countries'),
    ('Gini below 0.30 — comparable to Scandinavia and Germany', 'Gini below 0.30 — comparable to Scandinavian countries'),
    ('target: Gini below 0.30 — comparable to Scandinavian countries, enforced by the wealth cap at scale (bootstrap: sliding cap 5',
     'target: Gini below 0.30 — comparable to Scandinavian countries, enforced by wealth cap (bootstrap: 5'),
]:
    count = content.count(old)
    if count > 0:
        content = content.replace(old, new)
        fixes.append('Gini 0.35->0.30 x{}: {}'.format(count, old[:30]))

# 2. Gini history chart 0-1 scale
replacements = [
    ('var g0 = history[0].idx; // 0-100',
     'var g0 = history[0].gini || (history[0].idx/100); // 0-1 scale'),
    ("[[0,'0'],[35,'35'],[70,'70'],[100,'100']].forEach",
     "[[0,'0'],[0.30,'0.30'],[0.60,'0.60'],[1,'1.0']].forEach"),
    ("var zones=[[0,35,'rgba(0,255,100,0.5)'],[35,70,'rgba(245,158,11,0.5)'],[70,100,'rgba(239,68,68,0.5)']]",
     "var zones=[[0,0.30,'rgba(0,255,100,0.5)'],[0.30,0.60,'rgba(245,158,11,0.5)'],[0.60,1.0,'rgba(239,68,68,0.5)']]"),
    ('var fill=bw*g0/100;',
     'var fill=bw*g0/1.0;'),
    ('var tx=bx+bw*35/100;',
     'var tx=bx+bw*0.30/1.0;'),
    ('var px=bx+bw*g0/100;',
     'var px=bx+bw*g0/1.0;'),
    ("ctx.fillText('Gini: ' + g0.toFixed(2), W/2, by-26);",
     "ctx.fillText('Gini: ' + g0.toFixed(4), W/2, by-26);"),
    ("if(g0<35) label='Below target — excellent equality';\n      else if(g0<70)",
     "if(g0<0.30) label='Below target — excellent equality';\n      else if(g0<0.60)"),
    ('var vs=history.map(function(h){return h.idx;});',
     'var vs=history.map(function(h){return h.gini||(h.idx/100);});'),
    ('var lo=0, hi=100;',
     'var lo=0, hi=1.0;'),
]
for old, new in replacements:
    if old in content:
        content = content.replace(old, new, 1)
        fixes.append('Chart fix: {}'.format(old[:40]))

# 3. Swap chart: filter negative prices
old_prices = 'const prices = pts.map(function(p){return p.p;});'
new_prices = 'pts = pts.filter(function(p){return p.p>0;}); const prices = pts.map(function(p){return p.p;});'
if old_prices in content:
    content = content.replace(old_prices, new_prices, 1)
    fixes.append('Swap chart negative filter')

# 4. Gini score display 4 decimal places
old_gs = "document.getElementById('idx-gini').textContent = typeof d.gini === 'number' ? d.gini.toFixed(3) : '—';"
new_gs = "document.getElementById('idx-gini').textContent = typeof d.gini === 'number' ? d.gini.toFixed(4) : '—';"
if old_gs in content:
    content = content.replace(old_gs, new_gs, 1)
    fixes.append('Gini 4dp score')

# 5. PEER_NODES -> PRIMARY_NODE_URL
peer_fixes = [
    ('PEER_NODES=https://aequitas.digital', 'PRIMARY_NODE_URL=https://aequitas.digital'),
    ("Set in your environment: <span style=\"color:var(--purple);font-family:var(--font-mono)\">PEER_NODES=https://aequitas.digital</span>",
     "Set in your environment: <span style=\"color:var(--purple);font-family:var(--font-mono)\">PRIMARY_NODE_URL=https://aequitas.digital</span>"),
]
for old, new in peer_fixes:
    if old in content:
        content = content.replace(old, new)
        fixes.append('PEER_NODES fix: {}'.format(old[:40]))

# Also fix in troubleshooting arrays in all languages
content = content.replace(
    "'Atur PEER_NODES=https://aequitas.digital dan deploy ulang'",
    "'Atur PRIMARY_NODE_URL=https://aequitas.digital dan deploy ulang'")
content = content.replace("'PEER_NODES tidak diatur'", "'PRIMARY_NODE_URL tidak diatur'")
content = content.replace("'PEER_NODES non configurato'", "'PRIMARY_NODE_URL non configurata'")
content = content.replace("'PEER_NODES non defini'", "'PRIMARY_NODE_URL non definie'")
content = content.replace("'PEER_NODES nao definido'", "'PRIMARY_NODE_URL nao definida'")
content = content.replace("'PEER_NODES non configure'", "'PRIMARY_NODE_URL non configuree'")

# 6. Node wallet fatal error fix
old_wallet = ('<strong style="color:var(--text)">A dedicated node wallet:</strong> Your node needs its own Ethereum wallet '
              'to sign transactions. This is NOT your personal AEQ wallet. Install MetaMask (metamask.io), create a new '
              'account specifically for your node, then export its private key: MetaMask &rarr; click account icon &rarr; '
              'Account Details &rarr; Show Private Key &rarr; enter your MetaMask password &rarr; copy. Keep this key '
              'strictly private &mdash; treat it like a password.')
new_wallet = ('<strong style="color:var(--text)">Node signing key (RELAYER_PRIVATE_KEY):</strong> Your node needs a dedicated '
              'Ethereum wallet to sign on-chain registrations. This can be any MetaMask wallet. Export its private key: '
              'MetaMask &rarr; Account Details &rarr; Show Private Key &rarr; enter password &rarr; copy. Keep strictly private. '
              '<strong style="color:var(--gold)">IMPORTANT:</strong> To receive validator rewards you also need '
              'NODE_OPERATOR_WALLET set to your <strong style="color:var(--neon)">registered Aequitas human wallet</strong> '
              '(the one verified with AequitasBio). Only verified humans can earn validator rewards.')
if old_wallet in content:
    content = content.replace(old_wallet, new_wallet, 1)
    fixes.append('Node wallet doc fix')
else:
    fixes.append('Node wallet doc: NOT FOUND (check manually)')

# 7. Security warning fix
old_sec = 'Use a dedicated wallet just for the node &mdash; not your personal AEQ wallet.'
new_sec = ('Use a separate MetaMask wallet for RELAYER_PRIVATE_KEY (signing). '
           'NODE_OPERATOR_WALLET (for rewards) must be your registered Aequitas human wallet.')
if old_sec in content:
    content = content.replace(old_sec, new_sec, 1)
    fixes.append('Security warning fix')

# 8. Mobile CSS
old_mobile = "@media(max-width:700px){.ns{grid-template-columns:1fr}}"
new_mobile = ("@media(max-width:700px){.ns{grid-template-columns:1fr}}\n"
              "@media(max-width:600px){"
              ".stats-grid{grid-template-columns:1fr 1fr}"
              "canvas{max-width:100%!important}"
              ".tab{padding:10px 8px;font-size:0.58rem}"
              ".rhero{padding:12px 14px 0}"
              ".nc{padding:14px}"
              ".rs{padding:12px}"
              "}")
if old_mobile in content:
    content = content.replace(old_mobile, new_mobile, 1)
    fixes.append('Mobile CSS')

open('x/humanity/keeper/api_html.go', 'wb').write(content.encode('utf-8'))
print('File: {} -> {} chars'.format(orig_len, len(content)))
print('Fixes applied:')
for f in fixes:
    print(' -', f)
