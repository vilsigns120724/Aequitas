"""Audit Round 5 fixes."""
import re, sys

def r(path, old, new, lbl):
    c = open(path, 'rb').read().decode('utf-8')
    if old in c:
        c = c.replace(old, new, 1)
        open(path, 'wb').write(c.encode('utf-8'))
        sys.stderr.write('OK: ' + lbl + '\n')
        return True
    sys.stderr.write('MISS: ' + lbl + '\n')
    return False

def ra(path, old, new, lbl):
    """Replace all occurrences."""
    c = open(path, 'rb').read().decode('utf-8')
    count = c.count(old)
    if count > 0:
        c = c.replace(old, new)
        open(path, 'wb').write(c.encode('utf-8'))
        sys.stderr.write('OK (' + str(count) + 'x): ' + lbl + '\n')
        return True
    sys.stderr.write('MISS: ' + lbl + '\n')
    return False

HTML = 'x/humanity/keeper/api_html.go'
STATE = 'x/humanity/keeper/state.go'
EVM_STORAGE = 'x/humanity/keeper/evm_storage.go'

# ── Fix 1: panic in evm_storage.go → non-fatal log ───────────────────────────
r(EVM_STORAGE,
  'fmt.Printf("[STARTUP] FATAL: could not enforce UNIQUE(human_wallet) on validator_keys: %v\\n", err)\n\t\tpanic("validator_keys uniqueness constraint failed — inspect DB for duplicates")',
  'Log.Error("Could not enforce UNIQUE(human_wallet) on validator_keys — DB may have duplicates", "error", err)',
  'evm_storage panic -> error log')

# ── Fix 2: IS_PRIMARY_NODE in PDF guide — all languages ──────────────────────
# Section 3 env vars tables: IS_PRIMARY_NODE says "Enables daily UBI distributions"
# Reality: DB-lock handles it now, IS_PRIMARY_NODE is obsolete
fixes_primary = [
    ('"true" only on the designated primary. Enables daily UBI + Validator + LP distributions.',
     '"true" only on the designated primary. DEPRECATED: distributions now use a DB-level lock (TryLockDistribution) — not required. Safe to leave unset.'),
    ('"true" nur auf dem designierten Primär-Node. Aktiviert tägliche Verteilungen.',
     '"true" nur auf dem designierten Primär-Node. VERALTET: Verteilungen nutzen jetzt DB-Lock — nicht mehr erforderlich.'),
    ('"true" solo en el nodo primario designado. Activa distribuciones diarias.',
     '"true" solo en el nodo primario. OBSOLETO: las distribuciones usan DB-lock — no requerido.'),
    ('"true" solo sul nodo primario designato. Abilita distribuzioni giornaliere.',
     '"true" solo sul nodo primario. OBSOLETO: le distribuzioni usano DB-lock — non necessario.'),
    ('"true" hanya pada node primer yang ditunjuk. Mengaktifkan distribusi harian.',
     '"true" hanya pada node primer. USANG: distribusi kini menggunakan DB-lock — tidak wajib.'),
    ('"true" uniquement sur le noeud principal designe. Active les distributions quotidiennes.',
     '"true" uniquement sur le noeud principal. OBSOLETE: les distributions utilisent DB-lock — non requis.'),
    ('"true" apenas no node principal designado. Ativa distribuicoes diarias.',
     '"true" apenas no node principal. OBSOLETO: distribuicoes usam DB-lock — nao necessario.'),
    ('"true" yalnizca belirlenmis birincil dugumde "true". Gunluk dagitimi etkinlestirir.',
     '"true" yalnizca birincil dugumde. KULLANIM DISI: dagitimlar DB-lock kullanir — gerekli degil.'),
]
c = open(HTML, 'rb').read().decode('utf-8')
count = 0
for old, new in fixes_primary:
    if old in c:
        c = c.replace(old, new)
        count += 1
open(HTML, 'wb').write(c.encode('utf-8'))
sys.stderr.write(f'OK ({count}): IS_PRIMARY_NODE obsolete description\n')

# ── Fix 3: PDF Guide Section 7 — PEER_NODES → PRIMARY_NODE_URL in all langs ──
c = open(HTML, 'rb').read().decode('utf-8')

# EN
c = c.replace(
    "s7:'7. P2P Networking & Sync',p7:'Set PEER_NODES to at least one known bootstrap URL. The node connects and syncs the full chain automatically using libp2p gossip (real-time) plus periodic HTTP pulls from peers (fallback). The primary node libp2p multiaddress is:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'The HTTP URL in PEER_NODES is more stable for bootstrap. The libp2p multiaddress above may change if the primary node is redeployed on Railway. When in doubt, use the HTTPS URL.'",
    "s7:'7. P2P Networking & Sync',p7:'Set PRIMARY_NODE_URL=https://aequitas.digital in your environment. The node auto-registers with the primary on startup, receives the full peer list, and begins syncing. The libp2p multiaddress below is for advanced/manual setups:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'PRIMARY_NODE_URL (HTTPS) is the recommended method. The libp2p multiaddress may change if the primary node is redeployed.'"
)
# DE
c = c.replace(
    "s7:'7. P2P-Netzwerk & Synchronisation',p7:'PEER_NODES auf mindestens eine bekannte Bootstrap-URL setzen.",
    "s7:'7. P2P-Netzwerk & Synchronisation',p7:'PRIMARY_NODE_URL=https://aequitas.digital setzen."
)
c = c.replace(
    "pn:'Die HTTP-URL in PEER_NODES ist stabiler für Bootstrap. Die Multiaddresse kann sich ändern wenn der Primär-Node auf Railway neu bereitgestellt wird. Im Zweifelsfall die HTTPS-URL verwenden.'",
    "pn:'PRIMARY_NODE_URL (HTTPS) ist die empfohlene Methode. Die libp2p-Multiaddresse kann sich bei einem Neudeployment ändern.'"
)
# ES
c = c.replace(
    "s7:'7. Red P2P y Sincronizacion',p7:'Establece PEER_NODES en al menos una URL de bootstrap.",
    "s7:'7. Red P2P y Sincronizacion',p7:'Establece PRIMARY_NODE_URL=https://aequitas.digital en tu entorno."
)
c = c.replace(
    "pn:'La URL HTTP en PEER_NODES es mas estable para bootstrap. La multidireccion puede cambiar si el nodo primario se reimplementa.'",
    "pn:'PRIMARY_NODE_URL (HTTPS) es el metodo recomendado. La multidireccion libp2p puede cambiar.'"
)
# IT
c = c.replace(
    "s7:'7. Rete P2P e Sincronizzazione',p7:'Imposta PEER_NODES su almeno un URL di bootstrap noto.",
    "s7:'7. Rete P2P e Sincronizzazione',p7:'Imposta PRIMARY_NODE_URL=https://aequitas.digital nel tuo ambiente."
)
c = c.replace(
    "pn:'L\\'URL HTTP in PEER_NODES e piu stabile per il bootstrap. Il multiindirizzo puo cambiare se il nodo primario viene ridistribuito su Railway.'",
    "pn:'PRIMARY_NODE_URL (HTTPS) e il metodo raccomandato. Il multiindirizzo libp2p puo cambiare.'"
)
# ID
c = c.replace(
    "s7:'7. Jaringan P2P dan Sinkronisasi',p7:'Atur PEER_NODES ke setidaknya satu URL bootstrap yang diketahui.",
    "s7:'7. Jaringan P2P dan Sinkronisasi',p7:'Atur PRIMARY_NODE_URL=https://aequitas.digital di environment."
)
c = c.replace(
    "pn:'URL HTTP di PEER_NODES lebih stabil untuk bootstrap. Multialamat libp2p dapat berubah jika node primer di-redeploy di Railway.'",
    "pn:'PRIMARY_NODE_URL (HTTPS) adalah metode yang direkomendasikan. Multialamat libp2p dapat berubah.'"
)
# FR
c = c.replace(
    "s7:'7. Reseau P2P et Synchronisation',p7:'Definir PEER_NODES sur au moins une URL de bootstrap connue.",
    "s7:'7. Reseau P2P et Synchronisation',p7:'Definir PRIMARY_NODE_URL=https://aequitas.digital dans l\\'environnement."
)
c = c.replace(
    "pn:'L\\'URL HTTP dans PEER_NODES est plus stable pour le bootstrap. La multiadresse peut changer si le noeud principal est redeploy sur Railway.'",
    "pn:'PRIMARY_NODE_URL (HTTPS) est la methode recommandee. La multiadresse libp2p peut changer.'"
)
# PT
c = c.replace(
    "s7:'7. Rede P2P e Sincronizacao',p7:'Definir PEER_NODES para pelo menos uma URL de bootstrap.",
    "s7:'7. Rede P2P e Sincronizacao',p7:'Definir PRIMARY_NODE_URL=https://aequitas.digital no ambiente."
)
c = c.replace(
    "pn:'A URL HTTP no PEER_NODES e mais estavel para bootstrap. O multiendereco pode mudar se o node principal for redeploy no Railway.'",
    "pn:'PRIMARY_NODE_URL (HTTPS) e o metodo recomendado. O multiendereco libp2p pode mudar.'"
)
# TR
c = c.replace(
    "s7:'7. P2P Ag ve Senkronizasyon',p7:'PEER_NODES\\'u en az bir bilinen bootstrap URL\\'sine ayarlayin.",
    "s7:'7. P2P Ag ve Senkronizasyon',p7:'PRIMARY_NODE_URL=https://aequitas.digital ortama ayarlayin."
)
c = c.replace(
    "pn:'PEER_NODES\\'daki HTTP URL\\'si bootstrap icin daha stabildir. Multiiadres, ana dugum Railway\\'de yeniden dagitilirsa degisebilir.'",
    "pn:'PRIMARY_NODE_URL (HTTPS) onerilen yontemdir. Libp2p multiadresi degisebilir.'"
)

open(HTML, 'wb').write(c.encode('utf-8'))
sys.stderr.write('OK: PDF guide section 7 PEER_NODES -> PRIMARY_NODE_URL (all languages)\n')

# ── Fix 4: bootstrap-desc in IT and FR translations ──────────────────────────
ra(HTML,
   "bootstrap-desc':'Per eseguire il tuo node, imposta la variabile d\\'ambiente PEER_NODES sull\\'indirizzo bootstrap qui sotto.",
   "bootstrap-desc':'Per eseguire il tuo node, imposta PRIMARY_NODE_URL=https://aequitas.digital nel tuo ambiente.",
   'IT bootstrap-desc PEER_NODES')
ra(HTML,
   "bootstrap-desc':'Définissez PEER_NODES sur l\\'adresse du nœud bootstrap.",
   "bootstrap-desc':'Définissez PRIMARY_NODE_URL=https://aequitas.digital dans votre environnement.",
   'FR bootstrap-desc PEER_NODES')

# ── Fix 5: Section 3 env var table — IS_PRIMARY_NODE obsolete in older langs ─
# Already partially done above; add more language coverage
for old, new, lbl in [
    ('"true" solo en el nodo primario designado. Activa distribuciones diarias.',
     '"true" solo en el nodo primario. OBSOLETO: distribuciones usan DB-lock ahora.',
     'ES IS_PRIMARY'),
]:
    r(HTML, old, new, lbl)

# ── Fix 6: NODE_KEY description wrong in multiple languages ───────────────────
# EN section 3 says "32-byte hex" but actually it's base64 libp2p key
for old, new, lbl in [
    ("'32-byte hex for stable libp2p peer identity. Auto-generated if omitted (changes on restart).'",
     "'Base64-encoded libp2p private key for stable peer identity. If not set: auto-generated and printed to stderr as \"SAVE THIS AS NODE_KEY ENVIRONMENT VAR: <base64>\". Copy and set it.'",
     'EN NODE_KEY correct format'),
    ("'32-Byte-Hex für stabile libp2p-Identität. Auto-generiert wenn nicht gesetzt.'",
     "'Base64-kodierter libp2p-Private-Key. Wenn nicht gesetzt: wird generiert und in stderr ausgegeben als \"SAVE THIS AS NODE_KEY: <base64>\". Kopieren und setzen.'",
     'DE NODE_KEY correct format'),
    ("'Hex 32 bytes para identidad P2P estable. Se genera automaticamente si no se establece.'",
     "'Clave libp2p base64 para identidad P2P estable. Si no se establece: auto-generada en stderr como \"SAVE THIS AS NODE_KEY: <base64>\".'",
     'ES NODE_KEY correct format'),
    ("'Hex 32 bytes para identidade P2P estaval. Auto-gerado se omitido.'",
     "'Chave libp2p base64 para identidade P2P. Se omitida: auto-gerada em stderr como \"SAVE THIS AS NODE_KEY: <base64>\".'",
     'PT NODE_KEY correct format'),
    ("'32-byte hex per identita P2P stabile. Auto-generato se omesso.'",
     "'Chiave libp2p base64 per identita P2P stabile. Se omessa: auto-generata in stderr come \"SAVE THIS AS NODE_KEY: <base64>\".'",
     'IT NODE_KEY correct format'),
    ("'Hex 32 bytes para identidade P2P estaval. Auto-gerado se omitido.'",
     "'Chave libp2p base64. Auto-gerada em stderr como \"SAVE THIS AS NODE_KEY: <base64>\" se omitida.'",
     'PT NODE_KEY correct format 2'),
    ("'Hex 32 byte per identita P2P stabile. Auto-generato se omesso.'",
     "'Chiave libp2p base64. Auto-generata in stderr come \"SAVE THIS AS NODE_KEY: <base64>\" se omessa.'",
     'IT NODE_KEY correct format 2'),
    ("'Hex 32 bytes para identidad P2P estable. Se genera automaticamente si no se establece.'",
     "'Clave libp2p base64. Auto-generada en stderr como \"SAVE THIS AS NODE_KEY: <base64>\".'",
     'ES NODE_KEY correct format 2'),
    ("'Hex 32 byte per identita P2P stabile. Auto-generato se omesso.'",
     "'Chiave libp2p base64. Auto-generata come \"SAVE THIS AS NODE_KEY: <base64>\" in stderr.'",
     'IT NODE_KEY 3'),
    ("'Hex 32 octets pour identite P2P stable. Auto-genere si omis.'",
     "'Cle libp2p base64 pour identite P2P stable. Auto-generee en stderr comme \"SAVE THIS AS NODE_KEY: <base64>\" si omise.'",
     'FR NODE_KEY correct format'),
    ("'Hex 32 bytes untuk identitas P2P stabil. Dibuat otomatis jika tidak diatur.'",
     "'Kunci libp2p base64 untuk identitas P2P stabil. Jika tidak diatur: auto-dibuat di stderr sebagai \"SAVE THIS AS NODE_KEY: <base64>\".'",
     'ID NODE_KEY correct format'),
    ("'Kararli libp2p kimligi icin 32 bayt hex. Atlanirsa otomatik olusturulur.'",
     "'Kararli P2P kimligi icin libp2p base64 anahtari. Atlanirsa stderr\\'de \"SAVE THIS AS NODE_KEY: <base64>\" olarak olusturulur.'",
     'TR NODE_KEY correct format'),
    ("'32-byte hex for stable libp2p peer identity. Auto-generated if omitted (changes on restart).'",
     "'Base64 libp2p key for stable P2P identity. If not set: auto-generated, printed to stderr as \"SAVE THIS AS NODE_KEY: <base64>\".'",
     'EN NODE_KEY correct 2'),
]:
    r(HTML, old, new, lbl)

sys.stderr.write('Done.\n')
