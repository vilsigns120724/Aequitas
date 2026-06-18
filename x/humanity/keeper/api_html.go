package keeper

const explorerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0">
<title>Aequitas — Proof of Humanity Chain</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;900&family=DM+Serif+Display&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
<style>
*{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#06091A;--card:#0C1228;--card2:#111830;--border:#1C2B48;
  --green:#10B981;--blue:#60A5FA;--gold:#F5A623;--purple:#A78BFA;--red:#F87171;
  --text:#F0F4FF;--muted:#6B7EA8;--teal:#2DD4BF;
  --shadow:0 2px 12px rgba(0,0,0,0.4);--shadow-lg:0 8px 32px rgba(0,0,0,0.5);
  --radius:12px;--radius-sm:8px;
}
body{background:var(--bg);color:var(--text);font-family:'Inter',sans-serif;min-height:100vh;overflow-x:hidden;line-height:1.5}
header{background:rgba(6,9,26,0.95);backdrop-filter:blur(12px);border-bottom:1px solid var(--border);padding:0 24px;position:sticky;top:0;z-index:100;display:flex;align-items:center;justify-content:space-between;height:60px;gap:10px}
.logo-wrap{display:flex;align-items:center;gap:12px;flex-shrink:0}
.logo-icon{width:32px;height:32px;background:linear-gradient(135deg,var(--gold),#E67E00);border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:16px;box-shadow:0 2px 8px rgba(245,166,35,0.3)}
.logo-text{font-size:1rem;font-weight:900;color:var(--gold);letter-spacing:3px}
.logo-sub{font-size:0.5rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase}
.header-right{display:flex;gap:8px;align-items:center}
.badge{display:flex;align-items:center;gap:5px;padding:5px 10px;border-radius:20px;font-size:0.6rem;letter-spacing:0.5px;font-weight:600}
.badge-live{background:rgba(16,185,129,0.1);border:1px solid rgba(16,185,129,0.25);color:var(--green)}
.badge-dag{background:rgba(96,165,250,0.1);border:1px solid rgba(96,165,250,0.25);color:var(--blue)}
.pulse{width:5px;height:5px;border-radius:50%;background:var(--green);animation:pulse 2s infinite}
@keyframes pulse{0%,100%{opacity:1;transform:scale(1)}50%{opacity:0.4;transform:scale(0.8)}}
.lang-sel{background:var(--card);color:var(--muted);border:1px solid var(--border);border-radius:6px;padding:5px 10px;font-family:'Inter',sans-serif;font-size:0.62rem;outline:none;cursor:pointer}
.tabs{background:rgba(6,9,26,0.8);border-bottom:1px solid var(--border);padding:0 24px;display:flex;overflow-x:auto;-webkit-overflow-scrolling:touch;scrollbar-width:none;gap:4px}
.tabs::-webkit-scrollbar{display:none}
.tab{padding:16px 16px;font-size:0.65rem;color:var(--muted);cursor:pointer;border-bottom:2px solid transparent;letter-spacing:0.5px;font-weight:600;white-space:nowrap;transition:all 0.2s;flex-shrink:0}
.tab:hover{color:var(--text)}.tab.active{color:var(--blue);border-bottom-color:var(--blue)}
.tab-content{display:none}.tab-content.active{display:block}
.hero{padding:20px 20px 0}
.section-label{font-size:0.6rem;color:var(--muted);letter-spacing:3px;text-transform:uppercase;margin-bottom:14px;font-weight:600}
.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:1px;background:var(--border);border:1px solid var(--border);border-radius:var(--radius);overflow:hidden;margin-bottom:20px}
.stat{background:var(--card);padding:20px 16px;position:relative;transition:background 0.2s}
.stat:hover{background:var(--card2)}
.stat-accent{position:absolute;top:0;left:0;right:0;height:3px}
.stat-icon{font-size:1rem;margin-bottom:8px}
.stat-lbl{font-size:0.6rem;color:var(--muted);letter-spacing:1px;text-transform:uppercase;margin-bottom:6px;font-weight:500}
.stat-val{font-size:1.7rem;font-weight:900;line-height:1;margin-bottom:4px;font-family:'DM Serif Display',serif}
.stat-sub{font-size:0.58rem;color:var(--muted);line-height:1.5}
.c-green .stat-val{color:var(--green)}.c-green .stat-accent{background:linear-gradient(90deg,var(--green),transparent)}
.c-blue .stat-val{color:var(--blue)}.c-blue .stat-accent{background:linear-gradient(90deg,var(--blue),transparent)}
.c-gold .stat-val{color:var(--gold)}.c-gold .stat-accent{background:linear-gradient(90deg,var(--gold),transparent)}
.c-purple .stat-val{color:var(--purple)}.c-purple .stat-accent{background:linear-gradient(90deg,var(--purple),transparent)}
.c-teal .stat-val{color:var(--teal)}.c-teal .stat-accent{background:linear-gradient(90deg,var(--teal),transparent)}
.info-banner{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:20px;margin-bottom:20px;display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:20px;box-shadow:var(--shadow)}
.ib-icon{font-size:1.4rem;margin-bottom:8px}
.ib-title{font-size:0.7rem;color:var(--gold);font-weight:700;margin-bottom:8px;letter-spacing:0.5px}
.ib-text{font-size:0.65rem;color:var(--muted);line-height:1.8}
.main-grid{display:grid;grid-template-columns:1fr 310px;gap:16px;padding:0 20px 20px}
@media(max-width:800px){.main-grid{grid-template-columns:1fr}.right-col{display:none}}
.section{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);overflow:hidden;box-shadow:var(--shadow)}
.sec-head{padding:14px 18px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;background:var(--card2)}
.sec-title{font-size:0.65rem;color:var(--muted);letter-spacing:1px;text-transform:uppercase;display:flex;align-items:center;gap:8px;font-weight:600}
.sec-dot{width:6px;height:6px;border-radius:50%;background:var(--green);box-shadow:0 0 6px var(--green)}
.sec-count{font-size:0.6rem;color:var(--muted);background:var(--card);padding:3px 8px;border-radius:10px;border:1px solid var(--border)}
.sec-desc{padding:10px 18px;font-size:0.65rem;color:var(--muted);background:var(--card2);border-bottom:1px solid var(--border);line-height:1.7}
.block-item{padding:12px 18px;border-bottom:1px solid rgba(28,43,72,0.5);display:grid;grid-template-columns:60px 1fr auto;gap:10px;align-items:center;transition:background 0.15s}
.block-item:hover{background:var(--card2)}.block-item:last-child{border-bottom:none}
.block-num{font-size:0.8rem;font-weight:700;color:var(--blue);font-family:'JetBrains Mono',monospace}
.block-hash{font-size:0.63rem;color:var(--muted);margin-bottom:2px;display:flex;align-items:center;gap:4px;flex-wrap:wrap;font-family:'JetBrains Mono',monospace}
.block-parents{font-size:0.57rem;color:#3A5570}
.block-right{text-align:right}
.block-humans{font-size:0.65rem;color:var(--gold);margin-bottom:2px;font-weight:600}
.block-time{font-size:0.57rem;color:var(--green)}
.bm{background:rgba(167,139,250,0.1);color:var(--purple);font-size:0.53rem;padding:2px 6px;border-radius:4px;border:1px solid rgba(167,139,250,0.2)}
.bt{background:rgba(16,185,129,0.1);color:var(--green);font-size:0.53rem;padding:2px 6px;border-radius:4px;border:1px solid rgba(16,185,129,0.2)}
.empty{padding:40px;text-align:center;color:var(--muted);font-size:0.7rem;line-height:2.5}
.right-col{display:flex;flex-direction:column;gap:12px}
.ic{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:18px;box-shadow:var(--shadow)}
.ic-title{font-size:0.6rem;color:var(--muted);letter-spacing:1.5px;text-transform:uppercase;margin-bottom:14px;font-weight:600}
.ic-row{display:flex;justify-content:space-between;align-items:center;padding:8px 0;border-bottom:1px solid rgba(28,43,72,0.5)}
.ic-row:last-child{border-bottom:none}
.ic-key{font-size:0.63rem;color:var(--muted)}
.ic-val{font-size:0.63rem;color:var(--text);text-align:right;max-width:58%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;font-weight:500}
.ic-val.g{color:var(--green)}.ic-val.b{color:var(--blue)}.ic-val.go{color:var(--gold)}.ic-val.p{color:var(--purple)}
.mm-card{background:rgba(96,165,250,0.05);border:1px solid rgba(96,165,250,0.15);border-radius:var(--radius);padding:16px}
.mm-title{font-size:0.6rem;color:var(--blue);letter-spacing:1.5px;margin-bottom:12px;font-weight:700;text-transform:uppercase}
.mm-row{display:flex;justify-content:space-between;padding:6px 0;border-bottom:1px solid rgba(28,43,72,0.5)}
.mm-row:last-child{border-bottom:none}
.mm-key{font-size:0.6rem;color:var(--muted)}.mm-val{font-size:0.6rem;color:var(--purple);font-family:'JetBrains Mono',monospace}
.mm-btn{width:100%;margin-top:12px;padding:11px;background:var(--blue);color:#06091A;border:none;border-radius:var(--radius-sm);cursor:pointer;font-family:'Inter',sans-serif;font-size:0.68rem;font-weight:700;letter-spacing:0.5px;transition:all 0.2s}
.mm-btn:hover{opacity:0.87;transform:translateY(-1px)}
.phil-card{background:linear-gradient(135deg,rgba(245,166,35,0.08),rgba(12,18,40,0.9));border:1px solid rgba(245,166,35,0.2);border-radius:var(--radius);padding:22px;text-align:center}
.phil-quote{font-size:0.85rem;color:var(--gold);font-style:italic;line-height:2;margin-bottom:6px;font-family:'DM Serif Display',serif}
.phil-sub{font-size:0.58rem;color:var(--muted);letter-spacing:1.5px;text-transform:uppercase}
.hs{padding:20px;display:grid;grid-template-columns:1fr 290px;gap:16px}
@media(max-width:800px){.hs{grid-template-columns:1fr}}
.hi{padding:12px 18px;border-bottom:1px solid rgba(28,43,72,0.5);display:flex;align-items:center;gap:12px;transition:background 0.15s}
.hi:hover{background:var(--card2)}.hi:last-child{border-bottom:none}
.hav{width:36px;height:36px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:0.65rem;font-weight:bold;flex-shrink:0;border:2px solid}
.hbal{font-size:0.82rem;color:var(--gold);font-weight:700;margin-bottom:1px;font-family:'DM Serif Display',serif}
.hadr{font-size:0.6rem;color:var(--muted);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;font-family:'JetBrains Mono',monospace}
.hbdg{font-size:0.56rem;padding:3px 8px;border-radius:10px;flex-shrink:0;background:rgba(16,185,129,0.1);color:var(--green);border:1px solid rgba(16,185,129,0.2);font-weight:600}
.is{padding:20px;display:grid;grid-template-columns:1fr 1fr;gap:16px}
@media(max-width:700px){.is{grid-template-columns:1fr}}
.idx{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:24px;box-shadow:var(--shadow)}
.idx-title{font-size:0.6rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:10px;font-weight:600}
.idx-desc{font-size:0.67rem;color:var(--muted);line-height:1.8;margin-bottom:16px}
.idx-big{font-size:2.8rem;font-weight:900;color:var(--gold);line-height:1;font-family:'DM Serif Display',serif}
.idx-lbl{font-size:0.6rem;color:var(--muted);margin-top:4px}
.bar-bg{height:8px;background:var(--card2);border-radius:6px;overflow:hidden;margin:14px 0 6px}
.bar-fill{height:100%;border-radius:6px;background:linear-gradient(90deg,var(--green),var(--gold),var(--red));transition:width 1.5s ease}
.bar-lbl{display:flex;justify-content:space-between;font-size:0.55rem;color:var(--muted)}
.mrow{display:grid;grid-template-columns:repeat(2,1fr);gap:8px;margin-top:14px}
.mbox{background:var(--card2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:12px;text-align:center}
.mval{font-size:1.15rem;font-weight:700;color:var(--gold);font-family:'DM Serif Display',serif}
.mlbl{font-size:0.57rem;color:var(--muted);margin-top:3px;font-weight:500}
.story{font-size:0.7rem;line-height:2;color:var(--muted)}
.story p{margin-bottom:14px}
.hlbox{background:rgba(245,166,35,0.05);border-left:3px solid var(--gold);border-radius:0 var(--radius-sm) var(--radius-sm) 0;padding:14px 18px;margin:16px 0;font-size:0.67rem;color:var(--text);line-height:1.9}
.ns{padding:20px;display:grid;grid-template-columns:1fr 1fr;gap:16px}
@media(max-width:700px){.ns{grid-template-columns:1fr}}
.nc{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:20px;box-shadow:var(--shadow)}
.nc-title{font-size:0.6rem;color:var(--muted);letter-spacing:1.5px;text-transform:uppercase;margin-bottom:14px;font-weight:600}
.nbox{background:var(--card2);border-radius:var(--radius-sm);padding:14px;border:1px solid var(--border);margin-bottom:10px}
.nstat{display:flex;align-items:center;gap:6px;font-size:0.67rem;color:var(--green);margin-bottom:5px;font-weight:600}
.ndot{width:7px;height:7px;border-radius:50%;background:var(--green);box-shadow:0 0 6px var(--green)}
.nurl{font-size:0.58rem;color:var(--muted);word-break:break-all;margin-bottom:3px;font-family:'JetBrains Mono',monospace}
.ndesc{font-size:0.58rem;color:#3A5570}
.spect{width:100%;border-collapse:collapse}
.spect td{padding:8px 0;border-bottom:1px solid rgba(28,43,72,0.5);font-size:0.63rem}
.spect tr:last-child td{border-bottom:none}
.spect td:first-child{color:var(--muted);width:45%}
.spect td:last-child{text-align:right;font-family:'JetBrains Mono',monospace;font-size:0.6rem}
.bsbox{background:var(--card2);border-radius:var(--radius-sm);padding:12px;font-size:0.58rem;color:var(--purple);word-break:break-all;line-height:1.7;border:1px solid var(--border);font-family:'JetBrains Mono',monospace}
.rs{padding:20px;max-width:600px;margin:0 auto}
.rhero{background:linear-gradient(135deg,rgba(96,165,250,0.08),rgba(12,18,40,0.9));border:1px solid rgba(96,165,250,0.2);border-radius:var(--radius);padding:24px;margin-bottom:16px;text-align:center;box-shadow:var(--shadow)}
.rhero-title{font-size:1.05rem;font-weight:700;color:var(--text);margin-bottom:8px;font-family:'DM Serif Display',serif}
.rhero-sub{font-size:0.67rem;color:var(--muted);line-height:1.9}
.aonly{background:rgba(167,139,250,0.05);border:1px solid rgba(167,139,250,0.15);border-radius:var(--radius);padding:20px;text-align:center;margin-bottom:16px}
.aonly-icon{font-size:2rem;margin-bottom:8px}
.aonly-title{font-size:0.7rem;color:var(--purple);font-weight:700;letter-spacing:1px;margin-bottom:10px}
.aonly-text{font-size:0.65rem;color:var(--muted);line-height:1.9}
.rsteps{display:grid;grid-template-columns:repeat(4,1fr);gap:8px;margin-bottom:16px}
@media(max-width:520px){.rsteps{grid-template-columns:repeat(2,1fr)}}
.rstep{background:var(--card);border:1px solid var(--border);border-radius:var(--radius-sm);padding:16px;text-align:center;transition:border-color 0.2s}
.rstep:hover{border-color:var(--blue)}
.snum{width:28px;height:28px;background:var(--blue);border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 10px;font-weight:700;font-size:0.72rem;color:#06091A}
.stitle{font-size:0.63rem;color:var(--text);font-weight:600;margin-bottom:5px}
.sdesc{font-size:0.6rem;color:var(--muted);line-height:1.7}
.pbar{background:rgba(16,185,129,0.08);border:1px solid rgba(16,185,129,0.15);border-radius:var(--radius-sm);padding:10px 14px;margin-bottom:14px;font-size:0.63rem;color:var(--green);text-align:center;line-height:1.8}
.rcard{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:20px;margin-bottom:14px;box-shadow:var(--shadow)}
.wbox{background:rgba(16,185,129,0.06);border:1px solid rgba(16,185,129,0.15);border-radius:var(--radius-sm);padding:10px;margin-bottom:10px;display:none}
.wlbl{font-size:0.57rem;color:var(--muted);margin-bottom:2px;letter-spacing:1px;font-weight:500}
.wadr{font-size:0.72rem;color:var(--green);font-weight:700;font-family:'JetBrains Mono',monospace}
.pbox{background:var(--card2);border:1px solid rgba(245,166,35,0.15);border-radius:var(--radius-sm);padding:10px;margin-bottom:10px;display:none}
.plbl{font-size:0.57rem;color:var(--gold);margin-bottom:2px;letter-spacing:1px;font-weight:500}
.pval{font-size:0.63rem;color:var(--muted);font-family:'JetBrains Mono',monospace}
.rbtn{width:100%;padding:14px;border-radius:var(--radius-sm);border:none;cursor:pointer;font-family:'Inter',sans-serif;font-size:0.74rem;font-weight:700;letter-spacing:0.3px;transition:all 0.2s;margin-bottom:8px}
.rbtn:hover{opacity:0.87;transform:translateY(-1px)}
.bc{background:var(--blue);color:#06091A}.br{background:var(--gold);color:#06091A}
.rbtn:disabled{opacity:0.3;cursor:not-allowed;transform:none}
.rlog{background:var(--card2);border-radius:var(--radius-sm);padding:12px;font-size:0.63rem;line-height:2;min-height:52px;border:1px solid var(--border);font-family:'JetBrains Mono',monospace}
.rlog .ok{color:var(--green)}.rlog .err{color:var(--red)}.rlog .info{color:var(--gold)}
.ps{padding:20px;max-width:800px;margin:0 auto}
/* Swap tab */
.pct-row{display:flex;gap:6px;margin-bottom:8px}
.pctbtn{flex:1;padding:8px;font-size:12px;background:var(--card2);border:1px solid var(--border);color:var(--text);border-radius:var(--radius-sm);cursor:pointer;font-family:'Inter',sans-serif;font-weight:600;transition:all 0.2s}
.pctbtn:hover{border-color:var(--gold);color:var(--gold)}
/* Demurrage notice */
#demurrage-notice{font-size:13px;padding:12px 14px;border-radius:var(--radius-sm);background:rgba(245,166,35,0.08);border:1px solid rgba(245,166,35,0.2);color:var(--gold);margin:10px 0;line-height:1.7}
/* Swap specific */
.swap-dir{background:var(--card2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:8px;cursor:pointer;font-size:1rem;transition:all 0.2s;width:100%;margin:8px 0}
.swap-dir:hover{border-color:var(--blue)}
input[type=number]{background:var(--card2);border:1px solid var(--border);color:var(--text);border-radius:var(--radius-sm);padding:10px 12px;font-family:'Inter',sans-serif;font-size:0.8rem;outline:none;transition:border-color 0.2s}
input[type=number]:focus{border-color:var(--blue)}
input[type=number]::-webkit-inner-spin-button{opacity:0.5}
@media(max-width:480px){
  .stats-grid{grid-template-columns:repeat(2,1fr)}
  .stat-val{font-size:1.4rem}
  header{height:52px}
  .logo-text{font-size:0.85rem;letter-spacing:2px}
  .badge-dag{display:none}
  .main-grid{padding:0 12px 12px}
  .hero{padding:14px 12px 0}
}
</style>
<script src="https://cdnjs.cloudflare.com/ajax/libs/ethers/6.13.0/ethers.umd.min.js"></script>
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
  <div class="tab" onclick="showTab('swap',this)" data-i18n="tab-swap">🔄 Swap</div>
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
    <div><div class="ib-icon">🛡</div><div class="ib-title" data-i18n="h-sybil">Sybil Attack Prevention</div><div class="ib-text" data-i18n="h-sybil-t">Each biometric hash is stored permanently. Attempting to register twice with the same fingerprint is immediately rejected. One human, one wallet, forever. <strong style="color:var(--gold)">⚠ Test phase:</strong> Current biometric verification is device-bound. A hardware physiological sensor (MAX30102 PPG) is planned to provide device-independent identity verification in a future update.</div></div>
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

<!-- SWAP -->
<div id="tab-swap" class="tab-content">
<div class="rs">
  <div class="rhero">
    <div class="rhero-title" data-i18n="swap-title">🔄 Swap AEQ ↔ tUSD</div>
    <div class="rhero-sub" data-i18n="swap-sub">Exchange AEQ for tUSD (a simulated test-dollar) through the native liquidity pool. A 0.1% fee applies only to swaps — ordinary AEQ transfers between people remain completely free.</div>
  </div>
  <div class="pbar" data-i18n="swap-priv-bar">🔒 0.1% swap fee only · AEQ-to-AEQ transfers stay free · tUSD is a test currency with no real-world value</div>
  <div class="rcard">
    <div class="wbox" id="swap-wbox"><div class="wlbl" data-i18n="conn-wallet">CONNECTED WALLET</div><div class="wadr" id="swap-wadr">—</div></div>
    <div id="demurrage-notice" style="display:none;font-size:13px;padding:10px 12px;border-radius:8px;background:rgba(255,179,0,0.1);border:1px solid rgba(255,179,0,0.3);color:var(--gold);margin:10px 0"></div>
    <div class="ic-row" style="margin:8px 0"><span class="ic-key" data-i18n="swap-your-aeq">Your AEQ</span><span class="ic-val go" id="swap-bal-aeq">—</span></div>
    <div class="ic-row" style="margin-bottom:16px"><span class="ic-key" data-i18n="swap-your-tusd">Your tUSD</span><span class="ic-val go" id="swap-bal-tusd">—</span></div>

    <div class="swap-dir" id="swap-dir-box" style="display:flex;gap:8px;margin-bottom:12px">
      <button class="rbtn" id="swap-dir-a2t" onclick="setSwapDirection('aeq_to_tusd')" data-i18n="swap-aeq-to-tusd" style="flex:1">AEQ → tUSD</button>
      <button class="rbtn" id="swap-dir-t2a" onclick="setSwapDirection('tusd_to_aeq')" data-i18n="swap-tusd-to-aeq" style="flex:1">tUSD → AEQ</button>
    </div>
    <input type="number" id="swap-amount" placeholder="Amount" style="width:100%;padding:14px;border-radius:8px;border:1px solid var(--border);background:#0A1220;color:#E8EDF5;font-size:16px;margin-bottom:8px;box-sizing:border-box">
    <div class="ic-row" style="margin-bottom:8px"><span class="ic-key" data-i18n="swap-fee-est">Estimated 0.1% fee</span><span class="ic-val" id="swap-fee-est">—</span></div>
    <div id="swap-warn" style="display:none;font-size:13px;padding:10px 12px;border-radius:8px;background:rgba(255,179,0,0.1);border:1px solid rgba(255,179,0,0.3);color:var(--gold);margin-bottom:16px"></div>

    <button class="rbtn bc" id="swap-btn-conn" onclick="connectSwapWallet()" data-i18n="btn-conn">🦊 CONNECT METAMASK</button>
    <button class="rbtn br" id="swap-btn-go" onclick="doSwap()" disabled data-i18n="swap-btn-go">🔄 SWAP</button>
    <div class="rlog" id="swap-log"><span class="info" data-i18n="swap-log-hint">// Connect your wallet to swap...</span></div>

    <div class="ic" style="margin-top:20px">
      <div class="ic-title" data-i18n="swap-no-liquidity">No tUSD yet?</div>
      <div class="ic-row"><span class="ic-key" data-i18n="swap-faucet-desc">Registered humans can claim test-tUSD once</span></div>
      <button class="rbtn" id="swap-btn-faucet" onclick="claimFaucet()" disabled data-i18n="swap-btn-faucet" style="margin-top:8px">💧 CLAIM TEST-tUSD</button>
    </div>

    <div class="ic" style="margin-top:20px">
      <div class="ic-title" data-i18n="swap-addliq-title">Provide Liquidity</div>
      <div class="ic-row"><span class="ic-key" id="swap-addliq-desc" data-i18n="swap-addliq-desc">Be the first to deposit — your ratio sets the starting price.</span></div>
      <input type="number" id="addliq-aeq" placeholder="AEQ amount" oninput="updateLiquidityRatio('aeq')" style="width:100%;padding:12px;border-radius:8px;border:1px solid var(--border);background:#0A1220;color:#E8EDF5;font-size:15px;margin:8px 0 4px;box-sizing:border-box">
      <div class="pct-row" style="display:flex;gap:6px;margin-bottom:8px">
        <button class="rbtn pctbtn" onclick="setPctAmount('aeq',0.25)" style="flex:1;padding:8px;font-size:12px">25%</button>
        <button class="rbtn pctbtn" onclick="setPctAmount('aeq',0.5)" style="flex:1;padding:8px;font-size:12px">50%</button>
        <button class="rbtn pctbtn" onclick="setPctAmount('aeq',0.75)" style="flex:1;padding:8px;font-size:12px">75%</button>
        <button class="rbtn pctbtn" onclick="setPctAmount('aeq',1)" style="flex:1;padding:8px;font-size:12px">MAX</button>
      </div>
      <input type="number" id="addliq-tusd" placeholder="tUSD amount" oninput="updateLiquidityRatio('tusd')" style="width:100%;padding:12px;border-radius:8px;border:1px solid var(--border);background:#0A1220;color:#E8EDF5;font-size:15px;margin-bottom:4px;box-sizing:border-box">
      <div class="pct-row" style="display:flex;gap:6px;margin-bottom:8px">
        <button class="rbtn pctbtn" onclick="setPctAmount('tusd',0.25)" style="flex:1;padding:8px;font-size:12px">25%</button>
        <button class="rbtn pctbtn" onclick="setPctAmount('tusd',0.5)" style="flex:1;padding:8px;font-size:12px">50%</button>
        <button class="rbtn pctbtn" onclick="setPctAmount('tusd',0.75)" style="flex:1;padding:8px;font-size:12px">75%</button>
        <button class="rbtn pctbtn" onclick="setPctAmount('tusd',1)" style="flex:1;padding:8px;font-size:12px">MAX</button>
      </div>
      <button class="rbtn" id="swap-btn-addliq" onclick="doAddLiquidity()" disabled data-i18n="swap-btn-addliq" style="margin-top:4px">💧 ADD LIQUIDITY</button>
    </div>

    <div class="ic" id="lp-position-box" style="margin-top:20px;display:none">
      <div class="ic-title" data-i18n="swap-lp-title">Your LP Position</div>
      <div class="ic-row"><span class="ic-key" data-i18n="swap-lp-share">Pool Share</span><span class="ic-val go" id="lp-share-pct">—</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="swap-lp-withdrawable">Withdrawable</span><span class="ic-val" id="lp-withdrawable">—</span></div>
      <div style="display:flex;align-items:center;gap:8px;margin:10px 0 6px">
        <input type="number" id="remove-pct-input" min="0" max="100" step="0.1" placeholder="%" oninput="setRemovePctManual(this.value)" style="width:80px;padding:10px;border-radius:8px;border:1px solid var(--border);background:#0A1220;color:#E8EDF5;font-size:14px;box-sizing:border-box">
        <span style="color:var(--muted);font-size:13px" data-i18n="swap-lp-pct-label">% of your position</span>
      </div>
      <div class="pct-row" style="display:flex;gap:6px;margin-bottom:8px">
        <button class="rbtn pctbtn" onclick="setRemovePct(0.25,this)" style="flex:1;padding:8px;font-size:12px">25%</button>
        <button class="rbtn pctbtn" onclick="setRemovePct(0.5,this)" style="flex:1;padding:8px;font-size:12px">50%</button>
        <button class="rbtn pctbtn" onclick="setRemovePct(0.75,this)" style="flex:1;padding:8px;font-size:12px">75%</button>
        <button class="rbtn pctbtn" onclick="setRemovePct(1,this)" style="flex:1;padding:8px;font-size:12px">MAX</button>
      </div>
      <div class="ic-row" style="margin-bottom:8px"><span class="ic-key" data-i18n="swap-lp-youget">You will receive</span><span class="ic-val go" id="lp-remove-preview">—</span></div>
      <button class="rbtn br" id="swap-btn-removeliq" onclick="doRemoveLiquidity()" data-i18n="swap-btn-removeliq">🔥 REMOVE LIQUIDITY</button>
    </div>
  </div>

  <div class="ic">
    <div class="ic-title" data-i18n="swap-pool-title">Pool Status</div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-pool-aeq">Pool AEQ Reserve</span><span class="ic-val" id="pool-reserve-aeq">—</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-pool-tusd">Pool tUSD Reserve</span><span class="ic-val" id="pool-reserve-tusd">—</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-pool-price">Current Price</span><span class="ic-val go" id="pool-price">—</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-fee-bps">Swap Fee</span><span class="ic-val g">0.1%</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-fee-split">Fee Distribution</span><span class="ic-val" data-i18n="swap-fee-split-v">40% Validators / 30% LPs / 20% UBI / 10% Treasury</span></div>
  </div>
  <div class="ic">
    <div class="ic-title" data-i18n="swap-pools-addr-title">Tokenomics Pool Addresses</div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-validators">Validators (40%)</span><span class="ic-val p" style="font-size:11px">0x78c1...d2bA</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-lps">Liquidity Providers (30%)</span><span class="ic-val p" style="font-size:11px">0xc181...01EB</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-ubi">UBI Pool (20%)</span><span class="ic-val p" style="font-size:11px">0x4A9b...054A</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-treasury">Treasury (10%)</span><span class="ic-val p" style="font-size:11px">0x2273...3eb15</span></div>
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
    <div class="hlbox" data-i18n="ca-text">Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier (Groth16): 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Main):     0xE832Ac8Fa64F1AE2c6a5fE5d7DFbF0f9475ec0ae<br>V5 (Sepolia legacy):   0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5</div>
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
    <div class="hlbox" data-i18n="ubi-box">Sources: Transaction fees (20%) · Wealth cap overflow · Demurrage · Inactive escrow<br><br>Daily: UBI Pool divided equally among all registered humans. Pool resets to zero after each distribution and refills continuously from protocol fees.</div>
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
const V7_CONTRACT = '0xE832Ac8Fa64F1AE2c6a5fE5d7DFbF0f9475ec0ae';
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
  'h-sybil':'Sybil Attack Prevention','h-sybil-t':'Each biometric hash is stored permanently. Attempting to register twice is immediately rejected. One human, one wallet, forever. ⚠ Test phase: current verification is device-bound. A physiological sensor (MAX30102 PPG) is planned for device-independent identity in a future update.',
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
  'ca-title':'Contract Addresses','ca-text':'Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Main): 0xE832Ac8Fa64F1AE2c6a5fE5d7DFbF0f9475ec0ae<br>V5 Sepolia:  0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5',
  'poa-title':'1. PROOF OF ALIVE','poa-text':'<p>What happens to money when people die or disappear? In Bitcoin, millions of BTC are permanently lost. In Aequitas, if someone disappears, their AEQ eventually returns to the community.</p>',
  'poa-box':'Year 0-2: Normal usage<br>Year 2: Warning 1 — Guardian can respond<br>Year 2+60d: Warning 2<br>Year 2+120d: Warning 3<br>Year 2+180d: AEQ goes to PERSONAL ESCROW<br>Year 4: If still inactive — UBI Pool',
  'guard-title':'2. GUARDIAN SYSTEM','guard-text':'<p>What if someone cannot access their device for months? A trusted Guardian can confirm they are still alive — without any transaction rights.</p>',
  'guard-box':'1 Guardian per human (another verified human)<br>Guardian can ONLY call confirmAlive() — zero transaction rights<br>Guardian CANNOT move funds or transfer AEQ<br>Max 3 wards · 7-day timelock · No circular relationships',
  'dem-title':'3. DEMURRAGE — Anti-Hoarding',
  'dem-box':'1% annual fee on balance ABOVE fairShare goes to UBI Pool<br><br>Example: fairShare=1,000 · Balance=3,000 · Excess=2,000 · Monthly fee=1.67 AEQ',
  'dem-text':'<p>Historical precedent: Worgl, Austria (1932) — demurrage currency reduced unemployment 25% in one year.</p>',
  'cap-title':'4. WEALTH CAP','cap-box':'Phase 0: 50x fairShare · Phase 1: 20x · Phase 2: 10x · Phase 3: 5x · Phase 4: 3x<br><br>Always active from human #1. Excess instantly redistributed to ALL active humans.',
  'ubi-title':'5. UNIVERSAL BASIC INCOME','ubi-box':'Sources: Transaction fees (20%) · Wealth cap overflow · Demurrage · Inactive escrow<br><br>Daily: UBI Pool divided equally among all registered humans. Pool resets after each distribution.',
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
  'h-sybil':'Schutz vor Sybil-Angriffen','h-sybil-t':'Jeder biometrische Hash wird dauerhaft gespeichert. Doppelregistrierung wird sofort abgelehnt. Eine Person, eine Wallet, für immer. ⚠ Testphase: Aktuelle Verifizierung ist gerätegebunden. Ein physiologischer Sensor (MAX30102 PPG) ist für geräteunabhängige Identität in einem zukünftigen Update geplant.',
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
  'ca-title':'Contract-Adressen','ca-text':'Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Main): 0xE832Ac8Fa64F1AE2c6a5fE5d7DFbF0f9475ec0ae<br>V5 Sepolia:  0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5',
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
  'ca-title':'Direcciones de Contratos','ca-text':'Cadena: Aequitas Chain (Chain ID: 1926)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Main): 0xE832Ac8Fa64F1AE2c6a5fE5d7DFbF0f9475ec0ae',
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
  'ca-title':'Адреса Контрактов','ca-text':'Цепочка: Aequitas Chain (Chain ID: 1926)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Main): 0xE832Ac8Fa64F1AE2c6a5fE5d7DFbF0f9475ec0ae',
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
  'ca-title':'合约地址','ca-text':'链：Aequitas Chain（Chain ID: 1926）<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Main): 0xE832Ac8Fa64F1AE2c6a5fE5d7DFbF0f9475ec0ae',
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
  'ca-title':'Alamat Kontrak','ca-text':'Rantai: Aequitas Chain (Chain ID: 1926)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Main): 0xE832Ac8Fa64F1AE2c6a5fE5d7DFbF0f9475ec0ae',
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
  if (name === 'swap') loadPoolStatus();
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
        // AEQ is declared here as the chain's native currency (like ETH on
        // Ethereum) — MetaMask shows this automatically in the main
        // account balance display once eth_getBalance returns real
        // values, no further setup needed. We previously ALSO called
        // wallet_watchAsset below to add AEQ again as a separate ERC20
        // custom token. That meant AEQ showed up twice in MetaMask: once
        // correctly as the native balance, and once as an ERC20 entry
        // whose balance came from the V7 contract's balanceOf() mapping
        // instead — two numbers for "your AEQ" that could drift apart
        // (e.g. after a native transfer, only the native number changes,
        // while the ERC20 entry still shows the contract's value). Now
        // that registration and transfers write to the native balance,
        // the ERC20 entry no longer reflects the real, current state and
        // has been removed.
        nativeCurrency: { name: 'AEQ', symbol: 'AEQ', decimals: 18 },
        rpcUrls: ['https://aequitas-production-9fba.up.railway.app/rpc'],
        blockExplorerUrls: ['https://aequitas-production-9fba.up.railway.app']
      }]
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

    // Pool balances
    const fmtPool = v => v && v !== '0.0000' ? v + ' AEQ' : '—';
    document.getElementById('pool-v').textContent = fmtPool(d.pool_validators);
    document.getElementById('pool-l').textContent = fmtPool(d.pool_lp);
    document.getElementById('pool-u').textContent = fmtPool(d.pool_ubi);
    document.getElementById('pool-t').textContent = fmtPool(d.pool_treasury);

    // Fix stale subtitle now that demurrage/wealth-cap mean supply can drift
    const subEl = document.getElementById('s-supply-sub');
    if (subEl) subEl.textContent = 'Circulating across all accounts';

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

// ── SWAP TAB ─────────────────────────────────────────────────────────────
let swapWaddr = null;
let swapDirection = 'aeq_to_tusd';
let currentPoolAEQ = 0;
let currentPoolTUSD = 0;
let myAEQBalance = 0;
let myTUSDBalance = 0;

function swapLog(msg, type) {
  const el = document.getElementById('swap-log');
  el.innerHTML += '<div><span class="' + (type || 'info') + '">' + msg + '</span></div>';
  el.scrollTop = el.scrollHeight;
}

async function loadPoolStatus() {
  try {
    const d = await (await fetch('/api/pool')).json();
    currentPoolAEQ = d.reserve_aeq;
    currentPoolTUSD = d.reserve_tusd;
    document.getElementById('pool-reserve-aeq').textContent = fmt(d.reserve_aeq) + ' AEQ';
    document.getElementById('pool-reserve-tusd').textContent = fmt(d.reserve_tusd) + ' tUSD';
    document.getElementById('pool-price').textContent = d.reserve_aeq > 0
      ? ('1 AEQ ≈ ' + d.price_aeq_in_tusd.toFixed(4) + ' tUSD')
      : 'No liquidity yet';
    const desc = document.getElementById('swap-addliq-desc');
    if (desc) {
      desc.textContent = d.reserve_aeq > 0
        ? ('Pool ratio: 1 AEQ ≈ ' + d.price_aeq_in_tusd.toFixed(4) + ' tUSD — match this ratio when depositing')
        : 'Be the first to deposit — your ratio sets the starting price.';
    }
    updateFeeEstimate();
  } catch (e) {}
}

function setSwapDirection(dir) {
  swapDirection = dir;
  const a2t = document.getElementById('swap-dir-a2t');
  const t2a = document.getElementById('swap-dir-t2a');
  if (dir === 'aeq_to_tusd') {
    a2t.style.background = 'var(--gold)'; a2t.style.color = '#050A14';
    t2a.style.background = ''; t2a.style.color = '';
  } else {
    t2a.style.background = 'var(--gold)'; t2a.style.color = '#050A14';
    a2t.style.background = ''; a2t.style.color = '';
  }
  updateFeeEstimate();
}

// Mirrors the same constant-product math the server uses (see swapLocked
// in state.go), so the UI can warn BEFORE asking for a signature instead
// of after a wasted MetaMask popup. This is just for live feedback —
// the server still re-validates for real when the swap actually submits,
// since the pool could change between typing and submitting.
function estimateSwapOutput(amountIn, aeqToTusd) {
  if (amountIn <= 0 || currentPoolAEQ <= 0 || currentPoolTUSD <= 0) return null;
  const fee = amountIn * 0.001;
  const amountInAfterFee = amountIn - fee;
  let amountOut, reserveOut;
  if (aeqToTusd) {
    amountOut = (currentPoolTUSD * amountInAfterFee) / (currentPoolAEQ + amountInAfterFee);
    reserveOut = currentPoolTUSD;
  } else {
    amountOut = (currentPoolAEQ * amountInAfterFee) / (currentPoolTUSD + amountInAfterFee);
    reserveOut = currentPoolAEQ;
  }
  return { amountOut, fee, tooLarge: amountOut >= reserveOut };
}

function updateFeeEstimate() {
  const amt = parseFloat(document.getElementById('swap-amount').value || '0');
  const unit = swapDirection === 'aeq_to_tusd' ? 'AEQ' : 'tUSD';
  const outUnit = swapDirection === 'aeq_to_tusd' ? 'tUSD' : 'AEQ';
  const fee = amt * 0.001;
  document.getElementById('swap-fee-est').textContent = fee > 0 ? (fee.toFixed(6) + ' ' + unit) : '—';

  const goBtn = document.getElementById('swap-btn-go');
  const warnEl = document.getElementById('swap-warn');
  if (currentPoolAEQ <= 0 || currentPoolTUSD <= 0) {
    warnEl.textContent = '⚠ Pool has no liquidity yet — deposit some below before swapping.';
    warnEl.style.display = 'block';
    if (swapWaddr) goBtn.disabled = true;
    return;
  }
  if (amt <= 0) {
    warnEl.style.display = 'none';
    if (swapWaddr) goBtn.disabled = false;
    return;
  }
  const est = estimateSwapOutput(amt, swapDirection === 'aeq_to_tusd');
  if (est && est.tooLarge) {
    // Binary-search the largest input that stays safely under the
    // reserve, so the warning can suggest a concrete number instead of
    // just saying "too much" — 99% of the output reserve as a safety
    // margin, since the pool could shift slightly before this submits.
    let lo = 0, hi = amt;
    for (let i = 0; i < 30; i++) {
      const mid = (lo + hi) / 2;
      const midEst = estimateSwapOutput(mid, swapDirection === 'aeq_to_tusd');
      if (midEst && midEst.amountOut < (swapDirection === 'aeq_to_tusd' ? currentPoolTUSD : currentPoolAEQ) * 0.99) lo = mid;
      else hi = mid;
    }
    warnEl.innerHTML = '⚠ Too large for current pool liquidity. Try up to ~' + lo.toFixed(4) + ' ' + unit + '.';
    warnEl.style.display = 'block';
    if (swapWaddr) goBtn.disabled = true;
  } else if (est) {
    warnEl.innerHTML = 'You will receive ≈ ' + est.amountOut.toFixed(6) + ' ' + outUnit;
    warnEl.style.display = 'block';
    if (swapWaddr) goBtn.disabled = false;
  }
}

async function connectSwapWallet() {
  if (!window.ethereum) {
    swapLog('MetaMask not found. Please install MetaMask.', 'err');
    return;
  }
  try {
    await addToMetaMask();
    const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
    swapWaddr = accounts[0];
    document.getElementById('swap-wbox').style.display = 'block';
    document.getElementById('swap-wadr').textContent = swapWaddr;
    const btn = document.getElementById('swap-btn-conn');
    btn.textContent = swapWaddr.slice(0, 10) + '...' + swapWaddr.slice(-4);
    btn.style.background = 'var(--green)';
    btn.style.color = '#050A14';
    await refreshSwapBalances();
    await loadLPPosition();
    document.getElementById('swap-btn-go').disabled = false;
    document.getElementById('swap-btn-faucet').disabled = false;
    document.getElementById('swap-btn-addliq').disabled = false;
    setSwapDirection('aeq_to_tusd');
  } catch (e) {
    swapLog('Connection failed: ' + e.message, 'err');
  }
}

async function refreshSwapBalances() {
  if (!swapWaddr) return;
  try {
    const br = await fetch('/api/balance?wallet=' + swapWaddr);
    const bd = await br.json();
    myAEQBalance = bd.balance || 0;
    myTUSDBalance = bd.tusd_balance || 0;
    document.getElementById('swap-bal-aeq').textContent = fmt(bd.balance) + ' AEQ';
    document.getElementById('swap-bal-tusd').textContent = fmt(bd.tusd_balance) + ' tUSD';
    showDemurrageNotice(bd);
  } catch (e) {}
}

// Surfaces the demurrage warning at "login" time (i.e. whenever the
// wallet connects/refreshes its balance) per the two-stage design: a
// one-time notice once the account enters the 14-day window (the server
// tracks whether this has already fired and won't repeat it), and a
// notice on every check once inside the final 7 days before decay
// actually starts. Once decay is active, a different, ongoing message
// is shown instead of either warning.
function showDemurrageNotice(bd) {
  const box = document.getElementById('demurrage-notice');
  if (!box) return;
  if (bd.demurrage_active) {
    box.style.display = 'block';
    box.innerHTML = '⏳ Part of your idle AEQ balance is now slowly decaying (0.5%/month) because it hasn\'t been used in over 3 months. Send, swap, or deposit any amount to reset the clock.';
  } else if (bd.show_7_day_notice) {
    box.style.display = 'block';
    box.innerHTML = '⏳ Your AEQ balance will start decaying in ' + bd.demurrage_days_until_start.toFixed(1) + ' days unless you send, swap, or deposit some of it.';
  } else if (bd.show_14_day_notice) {
    box.style.display = 'block';
    box.innerHTML = '💡 Heads up: if this balance stays untouched, it will start slowly decaying in about 2 weeks. Any transfer, swap, or deposit resets the countdown.';
  } else {
    box.style.display = 'none';
  }
}

// Fills the AddLiquidity input for side ('aeq' or 'tusd') with pct of
// the user's own balance for that currency (0.25/0.5/0.75/1 = 25/50/75/
// 100%). Triggers the existing ratio-matching logic afterward so the
// OTHER field auto-fills too, exactly as if the user had typed it
// themselves — same behavior, just one click instead of a calculator.
function setPctAmount(side, pct) {
  if (side === 'aeq') {
    const amt = myAEQBalance * pct;
    document.getElementById('addliq-aeq').value = amt > 0 ? amt.toFixed(6) : '';
    updateLiquidityRatio('aeq');
  } else {
    const amt = myTUSDBalance * pct;
    document.getElementById('addliq-tusd').value = amt > 0 ? amt.toFixed(6) : '';
    updateLiquidityRatio('tusd');
  }
}

// Signs a fixed, human-readable message describing exactly what's being
// authorized — the wallet owner sees this in MetaMask's signing prompt
// before approving, and the server checks the signature matches both the
// claimed wallet AND this exact message (see verifyPersonalSign in swap.go).
async function signMessage(message) {
  return await window.ethereum.request({
    method: 'personal_sign',
    params: [message, swapWaddr]
  });
}

async function doSwap() {
  if (!swapWaddr) return;
  const amount = parseFloat(document.getElementById('swap-amount').value || '0');
  if (amount <= 0) { swapLog('Enter a valid amount', 'err'); return; }

  document.getElementById('swap-btn-go').disabled = true;
  try {
    const message = 'Aequitas Swap: ' + swapDirection + ' ' + amount.toFixed(8);
    swapLog('Sign the message in MetaMask to confirm this swap...', 'info');
    const signature = await signMessage(message);

    const resp = await fetch('/api/swap', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ wallet: swapWaddr, direction: swapDirection, amount, signature })
    });
    const data = await resp.json();
    if (data.success) {
      swapLog('✓ Swapped! Received ' + data.amount_out.toFixed(6) + ' ' + (swapDirection === 'aeq_to_tusd' ? 'tUSD' : 'AEQ'), 'ok');
      document.getElementById('swap-bal-aeq').textContent = fmt(data.new_aeq_balance) + ' AEQ';
      document.getElementById('swap-bal-tusd').textContent = fmt(data.new_tusd_balance) + ' tUSD';
      loadPoolStatus();
    } else {
      swapLog('✗ Swap failed: ' + data.message, 'err');
    }
  } catch (e) {
    swapLog('✗ Error: ' + e.message, 'err');
  }
  document.getElementById('swap-btn-go').disabled = false;
}

async function claimFaucet() {
  if (!swapWaddr) return;
  document.getElementById('swap-btn-faucet').disabled = true;
  try {
    const message = 'Aequitas tUSD Faucet Claim: ' + swapWaddr.toLowerCase();
    swapLog('Sign the message in MetaMask to claim test-tUSD...', 'info');
    const signature = await signMessage(message);

    const resp = await fetch('/api/faucet', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ wallet: swapWaddr, signature })
    });
    const data = await resp.json();
    if (data.success) {
      swapLog('✓ Claimed ' + data.granted + ' test-tUSD', 'ok');
      document.getElementById('swap-bal-tusd').textContent = fmt(data.granted) + ' tUSD';
    } else {
      swapLog('✗ Faucet claim failed: ' + data.message, 'err');
      document.getElementById('swap-btn-faucet').disabled = false;
    }
  } catch (e) {
    swapLog('✗ Error: ' + e.message, 'err');
    document.getElementById('swap-btn-faucet').disabled = false;
  }
}

// When the pool already has liquidity, typing one amount auto-fills the
// other at the pool's current ratio — matches what AddLiquidity itself
// requires (within 1% tolerance), so users don't have to calculate it
// by hand and then get rejected for a slightly-off ratio.
function updateLiquidityRatio(changed) {
  if (currentPoolAEQ <= 0 || currentPoolTUSD <= 0) return; // first depositor sets any ratio
  const aeqInput = document.getElementById('addliq-aeq');
  const tusdInput = document.getElementById('addliq-tusd');
  if (changed === 'aeq') {
    const aeq = parseFloat(aeqInput.value || '0');
    if (aeq > 0) tusdInput.value = (aeq * (currentPoolTUSD / currentPoolAEQ)).toFixed(6);
  } else {
    const tusd = parseFloat(tusdInput.value || '0');
    if (tusd > 0) aeqInput.value = (tusd * (currentPoolAEQ / currentPoolTUSD)).toFixed(6);
  }
}

async function doAddLiquidity() {
  if (!swapWaddr) return;
  const amountAEQ = parseFloat(document.getElementById('addliq-aeq').value || '0');
  const amountTUSD = parseFloat(document.getElementById('addliq-tusd').value || '0');
  if (amountAEQ <= 0 || amountTUSD <= 0) { swapLog('Enter both AEQ and tUSD amounts', 'err'); return; }

  document.getElementById('swap-btn-addliq').disabled = true;
  try {
    const message = 'Aequitas Add Liquidity: ' + amountAEQ.toFixed(8) + ' AEQ + ' + amountTUSD.toFixed(8) + ' tUSD';
    swapLog('Sign the message in MetaMask to confirm this deposit...', 'info');
    const signature = await signMessage(message);

    const resp = await fetch('/api/add-liquidity', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ wallet: swapWaddr, amount_aeq: amountAEQ, amount_tusd: amountTUSD, signature })
    });
    const data = await resp.json();
    if (data.success) {
      swapLog('✓ Liquidity added: ' + amountAEQ + ' AEQ + ' + amountTUSD + ' tUSD', 'ok');
      document.getElementById('addliq-aeq').value = '';
      document.getElementById('addliq-tusd').value = '';
      await refreshSwapBalances();
      await loadPoolStatus();
      await loadLPPosition();
    } else {
      swapLog('✗ Add liquidity failed: ' + data.message, 'err');
    }
  } catch (e) {
    swapLog('✗ Error: ' + e.message, 'err');
  }
  document.getElementById('swap-btn-addliq').disabled = false;
}

// ── LP POSITION / REMOVE LIQUIDITY ──────────────────────────────────────
let myLPShares = 0;
let myFullWithdrawableAEQ = 0;
let myFullWithdrawableTUSD = 0;

async function loadLPPosition() {
  if (!swapWaddr) return;
  try {
    const d = await (await fetch('/api/lp-position?wallet=' + swapWaddr)).json();
    myLPShares = d.shares || 0;
    myFullWithdrawableAEQ = d.withdrawable_aeq || 0;
    myFullWithdrawableTUSD = d.withdrawable_tusd || 0;
    const box = document.getElementById('lp-position-box');
    if (myLPShares > 0) {
      box.style.display = 'block';
      document.getElementById('lp-share-pct').textContent = d.pool_share_pct.toFixed(4) + '%';
      document.getElementById('lp-withdrawable').textContent =
        d.withdrawable_aeq.toFixed(4) + ' AEQ + ' + d.withdrawable_tusd.toFixed(4) + ' tUSD';
      updateRemovePreview();
    } else {
      box.style.display = 'none';
    }
  } catch (e) {}
}

// Recomputes "you will receive" from the currently selected removePct —
// called whenever removePct changes, whether from a percentage button or
// the manual input field, so both paths stay in sync with the same preview.
function updateRemovePreview() {
  const aeq = myFullWithdrawableAEQ * removePct;
  const tusd = myFullWithdrawableTUSD * removePct;
  document.getElementById('lp-remove-preview').textContent =
    aeq.toFixed(6) + ' AEQ + ' + tusd.toFixed(6) + ' tUSD';
}

// Manual percentage input — lets someone type e.g. "37.5" instead of only
// having the 25/50/75/100 quick buttons. Clears the active button
// highlighting since a manual value generally won't match one exactly.
function setRemovePctManual(value) {
  const pct = parseFloat(value || '0');
  if (pct < 0 || pct > 100 || isNaN(pct)) return;
  removePct = pct / 100;
  document.querySelectorAll('#lp-position-box .pctbtn').forEach(b => { b.style.background = ''; b.style.color = ''; });
  updateRemovePreview();
}

// Stores the chosen withdrawal fraction (set by the 25/50/75/MAX buttons)
// so doRemoveLiquidity knows how many shares to burn without needing a
// raw share-count input field — most people think in "withdraw half my
// position", not in the underlying share units.
let removePct = 1;
function setRemovePct(pct, btn) {
  removePct = pct;
  document.querySelectorAll('#lp-position-box .pctbtn').forEach(b => { b.style.background = ''; b.style.color = ''; });
  if (btn) { btn.style.background = 'var(--gold)'; btn.style.color = '#050A14'; }
  document.getElementById('remove-pct-input').value = (pct * 100).toString();
  updateRemovePreview();
}

async function doRemoveLiquidity() {
  if (!swapWaddr || myLPShares <= 0) return;
  const sharesToBurn = myLPShares * removePct;

  document.getElementById('swap-btn-removeliq').disabled = true;
  try {
    const message = 'Aequitas Remove Liquidity: ' + sharesToBurn.toFixed(8) + ' shares';
    swapLog('Sign the message in MetaMask to confirm this withdrawal...', 'info');
    const signature = await signMessage(message);

    const resp = await fetch('/api/remove-liquidity', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ wallet: swapWaddr, shares: sharesToBurn, signature })
    });
    const data = await resp.json();
    if (data.success) {
      swapLog('✓ Removed liquidity: received ' + data.amount_aeq.toFixed(4) + ' AEQ + ' + data.amount_tusd.toFixed(4) + ' tUSD', 'ok');
      await refreshSwapBalances();
      await loadPoolStatus();
      await loadLPPosition();
    } else {
      swapLog('✗ Remove liquidity failed: ' + data.message, 'err');
    }
  } catch (e) {
    swapLog('✗ Error: ' + e.message, 'err');
  }
  document.getElementById('swap-btn-removeliq').disabled = false;
}

document.addEventListener('DOMContentLoaded', () => {
  const amtInput = document.getElementById('swap-amount');
  if (amtInput) amtInput.addEventListener('input', updateFeeEstimate);
});

function checkProofParams() {
  const p = new URLSearchParams(window.location.search);
  const proofId = p.get('proofId');
  const proof = p.get('proof');
  const bioHash = p.get('bioHash');
  if (bioHash) {
    // NEW flow: the app only sent its biometric identity hash, not a
    // pre-made proof. We connect the wallet FIRST, then generate the ZK
    // proof ourselves with the now-known real wallet address — this is
    // what actually binds the proof to a specific wallet cryptographically
    // (previously the app called /prove with the zero address before any
    // wallet was even chosen, so the proof was never really tied to one).
    pendingBioHash = bioHash;
    document.querySelectorAll('.tab')[0].click();
    setTimeout(() => connectWalletAndProve(), 600);
  } else if (proofId) {
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

// Holds the biometric identity hash from the app while we wait for the
// wallet to connect — only used by the new bioHash flow above.
let pendingBioHash = null;

// New-flow counterpart to connectWallet(): connects MetaMask, and THEN
// calls /prove with the real wallet address now that we have one,
// instead of expecting an already-made proof to exist. This is the
// piece that actually closes the wallet-binding gap, since the proof's
// commitment now genuinely depends on which wallet asked for it.
async function connectWalletAndProve() {
  if (!window.ethereum) {
    addLog('MetaMask not found. Please install MetaMask.', 'err');
    return;
  }
  if (!pendingBioHash) {
    addLog('No biometric identity hash to prove — please retry from the app.', 'err');
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

    const br = await fetch('/api/balance?wallet=' + waddr);
    const bd = await br.json();
    if (bd.is_human) {
      addLog('Already registered! Balance: ' + bd.balance + ' AEQ', 'ok');
      document.getElementById('btn-reg').disabled = true;
      document.getElementById('btn-reg').textContent = 'ALREADY REGISTERED';
      return;
    }

    addLog('Wallet connected. Generating ZK proof for this wallet...', 'info');
    // salt generated here (browser, with crypto.getRandomValues — far
    // stronger than the app's old Math.random()-based salt) since this
    // is where the proof is now actually made.
    const saltBytes = new Uint8Array(32);
    crypto.getRandomValues(saltBytes);
    let saltBig = BigInt(0);
    for (let i = 0; i < saltBytes.length; i++) saltBig = (saltBig << BigInt(8)) | BigInt(saltBytes[i]);
    const FIELD_SIZE = BigInt("21888242871839275222246405745257275088548364400416034343698204186575808495617");
    const salt = (saltBig % FIELD_SIZE).toString();

    const proveResp = await fetch(PS + '/prove', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ bio: pendingBioHash, salt: salt, wallet: waddr })
    });
    if (!proveResp.ok) {
      const err = await proveResp.json();
      if (err.registered) {
        addLog('This identity is already registered.', 'ok');
        document.getElementById('btn-reg').disabled = true;
        document.getElementById('btn-reg').textContent = 'ALREADY REGISTERED';
        return;
      }
      addLog('Proof generation failed: ' + (err.error || 'unknown error'), 'err');
      return;
    }
    proofData = await proveResp.json();
    document.getElementById('pbox').style.display = 'block';
    document.getElementById('pval').textContent = 'Proof ready for ' + waddr.slice(0, 10) + '...';
    document.getElementById('btn-reg').disabled = false;
    document.getElementById('btn-reg').textContent = 'PROOF READY — CLICK TO REGISTER';
    addLog('Proof generated for your wallet. Click REGISTER to continue.', 'ok');
  } catch (e) {
    addLog('Connection failed: ' + e.message, 'err');
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
    addLog('Preparing signature...', 'info');
    document.getElementById('btn-reg').disabled = true;

    // commitment is pubSignals[0] — must match exactly what the contract reads
    const commitment = proofData.pubSignals[0];

    // Build the EXACT same hash the contract computes:
    // keccak256(abi.encodePacked(block.chainid, address(this), "register", commitment))
    const messageHash = ethers.solidityPackedKeccak256(
      ['uint256', 'address', 'string', 'uint256'],
      [1926, V7_CONTRACT, 'register', commitment]
    );

    addLog('Please sign the message in MetaMask to prove this wallet is yours (no gas, no cost)...', 'info');
    // personal_sign automatically adds the "\x19Ethereum Signed Message:\n32" prefix
    const signature = await window.ethereum.request({
      method: 'personal_sign',
      params: [messageHash, waddr]
    });

    addLog('Registering on Aequitas V7...', 'info');
    const r = await fetch('/api/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        wallet: waddr,
        pA: proofData.pA, pB: proofData.pB, pC: proofData.pC, pubSignals: proofData.pubSignals,
        signature: signature,
        bioHash: pendingBioHash || ''
      })
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
setInterval(loadPoolStatus, 8000);
</script>
</body>
</html>`
