"""
Aequitas Node Operator Guide — PDF Generator
Professional light-background design with dark readable text.
Supports 8 Latin-script languages.
"""
from reportlab.lib.pagesizes import A4
from reportlab.lib.styles import ParagraphStyle
from reportlab.lib.units import cm, mm
from reportlab.lib.colors import HexColor, white, black
from reportlab.platypus import (SimpleDocTemplate, Paragraph, Spacer, Table,
                                 TableStyle, HRFlowable, KeepTogether,
                                 PageBreak)
from reportlab.lib.enums import TA_LEFT, TA_CENTER, TA_RIGHT
from reportlab.platypus import Flowable

W, H = A4

# Color palette — works on white background
C_PURPLE  = HexColor('#6D28D9')   # deep purple
C_GOLD    = HexColor('#D97706')   # amber/gold
C_TEAL    = HexColor('#0D9488')   # teal
C_GREEN   = HexColor('#059669')   # green
C_RED     = HexColor('#DC2626')   # red
C_BODY    = HexColor('#1E1B4B')   # dark navy for body text
C_MUTED   = HexColor('#4B5563')   # gray
C_LIGHT   = HexColor('#F5F3FF')   # very light purple bg
C_BORDER  = HexColor('#DDD6FE')   # light purple border
C_BG_GOLD = HexColor('#FFFBEB')   # very light gold bg
C_BG_RED  = HexColor('#FEF2F2')   # very light red bg
C_BG_TEAL = HexColor('#F0FDFA')   # very light teal bg
C_LINE    = HexColor('#E5E7EB')   # horizontal rule

def PS(name, **kw):
    defaults = dict(fontName='Helvetica', textColor=C_BODY, leading=15, spaceAfter=4)
    defaults.update(kw)
    return ParagraphStyle(name, **defaults)

S = {
    'title':    PS('TI', fontName='Helvetica-Bold', fontSize=22, textColor=C_PURPLE,
                   leading=28, spaceAfter=2, alignment=TA_CENTER),
    'subtitle': PS('ST', fontSize=10, textColor=C_MUTED, alignment=TA_CENTER, spaceAfter=6),
    'tagline':  PS('TL', fontSize=8,  textColor=C_MUTED, alignment=TA_CENTER, spaceAfter=16),
    'h1':       PS('H1', fontName='Helvetica-Bold', fontSize=13, textColor=C_PURPLE,
                   spaceBefore=18, spaceAfter=6, borderPad=0,
                   backColor=C_LIGHT, borderColor=C_BORDER, borderWidth=0,
                   leftIndent=0, rightIndent=0),
    'h2':       PS('H2', fontName='Helvetica-Bold', fontSize=10, textColor=C_GOLD,
                   spaceBefore=10, spaceAfter=4),
    'body':     PS('BO', fontSize=9.5, textColor=C_BODY, leading=15, spaceAfter=6),
    'body_sm':  PS('BS', fontSize=8.5, textColor=C_MUTED, leading=13, spaceAfter=4),
    'code':     PS('CO', fontName='Courier', fontSize=8.5, textColor=C_PURPLE,
                   backColor=C_LIGHT, leading=13, leftIndent=8, rightIndent=8,
                   spaceAfter=8, spaceBefore=4),
    'warn':     PS('WA', fontSize=8.5, textColor=C_RED, leading=13, spaceAfter=6,
                   leftIndent=10),
    'note':     PS('NO', fontSize=8.5, textColor=C_TEAL, leading=13, spaceAfter=6,
                   leftIndent=10),
    'bullet':   PS('BU', fontSize=9.5, textColor=C_BODY, leading=14,
                   leftIndent=14, spaceAfter=3),
    'footer':   PS('FO', fontSize=7.5, textColor=C_MUTED, alignment=TA_CENTER,
                   leading=11),
}

def HR():
    return HRFlowable(width='100%', thickness=0.5, color=C_LINE,
                      spaceAfter=8, spaceBefore=4)

def section_title(text):
    """Section header with purple left border effect via table."""
    data = [[Paragraph(text, PS('sh', fontName='Helvetica-Bold', fontSize=12,
                                 textColor=C_PURPLE, leading=15))]]
    t = Table(data, colWidths=[16*cm])
    t.setStyle(TableStyle([
        ('BACKGROUND', (0,0), (-1,-1), C_LIGHT),
        ('LINEAFTER', (0,0), (0,-1), 0, white),
        ('LINEBEFORE', (0,0), (0,-1), 4, C_PURPLE),
        ('LEFTPADDING', (0,0), (-1,-1), 10),
        ('RIGHTPADDING', (0,0), (-1,-1), 8),
        ('TOPPADDING', (0,0), (-1,-1), 6),
        ('BOTTOMPADDING', (0,0), (-1,-1), 6),
    ]))
    return t

def info_box(text, style='teal'):
    colors = {
        'teal': (C_BG_TEAL, C_TEAL, C_TEAL),
        'gold':  (C_BG_GOLD, C_GOLD, C_GOLD),
        'red':   (C_BG_RED,  C_RED,  C_RED),
    }
    bg, border, tc = colors.get(style, colors['teal'])
    data = [[Paragraph(text, PS('ib', fontSize=8.5, textColor=tc, leading=13))]]
    t = Table(data, colWidths=[16*cm])
    t.setStyle(TableStyle([
        ('BACKGROUND', (0,0), (-1,-1), bg),
        ('BOX', (0,0), (-1,-1), 0.5, border),
        ('LEFTPADDING', (0,0), (-1,-1), 10),
        ('RIGHTPADDING', (0,0), (-1,-1), 10),
        ('TOPPADDING', (0,0), (-1,-1), 6),
        ('BOTTOMPADDING', (0,0), (-1,-1), 6),
        ('ROUNDEDCORNERS', [4]),
    ]))
    return t

# ── TRANSLATIONS ──────────────────────────────────────────────────────────────

LANGS = {

'en': {
'title':       'AEQUITAS NODE OPERATOR GUIDE',
'version':     'v1.0 · June 2026 · aequitas.digital',
'tagline':     'Complete step-by-step guide · No prior blockchain experience needed · ~20–30 min',
'chain':       'Aequitas Chain · Chain ID 1926 · EVM Compatible · Go 1.24',
'what_h':      'What is an Aequitas Node?',
'what_b':      'An Aequitas node runs in the cloud, validates human registrations, produces blocks, and keeps the blockchain alive. Node operators earn a daily share of protocol fees — automatically every day at <b>20:00 Berlin time (CEST/CET)</b>. No minimum uptime required.',
'warn_wallet': '⚠ IMPORTANT: To receive validator rewards, <b>NODE_OPERATOR_WALLET must be your registered Aequitas human wallet</b> — the one verified with AequitasBio. Only verified humans earn validator rewards.',
'pre_h':       'Before You Start',
'pre': [
    ('<b>Aequitas account:</b>', 'Register via the AequitasBio Android app. You need your human wallet address to receive rewards. If not yet registered, do that first.'),
    ('<b>GitHub account (free):</b>', 'Create at github.com — you need it to fork the Aequitas repository.'),
    ('<b>Railway account (free):</b>', 'Sign up at railway.app using GitHub. No server or command line needed.'),
    ('<b>Signing wallet (MetaMask):</b>', 'Create a <b>separate</b> MetaMask wallet for your node to sign transactions. Export its private key: MetaMask → Account Details → Show Private Key. Keep it strictly private.'),
    ('<b>10–30 minutes:</b>', 'Railway automates the build and deployment.'),
],
'vars_h':      'Environment Variables',
'vars_warn':   '⚠ RELAYER_PRIVATE_KEY is a private key. Never share it, never paste it in chat or email. Use a dedicated wallet — not your personal human wallet.',
'vars_cols':   ['Variable', 'Required?', 'Description'],
'vars': [
    ('DATABASE_URL',        'YES',          'PostgreSQL connection string. Railway auto-injects when PostgreSQL is in the same project.'),
    ('RELAYER_PRIVATE_KEY', 'YES',          'Private key of your signing wallet (0x…, 66 chars). MetaMask → Account Details → Show Private Key.'),
    ('RELAYER_ADDRESS',     'Recommended',  'Public address matching RELAYER_PRIVATE_KEY (0x…, 42 chars). Prevents startup errors.'),
    ('NODE_OPERATOR_WALLET','For rewards',  'Your registered HUMAN wallet on Aequitas. Receives daily validator rewards at 20:00 Berlin. MUST be a verified human.'),
    ('PEER_SECRET',         'Multi-node',   'Shared secret — all nodes in the network must use the SAME value. Obtain from the network operator. Do not publish.'),
    ('SELF_URL',            'Multi-node',   'Your node\'s public URL: https://YOUR-NAME.up.railway.app'),
    ('PRIMARY_NODE_URL',    'Multi-node',   'Set to: https://aequitas.digital — your node auto-registers here on startup.'),
    ('NODE_KEY',            'Optional',     'Base64 libp2p key for stable peer ID. If not set: auto-generated and printed to stderr as "SAVE THIS AS NODE_KEY: <base64>". Copy and set it.'),
    ('IS_PRIMARY_NODE',     'NO',           'Leave unset. Distribution now uses a DB-level lock — any node can run it without this variable.'),
    ('RESET_STATE',         'NO',           'DANGEROUS: wipes the entire database on restart. Development only. Never in production.'),
],
'steps_h':     'Step-by-Step Deployment (Railway)',
'steps': [
    ('Fork the Repository',
     'Open <b>github.com/hanoi96international-gif/Aequitas</b> in your browser. Click <b>Fork</b> in the top-right corner, then <b>Create fork</b>. GitHub creates a copy under your account.'),
    ('Create PostgreSQL Database',
     'Go to <b>railway.app</b>, sign in with GitHub. Click <b>New Project</b> → <b>+ New</b> → <b>Database</b> → <b>Add PostgreSQL</b>. Railway creates the database and will auto-inject DATABASE_URL.'),
    ('Deploy Your Node',
     'In the same Railway project: <b>+ New</b> → <b>GitHub Repo</b> → select your Aequitas fork. Railway detects the Dockerfile automatically. Click <b>Deploy Now</b>.'),
    ('Set Environment Variables',
     'Click your Aequitas service → <b>Variables</b>. Add the required variables (see table above). Minimum required:\n• RELAYER_PRIVATE_KEY\n• RELAYER_ADDRESS\n• NODE_OPERATOR_WALLET\n• PEER_SECRET\n• SELF_URL\n• PRIMARY_NODE_URL = https://aequitas.digital\nSave — Railway auto-redeploys.'),
    ('Get Your Public URL',
     'Settings → Networking → <b>Generate Domain</b>. Open <b>https://YOUR-URL/api/status</b> in your browser. You should see JSON with <b>"height"</b> climbing every ~6 seconds.'),
    ('Verify in Logs',
     'Check Deploy Logs for:\n✓  API Server listening on port 8080\n✓  [SYNC] Connected to peer https://aequitas.digital\n✓  [NODE] Registered node operator wallet: 0x…'),
    ('Earning Rewards',
     'Rewards are distributed automatically every day at <b>20:00 Berlin time (CEST/CET)</b>. The Validators Pool (40% of all protocol fees: swap fees, demurrage, wealth cap overflow) is split proportionally among all registered node operators by blocks produced. Keep your node running — no manual action needed.'),
],
'docker_h':    'Alternative: Docker Deployment',
'docker_cmd':  ('git clone https://github.com/hanoi96international-gif/Aequitas && cd Aequitas\n'
                'docker build -t aequitas-node .\n'
                'docker run -d --name aequitas --restart unless-stopped \\\n'
                '  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n'
                '  -e RELAYER_PRIVATE_KEY="0xYOUR_KEY" \\\n'
                '  -e RELAYER_ADDRESS="0xYOUR_ADDRESS" \\\n'
                '  -e NODE_OPERATOR_WALLET="0xYOUR_HUMAN_WALLET" \\\n'
                '  -e PEER_SECRET="from-network-operator" \\\n'
                '  -e SELF_URL="https://YOUR-URL" \\\n'
                '  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n'
                '  -p 8080:8080 aequitas-node'),
'trouble_h':   'Troubleshooting',
'trouble_cols':['Symptom', 'Cause', 'Solution'],
'trouble': [
    ('Height stays at 0',        'PRIMARY_NODE_URL missing',   'Set PRIMARY_NODE_URL=https://aequitas.digital and also set SELF_URL. Redeploy.'),
    ('DATABASE_URL error',       'Wrong connection string',    'Format: postgres://user:pass@host:5432/dbname'),
    ('"no code at address"',     'V7 not deployed yet',        'Normal on first start — node deploys V7 automatically. Wait 10–15 seconds.'),
    ('No rewards received',      'NODE_OPERATOR_WALLET wrong', 'Must be your registered HUMAN wallet on Aequitas (verified with AequitasBio).'),
    ('"Application error" Rail.',  'Build or startup failed',  'Check Deploy Logs. Common causes: missing DATABASE_URL or wrong key format.'),
    ('Port 8080 unreachable',    'Firewall config (Docker)',   'Open TCP port 8080 inbound in your firewall or cloud security group.'),
],
'footer':      'Aequitas Chain · Chain ID 1926 · aequitas.digital · UBI: daily 20:00 Berlin (CEST/CET)',
},

'de': {
'title':       'AEQUITAS NODE-BETREIBER-ANLEITUNG',
'version':     'v1.0 · Juni 2026 · aequitas.digital',
'tagline':     'Vollstaendige Schritt-fuer-Schritt-Anleitung · Keine Blockchain-Vorkenntnisse noetig · ca. 20–30 Min.',
'chain':       'Aequitas Chain · Chain ID 1926 · EVM-kompatibel · Go 1.24',
'what_h':      'Was ist ein Aequitas-Node?',
'what_b':      'Ein Aequitas-Node laeuft in der Cloud, validiert Menschenregistrierungen, produziert Bloecke und haelt die Blockchain am Leben. Node-Betreiber erhalten taeglich einen Anteil der Protokollgebuehren — automatisch jeden Tag um <b>20:00 Uhr Berliner Zeit (CEST/CET)</b>.',
'warn_wallet': '⚠ WICHTIG: Damit du Validator-Belohnungen erhaeltst, muss <b>NODE_OPERATOR_WALLET deine registrierte Aequitas-Mensch-Wallet sein</b> — die mit AequitasBio verifizierte. Nur verifizierte Menschen erhalten Validator-Belohnungen.',
'pre_h':       'Bevor du beginnst',
'pre': [
    ('<b>Aequitas-Konto:</b>', 'Registriere dich ueber die AequitasBio Android-App. Du brauchst deine Mensch-Wallet-Adresse fuer Belohnungen.'),
    ('<b>GitHub-Konto (kostenlos):</b>', 'Erstelle eines auf github.com — noetig um das Aequitas-Repository zu forken.'),
    ('<b>Railway-Konto (kostenlos):</b>', 'Melde dich auf railway.app mit GitHub an. Kein eigener Server oder Terminal noetig.'),
    ('<b>Signing-Wallet (MetaMask):</b>', 'Erstelle eine <b>separate</b> MetaMask-Wallet fuer deinen Node zum Signieren. Exportiere den privaten Schluessel: MetaMask → Kontodetails → Privaten Schluessel anzeigen. Streng geheimhalten.'),
    ('<b>10–30 Minuten:</b>', 'Railway automatisiert Build und Deployment.'),
],
'vars_h':      'Umgebungsvariablen',
'vars_warn':   '⚠ RELAYER_PRIVATE_KEY ist ein privater Schluessel. Niemals teilen, nicht in Chats einfuegen. Verwende eine dedizierte Wallet — nicht deine persoenliche Mensch-Wallet.',
'vars_cols':   ['Variable', 'Erforderlich?', 'Beschreibung'],
'vars': [
    ('DATABASE_URL',        'JA',           'PostgreSQL-Verbindungsstring. Railway setzt dies automatisch wenn PostgreSQL im gleichen Projekt ist.'),
    ('RELAYER_PRIVATE_KEY', 'JA',           'Privater Schluessel deiner Signing-Wallet (0x…, 66 Zeichen). MetaMask → Kontodetails → Privaten Schluessel anzeigen.'),
    ('RELAYER_ADDRESS',     'Empfohlen',    'Oeffentliche Adresse passend zu RELAYER_PRIVATE_KEY (0x…, 42 Zeichen). Verhindert Startfehler.'),
    ('NODE_OPERATOR_WALLET','Fuer Bel.',    'Deine registrierte MENSCH-Wallet auf Aequitas. Erhaelt taeglich Validator-Bel. um 20:00 Berliner Zeit. MUSS eine verifizierte Mensch-Wallet sein.'),
    ('PEER_SECRET',         'Multi-Node',   'Gemeinsames Geheimnis — alle Nodes muessen denselben Wert nutzen. Vom Netzwerkbetreiber erhalten.'),
    ('SELF_URL',            'Multi-Node',   'Oeffentliche URL deines Nodes: https://DEIN-NAME.up.railway.app'),
    ('PRIMARY_NODE_URL',    'Multi-Node',   'Auf https://aequitas.digital setzen — dein Node registriert sich hier automatisch beim Start.'),
    ('NODE_KEY',            'Optional',     'Base64-libp2p-Schluessel fuer stabile Peer-ID. Wenn nicht gesetzt: automatisch generiert, in stderr ausgegeben als "SAVE THIS AS NODE_KEY: <base64>". Kopieren und setzen.'),
    ('IS_PRIMARY_NODE',     'NEIN',         'Nicht setzen. Die Ausschuettung nutzt jetzt einen DB-Lock — jeder Node kann sie ausfuehren.'),
    ('RESET_STATE',         'NEIN',         'GEFAEHRLICH: Loescht die gesamte DB beim Neustart. Nur fuer Entwicklung. Niemals in Produktion.'),
],
'steps_h':     'Schritt-fuer-Schritt-Deployment (Railway)',
'steps': [
    ('Repository forken',
     'Oeffne <b>github.com/hanoi96international-gif/Aequitas</b>. Klicke <b>Fork</b> oben rechts, dann <b>Create fork</b>. GitHub erstellt eine Kopie unter deinem Konto.'),
    ('PostgreSQL-Datenbank erstellen',
     'Gehe zu <b>railway.app</b>, melde dich mit GitHub an. Klicke <b>New Project</b> → <b>+ New</b> → <b>Database</b> → <b>Add PostgreSQL</b>. Railway erstellt die Datenbank und setzt DATABASE_URL automatisch.'),
    ('Node deployen',
     'Im gleichen Railway-Projekt: <b>+ New</b> → <b>GitHub Repo</b> → deinen Aequitas-Fork auswaehlen. Railway erkennt das Dockerfile automatisch. Klicke <b>Deploy Now</b>.'),
    ('Umgebungsvariablen setzen',
     'Klicke deinen Aequitas-Service → <b>Variables</b>. Mindestens erforderlich:\n• RELAYER_PRIVATE_KEY\n• RELAYER_ADDRESS\n• NODE_OPERATOR_WALLET\n• PEER_SECRET\n• SELF_URL\n• PRIMARY_NODE_URL = https://aequitas.digital\nSpeichern → Railway deployt automatisch neu.'),
    ('Oeffentliche URL erhalten',
     'Settings → Networking → <b>Generate Domain</b>. Oeffne <b>https://DEINE-URL/api/status</b>. Du siehst JSON mit <b>"height"</b> der alle ~6 Sekunden steigt.'),
    ('Logs pruefen',
     'Pruefe die Deploy-Logs:\n✓  API Server listening on port 8080\n✓  [SYNC] Connected to peer https://aequitas.digital\n✓  [NODE] Registered node operator wallet: 0x…'),
    ('Belohnungen erhalten',
     'Belohnungen werden taeglich automatisch um <b>20:00 Uhr Berliner Zeit (CEST/CET)</b> verteilt. Der Validators-Pool (40% aller Protokollgebuehren) wird proportional nach produzierten Bloecken aufgeteilt. Node laufen lassen — kein manueller Eingriff noetig.'),
],
'docker_h':    'Alternative: Docker-Deployment',
'docker_cmd':  ('git clone https://github.com/hanoi96international-gif/Aequitas && cd Aequitas\n'
                'docker build -t aequitas-node .\n'
                'docker run -d --name aequitas --restart unless-stopped \\\n'
                '  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n'
                '  -e RELAYER_PRIVATE_KEY="0xDEIN_SCHLUESSEL" \\\n'
                '  -e RELAYER_ADDRESS="0xDEINE_ADRESSE" \\\n'
                '  -e NODE_OPERATOR_WALLET="0xDEINE_MENSCH_WALLET" \\\n'
                '  -e PEER_SECRET="vom-Netzwerkbetreiber" \\\n'
                '  -e SELF_URL="https://DEINE-URL" \\\n'
                '  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n'
                '  -p 8080:8080 aequitas-node'),
'trouble_h':   'Fehlerbehebung',
'trouble_cols':['Symptom', 'Ursache', 'Loesung'],
'trouble': [
    ('Height bleibt bei 0',       'PRIMARY_NODE_URL fehlt',     'PRIMARY_NODE_URL=https://aequitas.digital und SELF_URL setzen. Neu deployen.'),
    ('DATABASE_URL-Fehler',       'Falscher Verbindungsstring', 'Format: postgres://user:pass@host:5432/dbname'),
    ('"no code at address"',      'V7 noch nicht deployed',     'Normal beim ersten Start — Node deployed V7 automatisch. 10–15 Sekunden warten.'),
    ('Keine Belohnungen',         'NODE_OPERATOR_WALLET falsch','Muss deine registrierte MENSCH-Wallet auf Aequitas sein (mit AequitasBio verifiziert).'),
    ('"Application error" Rail.', 'Build- oder Startfehler',   'Deploy-Logs pruefen. Haeufige Ursachen: fehlende DATABASE_URL oder falsches Schluessel-Format.'),
    ('Port 8080 nicht erreichbar','Firewall-Config (Docker)',   'TCP-Port 8080 eingehend in Firewall oder Cloud-Security-Gruppe oeffnen.'),
],
'footer':      'Aequitas Chain · Chain ID 1926 · aequitas.digital · UBI: taeglich 20:00 Berliner Zeit (CEST/CET)',
},

} # end LANGS — add more languages here by following the same structure

# ── PDF GENERATION ────────────────────────────────────────────────────────────

def make_pdf(path, lang_key):
    L = LANGS[lang_key]

    doc = SimpleDocTemplate(
        path, pagesize=A4,
        leftMargin=2.2*cm, rightMargin=2.2*cm,
        topMargin=2*cm, bottomMargin=2*cm,
        title=L['title'], author='Aequitas Network',
        subject='Node Operator Guide'
    )

    story = []

    # ── COVER SECTION ──────────────────────────────────────────────────────────
    story.append(Spacer(1, 8))
    story.append(Paragraph(L['title'], S['title']))
    story.append(Paragraph(L['version'], S['subtitle']))
    story.append(Paragraph(L['tagline'], S['tagline']))

    # Chain info bar
    bar_data = [[Paragraph(L['chain'], PS('ci', fontSize=8, textColor=C_PURPLE,
                                           fontName='Courier', leading=11))]]
    bar = Table(bar_data, colWidths=[16*cm])
    bar.setStyle(TableStyle([
        ('BACKGROUND', (0,0), (-1,-1), C_LIGHT),
        ('BOX', (0,0), (-1,-1), 0.5, C_BORDER),
        ('LEFTPADDING', (0,0), (-1,-1), 12),
        ('RIGHTPADDING', (0,0), (-1,-1), 12),
        ('TOPPADDING', (0,0), (-1,-1), 6),
        ('BOTTOMPADDING', (0,0), (-1,-1), 6),
    ]))
    story.append(bar)
    story.append(Spacer(1, 12))
    story.append(HR())

    # ── WHAT IS A NODE ─────────────────────────────────────────────────────────
    story.append(KeepTogether([
        section_title(L['what_h']),
        Spacer(1, 6),
        Paragraph(L['what_b'], S['body']),
        Spacer(1, 4),
        info_box(L['warn_wallet'], 'gold'),
    ]))
    story.append(Spacer(1, 8))
    story.append(HR())

    # ── PREREQUISITES ──────────────────────────────────────────────────────────
    story.append(section_title(L['pre_h']))
    story.append(Spacer(1, 6))

    for i, (bold_part, desc) in enumerate(L['pre'], 1):
        row_data = [[
            Paragraph(f'<font color="#6D28D9"><b>{i}</b></font>',
                      PS('pn', fontSize=11, textColor=C_PURPLE, fontName='Helvetica-Bold',
                         alignment=TA_CENTER, leading=14)),
            Paragraph(f'{bold_part} {desc}', S['body']),
        ]]
        row = Table(row_data, colWidths=[0.8*cm, 15.2*cm])
        row.setStyle(TableStyle([
            ('VALIGN', (0,0), (-1,-1), 'TOP'),
            ('TOPPADDING', (0,0), (-1,-1), 3),
            ('BOTTOMPADDING', (0,0), (-1,-1), 3),
            ('LEFTPADDING', (0,0), (0,0), 0),
            ('LEFTPADDING', (0,0), (1,0), 8),
        ]))
        story.append(row)

    story.append(Spacer(1, 8))
    story.append(HR())

    # ── ENVIRONMENT VARIABLES TABLE ────────────────────────────────────────────
    story.append(section_title(L['vars_h']))
    story.append(Spacer(1, 6))
    story.append(info_box(L['vars_warn'], 'red'))
    story.append(Spacer(1, 8))

    # Table header
    def make_th(txt):
        return Paragraph(f'<b>{txt}</b>',
                         PS('th', fontName='Helvetica-Bold', fontSize=8.5,
                             textColor=white, leading=11))
    def make_td_var(txt):
        return Paragraph(txt, PS('tv', fontName='Courier', fontSize=8,
                                  textColor=C_PURPLE, leading=11))
    def make_td_req(txt):
        color = C_RED if txt in ('YES','JA','SI','OUI','SIM','EVET','YA') else \
                C_GREEN if 'reward' in txt.lower() or 'Bel' in txt or 'reward' in txt.lower() or 'recomp' in txt.lower() else \
                C_GOLD if 'Rec' in txt or 'Emp' in txt or 'Cons' in txt or 'Multi' in txt else C_MUTED
        return Paragraph(f'<b>{txt}</b>',
                         PS('tr', fontName='Helvetica-Bold', fontSize=8,
                             textColor=color, leading=11))
    def make_td_desc(txt):
        return Paragraph(txt, PS('td', fontSize=8, textColor=C_BODY, leading=12))

    hdr = [make_th(h) for h in L['vars_cols']]
    tdata = [hdr]
    for var, req, desc in L['vars']:
        tdata.append([make_td_var(var), make_td_req(req), make_td_desc(desc)])

    vtable = Table(tdata, colWidths=[4.2*cm, 2.4*cm, 9.4*cm], repeatRows=1)
    vtable.setStyle(TableStyle([
        ('BACKGROUND', (0,0), (-1,0), C_PURPLE),
        ('ROWBACKGROUNDS', (0,1), (-1,-1), [white, C_LIGHT]),
        ('GRID', (0,0), (-1,-1), 0.4, C_LINE),
        ('VALIGN', (0,0), (-1,-1), 'TOP'),
        ('TOPPADDING', (0,0), (-1,-1), 5),
        ('BOTTOMPADDING', (0,0), (-1,-1), 5),
        ('LEFTPADDING', (0,0), (-1,-1), 7),
        ('RIGHTPADDING', (0,0), (-1,-1), 7),
        ('LINEABOVE', (0,0), (-1,0), 0, C_PURPLE),
    ]))
    story.append(vtable)
    story.append(Spacer(1, 8))
    story.append(HR())

    # ── DEPLOYMENT STEPS ───────────────────────────────────────────────────────
    story.append(section_title(L['steps_h']))
    story.append(Spacer(1, 6))

    for i, (title, text) in enumerate(L['steps'], 1):
        color = C_GOLD if i == 7 else C_PURPLE
        step_content = [
            Paragraph(f'<b>{title}</b>',
                      PS('st', fontName='Helvetica-Bold', fontSize=10,
                          textColor=color, leading=13, spaceAfter=3)),
            Paragraph(text.replace('\n', '<br/>'), S['body']),
        ]
        step_data = [[
            Paragraph(f'{i}', PS('sn', fontName='Helvetica-Bold', fontSize=12,
                                   textColor=white, alignment=TA_CENTER, leading=16)),
            step_content,
        ]]
        step = Table(step_data, colWidths=[0.9*cm, 15.1*cm])
        step.setStyle(TableStyle([
            ('BACKGROUND', (0,0), (0,0), color),
            ('VALIGN', (0,0), (-1,-1), 'TOP'),
            ('TOPPADDING', (0,0), (0,0), 3),
            ('BOTTOMPADDING', (0,0), (0,0), 3),
            ('TOPPADDING', (0,0), (1,0), 0),
            ('BOTTOMPADDING', (0,0), (1,0), 6),
            ('LEFTPADDING', (0,0), (-1,-1), 0),
            ('LEFTPADDING', (0,0), (1,0), 10),
            ('ROUNDEDCORNERS', [4]),
        ]))
        story.append(step)
        story.append(Spacer(1, 4))

    story.append(HR())

    # ── DOCKER ALTERNATIVE ─────────────────────────────────────────────────────
    story.append(section_title(L['docker_h']))
    story.append(Spacer(1, 6))
    story.append(Paragraph(L['docker_cmd'].replace('\n','<br/>').replace(' ', '&nbsp;'), S['code']))
    story.append(HR())

    # ── TROUBLESHOOTING ────────────────────────────────────────────────────────
    story.append(section_title(L['trouble_h']))
    story.append(Spacer(1, 6))

    th2 = [Paragraph(f'<b>{h}</b>', PS('th2', fontName='Helvetica-Bold', fontSize=8.5,
                                        textColor=white, leading=11))
           for h in L['trouble_cols']]
    tdata2 = [th2]
    for symptom, cause, sol in L['trouble']:
        tdata2.append([
            Paragraph(symptom, PS('ts', fontSize=8, textColor=C_RED, leading=12, fontName='Helvetica-Bold')),
            Paragraph(cause,   PS('tc', fontSize=8, textColor=C_MUTED, leading=12)),
            Paragraph(sol,     PS('ts2',fontSize=8, textColor=C_BODY, leading=12)),
        ])
    ttable = Table(tdata2, colWidths=[4.5*cm, 4*cm, 7.5*cm], repeatRows=1)
    ttable.setStyle(TableStyle([
        ('BACKGROUND', (0,0), (-1,0), C_PURPLE),
        ('ROWBACKGROUNDS', (0,1), (-1,-1), [white, C_LIGHT]),
        ('GRID', (0,0), (-1,-1), 0.4, C_LINE),
        ('VALIGN', (0,0), (-1,-1), 'TOP'),
        ('TOPPADDING', (0,0), (-1,-1), 5),
        ('BOTTOMPADDING', (0,0), (-1,-1), 5),
        ('LEFTPADDING', (0,0), (-1,-1), 7),
        ('RIGHTPADDING', (0,0), (-1,-1), 7),
    ]))
    story.append(ttable)
    story.append(Spacer(1, 12))
    story.append(HR())

    # ── FOOTER ─────────────────────────────────────────────────────────────────
    story.append(Paragraph(L['footer'], S['footer']))

    doc.build(story)
    print(f'  Generated: {path}')


if __name__ == '__main__':
    import os
    out = 'C:/Users/aequitas-chain/downloads'
    os.makedirs(out, exist_ok=True)

    print('Generating Node Operator Guides...')
    for lang in LANGS:
        fname = f'{out}/Aequitas_Node_Guide_{lang.upper()}.pdf'
        make_pdf(fname, lang)
        print(f'    -> {lang.upper()} done')

    print(f'\nAll PDFs saved to {out}/')
