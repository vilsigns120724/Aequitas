"""Generate Aequitas Node Operator Guide PDFs in 8 Latin-script languages."""
from reportlab.lib.pagesizes import A4
from reportlab.lib.styles import ParagraphStyle
from reportlab.lib.units import cm
from reportlab.lib.colors import HexColor
from reportlab.platypus import (SimpleDocTemplate, Paragraph, Spacer, Table,
                                 TableStyle, HRFlowable, KeepTogether)
from reportlab.lib.enums import TA_CENTER

GOLD=HexColor('#F0B429'); PURPLE=HexColor('#9B72F6'); TEAL=HexColor('#22D3EE')
NEON=HexColor('#34D399'); RED=HexColor('#F87171'); MUTED=HexColor('#8892A4')
TEXT=HexColor('#E8EDF5'); CARD=HexColor('#131620'); BG=HexColor('#0C0E16')

def P(name,**kw):
    d=dict(fontName='Helvetica',textColor=TEXT,leading=14,spaceAfter=4)
    d.update(kw)
    return ParagraphStyle(name,**d)

S={
    'title':P('T',fontName='Helvetica-Bold',fontSize=20,textColor=GOLD,leading=26,spaceAfter=4,alignment=TA_CENTER),
    'sub':P('ST',fontSize=9,textColor=MUTED,alignment=TA_CENTER,spaceAfter=12),
    'h1':P('H1',fontName='Helvetica-Bold',fontSize=12,textColor=PURPLE,spaceBefore=14,spaceAfter=5),
    'h2':P('H2',fontName='Helvetica-Bold',fontSize=9,textColor=GOLD,spaceBefore=8,spaceAfter=3),
    'body':P('B',fontSize=8.5,textColor=TEXT,leading=13,spaceAfter=5),
    'code':P('C',fontName='Courier',fontSize=7.5,textColor=NEON,backColor=HexColor('#070B16'),leading=11,leftIndent=6,spaceAfter=6),
    'warn':P('W',fontSize=8,textColor=RED,leading=11,spaceAfter=5,leftIndent=6),
    'note':P('N',fontSize=8,textColor=TEAL,leading=11,spaceAfter=5,leftIndent=6),
    'muted':P('M',fontSize=8,textColor=MUTED,leading=11,spaceAfter=3),
    'bullet':P('BU',fontSize=8.5,textColor=TEXT,leading=12,leftIndent=12,spaceAfter=3),
}

def hr():
    return HRFlowable(width='100%',thickness=0.4,color=HexColor('#1E2D45'),spaceAfter=8,spaceBefore=3)

# Translations for all 8 Latin-script languages
LANGS = {
'en': {
    'title':'AEQUITAS NODE OPERATOR GUIDE',
    'version':'Version 1.0 · June 2026 · aequitas.digital',
    'tagline':'Complete step-by-step guide · No prior blockchain experience required · ~20-30 min',
    'what_h':'What is an Aequitas Node?',
    'what_b':'An Aequitas node runs in the cloud and participates in the network. It validates human registrations, produces blocks, and keeps the blockchain alive. Node operators earn a daily share of all protocol fees — automatically at <b>20:00 Berlin time (CEST/CET)</b>.',
    'pre_h':'Before You Start — What You Need',
    'pre':[(1,'<b>Aequitas account:</b> Register via the Android app. You need a wallet address to receive rewards.'),
           (2,'<b>GitHub account (free):</b> github.com — needed to fork the Aequitas code.'),
           (3,'<b>Railway account (free):</b> railway.app, sign in with GitHub. No server needed.'),
           (4,'<b>Dedicated node wallet (MetaMask):</b> A separate wallet just for your node. Export its private key: MetaMask → Account Details → Show Private Key. Keep this strictly private.'),
           (5,'<b>10-30 minutes</b> — Railway does most of the work automatically.')],
    'vars_h':'Environment Variables — Complete Reference',
    'vars_warn':'SECURITY: Your RELAYER_PRIVATE_KEY is like a master password. Never share it.',
    'vars_cols':['Variable','Required?','What to set'],
    'vars':[
        ('DATABASE_URL','YES','Auto-injected by Railway when PostgreSQL is in the same project.'),
        ('RELAYER_PRIVATE_KEY','YES','Private key of your node wallet (0x..., 66 chars). MetaMask → Account Details → Show Private Key.'),
        ('RELAYER_ADDRESS','Recommended','Public address matching RELAYER_PRIVATE_KEY (0x..., 42 chars).'),
        ('NODE_OPERATOR_WALLET','For rewards','Your Aequitas HUMAN wallet. Receives daily validator rewards at 20:00 Berlin.'),
        ('PEER_SECRET','Multi-node','Shared secret — ALL nodes must use the SAME value. Get from network operator.'),
        ('SELF_URL','Multi-node','Your node public URL: https://YOUR-NAME.up.railway.app'),
        ('PRIMARY_NODE_URL','Multi-node','Set to: https://aequitas.digital'),
        ('NODE_KEY','Optional','Base64 libp2p key for stable peer ID. If unset: auto-generated, printed to stderr as "SAVE THIS AS NODE_KEY ENVIRONMENT VAR: <base64>". Copy and set it.'),
        ('IS_PRIMARY_NODE','NO','Leave false. Only the official primary node uses true.'),
        ('RESET_STATE','NO','DANGEROUS: wipes database on restart. Development only.'),
    ],
    'steps_h':'Step-by-Step Deployment on Railway',
    'steps':[
        ('Step 1 — Fork Repository','Open github.com/hanoi96international-gif/Aequitas → click <b>Fork</b> → <b>Create fork</b>.'),
        ('Step 2 — Create PostgreSQL Database','railway.app → <b>New Project</b> → <b>+ New</b> → <b>Database</b> → <b>Add PostgreSQL</b>. Railway auto-injects DATABASE_URL.'),
        ('Step 3 — Deploy Node','In same project: <b>+ New</b> → <b>GitHub Repo</b> → select your fork → <b>Deploy Now</b>.'),
        ('Step 4 — Set Variables','Service → <b>Variables</b>. Set: RELAYER_PRIVATE_KEY, RELAYER_ADDRESS, NODE_OPERATOR_WALLET, PEER_SECRET, SELF_URL, PRIMARY_NODE_URL=https://aequitas.digital. Save → auto-redeploy.'),
        ('Step 5 — Get Public URL','Settings → Networking → <b>Generate Domain</b>. Open https://YOUR-URL/api/status — height should climb.'),
        ('Step 6 — Verify Logs','Look for: "API Server listening on port 8080", "[SYNC] Connected to peer", "[NODE] Registered node operator wallet".'),
        ('Step 7 — Earn Rewards','Rewards distributed automatically every day at <b>20:00 Berlin time</b>. 40% of all protocol fees split among node operators. Just keep your node running.'),
    ],
    'trouble_h':'Troubleshooting',
    'trouble_cols':['Symptom','Likely Cause','Solution'],
    'trouble':[
        ('Height stays at 0','PRIMARY_NODE_URL missing','Set PRIMARY_NODE_URL=https://aequitas.digital + SELF_URL. Redeploy.'),
        ('DATABASE_URL error','Wrong connection string','Format: postgres://user:pass@host:5432/dbname'),
        ('"no code at address"','V7 not deployed yet','Normal on first start — node auto-deploys. Wait a few seconds.'),
        ('No rewards','NODE_OPERATOR_WALLET missing','Add NODE_OPERATOR_WALLET=0xYOUR_HUMAN_WALLET.'),
        ('"Application error" (Railway)','Build/startup failure','Check Deploy Logs. Common: missing DATABASE_URL or wrong key format.'),
    ],
    'footer':'Questions: github.com/hanoi96international-gif/Aequitas · aequitas.digital · UBI: 20:00 Berlin time daily',
},
'de': {
    'title':'AEQUITAS NODE-BETREIBER-ANLEITUNG',
    'version':'Version 1.0 · Juni 2026 · aequitas.digital',
    'tagline':'Vollstaendige Schritt-fuer-Schritt-Anleitung · Keine Blockchain-Vorkenntnisse · ca. 20-30 Min.',
    'what_h':'Was ist ein Aequitas-Node?',
    'what_b':'Ein Aequitas-Node laeuft in der Cloud und nimmt am Netzwerk teil. Er validiert Menschenregistrierungen, produziert Bloecke und haelt die Blockchain am Leben. Node-Betreiber erhalten taeglich einen Anteil aller Protokollgebuehren — automatisch um <b>20:00 Uhr Berliner Zeit (CEST/CET)</b>.',
    'pre_h':'Vor dem Start — Was du brauchst',
    'pre':[(1,'<b>Aequitas-Konto:</b> Registriere dich ueber die Android-App. Du brauchst eine Wallet-Adresse fuer Belohnungen.'),
           (2,'<b>GitHub-Konto (kostenlos):</b> github.com — noetig um den Aequitas-Code zu forken.'),
           (3,'<b>Railway-Konto (kostenlos):</b> railway.app, mit GitHub anmelden. Kein eigener Server noetig.'),
           (4,'<b>Eigene Node-Wallet (MetaMask):</b> Eine separate Wallet nur fuer den Node. Privaten Schluessel exportieren: MetaMask → Kontodetails → Privaten Schluessel anzeigen. Streng geheimhalten.'),
           (5,'<b>10-30 Minuten</b> — Railway erledigt den Grossteil automatisch.')],
    'vars_h':'Umgebungsvariablen — Vollstaendige Referenz',
    'vars_warn':'SICHERHEIT: Dein RELAYER_PRIVATE_KEY ist wie ein Master-Passwort. Niemals teilen.',
    'vars_cols':['Variable','Erforderlich?','Was eintragen'],
    'vars':[
        ('DATABASE_URL','JA','Von Railway automatisch gesetzt wenn PostgreSQL im selben Projekt ist.'),
        ('RELAYER_PRIVATE_KEY','JA','Privater Schluessel der Node-Wallet (0x..., 66 Zeichen). MetaMask → Kontodetails → Privaten Schluessel anzeigen.'),
        ('RELAYER_ADDRESS','Empfohlen','Oeffentliche Adresse passend zu RELAYER_PRIVATE_KEY (0x..., 42 Zeichen).'),
        ('NODE_OPERATOR_WALLET','Fuer Bel.','Deine Aequitas-Mensch-Wallet. Erhaelt taeglich Validator-Belohnungen um 20:00 Uhr Berliner Zeit.'),
        ('PEER_SECRET','Multi-Node','Gemeinsames Geheimnis — ALLE Nodes muessen denselben Wert nutzen. Vom Netzwerkbetreiber erhalten.'),
        ('SELF_URL','Multi-Node','Oeffentliche URL des Nodes: https://DEIN-NAME.up.railway.app'),
        ('PRIMARY_NODE_URL','Multi-Node','Setzen auf: https://aequitas.digital'),
        ('NODE_KEY','Optional','Base64-libp2p-Schluessel. Wenn nicht gesetzt: wird automatisch generiert und in stderr ausgegeben als "SAVE THIS AS NODE_KEY ENVIRONMENT VAR: <base64>". Kopieren und setzen.'),
        ('IS_PRIMARY_NODE','NEIN','Nicht setzen (false). Nur der offizielle Primaer-Node nutzt true.'),
        ('RESET_STATE','NEIN','GEFAEHRLICH: loescht DB beim Neustart. Nur Entwicklung.'),
    ],
    'steps_h':'Schritt-fuer-Schritt-Deployment auf Railway',
    'steps':[
        ('Schritt 1 — Repository forken','github.com/hanoi96international-gif/Aequitas oeffnen → <b>Fork</b> klicken → <b>Create fork</b>.'),
        ('Schritt 2 — PostgreSQL-Datenbank erstellen','railway.app → <b>New Project</b> → <b>+ New</b> → <b>Database</b> → <b>Add PostgreSQL</b>. Railway setzt DATABASE_URL automatisch.'),
        ('Schritt 3 — Node deployen','Im selben Projekt: <b>+ New</b> → <b>GitHub Repo</b> → Fork auswaehlen → <b>Deploy Now</b>.'),
        ('Schritt 4 — Variablen setzen','Service → <b>Variables</b>. Eintragen: RELAYER_PRIVATE_KEY, RELAYER_ADDRESS, NODE_OPERATOR_WALLET, PEER_SECRET, SELF_URL, PRIMARY_NODE_URL=https://aequitas.digital. Speichern → auto-Redeploy.'),
        ('Schritt 5 — Oeffentliche URL erhalten','Settings → Networking → <b>Generate Domain</b>. https://DEINE-URL/api/status oeffnen — height sollte steigen.'),
        ('Schritt 6 — Logs pruefen','Auf: "API Server listening on port 8080", "[SYNC] Connected to peer", "[NODE] Registered node operator wallet" pruefen.'),
        ('Schritt 7 — Belohnungen erhalten','Belohnungen werden taeglich automatisch um <b>20:00 Uhr Berliner Zeit</b> verteilt. 40% aller Protokollgebuehren auf alle Node-Betreiber aufgeteilt. Node laufen lassen.'),
    ],
    'trouble_h':'Fehlerbehebung',
    'trouble_cols':['Symptom','Wahrscheinliche Ursache','Loesung'],
    'trouble':[
        ('Hoehe bleibt bei 0','PRIMARY_NODE_URL fehlt','PRIMARY_NODE_URL=https://aequitas.digital + SELF_URL setzen. Neu deployen.'),
        ('DATABASE_URL-Fehler','Falscher Connection-String','Format: postgres://user:pass@host:5432/dbname'),
        ('"no code at address"','V7 noch nicht deployed','Normal beim ersten Start — Node deployed automatisch. Kurz warten.'),
        ('Keine Belohnungen','NODE_OPERATOR_WALLET fehlt','NODE_OPERATOR_WALLET=0xDEINE_MENSCH_WALLET hinzufuegen.'),
        ('"Application error" (Railway)','Build-/Startfehler','Deploy-Logs pruefen. Haeufig: fehlende DATABASE_URL oder falsches Key-Format.'),
    ],
    'footer':'Fragen: github.com/hanoi96international-gif/Aequitas · aequitas.digital · UBI: taeglich 20:00 Uhr Berliner Zeit',
},
'es': {
    'title':'GUIA DEL OPERADOR DE NODO AEQUITAS',
    'version':'Version 1.0 · Junio 2026 · aequitas.digital',
    'tagline':'Guia completa paso a paso · Sin experiencia previa en blockchain · ~20-30 min',
    'what_h':'Que es un Nodo Aequitas?',
    'what_b':'Un nodo Aequitas se ejecuta en la nube y participa en la red. Valida registros humanos, produce bloques y mantiene viva la blockchain. Los operadores de nodo ganan una parte diaria de todas las tarifas del protocolo — automaticamente a las <b>20:00 hora de Berlin (CEST/CET)</b>.',
    'pre_h':'Antes de Empezar — Lo que Necesitas',
    'pre':[(1,'<b>Cuenta Aequitas:</b> Registrate mediante la app Android. Necesitas una direccion de wallet para recibir recompensas.'),
           (2,'<b>Cuenta GitHub (gratis):</b> github.com — necesaria para hacer fork del codigo Aequitas.'),
           (3,'<b>Cuenta Railway (gratis):</b> railway.app, inicia sesion con GitHub. No se necesita servidor propio.'),
           (4,'<b>Wallet dedicada para el nodo (MetaMask):</b> Una wallet separada solo para tu nodo. Exporta la clave privada: MetaMask → Detalles de cuenta → Mostrar clave privada. Mantenerla estrictamente privada.'),
           (5,'<b>10-30 minutos</b> — Railway hace la mayor parte del trabajo automaticamente.')],
    'vars_h':'Variables de Entorno — Referencia Completa',
    'vars_warn':'SEGURIDAD: Tu RELAYER_PRIVATE_KEY es como una contrasena maestra. Nunca la compartas.',
    'vars_cols':['Variable','Requerida?','Que configurar'],
    'vars':[
        ('DATABASE_URL','SI','Inyectada automaticamente por Railway cuando PostgreSQL esta en el mismo proyecto.'),
        ('RELAYER_PRIVATE_KEY','SI','Clave privada de tu wallet de nodo (0x..., 66 caracteres). MetaMask → Detalles de cuenta → Mostrar clave privada.'),
        ('RELAYER_ADDRESS','Recomendado','Direccion publica que coincide con RELAYER_PRIVATE_KEY (0x..., 42 caracteres).'),
        ('NODE_OPERATOR_WALLET','Para recompensas','Tu wallet humana de Aequitas. Recibe recompensas diarias de validador a las 20:00 hora de Berlin.'),
        ('PEER_SECRET','Multi-nodo','Secreto compartido — TODOS los nodos deben usar el MISMO valor. Obtenerlo del operador de red.'),
        ('SELF_URL','Multi-nodo','URL publica de tu nodo: https://TU-NOMBRE.up.railway.app'),
        ('PRIMARY_NODE_URL','Multi-nodo','Establecer en: https://aequitas.digital'),
        ('NODE_KEY','Opcional','Clave libp2p base64 para ID estable. Si no se establece: generada automaticamente, impresa en stderr como "SAVE THIS AS NODE_KEY ENVIRONMENT VAR: <base64>". Copiar y establecer.'),
        ('IS_PRIMARY_NODE','NO','Dejar en false. Solo el nodo primario oficial usa true.'),
        ('RESET_STATE','NO','PELIGROSO: borra la base de datos al reiniciar. Solo desarrollo.'),
    ],
    'steps_h':'Despliegue Paso a Paso en Railway',
    'steps':[
        ('Paso 1 — Hacer Fork del Repositorio','Abrir github.com/hanoi96international-gif/Aequitas → clic en <b>Fork</b> → <b>Create fork</b>.'),
        ('Paso 2 — Crear Base de Datos PostgreSQL','railway.app → <b>New Project</b> → <b>+ New</b> → <b>Database</b> → <b>Add PostgreSQL</b>. Railway inyecta DATABASE_URL automaticamente.'),
        ('Paso 3 — Desplegar Nodo','En el mismo proyecto: <b>+ New</b> → <b>GitHub Repo</b> → seleccionar tu fork → <b>Deploy Now</b>.'),
        ('Paso 4 — Configurar Variables','Servicio → <b>Variables</b>. Configurar: RELAYER_PRIVATE_KEY, RELAYER_ADDRESS, NODE_OPERATOR_WALLET, PEER_SECRET, SELF_URL, PRIMARY_NODE_URL=https://aequitas.digital. Guardar → redespliegue automatico.'),
        ('Paso 5 — Obtener URL Publica','Settings → Networking → <b>Generate Domain</b>. Abrir https://TU-URL/api/status — height debe subir.'),
        ('Paso 6 — Verificar Logs','Buscar: "API Server listening on port 8080", "[SYNC] Connected to peer", "[NODE] Registered node operator wallet".'),
        ('Paso 7 — Ganar Recompensas','Las recompensas se distribuyen automaticamente cada dia a las <b>20:00 hora de Berlin</b>. 40% de todas las tarifas del protocolo repartidas entre operadores. Solo mantener el nodo en funcionamiento.'),
    ],
    'trouble_h':'Solucion de Problemas',
    'trouble_cols':['Sintoma','Causa Probable','Solucion'],
    'trouble':[
        ('Height se queda en 0','PRIMARY_NODE_URL no configurada','Configurar PRIMARY_NODE_URL=https://aequitas.digital + SELF_URL. Redesplegar.'),
        ('Error DATABASE_URL','Cadena de conexion incorrecta','Formato: postgres://usuario:contrasena@host:5432/nombrebd'),
        ('"no code at address"','V7 aun no desplegado','Normal al primer inicio — el nodo lo despliega automaticamente. Esperar.'),
        ('Sin recompensas','NODE_OPERATOR_WALLET no configurada','Anadir NODE_OPERATOR_WALLET=0xTU_WALLET_HUMANA.'),
        ('"Application error" (Railway)','Fallo de build/inicio','Revisar Deploy Logs. Comun: DATABASE_URL faltante o formato de clave incorrecto.'),
    ],
    'footer':'Preguntas: github.com/hanoi96international-gif/Aequitas · aequitas.digital · UBI: 20:00 hora de Berlin diariamente',
},
'fr': {
    'title':'GUIDE DE L\'OPERATEUR DE NOEUD AEQUITAS',
    'version':'Version 1.0 · Juin 2026 · aequitas.digital',
    'tagline':'Guide complet etape par etape · Aucune experience blockchain requise · ~20-30 min',
    'what_h':'Qu\'est-ce qu\'un Noeud Aequitas?',
    'what_b':'Un noeud Aequitas s\'execute dans le cloud et participe au reseau. Il valide les enregistrements humains, produit des blocs et maintient la blockchain en vie. Les operateurs de noeuds gagnent une part quotidienne de tous les frais du protocole — automatiquement a <b>20h00 heure de Berlin (CEST/CET)</b>.',
    'pre_h':'Avant de Commencer — Ce dont Vous Avez Besoin',
    'pre':[(1,'<b>Compte Aequitas:</b> Inscrivez-vous via l\'app Android. Vous avez besoin d\'une adresse de wallet pour recevoir des recompenses.'),
           (2,'<b>Compte GitHub (gratuit):</b> github.com — necessaire pour forker le code Aequitas.'),
           (3,'<b>Compte Railway (gratuit):</b> railway.app, connexion avec GitHub. Aucun serveur requis.'),
           (4,'<b>Wallet dedié pour le noeud (MetaMask):</b> Un wallet separe uniquement pour votre noeud. Exportez la cle privee: MetaMask → Details du compte → Afficher la cle privee. A garder strictement prive.'),
           (5,'<b>10-30 minutes</b> — Railway fait l\'essentiel du travail automatiquement.')],
    'vars_h':'Variables d\'Environnement — Reference Complete',
    'vars_warn':'SECURITE: Votre RELAYER_PRIVATE_KEY est comme un mot de passe maitre. Ne jamais partager.',
    'vars_cols':['Variable','Requis?','Quoi configurer'],
    'vars':[
        ('DATABASE_URL','OUI','Injectee automatiquement par Railway quand PostgreSQL est dans le meme projet.'),
        ('RELAYER_PRIVATE_KEY','OUI','Cle privee de votre wallet de noeud (0x..., 66 caracteres). MetaMask → Details du compte → Afficher la cle privee.'),
        ('RELAYER_ADDRESS','Recommande','Adresse publique correspondant a RELAYER_PRIVATE_KEY (0x..., 42 caracteres).'),
        ('NODE_OPERATOR_WALLET','Pour recomp.','Votre wallet humain Aequitas. Recoit des recompenses de validateur quotidiennes a 20h00 Berlin.'),
        ('PEER_SECRET','Multi-noeud','Secret partage — TOUS les noeuds doivent utiliser la MEME valeur. Obtenir de l\'operateur reseau.'),
        ('SELF_URL','Multi-noeud','URL publique de votre noeud: https://VOTRE-NOM.up.railway.app'),
        ('PRIMARY_NODE_URL','Multi-noeud','Definir sur: https://aequitas.digital'),
        ('NODE_KEY','Optionnel','Cle libp2p base64 pour ID stable. Si non defini: genere automatiquement, imprime en stderr comme "SAVE THIS AS NODE_KEY ENVIRONMENT VAR: <base64>". Copier et definir.'),
        ('IS_PRIMARY_NODE','NON','Laisser a false. Seul le noeud primaire officiel utilise true.'),
        ('RESET_STATE','NON','DANGEREUX: efface la base de donnees au redemarrage. Developpement uniquement.'),
    ],
    'steps_h':'Deploiement Etape par Etape sur Railway',
    'steps':[
        ('Etape 1 — Forker le Depot','Ouvrir github.com/hanoi96international-gif/Aequitas → cliquer <b>Fork</b> → <b>Create fork</b>.'),
        ('Etape 2 — Creer une Base de Donnees PostgreSQL','railway.app → <b>New Project</b> → <b>+ New</b> → <b>Database</b> → <b>Add PostgreSQL</b>. Railway injecte DATABASE_URL automatiquement.'),
        ('Etape 3 — Deployer le Noeud','Dans le meme projet: <b>+ New</b> → <b>GitHub Repo</b> → selectionner votre fork → <b>Deploy Now</b>.'),
        ('Etape 4 — Configurer les Variables','Service → <b>Variables</b>. Configurer: RELAYER_PRIVATE_KEY, RELAYER_ADDRESS, NODE_OPERATOR_WALLET, PEER_SECRET, SELF_URL, PRIMARY_NODE_URL=https://aequitas.digital. Sauvegarder → redeploiement auto.'),
        ('Etape 5 — Obtenir l\'URL Publique','Settings → Networking → <b>Generate Domain</b>. Ouvrir https://VOTRE-URL/api/status — height doit monter.'),
        ('Etape 6 — Verifier les Logs','Chercher: "API Server listening on port 8080", "[SYNC] Connected to peer", "[NODE] Registered node operator wallet".'),
        ('Etape 7 — Gagner des Recompenses','Les recompenses sont distribuees automatiquement chaque jour a <b>20h00 heure de Berlin</b>. 40% de tous les frais du protocole repartis entre les operateurs. Garder simplement le noeud en marche.'),
    ],
    'trouble_h':'Depannage',
    'trouble_cols':['Symptome','Cause Probable','Solution'],
    'trouble':[
        ('Height reste a 0','PRIMARY_NODE_URL non configure','Configurer PRIMARY_NODE_URL=https://aequitas.digital + SELF_URL. Redeployer.'),
        ('Erreur DATABASE_URL','Chaine de connexion incorrecte','Format: postgres://utilisateur:motdepasse@host:5432/nombd'),
        ('"no code at address"','V7 pas encore deploye','Normal au premier demarrage — le noeud le deploie automatiquement. Attendre.'),
        ('Pas de recompenses','NODE_OPERATOR_WALLET manquant','Ajouter NODE_OPERATOR_WALLET=0xVOTRE_WALLET_HUMAIN.'),
        ('"Application error" (Railway)','Echec build/demarrage','Verifier Deploy Logs. Courant: DATABASE_URL manquant ou mauvais format de cle.'),
    ],
    'footer':'Questions: github.com/hanoi96international-gif/Aequitas · aequitas.digital · UBI: 20h00 Berlin quotidiennement',
},
'pt': {
    'title':'GUIA DO OPERADOR DE NO AEQUITAS',
    'version':'Versao 1.0 · Junho 2026 · aequitas.digital',
    'tagline':'Guia completo passo a passo · Sem experiencia previa em blockchain · ~20-30 min',
    'what_h':'O que e um No Aequitas?',
    'what_b':'Um no Aequitas executa na nuvem e participa na rede. Valida registros humanos, produz blocos e mantém a blockchain viva. Operadores de no ganham uma parte diaria de todas as taxas do protocolo — automaticamente as <b>20:00 horario de Berlim (CEST/CET)</b>.',
    'pre_h':'Antes de Comecar — O que Voce Precisa',
    'pre':[(1,'<b>Conta Aequitas:</b> Registre-se pelo app Android. Voce precisa de um endereco de carteira para receber recompensas.'),
           (2,'<b>Conta GitHub (gratis):</b> github.com — necessaria para fazer fork do codigo Aequitas.'),
           (3,'<b>Conta Railway (gratis):</b> railway.app, entre com GitHub. Nenhum servidor proprio necessario.'),
           (4,'<b>Carteira dedicada para o no (MetaMask):</b> Uma carteira separada so para seu no. Exporte a chave privada: MetaMask → Detalhes da conta → Mostrar chave privada. Mantenha estritamente privada.'),
           (5,'<b>10-30 minutos</b> — Railway faz a maior parte do trabalho automaticamente.')],
    'vars_h':'Variaveis de Ambiente — Referencia Completa',
    'vars_warn':'SEGURANCA: Sua RELAYER_PRIVATE_KEY e como uma senha mestra. Nunca compartilhe.',
    'vars_cols':['Variavel','Necessario?','O que configurar'],
    'vars':[
        ('DATABASE_URL','SIM','Injetada automaticamente pelo Railway quando PostgreSQL esta no mesmo projeto.'),
        ('RELAYER_PRIVATE_KEY','SIM','Chave privada da sua carteira de no (0x..., 66 caracteres). MetaMask → Detalhes da conta → Mostrar chave privada.'),
        ('RELAYER_ADDRESS','Recomendado','Endereco publico correspondente a RELAYER_PRIVATE_KEY (0x..., 42 caracteres).'),
        ('NODE_OPERATOR_WALLET','Para recomp.','Sua carteira humana Aequitas. Recebe recompensas diarias de validador as 20:00 Berlim.'),
        ('PEER_SECRET','Multi-no','Segredo compartilhado — TODOS os nos devem usar o MESMO valor. Obter do operador da rede.'),
        ('SELF_URL','Multi-no','URL publica do seu no: https://SEU-NOME.up.railway.app'),
        ('PRIMARY_NODE_URL','Multi-no','Definir como: https://aequitas.digital'),
        ('NODE_KEY','Opcional','Chave libp2p base64 para ID estavel. Se nao definida: gerada automaticamente, impressa em stderr como "SAVE THIS AS NODE_KEY ENVIRONMENT VAR: <base64>". Copiar e definir.'),
        ('IS_PRIMARY_NODE','NAO','Deixar como false. Apenas o no primario oficial usa true.'),
        ('RESET_STATE','NAO','PERIGOSO: apaga o banco de dados ao reiniciar. Somente desenvolvimento.'),
    ],
    'steps_h':'Implantacao Passo a Passo no Railway',
    'steps':[
        ('Passo 1 — Fazer Fork do Repositorio','Abrir github.com/hanoi96international-gif/Aequitas → clicar <b>Fork</b> → <b>Create fork</b>.'),
        ('Passo 2 — Criar Banco de Dados PostgreSQL','railway.app → <b>New Project</b> → <b>+ New</b> → <b>Database</b> → <b>Add PostgreSQL</b>. Railway injeta DATABASE_URL automaticamente.'),
        ('Passo 3 — Implantar No','No mesmo projeto: <b>+ New</b> → <b>GitHub Repo</b> → selecionar seu fork → <b>Deploy Now</b>.'),
        ('Passo 4 — Configurar Variaveis','Servico → <b>Variables</b>. Configurar: RELAYER_PRIVATE_KEY, RELAYER_ADDRESS, NODE_OPERATOR_WALLET, PEER_SECRET, SELF_URL, PRIMARY_NODE_URL=https://aequitas.digital. Salvar → reimplantacao automatica.'),
        ('Passo 5 — Obter URL Publica','Settings → Networking → <b>Generate Domain</b>. Abrir https://SUA-URL/api/status — height deve subir.'),
        ('Passo 6 — Verificar Logs','Procurar: "API Server listening on port 8080", "[SYNC] Connected to peer", "[NODE] Registered node operator wallet".'),
        ('Passo 7 — Ganhar Recompensas','Recompensas distribuidas automaticamente todos os dias as <b>20:00 horario de Berlim</b>. 40% de todas as taxas do protocolo divididas entre operadores. Apenas mantenha o no em funcionamento.'),
    ],
    'trouble_h':'Solucao de Problemas',
    'trouble_cols':['Sintoma','Causa Provavel','Solucao'],
    'trouble':[
        ('Height fica em 0','PRIMARY_NODE_URL nao configurada','Configurar PRIMARY_NODE_URL=https://aequitas.digital + SELF_URL. Reimplantar.'),
        ('Erro DATABASE_URL','String de conexao incorreta','Formato: postgres://usuario:senha@host:5432/nomebd'),
        ('"no code at address"','V7 ainda nao implantado','Normal no primeiro inicio — o no implanta automaticamente. Aguardar.'),
        ('Sem recompensas','NODE_OPERATOR_WALLET nao configurada','Adicionar NODE_OPERATOR_WALLET=0xSUA_CARTEIRA_HUMANA.'),
        ('"Application error" (Railway)','Falha de build/inicio','Verificar Deploy Logs. Comum: DATABASE_URL ausente ou formato de chave incorreto.'),
    ],
    'footer':'Perguntas: github.com/hanoi96international-gif/Aequitas · aequitas.digital · UBI: 20:00 Berlin diariamente',
},
'it': {
    'title':'GUIDA DELL\'OPERATORE DI NODO AEQUITAS',
    'version':'Versione 1.0 · Giugno 2026 · aequitas.digital',
    'tagline':'Guida completa passo dopo passo · Nessuna esperienza blockchain richiesta · ~20-30 min',
    'what_h':'Cos\'e un Nodo Aequitas?',
    'what_b':'Un nodo Aequitas viene eseguito nel cloud e partecipa alla rete. Valida le registrazioni umane, produce blocchi e mantiene viva la blockchain. Gli operatori di nodo guadagnano una quota giornaliera di tutte le commissioni del protocollo — automaticamente alle <b>20:00 ora di Berlino (CEST/CET)</b>.',
    'pre_h':'Prima di Iniziare — Cosa ti Serve',
    'pre':[(1,'<b>Account Aequitas:</b> Registrati tramite l\'app Android. Hai bisogno di un indirizzo wallet per ricevere ricompense.'),
           (2,'<b>Account GitHub (gratuito):</b> github.com — necessario per fare il fork del codice Aequitas.'),
           (3,'<b>Account Railway (gratuito):</b> railway.app, accedi con GitHub. Nessun server proprio richiesto.'),
           (4,'<b>Wallet dedicato per il nodo (MetaMask):</b> Un wallet separato solo per il tuo nodo. Esporta la chiave privata: MetaMask → Dettagli account → Mostra chiave privata. Tenerla strettamente privata.'),
           (5,'<b>10-30 minuti</b> — Railway fa la maggior parte del lavoro automaticamente.')],
    'vars_h':'Variabili d\'Ambiente — Riferimento Completo',
    'vars_warn':'SICUREZZA: La tua RELAYER_PRIVATE_KEY e come una password principale. Non condividerla mai.',
    'vars_cols':['Variabile','Richiesta?','Cosa configurare'],
    'vars':[
        ('DATABASE_URL','SI','Iniettata automaticamente da Railway quando PostgreSQL e nello stesso progetto.'),
        ('RELAYER_PRIVATE_KEY','SI','Chiave privata del tuo wallet di nodo (0x..., 66 caratteri). MetaMask → Dettagli account → Mostra chiave privata.'),
        ('RELAYER_ADDRESS','Consigliato','Indirizzo pubblico corrispondente a RELAYER_PRIVATE_KEY (0x..., 42 caratteri).'),
        ('NODE_OPERATOR_WALLET','Per ricomp.','Il tuo wallet umano Aequitas. Riceve ricompense quotidiane da validatore alle 20:00 Berlino.'),
        ('PEER_SECRET','Multi-nodo','Segreto condiviso — TUTTI i nodi devono usare lo STESSO valore. Ottenere dall\'operatore di rete.'),
        ('SELF_URL','Multi-nodo','URL pubblica del tuo nodo: https://TUO-NOME.up.railway.app'),
        ('PRIMARY_NODE_URL','Multi-nodo','Impostare su: https://aequitas.digital'),
        ('NODE_KEY','Opzionale','Chiave libp2p base64 per ID stabile. Se non impostata: generata automaticamente, stampata in stderr come "SAVE THIS AS NODE_KEY ENVIRONMENT VAR: <base64>". Copiare e impostare.'),
        ('IS_PRIMARY_NODE','NO','Lasciare a false. Solo il nodo primario ufficiale usa true.'),
        ('RESET_STATE','NO','PERICOLOSO: cancella il database al riavvio. Solo sviluppo.'),
    ],
    'steps_h':'Distribuzione Passo dopo Passo su Railway',
    'steps':[
        ('Passo 1 — Fork del Repository','Aprire github.com/hanoi96international-gif/Aequitas → cliccare <b>Fork</b> → <b>Create fork</b>.'),
        ('Passo 2 — Creare Database PostgreSQL','railway.app → <b>New Project</b> → <b>+ New</b> → <b>Database</b> → <b>Add PostgreSQL</b>. Railway inietta DATABASE_URL automaticamente.'),
        ('Passo 3 — Distribuire il Nodo','Nello stesso progetto: <b>+ New</b> → <b>GitHub Repo</b> → selezionare il fork → <b>Deploy Now</b>.'),
        ('Passo 4 — Impostare le Variabili','Servizio → <b>Variables</b>. Configurare: RELAYER_PRIVATE_KEY, RELAYER_ADDRESS, NODE_OPERATOR_WALLET, PEER_SECRET, SELF_URL, PRIMARY_NODE_URL=https://aequitas.digital. Salvare → ridistribuzione automatica.'),
        ('Passo 5 — Ottenere URL Pubblica','Settings → Networking → <b>Generate Domain</b>. Aprire https://TUO-URL/api/status — height deve salire.'),
        ('Passo 6 — Verificare i Log','Cercare: "API Server listening on port 8080", "[SYNC] Connected to peer", "[NODE] Registered node operator wallet".'),
        ('Passo 7 — Guadagnare Ricompense','Le ricompense vengono distribuite automaticamente ogni giorno alle <b>20:00 ora di Berlino</b>. Il 40% di tutte le commissioni del protocollo diviso tra gli operatori. Tenere semplicemente il nodo in esecuzione.'),
    ],
    'trouble_h':'Risoluzione dei Problemi',
    'trouble_cols':['Sintomo','Causa Probabile','Soluzione'],
    'trouble':[
        ('Height rimane a 0','PRIMARY_NODE_URL non configurata','Configurare PRIMARY_NODE_URL=https://aequitas.digital + SELF_URL. Ridistribuire.'),
        ('Errore DATABASE_URL','Stringa di connessione errata','Formato: postgres://utente:password@host:5432/nomedb'),
        ('"no code at address"','V7 non ancora distribuito','Normale al primo avvio — il nodo distribuisce automaticamente. Attendere.'),
        ('Nessuna ricompensa','NODE_OPERATOR_WALLET non configurato','Aggiungere NODE_OPERATOR_WALLET=0xTUO_WALLET_UMANO.'),
        ('"Application error" (Railway)','Errore build/avvio','Controllare Deploy Logs. Comune: DATABASE_URL mancante o formato chiave errato.'),
    ],
    'footer':'Domande: github.com/hanoi96international-gif/Aequitas · aequitas.digital · UBI: 20:00 ora di Berlino giornalmente',
},
'tr': {
    'title':'AEQUITAS DUGUM OPERATORU KILAVUZU',
    'version':'Surum 1.0 · Haziran 2026 · aequitas.digital',
    'tagline':'Eksiksiz adim adim kilavuz · Onceki blockchain deneyimi gerekmez · ~20-30 dak.',
    'what_h':'Aequitas Dugumu Nedir?',
    'what_b':'Bir Aequitas dugumu bulutta calisir ve aga katilir. Insan kayitlarini dogrular, bloklar uretir ve blockchain\'i canli tutar. Dugum operatorleri tum protokol ucretlerinin gunluk bir payini kazanir — otomatik olarak <b>Berlin saati 20:00\'de (CEST/CET)</b>.',
    'pre_h':'Baslamadan Once — Ihtiyaciniz Olanlar',
    'pre':[(1,'<b>Aequitas hesabi:</b> Android uygulamasi araciligiyla kaydolun. Odulleri almak icin bir cuzdan adresine ihtiyaciniz var.'),
           (2,'<b>GitHub hesabi (ucretsiz):</b> github.com — Aequitas kodunu forkladiniz icin gerekli.'),
           (3,'<b>Railway hesabi (ucretsiz):</b> railway.app, GitHub ile giris yapin. Kendi sunucunuza gerek yok.'),
           (4,'<b>Dugum icin ayri cuzdan (MetaMask):</b> Dugumunuz icin ayri bir cuzdan. Ozel anahtari disa aktarin: MetaMask → Hesap Detaylari → Ozel Anahtari Goster. Kesinlikle gizli tutun.'),
           (5,'<b>10-30 dakika</b> — Railway islerin cogunlugu otomatik olarak halleder.')],
    'vars_h':'Ortam Degiskenleri — Tam Referans',
    'vars_warn':'GUVENLIK: RELAYER_PRIVATE_KEY ana sifre gibidir. Asla paylasmayin.',
    'vars_cols':['Degisken','Gerekli?','Ne ayarlanmali'],
    'vars':[
        ('DATABASE_URL','EVET','Railway, PostgreSQL ayni projede oldugundan otomatik enjekte eder.'),
        ('RELAYER_PRIVATE_KEY','EVET','Dugum cuzdan ozel anahtari (0x..., 66 karakter). MetaMask → Hesap Detaylari → Ozel Anahtari Goster.'),
        ('RELAYER_ADDRESS','Onerilir','RELAYER_PRIVATE_KEY ile eslesen genel adres (0x..., 42 karakter).'),
        ('NODE_OPERATOR_WALLET','Odullar icin','Aequitas insan cuzdan adresiniz. Berlin saati 20:00\'de gunluk validator odulleri alir.'),
        ('PEER_SECRET','Cok dugum','Paylasilan sir — TUM dugumler AYNI degeri kullanmalidir. Ag operatorunden alin.'),
        ('SELF_URL','Cok dugum','Dugumunuzun genel URL\'si: https://ADINIZ.up.railway.app'),
        ('PRIMARY_NODE_URL','Cok dugum','Su sekilde ayarlayin: https://aequitas.digital'),
        ('NODE_KEY','Ihtiyari','Istikrarli P2P ID icin base64 libp2p anahtari. Ayarlanmamissa: otomatik uretilir, stderr\'de "SAVE THIS AS NODE_KEY ENVIRONMENT VAR: <base64>" olarak yazdirilir. Kopyalayin ve ayarlayin.'),
        ('IS_PRIMARY_NODE','HAYIR','False birakın. Yalnizca resmi birincil dugum true kullanir.'),
        ('RESET_STATE','HAYIR','TEHLIKELI: yeniden baslatmada veritabanini siler. Yalnizca gelistirme.'),
    ],
    'steps_h':'Railway\'de Adim Adim Dagitim',
    'steps':[
        ('Adim 1 — Depoyu Forklayin','github.com/hanoi96international-gif/Aequitas ayin → <b>Fork</b>\'a tiklayin → <b>Create fork</b>.'),
        ('Adim 2 — PostgreSQL Veritabani Olusturun','railway.app → <b>New Project</b> → <b>+ New</b> → <b>Database</b> → <b>Add PostgreSQL</b>. Railway DATABASE_URL\'yi otomatik enjekte eder.'),
        ('Adim 3 — Dugumu Dagatin','Ayni projede: <b>+ New</b> → <b>GitHub Repo</b> → forkunuzu secin → <b>Deploy Now</b>.'),
        ('Adim 4 — Degiskenleri Ayarlayin','Servis → <b>Variables</b>. Ayarlayin: RELAYER_PRIVATE_KEY, RELAYER_ADDRESS, NODE_OPERATOR_WALLET, PEER_SECRET, SELF_URL, PRIMARY_NODE_URL=https://aequitas.digital. Kaydet → otomatik yeniden dagitim.'),
        ('Adim 5 — Genel URL Edinin','Settings → Networking → <b>Generate Domain</b>. https://URL\'NIZ/api/status acin — height artmali.'),
        ('Adim 6 — Gunlukleri Kontrol Edin','Arayan: "API Server listening on port 8080", "[SYNC] Connected to peer", "[NODE] Registered node operator wallet".'),
        ('Adim 7 — Odullerinizi Kazanin','Odullar her gun Berlin saati <b>20:00\'de</b> otomatik olarak dagitilir. Tum protokol ucretlerinin %40\'i operatorler arasinda paylasilir. Sadece dugumunuzu calisir tutun.'),
    ],
    'trouble_h':'Sorun Giderme',
    'trouble_cols':['Belirti','Olasi Neden','Cozum'],
    'trouble':[
        ('Height 0\'da kaliyor','PRIMARY_NODE_URL ayarlanmamis','PRIMARY_NODE_URL=https://aequitas.digital + SELF_URL ayarlayin. Yeniden dagatin.'),
        ('DATABASE_URL hatasi','Yanlis baglanti dizisi','Format: postgres://kullanici:sifre@host:5432/veritabaniadi'),
        ('"no code at address"','V7 henuz dagitilmamis','Ilk baslatmada normal — dugum otomatik dagitir. Bekleyin.'),
        ('Odul yok','NODE_OPERATOR_WALLET ayarlanmamis','NODE_OPERATOR_WALLET=0xINSAN_CUZDAN_ADRESINIZ ekleyin.'),
        ('"Application error" (Railway)','Build/baslatma hatasi','Deploy Logs\'u kontrol edin. Yaygin: DATABASE_URL eksik veya yanlis anahtar formati.'),
    ],
    'footer':'Sorular: github.com/hanoi96international-gif/Aequitas · aequitas.digital · UBI: Berlin saati 20:00 gunluk',
},
'id': {
    'title':'PANDUAN OPERATOR NODE AEQUITAS',
    'version':'Versi 1.0 · Juni 2026 · aequitas.digital',
    'tagline':'Panduan langkah demi langkah · Tidak perlu pengalaman blockchain sebelumnya · ~20-30 menit',
    'what_h':'Apa itu Node Aequitas?',
    'what_b':'Node Aequitas berjalan di cloud dan berpartisipasi dalam jaringan. Node memvalidasi pendaftaran manusia, memproduksi blok, dan menjaga blockchain tetap hidup. Operator node mendapatkan bagian harian dari semua biaya protokol — secara otomatis pukul <b>20:00 waktu Berlin (CEST/CET)</b>.',
    'pre_h':'Sebelum Memulai — Apa yang Anda Butuhkan',
    'pre':[(1,'<b>Akun Aequitas:</b> Daftar melalui aplikasi Android. Anda memerlukan alamat dompet untuk menerima hadiah.'),
           (2,'<b>Akun GitHub (gratis):</b> github.com — diperlukan untuk mem-fork kode Aequitas.'),
           (3,'<b>Akun Railway (gratis):</b> railway.app, masuk dengan GitHub. Tidak perlu server sendiri.'),
           (4,'<b>Dompet khusus untuk node (MetaMask):</b> Dompet terpisah khusus untuk node Anda. Ekspor kunci privat: MetaMask → Detail Akun → Tampilkan Kunci Privat. Jaga ketat kerahasiaannya.'),
           (5,'<b>10-30 menit</b> — Railway melakukan sebagian besar pekerjaan secara otomatis.')],
    'vars_h':'Variabel Lingkungan — Referensi Lengkap',
    'vars_warn':'KEAMANAN: RELAYER_PRIVATE_KEY Anda seperti kata sandi utama. Jangan pernah dibagikan.',
    'vars_cols':['Variabel','Diperlukan?','Apa yang dikonfigurasi'],
    'vars':[
        ('DATABASE_URL','YA','Diinjeksi otomatis oleh Railway ketika PostgreSQL berada dalam proyek yang sama.'),
        ('RELAYER_PRIVATE_KEY','YA','Kunci privat dompet node Anda (0x..., 66 karakter). MetaMask → Detail Akun → Tampilkan Kunci Privat.'),
        ('RELAYER_ADDRESS','Disarankan','Alamat publik yang sesuai dengan RELAYER_PRIVATE_KEY (0x..., 42 karakter).'),
        ('NODE_OPERATOR_WALLET','Untuk hadiah','Dompet manusia Aequitas Anda. Menerima hadiah validator harian pukul 20:00 Berlin.'),
        ('PEER_SECRET','Multi-node','Rahasia bersama — SEMUA node harus menggunakan nilai YANG SAMA. Dapatkan dari operator jaringan.'),
        ('SELF_URL','Multi-node','URL publik node Anda: https://NAMA-ANDA.up.railway.app'),
        ('PRIMARY_NODE_URL','Multi-node','Atur ke: https://aequitas.digital'),
        ('NODE_KEY','Opsional','Kunci libp2p base64 untuk ID stabil. Jika tidak diatur: dibuat otomatis, dicetak di stderr sebagai "SAVE THIS AS NODE_KEY ENVIRONMENT VAR: <base64>". Salin dan atur.'),
        ('IS_PRIMARY_NODE','TIDAK','Biarkan false. Hanya node primer resmi yang menggunakan true.'),
        ('RESET_STATE','TIDAK','BERBAHAYA: menghapus database saat restart. Hanya untuk pengembangan.'),
    ],
    'steps_h':'Deployment Langkah demi Langkah di Railway',
    'steps':[
        ('Langkah 1 — Fork Repositori','Buka github.com/hanoi96international-gif/Aequitas → klik <b>Fork</b> → <b>Create fork</b>.'),
        ('Langkah 2 — Buat Database PostgreSQL','railway.app → <b>New Project</b> → <b>+ New</b> → <b>Database</b> → <b>Add PostgreSQL</b>. Railway menginjeksi DATABASE_URL secara otomatis.'),
        ('Langkah 3 — Deploy Node','Di proyek yang sama: <b>+ New</b> → <b>GitHub Repo</b> → pilih fork Anda → <b>Deploy Now</b>.'),
        ('Langkah 4 — Atur Variabel','Layanan → <b>Variables</b>. Atur: RELAYER_PRIVATE_KEY, RELAYER_ADDRESS, NODE_OPERATOR_WALLET, PEER_SECRET, SELF_URL, PRIMARY_NODE_URL=https://aequitas.digital. Simpan → redeploy otomatis.'),
        ('Langkah 5 — Dapatkan URL Publik','Settings → Networking → <b>Generate Domain</b>. Buka https://URL-ANDA/api/status — height harus naik.'),
        ('Langkah 6 — Verifikasi Log','Cari: "API Server listening on port 8080", "[SYNC] Connected to peer", "[NODE] Registered node operator wallet".'),
        ('Langkah 7 — Dapatkan Hadiah','Hadiah didistribusikan otomatis setiap hari pukul <b>20:00 waktu Berlin</b>. 40% dari semua biaya protokol dibagi di antara operator. Cukup jaga node Anda tetap berjalan.'),
    ],
    'trouble_h':'Pemecahan Masalah',
    'trouble_cols':['Gejala','Kemungkinan Penyebab','Solusi'],
    'trouble':[
        ('Height tetap di 0','PRIMARY_NODE_URL tidak diatur','Atur PRIMARY_NODE_URL=https://aequitas.digital + SELF_URL. Deploy ulang.'),
        ('Error DATABASE_URL','String koneksi salah','Format: postgres://pengguna:kata_sandi@host:5432/nama_db'),
        ('"no code at address"','V7 belum di-deploy','Normal saat pertama kali mulai — node men-deploy secara otomatis. Tunggu sebentar.'),
        ('Tidak ada hadiah','NODE_OPERATOR_WALLET tidak diatur','Tambahkan NODE_OPERATOR_WALLET=0xALAMAT_DOMPET_MANUSIA_ANDA.'),
        ('"Application error" (Railway)','Kegagalan build/startup','Periksa Deploy Logs. Umum: DATABASE_URL hilang atau format kunci salah.'),
    ],
    'footer':'Pertanyaan: github.com/hanoi96international-gif/Aequitas · aequitas.digital · UBI: 20:00 Berlin setiap hari',
},
}

def make_pdf(path, lang_key):
    L = LANGS[lang_key]
    doc = SimpleDocTemplate(path, pagesize=A4,
                            leftMargin=2*cm, rightMargin=2*cm,
                            topMargin=2*cm, bottomMargin=2*cm,
                            title=L['title'])
    story = []

    # Header
    story += [
        Paragraph(L['title'], S['title']),
        Paragraph(L['version'], S['sub']),
        Paragraph(L['tagline'], S['muted']),
        hr(),
    ]

    # What is a node
    story += [Paragraph(L['what_h'], S['h1']), Paragraph(L['what_b'], S['body'])]
    story.append(hr())

    # Prerequisites
    story.append(Paragraph(L['pre_h'], S['h1']))
    for num, text in L['pre']:
        story.append(Paragraph(f'<font color="#F0B429"><b>{num}.</b></font>  {text}', S['bullet']))
    story.append(hr())

    # Env vars
    story.append(Paragraph(L['vars_h'], S['h1']))
    story.append(Paragraph(L['vars_warn'], S['warn']))
    hdr = [Paragraph(f'<b>{h}</b>', ParagraphStyle('th', fontName='Helvetica-Bold',
            fontSize=8, textColor=TEXT, leading=10)) for h in L['vars_cols']]
    tdata = [hdr]
    for var, req, desc in L['vars']:
        req_c = '#F87171' if req in ('YES','SI','JA','SIM','OUI','EVET','YA') else \
                '#22D3EE' if 'Rec' in req or 'Emp' in req or 'Cons' in req or 'Oner' in req else \
                '#8892A4'
        tdata.append([
            Paragraph(f'<font name="Courier" color="#34D399">{var}</font>',
                      ParagraphStyle('v', fontName='Courier', fontSize=7.5, textColor=NEON, leading=10)),
            Paragraph(f'<font color="{req_c}"><b>{req}</b></font>',
                      ParagraphStyle('r', fontSize=7.5, textColor=TEXT, leading=10)),
            Paragraph(desc, ParagraphStyle('d', fontSize=7.5, textColor=MUTED, leading=11)),
        ])
    t = Table(tdata, colWidths=[3.8*cm, 2.2*cm, 10.5*cm], repeatRows=1)
    t.setStyle(TableStyle([
        ('BACKGROUND',(0,0),(-1,0),HexColor('#1A1D2B')),
        ('ROWBACKGROUNDS',(0,1),(-1,-1),[HexColor('#0C0E16'),HexColor('#131620')]),
        ('GRID',(0,0),(-1,-1),0.3,HexColor('#1E2D45')),
        ('VALIGN',(0,0),(-1,-1),'TOP'),
        ('TOPPADDING',(0,0),(-1,-1),4),('BOTTOMPADDING',(0,0),(-1,-1),4),
        ('LEFTPADDING',(0,0),(-1,-1),5),('RIGHTPADDING',(0,0),(-1,-1),5),
    ]))
    story.append(t)
    story.append(Spacer(1,8))
    story.append(hr())

    # Steps
    story.append(Paragraph(L['steps_h'], S['h1']))
    for title, text in L['steps']:
        story.append(KeepTogether([
            Paragraph(title, S['h2']),
            Paragraph(text, S['body']),
        ]))
    story.append(hr())

    # Troubleshooting
    story.append(Paragraph(L['trouble_h'], S['h1']))
    hdr2 = [Paragraph(f'<b>{h}</b>', ParagraphStyle('th2', fontName='Helvetica-Bold',
             fontSize=8, textColor=TEXT, leading=10)) for h in L['trouble_cols']]
    tdata2 = [hdr2]
    for s, c, sol in L['trouble']:
        tdata2.append([
            Paragraph(s,   ParagraphStyle('s2', fontSize=8, textColor=MUTED, leading=11)),
            Paragraph(c,   ParagraphStyle('c2', fontSize=8, textColor=MUTED, leading=11)),
            Paragraph(sol, ParagraphStyle('so2',fontSize=8, textColor=TEXT,  leading=11)),
        ])
    t2 = Table(tdata2, colWidths=[4.5*cm, 4*cm, 8*cm], repeatRows=1)
    t2.setStyle(TableStyle([
        ('BACKGROUND',(0,0),(-1,0),HexColor('#1A1D2B')),
        ('ROWBACKGROUNDS',(0,1),(-1,-1),[HexColor('#0C0E16'),HexColor('#131620')]),
        ('GRID',(0,0),(-1,-1),0.3,HexColor('#1E2D45')),
        ('VALIGN',(0,0),(-1,-1),'TOP'),
        ('TOPPADDING',(0,0),(-1,-1),4),('BOTTOMPADDING',(0,0),(-1,-1),4),
        ('LEFTPADDING',(0,0),(-1,-1),5),('RIGHTPADDING',(0,0),(-1,-1),5),
    ]))
    story.append(t2)
    story.append(hr())

    # Footer
    story.append(Paragraph(L['footer'], S['muted']))
    doc.build(story)
    print(f'Generated: {path}')

OUT = 'C:/Users/aequitas-chain/downloads'
for lang in LANGS:
    make_pdf(f'{OUT}/Aequitas_Node_Guide_{lang.upper()}.pdf', lang)
print('All PDFs generated.')
