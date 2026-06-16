package keeper

const explorerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0">
<title>Aequitas — Proof of Humanity Chain</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
:root{--bg:#050A14;--card:#0D1421;--card2:#111E2E;--border:#1A2D45;--green:#00E676;--blue:#4FC3F7;--gold:#FFB300;--purple:#CE93D8;--red:#EF5350;--text:#E8F4FD;--muted:#6B8CAE;--teal:#4DD0E1}
body{background:var(--bg);color:var(--text);font-family:'Courier New',monospace;min-height:100vh;overflow-x:hidden}
header{background:#080F1E;border-bottom:1px solid var(--border);padding:0 20px;position:sticky;top:0;z-index:100;display:flex;align-items:center;justify-content:space-between;height:56px;gap:10px}
.logo-wrap{display:flex;align-items:center;gap:10px;flex-shrink:0}
.logo-icon{width:28px;height:28px;background:var(--gold);border-radius:6px;display:flex;align-items:center;justify-content:center;font-size:15px}
.logo-text{font-size:1rem;font-weight:900;color:var(--gold);letter-spacing:4px}
.logo-sub{font-size:0.5rem;color:var(--muted);letter-spacing:2px}
.header-right{display:flex;gap:8px;align-items:center}
.badge{display:flex;align-items:center;gap:4px;padding:4px 8px;border-radius:12px;font-size:0.6rem;letter-spacing:1px}
.badge-live{background:#00E67612;border:1px solid #00E67628;color:var(--green)}
.badge-dag{background:#4FC3F712;border:1px solid #4FC3F728;color:var(--blue)}
.pulse{width:5px;height:5px;border-radius:50%;background:var(--green);animation:pulse 2s infinite}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:0.3}}
.lang-sel{background:#080F1E;color:var(--muted);border:1px solid var(--border);border-radius:5px;padding:4px 8px;font-family:monospace;font-size:0.62rem;outline:none;cursor:pointer}
.tabs{background:#080F1E;border-bottom:1px solid var(--border);padding:0 20px;display:flex;overflow-x:auto;-webkit-overflow-scrolling:touch;scrollbar-width:none}
.tabs::-webkit-scrollbar{display:none}
.tab{padding:12px 14px;font-size:0.62rem;color:var(--muted);cursor:pointer;border-bottom:2px solid transparent;letter-spacing:1px;text-transform:uppercase;white-space:nowrap;transition:all 0.2s;flex-shrink:0}
.tab:hover{color:var(--text)}.tab.active{color:var(--blue);border-bottom-color:var(--blue)}
.tab-content{display:none}.tab-content.active{display:block}
.hero{padding:16px 16px 0}
.section-label{font-size:0.55rem;color:var(--muted);letter-spacing:4px;text-transform:uppercase;margin-bottom:12px}
.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(140px,1fr));gap:1px;background:var(--border);border:1px solid var(--border);border-radius:10px;overflow:hidden;margin-bottom:16px}
.stat{background:var(--card);padding:16px 14px;position:relative}
.stat-accent{position:absolute;top:0;left:0;right:0;height:2px}
.stat-icon{font-size:0.9rem;margin-bottom:6px}
.stat-lbl{font-size:0.55rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:5px}
.stat-val{font-size:1.6rem;font-weight:900;line-height:1;margin-bottom:3px}
.stat-sub{font-size:0.56rem;color:var(--muted);line-height:1.5}
.c-green .stat-val{color:var(--green)}.c-green .stat-accent{background:var(--green)}
.c-blue .stat-val{color:var(--blue)}.c-blue .stat-accent{background:var(--blue)}
.c-gold .stat-val{color:var(--gold)}.c-gold .stat-accent{background:var(--gold)}
.c-purple .stat-val{color:var(--purple)}.c-purple .stat-accent{background:var(--purple)}
.c-teal .stat-val{color:var(--teal)}.c-teal .stat-accent{background:var(--teal)}
.info-banner{background:#0D1E3A;border:1px solid #1A3A5C;border-radius:10px;padding:16px;margin-bottom:16px;display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:14px}
.ib-icon{font-size:1.2rem;margin-bottom:6px}
.ib-title{font-size:0.68rem;color:var(--gold);font-weight:bold;margin-bottom:6px;letter-spacing:1px}
.ib-text{font-size:0.63rem;color:var(--muted);line-height:1.8}
.main-grid{display:grid;grid-template-columns:1fr 300px;gap:12px;padding:0 16px 16px}
@media(max-width:800px){.main-grid{grid-template-columns:1fr}.right-col{display:none}}
.section{background:var(--card);border:1px solid var(--border);border-radius:10px;overflow:hidden}
.sec-head{padding:11px 16px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;background:#080F1E}
.sec-title{font-size:0.62rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;display:flex;align-items:center;gap:6px}
.sec-dot{width:5px;height:5px;border-radius:50%;background:var(--green)}
.sec-count{font-size:0.58rem;color:var(--muted);background:var(--card2);padding:2px 7px;border-radius:8px;border:1px solid var(--border)}
.sec-desc{padding:9px 16px;font-size:0.62rem;color:var(--muted);background:#080F1E;border-bottom:1px solid var(--border);line-height:1.7}
.block-item{padding:10px 16px;border-bottom:1px solid #0A1220;display:grid;grid-template-columns:56px 1fr auto;gap:8px;align-items:center}
.block-item:hover{background:#0D1421}.block-item:last-child{border-bottom:none}
.block-num{font-size:0.78rem;font-weight:bold;color:var(--blue)}
.block-hash{font-size:0.63rem;color:var(--muted);margin-bottom:2px;display:flex;align-items:center;gap:4px;flex-wrap:wrap}
.block-parents{font-size:0.57rem;color:#3A5570}
.block-right{text-align:right}
.block-humans{font-size:0.65rem;color:var(--gold);margin-bottom:2px}
.block-time{font-size:0.57rem;color:var(--green)}
.bm{background:#2D1B4E;color:var(--purple);font-size:0.53rem;padding:1px 4px;border-radius:3px;border:1px solid #4A2D7A}
.bt{background:#0D2A1A;color:var(--green);font-size:0.53rem;padding:1px 4px;border-radius:3px;border:1px solid #1A4A2A}
.empty{padding:32px;text-align:center;color:var(--muted);font-size:0.68rem;line-height:2.2}
.right-col{display:flex;flex-direction:column;gap:10px}
.ic{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:16px}
.ic-title{font-size:0.58rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:12px}
.ic-row{display:flex;justify-content:space-between;align-items:center;padding:7px 0;border-bottom:1px solid #0A1220}
.ic-row:last-child{border-bottom:none}
.ic-key{font-size:0.62rem;color:var(--muted)}
.ic-val{font-size:0.62rem;color:var(--text);text-align:right;max-width:58%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.ic-val.g{color:var(--green)}.ic-val.b{color:var(--blue)}.ic-val.go{color:var(--gold)}.ic-val.p{color:var(--purple)}
.mm-card{background:#0D1E3A;border:1px solid #1A3A5C;border-radius:10px;padding:14px}
.mm-title{font-size:0.58rem;color:var(--blue);letter-spacing:2px;margin-bottom:10px;font-weight:bold}
.mm-row{display:flex;justify-content:space-between;padding:5px 0;border-bottom:1px solid #1A2D45}
.mm-row:last-child{border-bottom:none}
.mm-key{font-size:0.58rem;color:var(--muted)}.mm-val{font-size:0.58rem;color:var(--purple)}
.mm-btn{width:100%;margin-top:10px;padding:9px;background:var(--blue);color:#050A14;border:none;border-radius:7px;cursor:pointer;font-family:monospace;font-size:0.65rem;font-weight:bold;letter-spacing:1px}
.phil-card{background:linear-gradient(135deg,#1A1200,#0D1421);border:1px solid #3A2800;border-radius:10px;padding:18px;text-align:center}
.phil-quote{font-size:0.78rem;color:var(--gold);font-style:italic;line-height:2;margin-bottom:5px}
.phil-sub{font-size:0.57rem;color:var(--muted);letter-spacing:2px}
.hs{padding:16px;display:grid;grid-template-columns:1fr 280px;gap:12px}
@media(max-width:800px){.hs{grid-template-columns:1fr}}
.hi{padding:11px 16px;border-bottom:1px solid #0A1220;display:flex;align-items:center;gap:10px}
.hi:hover{background:#0D1421}.hi:last-child{border-bottom:none}
.hav{width:34px;height:34px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:0.65rem;font-weight:bold;flex-shrink:0;border:2px solid}
.hbal{font-size:0.78rem;color:var(--gold);font-weight:bold;margin-bottom:1px}
.hadr{font-size:0.6rem;color:var(--muted);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.hbdg{font-size:0.55rem;padding:2px 7px;border-radius:8px;flex-shrink:0;background:#0D2A1A;color:var(--green);border:1px solid #1A4A2A}
.is{padding:16px;display:grid;grid-template-columns:1fr 1fr;gap:12px}
@media(max-width:700px){.is{grid-template-columns:1fr}}
.idx{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px}
.idx-title{font-size:0.58rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:8px}
.idx-desc{font-size:0.65rem;color:var(--muted);line-height:1.8;margin-bottom:14px}
.idx-big{font-size:2.6rem;font-weight:900;color:var(--gold);line-height:1}
.idx-lbl{font-size:0.58rem;color:var(--muted);margin-top:3px}
.bar-bg{height:7px;background:#0D1421;border-radius:4px;overflow:hidden;margin:12px 0 5px}
.bar-fill{height:100%;border-radius:4px;background:linear-gradient(90deg,var(--green),var(--gold),var(--red));transition:width 1.5s}
.bar-lbl{display:flex;justify-content:space-between;font-size:0.53rem;color:var(--muted)}
.mrow{display:grid;grid-template-columns:repeat(2,1fr);gap:7px;margin-top:12px}
.mbox{background:#080F1E;border-radius:6px;padding:10px;text-align:center}
.mval{font-size:1.1rem;font-weight:bold;color:var(--gold)}
.mlbl{font-size:0.55rem;color:var(--muted);margin-top:2px}
.story{font-size:0.67rem;line-height:2;color:var(--muted)}
.story p{margin-bottom:12px}
.hlbox{background:#080F1E;border-left:3px solid var(--gold);border-radius:0 8px 8px 0;padding:12px 16px;margin:14px 0;font-size:0.65rem;color:var(--text);line-height:1.8}
.ns{padding:16px;display:grid;grid-template-columns:1fr 1fr;gap:12px}
@media(max-width:700px){.ns{grid-template-columns:1fr}}
.nc{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:18px}
.nc-title{font-size:0.58rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:12px}
.nbox{background:#080F1E;border-radius:7px;padding:12px;border:1px solid var(--border);margin-bottom:8px}
.nstat{display:flex;align-items:center;gap:5px;font-size:0.65rem;color:var(--green);margin-bottom:4px;font-weight:bold}
.ndot{width:6px;height:6px;border-radius:50%;background:var(--green);box-shadow:0 0 5px var(--green)}
.nurl{font-size:0.57rem;color:var(--muted);word-break:break-all;margin-bottom:3px}
.ndesc{font-size:0.57rem;color:#3A5570}
.spect{width:100%;border-collapse:collapse}
.spect td{padding:7px 0;border-bottom:1px solid #0A1220;font-size:0.62rem}
.spect tr:last-child td{border-bottom:none}
.spect td:first-child{color:var(--muted);width:45%}
.spect td:last-child{text-align:right}
.bsbox{background:#080F1E;border-radius:7px;padding:10px;font-size:0.58rem;color:var(--purple);word-break:break-all;line-height:1.7;border:1px solid var(--border)}
.rs{padding:16px;max-width:600px;margin:0 auto}
.rhero{background:#0D1E3A;border:1px solid #1A3A5C;border-radius:10px;padding:20px;margin-bottom:14px;text-align:center}
.rhero-title{font-size:0.95rem;font-weight:bold;color:var(--text);margin-bottom:7px}
.rhero-sub{font-size:0.65rem;color:var(--muted);line-height:1.8}
.aonly{background:#0D1220;border:1px solid #1A2040;border-radius:10px;padding:18px;text-align:center;margin-bottom:14px}
.aonly-icon{font-size:1.8rem;margin-bottom:7px}
.aonly-title{font-size:0.68rem;color:var(--purple);font-weight:bold;letter-spacing:2px;margin-bottom:8px}
.aonly-text{font-size:0.63rem;color:var(--muted);line-height:1.8}
.rsteps{display:grid;grid-template-columns:repeat(4,1fr);gap:7px;margin-bottom:14px}
@media(max-width:520px){.rsteps{grid-template-columns:repeat(2,1fr)}}
.rstep{background:var(--card);border:1px solid var(--border);border-radius:8px;padding:14px;text-align:center}
.snum{width:26px;height:26px;background:var(--blue);border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 8px;font-weight:bold;font-size:0.7rem;color:#050A14}
.stitle{font-size:0.62rem;color:var(--text);font-weight:bold;margin-bottom:4px}
.sdesc{font-size:0.58rem;color:var(--muted);line-height:1.6}
.pbar{background:#0D1A0D;border:1px solid #1A3020;border-radius:7px;padding:9px 12px;margin-bottom:12px;font-size:0.62rem;color:var(--green);text-align:center;line-height:1.7}
.rcard{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:18px;margin-bottom:12px}
.wbox{background:#0D1A0D;border:1px solid #1A3020;border-radius:7px;padding:9px;margin-bottom:9px;display:none}
.wlbl{font-size:0.55rem;color:var(--muted);margin-bottom:2px;letter-spacing:1px}
.wadr{font-size:0.7rem;color:var(--green);font-weight:bold}
.pbox{background:var(--card2);border:1px solid #3A2800;border-radius:7px;padding:9px;margin-bottom:9px;display:none}
.plbl{font-size:0.55rem;color:var(--gold);margin-bottom:2px;letter-spacing:1px}
.pval{font-size:0.62rem;color:var(--muted)}
.rbtn{width:100%;padding:13px;border-radius:7px;border:none;cursor:pointer;font-family:monospace;font-size:0.72rem;font-weight:bold;letter-spacing:1px;transition:all 0.2s;margin-bottom:7px}
.bc{background:var(--blue);color:#050A14}.bc:hover{opacity:0.87}
.br{background:var(--gold);color:#050A14}.br:hover{opacity:0.87}
.rbtn:disabled{opacity:0.3;cursor:not-allowed}
.rlog{background:#080F1E;border-radius:7px;padding:10px;font-size:0.63rem;line-height:1.9;min-height:50px;border:1px solid var(--border)}
.rlog .ok{color:var(--green)}.rlog .err{color:var(--red)}.rlog .info{color:var(--gold)}
.ps{padding:16px;max-width:800px;margin:0 auto}
@media(max-width:480px){
  .stats-grid{grid-template-columns:repeat(2,1fr)}
  .stat-val{font-size:1.3rem}
  header{height:50px}
  .logo-text{font-size:0.85rem;letter-spacing:2px}
  .badge-dag{display:none}
}
</style>
</head>
<body>
<header>
  <div class="logo-wrap">
    <div class="logo-icon">⚖</div>
    <div><div class="logo-text">AEQUITAS</div><div class="logo-sub" data-i18n="logo-sub">PROOF OF HUMANITY</div></div>
  </div>
  <select class="lang-sel" id="lang-sel" onchange="setLang(this.value)">
    <option value="en">🌐 EN</option>
    <option value="de">🌐 DE</option>
    <option value="es">🌐 ES</option>
    <option value="ru">🌐 RU</option>
    <option value="zh">🌐 ZH</option>
    <option value="id">🌐 ID</option>
  </select>
  <div class="header-right">
    <div class="badge badge-live"><span class="pulse"></span><span data-i18n="live">LIVE</span></div>
    <div class="badge badge-dag">● BLOCKDAG</div>
  </div>
</header>
<div class="tabs">
  <div class="tab active" onclick="showTab('register',this)" data-i18n="tab-register">🔐 Register</div>
  <div class="tab" onclick="showTab('explorer',this)" data-i18n="tab-explorer">🔍 Explorer</div>
  <div class="tab" onclick="showTab('humans',this)" data-i18n="tab-humans">👥 Humans</div>
  <div class="tab" onclick="showTab('index',this)" data-i18n="tab-index">📊 Index</div>
  <div class="tab" onclick="showTab('network',this)" data-i18n="tab-network">🌐 Network</div>
  <div class="tab" onclick="showTab('protocol',this)" data-i18n="tab-protocol">📜 Protocol V7</div>
</div>

<!-- REGISTER -->
<div id="tab-register" class="tab-content active">
<div class="rs">
  <div class="rhero">
    <div class="rhero-title" data-i18n="reg-title">🔐 Register as a Verified Human</div>
    <div class="rhero-sub" data-i18n="reg-sub">Join the Aequitas network and receive 1,000 AEQ. One-time, permanent, gasless. No personal data stored.</div>
  </div>
  <div class="aonly">
    <div class="aonly-icon">📱</div>
    <div class="aonly-title" data-i18n="app-title">REGISTRATION VIA ANDROID APP</div>
    <div class="aonly-text" data-i18n="app-text">Proof of Humanity requires biometric verification on your personal device. Your fingerprint is processed by the Hardware Secure Element — raw data never leaves your phone. Download the app, scan your fingerprint, connect your wallet, and your <strong style="color:var(--gold)">1,000 AEQ will be granted automatically</strong>.</div>
  </div>
  <div class="rsteps">
    <div class="rstep"><div class="snum">1</div><div class="stitle" data-i18n="s1t">Biometric Scan</div><div class="sdesc" data-i18n="s1d">Open app · scan fingerprint · HSE processes · data never leaves device</div></div>
    <div class="rstep"><div class="snum">2</div><div class="stitle" data-i18n="s2t">ZKP Generation</div><div class="sdesc" data-i18n="s2d">Groth16 proof generated · uniqueness verified · hash never revealed</div></div>
    <div class="rstep"><div class="snum">3</div><div class="stitle" data-i18n="s3t">Connect Wallet</div><div class="sdesc" data-i18n="s3d">App opens MetaMask · connect wallet · address receives 1,000 AEQ</div></div>
    <div class="rstep"><div class="snum">4</div><div class="stitle" data-i18n="s4t">1,000 AEQ</div><div class="sdesc" data-i18n="s4d">Registered on V6 · confirmed in next block · app notifies automatically</div></div>
  </div>
  <div class="pbar" data-i18n="priv-bar">🔒 Hardware Secure Element · Groth16 ZKP · Data never leaves device · No gas fees · Permanent Sybil protection</div>
  <div class="rcard">
    <div class="wbox" id="wbox"><div class="wlbl" data-i18n="conn-wallet">CONNECTED WALLET</div><div class="wadr" id="wadr">—</div></div>
    <div class="pbox" id="pbox"><div class="plbl" data-i18n="proof-recv">⚡ ZK PROOF RECEIVED</div><div class="pval" id="pval" data-i18n="proof-hint">Connect wallet to register</div></div>
    <button class="rbtn bc" id="btn-conn" onclick="connectWallet()" data-i18n="btn-conn">🦊 CONNECT METAMASK</button>
    <button class="rbtn br" id="btn-reg" onclick="doRegister()" disabled data-i18n="btn-reg">🔐 REGISTER ON-CHAIN</button>
    <div class="rlog" id="rlog"><span class="info" data-i18n="reg-log-hint">// Open Aequitas Android App to generate your proof, then return here...</span></div>
  </div>
  <div class="ic">
    <div class="ic-title" data-i18n="reg-details">Registration Details</div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-network">Network</span><span class="ic-val p">Aequitas Chain (BlockDAG)</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-chainid">Chain ID</span><span class="ic-val b">1926</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-grant">Grant Amount</span><span class="ic-val go">1,000 AEQ</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-fee">Gas Fee</span><span class="ic-val g" data-i18n="free">FREE (gasless)</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-limit">Registrations</span><span class="ic-val" data-i18n="k-limit-v">Once per human · permanent · immutable</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-bio">Biometric Data</span><span class="ic-val g" data-i18n="never-stored">Never stored anywhere</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-conf">Confirmation</span><span class="ic-val" data-i18n="k-conf-v">Within 6 seconds</span></div>
  </div>
</div>
</div>

<!-- EXPLORER -->
<div id="tab-explorer" class="tab-content">
<div class="hero">
  <div class="section-label" data-i18n="live-stats">Live Chain Statistics</div>
  <div class="stats-grid">
    <div class="stat c-blue"><div class="stat-accent"></div><div class="stat-icon">🔗</div><div class="stat-lbl" data-i18n="s-height">Block Height</div><div class="stat-val" id="s-height">—</div><div class="stat-sub" data-i18n="s-height-sub">New block every 6s · BlockDAG · Two nodes parallel</div></div>
    <div class="stat c-green"><div class="stat-accent"></div><div class="stat-icon">🧬</div><div class="stat-lbl" data-i18n="s-humans">Verified Humans</div><div class="stat-val" id="s-humans">—</div><div class="stat-sub" data-i18n="s-humans-sub">Biometric ZKP · One person, one wallet, forever</div></div>
    <div class="stat c-gold"><div class="stat-accent"></div><div class="stat-icon">🪙</div><div class="stat-lbl" data-i18n="s-supply">Total Supply</div><div class="stat-val" id="s-supply">—</div><div class="stat-sub" data-i18n="s-supply-sub">Always = Humans x 1,000 AEQ</div></div>
    <div class="stat c-purple"><div class="stat-accent"></div><div class="stat-icon">⚖</div><div class="stat-lbl" data-i18n="s-index">Aequitas Index</div><div class="stat-val" id="s-index">—</div><div class="stat-sub" data-i18n="s-index-sub">0 = perfect equality · 100 = max inequality</div></div>
    <div class="stat c-teal"><div class="stat-accent"></div><div class="stat-icon">⚡</div><div class="stat-lbl" data-i18n="s-uptime">Uptime</div><div class="stat-val" id="s-uptime" style="font-size:1rem">—</div><div class="stat-sub" data-i18n="s-uptime-sub">Node v0.3.0 · Railway + Render · PostgreSQL</div></div>
  </div>
  <div class="info-banner">
    <div><div class="ib-icon">🧬</div><div class="ib-title" data-i18n="ib-poh">Proof of Humanity</div><div class="ib-text" data-i18n="ib-poh-t">Every AEQ holder must prove they are a unique living human. No bots, no corporations, no AI can hold AEQ. Only real humans. Biometric data never leaves your device.</div></div>
    <div><div class="ib-icon">⚖</div><div class="ib-title" data-i18n="ib-fair">Radically Fair</div><div class="ib-text" data-i18n="ib-fair-t">Every verified human receives exactly 1,000 AEQ. No pre-mine, no founder allocation. Total supply always equals verified humans x 1,000.</div></div>
    <div><div class="ib-icon">🔗</div><div class="ib-title" data-i18n="ib-dag">BlockDAG Architecture</div><div class="ib-text" data-i18n="ib-dag-t">Multiple blocks can be produced simultaneously and merged. Higher throughput, lower latency, better fault tolerance. Merge events marked with a special badge in the explorer.</div></div>
    <div><div class="ib-icon">⛽</div><div class="ib-title" data-i18n="ib-gas">Truly Gasless</div><div class="ib-text" data-i18n="ib-gas-t">Registration costs absolutely nothing. No ETH, BNB, or MATIC required. No credit card, no bank account. If you are a human with a smartphone, you can register.</div></div>
  </div>
</div>
<div class="main-grid">
  <div class="section">
    <div class="sec-head"><div class="sec-title"><span class="sec-dot"></span><span data-i18n="recent-blocks">Recent Blocks</span></div><div class="sec-count" id="block-count">—</div></div>
    <div class="sec-desc" data-i18n="blocks-desc">MERGE = multiple parents (BlockDAG feature). TX = registration transaction. Block time: ~6 seconds average.</div>
    <div id="blocks-list"><div class="empty" data-i18n="loading">Loading blocks...</div></div>
  </div>
  <div class="right-col">
    <div class="ic">
      <div class="ic-title" data-i18n="net-info">Network Info</div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-chain">Chain Name</span><span class="ic-val go">Aequitas Chain</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-chainid">Chain ID</span><span class="ic-val b">1926</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-symbol">Symbol</span><span class="ic-val go">AEQ</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-btime">Block Time</span><span class="ic-val">6 seconds</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-cons">Consensus</span><span class="ic-val p">BlockDAG + PoH</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-nodes">Active Nodes</span><span class="ic-val g">2 Online</span></div>
      <div class="ic-row"><span class="ic-key">ZKP</span><span class="ic-val">Groth16</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-storage">Storage</span><span class="ic-val g">PostgreSQL</span></div>
    </div>
    <div class="mm-card">
      <div class="mm-title" data-i18n="add-mm">ADD TO METAMASK</div>
      <div class="mm-row"><span class="mm-key" data-i18n="k-chain">Network Name</span><span class="mm-val">Aequitas Chain</span></div>
      <div class="mm-row"><span class="mm-key">RPC URL</span><span class="mm-val" style="font-size:0.5rem">...9fba.up.railway.app/rpc</span></div>
      <div class="mm-row"><span class="mm-key" data-i18n="k-chainid">Chain ID</span><span class="mm-val">1926</span></div>
      <div class="mm-row"><span class="mm-key" data-i18n="k-symbol">Symbol</span><span class="mm-val">AEQ</span></div>
      <div class="mm-row"><span class="mm-key" data-i18n="k-dec">Decimals</span><span class="mm-val">18</span></div>
      <button class="mm-btn" onclick="addToMetaMask()" data-i18n="btn-add-mm">+ ADD AEQUITAS NETWORK</button>
    </div>
    <div class="phil-card">
      <div class="phil-quote" data-i18n="phil">"Money exists because people exist.<br>Nothing more, nothing less."</div>
      <div class="phil-sub" data-i18n="phil-sub">— THE AEQUITAS PRINCIPLE —</div>
    </div>
  </div>
</div>
</div>

<!-- HUMANS -->
<div id="tab-humans" class="tab-content">
<div class="hero">
  <div class="section-label" data-i18n="humans-title">Verified Humans on Aequitas Chain</div>
  <div class="info-banner">
    <div><div class="ib-icon">🔒</div><div class="ib-title" data-i18n="h-what">What is a Verified Human?</div><div class="ib-text" data-i18n="h-what-t">A Verified Human is a wallet address cryptographically proven to belong to a unique living human. Biometric data is never transmitted or stored. Only a Zero-Knowledge Proof is used.</div></div>
    <div><div class="ib-icon">🧮</div><div class="ib-title" data-i18n="h-zkp">Zero-Knowledge Proof System</div><div class="ib-text" data-i18n="h-zkp-t">Aequitas uses the Groth16 proving system over BN128 elliptic curve. Proof size: ~200 bytes. Verification: ~10ms. The same system used by Zcash.</div></div>
    <div><div class="ib-icon">🛡</div><div class="ib-title" data-i18n="h-sybil">Sybil Attack Prevention</div><div class="ib-text" data-i18n="h-sybil-t">Each biometric hash is stored permanently. Attempting to register twice with the same fingerprint is immediately rejected. One human, one wallet, forever.</div></div>
    <div><div class="ib-icon">🌍</div><div class="ib-title" data-i18n="h-global">Global Inclusion</div><div class="ib-text" data-i18n="h-global-t">No bank account, no credit card, no cryptocurrency required. Just an Android smartphone with a fingerprint sensor — a device over 3 billion people already own.</div></div>
  </div>
</div>
<div class="hs">
  <div class="section">
    <div class="sec-head"><div class="sec-title"><span class="sec-dot"></span><span data-i18n="reg-humans">Registered Humans</span></div><div class="sec-count" id="h-count">0</div></div>
    <div class="sec-desc" data-i18n="h-desc">Every address verified as unique human through biometric ZKP. Each received 1,000 AEQ. Permanent, immutable, on-chain.</div>
    <div id="humans-list"><div class="empty" data-i18n="no-humans">No humans registered yet. Download the Aequitas Android App and be the first!</div></div>
  </div>
  <div class="right-col">
    <div class="ic">
      <div class="ic-title" data-i18n="reg-stats">Registry Stats</div>
      <div class="ic-row"><span class="ic-key" data-i18n="total-humans">Total Humans</span><span class="ic-val g" id="stat-humans">0</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="s-supply">Total Supply</span><span class="ic-val go" id="stat-supply">0 AEQ</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-grant">Grant per Human</span><span class="ic-val go">1,000 AEQ</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-fee">Registration Fee</span><span class="ic-val g" data-i18n="free">FREE</span></div>
      <div class="ic-row"><span class="ic-key">ZKP System</span><span class="ic-val">Groth16 / BN128</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-bio">Biometric Storage</span><span class="ic-val g" data-i18n="never-stored">Never stored</span></div>
    </div>
  </div>
</div>
</div>

<!-- INDEX -->
<div id="tab-index" class="tab-content">
<div class="is">
  <div class="idx" style="grid-column:1/-1">
    <div class="idx-title" data-i18n="idx-title">Aequitas Index — Real-Time Economic Equality Score</div>
    <div class="idx-desc" data-i18n="idx-desc">Calculated from the on-chain balance distribution of all verified humans. 0 = perfect equality. 100 = maximum inequality. The protocol uses this to automatically trigger redistribution.</div>
    <div style="display:grid;grid-template-columns:auto 1fr;gap:20px;align-items:center;margin-top:12px">
      <div><div class="idx-big" id="idx-score">—</div><div class="idx-lbl" data-i18n="curr-idx">Current Index</div></div>
      <div>
        <div class="bar-bg"><div class="bar-fill" id="idx-bar" style="width:0%"></div></div>
        <div class="bar-lbl"><span data-i18n="bar-0">0 — Perfect Equality</span><span>50</span><span data-i18n="bar-100">100 — Max Inequality</span></div>
        <div style="margin-top:8px;font-size:0.63rem;color:var(--muted);background:#080F1E;padding:8px;border-radius:6px" id="idx-phase-desc">—</div>
      </div>
    </div>
    <div class="mrow" style="grid-template-columns:repeat(4,1fr)">
      <div class="mbox"><div class="mval" id="idx-gini">—</div><div class="mlbl" data-i18n="gini">Gini Coefficient</div></div>
      <div class="mbox"><div class="mval" id="idx-supply2">—</div><div class="mlbl" data-i18n="s-supply">Total Supply</div></div>
      <div class="mbox"><div class="mval" id="idx-phase">—</div><div class="mlbl" data-i18n="phase">Protocol Phase</div></div>
      <div class="mbox"><div class="mval" id="idx-humans2">—</div><div class="mlbl" data-i18n="s-humans">Verified Humans</div></div>
    </div>
  </div>
  <div class="idx">
    <div class="idx-title" data-i18n="pools-title">Redistribution Pools</div>
    <div class="idx-desc" data-i18n="pools-desc">When inequality thresholds are exceeded, AEQ is automatically redirected. Controlled entirely by protocol logic.</div>
    <div class="mrow">
      <div class="mbox"><div class="mval" id="pool-v">—</div><div class="mlbl" data-i18n="vel-pool">Velocity Pool</div></div>
      <div class="mbox"><div class="mval" id="pool-l">—</div><div class="mlbl" data-i18n="liq-pool">Liquidity Pool</div></div>
      <div class="mbox"><div class="mval" id="pool-u">—</div><div class="mlbl" data-i18n="ubi-pool">UBI Pool</div></div>
      <div class="mbox"><div class="mval" id="pool-t">—</div><div class="mlbl" data-i18n="treasury">Treasury</div></div>
    </div>
  </div>
  <div class="idx">
    <div class="idx-title" data-i18n="phases-title">Protocol Phases</div>
    <div class="idx-desc" data-i18n="phases-desc">Transitions happen automatically — no governance vote required.</div>
    <table class="spect">
      <tr><td>Phase 0</td><td style="color:var(--green)" data-i18n="p0">Bootstrap · &lt;100 humans · Cap: 50x fairShare</td></tr>
      <tr><td>Phase 1</td><td style="color:var(--blue)" data-i18n="p1">Growth · 100-10,000 · Cap: 20x</td></tr>
      <tr><td>Phase 2</td><td style="color:var(--gold)" data-i18n="p2">Stability · 10k-1M · Cap: 10x</td></tr>
      <tr><td>Phase 3</td><td style="color:var(--purple)" data-i18n="p3">Maturity · 1M+ · Cap: 3x</td></tr>
    </table>
  </div>
  <div class="idx" style="grid-column:1/-1">
    <div class="idx-title" data-i18n="story-title">The Story of Aequitas — Why This Exists</div>
    <div class="story" data-i18n="story-text"><p>The year is 2009. Satoshi Nakamoto releases Bitcoin. For the first time, value can transfer between any two people without a bank. A genuine revolution. But something goes wrong almost immediately.</p><p>Early miners accumulate millions of coins at almost zero cost. By 2021, the top 1% of Bitcoin addresses control over 90% of all Bitcoin. Bitcoin's estimated Gini exceeds 0.85 — higher than any country on Earth. The cryptocurrency that was supposed to democratize finance created the most extreme wealth concentration in human history.</p><p><span style="color:var(--gold)">Aequitas</span> — Latin for "fairness" and "equality" — was created to answer: <em style="color:var(--gold)">"What would a cryptocurrency look like if designed from first principles to be fair to every human being?"</em></p><p>The answer: <strong style="color:var(--text)">Money exists because people exist. Therefore, every person should have an equal share of money simply by virtue of being human.</strong></p><p>The Aequitas network launched in June 2026. Currently in Phase 0. The goal: demonstrate that money can be distributed fairly, equality maintained through mathematical governance, and financial inclusion achieved at global scale.</p><p><em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em></p></div>
  </div>
</div>
</div>

<!-- NETWORK -->
<div id="tab-network" class="tab-content">
<div class="ns">
  <div class="nc" style="grid-column:1/-1">
    <div class="nc-title" data-i18n="nodes-title">Active Nodes — Current Network Topology</div>
    <div style="font-size:0.65rem;color:var(--muted);line-height:1.8;margin-bottom:12px" data-i18n="nodes-desc">The Aequitas network operates on two nodes in geographically distributed cloud environments. Both participate in block production, state synchronization, and API serving. They communicate via libp2p and sync via HTTP. Both share the same PostgreSQL database.</div>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:8px">
      <div class="nbox"><div class="nstat"><span class="ndot"></span><span data-i18n="node1">Node 1 — Railway (Primary)</span></div><div class="nurl">aequitas-production-9fba.up.railway.app</div><div class="ndesc" data-i18n="node1-desc">Primary API · Block producer · P2P bootstrap · PostgreSQL · RPC for MetaMask</div></div>
      <div class="nbox"><div class="nstat"><span class="ndot"></span><span data-i18n="node2">Node 2 — Render (Secondary)</span></div><div class="nurl">aequitas-node-2.onrender.com</div><div class="ndesc" data-i18n="node2-desc">Secondary API · Block producer · P2P peer · HTTP sync · Shared PostgreSQL</div></div>
    </div>
  </div>
  <div class="nc">
    <div class="nc-title" data-i18n="bootstrap-title">Bootstrap Node</div>
    <div style="font-size:0.63rem;color:var(--muted);line-height:1.8;margin-bottom:10px" data-i18n="bootstrap-desc">To run your own Aequitas node, connect to the bootstrap node using the libp2p multiaddress below.</div>
    <div class="bsbox">/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R</div>
  </div>
  <div class="nc">
    <div class="nc-title" data-i18n="tech-title">Technical Specifications</div>
    <table class="spect">
      <tr><td data-i18n="k-chainid">Chain ID</td><td>1926 (0x786)</td></tr>
      <tr><td>EVM</td><td style="color:var(--green)" data-i18n="evm-yes">Yes — JSON-RPC /rpc · MetaMask</td></tr>
      <tr><td data-i18n="k-btime">Block Time</td><td>~6 seconds</td></tr>
      <tr><td data-i18n="k-cons">Consensus</td><td style="color:var(--purple)">BlockDAG + PoH</td></tr>
      <tr><td>P2P</td><td>libp2p (Go)</td></tr>
      <tr><td>ZKP</td><td>Groth16 / snarkjs / circom</td></tr>
      <tr><td>Curve</td><td>BN128 (alt-bn128)</td></tr>
      <tr><td data-i18n="k-storage">Storage</td><td style="color:var(--green)">PostgreSQL</td></tr>
      <tr><td data-i18n="k-lang">Language</td><td>Go 1.24</td></tr>
      <tr><td data-i18n="k-src">Source</td><td><a href="https://github.com/hanoi96international-gif/Aequitas" target="_blank" style="color:var(--blue)">GitHub</a></td></tr>
    </table>
  </div>
  <div class="nc">
    <div class="nc-title" data-i18n="mm-config">MetaMask Configuration</div>
    <table class="spect">
      <tr><td data-i18n="k-chain">Network Name</td><td style="color:var(--gold)">Aequitas Chain</td></tr>
      <tr><td>RPC URL</td><td style="color:var(--blue);font-size:0.55rem">https://aequitas-production-9fba.up.railway.app/rpc</td></tr>
      <tr><td data-i18n="k-chainid">Chain ID</td><td style="color:var(--gold)">1926</td></tr>
      <tr><td data-i18n="k-symbol">Symbol</td><td style="color:var(--gold)">AEQ</td></tr>
      <tr><td data-i18n="k-dec">Decimals</td><td>18</td></tr>
    </table>
    <button class="mm-btn" onclick="addToMetaMask()" style="margin-top:12px" data-i18n="btn-add-mm">+ ADD TO METAMASK</button>
  </div>
</div>
</div>

<!-- PROTOCOL V6 -->
<div id="tab-protocol" class="tab-content">
<div class="ps">
  <div class="section-label" data-i18n="proto-label">Aequitas V6 Protocol — Technical Documentation</div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="ca-title">Contract Addresses</div>
    <div class="hlbox" data-i18n="ca-text">Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier (Groth16): 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Main):     0xD487544fE06DeD5025DF7bD45bdFba5e9ffadd3f<br>V5 (Sepolia legacy):   0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="poa-title">1. PROOF OF ALIVE</div>
    <div class="story" data-i18n="poa-text"><p>What happens to money when people die or disappear? In Bitcoin, millions of BTC are permanently lost. In Aequitas, if someone disappears, their AEQ eventually returns to the community.</p></div>
    <div class="hlbox" data-i18n="poa-box">Year 0-2: Normal usage<br>Year 2: Warning 1 — Guardian can respond<br>Year 2+60d: Warning 2<br>Year 2+120d: Warning 3<br>Year 2+180d: AEQ goes to PERSONAL ESCROW<br>Year 4: If still inactive — UBI Pool — distributed equally</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="guard-title">2. GUARDIAN SYSTEM</div>
    <div class="story" data-i18n="guard-text"><p>What if someone cannot access their device for months? A trusted Guardian can confirm they are still alive — without any transaction rights.</p></div>
    <div class="hlbox" data-i18n="guard-box">1 Guardian per human (another verified human)<br>Guardian can ONLY call confirmAlive() — zero transaction rights<br>Guardian CANNOT move funds or transfer AEQ<br>Max 3 wards per Guardian<br>7-day timelock on assignment (prevents forced assignment)<br>No circular relationships</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="dem-title">3. DEMURRAGE — Anti-Hoarding</div>
    <div class="hlbox" data-i18n="dem-box">1% annual fee on balance ABOVE your fairShare goes to UBI Pool<br><br>Example: fairShare = 1,000 AEQ · Your balance = 3,000 AEQ<br>Excess: 2,000 AEQ · Monthly fee: 1.67 AEQ to UBI Pool</div>
    <div class="story" style="margin-top:12px" data-i18n="dem-text"><p>Historical precedent: Worgl, Austria (1932) — demurrage currency reduced unemployment 25% in one year. The Central Bank shut it down because it worked too well.</p></div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="cap-title">4. WEALTH CAP</div>
    <div class="hlbox" data-i18n="cap-box">Phase 0: 50x fairShare · Phase 1: 20x · Phase 2: 10x · Phase 3: 5x · Phase 4: 3x<br><br>Always active from human #1. Excess is instantly redistributed equally to ALL active humans.</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="ubi-title">5. UNIVERSAL BASIC INCOME</div>
    <div class="hlbox" data-i18n="ubi-box">Sources: Transaction fees · Wealth cap overflow · Demurrage · Inactive escrow<br><br>Monthly: UBI Pool divided equally among all active humans</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="inf-title">6. NO ALGORITHMIC INFLATION</div>
    <div class="hlbox" data-i18n="inf-box">The ONLY event that creates new AEQ: a new verified human registers<br><br>Total AEQ = Verified Active Humans x 1,000 (always true, always verifiable)</div>
  </div>
</div>
</div>

<script>
const PS = 'https://aequitas-proof-server-production.up.railway.app';
const CID = '0x786';
let waddr = '', proofData = null, curLang = 'en';

const T = {
en:{
  'logo-sub':'PROOF OF HUMANITY','live':'LIVE',
  'tab-register':'🔐 Register','tab-explorer':'🔍 Explorer','tab-humans':'👥 Humans','tab-index':'📊 Index','tab-network':'🌐 Network','tab-protocol':'📜 Protocol V7',
  'reg-title':'🔐 Register as a Verified Human',
  'reg-sub':'Join the Aequitas network and receive 1,000 AEQ. One-time, permanent, gasless. No personal data stored.',
  'app-title':'REGISTRATION VIA ANDROID APP',
  'app-text':'Proof of Humanity requires biometric verification on your personal device. Your fingerprint is processed by the Hardware Secure Element — raw data never leaves your phone. Download the app, scan your fingerprint, connect your wallet, and your <strong style="color:var(--gold)">1,000 AEQ will be granted automatically</strong>.',
  's1t':'Biometric Scan','s1d':'Open app · scan fingerprint · HSE processes · data never leaves device',
  's2t':'ZKP Generation','s2d':'Groth16 proof generated · uniqueness verified · hash never revealed',
  's3t':'Connect Wallet','s3d':'App opens MetaMask · connect wallet · address receives 1,000 AEQ',
  's4t':'1,000 AEQ','s4d':'Registered on V6 · confirmed in next block · app notifies automatically',
  'priv-bar':'🔒 Hardware Secure Element · Groth16 ZKP · Data never leaves device · No gas fees · Permanent Sybil protection',
  'conn-wallet':'CONNECTED WALLET','proof-recv':'⚡ ZK PROOF RECEIVED','proof-hint':'Connect wallet to register',
  'btn-conn':'🦊 CONNECT METAMASK','btn-reg':'🔐 REGISTER ON-CHAIN',
  'reg-log-hint':'// Open Aequitas Android App to generate your proof, then return here...',
  'reg-details':'Registration Details','k-network':'Network','k-chainid':'Chain ID','k-grant':'Grant Amount',
  'k-fee':'Gas Fee','free':'FREE (gasless)','k-limit':'Registrations','k-limit-v':'Once per human · permanent · immutable',
  'k-bio':'Biometric Data','never-stored':'Never stored anywhere','k-conf':'Confirmation','k-conf-v':'Within 6 seconds',
  'live-stats':'Live Chain Statistics',
  's-height':'Block Height','s-height-sub':'New block every 6s · BlockDAG · Two nodes parallel',
  's-humans':'Verified Humans','s-humans-sub':'Biometric ZKP · One person, one wallet, forever',
  's-supply':'Total Supply','s-supply-sub':'Always = Humans x 1,000 AEQ',
  's-index':'Aequitas Index','s-index-sub':'0 = perfect equality · 100 = max inequality',
  's-uptime':'Uptime','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Proof of Humanity','ib-poh-t':'Every AEQ holder must prove they are a unique living human. No bots, no corporations, no AI can hold AEQ. Only real humans. Biometric data never leaves your device.',
  'ib-fair':'Radically Fair','ib-fair-t':'Every verified human receives exactly 1,000 AEQ. No pre-mine, no founder allocation. Total supply always equals verified humans x 1,000.',
  'ib-dag':'BlockDAG Architecture','ib-dag-t':'Multiple blocks can be produced simultaneously and merged. Higher throughput, lower latency, better fault tolerance.',
  'ib-gas':'Truly Gasless','ib-gas-t':'Registration costs absolutely nothing. No ETH, BNB, or MATIC required. No credit card, no bank account.',
  'recent-blocks':'Recent Blocks','blocks-desc':'MERGE = multiple parents (BlockDAG). TX = registration transaction. Block time: ~6 seconds.',
  'loading':'Loading blocks...','net-info':'Network Info','k-chain':'Chain Name','k-symbol':'Symbol','k-btime':'Block Time',
  'k-cons':'Consensus','k-nodes':'Active Nodes','k-storage':'Storage','add-mm':'🦊 ADD TO METAMASK','k-dec':'Decimals',
  'btn-add-mm':'+ ADD AEQUITAS NETWORK',
  'phil':'"Money exists because people exist.<br>Nothing more, nothing less."','phil-sub':'— THE AEQUITAS PRINCIPLE —',
  'humans-title':'Verified Humans on Aequitas Chain',
  'h-what':'What is a Verified Human?','h-what-t':'A Verified Human is a wallet address cryptographically proven to belong to a unique living human. Biometric data is never transmitted or stored.',
  'h-zkp':'Zero-Knowledge Proof System','h-zkp-t':'Aequitas uses the Groth16 proving system over BN128 elliptic curve. Proof size: ~200 bytes. Verification: ~10ms.',
  'h-sybil':'Sybil Attack Prevention','h-sybil-t':'Each biometric hash is stored permanently. Attempting to register twice is immediately rejected. One human, one wallet, forever.',
  'h-global':'Global Inclusion','h-global-t':'No bank account, no credit card, no cryptocurrency required. Just an Android smartphone with a fingerprint sensor.',
  'reg-humans':'Registered Humans','h-desc':'Every address verified as unique human through biometric ZKP. Each received 1,000 AEQ. Permanent, immutable, on-chain.',
  'no-humans':'No humans registered yet.\n\nDownload the Aequitas Android App and be the first human on the chain!',
  'reg-stats':'Registry Stats','total-humans':'Total Humans',
  'idx-title':'Aequitas Index — Real-Time Economic Equality Score',
  'idx-desc':'Calculated from the on-chain balance distribution of all verified humans. 0 = perfect equality. 100 = maximum inequality.',
  'curr-idx':'Current Index','bar-0':'0 — Perfect Equality','bar-100':'100 — Max Inequality',
  'gini':'Gini Coefficient','phase':'Protocol Phase',
  'pools-title':'Redistribution Pools','pools-desc':'When inequality thresholds are exceeded, AEQ is automatically redirected. Controlled entirely by protocol logic.',
  'vel-pool':'Velocity Pool','liq-pool':'Liquidity Pool','ubi-pool':'UBI Pool','treasury':'Treasury',
  'phases-title':'Protocol Phases','phases-desc':'Transitions happen automatically — no governance vote required.',
  'p0':'Bootstrap · &lt;100 humans · Cap: 50x','p1':'Growth · 100-10,000 · Cap: 20x',
  'p2':'Stability · 10k-1M · Cap: 10x','p3':'Maturity · 1M+ · Cap: 3x',
  'story-title':'The Story of Aequitas — Why This Exists',
  'story-text':'<p>The year is 2009. Satoshi Nakamoto releases Bitcoin. For the first time, value can transfer between any two people without a bank. A genuine revolution. But something goes wrong almost immediately.</p><p>Early miners accumulate millions of coins at almost zero cost. By 2021, the top 1% of Bitcoin addresses control over 90% of all Bitcoin. Bitcoin\'s estimated Gini exceeds 0.85 — higher than any country on Earth.</p><p><span style="color:var(--gold)">Aequitas</span> was created to answer: <em style="color:var(--gold)">"What would a cryptocurrency look like if designed from first principles to be fair to every human being?"</em></p><p>The answer: <strong>Money exists because people exist. Therefore, every person should have an equal share of money simply by virtue of being human.</strong></p><p><em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em></p>',
  'nodes-title':'Active Nodes — Current Network Topology',
  'nodes-desc':'The Aequitas network operates on two nodes in geographically distributed cloud environments. Both participate in block production, state synchronization, and API serving.',
  'node1':'Node 1 — Railway (Primary)','node1-desc':'Primary API · Block producer · P2P bootstrap · PostgreSQL · RPC for MetaMask',
  'node2':'Node 2 — Render (Secondary)','node2-desc':'Secondary API · Block producer · P2P peer · HTTP sync · Shared PostgreSQL',
  'bootstrap-title':'Bootstrap Node','bootstrap-desc':'To run your own Aequitas node, connect to the bootstrap node using the libp2p multiaddress below.',
  'tech-title':'Technical Specifications','mm-config':'MetaMask Configuration',
  'k-lang':'Language','k-src':'Source','evm-yes':'Yes — JSON-RPC /rpc · MetaMask',
  'proto-label':'Aequitas V6 Protocol — Technical Documentation',
  'ca-title':'Contract Addresses','ca-text':'Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV6:  0xD487544fE06DeD5025DF7bD45bdFba5e9ffadd3f<br>V5 Sepolia:  0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5',
  'poa-title':'1. PROOF OF ALIVE','poa-text':'<p>What happens to money when people die or disappear? In Bitcoin, millions of BTC are permanently lost. In Aequitas, if someone disappears, their AEQ eventually returns to the community.</p>',
  'poa-box':'Year 0-2: Normal usage<br>Year 2: Warning 1 — Guardian can respond<br>Year 2+60d: Warning 2<br>Year 2+120d: Warning 3<br>Year 2+180d: AEQ goes to PERSONAL ESCROW<br>Year 4: If still inactive — UBI Pool',
  'guard-title':'2. GUARDIAN SYSTEM','guard-text':'<p>What if someone cannot access their device for months? A trusted Guardian can confirm they are still alive — without any transaction rights.</p>',
  'guard-box':'1 Guardian per human (another verified human)<br>Guardian can ONLY call confirmAlive() — zero transaction rights<br>Guardian CANNOT move funds or transfer AEQ<br>Max 3 wards · 7-day timelock · No circular relationships',
  'dem-title':'3. DEMURRAGE — Anti-Hoarding',
  'dem-box':'1% annual fee on balance ABOVE fairShare goes to UBI Pool<br><br>Example: fairShare=1,000 · Balance=3,000 · Excess=2,000 · Monthly fee=1.67 AEQ',
  'dem-text':'<p>Historical precedent: Worgl, Austria (1932) — demurrage currency reduced unemployment 25% in one year.</p>',
  'cap-title':'4. WEALTH CAP','cap-box':'Phase 0: 50x fairShare · Phase 1: 20x · Phase 2: 10x · Phase 3: 5x · Phase 4: 3x<br><br>Always active from human #1. Excess instantly redistributed to ALL active humans.',
  'ubi-title':'5. UNIVERSAL BASIC INCOME','ubi-box':'Sources: Transaction fees · Wealth cap overflow · Demurrage · Inactive escrow<br><br>Monthly: UBI Pool divided equally among all active humans',
  'inf-title':'6. NO ALGORITHMIC INFLATION','inf-box':'The ONLY event that creates new AEQ: a new verified human registers<br><br>Total AEQ = Verified Active Humans x 1,000'
},
de:{
  'logo-sub':'MENSCHLICHKEITSNACHWEIS','live':'LIVE',
  'tab-register':'🔐 Registrieren','tab-explorer':'🔍 Explorer','tab-humans':'👥 Menschen','tab-index':'📊 Index','tab-network':'🌐 Netzwerk','tab-protocol':'📜 Protokoll V7',
  'reg-title':'🔐 Als verifizierter Mensch registrieren',
  'reg-sub':'Tritt dem Aequitas-Netzwerk bei und erhalte 1.000 AEQ. Einmalig, permanent, gebührenfrei. Keine persönlichen Daten gespeichert.',
  'app-title':'REGISTRIERUNG NUR ÜBER ANDROID-APP',
  'app-text':'Der Menschlichkeitsnachweis erfordert biometrische Verifizierung auf deinem Gerät. Dein Fingerabdruck wird durch das Hardware Secure Element verarbeitet — rohe Daten verlassen niemals dein Telefon. Lade die App herunter, scanne deinen Fingerabdruck, verbinde deine Wallet, und deine <strong style="color:var(--gold)">1.000 AEQ werden automatisch gewährt</strong>.',
  's1t':'Biometrischer Scan','s1d':'App öffnen · Fingerabdruck scannen · HSE verarbeitet · Daten verlassen nie das Gerät',
  's2t':'ZKP-Erzeugung','s2d':'Groth16-Beweis generiert · Einzigartigkeit verifiziert · Hash nie enthüllt',
  's3t':'Wallet verbinden','s3d':'App öffnet MetaMask · Wallet verbinden · Adresse erhält 1.000 AEQ',
  's4t':'1.000 AEQ','s4d':'Auf V6 registriert · im nächsten Block bestätigt · App benachrichtigt automatisch',
  'priv-bar':'🔒 Hardware Secure Element · Groth16 ZKP · Daten verlassen nie das Gerät · Keine Gasgebühren · Permanenter Sybil-Schutz',
  'conn-wallet':'VERBUNDENE WALLET','proof-recv':'⚡ ZK-BEWEIS EMPFANGEN','proof-hint':'Wallet verbinden um zu registrieren',
  'btn-conn':'🦊 METAMASK VERBINDEN','btn-reg':'🔐 ON-CHAIN REGISTRIEREN',
  'reg-log-hint':'// Öffne die Aequitas Android-App um deinen Beweis zu generieren, dann kehre hierher zurück...',
  'reg-details':'Registrierungsdetails','k-network':'Netzwerk','k-chainid':'Chain-ID','k-grant':'Zuteilung',
  'k-fee':'Gasgebühr','free':'KOSTENLOS','k-limit':'Registrierungen','k-limit-v':'Einmalig · permanent · unveränderlich',
  'k-bio':'Biometrische Daten','never-stored':'Niemals gespeichert','k-conf':'Bestätigung','k-conf-v':'Innerhalb von 6 Sekunden',
  'live-stats':'Live Chain-Statistiken',
  's-height':'Blockhöhe','s-height-sub':'Neuer Block alle 6 Sek · BlockDAG · Zwei Nodes parallel',
  's-humans':'Verifizierte Menschen','s-humans-sub':'Biometrischer ZKP · Eine Person, eine Wallet, für immer',
  's-supply':'Gesamtmenge','s-supply-sub':'Immer = Menschen x 1.000 AEQ',
  's-index':'Aequitas-Index','s-index-sub':'0 = vollkommene Gleichheit · 100 = maximale Ungleichheit',
  's-uptime':'Betriebszeit','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Menschlichkeitsnachweis','ib-poh-t':'Jeder AEQ-Inhaber muss beweisen, dass er ein einzigartiger lebender Mensch ist. Keine Bots, keine Unternehmen, keine KI. Nur echte Menschen.',
  'ib-fair':'Radikal faire Verteilung','ib-fair-t':'Jeder verifizierte Mensch erhält genau 1.000 AEQ. Keine Vorzuteilung, keine Gründeranteile. Gesamtmenge immer = Menschen x 1.000.',
  'ib-dag':'BlockDAG-Architektur','ib-dag-t':'Mehrere Blöcke können gleichzeitig produziert und zusammengeführt werden. Höherer Durchsatz, niedrigere Latenz.',
  'ib-gas':'Wirklich gebührenfrei','ib-gas-t':'Registrierung kostet absolut nichts. Kein ETH, BNB oder MATIC. Kein Bankkonto erforderlich.',
  'recent-blocks':'Aktuelle Blöcke','blocks-desc':'MERGE = mehrere Elternblöcke (BlockDAG). TX = Registrierungstransaktion. Blockzeit: ~6 Sekunden.',
  'loading':'Blöcke werden geladen...','net-info':'Netzwerkinformationen','k-chain':'Netzwerkname','k-symbol':'Symbol','k-btime':'Blockzeit',
  'k-cons':'Konsens','k-nodes':'Aktive Nodes','k-storage':'Speicher','add-mm':'🦊 ZU METAMASK HINZUFÜGEN','k-dec':'Dezimalstellen',
  'btn-add-mm':'+ AEQUITAS-NETZWERK HINZUFÜGEN',
  'phil':'"Geld existiert weil Menschen existieren.<br>Nichts mehr, nichts weniger."','phil-sub':'— DAS AEQUITAS-PRINZIP —',
  'humans-title':'Verifizierte Menschen auf der Aequitas Chain',
  'h-what':'Was ist ein verifizierter Mensch?','h-what-t':'Ein verifizierter Mensch ist eine Wallet-Adresse, die kryptographisch bewiesen wurde, einem einzigartigen lebenden Menschen zu gehören.',
  'h-zkp':'Zero-Knowledge-Proof-System','h-zkp-t':'Aequitas verwendet das Groth16-Beweissystem über die BN128-elliptische Kurve. Beweisdauer: ~200 Bytes. Verifizierungszeit: ~10ms.',
  'h-sybil':'Schutz vor Sybil-Angriffen','h-sybil-t':'Jeder biometrische Hash wird dauerhaft gespeichert. Doppelregistrierung wird sofort abgelehnt. Eine Person, eine Wallet, für immer.',
  'h-global':'Globale Inklusion','h-global-t':'Kein Bankkonto, keine Kreditkarte, keine Kryptowährung erforderlich. Nur ein Android-Smartphone mit Fingerabdrucksensor.',
  'reg-humans':'Registrierte Menschen','h-desc':'Jede Adresse als einzigartiger Mensch durch biometrischen ZKP verifiziert. Jeder erhielt 1.000 AEQ. Dauerhaft, unveränderlich.',
  'no-humans':'Noch keine Menschen registriert.\n\nLade die Aequitas Android-App herunter und sei der erste!',
  'reg-stats':'Registrierungsstatistik','total-humans':'Gesamte Menschen',
  'idx-title':'Aequitas-Index — Wirtschaftlicher Gleichheitswert in Echtzeit',
  'idx-desc':'Berechnet aus der On-Chain-Bilanzverteilung aller verifizierten Menschen. 0 = vollkommene Gleichheit. 100 = maximale Ungleichheit.',
  'curr-idx':'Aktueller Index','bar-0':'0 — Vollkommene Gleichheit','bar-100':'100 — Max. Ungleichheit',
  'gini':'Gini-Koeffizient','phase':'Protokollphase',
  'pools-title':'Umverteilungspools','pools-desc':'Wenn Ungleichheitsschwellenwerte überschritten werden, wird AEQ automatisch umgeleitet.',
  'vel-pool':'Velocity-Pool','liq-pool':'Liquiditäts-Pool','ubi-pool':'UBI-Pool','treasury':'Tresor',
  'phases-title':'Protokollphasen','phases-desc':'Übergänge erfolgen automatisch — keine Abstimmung erforderlich.',
  'p0':'Bootstrap · &lt;100 Menschen · Cap: 50x','p1':'Wachstum · 100-10.000 · Cap: 20x',
  'p2':'Stabilität · 10k-1M · Cap: 10x','p3':'Reife · 1M+ · Cap: 3x',
  'story-title':'Die Geschichte von Aequitas — Warum das existiert',
  'story-text':'<p>Das Jahr ist 2009. Satoshi Nakamoto veröffentlicht Bitcoin. Zum ersten Mal können Werte zwischen zwei Menschen übertragen werden ohne Banken. Eine echte Revolution. Aber fast sofort geht etwas schief.</p><p>Frühe Miner häufen Millionen von Coins an die sie fast nichts kosten. Bis 2021 kontrolliert das oberste 1% der Bitcoin-Adressen über 90% aller Bitcoins.</p><p><span style="color:var(--gold)">Aequitas</span> — Lateinisch für "Fairness" und "Gleichheit" — wurde geschaffen um zu antworten: <em style="color:var(--gold)">"Wie würde eine Kryptowährung aussehen die von Grund auf fair für jeden Menschen konzipiert wurde?"</em></p><p>Die Antwort: <strong>Geld existiert weil Menschen existieren. Daher sollte jeder Mensch einfach aufgrund seiner Menschlichkeit einen gleichen Anteil am Geld haben.</strong></p><p><em style="color:var(--gold)">"Geld existiert weil Menschen existieren. Nichts mehr, nichts weniger."</em></p>',
  'nodes-title':'Aktive Nodes — Aktuelle Netzwerktopologie',
  'nodes-desc':'Das Aequitas-Netzwerk betreibt zwei Nodes in geografisch verteilten Cloud-Umgebungen. Beide nehmen an Blockproduktion, Statussynchronisation und API-Bereitstellung teil.',
  'node1':'Node 1 — Railway (Primär)','node1-desc':'Primärer API-Server · Blockproduzent · P2P-Bootstrap · PostgreSQL · RPC für MetaMask',
  'node2':'Node 2 — Render (Sekundär)','node2-desc':'Sekundärer API-Server · Blockproduzent · P2P-Peer · HTTP-Sync · Geteiltes PostgreSQL',
  'bootstrap-title':'Bootstrap-Node','bootstrap-desc':'Um deinen eigenen Aequitas-Node zu betreiben, verbinde dich über die unten stehende libp2p-Multiadresse.',
  'tech-title':'Technische Spezifikationen','mm-config':'MetaMask-Konfiguration',
  'k-lang':'Sprache','k-src':'Quellcode','evm-yes':'Ja — JSON-RPC /rpc · MetaMask',
  'proto-label':'Aequitas V6 Protokoll — Technische Dokumentation',
  'ca-title':'Contract-Adressen','ca-text':'Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV6:  0xD487544fE06DeD5025DF7bD45bdFba5e9ffadd3f<br>V5 Sepolia:  0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5',
  'poa-title':'1. LEBENSNACHWEIS','poa-text':'<p>Was passiert mit Geld wenn Menschen sterben oder verschwinden? Bei Bitcoin sind Millionen BTC dauerhaft verloren. Bei Aequitas kehrt das AEQ einer verschwundenen Person schließlich zur Gemeinschaft zurück.</p>',
  'poa-box':'Jahr 0-2: Normale Nutzung<br>Jahr 2: Warnung 1 — Guardian kann antworten<br>Jahr 2+60T: Warnung 2<br>Jahr 2+120T: Warnung 3<br>Jahr 2+180T: AEQ ins persönliche Treuhand-Konto<br>Jahr 4: Bei weiter Inaktivität — UBI-Pool',
  'guard-title':'2. GUARDIAN-SYSTEM','guard-text':'<p>Was wenn jemand monatelang nicht auf sein Gerät zugreifen kann? Ein vertrauenswürdiger Guardian kann bestätigen dass sie noch am Leben sind — ohne Transaktionsrechte.</p>',
  'guard-box':'1 Guardian pro Mensch · Nur confirmAlive() erlaubt · Keine Transaktionsrechte<br>Max 3 Schutzbefohlene · 7-Tage-Zeitsperre · Keine Kreisbeziehungen',
  'dem-title':'3. DEMURRAGE — Anti-Hortung',
  'dem-box':'1% jährliche Gebühr auf Guthaben ÜBER fairShare geht in UBI-Pool<br><br>Beispiel: fairShare=1.000 · Guthaben=3.000 · Überschuss=2.000 · Monatliche Gebühr=1,67 AEQ',
  'dem-text':'<p>Historisches Beispiel: Wörgl, Österreich (1932) — Demurrage-Währung reduzierte die Arbeitslosigkeit in einem Jahr um 25%.</p>',
  'cap-title':'4. VERMÖGENSOBERGRENZE',
  'cap-box':'Phase 0: 50x fairShare · Phase 1: 20x · Phase 2: 10x · Phase 3: 5x · Phase 4: 3x<br><br>Immer aktiv ab Mensch #1. Überschuss wird sofort gleichmäßig an ALLE aktiven Menschen verteilt.',
  'ubi-title':'5. UNIVERSELLES GRUNDEINKOMMEN',
  'ubi-box':'Quellen: Transaktionsgebühren · Vermögensobergrenze-Überschuss · Demurrage · Inaktive Treuhand<br><br>Monatlich: UBI-Pool gleichmäßig unter allen aktiven Menschen verteilt',
  'inf-title':'6. KEINE ALGORITHMISCHE INFLATION',
  'inf-box':'Das EINZIGE Ereignis das neues AEQ schafft: ein neuer verifizierter Mensch registriert sich<br><br>Gesamt-AEQ = Verifizierte aktive Menschen x 1.000'
},
es:{
  'logo-sub':'PRUEBA DE HUMANIDAD','live':'EN VIVO',
  'tab-register':'🔐 Registrar','tab-explorer':'🔍 Explorador','tab-humans':'👥 Humanos','tab-index':'📊 Índice','tab-network':'🌐 Red','tab-protocol':'📜 Protocolo V7',
  'reg-title':'🔐 Regístrate como Humano Verificado',
  'reg-sub':'Únete a la red Aequitas y recibe 1,000 AEQ. Único, permanente, sin gas. Sin datos personales.',
  'app-title':'REGISTRO SOLO VÍA APP ANDROID',
  'app-text':'La Prueba de Humanidad requiere verificación biométrica en tu dispositivo. Tu huella se procesa por el Elemento Seguro de Hardware — los datos brutos nunca salen de tu teléfono. Descarga la app, escanea tu huella, conecta tu wallet, y tus <strong style="color:var(--gold)">1,000 AEQ serán otorgados automáticamente</strong>.',
  's1t':'Escaneo Biométrico','s1d':'Abrir app · escanear huella · HSE procesa · datos nunca salen del dispositivo',
  's2t':'Generación ZKP','s2d':'Prueba Groth16 generada · unicidad verificada · hash nunca revelado',
  's3t':'Conectar Wallet','s3d':'App abre MetaMask · conectar wallet · dirección recibe 1,000 AEQ',
  's4t':'1,000 AEQ','s4d':'Registrado en V6 · confirmado en próximo bloque · app notifica automáticamente',
  'priv-bar':'🔒 Elemento Seguro de Hardware · ZKP Groth16 · Datos nunca salen del dispositivo · Sin tarifas de gas',
  'conn-wallet':'WALLET CONECTADA','proof-recv':'⚡ PRUEBA ZK RECIBIDA','proof-hint':'Conecta wallet para registrar',
  'btn-conn':'🦊 CONECTAR METAMASK','btn-reg':'🔐 REGISTRAR ON-CHAIN',
  'reg-log-hint':'// Abre la App Android Aequitas para generar tu prueba, luego regresa aquí...',
  'reg-details':'Detalles del Registro','k-network':'Red','k-chainid':'ID de Cadena','k-grant':'Bono',
  'k-fee':'Tarifa de Gas','free':'GRATIS','k-limit':'Registros','k-limit-v':'Una vez · permanente · inmutable',
  'k-bio':'Datos Biométricos','never-stored':'Nunca almacenados','k-conf':'Confirmación','k-conf-v':'En 6 segundos',
  'live-stats':'Estadísticas de Cadena en Vivo',
  's-height':'Altura de Bloque','s-height-sub':'Nuevo bloque cada 6s · BlockDAG · Dos nodos paralelos',
  's-humans':'Humanos Verificados','s-humans-sub':'ZKP biométrico · Una persona, una wallet, siempre',
  's-supply':'Suministro Total','s-supply-sub':'Siempre = Humanos x 1,000 AEQ',
  's-index':'Índice Aequitas','s-index-sub':'0 = igualdad perfecta · 100 = desigualdad máxima',
  's-uptime':'Tiempo Activo','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Prueba de Humanidad','ib-poh-t':'Cada titular de AEQ debe probar que es un humano único vivo. Sin bots, sin corporaciones, sin IA. Solo humanos reales.',
  'ib-fair':'Distribución Radicalmente Justa','ib-fair-t':'Cada humano verificado recibe exactamente 1,000 AEQ. Sin pre-minado, sin asignación a fundadores.',
  'ib-dag':'Arquitectura BlockDAG','ib-dag-t':'Múltiples bloques pueden producirse simultáneamente y fusionarse. Mayor rendimiento, menor latencia.',
  'ib-gas':'Verdaderamente Sin Gas','ib-gas-t':'El registro no cuesta nada. No se necesita ETH, BNB ni MATIC. Sin cuenta bancaria.',
  'recent-blocks':'Bloques Recientes','blocks-desc':'MERGE = múltiples padres (BlockDAG). TX = transacción de registro. Tiempo de bloque: ~6 segundos.',
  'loading':'Cargando bloques...','net-info':'Información de Red','k-chain':'Nombre de Red','k-symbol':'Símbolo','k-btime':'Tiempo de Bloque',
  'k-cons':'Consenso','k-nodes':'Nodos Activos','k-storage':'Almacenamiento','add-mm':'🦊 AGREGAR A METAMASK','k-dec':'Decimales',
  'btn-add-mm':'+ AGREGAR RED AEQUITAS',
  'phil':'"El dinero existe porque las personas existen.<br>Nada más, nada menos."','phil-sub':'— EL PRINCIPIO AEQUITAS —',
  'humans-title':'Humanos Verificados en Aequitas Chain',
  'h-what':'¿Qué es un Humano Verificado?','h-what-t':'Un Humano Verificado es una dirección wallet demostrada criptográficamente que pertenece a un humano único vivo.',
  'h-zkp':'Sistema ZKP','h-zkp-t':'Aequitas usa Groth16 sobre BN128. Tamaño de prueba: ~200 bytes. Verificación: ~10ms.',
  'h-sybil':'Prevención de Ataques Sybil','h-sybil-t':'Cada hash biométrico se almacena permanentemente. Intentar registrarse dos veces se rechaza inmediatamente.',
  'h-global':'Inclusión Global','h-global-t':'Sin cuenta bancaria, tarjeta de crédito ni criptomoneda. Solo un smartphone Android con sensor de huella.',
  'reg-humans':'Humanos Registrados','h-desc':'Cada dirección verificada como humano único mediante ZKP biométrico. Cada uno recibió 1,000 AEQ. Permanente.',
  'no-humans':'No hay humanos registrados aún.\n\n¡Descarga la App Android Aequitas y sé el primero!',
  'reg-stats':'Estadísticas del Registro','total-humans':'Total de Humanos',
  'idx-title':'Índice Aequitas — Puntuación de Igualdad Económica en Tiempo Real',
  'idx-desc':'Calculado desde la distribución de saldos on-chain. 0 = igualdad perfecta. 100 = desigualdad máxima.',
  'curr-idx':'Índice Actual','bar-0':'0 — Igualdad Perfecta','bar-100':'100 — Máx. Desigualdad',
  'gini':'Coeficiente Gini','phase':'Fase del Protocolo',
  'pools-title':'Pools de Redistribución','pools-desc':'Cuando se superan los umbrales de desigualdad, AEQ se redirige automáticamente.',
  'vel-pool':'Pool Velocidad','liq-pool':'Pool Liquidez','ubi-pool':'Pool UBI','treasury':'Tesorería',
  'phases-title':'Fases del Protocolo','phases-desc':'Las transiciones ocurren automáticamente — no se requiere votación.',
  'p0':'Bootstrap · &lt;100 humanos · Cap: 50x','p1':'Crecimiento · 100-10,000 · Cap: 20x',
  'p2':'Estabilidad · 10k-1M · Cap: 10x','p3':'Madurez · 1M+ · Cap: 3x',
  'story-title':'La Historia de Aequitas','story-text':'<p>El año es 2009. Satoshi Nakamoto lanza Bitcoin. Por primera vez, el valor puede transferirse sin bancos. Una revolución genuina. Pero casi de inmediato algo sale mal.</p><p>Los primeros mineros acumulan millones de monedas casi gratis. Para 2021, el 1% superior controla más del 90% de todo el Bitcoin.</p><p><span style="color:var(--gold)">Aequitas</span> fue creado para responder: <em style="color:var(--gold)">"¿Cómo sería una criptomoneda diseñada para ser justa con todo ser humano?"</em></p><p>La respuesta: <strong>El dinero existe porque las personas existen. Por lo tanto, cada persona debería tener una parte igual del dinero.</strong></p><p><em style="color:var(--gold)">"El dinero existe porque las personas existen. Nada más, nada menos."</em></p>',
  'nodes-title':'Nodos Activos','nodes-desc':'La red Aequitas opera en dos nodos en entornos cloud distribuidos.',
  'node1':'Nodo 1 — Railway (Primario)','node1-desc':'Servidor API primario · Productor de bloques · P2P bootstrap · PostgreSQL · RPC para MetaMask',
  'node2':'Nodo 2 — Render (Secundario)','node2-desc':'Servidor API secundario · Productor de bloques · Par P2P · Sincronización HTTP · PostgreSQL compartido',
  'bootstrap-title':'Nodo Bootstrap','bootstrap-desc':'Para ejecutar tu propio nodo Aequitas, conéctate al nodo bootstrap usando la multidirección libp2p.',
  'tech-title':'Especificaciones Técnicas','mm-config':'Configuración MetaMask',
  'k-lang':'Lenguaje','k-src':'Código Fuente','evm-yes':'Sí — JSON-RPC /rpc · MetaMask',
  'proto-label':'Protocolo Aequitas V6 — Documentación Técnica',
  'ca-title':'Direcciones de Contratos','ca-text':'Cadena: Aequitas Chain (Chain ID: 1926)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV6:  0xD487544fE06DeD5025DF7bD45bdFba5e9ffadd3f',
  'poa-title':'1. PRUEBA DE VIDA','poa-text':'<p>En Aequitas, si alguien desaparece, su AEQ eventualmente regresa a la comunidad.</p>',
  'poa-box':'Año 0-2: Uso normal<br>Año 2: Advertencia 1<br>Año 2+180d: AEQ a depósito personal<br>Año 4: Si inactivo — Pool UBI',
  'guard-title':'2. SISTEMA GUARDIAN','guard-text':'<p>Un Guardian de confianza puede confirmar que alguien está vivo — sin derechos de transacción.</p>',
  'guard-box':'1 Guardian por humano · Solo confirmAlive() · Sin derechos de transacción<br>Máx 3 pupilos · Bloqueo 7 días · Sin relaciones circulares',
  'dem-title':'3. DEMURRAGE','dem-box':'1% anual sobre saldo POR ENCIMA de fairShare va al Pool UBI',
  'dem-text':'<p>Precedente: Wörgl, Austria (1932) — redujo el desempleo 25% en un año.</p>',
  'cap-title':'4. LÍMITE DE RIQUEZA','cap-box':'Fase 0: 50x · Fase 1: 20x · Fase 2: 10x · Fase 3: 5x · Fase 4: 3x fairShare',
  'ubi-title':'5. INGRESO BÁSICO UNIVERSAL','ubi-box':'Fuentes: Comisiones · Desbordamiento de límite · Demurrage · Custodia inactiva',
  'inf-title':'6. SIN INFLACIÓN ALGORÍTMICA','inf-box':'El ÚNICO evento que crea AEQ: un nuevo humano verificado se registra'
},
ru:{
  'logo-sub':'ДОКАЗАТЕЛЬСТВО ЧЕЛОВЕЧНОСТИ','live':'В ЭФИРЕ',
  'tab-register':'🔐 Регистрация','tab-explorer':'🔍 Проводник','tab-humans':'👥 Люди','tab-index':'📊 Индекс','tab-network':'🌐 Сеть','tab-protocol':'📜 Протокол V7',
  'reg-title':'🔐 Зарегистрируйтесь как Верифицированный Человек',
  'reg-sub':'Присоединитесь к сети Aequitas и получите 1 000 AEQ. Одноразово, постоянно, бесплатно. Никаких личных данных.',
  'app-title':'РЕГИСТРАЦИЯ ТОЛЬКО ЧЕРЕЗ ANDROID-ПРИЛОЖЕНИЕ',
  'app-text':'Доказательство человечности требует биометрической верификации на вашем устройстве. Ваш отпечаток обрабатывается Hardware Secure Element — сырые данные никогда не покидают ваш телефон. Скачайте приложение, отсканируйте отпечаток, подключите кошелёк, и ваши <strong style="color:var(--gold)">1 000 AEQ будут начислены автоматически</strong>.',
  's1t':'Биометрический Скан','s1d':'Открыть приложение · сканировать отпечаток · HSE обрабатывает · данные не покидают устройство',
  's2t':'Генерация ZKP','s2d':'Сгенерировано доказательство Groth16 · уникальность верифицирована · хэш не раскрывается',
  's3t':'Подключить Кошелёк','s3d':'Приложение открывает MetaMask · подключить кошелёк · адрес получает 1 000 AEQ',
  's4t':'1 000 AEQ','s4d':'Зарегистрировано на V6 · подтверждено в следующем блоке · приложение уведомляет автоматически',
  'priv-bar':'🔒 Hardware Secure Element · ZKP Groth16 · Данные не покидают устройство · Без комиссий',
  'conn-wallet':'ПОДКЛЮЧЁННЫЙ КОШЕЛЁК','proof-recv':'⚡ ZK-ДОКАЗАТЕЛЬСТВО ПОЛУЧЕНО','proof-hint':'Подключите кошелёк для регистрации',
  'btn-conn':'🦊 ПОДКЛЮЧИТЬ METAMASK','btn-reg':'🔐 ЗАРЕГИСТРИРОВАТЬСЯ ON-CHAIN',
  'reg-log-hint':'// Откройте приложение Aequitas для генерации доказательства, затем вернитесь...',
  'reg-details':'Детали Регистрации','k-network':'Сеть','k-chainid':'ID Цепочки','k-grant':'Грант',
  'k-fee':'Комиссия','free':'БЕСПЛАТНО','k-limit':'Регистрации','k-limit-v':'Один раз · постоянно · неизменно',
  'k-bio':'Биометрические Данные','never-stored':'Никогда не хранятся','k-conf':'Подтверждение','k-conf-v':'В течение 6 секунд',
  'live-stats':'Статистика в реальном времени',
  's-height':'Высота Блока','s-height-sub':'Новый блок каждые 6 сек · BlockDAG · Два узла параллельно',
  's-humans':'Верифицированных Людей','s-humans-sub':'Биометрический ZKP · Один человек, один кошелёк, навсегда',
  's-supply':'Общее Предложение','s-supply-sub':'Всегда = Люди x 1 000 AEQ',
  's-index':'Индекс Aequitas','s-index-sub':'0 = полное равенство · 100 = максимальное неравенство',
  's-uptime':'Время Работы','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Доказательство Человечности','ib-poh-t':'Каждый владелец AEQ должен доказать что он уникальный живой человек. Никаких ботов, корпораций, ИИ. Только настоящие люди.',
  'ib-fair':'Радикально Справедливое Распределение','ib-fair-t':'Каждый верифицированный человек получает ровно 1 000 AEQ. Общее предложение всегда = люди x 1 000.',
  'ib-dag':'Архитектура BlockDAG','ib-dag-t':'Несколько блоков могут производиться одновременно и объединяться. Более высокая пропускная способность.',
  'ib-gas':'По-Настоящему Бесплатно','ib-gas-t':'Регистрация не стоит ничего. Не нужен ETH, BNB или MATIC. Не нужен банковский счёт.',
  'recent-blocks':'Последние Блоки','blocks-desc':'MERGE = несколько родителей (BlockDAG). TX = транзакция регистрации. Время блока: ~6 секунд.',
  'loading':'Загрузка блоков...','net-info':'Информация о Сети','k-chain':'Название Сети','k-symbol':'Символ','k-btime':'Время Блока',
  'k-cons':'Консенсус','k-nodes':'Активные Узлы','k-storage':'Хранение','add-mm':'🦊 ДОБАВИТЬ В METAMASK','k-dec':'Знаков',
  'btn-add-mm':'+ ДОБАВИТЬ СЕТЬ AEQUITAS',
  'phil':'"Деньги существуют потому что существуют люди.<br>Ничего больше, ничего меньше."','phil-sub':'— ПРИНЦИП AEQUITAS —',
  'humans-title':'Верифицированные Люди на Aequitas Chain',
  'h-what':'Что такое Верифицированный Человек?','h-what-t':'Верифицированный Человек — это адрес кошелька, доказанно принадлежащий уникальному живому человеку.',
  'h-zkp':'Система ZKP','h-zkp-t':'Aequitas использует Groth16 над BN128. Размер доказательства: ~200 байт. Верификация: ~10мс.',
  'h-sybil':'Защита от Атак Сивиллы','h-sybil-t':'Каждый биометрический хэш хранится постоянно. Двойная регистрация немедленно отклоняется.',
  'h-global':'Глобальное Включение','h-global-t':'Не нужен банковский счёт, кредитная карта или криптовалюта. Только Android-смартфон с сенсором отпечатка.',
  'reg-humans':'Зарегистрированных Людей','h-desc':'Каждый адрес верифицирован через биометрический ZKP. Каждый получил 1 000 AEQ. Постоянно, неизменно.',
  'no-humans':'Людей ещё нет.\n\nСкачай приложение Aequitas Android и стань первым!',
  'reg-stats':'Статистика Реестра','total-humans':'Всего Людей',
  'idx-title':'Индекс Aequitas — Оценка Экономического Равенства в Реальном Времени',
  'idx-desc':'Рассчитывается из распределения балансов on-chain. 0 = полное равенство. 100 = максимальное неравенство.',
  'curr-idx':'Текущий Индекс','bar-0':'0 — Полное Равенство','bar-100':'100 — Макс. Неравенство',
  'gini':'Коэффициент Джини','phase':'Фаза Протокола',
  'pools-title':'Пулы Перераспределения','pools-desc':'Когда пороги неравенства превышены, AEQ автоматически перенаправляется.',
  'vel-pool':'Пул Скорости','liq-pool':'Пул Ликвидности','ubi-pool':'Пул UBI','treasury':'Казначейство',
  'phases-title':'Фазы Протокола','phases-desc':'Переходы происходят автоматически — голосование не требуется.',
  'p0':'Загрузка · &lt;100 людей · Cap: 50x','p1':'Рост · 100-10 000 · Cap: 20x',
  'p2':'Стабильность · 10k-1M · Cap: 10x','p3':'Зрелость · 1M+ · Cap: 3x',
  'story-title':'История Aequitas — Почему Это Существует',
  'story-text':'<p>2009 год. Сатоши Накамото выпускает Биткоин. Впервые ценность можно передавать без банков. Революция. Но почти сразу что-то идёт не так.</p><p>Ранние майнеры накапливают миллионы монет почти бесплатно. К 2021 году верхний 1% адресов контролирует более 90% всех Биткоинов.</p><p><span style="color:var(--gold)">Aequitas</span> был создан чтобы ответить: <em style="color:var(--gold)">"Как выглядела бы криптовалюта, разработанная для справедливости к каждому человеку?"</em></p><p>Ответ: <strong>Деньги существуют потому что существуют люди. Поэтому каждый человек должен иметь равную долю денег просто будучи человеком.</strong></p><p><em style="color:var(--gold)">"Деньги существуют потому что существуют люди. Ничего больше, ничего меньше."</em></p>',
  'nodes-title':'Активные Узлы','nodes-desc':'Сеть Aequitas работает на двух узлах в географически распределённых облачных средах.',
  'node1':'Узел 1 — Railway (Основной)','node1-desc':'Основной API · Производитель блоков · P2P-bootstrap · PostgreSQL · RPC для MetaMask',
  'node2':'Узел 2 — Render (Вторичный)','node2-desc':'Вторичный API · Производитель блоков · P2P-пир · HTTP-синхронизация · Общий PostgreSQL',
  'bootstrap-title':'Bootstrap-Узел','bootstrap-desc':'Для запуска собственного узла подключитесь через libp2p мультиадрес.',
  'tech-title':'Технические Характеристики','mm-config':'Настройка MetaMask',
  'k-lang':'Язык','k-src':'Исходный Код','evm-yes':'Да — JSON-RPC /rpc · MetaMask',
  'proto-label':'Протокол Aequitas V6 — Техническая Документация',
  'ca-title':'Адреса Контрактов','ca-text':'Цепочка: Aequitas Chain (Chain ID: 1926)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV6:  0xD487544fE06DeD5025DF7bD45bdFba5e9ffadd3f',
  'poa-title':'1. ДОКАЗАТЕЛЬСТВО ЖИЗНИ','poa-text':'<p>В Aequitas если кто-то исчезает, его AEQ в конечном итоге возвращается сообществу.</p>',
  'poa-box':'Год 0-2: Нормальное использование<br>Год 2: Предупреждение 1<br>Год 2+180д: AEQ на персональный эскроу<br>Год 4: При инактивности — Пул UBI',
  'guard-title':'2. СИСТЕМА GUARDIAN','guard-text':'<p>Доверенный Guardian может подтвердить что человек жив — без прав на транзакции.</p>',
  'guard-box':'1 Guardian на человека · Только confirmAlive() · Без прав транзакций<br>Макс 3 подопечных · Блокировка 7 дней · Без круговых отношений',
  'dem-title':'3. ДЕМУРРЕДЖ','dem-box':'1% годовых на баланс ВЫШЕ fairShare идёт в Пул UBI',
  'dem-text':'<p>Исторический прецедент: Вёргль, Австрия (1932) — сократил безработицу на 25% за год.</p>',
  'cap-title':'4. ОГРАНИЧЕНИЕ БОГАТСТВА','cap-box':'Фаза 0: 50x · Фаза 1: 20x · Фаза 2: 10x · Фаза 3: 5x · Фаза 4: 3x fairShare',
  'ubi-title':'5. БАЗОВЫЙ ДОХОД','ubi-box':'Источники: Комиссии · Переполнение лимита · Демурредж · Неактивный эскроу',
  'inf-title':'6. БЕЗ АЛГОРИТМИЧЕСКОЙ ИНФЛЯЦИИ','inf-box':'Единственное событие создающее AEQ: регистрация нового верифицированного человека'
},
zh:{
  'logo-sub':'人类证明','live':'直播',
  'tab-register':'🔐 注册','tab-explorer':'🔍 浏览器','tab-humans':'👥 人类','tab-index':'📊 指数','tab-network':'🌐 网络','tab-protocol':'📜 协议 V7',
  'reg-title':'🔐 注册为已验证人类',
  'reg-sub':'加入Aequitas网络并接收1,000 AEQ。一次性、永久、无Gas费。不存储个人数据。',
  'app-title':'仅通过ANDROID应用注册',
  'app-text':'人类证明需要在您的个人设备上进行生物特征验证。您的指纹由硬件安全元件处理——原始数据永远不会离开您的手机。下载应用，扫描指纹，连接钱包，您的<strong style="color:var(--gold)">1,000 AEQ将自动发放</strong>。',
  's1t':'生物特征扫描','s1d':'打开应用 · 扫描指纹 · HSE处理 · 数据永不离开设备',
  's2t':'ZKP生成','s2d':'生成Groth16证明 · 验证唯一性 · 哈希从不泄露',
  's3t':'连接钱包','s3d':'应用打开MetaMask · 连接钱包 · 地址接收1,000 AEQ',
  's4t':'1,000 AEQ','s4d':'在V6上注册 · 下一个区块内确认 · 应用自动通知',
  'priv-bar':'🔒 硬件安全元件 · Groth16 ZKP · 数据永不离开设备 · 无Gas费',
  'conn-wallet':'已连接钱包','proof-recv':'⚡ 已收到ZK证明','proof-hint':'连接钱包以注册',
  'btn-conn':'🦊 连接METAMASK','btn-reg':'🔐 链上注册',
  'reg-log-hint':'// 打开Aequitas Android应用生成您的证明，然后返回此处...',
  'reg-details':'注册详情','k-network':'网络','k-chainid':'链ID','k-grant':'补助金',
  'k-fee':'Gas费','free':'免费','k-limit':'注册次数','k-limit-v':'每人一次 · 永久 · 不可变',
  'k-bio':'生物特征数据','never-stored':'从不存储','k-conf':'确认','k-conf-v':'6秒内',
  'live-stats':'链上实时统计',
  's-height':'区块高度','s-height-sub':'每6秒新区块 · BlockDAG · 两个节点并行',
  's-humans':'已验证人类','s-humans-sub':'生物特征ZKP · 一人一钱包，永久',
  's-supply':'总供应量','s-supply-sub':'始终 = 人类 x 1,000 AEQ',
  's-index':'Aequitas指数','s-index-sub':'0 = 完全平等 · 100 = 最大不平等',
  's-uptime':'运行时间','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'人类证明','ib-poh-t':'每个AEQ持有者必须证明自己是唯一的活人。没有机器人、公司或AI。只有真实的人类。',
  'ib-fair':'根本公平的分配','ib-fair-t':'每个经过验证的人类获得恰好1,000 AEQ。总供应量始终等于已验证人类 x 1,000。',
  'ib-dag':'BlockDAG架构','ib-dag-t':'多个区块可以同时产生并合并。更高吞吐量，更低延迟。',
  'ib-gas':'真正无Gas费','ib-gas-t':'注册绝对免费。不需要ETH、BNB或MATIC。不需要银行账户。',
  'recent-blocks':'最近区块','blocks-desc':'MERGE = 多个父区块（BlockDAG）。TX = 注册交易。区块时间：约6秒。',
  'loading':'加载区块中...','net-info':'网络信息','k-chain':'网络名称','k-symbol':'符号','k-btime':'出块时间',
  'k-cons':'共识','k-nodes':'活跃节点','k-storage':'存储','add-mm':'🦊 添加到METAMASK','k-dec':'小数位',
  'btn-add-mm':'+ 添加AEQUITAS网络',
  'phil':'"货币存在是因为人类存在。<br>仅此而已，不多也不少。"','phil-sub':'— AEQUITAS原则 —',
  'humans-title':'Aequitas Chain上的已验证人类',
  'h-what':'已验证人类是什么？','h-what-t':'已验证人类是一个加密证明属于独特活人的钱包地址。生物特征数据从不传输或存储。',
  'h-zkp':'零知识证明系统','h-zkp-t':'Aequitas使用BN128上的Groth16。证明大小：~200字节。验证：~10ms。',
  'h-sybil':'女巫攻击防护','h-sybil-t':'每个生物特征哈希永久存储。尝试用同一指纹注册两次立即被拒绝。',
  'h-global':'全球包容','h-global-t':'不需要银行账户、信用卡或加密货币。只需Android智能手机。',
  'reg-humans':'已注册人类','h-desc':'每个地址通过生物特征ZKP验证为唯一人类。每人获得1,000 AEQ。永久，不可变。',
  'no-humans':'还没有人类注册。\n\n下载Aequitas Android应用成为第一个！',
  'reg-stats':'注册统计','total-humans':'总人类数',
  'idx-title':'Aequitas指数 — 实时经济平等分数',
  'idx-desc':'从所有已验证人类的链上余额分布计算。0 = 完全平等。100 = 最大不平等。',
  'curr-idx':'当前指数','bar-0':'0 — 完全平等','bar-100':'100 — 最大不平等',
  'gini':'基尼系数','phase':'协议阶段',
  'pools-title':'再分配池','pools-desc':'当不平等阈值被超过时，AEQ自动重定向。',
  'vel-pool':'速度池','liq-pool':'流动性池','ubi-pool':'UBI池','treasury':'国库',
  'phases-title':'协议阶段','phases-desc':'过渡自动发生 — 不需要投票。',
  'p0':'引导期 · &lt;100人 · Cap: 50x','p1':'增长期 · 100-10,000 · Cap: 20x',
  'p2':'稳定期 · 10k-1M · Cap: 10x','p3':'成熟期 · 1M+ · Cap: 3x',
  'story-title':'Aequitas的故事','story-text':'<p>2009年。中本聪发布比特币。有史以来第一次无需银行即可传递价值。真正的革命。但几乎立即就出现了问题。</p><p>早期矿工以几乎为零的成本积累了数百万枚比特币。到2021年，前1%控制了超过90%的所有比特币。</p><p><span style="color:var(--gold)">Aequitas</span>被创建来回答：<em style="color:var(--gold)">"如果一种加密货币从第一原则出发设计，对每个人都公平，它会是什么样子？"</em></p><p>答案：<strong>货币存在是因为人类存在。因此，每个人仅凭其是人类这一事实，就应该拥有等额的货币。</strong></p><p><em style="color:var(--gold)">"货币存在是因为人类存在。仅此而已，不多也不少。"</em></p>',
  'nodes-title':'活跃节点','nodes-desc':'Aequitas网络在地理分布的云环境中运行两个节点。',
  'node1':'节点1 — Railway（主要）','node1-desc':'主要API服务器 · 区块生产者 · P2P引导 · PostgreSQL · MetaMask的RPC',
  'node2':'节点2 — Render（次要）','node2-desc':'次要API服务器 · 区块生产者 · P2P对等节点 · HTTP同步 · 共享PostgreSQL',
  'bootstrap-title':'引导节点','bootstrap-desc':'要运行您自己的Aequitas节点，请使用libp2p多地址连接到引导节点。',
  'tech-title':'技术规格','mm-config':'MetaMask配置',
  'k-lang':'语言','k-src':'源代码','evm-yes':'是 — JSON-RPC /rpc · MetaMask',
  'proto-label':'Aequitas V6协议 — 技术文档',
  'ca-title':'合约地址','ca-text':'链：Aequitas Chain（Chain ID: 1926）<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV6:  0xD487544fE06DeD5025DF7bD45bdFba5e9ffadd3f',
  'poa-title':'1. 生命证明','poa-text':'<p>在Aequitas中，如果有人消失，他们的AEQ最终会返回社区。</p>',
  'poa-box':'第0-2年：正常使用<br>第2年：警告1<br>第2年+180天：AEQ进入个人托管<br>第4年：如果仍不活跃 — UBI池',
  'guard-title':'2. 监护人系统','guard-text':'<p>受信任的监护人可以确认某人仍然活着 — 没有任何交易权限。</p>',
  'guard-box':'每人1个监护人 · 仅限confirmAlive() · 无交易权限<br>最多3个被监护人 · 7天时间锁 · 无循环关系',
  'dem-title':'3. 滞留费','dem-box':'超出fairShare部分的余额每年1%费用进入UBI池',
  'dem-text':'<p>历史先例：奥地利沃尔格尔（1932年）— 一年内将失业率降低了25%。</p>',
  'cap-title':'4. 财富上限','cap-box':'第0阶段：50x · 第1阶段：20x · 第2阶段：10x · 第3阶段：5x · 第4阶段：3x fairShare',
  'ubi-title':'5. 全民基本收入','ubi-box':'来源：交易费用 · 财富上限溢出 · 滞留费 · 非活跃托管',
  'inf-title':'6. 无算法通胀','inf-box':'创造新AEQ的唯一事件：新验证人类注册'
},
id:{
  'logo-sub':'BUKTI KEMANUSIAAN','live':'SIARAN LANGSUNG',
  'tab-register':'🔐 Daftar','tab-explorer':'🔍 Penjelajah','tab-humans':'👥 Manusia','tab-index':'📊 Indeks','tab-network':'🌐 Jaringan','tab-protocol':'📜 Protokol V7',
  'reg-title':'🔐 Daftar sebagai Manusia Terverifikasi',
  'reg-sub':'Bergabunglah dengan jaringan Aequitas dan terima 1.000 AEQ. Sekali, permanen, tanpa gas. Tidak ada data pribadi.',
  'app-title':'PENDAFTARAN HANYA MELALUI APLIKASI ANDROID',
  'app-text':'Bukti Kemanusiaan memerlukan verifikasi biometrik di perangkat Anda. Sidik jari Anda diproses oleh Hardware Secure Element — data mentah tidak pernah meninggalkan ponsel Anda. Unduh aplikasinya, pindai sidik jari, hubungkan wallet, dan <strong style="color:var(--gold)">1.000 AEQ Anda akan diberikan otomatis</strong>.',
  's1t':'Pemindaian Biometrik','s1d':'Buka aplikasi · pindai sidik jari · HSE memproses · data tidak pernah meninggalkan perangkat',
  's2t':'Pembuatan ZKP','s2d':'Bukti Groth16 dihasilkan · keunikan diverifikasi · hash tidak pernah terungkap',
  's3t':'Hubungkan Wallet','s3d':'Aplikasi membuka MetaMask · hubungkan wallet · alamat menerima 1.000 AEQ',
  's4t':'1.000 AEQ','s4d':'Terdaftar di V6 · dikonfirmasi di blok berikutnya · aplikasi memberi tahu otomatis',
  'priv-bar':'🔒 Hardware Secure Element · ZKP Groth16 · Data tidak pernah meninggalkan perangkat · Tanpa biaya gas',
  'conn-wallet':'DOMPET TERHUBUNG','proof-recv':'⚡ BUKTI ZK DITERIMA','proof-hint':'Hubungkan wallet untuk mendaftar',
  'btn-conn':'🦊 HUBUNGKAN METAMASK','btn-reg':'🔐 DAFTAR ON-CHAIN',
  'reg-log-hint':'// Buka Aplikasi Android Aequitas untuk menghasilkan bukti Anda, lalu kembali ke sini...',
  'reg-details':'Detail Pendaftaran','k-network':'Jaringan','k-chainid':'ID Rantai','k-grant':'Hibah',
  'k-fee':'Biaya Gas','free':'GRATIS','k-limit':'Pendaftaran','k-limit-v':'Sekali per manusia · permanen · tidak dapat diubah',
  'k-bio':'Data Biometrik','never-stored':'Tidak pernah disimpan','k-conf':'Konfirmasi','k-conf-v':'Dalam 6 detik',
  'live-stats':'Statistik Rantai Langsung',
  's-height':'Tinggi Blok','s-height-sub':'Blok baru setiap 6 detik · BlockDAG · Dua node paralel',
  's-humans':'Manusia Terverifikasi','s-humans-sub':'ZKP biometrik · Satu orang, satu wallet, selamanya',
  's-supply':'Total Pasokan','s-supply-sub':'Selalu = Manusia x 1.000 AEQ',
  's-index':'Indeks Aequitas','s-index-sub':'0 = kesetaraan sempurna · 100 = ketidaksetaraan maksimum',
  's-uptime':'Waktu Aktif','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Bukti Kemanusiaan','ib-poh-t':'Setiap pemegang AEQ harus membuktikan bahwa mereka adalah manusia unik yang hidup. Tidak ada bot, korporasi, atau AI. Hanya manusia nyata.',
  'ib-fair':'Distribusi yang Benar-Benar Adil','ib-fair-t':'Setiap manusia terverifikasi menerima tepat 1.000 AEQ. Total pasokan selalu sama dengan manusia x 1.000.',
  'ib-dag':'Arsitektur BlockDAG','ib-dag-t':'Beberapa blok dapat diproduksi secara bersamaan dan digabungkan. Throughput lebih tinggi, latensi lebih rendah.',
  'ib-gas':'Benar-Benar Tanpa Gas','ib-gas-t':'Pendaftaran tidak memerlukan biaya sama sekali. Tidak perlu ETH, BNB, atau MATIC. Tidak perlu rekening bank.',
  'recent-blocks':'Blok Terbaru','blocks-desc':'MERGE = beberapa induk (BlockDAG). TX = transaksi registrasi. Waktu blok: ~6 detik.',
  'loading':'Memuat blok...','net-info':'Informasi Jaringan','k-chain':'Nama Jaringan','k-symbol':'Simbol','k-btime':'Waktu Blok',
  'k-cons':'Konsensus','k-nodes':'Node Aktif','k-storage':'Penyimpanan','add-mm':'🦊 TAMBAHKAN KE METAMASK','k-dec':'Desimal',
  'btn-add-mm':'+ TAMBAHKAN JARINGAN AEQUITAS',
  'phil':'"Uang ada karena manusia ada.<br>Tidak lebih, tidak kurang."','phil-sub':'— PRINSIP AEQUITAS —',
  'humans-title':'Manusia Terverifikasi di Aequitas Chain',
  'h-what':'Apa itu Manusia Terverifikasi?','h-what-t':'Manusia Terverifikasi adalah alamat wallet yang terbukti secara kriptografis milik manusia unik yang hidup.',
  'h-zkp':'Sistem Bukti Zero-Knowledge','h-zkp-t':'Aequitas menggunakan Groth16 atas BN128. Ukuran bukti: ~200 byte. Verifikasi: ~10ms.',
  'h-sybil':'Pencegahan Serangan Sybil','h-sybil-t':'Setiap hash biometrik disimpan secara permanen. Mencoba mendaftar dua kali langsung ditolak.',
  'h-global':'Inklusi Global','h-global-t':'Tidak perlu rekening bank, kartu kredit, atau cryptocurrency. Hanya smartphone Android.',
  'reg-humans':'Manusia Terdaftar','h-desc':'Setiap alamat diverifikasi sebagai manusia unik melalui ZKP biometrik. Masing-masing menerima 1.000 AEQ. Permanen.',
  'no-humans':'Belum ada manusia terdaftar.\n\nUnduh Aplikasi Android Aequitas dan jadilah yang pertama!',
  'reg-stats':'Statistik Registri','total-humans':'Total Manusia',
  'idx-title':'Indeks Aequitas — Skor Kesetaraan Ekonomi Real-Time',
  'idx-desc':'Dihitung dari distribusi saldo on-chain semua manusia terverifikasi. 0 = kesetaraan sempurna. 100 = ketidaksetaraan maksimum.',
  'curr-idx':'Indeks Saat Ini','bar-0':'0 — Kesetaraan Sempurna','bar-100':'100 — Ketidaksetaraan Maks.',
  'gini':'Koefisien Gini','phase':'Fase Protokol',
  'pools-title':'Pool Redistribusi','pools-desc':'Ketika ambang ketidaksetaraan terlampaui, AEQ secara otomatis diarahkan.',
  'vel-pool':'Pool Kecepatan','liq-pool':'Pool Likuiditas','ubi-pool':'Pool UBI','treasury':'Perbendaharaan',
  'phases-title':'Fase Protokol','phases-desc':'Transisi terjadi secara otomatis — tidak diperlukan pemungutan suara.',
  'p0':'Bootstrap · &lt;100 manusia · Cap: 50x','p1':'Pertumbuhan · 100-10.000 · Cap: 20x',
  'p2':'Stabilitas · 10k-1M · Cap: 10x','p3':'Kedewasaan · 1M+ · Cap: 3x',
  'story-title':'Kisah Aequitas','story-text':'<p>Tahun 2009. Satoshi Nakamoto merilis Bitcoin. Untuk pertama kalinya nilai dapat ditransfer tanpa bank. Sebuah revolusi sejati. Tetapi sesuatu segera berjalan salah.</p><p>Penambang awal mengumpulkan jutaan koin dengan biaya hampir nol. Pada 2021, 1% teratas mengendalikan lebih dari 90% semua Bitcoin.</p><p><span style="color:var(--gold)">Aequitas</span> diciptakan untuk menjawab: <em style="color:var(--gold)">"Seperti apa cryptocurrency jika dirancang untuk adil bagi setiap manusia?"</em></p><p>Jawabannya: <strong>Uang ada karena manusia ada. Oleh karena itu, setiap orang harus memiliki bagian yang sama dari uang.</strong></p><p><em style="color:var(--gold)">"Uang ada karena manusia ada. Tidak lebih, tidak kurang."</em></p>',
  'nodes-title':'Node Aktif','nodes-desc':'Jaringan Aequitas beroperasi pada dua node di lingkungan cloud yang didistribusikan secara geografis.',
  'node1':'Node 1 — Railway (Utama)','node1-desc':'Server API utama · Produsen blok · P2P bootstrap · PostgreSQL · RPC untuk MetaMask',
  'node2':'Node 2 — Render (Sekunder)','node2-desc':'Server API sekunder · Produsen blok · P2P peer · Sinkronisasi HTTP · PostgreSQL bersama',
  'bootstrap-title':'Node Bootstrap','bootstrap-desc':'Untuk menjalankan node Aequitas Anda sendiri, hubungkan ke node bootstrap menggunakan alamat libp2p.',
  'tech-title':'Spesifikasi Teknis','mm-config':'Konfigurasi MetaMask',
  'k-lang':'Bahasa','k-src':'Kode Sumber','evm-yes':'Ya — JSON-RPC /rpc · MetaMask',
  'proto-label':'Protokol Aequitas V6 — Dokumentasi Teknis',
  'ca-title':'Alamat Kontrak','ca-text':'Rantai: Aequitas Chain (Chain ID: 1926)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV6:  0xD487544fE06DeD5025DF7bD45bdFba5e9ffadd3f',
  'poa-title':'1. BUKTI HIDUP','poa-text':'<p>Di Aequitas, jika seseorang menghilang, AEQ mereka akhirnya kembali ke komunitas.</p>',
  'poa-box':'Tahun 0-2: Penggunaan normal<br>Tahun 2: Peringatan 1<br>Tahun 2+180h: AEQ ke escrow pribadi<br>Tahun 4: Jika masih tidak aktif — Pool UBI',
  'guard-title':'2. SISTEM GUARDIAN','guard-text':'<p>Guardian tepercaya dapat mengkonfirmasi seseorang masih hidup — tanpa hak transaksi.</p>',
  'guard-box':'1 Guardian per manusia · Hanya confirmAlive() · Tanpa hak transaksi<br>Maks 3 wali · Kunci 7 hari · Tanpa hubungan melingkar',
  'dem-title':'3. DEMURRAGE','dem-box':'1% biaya tahunan atas saldo DI ATAS fairShare ke Pool UBI',
  'dem-text':'<p>Preseden: Worgl, Austria (1932) — mengurangi pengangguran 25% dalam satu tahun.</p>',
  'cap-title':'4. BATAS KEKAYAAN','cap-box':'Fase 0: 50x · Fase 1: 20x · Fase 2: 10x · Fase 3: 5x · Fase 4: 3x fairShare',
  'ubi-title':'5. PENDAPATAN DASAR UNIVERSAL','ubi-box':'Sumber: Biaya transaksi · Kelebihan batas · Demurrage · Escrow tidak aktif',
  'inf-title':'6. TANPA INFLASI ALGORITMIK','inf-box':'Satu-satunya peristiwa yang menciptakan AEQ baru: manusia terverifikasi baru mendaftar'
}
};

function showTab(name, el) {
  document.querySelectorAll('.tab-content').forEach(t => t.classList.remove('active'));
  document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
  document.getElementById('tab-' + name).classList.add('active');
  el.classList.add('active');
}

function setLang(lang) {
  curLang = lang;
  document.getElementById('lang-sel').value = lang;
  const t = T[lang];
  if (!t) return;
  document.querySelectorAll('[data-i18n]').forEach(el => {
    const key = el.getAttribute('data-i18n');
    if (t[key] !== undefined) el.innerHTML = t[key];
  });
}

function fmt(n) {
  if (n === undefined || n === null) return '—';
  if (typeof n === 'number') return n.toLocaleString();
  return n;
}

function timeAgo(ts) {
  const d = Math.floor(Date.now() / 1000) - ts;
  if (d < 60) return d + 's ago';
  if (d < 3600) return Math.floor(d / 60) + 'm ago';
  return Math.floor(d / 3600) + 'h ago';
}

function short(h, s, e) {
  s = s || 8; e = e || 6;
  return h ? h.slice(0, s) + '...' + h.slice(-e) : '—';
}

function avatarColor(a) {
  const c = ['#4FC3F7', '#00E676', '#FFB300', '#CE93D8', '#EF5350', '#4DD0E1'];
  return c[parseInt((a || '0x00').slice(2, 4), 16) % c.length];
}

async function addToMetaMask() {
  if (!window.ethereum) { alert('MetaMask not found. Please install MetaMask.'); return; }
  try {
    await window.ethereum.request({
      method: 'wallet_addEthereumChain',
      params: [{
        chainId: CID,
        chainName: 'Aequitas Chain',
        nativeCurrency: { name: 'AEQ', symbol: 'AEQ', decimals: 18 },
        rpcUrls: ['https://aequitas-production-9fba.up.railway.app/rpc'],
        blockExplorerUrls: ['https://aequitas-production-9fba.up.railway.app']
      }]
    });
    // Add AEQ as watchable token
    await window.ethereum.request({
      method: 'wallet_watchAsset',
      params: {
        type: 'ERC20',
        options: {
          address: '0xD487544fE06DeD5025DF7bD45bdFba5e9ffadd3f',
          symbol: 'AEQ',
          decimals: 18,
          name: 'Aequitas'
        }
      }
    });
  } catch (e) { console.error('MetaMask error:', e); }
}

async function loadStatus() {
  try {
    const d = await (await fetch('/api/status')).json();
    document.getElementById('s-height').textContent = fmt(d.height);
    document.getElementById('s-humans').textContent = fmt(d.total_humans);
    document.getElementById('s-supply').textContent = d.total_supply || '—';
    document.getElementById('s-index').textContent = fmt(d.index);
    const up = d.uptime || 0;
    document.getElementById('s-uptime').textContent = Math.floor(up/3600) + 'h ' + Math.floor((up%3600)/60) + 'm';
    document.getElementById('idx-score').textContent = fmt(d.index);
    document.getElementById('idx-gini').textContent = typeof d.gini === 'number' ? d.gini.toFixed(3) : '—';
    document.getElementById('idx-supply2').textContent = d.total_supply || '—';
    document.getElementById('idx-phase').textContent = fmt(d.phase);
    document.getElementById('idx-humans2').textContent = fmt(d.total_humans);
    document.getElementById('stat-humans').textContent = fmt(d.total_humans);
    document.getElementById('stat-supply').textContent = d.total_supply || '—';
    if (d.index !== undefined) {
      document.getElementById('idx-bar').style.width = Math.min(d.index, 100) + '%';
      const phases = ['Phase 0: Bootstrap — building the network', 'Phase 1: Growth — expanding human registry', 'Phase 2: Stability — redistribution active', 'Phase 3: Maturity — full decentralization'];
      document.getElementById('idx-phase-desc').textContent = phases[d.phase || 0] || 'Phase ' + (d.phase || 0);
    }
  } catch (e) {}
}

async function loadBlocks() {
  try {
    const blocks = await (await fetch('/api/blocks')).json();
    const list = document.getElementById('blocks-list');
    if (!blocks || !blocks.length) { list.innerHTML = '<div class="empty">No blocks yet</div>'; return; }
    document.getElementById('block-count').textContent = blocks.length + ' blocks';
    list.innerHTML = blocks.map(b => {
      const merge = b.parent_hashes && b.parent_hashes.length > 1;
      const hasTx = b.transactions && b.transactions.length > 0;
      return '<div class="block-item"><div class="block-num">#' + b.height + '</div><div><div class="block-hash">' + short(b.hash) + (merge ? ' <span class="bm">MERGE</span>' : '') + (hasTx ? ' <span class="bt">TX</span>' : '') + '</div><div class="block-parents">' + (b.parent_hashes ? b.parent_hashes.length + ' parent(s)' : '') + '</div></div><div class="block-right"><div class="block-humans">' + (b.humans || 0) + ' humans</div><div class="block-time">' + timeAgo(b.timestamp) + '</div></div></div>';
    }).join('');
  } catch (e) {}
}

async function loadHumans() {
  try {
    const d = await (await fetch('/api/humans')).json();
    document.getElementById('h-count').textContent = fmt(d.total);
    const list = document.getElementById('humans-list');
    if (!d.humans || !d.humans.length) { list.innerHTML = '<div class="empty">No humans registered yet.<br><br>Download the Aequitas Android App and be the first!</div>'; return; }
    list.innerHTML = d.humans.map(h => {
      const color = avatarColor(h.address || '0x00');
      const init = (h.address || '??').slice(2, 4).toUpperCase();
      return '<div class="hi"><div class="hav" style="background:' + color + '20;color:' + color + ';border-color:' + color + '50">' + init + '</div><div style="flex:1;min-width:0"><div class="hbal">' + fmt(h.balance) + ' AEQ</div><div class="hadr">' + h.address + '</div></div><div class="hbdg">HUMAN</div></div>';
    }).join('');
  } catch (e) {}
}

function checkProofParams() {
  const p = new URLSearchParams(window.location.search);
  const proofId = p.get('proofId');
  const proof = p.get('proof');
  if (proofId) {
    fetch(PS + '/get/' + proofId).then(r => r.json()).then(pd => {
      proofData = pd;
      document.getElementById('pbox').style.display = 'block';
      document.getElementById('pval').textContent = 'Proof ID: ' + proofId + ' — Connect wallet to register';
      document.querySelectorAll('.tab')[0].click();
      setTimeout(() => connectWallet(), 600);
    }).catch(e => console.error(e));
  } else if (proof) {
    try {
      proofData = JSON.parse(decodeURIComponent(proof));
      document.getElementById('pbox').style.display = 'block';
      document.getElementById('pval').textContent = 'Proof received — Connect wallet to register';
      document.querySelectorAll('.tab')[0].click();
      setTimeout(() => connectWallet(), 600);
    } catch (e) {}
  }
}

async function connectWallet() {
  if (!window.ethereum) {
    addLog('MetaMask not found. Please install MetaMask.', 'err');
    return;
  }
  try {
    await addToMetaMask();
    const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
    waddr = accounts[0];
    document.getElementById('wbox').style.display = 'block';
    document.getElementById('wadr').textContent = waddr;
    const btn = document.getElementById('btn-conn');
    btn.textContent = waddr.slice(0, 10) + '...' + waddr.slice(-4);
    btn.style.background = 'var(--green)';
    btn.style.color = '#050A14';
    try {
      const br = await fetch('/api/balance?wallet=' + waddr);
      const bd = await br.json();
      if (bd.is_human) {
        addLog('Already registered! Balance: ' + bd.balance + ' AEQ', 'ok');
        document.getElementById('btn-reg').disabled = true;
        document.getElementById('btn-reg').textContent = 'ALREADY REGISTERED';
      } else if (proofData) {
        document.getElementById('btn-reg').disabled = false;
        document.getElementById('btn-reg').textContent = 'PROOF READY — CLICK TO REGISTER';
      } else {
        document.getElementById('btn-reg').disabled = true;
      }
    } catch (e) {
      document.getElementById('btn-reg').disabled = !proofData;
    }
  } catch (e) {
    addLog('Connection failed: ' + e.message, 'err');
  }
}

function addLog(msg, type) {
  const el = document.getElementById('rlog');
  el.innerHTML += '<div><span class="' + type + '">' + msg + '</span></div>';
}

async function doRegister() {
  if (!waddr || !proofData) return;
  try {
    addLog('Registering on Aequitas V6...', 'info');
    document.getElementById('btn-reg').disabled = true;
    const r = await fetch('/api/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ wallet: waddr, pA: proofData.pA, pB: proofData.pB, pC: proofData.pC, pubSignals: proofData.pubSignals })
    });
    const d = await r.json();
    if (!d.success) { addLog('Error: ' + d.message, 'err'); document.getElementById('btn-reg').disabled = false; return; }
    addLog('Registered! ' + d.message, 'ok');
    setTimeout(() => { window.location.href = '/registered?wallet=' + waddr; }, 1500);
  } catch (e) { addLog('Error: ' + e.message, 'err'); document.getElementById('btn-reg').disabled = false; }
}

window.ethereum && window.ethereum.on('accountsChanged', function(a) {
  waddr = a[0] || '';
  if (waddr) {
    document.getElementById('wbox').style.display = 'block';
    document.getElementById('wadr').textContent = waddr;
    const btn = document.getElementById('btn-conn');
    btn.textContent = waddr.slice(0, 10) + '...' + waddr.slice(-4);
    btn.style.background = 'var(--green)';
    btn.style.color = '#050A14';
    fetch('/api/balance?wallet=' + waddr).then(function(r) { return r.json(); }).then(function(bd) {
      if (bd.is_human) {
        document.getElementById('btn-reg').disabled = true;
        document.getElementById('btn-reg').textContent = 'ALREADY REGISTERED';
        addLog('Already registered! Balance: ' + bd.balance + ' AEQ', 'ok');
      } else {
        document.getElementById('btn-reg').disabled = !proofData;
        if (proofData) document.getElementById('btn-reg').textContent = 'PROOF READY — CLICK TO REGISTER';
      }
    }).catch(function() { document.getElementById('btn-reg').disabled = !proofData; });
  }
});

checkProofParams();
loadStatus();
loadBlocks();
loadHumans();
setInterval(loadStatus, 6000);
setInterval(loadBlocks, 6000);
setInterval(loadHumans, 10000);
</script>
</body>
</html>`
