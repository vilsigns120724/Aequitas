"""
Aequitas Node Operator Guide — Complete PDF Generator
Matches the website inline guide 100%. White background, readable dark text.
"""
from reportlab.lib.pagesizes import A4
from reportlab.lib.styles import ParagraphStyle
from reportlab.lib.units import cm
from reportlab.lib.colors import HexColor, white
from reportlab.platypus import (SimpleDocTemplate, Paragraph, Spacer, Table,
                                 TableStyle, HRFlowable, KeepTogether, PageBreak)
from reportlab.lib.enums import TA_CENTER, TA_LEFT

# Colors — dark text on white background
PURPLE  = HexColor('#5B21B6')
GOLD    = HexColor('#B45309')
TEAL    = HexColor('#0F766E')
GREEN   = HexColor('#047857')
RED     = HexColor('#B91C1C')
NAVY    = HexColor('#1E1B4B')
GRAY    = HexColor('#374151')
MUTED   = HexColor('#6B7280')
LPUR    = HexColor('#EDE9FE')  # light purple bg
LGOLD   = HexColor('#FFFBEB')  # light gold bg
LRED    = HexColor('#FEF2F2')  # light red bg
LTEAL   = HexColor('#F0FDFA')  # light teal bg
LBORDER = HexColor('#DDD6FE')
LINE    = HexColor('#E5E7EB')
ALTROW  = HexColor('#F9FAFB')

def S(name, **kw):
    d = dict(fontName='Helvetica', textColor=NAVY, leading=15, spaceAfter=4,
             fontSize=9.5)
    d.update(kw)
    return ParagraphStyle(name, **d)

STYLES = {
    'title':  S('T',  fontName='Helvetica-Bold', fontSize=20, textColor=PURPLE,
                 leading=26, spaceAfter=2, alignment=TA_CENTER),
    'sub':    S('SU', fontSize=9, textColor=MUTED, alignment=TA_CENTER, spaceAfter=4),
    'tag':    S('TG', fontSize=8, textColor=MUTED, alignment=TA_CENTER, spaceAfter=14),
    'h1':     S('H1', fontName='Helvetica-Bold', fontSize=11, textColor=PURPLE,
                 spaceBefore=16, spaceAfter=6),
    'h2':     S('H2', fontName='Helvetica-Bold', fontSize=10, textColor=GOLD,
                 spaceBefore=10, spaceAfter=4),
    'body':   S('BO', spaceAfter=6, leading=15),
    'sm':     S('SM', fontSize=8.5, textColor=GRAY, leading=13, spaceAfter=4),
    'code':   S('CO', fontName='Courier', fontSize=7.5, textColor=PURPLE,
                 backColor=LPUR, leading=11, leftIndent=8, rightIndent=8,
                 spaceAfter=8, spaceBefore=2),
    'warn':   S('WN', fontSize=8.5, textColor=RED, leading=13, spaceAfter=6,
                 leftIndent=8, fontName='Helvetica-Bold'),
    'info':   S('IN', fontSize=8.5, textColor=TEAL, leading=13, spaceAfter=6,
                 leftIndent=8),
    'bullet': S('BU', leftIndent=14, spaceAfter=3, leading=14),
    'foot':   S('FO', fontSize=7.5, textColor=MUTED, alignment=TA_CENTER, leading=11),
}

def HR():
    return HRFlowable(width='100%', thickness=0.4, color=LINE, spaceAfter=8, spaceBefore=4)

def box(text, color=TEAL, bg=LTEAL):
    d = [[Paragraph(text, S('bx', fontSize=8.5, textColor=color, leading=13))]]
    t = Table(d, colWidths=[16.6*cm])
    t.setStyle(TableStyle([
        ('BACKGROUND', (0,0), (-1,-1), bg),
        ('BOX', (0,0), (-1,-1), 0.5, color),
        ('LEFTPADDING', (0,0), (-1,-1), 10),
        ('RIGHTPADDING', (0,0), (-1,-1), 10),
        ('TOPPADDING', (0,0), (-1,-1), 6),
        ('BOTTOMPADDING', (0,0), (-1,-1), 6),
    ]))
    return t

def step_block(num, title, content_items, color=PURPLE):
    """Render a numbered step with title and content."""
    story = []
    # Number badge + title in one row
    badge = Paragraph(str(num), S('bd', fontName='Helvetica-Bold', fontSize=11,
                                   textColor=white, alignment=TA_CENTER, leading=14))
    head  = Paragraph(f'<b>{title}</b>', S('sh', fontName='Helvetica-Bold',
                                             fontSize=10, textColor=color, leading=13))
    row = Table([[badge, head]], colWidths=[0.8*cm, 15.8*cm])
    row.setStyle(TableStyle([
        ('BACKGROUND', (0,0), (0,0), color),
        ('VALIGN', (0,0), (-1,-1), 'MIDDLE'),
        ('TOPPADDING', (0,0), (0,0), 4),
        ('BOTTOMPADDING', (0,0), (0,0), 4),
        ('TOPPADDING', (0,0), (1,0), 2),
        ('BOTTOMPADDING', (0,0), (1,0), 2),
        ('LEFTPADDING', (0,0), (1,0), 8),
        ('LEFTPADDING', (0,0), (0,0), 0),
        ('RIGHTPADDING', (0,0), (-1,-1), 0),
    ]))
    story.append(row)
    for item in content_items:
        story.append(item)
    story.append(Spacer(1, 6))
    return story

def var_table(rows, cols):
    def th(t): return Paragraph(f'<b>{t}</b>', S('th', fontName='Helvetica-Bold',
                                                   fontSize=8, textColor=white, leading=11))
    def tv(t): return Paragraph(t, S('tv', fontName='Courier', fontSize=7.5,
                                      textColor=PURPLE, leading=10))
    def tr_req(t):
        c = RED if t in ('YES','JA','SI','OUI','SIM','EVET','YA') else \
            GREEN if 'reward' in t.lower() or 'Bel' in t or 'Bel.' in t else \
            GOLD if 'Rec' in t or 'Emp' in t or 'Multi' in t or 'Opt' in t else MUTED
        return Paragraph(f'<b>{t}</b>', S('tq', fontName='Helvetica-Bold', fontSize=8,
                                           textColor=c, leading=10))
    def td(t): return Paragraph(t, S('td', fontSize=8, textColor=GRAY, leading=12))

    data = [[th(c) for c in cols]]
    for r in rows:
        data.append([tv(r[0]), tr_req(r[1]), td(r[2])])
    cw = [4*cm, 2.2*cm, 10.4*cm]
    t = Table(data, colWidths=cw, repeatRows=1)
    t.setStyle(TableStyle([
        ('BACKGROUND', (0,0), (-1,0), PURPLE),
        ('ROWBACKGROUNDS', (0,1), (-1,-1), [white, ALTROW]),
        ('GRID', (0,0), (-1,-1), 0.3, LINE),
        ('VALIGN', (0,0), (-1,-1), 'TOP'),
        ('TOPPADDING', (0,0), (-1,-1), 5),
        ('BOTTOMPADDING', (0,0), (-1,-1), 5),
        ('LEFTPADDING', (0,0), (-1,-1), 6),
        ('RIGHTPADDING', (0,0), (-1,-1), 6),
    ]))
    return t

def trouble_table(rows, cols):
    def th(t): return Paragraph(f'<b>{t}</b>', S('th2', fontName='Helvetica-Bold',
                                                   fontSize=8, textColor=white, leading=11))
    def td1(t): return Paragraph(t, S('t1', fontSize=8, textColor=RED,
                                       fontName='Helvetica-Bold', leading=12))
    def td2(t): return Paragraph(t, S('t2', fontSize=8, textColor=MUTED, leading=12))
    def td3(t): return Paragraph(t, S('t3', fontSize=8, textColor=GRAY, leading=12))
    data = [[th(c) for c in cols]]
    for r in rows:
        data.append([td1(r[0]), td2(r[1]), td3(r[2])])
    t = Table(data, colWidths=[4.5*cm, 4*cm, 8.1*cm], repeatRows=1)
    t.setStyle(TableStyle([
        ('BACKGROUND', (0,0), (-1,0), PURPLE),
        ('ROWBACKGROUNDS', (0,1), (-1,-1), [white, ALTROW]),
        ('GRID', (0,0), (-1,-1), 0.3, LINE),
        ('VALIGN', (0,0), (-1,-1), 'TOP'),
        ('TOPPADDING', (0,0), (-1,-1), 5),
        ('BOTTOMPADDING', (0,0), (-1,-1), 5),
        ('LEFTPADDING', (0,0), (-1,-1), 6),
        ('RIGHTPADDING', (0,0), (-1,-1), 6),
    ]))
    return t

# ── CONTENT ───────────────────────────────────────────────────────────────────

EN = dict(
title    = 'AEQUITAS NODE OPERATOR GUIDE',
version  = 'v1.0 · June 2026 · aequitas.digital',
tagline  = 'Complete step-by-step guide · No prior blockchain experience required · ~20–30 min',

prereq_title = 'Before You Start — What You Need',
prereqs = [
    ('1.', '<b>An Aequitas account:</b> You must first be registered as a human on Aequitas. Install the Android app, complete biometric registration, and note your wallet address. Without this, you cannot receive validator rewards.'),
    ('2.', '<b>A GitHub account (free):</b> Go to github.com and create a free account. You need this to copy (fork) the Aequitas code so Railway can deploy it.'),
    ('3.', '<b>A Railway account (free):</b> Go to railway.app and sign in with GitHub. Railway is a hosting platform that runs your node in the cloud — no server or command line required.'),
    ('4.', '<b>Node signing key (RELAYER_PRIVATE_KEY):</b> Your node needs a dedicated Ethereum wallet to sign on-chain registrations. This can be any MetaMask wallet. Export its private key: MetaMask → Account Details → Show Private Key → enter password → copy. Keep strictly private. <b>IMPORTANT:</b> To receive validator rewards you also need NODE_OPERATOR_WALLET set to your <b>registered Aequitas human wallet</b> (the one verified with AequitasBio). Only verified humans can earn validator rewards.'),
    ('5.', '<b>10–30 minutes of your time.</b> Railway does most of the work automatically.'),
],

vars_title = 'Step 3 — Environment Variables',
vars_warn  = 'Security Warning: Your RELAYER_PRIVATE_KEY is like a master password. Anyone who has it controls your node wallet. Never share it publicly, never paste it in chat or email. Use a separate MetaMask wallet for RELAYER_PRIVATE_KEY (signing). NODE_OPERATOR_WALLET (for rewards) must be your registered Aequitas human wallet.',
var_cols   = ['Variable', 'Required?', 'What to set'],
vars = [
    ('DATABASE_URL',        'YES',         'Your PostgreSQL connection string. On Railway: auto-injected when PostgreSQL is in the same project. Format: postgres://user:pass@host:5432/dbname'),
    ('RELAYER_PRIVATE_KEY', 'YES',         'The private key (0x…, 66 chars) of your dedicated node wallet. MetaMask: Account Details → Show Private Key → enter password → copy.'),
    ('RELAYER_ADDRESS',     'Recommended', 'The wallet address (0x…, 42 chars) matching RELAYER_PRIVATE_KEY. Copy from MetaMask. A fallback exists but setting this explicitly prevents startup errors.'),
    ('NODE_OPERATOR_WALLET','For rewards', 'Your Aequitas human wallet address — registered via the Android app. Receives your daily validator rewards (40% of all protocol fees). Must be a registered human.'),
    ('PEER_SECRET',         'Multi-node',  'A shared secret that authorises your node as a validator. Every node in the same network must use the identical value. Get from the network operator — do not share publicly.'),
    ('SELF_URL',            'Multi-node',  'Your node\'s own public HTTPS URL (e.g. https://my-node.up.railway.app). Required for peer discovery self-exclusion. Find in Railway: Settings → Networking → Public Networking.'),
    ('PRIMARY_NODE_URL',    'Multi-node',  'Set to: https://aequitas.digital — the primary node your node registers with for automatic peer discovery. On startup your node posts its URL + signing address to the primary, gets the full peer list back, and joins the network automatically.'),
    ('PORT',                'No',          'Leave unset on Railway — Railway sets this automatically. Default is 8080.'),
    ('NODE_KEY',            'No',          'Base64 libp2p key for stable P2P identity. Auto-generated if omitted, but changes on every restart. If not set, the node prints it to stderr: "SAVE THIS AS NODE_KEY ENVIRONMENT VAR: <base64>". Copy and paste it here.'),
    ('IS_PRIMARY_NODE',     'No',          'Leave unset or false. Distribution now uses a DB-level lock — any node can run it without this variable. Setting true on a secondary node is no longer necessary.'),
    ('RESET_STATE',         'No',          'DANGEROUS: Setting this to true wipes your entire database on every restart. Development use only. Never in production.'),
],

railway_title = 'Step 4 — Deploy on Railway (Recommended)',
railway_intro = 'Railway is the easiest way to run your node — no server setup, no command line required. The free tier covers all requirements. Total time: about 10–15 minutes.',
railway_steps = [
    'In your Railway project (from Step 2), click <b>+ New → GitHub Repo</b>',
    'Select your Aequitas fork (from Step 1) — Railway detects the Dockerfile automatically',
    'Click <b>Deploy Now</b> — a first build starts (may fail without env vars, that is normal)',
    'Click your Aequitas service → <b>Variables</b> → add each variable (see table above). Minimum required: RELAYER_PRIVATE_KEY, RELAYER_ADDRESS, NODE_OPERATOR_WALLET, PEER_SECRET, SELF_URL, PRIMARY_NODE_URL=https://aequitas.digital',
    'Click <b>Deploy</b> (or save variables to trigger auto-redeploy). Build takes ~3 minutes while Go compiles the node binary.',
    'Watch <b>Deploy Logs</b>. Success looks like: <font name="Courier" color="#5B21B6">Aequitas Node Running</font> and <font name="Courier" color="#0F766E">[NODE] Registered node operator wallet: 0x…</font>',
    'Go to <b>Settings → Networking → Generate Domain</b> to get your public URL',
    'Open <font name="Courier">https://YOUR-URL/api/status</font> — you should see JSON with <b>height</b> climbing every ~6 seconds',
],
railway_vars_code = (
    '# Railway auto-sets DATABASE_URL if PostgreSQL is in the same project\n'
    'RELAYER_PRIVATE_KEY    = 0xYOUR_PRIVATE_KEY\n'
    'RELAYER_ADDRESS        = 0xYOUR_NODE_WALLET_ADDRESS\n'
    'NODE_OPERATOR_WALLET   = 0xYOUR_HUMAN_WALLET\n'
    'PEER_SECRET            = get-this-from-network-operator\n'
    'SELF_URL               = https://YOUR-RAILWAY-DOMAIN.up.railway.app\n'
    'PRIMARY_NODE_URL       = https://aequitas.digital'
),

docker_title = 'Step 4b — Alternative: Deploy with Docker (Advanced)',
docker_intro = 'Use this if you have your own server (VPS, home server, cloud VM). Requires Docker and a PostgreSQL database.',
docker_code  = (
    '# 1. Download the code\n'
    'git clone https://github.com/hanoi96international-gif/Aequitas && cd Aequitas\n\n'
    '# 2. Build the node image (~3 min for Go compilation)\n'
    'docker build -t aequitas-node .\n\n'
    '# 3. Start the node\n'
    'docker run -d --name aequitas-node --restart unless-stopped \\\n'
    '  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n'
    '  -e RELAYER_PRIVATE_KEY="0xYOUR_PRIVATE_KEY" \\\n'
    '  -e RELAYER_ADDRESS="0xYOUR_NODE_WALLET_ADDRESS" \\\n'
    '  -e NODE_OPERATOR_WALLET="0xYOUR_HUMAN_WALLET" \\\n'
    '  -e PEER_SECRET="get-from-network-operator" \\\n'
    '  -e SELF_URL="https://YOUR-PUBLIC-URL" \\\n'
    '  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n'
    '  -p 8080:8080 aequitas-node\n\n'
    '# 4. Watch the live logs\n'
    'docker logs -f aequitas-node'
),

verify_title = 'Step 5 — Verify Your Node is Running',
verify_body  = 'Open these URLs in your browser. Replace YOUR-NODE-URL with your actual Railway domain or server address.',
verify_code  = (
    'https://YOUR-NODE-URL/api/status\n'
    ' → Expected: {"height": 1234, "total_humans": N, "aequitas_index": N}\n\n'
    'https://YOUR-NODE-URL/rpc\n'
    ' → Expected: {"jsonrpc":"2.0","error":"method not specified"} — RPC is alive'
),
verify_note  = 'The block height should match the primary node within 1–2 blocks within seconds of startup. If it stays at 0, check that PRIMARY_NODE_URL=https://aequitas.digital is set and reachable.',

valkey_title = 'Step 5b — Register Your Validator Key (Decentralized Auth)',
valkey_body  = 'Instead of a shared PEER_SECRET, register your node signing key with your human wallet. This cryptographically proves you control both keys. Get the signing key signature by running this on your server (SSH/Railway shell):',
valkey_code  = 'curl "http://localhost:8080/api/sign-validator-challenge?wallet=0xYOUR_HUMAN_WALLET"',
valkey_note  = 'Then use the website Network → Run a Node tab and click "Sign with MetaMask & Register" to complete the registration.',

mm_title = 'Step 6 — Connect MetaMask to Your Node (Optional)',
mm_body  = 'In MetaMask: click the network dropdown → Add network → Add a network manually, then enter:',
mm_rows  = [
    ('Network Name',    'Aequitas Chain'),
    ('RPC URL',         'https://YOUR-NODE-URL/rpc'),
    ('Chain ID',        '1926'),
    ('Currency Symbol', 'AEQ'),
    ('Decimals',        '18'),
    ('Block Explorer',  'https://aequitas.digital'),
],

rewards_title = 'Step 7 — Earning Validator Rewards',
rewards_box   = 'The Validators Pool collects 40% of all protocol fees (swap fees, demurrage, wealth cap overflow). Every day at 20:00 Berlin time (CEST/CET, handles DST automatically) the node distributes the pool balance to all registered node operators proportionally by blocks produced. The more consistently your node runs, the larger your share.',
rewards_steps = [
    'Make sure you are registered as a human on Aequitas. If not: install the Android app and complete biometric registration first. You will receive a wallet address and 1,000 AEQ.',
    'Set <font name="Courier">NODE_OPERATOR_WALLET</font> = your Aequitas human wallet address in your Railway Variables',
    'Save — Railway redeploys automatically. On Docker: <font name="Courier">docker restart aequitas-node</font>',
    'In your node logs, confirm: <font name="Courier" color="#0F766E">[NODE] Registered node operator wallet: 0x…</font>',
    'Rewards are distributed automatically every day at 20:00 Berlin time (CEST/CET). Just keep your node running — no further action needed.',
],

trouble_title = 'Troubleshooting',
trouble_cols  = ['Symptom', 'Likely Cause', 'Solution'],
trouble_rows  = [
    ('Block height stays at 0',        'PRIMARY_NODE_URL not set or wrong',        'Set PRIMARY_NODE_URL=https://aequitas.digital and redeploy. Also set SELF_URL to your node\'s public URL.'),
    ('DATABASE_URL error on startup',  'Wrong connection string',                  'Check format: postgres://user:pass@host:5432/dbname — make sure PostgreSQL is running and accessible.'),
    ('"no code at address" in logs',   'V7 contract not yet deployed',             'Normal on first start — node auto-deploys V7. Wait a few seconds and check again.'),
    ('"NODE_OPERATOR_WALLET not set"', 'Missing environment variable',             'Add NODE_OPERATOR_WALLET=0xYOUR_HUMAN_WALLET. Node runs fine without it but you won\'t receive rewards.'),
    ('Railway "Application error"',    'Build or startup failure',                 'Check Deploy Logs. Most common: DATABASE_URL missing or RELAYER_PRIVATE_KEY in wrong format (must start with 0x).'),
    ('Port 8080 not reachable (Docker)','Firewall or cloud provider config',       'Open TCP port 8080 inbound in your firewall or cloud security group settings.'),
    ('Docker build fails (module error)','No internet during build',               'Docker build needs outbound internet to download Go modules. Railway handles this automatically.'),
],

footer = 'Aequitas Chain · Chain ID 1926 · aequitas.digital · Validator rewards: daily at 20:00 Berlin time (CEST/CET)',
)

DE = dict(
title    = 'AEQUITAS NODE-BETREIBER-ANLEITUNG',
version  = 'v1.0 · Juni 2026 · aequitas.digital',
tagline  = 'Vollstaendige Schritt-fuer-Schritt-Anleitung · Keine Vorkenntnisse noetig · ca. 20–30 Min.',

prereq_title = 'Vor dem Start — Was du brauchst',
prereqs = [
    ('1.', '<b>Ein Aequitas-Konto:</b> Du musst zuerst als Mensch auf Aequitas registriert sein. Installiere die Android-App, schliesse die biometrische Registrierung ab und notiere deine Wallet-Adresse. Ohne dies kannst du keine Validator-Belohnungen erhalten.'),
    ('2.', '<b>Ein GitHub-Konto (kostenlos):</b> Erstelle eines auf github.com. Du brauchst es um den Aequitas-Code zu forken, damit Railway ihn deployen kann.'),
    ('3.', '<b>Ein Railway-Konto (kostenlos):</b> Gehe zu railway.app und melde dich mit GitHub an. Railway ist eine Hosting-Plattform die deinen Node in der Cloud betreibt — kein eigener Server oder Terminal erforderlich.'),
    ('4.', '<b>Node Signing-Key (RELAYER_PRIVATE_KEY):</b> Dein Node braucht eine dedizierte Ethereum-Wallet zum Signieren. Das kann jede MetaMask-Wallet sein. Exportiere den privaten Schluessel: MetaMask → Kontodetails → Privaten Schluessel anzeigen → Passwort eingeben → kopieren. Streng geheimhalten. <b>WICHTIG:</b> Um Validator-Belohnungen zu erhalten, muss NODE_OPERATOR_WALLET deine <b>registrierte Aequitas-Mensch-Wallet</b> sein (die mit AequitasBio verifizierte). Nur verifizierte Menschen koennen Validator-Belohnungen verdienen.'),
    ('5.', '<b>10–30 Minuten deiner Zeit.</b> Railway erledigt den Grossteil automatisch.'),
],

vars_title = 'Schritt 3 — Umgebungsvariablen',
vars_warn  = 'Sicherheitswarnung: Dein RELAYER_PRIVATE_KEY ist wie ein Master-Passwort. Wer ihn hat, kontrolliert deine Node-Wallet. Niemals oeffentlich teilen, niemals in Chat oder E-Mail einfuegen. Verwende fuer RELAYER_PRIVATE_KEY eine separate Wallet. NODE_OPERATOR_WALLET (fuer Belohnungen) muss deine registrierte Aequitas-Mensch-Wallet sein.',
var_cols   = ['Variable', 'Erforderlich?', 'Was eintragen'],
vars = [
    ('DATABASE_URL',        'JA',           'Dein PostgreSQL-Verbindungsstring. Auf Railway: automatisch gesetzt wenn PostgreSQL im gleichen Projekt. Format: postgres://user:pass@host:5432/dbname'),
    ('RELAYER_PRIVATE_KEY', 'JA',           'Privater Schluessel deiner Node-Wallet (0x…, 66 Zeichen). MetaMask: Kontodetails → Privaten Schluessel anzeigen → Passwort → kopieren.'),
    ('RELAYER_ADDRESS',     'Empfohlen',    'Wallet-Adresse (0x…, 42 Zeichen) passend zu RELAYER_PRIVATE_KEY. Aus MetaMask kopieren. Verhindert Startfehler.'),
    ('NODE_OPERATOR_WALLET','Fuer Bel.',    'Deine Aequitas-Mensch-Wallet — die via Android-App registrierte. Erhaelt taeglich Validator-Belohnungen (40% aller Protokollgebuehren). Muss ein registrierter Mensch sein.'),
    ('PEER_SECRET',         'Multi-Node',   'Gemeinsames Geheimnis das deinen Node als Validator autorisiert. Alle Nodes im Netzwerk muessen identischen Wert nutzen. Vom Netzwerkbetreiber erhalten — nicht teilen.'),
    ('SELF_URL',            'Multi-Node',   'Eigene oeffentliche HTTPS-URL des Nodes (z.B. https://mein-node.up.railway.app). In Railway: Settings → Networking → Public Networking.'),
    ('PRIMARY_NODE_URL',    'Multi-Node',   'Auf https://aequitas.digital setzen — der Primaer-Node bei dem sich dein Node registriert. Beim Start postet der Node URL + Signing-Adresse und bekommt die Peer-Liste zurueck.'),
    ('PORT',                'Nein',         'Auf Railway nicht setzen — wird automatisch gesetzt. Standard ist 8080.'),
    ('NODE_KEY',            'Nein',         'Base64 libp2p-Schluessel fuer stabile Peer-Identitaet. Auto-generiert wenn nicht gesetzt, aendert sich dann bei jedem Neustart. Beim ersten Start in stderr ausgegeben: "SAVE THIS AS NODE_KEY: <base64>". Kopieren und hier setzen.'),
    ('IS_PRIMARY_NODE',     'Nein',         'Nicht setzen oder false lassen. Die Ausschuettung nutzt jetzt einen DB-Lock — jeder Node kann sie ohne diese Variable ausfuehren.'),
    ('RESET_STATE',         'Nein',         'GEFAEHRLICH: True loescht die gesamte DB bei jedem Neustart. Nur fuer Entwicklung. Niemals in Produktion.'),
],

railway_title = 'Schritt 4 — Deployment auf Railway (Empfohlen)',
railway_intro = 'Railway ist der einfachste Weg deinen Node zu betreiben — kein Server-Setup, kein Terminal erforderlich. Der kostenlose Tarif deckt alle Anforderungen. Gesamtzeit: ca. 10–15 Minuten.',
railway_steps = [
    'In deinem Railway-Projekt (aus Schritt 2): <b>+ New → GitHub Repo</b> klicken',
    'Deinen Aequitas-Fork auswaehlen (aus Schritt 1) — Railway erkennt das Dockerfile automatisch',
    '<b>Deploy Now</b> klicken — ein erster Build startet (kann ohne Env Vars fehlschlagen, das ist normal)',
    'Aequitas-Service → <b>Variables</b> → Variablen hinzufuegen (siehe Tabelle oben). Mindest-Anforderung: RELAYER_PRIVATE_KEY, RELAYER_ADDRESS, NODE_OPERATOR_WALLET, PEER_SECRET, SELF_URL, PRIMARY_NODE_URL=https://aequitas.digital',
    '<b>Deploy</b> klicken (oder Variablen speichern fuer Auto-Redeploy). Build dauert ~3 Minuten fuer Go-Kompilierung.',
    'Deploy-Logs beobachten. Erfolg sieht so aus: <font name="Courier" color="#5B21B6">Aequitas Node Running</font> und <font name="Courier" color="#0F766E">[NODE] Registered node operator wallet: 0x…</font>',
    '<b>Settings → Networking → Generate Domain</b> fuer deine oeffentliche URL',
    '<font name="Courier">https://DEINE-URL/api/status</font> aufrufen — du siehst JSON mit <b>height</b> der alle ~6 Sekunden steigt',
],
railway_vars_code = (
    '# Railway setzt DATABASE_URL automatisch wenn PostgreSQL im gleichen Projekt\n'
    'RELAYER_PRIVATE_KEY    = 0xDEIN_PRIVATER_SCHLUESSEL\n'
    'RELAYER_ADDRESS        = 0xDEINE_NODE_WALLET_ADRESSE\n'
    'NODE_OPERATOR_WALLET   = 0xDEINE_MENSCH_WALLET\n'
    'PEER_SECRET            = vom-Netzwerkbetreiber-erhalten\n'
    'SELF_URL               = https://DEIN-RAILWAY-DOMAIN.up.railway.app\n'
    'PRIMARY_NODE_URL       = https://aequitas.digital'
),

docker_title = 'Schritt 4b — Alternative: Docker-Deployment (Fortgeschritten)',
docker_intro = 'Nutze dies wenn du einen eigenen Server hast (VPS, Heimserver, Cloud-VM). Erfordert Docker und eine PostgreSQL-Datenbank.',
docker_code  = (
    '# 1. Code herunterladen\n'
    'git clone https://github.com/hanoi96international-gif/Aequitas && cd Aequitas\n\n'
    '# 2. Node-Image erstellen (~3 Min fuer Go-Kompilierung)\n'
    'docker build -t aequitas-node .\n\n'
    '# 3. Node starten — alle Platzhalter ersetzen\n'
    'docker run -d --name aequitas-node --restart unless-stopped \\\n'
    '  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n'
    '  -e RELAYER_PRIVATE_KEY="0xDEIN_PRIVATER_SCHLUESSEL" \\\n'
    '  -e RELAYER_ADDRESS="0xDEINE_NODE_WALLET_ADRESSE" \\\n'
    '  -e NODE_OPERATOR_WALLET="0xDEINE_MENSCH_WALLET" \\\n'
    '  -e PEER_SECRET="vom-Netzwerkbetreiber" \\\n'
    '  -e SELF_URL="https://DEINE-OEFFENTLICHE-URL" \\\n'
    '  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n'
    '  -p 8080:8080 aequitas-node\n\n'
    '# 4. Live-Logs beobachten\n'
    'docker logs -f aequitas-node'
),

verify_title = 'Schritt 5 — Node-Betrieb pruefen',
verify_body  = 'Oeffne diese URLs im Browser. Ersetze DEINE-NODE-URL durch deine Railway-Domain oder Server-Adresse.',
verify_code  = (
    'https://DEINE-NODE-URL/api/status\n'
    ' → Erwartet: {"height": 1234, "total_humans": N, "aequitas_index": N}\n\n'
    'https://DEINE-NODE-URL/rpc\n'
    ' → Erwartet: {"jsonrpc":"2.0","error":"method not specified"} — RPC laeuft'
),
verify_note  = 'Die Blockhoehe sollte innerhalb von Sekunden mit dem Primaer-Node uebereinstimmen (1–2 Bloecke). Bleibt sie bei 0: PRIMARY_NODE_URL=https://aequitas.digital pruefen.',

valkey_title = 'Schritt 5b — Validator-Schluessel registrieren (Dezentrale Auth)',
valkey_body  = 'Statt eines gemeinsamen PEER_SECRET kannst du deinen Node-Signing-Key mit deiner Mensch-Wallet registrieren. Fuhre diesen Befehl auf deinem Server aus (SSH/Railway Shell):',
valkey_code  = 'curl "http://localhost:8080/api/sign-validator-challenge?wallet=0xDEINE_MENSCH_WALLET"',
valkey_note  = 'Dann auf der Website unter Network → Run a Node den Button "Sign with MetaMask & Register" nutzen um die Registrierung abzuschliessen.',

mm_title = 'Schritt 6 — MetaMask mit deinem Node verbinden (Optional)',
mm_body  = 'In MetaMask: Netzwerk-Dropdown → Netzwerk hinzufuegen → Netzwerk manuell hinzufuegen:',
mm_rows  = [
    ('Network Name',    'Aequitas Chain'),
    ('RPC URL',         'https://DEINE-NODE-URL/rpc'),
    ('Chain ID',        '1926'),
    ('Currency Symbol', 'AEQ'),
    ('Decimals',        '18'),
    ('Block Explorer',  'https://aequitas.digital'),
],

rewards_title = 'Schritt 7 — Validator-Belohnungen erhalten',
rewards_box   = 'Der Validators-Pool sammelt 40% aller Protokollgebuehren (Swap-Gebuehren, Demurrage, Wealth-Cap-Ueberschuss). Jeden Tag um 20:00 Uhr Berliner Zeit (CEST/CET, DST automatisch) verteilt der Node den Pool-Saldo proportional nach produzierten Bloecken an alle registrierten Node-Betreiber. Je laenger dein Node laeuft, desto groesser dein Anteil.',
rewards_steps = [
    'Stelle sicher, dass du als Mensch auf Aequitas registriert bist. Falls nicht: Android-App installieren und biometrische Registrierung abschliessen. Du erhaeltst eine Wallet-Adresse und 1.000 AEQ.',
    '<font name="Courier">NODE_OPERATOR_WALLET</font> = deine Aequitas-Mensch-Wallet-Adresse in Railway Variables setzen',
    'Speichern — Railway redeployt automatisch. Mit Docker: <font name="Courier">docker restart aequitas-node</font>',
    'In den Node-Logs bestaetigen: <font name="Courier" color="#0F766E">[NODE] Registered node operator wallet: 0x…</font>',
    'Belohnungen werden automatisch jeden Tag um 20:00 Uhr Berliner Zeit (CEST/CET) verteilt. Node laufen lassen — kein weiterer Eingriff noetig.',
],

trouble_title = 'Fehlerbehebung',
trouble_cols  = ['Symptom', 'Wahrscheinliche Ursache', 'Loesung'],
trouble_rows  = [
    ('Blockhoehe bleibt bei 0',         'PRIMARY_NODE_URL nicht gesetzt',             'PRIMARY_NODE_URL=https://aequitas.digital setzen und neu deployen. SELF_URL auf Node-URL setzen.'),
    ('DATABASE_URL-Fehler beim Start',   'Falscher Connection-String',                'Format pruefen: postgres://user:pass@host:5432/dbname — PostgreSQL muss erreichbar sein.'),
    ('"no code at address" in Logs',     'V7-Contract noch nicht deployed',            'Normal beim ersten Start — Node deployed V7 automatisch. Kurz warten.'),
    ('"NODE_OPERATOR_WALLET not set"',   'Fehlende Umgebungsvariable',                'NODE_OPERATOR_WALLET=0xDEINE_MENSCH_WALLET hinzufuegen. Node laeuft ohne, aber keine Belohnungen.'),
    ('Railway "Application error"',      'Build- oder Startfehler',                   'Deploy-Logs pruefen. Haeufigste Ursache: fehlende DATABASE_URL oder falsches Schluessel-Format.'),
    ('Port 8080 nicht erreichbar (Docker)','Firewall oder Cloud-Konfiguration',        'TCP-Port 8080 eingehend in Firewall oder Cloud-Security-Gruppe oeffnen.'),
    ('Docker Build scheitert (Module)',   'Kein Internet beim Build',                  'Docker Build benoetigt ausgehenden Internetzugang. Railway erledigt das automatisch.'),
],

footer = 'Aequitas Chain · Chain ID 1926 · aequitas.digital · Validator-Belohnungen: taeglich 20:00 Uhr Berliner Zeit (CEST/CET)',
)

# ── PDF BUILDER ───────────────────────────────────────────────────────────────

def build_pdf(path, L):
    doc = SimpleDocTemplate(path, pagesize=A4,
                            leftMargin=2*cm, rightMargin=2*cm,
                            topMargin=2*cm, bottomMargin=2*cm,
                            title=L['title'])
    fl = []

    # HEADER
    fl += [Paragraph(L['title'],   STYLES['title']),
           Paragraph(L['version'], STYLES['sub']),
           Paragraph(L['tagline'], STYLES['tag']),
           HR()]

    # BEFORE YOU START
    fl += [Paragraph(L['prereq_title'], STYLES['h1'])]
    for num, text in L['prereqs']:
        row = Table([[Paragraph(f'<b><font color="#5B21B6">{num}</font></b>',
                                 S('pn', fontName='Helvetica-Bold', fontSize=10,
                                   textColor=PURPLE, alignment=TA_CENTER, leading=14)),
                      Paragraph(text, STYLES['body'])]],
                    colWidths=[0.7*cm, 15.9*cm])
        row.setStyle(TableStyle([('VALIGN',(0,0),(-1,-1),'TOP'),
                                  ('TOPPADDING',(0,0),(-1,-1),2),
                                  ('BOTTOMPADDING',(0,0),(-1,-1),2),
                                  ('LEFTPADDING',(0,0),(1,0),8),
                                  ('LEFTPADDING',(0,0),(0,0),0)]))
        fl.append(row)
    fl.append(HR())

    # STEP 3 VARS
    fl += [Paragraph(L['vars_title'], STYLES['h1']),
           box(L['vars_warn'], RED, LRED),
           Spacer(1,8),
           var_table(L['vars'], L['var_cols']),
           HR()]

    # STEP 4 RAILWAY
    fl += [Paragraph(L['railway_title'], STYLES['h1']),
           Paragraph(L['railway_intro'], STYLES['body'])]
    for i, step in enumerate(L['railway_steps'], 1):
        color = GOLD if i in (4, 7) else PURPLE
        row = Table([[Paragraph(str(i), S('sn', fontName='Helvetica-Bold', fontSize=10,
                                           textColor=white, alignment=TA_CENTER, leading=13)),
                      Paragraph(step, STYLES['body'])]],
                    colWidths=[0.7*cm, 15.9*cm])
        row.setStyle(TableStyle([('BACKGROUND',(0,0),(0,0),color),
                                  ('VALIGN',(0,0),(-1,-1),'TOP'),
                                  ('TOPPADDING',(0,0),(0,0),3),
                                  ('BOTTOMPADDING',(0,0),(0,0),3),
                                  ('TOPPADDING',(0,0),(1,0),1),
                                  ('BOTTOMPADDING',(0,0),(1,0),4),
                                  ('LEFTPADDING',(0,0),(0,0),0),
                                  ('LEFTPADDING',(0,0),(1,0),8)]))
        fl.append(row)
    fl += [Spacer(1,6),
           Paragraph(L['railway_vars_code'].replace('\n','<br/>').replace(' ','&nbsp;'),
                     STYLES['code']),
           HR()]

    # STEP 4B DOCKER
    fl += [Paragraph(L['docker_title'], STYLES['h1']),
           Paragraph(L['docker_intro'], STYLES['body']),
           Paragraph(L['docker_code'].replace('\n','<br/>').replace(' ','&nbsp;'),
                     STYLES['code']),
           HR()]

    # STEP 5 VERIFY
    fl += [Paragraph(L['verify_title'], STYLES['h1']),
           Paragraph(L['verify_body'], STYLES['body']),
           Paragraph(L['verify_code'].replace('\n','<br/>').replace(' ','&nbsp;'),
                     STYLES['code']),
           box(L['verify_note'], TEAL, LTEAL),
           HR()]

    # STEP 5B VALIDATOR KEY
    fl += [Paragraph(L['valkey_title'], STYLES['h1']),
           Paragraph(L['valkey_body'], STYLES['body']),
           Paragraph(L['valkey_code'], STYLES['code']),
           Paragraph(L['valkey_note'], STYLES['info']),
           HR()]

    # STEP 6 METAMASK
    fl += [Paragraph(L['mm_title'], STYLES['h1']),
           Paragraph(L['mm_body'], STYLES['body'])]
    mm_data = []
    for k, v in L['mm_rows']:
        mm_data.append([Paragraph(k, S('mk', fontSize=8.5, textColor=GRAY,
                                        fontName='Helvetica-Bold', leading=12)),
                         Paragraph(v, S('mv', fontName='Courier', fontSize=8.5,
                                         textColor=PURPLE, leading=12))])
    mm = Table(mm_data, colWidths=[5*cm, 11.6*cm])
    mm.setStyle(TableStyle([('ROWBACKGROUNDS',(0,0),(-1,-1),[white,ALTROW]),
                             ('GRID',(0,0),(-1,-1),0.3,LINE),
                             ('TOPPADDING',(0,0),(-1,-1),5),
                             ('BOTTOMPADDING',(0,0),(-1,-1),5),
                             ('LEFTPADDING',(0,0),(-1,-1),8)]))
    fl += [mm, HR()]

    # STEP 7 REWARDS
    fl += [Paragraph(L['rewards_title'], STYLES['h1']),
           box(L['rewards_box'], GOLD, LGOLD)]
    for i, step in enumerate(L['rewards_steps'], 1):
        row = Table([[Paragraph(str(i), S('rn', fontName='Helvetica-Bold', fontSize=9,
                                           textColor=white, alignment=TA_CENTER, leading=12)),
                      Paragraph(step, STYLES['body'])]],
                    colWidths=[0.65*cm, 15.95*cm])
        row.setStyle(TableStyle([('BACKGROUND',(0,0),(0,0),GOLD),
                                  ('VALIGN',(0,0),(-1,-1),'TOP'),
                                  ('TOPPADDING',(0,0),(-1,-1),3),
                                  ('BOTTOMPADDING',(0,0),(-1,-1),3),
                                  ('LEFTPADDING',(0,0),(0,0),0),
                                  ('LEFTPADDING',(0,0),(1,0),8)]))
        fl.append(row)
    fl.append(HR())

    # TROUBLESHOOTING
    fl += [Paragraph(L['trouble_title'], STYLES['h1']),
           trouble_table(L['trouble_rows'], L['trouble_cols']),
           Spacer(1,12), HR()]

    # FOOTER
    fl.append(Paragraph(L['footer'], STYLES['foot']))

    doc.build(fl)

# ── GENERATE ──────────────────────────────────────────────────────────────────

if __name__ == '__main__':
    import os
    out = 'C:/Users/aequitas-chain/downloads'
    os.makedirs(out, exist_ok=True)

    for lang_key, L in [('EN', EN), ('DE', DE)]:
        path = f'{out}/Aequitas_Node_Guide_{lang_key}.pdf'
        build_pdf(path, L)
        print(f'Generated: {path}')
    print('Done.')
