package keeper

const explorerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0">
<meta name="google" content="notranslate">
<meta name="translation" content="notranslate">
<title>Aequitas — Proof of Humanity Chain</title>
<meta name="description" content="Aequitas Chain — a Proof of Humanity blockchain with built-in Universal Basic Income, demurrage, and wealth cap. Chain ID 1926.">
<meta name="theme-color" content="#0C0E16">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
<meta name="apple-mobile-web-app-title" content="Aequitas">
<link rel="preconnect" href="https://fonts.bunny.net">
<link href="https://fonts.bunny.net/css?family=inter:400,500,600,700,900|dm-serif-display:400|jetbrains-mono:400,600&display=swap" rel="stylesheet">
<style>
:root{
  --font-body:'Inter',-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;
  --font-display:'DM Serif Display',Georgia,serif;
  --font-mono:'JetBrains Mono','Fira Code',Consolas,monospace;
  --bg:#0C0E16;--card:#131620;--card2:#1A1D2B;--border:rgba(255,255,255,0.07);
  --purple:#9B72F6;--teal:#22D3EE;--neon:#34D399;--gold:#F0B429;
  --green:#34D399;--red:#F87171;--blue:#60A5FA;
  --text:#E8EDF5;--muted:#8892A4;
  --glow-purple:0 0 24px rgba(155,114,246,0.18);
  --glow-teal:0 0 20px rgba(34,211,238,0.14);
  --glow-strong:0 4px 32px rgba(155,114,246,0.28);
  --grad:linear-gradient(135deg,var(--purple),var(--teal));
  --radius:12px;--radius-sm:8px;
  --shadow:0 2px 16px rgba(0,0,0,0.4);
}
*{box-sizing:border-box;margin:0;padding:0}
body{background:var(--bg);color:var(--text);font-family:var(--font-body);min-height:100vh;overflow-x:hidden;line-height:1.5;background-image:radial-gradient(ellipse 80% 50% at 20% 0%,rgba(155,114,246,0.06) 0%,transparent 60%),radial-gradient(ellipse 60% 40% at 80% 100%,rgba(34,211,238,0.04) 0%,transparent 60%)}
body::before{display:none}
@keyframes starFloat{0%{transform:translateY(0)}100%{transform:translateY(-8px)}}
header{background:rgba(13,15,22,0.96);border-bottom:1px solid rgba(255,255,255,0.07);padding:0 24px;position:sticky;top:0;z-index:100;display:flex;align-items:center;justify-content:space-between;height:60px;gap:10px;box-shadow:0 1px 24px rgba(0,0,0,0.5);backdrop-filter:blur(12px)}
header::before{content:'';position:absolute;top:0;left:0;right:0;height:2px;background:var(--grad);opacity:0.8}
.logo-wrap{display:flex;align-items:center;gap:12px;flex-shrink:0;position:relative;z-index:1}
.logo-icon{width:34px;height:34px;border-radius:9px;background:var(--grad);display:flex;align-items:center;justify-content:center;font-size:17px;box-shadow:var(--glow-purple)}
.logo-text{font-size:1rem;font-weight:900;letter-spacing:3px;background:var(--grad);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text}
.logo-sub{font-size:0.48rem;color:var(--muted);letter-spacing:2.5px;text-transform:uppercase}
.header-right{display:flex;gap:8px;align-items:center;position:relative;z-index:1}
.badge{display:flex;align-items:center;gap:5px;padding:5px 11px;border-radius:20px;font-size:0.58rem;letter-spacing:0.5px;font-weight:600}
.badge-live{background:rgba(4,120,87,0.08);border:1px solid rgba(4,120,87,0.25);color:var(--neon)}
.badge-dag{background:rgba(107,70,193,0.08);border:1px solid rgba(107,70,193,0.2);color:var(--purple)}
.pulse{width:5px;height:5px;border-radius:50%;background:var(--neon);box-shadow:0 0 6px var(--neon);animation:pulse 2s infinite}
@keyframes pulse{0%,100%{opacity:1;transform:scale(1)}50%{opacity:0.4;transform:scale(0.7)}}
.lang-sel{background:rgba(255,255,255,0.06);color:var(--muted);border:1px solid rgba(255,255,255,0.1);border-radius:6px;padding:5px 10px;font-family:var(--font-body);font-size:0.62rem;outline:none;cursor:pointer}
.tabs{background:rgba(13,15,22,0.9);border-bottom:1px solid rgba(255,255,255,0.08);padding:0 24px;display:flex;overflow-x:auto;-webkit-overflow-scrolling:touch;scrollbar-width:none;gap:2px;position:relative;z-index:1}
.tabs::-webkit-scrollbar{display:none}
.tab{padding:16px 16px;font-size:0.65rem;color:var(--muted);cursor:pointer;border-bottom:2px solid transparent;letter-spacing:0.5px;font-weight:600;white-space:nowrap;transition:all 0.2s;flex-shrink:0}
.tab:hover{color:var(--purple)}.tab.active{color:var(--purple);border-bottom-color:var(--purple);background:rgba(155,114,246,0.08)}
.tab-content{display:none;position:relative;z-index:1}.tab-content.active{display:block}
html[data-active=network]  #tab-network {display:block!important}
html[data-active=explorer] #tab-explorer{display:block!important}
html[data-active=exchange] #tab-exchange{display:block!important}
html[data-active=index]    #tab-index   {display:block!important}
.hero{padding:20px 20px 0;position:relative;z-index:1}
.section-label{font-size:0.6rem;color:var(--muted);letter-spacing:3px;text-transform:uppercase;margin-bottom:14px;font-weight:600}
.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:1px;background:#E0D9D0;border:1px solid var(--border);border-radius:var(--radius);overflow:hidden;margin-bottom:20px;box-shadow:0 2px 12px rgba(0,0,0,0.06)}
.stat{background:var(--card);padding:20px 16px;position:relative;transition:all 0.2s;cursor:default}
.stat:hover{background:var(--card2);box-shadow:inset 0 0 20px rgba(107,70,193,0.04)}
.stat-accent{position:absolute;top:0;left:0;right:0;height:2px}
.stat-icon{font-size:1rem;margin-bottom:8px}
.stat-lbl{font-size:0.58rem;color:var(--muted);letter-spacing:1.5px;text-transform:uppercase;margin-bottom:6px;font-weight:500}
.stat-val{font-size:1.7rem;font-weight:900;line-height:1;margin-bottom:4px;font-family:var(--font-display)}
.stat-sub{font-size:0.57rem;color:var(--muted);line-height:1.5}
.c-green .stat-val{color:var(--neon)}.c-green .stat-accent{background:linear-gradient(90deg,var(--neon),transparent)}
.c-blue .stat-val{color:var(--teal)}.c-blue .stat-accent{background:linear-gradient(90deg,var(--teal),transparent)}
.c-gold .stat-val{color:var(--gold)}.c-gold .stat-accent{background:linear-gradient(90deg,var(--gold),transparent)}
.c-purple .stat-val{color:var(--purple)}.c-purple .stat-accent{background:linear-gradient(90deg,var(--purple),transparent)}
.c-teal .stat-val{color:var(--teal)}.c-teal .stat-accent{background:linear-gradient(90deg,var(--teal),transparent)}
.info-banner{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:20px;margin-bottom:20px;display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:20px;box-shadow:0 2px 12px rgba(0,0,0,0.06)}
.ib-icon{font-size:1.4rem;margin-bottom:8px}
.ib-title{font-size:0.7rem;color:var(--gold);font-weight:700;margin-bottom:8px;letter-spacing:0.5px}
.ib-text{font-size:0.65rem;color:var(--muted);line-height:1.8}
.main-grid{display:grid;grid-template-columns:1fr 310px;gap:16px;padding:0 20px 20px;position:relative;z-index:1}
@media(max-width:800px){.main-grid{grid-template-columns:1fr}.right-col{display:none}}
.section{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.06)}
.sec-head{padding:14px 18px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;background:var(--card2)}
.sec-title{font-size:0.65rem;color:var(--muted);letter-spacing:1px;text-transform:uppercase;display:flex;align-items:center;gap:8px;font-weight:600}
.sec-dot{width:6px;height:6px;border-radius:50%;background:var(--neon);box-shadow:0 0 8px var(--neon)}
.sec-count{font-size:0.6rem;color:var(--muted);background:var(--card2);padding:3px 8px;border-radius:10px;border:1px solid var(--border)}
.sec-desc{padding:10px 18px;font-size:0.65rem;color:var(--muted);background:rgba(139,92,246,0.03);border-bottom:1px solid var(--border);line-height:1.7}
.block-item{padding:12px 18px;border-bottom:1px solid rgba(139,92,246,0.08);display:grid;grid-template-columns:60px 1fr auto;gap:10px;align-items:center;transition:all 0.15s;cursor:pointer}
.block-item:hover{background:rgba(139,92,246,0.09)}.block-item:last-child{border-bottom:none}
.block-num{font-size:0.8rem;font-weight:700;color:var(--purple);font-family:var(--font-mono);text-shadow:0 0 8px rgba(139,92,246,0.4)}
.block-hash{font-size:0.63rem;color:var(--muted);margin-bottom:2px;display:flex;align-items:center;gap:4px;flex-wrap:wrap;font-family:var(--font-mono)}
.block-parents{font-size:0.57rem;color:rgba(139,92,246,0.3)}
.block-right{text-align:right}
.block-humans{font-size:0.65rem;color:var(--gold);margin-bottom:2px;font-weight:600}
.block-time{font-size:0.57rem;color:var(--neon)}
.block-detail-overlay{display:none;position:fixed;inset:0;background:rgba(0,0,0,0.8);z-index:1000;padding:20px;overflow-y:auto;backdrop-filter:blur(4px)}
.block-detail-overlay.open{display:flex;align-items:flex-start;justify-content:center;padding-top:50px}
.bdc{background:var(--card);border:1px solid rgba(139,92,246,0.3);border-radius:12px;width:100%;max-width:620px;overflow:hidden;box-shadow:0 0 40px rgba(139,92,246,0.15)}
.bdc-hdr{background:rgba(139,92,246,0.1);padding:14px 18px;display:flex;justify-content:space-between;align-items:center;border-bottom:1px solid rgba(139,92,246,0.15)}
.bdc-close{cursor:pointer;color:var(--muted);font-size:1.1rem;padding:4px 10px;border-radius:6px;background:rgba(139,92,246,0.1);border:1px solid var(--border);transition:all 0.15s}
.bdc-close:hover{color:var(--text);background:rgba(139,92,246,0.25)}
.bdc-row{padding:9px 18px;border-bottom:1px solid rgba(139,92,246,0.06);display:grid;grid-template-columns:130px 1fr;gap:8px;font-size:0.62rem}
.bdc-k{color:var(--muted);font-weight:600;padding-top:1px}
.bdc-v{color:var(--text);font-family:var(--font-mono);word-break:break-all;line-height:1.5}
.bdc-tx{margin:12px 18px;padding:9px 12px;background:rgba(0,230,118,0.04);border-radius:6px;border:1px solid rgba(0,230,118,0.15);font-size:0.59rem;font-family:var(--font-mono);color:var(--neon);word-break:break-all;line-height:1.6}
.bdc-tx-hdr{padding:10px 18px 4px;font-size:0.6rem;font-weight:700;color:var(--neon);text-transform:uppercase;letter-spacing:1px}
.bm{background:rgba(139,92,246,0.1);color:var(--purple);font-size:0.53rem;padding:2px 6px;border-radius:4px;border:1px solid rgba(139,92,246,0.2)}
.bt{background:rgba(0,255,209,0.08);color:var(--neon);font-size:0.53rem;padding:2px 6px;border-radius:4px;border:1px solid rgba(0,255,209,0.15)}
.empty{padding:40px;text-align:center;color:var(--muted);font-size:0.7rem;line-height:2.5}
.right-col{display:flex;flex-direction:column;gap:12px;position:relative;z-index:1}
.ic{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:18px;box-shadow:0 2px 12px rgba(0,0,0,0.06)}
.ic-title{font-size:0.6rem;color:var(--purple);letter-spacing:1.5px;text-transform:uppercase;margin-bottom:14px;font-weight:600}
.ic-row{display:flex;justify-content:space-between;align-items:center;padding:8px 0;border-bottom:1px solid rgba(139,92,246,0.08)}
.ic-row:last-child{border-bottom:none}
.ic-key{font-size:0.63rem;color:var(--muted)}
.ic-val{font-size:0.63rem;color:var(--text);text-align:right;max-width:58%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;font-weight:500}
.ic-val.g{color:var(--neon)}.ic-val.b{color:var(--teal)}.ic-val.go{color:var(--gold)}.ic-val.p{color:var(--purple)}
.mm-card{background:rgba(8,145,178,0.04);border:1px solid rgba(8,145,178,0.18);border-radius:var(--radius);padding:16px;box-shadow:0 2px 12px rgba(8,145,178,0.08)}
.mm-title{font-size:0.6rem;color:var(--teal);letter-spacing:1.5px;margin-bottom:12px;font-weight:700;text-transform:uppercase}
.mm-row{display:flex;justify-content:space-between;padding:6px 0;border-bottom:1px solid rgba(6,182,212,0.08)}
.mm-row:last-child{border-bottom:none}
.mm-key{font-size:0.6rem;color:var(--muted)}.mm-val{font-size:0.6rem;color:var(--purple);font-family:var(--font-mono)}
.mm-btn{width:100%;margin-top:12px;padding:11px;background:var(--grad);color:#fff;border:none;border-radius:var(--radius-sm);cursor:pointer;font-family:var(--font-body);font-size:0.68rem;font-weight:700;letter-spacing:0.5px;transition:all 0.2s;box-shadow:var(--glow-purple)}
.mm-btn:hover{opacity:0.87;transform:translateY(-1px);box-shadow:var(--glow-strong)}
.phil-card{background:linear-gradient(135deg,rgba(107,70,193,0.06),rgba(8,145,178,0.03));border:1px solid rgba(107,70,193,0.15);border-radius:var(--radius);padding:22px;text-align:center;box-shadow:0 2px 12px rgba(107,70,193,0.08)}
.phil-quote{font-size:0.85rem;color:var(--gold);font-style:italic;line-height:2;margin-bottom:6px;font-family:var(--font-display)}
.phil-sub{font-size:0.58rem;color:var(--muted);letter-spacing:1.5px;text-transform:uppercase}
.hs{padding:20px;display:grid;grid-template-columns:1fr 290px;gap:16px;position:relative;z-index:1}
@media(max-width:800px){.hs{grid-template-columns:1fr}}
.hi{padding:12px 18px;border-bottom:1px solid rgba(139,92,246,0.08);display:flex;align-items:center;gap:12px;transition:all 0.15s}
.hi:hover{background:rgba(139,92,246,0.05)}.hi:last-child{border-bottom:none}
.hav{width:36px;height:36px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:0.65rem;font-weight:bold;flex-shrink:0;border:2px solid;box-shadow:0 0 10px rgba(139,92,246,0.2)}
.hbal{font-size:0.82rem;color:var(--gold);font-weight:700;margin-bottom:1px;font-family:var(--font-display);text-shadow:0 0 10px rgba(245,166,35,0.3)}
.hadr{font-size:0.6rem;color:var(--muted);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;font-family:var(--font-mono)}
.hbdg{font-size:0.56rem;padding:3px 8px;border-radius:10px;flex-shrink:0;background:rgba(0,255,209,0.08);color:var(--neon);border:1px solid rgba(0,255,209,0.2);font-weight:600}
.is{padding:20px;display:grid;grid-template-columns:1fr 1fr;gap:16px;position:relative;z-index:1}
@media(max-width:700px){.is{grid-template-columns:1fr}}
.idx{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:24px;box-shadow:0 2px 12px rgba(0,0,0,0.06);transition:border-color 0.25s,box-shadow 0.25s}
.idx:hover{border-color:rgba(107,70,193,0.3);box-shadow:0 4px 20px rgba(107,70,193,0.10)}
.idx-title{font-size:0.6rem;color:var(--purple);letter-spacing:2px;text-transform:uppercase;margin-bottom:12px;font-weight:700;display:flex;align-items:center;gap:8px}
.ci-btn{padding:2px 8px;font-size:0.58rem;font-family:JetBrains Mono,monospace;background:rgba(139,92,246,0.08);border:1px solid rgba(139,92,246,0.2);color:var(--muted);border-radius:4px;cursor:pointer;transition:all .15s}
.ci-btn:hover{background:rgba(139,92,246,0.18);color:#c4b5fd}
.ci-btn.ci-active{background:rgba(139,92,246,0.22);border-color:rgba(139,92,246,0.6);color:#c4b5fd}
.idx-title::before{content:'';display:inline-block;width:3px;height:12px;background:linear-gradient(180deg,var(--purple),var(--teal));border-radius:2px;flex-shrink:0}
.idx-desc{font-size:0.67rem;color:var(--muted);line-height:1.8;margin-bottom:16px}
.idx-big{font-size:2.8rem;font-weight:900;line-height:1;font-family:var(--font-display);background:var(--grad);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text}
.idx-lbl{font-size:0.6rem;color:var(--muted);margin-top:4px}
.bar-bg{height:8px;background:rgba(139,92,246,0.1);border-radius:6px;overflow:hidden;margin:14px 0 6px;border:1px solid rgba(139,92,246,0.1)}
.bar-fill{height:100%;border-radius:6px;background:linear-gradient(90deg,var(--neon),var(--teal),var(--purple),var(--gold),var(--red));transition:width 1.5s ease;box-shadow:0 0 8px rgba(139,92,246,0.4)}
.bar-lbl{display:flex;justify-content:space-between;font-size:0.55rem;color:var(--muted)}
.mrow{display:grid;grid-template-columns:repeat(2,1fr);gap:8px;margin-top:14px}
.mbox{background:var(--card2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:12px;text-align:center;transition:all 0.2s}
.mbox:hover{border-color:rgba(107,70,193,0.3);box-shadow:0 4px 16px rgba(107,70,193,0.10)}
.mval{font-size:1.15rem;font-weight:700;color:var(--teal);font-family:var(--font-display);text-shadow:0 0 10px rgba(6,182,212,0.3)}
.mlbl{font-size:0.57rem;color:var(--muted);margin-top:3px;font-weight:500}
.story{font-size:0.7rem;line-height:2;color:var(--muted)}
.story p{margin-bottom:14px}
.hlbox{background:rgba(139,92,246,0.05);border-left:3px solid var(--purple);border-radius:0 var(--radius-sm) var(--radius-sm) 0;padding:14px 18px;margin:16px 0;font-size:0.67rem;color:var(--text);line-height:1.9}
.ns{padding:20px;display:grid;grid-template-columns:1fr 1fr;gap:16px;position:relative;z-index:1}
@media(max-width:700px){.ns{grid-template-columns:1fr}}
@media(max-width:600px){.stats-grid{grid-template-columns:1fr 1fr}canvas{max-width:100%!important}.tab{padding:10px 8px;font-size:0.58rem}.rhero{padding:12px 14px 0}.nc{padding:14px}.rs{padding:12px}.swap-form{padding:12px}.pool-cards{grid-template-columns:1fr}}
.nc{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:20px;box-shadow:0 2px 12px rgba(0,0,0,0.06);transition:border-color 0.25s,box-shadow 0.25s}
.nc:hover{border-color:rgba(107,70,193,0.3);box-shadow:0 4px 20px rgba(107,70,193,0.10)}
.nc-title{font-size:0.6rem;color:var(--purple);letter-spacing:1.5px;text-transform:uppercase;margin-bottom:14px;font-weight:700;display:flex;align-items:center;gap:8px}
.nc-title::before{content:'';display:inline-block;width:3px;height:12px;background:linear-gradient(180deg,var(--purple),var(--teal));border-radius:2px;flex-shrink:0}
.nbox{background:var(--card2);border-radius:var(--radius-sm);padding:14px;border:1px solid var(--border);margin-bottom:10px}
.nstat{display:flex;align-items:center;gap:6px;font-size:0.67rem;color:var(--neon);margin-bottom:5px;font-weight:600}
.ndot{width:7px;height:7px;border-radius:50%;background:var(--neon);box-shadow:0 0 8px var(--neon)}
.nurl{font-size:0.58rem;color:var(--muted);word-break:break-all;margin-bottom:3px;font-family:var(--font-mono)}
.ndesc{font-size:0.58rem;color:rgba(139,92,246,0.4)}
.spect{width:100%;border-collapse:collapse}
.spect td{padding:9px 4px;border-bottom:1px solid rgba(139,92,246,0.08);font-size:0.63rem;transition:background 0.15s}
.spect tr:hover td{background:rgba(139,92,246,0.05)}
.spect tr:last-child td{border-bottom:none}
.spect td:first-child{color:var(--muted);width:45%;padding-left:2px}
.spect td:last-child{text-align:right;font-family:var(--font-mono);font-size:0.6rem;color:var(--purple);padding-right:2px}
.bsbox{background:var(--card2);border-radius:var(--radius-sm);padding:12px;font-size:0.58rem;color:var(--purple);word-break:break-all;line-height:1.7;border:1px solid var(--border);font-family:var(--font-mono)}
.rs{padding:20px;max-width:600px;margin:0 auto;position:relative;z-index:1}
.rhero{background:linear-gradient(135deg,rgba(139,92,246,0.1),rgba(6,182,212,0.05));border:1px solid rgba(139,92,246,0.25);border-radius:var(--radius);padding:24px;margin-bottom:16px;text-align:center;box-shadow:var(--glow-purple)}
.rhero-title{font-size:1.05rem;font-weight:700;color:var(--text);margin-bottom:8px;font-family:var(--font-display)}
.rhero-sub{font-size:0.67rem;color:var(--muted);line-height:1.9}
.aonly{background:rgba(139,92,246,0.06);border:1px solid rgba(139,92,246,0.2);border-radius:var(--radius);padding:20px;text-align:center;margin-bottom:16px}
.aonly-icon{font-size:2rem;margin-bottom:8px}
.aonly-title{font-size:0.7rem;color:var(--purple);font-weight:700;letter-spacing:1px;margin-bottom:10px;text-shadow:0 0 10px rgba(139,92,246,0.4)}
.aonly-text{font-size:0.65rem;color:var(--muted);line-height:1.9}
.rsteps{display:grid;grid-template-columns:repeat(4,1fr);gap:8px;margin-bottom:16px}
@media(max-width:520px){.rsteps{grid-template-columns:repeat(2,1fr)}}
.rstep{background:var(--card);border:1px solid var(--border);border-radius:var(--radius-sm);padding:16px;text-align:center;transition:all 0.2s}
.rstep:hover{border-color:rgba(139,92,246,0.4);box-shadow:var(--glow-purple);transform:translateY(-2px)}
.snum{width:28px;height:28px;background:var(--grad);border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 10px;font-weight:700;font-size:0.72rem;color:#fff;box-shadow:var(--glow-purple)}
.stitle{font-size:0.63rem;color:var(--text);font-weight:600;margin-bottom:5px}
.sdesc{font-size:0.6rem;color:var(--muted);line-height:1.7}
.pbar{background:rgba(0,255,209,0.06);border:1px solid rgba(0,255,209,0.15);border-radius:var(--radius-sm);padding:10px 14px;margin-bottom:14px;font-size:0.63rem;color:var(--neon);text-align:center;line-height:1.8}
.rcard{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:20px;margin-bottom:14px;box-shadow:0 2px 12px rgba(0,0,0,0.06)}
.wbox{background:rgba(0,255,209,0.05);border:1px solid rgba(0,255,209,0.15);border-radius:var(--radius-sm);padding:10px;margin-bottom:10px;display:none}
.wlbl{font-size:0.57rem;color:var(--muted);margin-bottom:2px;letter-spacing:1px;font-weight:500}
.wadr{font-size:0.72rem;color:var(--neon);font-weight:700;font-family:var(--font-mono);text-shadow:0 0 8px rgba(0,255,209,0.3)}
.pbox{background:var(--card2);border:1px solid rgba(245,166,35,0.15);border-radius:var(--radius-sm);padding:10px;margin-bottom:10px;display:none}
.plbl{font-size:0.57rem;color:var(--gold);margin-bottom:2px;letter-spacing:1px;font-weight:500}
.pval{font-size:0.63rem;color:var(--muted);font-family:var(--font-mono)}
.rbtn{width:100%;padding:14px;border-radius:var(--radius-sm);border:none;cursor:pointer;font-family:var(--font-body);font-size:0.74rem;font-weight:700;letter-spacing:0.3px;transition:all 0.2s;margin-bottom:8px}
.rbtn:hover{opacity:0.9;transform:translateY(-1px)}
.bc{background:var(--grad);color:#fff;box-shadow:var(--glow-purple)}.bc:hover{box-shadow:var(--glow-strong)}
.br{background:linear-gradient(135deg,var(--gold),#E67E00);color:#fff;box-shadow:0 0 15px rgba(245,166,35,0.3)}
.rbtn:disabled{opacity:0.3;cursor:not-allowed;transform:none;box-shadow:none}
.rlog{background:var(--card2);border-radius:var(--radius-sm);padding:12px;font-size:0.63rem;line-height:2;min-height:52px;border:1px solid var(--border);font-family:var(--font-mono)}
.rlog .ok{color:var(--neon)}.rlog .err{color:var(--red)}.rlog .info{color:var(--gold)}
.ps{padding:20px;max-width:800px;margin:0 auto;position:relative;z-index:1}
.pct-row{display:flex;gap:6px;margin-bottom:8px}
.pctbtn{flex:1;padding:8px;font-size:12px;background:var(--card2);border:1px solid var(--border);color:var(--text);border-radius:var(--radius-sm);cursor:pointer;font-family:var(--font-body);font-weight:600;transition:all 0.2s}
.pctbtn:hover{border-color:var(--purple);color:var(--purple);box-shadow:var(--glow-purple)}
#demurrage-notice{font-size:13px;padding:12px 14px;border-radius:var(--radius-sm);background:rgba(245,166,35,0.06);border:1px solid rgba(245,166,35,0.2);color:var(--gold);margin:10px 0;line-height:1.7}
.swap-dir{background:var(--card2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:8px;cursor:pointer;font-size:1rem;transition:all 0.2s;width:100%;margin:8px 0}
.swap-dir:hover{border-color:var(--purple);box-shadow:var(--glow-purple)}
input[type=number]{background:var(--card2);border:1px solid var(--border);color:var(--text);border-radius:var(--radius-sm);padding:10px 12px;font-family:var(--font-body);font-size:0.8rem;outline:none;transition:all 0.2s}
input[type=number]:focus{border-color:var(--purple);box-shadow:0 0 8px rgba(139,92,246,0.2)}
input[type=number]::-webkit-inner-spin-button{opacity:0.5}
@media(max-width:480px){.stats-grid{grid-template-columns:repeat(2,1fr)}.stat-val{font-size:1.4rem}header{height:52px}.logo-text{font-size:0.85rem;letter-spacing:2px}.badge-dag{display:none}.main-grid{padding:0 12px 12px}.hero{padding:14px 12px 0}.tab{padding:12px 10px;font-size:0.6rem}}@media(max-width:600px){.idx-grade-grid{grid-template-columns:repeat(2,1fr)!important}}
/* ── SWAP ENHANCEMENTS ────────────────────────────────────── */
.sd-panel{background:var(--card2);border:1px solid rgba(139,92,246,0.18);border-radius:var(--radius-sm);padding:13px;margin:8px 0;animation:sdIn 0.18s ease}
@keyframes sdIn{from{opacity:0;transform:translateY(-4px)}to{opacity:1;transform:translateY(0)}}
.sd-row{display:flex;justify-content:space-between;align-items:center;padding:6px 0;font-size:0.63rem;border-bottom:1px solid rgba(139,92,246,0.07)}
.sd-row:last-child{border-bottom:none}
.sd-key{color:var(--muted)}
.sd-val{font-weight:600}
.impact-low{color:var(--neon)}.impact-med{color:var(--gold)}.impact-hi{color:var(--red)}
.sd-header{font-size:0.54rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;font-weight:600;margin-bottom:10px;display:flex;align-items:center;gap:6px}
.sd-header::before{content:'';display:inline-block;width:3px;height:10px;background:var(--purple);border-radius:2px}
/* ── POOL DEPTH BAR ──────────────────────────────────────── */
.depth-track{height:14px;border-radius:7px;overflow:hidden;display:flex;border:1px solid var(--border)}
.depth-aeq-fill{background:linear-gradient(90deg,var(--purple),rgba(139,92,246,0.55));transition:width 1.2s ease}
.depth-tusd-fill{background:linear-gradient(90deg,rgba(6,182,212,0.55),var(--teal));flex:1}
.depth-lbls{display:flex;justify-content:space-between;font-size:0.56rem;color:var(--muted);margin-top:5px}
.amm-box{background:rgba(139,92,246,0.04);border:1px solid rgba(139,92,246,0.13);border-radius:var(--radius-sm);padding:13px;margin-top:10px}
.amm-formula{font-size:0.67rem;color:var(--purple);font-family:var(--font-mono);text-align:center;padding:9px;background:rgba(139,92,246,0.09);border-radius:6px;margin:8px 0;border:1px solid rgba(139,92,246,0.13);letter-spacing:0.5px}
.amm-text{font-size:0.6rem;color:var(--muted);line-height:1.88}
/* ── UBI HERO ────────────────────────────────────────────── */
.ubi-hero-section{background:linear-gradient(135deg,rgba(245,166,35,0.1),rgba(139,92,246,0.06),rgba(0,255,209,0.04));border:1px solid rgba(245,166,35,0.3);border-radius:var(--radius);padding:24px;margin:14px 0;text-align:center;position:relative;overflow:hidden;box-shadow:0 0 40px rgba(245,166,35,0.07)}
.ubi-hero-section::before{content:'';position:absolute;top:0;left:0;right:0;height:2px;background:linear-gradient(90deg,var(--gold),var(--neon),var(--purple),var(--neon),var(--gold))}
.ubi-big-timer{font-size:2rem;font-weight:900;font-family:var(--font-mono);color:var(--gold);text-shadow:0 0 25px rgba(245,166,35,0.5);letter-spacing:3px;margin:8px 0}
.ubi-pool-amount{font-size:1.5rem;font-weight:700;font-family:var(--font-display);color:var(--neon);text-shadow:0 0 15px rgba(0,255,209,0.35);margin:4px 0}
.ubi-fill-track{height:7px;background:rgba(245,166,35,0.1);border-radius:4px;overflow:hidden;margin:12px auto;max-width:320px;border:1px solid rgba(245,166,35,0.18)}
.ubi-fill-bar{height:100%;background:linear-gradient(90deg,var(--gold),var(--neon));border-radius:4px;transition:width 2s ease;box-shadow:0 0 8px rgba(245,166,35,0.4);width:0%}
.ubi-src-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:8px;margin:12px 0}
@media(max-width:580px){.ubi-src-grid{grid-template-columns:1fr}}
.ubi-src-card{background:var(--card2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:14px;text-align:center;transition:all 0.2s}
.ubi-src-card:hover{transform:translateY(-2px);box-shadow:var(--glow-purple)}
.ubi-src-pct{font-size:1.25rem;font-weight:700;font-family:var(--font-display);margin-bottom:3px}
.ubi-src-name{font-size:0.6rem;font-weight:700;margin-bottom:5px;letter-spacing:0.3px}
.ubi-src-desc{font-size:0.57rem;color:var(--muted);line-height:1.75}
.pools4-grid{display:grid;grid-template-columns:1fr 1fr;gap:10px}
@media(max-width:580px){.pools4-grid{grid-template-columns:1fr}}
.pool4-card{background:var(--card2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:16px;transition:all 0.2s}
.pool4-card:hover{transform:translateY(-2px);box-shadow:var(--glow-purple)}
.pool4-head{display:flex;justify-content:space-between;align-items:center;margin-bottom:10px}
.pool4-name{font-size:0.61rem;font-weight:700;letter-spacing:0.3px}
.pool4-badge{font-size:0.56rem;color:var(--muted);background:var(--card);border:1px solid var(--border);padding:2px 8px;border-radius:10px}
.pool4-amount{font-size:1.05rem;font-weight:700;font-family:var(--font-display);margin-bottom:3px}
.pool4-timer{font-size:0.59rem;font-weight:600;margin-bottom:7px}
.pool4-desc{font-size:0.57rem;color:var(--muted);line-height:1.75}
/* ── EXPLORE CARDS ───────────────────────────────────────────── */
.expl-card{background:var(--card);border:1px solid var(--border);border-radius:var(--radius-sm);padding:14px;cursor:pointer;transition:all 0.2s}
.expl-card:hover{border-color:rgba(139,92,246,0.4);background:rgba(139,92,246,0.06);transform:translateY(-2px);box-shadow:var(--glow-purple)}
.expl-icon{font-size:1.1rem;margin-bottom:6px}
.expl-name{font-size:0.63rem;font-weight:700;color:var(--text);margin-bottom:4px}
.expl-desc{font-size:0.57rem;color:var(--muted);line-height:1.7}
/* ── SUB-TAB NAVIGATION ─────────────────────────────────────── */
.stabs{display:flex;gap:2px;padding:8px 20px 0;overflow-x:auto;background:var(--card2);border-bottom:1px solid var(--border);-webkit-overflow-scrolling:touch;scrollbar-width:none}
.stabs::-webkit-scrollbar{display:none}
.stab{padding:7px 15px;font-size:0.6rem;color:var(--muted);cursor:pointer;border-radius:6px 6px 0 0;letter-spacing:0.5px;font-weight:600;white-space:nowrap;transition:all 0.2s;border:1px solid transparent;border-bottom:none;flex-shrink:0;position:relative}
.stab:hover{color:var(--text);background:rgba(139,92,246,0.1)}
.stab.active{color:var(--purple);background:rgba(155,114,246,0.12);border-color:rgba(155,114,246,0.35)}
.stab.active::after{content:'';position:absolute;bottom:-1px;left:0;right:0;height:2px;background:linear-gradient(90deg,var(--purple),var(--teal));border-radius:2px 2px 0 0}
.stab-panel{display:none}
.stab-panel.active{display:block}
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
    <option value="fr">🌐 FR</option>
    <option value="pt">🌐 PT</option>
    <option value="ru">🌐 RU</option>
    <option value="zh">🌐 ZH</option>
    <option value="ar">🌐 AR</option>
    <option value="hi">🌐 HI</option>
    <option value="id">🌐 ID</option>
    <option value="it">🌐 IT</option>
    <option value="tr">🌐 TR</option>
  </select>
  <div class="header-right">
    <div class="badge badge-live"><span class="pulse"></span><span data-i18n="live">LIVE</span></div>
    <div class="badge badge-dag">● BLOCKDAG</div>
  </div>
</header>
<div class="tabs">
  <div class="tab active" onclick="showTab('register',this)">🔐 Register</div>
  <div class="tab" onclick="showTab('explorer',this)">🔍 Explorer</div>
  <div class="tab" onclick="showTab('index',this)">⚖️ Equality</div>
  <div class="tab" onclick="showTab('network',this)">🌐 Network</div>
  <div class="tab" onclick="showTab('exchange',this);setTimeout(drawPriceChart,50)">🔄 Exchange</div>
</div>

<!-- REGISTER -->
<div id="tab-register" class="tab-content active">
<div class="rs">
  <div class="rhero">
    <div class="rhero-title" data-i18n="reg-title">🔐 Register as a Verified Human</div>
    <div class="rhero-sub" data-i18n="reg-sub">Join the Aequitas network and receive your 1,000 AEQ Universal Basic Income grant. Registration is one-time, permanent, and completely gasless. No personal data is ever stored — only a cryptographic proof that you are a unique human being.</div>
    <a href="/download/app.apk" style="display:inline-flex;align-items:center;gap:10px;margin-top:18px;background:var(--grad);color:#fff;padding:13px 28px;border-radius:10px;font-size:0.75rem;font-weight:700;text-decoration:none;letter-spacing:0.5px;box-shadow:var(--glow-purple);transition:all 0.2s" onmouseover="this.style.opacity='0.87';this.style.transform='translateY(-2px)'" onmouseout="this.style.opacity='1';this.style.transform='translateY(0)'">
      <span style="font-size:1.1rem">📱</span>
      <span data-i18n="btn-download-app">DOWNLOAD AEQUITASBIO APP</span>
    </a>
    <div style="font-size:0.55rem;color:rgba(255,255,255,0.35);margin-top:8px">Android APK · direct download</div>
  </div>
  <!-- USP / EQUAL START SECTION -->
  <div style="background:linear-gradient(135deg,rgba(0,255,209,0.06),rgba(139,92,246,0.04));border:1px solid rgba(0,255,209,0.2);border-radius:var(--radius);padding:22px;margin-bottom:16px;position:relative;overflow:hidden">
    <div style="position:absolute;top:0;left:0;right:0;height:2px;background:linear-gradient(90deg,var(--neon),var(--purple))"></div>
    <div style="font-size:0.78rem;font-weight:700;font-family:var(--font-display);color:var(--neon);margin-bottom:4px;line-height:1.4" data-i18n="usp-headline">For the first time in history — everyone starts equal</div>
    <div style="font-size:0.62rem;color:var(--muted);margin-bottom:16px;line-height:1.7" data-i18n="usp-sub">If you own an Android smartphone, you qualify. No bank, no crypto background, no investment needed.</div>
    <div style="display:grid;grid-template-columns:repeat(2,1fr);gap:8px">
      <div style="background:var(--card);border:1px solid rgba(0,255,209,0.15);border-radius:var(--radius-sm);padding:12px">
        <div style="font-size:0.67rem;font-weight:700;color:var(--neon);margin-bottom:4px" data-i18n="usp-c1-title">0.00 Start Investment</div>
        <div style="font-size:0.59rem;color:var(--muted);line-height:1.75" data-i18n="usp-c1-desc">Registration is completely gasless. No ETH, no MATIC, no credit card. The protocol pays all fees on your behalf.</div>
      </div>
      <div style="background:var(--card);border:1px solid rgba(245,166,35,0.15);border-radius:var(--radius-sm);padding:12px">
        <div style="font-size:0.67rem;font-weight:700;color:var(--gold);margin-bottom:4px" data-i18n="usp-c2-title">1,000 AEQ for every human</div>
        <div style="font-size:0.59rem;color:var(--muted);line-height:1.75" data-i18n="usp-c2-desc">Billionaire or subsistence farmer — everyone gets exactly 1,000 AEQ. Not more, not less. Equal start, guaranteed by math.</div>
      </div>
      <div style="background:var(--card);border:1px solid rgba(139,92,246,0.15);border-radius:var(--radius-sm);padding:12px">
        <div style="font-size:0.67rem;font-weight:700;color:var(--purple);margin-bottom:4px" data-i18n="usp-c3-title">Accessible to all</div>
        <div style="font-size:0.59rem;color:var(--muted);line-height:1.75" data-i18n="usp-c3-desc">No bank account, no credit card, no government ID. Registration uses an affordable biometric hardware kit (fingerprint scanner + pulse sensor, ~$15) — designed for global access.</div>
      </div>
      <div style="background:var(--card);border:1px solid rgba(6,182,212,0.15);border-radius:var(--radius-sm);padding:12px">
        <div style="font-size:0.67rem;font-weight:700;color:var(--teal);margin-bottom:4px" data-i18n="usp-c4-title">Daily UBI forever</div>
        <div style="font-size:0.59rem;color:var(--muted);line-height:1.75" data-i18n="usp-c4-desc">Once registered, you receive a daily share of UBI payouts automatically — every day, no action required.</div>
      </div>
    </div>
  </div>

  <div class="aonly">
    <div class="aonly-icon">📱</div>
    <div class="aonly-title" data-i18n="app-title">REGISTRATION VIA ANDROID APP</div>
    <div class="aonly-text" data-i18n="app-text">Proof of Humanity requires physical biometric verification via a dedicated hardware kit. <strong>Phase 1:</strong> The R503 optical sensor scans <strong>all 10 fingerprints</strong> — combined into a single biometric hash. A MAX30102 PPG sensor confirms liveness (pulse detection) to prevent replay attacks. <strong>Phase 2</strong> adds hand vein pattern via IR camera (1 in 10⁷ uniqueness). <strong>Phase 3</strong> adds iris scan — the gold standard with 1 in 10⁷⁸ uniqueness, fully device-independent. Raw biometric data never leaves the device. The AequitasBio app generates a Groth16 Zero-Knowledge Proof that mathematically proves you are a unique living human — without revealing any personal information. Your <strong style="color:var(--gold)">1,000 AEQ will be credited automatically</strong> upon successful verification.</div>
  </div>
  <div class="rsteps">
    <div class="rstep"><div class="snum">1</div><div class="stitle" data-i18n="s1t">Biometric Scan</div><div class="sdesc" data-i18n="s1d">The AequitasBio app scans all 10 fingerprints via the R503 optical sensor and confirms liveness via MAX30102 PPG (pulse detection). Phase 2 adds hand vein IR scan. Phase 3 adds iris. Raw data never leaves the device.</div></div>
    <div class="rstep"><div class="snum">2</div><div class="stitle" data-i18n="s2t">ZK Proof Generation</div><div class="sdesc" data-i18n="s2d">A Groth16 Zero-Knowledge Proof is generated binding all biometric factors into a single commitment and nullifier. This cryptographically proves you are a unique living human — without revealing any biometric data.</div></div>
    <div class="rstep"><div class="snum">3</div><div class="stitle" data-i18n="s3t">Connect Wallet</div><div class="sdesc" data-i18n="s3d">The app opens MetaMask on this page. Connect your Ethereum wallet — this is the address that will receive your 1,000 AEQ grant. The proof is cryptographically bound to your wallet.</div></div>
    <div class="rstep"><div class="snum">4</div><div class="stitle" data-i18n="s4t">1,000 AEQ Granted</div><div class="sdesc" data-i18n="s4d">Your registration is confirmed on the Aequitas BlockDAG within 6 seconds. 1,000 AEQ is credited to your wallet instantly, gasless. Your identity is permanently recorded as verified human.</div></div>
  </div>
  <div class="pbar" data-i18n="priv-bar">🔒 R503 All-10-Fingerprint · MAX30102 Liveness · Phase 2: Hand Vein IR · Phase 3: Iris (10⁷⁸) · Groth16 ZKP · Data never leaves device · One human · Forever</div>
  <div class="pbar" style="background:rgba(245,166,35,0.06);border:1px solid rgba(245,166,35,0.2);color:var(--gold)">📱 MetaMask Mobile: if AEQ balance shows 0 after registration, go to Settings → Networks → delete Aequitas Chain → re-add via this website</div>
  <div class="rcard">
    <div class="wbox" id="wbox"><div class="wlbl" data-i18n="conn-wallet">CONNECTED WALLET</div><div class="wadr" id="wadr" title="">—</div><button onclick="copyAddr('wadr',this)" style="margin-top:4px;padding:3px 10px;font-size:0.56rem;background:rgba(0,255,209,0.08);border:1px solid rgba(0,255,209,0.2);color:var(--neon);border-radius:4px;cursor:pointer">📋 Copy</button></div>
    <div id="demurrage-notice" style="display:none"></div>
    <div class="pbox" id="pbox"><div class="plbl" data-i18n="proof-recv">⚡ ZK PROOF RECEIVED</div><div class="pval" id="pval" data-i18n="proof-hint">Connect wallet to register</div></div>
    <button class="rbtn bc" id="btn-conn" onclick="connectWallet()" data-i18n="btn-conn">🦊 CONNECT METAMASK</button>
    <button id="btn-disconnect" onclick="disconnectWallet()" style="display:none;margin-top:6px;padding:8px 16px;font-size:0.6rem;letter-spacing:1px;border:1px solid rgba(248,113,113,0.4);background:rgba(248,113,113,0.08);color:var(--red);border-radius:6px;cursor:pointer;width:100%">⊘ DISCONNECT WALLET</button>
    <button class="rbtn br" id="btn-reg" onclick="doRegister()" disabled data-i18n="btn-reg">🔐 REGISTER ON-CHAIN</button>
    <button class="rbtn" id="btn-web-reg" onclick="registerViaBrowser()" style="background:linear-gradient(135deg,#0ea5e9,#6366f1);color:#fff;margin-top:8px" data-i18n="btn-web-reg">🌐 REGISTER VIA BROWSER (WebAuthn)</button>
    <div id="web-reg-warn" style="display:none;font-size:0.62rem;color:#f59e0b;background:rgba(245,158,11,0.08);border:1px solid rgba(245,158,11,0.3);border-radius:6px;padding:8px 10px;margin-top:6px" data-i18n="web-reg-warn">⚠ Device-bound: This identity is tied to this device and browser. You cannot transfer it to another device. For permanent multi-device identity, use the Aequitas Android App instead.<br><br>⚠ <strong>Important:</strong> WebAuthn proves device possession — NOT biological uniqueness. A person with two devices could theoretically register twice. If uniqueness is critical to you, use the Android App with biometric verification instead.</div>
    <div class="rlog" id="rlog"><span class="info" data-i18n="reg-log-hint">// Open Aequitas Android App to generate your proof, then return here...</span></div>
  </div>
  <div class="ic">
    <div class="ic-title" data-i18n="reg-details">Registration Details</div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-network">Network</span><span class="ic-val p">Aequitas Chain (BlockDAG)</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-chainid">Chain ID</span><span class="ic-val b">1926 (0x786)</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-grant">UBI Grant</span><span class="ic-val go">1,000 AEQ per human</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-fee">Gas Fee</span><span class="ic-val g" data-i18n="free">FREE — completely gasless</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-limit">Registrations</span><span class="ic-val" data-i18n="k-limit-v">Once per human · permanent · immutable</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-bio">Biometric Data</span><span class="ic-val g" data-i18n="never-stored">Never stored — stays on your device</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-proof">Proof System</span><span class="ic-val p">Groth16 ZKP (Zero-Knowledge)</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-conf">Confirmation</span><span class="ic-val" data-i18n="k-conf-v">Within 6 seconds (1 block)</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="k-sybil">Sybil Protection</span><span class="ic-val g" data-i18n="k-sybil-v">One identity per biometric · permanent lock</span></div>
  </div>

  <!-- EXPLORE SECTION -->
  <div style="margin-top:20px;background:linear-gradient(135deg,rgba(139,92,246,0.07),rgba(6,182,212,0.03));border:1px solid rgba(139,92,246,0.2);border-radius:var(--radius);padding:20px">
    <div style="font-size:0.57rem;color:var(--purple);letter-spacing:2.5px;text-transform:uppercase;font-weight:700;margin-bottom:14px" data-i18n="explore-title">Explore Aequitas</div>
    <div style="display:grid;grid-template-columns:repeat(2,1fr);gap:8px">
      <div class="expl-card" onclick="goTab('index','eqi-score')">
        <div class="expl-icon">⚖️</div>
        <div class="expl-name" data-i18n="expl-score">Equality Score</div>
        <div class="expl-desc" data-i18n="expl-score-d">Live Gini coefficient · Aequitas Index · wealth distribution in real time</div>
      </div>
      <div class="expl-card" onclick="goTab('index','eqi-economy')">
        <div class="expl-icon">💸</div>
        <div class="expl-name" data-i18n="expl-economy">UBI &amp; Redistribution Pools</div>
        <div class="expl-desc" data-i18n="expl-economy-d">Daily UBI countdown · 4 on-chain pools · demurrage · Protocol Phases</div>
      </div>
      <div class="expl-card" onclick="goTab('index','eqi-lorenz')">
        <div class="expl-icon">📈</div>
        <div class="expl-name" data-i18n="expl-charts">Charts &amp; History</div>
        <div class="expl-desc" data-i18n="expl-charts-d">Gini history · Lorenz curve · Wealth Cap bootstrap slider · The story of Aequitas</div>
      </div>
      <div class="expl-card" onclick="goTab('network','net-protocol')">
        <div class="expl-icon">📜</div>
        <div class="expl-name" data-i18n="expl-v7">Protocol V7 Docs</div>
        <div class="expl-desc" data-i18n="expl-v7-d">AequitasV7 contract · 6 mechanisms · ZK proof · wealth cap · demurrage · immutable code</div>
      </div>
      <div class="expl-card" onclick="goTab('explorer','sep-blocks')">
        <div class="expl-icon">🔍</div>
        <div class="expl-name" data-i18n="expl-explorer">Block Explorer</div>
        <div class="expl-desc" data-i18n="expl-explorer-d">Live BlockDAG · click any block to see validator, hash, transactions, parent hashes</div>
      </div>
      <div class="expl-card" onclick="goTab('network','net-overview')">
        <div class="expl-icon">🌐</div>
        <div class="expl-name" data-i18n="expl-network">Network &amp; Nodes</div>
        <div class="expl-desc" data-i18n="expl-network-d">Node topology · run your own node · technical specs · Chain ID 1926</div>
      </div>
    </div>
  </div>
</div>
</div>

<!-- EXPLORER + HUMANS -->
<div id="tab-explorer" class="tab-content">
<nav class="stabs">
  <div class="stab active" onclick="showStab('tab-explorer','sep-blocks',this)">📦 Blocks</div>
  <div class="stab" onclick="showStab('tab-explorer','sep-humans',this)">👥 Humans</div>
</nav>
<div id="sep-blocks" class="stab-panel active">
<div class="hero">
  <div class="section-label" data-i18n="live-stats">Live Chain Statistics</div>
  <div class="stats-grid">
    <div class="stat c-blue"><div class="stat-accent"></div><div class="stat-icon">🔗</div><div class="stat-lbl" data-i18n="s-height">Block Height</div><div class="stat-val" id="s-height">—</div><div class="stat-sub" data-i18n="s-height-sub">New block every ~6s · BlockDAG · Parallel production</div></div>
    <div class="stat c-green"><div class="stat-accent"></div><div class="stat-icon">🧬</div><div class="stat-lbl" data-i18n="s-humans">Verified Humans</div><div class="stat-val" id="s-humans">—</div><div class="stat-sub" data-i18n="s-humans-sub">Biometric ZKP · One person, one wallet, forever</div></div>
    <div class="stat c-gold"><div class="stat-accent"></div><div class="stat-icon">🪙</div><div class="stat-lbl" data-i18n="s-supply">Total Supply</div><div class="stat-val" id="s-supply">—</div><div class="stat-sub" data-i18n="s-supply-sub">Always = Humans × 1,000 AEQ</div></div>
    <div class="stat c-purple"><div class="stat-accent"></div><div class="stat-icon">⚖</div><div class="stat-lbl" data-i18n="s-index">Aequitas Index</div><div class="stat-val" id="s-index">—</div><div class="stat-sub" data-i18n="s-index-sub">0 = perfect equality · 100 = max inequality</div></div>
    <div class="stat c-teal"><div class="stat-accent"></div><div class="stat-icon">⚡</div><div class="stat-lbl" data-i18n="s-uptime">Uptime</div><div class="stat-val" id="s-uptime" style="font-size:1rem">—</div><div class="stat-sub" data-i18n="s-uptime-sub">Node v0.3.0 · Railway + Render · PostgreSQL</div></div>
  </div>
  <div class="info-banner">
    <div>
      <div class="ib-icon">🧬</div>
      <div class="ib-title" data-i18n="ib-poh">Proof of Humanity</div>
      <div class="ib-text" data-i18n="ib-poh-t">Every AEQ holder must cryptographically prove they are a unique living human. No bots, no corporations, no AI, no duplicates. Biometric data never leaves your device — only a mathematical proof of uniqueness is transmitted. This means AEQ is the first currency that is exclusively human.</div>
    </div>
    <div>
      <div class="ib-icon">⚖</div>
      <div class="ib-title" data-i18n="ib-fair">Radically Fair Distribution</div>
      <div class="ib-text" data-i18n="ib-fair-t">Every verified human receives exactly 1,000 AEQ upon registration — no more, no less. No pre-mine, no founder allocation, no investor rounds. The total supply always and exactly equals the number of verified humans multiplied by 1,000. This is enforced mathematically, not by policy.</div>
    </div>
    <div>
      <div class="ib-icon">🔗</div>
      <div class="ib-title" data-i18n="ib-dag">BlockDAG Architecture</div>
      <div class="ib-text" data-i18n="ib-dag-t">Unlike traditional blockchains where only one block can exist per height, Aequitas uses a Directed Acyclic Graph (DAG) structure. Multiple blocks can be produced simultaneously by different nodes and later merged into the DAG. This enables higher throughput, lower latency, and eliminates single-node bottlenecks. Merge events are marked with a special badge in the explorer below.</div>
    </div>
    <div>
      <div class="ib-icon">⛽</div>
      <div class="ib-title" data-i18n="ib-gas">Truly Gasless</div>
      <div class="ib-text" data-i18n="ib-gas-t">All registrations and AEQ transfers cost absolutely nothing. No ETH, BNB, or MATIC required. No credit card, no bank account, no prior cryptocurrency needed. The relayer covers all transaction costs on behalf of users. If you are a human with a smartphone, you can participate — regardless of your economic situation.</div>
    </div>
  </div>
</div>
<div class="main-grid">
  <div class="section">
    <div class="sec-head"><div class="sec-title"><span class="sec-dot"></span><span data-i18n="recent-blocks">Recent Blocks</span></div><div class="sec-count" id="block-count">—</div></div>
    <div class="sec-desc" data-i18n="blocks-desc">Each row represents one block in the Aequitas BlockDAG. MERGE = this block has multiple parents, meaning two blocks were produced in parallel and later merged — the core feature of BlockDAG. TX = this block contains a human registration transaction. Block time averages ~6 seconds.</div>
    <div id="blocks-list"><div class="empty" data-i18n="loading">Loading blocks...</div></div>
  </div>
  <!-- Block detail overlay -->
  <div class="block-detail-overlay" id="block-detail-overlay" onclick="if(event.target===this)closeBlock()">
    <div class="bdc">
      <div class="bdc-hdr">
        <div style="font-size:0.75rem;font-weight:700;color:var(--purple);font-family:var(--font-mono)" id="bdc-title">Block #—</div>
        <div class="bdc-close" onclick="closeBlock()">✕ Close</div>
      </div>
      <div id="bdc-content"></div>
    </div>
  </div>
  <div class="right-col">
    <div class="ic">
      <div class="ic-title" data-i18n="net-info">Network Info</div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-chain">Chain Name</span><span class="ic-val go">Aequitas Chain</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-chainid">Chain ID</span><span class="ic-val b">1926 (0x786)</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-symbol">Symbol</span><span class="ic-val go">AEQ</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-btime">Block Time</span><span class="ic-val">~6 seconds</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-cons">Consensus</span><span class="ic-val p">BlockDAG + PoH</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-nodes">Active Nodes</span><span class="ic-val g">2 Online</span></div>
      <div class="ic-row"><span class="ic-key">ZKP System</span><span class="ic-val p">Groth16 / BN128</span></div>
      <div class="ic-row"><span class="ic-key">EVM Compatible</span><span class="ic-val g">Yes (Chain ID 1926)</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-storage">Storage</span><span class="ic-val g">PostgreSQL</span></div>
      <div class="ic-row"><span class="ic-key">P2P</span><span class="ic-val">libp2p (Go)</span></div>
    </div>
    <div class="mm-card">
      <div class="mm-title" data-i18n="add-mm">ADD TO METAMASK</div>
      <div style="font-size:0.6rem;color:var(--muted);margin-bottom:10px;line-height:1.7">Add Aequitas Chain to MetaMask to view your AEQ balance and interact with the network directly from your browser or mobile wallet.</div>
      <div class="mm-row"><span class="mm-key" data-i18n="k-chain">Network Name</span><span class="mm-val">Aequitas Chain</span></div>
      <div class="mm-row"><span class="mm-key">RPC URL</span><span class="mm-val" style="font-size:0.5rem">aequitas.digital/rpc</span></div>
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
<div id="sep-humans" class="stab-panel">
<div class="hero">
  <div class="section-label" data-i18n="humans-title">Verified Humans on Aequitas Chain</div>
  <div class="info-banner">
    <div>
      <div class="ib-icon">🔒</div>
      <div class="ib-title" data-i18n="h-what">What is a Verified Human?</div>
      <div class="ib-text" data-i18n="h-what-t">A Verified Human is a wallet address cryptographically proven to belong to a unique living human being. Verification uses a physical hardware kit: the R503 optical sensor scans all 10 fingerprints; the MAX30102 PPG sensor detects a live pulse. Phase 2 adds hand vein IR imaging; Phase 3 adds an iris scan — the gold standard with 1-in-10⁷⁸ uniqueness. Only a Groth16 Zero-Knowledge Proof is ever transmitted. Raw biometric data never leaves the device, never touches a server, and is never stored anywhere.</div>
    </div>
    <div>
      <div class="ib-icon">🧮</div>
      <div class="ib-title" data-i18n="h-zkp">Zero-Knowledge Proof System</div>
      <div class="ib-text" data-i18n="h-zkp-t">Aequitas uses Groth16 proofs on the BN128 elliptic curve — the same curve used by Ethereum and Zcash. Proof size: ~200 bytes. Verification: ~10ms. The multi-factor commitment binds all biometric inputs into a single proof: <code>commitment = keccak256(iris_hash ‖ vein_hash ‖ fingers_hash ‖ wallet)</code>. The nullifier is bound exclusively to physical body features — not the device — so a person who loses their phone can re-verify with the same iris and hand veins without creating a duplicate identity. No biometric data is ever revealed or stored.</div>
    </div>
    <div>
      <div class="ib-icon">🛡</div>
      <div class="ib-title" data-i18n="h-sybil">Sybil Attack Prevention</div>
      <div class="ib-text" data-i18n="h-sybil-t">A Sybil attack is when one person creates multiple identities to gain an unfair advantage. Aequitas prevents this with a permanent nullifier stored on-chain — derived from physical biometrics, not from a device. <strong>Phase 1</strong> uses all 10 fingerprints + MAX30102 liveness (PPG pulse). <strong>Phase 2</strong> adds hand vein IR pattern (1 in 10⁷, uncopyable — an internal body feature). <strong>Phase 3</strong> adds iris scan (1 in 10⁷⁸ — the global gold standard). From Phase 3, the nullifier is <code>keccak256(iris_hash ‖ vein_hash ‖ domain)</code> — bound to the body, not the phone. Losing a device does not create a second identity. One human, one wallet, forever.</div>
    </div>
    <div>
      <div class="ib-icon">🌍</div>
      <div class="ib-title" data-i18n="h-global">Global Financial Inclusion</div>
      <div class="ib-text" data-i18n="h-global-t">1.4 billion adults worldwide have no bank account. Aequitas requires nothing more than an Android smartphone with a fingerprint or face sensor — a device over 3 billion people already own. No bank account, no credit card, no prior cryptocurrency, no government ID. Just being human is enough to participate in the Aequitas economy.</div>
    </div>
    <div>
      <div class="ib-icon">🫁</div>
      <div class="ib-title" data-i18n="h-bio-hw">Biometric Hardware Roadmap</div>
      <div class="ib-text" data-i18n="h-bio-hw-t">
        <strong style="color:var(--neon)">Phase 1 (active)</strong> — <strong>R503</strong> optical fingerprint scanner (GROW, UART): all 10 fingers, combined hash. <strong>MAX30102</strong> PPG sensor: confirms live pulse, rejects replays and casts.<br><br>
        <strong style="color:var(--blue)">Phase 2 (planned)</strong> — <strong>ESP32-CAM + 850 nm IR LED</strong>: infrared hand vein imaging. Veins are an internal body feature — impossible to copy, immutable over a lifetime, unique 1 in 10⁷. Different even in identical twins.<br><br>
        <strong style="color:var(--gold)">Phase 3 (planned)</strong> — <strong>IR Iris module</strong>: 240+ degrees of freedom, uniqueness 1 in 10⁷⁸. The global gold standard. Different in identical twins. Immutable from birth. From Phase 3, the Sybil nullifier is bound exclusively to iris + vein — phone loss cannot create a duplicate identity.
      </div>
    </div>
  </div>
</div>
<div class="hs">
  <div class="section">
    <div class="sec-head"><div class="sec-title"><span class="sec-dot"></span><span data-i18n="reg-humans">Registered Humans</span></div><div class="sec-count" id="h-count">0</div></div>
    <div class="sec-desc" data-i18n="h-desc">Every address below has been verified as a unique human through biometric Zero-Knowledge Proof. Each received exactly 1,000 AEQ upon registration. The registry is permanent, immutable, and on-chain — no entry can ever be deleted or modified.</div>
    <div id="humans-list"><div class="empty" data-i18n="no-humans">No humans registered yet. Download the Aequitas Android App and be the first!</div></div>
  </div>
  <div class="right-col">
    <div class="ic">
      <div class="ic-title" data-i18n="reg-stats">Registry Stats</div>
      <div class="ic-row"><span class="ic-key" data-i18n="total-humans">Total Humans</span><span class="ic-val g" id="stat-humans">0</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="s-supply">Total Supply</span><span class="ic-val go" id="stat-supply">0 AEQ</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-grant">Grant per Human</span><span class="ic-val go">1,000 AEQ</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-fee">Registration Fee</span><span class="ic-val g" data-i18n="free">FREE — gasless</span></div>
      <div class="ic-row"><span class="ic-key">ZKP System</span><span class="ic-val p">Groth16 / BN128</span></div>
      <div class="ic-row"><span class="ic-key">Hash System</span><span class="ic-val p">keccak256</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="k-bio">Biometric Storage</span><span class="ic-val g" data-i18n="never-stored">Never stored</span></div>
      <div class="ic-row"><span class="ic-key">Sybil Protection</span><span class="ic-val g">Permanent · On-chain</span></div>
    </div>
    <div class="ic" style="margin-top:12px">
      <div class="ic-title">❓ FAQ</div>
      <div style="font-size:0.63rem;color:var(--text);padding:8px 0;border-bottom:1px solid var(--border);font-weight:600">Is my biometric data safe?</div>
      <div style="font-size:0.62rem;color:var(--muted);padding:6px 0 10px;border-bottom:1px solid var(--border);line-height:1.7">Yes. Your fingerprint or face scan never leaves your device. The Hardware Secure Element processes the biometric and produces a cryptographic key. Only a mathematical proof derived from that key is ever transmitted.</div>
      <div style="font-size:0.63rem;color:var(--text);padding:8px 0;border-bottom:1px solid var(--border);font-weight:600">Can I register with a different wallet later?</div>
      <div style="font-size:0.62rem;color:var(--muted);padding:6px 0 10px;border-bottom:1px solid var(--border);line-height:1.7">No. Registration is permanently bound to one wallet address per biometric identity. This is by design — it prevents Sybil attacks and ensures the one-person-one-wallet guarantee.</div>
      <div style="font-size:0.63rem;color:var(--text);padding:8px 0;border-bottom:1px solid var(--border);font-weight:600">What happens if I lose my phone?</div>
      <div style="font-size:0.62rem;color:var(--muted);padding:6px 0 10px;line-height:1.7">Your AEQ remains in your wallet — it is tied to your private key, not your phone. You can still access your wallet via MetaMask with your seed phrase. Wallet recovery is independent of the biometric registration.</div>
    </div>
  </div>
</div>
</div>
</div>

<!-- EXCHANGE -->
<div id="tab-exchange" class="tab-content">
<nav class="stabs">
  <div class="stab active" onclick="showStab('tab-exchange','exch-swap',this);setTimeout(drawPriceChart,50)">🔄 Swap</div>
  <div class="stab" onclick="showStab('tab-exchange','exch-liquidity',this)">💧 Liquidity</div>
</nav>
<div id="exch-swap" class="stab-panel active">
<div style="padding:16px 20px 0">
  <div class="idx">
    <div class="idx-title">AEQ / tUSD — Live Price</div>
    <div style="font-size:0.63rem;color:var(--muted);margin-bottom:12px">Real-time price derived from pool reserves (x·y=k). Updates every 8 seconds as new pool data arrives.</div>
    <div style="display:flex;gap:4px;margin-bottom:6px">
      <button onclick="setChartInterval(60000)" id="ci-1m" class="ci-btn ci-active">1m</button>
      <button onclick="setChartInterval(300000)" id="ci-5m" class="ci-btn">5m</button>
      <button onclick="setChartInterval(1800000)" id="ci-30m" class="ci-btn">30m</button>
      <button onclick="setChartInterval(3600000)" id="ci-1h" class="ci-btn">1h</button>
      <button onclick="setChartInterval(14400000)" id="ci-4h" class="ci-btn">4h</button>
      <button onclick="setChartInterval(0)" id="ci-all" class="ci-btn">All</button>
    </div>
    <canvas id="price-chart" height="160" style="width:100%;border-radius:6px;background:var(--card2)"></canvas>
    <div id="price-chart-empty" style="display:none;text-align:center;padding:24px;color:var(--muted);font-size:0.63rem">No pool data yet — add liquidity to see the price chart.</div>
  </div>
</div>
<div class="rs">
  <div class="rhero">
    <div class="rhero-title" data-i18n="swap-title">🔄 Swap AEQ ↔ tUSD</div>
    <div class="rhero-sub" data-i18n="swap-sub">Exchange AEQ for tUSD (a simulated test-dollar) through the native liquidity pool. A 0.1% fee applies only to swaps — ordinary AEQ transfers between people remain completely free.</div>
  </div>
  <div class="pbar" data-i18n="swap-priv-bar">🔒 0.1% swap fee only · AEQ-to-AEQ transfers stay free · tUSD is a test currency with no real-world value</div>
  <div class="rcard">
    <div class="wbox" id="swap-wbox"><div class="wlbl" data-i18n="conn-wallet">CONNECTED WALLET</div><div class="wadr" id="swap-wadr" title="">—</div><button onclick="copyAddr('swap-wadr',this)" style="margin-top:4px;padding:3px 10px;font-size:0.56rem;background:rgba(0,255,209,0.08);border:1px solid rgba(0,255,209,0.2);color:var(--neon);border-radius:4px;cursor:pointer">📋 Copy</button></div>
    <div id="demurrage-notice" style="display:none;font-size:13px;padding:10px 12px;border-radius:8px;background:rgba(255,179,0,0.1);border:1px solid rgba(255,179,0,0.3);color:var(--gold);margin:10px 0"></div>
    <div class="ic-row" style="margin:8px 0"><span class="ic-key" data-i18n="swap-your-aeq">Your AEQ</span><span class="ic-val go" id="swap-bal-aeq">—</span></div>
    <div class="ic-row" style="margin-bottom:16px"><span class="ic-key" data-i18n="swap-your-tusd">Your tUSD</span><span class="ic-val go" id="swap-bal-tusd">—</span></div>

    <!-- DEX-style Sell panel -->
    <div style="background:var(--card2);border:1px solid var(--border);border-radius:10px;padding:14px;margin-bottom:2px">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:8px">
        <div style="font-size:0.54rem;color:var(--muted);text-transform:uppercase;letter-spacing:1px" data-i18n="swap-sell-label">Sell</div>
        <div style="font-size:0.58rem;color:var(--muted)">Bal: <span id="swap-from-bal" style="color:var(--neon)">—</span></div>
      </div>
      <div style="display:flex;gap:8px;align-items:center">
        <button id="swap-from-btn" onclick="reverseSwapDir()" style="display:flex;align-items:center;gap:5px;background:rgba(139,92,246,0.12);border:1px solid rgba(139,92,246,0.3);border-radius:8px;padding:8px 12px;cursor:pointer;min-width:86px;font-size:0.68rem;font-weight:700;color:var(--text);white-space:nowrap">
          <span id="swap-from-icon">🔶</span><span id="swap-from-sym">AEQ</span><span style="color:var(--muted);font-size:0.55rem;margin-left:auto">⇄</span>
        </button>
        <input type="number" id="swap-amount" placeholder="0.00" oninput="updateFeeEstimate()" style="flex:1;padding:12px;border-radius:8px;border:1px solid var(--border);background:#0A1220;color:#E8EDF5;font-size:16px;min-width:0;box-sizing:border-box">
      </div>
      <div style="display:flex;gap:5px;margin-top:8px">
        <button class="rbtn pctbtn" onclick="setSwapPct(0.25)" style="flex:1;padding:6px;font-size:11px">25%</button>
        <button class="rbtn pctbtn" onclick="setSwapPct(0.5)" style="flex:1;padding:6px;font-size:11px">50%</button>
        <button class="rbtn pctbtn" onclick="setSwapPct(0.75)" style="flex:1;padding:6px;font-size:11px">75%</button>
        <button class="rbtn pctbtn" onclick="setSwapPct(1)" style="flex:1;padding:6px;font-size:11px">MAX</button>
      </div>
    </div>
    <!-- Reverse direction arrow -->
    <div style="display:flex;justify-content:center;margin:4px 0">
      <button onclick="reverseSwapDir()" style="background:var(--card2);border:1px solid var(--border);border-radius:50%;width:32px;height:32px;display:flex;align-items:center;justify-content:center;cursor:pointer;font-size:1rem;color:var(--muted)" title="Reverse direction">⇅</button>
    </div>
    <!-- DEX-style Receive panel -->
    <div style="background:var(--card2);border:1px solid var(--border);border-radius:10px;padding:14px;margin-bottom:8px">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:8px">
        <div style="font-size:0.54rem;color:var(--muted);text-transform:uppercase;letter-spacing:1px" data-i18n="swap-receive-label">Receive</div>
        <div style="font-size:0.58rem;color:var(--muted)">Bal: <span id="swap-to-bal" style="color:var(--neon)">—</span></div>
      </div>
      <div style="display:flex;gap:8px;align-items:center">
        <div style="display:flex;align-items:center;gap:5px;background:rgba(6,182,212,0.08);border:1px solid rgba(6,182,212,0.2);border-radius:8px;padding:8px 12px;min-width:86px;font-size:0.68rem;font-weight:700;color:var(--text)">
          <span id="swap-to-icon">💵</span><span id="swap-to-sym">tUSD</span>
        </div>
        <div id="swap-out-est-dex" style="flex:1;padding:12px;border-radius:8px;border:1px solid rgba(255,255,255,0.05);background:rgba(0,0,0,0.15);color:var(--neon);font-size:16px;font-family:monospace;min-width:0">—</div>
      </div>
    </div>
    <div id="swap-details-panel" class="sd-panel" style="display:none">
      <div class="sd-header" data-i18n="swap-details-hdr">Swap Details</div>
      <div class="sd-row"><span class="sd-key" data-i18n="swap-out-lbl">You receive (est.)</span><span class="sd-val" id="swap-out-est" style="color:var(--neon)">—</span></div>
      <div class="sd-row"><span class="sd-key" data-i18n="swap-impact-lbl">Price impact</span><span class="sd-val" id="swap-price-impact">—</span></div>
      <div class="sd-row"><span class="sd-key" data-i18n="swap-fee-est">Protocol fee (0.1%)</span><span class="sd-val" id="swap-fee-est" style="color:var(--muted)">—</span></div>
      <div class="sd-row"><span class="sd-key" data-i18n="swap-rate-lbl">Exchange rate</span><span class="sd-val" id="swap-rate-display" style="color:var(--purple)">—</span></div>
    </div>
    <div id="swap-warn" style="display:none;font-size:13px;padding:10px 12px;border-radius:8px;background:rgba(255,179,0,0.1);border:1px solid rgba(255,179,0,0.3);color:var(--gold);margin-bottom:10px"></div>

    <button class="rbtn bc" id="swap-btn-conn" onclick="connectSwapWallet()" data-i18n="btn-conn">🦊 CONNECT METAMASK</button>
    <button id="swap-btn-disconnect" onclick="disconnectWallet()" style="display:none;margin-top:6px;padding:8px 16px;font-size:0.6rem;letter-spacing:1px;border:1px solid rgba(248,113,113,0.4);background:rgba(248,113,113,0.08);color:var(--red);border-radius:6px;cursor:pointer;width:100%">⊘ DISCONNECT WALLET</button>
    <button class="rbtn br" id="swap-btn-go" onclick="doSwap()" disabled data-i18n="swap-btn-go">🔄 SWAP</button>
    <div class="rlog" id="swap-log"><span class="info" data-i18n="swap-log-hint">// Connect your wallet to swap...</span></div>

    <div class="ic" style="margin-top:20px">
      <div class="ic-title" data-i18n="swap-no-liquidity">No tUSD yet?</div>
      <div class="ic-row"><span class="ic-key" data-i18n="swap-faucet-desc">Registered humans can claim test-tUSD once</span></div>
      <button class="rbtn" id="swap-btn-faucet" onclick="claimFaucet()" disabled data-i18n="swap-btn-faucet" style="margin-top:8px">💧 CLAIM TEST-tUSD</button>
    </div>
</div>
</div>
</div>
<div id="exch-liquidity" class="stab-panel">
<div class="rs">
  <div class="rhero">
    <div class="rhero-title">💧 Liquidity</div>
    <div class="rhero-sub">Provide AEQ / tUSD liquidity to earn 30% of all swap fees, distributed daily.</div>
  </div>

<div class="ic">
    <div class="ic-title" data-i18n="swap-pool-title">AEQ / tUSD — Pool Status</div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-pool-price">Spot Price</span><span class="ic-val go" id="pool-price">—</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-pool-aeq">AEQ Reserve</span><span class="ic-val p" id="pool-reserve-aeq">—</span></div>
    <div class="ic-row" style="margin-bottom:4px"><span class="ic-key" data-i18n="swap-pool-tusd">tUSD Reserve</span><span class="ic-val b" id="pool-reserve-tusd">—</span></div>
    <div style="margin:12px 0 4px">
      <div style="font-size:0.54rem;color:var(--muted);margin-bottom:6px;font-weight:600;letter-spacing:1.5px;text-transform:uppercase" data-i18n="swap-depth-lbl">Pool Composition</div>
      <div class="depth-track">
        <div id="depth-aeq-fill" class="depth-aeq-fill" style="width:50%"></div>
        <div class="depth-tusd-fill"></div>
      </div>
      <div class="depth-lbls">
        <span style="color:var(--purple)">AEQ <span id="depth-aeq-pct">50%</span></span>
        <span style="color:var(--teal)"><span id="depth-tusd-pct">50%</span> tUSD</span>
      </div>
    </div>
    <div class="ic-row" style="padding-top:4px"><span class="ic-key" data-i18n="swap-fee-bps">Swap Fee</span><span class="ic-val g">0.1% · split 40/30/20/10</span></div>
    <div class="amm-box">
      <div style="font-size:0.54rem;color:var(--purple);font-weight:700;letter-spacing:1.2px;text-transform:uppercase;margin-bottom:6px" data-i18n="amm-title">x × y = k — Constant Product AMM</div>
      <div class="amm-formula">AEQ_reserve × tUSD_reserve = k (constant)</div>
      <div class="amm-text" data-i18n="amm-text">When you swap AEQ for tUSD, AEQ reserve grows and tUSD reserve shrinks — their product always stays equal to k. Every swap moves the price. Larger swaps relative to pool size cause greater price impact. The 0.1% fee is taken from the input before the formula is applied, ensuring the pool earns on every trade.</div>
    </div>
  </div>
  <div class="ic">
    <div class="ic-title" data-i18n="swap-pools-addr-title">Tokenomics Pool Addresses</div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-validators">Validators (40%)</span><span class="ic-val p" style="font-size:11px">0x78c1...d2bA</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-lps">Liquidity Providers (30%)</span><span class="ic-val p" style="font-size:11px">0xc181...01EB</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-ubi">UBI Pool (20%)</span><span class="ic-val p" style="font-size:11px">0x4A9b...054A</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-treasury">Treasury (10%)</span><span class="ic-val p" style="font-size:11px">0x2273...3eb15</span></div>
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
      <div id="lp-remove-inline" style="background:rgba(0,255,209,0.07);border:1px solid rgba(0,255,209,0.2);border-radius:6px;padding:8px 12px;margin-bottom:8px;font-size:0.66rem;color:var(--teal);text-align:center;display:none">≈ <span id="lp-inline-aeq">—</span> AEQ + <span id="lp-inline-tusd">—</span> tUSD</div>
      <div class="ic-row" style="margin-bottom:8px"><span class="ic-key" data-i18n="swap-lp-youget">You will receive</span><span class="ic-val go" id="lp-remove-preview">—</span></div>
      <button class="rbtn br" id="swap-btn-removeliq" onclick="doRemoveLiquidity()" data-i18n="swap-btn-removeliq">🔥 REMOVE LIQUIDITY</button>
    </div>
  </div>
</div>
</div>

<!-- INDEX (Equality) -->
<div id="tab-index" class="tab-content">
<nav class="stabs">
  <div class="stab active" onclick="showStab('tab-index','eqi-score',this)">📊 Score</div>
  <div class="stab" onclick="showStab('tab-index','eqi-lorenz',this);setTimeout(drawLorenzCurve,60)">📈 Distribution</div>
  <div class="stab" onclick="showStab('tab-index','eqi-economy',this)">💸 Economy</div>
  <div class="stab" onclick="showStab('tab-index','eqi-story',this)">📖 Story</div>
</nav>
<div id="eqi-score" class="stab-panel active">
<div class="is">
  <div class="idx" style="grid-column:1/-1">
    <div class="idx-title" data-i18n="idx-title">Aequitas Index — Real-Time Economic Equality Score</div>
    <div class="idx-desc" data-i18n="idx-desc">The Aequitas Index is derived from the <strong style="color:var(--teal)">Gini coefficient</strong> — the international standard for measuring wealth inequality, adopted by the World Bank, OECD, and UN. Unlike a simple richest-vs-poorest ratio, the Gini coefficient captures the <em style="color:var(--text)">entire distribution</em> across every verified human simultaneously, in a single number. <strong style="color:var(--neon)">0 = perfect equality</strong> (every wallet holds the same AEQ). <strong style="color:var(--red)">100 = maximum concentration</strong> (one wallet holds all AEQ in existence). For context: Bitcoin Gini ≈ 0.85 (Index 85) · most unequal country on Earth (South Africa) ≈ 0.63 · Scandinavia ≈ 0.27. Aequitas targets Gini below 0.30 at scale — comparable to the most equal developed economies — enforced automatically by the wealth cap and redistribution pools, no governance vote required.</div>
    <div style="display:grid;grid-template-columns:auto 1fr;gap:20px;align-items:center;margin-top:12px">
      <div><div class="idx-big" id="idx-score">—</div><div class="idx-lbl" data-i18n="curr-idx">Current Index</div></div>
      <div>
        <div class="bar-bg"><div class="bar-fill" id="idx-bar" style="width:0%"></div></div>
        <div class="bar-lbl"><span data-i18n="bar-0">0 — Perfect Equality</span><span>50</span><span data-i18n="bar-100">100 — Max Inequality</span></div>
        <div style="margin-top:8px;font-size:0.63rem;color:var(--muted);background:var(--card2);padding:8px 12px;border-radius:6px;border:1px solid var(--border)" id="idx-phase-desc">—</div>
      </div>
    </div>
    <div class="mrow" style="grid-template-columns:repeat(4,1fr)">
      <div class="mbox">
        <div class="mval" id="idx-gini">—</div>
        <div class="mlbl" data-i18n="gini">Gini Coefficient</div>
        <div style="font-size:0.55rem;color:var(--muted);margin-top:4px" data-i18n="gini-desc">0 = equal · 1 = unequal</div>
        <div id="gini-n-warn" style="display:none;font-size:0.5rem;color:var(--gold);margin-top:2px">⚠ N&lt;10: not yet significant</div>
      </div>
      <div class="mbox">
        <div class="mval" id="idx-supply2">—</div>
        <div class="mlbl" data-i18n="s-supply">Total Supply</div>
        <div style="font-size:0.55rem;color:var(--muted);margin-top:4px" data-i18n="supply-desc">Always = Humans × 1,000 AEQ</div>
      </div>
      <div class="mbox">
        <div class="mval" id="idx-phase">—</div>
        <div class="mlbl" data-i18n="phase">Protocol Phase</div>
        <div style="font-size:0.55rem;color:var(--muted);margin-top:4px" data-i18n="phase-desc">Auto-advances by human count</div>
      </div>
      <div class="mbox">
        <div class="mval" id="idx-humans2">—</div>
        <div class="mlbl" data-i18n="s-humans">Verified Humans</div>
        <div style="font-size:0.55rem;color:var(--muted);margin-top:4px" data-i18n="humans-desc">Biometrically verified unique humans</div>
      </div>
    </div>
    <div style="margin-top:20px;display:grid;grid-template-columns:1fr 1fr;gap:10px">
      <div style="background:var(--card2);border:1px solid rgba(6,182,212,0.2);border-radius:var(--radius-sm);padding:16px">
        <div style="font-size:0.6rem;color:var(--teal);letter-spacing:1.5px;text-transform:uppercase;margin-bottom:10px;font-weight:600" data-i18n="gini-what-title">What is the Gini Coefficient?</div>
        <div style="font-size:0.64rem;color:var(--muted);line-height:1.9" data-i18n="gini-what-text">Developed by Italian statistician Corrado Gini (1912). It measures wealth distribution by comparing the actual balance distribution against a hypothetical perfectly equal baseline — visualized as the Lorenz curve. The coefficient equals the ratio of the area between the Lorenz curve and the diagonal of equality to the total area below that diagonal. Scale: 0 means every person holds identical wealth. 1 means one person holds all wealth in existence. Used by the World Bank, OECD, and UN to compare countries. Reference values: Bitcoin ≈ 0.85 · South Africa (world record) ≈ 0.63 · Brazil ≈ 0.53 · USA ≈ 0.41 · Germany ≈ 0.31 · Scandinavia ≈ 0.27 · Aequitas long-term target: Gini below 0.30 — comparable to Scandinavian countries, enforced by wealth cap (bootstrap: 5×→25× per human).</div>
      </div>
      <div style="background:var(--card2);border:1px solid rgba(139,92,246,0.2);border-radius:var(--radius-sm);padding:16px">
        <div style="font-size:0.6rem;color:var(--purple);letter-spacing:1.5px;text-transform:uppercase;margin-bottom:10px;font-weight:600" data-i18n="gini-calc-title">How is the Aequitas Index calculated?</div>
        <div style="font-size:0.62rem;color:var(--purple);font-family:var(--font-mono);text-align:center;margin:8px 0;padding:10px;background:rgba(139,92,246,0.08);border-radius:6px;border:1px solid rgba(139,92,246,0.15)">G = Σ|xi − xj| / (2 × n² × x̄)<br><span style="color:var(--muted);font-size:0.58rem">Aequitas Index = G × 100</span></div>
        <div style="font-size:0.64rem;color:var(--muted);line-height:1.9;margin-top:8px" data-i18n="gini-calc-text">All AEQ balances of verified humans are collected (x₁ through xₙ). The formula computes the mean absolute difference between every possible pair of balances, normalized by the number of people squared (n²) and the mean balance (x̄). The result ranges 0–1 and is multiplied by 100 to produce the Aequitas Index. Updated on-chain after every registration, every monthly demurrage run, every pool payout, and every wealth cap enforcement event — via the keeper calling updateGini().</div>
      </div>
    </div>
    <div class="idx-grade-grid" style="margin-top:10px;display:grid;grid-template-columns:repeat(4,1fr);gap:8px">
      <div style="background:rgba(0,255,209,0.06);border:1px solid rgba(0,255,209,0.25);border-radius:var(--radius-sm);padding:14px;text-align:center">
        <div style="font-size:1.05rem;font-weight:700;color:var(--neon);font-family:var(--font-display)">0 – 35</div>
        <div style="font-size:0.6rem;color:var(--neon);margin-top:5px;font-weight:700;letter-spacing:0.5px">IDEAL</div>
        <div style="font-size:0.56rem;color:var(--muted);margin-top:5px;line-height:1.7">Healthier than most nations on Earth. Comparable to Scandinavia (0.27) and Germany (0.31). Wealth cap and demurrage successfully maintaining fair distribution.</div>
      </div>
      <div style="background:rgba(96,165,250,0.06);border:1px solid rgba(96,165,250,0.25);border-radius:var(--radius-sm);padding:14px;text-align:center">
        <div style="font-size:1.05rem;font-weight:700;color:var(--blue);font-family:var(--font-display)">35 – 50</div>
        <div style="font-size:0.6rem;color:var(--blue);margin-top:5px;font-weight:700;letter-spacing:0.5px">GOOD</div>
        <div style="font-size:0.56rem;color:var(--muted);margin-top:5px;line-height:1.7">Comparable to the USA (0.41) or France (0.32). Within the range of most developed economies. Redistribution mechanisms actively flattening the curve.</div>
      </div>
      <div style="background:rgba(245,166,35,0.06);border:1px solid rgba(245,166,35,0.25);border-radius:var(--radius-sm);padding:14px;text-align:center">
        <div style="font-size:1.05rem;font-weight:700;color:var(--gold);font-family:var(--font-display)">50 – 70</div>
        <div style="font-size:0.6rem;color:var(--gold);margin-top:5px;font-weight:700;letter-spacing:0.5px">WARNING</div>
        <div style="font-size:0.56rem;color:var(--muted);margin-top:5px;line-height:1.7">Higher than most European nations — comparable to Brazil (0.53) or Russia. Protocol redistribution at elevated intensity.</div>
      </div>
      <div style="background:rgba(248,113,113,0.06);border:1px solid rgba(248,113,113,0.25);border-radius:var(--radius-sm);padding:14px;text-align:center">
        <div style="font-size:1.05rem;font-weight:700;color:var(--red);font-family:var(--font-display)">70 – 100</div>
        <div style="font-size:0.6rem;color:var(--red);margin-top:5px;font-weight:700;letter-spacing:0.5px">CRITICAL</div>
        <div style="font-size:0.56rem;color:var(--muted);margin-top:5px;line-height:1.7">Worse than any country on Earth (South Africa record: 0.63). Approaching Bitcoin (0.85). Protocol at maximum intervention — wealth cap and redistribution at full force.</div>
      </div>
    </div>
    <div id="wealth-cap-info" style="margin-top:10px;background:var(--card2);border:1px solid rgba(0,255,209,0.2);border-radius:var(--radius-sm);padding:12px 16px;font-size:0.63rem;color:var(--muted);line-height:1.8">
      <span style="color:var(--neon);font-weight:700" data-i18n="wcap-lbl">Current Wealth Cap:</span>
      <span id="live-cap-aeq" style="color:var(--gold);font-weight:700;margin:0 6px">—</span>AEQ
      <span style="margin:0 8px;opacity:0.4">·</span>
      <span data-i18n="wcap-mult">Multiplier:</span>
      <span id="live-cap-mult" style="color:var(--teal);font-weight:700;margin-left:4px">—</span>
      <span style="margin:0 8px;opacity:0.4">·</span>
      <span data-i18n="wcap-avg">Avg balance:</span>
      <span id="live-cap-avg" style="color:var(--purple);font-weight:700;margin-left:4px">—</span> AEQ
    </div>
    <div style="margin-top:10px;background:rgba(245,166,35,0.04);border:1px solid rgba(245,166,35,0.15);border-radius:var(--radius-sm);padding:16px">
      <div style="font-size:0.6rem;color:var(--gold);letter-spacing:1.5px;text-transform:uppercase;margin-bottom:10px;font-weight:600" data-i18n="gini-why-title">Why the Gini coefficient — and not a simpler metric?</div>
      <div style="font-size:0.63rem;color:var(--muted);line-height:1.9" data-i18n="gini-why-text">A simple "richest vs. poorest" ratio is easy to game and misses what happens in the middle: a network could have 10,000 people, a low min/max spread, yet 90% of all AEQ concentrated in 100 wallets. The Gini coefficient detects this — a ratio does not. It captures the complete distribution across all verified humans in a single auditable number. Because Aequitas publishes this number on-chain (via updateGini), it is transparent, tamper-evident, and globally verifiable. The protocol uses it as the primary input signal for automatic phase transitions, wealth cap multiplier selection, and redistribution intensity — creating a self-correcting economic system governed entirely by mathematics. No human, no committee, no foundation can override the index reading or the mechanisms it triggers.</div>
    </div>
  </div>
  <div class="idx" style="grid-column:1/-1">
    <div class="idx-title">Gini Index History</div>
    <div style="font-size:0.63rem;color:var(--muted);margin-bottom:12px">Recorded after each UBI distribution. Shows how equality evolves as the network grows. Lower is better — target is Gini below 0.30.</div>
    <canvas id="gini-history-chart" height="160" style="width:100%;border-radius:6px;background:var(--card2)"></canvas>
    <div id="gini-history-empty" style="display:none;text-align:center;padding:24px;color:var(--muted);font-size:0.63rem">No snapshots yet — first one saved after the next UBI distribution.</div>
  </div>
</div>
</div>
<div id="eqi-lorenz" class="stab-panel">
<div style="padding:24px 20px;max-width:1100px;margin:0 auto">
  <div style="margin-bottom:20px">
    <div style="font-size:0.58rem;color:var(--purple);letter-spacing:2.5px;text-transform:uppercase;font-weight:700;margin-bottom:6px">Wealth Distribution Analysis</div>
    <div style="font-size:1.3rem;font-weight:800;color:var(--text);letter-spacing:-0.5px;margin-bottom:10px">Lorenz Curve — AEQ Distribution Across Humans</div>
    <div style="font-size:0.68rem;color:var(--muted);line-height:1.9;max-width:780px">
      The <strong style="color:var(--text)">Lorenz Curve</strong> visualizes how AEQ wealth is distributed among registered humans. The <span style="color:var(--purple)">diagonal line = perfect equality</span> — every human holds the same share. The further the gold curve bows <em>below</em> the diagonal, the more unequal the distribution. Reference curves show inequality in real countries. Aequitas targets a Gini coefficient <span style="color:var(--neon)">below 0.30</span> — comparable to Scandinavia.
    </div>
  </div>
  <div style="background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:20px;margin-bottom:16px">
    <canvas id="lorenz-chart" height="460" style="width:100%;border-radius:6px;background:#070B16;display:block"></canvas>
  </div>
  <div style="display:grid;grid-template-columns:repeat(3,1fr);gap:12px;margin-bottom:16px">
    <div style="background:var(--card2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:14px;text-align:center">
      <div style="font-size:0.57rem;color:var(--muted);text-transform:uppercase;letter-spacing:1px;margin-bottom:6px">Aequitas Now</div>
      <div id="lorenz-gini-val" style="font-size:1.4rem;font-weight:800;color:var(--neon);font-family:var(--font-mono)">—</div>
      <div style="font-size:0.58rem;color:var(--muted);margin-top:4px">Gini coefficient (0–1)</div>
    </div>
    <div style="background:var(--card2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:14px;text-align:center">
      <div style="font-size:0.57rem;color:var(--muted);text-transform:uppercase;letter-spacing:1px;margin-bottom:6px">Target</div>
      <div style="font-size:1.4rem;font-weight:800;color:var(--purple);font-family:var(--font-mono)">&lt; 0.30</div>
      <div style="font-size:0.58rem;color:var(--muted);margin-top:4px">Like Scandinavia (~0.27)</div>
    </div>
    <div style="background:var(--card2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:14px;text-align:center">
      <div style="font-size:0.57rem;color:var(--muted);text-transform:uppercase;letter-spacing:1px;margin-bottom:6px">Bitcoin Gini</div>
      <div style="font-size:1.4rem;font-weight:800;color:var(--red);font-family:var(--font-mono)">~0.85</div>
      <div style="font-size:0.58rem;color:var(--muted);margin-top:4px">Most unequal currency ever</div>
    </div>
  </div>
  <div style="background:var(--card2);border:1px solid rgba(155,114,246,0.2);border-radius:var(--radius-sm);padding:14px;font-size:0.62rem;color:var(--muted);line-height:1.85">
    <strong style="color:var(--purple)">How to read this chart:</strong> The X-axis shows the population from poorest (left) to richest (right). The Y-axis shows cumulative AEQ wealth. A point at (50%, 50%) = the poorest half of humans hold exactly half the AEQ. In perfect equality, the curve IS the diagonal. Aequitas enforces equality through automatic wealth cap, demurrage (0.5%/month decay), and daily UBI redistribution — keeping the curve close to the diagonal as the network grows.
  </div>
</div>
</div>
<div id="eqi-economy" class="stab-panel">
<div class="is">
<div class="idx" style="grid-column:1/-1">
    <div class="idx-title" data-i18n="pools-title">Redistribution Pools — Daily Economic Rebalancing</div>
    <div class="idx-desc" data-i18n="pools-desc">Every swap fee, demurrage charge, and wealth cap overflow flows automatically into four on-chain pools. No manual intervention, no admin key, no governance vote — the protocol distributes everything through code. Each pool pays out once per 24 hours.</div>

    <!-- UBI HERO SECTION -->
    <div class="ubi-hero-section">
      <div style="font-size:0.58rem;color:var(--gold);letter-spacing:3px;text-transform:uppercase;font-weight:700;margin-bottom:6px" data-i18n="ubi-hero-title">Universal Basic Income Pool</div>
      <div style="font-size:0.62rem;color:var(--muted);margin-bottom:10px" data-i18n="ubi-hero-sub">Accumulating — next payout distributed equally to all verified humans in:</div>
      <div id="ubi-timer" class="ubi-big-timer">—</div>
      <div style="font-size:0.6rem;color:var(--muted);margin-bottom:6px" data-i18n="ubi-bal-lbl">current pool balance</div>
      <div id="pool-u" class="ubi-pool-amount">0.0000 AEQ</div>
      <div class="ubi-fill-track"><div id="ubi-fill-bar" class="ubi-fill-bar"></div></div>
      <div style="font-size:0.61rem;color:var(--muted);line-height:1.85;margin-top:6px" data-i18n="ubi-hero-desc">Split equally among all verified humans · paid every 24 h · pool resets to zero after each payout · no minimum balance required to receive</div>
    </div>

    <!-- UBI SOURCE BREAKDOWN -->
    <div style="font-size:0.54rem;color:var(--muted);letter-spacing:2.5px;text-transform:uppercase;font-weight:600;margin:16px 0 8px" data-i18n="ubi-how-fills">How the UBI Pool fills up</div>
    <div class="ubi-src-grid">
      <div class="ubi-src-card" style="border-color:rgba(6,182,212,0.2)">
        <div class="ubi-src-pct" style="color:var(--teal)">20%</div>
        <div class="ubi-src-name" style="color:var(--teal)" data-i18n="ubi-src-swap">Swap Fees</div>
        <div class="ubi-src-desc" data-i18n="ubi-src-swap-d">Every AEQ↔tUSD swap contributes 20% of its 0.1% fee here. More trading activity = faster pool fill.</div>
      </div>
      <div class="ubi-src-card" style="border-color:rgba(245,166,35,0.2)">
        <div class="ubi-src-pct" style="color:var(--gold)">variable</div>
        <div class="ubi-src-name" style="color:var(--gold)" data-i18n="ubi-src-dem">Demurrage</div>
        <div class="ubi-src-desc" data-i18n="ubi-src-dem-d">Idle AEQ (3+ months inactive) decays at 0.5%/month. The decayed amount enters the 40/30/20/10 split — 20% goes to UBI.</div>
      </div>
      <div class="ubi-src-card" style="border-color:rgba(139,92,246,0.2)">
        <div class="ubi-src-pct" style="color:var(--purple)">variable</div>
        <div class="ubi-src-name" style="color:var(--purple)" data-i18n="ubi-src-cap">Wealth Cap Overflow</div>
        <div class="ubi-src-desc" data-i18n="ubi-src-cap-d">Wallets exceeding 25× average balance have the excess confiscated instantly. 20% flows to UBI immediately.</div>
      </div>
    </div>

    <!-- ALL FOUR POOLS GRID -->
    <div style="font-size:0.54rem;color:var(--muted);letter-spacing:2.5px;text-transform:uppercase;font-weight:600;margin:16px 0 10px" data-i18n="pools4-header">All four redistribution pools</div>
    <div class="pools4-grid">
      <div class="pool4-card" style="border-color:rgba(139,92,246,0.2)" onmouseover="this.style.borderColor='rgba(139,92,246,0.4)'" onmouseout="this.style.borderColor='rgba(139,92,246,0.2)'">
        <div class="pool4-head">
          <span class="pool4-name" style="color:var(--purple)" data-i18n="vel-pool">VALIDATORS</span>
          <span class="pool4-badge">40% of fees</span>
        </div>
        <div id="pool-v" class="pool4-amount" style="color:var(--purple)">0.0000 AEQ</div>
        <div class="pool4-timer" style="color:var(--purple)">⏰ Next: <span id="validators-timer">—</span></div>
        <div class="pool4-desc" data-i18n="vel-pool-desc">Node operators who produce blocks, validate ZK registrations, and secure the Aequitas BlockDAG. Paid daily, proportional to block production.</div>
      </div>
      <div class="pool4-card" style="border-color:rgba(6,182,212,0.2)" onmouseover="this.style.borderColor='rgba(6,182,212,0.4)'" onmouseout="this.style.borderColor='rgba(6,182,212,0.2)'">
        <div class="pool4-head">
          <span class="pool4-name" style="color:var(--teal)" data-i18n="liq-pool">LIQUIDITY PROVIDERS</span>
          <span class="pool4-badge">30% of fees</span>
        </div>
        <div id="pool-l" class="pool4-amount" style="color:var(--teal)">0.0000 AEQ</div>
        <div class="pool4-timer" style="color:var(--teal)">⏰ Next: <span id="lp-timer">—</span></div>
        <div class="pool4-desc" data-i18n="liq-pool-desc">Providers of AEQ/tUSD liquidity to the AMM pool receive 30% of all fees, proportional to their LP share. Deeper liquidity = lower price impact for every trader.</div>
      </div>
      <div class="pool4-card" style="border:1px solid rgba(245,166,35,0.3);background:linear-gradient(135deg,rgba(245,166,35,0.06),var(--card2))" onmouseover="this.style.borderColor='rgba(245,166,35,0.5)'" onmouseout="this.style.borderColor='rgba(245,166,35,0.3)'">
        <div class="pool4-head">
          <span class="pool4-name" style="color:var(--gold)" data-i18n="ubi-pool">UBI POOL</span>
          <span class="pool4-badge">20% of fees</span>
        </div>
        <div class="pool4-amount" style="color:var(--gold)" data-i18n="ubi-see-above">see countdown above</div>
        <div class="pool4-timer" style="color:var(--gold)" data-i18n="ubi-timer-above">⏰ countdown displayed above</div>
        <div class="pool4-desc" data-i18n="ubi-pool-desc">20% of swap fees + demurrage + wealth cap overflow → divided equally among all verified humans every 24 hours. Even with zero trading, demurrage and wealth cap ensure the pool always fills.</div>
      </div>
      <div class="pool4-card" style="border-color:rgba(96,165,250,0.2)" onmouseover="this.style.borderColor='rgba(96,165,250,0.4)'" onmouseout="this.style.borderColor='rgba(96,165,250,0.2)'">
        <div class="pool4-head">
          <span class="pool4-name" style="color:var(--blue)" data-i18n="treasury">TREASURY</span>
          <span class="pool4-badge">10% of fees</span>
        </div>
        <div id="pool-t" class="pool4-amount" style="color:var(--blue)">0.0000 AEQ</div>
        <div class="pool4-timer" style="color:var(--blue)" data-i18n="pool-t-timer">Accumulates — no timer</div>
        <div class="pool4-desc" data-i18n="treasury-desc">Protocol development, infrastructure, security audits, and future upgrades. Governed by the Aequitas team with full on-chain transparency.</div>
      </div>
    </div>
  </div>
  <div class="idx">
    <div class="idx-title" data-i18n="phases-title">Protocol Phases</div>
    <div class="idx-desc" data-i18n="phases-desc">The wealth cap uses a bootstrap multiplier during Phase 0: max(5, min(N, 25))× average balance. With 1–4 humans: 5× average. Each new human adds 1×. At 25+ humans: locks permanently at 25×. Phase 1+ maintains 25× fixed. All transitions trigger automatically by human count — no governance vote, no admin key required.</div>
    <table class="spect">
      <tr><td><strong style="color:var(--neon)">Phase 0</strong></td><td style="color:var(--neon)" data-i18n="p0">Bootstrap · &lt;100 humans · Wealth Cap: max(5,min(N,25))× average · Slides 5×→25× until 25th human · Currently active</td></tr>
      <tr><td><strong style="color:var(--blue)">Phase 1</strong></td><td style="color:var(--blue)" data-i18n="p1">Growth · 100–10,000 humans · Wealth Cap: 25× average balance</td></tr>
      <tr><td><strong style="color:var(--gold)">Phase 2</strong></td><td style="color:var(--gold)" data-i18n="p2">Stability · 10,000–1M humans · Wealth Cap: 25× average balance</td></tr>
      <tr><td><strong style="color:var(--purple)">Phase 3</strong></td><td style="color:var(--purple)" data-i18n="p3">Maturity · 1M+ humans · Wealth Cap: 25× average balance</td></tr>
    </table>
    <div class="hlbox" data-i18n="wealth-cap-explain">The <strong>Wealth Cap</strong> during Phase 0 (Bootstrap) uses the formula <strong>max(5, min(N, 25))× average AEQ balance</strong>, where N = registered humans. With 1–4 humans: cap = 5× average. Each new human adds 1×. At 25+ humans: the multiplier locks permanently at 25×. The cap always scales with the live average balance — automatically adjusting as the network grows.</div>
  </div>
  <div class="idx">
    <div class="idx-title" data-i18n="demurrage-title">Demurrage — Incentive to Circulate</div>
    <div class="idx-desc" data-i18n="demurrage-desc">Aequitas implements a demurrage mechanism inspired by historical complementary currencies like the Wörgl experiment (1932) and the Chiemgauer (2003). Idle AEQ balances slowly lose value to discourage hoarding and incentivize economic participation.</div>
    <table class="spect">
      <tr><td data-i18n="dem-rate-k">Decay Rate</td><td data-i18n="dem-rate-v">0.5% per month (continuous, not stepped)</td></tr>
      <tr><td data-i18n="dem-grace-k">Grace Period</td><td data-i18n="dem-grace-v">3 months of inactivity before decay begins</td></tr>
      <tr><td data-i18n="dem-reset-k">Clock Reset</td><td data-i18n="dem-reset-v">Any transfer, swap, or liquidity action resets the timer</td></tr>
      <tr><td data-i18n="dem-dest-k">Decayed AEQ goes to</td><td data-i18n="dem-dest-v">Redistribution pools (same 40/30/20/10 split)</td></tr>
      <tr><td data-i18n="dem-warn-k">Warning System</td><td data-i18n="dem-warn-v">14-day notice (once) + 7-day repeated reminder at login</td></tr>
    </table>
  </div>
  <div class="idx" style="grid-column:1/-1">
    <div class="idx-title">Wealth Cap Multiplier — Bootstrap Slider</div>
    <div style="font-size:0.63rem;color:var(--muted);margin-bottom:12px">Formula: <code style="color:var(--teal)">max(5, min(N, 25))×</code> average AEQ balance. Each new human slides the cap up by 1×, until the 25th human locks it at 25× permanently.</div>
    <canvas id="wcap-slide-chart" height="120" style="width:100%;border-radius:6px;background:var(--card2)"></canvas>
  </div>
</div>
</div>

<div id="eqi-story" class="stab-panel">
<div style="padding:24px;max-width:900px;margin:0 auto">

  <!-- ORIGIN STORY -->
  <div class="idx" style="margin-bottom:16px;background:linear-gradient(135deg,rgba(107,70,193,0.06),rgba(8,145,178,0.03));border-color:rgba(107,70,193,0.2)">
    <div class="idx-title" data-i18n="story-title">The Story of Aequitas &#8212; Why This Exists</div>
    <div class="story" style="font-size:0.7rem;line-height:2;color:var(--muted)">
      <p>The year is 2009. Satoshi Nakamoto releases Bitcoin. For the first time, value can transfer between any two people without a bank. A genuine revolution. But something goes wrong almost immediately.</p>
      <p>Early miners accumulate millions of coins at near-zero cost. By 2021, the top 1% of Bitcoin addresses control over 90% of all Bitcoin. Bitcoin&#8217;s Gini coefficient exceeds <strong style="color:var(--red)">0.85</strong> &#8212; higher than any country on Earth. The technology meant to democratize finance created the most extreme wealth concentration in history.</p>
      <div class="hlbox" style="margin:16px 0">
        <strong style="color:var(--purple)">Aequitas</strong> &#8212; Latin for <em>fairness</em> and <em>equity</em> &#8212; was built to answer one question:<br>
        <em style="color:var(--gold);font-size:0.85rem">"What would a cryptocurrency look like if designed from first principles to be fair to every human being?"</em>
      </div>
      <p>The answer: <strong style="color:var(--text)">Money exists because people exist. Therefore every person should have an equal share of money simply by virtue of being human.</strong></p>
      <p>Aequitas implements this mathematically. Every verified human receives exactly 1,000 AEQ &#8212; billionaire or subsistence farmer, no exceptions. Four redistribution mechanisms ensure inequality cannot accumulate indefinitely. The Gini coefficient is tracked on-chain in real time.</p>
      <p><em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em></p>
    </div>
  </div>

  <!-- HOW IT WORKS -->
  <div style="display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:16px">
    <div class="idx">
      <div class="idx-title">The Core Innovation</div>
      <div style="font-size:0.65rem;color:var(--muted);line-height:1.9">
        <div style="margin-bottom:10px"><span style="color:var(--purple);font-weight:700">ZK Biometric Proof</span><br>Your fingerprint is processed on your phone&#8217;s Hardware Secure Element. A Groth16 Zero-Knowledge Proof proves you are a unique human being &#8212; without any biometric data ever leaving your device or touching any server.</div>
        <div style="margin-bottom:10px"><span style="color:var(--teal);font-weight:700">No-Stake Blockchain</span><br>No mining. No staking. No proof-of-work. Block production is open to any node operator. Validators earn from the 40% fee pool &#8212; incentivized by fairness, not capital.</div>
        <div><span style="color:var(--neon);font-weight:700">One Human = One Wallet = 1,000 AEQ</span><br>Supply formula: Total AEQ = Verified Humans x 1,000. No pre-mine. No admin keys. No governance vote can change this.</div>
      </div>
    </div>
    <div class="idx">
      <div class="idx-title">The 4 Redistribution Mechanisms</div>
      <div style="font-size:0.65rem;color:var(--muted);line-height:1.9">
        <div style="padding:8px;background:rgba(4,120,87,0.06);border-radius:6px;border-left:3px solid var(--neon);margin-bottom:8px"><strong style="color:var(--neon)">UBI Pool (20%)</strong><br>Every 24 hours, the pool divides equally among all verified humans. Funded by swap fees + demurrage + wealth cap overflows.</div>
        <div style="padding:8px;background:rgba(107,70,193,0.06);border-radius:6px;border-left:3px solid var(--purple);margin-bottom:8px"><strong style="color:var(--purple)">Validators Pool (40%)</strong><br>Node operators earn from all protocol fees. More nodes = more decentralization.</div>
        <div style="padding:8px;background:rgba(8,145,178,0.06);border-radius:6px;border-left:3px solid var(--teal);margin-bottom:8px"><strong style="color:var(--teal)">Liquidity Pool (30%)</strong><br>Liquidity providers earn proportionally from all swap activity.</div>
        <div style="padding:8px;background:rgba(146,64,14,0.06);border-radius:6px;border-left:3px solid var(--gold)"><strong style="color:var(--gold)">Treasury (10%)</strong><br>Protocol development. No VC allocation. No founder bonus.</div>
      </div>
    </div>
  </div>

  <!-- PHASE ROADMAP -->
  <div class="idx" style="margin-bottom:16px">
    <div class="idx-title">Phase Roadmap &#8212; The Path to Global Scale</div>
    <div style="display:grid;grid-template-columns:repeat(4,1fr);gap:8px;margin-top:12px">
      <div style="background:rgba(107,70,193,0.08);border:2px solid var(--purple);border-radius:10px;padding:14px;position:relative">
        <div style="position:absolute;top:-1px;left:50%;transform:translateX(-50%);background:var(--purple);color:#fff;font-size:0.55rem;font-weight:700;padding:2px 8px;border-radius:0 0 6px 6px">ACTIVE NOW</div>
        <div style="font-size:0.72rem;font-weight:700;color:var(--purple);margin:10px 0 6px">Phase 0</div>
        <div style="font-size:0.62rem;font-weight:600;color:var(--text);margin-bottom:6px">Bootstrap</div>
        <div style="font-size:0.6rem;color:var(--muted);line-height:1.75">0 &#8211; 100 humans. Sliding wealth cap 5x &#8594; 25x. Foundation building.</div>
      </div>
      <div style="background:var(--card2);border:1px solid var(--border);border-radius:10px;padding:14px">
        <div style="font-size:0.72rem;font-weight:700;color:var(--teal);margin-bottom:6px">Phase 1</div>
        <div style="font-size:0.62rem;font-weight:600;color:var(--text);margin-bottom:6px">Growth</div>
        <div style="font-size:0.6rem;color:var(--muted);line-height:1.75">100 &#8211; 10,000 humans. Fixed cap 25x. Open node joining.</div>
      </div>
      <div style="background:var(--card2);border:1px solid var(--border);border-radius:10px;padding:14px">
        <div style="font-size:0.72rem;font-weight:700;color:var(--gold);margin-bottom:6px">Phase 2</div>
        <div style="font-size:0.62rem;font-weight:600;color:var(--text);margin-bottom:6px">Stability</div>
        <div style="font-size:0.6rem;color:var(--muted);line-height:1.75">10,000 &#8211; 1M humans. Min 10 nodes. Fully decentralized.</div>
      </div>
      <div style="background:var(--card2);border:1px solid var(--border);border-radius:10px;padding:14px">
        <div style="font-size:0.72rem;font-weight:700;color:var(--neon);margin-bottom:6px">Phase 3</div>
        <div style="font-size:0.62rem;font-weight:600;color:var(--text);margin-bottom:6px">Maturity</div>
        <div style="font-size:0.6rem;color:var(--muted);line-height:1.75">1M+ humans. Global UBI at scale. Gini target &lt;0.30.</div>
      </div>
    </div>
    <div style="margin-top:14px;padding:12px 16px;background:rgba(107,70,193,0.05);border-radius:8px;font-size:0.63rem;color:var(--muted);line-height:1.75;border:1px solid rgba(107,70,193,0.12)">Phase transitions are automatic &#8212; triggered by human count thresholds, enforced by the smart contract. No governance vote, no admin key.</div>
  </div>

  <!-- GUARDIAN SYSTEM -->
  <div class="idx" style="margin-bottom:16px;background:linear-gradient(135deg,rgba(8,145,178,0.05),rgba(107,70,193,0.03));border-color:rgba(8,145,178,0.2)">
    <div class="idx-title" style="color:var(--teal)">Guardian System &#8212; Human Failsafe for Lost Wallets</div>
    <div style="font-size:0.65rem;color:var(--muted);line-height:1.9;margin-bottom:14px">What happens when someone is hospitalized, incarcerated, or dies? In most crypto systems, lost wallets mean lost coins forever. Aequitas has a three-layer inactivity recovery system.</div>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:10px;margin-bottom:14px">
      <div>
        <div style="font-size:0.6rem;color:var(--teal);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:8px">What is a Guardian?</div>
        <div style="font-size:0.63rem;color:var(--muted);line-height:1.85">A Guardian is a <strong style="color:var(--text)">trusted verified human</strong> you designate. They have exactly one power: confirming you are still alive. They cannot move funds, transfer AEQ, or access your wallet under any circumstances. Maximum 3 wards per Guardian prevents centralization of trust.</div>
      </div>
      <div>
        <div style="font-size:0.6rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:8px">Inactivity Timeline</div>
        <div style="font-size:0.62rem;color:var(--muted);line-height:2">
          <div style="display:flex;gap:8px;align-items:center;margin-bottom:4px"><span style="min-width:80px;color:var(--text);font-weight:600">0 &#8211; 2 years</span>Normal usage, no restrictions</div>
          <div style="display:flex;gap:8px;align-items:center;margin-bottom:4px"><span style="min-width:80px;color:var(--gold);font-weight:600">Year 2</span>Warning 1 &#8212; Guardian can respond</div>
          <div style="display:flex;gap:8px;align-items:center;margin-bottom:4px"><span style="min-width:80px;color:var(--gold);font-weight:600">+60 days</span>Warning 2 &#8212; escalating urgency</div>
          <div style="display:flex;gap:8px;align-items:center;margin-bottom:4px"><span style="min-width:80px;color:var(--red);font-weight:600">+180 days</span>AEQ moved to escrow (recoverable)</div>
          <div style="display:flex;gap:8px;align-items:center"><span style="min-width:80px;color:var(--red);font-weight:600">Year 4</span>Escrow released to UBI Pool</div>
        </div>
      </div>
    </div>
    <div style="padding:12px 16px;background:rgba(8,145,178,0.06);border-radius:8px;border:1px solid rgba(8,145,178,0.15);font-size:0.62rem;color:var(--muted);line-height:1.8"><strong style="color:var(--teal)">Key protections:</strong> 7-day timelock on Guardian assignment. No circular Guardian relationships. Guardian assignment is public and on-chain.</div>
  </div>

  <!-- PROOF OF ALIVE -->
  <div class="idx" style="margin-bottom:16px">
    <div class="idx-title">Proof of Alive &#8212; Biometric Without Surveillance</div>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:10px;margin-top:12px">
      <div>
        <div style="font-size:0.6rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:8px">Why Biometrics?</div>
        <div style="font-size:0.63rem;color:var(--muted);line-height:1.85">Sybil resistance &#8212; preventing one person from registering multiple wallets &#8212; is the core unsolved problem of fair money distribution. Aequitas uses biometric Zero-Knowledge Proofs: mathematically proven uniqueness with <strong style="color:var(--text)">zero personal data exposed</strong>.</div>
      </div>
      <div>
        <div style="font-size:0.6rem;color:var(--teal);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:8px">How It Works</div>
        <div style="font-size:0.63rem;color:var(--muted);line-height:1.85">1. Fingerprint scanned by Hardware Secure Element (HSE)<br>2. HSE derives a device-bound cryptographic key<br>3. keccak256(bioKey) = bioHash (stored on-chain)<br>4. Groth16 ZK proof: proves uniqueness without revealing bioHash<br>5. Nullifier system: attempting to register twice is mathematically impossible</div>
      </div>
    </div>
    <div style="margin-top:12px;padding:12px 16px;background:rgba(4,120,87,0.06);border-radius:8px;border:1px solid rgba(4,120,87,0.15);font-size:0.63rem;color:var(--muted);line-height:1.8"><strong style="color:var(--neon)">Privacy guarantee:</strong> Your raw biometric data never leaves your device. The proof server receives only the ZK proof (200 bytes). The chain stores only a nullifier hash. No government, company, or attacker can reverse-engineer your identity. <em>Future: PPG sensor (MAX30102) for device-independent biometric verification.</em></div>
  </div>

  <!-- THE VISION -->
  <div class="idx" style="background:linear-gradient(135deg,rgba(146,64,14,0.06),rgba(107,70,193,0.04));border-color:rgba(146,64,14,0.2)">
    <div class="idx-title" style="color:var(--gold)">The Vision &#8212; A Global Basic Income Protocol</div>
    <div style="font-size:0.72rem;color:var(--muted);line-height:2;font-style:italic;text-align:center;padding:10px 0;font-family:var(--font-display)">"Imagine a world where every person on Earth &#8212; regardless of where they were born, what language they speak, or how much money their parents had &#8212; receives a guaranteed daily income simply for being human. Not as charity. As a mathematical right, enforced by code that no government or corporation can override."</div>
    <div style="display:grid;grid-template-columns:repeat(3,1fr);gap:8px;margin-top:12px">
      <div style="text-align:center;padding:14px;background:rgba(107,70,193,0.06);border-radius:8px;border:1px solid rgba(107,70,193,0.12)">
        <div style="font-size:1.4rem;font-weight:900;color:var(--purple);font-family:var(--font-display)">8B</div>
        <div style="font-size:0.58rem;color:var(--muted);margin-top:4px">humans could register</div>
      </div>
      <div style="text-align:center;padding:14px;background:rgba(4,120,87,0.06);border-radius:8px;border:1px solid rgba(4,120,87,0.12)">
        <div style="font-size:1.4rem;font-weight:900;color:var(--neon);font-family:var(--font-display)">&lt;0.30</div>
        <div style="font-size:0.58rem;color:var(--muted);margin-top:4px">Gini target (Scandinavian level)</div>
      </div>
      <div style="text-align:center;padding:14px;background:rgba(146,64,14,0.06);border-radius:8px;border:1px solid rgba(146,64,14,0.12)">
        <div style="font-size:1.4rem;font-weight:900;color:var(--gold);font-family:var(--font-display)">0</div>
        <div style="font-size:0.58rem;color:var(--muted);margin-top:4px">admin keys or governance votes</div>
      </div>
    </div>
  </div>

</div>
</div>

<!-- closes tab-index (one div was unclosed, tab-network was nested inside it) -->
</div>

<!-- NETWORK (merged) -->
<div id="tab-network" class="tab-content">
<nav class="stabs">
  <div class="stab active" onclick="showStab('tab-network','net-overview',this)">🌐 Overview</div>
  <div class="stab" onclick="showStab('tab-network','net-story',this)">📖 Story</div>
  <div class="stab" onclick="showStab('tab-network','net-runnode',this)">⚙️ Run a Node</div>
  <div class="stab" onclick="showStab('tab-network','net-protocol',this)">📜 Protocol V7</div>
</nav>
<div id="net-overview" class="stab-panel active">
<div class="ns">
<div class="nc" style="grid-column:1/-1">
    <div class="nc-title" data-i18n="nodes-title">Active Nodes — Current Network Topology</div>
    <div style="font-size:0.65rem;color:var(--muted);line-height:1.8;margin-bottom:12px" data-i18n="nodes-desc">The Aequitas network currently operates on two geographically distributed nodes. Node 1 runs on Railway (primary, aequitas.digital) and Node 2 runs on a Contabo VPS (173.249.37.118:8080). Both participate in block production, state synchronization, and API serving with active MERGE events. They communicate peer-to-peer via libp2p and synchronize block state via HTTP. The network is designed to support additional nodes — any registered human can run a validator node by setting the required environment variables and obtaining PEER_SECRET and SNAPSHOT_TOKEN from the network operator.</div>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:8px">
      <div class="nbox">
        <div class="nstat"><span class="ndot"></span><span data-i18n="node1">Node 1 — Railway (Primary)</span></div>
        <div class="nurl">aequitas-production-9fba.up.railway.app</div>
        <div class="ndesc" data-i18n="node1-desc">Primary API · Block producer · UBI distribution · P2P bootstrap · PostgreSQL · RPC for MetaMask</div>
        <div style="margin-top:6px;font-size:0.57rem;color:rgba(0,255,209,0.5)">DB-Lock (TryLockDistribution) · Daily pool distributions at 20:00 Berlin</div>
      </div>
      <div class="nbox">
        <div class="nstat"><span class="ndot"></span><span data-i18n="node2">Node 2 — Contabo VPS (Secondary)</span></div>
        <div class="nurl">173.249.37.118:8080</div>
        <div class="ndesc" data-i18n="node2-desc">Secondary API · Block producer · P2P peer · HTTP sync · Shared PostgreSQL state</div>
        <div style="margin-top:6px;font-size:0.57rem;color:rgba(139,92,246,0.5)">Redundancy · Geographic distribution · Contabo VPS</div>
      </div>
    </div>
  </div>
  <div class="nc">
    <div class="nc-title" data-i18n="bootstrap-title">Connect a New Node</div>
    <div style="font-size:0.63rem;color:var(--muted);line-height:1.8;margin-bottom:10px" data-i18n="bootstrap-desc">To run your own Aequitas node, set PRIMARY_NODE_URL=https://aequitas.digital in your environment. Your node registers automatically, syncs the full chain state, and begins participating in block production.</div>
    <div style="font-size:0.6rem;color:var(--muted);margin-bottom:6px;letter-spacing:1px">LIBP2P MULTIADDRESS</div>
    <div class="bsbox">/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R</div>
    <div style="font-size:0.6rem;color:var(--muted);margin-top:10px;line-height:1.7">Set in your environment: <span style="color:var(--purple);font-family:var(--font-mono)">PRIMARY_NODE_URL=https://aequitas.digital</span></div>
  </div>
  <div class="nc">
    <div class="nc-title" data-i18n="tech-title">Technical Specifications</div>
    <table class="spect">
      <tr><td data-i18n="k-chainid">Chain ID</td><td>1926 (0x786)</td></tr>
      <tr><td>Architecture</td><td style="color:var(--purple)">BlockDAG (Directed Acyclic Graph)</td></tr>
      <tr><td>EVM Compatible</td><td style="color:var(--green)" data-i18n="evm-yes">Yes — JSON-RPC /rpc · MetaMask</td></tr>
      <tr><td data-i18n="k-btime">Block Time</td><td>~6 seconds average</td></tr>
      <tr><td data-i18n="k-cons">Consensus</td><td style="color:var(--purple)">BlockDAG + Proof of Humanity</td></tr>
      <tr><td>P2P Protocol</td><td>libp2p (Go implementation)</td></tr>
      <tr><td>ZKP System</td><td>Groth16 / snarkjs / circom</td></tr>
      <tr><td>Elliptic Curve</td><td>BN128 (alt-bn128)</td></tr>
      <tr><td>Bio Hash</td><td style="color:var(--teal)">keccak256 (post-quantum safe)</td></tr>
      <tr><td data-i18n="k-storage">Storage</td><td style="color:var(--green)">PostgreSQL (persistent)</td></tr>
      <tr><td data-i18n="k-lang">Language</td><td>Go 1.24 (chain) · Node.js (proof server)</td></tr>
      <tr><td data-i18n="k-src">Source Code</td><td><a href="https://github.com/hanoi96international-gif/Aequitas" target="_blank" style="color:var(--blue)">GitHub — Open Source</a></td></tr>
    </table>
  </div>
  <div class="nc">
    <div class="nc-title" data-i18n="mm-config">MetaMask Configuration</div>
    <div style="font-size:0.62rem;color:var(--muted);line-height:1.7;margin-bottom:12px">Add Aequitas Chain to MetaMask to view your AEQ balance, send transactions, and interact with the V7 contract directly from your browser or mobile wallet.</div>
    <table class="spect">
      <tr><td data-i18n="k-chain">Network Name</td><td style="color:var(--gold)">Aequitas Chain</td></tr>
      <tr><td>RPC URL</td><td style="color:var(--blue);font-size:0.52rem">https://aequitas.digital/rpc</td></tr>
      <tr><td data-i18n="k-chainid">Chain ID</td><td style="color:var(--gold)">1926</td></tr>
      <tr><td data-i18n="k-symbol">Currency Symbol</td><td style="color:var(--gold)">AEQ</td></tr>
      <tr><td data-i18n="k-dec">Decimals</td><td>18</td></tr>
    </table>
    <button class="mm-btn" onclick="addToMetaMask()" style="margin-top:12px" data-i18n="btn-add-mm">+ ADD TO METAMASK</button>
    <div style="font-size:0.58rem;color:var(--muted);margin-top:8px;line-height:1.6">📱 MetaMask Mobile: if AEQ shows 0 after adding, delete the network and re-add it using the button above.</div>
</div>
</div>
</div>
<div id="net-story" class="stab-panel">
<div style="padding:24px;max-width:900px;margin:0 auto">

  <!-- HOW TO JOIN -->
  <div style="margin-bottom:20px">
    <div style="font-size:0.57rem;color:var(--purple);letter-spacing:2.5px;text-transform:uppercase;font-weight:700;margin-bottom:14px">Choose Your Path</div>
    <div style="display:grid;grid-template-columns:repeat(3,1fr);gap:12px">
      <div style="background:var(--card2);border:1px solid var(--border);border-radius:var(--radius);padding:20px;box-shadow:var(--shadow)">
        <div style="font-size:1.5rem;margin-bottom:10px">&#128100;</div>
        <div style="font-size:0.72rem;font-weight:700;color:var(--text);margin-bottom:8px">I am a Human</div>
        <div style="font-size:0.62rem;color:var(--muted);line-height:1.75;margin-bottom:14px">I want to register, receive 1,000 AEQ, and join the basic income network.</div>
        <div style="font-size:0.6rem;color:var(--purple);font-weight:600;line-height:1.9">1. Download AequitasBio app<br>2. Scan your biometric<br>3. Connect MetaMask<br>4. Receive 1,000 AEQ instantly</div>
      </div>
      <div style="background:var(--card2);border:1px solid var(--border);border-radius:var(--radius);padding:20px;box-shadow:var(--shadow)">
        <div style="font-size:1.5rem;margin-bottom:10px">&#128421;&#65039;</div>
        <div style="font-size:0.72rem;font-weight:700;color:var(--text);margin-bottom:8px">I am a Node Operator</div>
        <div style="font-size:0.62rem;color:var(--muted);line-height:1.75;margin-bottom:14px">I want to run a full node, participate in block production, and earn from the 40% validator pool.</div>
        <div style="font-size:0.6rem;color:var(--teal);font-weight:600;line-height:1.9">1. Register as a human (required)<br>2. Set PRIMARY_NODE_URL=https://aequitas.digital<br>3. Deploy on Railway/Contabo/VPS<br>4. Earn daily from validator pool</div>
      </div>
      <div style="background:var(--card2);border:1px solid var(--border);border-radius:var(--radius);padding:20px;box-shadow:var(--shadow)">
        <div style="font-size:1.5rem;margin-bottom:10px">&#128187;</div>
        <div style="font-size:0.72rem;font-weight:700;color:var(--text);margin-bottom:8px">I am a Developer</div>
        <div style="font-size:0.62rem;color:var(--muted);line-height:1.75;margin-bottom:14px">I want to build on Aequitas, integrate the API, or contribute to the protocol.</div>
        <div style="font-size:0.6rem;color:var(--gold);font-weight:600;line-height:1.9">1. EVM-compatible JSON-RPC<br>2. Chain ID: 1926 &#183; RPC: /rpc<br>3. OpenAPI: /api/* endpoints<br>4. Metrics: /metrics (Prometheus)</div>
      </div>
    </div>
  </div>

  <!-- TOKEN FLOW DIAGRAM -->
  <div class="idx" style="margin-bottom:16px">
    <div class="idx-title">AEQ Token Flow Diagram</div>
    <div style="overflow-x:auto;padding:8px 0">
    <svg viewBox="0 0 680 240" xmlns="http://www.w3.org/2000/svg" style="width:100%;max-width:680px;display:block;margin:0 auto;font-family:Inter,sans-serif">
      <rect width="680" height="240" fill="rgba(255,255,255,0.03)" rx="10"/>
      <rect x="20" y="90" width="110" height="60" rx="8" fill="rgba(107,70,193,0.1)" stroke="#6B46C1" stroke-width="1.5"/>
      <text x="75" y="115" text-anchor="middle" font-size="11" font-weight="700" fill="rgba(155,114,246,0.9)">HUMAN</text>
      <text x="75" y="130" text-anchor="middle" font-size="9" fill="rgba(136,146,164,0.9)">registers</text>
      <text x="75" y="143" text-anchor="middle" font-size="9" fill="rgba(136,146,164,0.9)">+1,000 AEQ minted</text>
      <line x1="130" y1="120" x2="190" y2="120" stroke="#6B46C1" stroke-width="1.5" marker-end="url(#arr)"/>
      <rect x="190" y="70" width="130" height="100" rx="8" fill="rgba(8,145,178,0.08)" stroke="#0891B2" stroke-width="1.5"/>
      <text x="255" y="95" text-anchor="middle" font-size="11" font-weight="700" fill="rgba(34,211,238,0.9)">AEQ ACTIVITY</text>
      <text x="255" y="112" text-anchor="middle" font-size="9" fill="rgba(136,146,164,0.9)">Swap fees (0.1%)</text>
      <text x="255" y="126" text-anchor="middle" font-size="9" fill="rgba(136,146,164,0.9)">Demurrage (0.5%/mo)</text>
      <text x="255" y="140" text-anchor="middle" font-size="9" fill="rgba(136,146,164,0.9)">Wealth cap overflow</text>
      <text x="255" y="154" text-anchor="middle" font-size="9" fill="rgba(136,146,164,0.9)">Inactive escrow</text>
      <line x1="320" y1="120" x2="380" y2="120" stroke="#0891B2" stroke-width="1.5" marker-end="url(#arr2)"/>
      <rect x="380" y="50" width="120" height="140" rx="8" fill="rgba(146,64,14,0.06)" stroke="#92400E" stroke-width="1.5"/>
      <text x="440" y="75" text-anchor="middle" font-size="10" font-weight="700" fill="rgba(240,180,41,0.9)">REDISTRIBUTION</text>
      <text x="440" y="95" text-anchor="middle" font-size="9" fill="rgba(52,211,153,0.9)">&#9679; UBI Pool 20%</text>
      <text x="440" y="112" text-anchor="middle" font-size="9" fill="rgba(155,114,246,0.9)">&#9679; Validators 40%</text>
      <text x="440" y="129" text-anchor="middle" font-size="9" fill="rgba(34,211,238,0.9)">&#9679; Liquidity LP 30%</text>
      <text x="440" y="146" text-anchor="middle" font-size="9" fill="rgba(240,180,41,0.9)">&#9679; Treasury 10%</text>
      <text x="440" y="165" text-anchor="middle" font-size="8" fill="rgba(136,146,164,0.9)">paid out daily</text>
      <text x="440" y="178" text-anchor="middle" font-size="8" fill="rgba(136,146,164,0.9)">automatic on-chain</text>
      <path d="M500 120 Q580 120 580 200 Q580 220 320 220 Q75 220 75 155" fill="none" stroke="#047857" stroke-width="1.5" stroke-dasharray="5,3" marker-end="url(#arr3)"/>
      <text x="340" y="215" text-anchor="middle" font-size="8" fill="rgba(52,211,153,0.9)">daily UBI returns to all verified humans</text>
      <defs>
        <marker id="arr" markerWidth="8" markerHeight="6" refX="6" refY="3" orient="auto"><polygon points="0 0, 8 3, 0 6" fill="rgba(155,114,246,0.9)"/></marker>
        <marker id="arr2" markerWidth="8" markerHeight="6" refX="6" refY="3" orient="auto"><polygon points="0 0, 8 3, 0 6" fill="rgba(34,211,238,0.9)"/></marker>
        <marker id="arr3" markerWidth="8" markerHeight="6" refX="6" refY="3" orient="auto"><polygon points="0 0, 8 3, 0 6" fill="rgba(52,211,153,0.9)"/></marker>
      </defs>
    </svg>
    </div>
  </div>

  <!-- NETWORK TOPOLOGY -->
  <div class="idx">
    <div class="idx-title">Network Topology &#8212; Current State</div>
    <div style="overflow-x:auto;padding:8px 0">
    <svg viewBox="0 0 500 185" xmlns="http://www.w3.org/2000/svg" style="width:100%;max-width:500px;display:block;margin:0 auto;font-family:Inter,sans-serif">
      <rect width="500" height="185" fill="rgba(255,255,255,0.03)" rx="10"/>
      <ellipse cx="250" cy="95" rx="60" ry="30" fill="rgba(107,70,193,0.08)" stroke="rgba(107,70,193,0.3)" stroke-width="1" stroke-dasharray="4,3"/>
      <text x="250" y="91" text-anchor="middle" font-size="9" fill="rgba(155,114,246,0.9)">P2P libp2p</text>
      <text x="250" y="105" text-anchor="middle" font-size="8" fill="rgba(136,146,164,0.9)">BlockDAG sync</text>
      <rect x="30" y="60" width="120" height="70" rx="8" fill="rgba(4,120,87,0.08)" stroke="#047857" stroke-width="1.5"/>
      <text x="90" y="85" text-anchor="middle" font-size="10" font-weight="700" fill="rgba(52,211,153,0.9)">Node 1</text>
      <text x="90" y="100" text-anchor="middle" font-size="8" fill="rgba(136,146,164,0.9)">Railway (Primary)</text>
      <text x="90" y="116" text-anchor="middle" font-size="7.5" fill="rgba(52,211,153,0.9)">&#9679; Primary API &#183; PostgreSQL</text>
      <text x="90" y="126" text-anchor="middle" font-size="7.5" fill="rgba(52,211,153,0.9)">&#9679; UBI distribution</text>
      <rect x="350" y="60" width="120" height="70" rx="8" fill="rgba(8,145,178,0.08)" stroke="#0891B2" stroke-width="1.5"/>
      <text x="410" y="85" text-anchor="middle" font-size="10" font-weight="700" fill="rgba(34,211,238,0.9)">Node 2</text>
      <text x="410" y="100" text-anchor="middle" font-size="8" fill="rgba(136,146,164,0.9)">Contabo VPS (Secondary)</text>
      <text x="410" y="116" text-anchor="middle" font-size="7.5" fill="rgba(34,211,238,0.9)">&#9679; Secondary API</text>
      <text x="410" y="126" text-anchor="middle" font-size="7.5" fill="rgba(34,211,238,0.9)">&#9679; HTTP sync</text>
      <rect x="175" y="10" width="150" height="38" rx="8" fill="rgba(146,64,14,0.06)" stroke="#92400E" stroke-width="1"/>
      <text x="250" y="27" text-anchor="middle" font-size="9" font-weight="700" fill="rgba(240,180,41,0.9)">MetaMask / Users</text>
      <text x="250" y="41" text-anchor="middle" font-size="7.5" fill="rgba(136,146,164,0.9)">JSON-RPC &#183; Chain ID 1926</text>
      <rect x="165" y="143" width="170" height="34" rx="8" fill="rgba(107,70,193,0.04)" stroke="rgba(107,70,193,0.2)" stroke-width="1" stroke-dasharray="4,3"/>
      <text x="250" y="158" text-anchor="middle" font-size="8" fill="rgba(136,146,164,0.9)">Future nodes (Phase 1+)</text>
      <text x="250" y="170" text-anchor="middle" font-size="7.5" fill="rgba(107,70,193,0.6)">registered humans can run validator nodes</text>
      <line x1="150" y1="95" x2="190" y2="95" stroke="#6B46C1" stroke-width="1.5" stroke-dasharray="4,3"/>
      <line x1="310" y1="95" x2="350" y2="95" stroke="#6B46C1" stroke-width="1.5" stroke-dasharray="4,3"/>
      <line x1="250" y1="48" x2="250" y2="65" stroke="#6B46C1" stroke-width="1.5" stroke-dasharray="4,3"/>
    </svg>
    </div>
  </div>

</div>
</div>

<div id="net-runnode" class="stab-panel">
<div class="ns" style="grid-template-columns:1fr">
<div class="nc" style="grid-column:1/-1;background:linear-gradient(135deg,rgba(245,166,35,0.06),rgba(13,8,32,0.9));border-color:rgba(245,166,35,0.2)">
    <div class="nc-title" style="color:var(--gold)" data-i18n="run-node-title">Run Your Own Node — Help Secure the Network</div>
    <div style="font-size:0.67rem;color:var(--muted);line-height:1.9;margin-bottom:16px" data-i18n="run-node-desc">Registered humans can run an Aequitas validator node — no stake, no application required, but <strong style="color:var(--gold)">NODE_OPERATOR_WALLET must be a registered Aequitas human</strong> (this is a protocol security requirement). Nodes participate in block production, validate the human registry, and synchronize the BlockDAG. Node operators earn a share of protocol fees via the Validators Pool (40% of all protocol fees, distributed daily). The more nodes that run, the more decentralized and resilient the network becomes.</div>
    <div style="display:flex;gap:12px;flex-wrap:wrap;margin-bottom:16px">
      <a id="node-guide-pdf-btn" href="/download/node-guide-en.pdf" download style="display:inline-flex;align-items:center;gap:8px;background:var(--gold);color:#06091A;padding:12px 20px;border-radius:8px;font-size:0.7rem;font-weight:700;text-decoration:none;transition:opacity 0.2s" onmouseover="this.style.opacity=0.87" onmouseout="this.style.opacity=1">
        📄 Node Operator Guide (PDF)
      </a>
      <a href="https://github.com/hanoi96international-gif/Aequitas" target="_blank" style="display:inline-flex;align-items:center;gap:8px;background:rgba(139,92,246,0.15);color:var(--purple);border:1px solid rgba(139,92,246,0.3);padding:12px 20px;border-radius:8px;font-size:0.7rem;font-weight:700;text-decoration:none;transition:all 0.2s" onmouseover="this.style.opacity=0.87" onmouseout="this.style.opacity=1">
        🐙 View Source on GitHub
      </a>
    </div>
    <!-- INLINE NODE GUIDE -->
    <div id="node-guide" style="display:block;background:var(--card);border:1px solid rgba(245,166,35,0.2);border-radius:var(--radius);padding:24px;margin-top:4px">

      <!-- Header -->
      <div style="display:flex;align-items:center;justify-content:space-between;flex-wrap:wrap;gap:8px;margin-bottom:20px">
        <div>
          <div style="font-size:0.58rem;color:var(--gold);letter-spacing:2.5px;text-transform:uppercase;font-weight:700;display:flex;align-items:center;gap:8px">
            AEQUITAS NODE OPERATOR GUIDE
            <span style="font-size:0.52rem;background:rgba(245,166,35,0.12);border:1px solid rgba(245,166,35,0.3);color:var(--gold);padding:2px 8px;border-radius:10px">v1.0 · June 2026</span>
          </div>
          <div style="font-size:0.6rem;color:var(--muted);margin-top:4px">Complete step-by-step guide &middot; No prior blockchain experience required &middot; Estimated time: 20&ndash;30 min</div>
        </div>
      </div>

      <!-- What is a node -->
      <div style="background:rgba(139,92,246,0.06);border:1px solid rgba(139,92,246,0.2);border-radius:8px;padding:14px;margin-bottom:20px">
        <div style="font-size:0.6rem;color:var(--purple);font-weight:700;margin-bottom:6px">What is an Aequitas Node?</div>
        <div style="font-size:0.62rem;color:var(--muted);line-height:1.9">An Aequitas node is a program that runs in the cloud and participates in the Aequitas network. It keeps a copy of the entire blockchain, validates who is a registered human, and produces new blocks (like new pages in the global ledger). The more nodes exist, the more decentralized and resilient the network becomes. As a reward for running a node, you receive a daily share of all protocol fees &mdash; automatically, with no further action required on your part.</div>
      </div>

      <!-- Checklist -->
      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">Before You Start &mdash; What You Need</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:2.1;margin-bottom:18px">
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="color:var(--gold);font-weight:700;min-width:16px">1.</span><span><strong style="color:var(--text)">An Aequitas account:</strong> You must first be registered as a human on Aequitas. Install the Android app, complete biometric registration, and note your wallet address. Without this, you cannot receive validator rewards.</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="color:var(--gold);font-weight:700;min-width:16px">2.</span><span><strong style="color:var(--text)">A GitHub account (free):</strong> Go to github.com and create a free account. You need this to copy (fork) the Aequitas code so Railway can deploy it.</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="color:var(--gold);font-weight:700;min-width:16px">3.</span><span><strong style="color:var(--text)">A Railway account (free):</strong> Go to railway.app and sign in with GitHub. Railway is a hosting platform that runs your node in the cloud &mdash; no server or command line required.</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="color:var(--gold);font-weight:700;min-width:16px">4.</span><span><strong style="color:var(--text)">Node signing key (RELAYER_PRIVATE_KEY):</strong> Your node needs a dedicated Ethereum wallet to sign on-chain registrations. This can be any MetaMask wallet. Export its private key: MetaMask &rarr; Account Details &rarr; Show Private Key &rarr; enter password &rarr; copy. Keep strictly private. <strong style="color:var(--gold)">IMPORTANT:</strong> To receive validator rewards you also need NODE_OPERATOR_WALLET set to your <strong style="color:var(--neon)">registered Aequitas human wallet</strong> (the one verified with AequitasBio). Only verified humans can earn validator rewards.</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px"><span style="color:var(--gold);font-weight:700;min-width:16px">5.</span><span><strong style="color:var(--text)">10&ndash;30 minutes of your time.</strong> Railway does most of the work automatically.</span></div>
      </div>

      <!-- Step 1: Fork -->
      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">Step 1 &mdash; Fork the Aequitas Repository on GitHub</div>
      <div style="background:rgba(0,220,170,0.05);border:1px solid rgba(0,220,170,0.15);border-radius:6px;padding:10px 14px;margin-bottom:10px;font-size:0.6rem"><span style="color:var(--teal);font-weight:700">What is a fork?</span> <span style="color:var(--muted)">A fork is your own personal copy of the Aequitas code on GitHub. Railway deploys directly from your fork, so you need one first.</span></div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:2.2;margin-bottom:18px">
        <div><span style="color:var(--gold);font-weight:700">a)</span> Open <span style="font-family:var(--font-mono);color:var(--neon)">github.com/hanoi96international-gif/Aequitas</span> in your browser</div>
        <div><span style="color:var(--gold);font-weight:700">b)</span> Click the <strong style="color:var(--text)">Fork</strong> button in the top-right corner of the page</div>
        <div><span style="color:var(--gold);font-weight:700">c)</span> Click <strong style="color:var(--text)">Create fork</strong> &mdash; GitHub creates a copy under your own account (e.g. github.com/YOUR-NAME/Aequitas)</div>
        <div><span style="color:var(--gold);font-weight:700">d)</span> Done &mdash; you now have your own copy of the Aequitas node code</div>
      </div>

      <!-- Step 2: Database -->
      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">Step 2 &mdash; Create a PostgreSQL Database</div>
      <div style="background:rgba(0,220,170,0.05);border:1px solid rgba(0,220,170,0.15);border-radius:6px;padding:10px 14px;margin-bottom:10px;font-size:0.6rem"><span style="color:var(--teal);font-weight:700">What is a database?</span> <span style="color:var(--muted)">Your node needs permanent storage for all block data and human registrations. Think of it like a cloud hard drive. Without a database, your node loses all data every time it restarts. PostgreSQL is the storage system Aequitas uses.</span></div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:2.2;margin-bottom:18px">
        <div><span style="color:var(--gold);font-weight:700">a)</span> Go to <strong style="color:var(--text)">railway.app</strong> and sign in with your GitHub account</div>
        <div><span style="color:var(--gold);font-weight:700">b)</span> Click <strong style="color:var(--text)">New Project</strong></div>
        <div><span style="color:var(--gold);font-weight:700">c)</span> Inside your new project, click <strong style="color:var(--text)">+ New</strong> &rarr; <strong style="color:var(--text)">Database</strong> &rarr; <strong style="color:var(--text)">Add PostgreSQL</strong></div>
        <div><span style="color:var(--gold);font-weight:700">d)</span> Railway creates your database automatically. When you add the node service to the same project in Step 4, Railway injects the DATABASE_URL connection string for you &mdash; no manual configuration needed.</div>
      </div>

      <!-- Step 3: Env vars -->
      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">Step 3 &mdash; Understand the Environment Variables</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9;margin-bottom:10px">Environment variables are configuration settings you pass to your node before it starts. Think of them like a settings file. Collect these values before deploying &mdash; you will enter them in Step 4.</div>
      <div style="background:rgba(220,50,50,0.06);border:1px solid rgba(220,50,50,0.2);border-radius:6px;padding:10px 14px;margin-bottom:12px;font-size:0.6rem"><span style="color:#f87171;font-weight:700">Security Warning: </span><span style="color:var(--muted)">Your RELAYER_PRIVATE_KEY is like a master password. Anyone who has it controls your node wallet. Never share it publicly, never paste it in chat or email. Use a separate MetaMask wallet for RELAYER_PRIVATE_KEY (signing). NODE_OPERATOR_WALLET (for rewards) must be your registered Aequitas human wallet.</span></div>
      <table style="width:100%;border-collapse:collapse;margin-bottom:18px">
        <tr style="background:rgba(139,92,246,0.15)">
          <td style="font-size:0.6rem;color:var(--text);padding:8px;font-weight:700;width:36%">Variable</td>
          <td style="font-size:0.6rem;color:var(--text);padding:8px;font-weight:700;width:12%">Required?</td>
          <td style="font-size:0.6rem;color:var(--text);padding:8px;font-weight:700">What to enter and where to find it</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--neon);padding:8px">DATABASE_URL</td>
          <td style="font-size:0.6rem;color:#f87171;padding:8px;font-weight:700">YES</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Your PostgreSQL connection string. On Railway: auto-injected when PostgreSQL is in the same project. Format: <span style="font-family:var(--font-mono)">postgres://user:pass@host:5432/dbname</span></td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--neon);padding:8px">RELAYER_PRIVATE_KEY</td>
          <td style="font-size:0.6rem;color:#f87171;padding:8px;font-weight:700">YES</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">The private key (starts with 0x, 66 characters total) of your dedicated node wallet. In MetaMask: click account icon &rarr; Account Details &rarr; Show Private Key &rarr; enter your MetaMask password &rarr; copy the key</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--neon);padding:8px">RELAYER_ADDRESS</td>
          <td style="font-size:0.6rem;color:var(--teal);padding:8px">Recommended</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">The wallet address (starts with 0x, 42 characters) matching RELAYER_PRIVATE_KEY. This is the public address &mdash; safe to share. Copy it from MetaMask. A fallback exists in the node code, but setting this explicitly prevents startup errors.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--neon);padding:8px">NODE_OPERATOR_WALLET</td>
          <td style="font-size:0.6rem;color:var(--teal);padding:8px">For rewards</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Your Aequitas human wallet address &mdash; the one you registered with via the Android app. This wallet receives your daily validator rewards (40% of all protocol fees). Must be a registered human on Aequitas. Find it in the app under your profile.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--neon);padding:8px">PEER_SECRET</td>
          <td style="font-size:0.6rem;color:#f87171;padding:8px;font-weight:700">For multi-node</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">A shared secret string that authorises your node as a validator in the network. Every node in the same network must use the <strong style="color:var(--text)">identical value</strong>. Without it your node syncs blocks but its own blocks are rejected by the primary. Get this value from the network operator — do not share publicly.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--neon);padding:8px">SELF_URL</td>
          <td style="font-size:0.6rem;color:var(--teal);padding:8px">For multi-node</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Your node's own public HTTPS URL (e.g. <span style="font-family:var(--font-mono);color:var(--neon)">https://my-node.up.railway.app</span>). Required for peer discovery self-exclusion — without it your node may try to sync from itself. Find your URL in Railway: Settings &rarr; Networking &rarr; Public Networking.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--neon);padding:8px">PRIMARY_NODE_URL</td>
          <td style="font-size:0.6rem;color:var(--teal);padding:8px">For multi-node</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Set to: <span style="font-family:var(--font-mono);color:var(--neon)">https://aequitas.digital</span> &mdash; the primary node your node registers with for automatic peer discovery. On startup your node posts its URL + signing address to the primary, gets the full peer list back, and joins the network automatically. No manual PEER_NODES list needed.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--muted);padding:8px">PORT</td>
          <td style="font-size:0.6rem;color:var(--muted);padding:8px">No</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Leave unset on Railway &mdash; Railway sets this automatically. Default is 8080.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--muted);padding:8px">NODE_KEY</td>
          <td style="font-size:0.6rem;color:var(--muted);padding:8px">No</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Base64-encoded libp2p private key for stable P2P identity. If not set: auto-generated on first start and printed to stderr. Copy the base64 string from <span style="font-family:var(--font-mono);color:var(--gold)">SAVE THIS AS NODE_KEY ENVIRONMENT VAR: &lt;base64&gt;</span> and paste it here to keep a stable peer ID across restarts.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--muted);padding:8px">IS_PRIMARY_NODE</td>
          <td style="font-size:0.6rem;color:var(--muted);padding:8px">No</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">No longer needed. Distribution is now coordinated via a PostgreSQL CAS lock (TryLockDistribution) — every node schedules at 20:00 Berlin time but only the first to win the DB lock runs the distribution. Leave unset on all nodes.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--neon);padding:8px">BOOTSTRAP_SNAPSHOT_URL</td>
          <td style="font-size:0.6rem;color:var(--muted);padding:8px">Multi-node</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Set to <code style="color:var(--gold)">https://aequitas.digital/api/snapshot</code> on a fresh node. If the local DB has 0 humans at startup, the node automatically downloads and imports the full state from this URL — fixing StateRoot mismatches immediately. Also set BOOTSTRAP_SIGNER and SNAPSHOT_TOKEN.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--neon);padding:8px">BOOTSTRAP_SIGNER</td>
          <td style="font-size:0.6rem;color:var(--muted);padding:8px">Multi-node</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Ethereum signing address of the primary node. Required when BOOTSTRAP_SNAPSHOT_URL is set. Get it from <code style="color:var(--gold)">https://aequitas.digital/api/status</code> → <code>signing_address</code> field. Prevents importing a tampered snapshot.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--neon);padding:8px">SNAPSHOT_TOKEN</td>
          <td style="font-size:0.6rem;color:var(--muted);padding:8px">Multi-node</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Bearer token to authenticate the /api/snapshot endpoint. Must match the value set on the primary node. Get it from the network operator.</td>
        </tr>
        <tr>
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:#f87171;padding:8px">RESET_STATE</td>
          <td style="font-size:0.6rem;color:var(--muted);padding:8px">No</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">DANGEROUS: Setting this to true wipes your entire database on every restart. Development use only. Never in production.</td>
        </tr>
      </table>

      <!-- Step 4 Railway -->
      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">Step 4 &mdash; Deploy on Railway (Recommended)</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9;margin-bottom:12px">Railway is the easiest way to run your node &mdash; no server setup, no command line required. The free tier covers all requirements. Total time: about 10&ndash;15 minutes.</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:2.2;margin-bottom:18px">
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="display:inline-flex;align-items:center;justify-content:center;background:rgba(139,92,246,0.2);color:var(--purple);font-weight:700;font-size:0.58rem;min-width:22px;height:22px;border-radius:50%">1</span><span>In your Railway project (from Step 2), click <strong style="color:var(--text)">+ New</strong> &rarr; <strong style="color:var(--text)">GitHub Repo</strong></span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="display:inline-flex;align-items:center;justify-content:center;background:rgba(139,92,246,0.2);color:var(--purple);font-weight:700;font-size:0.58rem;min-width:22px;height:22px;border-radius:50%">2</span><span>Select your Aequitas fork (from Step 1) &mdash; Railway detects the Dockerfile automatically</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="display:inline-flex;align-items:center;justify-content:center;background:rgba(139,92,246,0.2);color:var(--purple);font-weight:700;font-size:0.58rem;min-width:22px;height:22px;border-radius:50%">3</span><span>Click <strong style="color:var(--text)">Deploy Now</strong> &mdash; a first build starts (may fail without env vars, that is normal)</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="display:inline-flex;align-items:center;justify-content:center;background:rgba(139,92,246,0.2);color:var(--purple);font-weight:700;font-size:0.58rem;min-width:22px;height:22px;border-radius:50%">4</span><span>Click your Aequitas service &rarr; <strong style="color:var(--text)">Variables</strong> &rarr; add each variable:</span></div>
        <div style="font-family:var(--font-mono);background:rgba(0,0,0,0.3);border:1px solid rgba(139,92,246,0.15);border-radius:6px;padding:12px;margin:4px 0 8px 32px;font-size:0.61rem;line-height:2.1;overflow-x:auto">
          <span style="color:var(--muted)"># Railway auto-sets DATABASE_URL if PostgreSQL is in the same project</span><br>
          <span style="color:var(--neon)">RELAYER_PRIVATE_KEY</span> = <span style="color:var(--gold)">0xYOUR_PRIVATE_KEY</span><br>
          <span style="color:var(--neon)">RELAYER_ADDRESS</span> = <span style="color:var(--gold)">0xYOUR_NODE_WALLET_ADDRESS</span><br>
          <span style="color:var(--neon)">NODE_OPERATOR_WALLET</span> = <span style="color:var(--gold)">0xYOUR_HUMAN_WALLET</span><br>
          <span style="color:var(--neon)">PEER_SECRET</span> = <span style="color:var(--gold)">get-this-from-network-operator</span><br>
          <span style="color:var(--neon)">SELF_URL</span> = <span style="color:var(--gold)">https://YOUR-RAILWAY-DOMAIN.up.railway.app</span><br>
          <span style="color:var(--neon)">PRIMARY_NODE_URL</span> = <span style="color:var(--gold)">https://aequitas.digital</span>
        </div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="display:inline-flex;align-items:center;justify-content:center;background:rgba(139,92,246,0.2);color:var(--purple);font-weight:700;font-size:0.58rem;min-width:22px;height:22px;border-radius:50%">5</span><span>Click <strong style="color:var(--text)">Deploy</strong> (or save variables to trigger auto-redeploy). Build takes ~3 minutes while Go compiles the node binary.</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="display:inline-flex;align-items:center;justify-content:center;background:rgba(139,92,246,0.2);color:var(--purple);font-weight:700;font-size:0.58rem;min-width:22px;height:22px;border-radius:50%">6</span><span>Watch <strong style="color:var(--text)">Deploy Logs</strong>. Success looks like: <span style="font-family:var(--font-mono);color:var(--teal)">Aequitas Node Running V</span> and <span style="font-family:var(--font-mono);color:var(--teal)">[NODE] Registered node operator wallet: 0x...</span></span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="display:inline-flex;align-items:center;justify-content:center;background:rgba(139,92,246,0.2);color:var(--purple);font-weight:700;font-size:0.58rem;min-width:22px;height:22px;border-radius:50%">7</span><span>Go to <strong style="color:var(--text)">Settings</strong> &rarr; <strong style="color:var(--text)">Networking</strong> &rarr; <strong style="color:var(--text)">Generate Domain</strong> to get your public URL</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px"><span style="display:inline-flex;align-items:center;justify-content:center;background:rgba(139,92,246,0.2);color:var(--purple);font-weight:700;font-size:0.58rem;min-width:22px;height:22px;border-radius:50%">8</span><span>Open <span style="font-family:var(--font-mono);color:var(--neon)">https://YOUR-URL/api/status</span> in your browser &mdash; you should see JSON with <strong style="color:var(--text)">height</strong> climbing every ~6 seconds</span></div>
      </div>

      <!-- Step 4b Docker -->
      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">Step 4b &mdash; Alternative: Deploy with Docker / VPS</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9;margin-bottom:8px">For your own server (VPS like Contabo, Hetzner, DigitalOcean). Requires Docker and a PostgreSQL database. You can share Railway's PostgreSQL as an external DB — use the external proxy URL from Railway's Postgres service variables. <strong style="color:var(--gold)">Important:</strong> NODE_OPERATOR_WALLET must be a registered Aequitas human wallet — this is required by the protocol to run a validator node.</div>
      <div style="font-family:var(--font-mono);background:rgba(0,0,0,0.3);border:1px solid rgba(139,92,246,0.15);border-radius:6px;padding:14px;margin-bottom:8px;font-size:0.61rem;line-height:2.2;overflow-x:auto">
        <span style="color:var(--muted)"># 1. Clone and build (~3 min)</span><br>
        git clone https://github.com/hanoi96international-gif/Aequitas &amp;&amp; cd Aequitas<br>
        docker build -t aequitas-node .<br><br>
        <span style="color:var(--muted)"># 2. First start (without NODE_KEY — copy it from the logs after)</span><br>
        docker run -d --name aequitas-node --restart unless-stopped \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">DATABASE_URL</span>="<span style="color:var(--gold)">postgresql://user:pass@host:5432/dbname</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">RELAYER_PRIVATE_KEY</span>="<span style="color:var(--gold)">YOUR_HEX_PRIVATE_KEY</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">RELAYER_ADDRESS</span>="<span style="color:var(--gold)">0xYOUR_NODE_SIGNING_ADDRESS</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">NODE_OPERATOR_WALLET</span>="<span style="color:var(--gold)">0xYOUR_REGISTERED_HUMAN_WALLET</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">PEER_SECRET</span>="<span style="color:var(--gold)">get-from-network-operator</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">SELF_URL</span>="<span style="color:var(--gold)">http://YOUR-SERVER-IP:8080</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">PRIMARY_NODE_URL</span>="https://aequitas.digital" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">BOOTSTRAP_SNAPSHOT_URL</span>="https://aequitas.digital/api/snapshot" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">BOOTSTRAP_SIGNER</span>="0x92cbedec9d348b4762cb9af99500ee6139c5b671" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">SNAPSHOT_TOKEN</span>="<span style="color:var(--gold)">get-from-network-operator</span>" \<br>
        &nbsp;&nbsp;-p 8080:8080 -p 4001:4001 aequitas-node<br><br>
        <span style="color:var(--muted)"># 3. Get NODE_KEY from logs (one-time):</span><br>
        docker logs aequitas-node 2&gt;&amp;1 | grep "SAVE THIS AS NODE_KEY"<br><br>
        <span style="color:var(--muted)"># 4. Final command with NODE_KEY (use this from now on):</span><br>
        docker stop aequitas-node &amp;&amp; docker rm aequitas-node<br>
        docker run -d --name aequitas-node --restart unless-stopped \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">DATABASE_URL</span>="<span style="color:var(--gold)">postgresql://user:pass@host:5432/dbname</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">RELAYER_PRIVATE_KEY</span>="<span style="color:var(--gold)">YOUR_HEX_PRIVATE_KEY</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">RELAYER_ADDRESS</span>="<span style="color:var(--gold)">0xYOUR_NODE_SIGNING_ADDRESS</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">NODE_OPERATOR_WALLET</span>="<span style="color:var(--gold)">0xYOUR_REGISTERED_HUMAN_WALLET</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">NODE_KEY</span>="<span style="color:var(--gold)">base64-from-step-3</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">PEER_SECRET</span>="<span style="color:var(--gold)">get-from-network-operator</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">SELF_URL</span>="<span style="color:var(--gold)">http://YOUR-SERVER-IP:8080</span>" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">PRIMARY_NODE_URL</span>="https://aequitas.digital" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">BOOTSTRAP_SNAPSHOT_URL</span>="https://aequitas.digital/api/snapshot" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">BOOTSTRAP_SIGNER</span>="0x92cbedec9d348b4762cb9af99500ee6139c5b671" \<br>
        &nbsp;&nbsp;-e <span style="color:var(--neon)">SNAPSHOT_TOKEN</span>="<span style="color:var(--gold)">get-from-network-operator</span>" \<br>
        &nbsp;&nbsp;-p 8080:8080 -p 4001:4001 aequitas-node
      </div>
      <div style="background:rgba(245,166,35,0.06);border:1px solid rgba(245,166,35,0.2);border-radius:6px;padding:8px 12px;margin-bottom:18px;font-size:0.6rem;color:var(--muted)">
        <strong style="color:var(--gold)">Tip:</strong> Save all vars in <code>/root/.aequitas.env</code> (chmod 600) and use <code>--env-file /root/.aequitas.env</code> instead of listing each <code>-e</code> — keeps secrets out of shell history and simplifies updates.
      </div>
      <div style="background:rgba(0,220,170,0.05);border:1px solid rgba(0,220,170,0.15);border-radius:6px;padding:10px 14px;margin-bottom:18px;font-size:0.62rem;color:var(--muted)">
        <strong style="color:var(--neon)">Port requirements:</strong> TCP 8080 must be open inbound (API + RPC). TCP 4001 is optional (P2P — enables direct node-to-node connections). If P2P is firewalled, HTTP sync still works. On Linux: <code style="color:var(--teal)">ufw allow 8080/tcp</code>
      </div>

      <!-- Step 5: Verify -->
      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">Step 5 &mdash; Verify Your Node is Running</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9;margin-bottom:8px">Open these URLs in your browser. Replace <span style="font-family:var(--font-mono);color:var(--neon)">YOUR-NODE-URL</span> with your actual Railway domain or server address.</div>
      <div style="font-family:var(--font-mono);background:rgba(0,0,0,0.3);border:1px solid rgba(139,92,246,0.15);border-radius:6px;padding:14px;margin-bottom:8px;font-size:0.61rem;line-height:2.2">
        https://YOUR-NODE-URL/api/status<br>
        <span style="color:var(--muted)">&nbsp;&rarr; Expected: {"height": 1234, "total_humans": N, "aequitas_index": N}</span><br><br>
        https://YOUR-NODE-URL/rpc<br>
        <span style="color:var(--muted)">&nbsp;&rarr; Expected: {"jsonrpc":"2.0","error":"method not specified"} &mdash; this confirms RPC is alive</span>
      </div>
      <div style="background:rgba(0,220,170,0.05);border:1px solid rgba(0,220,170,0.15);border-radius:6px;padding:10px 14px;margin-bottom:18px;font-size:0.62rem;color:var(--muted)">The block height should match the primary node within 1&ndash;2 blocks within seconds of startup. If it stays at 0, check that PRIMARY_NODE_URL=https://aequitas.digital is set and reachable.</div>

      <!-- Step 5b: Register Validator Key -->
      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">Step 5b &mdash; Register Your Validator Key (Required for Block Rewards)</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9;margin-bottom:10px">Link your node's signing key to your registered human wallet. This proves you control both keys and authorises your node to earn validator rewards. Call this from <strong>any machine</strong> — your laptop, the VPS itself, wherever. Requires your <code style="color:var(--gold)">SNAPSHOT_TOKEN</code> in the Authorization header:</div>
      <div style="font-family:var(--font-mono);background:rgba(0,0,0,0.35);border:1px solid rgba(139,92,246,0.15);border-radius:6px;padding:10px 14px;font-size:0.6rem;color:var(--teal);margin-bottom:12px;overflow-x:auto">
        <span style="color:var(--muted)"># From any machine — replace values with yours:</span><br>
        curl -H "Authorization: Bearer <span style="color:var(--gold)">YOUR-SNAPSHOT-TOKEN</span>" \<br>
        &nbsp;&nbsp;"http://<span style="color:var(--gold)">YOUR-NODE-IP:8080</span>/api/sign-validator-challenge?wallet=<span style="color:var(--gold)">0xYOUR_HUMAN_WALLET</span>"<br><br>
        <span style="color:var(--muted)"># On the VPS server itself (no token needed from localhost):</span><br>
        curl -H "Authorization: Bearer <span style="color:var(--gold)">YOUR-SNAPSHOT-TOKEN</span>" \<br>
        &nbsp;&nbsp;"http://localhost:8080/api/sign-validator-challenge?wallet=<span style="color:var(--gold)">0xYOUR_HUMAN_WALLET</span>"
      </div>
      <div style="background:rgba(245,166,35,0.06);border:1px solid rgba(245,166,35,0.2);border-radius:6px;padding:8px 12px;font-size:0.6rem;color:var(--gold);margin-bottom:12px">⚠ The endpoint requires <code>SNAPSHOT_TOKEN</code> if set on the node. Get the token from the network operator (it is the same value you set as SNAPSHOT_TOKEN on your node).</div>
      <div id="vk-reg-box" style="background:rgba(139,92,246,0.05);border:1px solid rgba(139,92,246,0.2);border-radius:8px;padding:16px;margin-bottom:18px">
        <div style="font-size:0.6rem;color:var(--muted);margin-bottom:8px">Enter your node RELAYER_ADDRESS and the signature from the command above:</div>
        <input id="vk-signing-addr" placeholder="0x... (RELAYER_ADDRESS — your node signing address)" style="width:100%;box-sizing:border-box;background:rgba(0,0,0,0.3);border:1px solid rgba(139,92,246,0.3);color:var(--text);border-radius:6px;padding:8px 12px;font-family:var(--font-mono);font-size:0.62rem;margin-bottom:6px">
        <input id="vk-signing-sig" placeholder='Signing key signature (from curl output, the "signature" field)' style="width:100%;box-sizing:border-box;background:rgba(0,0,0,0.3);border:1px solid rgba(139,92,246,0.3);color:var(--text);border-radius:6px;padding:8px 12px;font-family:var(--font-mono);font-size:0.62rem;margin-bottom:8px">
        <button onclick="registerValidatorKey()" style="background:rgba(139,92,246,0.8);color:#fff;border:none;border-radius:6px;padding:10px 20px;font-size:0.65rem;cursor:pointer;font-weight:700">&#x1F511; Sign with MetaMask &amp; Register</button>
        <div id="vk-status" style="margin-top:8px;font-size:0.6rem;color:var(--muted)"></div>
      </div>

      <!-- Step 6: MetaMask -->
      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">Step 6 &mdash; Connect MetaMask to Your Node (Optional)</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9;margin-bottom:8px">You can use your own node as a custom RPC in MetaMask so your wallet connects through your node instead of the shared public node. In MetaMask: click the network dropdown at the top &rarr; <strong style="color:var(--text)">Add network</strong> &rarr; <strong style="color:var(--text)">Add a network manually</strong>, then enter:</div>
      <table style="width:100%;border-collapse:collapse;margin-bottom:18px">
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)"><td style="font-size:0.62rem;color:var(--muted);padding:7px 0;width:40%">Network Name</td><td style="font-size:0.62rem;font-family:var(--font-mono);color:var(--text);padding:7px 0">Aequitas Chain</td></tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)"><td style="font-size:0.62rem;color:var(--muted);padding:7px 0">RPC URL</td><td style="font-size:0.62rem;font-family:var(--font-mono);color:var(--neon);padding:7px 0">https://YOUR-NODE-URL/rpc</td></tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)"><td style="font-size:0.62rem;color:var(--muted);padding:7px 0">Chain ID</td><td style="font-size:0.62rem;font-family:var(--font-mono);color:var(--text);padding:7px 0">1926</td></tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)"><td style="font-size:0.62rem;color:var(--muted);padding:7px 0">Currency Symbol</td><td style="font-size:0.62rem;font-family:var(--font-mono);color:var(--text);padding:7px 0">AEQ</td></tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)"><td style="font-size:0.62rem;color:var(--muted);padding:7px 0">Decimals</td><td style="font-size:0.62rem;font-family:var(--font-mono);color:var(--text);padding:7px 0">18</td></tr>
        <tr><td style="font-size:0.62rem;color:var(--muted);padding:7px 0">Block Explorer</td><td style="font-size:0.62rem;font-family:var(--font-mono);color:var(--purple);padding:7px 0">https://aequitas.digital</td></tr>
      </table>

      <!-- Step 7: Rewards -->
      <div style="font-size:0.58rem;color:var(--gold);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid rgba(245,166,35,0.2);padding-bottom:6px">Step 7 &mdash; Earning Validator Rewards</div>
      <div style="background:rgba(245,166,35,0.05);border:1px solid rgba(245,166,35,0.2);border-radius:6px;padding:12px 14px;margin-bottom:12px;font-size:0.62rem;color:var(--muted);line-height:1.9">The Validators Pool collects 40% of all protocol fees (swap fees, demurrage, wealth cap overflow). Every day at <strong style="color:var(--gold)">20:00 Berlin time</strong> (CEST/CET, handles DST automatically) the primary node distributes the pool balance to all registered node operator wallets proportionally. The more consistently your node runs, the larger your share.</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:2.2;margin-bottom:18px">
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:4px"><span style="color:var(--gold);font-weight:700;min-width:16px">1.</span><span>Make sure you are registered as a human on Aequitas. If not: install the Android app and complete biometric registration first. You will receive a wallet address and 1,000 AEQ.</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:4px"><span style="color:var(--gold);font-weight:700;min-width:16px">2.</span><span>Set <span style="font-family:var(--font-mono);color:var(--neon)">NODE_OPERATOR_WALLET</span> = your Aequitas human wallet address in your Railway Variables</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:4px"><span style="color:var(--gold);font-weight:700;min-width:16px">3.</span><span>Save &mdash; Railway redeploys automatically. On Docker: <span style="font-family:var(--font-mono);color:var(--teal)">docker restart aequitas-node</span></span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:4px"><span style="color:var(--gold);font-weight:700;min-width:16px">4.</span><span>In your node logs, confirm: <span style="font-family:var(--font-mono);color:var(--teal)">[NODE] Registered node operator wallet: 0x...</span></span></div>
        <div style="display:flex;align-items:flex-start;gap:10px"><span style="color:var(--gold);font-weight:700;min-width:16px">5.</span><span>Rewards are distributed automatically every day at 20:00 Berlin time (CEST/CET). Just keep your node running &mdash; no further action needed.</span></div>
      </div>

      <!-- Troubleshooting -->
      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">Troubleshooting</div>
      <table style="width:100%;border-collapse:collapse;margin-bottom:18px">
        <tr style="background:rgba(139,92,246,0.12)">
          <td style="font-size:0.6rem;color:var(--text);padding:8px;font-weight:700;width:32%">Symptom</td>
          <td style="font-size:0.6rem;color:var(--text);padding:8px;font-weight:700;width:32%">Likely cause</td>
          <td style="font-size:0.6rem;color:var(--text);padding:8px;font-weight:700">Solution</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Block height stays at 0</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">PRIMARY_NODE_URL not set or wrong</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Set PRIMARY_NODE_URL=https://aequitas.digital and redeploy. Also set SELF_URL to your own node's public URL.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">DATABASE_URL error on startup</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Wrong connection string or PostgreSQL unreachable</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Check format: postgres://user:pass@host:5432/dbname &mdash; make sure PostgreSQL is running and accessible</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">"no code at address" in logs</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">V7 contract not yet deployed in this EVM</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Normal on first start when RELAYER_ADDRESS is set &mdash; node auto-deploys V7. Wait a few seconds and check again.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">"NODE_OPERATOR_WALLET not set" in logs</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Missing environment variable</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Add NODE_OPERATOR_WALLET=0xYOUR_HUMAN_WALLET to your variables. Node runs fine without it but you won't receive rewards.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Railway shows "Application error"</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Build or startup failure</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Check Deploy Logs in Railway for the error message. Most common cause: DATABASE_URL missing or RELAYER_PRIVATE_KEY in wrong format (must start with 0x).</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Port 8080 not reachable (Docker)</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Firewall or cloud provider config</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Open TCP port 8080 inbound in your firewall or cloud security group settings.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Docker build fails with module error</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">No internet access during build</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Docker build needs outbound internet to download Go modules. Railway handles this automatically.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">⚠ P2P bootstrap unreachable (HTTP sync still works)</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">libp2p port 4001 firewalled (very common)</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Not critical — HTTP block sync is the primary mechanism and runs automatically. Add <code>-p 4001:4001</code> and <code>ufw allow 4001/tcp</code> to enable P2P as well.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Bootstrap snapshot failed / StateRoot mismatch</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">SNAPSHOT_TOKEN not set on primary, or BOOTSTRAP_SIGNER wrong</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Set BOOTSTRAP_SNAPSHOT_URL=https://aequitas.digital/api/snapshot, BOOTSTRAP_SIGNER=0x92cbedec9d348b4762cb9af99500ee6139c5b671, and SNAPSHOT_TOKEN (from network operator). Restart — node imports state automatically if DB is empty.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Node not in block explorer / no MERGE blocks</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Port 8080 not reachable from outside OR Step 5b not done</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">1) Open port 8080 inbound (<code>ufw allow 8080/tcp</code>). 2) Set SELF_URL=http://YOUR-IP:8080. 3) Complete Step 5b to register your signing key. Then the primary node syncs your blocks and MERGE events appear.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">MetaMask shows 0 AEQ or wrong balance after registration</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Stale network config in MetaMask (cached old RPC data)</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">MetaMask → Settings → Networks → delete all "Aequitas Chain" entries → re-add via the "+ ADD AEQUITAS NETWORK" button on this website. Balance will update immediately.</td>
        </tr>
        <tr>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">NODE_KEY generating new key on every restart</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">NODE_KEY env var not set</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">On first start, look for <code>SAVE THIS AS NODE_KEY ENVIRONMENT VAR: &lt;base64&gt;</code> in logs. Copy that value and add it as NODE_KEY environment variable. Restart once — P2P identity is now stable across all future restarts.</td>
        </tr>
      </table>

      <!-- Footer -->
      <div style="font-size:0.58rem;color:var(--gold);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid rgba(245,166,35,0.2);padding-bottom:6px">Questions / Feedback</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9">Open an issue on <a href="https://github.com/hanoi96international-gif/Aequitas" target="_blank" style="color:var(--purple)">GitHub</a> or reach the Aequitas team via the repository. Feedback on node setup, performance, and documentation gaps is especially welcome. Download this guide as a PDF in your selected language using the button above.</div>
    </div>
  </div>
</div>
</div>
<div id="net-protocol" class="stab-panel">
<div class="ps">
  <div class="section-label" data-i18n="proto-label">Aequitas V7 Protocol — Technical Documentation</div>

  <!-- V7 INTRO CARD -->
  <div class="idx" style="margin-bottom:12px;background:linear-gradient(135deg,rgba(139,92,246,0.08),rgba(6,182,212,0.04));border-color:rgba(139,92,246,0.25)">
    <div class="idx-title" data-i18n="v7-intro-title">What is AequitasV7?</div>
    <div style="font-size:0.65rem;color:var(--muted);line-height:1.9;margin-bottom:14px" data-i18n="v7-intro-text">AequitasV7 is the central smart contract of the Aequitas protocol. "V7" refers to the 7th major version of the fairness contract — the result of iterative design refinement focused on mathematical correctness, gas efficiency, and attack resistance. It is deployed on Aequitas Chain (Chain ID 1926) and handles every aspect of the protocol: human registration, ZK proof verification, balance management, wealth cap enforcement, UBI distribution, swap fees, and all governance parameters. No admin can upgrade or replace the contract — it is the immutable law of the Aequitas economy. The six mechanisms below do not work in isolation. They form a self-reinforcing system: demurrage feeds the UBI pool, wealth cap overflows add to UBI, swap fees distribute to all four pools simultaneously. Every economic activity strengthens redistribution.</div>
    <div style="display:grid;grid-template-columns:repeat(3,1fr);gap:8px">
      <div style="background:var(--card);border:1px solid var(--border);border-radius:6px;padding:10px;text-align:center">
        <div style="font-size:0.85rem;font-weight:700;color:var(--purple);font-family:var(--font-display)">6</div>
        <div style="font-size:0.55rem;color:var(--muted);margin-top:3px">Protocol Mechanisms</div>
      </div>
      <div style="background:var(--card);border:1px solid var(--border);border-radius:6px;padding:10px;text-align:center">
        <div style="font-size:0.85rem;font-weight:700;color:var(--neon);font-family:var(--font-display)">0</div>
        <div style="font-size:0.55rem;color:var(--muted);margin-top:3px">Admin Keys</div>
      </div>
      <div style="background:var(--card);border:1px solid var(--border);border-radius:6px;padding:10px;text-align:center">
        <div style="font-size:0.7rem;font-weight:700;color:var(--gold);font-family:var(--font-display)">immutable</div>
        <div style="font-size:0.55rem;color:var(--muted);margin-top:3px">Contract Code</div>
      </div>
    </div>
  </div>

  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="ca-title">Contract &amp; Network Addresses</div>
    <div style="font-size:0.65rem;color:var(--muted);line-height:1.9;margin-bottom:12px" data-i18n="ca-desc">AequitasV7 is the single source of truth for the entire Aequitas economy. Every AEQ balance, every human registration, every UBI payout, and every wealth cap enforcement is governed by this one immutable contract — deployed on Aequitas Chain, a custom EVM-compatible blockchain running a BlockDAG consensus engine. There is no admin key, no upgrade proxy, no governance vote that can change a single line of its logic. The code that runs today is the code that will run in ten years.<br><br>The BioVerifier contract receives Groth16 zero-knowledge proofs generated entirely on the user's Android device. It verifies mathematically on-chain in ~10 ms that a new registrant is a unique living human — without ever learning their name, identity, or biometric data. This is what makes gasless, investment-free registration possible: the proof is the only thing that ever leaves the device.<br><br>Together, these two contracts make possible something that has never existed in any currency system in history: a money supply whose rules — who gets it, how much exists, how it redistributes — cannot be altered by any person, company, or government. Ever.</div>
    <div class="hlbox" data-i18n="ca-text">Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier (Groth16 on-chain verifier): 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Main contract): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="poa-title">1. PROOF OF ALIVE — Inactive Balance Recovery</div>
    <div class="story" data-i18n="poa-text"><p>What happens to AEQ when people die or become permanently incapacitated? In Bitcoin and most cryptocurrencies, lost wallets mean permanently lost supply — millions of BTC are estimated to be inaccessible forever. Aequitas solves this through a multi-stage inactivity recovery system: if a wallet shows no activity for an extended period, its balance is gradually returned to the community through the UBI pool, ensuring the total effective supply remains meaningful.</p></div>
    <div class="hlbox" data-i18n="poa-box">Year 0–2: Normal usage — no restrictions<br>Year 2: Warning 1 sent — Guardian can respond on behalf<br>Year 2+60d: Warning 2 — escalating urgency<br>Year 2+120d: Warning 3 — final notice<br>Year 2+180d: AEQ moved to personal ESCROW (still recoverable)<br>Year 4: If still inactive — ESCROW released to UBI Pool</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="guard-title">2. GUARDIAN SYSTEM — Human Failsafe</div>
    <div class="story" data-i18n="guard-text"><p>What if someone is hospitalized, incarcerated, or otherwise unable to access their device for months? The Guardian system allows a trusted person — another verified human — to confirm that the wallet owner is still alive, preventing their AEQ from being moved to escrow. The Guardian has strictly zero financial access: they can only call a single function that resets the inactivity clock. They cannot move, spend, or access any funds under any circumstances.</p></div>
    <div class="hlbox" data-i18n="guard-box">1 Guardian per human · must be a verified human on Aequitas<br>Guardian can ONLY call confirmAlive() — zero transaction rights<br>Guardian CANNOT move funds, transfer AEQ, or access the wallet<br>Maximum 3 wards per Guardian (prevents centralization of trust)<br>7-day timelock on Guardian assignment (prevents forced assignment)<br>No circular guardian relationships allowed</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="dem-title">3. DEMURRAGE — Anti-Hoarding Mechanism</div>
    <div class="story" data-i18n="dem-text"><p>Demurrage is a holding cost on money — a negative interest rate that makes hoarding expensive and circulation attractive. It has historical precedent: the Wörgl experiment (Austria, 1932) used a demurrage currency and reduced local unemployment by 25% within one year. The Central Bank of Austria shut it down precisely because it worked too well and threatened the banking monopoly. The Chiemgauer (Germany, 2003) operates on the same principle and has circulated successfully for over 20 years. Aequitas implements continuous demurrage at 0.5% per month, applied only after a 3-month grace period of inactivity.</p></div>
    <div class="hlbox" data-i18n="dem-box">Rate: 0.5% per month after 3 months of inactivity (continuous, not stepped)<br>Clock resets automatically on any transfer, swap, or liquidity action<br>Decayed AEQ is redistributed to the four pools — never burned<br>14-day warning shown once · 7-day warning repeated on each active session</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="cap-title">4. WEALTH CAP — Mathematical Fairness Enforcement</div>
    <div class="hlbox" data-i18n="cap-box">Bootstrap cap: max(5,min(N,25))× current average AEQ balance<br>1–4 humans: 5× · grows +1× per new human · 25+ humans: 25× permanently<br>Applies to ALL addresses except the 4 protocol pool addresses<br>Excess AEQ instantly redistributed · No manual intervention required</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="ubi-title">5. UNIVERSAL BASIC INCOME — Daily Redistribution</div>
    <div class="hlbox" data-i18n="ubi-box">Sources of UBI Pool income:<br>· 20% of all swap fees from the AEQ↔tUSD AMM pool<br>· Overflow from wealth cap enforcement<br>· Demurrage charges from inactive accounts<br>· Inactive escrow released after 4 years<br><br>Distribution: Every 24 hours, the entire UBI pool balance is divided equally among all registered verified humans. The pool resets to zero and begins filling again immediately from ongoing protocol activity.</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="inf-title">6. NO ALGORITHMIC INFLATION — Fixed Supply Formula</div>
    <div class="hlbox" data-i18n="inf-box">The ONLY event that creates new AEQ: a new verified human registers.<br><br>Total Supply = Verified Humans × 1,000 AEQ<br><br>This is not a policy — it is enforced by the protocol. No admin can mint additional AEQ, no governance vote can change the issuance, no founder allocation was pre-mined. AEQ is the only cryptocurrency where the total supply is determined solely by the number of verified living humans.</div>
  </div>
  <div class="idx" style="margin-bottom:12px;background:linear-gradient(135deg,rgba(6,182,212,0.06),rgba(13,8,32,0.9));border-color:rgba(6,182,212,0.2)">
    <div class="idx-title" style="color:var(--teal)">Open Source Chain Logic</div>
    <div style="font-size:0.63rem;color:var(--muted);line-height:1.9">The Aequitas chain core — consensus engine, state machine, redistribution logic, wealth cap formula, and ZK proof verification — is written in Go. The redistribution algorithms (CalcGini, enforceWealthCap, DistributeUBIPool, settleDemurrage) are open for review.<br><br>Smart contract source code for AequitasV7 and BioVerifier is embedded in the chain binary and verifiable via the contract addresses above. Chain ID 1926, RPC: <span style="color:var(--teal);font-family:var(--font-mono)">https://aequitas.digital/rpc</span></div>
    <div style="margin-top:12px;padding:10px 14px;background:rgba(6,182,212,0.06);border:1px solid rgba(6,182,212,0.15);border-radius:6px;font-size:0.6rem;color:var(--teal);font-family:var(--font-mono)">
      /metrics — Prometheus endpoint (gini, humans, pools, block height)<br>
      /api/gini/history — Gini snapshots after each UBI distribution<br>
      /api/humans — All verified human balances (Lorenz curve source)<br>
      /api/wealth-cap — Live cap, multiplier, average balance
    </div>
  </div>
  <div class="idx" style="margin-bottom:12px;background:linear-gradient(135deg,rgba(139,92,246,0.06),rgba(13,8,32,0.9));border-color:rgba(139,92,246,0.2)">
    <div class="idx-title" style="color:var(--purple)">Node Decentralization Roadmap</div>
    <div style="font-size:0.63rem;color:var(--muted);line-height:1.9">Currently the network runs on 2 active nodes (Railway primary + Contabo VPS secondary) with MERGE events from both. Only registered humans can run validator nodes — this is a security requirement. Decentralization is a staged process:<br><br>
    <span style="color:var(--neon)">Phase 0 (now):</span> 2-node bootstrapping — Railway (aequitas.digital) + Contabo VPS. Trust established through code transparency. Registered humans can apply to join as validators.<br>
    <span style="color:var(--blue)">Phase 1 (100+ humans):</span> Open node join — any registered human can run a full node and earn validator rewards from the 40% pool.<br>
    <span style="color:var(--gold)">Phase 2 (1,000+ humans):</span> Minimum 10 independent node operators required. Node diversity enforced by smart contract.<br>
    <span style="color:var(--purple)">Phase 3 (10,000+ humans):</span> Fully decentralized BlockDAG. No single operator can censor or halt the chain.<br><br>
    The node operator guide (PDF) is available on the Network tab. Each new node operator earns from the 40% validator pool — the more nodes, the more resilient the network.</div>
  </div>
</div>
</div>

<script>
const PS = '/api'; // proof calls proxied via /api/prove on this node (avoids browser CORS)
const CID = '0x786';
const V7_CONTRACT = '0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78';
let waddr = '', proofData = null, curLang = 'en';

const T = {
en:{
  'logo-sub':'PROOF OF HUMANITY','live':'LIVE',
  'tab-register':'🔐 Register','tab-explorer':'🔍 Explorer','tab-humans':'👥 Humans','tab-index':'📊 Index','tab-network':'🌐 Network','tab-protocol':'📜 Protocol V7','tab-swap':'🔄 Swap',
  'reg-title':'🔐 Register as a Verified Human',
  'reg-sub':'Join the Aequitas network and receive your 1,000 AEQ Universal Basic Income grant. Registration is one-time, permanent, and completely gasless. No personal data is ever stored.',
  'app-title':'REGISTRATION VIA ANDROID APP',
  'app-text':'Proof of Humanity uses a physical 3-factor biometric system. <strong>Phase 1:</strong> R503 optical sensor scans all 10 fingerprints + MAX30102 PPG confirms liveness. <strong>Phase 2:</strong> Hand vein IR pattern (1 in 10⁷). <strong>Phase 3:</strong> Iris scan (1 in 10⁷⁸, gold standard). A Groth16 ZK proof commits all factors without revealing any biometric data. Your <strong style="color:var(--gold)">1,000 AEQ credited automatically</strong> upon verification.',
  's1t':'Biometric Scan','s1d':'AequitasBio scans all 10 fingerprints (R503 optical sensor) + MAX30102 PPG pulse confirms liveness. Phase 2: hand vein IR. Phase 3: iris. Raw data never leaves the device.',
  's2t':'ZK Proof Generation','s2d':'Groth16 ZK proof commits all biometric factors (fingers + vein + iris) into commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier bound to body, not phone.',
  's3t':'Connect Wallet','s3d':'The app opens MetaMask on this page · connect your Ethereum wallet · the proof is cryptographically bound to your wallet address',
  's4t':'1,000 AEQ Granted','s4d':'Registration confirmed on Aequitas BlockDAG within 6 seconds · 1,000 AEQ credited instantly · your identity is permanently recorded as a verified human',
  'priv-bar':'🔒 R503 All-10-Fingerprint · MAX30102 Liveness · Phase 2: Hand Vein IR · Phase 3: Iris (10⁷⁸) · Groth16 ZKP · Data never leaves device · One human · Forever',
  'conn-wallet':'CONNECTED WALLET','proof-recv':'⚡ ZK PROOF RECEIVED','proof-hint':'Connect wallet to register',
  'btn-conn':'🦊 CONNECT METAMASK','btn-reg':'🔐 REGISTER ON-CHAIN',
  'btn-web-reg':'🌐 REGISTER VIA BROWSER (WebAuthn)',
  'web-reg-warn':'⚠ Device-bound: This identity is tied to this device and browser. You cannot transfer it to another device. For permanent multi-device identity, use the Aequitas Android App instead.',
  'reg-log-hint':'// Open Aequitas Android App to generate your proof, then return here...',
  'reg-details':'Registration Details','k-network':'Network','k-chainid':'Chain ID','k-grant':'UBI Grant',
  'k-fee':'Gas Fee','free':'FREE — completely gasless','k-limit':'Registrations','k-limit-v':'Once per human · permanent · immutable',
  'k-bio':'Biometric Data','never-stored':'Never stored — stays on your device',
  'k-proof':'Proof System','k-conf':'Confirmation','k-conf-v':'Within 6 seconds (1 block)',
  'k-sybil':'Sybil Protection','k-sybil-v':'One identity per biometric · permanent lock',
  'live-stats':'Live Chain Statistics',
  's-height':'Block Height','s-height-sub':'New block every ~6s · BlockDAG · Parallel production',
  's-humans':'Verified Humans','s-humans-sub':'Biometric ZKP · One person, one wallet, forever',
  's-supply':'Total Supply','s-supply-sub':'Always = Humans × 1,000 AEQ',
  's-index':'Aequitas Index','s-index-sub':'0 = perfect equality · 100 = max inequality',
  's-uptime':'Uptime','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Proof of Humanity','ib-poh-t':'Every AEQ holder must cryptographically prove they are a unique living human. No bots, no corporations, no AI, no duplicates. Biometric data never leaves your device — only a mathematical proof is transmitted.',
  'ib-fair':'Radically Fair Distribution','ib-fair-t':'Every verified human receives exactly 1,000 AEQ upon registration — no more, no less. No pre-mine, no founder allocation, no investor rounds. Total supply always equals verified humans × 1,000.',
  'ib-dag':'BlockDAG Architecture','ib-dag-t':'Multiple blocks can be produced simultaneously and merged into the DAG. Higher throughput, lower latency, better fault tolerance than traditional linear blockchains.',
  'ib-gas':'Truly Gasless','ib-gas-t':'Registration and AEQ transfers cost absolutely nothing. No ETH, BNB, or MATIC required. No credit card, no bank account, no prior cryptocurrency needed.',
  'recent-blocks':'Recent Blocks','blocks-desc':'MERGE = multiple parents merged (BlockDAG). TX = registration transaction. Block time: ~6 seconds. Two nodes produce blocks in parallel.',
  'loading':'Loading blocks...','net-info':'Network Info','k-chain':'Chain Name','k-symbol':'Symbol','k-btime':'Block Time',
  'k-cons':'Consensus','k-nodes':'Active Nodes','k-storage':'Storage','add-mm':'🦊 ADD TO METAMASK','k-dec':'Decimals',
  'btn-add-mm':'+ ADD AEQUITAS NETWORK',
  'phil':'"Money exists because people exist.<br>Nothing more, nothing less."','phil-sub':'— THE AEQUITAS PRINCIPLE —',
  'humans-title':'Verified Humans on Aequitas Chain',
  'h-what':'What is a Verified Human?','h-what-t':'A Verified Human is a wallet address cryptographically proven to belong to a unique living human. Verification uses a 3-factor hardware system: R503 sensor scans all 10 fingerprints; MAX30102 PPG confirms live pulse; Phase 2 adds hand vein IR (1 in 10⁷); Phase 3 adds iris (1 in 10⁷⁸). Only a Groth16 ZK proof is transmitted. No biometric data ever leaves the device.',
  'h-zkp':'Zero-Knowledge Proof System','h-zkp-t':'Aequitas uses Groth16 on BN128 — same curve as Ethereum and Zcash. Proof: ~200 bytes. Verification: ~10ms. commitment = keccak256(iris‖vein‖fingers‖wallet). The nullifier is body-bound: losing your phone does not create a second identity. No biometric data is ever revealed or stored.',
  'h-sybil':'Sybil Attack Prevention','h-sybil-t':'Phase 1: all 10 fingerprints + MAX30102 liveness (PPG pulse, rejects casts/replays). Phase 2: hand vein IR — internal body feature, impossible to copy, unique 1 in 10⁷, different in identical twins. Phase 3: iris — 1 in 10⁷⁸, the global gold standard. Nullifier = keccak256(iris‖vein‖domain). One human, one wallet, forever.',
  'h-global':'Global Financial Inclusion','h-global-t':'No bank account, no credit card, no prior cryptocurrency required. Just an Android smartphone with a fingerprint or face sensor. Aequitas is designed to be accessible to every human on Earth.',
  'h-bio-hw':'Biometric Hardware Roadmap','h-bio-hw-t':'Phase 1 (active): R503 optical fingerprint scanner — all 10 fingers combined hash. MAX30102 PPG liveness. Phase 2 (planned): ESP32-CAM + 850 nm IR LED — hand vein imaging, 1 in 10⁷ uniqueness. Phase 3 (planned): IR iris module — 240+ degrees of freedom, 1 in 10⁷⁸, fully device-independent, identical twins differ.',
  'reg-humans':'Registered Humans','h-desc':'Every address below has been verified as a unique human through biometric ZKP. Each received exactly 1,000 AEQ. The registry is permanent, immutable, and on-chain.',
  'no-humans':'No humans registered yet.\n\nDownload the Aequitas Android App and be the first human on the chain!',
  'reg-stats':'Registry Stats','total-humans':'Total Humans',
  'idx-title':'Aequitas Index — Real-Time Economic Equality Score',
  'idx-desc':'The Aequitas Index is derived from the <strong style="color:var(--teal)">Gini coefficient</strong> — the international standard for measuring wealth inequality, adopted by the World Bank, OECD, and UN. It captures the complete balance distribution across every verified human simultaneously. <strong style="color:var(--neon)">0 = perfect equality</strong> (every wallet holds the same AEQ). <strong style="color:var(--red)">100 = total concentration</strong> (one wallet holds all AEQ). Bitcoin Gini ≈ 0.85 (Index 85) · South Africa (world record) ≈ 0.63 · Scandinavia ≈ 0.27 · Aequitas long-term target: Gini below 0.30 — comparable to the most equal developed economies, enforced by the wealth cap and redistribution pools.',
  'gini-what-title':'What is the Gini Coefficient?',
  'gini-what-text':'Developed by Italian statistician Corrado Gini (1912). Measures wealth distribution by comparing actual balances against a hypothetical perfectly equal baseline — visualized as the Lorenz curve. Scale: 0 (everyone holds the same) to 1 (one person holds everything). Used by World Bank, OECD, UN to compare countries. Reference values: Bitcoin ≈ 0.85 · South Africa (world record) ≈ 0.63 · USA ≈ 0.41 · Germany ≈ 0.31 · Scandinavia ≈ 0.27 · Aequitas long-term target: Gini below 0.30 — comparable to Scandinavian countries, enforced by wealth cap (bootstrap: 5×→25× per human).',
  'gini-calc-title':'How is the Aequitas Index calculated?',
  'gini-calc-text':'All AEQ balances of verified humans are collected. The formula computes the mean absolute difference between every possible pair of balances, normalized by population squared (n²) and the mean balance (x̄). Result 0–1 multiplied by 100 = Aequitas Index. Updated on-chain after every registration, monthly demurrage run, pool payout, and wealth cap event — via keeper calling updateGini().',
  'gini-why-title':'Why Gini — and not a simpler metric?',
  'gini-why-text':'A simple richest-vs-poorest ratio is easy to game: 10,000 wallets could show a low spread but 90% of AEQ concentrated in 100 hands — Gini detects this, a ratio does not. The coefficient captures the complete distribution across all verified humans in one auditable number. Aequitas publishes this on-chain — transparent, tamper-evident, globally verifiable. It is the primary signal for automatic phase transitions, wealth cap calibration, and redistribution intensity. No human can override the index reading or the mechanisms it triggers.',
  'curr-idx':'Current Index','bar-0':'0 — Perfect Equality','bar-100':'100 — Max Inequality','wcap-lbl':'Current Wealth Cap:','wcap-mult':'Multiplier:','wcap-avg':'Avg balance:',
  'gini':'Gini Coefficient','gini-desc':'0 = equal · 1 = unequal',
  'supply-desc':'Always = Humans × 1,000 AEQ',
  'phase':'Protocol Phase','phase-desc':'Auto-advances by human count',
  'humans-desc':'Biometrically verified unique humans',
  'pools-title':'Redistribution Pools',
  'pools-desc':'Every swap fee, demurrage charge, and wealth cap overflow is automatically split across four pools. No manual intervention — the protocol handles all redistribution through code alone. All pools pay out daily.',
  'vel-pool':'Validators Pool','vel-pool-desc':'40% of all fees → node operators who secure the network',
  'liq-pool':'Liquidity Pool','liq-pool-desc':'30% of all fees → liquidity providers, proportional to LP shares',
  'ubi-pool':'UBI Pool','ubi-pool-desc':'20% of all fees → all verified humans equally, every 24 hours',
  'treasury':'Treasury','treasury-desc':'10% of all fees → protocol development and maintenance',
  'phases-title':'Protocol Phases',
  'phases-desc':'The wealth cap uses a bootstrap multiplier during Phase 0: max(5, min(N, 25))× average balance. With 1–4 humans: 5× average. Each new human adds 1×. At 25+ humans: locks permanently at 25×. Phase 1+ maintains 25× fixed. All transitions trigger automatically by human count — no governance, no admin key.',
  'p0':'Bootstrap · &lt;100 humans · Wealth Cap: max(5,min(N,25))× average · Slides 5×→25× until 25th human · Currently active',
  'p1':'Growth · 100–10,000 humans · Wealth Cap: 25× average balance',
  'p2':'Stability · 10,000–1M humans · Wealth Cap: 25× average balance',
  'p3':'Maturity · 1M+ humans · Wealth Cap: 25× average balance',
  'wealth-cap-explain':'The Wealth Cap in Phase 0 (Bootstrap) uses max(5, min(N, 25))× average AEQ balance, where N = registered humans. 1–4 humans: cap = 5× average. Each new human adds 1×. 25+ humans: locked permanently at 25×. The cap always scales with the live average balance.',
  'demurrage-title':'Demurrage — Incentive to Circulate',
  'demurrage-desc':'Aequitas implements a demurrage mechanism inspired by historical complementary currencies. Idle AEQ balances slowly lose value to discourage hoarding and incentivize economic participation.',
  'dem-rate-k':'Decay Rate','dem-rate-v':'0.5% per month (continuous, not stepped)',
  'dem-grace-k':'Grace Period','dem-grace-v':'3 months of inactivity before decay begins',
  'dem-reset-k':'Clock Reset','dem-reset-v':'Any transfer, swap, or liquidity action resets the timer to zero',
  'dem-dest-k':'Decayed AEQ goes to','dem-dest-v':'Redistribution pools (40/30/20/10 split)',
  'dem-warn-k':'Warning System','dem-warn-v':'14-day notice (shown once) + 7-day repeated reminder at each login',
  'story-title':'The Story of Aequitas — Why This Exists',
  'story-text':'<p>The year is 2009. Satoshi Nakamoto releases Bitcoin. For the first time, value can transfer between any two people without a bank. A genuine revolution. But something goes wrong almost immediately.</p><p>Early miners accumulate millions of coins at almost zero cost. By 2021, the top 1% of Bitcoin addresses control over 90% of all Bitcoin. Bitcoin\'s estimated Gini coefficient exceeds 0.85 — higher than any country on Earth. The cryptocurrency that was supposed to democratize finance created the most extreme wealth concentration in human history.</p><p><span style="color:var(--gold)">Aequitas</span> — Latin for "fairness" and "equality" — was created to answer a single question: <em style="color:var(--gold)">"What would a cryptocurrency look like if designed from first principles to be fair to every human being?"</em></p><p>The answer is simple: <strong style="color:var(--text)">Money exists because people exist. Therefore, every person should have an equal share of money simply by virtue of being human.</strong></p><p>Aequitas implements this mathematically. Every verified human receives 1,000 AEQ. No mining, no staking, no early-adopter advantage. The wealth cap, demurrage, and redistribution pools ensure inequality cannot accumulate indefinitely. The protocol adjusts automatically as the network grows.</p><p>The Aequitas network launched in June 2026. Currently in Phase 0. The goal: demonstrate that money can be distributed fairly, Gini coefficient held below 0.30 (comparable to the most equal developed nations), and financial inclusion achieved at global scale — without any central authority.</p><p><em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em></p>',
  'nodes-title':'Active Nodes — Current Network Topology',
  'nodes-desc':'The Aequitas network currently operates on two geographically distributed nodes. Both participate in block production, state synchronization, and API serving. They communicate peer-to-peer via libp2p and synchronize block state via HTTP. Both share access to the same PostgreSQL database for persistent state. The network is designed to support additional nodes — any operator can join.',
  'node1':'Node 1 — Railway (Primary)','node1-desc':'Primary API · Block producer · UBI distribution · P2P bootstrap · PostgreSQL · RPC for MetaMask',
  'node2':'Node 2 — Render (Secondary)','node2-desc':'Secondary API · Block producer · P2P peer · HTTP sync · Shared PostgreSQL state',
  'run-node-title':'Run Your Own Node — Help Secure the Network',
  'run-node-desc':'Anyone can run an Aequitas node — no permission, no stake, no application required. Nodes participate in block production, validate the human registry, and synchronize the BlockDAG. Node operators earn a share of protocol fees via the Validators Pool (40% of all swap fees, distributed daily).',
  'bootstrap-title':'Connect a New Node','bootstrap-desc':'To run your own Aequitas node, set PRIMARY_NODE_URL=https://aequitas.digital in your environment. Your node registers automatically, syncs the full chain state, and begins participating in block production.',
  'tech-title':'Technical Specifications','mm-config':'MetaMask Configuration',
  'k-lang':'Language','k-src':'Source','evm-yes':'Yes — JSON-RPC /rpc · MetaMask compatible',
  'proto-label':'Aequitas V7 Protocol — Technical Documentation',
  'ca-title':'Contract Addresses',
  'ca-text':'Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Main): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7 is the single source of truth for the entire Aequitas economy. Every AEQ balance, every human registration, every UBI payout, and every wealth cap enforcement is governed by this one immutable contract — deployed on Aequitas Chain, a custom EVM-compatible blockchain running a BlockDAG consensus engine. There is no admin key, no upgrade proxy, no governance vote that can change a single line of its logic. The code that runs today is the code that will run in ten years.<br><br>The BioVerifier contract receives Groth16 zero-knowledge proofs generated entirely on the user\'s Android device. It verifies mathematically on-chain in ~10 ms that a new registrant is a unique living human — without ever learning their name, identity, or biometric data. This is what makes gasless, investment-free registration possible: the proof is the only thing that ever leaves the device.<br><br>Together, these two contracts make possible something that has never existed in any currency system in history: a money supply whose rules — who gets it, how much exists, how it redistributes — cannot be altered by any person, company, or government. Ever.',
  'poa-title':'1. PROOF OF ALIVE','poa-text':'<p>What happens to AEQ when people die or disappear? In Bitcoin, millions of BTC are permanently lost. In Aequitas, if someone is inactive for an extended period, their AEQ eventually returns to the community through the UBI pool.</p>',
  'poa-box':'Year 0-2: Normal usage<br>Year 2: Warning 1 — Guardian can respond<br>Year 2+60d: Warning 2<br>Year 2+120d: Warning 3<br>Year 2+180d: AEQ goes to PERSONAL ESCROW<br>Year 4: If still inactive — returns to UBI Pool',
  'guard-title':'2. GUARDIAN SYSTEM','guard-text':'<p>What if someone cannot access their device for months? A trusted Guardian — another verified human — can confirm they are still alive, without any transaction rights.</p>',
  'guard-box':'1 Guardian per human (must be another verified human)<br>Guardian can ONLY call confirmAlive() — zero transaction rights<br>Guardian CANNOT move funds or transfer AEQ<br>Max 3 wards · 7-day timelock · No circular relationships allowed',
  'dem-title':'3. DEMURRAGE — Anti-Hoarding Mechanism',
  'dem-box':'Rate: 0.5%/month after 3 months grace period<br>Clock resets on any transfer, swap, or liquidity action<br>Decayed AEQ redistributed to pools (not burned)',
  'dem-text':'<p>Historical precedent: The Wörgl experiment (Austria, 1932) used a demurrage currency and reduced unemployment by 25% in one year. The Chiemgauer (Germany, 2003) has operated successfully for over 20 years using a similar mechanism.</p>',
  'cap-title':'4. WEALTH CAP — Mathematical Fairness','cap-box':'Bootstrap cap: max(5,min(N,25))× current average AEQ balance<br>1–4 humans: 5× · +1× per human · 25+: 25× permanently<br>Excess AEQ instantly redistributed · No manual intervention',
  'ubi-title':'5. UNIVERSAL BASIC INCOME','ubi-box':'Sources: Swap fees (20%) · Wealth cap overflow · Demurrage · Inactive escrow<br><br>Daily: UBI Pool divided equally among all registered humans. Pool resets to zero after each distribution and refills continuously.',
  'inf-title':'6. NO ALGORITHMIC INFLATION','inf-box':'The ONLY event that creates new AEQ: a new verified human registers<br><br>Total Supply = Verified Humans × 1,000 AEQ — always, exactly.',
  'explore-title':'Explore Aequitas',
  'expl-score':'Equality Score','expl-score-d':'Live Gini coefficient · Aequitas Index · wealth distribution in real time',
  'expl-economy':'UBI &amp; Redistribution Pools','expl-economy-d':'Daily UBI countdown · 4 on-chain pools · demurrage · Protocol Phases',
  'expl-charts':'Charts &amp; History','expl-charts-d':'Gini history · Lorenz curve · Wealth Cap bootstrap slider · The story of Aequitas',
  'expl-v7':'Protocol V7 Docs','expl-v7-d':'AequitasV7 contract · 6 mechanisms · ZK proof · wealth cap · demurrage · immutable code',
  'expl-explorer':'Block Explorer','expl-explorer-d':'Live BlockDAG · click any block to see validator, hash, transactions, parent hashes',
    'btn-download-app':'DOWNLOAD AEQUITASBIO APP',
  'usp-headline':'For the first time in history — everyone starts equal',
  'usp-sub':'If you own an Android smartphone, you qualify. No bank, no crypto background, no investment needed.',
  'usp-c1-title':'0.00 Start Investment','usp-c1-desc':'Registration is completely gasless. No ETH, no MATIC, no credit card. The protocol pays all fees on your behalf.',
  'usp-c2-title':'1,000 AEQ for every human','usp-c2-desc':'Billionaire or subsistence farmer — everyone gets exactly 1,000 AEQ. Not more, not less. Equal start, guaranteed by math.',
  'usp-c3-title':'Accessible to all','usp-c3-desc':'No bank account, no credit card, no government ID. Registration uses an affordable biometric hardware kit (fingerprint scanner + pulse sensor, ~$15) — designed for global access.',
  'usp-c4-title':'Daily UBI forever','usp-c4-desc':'Once registered, you receive a daily share of UBI payouts automatically — every day, no action required.',
  'ubi-hero-title':'UNIVERSAL BASIC INCOME POOL','ubi-hero-sub':'Accumulating — next payout distributed equally to all verified humans in:',
  'ubi-hero-desc':'Split equally among all verified humans · paid every 24h · pool resets to zero after each payout · no minimum balance required',
  'ubi-bal-lbl':'current pool balance','ubi-how-fills':'HOW THE UBI POOL FILLS UP',
  'ubi-see-above':'see countdown above','ubi-timer-above':'⏰ countdown displayed above',
  'ubi-src-swap':'20% Swap Fees','ubi-src-swap-d':'Every AEQ↔tUSD swap contributes 20% of its 0.1% fee here. More trading activity = faster pool fill.',
  'ubi-src-dem':'variable Demurrage','ubi-src-dem-d':'Idle AEQ (3+ months inactive) decays at 0.5%/month. The decayed amount enters the 40/30/20/10 split — 20% goes to UBI.',
  'ubi-src-cap':'variable Wealth Cap Overflow','ubi-src-cap-d':'Wallets exceeding 25+ average balance have the excess confiscated instantly. 20% flows to UBI immediately.',
  'ubi-pool-desc':'20% of swap fees + demurrage + wealth cap overflow → divided equally among all verified humans every 24 hours. Even with zero trading, demurrage and wealth cap ensure the pool always fills.',
  'pool-t-timer':'Accumulates — no timer',
  'pools4-header':'ALL FOUR REDISTRIBUTION POOLS',
  'swap-title':'🔄 Swap AEQ ↔ tUSD',
  'swap-sub':'Exchange AEQ for tUSD (a simulated test-dollar) through the native liquidity pool. A 0.1% fee applies only to swaps — ordinary AEQ transfers between people remain completely free.',
  'swap-faucet-desc':'Claim 1,000 tUSD once to pair with your AEQ — for your first liquidity deposit.',
  'swap-btn-faucet':'CLAIM TEST tUSD (once)','swap-btn-conn':'🦊 CONNECT METAMASK','swap-btn-go':'🔄 SWAP',
  'swap-rate-lbl':'Rate','swap-fee-bps':'Fee','swap-out-lbl':'You receive approx.','swap-impact-lbl':'Price Impact',
  'swap-depth-lbl':'Pool Depth','swap-pool-aeq':'Pool AEQ','swap-pool-tusd':'Pool tUSD','swap-pool-price':'Price',
  'swap-pool-title':'AMM Liquidity Pool','swap-no-liquidity':'No liquidity yet','swap-details-hdr':'Swap Details',
  'swap-lp-title':'Your LP Position','swap-lp-share':'Pool Share','swap-lp-withdrawable':'Withdrawable',
  'swap-lp-youget':'You get approx.','swap-lp-pct-label':'of pool','swap-lps':'LP Shares',
  'swap-your-aeq':'Your AEQ','swap-your-tusd':'Your tUSD',
  'swap-addliq-title':'Add Liquidity','swap-addliq-desc':'Deposit AEQ and tUSD to earn 30% of all swap fees proportional to your share.',
  'swap-btn-addliq':'+ ADD LIQUIDITY','swap-btn-removeliq':'− REMOVE LIQUIDITY',
  'swap-fee-est':'Estimated fee','swap-log-hint':'// Connect wallet to swap AEQ ↔ tUSD...',
  'swap-ubi':'20% UBI','swap-validators':'40% Validators','swap-treasury':'10% Treasury',
  'amm-title':'How the AMM works','amm-text':'Automated Market Maker using the x·y=k formula. Price is determined by pool ratio. Deeper pools = lower price impact per swap.',
  'pools-addr-title':'Pool Contract Addresses','swap-pools-addr-title':'Pool Addresses','swap-priv-bar':'🔒 Non-custodial · AMM x·y=k · 0.1% fee · Instant settlement · No slippage protection needed at small sizes',
  'v7-intro-title':'What is AequitasV7?',
  'v7-intro-text':'AequitasV7 is the single source of truth for the entire Aequitas economy. Every AEQ balance, every human registration, every UBI payout, and every wealth cap enforcement is governed by this one immutable contract.',
'expl-network':'Network &amp; Nodes','expl-network-d':'Node topology · run your own node · technical specs · Chain ID 1926'
,'swap-sell-label':'Sell','swap-receive-label':'Receive',
},
de:{
  'logo-sub':'MENSCHLICHKEITSNACHWEIS','live':'LIVE',
  'tab-register':'🔐 Registrieren','tab-explorer':'🔍 Explorer','tab-humans':'👥 Menschen','tab-index':'📊 Index','tab-network':'🌐 Netzwerk','tab-protocol':'📜 Protokoll V7','tab-swap':'🔄 Tauschen',
  'reg-title':'🔐 Als verifizierter Mensch registrieren',
  'reg-sub':'Tritt dem Aequitas-Netzwerk bei und erhalte dein Universelles Grundeinkommen von 1.000 AEQ. Einmalig, permanent und vollständig gebührenfrei. Keine persönlichen Daten werden jemals gespeichert.',
  'app-title':'REGISTRIERUNG NUR ÜBER ANDROID-APP',
  'app-text':'Der Menschlichkeitsnachweis verwendet ein physisches 3-Faktor-Biometrie-System. <strong>Phase 1:</strong> R503-Optischer-Sensor scannt alle 10 Fingerabdrücke + MAX30102 PPG bestätigt lebenden Puls. <strong>Phase 2:</strong> Handvenen-IR-Kamera (1 von 10⁷). <strong>Phase 3:</strong> Iris-Scan (1 von 10⁷⁸, Goldstandard). Ein Groth16-ZK-Beweis bindet alle Faktoren ohne biometrische Daten preiszugeben. Deine <strong style="color:var(--gold)">1.000 AEQ werden automatisch gutgeschrieben</strong>.',
  's1t':'Biometrischer Scan','s1d':'AequitasBio scannt alle 10 Fingerabdrücke (R503) + MAX30102 PPG-Puls bestätigt Lebenderkennung. Phase 2: Handvenen-IR. Phase 3: Iris. Rohdaten verlassen nie das Gerät.',
  's2t':'ZK-Beweis-Erzeugung','s2d':'Groth16-ZK-Beweis bindet alle biometrischen Faktoren: commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier ist körpergebunden — Telefonverlust erzeugt keine zweite Identität.',
  's3t':'Wallet verbinden','s3d':'Die App öffnet MetaMask auf dieser Seite · verbinde deine Ethereum-Wallet · der Beweis ist kryptografisch an deine Wallet-Adresse gebunden',
  's4t':'1.000 AEQ gutgeschrieben','s4d':'Registrierung auf Aequitas BlockDAG innerhalb von 6 Sekunden bestätigt · 1.000 AEQ sofort gutgeschrieben · deine Identität ist dauerhaft als verifizierter Mensch gespeichert',
  'priv-bar':'🔒 R503 Alle-10-Finger · MAX30102 Lebenderkennung · Phase 2: Handvenen-IR · Phase 3: Iris (10⁷⁸) · Groth16 ZKP · Daten verlassen nie Gerät · Ein Mensch · Für immer',
  'conn-wallet':'VERBUNDENE WALLET','proof-recv':'⚡ ZK-BEWEIS EMPFANGEN','proof-hint':'Wallet verbinden um zu registrieren',
  'btn-conn':'🦊 METAMASK VERBINDEN','btn-reg':'🔐 ON-CHAIN REGISTRIEREN',
  'btn-web-reg':'🌐 IM BROWSER REGISTRIEREN (WebAuthn)',
  'web-reg-warn':'⚠ Gerätgebunden: Diese Identität ist an dieses Gerät und diesen Browser gebunden. Sie kann nicht auf ein anderes Gerät übertragen werden. Für dauerhafte Geräteunabhängigkeit nutze die Aequitas Android App.',
  'reg-log-hint':'// Öffne die Aequitas Android App um deinen Beweis zu erstellen, dann kehre hierher zurück...',
  'reg-details':'Registrierungsdetails','k-network':'Netzwerk','k-chainid':'Chain-ID','k-grant':'UBI-Zuteilung',
  'k-fee':'Gasgebühr','free':'KOSTENLOS — vollständig gebührenfrei','k-limit':'Registrierungen','k-limit-v':'Einmal pro Mensch · permanent · unveränderlich',
  'k-bio':'Biometrische Daten','never-stored':'Nie gespeichert — bleibt auf deinem Gerät',
  'k-proof':'Beweissystem','k-conf':'Bestätigung','k-conf-v':'Innerhalb von 6 Sekunden (1 Block)',
  'k-sybil':'Sybil-Schutz','k-sybil-v':'Eine Identität pro Biometrie · dauerhaft gesperrt',
  'live-stats':'Live-Chain-Statistiken',
  's-height':'Blockhöhe','s-height-sub':'Neuer Block alle ~6s · BlockDAG · Parallele Produktion',
  's-humans':'Verifizierte Menschen','s-humans-sub':'Biometrisches ZKP · Eine Person, eine Wallet, für immer',
  's-supply':'Gesamtmenge','s-supply-sub':'Immer = Menschen × 1.000 AEQ',
  's-index':'Aequitas-Index','s-index-sub':'0 = perfekte Gleichheit · 100 = maximale Ungleichheit',
  's-uptime':'Laufzeit','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Menschlichkeitsnachweis','ib-poh-t':'Jeder AEQ-Inhaber muss kryptografisch beweisen dass er ein einzigartiger lebender Mensch ist. Keine Bots, keine Unternehmen, keine KI. Biometrische Daten verlassen nie dein Gerät.',
  'ib-fair':'Radikal gerechte Verteilung','ib-fair-t':'Jeder verifizierte Mensch erhält genau 1.000 AEQ bei der Registrierung. Kein Pre-Mining, keine Gründerzuteilung. Gesamtmenge entspricht immer Verifizierte Menschen × 1.000.',
  'ib-dag':'BlockDAG-Architektur','ib-dag-t':'Mehrere Blöcke können gleichzeitig produziert und zusammengeführt werden. Höherer Durchsatz, geringere Latenz als lineare Blockchains.',
  'ib-gas':'Wirklich gebührenfrei','ib-gas-t':'Registrierung und AEQ-Transfers kosten absolut nichts. Kein ETH, BNB oder MATIC erforderlich. Kein Bankkonto, keine Kreditkarte nötig.',
  'recent-blocks':'Aktuelle Blöcke','blocks-desc':'MERGE = mehrere Eltern zusammengeführt (BlockDAG). TX = Registrierungstransaktion. Blockzeit: ~6 Sekunden.',
  'loading':'Blöcke werden geladen...','net-info':'Netzwerkinformationen','k-chain':'Chain-Name','k-symbol':'Symbol','k-btime':'Blockzeit',
  'k-cons':'Konsens','k-nodes':'Aktive Nodes','k-storage':'Speicher','add-mm':'🦊 ZU METAMASK HINZUFÜGEN','k-dec':'Dezimalstellen',
  'btn-add-mm':'+ AEQUITAS-NETZWERK HINZUFÜGEN',
  'phil':'"Geld existiert weil Menschen existieren.<br>Nichts mehr, nichts weniger."','phil-sub':'— DAS AEQUITAS-PRINZIP —',
  'humans-title':'Verifizierte Menschen auf der Aequitas Chain',
  'h-what':'Was ist ein verifizierter Mensch?','h-what-t':'Ein verifizierter Mensch ist eine Wallet-Adresse, die kryptografisch bewiesen einem einzigartigen lebenden Menschen gehört. Die Verifikation erfolgt über ein physisches 3-Faktor-System: R503 scannt alle 10 Fingerabdrücke; MAX30102 bestätigt Puls; Phase 2: Handvenen-IR (1 von 10⁷); Phase 3: Iris (1 von 10⁷⁸). Nur ein Groth16-ZK-Beweis wird übertragen.',
  'h-zkp':'Zero-Knowledge-Beweissystem','h-zkp-t':'Aequitas verwendet Groth16 auf BN128 — dieselbe Kurve wie Ethereum und Zcash. ~200 Bytes, ~10ms. commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier körpergebunden — Telefonverlust erzeugt keine zweite Identität.',
  'h-sybil':'Sybil-Angriff-Prävention','h-sybil-t':'Phase 1: alle 10 Fingerabdrücke + MAX30102-Lebenderkennung (verhindert Gips-Abdrücke/Replays). Phase 2: Handvenen-IR — inneres Körpermerkmal, 1 von 10⁷, bei eineiigen Zwillingen verschieden. Phase 3: Iris — 1 von 10⁷⁸, Goldstandard. Nullifier = keccak256(iris‖vein‖domain). Ein Mensch, eine Wallet, für immer.',
  'h-global':'Globale finanzielle Inklusion','h-global-t':'Kein Bankkonto, keine Kreditkarte, keine Kryptowährung erforderlich. Nur ein Android-Smartphone mit Fingerabdruck- oder Gesichtssensor.',
  'h-bio-hw':'Biometrische Hardware-Roadmap','h-bio-hw-t':'Phase 1 (aktiv): R503 Fingerabdruckscanner — alle 10 Finger. MAX30102 PPG Lebenderkennung. Phase 2 (geplant): ESP32-CAM + 850 nm IR-LED — Handvenen, 1 von 10⁷. Phase 3 (geplant): IR-Iris-Modul — 240+ Freiheitsgrade, 1 von 10⁷⁸, geräteunabhängig, Zwillinge unterschiedlich.',
  'reg-humans':'Registrierte Menschen','h-desc':'Jede Adresse wurde als einzigartiger Mensch durch biometrisches ZKP verifiziert. Jeder erhielt genau 1.000 AEQ. Dauerhaft, unveränderlich, on-chain.',
  'no-humans':'Noch keine Menschen registriert.\n\nLade die Aequitas Android App herunter und sei der erste Mensch auf der Chain!',
  'reg-stats':'Registrierungsstatistiken','total-humans':'Gesamtmenschen',
  'idx-title':'Aequitas-Index — Echtzeit-Wirtschaftsgleichheits-Score',
  'idx-desc':'Der Aequitas-Index wird aus dem <strong style="color:var(--teal)">Gini-Koeffizienten</strong> abgeleitet — dem internationalen Standard zur Messung wirtschaftlicher Ungleichheit, genutzt von Weltbank, OECD und UN. Er erfasst die vollständige Bilanzverteilung aller verifizierten Menschen gleichzeitig. <strong style="color:var(--neon)">0 = perfekte Gleichheit</strong> (jede Wallet hält gleich viel AEQ). <strong style="color:var(--red)">100 = totale Konzentration</strong> (eine Wallet hält alles). Bitcoin-Gini ≈ 0,85 (Index 85) · Südafrika (Weltrekord) ≈ 0,63 · Skandinavien ≈ 0,27 · Aequitas-Langzeitziel: Gini unter 0,30 (Index unter 30) — vergleichbar mit den gleichheitsstärksten Industrieländern, automatisch durchgesetzt durch den Vermögensobergrenze-Mechanismus.',
  'gini-what-title':'Was ist der Gini-Koeffizient?',
  'gini-what-text':'Entwickelt vom italienischen Statistiker Corrado Gini (1912). Misst die Vermögensverteilung durch Vergleich mit einer perfekt gleichen Verteilung — visualisiert als Lorenz-Kurve. Skala: 0 (alle halten gleich viel) bis 1 (eine Person hält alles). Genutzt von Weltbank, OECD, UN. Referenzwerte: Bitcoin ≈ 0,85 · Südafrika (Weltrekord) ≈ 0,63 · USA ≈ 0,41 · Deutschland ≈ 0,31 · Schweden ≈ 0,27 · Aequitas-Langzeitziel: Gini unter 0,30 — vergleichbar mit Skandinavien und Deutschland, durchgesetzt durch den Vermögensdeckel (Bootstrap: gleitender Deckel 5×→25× pro Mensch).',
  'gini-calc-title':'Wie wird der Aequitas-Index berechnet?',
  'gini-calc-text':'Alle AEQ-Salden verifizierter Menschen werden erfasst. Die Formel berechnet die mittlere absolute Differenz zwischen allen Saldo-Paaren, normiert durch Bevölkerungsgröße im Quadrat (n²) und Durchschnittssaldo (x̄). Ergebnis 0–1 multipliziert mit 100 = Aequitas-Index. Aktualisiert On-Chain nach jeder Registrierung, jedem monatlichen Demurrage-Lauf, jeder Pool-Ausschüttung und jedem Vermögensobergrenze-Ereignis — via Keeper-Aufruf updateGini().',
  'gini-why-title':'Warum Gini — und nicht eine einfachere Kennzahl?',
  'gini-why-text':'Ein "Reich-Arm-Verhältnis" ist leicht manipulierbar: 10.000 Wallets könnten eine geringe Spanne zeigen, aber 90% des AEQ in 100 Händen halten — Gini erkennt das, ein Verhältnis nicht. Der Koeffizient erfasst die vollständige Verteilung aller verifizierten Menschen in einer einzigen prüfbaren Zahl. Aequitas veröffentlicht diese On-Chain — transparent, manipulationssicher, weltweit verifizierbar. Sie ist das Hauptsignal für automatische Phasenübergänge, Vermögensobergrenze-Kalibrierung und Umverteilungsintensität. Kein Mensch kann den Index-Wert oder die von ihm ausgelösten Mechanismen überschreiben.',
  'curr-idx':'Aktueller Index','bar-0':'0 — Perfekte Gleichheit','bar-100':'100 — Max. Ungleichheit',
  'gini':'Gini-Koeffizient','gini-desc':'0 = gleich · 1 = ungleich',
  'supply-desc':'Immer = Menschen × 1.000 AEQ',
  'phase':'Protokollphase','phase-desc':'Automatisch nach Menschenanzahl',
  'humans-desc':'Biometrisch verifizierte einzigartige Menschen',
  'pools-title':'Umverteilungspools',
  'pools-desc':'Jede Swap-Gebühr, Demurrage-Belastung und Vermögensobergrenze-Überschuss wird automatisch auf vier Pools aufgeteilt. Keine manuelle Eingriffe. Alle Pools zahlen täglich aus.',
  'vel-pool':'Validatoren-Pool','vel-pool-desc':'40% aller Gebühren → Node-Betreiber die das Netzwerk sichern',
  'liq-pool':'Liquiditäts-Pool','liq-pool-desc':'30% aller Gebühren → Liquiditätsanbieter, proportional zu LP-Anteilen',
  'ubi-pool':'UBI-Pool','ubi-pool-desc':'20% aller Gebühren → alle verifizierten Menschen gleichmäßig, alle 24 Stunden',
  'treasury':'Schatzkammer','treasury-desc':'10% aller Gebühren → Protokollentwicklung und -wartung',
  'phases-title':'Protokollphasen',
  'phases-desc':'In Phase 0 verwendet die Vermögensobergrenze einen Bootstrap-Multiplikator: max(5, min(N, 25))× Durchschnittsguthaben. Mit 1–4 Menschen: 5× Durchschnitt. Jeder neue Mensch erhöht um 1×. Ab 25+ Menschen: dauerhaft auf 25× fixiert. Phase 1+ behält 25× fest. Alle Übergänge erfolgen automatisch — kein Governance-Vote, kein Admin-Key.',
  'p0':'Bootstrap · &lt;100 Menschen · Vermögensobergrenze: max(5,min(N,25))× Durchschnitt · Gleitet 5×→25× bis zum 25. Menschen · Derzeit aktiv',
  'p1':'Wachstum · 100–10.000 Menschen · Vermögensobergrenze: 25× Durchschnittsguthaben',
  'p2':'Stabilität · 10.000–1M Menschen · Vermögensobergrenze: 25× Durchschnittsguthaben',
  'p3':'Reife · 1M+ Menschen · Vermögensobergrenze: 25× Durchschnittsguthaben',
  'wealth-cap-explain':'Die Vermögensobergrenze in Phase 0 (Bootstrap) verwendet max(5, min(N, 25))× Durchschnittsguthaben, wobei N = registrierte Menschen. 1–4 Menschen: 5× Durchschnitt. Jeder neue Mensch erhöht um 1×. Ab 25+ Menschen: dauerhaft 25×. Die Obergrenze skaliert stets mit dem Live-Durchschnittsguthaben.',
  'demurrage-title':'Demurrage — Anreiz zum Zirkulieren',
  'demurrage-desc':'Aequitas implementiert einen Demurrage-Mechanismus inspiriert von historischen Komplementärwährungen. Inaktive AEQ-Guthaben verlieren langsam an Wert um Hortung zu entmutigen.',
  'dem-rate-k':'Verfallsrate','dem-rate-v':'0,5% pro Monat (kontinuierlich, nicht gestuft)',
  'dem-grace-k':'Schonfrist','dem-grace-v':'3 Monate Inaktivität bevor der Verfall beginnt',
  'dem-reset-k':'Uhr-Reset','dem-reset-v':'Jede Überweisung, Swap oder Liquiditätsaktion setzt den Timer zurück',
  'dem-dest-k':'Verfallenes AEQ geht an','dem-dest-v':'Umverteilungspools (40/30/20/10 Aufteilung)',
  'dem-warn-k':'Warnsystem','dem-warn-v':'14-Tage-Hinweis (einmal) + 7-Tage-Wiederholung bei jedem Login',
  'story-title':'Die Geschichte von Aequitas — Warum es das gibt',
  'story-text':'<p>Das Jahr ist 2009. Satoshi Nakamoto veröffentlicht Bitcoin. Zum ersten Mal kann Wert zwischen zwei Menschen ohne eine Bank übertragen werden. Eine echte Revolution. Aber fast sofort läuft etwas schief.</p><p>Frühe Miner häufen Millionen von Coins zu fast null Kosten an. Bis 2021 kontrollieren die obersten 1% der Bitcoin-Adressen über 90% aller Bitcoin. Bitcoins geschätzter Gini-Koeffizient übersteigt 0,85 — höher als in jedem Land auf der Erde.</p><p><span style="color:var(--gold)">Aequitas</span> — Lateinisch für "Fairness" und "Gleichheit" — wurde geschaffen um eine einzige Frage zu beantworten: <em style="color:var(--gold)">"Wie würde eine Kryptowährung aussehen die von Grund auf fair für jeden Menschen konzipiert wurde?"</em></p><p>Die Antwort ist einfach: <strong style="color:var(--text)">Geld existiert weil Menschen existieren. Daher sollte jeder Mensch einfach durch seine Existenz einen gleichen Anteil am Geld haben.</strong></p><p>Aequitas setzt dies mathematisch um. Jeder verifizierte Mensch erhält 1.000 AEQ. Kein Mining, kein Staking, kein Frühanwender-Vorteil. Die Vermögensobergrenze, Demurrage und Umverteilungspools stellen sicher dass sich Ungleichheit nicht unbegrenzt anhäufen kann.</p><p><em style="color:var(--gold)">"Geld existiert weil Menschen existieren. Nichts mehr, nichts weniger."</em></p>',
  'nodes-title':'Aktive Nodes — Aktuelle Netzwerktopologie',
  'nodes-desc':'Das Aequitas-Netzwerk betreibt derzeit zwei geografisch verteilte Nodes. Beide nehmen an Blockproduktion, Statussynchronisation und API-Bereitstellung teil. Sie kommunizieren per libp2p und synchronisieren Blockzustände via HTTP. Das Netzwerk ist für zusätzliche Nodes ausgelegt — jeder Betreiber kann beitreten.',
  'node1':'Node 1 — Railway (Primär)','node1-desc':'Primärer API-Server · Blockproduzent · UBI-Verteilung · P2P-Bootstrap · PostgreSQL · RPC für MetaMask',
  'node2':'Node 2 — Render (Sekundär)','node2-desc':'Sekundärer API-Server · Blockproduzent · P2P-Peer · HTTP-Sync · Geteilter PostgreSQL-Status',
  'run-node-title':'Eigenen Node betreiben — Das Netzwerk sichern',
  'run-node-desc':'Jeder kann einen Aequitas-Node betreiben — keine Genehmigung, kein Stake, keine Bewerbung erforderlich. Nodes nehmen an der Blockproduktion teil und validieren die Menschenregistrierung. Node-Betreiber erhalten täglich einen Anteil der Protokollgebühren über den Validators-Pool (40% aller Swap-Gebühren).',
  'bootstrap-title':'Neuen Node verbinden','bootstrap-desc':'Um einen eigenen Aequitas-Node zu betreiben, setze die PRIMARY_NODE_URL=https://aequitas.digital in deiner Umgebung. Dein Node synchronisiert automatisch den vollständigen Chain-Zustand und beginnt mit der Blockproduktion.',
  'tech-title':'Technische Spezifikationen','mm-config':'MetaMask-Konfiguration',
  'k-lang':'Sprache','k-src':'Quellcode','evm-yes':'Ja — JSON-RPC /rpc · MetaMask-kompatibel',
  'proto-label':'Aequitas V7 Protokoll — Technische Dokumentation',
  'ca-title':'Contract- & Netzwerk-Adressen','ca-text':'Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier (Groth16 On-Chain-Verifier): 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Haupt-Contract): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7 ist die einzige Wahrheitsquelle der gesamten Aequitas-Wirtschaft. Jedes AEQ-Guthaben, jede Menschenregistrierung, jede UBI-Auszahlung und jede Durchsetzung der Vermögensobergrenze wird durch diesen einen unveränderlichen Contract geregelt — deployed auf der Aequitas Chain, einer maßgeschneiderten EVM-kompatiblen Blockchain mit BlockDAG-Konsens. Es gibt keinen Admin-Schlüssel, keinen Upgrade-Proxy, keine Governance-Abstimmung die eine einzige Zeile seiner Logik ändern könnte. Der Code der heute läuft ist der Code der in zehn Jahren läuft.<br><br>Der BioVerifier-Contract empfängt Groth16-Zero-Knowledge-Beweise die vollständig auf dem Android-Gerät des Nutzers erzeugt werden. Er verifiziert mathematisch on-chain in ~10 ms dass ein neuer Registrierungskandidat ein einzigartiger lebender Mensch ist — ohne jemals seinen Namen, seine Identität oder seine biometrischen Daten zu erfahren. Das ist es was die gasfreie, investitionsfreie Registrierung möglich macht: Der Beweis ist das Einzige was das Gerät je verlässt.<br><br>Zusammen machen diese zwei Contracts etwas möglich das in keinem Währungssystem der Geschichte je existiert hat: eine Geldmenge deren Regeln — wie viel existiert, wer es bekommt, wie es umverteilt wird — von keiner Person, keinem Unternehmen und keiner Regierung je geändert werden können. Niemals.',
  'ib-poh':'Menschlichkeitsnachweis','ib-poh-t':'Jeder AEQ-Inhaber muss kryptographisch beweisen dass er ein einzigartiger lebender Mensch ist. Keine Bots, keine Unternehmen, keine KI, keine Duplikate. Biometrische Daten verlassen niemals dein Gerät — nur ein mathematischer Einzigartigkeitsbeweis wird übertragen. Das bedeutet: AEQ ist die erste Währung die ausschließlich menschlich ist.',
  'ib-fair':'Radikal faire Verteilung','ib-fair-t':'Jeder verifizierte Mensch erhält bei der Registrierung genau 1.000 AEQ — nicht mehr, nicht weniger. Kein Pre-Mining, keine Gründer-Zuteilung, keine Investorenrunden. Die Gesamtmenge ist immer und exakt gleich der Anzahl verifizierter Menschen multipliziert mit 1.000. Dies wird mathematisch erzwungen, nicht durch Richtlinien.',
  'ib-dag':'BlockDAG-Architektur','ib-dag-t':'Im Gegensatz zu traditionellen Blockchains wo nur ein Block pro Höhe existieren kann, verwendet Aequitas eine DAG-Struktur. Mehrere Blöcke können gleichzeitig von verschiedenen Nodes produziert und später in den DAG zusammengeführt werden. Dies ermöglicht höheren Durchsatz, niedrigere Latenz und eliminiert Einzelknoten-Engpässe. Merge-Ereignisse werden im Explorer mit einem speziellen Badge markiert.',
  'ib-gas':'Wirklich gebührenfrei','ib-gas-t':'Alle Registrierungen und AEQ-Übertragungen kosten absolut nichts. Kein ETH, BNB oder MATIC erforderlich. Keine Kreditkarte, kein Bankkonto, keine vorherige Kryptowährung nötig. Der Relayer übernimmt alle Transaktionskosten. Wenn du ein Mensch mit einem Smartphone bist, kannst du teilnehmen — unabhängig von deiner wirtschaftlichen Situation.',
  'h-what':'Was ist ein verifizierter Mensch?','h-what-t':'Ein verifizierter Mensch ist eine Wallet-Adresse, die kryptographisch bewiesen einem einzigartigen lebenden Menschen gehört. Die Verifikation erfolgt über ein physisches 3-Faktor-Hardware-System: R503-Sensor scannt alle 10 Fingerabdrücke; MAX30102-PPG bestätigt den lebenden Puls; Phase 2 ergänzt Handvenen-IR (1 von 10⁷); Phase 3 ergänzt Iris (1 von 10⁷⁸). Nur ein Groth16-ZK-Beweis wird übertragen. Keine biometrischen Daten verlassen das Gerät.',
  'h-zkp':'Zero-Knowledge-Proof-System','h-zkp-t':'Aequitas verwendet Groth16 auf BN128 — dieselbe Kurve wie Ethereum und Zcash. Beweisgröße: ~200 Byte. Verifikationszeit: ~10ms. commitment = keccak256(iris‖vein‖fingers‖wallet). Der Nullifier ist körpergebunden: Telefonverlust erzeugt keine zweite Identität. Keine biometrischen Daten werden je gespeichert oder übertragen.',
  'h-sybil':'Sybil-Angriff-Prävention','h-sybil-t':'Phase 1: alle 10 Fingerabdrücke + MAX30102-Lebenderkennung (PPG-Puls, verhindert Abdruckfälschungen). Phase 2: Handvenen-IR — inneres Körpermerkmal, nicht kopierbar, einzigartig 1 von 10⁷, bei eineiigen Zwillingen verschieden. Phase 3: Iris — 1 von 10⁷⁸, der weltweite Goldstandard. Nullifier = keccak256(iris‖vein‖domain). Ein Mensch, eine Wallet, für immer.',
  'h-global':'Globale finanzielle Inklusion','h-global-t':'1,4 Milliarden Erwachsene weltweit haben kein Bankkonto. Aequitas benötigt nur ein Android-Smartphone mit einem Fingerabdruck- oder Gesichtssensor — ein Gerät das über 3 Milliarden Menschen bereits besitzen. Kein Bankkonto, keine Kreditkarte, keine vorherige Kryptowährung, kein Personalausweis. Einfach Mensch zu sein reicht aus.',
  'h-bio-hw':'Biometrische Hardware-Roadmap','h-bio-hw-t':'Phase 1 (aktiv): R503 optischer Fingerabdruckscanner — alle 10 Finger kombinierter Hash. MAX30102 PPG Lebenderkennung. Phase 2 (geplant): ESP32-CAM + 850 nm IR-LED — Handvenen-Bildgebung, 1 von 10⁷ Einzigartigkeit. Phase 3 (geplant): IR-Iris-Modul — 240+ Freiheitsgrade, 1 von 10⁷⁸, vollständig geräteunabhängig, eineiige Zwillinge unterschiedlich.',
  'poa-title':'1. LEBENSNACHWEIS — Inaktive Guthaben-Rückgewinnung','poa-text':'<p>Was passiert mit AEQ wenn Menschen sterben oder dauerhaft handlungsunfähig werden? Bei Bitcoin und den meisten Kryptowährungen bedeuten verlorene Wallets dauerhaft verlorene Menge. Aequitas löst dies durch ein mehrstufiges Inaktivitäts-Rückgewinnungssystem: Wenn eine Wallet über einen längeren Zeitraum keine Aktivität zeigt, wird ihr Guthaben schrittweise über den UBI-Pool zur Gemeinschaft zurückgeführt.</p>',
  'poa-box':'Jahr 0–2: Normale Nutzung — keine Einschränkungen<br>Jahr 2: Warnung 1 — Guardian kann im Namen antworten<br>Jahr 2+60T: Warnung 2 — steigende Dringlichkeit<br>Jahr 2+120T: Warnung 3 — letzte Benachrichtigung<br>Jahr 2+180T: AEQ in persönliches TREUHANDKONTO verschoben (noch rückgewinnbar)<br>Jahr 4: Bei weiter Inaktivität — Treuhand an UBI-Pool freigegeben',
  'guard-title':'2. GUARDIAN-SYSTEM — Menschliche Absicherung','guard-text':'<p>Was wenn jemand hospitalisiert, inhaftiert oder anderweitig monatelang nicht in der Lage ist auf sein Gerät zuzugreifen? Das Guardian-System erlaubt einer vertrauenswürdigen Person — einem anderen verifizierten Menschen — zu bestätigen dass der Wallet-Inhaber noch lebt, wodurch verhindert wird dass sein AEQ ins Treuhandkonto verschoben wird. Der Guardian hat strikt null finanziellen Zugang: Er kann nur eine einzige Funktion aufrufen die den Inaktivitätstimer zurücksetzt. Er kann unter keinen Umständen Gelder verschieben, ausgeben oder darauf zugreifen.</p>',
  'guard-box':'1 Guardian pro Mensch · muss ein verifizierter Mensch auf Aequitas sein<br>Guardian kann NUR confirmAlive() aufrufen — null Transaktionsrechte<br>Guardian KANN KEINE Gelder verschieben, AEQ übertragen oder auf die Wallet zugreifen<br>Maximal 3 Schutzbefohlene pro Guardian (verhindert Zentralisierung des Vertrauens)<br>7-Tage-Zeitsperre bei Guardian-Zuweisung (verhindert erzwungene Zuweisung)<br>Keine zirkulären Guardian-Beziehungen erlaubt',
  'dem-title':'3. DEMURRAGE — Anti-Hortungs-Mechanismus',
  'dem-box':'Rate: 0,5% pro Monat nach 3 Monaten Inaktivität (kontinuierlich, nicht gestuft)<br>Uhr setzt sich automatisch zurück bei jeder Überweisung, Swap oder Liquiditätsaktion<br>Verfallenes AEQ wird an die vier Pools umverteilt — niemals vernichtet<br>14-Tage-Warnung einmalig angezeigt · 7-Tage-Warnung bei jeder aktiven Sitzung wiederholt',
  'dem-text':'<p>Demurrage ist ein Haltungskosten auf Geld — ein negativer Zinssatz der Horten teuer und Zirkulation attraktiv macht. Historisches Beispiel: Das Wörgl-Experiment (Österreich, 1932) verwendete eine Demurrage-Währung und reduzierte die lokale Arbeitslosigkeit innerhalb eines Jahres um 25%. Die Österreichische Zentralbank stellte es genau deshalb ein weil es zu gut funktionierte. Der Chiemgauer (Deutschland, 2003) arbeitet nach demselben Prinzip und zirkuliert seit über 20 Jahren erfolgreich.</p>',
  'cap-title':'4. VERMÖGENSOBERGRENZE — Mathematische Fairness-Durchsetzung','cap-box':'Bootstrap-Deckel: max(5,min(N,25))× aktuelles Durchschnittsguthaben<br>1–4 Menschen: 5× · +1× pro Mensch · 25+: dauerhaft 25×<br>Gilt für ALLE Adressen außer den 4 Protokoll-Pool-Adressen<br>Überschuss-AEQ sofort weitergeleitet · Keine manuellen Eingriffe',
  'ubi-title':'5. UNIVERSELLES GRUNDEINKOMMEN — Tägliche Umverteilung','ubi-box':'Quellen des UBI-Pool-Einkommens:<br>· 20% aller Swap-Gebühren aus dem AEQ↔tUSD AMM-Pool<br>· Überschuss aus der Vermögensobergrenze-Durchsetzung<br>· Demurrage-Gebühren von inaktiven Konten<br>· Inaktive Treuhand nach 4 Jahren freigegeben<br><br>Ausschüttung: Alle 24 Stunden wird der gesamte UBI-Pool-Saldo gleichmäßig unter allen registrierten verifizierten Menschen aufgeteilt. Der Pool setzt sich auf null zurück und beginnt sofort wieder aus der laufenden Protokollaktivität aufzufüllen.',
  'inf-title':'6. KEINE ALGORITHMISCHE INFLATION — Feste Mengenformel','inf-box':'Das EINZIGE Ereignis das neues AEQ schafft: ein neuer verifizierter Mensch registriert sich.<br><br>Gesamtmenge = Verifizierte Menschen × 1.000 AEQ<br><br>Dies ist keine Richtlinie — es wird durch das Protokoll erzwungen. Kein Admin kann zusätzliches AEQ prägen, kein Governance-Votum kann die Ausgabe ändern, keine Gründer-Zuteilung wurde vorab gemint. AEQ ist die einzige Kryptowährung bei der die Gesamtmenge ausschließlich durch die Anzahl verifizierter lebender Menschen bestimmt wird.',
  'btn-download-app':'AEQUITASBIO APP HERUNTERLADEN',
  'swap-title':'🔄 Tausche AEQ ↔ tUSD',
  'swap-sub':'Tausche AEQ gegen tUSD (ein simulierter Test-Dollar) über den nativen Liquiditäts-Pool. 0,1% Gebühr gilt nur für Swaps — gewöhnliche AEQ-Transfers zwischen Menschen bleiben vollständig kostenlos.',
  'swap-priv-bar':'🔒 Nur 0,1% Swap-Gebühr · AEQ-zu-AEQ-Transfers kostenlos · tUSD ist eine Testwährung ohne realen Wert',
  'swap-your-aeq':'Dein AEQ','swap-your-tusd':'Dein tUSD',
  'swap-aeq-to-tusd':'AEQ → tUSD','swap-tusd-to-aeq':'tUSD → AEQ',
  'swap-fee-est':'Protokollgebühr (0,1%)','swap-details-hdr':'Swap-Details',
  'swap-out-lbl':'Du erhältst (ca.)','swap-impact-lbl':'Preisauswirkung','swap-rate-lbl':'Wechselkurs',
  'swap-btn-conn':'🦊 METAMASK VERBINDEN','swap-btn-go':'🔄 TAUSCHEN',
  'swap-log-hint':'// Wallet verbinden um zu tauschen...',
  'swap-no-liquidity':'Noch kein tUSD?','swap-faucet-desc':'Registrierte Menschen können einmalig Test-tUSD beanspruchen',
  'swap-btn-faucet':'💧 TEST-tUSD BEANSPRUCHEN',
  'swap-addliq-title':'Liquidität bereitstellen','swap-addliq-desc':'Sei der Erste der einzahlt — dein Verhältnis legt den Startpreis fest.',
  'swap-btn-addliq':'💧 LIQUIDITÄT HINZUFÜGEN',
  'swap-lp-title':'Deine LP-Position','swap-lp-share':'Pool-Anteil','swap-lp-withdrawable':'Auszahlbar',
  'swap-lp-pct-label':'% deiner Position','swap-lp-youget':'Du erhältst','swap-btn-removeliq':'🔥 LIQUIDITÄT ENTFERNEN',
  'swap-pool-title':'AEQ / tUSD — Pool-Status',
  'swap-pool-aeq':'AEQ-Reserve','swap-pool-tusd':'tUSD-Reserve','swap-pool-price':'Spot-Preis',
  'swap-depth-lbl':'Pool-Zusammensetzung',
  'amm-title':'x × y = k — Konstantprodukt-AMM',
  'amm-text':'Wenn du AEQ gegen tUSD tauschst, wächst die AEQ-Reserve und die tUSD-Reserve schrumpft — ihr Produkt bleibt immer gleich k. Jeder Swap bewegt den Preis. Größere Swaps relativ zur Pool-Größe führen zu größerer Preisauswirkung. Die 0,1% Gebühr wird vor Anwendung der Formel abgezogen — so verdient der Pool an jedem Trade.',
  'swap-fee-bps':'Swap-Gebühr','swap-fee-split':'Gebührenverteilung','swap-fee-split-v':'40% Validatoren / 30% LPs / 20% UBI / 10% Schatzkammer',
  'swap-pools-addr-title':'Tokenomics-Pool-Adressen',
  'swap-validators':'Validatoren (40%)','swap-lps':'Liquiditätsanbieter (30%)','swap-ubi':'UBI-Pool (20%)','swap-treasury':'Schatzkammer (10%)',
  'ubi-hero-title':'UNIVERSELLES GRUNDEINKOMMEN — UBI-POOL',
  'ubi-hero-sub':'Akkumuliert — nächste Ausschüttung gleichmäßig an alle verifizierten Menschen in:',
  'ubi-bal-lbl':'aktuelles Pool-Guthaben',
  'ubi-hero-desc':'Gleichmäßig unter allen verifizierten Menschen aufgeteilt · alle 24h ausgezahlt · Pool setzt auf null zurück · kein Mindestguthaben nötig',
  'ubi-how-fills':'Wie der UBI-Pool sich füllt',
  'ubi-src-swap':'Swap-Gebühren','ubi-src-swap-d':'Jeder AEQ↔tUSD-Swap trägt 20% seiner 0,1% Gebühr bei. Mehr Handelsaktivität = schnelleres Auffüllen.',
  'ubi-src-dem':'Demurrage','ubi-src-dem-d':'Inaktives AEQ (3+ Monate) verfällt mit 0,5%/Monat. Der verfallene Betrag geht in die 40/30/20/10-Aufteilung — 20% an UBI.',
  'ubi-src-cap':'Vermögensobergrenze-Überschuss','ubi-src-cap-d':'Wallets die den Vermögensdeckel (max(5,min(N,25))× Durchschnitt) überschreiten werden sofort gekappt. 20% fließt direkt an UBI.',
  'pools4-header':'Alle vier Umverteilungs-Pools',
  'vel-pool-desc':'Node-Betreiber die Blöcke produzieren, ZK-Registrierungen validieren und den BlockDAG sichern. Täglich ausgezahlt proportional zur Blockproduktion.',
  'liq-pool-desc':'Anbieter von AEQ/tUSD-Liquidität erhalten 30% aller Gebühren proportional zu ihrem LP-Anteil. Tiefere Liquidität = geringere Preisauswirkung für alle Nutzer.',
  'ubi-pool-desc':'20% der Swap-Gebühren + Demurrage + Vermögensobergrenze-Überschuss → gleichmäßig unter allen verifizierten Menschen alle 24 Stunden. Auch ohne Trading füllt sich der Pool durch Demurrage und Vermögensobergrenze.',
  'treasury-desc':'Protokollentwicklung, Infrastruktur, Sicherheitsprüfungen und zukünftige Upgrades. Vollständige On-Chain-Transparenz.',
  'ubi-see-above':'siehe Countdown oben','ubi-timer-above':'⏰ Countdown oben angezeigt','pool-t-timer':'Akkumuliert — kein Timer',
  'usp-headline':'Zum ersten Mal in der Geschichte — alle starten gleich',
  'usp-sub':'Ein Android-Smartphone genügt. Kein Bankkonto, keine Kreditkarte, keine Vorkenntnisse, keine Investition.',
  'usp-c1-title':'0,00 € Startinvestition','usp-c1-desc':'Die Registrierung ist vollständig gebührenfrei. Kein ETH, kein BNB, keine Kreditkarte. Das Protokoll übernimmt alle Transaktionskosten — du startest bei null.',
  'usp-c2-title':'1.000 AEQ für jeden Menschen','usp-c2-desc':'Millionär oder Subsistenzlandwirt — jeder erhält exakt 1.000 AEQ. Nicht mehr, nicht weniger. Gleicher Start, mathematisch garantiert.',
  'usp-c3-title':'Für alle zugänglich','usp-c3-desc':'Kein Bankkonto, keine Kreditkarte, kein Personalausweis. Die Registrierung erfolgt über ein günstiges Biometrie-Hardware-Kit (Fingerabdruckscanner + Pulssensor, ~15 €) — designed für globalen Zugang.',
  'usp-c4-title':'Täglich UBI empfangen','usp-c4-desc':'Nach der Registrierung erhältst du automatisch täglich einen Anteil der UBI-Ausschüttung — jeden Tag, ohne Aktion, solange du AEQ hältst.',
  'v7-intro-title':'Was ist AequitasV7?',
  'v7-intro-text':'AequitasV7 ist der zentrale Smart Contract des Aequitas-Protokolls. "V7" steht für die 7. Hauptversion des Fairness-Contracts — das Ergebnis iterativer Designverbesserung. Er ist unveränderlich auf der Aequitas Chain (Chain ID 1926) deployed und regelt jeden Aspekt des Protokolls: Menschenregistrierung, ZK-Beweisverifizierung, Guthabenverwaltung, Vermögensobergrenze, UBI-Ausschüttung, Swap-Gebühren und alle Governance-Parameter. Kein Admin kann den Contract upgraden oder ersetzen — er ist das unveränderliche Gesetz der Aequitas-Wirtschaft.',
  'explore-title':'Aequitas entdecken',
  'expl-score':'Gleichheits-Score','expl-score-d':'Live-Gini-Koeffizient · Aequitas-Index · Vermögensverteilung in Echtzeit',
  'expl-economy':'UBI &amp; Umverteilungspools','expl-economy-d':'Täglicher UBI-Countdown · 4 On-Chain-Pools · Demurrage · Protokollphasen',
  'expl-charts':'Diagramme &amp; Verlauf','expl-charts-d':'Gini-Verlauf · Lorenz-Kurve · Vermögensobergrenze-Bootstrap-Slider · Die Geschichte von Aequitas',
  'expl-v7':'Protokoll V7 Dokumentation','expl-v7-d':'AequitasV7-Contract · 6 Mechanismen · ZK-Beweis · Vermögensobergrenze · Demurrage · unveränderlicher Code',
  'expl-explorer':'Block-Explorer','expl-explorer-d':'Live-BlockDAG · Block anklicken um Validator, Hash, Transaktionen, Eltern-Hashes zu sehen',
  'swap-sell-label':'Verkaufen','swap-receive-label':'Erhalten',
  'expl-network':'Netzwerk &amp; Nodes','expl-network-d':'Node-Topologie · eigenen Node betreiben · technische Spezifikationen · Chain-ID 1926'
},
es:{
  'logo-sub':'PRUEBA DE HUMANIDAD','live':'EN VIVO',
  'tab-register':'🔐 Registrar','tab-explorer':'🔍 Explorador','tab-humans':'👥 Humanos','tab-index':'📊 Índice','tab-network':'🌐 Red','tab-protocol':'📜 Protocolo V7','tab-swap':'🔄 Intercambiar',
  'reg-title':'🔐 Regístrate como Humano Verificado',
  'reg-sub':'Únete a la red Aequitas y recibe tu subsidio de Renta Básica Universal de 1,000 AEQ. Único, permanente y completamente gratuito. Ningún dato personal es almacenado.',
  'app-title':'REGISTRO SOLO VÍA APP ANDROID',
  'app-text':'La Prueba de Humanidad usa un sistema biométrico físico de 3 factores. Fase 1: sensor óptico R503 escanea las 10 huellas + MAX30102 PPG confirma pulso vivo. Fase 2: cámara IR de venas de la mano (1 en 10⁷ de unicidad). Fase 3: escaneo de iris — estándar de oro, 1 en 10⁷⁸, totalmente independiente del dispositivo. Una prueba ZK Groth16 compromete todos los factores sin revelar datos biométricos. Tus 1,000 AEQ se acreditan automáticamente al verificar.',
  's1t':'Escaneo Biométrico','s1d':'AequitasBio escanea las 10 huellas (sensor óptico R503) + pulso MAX30102 PPG confirma vivacidad. Fase 2: IR de venas de mano. Fase 3: iris. Los datos brutos nunca salen del dispositivo.',
  's2t':'Generación de Prueba ZK','s2d':'La prueba ZK Groth16 compromete todos los factores biométricos: commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier vinculado al cuerpo, no al teléfono — perder el dispositivo no puede crear una segunda identidad.',
  's3t':'Conectar Wallet','s3d':'La app abre MetaMask en esta página · conecta tu wallet Ethereum · la prueba está criptográficamente vinculada a tu dirección',
  's4t':'1,000 AEQ Acreditados','s4d':'Registro confirmado en el BlockDAG de Aequitas en 6 segundos · 1,000 AEQ acreditados instantáneamente · tu identidad queda permanentemente registrada',
  'priv-bar':'🔒 R503 10 Huellas · MAX30102 Vivacidad · Fase 2: IR de Venas de Mano · Fase 3: Iris (10⁷⁸) · Groth16 ZKP · Datos nunca salen del dispositivo · Un humano · Para siempre',
  'conn-wallet':'WALLET CONECTADA','proof-recv':'⚡ PRUEBA ZK RECIBIDA','proof-hint':'Conecta wallet para registrar',
  'btn-conn':'🦊 CONECTAR METAMASK','btn-reg':'🔐 REGISTRAR ON-CHAIN',
  'btn-web-reg':'🌐 REGISTRAR VIA NAVEGADOR (WebAuthn)',
  'web-reg-warn':'⚠ Vinculado al dispositivo: Esta identidad está vinculada a este dispositivo y navegador. No puedes transferirla a otro dispositivo. Para identidad permanente multidispositivo, usa la App Android de Aequitas.',
  'reg-log-hint':'// Abre la App Android Aequitas para generar tu prueba, luego regresa aquí...',
  'reg-details':'Detalles del Registro','k-network':'Red','k-chainid':'ID de Cadena','k-grant':'Subsidio UBI',
  'k-fee':'Tarifa de Gas','free':'GRATIS — completamente sin gas','k-limit':'Registros','k-limit-v':'Una vez · permanente · inmutable',
  'k-bio':'Datos Biométricos','never-stored':'Nunca almacenados — permanece en tu dispositivo',
  'k-proof':'Sistema de Prueba','k-conf':'Confirmación','k-conf-v':'En 6 segundos (1 bloque)',
  'k-sybil':'Protección Sybil','k-sybil-v':'Una identidad por biometría · bloqueo permanente',
  'live-stats':'Estadísticas de Cadena en Vivo',
  's-height':'Altura de Bloque','s-height-sub':'Nuevo bloque cada ~6s · BlockDAG · Producción paralela',
  's-humans':'Humanos Verificados','s-humans-sub':'ZKP biométrico · Una persona, una wallet, siempre',
  's-supply':'Suministro Total','s-supply-sub':'Siempre = Humanos × 1,000 AEQ',
  's-index':'Índice Aequitas','s-index-sub':'0 = igualdad perfecta · 100 = desigualdad máxima',
  's-uptime':'Tiempo Activo','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Prueba de Humanidad','ib-poh-t':'Cada titular de AEQ debe probar criptográficamente que es un humano único vivo. Sin bots, sin corporaciones, sin IA. Los datos biométricos nunca salen de tu dispositivo.',
  'ib-fair':'Distribución Radicalmente Justa','ib-fair-t':'Cada humano verificado recibe exactamente 1,000 AEQ al registrarse. Sin pre-minado, sin asignación a fundadores. El suministro total siempre equivale a humanos verificados × 1,000.',
  'ib-dag':'Arquitectura BlockDAG','ib-dag-t':'Múltiples bloques pueden producirse simultáneamente y fusionarse. Mayor rendimiento, menor latencia que las blockchains lineales.',
  'ib-gas':'Verdaderamente Sin Gas','ib-gas-t':'El registro y las transferencias no cuestan nada. No se necesita ETH, BNB ni MATIC. Sin cuenta bancaria ni tarjeta de crédito.',
  'recent-blocks':'Bloques Recientes','blocks-desc':'MERGE = múltiples padres fusionados (BlockDAG). TX = transacción de registro. Tiempo de bloque: ~6 segundos.',
  'loading':'Cargando bloques...','net-info':'Información de Red','k-chain':'Nombre de Cadena','k-symbol':'Símbolo','k-btime':'Tiempo de Bloque',
  'k-cons':'Consenso','k-nodes':'Nodos Activos','k-storage':'Almacenamiento','add-mm':'🦊 AGREGAR A METAMASK','k-dec':'Decimales',
  'btn-add-mm':'+ AGREGAR RED AEQUITAS',
  'phil':'"El dinero existe porque las personas existen.<br>Nada más, nada menos."','phil-sub':'— EL PRINCIPIO AEQUITAS —',
  'humans-title':'Humanos Verificados en Aequitas Chain',
  'h-what':'¿Qué es un Humano Verificado?','h-what-t':'Un Humano Verificado es una dirección wallet demostrada criptográficamente que pertenece a un humano único vivo. La verificación usa un sistema de hardware de 3 factores: R503 escanea las 10 huellas; MAX30102 PPG confirma pulso vivo; Fase 2: IR de venas de mano (1 en 10⁷); Fase 3: iris (1 en 10⁷⁸). Solo se transmite una prueba ZK Groth16. Ningún dato biométrico abandona el dispositivo.',
  'h-zkp':'Sistema de Prueba ZK','h-zkp-t':'Aequitas usa Groth16 en BN128 — misma curva que Ethereum y Zcash. ~200 bytes, ~10ms. commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier vinculado al cuerpo: perder tu teléfono no crea una segunda identidad. No se almacena ningún dato biométrico.',
  'h-sybil':'Prevención de Ataques Sybil','h-sybil-t':'Fase 1: las 10 huellas + vivacidad MAX30102 (pulso PPG, rechaza moldes/repeticiones). Fase 2: IR de venas de mano — característica corporal interna, imposible de copiar, 1 en 10⁷, diferente en gemelos idénticos. Fase 3: iris — 1 en 10⁷⁸, estándar de oro global. Nullifier = keccak256(iris‖vein‖domain). Un humano, una wallet, para siempre.',
  'h-global':'Inclusión Financiera Global','h-global-t':'Sin cuenta bancaria, tarjeta de crédito ni criptomoneda previa. Solo un smartphone Android con sensor biométrico.',
  'h-bio-hw':'Hoja de Ruta del Hardware Biométrico','h-bio-hw-t':'Fase 1 (activa): escáner de huellas óptico R503 — hash combinado de los 10 dedos. Vivacidad MAX30102 PPG. Fase 2 (planificada): ESP32-CAM + LED IR de 850 nm — imagen de venas de mano, 1 en 10⁷ de unicidad. Fase 3 (planificada): módulo de iris IR — 240+ grados de libertad, 1 en 10⁷⁸, totalmente independiente del dispositivo, gemelos idénticos difieren.',
  'reg-humans':'Humanos Registrados','h-desc':'Cada dirección verificada como humano único mediante ZKP biométrico. Cada uno recibió exactamente 1,000 AEQ. Permanente, inmutable, on-chain.',
  'no-humans':'No hay humanos registrados aún.\n\n¡Descarga la App Android Aequitas y sé el primero!',
  'reg-stats':'Estadísticas del Registro','total-humans':'Total de Humanos',
  'idx-title':'Índice Aequitas — Puntuación de Igualdad Económica en Tiempo Real',
  'idx-desc':'El Índice Aequitas mide la desigualdad económica de todos los humanos verificados en tiempo real. Se calcula desde el coeficiente Gini de la distribución de saldos on-chain. 0 = igualdad perfecta. 100 = desigualdad máxima.',
  'curr-idx':'Índice Actual','bar-0':'0 — Igualdad Perfecta','bar-100':'100 — Máx. Desigualdad',
  'gini':'Coeficiente Gini','gini-desc':'0 = igual · 1 = desigual',
  'supply-desc':'Siempre = Humanos × 1,000 AEQ',
  'phase':'Fase del Protocolo','phase-desc':'Avanza automáticamente por recuento humano',
  'humans-desc':'Humanos únicos verificados biométricamente',
  'pools-title':'Pools de Redistribución',
  'pools-desc':'Cada tarifa de swap, cargo de demurrage y desbordamiento del límite de riqueza se divide automáticamente entre cuatro pools. Sin intervención manual. Todos los pools pagan diariamente.',
  'vel-pool':'Pool Validadores','vel-pool-desc':'40% de todas las tarifas → operadores de nodos que aseguran la red',
  'liq-pool':'Pool Liquidez','liq-pool-desc':'30% de todas las tarifas → proveedores de liquidez, proporcional a participaciones LP',
  'ubi-pool':'Pool UBI','ubi-pool-desc':'20% de todas las tarifas → todos los humanos verificados por igual, cada 24 horas',
  'treasury':'Tesorería','treasury-desc':'10% de todas las tarifas → desarrollo y mantenimiento del protocolo',
  'phases-title':'Fases del Protocolo',
  'phases-desc':'En Fase 0, el límite de riqueza usa un multiplicador de arranque: max(5, min(N, 25))× saldo promedio. Con 1–4 humanos: 5× promedio. Cada nuevo humano añade 1×. A 25+ humanos: fijado permanentemente en 25×. Fase 1+ mantiene 25× fijo. Todas las transiciones son automáticas — sin voto de gobernanza, sin clave de administrador.',
  'p0':'Bootstrap · &lt;100 humanos · Límite de Riqueza: max(5,min(N,25))× promedio · Deslizamiento 5×→25× hasta el 25.º humano · Actualmente activo',
  'p1':'Crecimiento · 100–10,000 humanos · Límite de Riqueza: 25× saldo promedio',
  'p2':'Estabilidad · 10,000–1M humanos · Límite de Riqueza: 25× saldo promedio',
  'p3':'Madurez · 1M+ humanos · Límite de Riqueza: 25× saldo promedio',
  'wealth-cap-explain':'El Límite de Riqueza en Fase 0 (Bootstrap) usa max(5, min(N, 25))× saldo promedio, donde N = humanos registrados. 1–4 humanos: 5× promedio. Cada nuevo humano añade 1×. 25+ humanos: bloqueado en 25× permanentemente. El límite siempre se escala con el saldo promedio actual.',
  'btn-download-app':'DESCARGAR APP AEQUITASBIO',
  'swap-title':'🔄 Intercambiar AEQ ↔ tUSD','swap-sub':'Intercambia AEQ por tUSD (un dólar de prueba simulado) a través del pool de liquidez nativo. Se aplica una comisión del 0,1% solo a los intercambios — las transferencias ordinarias de AEQ entre personas permanecen completamente gratuitas.',
  'swap-priv-bar':'🔒 Solo 0,1% de comisión de swap · Transferencias AEQ a AEQ gratuitas · tUSD es una moneda de prueba sin valor real',
  'swap-your-aeq':'Tu AEQ','swap-your-tusd':'Tu tUSD','swap-aeq-to-tusd':'AEQ → tUSD','swap-tusd-to-aeq':'tUSD → AEQ',
  'swap-fee-est':'Comisión de protocolo (0,1%)','swap-details-hdr':'Detalles del Swap',
  'swap-out-lbl':'Recibes (est.)','swap-impact-lbl':'Impacto en precio','swap-rate-lbl':'Tipo de cambio',
  'swap-depth-lbl':'Composición del Pool','amm-title':'x × y = k — AMM de Producto Constante',
  'amm-text':'Cuando intercambias AEQ por tUSD, la reserva de AEQ crece y la de tUSD decrece — su producto siempre permanece igual a k. Swaps más grandes causan mayor impacto en precio. La comisión del 0,1% se descuenta antes de aplicar la fórmula.',
  'swap-btn-conn':'🦊 CONECTAR METAMASK','swap-btn-go':'🔄 INTERCAMBIAR',
  'swap-log-hint':'// Conecta tu wallet para intercambiar...',
  'swap-no-liquidity':'¿Sin tUSD todavía?','swap-faucet-desc':'Los humanos registrados pueden reclamar tUSD de prueba una vez','swap-btn-faucet':'💧 RECLAMAR tUSD DE PRUEBA',
  'swap-addliq-title':'Proporcionar Liquidez','swap-addliq-desc':'Sé el primero en depositar — tu ratio establece el precio inicial.','swap-btn-addliq':'💧 AGREGAR LIQUIDEZ',
  'swap-lp-title':'Tu Posición LP','swap-lp-share':'Participación del Pool','swap-lp-withdrawable':'Retirable',
  'swap-lp-pct-label':'% de tu posición','swap-lp-youget':'Recibirás','swap-btn-removeliq':'🔥 RETIRAR LIQUIDEZ',
  'swap-pool-title':'AEQ / tUSD — Estado del Pool',
  'swap-pool-aeq':'Reserva AEQ','swap-pool-tusd':'Reserva tUSD','swap-pool-price':'Precio Spot',
  'swap-fee-bps':'Comisión de Swap','swap-fee-split':'Distribución de comisiones','swap-fee-split-v':'40% Validadores / 30% LPs / 20% UBI / 10% Tesorería',
  'swap-pools-addr-title':'Direcciones de Pools Tokenomics',
  'swap-validators':'Validadores (40%)','swap-lps':'Proveedores de Liquidez (30%)','swap-ubi':'Pool UBI (20%)','swap-treasury':'Tesorería (10%)',
  'ubi-hero-title':'RENTA BÁSICA UNIVERSAL — POOL UBI',
  'ubi-hero-sub':'Acumulando — próximo pago distribuido por igual a todos los humanos verificados en:',
  'ubi-bal-lbl':'saldo actual del pool','ubi-hero-desc':'Dividido por igual entre todos · pagado cada 24h · el pool se reinicia a cero · sin saldo mínimo requerido',
  'ubi-how-fills':'Cómo se llena el Pool UBI',
  'ubi-src-swap':'Comisiones de Swap','ubi-src-swap-d':'Cada swap AEQ↔tUSD contribuye el 20% de su comisión de 0,1%. Más actividad = llenado más rápido.',
  'ubi-src-dem':'Demurrage','ubi-src-dem-d':'AEQ inactivo (3+ meses) decae al 0,5%/mes. El 20% del importe decaído va al UBI.',
  'ubi-src-cap':'Desbordamiento del Límite','ubi-src-cap-d':'Wallets que superan el límite de riqueza (max(5,min(N,25))× promedio) son confiscadas al instante. El 20% fluye al UBI.',
  'pools4-header':'Los cuatro pools de redistribución',
  'ubi-see-above':'ver countdown arriba','ubi-timer-above':'⏰ countdown mostrado arriba','pool-t-timer':'Acumula — sin temporizador',
  'usp-headline':'Por primera vez en la historia — todos empiezan igual',
  'usp-sub':'Si tienes un smartphone Android, calificas. Sin banco, sin conocimientos cripto, sin inversión.',
  'usp-c1-title':'0,00 Inversión Inicial','usp-c1-desc':'El registro es completamente sin gas. Sin ETH, sin MATIC, sin tarjeta de crédito. El protocolo paga todas las comisiones.',
  'usp-c2-title':'1.000 AEQ para cada humano','usp-c2-desc':'Millonario o agricultor — todos reciben exactamente 1.000 AEQ. Inicio igual, garantizado matemáticamente.',
  'usp-c3-title':'Accesible para todos','usp-c3-desc':'Sin cuenta bancaria, tarjeta de crédito ni documento de identidad. El registro usa un kit biométrico asequible (escáner de huella + sensor de pulso, ~15 €) — diseñado para acceso global.',
  'usp-c4-title':'UBI diario para siempre','usp-c4-desc':'Tras registrarte recibes automáticamente una parte diaria de los pagos UBI — cada día, sin ninguna acción requerida.',
  'v7-intro-title':'¿Qué es AequitasV7?',
  'v7-intro-text':'AequitasV7 es el contrato inteligente central del protocolo Aequitas. "V7" es la 7ª versión mayor del contrato de equidad. Es inmutable en Aequitas Chain (ID 1926) y gestiona todo: registro humano, verificación ZK, gestión de saldos, límite de riqueza, distribución UBI, comisiones de swap. Ningún administrador puede actualizarlo. Los seis mecanismos forman un sistema autorreforzante: el demurrage alimenta el UBI, el desbordamiento del límite suma al UBI, las comisiones se distribuyen entre los cuatro pools simultáneamente.',
  'demurrage-title':'Demurrage — Incentivo para Circular',
  'demurrage-desc':'Aequitas implementa un mecanismo de demurrage inspirado en monedas complementarias históricas. Los saldos AEQ inactivos pierden valor lentamente para desalentar el acaparamiento.',
  'dem-rate-k':'Tasa de Decaimiento','dem-rate-v':'0.5% por mes (continuo, no escalonado)',
  'dem-grace-k':'Período de Gracia','dem-grace-v':'3 meses de inactividad antes de que comience el decaimiento',
  'dem-reset-k':'Reinicio del Reloj','dem-reset-v':'Cualquier transferencia, swap o acción de liquidez reinicia el temporizador',
  'dem-dest-k':'AEQ decaído va a','dem-dest-v':'Pools de redistribución (división 40/30/20/10)',
  'dem-warn-k':'Sistema de Advertencia','dem-warn-v':'Aviso de 14 días (una vez) + recordatorio de 7 días repetido en cada inicio',
  'story-title':'La Historia de Aequitas','story-text':'<p>El año es 2009. Satoshi Nakamoto lanza Bitcoin. Por primera vez el valor puede transferirse sin bancos. Una revolución genuina. Pero casi de inmediato algo sale mal.</p><p>Los primeros mineros acumulan millones de monedas a costo casi cero. Para 2021, el 1% superior controla más del 90% de todo el Bitcoin. El coeficiente Gini estimado de Bitcoin supera 0.85 — más alto que cualquier país en la Tierra.</p><p><span style="color:var(--gold)">Aequitas</span> fue creado para responder: <em style="color:var(--gold)">"¿Cómo sería una criptomoneda diseñada para ser justa con todo ser humano?"</em></p><p>La respuesta: <strong style="color:var(--text)">El dinero existe porque las personas existen. Por lo tanto, cada persona debería tener una parte igual del dinero por el simple hecho de ser humana.</strong></p><p><em style="color:var(--gold)">"El dinero existe porque las personas existen. Nada más, nada menos."</em></p>',
  'nodes-title':'Nodos Activos — Topología Actual de la Red',
  'nodes-desc':'La red Aequitas opera actualmente en dos nodos distribuidos geográficamente. Ambos participan en la producción de bloques, sincronización de estado y servicio de API. Se comunican peer-to-peer via libp2p y sincronizan el estado de bloques via HTTP. La red está diseñada para soportar nodos adicionales.',
  'node1':'Nodo 1 — Railway (Primario)','node1-desc':'API primario · Productor de bloques · Distribución UBI · Bootstrap P2P · PostgreSQL · RPC para MetaMask',
  'node2':'Nodo 2 — Render (Secundario)','node2-desc':'API secundario · Productor de bloques · Par P2P · Sincronización HTTP · Estado PostgreSQL compartido',
  'run-node-title':'Ejecuta Tu Propio Nodo — Ayuda a Asegurar la Red',
  'run-node-desc':'Cualquiera puede ejecutar un nodo de Aequitas — sin permiso, sin stake, sin solicitud requerida. Los nodos participan en la producción de bloques y validan el registro humano. Los operadores de nodos ganan una parte de las comisiones del protocolo via el Pool de Validadores (40% de todas las comisiones de swap, distribuidas diariamente).',
  'bootstrap-title':'Conectar un Nuevo Nodo','bootstrap-desc':'Para ejecutar tu propio nodo, establece PRIMARY_NODE_URL=https://aequitas.digital en tu entorno. Tu nodo sincronizará automáticamente el estado completo de la cadena.',
  'tech-title':'Especificaciones Técnicas','mm-config':'Configuración MetaMask',
  'k-lang':'Idioma','k-src':'Código Fuente','evm-yes':'Sí — JSON-RPC /rpc · Compatible con MetaMask',
  'proto-label':'Protocolo Aequitas V7 — Documentación Técnica',
  'ca-title':'Contratos y Direcciones de Red','ca-text':'Cadena: Aequitas Chain (ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier (verificador Groth16 on-chain): 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (contrato principal): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7 es la única fuente de verdad para toda la economía Aequitas. Cada saldo AEQ, cada registro humano, cada pago UBI y cada aplicación del límite de riqueza está gobernado por este único contrato inmutable — desplegado en Aequitas Chain, una blockchain personalizada compatible con EVM que ejecuta un motor de consenso BlockDAG. No hay clave de administrador, no hay proxy de actualización, no hay votación de gobernanza que pueda cambiar una sola línea de su lógica. El código que funciona hoy es el código que funcionará en diez años.<br><br>El contrato BioVerifier recibe pruebas de conocimiento cero Groth16 generadas completamente en el dispositivo Android del usuario. Verifica matemáticamente on-chain en ~10 ms que un nuevo registrante es un ser humano único y vivo — sin conocer jamás su nombre, identidad o datos biométricos. Esto es lo que hace posible el registro sin gas y sin inversión: la prueba es lo único que sale del dispositivo.<br><br>Juntos, estos dos contratos hacen posible algo que nunca ha existido en ningún sistema monetario de la historia: una oferta monetaria cuyas reglas — quién la recibe, cuánto existe, cómo se redistribuye — no puede ser alterada por ninguna persona, empresa o gobierno. Jamás.',
  'ib-poh':'Prueba de Humanidad','ib-poh-t':'Cada titular de AEQ debe probar criptográficamente que es un ser humano único y vivo. Sin bots, sin corporaciones, sin IA, sin duplicados. Los datos biométricos nunca salen de tu dispositivo — solo se transmite una prueba matemática de unicidad. AEQ es la primera moneda que es exclusivamente humana.',
  'ib-fair':'Distribución Radicalmente Justa','ib-fair-t':'Cada humano verificado recibe exactamente 1.000 AEQ al registrarse. Sin pre-minado, sin asignación a fundadores, sin rondas de inversores. El suministro total es siempre y exactamente igual al número de humanos verificados multiplicado por 1.000. Esto se aplica matemáticamente, no por política.',
  'ib-dag':'Arquitectura BlockDAG','ib-dag-t':'A diferencia de las blockchains tradicionales donde solo puede existir un bloque por altura, Aequitas usa una estructura DAG. Múltiples bloques pueden producirse simultáneamente por diferentes nodos y luego fusionarse en el DAG. Esto permite mayor rendimiento, menor latencia y elimina cuellos de botella. Los eventos de fusión se marcan con una insignia especial en el explorador.',
  'ib-gas':'Verdaderamente Sin Gas','ib-gas-t':'Todos los registros y transferencias de AEQ no cuestan absolutamente nada. No se necesita ETH, BNB ni MATIC. Sin tarjeta de crédito, sin cuenta bancaria, sin criptomoneda previa. El relayer cubre todos los costos de transacción. Si eres humano con un smartphone, puedes participar independientemente de tu situación económica.',
  'h-what':'¿Qué es un Humano Verificado?','h-what-t':'Un Humano Verificado es una dirección wallet demostrada criptográficamente que pertenece a un ser humano único y vivo. La verificación usa un sistema de hardware de 3 factores: R503 escanea las 10 huellas; MAX30102 PPG confirma pulso vivo; Fase 2: IR de venas de mano (1 en 10⁷); Fase 3: iris (1 en 10⁷⁸). Solo se transmite una prueba ZK Groth16. Ningún dato biométrico abandona el dispositivo.',
  'h-zkp':'Sistema de Prueba ZK','h-zkp-t':'Aequitas usa Groth16 en BN128 — misma curva que Ethereum y Zcash. ~200 bytes, ~10ms. commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier vinculado al cuerpo: perder tu teléfono no crea una segunda identidad. No se almacena ningún dato biométrico.',
  'h-sybil':'Prevención de Ataques Sybil','h-sybil-t':'Fase 1: las 10 huellas + vivacidad MAX30102 (pulso PPG, rechaza moldes/repeticiones). Fase 2: IR de venas de mano — característica corporal interna, imposible de copiar, 1 en 10⁷, diferente en gemelos idénticos. Fase 3: iris — 1 en 10⁷⁸, estándar de oro global. Nullifier = keccak256(iris‖vein‖domain). Un humano, una wallet, para siempre.',
  'h-global':'Inclusión Financiera Global','h-global-t':'1.400 millones de adultos en todo el mundo no tienen cuenta bancaria. Aequitas solo requiere un smartphone Android con sensor biométrico — un dispositivo que más de 3.000 millones de personas ya poseen. Sin cuenta bancaria, sin tarjeta de crédito, sin criptomoneda previa, sin documento de identidad. Simplemente ser humano es suficiente.',
  'h-bio-hw':'Hoja de Ruta del Hardware Biométrico','h-bio-hw-t':'Fase 1 (activa): escáner de huellas óptico R503 — hash combinado de los 10 dedos. Vivacidad MAX30102 PPG. Fase 2 (planificada): ESP32-CAM + LED IR de 850 nm — imagen de venas de mano, 1 en 10⁷ de unicidad. Fase 3 (planificada): módulo de iris IR — 240+ grados de libertad, 1 en 10⁷⁸, totalmente independiente del dispositivo, gemelos idénticos difieren.',
  'poa-title':'1. PRUEBA DE VIDA — Recuperación de Saldos Inactivos','poa-text':'<p>¿Qué pasa con AEQ cuando las personas mueren o quedan permanentemente incapacitadas? En Bitcoin, las wallets perdidas significan suministro perdido permanentemente. Aequitas soluciona esto mediante un sistema de recuperación por inactividad de múltiples etapas: si una wallet no muestra actividad durante un período prolongado, su saldo se devuelve gradualmente a la comunidad a través del pool UBI.</p>',
  'poa-box':'Año 0–2: Uso normal — sin restricciones<br>Año 2: Aviso 1 — el Guardian puede responder en nombre<br>Año 2+60d: Aviso 2 — urgencia creciente<br>Año 2+120d: Aviso 3 — aviso final<br>Año 2+180d: AEQ movido a CUSTODIA personal (aún recuperable)<br>Año 4: Si aún inactivo — CUSTODIA liberada al Pool UBI',
  'guard-title':'2. SISTEMA GUARDIAN — Salvaguarda Humana','guard-text':'<p>¿Y si alguien está hospitalizado, encarcelado o de algún modo incapaz de acceder a su dispositivo por meses? El sistema Guardian permite a una persona de confianza — otro humano verificado — confirmar que el propietario de la wallet sigue vivo. El Guardian tiene estrictamente cero acceso financiero: solo puede llamar una función que reinicia el temporizador de inactividad.</p>',
  'guard-box':'1 Guardian por humano · debe ser un humano verificado en Aequitas<br>Guardian SOLO puede llamar confirmAlive() — cero derechos de transacción<br>Guardian NO PUEDE mover fondos, transferir AEQ ni acceder a la wallet<br>Máximo 3 tutelados por Guardian (evita centralización de confianza)<br>Bloqueo de 7 días en asignación de Guardian (evita asignación forzada)<br>No se permiten relaciones Guardian circulares',
  'dem-title':'3. DEMURRAGE — Mecanismo Anti-Acaparamiento',
  'dem-box':'Tasa: 0,5% por mes después de 3 meses de inactividad (continuo, no escalonado)<br>El reloj se reinicia automáticamente con cualquier transferencia, swap o acción de liquidez<br>AEQ decaído redistribuido a los cuatro pools — nunca destruido<br>Aviso de 14 días mostrado una vez · aviso de 7 días repetido en cada sesión activa',
  'dem-text':'<p>El demurrage es un costo de tenencia sobre el dinero — una tasa de interés negativa que hace costoso acumular y atractivo circular. El experimento de Wörgl (Austria, 1932) usó una moneda con demurrage y redujo el desempleo local un 25% en un año. El Banco Central de Austria lo cerró precisamente porque funcionó demasiado bien. El Chiemgauer (Alemania, 2003) opera según el mismo principio con éxito desde hace más de 20 años.</p>',
  'cap-title':'4. LÍMITE DE RIQUEZA — Aplicación de Justicia Matemática','cap-box':'Límite bootstrap: max(5,min(N,25))× saldo promedio actual<br>1–4 humanos: 5× · +1× por humano · 25+: 25× permanente<br>Se aplica a TODAS las direcciones excepto las 4 pools del protocolo<br>Exceso AEQ redistribuido instantáneamente · Sin intervención manual',
  'ubi-title':'5. RENTA BÁSICA UNIVERSAL — Redistribución Diaria','ubi-box':'Fuentes de ingresos del Pool UBI:<br>· 20% de todas las comisiones de swap del pool AMM AEQ↔tUSD<br>· Desbordamiento de la aplicación del límite de riqueza<br>· Cargos de demurrage de cuentas inactivas<br>· Custodia inactiva liberada después de 4 años<br><br>Distribución: Cada 24 horas, todo el saldo del pool UBI se divide igualmente entre todos los humanos verificados registrados. El pool se reinicia a cero y comienza a llenarse inmediatamente de la actividad continua del protocolo.',
  'inf-title':'6. SIN INFLACIÓN ALGORÍTMICA — Fórmula de Suministro Fijo','inf-box':'El ÚNICO evento que crea nuevo AEQ: un nuevo humano verificado se registra.<br><br>Suministro Total = Humanos Verificados × 1.000 AEQ<br><br>Esto no es una política — es aplicado por el protocolo. Ningún administrador puede acuñar AEQ adicional, ningún voto de gobernanza puede cambiar la emisión. AEQ es la única criptomoneda donde el suministro total está determinado únicamente por el número de humanos vivos verificados.',
  'explore-title':'Explorar Aequitas',
  'expl-score':'Puntuación de Igualdad','expl-score-d':'Coeficiente Gini en vivo · Índice Aequitas · distribución de riqueza en tiempo real',
  'expl-economy':'UBI y Pools de Redistribución','expl-economy-d':'Cuenta regresiva UBI diaria · 4 pools on-chain · demurrage · Fases del Protocolo',
  'expl-charts':'Gráficos e Historial','expl-charts-d':'Historial Gini · curva de Lorenz · slider bootstrap del límite de riqueza · La historia de Aequitas',
  'expl-v7':'Documentación Protocolo V7','expl-v7-d':'Contrato AequitasV7 · 6 mecanismos · prueba ZK · límite de riqueza · demurrage · código inmutable',
  'expl-explorer':'Explorador de Bloques','expl-explorer-d':'BlockDAG en vivo · haz clic en cualquier bloque para ver validador, hash, transacciones, hashes padres',
  'swap-sell-label':'Vender','swap-receive-label':'Recibir',
  'expl-network':'Red y Nodos','expl-network-d':'Topología de nodos · ejecutar tu propio nodo · especificaciones técnicas · Chain ID 1926'
},
ru:{
  'logo-sub':'ДОКАЗАТЕЛЬСТВО ЧЕЛОВЕЧНОСТИ','live':'ОНЛАЙН',
  'tab-register':'🔐 Регистрация','tab-explorer':'🔍 Проводник','tab-humans':'👥 Люди','tab-index':'📊 Индекс','tab-network':'🌐 Сеть','tab-protocol':'📜 Протокол V7','tab-swap':'🔄 Обмен',
  'reg-title':'🔐 Зарегистрируйтесь как Верифицированный Человек',
  'reg-sub':'Присоединитесь к сети Aequitas и получите 1 000 AEQ в качестве Универсального Базового Дохода. Однократно, постоянно и полностью бесплатно. Никакие личные данные никогда не сохраняются.',
  'app-title':'РЕГИСТРАЦИЯ ТОЛЬКО ЧЕРЕЗ ANDROID-ПРИЛОЖЕНИЕ',
  'app-text':'Доказательство Человечности использует физическую биометрическую систему с 3 факторами. Фаза 1: оптический сенсор R503 сканирует все 10 отпечатков + MAX30102 PPG подтверждает живой пульс. Фаза 2: ИК-камера вен руки (1 из 10⁷ уникальности). Фаза 3: сканирование радужной оболочки — золотой стандарт, 1 из 10⁷⁸, полностью независимо от устройства. Доказательство Groth16 ZK фиксирует все факторы без раскрытия биометрических данных. 1 000 AEQ зачисляются автоматически после верификации.',
  's1t':'Биометрическое Сканирование','s1d':'AequitasBio сканирует все 10 отпечатков (оптический сенсор R503) + пульс MAX30102 PPG подтверждает живость. Фаза 2: ИК вен руки. Фаза 3: радужная оболочка. Необработанные данные никогда не покидают устройство.',
  's2t':'Создание ZK-Доказательства','s2d':'Доказательство Groth16 ZK фиксирует все биометрические факторы: commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier привязан к телу, не к телефону — потеря устройства не создаёт вторую идентичность.',
  's3t':'Подключение Кошелька','s3d':'Приложение открывает MetaMask на этой странице · подключите кошелёк Ethereum · доказательство криптографически привязано к вашему адресу',
  's4t':'1 000 AEQ Зачислены','s4d':'Регистрация подтверждена на BlockDAG Aequitas за 6 секунд · 1 000 AEQ зачислены мгновенно · личность навсегда записана как верифицированный человек',
  'priv-bar':'🔒 R503 Все 10 Отпечатков · MAX30102 Живость · Фаза 2: ИК Вен Руки · Фаза 3: Радужная Оболочка (10⁷⁸) · Groth16 ZKP · Данные никогда не покидают устройство · Один человек · Навсегда',
  'conn-wallet':'ПОДКЛЮЧЁННЫЙ КОШЕЛЁК','proof-recv':'⚡ ZK-ДОКАЗАТЕЛЬСТВО ПОЛУЧЕНО','proof-hint':'Подключите кошелёк для регистрации',
  'btn-conn':'🦊 ПОДКЛЮЧИТЬ METAMASK','btn-reg':'🔐 ЗАРЕГИСТРИРОВАТЬ ОН-ЧЕЙН',
  'btn-web-reg':'🌐 РЕГИСТРАЦИЯ ЧЕРЕЗ БРАУЗЕР (WebAuthn)',
  'web-reg-warn':'⚠ Привязано к устройству: Эта личность привязана к данному устройству и браузеру. Перенести её на другое устройство невозможно. Для постоянной кроссплатформенной личности используйте Android-приложение Aequitas.',
  'reg-log-hint':'// Откройте Android-приложение Aequitas для создания доказательства, затем вернитесь сюда...',
  'reg-details':'Детали Регистрации','k-network':'Сеть','k-chainid':'ID Цепи','k-grant':'Субсидия UBI',
  'k-fee':'Комиссия Gas','free':'БЕСПЛАТНО — полностью без комиссий','k-limit':'Регистрации','k-limit-v':'Один раз · постоянно · неизменно',
  'k-bio':'Биометрические Данные','never-stored':'Никогда не сохраняются — остаются на устройстве',
  'k-proof':'Система Доказательств','k-conf':'Подтверждение','k-conf-v':'В течение 6 секунд (1 блок)',
  'k-sybil':'Защита от Сибилл','k-sybil-v':'Одна идентичность на биометрию · постоянная блокировка',
  'live-stats':'Статистика Цепи в Реальном Времени',
  's-height':'Высота Блока','s-height-sub':'Новый блок каждые ~6с · BlockDAG · Параллельное производство',
  's-humans':'Верифицированные Люди','s-humans-sub':'Биометрический ZKP · Один человек, один кошелёк, навсегда',
  's-supply':'Общий Объём','s-supply-sub':'Всегда = Люди × 1 000 AEQ',
  's-index':'Индекс Aequitas','s-index-sub':'0 = идеальное равенство · 100 = максимальное неравенство',
  's-uptime':'Время Работы','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Доказательство Человечности','ib-poh-t':'Каждый владелец AEQ должен криптографически доказать что является уникальным живым человеком. Никаких ботов, корпораций, ИИ. Биометрические данные никогда не покидают устройство.',
  'ib-fair':'Радикально Справедливое Распределение','ib-fair-t':'Каждый верифицированный человек получает ровно 1 000 AEQ при регистрации. Никакого предварительного майнинга, никаких аллокаций основателям. Общий объём всегда равен верифицированные люди × 1 000.',
  'ib-dag':'Архитектура BlockDAG','ib-dag-t':'Несколько блоков могут производиться одновременно и объединяться. Более высокая пропускная способность, меньшая задержка.',
  'ib-gas':'Действительно Без Комиссий','ib-gas-t':'Регистрация и переводы AEQ не стоят ничего. ETH, BNB или MATIC не требуются. Банковский счёт и кредитная карта не нужны.',
  'recent-blocks':'Последние Блоки','blocks-desc':'MERGE = объединение нескольких родителей (BlockDAG). TX = транзакция регистрации. Время блока: ~6 секунд.',
  'loading':'Загрузка блоков...','net-info':'Информация о Сети','k-chain':'Имя Цепи','k-symbol':'Символ','k-btime':'Время Блока',
  'k-cons':'Консенсус','k-nodes':'Активные Ноды','k-storage':'Хранилище','add-mm':'🦊 ДОБАВИТЬ В METAMASK','k-dec':'Десятичные',
  'btn-add-mm':'+ ДОБАВИТЬ СЕТЬ AEQUITAS',
  'phil':'"Деньги существуют потому что существуют люди.<br>Ничего более, ничего менее."','phil-sub':'— ПРИНЦИП AEQUITAS —',
  'humans-title':'Верифицированные Люди в Aequitas Chain',
  'h-what':'Что такое Верифицированный Человек?','h-what-t':'Верифицированный Человек — адрес кошелька, криптографически доказанный принадлежащим уникальному живому человеку. Верификация использует аппаратную систему из 3 факторов: R503 сканирует все 10 отпечатков; MAX30102 PPG подтверждает живой пульс; Фаза 2: ИК вен руки (1 из 10⁷); Фаза 3: радужная оболочка (1 из 10⁷⁸). Передаётся только доказательство Groth16 ZK. Биометрические данные никогда не покидают устройство.',
  'h-zkp':'Система ZK-Доказательств','h-zkp-t':'Aequitas использует Groth16 на BN128 — та же кривая, что Ethereum и Zcash. ~200 байт, ~10мс. commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier привязан к телу: потеря телефона не создаёт вторую идентичность. Биометрические данные никогда не хранятся.',
  'h-sybil':'Защита от Атак Сибиллы','h-sybil-t':'Фаза 1: все 10 отпечатков + живость MAX30102 (пульс PPG, отклоняет слепки/повторы). Фаза 2: ИК вен руки — внутренняя характеристика тела, невозможно скопировать, 1 из 10⁷, различается у однояйцевых близнецов. Фаза 3: радужная оболочка — 1 из 10⁷⁸, мировой золотой стандарт. Nullifier = keccak256(iris‖vein‖domain). Один человек, один кошелёк, навсегда.',
  'h-global':'Глобальная Финансовая Инклюзия','h-global-t':'Банковский счёт, кредитная карта или криптовалюта не требуются. Только Android-смартфон с биометрическим датчиком.',
  'h-bio-hw':'Дорожная Карта Биометрического Оборудования','h-bio-hw-t':'Фаза 1 (активна): оптический сканер отпечатков R503 — комбинированный хэш всех 10 пальцев. Живость MAX30102 PPG. Фаза 2 (планируется): ESP32-CAM + ИК LED 850 нм — визуализация вен руки, 1 из 10⁷ уникальности. Фаза 3 (планируется): ИК модуль радужной оболочки — 240+ степеней свободы, 1 из 10⁷⁸, полностью независимо от устройства, однояйцевые близнецы различаются.',
  'reg-humans':'Зарегистрированные Люди','h-desc':'Каждый адрес верифицирован как уникальный человек через биометрический ZKP. Каждый получил ровно 1 000 AEQ. Постоянно, неизменно, он-чейн.',
  'no-humans':'Люди ещё не зарегистрированы.\n\nСкачайте Android-приложение Aequitas и будьте первым!',
  'reg-stats':'Статистика Реестра','total-humans':'Всего Людей',
  'idx-title':'Индекс Aequitas — Оценка Экономического Равенства в Реальном Времени',
  'idx-desc':'Индекс Aequitas измеряет экономическое неравенство всех верифицированных людей в реальном времени. Рассчитывается из коэффициента Джини распределения балансов он-чейн. 0 = идеальное равенство. 100 = максимальное неравенство.',
  'curr-idx':'Текущий Индекс','bar-0':'0 — Идеальное Равенство','bar-100':'100 — Макс. Неравенство',
  'gini':'Коэффициент Джини','gini-desc':'0 = равно · 1 = неравно',
  'supply-desc':'Всегда = Люди × 1 000 AEQ',
  'phase':'Фаза Протокола','phase-desc':'Автоматически по количеству людей',
  'humans-desc':'Биометрически верифицированные уникальные люди',
  'pools-title':'Пулы Перераспределения',
  'pools-desc':'Каждая комиссия свопа, плата за демередж и превышение лимита богатства автоматически делится между четырьмя пулами. Все пулы выплачивают ежедневно.',
  'vel-pool':'Пул Валидаторов','vel-pool-desc':'40% всех комиссий → операторы нод, обеспечивающие сеть',
  'liq-pool':'Пул Ликвидности','liq-pool-desc':'30% всех комиссий → поставщики ликвидности, пропорционально LP-долям',
  'ubi-pool':'Пул UBI','ubi-pool-desc':'20% всех комиссий → все верифицированные люди поровну, каждые 24 часа',
  'treasury':'Казначейство','treasury-desc':'10% всех комиссий → разработка и обслуживание протокола',
  'phases-title':'Фазы Протокола',
  'phases-desc':'В Фазе 0 (Bootstrap) применяется скользящий множитель: max(5, min(N, 25))× средний баланс. При 1–4 людях: 5× средний. Каждый новый человек прибавляет 1×. При 25+ людях: фиксируется навсегда на 25×. Фаза 1+ сохраняет 25× фиксированным. Переходы автоматические — без голосования, без административных ключей.',
  'p0':'Bootstrap · &lt;100 людей · Лимит богатства: max(5,min(N,25))× средний · Скользит 5×→25× до 25-го человека · Сейчас активен',
  'p1':'Рост · 100–10 000 людей · Лимит богатства: 25× средний баланс',
  'p2':'Стабильность · 10 000–1М людей · Лимит богатства: 25× средний баланс',
  'p3':'Зрелость · 1М+ людей · Лимит богатства: 25× средний баланс',
  'wealth-cap-explain':'В Фазе 0 (Bootstrap) Лимит Богатства = max(5, min(N, 25))× средний баланс AEQ, где N = количество зарегистрированных людей. 1–4 человека: 5× средний. Каждый новый человек прибавляет 1×. 25+ людей: фиксируется навсегда на 25×. Лимит всегда привязан к актуальному среднему балансу.',
  'demurrage-title':'Демередж — Стимул к Обращению',
  'demurrage-desc':'Aequitas реализует механизм демереджа, вдохновлённый историческими дополнительными валютами. Бездействующие балансы AEQ постепенно теряют стоимость для предотвращения накопления.',
  'dem-rate-k':'Скорость Распада','dem-rate-v':'0,5% в месяц (непрерывно)',
  'dem-grace-k':'Льготный Период','dem-grace-v':'3 месяца бездействия до начала распада',
  'dem-reset-k':'Сброс Таймера','dem-reset-v':'Любой перевод, своп или операция с ликвидностью сбрасывает таймер',
  'dem-dest-k':'Распавшийся AEQ идёт в','dem-dest-v':'Пулы перераспределения (40/30/20/10)',
  'dem-warn-k':'Система Предупреждений','dem-warn-v':'14-дневное уведомление (один раз) + 7-дневное повторение при каждом входе',
  'story-title':'История Aequitas — Почему это существует',
  'story-text':'<p>Год 2009. Сатоши Накамото выпускает Bitcoin. Впервые ценность может передаваться между людьми без банка. Настоящая революция. Но почти сразу что-то идёт не так.</p><p>Ранние майнеры накапливают миллионы монет почти бесплатно. К 2021 году топ 1% адресов Bitcoin контролирует более 90% всех Bitcoin. Коэффициент Джини Bitcoin превышает 0,85 — выше чем в любой стране мира.</p><p><span style="color:var(--gold)">Aequitas</span> был создан чтобы ответить на один вопрос: <em style="color:var(--gold)">"Как выглядела бы криптовалюта, спроектированная с нуля чтобы быть справедливой для каждого человека?"</em></p><p>Ответ прост: <strong style="color:var(--text)">Деньги существуют потому что существуют люди. Поэтому каждый человек должен иметь равную долю денег просто потому что он человек.</strong></p><p><em style="color:var(--gold)">"Деньги существуют потому что существуют люди. Ничего более, ничего менее."</em></p>',
  'nodes-title':'Активные Ноды — Текущая Топология Сети','nodes-desc':'Сеть Aequitas работает на двух географически распределённых нодах. Обе участвуют в производстве блоков и синхронизации. Сеть рассчитана на дополнительные ноды.',
  'node1':'Нода 1 — Railway (Основная)','node1-desc':'Основной API · Производитель блоков · Распределение UBI · P2P Bootstrap · PostgreSQL · RPC для MetaMask',
  'node2':'Нода 2 — Render (Вторичная)','node2-desc':'Вторичный API · Производитель блоков · P2P-пир · HTTP-синхронизация · Общее состояние PostgreSQL',
  'run-node-title':'Запустите Свою Ноду — Помогите Защитить Сеть',
  'run-node-desc':'Любой может запустить ноду без разрешения. Операторы нод получают 40% всех комиссий свопа ежедневно через Пул Валидаторов.',
  'bootstrap-title':'Подключить Новую Ноду','bootstrap-desc':'Установите PRIMARY_NODE_URL=https://aequitas.digital в вашей среде. Нода автоматически синхронизируется и начнёт производство блоков.',
  'tech-title':'Технические Характеристики','mm-config':'Конфигурация MetaMask',
  'k-lang':'Язык','k-src':'Исходный Код','evm-yes':'Да — JSON-RPC /rpc · Совместимо с MetaMask',
  'proto-label':'Протокол Aequitas V7 — Техническая Документация',
  'ca-title':'Адреса Контрактов','ca-text':'Цепь: Aequitas Chain (ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7: 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7 является единственным источником истины для всей экономики Aequitas. Каждый баланс AEQ, каждая регистрация человека, каждая выплата UBI и каждое применение ограничения богатства управляется этим одним неизменяемым контрактом — развёрнутым на Aequitas Chain, специализированном блокчейне совместимом с EVM работающем на механизме консенсуса BlockDAG. Нет ключа администратора, нет прокси обновления, нет голосования по управлению которое могло бы изменить хотя бы одну строку его логики. Код работающий сегодня — это код который будет работать через десять лет.<br><br>Контракт BioVerifier получает доказательства с нулевым разглашением Groth16 сгенерированные полностью на Android-устройстве пользователя. Он математически проверяет on-chain примерно за 10 мс что новый регистрант является уникальным живым человеком — не узнавая никогда его имени, личности или биометрических данных. Именно это делает возможной безгазовую регистрацию без инвестиций: доказательство — единственное что когда-либо покидает устройство.<br><br>Вместе эти два контракта делают возможным то чего никогда не существовало ни в одной денежной системе в истории: денежное предложение правила которого — кто его получает, сколько существует, как оно перераспределяется — не могут быть изменены ни одним человеком, компанией или правительством. Никогда.',
  'poa-title':'1. ДОКАЗАТЕЛЬСТВО ЖИЗНИ — Восстановление Неактивных Балансов','poa-text':'<p>Что происходит с AEQ когда люди умирают или становятся недееспособными? В Bitcoin потерянные кошельки означают навсегда потерянный объём. Aequitas решает это через многоуровневую систему: если кошелёк не проявляет активности в течение длительного периода, его баланс постепенно возвращается сообществу через пул UBI.</p>',
  'poa-box':'Год 0–2: Обычное использование — без ограничений<br>Год 2: Предупреждение 1 — Guardian может ответить от имени<br>Год 2+60д: Предупреждение 2 — нарастающая срочность<br>Год 2+120д: Предупреждение 3 — последнее уведомление<br>Год 2+180д: AEQ перемещён в личный ЭСКРОУ (ещё восстановимо)<br>Год 4: При сохранении бездействия — ЭСКРОУ в Пул UBI',
  'guard-title':'2. СИСТЕМА GUARDIAN — Человеческая Защита','guard-text':'<p>Что если кто-то госпитализирован или иначе не может получить доступ к устройству месяцами? Система Guardian позволяет доверенному лицу — другому верифицированному человеку — подтвердить что владелец кошелька жив. Guardian имеет строго нулевой финансовый доступ: он может только сбросить таймер бездействия.</p>',
  'guard-box':'1 Guardian на человека · должен быть верифицированным человеком в Aequitas<br>Guardian может ТОЛЬКО вызывать confirmAlive() — ноль прав транзакций<br>Guardian НЕ МОЖЕТ перемещать средства, переводить AEQ или получать доступ к кошельку<br>Максимум 3 подопечных · Блокировка 7 дней при назначении · Без круговых отношений',
  'dem-title':'3. ДЕМЕРЕДЖ — Механизм Против Накопления',
  'dem-box':'Ставка: 0,5%/месяц после 3 месяцев бездействия (непрерывно, не ступенчато)<br>Таймер сбрасывается при любом переводе, свопе или операции с ликвидностью<br>Decayed AEQ перераспределяется в пулы — никогда не сжигается',
  'dem-text':'<p>Демередж — стоимость хранения денег. Эксперимент Вёрглена (Австрия, 1932) сократил местную безработицу на 25% за год. Chiemgauer (Германия, 2003) работает по тому же принципу уже более 20 лет.</p>',
  'cap-title':'4. ЛИМИТ БОГАТСТВА — Математическое Обеспечение Справедливости','cap-box':'Bootstrap-лимит: max(5,min(N,25))× текущий средний баланс<br>1–4 людей: 5× · +1× за человека · 25+: 25× навсегда<br>Применяется ко всем адресам кроме 4 протокольных пулов<br>Избыток AEQ мгновенно перераспределяется · Без ручного вмешательства',
  'ubi-title':'5. УНИВЕРСАЛЬНЫЙ БАЗОВЫЙ ДОХОД — Ежедневное Перераспределение','ubi-box':'Источники: Комиссии свопов (20%) · Превышение лимита богатства · Демередж · Эскроу после 4 лет<br><br>Ежедневно: весь пул UBI делится поровну между всеми зарегистрированными людьми. Пул сбрасывается и сразу наполняется снова.',
  'inf-title':'6. НИКАКОЙ АЛГОРИТМИЧЕСКОЙ ИНФЛЯЦИИ — Фиксированная Формула','inf-box':'ЕДИНСТВЕННОЕ событие создающее новый AEQ: регистрируется новый верифицированный человек.<br><br>Общий Объём = Верифицированные Люди × 1 000 AEQ<br><br>Это не политика — обеспечивается протоколом. AEQ — единственная криптовалюта где объём определяется исключительно числом верифицированных живых людей.',
  'phases-desc':'В Фазе 0 лимит богатства использует скользящий Bootstrap-множитель: max(5, min(N, 25))× средний баланс. При 1–4 людях: 5× средний. Каждый новый человек прибавляет 1×. При 25+ людях: фиксируется навсегда на 25×. Фаза 1+ сохраняет 25× фиксированным. Переходы автоматические — без голосования, без административных ключей.',
  'p0':'Bootstrap · &lt;100 людей · Лимит богатства: max(5,min(N,25))× средний · Скользит 5×→25× до 25-го человека · Сейчас активен',
  'p1':'Рост · 100–10 000 людей · Лимит богатства: 25× средний баланс',
  'p2':'Стабильность · 10 000–1M людей · Лимит богатства: 25× (планируемое снижение: 10×)',
  'p3':'Зрелость · 1M+ людей · Лимит богатства: 25× (планируемое снижение: 5×)',
  'wealth-cap-explain':'Лимит богатства в настоящее время установлен на 25× среднего баланса AEQ всех верифицированных людей. Это фиксированная константа в живом коде Go. Поскольку значение всегда относительно текущего среднего, лимит автоматически масштабируется по мере роста сети.',
  'btn-download-app':'СКАЧАТЬ ПРИЛОЖЕНИЕ AEQUITASBIO',
  'swap-title':'🔄 Обмен AEQ ↔ tUSD','swap-sub':'Обменивайте AEQ на tUSD (симулированный тестовый доллар) через нативный пул ликвидности. Комиссия 0,1% применяется только к свопам — обычные переводы AEQ между людьми остаются полностью бесплатными.',
  'swap-priv-bar':'🔒 Только 0,1% комиссия свопа · Переводы AEQ-AEQ бесплатны · tUSD — тестовая валюта без реальной стоимости',
  'swap-your-aeq':'Ваш AEQ','swap-your-tusd':'Ваш tUSD','swap-aeq-to-tusd':'AEQ → tUSD','swap-tusd-to-aeq':'tUSD → AEQ',
  'swap-fee-est':'Комиссия протокола (0,1%)','swap-details-hdr':'Детали Свопа',
  'swap-out-lbl':'Вы получите (прим.)','swap-impact-lbl':'Влияние на цену','swap-rate-lbl':'Обменный курс',
  'swap-depth-lbl':'Состав Пула','amm-title':'x × y = k — AMM с Постоянным Произведением',
  'amm-text':'Когда вы обмениваете AEQ на tUSD, резерв AEQ растёт, а резерв tUSD уменьшается — их произведение всегда равно k. Более крупные свопы вызывают большее влияние на цену. Комиссия 0,1% вычитается до применения формулы.',
  'swap-btn-conn':'🦊 ПОДКЛЮЧИТЬ METAMASK','swap-btn-go':'🔄 ОБМЕНЯТЬ',
  'swap-log-hint':'// Подключите кошелёк для обмена...',
  'swap-no-liquidity':'Нет tUSD?','swap-faucet-desc':'Зарегистрированные люди могут получить тестовый tUSD один раз','swap-btn-faucet':'💧 ПОЛУЧИТЬ ТЕСТОВЫЙ tUSD',
  'swap-addliq-title':'Предоставить Ликвидность','swap-addliq-desc':'Будьте первым кто внесёт — ваше соотношение устанавливает начальную цену.','swap-btn-addliq':'💧 ДОБАВИТЬ ЛИКВИДНОСТЬ',
  'swap-lp-title':'Ваша LP-Позиция','swap-lp-share':'Доля в Пуле','swap-lp-withdrawable':'Доступно к выводу',
  'swap-lp-pct-label':'% вашей позиции','swap-lp-youget':'Вы получите','swap-btn-removeliq':'🔥 ВЫВЕСТИ ЛИКВИДНОСТЬ',
  'swap-pool-title':'AEQ / tUSD — Статус Пула',
  'swap-pool-aeq':'Резерв AEQ','swap-pool-tusd':'Резерв tUSD','swap-pool-price':'Спотовая Цена',
  'swap-fee-bps':'Комиссия Свопа','swap-fee-split':'Распределение комиссий','swap-fee-split-v':'40% Валидаторы / 30% LP / 20% UBI / 10% Казначейство',
  'swap-pools-addr-title':'Адреса Пулов Токеномики',
  'swap-validators':'Валидаторы (40%)','swap-lps':'Провайдеры Ликвидности (30%)','swap-ubi':'Пул UBI (20%)','swap-treasury':'Казначейство (10%)',
  'ubi-hero-title':'УНИВЕРСАЛЬНЫЙ БАЗОВЫЙ ДОХОД — ПУЛ UBI',
  'ubi-hero-sub':'Накапливается — следующая выплата поровну всем верифицированным людям через:',
  'ubi-bal-lbl':'текущий баланс пула','ubi-hero-desc':'Делится поровну между всеми · выплачивается каждые 24ч · пул обнуляется после выплаты · минимальный баланс не требуется',
  'ubi-how-fills':'Как заполняется Пул UBI',
  'ubi-src-swap':'Комиссии Свопов','ubi-src-swap-d':'Каждый своп AEQ↔tUSD вносит 20% своей комиссии 0,1%. Больше торговли = быстрее заполнение.',
  'ubi-src-dem':'Демередж','ubi-src-dem-d':'Неактивный AEQ (3+ месяца) убывает со скоростью 0,5%/месяц. 20% убывшей суммы идёт в UBI.',
  'ubi-src-cap':'Превышение Лимита Богатства','ubi-src-cap-d':'Кошельки превышающие лимит (max(5,min(N,25))× средний) конфискуются мгновенно. 20% поступает в UBI немедленно.',
  'pools4-header':'Все четыре пула перераспределения',
  'ubi-see-above':'см. обратный отсчёт выше','ubi-timer-above':'⏰ обратный отсчёт показан выше','pool-t-timer':'Накапливается — без таймера',
  'usp-headline':'Впервые в истории — все начинают на равных',
  'usp-sub':'Если у вас есть Android-смартфон — вы квалифицируетесь. Без банка, без знаний крипто, без инвестиций.',
  'usp-c1-title':'0,00 стартовых инвестиций','usp-c1-desc':'Регистрация полностью без газа. Без ETH, без MATIC, без кредитной карты. Протокол оплачивает все транзакционные сборы.',
  'usp-c2-title':'1 000 AEQ для каждого человека','usp-c2-desc':'Миллиардер или фермер — все получают ровно 1 000 AEQ. Не больше, не меньше. Равный старт, гарантированный математически.',
  'usp-c3-title':'Доступно для всех','usp-c3-desc':'Без банковского счёта, кредитной карты и документов. Регистрация через доступный биометрический комплект (сканер отпечатков + датчик пульса, ~$15) — для глобального доступа.',
  'usp-c4-title':'Ежедневный UBI навсегда','usp-c4-desc':'После регистрации вы автоматически получаете ежедневную долю выплат UBI — каждый день, без каких-либо действий.',
  'v7-intro-title':'Что такое AequitasV7?',
  'v7-intro-text':'AequitasV7 — центральный смарт-контракт протокола Aequitas. "V7" — 7-я основная версия контракта справедливости. Развёрнут неизменяемым образом в Aequitas Chain (ID 1926) и управляет всем: регистрация людей, верификация ZK, управление балансами, лимит богатства, распределение UBI, комиссии свопов. Ни один администратор не может обновить его. Шесть механизмов образуют самоусиливающуюся систему.',
  'explore-title':'Исследовать Aequitas',
  'expl-score':'Индекс равенства','expl-score-d':'Коэффициент Джини · Индекс Aequitas · распределение богатства в реальном времени',
  'expl-economy':'UBI и пулы перераспределения','expl-economy-d':'Ежедневный обратный отсчёт UBI · 4 on-chain пула · демерредж · Фазы протокола',
  'expl-charts':'Графики и история','expl-charts-d':'История Джини · кривая Лоренца · ползунок начального загрузчика богатства · История Aequitas',
  'expl-v7':'Документация Протокола V7','expl-v7-d':'Контракт AequitasV7 · 6 механизмов · ZK-доказательство · лимит богатства · демерредж · неизменяемый код',
  'expl-explorer':'Обозреватель блоков','expl-explorer-d':'Живой BlockDAG · нажмите на блок чтобы увидеть валидатора, хэш, транзакции, родительские хэши',
  'swap-sell-label':'Продать','swap-receive-label':'Получить',
  'gini-calc-title':'Как рассчитывается Индекс Aequitas','gini-calc-text':'Собираются все балансы AEQ. Формула вычисляет среднее абсолютное отклонение нормализованное на n2. Результат 0-1 x 100 = Индекс.','gini-why-title':'Почему Gini','gini-why-text':'Gini учитывает полное распределение среди всех людей в одном числе.','expl-network':'Сеть и узлы','expl-network-d':'Топология узлов · запустить собственный узел · технические характеристики · Chain ID 1926'
},
zh:{
  'logo-sub':'人类证明','live':'实时',
  'tab-register':'🔐 注册','tab-explorer':'🔍 浏览器','tab-humans':'👥 人类','tab-index':'📊 指数','tab-network':'🌐 网络','tab-protocol':'📜 协议 V7','tab-swap':'🔄 兑换',
  'reg-title':'🔐 注册成为经过验证的人类',
  'reg-sub':'加入Aequitas网络并获得1,000 AEQ的普遍基本收入补贴。一次性、永久性且完全免费。永远不会存储任何个人数据。',
  'app-title':'仅通过安卓应用注册',
  'app-text':'人类证明使用物理3因素生物特征系统。第1阶段：R503光学传感器扫描全部10枚指纹 + MAX30102 PPG确认活体脉搏。第2阶段：手静脉红外摄像头（独特性1/10⁷）。第3阶段：虹膜扫描——金标准，1/10⁷⁸，完全设备无关。Groth16 ZK证明提交所有因素而不泄露任何生物特征数据。验证后自动记入1,000 AEQ。',
  's1t':'生物特征扫描','s1d':'AequitasBio扫描全部10枚指纹（R503光学传感器）+ MAX30102 PPG脉搏确认活体性。第2阶段：手静脉红外。第3阶段：虹膜。原始数据永不离开设备。',
  's2t':'ZK证明生成','s2d':'Groth16 ZK证明提交所有生物特征因素：commitment = keccak256(iris‖vein‖fingers‖wallet)。Nullifier绑定到身体而非手机——丢失设备无法创建第二身份。',
  's3t':'连接钱包','s3d':'应用在此页面打开MetaMask · 连接您的以太坊钱包 · 证明与您的地址密码绑定',
  's4t':'获得1,000 AEQ','s4d':'在6秒内在Aequitas BlockDAG上确认注册 · 立即记入1,000 AEQ · 身份永久记录为经过验证的人类',
  'priv-bar':'🔒 R503全10指纹 · MAX30102活体性 · 第2阶段：手静脉红外 · 第3阶段：虹膜（10⁷⁸）· Groth16 ZKP · 数据永不离开设备 · 一个人类 · 永远',
  'conn-wallet':'已连接钱包','proof-recv':'⚡ 已收到ZK证明','proof-hint':'连接钱包以注册',
  'btn-conn':'🦊 连接 METAMASK','btn-reg':'🔐 链上注册',
  'btn-web-reg':'🌐 通过浏览器注册 (WebAuthn)',
  'web-reg-warn':'⚠ 设备绑定：此身份绑定到当前设备和浏览器，无法转移到其他设备。如需永久性多设备身份，请使用Aequitas安卓应用。',
  'reg-log-hint':'// 打开Aequitas安卓应用生成您的证明，然后返回此处...',
  'reg-details':'注册详情','k-network':'网络','k-chainid':'链ID','k-grant':'UBI补贴',
  'k-fee':'Gas费','free':'免费——完全无Gas','k-limit':'注册','k-limit-v':'每人一次 · 永久 · 不可更改',
  'k-bio':'生物特征数据','never-stored':'从不存储——保留在您的设备上',
  'k-proof':'证明系统','k-conf':'确认','k-conf-v':'6秒内（1个区块）',
  'k-sybil':'女巫攻击防护','k-sybil-v':'每个生物特征一个身份 · 永久锁定',
  'live-stats':'实时链统计',
  's-height':'区块高度','s-height-sub':'每约6秒新区块 · BlockDAG · 并行生产',
  's-humans':'已验证人类','s-humans-sub':'生物特征ZKP · 一人一钱包，永久',
  's-supply':'总供应量','s-supply-sub':'始终 = 人类 × 1,000 AEQ',
  's-index':'Aequitas指数','s-index-sub':'0 = 完全平等 · 100 = 最大不平等',
  's-uptime':'运行时间','s-uptime-sub':'节点 v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'人类证明','ib-poh-t':'每个AEQ持有者必须密码学证明其是独特的活人。没有机器人、公司、人工智能。生物特征数据永不离开设备。',
  'ib-fair':'彻底公平的分配','ib-fair-t':'每个经过验证的人类注册时恰好获得1,000 AEQ。没有预挖矿，没有创始人分配。总供应量始终等于已验证人类 × 1,000。',
  'ib-dag':'BlockDAG架构','ib-dag-t':'多个区块可以同时生产并合并。比线性区块链更高吞吐量、更低延迟。',
  'ib-gas':'真正无Gas','ib-gas-t':'注册和AEQ转账完全免费。不需要ETH、BNB或MATIC。无需银行账户或信用卡。',
  'recent-blocks':'最近区块','blocks-desc':'MERGE = 多个父区块合并（BlockDAG）。TX = 注册交易。区块时间：约6秒。',
  'loading':'加载区块中...','net-info':'网络信息','k-chain':'链名称','k-symbol':'符号','k-btime':'区块时间',
  'k-cons':'共识','k-nodes':'活跃节点','k-storage':'存储','add-mm':'🦊 添加到METAMASK','k-dec':'小数位',
  'btn-add-mm':'+ 添加AEQUITAS网络',
  'phil':'"货币存在是因为人类存在。<br>仅此而已，别无其他。"','phil-sub':'— AEQUITAS原则 —',
  'humans-title':'Aequitas链上的已验证人类',
  'h-what':'什么是已验证人类？','h-what-t':'已验证人类是密码学证明属于独特活人的钱包地址。验证使用3因素硬件系统：R503扫描全部10枚指纹；MAX30102 PPG确认活体脉搏；第2阶段：手静脉红外（1/10⁷）；第3阶段：虹膜（1/10⁷⁸）。仅传输Groth16 ZK证明。任何生物特征数据永不离开设备。',
  'h-zkp':'零知识证明系统','h-zkp-t':'Aequitas在BN128上使用Groth16——与Ethereum和Zcash相同的曲线。约200字节，约10毫秒。commitment = keccak256(iris‖vein‖fingers‖wallet)。Nullifier绑定到身体：丢失手机不会创建第二身份。不存储任何生物特征数据。',
  'h-sybil':'女巫攻击防护','h-sybil-t':'第1阶段：全部10枚指纹 + MAX30102活体性（PPG脉搏，拒绝模型/重放）。第2阶段：手静脉红外——内部身体特征，无法复制，1/10⁷，同卵双胞胎各不相同。第3阶段：虹膜——1/10⁷⁸，全球金标准。Nullifier = keccak256(iris‖vein‖domain)。一个人类，一个钱包，永远。',
  'h-global':'全球金融包容','h-global-t':'无需银行账户、信用卡或加密货币。只需一台带生物特征传感器的安卓手机。',
  'h-bio-hw':'生物特征硬件路线图','h-bio-hw-t':'第1阶段（活跃）：R503光学指纹扫描仪——全部10根手指组合哈希。MAX30102 PPG活体性。第2阶段（计划中）：ESP32-CAM + 850 nm红外LED——手静脉成像，1/10⁷唯一性。第3阶段（计划中）：红外虹膜模块——240+自由度，1/10⁷⁸，完全设备无关，同卵双胞胎各不相同。',
  'reg-humans':'已注册人类','h-desc':'每个地址通过生物特征ZKP验证为独特人类。每人恰好获得1,000 AEQ。永久、不可更改、链上。',
  'no-humans':'尚未注册人类。\n\n下载Aequitas安卓应用，成为链上第一个人类！',
  'reg-stats':'注册统计','total-humans':'总人数',
  'idx-title':'Aequitas指数——实时经济平等评分',
  'idx-desc':'Aequitas指数实时衡量所有经过验证的人类的经济不平等。从链上余额分布的基尼系数导出。0 = 完全平等。100 = 最大不平等。',
  'curr-idx':'当前指数','bar-0':'0 — 完全平等','bar-100':'100 — 最大不平等',
  'gini':'基尼系数','gini-desc':'0 = 平等 · 1 = 不平等',
  'supply-desc':'始终 = 人类 × 1,000 AEQ',
  'phase':'协议阶段','phase-desc':'按人类数量自动推进',
  'humans-desc':'经过生物特征验证的独特人类',
  'pools-title':'再分配池',
  'pools-desc':'每笔兑换费用、滞期费和财富上限溢出自动在四个池之间分配。无需人工干预。所有池每日分配。',
  'vel-pool':'验证者池','vel-pool-desc':'所有费用的40% → 保障网络安全的节点运营商',
  'liq-pool':'流动性池','liq-pool-desc':'所有费用的30% → 流动性提供者，按LP份额比例',
  'ubi-pool':'UBI池','ubi-pool-desc':'所有费用的20% → 所有经过验证的人类均等，每24小时',
  'treasury':'国库','treasury-desc':'所有费用的10% → 协议开发和维护',
  'phases-title':'协议阶段',
  'phases-desc':'阶段转换由人类数量自动触发——无需投票、治理或管理员密钥。',
  'p0':'启动 · &lt;100人类 · 财富上限：50×平均余额 · 当前活跃',
  'p1':'增长 · 100–10,000人类 · 财富上限：20×平均余额',
  'p2':'稳定 · 10,000–100万人类 · 财富上限：10×平均余额',
  'p3':'成熟 · 100万+人类 · 财富上限：3×平均余额 · 最大再分配',
  'wealth-cap-explain':'财富上限设定为所有经过验证的人类当前平均余额的倍数——而非固定数字。随着网络增长自动调整。',
  'demurrage-title':'滞期费——流通激励',
  'demurrage-desc':'Aequitas实施受历史互补货币启发的滞期费机制。闲置AEQ余额缓慢贬值以阻止囤积。',
  'dem-rate-k':'衰减率','dem-rate-v':'每月0.5%（连续，非阶梯式）',
  'dem-grace-k':'宽限期','dem-grace-v':'衰减开始前3个月不活动',
  'dem-reset-k':'计时器重置','dem-reset-v':'任何转账、兑换或流动性操作重置计时器',
  'dem-dest-k':'衰减的AEQ去往','dem-dest-v':'再分配池（40/30/20/10分配）',
  'dem-warn-k':'警告系统','dem-warn-v':'14天通知（一次）+ 每次登录7天重复提醒',
  'story-title':'Aequitas的故事——为何而生',
  'story-text':'<p>2009年。中本聪发布比特币。有史以来第一次，价值可以在不经过银行的情况下在两人之间转移。一场真正的革命。但几乎立刻出现了问题。</p><p>早期矿工以接近零的成本积累了数百万枚代币。到2021年，比特币地址中的前1%控制了90%以上的比特币。比特币的基尼系数超过0.85——高于地球上任何国家。</p><p><span style="color:var(--gold)">Aequitas</span>——拉丁语"公平"和"平等"——的创建是为了回答一个问题：<em style="color:var(--gold)">"如果从第一原则设计一种对每个人都公平的加密货币会是什么样？"</em></p><p>答案很简单：<strong style="color:var(--text)">货币存在是因为人类存在。因此，每个人仅凭成为人类就应该拥有等份的货币。</strong></p><p><em style="color:var(--gold)">"货币存在是因为人类存在。仅此而已，别无其他。"</em></p>',
  'nodes-title':'活跃节点 — 当前网络拓扑','nodes-desc':'Aequitas网络目前在两个地理分布的节点上运行。两者均参与区块生产、状态同步和API服务。通过libp2p点对点通信，通过HTTP同步区块状态。网络设计支持更多节点——任何运营商均可加入。',
  'run-node-title':'运行您自己的节点 — 帮助保护网络',
  'run-node-desc':'任何人都可以运行Aequitas节点——无需许可、无需质押、无需申请。节点参与区块生产并验证人类注册表。节点运营商通过验证者池（每日分配的所有互换费用的40%）赚取协议费用份额。',
  'run-node-title':'运行您自己的节点 — 帮助保护网络',
  'run-node-desc':'任何人都可以运行Aequitas节点——无需许可、无需质押、无需申请。节点参与区块生产并验证人类注册表。节点运营商通过验证者池（每日分配的所有互换费用的40%）赚取协议费用份额。',
  'node1':'节点1 — Railway（主要）','node1-desc':'主要API · 区块生产者 · UBI分配 · P2P引导 · PostgreSQL · MetaMask的RPC',
  'node2':'节点2 — Render（次要）','node2-desc':'次要API · 区块生产者 · P2P对等 · HTTP同步 · 共享PostgreSQL状态',
  'bootstrap-title':'运行自己的节点','bootstrap-desc':'任何人都可以通过运行节点加入Aequitas网络。下载节点指南获取分步说明。',
  'tech-title':'技术规格','mm-config':'MetaMask配置',
  'k-lang':'语言','k-src':'源代码','evm-yes':'是 — JSON-RPC /rpc · MetaMask兼容',
  'proto-label':'Aequitas V7协议——技术文档',
  'ca-title':'合约地址','ca-text':'链：Aequitas Chain（链ID：1926 · 0x786）<br>RPC：https://aequitas.digital/rpc<br><br>BioVerifier：0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7：0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7是整个Aequitas经济体系的唯一真实来源。每一个AEQ余额、每一次人类注册、每一次UBI支付以及每一次财富上限执行，都由这一个不可变合约管理——部署在Aequitas Chain上，这是一个运行BlockDAG共识引擎的定制EVM兼容区块链。没有管理员密钥、没有升级代理、没有任何治理投票能够改变其逻辑中的任何一行代码。今天运行的代码就是十年后运行的代码。<br><br>BioVerifier合约接收完全在用户Android设备上生成的Groth16零知识证明。它在约10毫秒内在链上数学验证新注册者是唯一的活体人类——而不会泄露他们的姓名、身份或生物特征数据。这使得无gas、无需投资的注册成为可能：证明是唯一离开设备的东西。<br><br>这两个合约共同使在历史上任何货币体系中从未存在过的事情成为可能：一种货币供应，其规则——谁获得它、有多少存在、如何重新分配——永远无法被任何人、公司或政府改变。永远。',
  'poa-title':'1. 生存证明 — 非活跃余额恢复','poa-text':'<p>当人们死亡或永久失去行为能力时AEQ会怎样？在比特币中，丢失的钱包意味着永久丢失的供应量。Aequitas通过多阶段非活跃恢复系统解决这个问题：如果一个钱包长时间没有活动，其余额会逐渐通过UBI池返回社区。</p>',
  'poa-box':'第0–2年：正常使用 — 无限制<br>第2年：警告1 — 监护人可以代表回应<br>第2年+60天：警告2 — 紧迫性增加<br>第2年+120天：警告3 — 最终通知<br>第2年+180天：AEQ移至个人托管（仍可恢复）<br>第4年：如果仍不活跃 — 托管释放至UBI池',
  'guard-title':'2. 监护人系统 — 人类安全保障','guard-text':'<p>如果有人住院或因其他原因数月无法访问其设备怎么办？监护人系统允许可信任的人——另一个经过验证的人类——确认钱包所有者仍然活着。监护人拥有严格为零的财务访问权限：只能调用重置非活跃计时器的单一函数。在任何情况下都不能移动、花费或访问资金。</p>',
  'guard-box':'每人1个监护人 · 必须是Aequitas上的经过验证的人类<br>监护人只能调用confirmAlive() — 零交易权限<br>监护人不能移动资金、转移AEQ或访问钱包<br>每个监护人最多3名受监护人 · 分配7天时间锁 · 不允许循环关系',
  'dem-title':'3. 滞期费 — 防囤积机制',
  'dem-box':'费率：3个月非活跃后每月0.5%（连续，非分步）<br>任何转账、互换或流动性操作会自动重置计时器<br>衰减的AEQ重新分配到四个池中 — 从不销毁<br>14天通知显示一次 · 每次活跃会话重复7天提醒',
  'dem-text':'<p>滞期费是货币的持有成本——一种使囤积变得昂贵、流通变得有吸引力的负利率。沃尔格实验（奥地利，1932年）使用滞期费货币在一年内将当地失业率降低了25%。奥地利中央银行正因为它运作得太好而关闭了它。Chiemgauer（德国，2003年）按照相同原则成功运营了20多年。</p>',
  'cap-title':'4. 财富上限 — 数学公平执行','cap-box':'启动上限：max(5,min(N,25))× 平均AEQ余额<br>1–4人：5×（5,000 AEQ）· 每增1人加1× · 25+人：25×（25,000 AEQ）永久<br>适用于除4个协议池外的所有地址<br>超额AEQ立即重新分配 · 无需手动干预',
  'ubi-title':'5. 普遍基本收入 — 每日再分配','ubi-box':'UBI池收入来源：<br>· AEQ↔tUSD AMM池所有互换费用的20%<br>· 财富上限执行的溢出<br>· 非活跃账户的滞期费<br>· 4年后释放的非活跃托管<br><br>分配：每24小时，整个UBI池余额在所有注册的经过验证的人类中平均分配。池重置为零并立即开始从持续的协议活动中重新填充。',
  'inf-title':'6. 无算法通胀 — 固定供应公式','inf-box':'创建新AEQ的唯一事件：新的经过验证的人类注册。<br><br>总供应量 = 经过验证的人类 × 1,000 AEQ<br><br>这不是政策——它由协议执行。没有管理员可以铸造额外的AEQ，没有治理投票可以改变发行，没有预挖矿的创始人分配。AEQ是唯一一种总供应量完全由经过验证的活人数量决定的加密货币。',
  'phases-desc':'阶段边界定义网络增长里程碑。启动阶段（&lt;25名注册人类）财富上限使用滑动乘数：max(5,min(N,25))×平均余额 — 1–4人时为5×，每增加1人加1×，25+人时达到完整25×。防止早期参与者在真正参与形成前集中财富。',
  'p0':'引导期 · 不足100人 · 上限：max(5,min(N,25))×平均 · 滑动5×→25×直至25人 · 当前激活',
  'p1':'增长期 · 100–10,000人 · 财富上限：25×平均余额',
  'p2':'稳定期 · 10,000–1M人 · 财富上限：25×平均余额',
  'p3':'成熟期 · 1M+人 · 财富上限：25×平均余额',
  'wealth-cap-explain':'财富上限在启动阶段动态调整：max(5, min(N, 25)) × 平均余额，N为已注册人类数。1–4人时：5×（5,000 AEQ）。每新增1人多1×。25+人时：永久25×（25,000 AEQ）。防止早期采用者在真实参与形成前过度积累。始终相对于当前平均余额。',
  'btn-download-app':'下载 AEQUITASBIO 应用',
  'swap-title':'🔄 兑换 AEQ ↔ tUSD','swap-sub':'通过原生流动性池将AEQ兑换为tUSD（模拟测试美元）。0.1%手续费仅适用于兑换 — 人与人之间的普通AEQ转账完全免费。',
  'swap-priv-bar':'🔒 仅0.1%兑换费 · AEQ到AEQ转账免费 · tUSD是无实际价值的测试货币',
  'swap-your-aeq':'你的 AEQ','swap-your-tusd':'你的 tUSD','swap-aeq-to-tusd':'AEQ → tUSD','swap-tusd-to-aeq':'tUSD → AEQ',
  'swap-fee-est':'协议手续费 (0.1%)','swap-details-hdr':'兑换详情',
  'swap-out-lbl':'你获得（估算）','swap-impact-lbl':'价格影响','swap-rate-lbl':'汇率',
  'swap-depth-lbl':'池子构成','amm-title':'x × y = k — 恒定乘积 AMM',
  'amm-text':'当你用AEQ兑换tUSD时，AEQ储备增加，tUSD储备减少——它们的乘积始终等于k。更大的兑换造成更大的价格影响。0.1%手续费在应用公式前从输入中扣除。',
  'swap-btn-conn':'🦊 连接 METAMASK','swap-btn-go':'🔄 兑换',
  'swap-log-hint':'// 连接钱包以兑换...',
  'swap-no-liquidity':'还没有 tUSD？','swap-faucet-desc':'已注册的人类可以申领一次测试 tUSD','swap-btn-faucet':'💧 申领测试 tUSD',
  'swap-addliq-title':'提供流动性','swap-addliq-desc':'成为第一个存款者 — 你的比例设定起始价格。','swap-btn-addliq':'💧 添加流动性',
  'swap-lp-title':'你的 LP 仓位','swap-lp-share':'池子份额','swap-lp-withdrawable':'可提取',
  'swap-lp-pct-label':'% 你的仓位','swap-lp-youget':'你将收到','swap-btn-removeliq':'🔥 移除流动性',
  'swap-pool-title':'AEQ / tUSD — 池子状态',
  'swap-pool-aeq':'AEQ 储备','swap-pool-tusd':'tUSD 储备','swap-pool-price':'现货价格',
  'swap-fee-bps':'兑换手续费','swap-fee-split':'手续费分配','swap-fee-split-v':'40% 验证者 / 30% LP / 20% UBI / 10% 国库',
  'swap-pools-addr-title':'代币经济池地址',
  'swap-validators':'验证者 (40%)','swap-lps':'流动性提供者 (30%)','swap-ubi':'UBI 池 (20%)','swap-treasury':'国库 (10%)',
  'ubi-hero-title':'普遍基本收入 — UBI 池',
  'ubi-hero-sub':'累积中 — 下次平等分配给所有验证人类：',
  'ubi-bal-lbl':'当前池余额','ubi-hero-desc':'在所有验证人类中平等分配 · 每24小时支付 · 支付后池归零 · 无最低余额要求',
  'ubi-how-fills':'UBI 池如何填充',
  'ubi-src-swap':'兑换手续费','ubi-src-swap-d':'每次AEQ↔tUSD兑换贡献其0.1%手续费的20%。更多交易 = 更快填充。',
  'ubi-src-dem':'滞期费','ubi-src-dem-d':'不活跃AEQ（3+个月）以0.5%/月衰减。衰减金额的20%进入UBI。',
  'ubi-src-cap':'财富上限溢出','ubi-src-cap-d':'超过max(5,min(N,25))×平均余额的钱包立即被没收超额部分。20%立即流入UBI。',
  'pools4-header':'所有四个再分配池',
  'ubi-see-above':'见上方倒计时','ubi-timer-above':'⏰ 倒计时显示在上方','pool-t-timer':'累积中 — 无计时器',
  'usp-headline':'历史上首次 — 所有人在平等条件下起步',
  'usp-sub':'只需拥有一部Android智能手机即可参与。无需银行账户，无需加密货币知识，无需任何投资。',
  'usp-c1-title':'0元启动投资','usp-c1-desc':'注册完全免gas。无需ETH、无需MATIC、无需信用卡。协议代您支付所有交易费用。',
  'usp-c2-title':'每人1,000 AEQ','usp-c2-desc':'亿万富翁还是贫困农民——每人恰好获得1,000 AEQ。不多不少。平等起点，数学保证。',
  'usp-c3-title':'人人可参与','usp-c3-desc':'无需银行账户、信用卡或身份证件。注册使用经济实惠的生物识别硬件套件（指纹扫描仪+脉搏传感器，约$15）——为全球访问而设计。',
  'usp-c4-title':'永久每日UBI','usp-c4-desc':'注册后，您每天自动获得UBI支付份额——每天，无需任何操作。',
  'v7-intro-title':'什么是 AequitasV7？',
  'v7-intro-text':'AequitasV7是Aequitas协议的核心智能合约。"V7"指公平合约的第7个主要版本。它不可更改地部署在Aequitas Chain（链ID 1926）上，处理所有方面：人类注册、ZK证明验证、余额管理、财富上限、UBI分配、兑换手续费。没有管理员可以升级它。六个机制形成自我强化系统。',
  'explore-title':'探索 Aequitas',
  'expl-score':'平等指数','expl-score-d':'实时基尼系数 · Aequitas指数 · 实时财富分配',
  'expl-economy':'UBI与再分配池','expl-economy-d':'每日UBI倒计时 · 4个链上池 · 货币持有税 · 协议阶段',
  'expl-charts':'图表与历史','expl-charts-d':'基尼历史 · 洛伦兹曲线 · 财富上限启动滑块 · Aequitas的故事',
  'expl-v7':'协议V7文档','expl-v7-d':'AequitasV7合约 · 6个机制 · ZK证明 · 财富上限 · 货币持有税 · 不可更改代码',
  'expl-explorer':'区块浏览器','expl-explorer-d':'实时BlockDAG · 点击任意区块查看验证者、哈希、交易、父哈希',
  'swap-sell-label':'卖出','swap-receive-label':'接收',
  'gini-calc-title':'如何计算Aequitas指数','gini-calc-text':'收集所有AEQ余额。公式计算每对余额之间的平均绝对差，结果0-1乘以100=Aequitas指数。','gini-why-title':'为什么选择基尼系数','gini-why-text':'基尼系数捕捉所有已验证人类的完整分布。Aequitas将此数据发布在链上。','expl-network':'网络与节点','expl-network-d':'节点拓扑 · 运行自己的节点 · 技术规格 · Chain ID 1926'
},
id:{
  'logo-sub':'BUKTI KEMANUSIAAN','live':'LANGSUNG',
  'tab-register':'🔐 Daftar','tab-explorer':'🔍 Penjelajah','tab-humans':'👥 Manusia','tab-index':'📊 Indeks','tab-network':'🌐 Jaringan','tab-protocol':'📜 Protokol V7','tab-swap':'🔄 Tukar',
  'reg-title':'🔐 Daftar sebagai Manusia Terverifikasi',
  'reg-sub':'Bergabunglah dengan jaringan Aequitas dan terima hibah Pendapatan Dasar Universal sebesar 1.000 AEQ. Satu kali, permanen, dan sepenuhnya gratis. Tidak ada data pribadi yang pernah disimpan.',
  'app-title':'PENDAFTARAN HANYA MELALUI APLIKASI ANDROID',
  'app-text':'Bukti Kemanusiaan menggunakan sistem biometrik fisik 3 faktor. Fase 1: sensor optik R503 memindai semua 10 sidik jari + MAX30102 PPG mengonfirmasi denyut nadi hidup. Fase 2: kamera IR vena tangan (keunikan 1 dari 10⁷). Fase 3: pemindaian iris — standar emas, 1 dari 10⁷⁸, sepenuhnya independen dari perangkat. Bukti ZK Groth16 mengkomit semua faktor tanpa mengungkapkan data biometrik apa pun. 1.000 AEQ Anda dikreditkan otomatis setelah verifikasi.',
  's1t':'Pemindaian Biometrik','s1d':'AequitasBio memindai semua 10 sidik jari (sensor optik R503) + denyut nadi MAX30102 PPG mengonfirmasi keaktifan. Fase 2: IR vena tangan. Fase 3: iris. Data mentah tidak pernah meninggalkan perangkat.',
  's2t':'Pembuatan Bukti ZK','s2d':'Bukti ZK Groth16 mengkomit semua faktor biometrik: commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier terikat ke tubuh, bukan ponsel — kehilangan perangkat tidak dapat membuat identitas kedua.',
  's3t':'Hubungkan Dompet','s3d':'Aplikasi membuka MetaMask di halaman ini · hubungkan dompet Ethereum Anda · bukti terikat secara kriptografis ke alamat Anda',
  's4t':'1.000 AEQ Dikreditkan','s4d':'Pendaftaran dikonfirmasi di BlockDAG Aequitas dalam 6 detik · 1.000 AEQ dikreditkan seketika · identitas Anda dicatat permanen sebagai manusia terverifikasi',
  'priv-bar':'🔒 R503 10 Sidik Jari · MAX30102 Keaktifan · Fase 2: IR Vena Tangan · Fase 3: Iris (10⁷⁸) · Groth16 ZKP · Data tidak pernah meninggalkan perangkat · Satu manusia · Selamanya',
  'conn-wallet':'DOMPET TERHUBUNG','proof-recv':'⚡ BUKTI ZK DITERIMA','proof-hint':'Hubungkan dompet untuk mendaftar',
  'btn-conn':'🦊 HUBUNGKAN METAMASK','btn-reg':'🔐 DAFTAR ON-CHAIN',
  'btn-web-reg':'🌐 DAFTAR VIA BROWSER (WebAuthn)',
  'web-reg-warn':'⚠ Terikat perangkat: Identitas ini terikat pada perangkat dan browser ini. Tidak dapat dipindahkan ke perangkat lain. Untuk identitas permanen multi-perangkat, gunakan Aplikasi Android Aequitas.',
  'reg-log-hint':'// Buka Aplikasi Android Aequitas untuk membuat bukti Anda, lalu kembali ke sini...',
  'reg-details':'Detail Pendaftaran','k-network':'Jaringan','k-chainid':'ID Rantai','k-grant':'Hibah UBI',
  'k-fee':'Biaya Gas','free':'GRATIS — sepenuhnya tanpa gas','k-limit':'Pendaftaran','k-limit-v':'Satu kali · permanen · tidak dapat diubah',
  'k-bio':'Data Biometrik','never-stored':'Tidak pernah disimpan — tetap di perangkat Anda',
  'k-proof':'Sistem Bukti','k-conf':'Konfirmasi','k-conf-v':'Dalam 6 detik (1 blok)',
  'k-sybil':'Perlindungan Sybil','k-sybil-v':'Satu identitas per biometrik · kunci permanen',
  'live-stats':'Statistik Rantai Langsung',
  's-height':'Tinggi Blok','s-height-sub':'Blok baru setiap ~6d · BlockDAG · Produksi paralel',
  's-humans':'Manusia Terverifikasi','s-humans-sub':'ZKP biometrik · Satu orang, satu dompet, selamanya',
  's-supply':'Total Pasokan','s-supply-sub':'Selalu = Manusia × 1.000 AEQ',
  's-index':'Indeks Aequitas','s-index-sub':'0 = kesetaraan sempurna · 100 = ketidaksetaraan maksimum',
  's-uptime':'Waktu Aktif','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Bukti Kemanusiaan','ib-poh-t':'Setiap pemegang AEQ harus membuktikan secara kriptografis bahwa mereka adalah manusia hidup yang unik. Tidak ada bot, korporasi, AI. Data biometrik tidak pernah meninggalkan perangkat.',
  'ib-fair':'Distribusi yang Benar-benar Adil','ib-fair-t':'Setiap manusia terverifikasi menerima tepat 1.000 AEQ saat pendaftaran. Tanpa pre-mining, tanpa alokasi pendiri. Total pasokan selalu sama dengan manusia terverifikasi × 1.000.',
  'ib-dag':'Arsitektur BlockDAG','ib-dag-t':'Beberapa blok dapat diproduksi secara bersamaan dan digabungkan. Throughput lebih tinggi, latensi lebih rendah.',
  'ib-gas':'Benar-benar Tanpa Gas','ib-gas-t':'Pendaftaran dan transfer AEQ tidak memerlukan biaya. Tidak perlu ETH, BNB, atau MATIC. Tidak perlu rekening bank atau kartu kredit.',
  'recent-blocks':'Blok Terbaru','blocks-desc':'MERGE = beberapa induk digabung (BlockDAG). TX = transaksi pendaftaran. Waktu blok: ~6 detik.',
  'loading':'Memuat blok...','net-info':'Informasi Jaringan','k-chain':'Nama Rantai','k-symbol':'Simbol','k-btime':'Waktu Blok',
  'k-cons':'Konsensus','k-nodes':'Node Aktif','k-storage':'Penyimpanan','add-mm':'🦊 TAMBAHKAN KE METAMASK','k-dec':'Desimal',
  'btn-add-mm':'+ TAMBAHKAN JARINGAN AEQUITAS',
  'phil':'"Uang ada karena manusia ada.<br>Tidak lebih, tidak kurang."','phil-sub':'— PRINSIP AEQUITAS —',
  'humans-title':'Manusia Terverifikasi di Aequitas Chain',
  'h-what':'Apa itu Manusia Terverifikasi?','h-what-t':'Manusia Terverifikasi adalah alamat dompet yang terbukti secara kriptografis milik manusia hidup yang unik. Verifikasi menggunakan sistem perangkat keras 3 faktor: R503 memindai semua 10 sidik jari; MAX30102 PPG mengonfirmasi denyut nadi hidup; Fase 2: IR vena tangan (1 dari 10⁷); Fase 3: iris (1 dari 10⁷⁸). Hanya bukti ZK Groth16 yang ditransmisikan. Tidak ada data biometrik yang meninggalkan perangkat.',
  'h-zkp':'Sistem Bukti ZK','h-zkp-t':'Aequitas menggunakan Groth16 pada BN128 — kurva yang sama dengan Ethereum dan Zcash. ~200 byte, ~10ms. commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier terikat ke tubuh: kehilangan ponsel tidak membuat identitas kedua. Tidak ada data biometrik yang pernah disimpan.',
  'h-sybil':'Pencegahan Serangan Sybil','h-sybil-t':'Fase 1: semua 10 sidik jari + keaktifan MAX30102 (denyut PPG, menolak cetakan/ulangan). Fase 2: IR vena tangan — fitur tubuh internal, tidak mungkin disalin, 1 dari 10⁷, berbeda pada kembar identik. Fase 3: iris — 1 dari 10⁷⁸, standar emas global. Nullifier = keccak256(iris‖vein‖domain). Satu manusia, satu dompet, selamanya.',
  'h-global':'Inklusi Keuangan Global','h-global-t':'Tidak perlu rekening bank, kartu kredit, atau cryptocurrency sebelumnya. Hanya smartphone Android dengan sensor biometrik.',
  'h-bio-hw':'Peta Jalan Perangkat Keras Biometrik','h-bio-hw-t':'Fase 1 (aktif): pemindai sidik jari optik R503 — hash gabungan semua 10 jari. Keaktifan MAX30102 PPG. Fase 2 (direncanakan): ESP32-CAM + LED IR 850 nm — pencitraan vena tangan, keunikan 1 dari 10⁷. Fase 3 (direncanakan): modul iris IR — 240+ derajat kebebasan, 1 dari 10⁷⁸, sepenuhnya independen dari perangkat, kembar identik berbeda.',
  'reg-humans':'Manusia Terdaftar','h-desc':'Setiap alamat diverifikasi sebagai manusia unik melalui ZKP biometrik. Masing-masing menerima tepat 1.000 AEQ. Permanen, tidak dapat diubah, on-chain.',
  'no-humans':'Belum ada manusia terdaftar.\n\nUnduh Aplikasi Android Aequitas dan jadilah yang pertama!',
  'reg-stats':'Statistik Registri','total-humans':'Total Manusia',
  'idx-title':'Indeks Aequitas — Skor Kesetaraan Ekonomi Real-Time',
  'idx-desc':'Indeks Aequitas mengukur ketidaksetaraan ekonomi semua manusia terverifikasi secara real-time. Diturunkan dari koefisien Gini distribusi saldo on-chain. 0 = kesetaraan sempurna. 100 = ketidaksetaraan maksimum.',
  'curr-idx':'Indeks Saat Ini','bar-0':'0 — Kesetaraan Sempurna','bar-100':'100 — Maks. Ketidaksetaraan',
  'gini':'Koefisien Gini','gini-desc':'0 = setara · 1 = tidak setara',
  'supply-desc':'Selalu = Manusia × 1.000 AEQ',
  'phase':'Fase Protokol','phase-desc':'Otomatis berdasarkan jumlah manusia',
  'humans-desc':'Manusia unik yang terverifikasi secara biometrik',
  'pools-title':'Pool Redistribusi',
  'pools-desc':'Setiap biaya swap, biaya demurrage, dan kelebihan batas kekayaan secara otomatis dibagi ke empat pool. Tanpa intervensi manual. Semua pool membayar setiap hari.',
  'vel-pool':'Pool Validator','vel-pool-desc':'40% semua biaya → operator node yang mengamankan jaringan',
  'liq-pool':'Pool Likuiditas','liq-pool-desc':'30% semua biaya → penyedia likuiditas, proporsional dengan saham LP',
  'ubi-pool':'Pool UBI','ubi-pool-desc':'20% semua biaya → semua manusia terverifikasi secara merata, setiap 24 jam',
  'treasury':'Perbendaharaan','treasury-desc':'10% semua biaya → pengembangan dan pemeliharaan protokol',
  'phases-title':'Fase Protokol',
  'demurrage-title':'Demurrage — Insentif untuk Bersirkulasi',
  'demurrage-desc':'Aequitas mengimplementasikan mekanisme demurrage yang terinspirasi dari mata uang komplementer historis. Saldo AEQ yang tidak aktif perlahan kehilangan nilai untuk mencegah penimbunan.',
  'dem-rate-k':'Tingkat Peluruhan','dem-rate-v':'0,5% per bulan (berkelanjutan, tidak bertahap)',
  'dem-grace-k':'Masa Tenggang','dem-grace-v':'3 bulan tidak aktif sebelum peluruhan dimulai',
  'dem-reset-k':'Reset Timer','dem-reset-v':'Setiap transfer, swap, atau tindakan likuiditas mereset timer',
  'dem-dest-k':'AEQ yang meluruh pergi ke','dem-dest-v':'Pool redistribusi (pembagian 40/30/20/10)',
  'dem-warn-k':'Sistem Peringatan','dem-warn-v':'Pemberitahuan 14 hari (sekali) + pengingat 7 hari berulang setiap login',
  'story-title':'Kisah Aequitas — Mengapa Ini Ada',
  'story-text':'<p>Tahun 2009. Satoshi Nakamoto merilis Bitcoin. Untuk pertama kalinya, nilai dapat ditransfer antara dua orang tanpa bank. Sebuah revolusi sejati. Tetapi hampir segera sesuatu yang salah terjadi.</p><p>Para penambang awal mengumpulkan jutaan koin dengan biaya hampir nol. Pada 2021, 1% teratas alamat Bitcoin menguasai lebih dari 90% semua Bitcoin. Koefisien Gini Bitcoin melebihi 0,85 — lebih tinggi dari negara mana pun di Bumi.</p><p><span style="color:var(--gold)">Aequitas</span> — Latin untuk "keadilan" dan "kesetaraan" — diciptakan untuk menjawab: <em style="color:var(--gold)">"Seperti apa cryptocurrency yang dirancang dari prinsip pertama untuk adil bagi setiap manusia?"</em></p><p>Jawabannya sederhana: <strong style="color:var(--text)">Uang ada karena manusia ada. Oleh karena itu, setiap orang harus memiliki bagian yang sama dari uang hanya karena menjadi manusia.</strong></p><p><em style="color:var(--gold)">"Uang ada karena manusia ada. Tidak lebih, tidak kurang."</em></p>',
  'nodes-title':'Node Aktif — Topologi Jaringan Saat Ini','nodes-desc':'Jaringan Aequitas saat ini beroperasi pada dua node yang tersebar secara geografis. Keduanya berpartisipasi dalam produksi blok, sinkronisasi status, dan layanan API. Jaringan dirancang untuk mendukung node tambahan — operator mana pun dapat bergabung.',
  'run-node-title':'Jalankan Node Anda Sendiri — Bantu Amankan Jaringan',
  'run-node-desc':'Siapa pun dapat menjalankan node Aequitas — tanpa izin, tanpa stake, tanpa pendaftaran. Node berpartisipasi dalam produksi blok dan memvalidasi registri manusia. Operator node mendapatkan bagian biaya protokol melalui Pool Validator (40% semua biaya swap, didistribusikan setiap hari).',
  'node1':'Node 1 — Railway (Utama)','node1-desc':'API utama · Produsen blok · Distribusi UBI · Bootstrap P2P · PostgreSQL · RPC untuk MetaMask',
  'node2':'Node 2 — Render (Sekunder)','node2-desc':'API sekunder · Produsen blok · Peer P2P · Sinkronisasi HTTP · Status PostgreSQL bersama',
  'bootstrap-title':'Jalankan Node Anda Sendiri','bootstrap-desc':'Siapa pun dapat bergabung dengan jaringan Aequitas dengan menjalankan node. Unduh panduan node untuk instruksi langkah demi langkah.',
  'tech-title':'Spesifikasi Teknis','mm-config':'Konfigurasi MetaMask',
  'k-lang':'Bahasa','k-src':'Kode Sumber','evm-yes':'Ya — JSON-RPC /rpc · Kompatibel MetaMask',
  'proto-label':'Protokol Aequitas V7 — Dokumentasi Teknis',
  'ca-title':'Alamat Kontrak','ca-text':'Rantai: Aequitas Chain (ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7: 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7 adalah satu-satunya sumber kebenaran untuk seluruh ekonomi Aequitas. Setiap saldo AEQ, setiap registrasi manusia, setiap pembayaran UBI, dan setiap penegakan batas kekayaan diatur oleh satu kontrak yang tidak dapat diubah ini — dikerahkan di Aequitas Chain, blockchain khusus yang kompatibel dengan EVM yang menjalankan mesin konsensus BlockDAG. Tidak ada kunci admin, tidak ada proxy upgrade, tidak ada pemungutan suara tata kelola yang dapat mengubah satu baris pun logikanya. Kode yang berjalan hari ini adalah kode yang akan berjalan sepuluh tahun lagi.<br><br>Kontrak BioVerifier menerima bukti zero-knowledge Groth16 yang dihasilkan sepenuhnya di perangkat Android pengguna. Ia memverifikasi secara matematis on-chain dalam ~10 ms bahwa pendaftar baru adalah manusia hidup yang unik — tanpa pernah mengetahui nama, identitas, atau data biometrik mereka. Inilah yang membuat registrasi tanpa gas dan tanpa investasi menjadi mungkin: bukti adalah satu-satunya hal yang pernah meninggalkan perangkat.<br><br>Bersama-sama, dua kontrak ini memungkinkan sesuatu yang belum pernah ada dalam sistem mata uang manapun dalam sejarah: pasokan uang yang aturannya — siapa yang mendapatkannya, berapa banyak yang ada, bagaimana redistribusinya — tidak dapat diubah oleh siapapun, perusahaan manapun, atau pemerintah manapun. Selamanya.',
  'poa-title':'1. BUKTI KEHIDUPAN — Pemulihan Saldo Tidak Aktif','poa-text':'<p>Apa yang terjadi dengan AEQ ketika orang meninggal atau menjadi tidak mampu secara permanen? Di Bitcoin, dompet yang hilang berarti pasokan yang hilang selamanya. Aequitas menyelesaikan ini melalui sistem pemulihan ketidakaktifan multi-tahap: jika dompet tidak menunjukkan aktivitas untuk jangka waktu yang lama, saldonya secara bertahap dikembalikan ke komunitas melalui pool UBI.</p>',
  'poa-box':'Tahun 0–2: Penggunaan normal — tanpa batasan<br>Tahun 2: Peringatan 1 — Guardian dapat merespons atas nama<br>Tahun 2+60h: Peringatan 2 — urgensi meningkat<br>Tahun 2+120h: Peringatan 3 — pemberitahuan terakhir<br>Tahun 2+180h: AEQ dipindahkan ke ESCROW pribadi (masih dapat dipulihkan)<br>Tahun 4: Jika masih tidak aktif — ESCROW dirilis ke Pool UBI',
  'guard-title':'2. SISTEM GUARDIAN — Perlindungan Manusia','guard-text':'<p>Bagaimana jika seseorang dirawat di rumah sakit atau tidak dapat mengakses perangkatnya selama berbulan-bulan? Sistem Guardian memungkinkan orang terpercaya — manusia terverifikasi lainnya — mengonfirmasi bahwa pemilik dompet masih hidup. Guardian memiliki nol akses keuangan: hanya dapat memanggil satu fungsi yang mereset timer ketidakaktifan. Tidak dapat memindahkan, membelanjakan, atau mengakses dana dalam keadaan apapun.</p>',
  'guard-box':'1 Guardian per manusia · harus manusia terverifikasi di Aequitas<br>Guardian HANYA dapat memanggil confirmAlive() — nol hak transaksi<br>Guardian TIDAK DAPAT memindahkan dana, mentransfer AEQ, atau mengakses dompet<br>Maksimal 3 wali per Guardian · Kunci waktu 7 hari · Tanpa hubungan melingkar',
  'dem-title':'3. DEMURRAGE — Mekanisme Anti-Penimbunan',
  'dem-box':'Tingkat: 0,5%/bulan setelah 3 bulan ketidakaktifan (berkelanjutan, tidak bertahap)<br>Timer direset secara otomatis dengan transfer, swap, atau tindakan likuiditas apapun<br>AEQ yang meluruh didistribusikan ulang ke empat pool — tidak pernah dibakar<br>Pemberitahuan 14 hari ditampilkan sekali · 7 hari diulang di setiap sesi aktif',
  'dem-text':'<p>Demurrage adalah biaya kepemilikan uang — suku bunga negatif yang membuat penimbunan mahal dan sirkulasi menarik. Eksperimen Wörgl (Austria, 1932) mengurangi pengangguran lokal 25% dalam satu tahun. Bank Sentral Austria menutupnya justru karena bekerja terlalu baik. Chiemgauer (Jerman, 2003) beroperasi dengan prinsip yang sama dengan sukses selama lebih dari 20 tahun.</p>',
  'cap-title':'4. BATAS KEKAYAAN — Penerapan Keadilan Matematis','cap-box':'Batas bootstrap: max(5,min(N,25))× saldo rata-rata saat ini<br>1–4 manusia: 5× · +1× per manusia · 25+: 25× permanen<br>Berlaku untuk SEMUA alamat kecuali 4 pool protokol<br>Kelebihan AEQ langsung didistribusikan ulang · Tanpa intervensi manual',
  'ubi-title':'5. PENDAPATAN DASAR UNIVERSAL — Redistribusi Harian','ubi-box':'Sumber pendapatan Pool UBI:<br>· 20% semua biaya swap dari pool AMM AEQ↔tUSD<br>· Overflow dari penerapan batas kekayaan<br>· Biaya demurrage dari akun tidak aktif<br>· Escrow tidak aktif dirilis setelah 4 tahun<br><br>Distribusi: Setiap 24 jam, seluruh saldo pool UBI dibagi rata di antara semua manusia terverifikasi yang terdaftar. Pool direset ke nol dan segera mulai diisi ulang dari aktivitas protokol yang berkelanjutan.',
  'inf-title':'6. TANPA INFLASI ALGORITMIK — Formula Pasokan Tetap','inf-box':'SATU-SATUNYA peristiwa yang menciptakan AEQ baru: manusia terverifikasi baru mendaftar.<br><br>Total Pasokan = Manusia Terverifikasi × 1.000 AEQ<br><br>Ini bukan kebijakan — ini diterapkan oleh protokol. Tidak ada admin yang dapat mencetak AEQ tambahan, tidak ada suara tata kelola yang dapat mengubah penerbitan. AEQ adalah satu-satunya cryptocurrency di mana total pasokan ditentukan semata-mata oleh jumlah manusia hidup yang terverifikasi.',
  'phases-desc':'Pada Fase 0, batas kekayaan menggunakan pengganda bootstrap: max(5, min(N, 25))× saldo rata-rata. Dengan 1–4 manusia: 5× rata-rata. Setiap manusia baru menambah 1×. Pada 25+ manusia: terkunci permanen di 25×. Fase 1+ mempertahankan 25× tetap. Semua transisi otomatis — tanpa pemungutan suara, tanpa kunci admin.',
  'p0':'Bootstrap · &lt;100 manusia · Batas Kekayaan: max(5,min(N,25))× rata-rata · Meluncur 5×→25× hingga manusia ke-25 · Saat ini aktif',
  'p1':'Pertumbuhan · 100–10.000 manusia · Batas Kekayaan: 25× saldo rata-rata',
  'p2':'Stabilitas · 10.000–1M manusia · Batas Kekayaan: 25× saldo rata-rata',
  'p3':'Kematangan · 1M+ manusia · Batas Kekayaan: 25× saldo rata-rata',
  'wealth-cap-explain':'Batas Kekayaan pada Fase 0 (Bootstrap) menggunakan max(5, min(N, 25))× saldo AEQ rata-rata, di mana N = manusia terdaftar. 1–4 manusia: 5× rata-rata. Setiap manusia baru menambah 1×. 25+ manusia: terkunci permanen di 25×. Batas selalu mengikuti saldo rata-rata saat ini.',
  'btn-download-app':'UNDUH APLIKASI AEQUITASBIO',
  'swap-title':'🔄 Tukar AEQ ↔ tUSD','swap-sub':'Tukarkan AEQ dengan tUSD (dolar uji simulasi) melalui pool likuiditas asli. Biaya 0,1% hanya berlaku untuk pertukaran — transfer AEQ biasa antar orang tetap sepenuhnya gratis.',
  'swap-priv-bar':'🔒 Hanya 0,1% biaya swap · Transfer AEQ-ke-AEQ gratis · tUSD adalah mata uang uji tanpa nilai nyata',
  'swap-your-aeq':'AEQ Anda','swap-your-tusd':'tUSD Anda','swap-aeq-to-tusd':'AEQ → tUSD','swap-tusd-to-aeq':'tUSD → AEQ',
  'swap-fee-est':'Biaya protokol (0,1%)','swap-details-hdr':'Detail Pertukaran',
  'swap-out-lbl':'Anda terima (est.)','swap-impact-lbl':'Dampak harga','swap-rate-lbl':'Nilai tukar',
  'swap-depth-lbl':'Komposisi Pool','amm-title':'x × y = k — AMM Produk Konstan',
  'amm-text':'Saat Anda menukar AEQ dengan tUSD, cadangan AEQ bertambah dan cadangan tUSD berkurang — produknya selalu sama dengan k. Pertukaran lebih besar menyebabkan dampak harga lebih besar. Biaya 0,1% dipotong sebelum rumus diterapkan.',
  'swap-btn-conn':'🦊 HUBUNGKAN METAMASK','swap-btn-go':'🔄 TUKAR',
  'swap-log-hint':'// Hubungkan dompet untuk menukar...',
  'swap-no-liquidity':'Belum punya tUSD?','swap-faucet-desc':'Manusia terdaftar dapat klaim tUSD uji sekali','swap-btn-faucet':'💧 KLAIM tUSD UJI',
  'swap-addliq-title':'Sediakan Likuiditas','swap-addliq-desc':'Jadilah yang pertama menyetor — rasio Anda menetapkan harga awal.','swap-btn-addliq':'💧 TAMBAH LIKUIDITAS',
  'swap-lp-title':'Posisi LP Anda','swap-lp-share':'Bagian Pool','swap-lp-withdrawable':'Dapat Ditarik',
  'swap-lp-pct-label':'% posisi Anda','swap-lp-youget':'Anda akan terima','swap-btn-removeliq':'🔥 HAPUS LIKUIDITAS',
  'swap-pool-title':'AEQ / tUSD — Status Pool',
  'swap-pool-aeq':'Cadangan AEQ','swap-pool-tusd':'Cadangan tUSD','swap-pool-price':'Harga Spot',
  'swap-fee-bps':'Biaya Swap','swap-fee-split':'Distribusi biaya','swap-fee-split-v':'40% Validator / 30% LP / 20% UBI / 10% Perbendaharaan',
  'swap-pools-addr-title':'Alamat Pool Tokenomik',
  'swap-validators':'Validator (40%)','swap-lps':'Penyedia Likuiditas (30%)','swap-ubi':'Pool UBI (20%)','swap-treasury':'Perbendaharaan (10%)',
  'ubi-hero-title':'PENDAPATAN DASAR UNIVERSAL — POOL UBI',
  'ubi-hero-sub':'Mengumpulkan — pembayaran berikutnya dibagikan merata ke semua manusia terverifikasi dalam:',
  'ubi-bal-lbl':'saldo pool saat ini','ubi-hero-desc':'Dibagi merata di antara semua · dibayar setiap 24j · pool direset ke nol · tidak perlu saldo minimum',
  'ubi-how-fills':'Bagaimana Pool UBI terisi',
  'ubi-src-swap':'Biaya Swap','ubi-src-swap-d':'Setiap swap AEQ↔tUSD berkontribusi 20% dari biaya 0,1%-nya. Lebih banyak trading = pengisian lebih cepat.',
  'ubi-src-dem':'Demurrage','ubi-src-dem-d':'AEQ tidak aktif (3+ bulan) berkurang 0,5%/bulan. 20% jumlah yang berkurang masuk ke UBI.',
  'ubi-src-cap':'Overflow Batas Kekayaan','ubi-src-cap-d':'Dompet yang melebihi batas kekayaan (max(5,min(N,25))× rata-rata) langsung disita kelebihannya. 20% mengalir ke UBI segera.',
  'pools4-header':'Keempat pool redistribusi',
  'ubi-see-above':'lihat hitung mundur di atas','ubi-timer-above':'⏰ hitung mundur ditampilkan di atas','pool-t-timer':'Mengumpulkan — tanpa timer',
  'usp-headline':'Untuk pertama kalinya dalam sejarah — semua memulai dengan setara',
  'usp-sub':'Jika Anda memiliki smartphone Android, Anda memenuhi syarat. Tanpa bank, tanpa pengetahuan kripto, tanpa investasi.',
  'usp-c1-title':'Investasi Awal 0,00','usp-c1-desc':'Pendaftaran sepenuhnya tanpa gas. Tanpa ETH, tanpa MATIC, tanpa kartu kredit. Protokol membayar semua biaya atas nama Anda.',
  'usp-c2-title':'1.000 AEQ untuk setiap manusia','usp-c2-desc':'Miliarder atau petani subsisten — semua mendapat tepat 1.000 AEQ. Tidak lebih, tidak kurang. Start setara, dijamin matematika.',
  'usp-c3-title':'Dapat diakses semua orang','usp-c3-desc':'Tanpa rekening bank, kartu kredit, atau dokumen ID. Pendaftaran menggunakan kit biometrik terjangkau (pemindai sidik jari + sensor denyut, ~$15) — dirancang untuk akses global.',
  'usp-c4-title':'UBI harian selamanya','usp-c4-desc':'Setelah terdaftar, Anda secara otomatis menerima bagian harian dari pembayaran UBI — setiap hari, tanpa tindakan apa pun.',
  'v7-intro-title':'Apa itu AequitasV7?',
  'v7-intro-text':'AequitasV7 adalah kontrak pintar inti dari protokol Aequitas. "V7" mengacu pada versi utama ke-7 dari kontrak keadilan. Dikerahkan secara tidak dapat diubah di Aequitas Chain (ID 1926) dan menangani setiap aspek: pendaftaran manusia, verifikasi ZK, manajemen saldo, batas kekayaan, distribusi UBI, biaya swap. Tidak ada admin yang dapat memperbaruinya. Keenam mekanisme membentuk sistem yang saling memperkuat.',
  'explore-title':'Jelajahi Aequitas',
  'expl-score':'Skor Kesetaraan','expl-score-d':'Koefisien Gini langsung · Indeks Aequitas · distribusi kekayaan secara real time',
  'expl-economy':'UBI &amp; Pool Redistribusi','expl-economy-d':'Hitung mundur UBI harian · 4 pool on-chain · demurrage · Fase Protokol',
  'expl-charts':'Grafik &amp; Riwayat','expl-charts-d':'Riwayat Gini · kurva Lorenz · slider bootstrap batas kekayaan · Kisah Aequitas',
  'expl-v7':'Dokumentasi Protokol V7','expl-v7-d':'Kontrak AequitasV7 · 6 mekanisme · bukti ZK · batas kekayaan · demurrage · kode tak berubah',
  'expl-explorer':'Block Explorer','expl-explorer-d':'BlockDAG langsung · klik blok apapun untuk melihat validator, hash, transaksi, hash induk',
  'swap-sell-label':'Jual','swap-receive-label':'Terima',
  'gini-calc-title':'Bagaimana Indeks Aequitas dihitung','gini-calc-text':'Semua saldo AEQ dikumpulkan. Rumus menghitung perbedaan absolut rata-rata dinormalisasi dengan n2. Hasil 0-1 dikali 100 = Indeks Aequitas.','gini-why-title':'Mengapa Gini','gini-why-text':'Koefisien Gini menangkap distribusi lengkap semua manusia terverifikasi.','expl-network':'Jaringan &amp; Node','expl-network-d':'Topologi node · jalankan node sendiri · spesifikasi teknis · Chain ID 1926'
},
it:{
  'logo-sub':'PROVA DI UMANITÀ','live':'LIVE',
  'tab-register':'🔐 Registrati','tab-explorer':'🔍 Explorer','tab-humans':'👥 Umani','tab-index':'📊 Indice','tab-network':'🌐 Rete','tab-protocol':'📜 Protocollo V7','tab-swap':'🔄 Scambia',
  'reg-title':'🔐 Registrati come Umano Verificato',
  'reg-sub':'Unisciti alla rete Aequitas e ricevi il tuo sussidio di Reddito Universale di Base di 1.000 AEQ. Una tantum, permanente e completamente gratuito. Nessun dato personale viene mai memorizzato.',
  'app-title':'REGISTRAZIONE SOLO VIA APP ANDROID',
  'app-text':'La Prova di Umanità usa un sistema biometrico fisico a 3 fattori. Fase 1: sensore ottico R503 scansiona tutte e 10 le impronte + MAX30102 PPG conferma il polso vitale. Fase 2: telecamera IR delle vene della mano (unicità 1 su 10⁷). Fase 3: scansione dell\'iride — standard d\'oro, 1 su 10⁷⁸, completamente indipendente dal dispositivo. Una prova ZK Groth16 impegna tutti i fattori senza rivelare alcun dato biometrico. I tuoi 1.000 AEQ vengono accreditati automaticamente al momento della verifica.',
  's1t':'Scansione Biometrica','s1d':'AequitasBio scansiona tutte e 10 le impronte (sensore ottico R503) + polso MAX30102 PPG conferma la vivezza. Fase 2: IR delle vene della mano. Fase 3: iride. I dati grezzi non lasciano mai il dispositivo.',
  's2t':'Generazione Prova ZK','s2d':'La prova ZK Groth16 impegna tutti i fattori biometrici: commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier legato al corpo, non al telefono — perdere il dispositivo non può creare una seconda identità.',
  's3t':'Connetti Wallet','s3d':'L\'app apre MetaMask su questa pagina · connetti il tuo wallet Ethereum · la prova è crittograficamente legata al tuo indirizzo',
  's4t':'1.000 AEQ Accreditati','s4d':'Registrazione confermata su Aequitas BlockDAG entro 6 secondi · 1.000 AEQ accreditati istantaneamente · la tua identità è registrata permanentemente come umano verificato',
  'priv-bar':'🔒 R503 Tutte le 10 Impronte · MAX30102 Vivezza · Fase 2: IR Vene Mano · Fase 3: Iride (10⁷⁸) · Groth16 ZKP · Dati non lasciano mai il dispositivo · Un umano · Per sempre',
  'conn-wallet':'WALLET CONNESSO','proof-recv':'⚡ PROVA ZK RICEVUTA','proof-hint':'Connetti wallet per registrarti',
  'btn-conn':'🦊 CONNETTI METAMASK','btn-reg':'🔐 REGISTRA ON-CHAIN',
  'btn-web-reg':'🌐 REGISTRA VIA BROWSER (WebAuthn)',
  'web-reg-warn':'⚠ Legato al dispositivo: Questa identità è legata a questo dispositivo e browser. Non è trasferibile su un altro dispositivo. Per un\'identità permanente multi-dispositivo, usa l\'App Android Aequitas.',
  'reg-log-hint':'// Apri l\'App Android Aequitas per generare la tua prova, poi torna qui...',
  'reg-details':'Dettagli Registrazione','k-network':'Rete','k-chainid':'ID Catena','k-grant':'Sussidio UBI',
  'k-fee':'Commissione Gas','free':'GRATUITO — completamente senza gas','k-limit':'Registrazioni','k-limit-v':'Una volta · permanente · immutabile',
  'k-bio':'Dati Biometrici','never-stored':'Mai memorizzati — rimangono sul tuo dispositivo',
  'k-proof':'Sistema di Prova','k-conf':'Conferma','k-conf-v':'Entro 6 secondi (1 blocco)',
  'k-sybil':'Protezione Sybil','k-sybil-v':'Una identità per biometrica · blocco permanente',
  'live-stats':'Statistiche Chain in Tempo Reale',
  's-height':'Altezza Blocco','s-height-sub':'Nuovo blocco ogni ~6s · BlockDAG · Produzione parallela',
  's-humans':'Umani Verificati','s-humans-sub':'ZKP biometrico · Una persona, un wallet, per sempre',
  's-supply':'Offerta Totale','s-supply-sub':'Sempre = Umani × 1.000 AEQ',
  's-index':'Indice Aequitas','s-index-sub':'0 = perfetta uguaglianza · 100 = massima disuguaglianza',
  's-uptime':'Uptime','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Prova di Umanità','ib-poh-t':'Ogni detentore di AEQ deve dimostrare crittograficamente di essere un essere umano unico e vivente. Nessun bot, nessuna azienda, nessuna IA. I dati biometrici non lasciano mai il tuo dispositivo.',
  'ib-fair':'Distribuzione Radicalmente Equa','ib-fair-t':'Ogni umano verificato riceve esattamente 1.000 AEQ alla registrazione. Nessun pre-mining, nessuna allocazione ai fondatori. L\'offerta totale è sempre uguale a umani verificati × 1.000.',
  'ib-dag':'Architettura BlockDAG','ib-dag-t':'Più blocchi possono essere prodotti simultaneamente e uniti. Throughput più alto, latenza più bassa rispetto alle blockchain lineari tradizionali.',
  'ib-gas':'Veramente Senza Gas','ib-gas-t':'La registrazione e i trasferimenti AEQ non costano assolutamente nulla. Non servono ETH, BNB o MATIC. Nessun conto bancario, nessuna carta di credito.',
  'recent-blocks':'Blocchi Recenti','blocks-desc':'MERGE = più genitori uniti (BlockDAG). TX = transazione di registrazione. Tempo blocco: ~6 secondi.',
  'loading':'Caricamento blocchi...','net-info':'Info Rete','k-chain':'Nome Catena','k-symbol':'Simbolo','k-btime':'Tempo Blocco',
  'k-cons':'Consenso','k-nodes':'Node Attivi','k-storage':'Archiviazione','add-mm':'🦊 AGGIUNGI A METAMASK','k-dec':'Decimali',
  'btn-add-mm':'+ AGGIUNGI RETE AEQUITAS',
  'phil':'"Il denaro esiste perché le persone esistono.<br>Niente di più, niente di meno."','phil-sub':'— IL PRINCIPIO AEQUITAS —',
  'humans-title':'Umani Verificati su Aequitas Chain',
  'h-what':'Cos\'è un Umano Verificato?','h-what-t':'Un Umano Verificato è un indirizzo wallet dimostrato crittograficamente appartenere a un essere umano unico e vivente. La verifica usa un sistema hardware a 3 fattori: R503 scansiona tutte e 10 le impronte; MAX30102 PPG conferma il polso vitale; Fase 2: IR vene della mano (1 su 10⁷); Fase 3: iride (1 su 10⁷⁸). Viene trasmessa solo una prova ZK Groth16. Nessun dato biometrico lascia mai il dispositivo.',
  'h-zkp':'Sistema di Prova a Conoscenza Zero','h-zkp-t':'Aequitas usa Groth16 su BN128 — stessa curva di Ethereum e Zcash. ~200 byte, ~10ms. commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier legato al corpo: perdere il telefono non crea una seconda identità. Nessun dato biometrico viene mai memorizzato.',
  'h-sybil':'Prevenzione Attacchi Sybil','h-sybil-t':'Fase 1: tutte e 10 le impronte + vivezza MAX30102 (polso PPG, rifiuta calchi/replay). Fase 2: IR vene della mano — caratteristica interna del corpo, impossibile da copiare, 1 su 10⁷, diversa nei gemelli identici. Fase 3: iride — 1 su 10⁷⁸, standard d\'oro globale. Nullifier = keccak256(iris‖vein‖domain). Un umano, un wallet, per sempre.',
  'h-global':'Inclusione Finanziaria Globale','h-global-t':'Nessun conto bancario, nessuna carta di credito, nessuna criptovaluta precedente necessaria. Solo uno smartphone Android con sensore biometrico. Aequitas è progettato per essere accessibile a ogni essere umano sulla Terra.',
  'h-bio-hw':'Roadmap Hardware Biometrico','h-bio-hw-t':'Fase 1 (attiva): scanner di impronte ottiche R503 — hash combinato di tutte e 10 le dita. Vivezza MAX30102 PPG. Fase 2 (pianificata): ESP32-CAM + LED IR 850 nm — imaging vene della mano, unicità 1 su 10⁷. Fase 3 (pianificata): modulo iride IR — 240+ gradi di libertà, 1 su 10⁷⁸, completamente indipendente dal dispositivo, i gemelli identici differiscono.',
  'reg-humans':'Umani Registrati','h-desc':'Ogni indirizzo è stato verificato come umano unico tramite ZKP biometrico. Ognuno ha ricevuto esattamente 1.000 AEQ. Il registro è permanente, immutabile e on-chain.',
  'no-humans':'Nessun umano registrato ancora.\n\nScarica l\'App Android Aequitas e sii il primo umano sulla chain!',
  'reg-stats':'Statistiche Registro','total-humans':'Totale Umani',
  'idx-title':'Indice Aequitas — Punteggio di Uguaglianza Economica in Tempo Reale',
  'idx-desc':'L\'Indice Aequitas misura la disuguaglianza economica tra tutti gli umani verificati in tempo reale. È derivato dal coefficiente Gini della distribuzione dei saldi on-chain. 0 = perfetta uguaglianza. 100 = massima disuguaglianza. Il protocollo attiva automaticamente i meccanismi di redistribuzione quando l\'indice sale.',
  'curr-idx':'Indice Attuale','bar-0':'0 — Perfetta Uguaglianza','bar-100':'100 — Massima Disuguaglianza',
  'gini':'Coefficiente Gini','gini-desc':'0 = uguale · 1 = disuguale',
  'supply-desc':'Sempre = Umani × 1.000 AEQ',
  'phase':'Fase Protocollo','phase-desc':'Avanza automaticamente per numero di umani',
  'humans-desc':'Umani unici verificati biometricamente',
  'pools-title':'Pool di Redistribuzione',
  'pools-desc':'Ogni commissione di swap, addebito di demurrage e overflow del limite di ricchezza viene automaticamente suddiviso tra quattro pool. Nessun intervento manuale — il protocollo gestisce tutta la redistribuzione solo attraverso il codice. Tutti i pool pagano quotidianamente.',
  'vel-pool':'Pool Validatori','vel-pool-desc':'40% di tutte le commissioni → operatori node che proteggono la rete',
  'liq-pool':'Pool Liquidità','liq-pool-desc':'30% di tutte le commissioni → fornitori di liquidità, proporzionale alle quote LP',
  'ubi-pool':'Pool UBI','ubi-pool-desc':'20% di tutte le commissioni → tutti gli umani verificati equamente, ogni 24 ore',
  'treasury':'Tesoreria','treasury-desc':'10% di tutte le commissioni → sviluppo e manutenzione del protocollo',
  'phases-title':'Fasi del Protocollo',
  'demurrage-title':'Demurrage — Incentivo a Circolare',
  'demurrage-desc':'Aequitas implementa un meccanismo di demurrage ispirato alle valute complementari storiche. I saldi AEQ inattivi perdono lentamente valore per scoraggiare l\'accumulo e incentivare la partecipazione economica.',
  'dem-rate-k':'Tasso di Decadimento','dem-rate-v':'0,5% al mese (continuo, non a gradini)',
  'dem-grace-k':'Periodo di Grazia','dem-grace-v':'3 mesi di inattività prima che inizi il decadimento',
  'dem-reset-k':'Reset Timer','dem-reset-v':'Qualsiasi trasferimento, swap o azione di liquidità azzera il timer',
  'dem-dest-k':'AEQ decaduto va a','dem-dest-v':'Pool di redistribuzione (suddivisione 40/30/20/10)',
  'dem-warn-k':'Sistema di Avviso','dem-warn-v':'Avviso di 14 giorni (una volta) + promemoria di 7 giorni ripetuto ad ogni accesso',
  'story-title':'La Storia di Aequitas — Perché Esiste',
  'story-text':'<p>L\'anno è 2009. Satoshi Nakamoto rilascia Bitcoin. Per la prima volta, il valore può trasferirsi tra due persone senza una banca. Una vera rivoluzione. Ma quasi immediatamente qualcosa va storto.</p><p>I primi miner accumulano milioni di monete a costo quasi zero. Entro il 2021, l\'1% superiore degli indirizzi Bitcoin controlla oltre il 90% di tutti i Bitcoin. Il coefficiente Gini stimato di Bitcoin supera 0,85 — più alto di qualsiasi paese sulla Terra. La criptovaluta che avrebbe dovuto democratizzare la finanza ha creato la più estrema concentrazione di ricchezza nella storia umana.</p><p><span style="color:var(--gold)">Aequitas</span> — Latino per "equità" e "uguaglianza" — è stato creato per rispondere a una singola domanda: <em style="color:var(--gold)">"Come sarebbe una criptovaluta progettata dai principi fondamentali per essere equa per ogni essere umano?"</em></p><p>La risposta è semplice: <strong style="color:var(--text)">Il denaro esiste perché le persone esistono. Quindi ogni persona dovrebbe avere una quota uguale di denaro semplicemente in virtù di essere umana.</strong></p><p>Aequitas implementa questo matematicamente. Ogni umano verificato riceve 1.000 AEQ. Nessun mining, nessuno staking, nessun vantaggio per i primi adottanti. Il protocollo si adatta automaticamente man mano che la rete cresce.</p><p><em style="color:var(--gold)">"Il denaro esiste perché le persone esistono. Niente di più, niente di meno."</em></p>',
  'nodes-title':'Node Attivi — Topologia Attuale della Rete',
  'nodes-desc':'La rete Aequitas opera attualmente su due node distribuiti geograficamente. Entrambi partecipano alla produzione di blocchi, sincronizzazione dello stato e servizio API. Comunicano peer-to-peer via libp2p e sincronizzano lo stato dei blocchi via HTTP. La rete è progettata per supportare node aggiuntivi.',
  'node1':'Node 1 — Railway (Primario)','node1-desc':'API primario · Produttore blocchi · Distribuzione UBI · Bootstrap P2P · PostgreSQL · RPC per MetaMask',
  'node2':'Node 2 — Render (Secondario)','node2-desc':'API secondario · Produttore blocchi · Peer P2P · Sincronizzazione HTTP · Stato PostgreSQL condiviso',
  'run-node-title':'Esegui il Tuo Node — Aiuta a Proteggere la Rete',
  'run-node-desc':'Chiunque può eseguire un node Aequitas — senza permesso, senza stake, senza candidatura richiesta. I node partecipano alla produzione di blocchi e validano il registro umano. Gli operatori di node guadagnano una quota delle commissioni del protocollo tramite il Pool Validatori (40% di tutte le commissioni di swap, distribuite quotidianamente).',
  'bootstrap-title':'Connettere un Nuovo Node','bootstrap-desc':'Per eseguire il tuo node, imposta PRIMARY_NODE_URL=https://aequitas.digital nel tuo ambiente. Il tuo node si sincronizzerà automaticamente con lo stato completo della chain.',
  'tech-title':'Specifiche Tecniche','mm-config':'Configurazione MetaMask',
  'k-lang':'Lingua','k-src':'Codice Sorgente','evm-yes':'Sì — JSON-RPC /rpc · Compatibile MetaMask',
  'proto-label':'Protocollo Aequitas V7 — Documentazione Tecnica',
  'ca-title':'Indirizzi Contratto','ca-text':'Chain: Aequitas Chain (ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Principale): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7 è l\'unica fonte di verità per l\'intera economia Aequitas. Ogni saldo AEQ, ogni registrazione umana, ogni pagamento UBI e ogni applicazione del limite di ricchezza è governato da questo unico contratto immutabile — distribuito su Aequitas Chain, una blockchain personalizzata compatibile con EVM che esegue un motore di consenso BlockDAG. Non c\'è chiave amministratore, nessun proxy di aggiornamento, nessun voto di governance che possa cambiare una singola riga della sua logica. Il codice che funziona oggi è il codice che funzionerà tra dieci anni.<br><br>Il contratto BioVerifier riceve prove a conoscenza zero Groth16 generate interamente sul dispositivo Android dell\'utente. Verifica matematicamente on-chain in ~10 ms che un nuovo registrante è un essere umano unico e vivo — senza mai conoscere il suo nome, identità o dati biometrici. Questo è ciò che rende possibile la registrazione senza gas e senza investimenti: la prova è l\'unica cosa che lascia mai il dispositivo.<br><br>Insieme, questi due contratti rendono possibile qualcosa che non è mai esistito in nessun sistema monetario nella storia: un\'offerta monetaria le cui regole — chi la ottiene, quanta ne esiste, come si ridistribuisce — non può essere alterata da nessuna persona, azienda o governo. Mai.',
  'poa-title':'1. PROVA DI VITA — Recupero Saldi Inattivi','poa-text':'<p>Cosa succede all\'AEQ quando le persone muoiono o diventano permanentemente incapaci? In Bitcoin, i portafogli persi significano fornitura persa permanentemente. Aequitas risolve questo con un sistema di recupero dell\'inattività a più fasi: se un portafoglio non mostra attività per un periodo prolungato, il suo saldo viene gradualmente restituito alla comunità attraverso il pool UBI.</p>',
  'poa-box':'Anno 0–2: Uso normale — nessuna restrizione<br>Anno 2: Avviso 1 — il Guardian può rispondere a nome<br>Anno 2+60g: Avviso 2 — urgenza crescente<br>Anno 2+120g: Avviso 3 — avviso finale<br>Anno 2+180g: AEQ spostato in ESCROW personale (ancora recuperabile)<br>Anno 4: Se ancora inattivo — ESCROW rilasciato al Pool UBI',
  'guard-title':'2. SISTEMA GUARDIAN — Protezione Umana','guard-text':'<p>E se qualcuno è ricoverato in ospedale o non riesce ad accedere al proprio dispositivo per mesi? Il sistema Guardian permette a una persona di fiducia — un altro umano verificato — di confermare che il proprietario del portafoglio è ancora vivo. Il Guardian ha accesso finanziario strettamente nullo: può solo chiamare una singola funzione che reimposta il timer di inattività. Non può spostare, spendere o accedere ai fondi in nessuna circostanza.</p>',
  'guard-box':'1 Guardian per umano · deve essere un umano verificato su Aequitas<br>Il Guardian può SOLO chiamare confirmAlive() — zero diritti di transazione<br>Il Guardian NON PUÒ spostare fondi, trasferire AEQ o accedere al portafoglio<br>Massimo 3 tutelati per Guardian · Blocco di 7 giorni all\'assegnazione · Nessuna relazione circolare',
  'dem-title':'3. DEMURRAGE — Meccanismo Anti-Accumulo',
  'dem-box':'Tasso: 0,5%/mese dopo 3 mesi di inattività (continuo, non a gradini)<br>Il timer si azzera automaticamente con qualsiasi trasferimento, swap o azione di liquidità<br>AEQ decaduto ridistribuito ai quattro pool — mai bruciato<br>Avviso di 14 giorni mostrato una volta · 7 giorni ripetuto in ogni sessione attiva',
  'dem-text':'<p>Il demurrage è un costo di detenzione sul denaro — un tasso di interesse negativo che rende costoso accumulare e attraente la circolazione. L\'esperimento di Wörgl (Austria, 1932) usò una valuta con demurrage e ridusse la disoccupazione locale del 25% in un anno. La Banca Centrale austriaca lo chiuse proprio perché funzionava troppo bene. Il Chiemgauer (Germania, 2003) opera con lo stesso principio con successo da oltre 20 anni.</p>',
  'cap-title':'4. LIMITE DI RICCHEZZA — Applicazione dell\'Equità Matematica','cap-box':'Bootstrap: max(5,min(N,25))× saldo AEQ medio<br>1–4 umani: 5× (5.000 AEQ) · Cresce 1× per umano · 25+: 25× (25.000 AEQ) permanente<br>Si applica a TUTTI gli indirizzi tranne i 4 pool del protocollo<br>L\'eccesso di AEQ viene immediatamente ridistribuito · Nessun intervento manuale',
  'ubi-title':'5. REDDITO UNIVERSALE DI BASE — Ridistribuzione Giornaliera','ubi-box':'Fonti di reddito del Pool UBI:<br>· 20% di tutte le commissioni di swap del pool AMM AEQ↔tUSD<br>· Overflow dall\'applicazione del limite di ricchezza<br>· Addebiti di demurrage da account inattivi<br>· Escrow inattivo rilasciato dopo 4 anni<br><br>Distribuzione: Ogni 24 ore, l\'intero saldo del pool UBI viene diviso equamente tra tutti gli umani verificati registrati. Il pool si azzera e inizia immediatamente a riempirsi di nuovo dall\'attività continua del protocollo.',
  'inf-title':'6. NESSUNA INFLAZIONE ALGORITMICA — Formula di Fornitura Fissa','inf-box':'L\'UNICO evento che crea nuovo AEQ: un nuovo umano verificato si registra.<br><br>Offerta Totale = Umani Verificati × 1.000 AEQ<br><br>Questo non è una politica — è applicato dal protocollo. Nessun amministratore può coniare AEQ aggiuntivo, nessun voto di governance può modificare l\'emissione. AEQ è l\'unica criptovaluta in cui l\'offerta totale è determinata esclusivamente dal numero di esseri umani vivi verificati.',
  'phases-desc':'In Fase 0 (Bootstrap) il limite di ricchezza usa un moltiplicatore scorrevole: max(5, min(N, 25))× saldo medio. Con 1–4 umani: 5× media. Ogni nuovo umano aggiunge 1×. A 25+ umani: bloccato permanentemente a 25×. Fase 1+ mantiene 25× fisso. Tutte le transizioni sono automatiche — nessun voto, nessuna chiave admin.',
  'p0':'Bootstrap · &lt;100 umani · Limite di Ricchezza: max(5,min(N,25))× media · Scorre 5×→25× fino al 25° umano · Attualmente attivo',
  'p1':'Crescita · 100–10.000 umani · Limite di Ricchezza: 25× saldo medio',
  'p2':'Stabilità · 10.000–1M umani · Limite di Ricchezza: 25× saldo medio',
  'p3':'Maturità · 1M+ umani · Limite di Ricchezza: 25× saldo medio',
  'wealth-cap-explain':'Il Limite di Ricchezza in Fase 0 (Bootstrap) usa max(5, min(N, 25))× saldo AEQ medio, dove N = umani registrati. 1–4 umani: 5× media. Ogni nuovo umano aggiunge 1×. 25+ umani: bloccato permanentemente a 25×. Il limite si adatta sempre al saldo medio corrente.',
  'btn-download-app':'SCARICA L\'APP AEQUITASBIO',
  'swap-title':'🔄 Scambia AEQ ↔ tUSD','swap-sub':'Scambia AEQ con tUSD (un dollaro di test simulato) attraverso il pool di liquidità nativo. Una commissione dello 0,1% si applica solo agli scambi — i normali trasferimenti AEQ tra persone rimangono completamente gratuiti.',
  'swap-priv-bar':'🔒 Solo 0,1% commissione swap · Trasferimenti AEQ-AEQ gratuiti · tUSD è una valuta di test senza valore reale',
  'swap-your-aeq':'Il tuo AEQ','swap-your-tusd':'Il tuo tUSD','swap-aeq-to-tusd':'AEQ → tUSD','swap-tusd-to-aeq':'tUSD → AEQ',
  'swap-fee-est':'Commissione protocollo (0,1%)','swap-details-hdr':'Dettagli Scambio',
  'swap-out-lbl':'Ricevi (est.)','swap-impact-lbl':'Impatto sul prezzo','swap-rate-lbl':'Tasso di cambio',
  'swap-depth-lbl':'Composizione del Pool','amm-title':'x × y = k — AMM a Prodotto Costante',
  'amm-text':'Quando scambi AEQ con tUSD, la riserva AEQ cresce e quella tUSD diminuisce — il loro prodotto rimane sempre uguale a k. Scambi più grandi causano un maggiore impatto sul prezzo. La commissione dello 0,1% viene detratta prima di applicare la formula.',
  'swap-btn-conn':'🦊 COLLEGA METAMASK','swap-btn-go':'🔄 SCAMBIA',
  'swap-log-hint':'// Collega il wallet per scambiare...',
  'swap-no-liquidity':'Nessun tUSD ancora?','swap-faucet-desc':'Gli umani registrati possono richiedere tUSD di test una volta','swap-btn-faucet':'💧 RICHIEDI tUSD DI TEST',
  'swap-addliq-title':'Fornire Liquidità','swap-addliq-desc':'Sii il primo a depositare — il tuo rapporto imposta il prezzo iniziale.','swap-btn-addliq':'💧 AGGIUNGI LIQUIDITÀ',
  'swap-lp-title':'La tua Posizione LP','swap-lp-share':'Quota del Pool','swap-lp-withdrawable':'Prelevabile',
  'swap-lp-pct-label':'% della tua posizione','swap-lp-youget':'Riceverai','swap-btn-removeliq':'🔥 RIMUOVI LIQUIDITÀ',
  'swap-pool-title':'AEQ / tUSD — Stato del Pool',
  'swap-pool-aeq':'Riserva AEQ','swap-pool-tusd':'Riserva tUSD','swap-pool-price':'Prezzo Spot',
  'swap-fee-bps':'Commissione Swap','swap-fee-split':'Distribuzione commissioni','swap-fee-split-v':'40% Validatori / 30% LP / 20% UBI / 10% Tesoreria',
  'swap-pools-addr-title':'Indirizzi Pool Tokenomics',
  'swap-validators':'Validatori (40%)','swap-lps':'Fornitori di Liquidità (30%)','swap-ubi':'Pool UBI (20%)','swap-treasury':'Tesoreria (10%)',
  'ubi-hero-title':'REDDITO UNIVERSALE DI BASE — POOL UBI',
  'ubi-hero-sub':'Accumulando — prossimo pagamento distribuito equamente a tutti gli umani verificati in:',
  'ubi-bal-lbl':'saldo attuale del pool','ubi-hero-desc':'Diviso equamente tra tutti · pagato ogni 24h · il pool si azzera dopo ogni pagamento · nessun saldo minimo richiesto',
  'ubi-how-fills':'Come si riempie il Pool UBI',
  'ubi-src-swap':'Commissioni Swap','ubi-src-swap-d':'Ogni swap AEQ↔tUSD contribuisce il 20% della sua commissione dello 0,1%. Più trading = riempimento più rapido.',
  'ubi-src-dem':'Demurrage','ubi-src-dem-d':'AEQ inattivo (3+ mesi) decade dello 0,5%/mese. Il 20% dell\'importo decaduto va all\'UBI.',
  'ubi-src-cap':'Overflow Limite di Ricchezza','ubi-src-cap-d':'I wallet che superano max(5,min(N,25))× il saldo medio hanno l\'eccesso confiscato istantaneamente. Il 20% fluisce all\'UBI.',
  'pools4-header':'Tutti e quattro i pool di redistribuzione',
  'ubi-see-above':'vedi conto alla rovescia sopra','ubi-timer-above':'⏰ conto alla rovescia mostrato sopra','pool-t-timer':'Accumula — nessun timer',
  'usp-headline':'Per la prima volta nella storia — tutti iniziano alla pari',
  'usp-sub':'Se possiedi uno smartphone Android, sei idoneo. Senza banca, senza conoscenze crypto, senza investimento.',
  'usp-c1-title':'0,00 Investimento Iniziale','usp-c1-desc':'La registrazione è completamente senza gas. Senza ETH, senza MATIC, senza carta di credito. Il protocollo paga tutte le commissioni per te.',
  'usp-c2-title':'1.000 AEQ per ogni umano','usp-c2-desc':'Miliardario o agricoltore di sussistenza — tutti ricevono esattamente 1.000 AEQ. Non di più, non di meno. Inizio uguale, garantito dalla matematica.',
  'usp-c3-title':'Accessibile a tutti','usp-c3-desc':'Nessun conto bancario, carta di credito o documento d\'identità. La registrazione utilizza un kit biometrico economico (scanner di impronte + sensore di polso, ~15 €) — pensato per l\'accesso globale.',
  'usp-c4-title':'UBI quotidiano per sempre','usp-c4-desc':'Una volta registrato, ricevi automaticamente una quota giornaliera dei pagamenti UBI — ogni giorno, senza alcuna azione richiesta.',
  'v7-intro-title':'Cos\'è AequitasV7?',
  'v7-intro-text':'AequitasV7 è il contratto intelligente centrale del protocollo Aequitas. "V7" si riferisce alla 7ª versione principale del contratto di equità. È distribuito immutabilmente su Aequitas Chain (ID 1926) e gestisce ogni aspetto: registrazione umana, verifica ZK, gestione saldi, limite di ricchezza, distribuzione UBI, commissioni swap. Nessun amministratore può aggiornarlo. I sei meccanismi formano un sistema auto-rinforzante.',
  'explore-title':'Esplora Aequitas',
  'expl-score':'Punteggio Uguaglianza','expl-score-d':'Coefficiente Gini live · Indice Aequitas · distribuzione ricchezza in tempo reale',
  'expl-economy':'UBI e Pool di Redistribuzione','expl-economy-d':'Conto alla rovescia UBI giornaliero · 4 pool on-chain · demurrage · Fasi del Protocollo',
  'expl-charts':'Grafici e Storia','expl-charts-d':'Storia Gini · curva di Lorenz · slider bootstrap limite ricchezza · La storia di Aequitas',
  'expl-v7':'Documentazione Protocollo V7','expl-v7-d':'Contratto AequitasV7 · 6 meccanismi · prova ZK · limite ricchezza · demurrage · codice immutabile',
  'expl-explorer':'Block Explorer','expl-explorer-d':'BlockDAG live · clicca qualsiasi blocco per vedere validatore, hash, transazioni, hash genitori',
  'swap-sell-label':'Vendi','swap-receive-label':'Ricevi',
  'gini-calc-title':'Come si calcola l indice','gini-calc-text':'Vengono raccolti tutti i saldi AEQ. La formula calcola la differenza assoluta media normalizzata per n2. Risultato 0-1 x 100 = Indice Aequitas.','gini-why-title':'Perche Gini','gini-why-text':'Il coefficiente Gini cattura la distribuzione completa in un numero verificabile.','expl-network':'Rete e Nodi','expl-network-d':'Topologia nodi · esegui il tuo nodo · specifiche tecniche · Chain ID 1926'
},
tr:{
  'logo-sub':'İNSANLIK KANITI','live':'CANLI',
  'tab-register':'🔐 Kayıt','tab-explorer':'🔍 Gezgin','tab-humans':'👥 İnsanlar','tab-index':'📊 Endeks','tab-network':'🌐 Ağ','tab-protocol':'📜 Protokol V7','tab-swap':'🔄 Takas',
  'reg-title':'🔐 Doğrulanmış İnsan Olarak Kayıt Ol',
  'reg-sub':'Aequitas ağına katıl ve 1.000 AEQ Evrensel Temel Gelir hibeni al. Tek seferlik, kalıcı ve tamamen ücretsiz. Hiçbir kişisel veri asla saklanmaz.',
  'app-title':'KAYIT YALNIZCA ANDROİD UYGULAMASI İLE',
  'app-text':'İnsanlık Kanıtı fiziksel 3 faktörlü biyometrik sistem kullanır. Faz 1: R503 optik sensör tüm 10 parmak izini tarar + MAX30102 PPG canlı nabzı doğrular. Faz 2: el damarı IR kamerası (1/10⁷ benzersizlik). Faz 3: iris taraması — altın standart, 1/10⁷⁸, tamamen cihazdan bağımsız. Groth16 ZK kanıtı tüm faktörleri hiçbir biyometrik veri ifşa etmeden taahhüt eder. 1.000 AEQ\'n doğrulama sonrası otomatik olarak yatırılır.',
  's1t':'Biyometrik Tarama','s1d':'AequitasBio tüm 10 parmak izini tarar (R503 optik sensör) + MAX30102 PPG nabzı canlılığı doğrular. Faz 2: el damarı IR. Faz 3: iris. Ham veriler asla cihazı terk etmez.',
  's2t':'ZK Kanıtı Oluşturma','s2d':'Groth16 ZK kanıtı tüm biyometrik faktörleri taahhüt eder: commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier telefona değil bedene bağlıdır — cihaz kaybı ikinci kimlik oluşturamaz.',
  's3t':'Cüzdan Bağla','s3d':'Uygulama bu sayfada MetaMask\'ı açar · Ethereum cüzdanını bağla · kanıt kriptografik olarak adresine bağlanır',
  's4t':'1.000 AEQ Yatırıldı','s4d':'Kayıt 6 saniye içinde Aequitas BlockDAG\'da onaylandı · 1.000 AEQ anında yatırıldı · kimliğin kalıcı olarak doğrulanmış insan olarak kaydedildi',
  'priv-bar':'🔒 R503 Tüm 10 Parmak İzi · MAX30102 Canlılık · Faz 2: El Damarı IR · Faz 3: İris (10⁷⁸) · Groth16 ZKP · Veriler asla cihazı terk etmez · Bir insan · Sonsuza dek',
  'conn-wallet':'BAĞLI CÜZDAN','proof-recv':'⚡ ZK KANITI ALINDI','proof-hint':'Kayıt için cüzdan bağla',
  'btn-conn':'🦊 METAMASK BAĞLA','btn-reg':'🔐 ZİNCİRE KAYIT OL',
  'btn-web-reg':'🌐 TARAYICI ÜZERİNDEN KAYIT (WebAuthn)',
  'web-reg-warn':'⚠ Cihaza bağlı: Bu kimlik bu cihaza ve tarayıcıya bağlıdır. Başka bir cihaza aktarılamaz. Kalıcı çok cihazlı kimlik için Aequitas Android Uygulamasını kullan.',
  'reg-log-hint':'// Kanıtını oluşturmak için Aequitas Android Uygulamasını aç, ardından buraya dön...',
  'reg-details':'Kayıt Detayları','k-network':'Ağ','k-chainid':'Zincir ID','k-grant':'UBI Hibesi',
  'k-fee':'Gas Ücreti','free':'ÜCRETSİZ — tamamen gas\'sız','k-limit':'Kayıtlar','k-limit-v':'İnsan başına bir kez · kalıcı · değiştirilemez',
  'k-bio':'Biyometrik Veri','never-stored':'Asla saklanmaz — cihazında kalır',
  'k-proof':'Kanıt Sistemi','k-conf':'Onay','k-conf-v':'6 saniye içinde (1 blok)',
  'k-sybil':'Sybil Koruması','k-sybil-v':'Biyometri başına bir kimlik · kalıcı kilit',
  'live-stats':'Canlı Zincir İstatistikleri',
  's-height':'Blok Yüksekliği','s-height-sub':'Her ~6 saniyede yeni blok · BlockDAG · Paralel üretim',
  's-humans':'Doğrulanmış İnsanlar','s-humans-sub':'Biyometrik ZKP · Bir kişi, bir cüzdan, sonsuza dek',
  's-supply':'Toplam Arz','s-supply-sub':'Her zaman = İnsanlar × 1.000 AEQ',
  's-index':'Aequitas Endeksi','s-index-sub':'0 = mükemmel eşitlik · 100 = maksimum eşitsizlik',
  's-uptime':'Çalışma Süresi','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'İnsanlık Kanıtı','ib-poh-t':'Her AEQ sahibi, benzersiz bir yaşayan insan olduğunu kriptografik olarak kanıtlamak zorundadır. Robot yok, şirket yok, yapay zeka yok. Biyometrik veriler asla cihazı terk etmez.',
  'ib-fair':'Radikal Şekilde Adil Dağıtım','ib-fair-t':'Her doğrulanmış insan kayıt sırasında tam olarak 1.000 AEQ alır. Ön madencilik yok, kurucu tahsisi yok. Toplam arz her zaman doğrulanmış insanlar × 1.000 eşittir.',
  'ib-dag':'BlockDAG Mimarisi','ib-dag-t':'Birden fazla blok eş zamanlı olarak üretilebilir ve birleştirilebilir. Doğrusal blok zincirlerine kıyasla daha yüksek verim, daha düşük gecikme.',
  'ib-gas':'Gerçekten Gas\'sız','ib-gas-t':'Kayıt ve AEQ transferleri kesinlikle ücretsizdir. ETH, BNB veya MATIC gerekmez. Banka hesabı veya kredi kartı gerekmez.',
  'recent-blocks':'Son Bloklar','blocks-desc':'MERGE = birden fazla ebeveyn birleştirildi (BlockDAG). TX = kayıt işlemi. Blok süresi: ~6 saniye. Bloka tıklayarak detayları, doğrulayıcıyı ve işlemleri görüntüle.',
  'loading':'Bloklar yükleniyor...','net-info':'Ağ Bilgisi','k-chain':'Zincir Adı','k-symbol':'Sembol','k-btime':'Blok Süresi',
  'k-cons':'Konsensüs','k-nodes':'Aktif Node\'lar','k-storage':'Depolama','add-mm':'🦊 METAMASK\'A EKLE','k-dec':'Ondalık',
  'btn-add-mm':'+ AEQUITAS AĞINI EKLE',
  'phil':'"Para insanlar var olduğu için var.<br>Bundan fazlası değil, bundan azı değil."','phil-sub':'— AEQUİTAS İLKESİ —',
  'humans-title':'Aequitas Zincirindeki Doğrulanmış İnsanlar',
  'h-what':'Doğrulanmış İnsan Nedir?','h-what-t':'Doğrulanmış İnsan, benzersiz bir yaşayan insana ait olduğu kriptografik olarak kanıtlanmış bir cüzdan adresidir. Doğrulama 3 faktörlü donanım sistemi kullanır: R503 tüm 10 parmak izini tarar; MAX30102 PPG canlı nabzı doğrular; Faz 2: el damarı IR (1/10⁷); Faz 3: iris (1/10⁷⁸). Yalnızca Groth16 ZK kanıtı iletilir. Hiçbir biyometrik veri cihazı terk etmez.',
  'h-zkp':'Sıfır Bilgi Kanıtı Sistemi','h-zkp-t':'Aequitas, BN128 üzerinde Groth16 kullanır — Ethereum ve Zcash ile aynı eğri. ~200 bayt, ~10ms. commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier bedene bağlıdır: telefonu kaybetmek ikinci kimlik oluşturmaz. Hiçbir biyometrik veri asla saklanmaz.',
  'h-sybil':'Sybil Saldırısı Önleme','h-sybil-t':'Faz 1: tüm 10 parmak izi + MAX30102 canlılık (PPG nabzı, kalıp/tekrar oynatmayı reddeder). Faz 2: el damarı IR — iç vücut özelliği, kopyalanması imkânsız, 1/10⁷, özdeş ikizlerde farklı. Faz 3: iris — 1/10⁷⁸, küresel altın standart. Nullifier = keccak256(iris‖vein‖domain). Bir insan, bir cüzdan, sonsuza dek.',
  'h-global':'Küresel Finansal Kapsayıcılık','h-global-t':'Banka hesabı, kredi kartı veya önceden kripto para gerekmez. Yalnızca biyometrik sensörlü bir Android akıllı telefon yeterlidir.',
  'h-bio-hw':'Biyometrik Donanım Yol Haritası','h-bio-hw-t':'Faz 1 (aktif): R503 optik parmak izi tarayıcısı — tüm 10 parmağın birleşik hash\'i. MAX30102 PPG canlılık. Faz 2 (planlandı): ESP32-CAM + 850 nm IR LED — el damarı görüntüleme, 1/10⁷ benzersizlik. Faz 3 (planlandı): IR iris modülü — 240+ serbestlik derecesi, 1/10⁷⁸, tamamen cihazdan bağımsız, özdeş ikizler farklı.',
  'reg-humans':'Kayıtlı İnsanlar','h-desc':'Aşağıdaki her adres, biyometrik ZKP aracılığıyla benzersiz insan olarak doğrulandı. Her biri tam olarak 1.000 AEQ aldı. Kalıcı, değiştirilemez, zincir üzerinde.',
  'no-humans':'Henüz kayıtlı insan yok.\n\nAequitas Android Uygulamasını indir ve zincirdeki ilk insan ol!',
  'reg-stats':'Kayıt İstatistikleri','total-humans':'Toplam İnsan',
  'idx-title':'Aequitas Endeksi — Gerçek Zamanlı Ekonomik Eşitlik Puanı',
  'idx-desc':'Aequitas Endeksi, tüm doğrulanmış insanların ekonomik eşitsizliğini gerçek zamanlı olarak ölçer. Zincir üzerindeki bakiye dağılımının <strong style="color:var(--teal)">Gini katsayısından</strong> türetilir. <strong style="color:var(--neon)">0 = mükemmel eşitlik</strong>. <strong style="color:var(--red)">100 = maksimum eşitsizlik</strong>. Bitcoin Gini ≈ 0,85 · Güney Afrika ≈ 0,63 · İskandinavya ≈ 0,27 · Aequitas hedefi: Gini 0,35\'in altında.',
  'gini-what-title':'Gini Katsayısı Nedir?',
  'gini-what-text':'İtalyan istatistikçi Corrado Gini tarafından 1912\'de geliştirilmiştir. Lorenz eğrisi ile görselleştirilen gerçek dağılımı mükemmel eşit dağılımla karşılaştırarak servet dağılımını ölçer. Ölçek: 0 (herkes aynı miktarı tutar) ile 1 (bir kişi her şeyi tutar). Dünya Bankası, OECD ve BM tarafından kullanılır.',
  'gini-calc-title':'Aequitas Endeksi Nasıl Hesaplanır?',
  'gini-calc-text':'Tüm doğrulanmış insanların AEQ bakiyeleri toplanır. Formül, tüm bakiye çiftleri arasındaki ortalama mutlak farkı, nüfus karesi (n²) ve ortalama bakiye (x̄) ile normalleştirilmiş olarak hesaplar. Sonuç 0–1 ile 100 ile çarpılır = Aequitas Endeksi.',
  'gini-why-title':'Neden Gini — Daha Basit Bir Metrik Değil?',
  'gini-why-text':'Basit bir zengin-fakir oranı kolayca manipüle edilebilir: 10.000 cüzdan düşük bir spread gösterebilir ama AEQ\'nun %90\'ı 100 elde konsantre olabilir — Gini bunu tespit eder, bir oran etmez. Katsayı, tüm doğrulanmış insanlar arasındaki tam dağılımı tek bir denetlenebilir sayıda yakalar.',
  'curr-idx':'Mevcut Endeks','bar-0':'0 — Mükemmel Eşitlik','bar-100':'100 — Maks. Eşitsizlik',
  'wcap-lbl':'Mevcut Servet Tavanı:','wcap-mult':'Çarpan:','wcap-avg':'Ort. bakiye:',
  'gini':'Gini Katsayısı','gini-desc':'0 = eşit · 1 = eşitsiz',
  'supply-desc':'Her zaman = İnsanlar × 1.000 AEQ',
  'phase':'Protokol Aşaması','phase-desc':'İnsan sayısına göre otomatik ilerler',
  'humans-desc':'Biyometrik olarak doğrulanmış benzersiz insanlar',
  'pools-title':'Yeniden Dağıtım Havuzları',
  'pools-desc':'Her takas ücreti, gecikme ücreti ve servet tavanı taşması otomatik olarak dört havuza bölünür. Manuel müdahale yok. Tüm havuzlar günlük ödeme yapar.',
  'vel-pool':'Doğrulayıcı Havuzu','vel-pool-desc':'Tüm ücretlerin %40\'ı → ağı güvence altına alan node operatörleri',
  'liq-pool':'Likidite Havuzu','liq-pool-desc':'Tüm ücretlerin %30\'u → LP paylarıyla orantılı likidite sağlayıcıları',
  'ubi-pool':'UBI Havuzu','ubi-pool-desc':'Tüm ücretlerin %20\'si → her 24 saatte tüm doğrulanmış insanlar eşit olarak',
  'treasury':'Hazine','treasury-desc':'Tüm ücretlerin %10\'u → protokol geliştirme ve bakımı',
  'phases-title':'Protokol Aşamaları',
  'phases-desc':'Aşama 0\'da servet tavanı bir bootstrap çarpanı kullanır: max(5, min(N, 25))× ortalama bakiye. 1–4 insanla: 5× ortalama. Her yeni insan 1× ekler. 25+ insanda: kalıcı olarak 25×\'e sabitlenir. Aşama 1+ 25×\'i sabit tutar. Tüm geçişler otomatiktir — yönetişim oyu yok, yönetici anahtarı yok.',
  'p0':'Bootstrap · &lt;100 insan · Servet Tavanı: max(5,min(N,25))× ort. · 5×→25× arası kayar · Şu anda aktif',
  'p1':'Büyüme · 100–10.000 insan · Servet Tavanı: 25× ortalama bakiye',
  'p2':'Kararlılık · 10.000–1M insan · Servet Tavanı: 25× ortalama bakiye',
  'p3':'Olgunluk · 1M+ insan · Servet Tavanı: 25× ortalama bakiye',
  'wealth-cap-explain':'Aşama 0\'daki (Bootstrap) Servet Tavanı max(5, min(N, 25))× ortalama AEQ bakiyesi kullanır; burada N = kayıtlı insan sayısı. 1–4 insan: 5× ortalama. Her yeni insan 1× ekler. 25+ insan: kalıcı olarak 25×. Tavan her zaman mevcut ortalama bakiyeyle ölçeklenir.',
  'demurrage-title':'Gecikme Ücreti — Dolaşım Teşviki',
  'demurrage-desc':'Aequitas, tarihi tamamlayıcı para birimlerinden ilham alan bir gecikme ücreti mekanizması uygular. Atıl AEQ bakiyeleri, biriktirmeyi caydırmak için yavaşça değer kaybeder.',
  'dem-rate-k':'Bozunma Hızı','dem-rate-v':'Ayda %0,5 (sürekli, kademeli değil)',
  'dem-grace-k':'İzin Süresi','dem-grace-v':'Bozunma başlamadan önce 3 aylık hareketsizlik',
  'dem-reset-k':'Saat Sıfırlama','dem-reset-v':'Herhangi bir transfer, takas veya likidite işlemi zamanlayıcıyı sıfırlar',
  'dem-dest-k':'Bozunan AEQ şuraya gider','dem-dest-v':'Yeniden dağıtım havuzları (40/30/20/10 bölünmesi)',
  'dem-warn-k':'Uyarı Sistemi','dem-warn-v':'14 günlük bildirim (bir kez) + her girişte 7 günlük tekrarlayan hatırlatma',
  'story-title':'Aequitas\'ın Hikayesi — Neden Var Olduğu',
  'story-text':'<p>Yıl 2009. Satoshi Nakamoto Bitcoin\'i yayınlıyor. İlk kez, değer bir banka olmadan iki kişi arasında transfer edilebiliyor. Gerçek bir devrim. Ama neredeyse hemen bir şeyler ters gidiyor.</p><p>Erken madenciler neredeyse sıfır maliyetle milyonlarca coin biriktiriyor. 2021\'e kadar Bitcoin adreslerinin en üst %1\'i tüm Bitcoin\'in %90\'ından fazlasını kontrol ediyor. Bitcoin\'in tahmini Gini katsayısı 0,85\'i aşıyor — Dünya\'daki herhangi bir ülkeden daha yüksek.</p><p><span style="color:var(--gold)">Aequitas</span> — Latince "adalet" ve "eşitlik" anlamına gelir — tek bir soruyu yanıtlamak için yaratıldı: <em style="color:var(--gold)">"Her insana adil olacak şekilde ilk ilkelerden tasarlanmış bir kripto para nasıl görünürdü?"</em></p><p>Cevap basit: <strong style="color:var(--text)">Para insanlar var olduğu için var. Bu nedenle her insan, sadece insan olduğu için paradan eşit pay almalıdır.</strong></p><p><em style="color:var(--gold)">"Para insanlar var olduğu için var. Bundan fazlası değil, bundan azı değil."</em></p>',
  'nodes-title':'Aktif Node\'lar — Mevcut Ağ Topolojisi',
  'nodes-desc':'Aequitas ağı şu anda iki coğrafi olarak dağıtılmış node üzerinde çalışıyor. Her ikisi de blok üretimine, durum senkronizasyonuna ve API hizmetine katılıyor. libp2p aracılığıyla eşler arası iletişim kuruyor ve HTTP aracılığıyla blok durumunu senkronize ediyorlar. Ağ ek node\'ları desteklemek üzere tasarlanmıştır.',
  'node1':'Node 1 — Railway (Birincil)','node1-desc':'Birincil API · Blok üreticisi · UBI dağıtımı · P2P Bootstrap · PostgreSQL · MetaMask için RPC',
  'node2':'Node 2 — Render (İkincil)','node2-desc':'İkincil API · Blok üreticisi · P2P eşi · HTTP senkronizasyonu · Paylaşılan PostgreSQL durumu',
  'run-node-title':'Kendi Node\'unu Çalıştır — Ağı Güvence Altına Almaya Yardım Et',
  'run-node-desc':'Herkes bir Aequitas node\'u çalıştırabilir — izin, stake veya başvuru gerekmez. Node\'lar blok üretimine katılır ve insan kaydını doğrular. Node operatörleri, Doğrulayıcı Havuzu aracılığıyla protokol ücretlerinden pay kazanır (tüm takas ücretlerinin %40\'ı, günlük dağıtılır).',
  'bootstrap-title':'Yeni Node Bağla','bootstrap-desc':'Kendi Aequitas node\'unu çalıştırmak için PRIMARY_NODE_URL=https://aequitas.digital ortam değişkenini ayarla. Node\'un tam zincir durumunu otomatik olarak senkronize edecek ve blok üretimine başlayacak.',
  'tech-title':'Teknik Özellikler','mm-config':'MetaMask Yapılandırması',
  'k-lang':'Dil','k-src':'Kaynak Kodu','evm-yes':'Evet — JSON-RPC /rpc · MetaMask uyumlu',
  'proto-label':'Aequitas V7 Protokolü — Teknik Dokümantasyon',
  'ca-title':'Sözleşme Adresleri','ca-text':'Zincir: Aequitas Chain (Zincir ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Ana): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7, tüm Aequitas ekonomisinin tek gerçek kaynağıdır. Her AEQ bakiyesi, her insan kaydı, her UBI ödemesi ve her servet tavanı uygulaması, bu tek değiştirilemez sözleşme tarafından yönetilir. Yönetici anahtarı yok, yükseltme proxy\'si yok, mantığının tek bir satırını değiştirebilecek yönetişim oyu yok. Bugün çalışan kod on yıl sonra da çalışacak koddur.',
  'poa-title':'1. HAYAT KANITI — Hareketsiz Bakiye Kurtarma','poa-text':'<p>İnsanlar ölünce veya kalıcı olarak yetersiz hale gelince AEQ\'ya ne olur? Bitcoin\'de kaybedilen cüzdanlar, kalıcı olarak kaybedilen arz anlamına gelir. Aequitas bunu çok aşamalı bir hareketsizlik kurtarma sistemiyle çözer.</p>',
  'poa-box':'Yıl 0–2: Normal kullanım — kısıtlama yok<br>Yıl 2: Uyarı 1 — Vasi adına yanıt verebilir<br>Yıl 2+60g: Uyarı 2 — artan aciliyet<br>Yıl 2+120g: Uyarı 3 — son bildirim<br>Yıl 2+180g: AEQ kişisel EMANET\'e taşındı (hâlâ kurtarılabilir)<br>Yıl 4: Hâlâ hareketsizse — EMANET UBI Havuzuna serbest bırakıldı',
  'guard-title':'2. VASİ SİSTEMİ — İnsani Güvence','guard-text':'<p>Ya biri hastanede ya da başka bir nedenle aylarca cihazına erişemiyorsa? Vasi sistemi, güvenilen bir kişinin — başka bir doğrulanmış insanın — cüzdan sahibinin hâlâ hayatta olduğunu onaylamasına izin verir. Vasinin kesinlikle sıfır finansal erişimi vardır: yalnızca hareketsizlik zamanlayıcısını sıfırlayan tek bir işlevi çağırabilir.</p>',
  'guard-box':'İnsan başına 1 Vasi · Aequitas\'ta doğrulanmış insan olmalı<br>Vasi YALNIZCA confirmAlive() çağırabilir — sıfır işlem hakkı<br>Vasi fon taşıyamaz, AEQ transfer edemez veya cüzdana erişemez<br>Vasi başına en fazla 3 korunan · 7 günlük kilit · Döngüsel ilişkiye izin yok',
  'dem-title':'3. GECİKME ÜCRETİ — Biriktirme Karşıtı Mekanizma',
  'dem-box':'Hız: 3 aylık hareketsizlikten sonra ayda %0,5 (sürekli, kademeli değil)<br>Herhangi bir transfer, takas veya likidite işlemi zamanlayıcıyı otomatik olarak sıfırlar<br>Bozunan AEQ dört havuza yeniden dağıtılır — asla yakılmaz<br>14 günlük uyarı bir kez gösterilir · 7 günlük uyarı her aktif oturumda tekrarlanır',
  'dem-text':'<p>Gecikme ücreti, para üzerindeki bir tutma maliyetidir — biriktirmeyi pahalı, dolaşımı çekici kılan negatif bir faiz oranı. Wörgl Deneyi (Avusturya, 1932), gecikme ücretli bir para birimi kullandı ve bir yılda yerel işsizliği %25 azalttı.</p>',
  'cap-title':'4. SERVET TAVANI — Matematiksel Adalet Uygulaması','cap-box':'Bootstrap tavanı: max(5,min(N,25))× mevcut ortalama AEQ bakiyesi<br>1–4 insan: 5× · insan başına +1× · 25+: kalıcı 25×<br>4 protokol havuzu adresi dışındaki TÜM adresler için geçerli<br>Fazla AEQ anında yeniden dağıtılır · Manuel müdahale yok',
  'ubi-title':'5. EVRENSEL TEMEL GELİR — Günlük Yeniden Dağıtım','ubi-box':'UBI Havuzu Gelir Kaynakları:<br>· AEQ↔tUSD AMM havuzundan tüm takas ücretlerinin %20\'si<br>· Servet tavanı uygulamasından taşma<br>· Hareketsiz hesaplardan gecikme ücretleri<br>· 4 yıl sonra serbest bırakılan hareketsiz emanet<br><br>Dağıtım: Her 24 saatte bir, tüm UBI Havuzu bakiyesi tüm kayıtlı doğrulanmış insanlar arasında eşit olarak bölünür.',
  'inf-title':'6. ALGORİTMİK ENFLASYON YOK — Sabit Arz Formülü','inf-box':'Yeni AEQ yaratan TEK olay: yeni bir doğrulanmış insan kaydolur.<br><br>Toplam Arz = Doğrulanmış İnsanlar × 1.000 AEQ<br><br>Bu bir politika değil — protokol tarafından zorlanır. Hiçbir yönetici ek AEQ basamaz.',
  'btn-download-app':'AEQUİTASBİO UYGULAMASINI İNDİR',
  'swap-title':'🔄 AEQ ↔ tUSD Takas Et','swap-sub':'Yerel likidite havuzu üzerinden AEQ\'yu tUSD (simüle edilmiş test doları) ile takas et. %0,1 ücret yalnızca takaslar için geçerlidir — insanlar arasındaki normal AEQ transferleri tamamen ücretsiz kalır.',
  'swap-priv-bar':'🔒 Yalnızca %0,1 takas ücreti · AEQ\'dan AEQ\'ya transferler ücretsiz · tUSD gerçek değeri olmayan test para birimidir',
  'swap-your-aeq':'Senin AEQ','swap-your-tusd':'Senin tUSD','swap-aeq-to-tusd':'AEQ → tUSD','swap-tusd-to-aeq':'tUSD → AEQ',
  'swap-fee-est':'Protokol ücreti (%0,1)','swap-details-hdr':'Takas Detayları',
  'swap-out-lbl':'Alacaksın (tahmini)','swap-impact-lbl':'Fiyat etkisi','swap-rate-lbl':'Döviz kuru',
  'swap-depth-lbl':'Havuz Bileşimi','amm-title':'x × y = k — Sabit Çarpım AMM',
  'amm-text':'AEQ\'yu tUSD karşılığında takas ettiğinde, AEQ rezervi büyür ve tUSD rezervi küçülür — çarpımları her zaman k\'ya eşit kalır. Daha büyük takaslar daha fazla fiyat etkisine neden olur. %0,1 ücreti formül uygulanmadan önce düşülür.',
  'swap-btn-conn':'🦊 METAMASK BAĞLA','swap-btn-go':'🔄 TAKAS ET',
  'swap-log-hint':'// Takas yapmak için cüzdan bağla...',
  'swap-no-liquidity':'Henüz tUSD yok mu?','swap-faucet-desc':'Kayıtlı insanlar bir kez test tUSD talep edebilir','swap-btn-faucet':'💧 TEST tUSD TALEP ET',
  'swap-addliq-title':'Likidite Sağla','swap-addliq-desc':'İlk yatıran ol — oranın başlangıç fiyatını belirler.','swap-btn-addliq':'💧 LİKİDİTE EKLE',
  'swap-lp-title':'LP Pozisyonun','swap-lp-share':'Havuz Payı','swap-lp-withdrawable':'Çekilebilir',
  'swap-lp-pct-label':'% pozisyonun','swap-lp-youget':'Alacaksın','swap-btn-removeliq':'🔥 LİKİDİTE KALDIR',
  'swap-pool-title':'AEQ / tUSD — Havuz Durumu',
  'swap-pool-aeq':'AEQ Rezervi','swap-pool-tusd':'tUSD Rezervi','swap-pool-price':'Spot Fiyat',
  'swap-fee-bps':'Takas Ücreti','swap-fee-split':'Ücret Dağılımı','swap-fee-split-v':'%40 Doğrulayıcılar / %30 LP\'ler / %20 UBI / %10 Hazine',
  'swap-pools-addr-title':'Tokenomik Havuz Adresleri',
  'swap-validators':'Doğrulayıcılar (%40)','swap-lps':'Likidite Sağlayıcıları (%30)','swap-ubi':'UBI Havuzu (%20)','swap-treasury':'Hazine (%10)',
  'ubi-hero-title':'EVRENSEL TEMEL GELİR — UBI HAVUZU',
  'ubi-hero-sub':'Biriktirilmekte — bir sonraki ödeme tüm doğrulanmış insanlara eşit olarak dağıtılıyor:',
  'ubi-bal-lbl':'mevcut havuz bakiyesi','ubi-hero-desc':'Tümüne eşit bölünür · her 24 saatte ödenir · havuz sıfırlanır · minimum bakiye gerekmez',
  'ubi-how-fills':'UBI Havuzu Nasıl Dolar',
  'ubi-src-swap':'Takas Ücretleri','ubi-src-swap-d':'Her AEQ↔tUSD takası, %0,1 ücretinin %20\'sini katkıda bulunur. Daha fazla işlem = daha hızlı dolma.',
  'ubi-src-dem':'Gecikme Ücreti','ubi-src-dem-d':'Hareketsiz AEQ (3+ ay) ayda %0,5 bozunur. Bozunan miktarın %20\'si UBI\'ya gider.',
  'ubi-src-cap':'Servet Tavanı Taşması','ubi-src-cap-d':'Servet tavanını (max(5,min(N,25))× ortalama) aşan cüzdanlar anında kesilir. %20\'si UBI\'ya akar.',
  'pools4-header':'Dört yeniden dağıtım havuzunun tamamı',
  'ubi-see-above':'yukarıdaki geri sayımı gör','ubi-timer-above':'⏰ geri sayım yukarıda gösterildi','pool-t-timer':'Birikiyor — zamanlayıcı yok',
  'usp-headline':'Tarihte ilk kez — herkes eşit başlıyor',
  'usp-sub':'Android akıllı telefonun varsa katılabilirsin. Banka yok, kripto bilgisi yok, yatırım yok.',
  'usp-c1-title':'0,00 Başlangıç Yatırımı','usp-c1-desc':'Kayıt tamamen gas\'sız. ETH, MATIC veya kredi kartı gerekmez. Protokol tüm işlem maliyetlerini öder.',
  'usp-c2-title':'Her insan için 1.000 AEQ','usp-c2-desc':'Milyarder ya da geçimlik çiftçi — herkes tam olarak 1.000 AEQ alır. Fazlası değil, azı değil. Eşit başlangıç, matematiksel garanti.',
  'usp-c3-title':'Herkese erişilebilir','usp-c3-desc':'Banka hesabı, kredi kartı veya kimlik belgesi gerekmez. Kayıt, uygun fiyatlı biyometrik donanım kiti (parmak izi tarayıcı + nabız sensörü, ~$15) kullanır — küresel erişim için tasarlandı.',
  'usp-c4-title':'Sonsuza kadar günlük UBI','usp-c4-desc':'Kaydolduktan sonra, her gün otomatik olarak UBI ödemelerinden pay alırsın — her gün, hiçbir işlem gerektirmez.',
  'v7-intro-title':'AequitasV7 Nedir?',
  'v7-intro-text':'AequitasV7, Aequitas protokolünün merkezi akıllı sözleşmesidir. "V7", adalet sözleşmesinin 7. ana sürümüdür. Aequitas Chain\'de (Zincir ID 1926) değiştirilemez şekilde dağıtılmıştır ve her şeyi yönetir: insan kaydı, ZK doğrulaması, bakiye yönetimi, servet tavanı, UBI dağıtımı, takas ücretleri. Hiçbir yönetici onu güncelleyemez. Altı mekanizma kendi kendini güçlendiren bir sistem oluşturur.',
  'explore-title':'Aequitas\'ı Keşfet',
  'expl-score':'Eşitlik Skoru','expl-score-d':'Canlı Gini katsayısı · Aequitas Endeksi · gerçek zamanlı servet dağılımı',
  'expl-economy':'UBI ve Yeniden Dağıtım Havuzları','expl-economy-d':'Günlük UBI geri sayımı · 4 on-chain havuz · demurrage · Protokol Aşamaları',
  'expl-charts':'Grafikler ve Tarih','expl-charts-d':'Gini geçmişi · Lorenz eğrisi · servet tavanı bootstrap kaydırıcısı · Aequitas\'ın hikayesi',
  'expl-v7':'Protokol V7 Dokümantasyonu','expl-v7-d':'AequitasV7 sözleşmesi · 6 mekanizma · ZK kanıtı · servet tavanı · demurrage · değiştirilemez kod',
  'expl-explorer':'Blok Gezgini','expl-explorer-d':'Canlı BlockDAG · doğrulayıcıyı, hash\'i, işlemleri, üst hash\'leri görmek için herhangi bir bloğa tıklayın',
  'swap-sell-label':'Sat','swap-receive-label':'Al',
  'expl-network':'Ağ ve Düğümler','expl-network-d':'Düğüm topolojisi · kendi düğümünü çalıştır · teknik özellikler · Zincir ID 1926'
},
fr:{
  'logo-sub':'PREUVE D\'HUMANITÉ','live':'EN DIRECT',
  'tab-register':'🔐 S\'inscrire','tab-explorer':'🔍 Explorateur','tab-humans':'👥 Humains','tab-index':'📊 Index','tab-network':'🌐 Réseau','tab-protocol':'📜 Protocole V7','tab-swap':'🔄 Échanger',
  'reg-title':'🔐 S\'inscrire en tant qu\'humain vérifié',
  'reg-sub':'Rejoignez le réseau Aequitas et recevez 1 000 AEQ de Revenu de Base Universel. L\'inscription est unique, permanente et totalement sans frais. Aucune donnée personnelle n\'est stockée.',
  'app-title':'INSCRIPTION VIA L\'APPLICATION ANDROID',
  'app-text':'La Preuve d\'Humanité utilise un système biométrique physique à 3 facteurs. Phase 1 : capteur optique R503 scanne les 10 empreintes + MAX30102 PPG confirme le pouls vital. Phase 2 : caméra IR des veines de la main (unicité 1 sur 10⁷). Phase 3 : scan de l\'iris — étalon-or, 1 sur 10⁷⁸, entièrement indépendant de l\'appareil. Une preuve ZK Groth16 engage tous les facteurs sans révéler aucune donnée biométrique. Vos 1 000 AEQ sont crédités automatiquement à la vérification.',
  's1t':'Scan Biométrique','s1d':'AequitasBio scanne les 10 empreintes (capteur optique R503) + pouls MAX30102 PPG confirme la vivacité. Phase 2 : IR veines de la main. Phase 3 : iris. Les données brutes ne quittent jamais l\'appareil.',
  's2t':'Génération de Preuve ZK','s2d':'La preuve ZK Groth16 engage tous les facteurs biométriques : commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier lié au corps, pas au téléphone — perdre l\'appareil ne peut pas créer une seconde identité.',
  's3t':'Connecter le Portefeuille','s3d':'L\'app ouvre MetaMask · connectez votre portefeuille Ethereum · la preuve est liée cryptographiquement à votre adresse',
  's4t':'1 000 AEQ Accordés','s4d':'Inscription confirmée sur le BlockDAG en 6 secondes · 1 000 AEQ crédités instantanément · identité enregistrée en permanence',
  'priv-bar':'🔒 R503 10 Empreintes · MAX30102 Vivacité · Phase 2 : IR Veines Main · Phase 3 : Iris (10⁷⁸) · Groth16 ZKP · Données ne quittent jamais l\'appareil · Un humain · Pour toujours',
  'conn-wallet':'PORTEFEUILLE CONNECTÉ','proof-recv':'⚡ PREUVE ZK REÇUE','proof-hint':'Connecter un portefeuille pour s\'inscrire',
  'btn-conn':'🦊 CONNECTER METAMASK','btn-reg':'🔐 INSCRIPTION ON-CHAIN',
  'btn-web-reg':'🌐 INSCRIPTION VIA NAVIGATEUR (WebAuthn)',
  'web-reg-warn':'⚠ Lié à l\'appareil : Cette identité est liée à cet appareil et navigateur. Non transférable. Pour identité multi-appareils, utilisez l\'app Android Aequitas.',
  'reg-log-hint':'// Ouvrir l\'app Android Aequitas pour générer votre preuve, puis revenir ici...',
  'reg-details':'Détails d\'inscription','k-network':'Réseau','k-chainid':'ID de chaîne','k-grant':'Allocation UBI',
  'k-fee':'Frais de gaz','free':'GRATUIT — totalement sans frais','k-limit':'Inscriptions','k-limit-v':'Une fois par humain · permanent · immuable',
  'k-bio':'Données biométriques','never-stored':'Jamais stockées — restent sur votre appareil',
  'k-proof':'Système de preuve','k-conf':'Confirmation','k-conf-v':'En 6 secondes (1 bloc)',
  'k-sybil':'Protection Sybil','k-sybil-v':'Une identité par biométrie · verrouillage permanent',
  'live-stats':'Statistiques de la chaîne en direct',
  's-height':'Hauteur de bloc','s-height-sub':'Nouveau bloc toutes les ~6s · BlockDAG · Production parallèle',
  's-humans':'Humains vérifiés','s-humans-sub':'ZKP biométrique · Une personne, un portefeuille, pour toujours',
  's-supply':'Offre totale','s-supply-sub':'Toujours = Humains × 1 000 AEQ',
  's-index':'Index Aequitas','s-index-sub':'0 = égalité parfaite · 100 = inégalité maximale',
  's-uptime':'Disponibilité','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Preuve d\'Humanité','ib-poh-t':'Chaque détenteur d\'AEQ doit prouver qu\'il est un humain vivant unique. Pas de robots, sociétés ni IA. Données biométriques jamais partagées.',
  'ib-fair':'Distribution radicalement équitable','ib-fair-t':'Chaque humain vérifié reçoit exactement 1 000 AEQ. Pas de pré-minage ni d\'allocation fondateurs. Offre = Humains × 1 000.',
  'ib-dag':'Architecture BlockDAG','ib-dag-t':'Plusieurs blocs produits simultanément et fusionnés. Débit plus élevé, latence plus faible.',
  'ib-gas':'Vraiment sans frais','ib-gas-t':'Inscription et transferts AEQ gratuits. Pas d\'ETH, BNB ou MATIC. Pas de carte bancaire nécessaire.',
  'recent-blocks':'Blocs récents','blocks-desc':'MERGE = plusieurs parents fusionnés (BlockDAG). TX = transaction d\'inscription. Temps de bloc : ~6 secondes.',
  'loading':'Chargement des blocs...','net-info':'Informations réseau','k-chain':'Nom de chaîne','k-symbol':'Symbole','k-btime':'Temps de bloc',
  'k-cons':'Consensus','k-nodes':'Nœuds actifs','k-storage':'Stockage','add-mm':'🦊 AJOUTER À METAMASK','k-dec':'Décimales',
  'btn-add-mm':'+ AJOUTER LE RÉSEAU AEQUITAS',
  'phil':'"L\'argent existe parce que les gens existent.<br>Rien de plus, rien de moins."','phil-sub':'— LE PRINCIPE AEQUITAS —',
  'humans-title':'Humains vérifiés sur Aequitas Chain',
  'h-what':'Qu\'est-ce qu\'un humain vérifié ?','h-what-t':'Adresse de portefeuille cryptographiquement prouvée comme appartenant à un humain vivant unique. La vérification utilise un système matériel à 3 facteurs : R503 scanne les 10 empreintes ; MAX30102 PPG confirme le pouls vital ; Phase 2 : IR veines de la main (1 sur 10⁷) ; Phase 3 : iris (1 sur 10⁷⁸). Seule une preuve ZK Groth16 est transmise. Aucune donnée biométrique ne quitte l\'appareil.',
  'h-zkp':'Système de preuve ZK','h-zkp-t':'Aequitas utilise Groth16 sur BN128 — même courbe qu\'Ethereum et Zcash. ~200 octets, ~10ms. commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier lié au corps : perdre son téléphone ne crée pas une seconde identité. Aucune donnée biométrique n\'est jamais stockée.',
  'h-sybil':'Prévention Sybil','h-sybil-t':'Phase 1 : 10 empreintes + vivacité MAX30102 (pouls PPG, rejette moulages/replays). Phase 2 : IR veines de la main — caractéristique corporelle interne, impossible à copier, 1 sur 10⁷, différente chez les jumeaux identiques. Phase 3 : iris — 1 sur 10⁷⁸, étalon-or mondial. Nullifier = keccak256(iris‖vein‖domain). Un humain, un portefeuille, pour toujours.',
  'h-global':'Inclusion financière mondiale','h-global-t':'Pas de compte bancaire, carte de crédit ou crypto préalable requis. Un smartphone Android avec capteur biométrique suffit.',
  'h-bio-hw':'Feuille de Route Matériel Biométrique','h-bio-hw-t':'Phase 1 (active) : scanner d\'empreintes optique R503 — hash combiné des 10 doigts. Vivacité MAX30102 PPG. Phase 2 (prévue) : ESP32-CAM + LED IR 850 nm — imagerie des veines de la main, unicité 1 sur 10⁷. Phase 3 (prévue) : module iris IR — 240+ degrés de liberté, 1 sur 10⁷⁸, entièrement indépendant de l\'appareil, les jumeaux identiques diffèrent.',
  'reg-humans':'Humains inscrits','h-desc':'Chaque adresse vérifiée comme humain unique via ZKP biométrique. Chacun a reçu 1 000 AEQ. Permanent, immuable, on-chain.',
  'no-humans':'Aucun humain inscrit pour l\'instant.\n\nTéléchargez l\'application Android Aequitas et soyez le premier !',
  'reg-stats':'Statistiques du registre','total-humans':'Total d\'humains',
  'idx-title':'Index Aequitas — Score d\'égalité économique en temps réel',
  'idx-desc':'L\'Index Aequitas est dérivé du <strong style="color:var(--teal)">coefficient de Gini</strong> — la norme internationale pour mesurer les inégalités (Banque mondiale, OCDE, ONU). <strong style="color:var(--neon)">0 = égalité parfaite</strong>. <strong style="color:var(--red)">100 = concentration totale</strong>. Objectif : Gini sous 0,35.',
  'gini-what-title':'Qu\'est-ce que le coefficient de Gini ?',
  'gini-what-text':'Développé par Corrado Gini (1912). Mesure la distribution des richesses. Échelle : 0 (tous égaux) à 1 (une personne détient tout). Utilisé par la Banque mondiale, l\'OCDE, l\'ONU.',
  'gini-calc-title':'Comment l\'Index est-il calculé ?',
  'gini-calc-text':'Tous les soldes AEQ collectés. Différence absolue moyenne entre toutes les paires, normalisée par n² et le solde moyen. Résultat × 100 = Index Aequitas.',
  'gini-why-title':'Pourquoi le Gini ?',
  'gini-why-text':'Un simple ratio riche/pauvre est manipulable. Le Gini capture la distribution complète en un seul chiffre auditable, publié on-chain — transparent et vérifiable mondialement.',
  'curr-idx':'Index actuel','bar-0':'0 — Égalité parfaite','bar-100':'100 — Inégalité max','wcap-lbl':'Plafond de richesse :','wcap-mult':'Multiplicateur :','wcap-avg':'Solde moyen :',
  'gini':'Coefficient de Gini','gini-desc':'0 = égal · 1 = inégal',
  'supply-desc':'Toujours = Humains × 1 000 AEQ',
  'phase':'Phase du protocole','phase-desc':'Avance automatiquement par nombre d\'humains',
  'humans-desc':'Humains uniques vérifiés biométriquement',
  'pools-title':'Pools de redistribution',
  'pools-desc':'Chaque frais de swap, demurrage et dépassement du plafond est divisé entre quatre pools. Tous versent quotidiennement.',
  'vel-pool':'Pool des validateurs','vel-pool-desc':'40% de tous les frais → opérateurs de nœuds qui sécurisent le réseau',
  'liq-pool':'Pool de liquidité','liq-pool-desc':'30% de tous les frais → fournisseurs de liquidité, proportionnellement aux parts LP',
  'ubi-pool':'Pool UBI','ubi-pool-desc':'20% de tous les frais → tous les humains vérifiés également, toutes les 24 heures',
  'treasury':'Trésorerie','treasury-desc':'10% de tous les frais → développement et maintenance du protocole',
  'phases-title':'Phases du protocole',
  'phases-desc':'Plafond bootstrap Phase 0 : max(5, min(N, 25))× solde moyen. 1–4 humains : 5×. Chaque humain ajoute 1×. 25+ humains : verrouillé à 25×. Transitions automatiques.',
  'p0':'Bootstrap · &lt;100 humains · Plafond : max(5,min(N,25))× moyen · 5×→25× · Actuellement actif',
  'p1':'Croissance · 100–10 000 humains · Plafond : 25× solde moyen',
  'p2':'Stabilité · 10 000–1M humains · Plafond : 25× solde moyen',
  'p3':'Maturité · 1M+ humains · Plafond : 25× solde moyen',
  'wealth-cap-explain':'Plafond Phase 0 : max(5, min(N, 25))× solde moyen. 1–4 humains : 5×. Chaque humain +1×. 25+ : verrouillé à 25×.',
  'demurrage-title':'Demurrage — Incitation à la circulation',
  'demurrage-desc':'Les soldes AEQ inactifs perdent lentement de la valeur pour décourager l\'accumulation.',
  'dem-rate-k':'Taux de décroissance','dem-rate-v':'0,5 % par mois (continu)',
  'dem-grace-k':'Période de grâce','dem-grace-v':'3 mois d\'inactivité avant début de décroissance',
  'dem-reset-k':'Réinitialisation','dem-reset-v':'Tout transfert, swap ou action de liquidité remet le compteur à zéro',
  'dem-dest-k':'L\'AEQ décroissant va vers','dem-dest-v':'Pools de redistribution (40/30/20/10)',
  'dem-warn-k':'Système d\'avertissement','dem-warn-v':'Avis 14 jours (une fois) + rappel 7 jours à chaque connexion',
  'story-title':'L\'histoire d\'Aequitas',
  'story-text':'<p>En 2009, Satoshi Nakamoto publie Bitcoin. Révolution genuïne — mais les premiers mineurs accumulent des millions à coût quasi nul. En 2021, le top 1% contrôle plus de 90% du Bitcoin. Gini Bitcoin &gt; 0,85.</p><p><span style="color:var(--gold)">Aequitas</span> — latin pour « équité » — répond : <em style="color:var(--gold)">« Quelle serait une cryptomonnaie conçue pour être juste envers chaque humain ? »</em></p><p><strong style="color:var(--text)">L\'argent existe parce que les gens existent. Donc chaque personne devrait avoir une part égale.</strong></p><p><em style="color:var(--gold)">« L\'argent existe parce que les gens existent. Rien de plus, rien de moins. »</em></p>',
  'nodes-title':'Nœuds actifs — Topologie réseau actuelle',
  'nodes-desc':'Le réseau Aequitas fonctionne sur deux nœuds géographiquement distribués participant à la production de blocs, synchronisation d\'état et service API. Nœuds supplémentaires bienvenus.',
  'node1':'Nœud 1 — Railway (Principal)','node1-desc':'API principal · Producteur de blocs · Distribution UBI · Bootstrap P2P · PostgreSQL · RPC MetaMask',
  'node2':'Nœud 2 — Render (Secondaire)','node2-desc':'API secondaire · Producteur de blocs · Pair P2P · Sync HTTP · État PostgreSQL partagé',
  'run-node-title':'Exécuter votre propre nœud','run-node-desc':'N\'importe qui peut exécuter un nœud Aequitas — sans permission, sans stake. Opérateurs gagnent 40% des frais de swap distribués quotidiennement.',
  'bootstrap-title':'Connecter un nouveau nœud','bootstrap-desc':'Définissez PRIMARY_NODE_URL=https://aequitas.digital dans votre environnement. Votre nœud synchronise automatiquement l\'état complet.',
  'tech-title':'Spécifications techniques','mm-config':'Configuration MetaMask',
  'k-lang':'Langue','k-src':'Source','evm-yes':'Oui — JSON-RPC /rpc · Compatible MetaMask',
  'proto-label':'Protocole Aequitas V7 — Documentation technique',
  'ca-title':'Adresses des contrats',
  'ca-text':'Chaîne : Aequitas Chain (Chain ID : 1926 · 0x786)<br>RPC : https://aequitas.digital/rpc<br><br>BioVerifier : 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 : 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7 est l\'unique source de vérité pour toute l\'économie Aequitas. Aucune clé d\'administration ni vote de gouvernance ne peut modifier sa logique. Le code actuel fonctionnera dans dix ans.',
  'poa-title':'1. PREUVE DE VIE','poa-text':'<p>Quand les gens décèdent, leurs AEQ retournent progressivement à la communauté via le pool UBI plutôt que d\'être perdus comme dans Bitcoin.</p>',
  'poa-box':'Années 0–2 : Utilisation normale<br>Année 2 : Avertissement 1 — Gardien peut répondre<br>Année 2+60j : Avertissement 2<br>Année 2+120j : Avertissement 3<br>Année 2+180j : AEQ en séquestre personnel<br>Année 4 : Si inactif — retourne au Pool UBI',
  'guard-title':'2. SYSTÈME DE GARDIEN','guard-text':'<p>Un Gardien de confiance (autre humain vérifié) peut confirmer qu\'une personne est encore en vie, sans aucun droit de transaction.</p>',
  'guard-box':'1 Gardien par humain · doit être humain vérifié Aequitas<br>Gardien peut UNIQUEMENT appeler confirmAlive() · zéro droit financier<br>Gardien NE PEUT PAS déplacer des fonds · Max 3 protégés · Timelock 7j',
  'dem-title':'3. DEMURRAGE — Anti-accumulation',
  'dem-box':'Taux : 0,5%/mois après 3 mois de grâce<br>Réinitialisation à chaque transfert, swap ou action de liquidité<br>AEQ décroissant redistribué dans les pools (non brûlé)',
  'dem-text':'<p>Précédent : Wörgl (Autriche, 1932) — réduction du chômage de 25% en un an. Chiemgauer (Allemagne, 2003) — fonctionne depuis plus de 20 ans.</p>',
  'cap-title':'4. PLAFOND DE RICHESSE','cap-box':'Plafond : max(5,min(N,25))× solde moyen<br>1–4 humains : 5× · +1× par humain · 25+ : 25× permanent<br>Excès immédiatement redistribué · Aucune intervention manuelle',
  'ubi-title':'5. REVENU DE BASE UNIVERSEL','ubi-box':'Sources : Frais de swap (20%) · Dépassement du plafond · Demurrage<br><br>Quotidien : Pool UBI divisé également entre tous les humains. Pool remis à zéro après chaque distribution.',
  'inf-title':'6. PAS D\'INFLATION ALGORITHMIQUE','inf-box':'Seul événement créant de l\'AEQ : un nouvel humain vérifié s\'inscrit.<br><br>Offre totale = Humains vérifiés × 1 000 AEQ — toujours, exactement.',
  'btn-download-app':'TÉLÉCHARGER AEQUITASBIO',
  'swap-title':'🔄 Échanger AEQ ↔ tUSD','swap-sub':'Échangez AEQ contre tUSD (dollar test) via le pool de liquidité natif. Frais 0,1% uniquement pour les swaps — transferts AEQ ordinaires totalement gratuits.',
  'swap-priv-bar':'🔒 Seulement 0,1% de frais · Transferts AEQ→AEQ gratuits · tUSD est une monnaie test sans valeur réelle',
  'swap-your-aeq':'Votre AEQ','swap-your-tusd':'Votre tUSD','swap-aeq-to-tusd':'AEQ → tUSD','swap-tusd-to-aeq':'tUSD → AEQ',
  'swap-fee-est':'Frais de protocole (0,1%)','swap-details-hdr':'Détails de l\'échange',
  'swap-out-lbl':'Vous recevez (est.)','swap-impact-lbl':'Impact sur le prix','swap-rate-lbl':'Taux de change',
  'swap-depth-lbl':'Composition du Pool','amm-title':'x × y = k — AMM à produit constant',
  'amm-text':'Lors d\'un swap, les réserves AEQ augmentent et les réserves tUSD diminuent — produit toujours égal à k. Swaps plus grands = plus grand impact sur le prix.',
  'swap-btn-conn':'🦊 CONNECTER METAMASK','swap-btn-go':'🔄 ÉCHANGER',
  'swap-log-hint':'// Connecter un portefeuille pour échanger...',
  'swap-no-liquidity':'Pas encore de tUSD ?','swap-faucet-desc':'Humains inscrits peuvent réclamer du tUSD test une fois','swap-btn-faucet':'💧 RÉCLAMER tUSD TEST',
  'swap-addliq-title':'Fournir de la liquidité','swap-addliq-desc':'Soyez le premier à déposer — votre ratio fixe le prix initial.','swap-btn-addliq':'💧 AJOUTER LIQUIDITÉ',
  'swap-lp-title':'Votre position LP','swap-lp-share':'Part du Pool','swap-lp-withdrawable':'Retirable',
  'swap-lp-pct-label':'% de votre position','swap-lp-youget':'Vous recevrez','swap-btn-removeliq':'🔥 RETIRER LIQUIDITÉ',
  'swap-pool-title':'AEQ / tUSD — Statut du Pool',
  'swap-pool-aeq':'Réserve AEQ','swap-pool-tusd':'Réserve tUSD','swap-pool-price':'Prix Spot',
  'swap-fee-bps':'Frais de Swap','swap-fee-split':'Répartition des frais','swap-fee-split-v':'40% Validateurs / 30% LP / 20% UBI / 10% Trésorerie',
  'swap-pools-addr-title':'Adresses des Pools Tokenomiques',
  'swap-validators':'Validateurs (40%)','swap-lps':'Fournisseurs de Liquidité (30%)','swap-ubi':'Pool UBI (20%)','swap-treasury':'Trésorerie (10%)',
  'ubi-hero-title':'REVENU DE BASE UNIVERSEL — POOL UBI',
  'ubi-hero-sub':'Accumulation — prochain paiement distribué à tous les humains vérifiés dans :',
  'ubi-bal-lbl':'solde actuel du pool','ubi-hero-desc':'Divisé également · payé toutes les 24h · pool remis à zéro · solde minimum non requis',
  'ubi-how-fills':'Comment le Pool UBI se remplit',
  'ubi-src-swap':'Frais de Swap','ubi-src-swap-d':'Chaque swap AEQ↔tUSD contribue 20% de ses frais. Plus d\'échanges = remplissage plus rapide.',
  'ubi-src-dem':'Demurrage','ubi-src-dem-d':'AEQ inactif (3+ mois) décroît 0,5%/mois. 20% du décroissant va à l\'UBI.',
  'ubi-src-cap':'Dépassement du Plafond','ubi-src-cap-d':'Portefeuilles dépassant le plafond immédiatement rognés. 20% afflue vers l\'UBI.',
  'pools4-header':'Les quatre pools de redistribution',
  'ubi-see-above':'voir compte à rebours ci-dessus','ubi-timer-above':'⏰ compte à rebours affiché ci-dessus','pool-t-timer':'Accumulation — pas de minuterie',
  'usp-headline':'Pour la première fois dans l\'histoire — tout le monde commence à égalité',
  'usp-sub':'Si vous avez un smartphone Android, vous êtes éligible. Pas de banque, pas de crypto, pas d\'investissement.',
  'usp-c1-title':'0 € d\'investissement initial','usp-c1-desc':'Inscription totalement sans frais. Pas d\'ETH ni de carte bancaire. Le protocole paie tous les frais.',
  'usp-c2-title':'1 000 AEQ pour chaque humain','usp-c2-desc':'Milliardaire ou agriculteur — tous reçoivent exactement 1 000 AEQ. Égalité garantie mathématiquement.',
  'usp-c3-title':'Accessible à tous','usp-c3-desc':'Pas de compte bancaire, carte de crédit ni pièce d\'identité. L\'inscription utilise un kit biométrique abordable (scanner d\'empreintes + capteur de pouls, ~15 €) — conçu pour un accès mondial.',
  'usp-c4-title':'UBI quotidien pour toujours','usp-c4-desc':'Une fois inscrit, votre part des paiements UBI arrive automatiquement chaque jour — sans aucune action.',
  'v7-intro-title':'Qu\'est-ce qu\'AequitasV7 ?',
  'v7-intro-text':'AequitasV7 est le contrat intelligent central d\'Aequitas. Déployé de manière immuable sur Aequitas Chain (ID 1926). Gère tout : inscription humaine, vérification ZK, soldes, plafond de richesse, UBI, frais de swap. Aucun administrateur ne peut le modifier.',
  'explore-title':'Explorer Aequitas',
  'expl-score':'Score d\'égalité','expl-score-d':'Coefficient de Gini en direct · Index Aequitas · distribution des richesses en temps réel',
  'expl-economy':'UBI &amp; Redistribution','expl-economy-d':'Compte à rebours UBI · 4 pools on-chain · demurrage · Phases du protocole',
  'expl-charts':'Graphiques &amp; Historique','expl-charts-d':'Historique Gini · courbe de Lorenz · curseur du plafond · L\'histoire d\'Aequitas',
  'expl-v7':'Docs Protocole V7','expl-v7-d':'Contrat AequitasV7 · 6 mécanismes · preuve ZK · plafond · demurrage · code immuable',
  'expl-explorer':'Explorateur de blocs','expl-explorer-d':'BlockDAG en direct · cliquez sur un bloc pour voir validateur, hash, transactions',
  'swap-sell-label':'Vendre','swap-receive-label':'Recevoir',
  'expl-network':'Réseau &amp; Nœuds','expl-network-d':'Topologie des nœuds · exécuter votre propre nœud · spécifications · Chain ID 1926'
},
pt:{
  'logo-sub':'PROVA DE HUMANIDADE','live':'AO VIVO',
  'tab-register':'🔐 Registrar','tab-explorer':'🔍 Explorador','tab-humans':'👥 Humanos','tab-index':'📊 Índice','tab-network':'🌐 Rede','tab-protocol':'📜 Protocolo V7','tab-swap':'🔄 Trocar',
  'reg-title':'🔐 Registrar como Humano Verificado',
  'reg-sub':'Junte-se à rede Aequitas e receba 1.000 AEQ de Renda Básica Universal. Registro único, permanente e completamente sem taxas. Nenhum dado pessoal é armazenado.',
  'app-title':'REGISTRO VIA APLICATIVO ANDROID',
  'app-text':'A Prova de Humanidade usa um sistema biométrico físico de 3 fatores. Fase 1: sensor óptico R503 escaneia todas as 10 impressões + MAX30102 PPG confirma pulso vivo. Fase 2: câmera IR de veias da mão (unicidade 1 em 10⁷). Fase 3: scan de íris — padrão ouro, 1 em 10⁷⁸, totalmente independente do dispositivo. Uma prova ZK Groth16 compromete todos os fatores sem revelar dados biométricos. Seus 1.000 AEQ são creditados automaticamente na verificação.',
  's1t':'Scan Biométrico','s1d':'AequitasBio escaneia todas as 10 impressões (sensor óptico R503) + pulso MAX30102 PPG confirma vivacidade. Fase 2: IR de veias da mão. Fase 3: íris. Dados brutos nunca saem do dispositivo.',
  's2t':'Geração de Prova ZK','s2d':'A prova ZK Groth16 compromete todos os fatores biométricos: commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier vinculado ao corpo, não ao telefone — perder o dispositivo não pode criar uma segunda identidade.',
  's3t':'Conectar Carteira','s3d':'O app abre MetaMask · conecte sua carteira Ethereum · prova ligada criptograficamente ao seu endereço',
  's4t':'1.000 AEQ Concedidos','s4d':'Registro confirmado no BlockDAG em 6 segundos · 1.000 AEQ creditados instantaneamente · identidade registrada permanentemente',
  'priv-bar':'🔒 R503 Todas as 10 Impressões · MAX30102 Vivacidade · Fase 2: IR Veias da Mão · Fase 3: Íris (10⁷⁸) · Groth16 ZKP · Dados nunca saem do dispositivo · Um humano · Para sempre',
  'conn-wallet':'CARTEIRA CONECTADA','proof-recv':'⚡ PROVA ZK RECEBIDA','proof-hint':'Conectar carteira para registrar',
  'btn-conn':'🦊 CONECTAR METAMASK','btn-reg':'🔐 REGISTRAR ON-CHAIN',
  'btn-web-reg':'🌐 REGISTRAR VIA NAVEGADOR (WebAuthn)',
  'web-reg-warn':'⚠ Vinculado ao dispositivo: Esta identidade está vinculada a este dispositivo e navegador. Não transferível. Para identidade multi-dispositivo, use o App Android Aequitas.',
  'reg-log-hint':'// Abra o App Android Aequitas para gerar sua prova, depois retorne aqui...',
  'reg-details':'Detalhes do Registro','k-network':'Rede','k-chainid':'ID da Cadeia','k-grant':'Concessão UBI',
  'k-fee':'Taxa de Gás','free':'GRATUITO — completamente sem taxas','k-limit':'Registros','k-limit-v':'Uma vez por humano · permanente · imutável',
  'k-bio':'Dados Biométricos','never-stored':'Nunca armazenados — ficam no seu dispositivo',
  'k-proof':'Sistema de Prova','k-conf':'Confirmação','k-conf-v':'Em 6 segundos (1 bloco)',
  'k-sybil':'Proteção Sybil','k-sybil-v':'Uma identidade por biometria · bloqueio permanente',
  'live-stats':'Estatísticas ao Vivo da Cadeia',
  's-height':'Altura do Bloco','s-height-sub':'Novo bloco a cada ~6s · BlockDAG · Produção paralela',
  's-humans':'Humanos Verificados','s-humans-sub':'ZKP biométrico · Uma pessoa, uma carteira, para sempre',
  's-supply':'Oferta Total','s-supply-sub':'Sempre = Humanos × 1.000 AEQ',
  's-index':'Índice Aequitas','s-index-sub':'0 = igualdade perfeita · 100 = desigualdade máxima',
  's-uptime':'Disponibilidade','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Prova de Humanidade','ib-poh-t':'Cada detentor de AEQ deve provar criptograficamente que é um humano vivo único. Sem bots, corporações ou IA. Dados biométricos nunca saem do dispositivo.',
  'ib-fair':'Distribuição Radicalmente Justa','ib-fair-t':'Cada humano verificado recebe exatamente 1.000 AEQ no registro. Sem pré-mineração. Oferta = Humanos × 1.000.',
  'ib-dag':'Arquitetura BlockDAG','ib-dag-t':'Vários blocos produzidos simultaneamente e mesclados. Maior throughput, menor latência.',
  'ib-gas':'Verdadeiramente Sem Taxas','ib-gas-t':'Registro e transferências AEQ custam absolutamente nada. Sem ETH, BNB ou MATIC. Sem conta bancária.',
  'recent-blocks':'Blocos Recentes','blocks-desc':'MERGE = vários pais mesclados (BlockDAG). TX = transação de registro. Tempo de bloco: ~6 segundos.',
  'loading':'Carregando blocos...','net-info':'Informações de Rede','k-chain':'Nome da Cadeia','k-symbol':'Símbolo','k-btime':'Tempo de Bloco',
  'k-cons':'Consenso','k-nodes':'Nodes Ativos','k-storage':'Armazenamento','add-mm':'🦊 ADICIONAR AO METAMASK','k-dec':'Decimais',
  'btn-add-mm':'+ ADICIONAR REDE AEQUITAS',
  'phil':'"O dinheiro existe porque as pessoas existem.<br>Nada mais, nada menos."','phil-sub':'— O PRINCÍPIO AEQUITAS —',
  'humans-title':'Humanos Verificados na Aequitas Chain',
  'h-what':'O que é um Humano Verificado?','h-what-t':'Um Humano Verificado é um endereço de carteira criptograficamente provado como pertencendo a um humano vivo único. A verificação usa um sistema de hardware de 3 fatores: R503 escaneia todas as 10 impressões; MAX30102 PPG confirma pulso vivo; Fase 2: IR de veias da mão (1 em 10⁷); Fase 3: íris (1 em 10⁷⁸). Apenas uma prova ZK Groth16 é transmitida. Nenhum dado biométrico deixa o dispositivo.',
  'h-zkp':'Sistema de Prova ZK','h-zkp-t':'Aequitas usa Groth16 sobre BN128 — mesma curva do Ethereum e Zcash. ~200 bytes, ~10ms. commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier vinculado ao corpo: perder seu telefone não cria uma segunda identidade. Nenhum dado biométrico é armazenado.',
  'h-sybil':'Prevenção de Ataque Sybil','h-sybil-t':'Fase 1: todas as 10 impressões + vivacidade MAX30102 (pulso PPG, rejeita moldes/replays). Fase 2: IR de veias da mão — característica corporal interna, impossível de copiar, 1 em 10⁷, diferente em gêmeos idênticos. Fase 3: íris — 1 em 10⁷⁸, padrão ouro global. Nullifier = keccak256(iris‖vein‖domain). Um humano, uma carteira, para sempre.',
  'h-global':'Inclusão Financeira Global','h-global-t':'Sem conta bancária, cartão ou criptomoeda prévia. Apenas smartphone Android com sensor biométrico.',
  'h-bio-hw':'Roteiro de Hardware Biométrico','h-bio-hw-t':'Fase 1 (ativa): scanner de impressões óptico R503 — hash combinado de todos os 10 dedos. Vivacidade MAX30102 PPG. Fase 2 (planejada): ESP32-CAM + LED IR 850 nm — imageamento de veias da mão, unicidade 1 em 10⁷. Fase 3 (planejada): módulo de íris IR — 240+ graus de liberdade, 1 em 10⁷⁸, totalmente independente do dispositivo, gêmeos idênticos diferem.',
  'reg-humans':'Humanos Registrados','h-desc':'Cada endereço verificado como humano único via ZKP biométrico. Cada um recebeu 1.000 AEQ. Permanente, imutável, on-chain.',
  'no-humans':'Nenhum humano registrado ainda.\n\nBaixe o App Android Aequitas e seja o primeiro humano na cadeia!',
  'reg-stats':'Estatísticas do Registro','total-humans':'Total de Humanos',
  'idx-title':'Índice Aequitas — Pontuação de Igualdade Econômica em Tempo Real',
  'idx-desc':'O Índice Aequitas é derivado do <strong style="color:var(--teal)">coeficiente de Gini</strong> (Banco Mundial, OCDE, ONU). <strong style="color:var(--neon)">0 = igualdade perfeita</strong>. <strong style="color:var(--red)">100 = concentração total</strong>. Meta: Gini abaixo de 0,35.',
  'gini-what-title':'O que é o Coeficiente de Gini?',
  'gini-what-text':'Desenvolvido por Corrado Gini (1912). Mede a distribuição de riqueza. Escala: 0 (todos iguais) a 1 (uma pessoa detém tudo). Banco Mundial, OCDE, ONU.',
  'gini-calc-title':'Como o Índice é calculado?',
  'gini-calc-text':'Todos os saldos AEQ coletados. Diferença absoluta média entre todos os pares, normalizada por n² e saldo médio. Resultado × 100 = Índice Aequitas.',
  'gini-why-title':'Por que Gini?',
  'gini-why-text':'Um simples ratio rico/pobre é manipulável. O Gini captura a distribuição completa em um único número auditável, publicado on-chain — transparente e verificável globalmente.',
  'curr-idx':'Índice Atual','bar-0':'0 — Igualdade Perfeita','bar-100':'100 — Desigualdade Máx.','wcap-lbl':'Teto de Riqueza Atual:','wcap-mult':'Multiplicador:','wcap-avg':'Saldo médio:',
  'gini':'Coeficiente de Gini','gini-desc':'0 = igual · 1 = desigual',
  'supply-desc':'Sempre = Humanos × 1.000 AEQ',
  'phase':'Fase do Protocolo','phase-desc':'Avança automaticamente pelo número de humanos',
  'humans-desc':'Humanos únicos verificados biometricamente',
  'pools-title':'Pools de Redistribuição',
  'pools-desc':'Cada taxa de swap, demurrage e excesso do teto é dividido entre quatro pools. Todos pagam diariamente.',
  'vel-pool':'Pool de Validadores','vel-pool-desc':'40% de todas as taxas → operadores de nodes que protegem a rede',
  'liq-pool':'Pool de Liquidez','liq-pool-desc':'30% de todas as taxas → provedores de liquidez, proporcional às cotas LP',
  'ubi-pool':'Pool UBI','ubi-pool-desc':'20% de todas as taxas → todos os humanos verificados igualmente, a cada 24 horas',
  'treasury':'Tesouro','treasury-desc':'10% de todas as taxas → desenvolvimento e manutenção do protocolo',
  'phases-title':'Fases do Protocolo',
  'phases-desc':'Teto bootstrap Fase 0: max(5, min(N, 25))× saldo médio. 1–4 humanos: 5×. Cada humano +1×. 25+ humanos: travado em 25×. Transições automáticas.',
  'p0':'Bootstrap · &lt;100 humanos · Teto: max(5,min(N,25))× médio · 5×→25× · Ativo agora',
  'p1':'Crescimento · 100–10.000 humanos · Teto: 25× saldo médio',
  'p2':'Estabilidade · 10.000–1M humanos · Teto: 25× saldo médio',
  'p3':'Maturidade · 1M+ humanos · Teto: 25× saldo médio',
  'wealth-cap-explain':'Teto Fase 0: max(5, min(N, 25))× saldo médio. 1–4 humanos: 5×. Cada humano +1×. 25+: travado em 25×.',
  'demurrage-title':'Demurrage — Incentivo para Circular',
  'demurrage-desc':'Saldos AEQ inativos perdem lentamente valor para desencorajar acumulação.',
  'dem-rate-k':'Taxa de Decaimento','dem-rate-v':'0,5% por mês (contínuo)',
  'dem-grace-k':'Período de Graça','dem-grace-v':'3 meses de inatividade antes do decaimento começar',
  'dem-reset-k':'Reinicialização','dem-reset-v':'Qualquer transferência, swap ou liquidez reinicia o contador',
  'dem-dest-k':'AEQ decaído vai para','dem-dest-v':'Pools de redistribuição (40/30/20/10)',
  'dem-warn-k':'Sistema de Aviso','dem-warn-v':'Aviso 14 dias (uma vez) + lembrete 7 dias repetido em cada login',
  'story-title':'A História da Aequitas',
  'story-text':'<p>Em 2009, Satoshi Nakamoto lança o Bitcoin. Revolução genuína — mas os primeiros mineradores acumulam milhões a custo quase zero. Em 2021, top 1% controla mais de 90% do Bitcoin. Gini Bitcoin &gt; 0,85.</p><p><span style="color:var(--gold)">Aequitas</span> — latim para "equidade" — responde: <em style="color:var(--gold)">"Como seria uma criptomoeda projetada para ser justa com cada ser humano?"</em></p><p><strong style="color:var(--text)">O dinheiro existe porque as pessoas existem. Portanto, cada pessoa deveria ter uma parte igual.</strong></p><p><em style="color:var(--gold)">"O dinheiro existe porque as pessoas existem. Nada mais, nada menos."</em></p>',
  'nodes-title':'Nodes Ativos — Topologia de Rede Atual',
  'nodes-desc':'A rede Aequitas opera em dois nodes distribuídos geograficamente, participando da produção de blocos, sincronização e API. Nodes adicionais são bem-vindos.',
  'node1':'Node 1 — Railway (Principal)','node1-desc':'API principal · Produtor de blocos · Distribuição UBI · Bootstrap P2P · PostgreSQL · RPC MetaMask',
  'node2':'Node 2 — Render (Secundário)','node2-desc':'API secundário · Produtor de blocos · Par P2P · Sync HTTP · Estado PostgreSQL compartilhado',
  'run-node-title':'Execute seu Próprio Node','run-node-desc':'Qualquer um pode executar um node Aequitas — sem permissão, sem stake. Operadores ganham 40% das taxas de swap distribuídas diariamente.',
  'bootstrap-title':'Conectar um Novo Node','bootstrap-desc':'Defina PRIMARY_NODE_URL=https://aequitas.digital no seu ambiente. Seu node sincroniza automaticamente o estado completo da cadeia.',
  'tech-title':'Especificações Técnicas','mm-config':'Configuração MetaMask',
  'k-lang':'Idioma','k-src':'Fonte','evm-yes':'Sim — JSON-RPC /rpc · Compatível MetaMask',
  'proto-label':'Protocolo Aequitas V7 — Documentação Técnica',
  'ca-title':'Endereços dos Contratos',
  'ca-text':'Cadeia: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7: 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7 é a única fonte de verdade para toda a economia Aequitas. Nenhuma chave de administrador nem voto de governança pode alterar sua lógica. O código atual rodará em dez anos.',
  'poa-title':'1. PROVA DE VIDA','poa-text':'<p>AEQ de pessoas falecidas retorna gradualmente à comunidade via pool UBI, em vez de ser perdido para sempre como no Bitcoin.</p>',
  'poa-box':'Anos 0–2: Uso normal<br>Ano 2: Aviso 1 — Guardião pode responder<br>Ano 2+60d: Aviso 2<br>Ano 2+120d: Aviso 3<br>Ano 2+180d: AEQ em custódia pessoal<br>Ano 4: Se inativo — retorna ao Pool UBI',
  'guard-title':'2. SISTEMA DE GUARDIÃO','guard-text':'<p>Um Guardião de confiança (outro humano verificado) pode confirmar que alguém está vivo, sem nenhum direito de transação.</p>',
  'guard-box':'1 Guardião por humano · deve ser humano verificado Aequitas<br>Guardião pode APENAS chamar confirmAlive() · zero direitos financeiros<br>Guardião NÃO PODE mover fundos · Máx. 3 protegidos · Timelock 7d',
  'dem-title':'3. DEMURRAGE — Anti-Acumulação',
  'dem-box':'Taxa: 0,5%/mês após 3 meses de graça<br>Reinicialização a cada transferência, swap ou liquidez<br>AEQ decaído redistribuído nos pools (não queimado)',
  'dem-text':'<p>Precedente: Wörgl (Áustria, 1932) — desemprego reduziu 25% em um ano. Chiemgauer (Alemanha, 2003) — opera com sucesso há mais de 20 anos.</p>',
  'cap-title':'4. TETO DE RIQUEZA','cap-box':'Teto: max(5,min(N,25))× saldo médio AEQ<br>1–4 humanos: 5× · +1× por humano · 25+: 25× permanente<br>Excesso redistribuído imediatamente · Sem intervenção manual',
  'ubi-title':'5. RENDA BÁSICA UNIVERSAL','ubi-box':'Fontes: Taxas de swap (20%) · Excesso do teto · Demurrage<br><br>Diário: Pool UBI dividido igualmente entre todos os humanos. Pool zera após cada distribuição.',
  'inf-title':'6. SEM INFLAÇÃO ALGORÍTMICA','inf-box':'Único evento criando AEQ: novo humano verificado se registra.<br><br>Oferta Total = Humanos Verificados × 1.000 AEQ — sempre, exatamente.',
  'btn-download-app':'BAIXAR AEQUITASBIO',
  'swap-title':'🔄 Trocar AEQ ↔ tUSD','swap-sub':'Troque AEQ por tUSD (dólar de teste) via pool de liquidez nativo. Taxa 0,1% apenas para swaps — transferências AEQ comuns completamente gratuitas.',
  'swap-priv-bar':'🔒 Apenas 0,1% de taxa · Transferências AEQ→AEQ gratuitas · tUSD é moeda de teste sem valor real',
  'swap-your-aeq':'Seu AEQ','swap-your-tusd':'Seu tUSD','swap-aeq-to-tusd':'AEQ → tUSD','swap-tusd-to-aeq':'tUSD → AEQ',
  'swap-fee-est':'Taxa de protocolo (0,1%)','swap-details-hdr':'Detalhes da Troca',
  'swap-out-lbl':'Você recebe (est.)','swap-impact-lbl':'Impacto no preço','swap-rate-lbl':'Taxa de câmbio',
  'swap-depth-lbl':'Composição do Pool','amm-title':'x × y = k — AMM de Produto Constante',
  'amm-text':'No swap, reservas AEQ aumentam e reservas tUSD diminuem — produto sempre igual a k. Swaps maiores causam maior impacto no preço.',
  'swap-btn-conn':'🦊 CONECTAR METAMASK','swap-btn-go':'🔄 TROCAR',
  'swap-log-hint':'// Conectar carteira para trocar...',
  'swap-no-liquidity':'Ainda sem tUSD?','swap-faucet-desc':'Humanos registrados podem reivindicar tUSD de teste uma vez','swap-btn-faucet':'💧 REIVINDICAR tUSD TESTE',
  'swap-addliq-title':'Fornecer Liquidez','swap-addliq-desc':'Seja o primeiro a depositar — sua proporção define o preço inicial.','swap-btn-addliq':'💧 ADICIONAR LIQUIDEZ',
  'swap-lp-title':'Sua Posição LP','swap-lp-share':'Cota do Pool','swap-lp-withdrawable':'Retirável',
  'swap-lp-pct-label':'% da sua posição','swap-lp-youget':'Você receberá','swap-btn-removeliq':'🔥 REMOVER LIQUIDEZ',
  'swap-pool-title':'AEQ / tUSD — Status do Pool',
  'swap-pool-aeq':'Reserva AEQ','swap-pool-tusd':'Reserva tUSD','swap-pool-price':'Preço Spot',
  'swap-fee-bps':'Taxa de Swap','swap-fee-split':'Distribuição de taxas','swap-fee-split-v':'40% Validadores / 30% LP / 20% UBI / 10% Tesouro',
  'swap-pools-addr-title':'Endereços dos Pools Tokenômicos',
  'swap-validators':'Validadores (40%)','swap-lps':'Provedores de Liquidez (30%)','swap-ubi':'Pool UBI (20%)','swap-treasury':'Tesouro (10%)',
  'ubi-hero-title':'RENDA BÁSICA UNIVERSAL — POOL UBI',
  'ubi-hero-sub':'Acumulando — próximo pagamento distribuído a todos os humanos verificados em:',
  'ubi-bal-lbl':'saldo atual do pool','ubi-hero-desc':'Dividido igualmente · pago a cada 24h · pool zerado · saldo mínimo não necessário',
  'ubi-how-fills':'Como o Pool UBI se enche',
  'ubi-src-swap':'Taxas de Swap','ubi-src-swap-d':'Cada swap AEQ↔tUSD contribui 20% de suas taxas. Mais trading = enchimento mais rápido.',
  'ubi-src-dem':'Demurrage','ubi-src-dem-d':'AEQ inativo (3+ meses) decai 0,5%/mês. 20% do decaído vai para UBI.',
  'ubi-src-cap':'Excesso do Teto','ubi-src-cap-d':'Carteiras que excedem o teto são imediatamente cortadas. 20% flui para UBI.',
  'pools4-header':'Os quatro pools de redistribuição',
  'ubi-see-above':'ver contagem regressiva acima','ubi-timer-above':'⏰ contagem regressiva exibida acima','pool-t-timer':'Acumulando — sem temporizador',
  'usp-headline':'Pela primeira vez na história — todos começam em igualdade',
  'usp-sub':'Com um smartphone Android você é elegível. Sem banco, sem crypto, sem investimento.',
  'usp-c1-title':'R$ 0,00 de Investimento Inicial','usp-c1-desc':'Registro completamente sem taxas. Sem ETH, MATIC ou cartão. O protocolo paga todos os custos.',
  'usp-c2-title':'1.000 AEQ para cada humano','usp-c2-desc':'Bilionário ou agricultor — todos recebem exatamente 1.000 AEQ. Igualdade garantida matematicamente.',
  'usp-c3-title':'Acessível a todos','usp-c3-desc':'Sem conta bancária, cartão de crédito ou documento. O registro usa um kit biométrico acessível (leitor de impressão digital + sensor de pulso, ~$15) — projetado para acesso global.',
  'usp-c4-title':'UBI diário para sempre','usp-c4-desc':'Após registrado, sua parte do UBI chega automaticamente todos os dias — sem nenhuma ação.',
  'v7-intro-title':'O que é AequitasV7?',
  'v7-intro-text':'AequitasV7 é o contrato inteligente central do protocolo Aequitas. Implantado de forma imutável na Aequitas Chain (ID 1926). Gerencia tudo: registro humano, verificação ZK, saldos, teto de riqueza, UBI, taxas de swap. Nenhum administrador pode modificá-lo.',
  'explore-title':'Explorar Aequitas',
  'expl-score':'Pontuação de Igualdade','expl-score-d':'Coeficiente de Gini ao vivo · Índice Aequitas · distribuição de riqueza em tempo real',
  'expl-economy':'UBI &amp; Redistribuição','expl-economy-d':'Contagem regressiva UBI · 4 pools on-chain · demurrage · Fases do Protocolo',
  'expl-charts':'Gráficos &amp; Histórico','expl-charts-d':'Histórico Gini · curva de Lorenz · controle do teto · A história da Aequitas',
  'expl-v7':'Docs Protocolo V7','expl-v7-d':'Contrato AequitasV7 · 6 mecanismos · prova ZK · teto · demurrage · código imutável',
  'expl-explorer':'Explorador de Blocos','expl-explorer-d':'BlockDAG ao vivo · clique em qualquer bloco para ver validador, hash, transações',
  'swap-sell-label':'Vender','swap-receive-label':'Receber',
  'expl-network':'Rede &amp; Nodes','expl-network-d':'Topologia de nodes · executar seu próprio node · especificações · Chain ID 1926'
},
ar:{
  'logo-sub':'إثبات الإنسانية','live':'مباشر',
  'tab-register':'🔐 تسجيل','tab-explorer':'🔍 المستكشف','tab-humans':'👥 البشر','tab-index':'📊 المؤشر','tab-network':'🌐 الشبكة','tab-protocol':'📜 البروتوكول V7','tab-swap':'🔄 تبادل',
  'reg-title':'🔐 التسجيل كإنسان موثق',
  'reg-sub':'انضم إلى شبكة Aequitas واحصل على منحة دخل أساسي شامل تبلغ 1,000 AEQ. التسجيل لمرة واحدة، دائم، ومجاني تماماً. لا يتم تخزين أي بيانات شخصية.',
  'app-title':'التسجيل عبر تطبيق أندرويد',
  'app-text':'يستخدم إثبات الإنسانية نظاماً بيومترياً مادياً ثلاثي العوامل. المرحلة 1: مستشعر R503 البصري يمسح جميع بصمات الأصابع العشر + MAX30102 PPG يؤكد النبض الحي. المرحلة 2: كاميرا IR لأوردة اليد (تفرد 1 من 10⁷). المرحلة 3: مسح القزحية — المعيار الذهبي، 1 من 10⁷⁸، مستقل كلياً عن الجهاز. دليل ZK من نوع Groth16 يلتزم بجميع العوامل دون الكشف عن أي بيانات بيومترية. يُضاف 1,000 AEQ تلقائياً بعد التحقق.',
  's1t':'المسح البيومتري','s1d':'AequitasBio يمسح جميع البصمات العشر (مستشعر R503 البصري) + نبض MAX30102 PPG يؤكد الحيوية. المرحلة 2: IR أوردة اليد. المرحلة 3: القزحية. البيانات الخام لا تغادر الجهاز أبداً.',
  's2t':'توليد دليل ZK','s2d':'دليل ZK من نوع Groth16 يلتزم بجميع العوامل البيومترية: commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier مرتبط بالجسم وليس بالهاتف — فقدان الجهاز لا يُنشئ هوية ثانية.',
  's3t':'ربط المحفظة','s3d':'يفتح التطبيق MetaMask · ارتبط بمحفظة Ethereum · الدليل مرتبط تشفيرياً بعنوان محفظتك',
  's4t':'تم منح 1,000 AEQ','s4d':'تم تأكيد التسجيل على BlockDAG خلال 6 ثوانٍ · اعتماد 1,000 AEQ فوراً · هويتك مسجلة بشكل دائم',
  'priv-bar':'🔒 R503 جميع البصمات العشر · MAX30102 حيوية · المرحلة 2: IR أوردة اليد · المرحلة 3: القزحية (10⁷⁸) · Groth16 ZKP · البيانات لا تغادر الجهاز · إنسان واحد · إلى الأبد',
  'conn-wallet':'المحفظة المتصلة','proof-recv':'⚡ تم استلام دليل ZK','proof-hint':'ربط محفظة للتسجيل',
  'btn-conn':'🦊 ربط METAMASK','btn-reg':'🔐 التسجيل ON-CHAIN',
  'btn-web-reg':'🌐 التسجيل عبر المتصفح (WebAuthn)',
  'web-reg-warn':'⚠ مرتبط بالجهاز: هذه الهوية مرتبطة بهذا الجهاز والمتصفح. لا يمكن نقلها. للهوية متعددة الأجهزة، استخدم تطبيق Aequitas Android.',
  'reg-log-hint':'// افتح تطبيق Aequitas Android لتوليد دليلك، ثم عد هنا...',
  'reg-details':'تفاصيل التسجيل','k-network':'الشبكة','k-chainid':'معرّف السلسلة','k-grant':'منحة UBI',
  'k-fee':'رسوم الغاز','free':'مجاني — بدون رسوم تماماً','k-limit':'التسجيلات','k-limit-v':'مرة واحدة لكل إنسان · دائم · غير قابل للتغيير',
  'k-bio':'البيانات البيومترية','never-stored':'لا تُخزَّن أبداً — تبقى على جهازك',
  'k-proof':'نظام الأدلة','k-conf':'التأكيد','k-conf-v':'خلال 6 ثوانٍ (كتلة واحدة)',
  'k-sybil':'حماية Sybil','k-sybil-v':'هوية واحدة لكل بيومتري · قفل دائم',
  'live-stats':'إحصائيات السلسلة المباشرة',
  's-height':'ارتفاع الكتلة','s-height-sub':'كتلة جديدة كل ~6 ث · BlockDAG · إنتاج متوازٍ',
  's-humans':'البشر الموثقون','s-humans-sub':'ZKP بيومتري · شخص واحد، محفظة واحدة، إلى الأبد',
  's-supply':'إجمالي العرض','s-supply-sub':'دائماً = البشر × 1,000 AEQ',
  's-index':'مؤشر Aequitas','s-index-sub':'0 = مساواة مثالية · 100 = أقصى عدم مساواة',
  's-uptime':'وقت التشغيل','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'إثبات الإنسانية','ib-poh-t':'يجب على كل حامل AEQ إثبات أنه إنسان حي فريد. لا بوتات ولا شركات ولا ذكاء اصطناعي. البيانات البيومترية لا تغادر جهازك.',
  'ib-fair':'توزيع عادل جذرياً','ib-fair-t':'كل إنسان موثق يحصل على 1,000 AEQ بالضبط عند التسجيل. لا تعدين مسبق. الإجمالي = البشر × 1,000.',
  'ib-dag':'بنية BlockDAG','ib-dag-t':'يمكن إنتاج كتل متعددة في وقت واحد ودمجها. إنتاجية أعلى وزمن استجابة أقل.',
  'ib-gas':'مجاني حقاً','ib-gas-t':'التسجيل وتحويلات AEQ لا تكلف شيئاً. لا حاجة لـ ETH أو BNB أو MATIC أو حساب بنكي.',
  'recent-blocks':'الكتل الأخيرة','blocks-desc':'MERGE = دمج عدة والدين (BlockDAG). TX = معاملة تسجيل. وقت الكتلة: ~6 ثوانٍ.',
  'loading':'جارٍ تحميل الكتل...','net-info':'معلومات الشبكة','k-chain':'اسم السلسلة','k-symbol':'الرمز','k-btime':'وقت الكتلة',
  'k-cons':'التوافق','k-nodes':'العقد النشطة','k-storage':'التخزين','add-mm':'🦊 إضافة إلى METAMASK','k-dec':'الأرقام العشرية',
  'btn-add-mm':'+ إضافة شبكة AEQUITAS',
  'phil':'"المال موجود لأن البشر موجودون.<br>لا أكثر، ولا أقل."','phil-sub':'— مبدأ AEQUITAS —',
  'humans-title':'البشر الموثقون على Aequitas Chain',
  'h-what':'ما هو الإنسان الموثق؟','h-what-t':'الإنسان الموثق هو عنوان محفظة مُثبت تشفيرياً كأنه ينتمي لإنسان حي فريد. يستخدم التحقق نظام أجهزة ثلاثي العوامل: R503 يمسح جميع البصمات العشر؛ MAX30102 PPG يؤكد النبض الحي؛ المرحلة 2: IR أوردة اليد (1 من 10⁷)؛ المرحلة 3: القزحية (1 من 10⁷⁸). يُرسَل دليل Groth16 ZK فقط. لا تغادر أي بيانات بيومترية الجهاز.',
  'h-zkp':'نظام أدلة ZK','h-zkp-t':'Aequitas يستخدم Groth16 على BN128 — نفس المنحنى المستخدم في Ethereum وZcash. ~200 بايت، ~10ms. commitment = keccak256(iris‖vein‖fingers‖wallet). Nullifier مرتبط بالجسم: فقدان الهاتف لا يُنشئ هوية ثانية. لا تُخزَّن أي بيانات بيومترية.',
  'h-sybil':'منع هجمات Sybil','h-sybil-t':'المرحلة 1: جميع البصمات العشر + حيوية MAX30102 (نبض PPG، يرفض القوالب/إعادة التشغيل). المرحلة 2: IR أوردة اليد — ميزة داخلية في الجسم، مستحيلة النسخ، 1 من 10⁷، مختلفة في التوائم المتطابقة. المرحلة 3: القزحية — 1 من 10⁷⁸، المعيار الذهبي العالمي. Nullifier = keccak256(iris‖vein‖domain). إنسان واحد، محفظة واحدة، إلى الأبد.',
  'h-global':'الشمول المالي العالمي','h-global-t':'لا حاجة لحساب بنكي أو بطاقة ائتمان أو عملة مشفرة. هاتف أندرويد بمستشعر بيومتري يكفي.',
  'h-bio-hw':'خارطة طريق الأجهزة البيومترية','h-bio-hw-t':'المرحلة 1 (نشطة): ماسح بصمات بصري R503 — هاش مجمّع لجميع الأصابع العشر. حيوية MAX30102 PPG. المرحلة 2 (مخططة): ESP32-CAM + LED IR بطول موجي 850 nm — تصوير أوردة اليد، تفرد 1 من 10⁷. المرحلة 3 (مخططة): وحدة قزحية IR — أكثر من 240 درجة حرية، 1 من 10⁷⁸، مستقلة كلياً عن الجهاز، التوائم المتطابقة تختلف.',
  'reg-humans':'البشر المسجلون','h-desc':'كل عنوان تم التحقق منه كإنسان فريد. كل واحد حصل على 1,000 AEQ بالضبط. دائم وغير قابل للتغيير.',
  'no-humans':'لا يوجد بشر مسجلون بعد.\n\nحمّل تطبيق Aequitas Android وكن أول إنسان على السلسلة!',
  'reg-stats':'إحصائيات السجل','total-humans':'إجمالي البشر',
  'idx-title':'مؤشر Aequitas — درجة المساواة الاقتصادية في الوقت الفعلي',
  'idx-desc':'مؤشر Aequitas مشتق من <strong style="color:var(--teal)">معامل جيني</strong> — المعيار الدولي لقياس عدم المساواة (البنك الدولي، OECD، الأمم المتحدة). <strong style="color:var(--neon)">0 = مساواة مثالية</strong>. <strong style="color:var(--red)">100 = تركيز كامل</strong>. الهدف: جيني أقل من 0.30.',
  'gini-what-title':'ما هو معامل جيني؟',
  'gini-what-text':'طوّره كورادو جيني (1912). يقيس توزيع الثروة. المقياس: 0 (الجميع متساوون) إلى 1 (شخص واحد يملك كل شيء). يُستخدم من قِبل البنك الدولي وOECD والأمم المتحدة.',
  'curr-idx':'المؤشر الحالي','bar-0':'0 — مساواة مثالية','bar-100':'100 — أقصى عدم مساواة','wcap-lbl':'سقف الثروة الحالي:','wcap-mult':'المضاعف:','wcap-avg':'متوسط الرصيد:',
  'gini':'معامل جيني','gini-desc':'0 = متساوٍ · 1 = غير متساوٍ',
  'supply-desc':'دائماً = البشر × 1,000 AEQ',
  'phase':'مرحلة البروتوكول','phase-desc':'يتقدم تلقائياً بعدد البشر',
  'humans-desc':'بشر فريدون موثقون بيومترياً',
  'pools-title':'مجمعات إعادة التوزيع',
  'pools-desc':'كل رسوم المبادلة والتلاشي والفائض من سقف الثروة يُقسَّم تلقائياً بين أربعة مجمعات. جميعها تدفع يومياً.',
  'vel-pool':'مجمع المدققين','vel-pool-desc':'40% من جميع الرسوم ← مشغّلو العقد الذين يؤمّنون الشبكة',
  'liq-pool':'مجمع السيولة','liq-pool-desc':'30% من جميع الرسوم ← مزودو السيولة، بنسبة حصص LP',
  'ubi-pool':'مجمع UBI','ubi-pool-desc':'20% من جميع الرسوم ← جميع البشر الموثقين بالتساوي، كل 24 ساعة',
  'treasury':'الخزينة','treasury-desc':'10% من جميع الرسوم ← تطوير البروتوكول وصيانته',
  'phases-title':'مراحل البروتوكول',
  'demurrage-title':'التلاشي — حافز للتداول',
  'demurrage-desc':'أرصدة AEQ غير النشطة تفقد قيمتها ببطء لثني الاكتناز وتحفيز المشاركة الاقتصادية.',
  'dem-rate-k':'معدل التلاشي','dem-rate-v':'0.5% شهرياً (مستمر)',
  'dem-grace-k':'فترة السماح','dem-grace-v':'3 أشهر من الخمول قبل بدء التلاشي',
  'dem-reset-k':'إعادة التعيين','dem-reset-v':'أي تحويل أو مبادلة أو إجراء سيولة يعيد العداد إلى الصفر',
  'dem-dest-k':'AEQ المتلاشي يذهب إلى','dem-dest-v':'مجمعات إعادة التوزيع (40/30/20/10)',
  'dem-warn-k':'نظام التحذير','dem-warn-v':'إشعار 14 يوماً (مرة واحدة) + تذكير 7 أيام عند كل تسجيل دخول',
  'story-title':'قصة Aequitas',
  'story-text':'<p>عام 2009، أصدر ساتوشي ناكاموتو Bitcoin. ثورة حقيقية — لكن المنقبين الأوائل جمعوا الملايين بتكلفة شبه معدومة. في 2021، يتحكم أعلى 1% في أكثر من 90% من Bitcoin. جيني Bitcoin &gt; 0.85.</p><p><span style="color:var(--gold)">Aequitas</span> — لاتينية لـ "العدالة" — أُنشئ للإجابة على: <em style="color:var(--gold)">"كيف ستبدو عملة مشفرة صُمِّمت لتكون عادلة لكل إنسان؟"</em></p><p><strong style="color:var(--text)">المال موجود لأن البشر موجودون. لذا يجب أن يحصل كل شخص على حصة متساوية.</strong></p><p><em style="color:var(--gold)">"المال موجود لأن البشر موجودون. لا أكثر، ولا أقل."</em></p>',
  'nodes-title':'العقد النشطة — طوبولوجيا الشبكة الحالية',
  'nodes-desc':'تعمل شبكة Aequitas على عقدتين موزعتين جغرافياً، تشاركان في إنتاج الكتل والمزامنة وخدمة API.',
  'node1':'العقدة 1 — Railway (الأساسية)','node1-desc':'API أساسي · منتج كتل · توزيع UBI · P2P Bootstrap · PostgreSQL · RPC لـ MetaMask',
  'node2':'العقدة 2 — Render (الثانوية)','node2-desc':'API ثانوي · منتج كتل · نظير P2P · مزامنة HTTP · حالة PostgreSQL مشتركة',
  'run-node-title':'قم بتشغيل عقدتك الخاصة','run-node-desc':'يمكن لأي شخص تشغيل عقدة Aequitas — بدون إذن أو حصة. المشغّلون يكسبون 40% من رسوم المبادلة يومياً.',
  'bootstrap-title':'ربط عقدة جديدة','bootstrap-desc':'اضبط PRIMARY_NODE_URL=https://aequitas.digital في بيئتك. عقدتك ستزامن حالة السلسلة الكاملة تلقائياً.',
  'tech-title':'المواصفات التقنية','mm-config':'إعداد MetaMask',
  'k-lang':'اللغة','k-src':'المصدر','evm-yes':'نعم — JSON-RPC /rpc · متوافق مع MetaMask',
  'proto-label':'بروتوكول Aequitas V7 — وثائق تقنية',
  'ca-title':'عناوين العقود',
  'ca-text':'السلسلة: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7: 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7 هو المصدر الوحيد للحقيقة لاقتصاد Aequitas بأكمله. لا مفتاح إدارة ولا تصويت حوكمة يمكنه تغيير منطقه.',
  'poa-title':'1. إثبات الحياة','poa-text':'<p>عند وفاة الأشخاص، تعود AEQ الخاصة بهم تدريجياً إلى المجتمع عبر مجمع UBI بدلاً من ضياعها للأبد.</p>',
  'poa-box':'السنوات 0–2: استخدام طبيعي<br>السنة 2: تحذير 1 — الحارس يمكنه الرد<br>السنة 2+60 يوم: تحذير 2<br>السنة 2+120 يوم: تحذير 3<br>السنة 2+180 يوم: AEQ في ضمان شخصي<br>السنة 4: إذا لا يزال خاملاً — يعود لمجمع UBI',
  'guard-title':'2. نظام الحارس','guard-text':'<p>حارس موثوق (إنسان موثق آخر) يمكنه تأكيد أن شخصاً ما لا يزال حياً، دون أي حقوق مالية.</p>',
  'guard-box':'حارس واحد لكل إنسان · يجب أن يكون إنساناً موثقاً<br>الحارس يمكنه فقط استدعاء confirmAlive() · صفر حقوق مالية<br>الحارس لا يمكنه تحريك الأموال · الحد الأقصى 3 · Timelock 7 أيام',
  'dem-title':'3. التلاشي — آلية مكافحة الاكتناز',
  'dem-box':'المعدل: 0.5%/شهر بعد 3 أشهر سماح<br>إعادة تعيين عند أي تحويل أو مبادلة أو سيولة<br>AEQ المتلاشي يُعاد توزيعه في المجمعات (لا يُحرق)',
  'dem-text':'<p>سابقة تاريخية: تجربة Wörgl (النمسا، 1932) — خفض البطالة 25% في عام واحد. Chiemgauer (ألمانيا، 2003) — يعمل بنجاح منذ أكثر من 20 عاماً.</p>',
  'cap-title':'4. سقف الثروة','cap-box':'السقف: max(5,min(N,25))× متوسط الرصيد<br>1–4 بشر: 5× · +1× لكل إنسان · 25+: 25× دائم<br>الفائض يُعاد توزيعه فوراً · بدون تدخل يدوي',
  'ubi-title':'5. الدخل الأساسي الشامل','ubi-box':'المصادر: رسوم المبادلة (20%) · فائض السقف · التلاشي<br><br>يومياً: مجمع UBI مقسّم بالتساوي بين جميع البشر المسجلين. يُعاد ضبط المجمع بعد كل توزيع.',
  'inf-title':'6. لا تضخم خوارزمي','inf-box':'الحدث الوحيد الذي ينشئ AEQ جديداً: تسجيل إنسان موثق جديد.<br><br>إجمالي العرض = البشر الموثقون × 1,000 AEQ — دائماً، بالضبط.',
  'btn-download-app':'تحميل تطبيق AEQUITASBIO',
  'swap-title':'🔄 تبادل AEQ ↔ tUSD','swap-sub':'تبادل AEQ مع tUSD (دولار اختبار محاكى) عبر مجمع السيولة الأصلي. رسوم 0.1% فقط للمبادلات — التحويلات العادية مجانية تماماً.',
  'swap-priv-bar':'🔒 رسوم 0.1% فقط · تحويلات AEQ→AEQ مجانية · tUSD عملة اختبار بدون قيمة حقيقية',
  'swap-your-aeq':'AEQ لديك','swap-your-tusd':'tUSD لديك','swap-aeq-to-tusd':'AEQ → tUSD','swap-tusd-to-aeq':'tUSD → AEQ',
  'swap-fee-est':'رسوم البروتوكول (0.1%)','swap-details-hdr':'تفاصيل التبادل',
  'swap-out-lbl':'ستحصل على (تقريباً)','swap-impact-lbl':'تأثير السعر','swap-rate-lbl':'سعر الصرف',
  'swap-depth-lbl':'تكوين المجمع','amm-title':'x × y = k — AMM ذو الجداء الثابت',
  'amm-text':'عند التبادل، تزداد احتياطيات AEQ وتنخفض احتياطيات tUSD — جداؤها يبقى دائماً مساوياً لـ k. التبادلات الكبيرة تسبب تأثيراً أكبر على السعر.',
  'swap-btn-conn':'🦊 ربط METAMASK','swap-btn-go':'🔄 تبادل',
  'swap-log-hint':'// ربط محفظة للتبادل...',
  'swap-no-liquidity':'لا يوجد tUSD بعد?','swap-faucet-desc':'البشر المسجلون يمكنهم المطالبة بـ tUSD اختبار مرة واحدة','swap-btn-faucet':'💧 المطالبة بـ tUSD الاختبار',
  'swap-addliq-title':'توفير السيولة','swap-addliq-desc':'كن أول من يودع — نسبتك تحدد السعر الأولي.','swap-btn-addliq':'💧 إضافة سيولة',
  'swap-lp-title':'مركز LP الخاص بك','swap-lp-share':'حصة المجمع','swap-lp-withdrawable':'قابل للسحب',
  'swap-lp-pct-label':'% من مركزك','swap-lp-youget':'ستحصل على','swap-btn-removeliq':'🔥 سحب السيولة',
  'swap-pool-title':'AEQ / tUSD — حالة المجمع',
  'swap-pool-aeq':'احتياطي AEQ','swap-pool-tusd':'احتياطي tUSD','swap-pool-price':'السعر الفوري',
  'swap-fee-bps':'رسوم المبادلة','swap-fee-split':'توزيع الرسوم','swap-fee-split-v':'40% مدققون / 30% LP / 20% UBI / 10% خزينة',
  'swap-pools-addr-title':'عناوين مجمعات التوكينوميكس',
  'swap-validators':'المدققون (40%)','swap-lps':'مزودو السيولة (30%)','swap-ubi':'مجمع UBI (20%)','swap-treasury':'الخزينة (10%)',
  'ubi-hero-title':'الدخل الأساسي الشامل — مجمع UBI',
  'ubi-hero-sub':'يتراكم — الدفعة التالية توزَّع بالتساوي على جميع البشر الموثقين خلال:',
  'ubi-bal-lbl':'رصيد المجمع الحالي','ubi-hero-desc':'مقسَّم بالتساوي · يُدفع كل 24 ساعة · يُصفَّر المجمع · لا يشترط رصيد أدنى',
  'ubi-how-fills':'كيف يمتلئ مجمع UBI',
  'ubi-src-swap':'رسوم المبادلة','ubi-src-swap-d':'كل مبادلة AEQ↔tUSD تساهم بـ 20% من رسومها. المزيد من التداول = امتلاء أسرع.',
  'ubi-src-dem':'التلاشي','ubi-src-dem-d':'AEQ الخامل (3+ أشهر) يتلاشى 0.5%/شهر. 20% من المتلاشي يذهب لـ UBI.',
  'ubi-src-cap':'فائض السقف','ubi-src-cap-d':'المحافظ التي تتجاوز السقف تُقلَّص فوراً. 20% يتدفق إلى UBI.',
  'pools4-header':'المجمعات الأربعة لإعادة التوزيع',
  'ubi-see-above':'انظر العد التنازلي أعلاه','ubi-timer-above':'⏰ العد التنازلي معروض أعلاه','pool-t-timer':'يتراكم — لا عداد',
  'usp-headline':'لأول مرة في التاريخ — الجميع يبدأ على قدم المساواة',
  'usp-sub':'إذا كان لديك هاتف أندرويد فأنت مؤهل. بدون بنك، بدون معرفة بالعملات المشفرة، بدون استثمار.',
  'usp-c1-title':'استثمار أولي 0','usp-c1-desc':'التسجيل مجاني تماماً. لا ETH ولا بطاقة بنكية. البروتوكول يدفع جميع رسوم المعاملات.',
  'usp-c2-title':'1,000 AEQ لكل إنسان','usp-c2-desc':'مليارديرًا كان أم مزارعاً — الجميع يحصل على 1,000 AEQ بالضبط. مساواة مضمونة رياضياً.',
  'usp-c3-title':'متاح للجميع','usp-c3-desc':'لا حاجة لحساب بنكي أو بطاقة ائتمان أو وثيقة هوية. يستخدم التسجيل مجموعة أجهزة بيومترية بأسعار معقولة (ماسح بصمات + مستشعر نبض، ~15 دولار) — مصمم للوصول العالمي.',
  'usp-c4-title':'UBI يومي إلى الأبد','usp-c4-desc':'بعد التسجيل، تصل حصتك من UBI تلقائياً كل يوم — دون أي إجراء.',
  'v7-intro-title':'ما هو AequitasV7؟',
  'v7-intro-text':'AequitasV7 هو العقد الذكي المركزي لبروتوكول Aequitas. مُنشر بشكل غير قابل للتغيير على Aequitas Chain (ID 1926). يدير كل شيء: التسجيل البشري، التحقق ZK، الأرصدة، سقف الثروة، UBI، رسوم المبادلة. لا يمكن لأي مدير تعديله.',
  'explore-title':'استكشف Aequitas',
  'expl-score':'درجة المساواة','expl-score-d':'معامل جيني مباشر · مؤشر Aequitas · توزيع الثروة في الوقت الفعلي',
  'expl-economy':'UBI وإعادة التوزيع','expl-economy-d':'عد UBI التنازلي اليومي · 4 مجمعات on-chain · تلاشي · مراحل البروتوكول',
  'expl-charts':'الرسوم البيانية والتاريخ','expl-charts-d':'تاريخ جيني · منحنى لورينز · شريط سقف الثروة · قصة Aequitas',
  'expl-v7':'وثائق البروتوكول V7','expl-v7-d':'عقد AequitasV7 · 6 آليات · دليل ZK · سقف الثروة · تلاشي · كود غير قابل للتغيير',
  'expl-explorer':'مستكشف الكتل','expl-explorer-d':'BlockDAG مباشر · انقر على أي كتلة لرؤية المدقق والهاش والمعاملات',
  'swap-sell-label':'بيع','swap-receive-label':'استلام',
  'expl-network':'الشبكة والعقد','expl-network-d':'طوبولوجيا العقد · تشغيل عقدتك الخاصة · المواصفات التقنية · Chain ID 1926'
},
hi:{
  'logo-sub':'मानवता का प्रमाण','live':'लाइव',
  'tab-register':'🔐 रजिस्टर','tab-explorer':'🔍 एक्सप्लोरर','tab-humans':'👥 मनुष्य','tab-index':'📊 इंडेक्स','tab-network':'🌐 नेटवर्क','tab-protocol':'📜 प्रोटोकॉल V7','tab-swap':'🔄 स्वैप',
  'reg-title':'🔐 सत्यापित मानव के रूप में रजिस्टर करें',
  'reg-sub':'Aequitas नेटवर्क से जुड़ें और 1,000 AEQ का यूनिवर्सल बेसिक इनकम अनुदान प्राप्त करें। रजिस्ट्रेशन एक बार, स्थायी और पूरी तरह निःशुल्क है। कोई व्यक्तिगत डेटा संग्रहीत नहीं किया जाता।',
  'app-title':'एंड्रॉयड ऐप के माध्यम से रजिस्ट्रेशन',
  'app-text':'मानवता का प्रमाण एक भौतिक 3-कारक बायोमेट्रिक प्रणाली का उपयोग करता है। चरण 1: R503 ऑप्टिकल सेंसर सभी 10 उंगलियों के निशान स्कैन करता है + MAX30102 PPG जीवित नाड़ी की पुष्टि करता है। चरण 2: हाथ की नस IR कैमरा (1/10⁷ विशिष्टता)। चरण 3: आईरिस स्कैन — स्वर्ण मानक, 1/10⁷⁸, पूर्णतः डिवाइस-स्वतंत्र। एक Groth16 ZK प्रमाण सभी कारकों को बिना किसी बायोमेट्रिक डेटा के प्रकट किए प्रतिबद्ध करता है। सत्यापन पर आपके 1,000 AEQ स्वचालित रूप से जमा हो जाते हैं।',
  's1t':'बायोमेट्रिक स्कैन','s1d':'AequitasBio सभी 10 उंगलियों के निशान स्कैन करता है (R503 ऑप्टिकल सेंसर) + MAX30102 PPG नाड़ी जीवितता की पुष्टि करती है। चरण 2: हाथ की नस IR। चरण 3: आईरिस। कच्चा डेटा कभी डिवाइस नहीं छोड़ता।',
  's2t':'ZK प्रमाण जनरेशन','s2d':'Groth16 ZK प्रमाण सभी बायोमेट्रिक कारकों को प्रतिबद्ध करता है: commitment = keccak256(iris‖vein‖fingers‖wallet)। Nullifier शरीर से बंधा है, फोन से नहीं — डिवाइस खोने से दूसरी पहचान नहीं बन सकती।',
  's3t':'वॉलेट कनेक्ट करें','s3d':'ऐप इस पेज पर MetaMask खोलती है · अपना Ethereum वॉलेट कनेक्ट करें · प्रमाण आपके वॉलेट पते से क्रिप्टोग्राफिक रूप से जुड़ा है',
  's4t':'1,000 AEQ प्रदान','s4d':'Aequitas BlockDAG पर 6 सेकंड में रजिस्ट्रेशन की पुष्टि · 1,000 AEQ तुरंत जमा · आपकी पहचान स्थायी रूप से दर्ज',
  'priv-bar':'🔒 R503 सभी 10 उंगलियां · MAX30102 जीवितता · चरण 2: हाथ नस IR · चरण 3: आईरिस (10⁷⁸) · Groth16 ZKP · डेटा डिवाइस नहीं छोड़ता · एक मानव · हमेशा के लिए',
  'conn-wallet':'कनेक्टेड वॉलेट','proof-recv':'⚡ ZK प्रमाण प्राप्त','proof-hint':'रजिस्टर करने के लिए वॉलेट कनेक्ट करें',
  'btn-conn':'🦊 METAMASK कनेक्ट करें','btn-reg':'🔐 ON-CHAIN रजिस्टर करें',
  'btn-web-reg':'🌐 ब्राउज़र के माध्यम से रजिस्टर करें (WebAuthn)',
  'web-reg-warn':'⚠ डिवाइस-बाउंड: यह पहचान इस डिवाइस और ब्राउज़र से जुड़ी है। इसे किसी अन्य डिवाइस पर स्थानांतरित नहीं किया जा सकता। स्थायी मल्टी-डिवाइस पहचान के लिए Aequitas Android App उपयोग करें।',
  'reg-log-hint':'// अपना प्रमाण उत्पन्न करने के लिए Aequitas Android App खोलें, फिर यहाँ वापस आएं...',
  'reg-details':'रजिस्ट्रेशन विवरण','k-network':'नेटवर्क','k-chainid':'चेन ID','k-grant':'UBI अनुदान',
  'k-fee':'गैस शुल्क','free':'निःशुल्क — पूरी तरह गैसलेस','k-limit':'रजिस्ट्रेशन','k-limit-v':'प्रति मानव एक बार · स्थायी · अपरिवर्तनीय',
  'k-bio':'बायोमेट्रिक डेटा','never-stored':'कभी संग्रहीत नहीं — आपके डिवाइस पर रहता है',
  'k-proof':'प्रमाण प्रणाली','k-conf':'पुष्टि','k-conf-v':'6 सेकंड के भीतर (1 ब्लॉक)',
  'k-sybil':'Sybil सुरक्षा','k-sybil-v':'प्रति बायोमेट्रिक एक पहचान · स्थायी लॉक',
  'live-stats':'लाइव चेन सांख्यिकी',
  's-height':'ब्लॉक हाइट','s-height-sub':'हर ~6s में नया ब्लॉक · BlockDAG · समानांतर उत्पादन',
  's-humans':'सत्यापित मनुष्य','s-humans-sub':'बायोमेट्रिक ZKP · एक व्यक्ति, एक वॉलेट, हमेशा के लिए',
  's-supply':'कुल आपूर्ति','s-supply-sub':'हमेशा = मनुष्य × 1,000 AEQ',
  's-index':'Aequitas इंडेक्स','s-index-sub':'0 = पूर्ण समानता · 100 = अधिकतम असमानता',
  's-uptime':'अपटाइम','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'मानवता का प्रमाण','ib-poh-t':'प्रत्येक AEQ धारक को क्रिप्टोग्राफिक रूप से साबित करना होगा कि वे एक अद्वितीय जीवित मानव हैं। कोई बॉट, कंपनी या AI नहीं। बायोमेट्रिक डेटा कभी डिवाइस नहीं छोड़ता।',
  'ib-fair':'मौलिक रूप से उचित वितरण','ib-fair-t':'प्रत्येक सत्यापित मानव को रजिस्ट्रेशन पर बिल्कुल 1,000 AEQ मिलता है। कोई प्री-माइनिंग नहीं। कुल आपूर्ति = मनुष्य × 1,000।',
  'ib-dag':'BlockDAG आर्किटेक्चर','ib-dag-t':'कई ब्लॉक एक साथ उत्पन्न और मर्ज किए जा सकते हैं। उच्च थ्रूपुट, कम विलंबता।',
  'ib-gas':'सच में निःशुल्क','ib-gas-t':'रजिस्ट्रेशन और AEQ ट्रांसफर में कुछ भी खर्च नहीं होता। ETH, BNB या MATIC की जरूरत नहीं।',
  'recent-blocks':'हालिया ब्लॉक','blocks-desc':'MERGE = कई पेरेंट मर्ज (BlockDAG)। TX = रजिस्ट्रेशन ट्रांजेक्शन। ब्लॉक समय: ~6 सेकंड।',
  'loading':'ब्लॉक लोड हो रहे हैं...','net-info':'नेटवर्क जानकारी','k-chain':'चेन नाम','k-symbol':'प्रतीक','k-btime':'ब्लॉक समय',
  'k-cons':'सहमति','k-nodes':'सक्रिय नोड्स','k-storage':'स्टोरेज','add-mm':'🦊 METAMASK में जोड़ें','k-dec':'दशमलव',
  'btn-add-mm':'+ AEQUITAS नेटवर्क जोड़ें',
  'phil':'"पैसा इसलिए है क्योंकि लोग हैं।<br>इससे ज़्यादा नहीं, इससे कम नहीं।"','phil-sub':'— AEQUITAS सिद्धांत —',
  'humans-title':'Aequitas Chain पर सत्यापित मनुष्य',
  'h-what':'सत्यापित मानव क्या है?','h-what-t':'सत्यापित मानव एक वॉलेट पता है जो क्रिप्टोग्राफिक रूप से सिद्ध है कि यह एक अद्वितीय जीवित मानव का है। सत्यापन 3-कारक हार्डवेयर प्रणाली का उपयोग करता है: R503 सभी 10 उंगलियों के निशान स्कैन करता है; MAX30102 PPG जीवित नाड़ी की पुष्टि करता है; चरण 2: हाथ की नस IR (1/10⁷); चरण 3: आईरिस (1/10⁷⁸)। केवल Groth16 ZK प्रमाण प्रेषित होता है। कोई बायोमेट्रिक डेटा डिवाइस नहीं छोड़ता।',
  'h-zkp':'ZK प्रमाण प्रणाली','h-zkp-t':'Aequitas BN128 पर Groth16 उपयोग करता है — Ethereum और Zcash जैसा ही वक्र। ~200 बाइट, ~10ms। commitment = keccak256(iris‖vein‖fingers‖wallet)। Nullifier शरीर से बंधा है: फोन खोने से दूसरी पहचान नहीं बनती। कोई बायोमेट्रिक डेटा कभी संग्रहीत नहीं होता।',
  'h-sybil':'Sybil अटैक रोकथाम','h-sybil-t':'चरण 1: सभी 10 उंगलियों के निशान + MAX30102 जीवितता (PPG नाड़ी, सांचे/रीप्ले अस्वीकार)। चरण 2: हाथ की नस IR — आंतरिक शारीरिक विशेषता, नकल असंभव, 1/10⁷, समान जुड़वां में भी भिन्न। चरण 3: आईरिस — 1/10⁷⁸, वैश्विक स्वर्ण मानक। Nullifier = keccak256(iris‖vein‖domain)। एक मानव, एक वॉलेट, हमेशा के लिए।',
  'h-global':'वैश्विक वित्तीय समावेशन','h-global-t':'कोई बैंक खाता, क्रेडिट कार्ड या क्रिप्टोकरेंसी की जरूरत नहीं। बस बायोमेट्रिक सेंसर वाला Android स्मार्टफोन।',
  'h-bio-hw':'बायोमेट्रिक हार्डवेयर रोडमैप','h-bio-hw-t':'चरण 1 (सक्रिय): R503 ऑप्टिकल फिंगरप्रिंट स्कैनर — सभी 10 उंगलियों का संयुक्त हैश। MAX30102 PPG जीवितता। चरण 2 (नियोजित): ESP32-CAM + 850 nm IR LED — हाथ की नस इमेजिंग, 1/10⁷ विशिष्टता। चरण 3 (नियोजित): IR आईरिस मॉड्यूल — 240+ स्वतंत्रता की डिग्री, 1/10⁷⁸, पूर्णतः डिवाइस-स्वतंत्र, समान जुड़वां भिन्न।',
  'reg-humans':'रजिस्टर्ड मनुष्य','h-desc':'प्रत्येक पता बायोमेट्रिक ZKP के माध्यम से अद्वितीय मानव के रूप में सत्यापित। प्रत्येक को बिल्कुल 1,000 AEQ मिला। स्थायी, अपरिवर्तनीय, ऑन-चेन।',
  'no-humans':'अभी तक कोई मानव रजिस्टर्ड नहीं।\n\nAequitas Android App डाउनलोड करें और चेन पर पहले मानव बनें!',
  'reg-stats':'रजिस्ट्री आँकड़े','total-humans':'कुल मनुष्य',
  'idx-title':'Aequitas इंडेक्स — रियल-टाइम आर्थिक समानता स्कोर',
  'idx-desc':'Aequitas इंडेक्स <strong style="color:var(--teal)">जिनी गुणांक</strong> से लिया गया है — विश्व बैंक, OECD और UN द्वारा अपनाया गया अंतरराष्ट्रीय मानक। <strong style="color:var(--neon)">0 = पूर्ण समानता</strong>। <strong style="color:var(--red)">100 = अधिकतम एकाग्रता</strong>। लक्ष्य: जिनी 0.30 से कम।',
  'gini-what-title':'जिनी गुणांक क्या है?',
  'gini-what-text':'इतालवी सांख्यिकीविद् कोर्राडो जिनी (1912) द्वारा विकसित। धन वितरण मापता है। पैमाना: 0 (सब समान) से 1 (एक व्यक्ति के पास सब कुछ)। विश्व बैंक, OECD, UN उपयोग करते हैं।',
  'curr-idx':'वर्तमान इंडेक्स','bar-0':'0 — पूर्ण समानता','bar-100':'100 — अधिकतम असमानता','wcap-lbl':'वर्तमान धन सीमा:','wcap-mult':'गुणक:','wcap-avg':'औसत बैलेंस:',
  'gini':'जिनी गुणांक','gini-desc':'0 = समान · 1 = असमान',
  'supply-desc':'हमेशा = मनुष्य × 1,000 AEQ',
  'phase':'प्रोटोकॉल चरण','phase-desc':'मानवों की संख्या से स्वचालित रूप से आगे बढ़ता है',
  'humans-desc':'बायोमेट्रिक रूप से सत्यापित अद्वितीय मनुष्य',
  'pools-title':'पुनर्वितरण पूल',
  'pools-desc':'प्रत्येक स्वैप शुल्क, डेमरेज और धन सीमा अधिशेष स्वचालित रूप से चार पूलों में विभाजित होता है। सभी पूल दैनिक भुगतान करते हैं।',
  'vel-pool':'वैलिडेटर पूल','vel-pool-desc':'सभी शुल्कों का 40% → नोड ऑपरेटर जो नेटवर्क सुरक्षित करते हैं',
  'liq-pool':'लिक्विडिटी पूल','liq-pool-desc':'सभी शुल्कों का 30% → लिक्विडिटी प्रदाता, LP शेयर के अनुपात में',
  'ubi-pool':'UBI पूल','ubi-pool-desc':'सभी शुल्कों का 20% → सभी सत्यापित मनुष्यों को समान रूप से, हर 24 घंटे',
  'treasury':'ट्रेजरी','treasury-desc':'सभी शुल्कों का 10% → प्रोटोकॉल विकास और रखरखाव',
  'phases-title':'प्रोटोकॉल चरण',
  'demurrage-title':'डेमरेज — परिसंचरण के लिए प्रोत्साहन',
  'demurrage-desc':'निष्क्रिय AEQ बैलेंस धीरे-धीरे मूल्य खोते हैं ताकि संचय को हतोत्साहित किया जा सके।',
  'dem-rate-k':'क्षय दर','dem-rate-v':'0.5% प्रति माह (निरंतर)',
  'dem-grace-k':'ग्रेस पीरियड','dem-grace-v':'क्षय शुरू होने से पहले 3 महीने की निष्क्रियता',
  'dem-reset-k':'रीसेट','dem-reset-v':'कोई भी ट्रांसफर, स्वैप या लिक्विडिटी एक्शन टाइमर शून्य करता है',
  'dem-dest-k':'क्षयित AEQ जाता है','dem-dest-v':'पुनर्वितरण पूल में (40/30/20/10 विभाजन)',
  'dem-warn-k':'चेतावनी प्रणाली','dem-warn-v':'14 दिन की सूचना (एक बार) + हर लॉगिन पर 7 दिन का अनुस्मारक',
  'story-title':'Aequitas की कहानी',
  'story-text':'<p>2009 में सातोशी नाकामोतो ने Bitcoin जारी किया। पहली बार बैंक के बिना मूल्य हस्तांतरण संभव हुआ। एक सच्ची क्रांति। लेकिन लगभग तुरंत कुछ गलत हो गया।</p><p>शुरुआती माइनर्स ने लाखों सिक्के लगभग शून्य लागत पर जमा किए। 2021 में, शीर्ष 1% Bitcoin पते 90% से अधिक Bitcoin नियंत्रित करते हैं। Bitcoin का जिनी गुणांक 0.85 से अधिक है।</p><p><span style="color:var(--gold)">Aequitas</span> — "न्याय" के लिए लैटिन — एक प्रश्न का उत्तर देने के लिए बनाया गया: <em style="color:var(--gold)">"एक क्रिप्टोकरेंसी कैसी दिखेगी जो हर मानव के लिए न्यायपूर्ण हो?"</em></p><p><strong style="color:var(--text)">पैसा इसलिए है क्योंकि लोग हैं। इसलिए हर व्यक्ति को केवल मानव होने के कारण धन का समान हिस्सा मिलना चाहिए।</strong></p>',
  'nodes-title':'सक्रिय नोड्स — वर्तमान नेटवर्क टोपोलॉजी',
  'nodes-desc':'Aequitas नेटवर्क वर्तमान में दो भौगोलिक रूप से वितरित नोड्स पर चलता है। दोनों ब्लॉक उत्पादन, स्टेट सिंक्रोनाइज़ेशन और API सेवा में भाग लेते हैं।',
  'node1':'नोड 1 — Railway (प्राथमिक)','node1-desc':'प्राथमिक API · ब्लॉक उत्पादक · UBI वितरण · P2P Bootstrap · PostgreSQL · MetaMask के लिए RPC',
  'node2':'नोड 2 — Render (द्वितीयक)','node2-desc':'द्वितीयक API · ब्लॉक उत्पादक · P2P पीयर · HTTP सिंक · साझा PostgreSQL स्टेट',
  'run-node-title':'अपना नोड चलाएं','run-node-desc':'कोई भी Aequitas नोड चला सकता है — बिना अनुमति, बिना स्टेक। ऑपरेटर दैनिक वितरित स्वैप शुल्क का 40% कमाते हैं।',
  'bootstrap-title':'नया नोड कनेक्ट करें','bootstrap-desc':'PRIMARY_NODE_URL=https://aequitas.digital अपने environment में सेट करें। आपका नोड स्वचालित रूप से पूर्ण चेन स्टेट सिंक करेगा।',
  'tech-title':'तकनीकी विशिष्टताएं','mm-config':'MetaMask कॉन्फ़िगरेशन',
  'k-lang':'भाषा','k-src':'स्रोत','evm-yes':'हाँ — JSON-RPC /rpc · MetaMask संगत',
  'proto-label':'Aequitas V7 प्रोटोकॉल — तकनीकी दस्तावेज़ीकरण',
  'ca-title':'अनुबंध पते',
  'ca-text':'चेन: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7: 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7 पूरी Aequitas अर्थव्यवस्था के लिए एकमात्र सच्चाई का स्रोत है। कोई एडमिन की, अपग्रेड प्रॉक्सी या गवर्नेंस वोट इसका तर्क नहीं बदल सकता।',
  'poa-title':'1. जीवन का प्रमाण','poa-text':'<p>जब लोग मरते हैं, उनका AEQ धीरे-धीरे UBI पूल के माध्यम से समुदाय को वापस जाता है, बजाय Bitcoin की तरह हमेशा के लिए खोने के।</p>',
  'poa-box':'वर्ष 0–2: सामान्य उपयोग<br>वर्ष 2: चेतावनी 1 — Guardian जवाब दे सकता है<br>वर्ष 2+60 दिन: चेतावनी 2<br>वर्ष 2+120 दिन: चेतावनी 3<br>वर्ष 2+180 दिन: AEQ व्यक्तिगत एस्क्रो में<br>वर्ष 4: निष्क्रिय रहने पर — UBI पूल में वापस',
  'guard-title':'2. गार्जियन सिस्टम','guard-text':'<p>एक विश्वसनीय Guardian (दूसरा सत्यापित मानव) पुष्टि कर सकता है कि कोई अभी भी जीवित है, बिना किसी वित्तीय अधिकार के।</p>',
  'guard-box':'प्रति मानव 1 Guardian · दूसरा सत्यापित मानव होना चाहिए<br>Guardian केवल confirmAlive() कॉल कर सकता है · शून्य वित्तीय अधिकार<br>Guardian धन नहीं हिला सकता · अधिकतम 3 · Timelock 7 दिन',
  'dem-title':'3. डेमरेज — संचय-विरोधी तंत्र',
  'dem-box':'दर: 3 महीने की छूट के बाद 0.5%/माह<br>किसी भी ट्रांसफर, स्वैप या लिक्विडिटी पर रीसेट<br>क्षयित AEQ पूलों में पुनर्वितरित (जला नहीं जाता)',
  'dem-text':'<p>ऐतिहासिक उदाहरण: Wörgl प्रयोग (ऑस्ट्रिया, 1932) — एक वर्ष में बेरोजगारी 25% कम। Chiemgauer (जर्मनी, 2003) — 20+ वर्षों से सफलतापूर्वक चल रहा है।</p>',
  'cap-title':'4. धन सीमा — गणितीय निष्पक्षता','cap-box':'सीमा: max(5,min(N,25))× औसत AEQ बैलेंस<br>1–4 मनुष्य: 5× · प्रति मानव +1× · 25+: 25× स्थायी<br>अतिरिक्त AEQ तुरंत पुनर्वितरित · कोई हस्तक्षेप नहीं',
  'ubi-title':'5. यूनिवर्सल बेसिक इनकम','ubi-box':'स्रोत: स्वैप शुल्क (20%) · सीमा अधिशेष · डेमरेज<br><br>दैनिक: UBI पूल सभी पंजीकृत मनुष्यों में समान रूप से विभाजित। प्रत्येक वितरण के बाद पूल शून्य हो जाता है।',
  'inf-title':'6. कोई एल्गोरिदमिक मुद्रास्फीति नहीं','inf-box':'केवल एक घटना नया AEQ बनाती है: नया सत्यापित मानव पंजीकृत होता है।<br><br>कुल आपूर्ति = सत्यापित मनुष्य × 1,000 AEQ — हमेशा, बिल्कुल।',
  'btn-download-app':'AEQUITASBIO ऐप डाउनलोड करें',
  'swap-title':'🔄 AEQ ↔ tUSD स्वैप करें','swap-sub':'नेटिव लिक्विडिटी पूल के माध्यम से AEQ को tUSD (सिमुलेटेड टेस्ट डॉलर) से बदलें। स्वैप के लिए केवल 0.1% शुल्क — सामान्य AEQ ट्रांसफर पूरी तरह निःशुल्क।',
  'swap-priv-bar':'🔒 केवल 0.1% स्वैप शुल्क · AEQ→AEQ ट्रांसफर निःशुल्क · tUSD कोई वास्तविक मूल्य के बिना टेस्ट मुद्रा है',
  'swap-your-aeq':'आपका AEQ','swap-your-tusd':'आपका tUSD','swap-aeq-to-tusd':'AEQ → tUSD','swap-tusd-to-aeq':'tUSD → AEQ',
  'swap-fee-est':'प्रोटोकॉल शुल्क (0.1%)','swap-details-hdr':'स्वैप विवरण',
  'swap-out-lbl':'आप प्राप्त करेंगे (अनुमानित)','swap-impact-lbl':'मूल्य प्रभाव','swap-rate-lbl':'विनिमय दर',
  'swap-depth-lbl':'पूल संरचना','amm-title':'x × y = k — कॉन्स्टेंट प्रोडक्ट AMM',
  'amm-text':'AEQ स्वैप करते समय AEQ रिजर्व बढ़ता है और tUSD रिजर्व घटता है — उनका गुणनफल हमेशा k के बराबर रहता है। बड़े स्वैप से मूल्य पर अधिक प्रभाव।',
  'swap-btn-conn':'🦊 METAMASK कनेक्ट करें','swap-btn-go':'🔄 स्वैप करें',
  'swap-log-hint':'// स्वैप करने के लिए वॉलेट कनेक्ट करें...',
  'swap-no-liquidity':'अभी tUSD नहीं है?','swap-faucet-desc':'पंजीकृत मनुष्य एक बार टेस्ट tUSD का दावा कर सकते हैं','swap-btn-faucet':'💧 टेस्ट tUSD का दावा करें',
  'swap-addliq-title':'लिक्विडिटी प्रदान करें','swap-addliq-desc':'पहले डिपॉजिट करें — आपका अनुपात प्रारंभिक मूल्य तय करता है।','swap-btn-addliq':'💧 लिक्विडिटी जोड़ें',
  'swap-lp-title':'आपकी LP स्थिति','swap-lp-share':'पूल हिस्सा','swap-lp-withdrawable':'निकालने योग्य',
  'swap-lp-pct-label':'आपकी स्थिति का %','swap-lp-youget':'आप प्राप्त करेंगे','swap-btn-removeliq':'🔥 लिक्विडिटी हटाएं',
  'swap-pool-title':'AEQ / tUSD — पूल स्थिति',
  'swap-pool-aeq':'AEQ रिजर्व','swap-pool-tusd':'tUSD रिजर्व','swap-pool-price':'स्पॉट मूल्य',
  'swap-fee-bps':'स्वैप शुल्क','swap-fee-split':'शुल्क वितरण','swap-fee-split-v':'40% वैलिडेटर / 30% LP / 20% UBI / 10% ट्रेजरी',
  'swap-pools-addr-title':'टोकनोमिक्स पूल पते',
  'swap-validators':'वैलिडेटर (40%)','swap-lps':'लिक्विडिटी प्रदाता (30%)','swap-ubi':'UBI पूल (20%)','swap-treasury':'ट्रेजरी (10%)',
  'ubi-hero-title':'यूनिवर्सल बेसिक इनकम — UBI पूल',
  'ubi-hero-sub':'जमा हो रहा है — अगला भुगतान सभी सत्यापित मनुष्यों को समान रूप से वितरित:',
  'ubi-bal-lbl':'वर्तमान पूल बैलेंस','ubi-hero-desc':'समान रूप से विभाजित · हर 24 घंटे भुगतान · पूल शून्य होता है · न्यूनतम बैलेंस की जरूरत नहीं',
  'ubi-how-fills':'UBI पूल कैसे भरता है',
  'ubi-src-swap':'स्वैप शुल्क','ubi-src-swap-d':'प्रत्येक AEQ↔tUSD स्वैप अपने 0.1% शुल्क का 20% योगदान देता है।',
  'ubi-src-dem':'डेमरेज','ubi-src-dem-d':'निष्क्रिय AEQ (3+ माह) 0.5%/माह क्षय होता है। क्षयित राशि का 20% UBI में जाता है।',
  'ubi-src-cap':'सीमा अधिशेष','ubi-src-cap-d':'सीमा से अधिक वॉलेट तुरंत कटते हैं। 20% UBI में प्रवाहित होता है।',
  'pools4-header':'चारों पुनर्वितरण पूल',
  'ubi-see-above':'ऊपर काउंटडाउन देखें','ubi-timer-above':'⏰ काउंटडाउन ऊपर दिखाया गया','pool-t-timer':'जमा हो रहा है — कोई टाइमर नहीं',
  'usp-headline':'इतिहास में पहली बार — सब एक समान से शुरू करते हैं',
  'usp-sub':'अगर आपके पास Android स्मार्टफोन है तो आप पात्र हैं। बिना बैंक, बिना क्रिप्टो ज्ञान, बिना निवेश।',
  'usp-c1-title':'₹0 प्रारंभिक निवेश','usp-c1-desc':'रजिस्ट्रेशन पूरी तरह निःशुल्क। कोई ETH, MATIC या क्रेडिट कार्ड नहीं। प्रोटोकॉल सभी लागत वहन करता है।',
  'usp-c2-title':'प्रत्येक मानव के लिए 1,000 AEQ','usp-c2-desc':'अरबपति हो या किसान — सभी को बिल्कुल 1,000 AEQ मिलता है। गणितीय गारंटी के साथ समान शुरुआत।',
  'usp-c3-title':'सभी के लिए सुलभ','usp-c3-desc':'कोई बैंक खाता, क्रेडिट कार्ड या सरकारी ID नहीं चाहिए। पंजीकरण सस्ती बायोमेट्रिक हार्डवेयर किट (फिंगरप्रिंट स्कैनर + पल्स सेंसर, ~$15) से होता है — वैश्विक पहुंच के लिए डिज़ाइन।',
  'usp-c4-title':'हमेशा के लिए दैनिक UBI','usp-c4-desc':'पंजीकरण के बाद, आपका UBI हिस्सा हर दिन स्वचालित रूप से आता है — बिना किसी कार्रवाई के।',
  'v7-intro-title':'AequitasV7 क्या है?',
  'v7-intro-text':'AequitasV7, Aequitas प्रोटोकॉल का केंद्रीय स्मार्ट अनुबंध है। Aequitas Chain (ID 1926) पर अपरिवर्तनीय रूप से तैनात। सब कुछ प्रबंधित करता है: मानव पंजीकरण, ZK सत्यापन, बैलेंस प्रबंधन, धन सीमा, UBI वितरण, स्वैप शुल्क। कोई व्यवस्थापक इसे अपडेट नहीं कर सकता।',
  'explore-title':'Aequitas एक्सप्लोर करें',
  'expl-score':'समानता स्कोर','expl-score-d':'लाइव जिनी गुणांक · Aequitas इंडेक्स · रियल-टाइम धन वितरण',
  'expl-economy':'UBI और पुनर्वितरण','expl-economy-d':'दैनिक UBI काउंटडाउन · 4 ऑन-चेन पूल · डेमरेज · प्रोटोकॉल चरण',
  'expl-charts':'चार्ट और इतिहास','expl-charts-d':'जिनी इतिहास · लॉरेंज वक्र · धन सीमा स्लाइडर · Aequitas की कहानी',
  'expl-v7':'प्रोटोकॉल V7 दस्तावेज़','expl-v7-d':'AequitasV7 अनुबंध · 6 तंत्र · ZK प्रमाण · धन सीमा · डेमरेज · अपरिवर्तनीय कोड',
  'expl-explorer':'ब्लॉक एक्सप्लोरर','expl-explorer-d':'लाइव BlockDAG · वैलिडेटर, हैश, ट्रांजेक्शन देखने के लिए किसी भी ब्लॉक पर क्लिक करें',
  'swap-sell-label':'बेचें','swap-receive-label':'प्राप्त करें',
  'expl-network':'नेटवर्क और नोड्स','expl-network-d':'नोड टोपोलॉजी · अपना नोड चलाएं · तकनीकी विशिष्टताएं · Chain ID 1926'
}
};

function showStab(parentId, stabId, el) {
  const parent = document.getElementById(parentId);
  parent.querySelectorAll('.stab-panel').forEach(p => p.classList.remove('active'));
  parent.querySelectorAll('.stab').forEach(s => s.classList.remove('active'));
  document.getElementById(stabId).classList.add('active');
  el.classList.add('active');
  if (stabId === 'eqi-score') { setTimeout(function(){ drawGiniHistoryChart(); }, 30); }
  if (stabId === 'eqi-lorenz') { setTimeout(drawLorenzCurve, 30); }
  if (stabId === 'eqi-economy') { setTimeout(drawWcapSlideChart, 30); }
  // Push sub-route URL
  const tabSlugMap = {'tab-register':'register','tab-explorer':'explorer','tab-index':'index','tab-network':'network','tab-exchange':'exchange'};
  const stabSlugMap = {'sep-blocks':'blocks','sep-humans':'humans','eqi-score':'score','eqi-lorenz':'distribution','eqi-economy':'economy','eqi-charts':'charts','net-overview':'overview','net-runnode':'node','net-protocol':'protocol','exch-swap':'swap','exch-liquidity':'liquidity'};
  const tabSlug = tabSlugMap[parentId];
  const stabSlug = stabSlugMap[stabId];
  if (tabSlug && stabSlug) history.pushState(null, '', '/' + tabSlug + '/' + stabSlug);
}

function showTab(name, el) {
  // Remove the server-injected data-active attribute so CSS rules like
  // html[data-active=network] #tab-network {display:block!important} no
  // longer force the old tab visible after a JS-driven tab switch.
  document.documentElement.removeAttribute('data-active');
  document.querySelectorAll('.tab-content').forEach(t => {
    t.classList.remove('active');
    t.style.display = '';
  });
  document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
  const tabContent = document.getElementById('tab-' + name);
  if (!tabContent) return;
  tabContent.classList.add('active');
  el.classList.add('active');
  // Always activate first stab-panel when switching tabs
  var panels2 = tabContent.querySelectorAll('.stab-panel');
  var stabs2 = tabContent.querySelectorAll('.stab');
  panels2.forEach(function(p){p.classList.remove('active');});
  stabs2.forEach(function(s){s.classList.remove('active');});
  if (panels2.length) panels2[0].classList.add('active');
  if (stabs2.length) stabs2[0].classList.add('active');
  if (name === 'exchange') { loadPoolStatus(); preloadPriceHistory(); }
  history.pushState(null, '', '/' + name);
}

function goTab(name, stabId) {
  let el = null;
  document.querySelectorAll('.tab').forEach(t => {
    if ((t.getAttribute('onclick') || '').includes("'" + name + "'")) el = t;
  });
  if (el) showTab(name, el);
  if (stabId) {
    const stabEl = document.querySelector('#tab-' + name + ' .stab[onclick*="\'' + stabId + '\'"]');
    if (stabEl) showStab('tab-' + name, stabId, stabEl);
  }
}

function setLang(lang) {
  curLang = lang;
  document.getElementById('lang-sel').value = lang;
  document.documentElement.dir = lang === 'ar' ? 'rtl' : 'ltr';
  document.documentElement.lang = lang;
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
  if (!window.ethereum) { addLog('🦊 MetaMask not found — <a href="https://metamask.io/download/" target="_blank" style="color:var(--gold)">install MetaMask</a> to use this feature.', 'warn'); return; }
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
        rpcUrls: ['https://aequitas.digital/rpc'],
        blockExplorerUrls: ['https://aequitas.digital']
      }]
    });
  } catch (e) { console.error('MetaMask error:', e); }
}

// UBI countdown timer — counts down to the next daily distribution.
// secsRemaining comes from the server (uptime modulo 86400 subtracted from 86400).
// Once it reaches zero it resets to 24h and keeps ticking, since the
// distribution just ran and the next one is 24h away again.
let ubiTimerInterval = null;
function startUBITimer(secsRemaining) {
  if (ubiTimerInterval) clearInterval(ubiTimerInterval);
  let secs = secsRemaining;
  const els = [
    document.getElementById('ubi-timer'),
    document.getElementById('validators-timer'),
    document.getElementById('lp-timer'),
  ].filter(Boolean);
  if (!els.length) return;

  const fmt = s => {
    const h = Math.floor(s / 3600);
    const m = Math.floor((s % 3600) / 60);
    const sec = s % 60;
    return String(h).padStart(2,'0') + 'h ' + String(m).padStart(2,'0') + 'm ' + String(sec).padStart(2,'0') + 's';
  };

  els.forEach(el => el.textContent = fmt(secs));
  ubiTimerInterval = setInterval(() => {
    secs--;
    if (secs <= 0) {
      secs = 86400;
      els.forEach(el => { el.style.color = 'var(--green)'; });
      setTimeout(() => { els.forEach(el => { el.style.color = ''; }); }, 3000);
      // Refresh only pool balances after distribution — do NOT call loadStatus()
      // here because loadStatus() calls startUBITimer() which would restart the
      // timer from 0 again if next_ubi_at hasn't been written yet → 2s reset loop.
      setTimeout(() => { if (typeof loadPoolStatus === 'function') loadPoolStatus(); }, 3000);
    }
    els.forEach(el => el.textContent = fmt(secs));
  }, 1000);
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
    document.getElementById('idx-gini').textContent = typeof d.gini === 'number' ? d.gini.toFixed(4) : '—';
    const gniWarn = document.getElementById('gini-n-warn');
    if (gniWarn) gniWarn.style.display = (d.total_humans < 10) ? 'block' : 'none';
    document.getElementById('idx-supply2').textContent = d.total_supply || '—';
    document.getElementById('idx-phase').textContent = fmt(d.phase);
    document.getElementById('idx-humans2').textContent = fmt(d.total_humans);
    document.getElementById('stat-humans').textContent = fmt(d.total_humans);
    document.getElementById('stat-supply').textContent = d.total_supply || '—';

    // Pool balances — show 0.0000 instead of — when pool is empty
    const fmtPool = v => (v || '0.0000') + ' AEQ';
    document.getElementById('pool-v').textContent = fmtPool(d.pool_validators);
    document.getElementById('pool-l').textContent = fmtPool(d.pool_lp);
    document.getElementById('pool-u').textContent = fmtPool(d.pool_ubi);
    document.getElementById('pool-t').textContent = fmtPool(d.pool_treasury);

    // UBI countdown timer + fill bar
    // Only (re)start the timer when the server returns a positive value.
    // When ubi_next_payout_secs === 0 (IS_PRIMARY_NODE not set, or next_ubi_at
    // not yet written to DB), leave the running timer alone — restarting from 0
    // causes a reset loop because loadStatus fires every 6s.
    if (d.ubi_next_payout_secs > 0) {
      startUBITimer(d.ubi_next_payout_secs);
    }
    const fillSecs = d.ubi_next_payout_secs || 0;
    const fillPct = Math.min(100, Math.max(0, (86400 - fillSecs) / 86400 * 100));
    const fillBar = document.getElementById('ubi-fill-bar');
    if (fillBar) fillBar.style.width = fillPct.toFixed(1) + '%';

    // Fix stale subtitle now that demurrage/wealth-cap mean supply can drift
    const subEl = document.getElementById('s-supply-sub');
    if (subEl) subEl.textContent = 'Always = Humans × 1,000 AEQ';

    if (d.index !== undefined) {
      document.getElementById('idx-bar').style.width = Math.min(d.index, 100) + '%';
      const phases = ['Phase 0: Bootstrap — sliding wealth cap 5×→25× (active)', 'Phase 1: Growth — expanding human registry (cap: 25×)', 'Phase 2: Stability — redistribution active (cap: 25×)', 'Phase 3: Maturity — full decentralization (cap: 25×)'];
      document.getElementById('idx-phase-desc').textContent = phases[d.phase || 0] || 'Phase ' + (d.phase || 0);
    }
  } catch (e) {}
  // Populate live wealth-cap widget (non-blocking)
  try {
    const wc = await (await fetch('/api/wealth-cap')).json();
    const capEl = document.getElementById('live-cap-aeq');
    const multEl = document.getElementById('live-cap-mult');
    const avgEl = document.getElementById('live-cap-avg');
    if (capEl && wc.cap_aeq !== undefined) capEl.textContent = wc.cap_aeq.toFixed(2);
    if (multEl && wc.multiplier !== undefined) multEl.textContent = wc.multiplier.toFixed(0) + '×';
    if (avgEl && wc.average_aeq !== undefined) avgEl.textContent = wc.average_aeq.toFixed(2);
  } catch(_) {}
}

async function drawGiniHistoryChart() {
  const canvas = document.getElementById('gini-history-chart');
  if (!canvas || !canvas.offsetParent) return;
  canvas.width = canvas.offsetWidth;
  const ctx = canvas.getContext('2d');
  const W = canvas.width, H = canvas.height;
  ctx.clearRect(0, 0, W, H);
  try {
    const d = await (await fetch('/api/gini/history')).json();
    const history = (d.history || []).slice().reverse();
    const emptyEl = document.getElementById('gini-history-empty');
    if (!history.length) {
      if (emptyEl) { emptyEl.style.display = 'block'; canvas.style.display = 'none'; } return;
    }
    if (emptyEl) { emptyEl.style.display = 'none'; canvas.style.display = 'block'; }
    // Single data point — draw a gauge/meter visualization
    if (history.length === 1) {
      var g0 = history[0].gini || (history[0].idx/100); // 0-1 scale
      // Background
      ctx.fillStyle='rgba(8,10,22,0.7)'; ctx.fillRect(0,0,W,H);
      // Horizontal bar gauge
      var bx=40, by=H/2-18, bw=W-80, bh=28, r=6;
      // Track
      ctx.fillStyle='rgba(255,255,255,0.06)';
      ctx.beginPath(); ctx.roundRect(bx,by,bw,bh,r); ctx.fill();
      // Zone colors: green 0-0.30, amber 0.30-0.70, red 0.70-1.0 (Gini 0–1 scale)
      var zones=[[0,0.30,'rgba(0,255,100,0.5)'],[0.30,0.70,'rgba(245,158,11,0.5)'],[0.70,1.0,'rgba(239,68,68,0.5)']];
      zones.forEach(function(z){
        var x1=bx+bw*z[0], x2=bx+bw*z[1];
        ctx.fillStyle=z[2]; ctx.fillRect(x1,by,x2-x1,bh);
      });
      // Fill up to current value
      var fill=bw*g0/1.0;
      var grd=ctx.createLinearGradient(bx,0,bx+fill,0);
      grd.addColorStop(0,'rgba(0,255,200,0.9)'); grd.addColorStop(0.5,'rgba(245,158,11,0.9)'); grd.addColorStop(1,'rgba(239,68,68,0.9)');
      ctx.fillStyle=grd; ctx.beginPath(); ctx.roundRect(bx,by,fill,bh,r); ctx.fill();
      // Target marker at 0.30 (Gini target)
      var tx=bx+bw*0.30;
      ctx.strokeStyle='rgba(0,255,209,0.9)'; ctx.lineWidth=2;
      ctx.beginPath(); ctx.moveTo(tx,by-6); ctx.lineTo(tx,by+bh+6); ctx.stroke();
      ctx.fillStyle='rgba(0,255,209,0.9)'; ctx.font='bold 9px JetBrains Mono,monospace'; ctx.textAlign='center';
      ctx.fillText('0.30', tx, by-10);
      // Pointer
      var px=bx+bw*g0/1.0;
      ctx.fillStyle='#fff'; ctx.beginPath(); ctx.moveTo(px,by-2); ctx.lineTo(px-5,by-10); ctx.lineTo(px+5,by-10); ctx.fill();
      // Labels: 0, 0.30, 0.70, 1.0 (Gini 0–1 scale)
      [[0,'0'],[0.30,'0.30'],[0.70,'0.70'],[1,'1.0']].forEach(function(l){
        ctx.fillStyle='rgba(200,168,76,0.5)'; ctx.font='9px JetBrains Mono,monospace'; ctx.textAlign='center';
        ctx.fillText(l[1], bx+bw*l[0], by+bh+14);
      });
      // Big value
      ctx.fillStyle='rgba(200,168,76,0.95)'; ctx.font='bold 28px JetBrains Mono,monospace'; ctx.textAlign='center';
      ctx.fillText('Gini: ' + g0.toFixed(4), W/2, by-26);
      // Description (g0 is 0–1 Gini scale, target is < 0.30)
      var label;
      if(g0<0.30) label='Below target — excellent equality';
      else if(g0<0.70) label='Above target — redistribution active';
      else label='Critical — protocol at maximum intervention';
      ctx.font='11px Inter,sans-serif'; ctx.fillStyle='rgba(200,200,200,0.6)';
      ctx.fillText(label, W/2, by+bh+28);
      ctx.font='10px Inter,sans-serif'; ctx.fillStyle='rgba(0,255,209,0.5)';
      ctx.fillText('History chart grows after each daily UBI distribution', W/2, H-10);
      return;
    }
    const pad = {l:48,r:24,t:36,b:32};
    const cW = W-pad.l-pad.r, cH = H-pad.t-pad.b;
    const toX = (i) => pad.l + cW*i/Math.max(history.length-1,1);
    const toY = (v) => pad.t + cH*(1-v/100);
    // danger zone (>70) subtle red tint
    const dg = ctx.createLinearGradient(0,toY(100),0,toY(70));
    dg.addColorStop(0,'rgba(248,113,113,0.06)'); dg.addColorStop(1,'rgba(248,113,113,0)');
    ctx.fillStyle=dg; ctx.fillRect(pad.l,toY(100),cW,toY(70)-toY(100));
    // grid lines
    for (let i=0;i<=4;i++) {
      const v=i*25, y=toY(v);
      ctx.strokeStyle = v===0?'rgba(139,92,246,0.2)':'rgba(139,92,246,0.08)';
      ctx.lineWidth = v===0?1.5:1;
      ctx.beginPath(); ctx.moveTo(pad.l,y); ctx.lineTo(W-pad.r,y); ctx.stroke();
      ctx.fillStyle='rgba(200,168,76,0.75)'; ctx.font='10px JetBrains Mono,monospace'; ctx.textAlign='right';
      ctx.fillText(v+'', pad.l-6, y+4);
    }
    // target 0.30 line (idx = gini*100, so toY(30) = Gini 0.30)
    const targetY = toY(30);
    ctx.save(); ctx.shadowColor='rgba(0,255,209,0.7)'; ctx.shadowBlur=5;
    ctx.strokeStyle='rgba(0,255,209,0.55)'; ctx.lineWidth=1.5; ctx.setLineDash([6,5]);
    ctx.beginPath(); ctx.moveTo(pad.l,targetY); ctx.lineTo(W-pad.r,targetY); ctx.stroke();
    ctx.setLineDash([]); ctx.restore();
    ctx.fillStyle='rgba(4,120,87,0.85)'; ctx.font='bold 9px JetBrains Mono,monospace'; ctx.textAlign='right';
    ctx.fillText('TARGET 0.30', W-pad.r-2, targetY-5);
    // bezier path helper
    var pathBez = function(pts) {
      ctx.moveTo(toX(0), toY(pts[0].idx));
      if (pts.length<3) { for(var k=1;k<pts.length;k++) ctx.lineTo(toX(k),toY(pts[k].idx)); return; }
      for (var k=1;k<pts.length-1;k++) {
        var mx=(toX(k)+toX(k+1))/2, my=(toY(pts[k].idx)+toY(pts[k+1].idx))/2;
        ctx.quadraticCurveTo(toX(k),toY(pts[k].idx),mx,my);
      }
      ctx.lineTo(toX(pts.length-1), toY(pts[pts.length-1].idx));
    };
    // gradient fill
    var fg=ctx.createLinearGradient(0,pad.t,0,H-pad.b);
    fg.addColorStop(0,'rgba(200,168,76,0.28)'); fg.addColorStop(0.7,'rgba(200,168,76,0.07)'); fg.addColorStop(1,'rgba(200,168,76,0.01)');
    ctx.beginPath(); pathBez(history);
    ctx.lineTo(toX(history.length-1),H-pad.b); ctx.lineTo(toX(0),H-pad.b); ctx.closePath();
    ctx.fillStyle=fg; ctx.fill();
    // glowing line
    ctx.save(); ctx.shadowColor='rgba(200,168,76,0.6)'; ctx.shadowBlur=10;
    ctx.strokeStyle='#C9A84C'; ctx.lineWidth=2.5;
    ctx.beginPath(); pathBez(history); ctx.stroke(); ctx.restore();
    // dots
    history.forEach(function(pt,i){
      var x=toX(i), y=toY(pt.idx);
      ctx.save(); ctx.shadowColor='rgba(200,168,76,0.9)'; ctx.shadowBlur=12;
      ctx.beginPath(); ctx.arc(x,y,4.5,0,2*Math.PI); ctx.fillStyle='#C9A84C'; ctx.fill(); ctx.restore();
      ctx.beginPath(); ctx.arc(x,y,2,0,2*Math.PI); ctx.fillStyle='#fff'; ctx.fill();
    });
    // latest value label
    var lpt=history[history.length-1], lx=toX(history.length-1), ly=toY(lpt.idx);
    ctx.fillStyle='rgba(200,168,76,0.95)'; ctx.font='bold 11px JetBrains Mono,monospace';
    ctx.textAlign = lx>W*0.7?'right':'left';
    ctx.fillText('Gini: '+lpt.idx.toFixed(3), lx+(lx>W*0.7?-8:8), ly-9);
    // title
    ctx.fillStyle='rgba(107,70,193,0.55)'; ctx.font='10px Inter,sans-serif'; ctx.textAlign='left';
    ctx.fillText('GINI INDEX HISTORY  —  0 = perfect equality  ·  100 = max inequality', pad.l, 20);
  } catch(e) {}
}

async function drawLorenzCurve() {
  var canvas = document.getElementById('lorenz-chart');
  if (!canvas || !canvas.offsetParent) return;
  canvas.width = canvas.offsetWidth;
  var W = canvas.width;
  // Mobile: legend goes below chart → taller canvas; desktop: legend right
  var isMobile = W < 480;
  canvas.height = isMobile ? 540 : 460;
  var ctx = canvas.getContext('2d');
  var H = canvas.height;
  ctx.clearRect(0, 0, W, H);
  ctx.fillStyle = '#070B16'; ctx.fillRect(0, 0, W, H);

  // Mobile layout: no right panel, legend drawn below chart
  // Desktop layout: 252px right legend panel, 82px top header
  var legendW = isMobile ? 0 : 252;
  var pad = isMobile
    ? {l:36, r:8,  t:54, b:44}   // mobile: full-width chart
    : {l:62, r:legendW, t:82, b:62}; // desktop
  var cW = W - pad.l - pad.r;
  var cH = H - pad.t - pad.b;
  function px(f) { return pad.l + cW * f; }
  function py(f) { return pad.t + cH * (1 - f); }
  function rr(x,y,w,h,r) { if(ctx.roundRect)ctx.roundRect(x,y,w,h,r); else ctx.rect(x,y,w,h); }

  try {
    var d = await (await fetch('/api/humans')).json();
    var humans = d.humans || [];
    if (humans.length < 2) {
      ctx.fillStyle='rgba(155,114,246,0.6)'; ctx.font='13px Inter'; ctx.textAlign='center';
      ctx.fillText('Need 2+ registered humans', W/2, H/2); return;
    }

    var bals = humans.map(function(h){ return parseFloat(h.balance)||0; }).sort(function(a,b){return a-b;});
    var n = bals.length, total = bals.reduce(function(s,b){return s+b;},0);

    var lorenz = [{x:0,y:0}]; var cum=0;
    for(var i=0;i<n;i++){cum+=bals[i];lorenz.push({x:(i+1)/n,y:total>0?cum/total:(i+1)/n});}

    var area=0;
    for(var i=1;i<lorenz.length;i++){area+=(lorenz[i].x-lorenz[i-1].x)*(lorenz[i].y+lorenz[i-1].y)/2;}
    var gini=Math.max(0,1-2*area);
    // Apply same small-sample bias correction as Go's calcGiniFromBalances: gini * n/(n-1)
    // Without this the Lorenz Gini differs from the Score Gini by factor n/(n-1).
    // At n=7: 0.0841 * 7/6 = 0.0981 — matching the server value.
    if(n>1) gini=Math.min(1, gini * n/(n-1));

    var gEl=document.getElementById('lorenz-gini-val');
    if(gEl){gEl.textContent=gini.toFixed(4);gEl.style.color=gini<0.30?'#34D399':'#F0B429';}

    // Interpolate at exactly x=0.5 between the two bracketing Lorenz points
    // (nearest-point snap was biased by data density near 50%).
    var aqY50 = (function(){
      for(var i=1;i<lorenz.length;i++){
        if(lorenz[i].x>=0.5){
          var t=(0.5-lorenz[i-1].x)/(lorenz[i].x-lorenz[i-1].x);
          return lorenz[i-1].y+t*(lorenz[i].y-lorenz[i-1].y);
        }
      }
      return lorenz[lorenz.length-1].y;
    })();
    var gC = gini<0.30?'#34D399':'#F0B429';

    // ── HEADER ─────────────────────────────────────────────────────────────
    if(isMobile) {
      // Mobile: compact single-line header + one info bar
      ctx.fillStyle='rgba(232,237,245,0.85)'; ctx.font='bold 10px Inter'; ctx.textAlign='left';
      ctx.fillText('LORENZ CURVE', pad.l, 13);
      ctx.fillStyle='rgba(136,146,164,0.55)'; ctx.font='8px Inter';
      ctx.fillText('Diagonal = perfect equality. Below = inequality.', pad.l, 25);
      // Single compact bar: Aequitas vs World
      var barW = W - pad.l - pad.r - 2;
      ctx.fillStyle='rgba(7,11,22,0.97)'; ctx.strokeStyle=gC; ctx.lineWidth=1;
      ctx.beginPath(); rr(pad.l, 30, barW, 20, 4); ctx.fill(); ctx.stroke();
      ctx.font='bold 9px JetBrains Mono'; ctx.textAlign='left';
      ctx.fillStyle=gC; ctx.fillText('Aequitas: '+gini.toFixed(4), pad.l+8, 43);
      ctx.fillStyle='rgba(167,139,250,0.85)'; ctx.fillText('| World avg: 0.38', pad.l+100, 43);
    } else {
      // Desktop: full title + two info boxes
      ctx.fillStyle='rgba(232,237,245,0.88)'; ctx.font='bold 11px Inter'; ctx.textAlign='left';
      ctx.fillText('LORENZ CURVE — WEALTH DISTRIBUTION', pad.l, 14);
      ctx.fillStyle='rgba(136,146,164,0.6)'; ctx.font='8.5px Inter';
      ctx.fillText('Diagonal = perfect equality.  Curves bowing down = more inequality.  Shaded area = size of inequality gap.', pad.l, 27);
      var bw=Math.min(180, Math.floor((cW - 12) / 2)), bh=40;
      ctx.fillStyle='rgba(7,11,22,0.97)'; ctx.strokeStyle=gC; ctx.lineWidth=1.5;
      ctx.beginPath(); rr(pad.l, 34, bw, bh, 5); ctx.fill(); ctx.stroke();
      ctx.fillStyle='rgba(136,146,164,0.6)'; ctx.font='7px JetBrains Mono'; ctx.textAlign='center';
      ctx.fillText('AEQUITAS GINI COEFFICIENT', pad.l+bw/2, 46);
      ctx.fillStyle=gC; ctx.font='bold 17px JetBrains Mono';
      ctx.fillText(gini.toFixed(4), pad.l+58, 65);
      ctx.fillStyle='rgba(200,200,200,0.65)'; ctx.font='9px JetBrains Mono'; ctx.textAlign='left';
      ctx.fillText('= '+gini.toFixed(4), pad.l+105, 65);
      var b2x=pad.l+bw+12;
      ctx.fillStyle='rgba(7,11,22,0.97)'; ctx.strokeStyle='rgba(167,139,250,0.7)'; ctx.lineWidth=1.5;
      ctx.beginPath(); rr(b2x, 34, bw, bh, 5); ctx.fill(); ctx.stroke();
      ctx.fillStyle='rgba(136,146,164,0.6)'; ctx.font='7px JetBrains Mono'; ctx.textAlign='center';
      ctx.fillText('WORLD AVERAGE GINI 2024', b2x+bw/2, 46);
      ctx.fillStyle='rgba(167,139,250,0.9)'; ctx.font='bold 17px JetBrains Mono';
      ctx.fillText('38.0%', b2x+58, 65);
      ctx.fillStyle='rgba(200,200,200,0.65)'; ctx.font='9px JetBrains Mono'; ctx.textAlign='left';
      ctx.fillText('= 0.380', b2x+108, 65);
    }

    // ── GRID ──────────────────────────────────────────────────────────────
    ctx.strokeStyle='rgba(255,255,255,0.04)'; ctx.lineWidth=1;
    for(var i=1;i<4;i++){
      ctx.beginPath();ctx.moveTo(pad.l,py(i/4));ctx.lineTo(pad.l+cW,py(i/4));ctx.stroke();
      ctx.beginPath();ctx.moveTo(px(i/4),pad.t);ctx.lineTo(px(i/4),pad.t+cH);ctx.stroke();
    }

    // ── AXIS ──────────────────────────────────────────────────────────────
    var axFontSz = isMobile ? 8 : 10;
    ctx.fillStyle='rgba(136,146,164,0.7)'; ctx.font=axFontSz+'px JetBrains Mono';
    // On mobile only show 0%, 50%, 100% to save space
    var tl = isMobile ? ['0%','50%','100%'] : ['0%','25%','50%','75%','100%'];
    var tlIdx = isMobile ? [0,2,4] : [0,1,2,3,4];
    for(var i=0;i<tl.length;i++){
      ctx.textAlign='center'; ctx.fillText(tl[i],px(tlIdx[i]/4),pad.t+cH+16);
      ctx.textAlign='right';  ctx.fillText(tl[i],pad.l-(isMobile?4:6),py(tlIdx[i]/4)+4);
    }
    if(!isMobile) {
      ctx.save();ctx.translate(12,pad.t+cH/2);ctx.rotate(-Math.PI/2);
      ctx.fillStyle='rgba(155,114,246,0.7)';ctx.font='10px Inter';ctx.textAlign='center';
      ctx.fillText('Cumulative % of AEQ wealth',0,0);ctx.restore();
      ctx.fillStyle='rgba(155,114,246,0.6)';ctx.font='10px Inter';ctx.textAlign='center';
      ctx.fillText('% of Population (poorest left → richest right)',px(0.5),pad.t+cH+36);
    } else {
      ctx.fillStyle='rgba(155,114,246,0.5)';ctx.font='8px Inter';ctx.textAlign='center';
      ctx.fillText('Population % →',px(0.5),pad.t+cH+30);
    }

    // ── 50% GUIDE LINE ─────────────────────────────────────────────────────
    ctx.beginPath();ctx.moveTo(px(0.5),pad.t);ctx.lineTo(px(0.5),pad.t+cH);
    ctx.strokeStyle='rgba(255,255,255,0.09)';ctx.lineWidth=1;ctx.setLineDash([4,3]);ctx.stroke();ctx.setLineDash([]);
    ctx.fillStyle='rgba(136,146,164,0.45)';ctx.font='8px JetBrains Mono';ctx.textAlign='center';
    ctx.fillText('50% mark',px(0.5),pad.t-4);

    // ── REFERENCE COUNTRIES (most unequal first → fills stack correctly) ───
    var refs = [
      {label:'South Africa', g:0.63, lc:'#F87171', fc:'rgba(239,68,68,0.18)', tag:'Extreme inequality'},
      {label:'Brazil',       g:0.53, lc:'#FB923C', fc:'rgba(251,146,60,0.14)', tag:'High inequality'},
      {label:'USA',          g:0.41, lc:'#FCD34D', fc:'rgba(252,211,77,0.11)', tag:'Moderate'},
      {label:'World Avg',    g:0.38, lc:'#A78BFA', fc:'rgba(167,139,250,0.09)', tag:'Global average'},
      {label:'Germany',      g:0.31, lc:'#34D399', fc:'rgba(52,211,153,0.08)', tag:'Low inequality'},
      {label:'Scandinavia',  g:0.27, lc:'#60A5FA', fc:'rgba(96,165,250,0.07)', tag:'Very low — target'}
    ];

    refs.forEach(function(ref){
      var rpts=[];
      for(var j=0;j<=120;j++){var xf=j/120;rpts.push({x:xf,y:Math.pow(xf,1+2*ref.g)});}
      ctx.beginPath();ctx.moveTo(px(0),py(0));
      rpts.forEach(function(p){ctx.lineTo(px(p.x),py(p.y));});
      for(var j=120;j>=0;j--){ctx.lineTo(px(j/120),py(j/120));}
      ctx.closePath();ctx.fillStyle=ref.fc;ctx.fill();

      ctx.beginPath();
      rpts.forEach(function(p,i){if(i===0)ctx.moveTo(px(p.x),py(p.y));else ctx.lineTo(px(p.x),py(p.y));});
      ctx.strokeStyle=ref.lc;
      ctx.lineWidth=ref.label==='World Avg'?1.9:1.2;
      ctx.setLineDash(ref.label==='World Avg'?[7,3]:[5,4]);ctx.stroke();ctx.setLineDash([]);
    });

    // ── EQUALITY DIAGONAL ──────────────────────────────────────────────────
    var diag=ctx.createLinearGradient(px(0),py(0),px(1),py(1));
    diag.addColorStop(0,'rgba(155,114,246,0.9)');diag.addColorStop(1,'rgba(34,211,238,0.9)');
    ctx.beginPath();ctx.moveTo(px(0),py(0));ctx.lineTo(px(1),py(1));
    ctx.strokeStyle=diag;ctx.lineWidth=2;ctx.setLineDash([8,5]);ctx.stroke();ctx.setLineDash([]);

    // ── AEQUITAS CURVE ─────────────────────────────────────────────────────
    ctx.beginPath();
    lorenz.forEach(function(p,i){if(i===0)ctx.moveTo(px(p.x),py(p.y));else ctx.lineTo(px(p.x),py(p.y));});
    for(var j=lorenz.length-1;j>=0;j--){ctx.lineTo(px(lorenz[j].x),py(lorenz[j].x));}
    ctx.closePath();
    var aqFill=ctx.createLinearGradient(0,py(0.5),0,py(0));
    aqFill.addColorStop(0,'rgba(240,180,41,0.48)');aqFill.addColorStop(1,'rgba(240,180,41,0.04)');
    ctx.fillStyle=aqFill;ctx.fill();

    ctx.beginPath();
    lorenz.forEach(function(p,i){if(i===0)ctx.moveTo(px(p.x),py(p.y));else ctx.lineTo(px(p.x),py(p.y));});
    ctx.save();ctx.shadowColor='rgba(240,180,41,0.8)';ctx.shadowBlur=12;
    ctx.strokeStyle='#F0B429';ctx.lineWidth=3;ctx.stroke();ctx.restore();
    lorenz.slice(1).forEach(function(p){
      ctx.beginPath();ctx.arc(px(p.x),py(p.y),4,0,2*Math.PI);
      ctx.fillStyle='#F0B429';ctx.fill();
      ctx.strokeStyle='rgba(0,0,0,0.6)';ctx.lineWidth=1;ctx.stroke();
    });

    // ── LEGEND ────────────────────────────────────────────────────────────
    var legendItems = [
      {label:'Aequitas',         gStr:gini.toFixed(4), color:'#F0B429', bold:true},
      {label:'Perfect Equality', gStr:'0.00',          color:'rgba(155,114,246,0.9)', bold:false}
    ];
    refs.slice().sort(function(a,b){return a.g-b.g;}).forEach(function(ref){
      legendItems.push({label:ref.label, gStr:ref.g.toFixed(2), color:ref.lc, bold:false});
    });

    // Dots at x=50% in chart (both mobile and desktop)
    legendItems.forEach(function(item){
      var dotY;
      if(item.bold) { dotY = py(aqY50); }
      else if(item.label==='Perfect Equality') { dotY = py(0.5); }
      else {
        var rm = refs.filter(function(r){return r.label===item.label;})[0];
        dotY = rm ? py(Math.pow(0.5,1+2*rm.g)) : null;
      }
      if(dotY != null) {
        ctx.beginPath(); ctx.arc(px(0.5), dotY, item.bold?5:3, 0, 2*Math.PI);
        ctx.fillStyle=item.color; ctx.fill();
        if(item.bold){ctx.strokeStyle='rgba(0,0,0,0.7)';ctx.lineWidth=1;ctx.stroke();}
      }
    });

    if(isMobile) {
      // ── MOBILE LEGEND: compact 2-column grid below chart ──────────────
      var legTop = pad.t + cH + 44;
      var colW = Math.floor((W - pad.l - pad.r) / 2);
      var rowH = 26;
      legendItems.forEach(function(item, idx){
        var col = idx % 2, row = Math.floor(idx / 2);
        var lx2 = pad.l + col * colW;
        var ly2 = legTop + row * rowH;
        // color dot
        ctx.beginPath(); ctx.arc(lx2+6, ly2+7, 5, 0, 2*Math.PI);
        ctx.fillStyle = item.color; ctx.fill();
        // label
        ctx.fillStyle = item.bold ? item.color : 'rgba(232,237,245,0.85)';
        ctx.font = (item.bold ? 'bold ' : '') + '9px Inter';
        ctx.textAlign='left';
        ctx.fillText(item.label, lx2+16, ly2+8);
        // gini
        ctx.fillStyle = 'rgba(136,146,164,0.7)';
        ctx.font = '8.5px JetBrains Mono';
        ctx.fillText('G='+item.gStr, lx2+16, ly2+19);
      });
    } else {
      // ── DESKTOP LEGEND: stacked right panel ───────────────────────────
      var lx = pad.l + cW + 14;
      var lw = pad.r - 20;
      var itemH = Math.min(40, cH / legendItems.length);
      var totalH = itemH * legendItems.length;
      var startY = pad.t + (cH - totalH) / 2 + itemH / 2;
      legendItems.forEach(function(item, idx){
        var cy = startY + idx * itemH;
        ctx.globalAlpha = item.bold ? 1.0 : 0.85;
        ctx.fillStyle = item.color;
        ctx.fillRect(lx, cy - Math.min(itemH*0.38,14), 3, Math.min(itemH*0.76,28));
        ctx.globalAlpha = 1.0;
        ctx.fillStyle = item.color;
        ctx.font = (item.bold?'bold ':'')+' 11px Inter'; ctx.textAlign='left';
        ctx.fillText(item.label, lx+9, cy-2);
        ctx.fillStyle = item.bold ? item.color : 'rgba(232,237,245,0.88)';
        ctx.font = (item.bold?'bold ':'')+' 11.5px JetBrains Mono';
        ctx.fillText('G='+item.gStr, lx+9, cy+11);
        if(itemH>=32){
          ctx.fillStyle='rgba(136,146,164,0.5)'; ctx.font='8px Inter';
          var rm2 = refs.filter(function(r){return r.label===item.label;})[0];
          var owns = item.bold ? '50% own '+(aqY50*100).toFixed(1)+'%'
            : item.label==='Perfect Equality' ? '50% own 50%'
            : rm2 ? '50% own '+Math.round(Math.pow(0.5,1+2*rm2.g)*100)+'%' : '';
          if(owns) ctx.fillText(owns, lx+9, cy+22);
        }
      });
    }

    // ── BOTTOM NOTE ────────────────────────────────────────────────────────
    var noteY = pad.t + cH + 50;
    if(noteY < H - 4) {
      ctx.fillStyle = gini<0.10 ? 'rgba(52,211,153,0.8)' : 'rgba(136,146,164,0.5)';
      ctx.font = (gini<0.10?'bold ':'') + 'italic 8.5px Inter'; ctx.textAlign='center';
      var noteText = gini<0.10
        ? 'Aequitas Gini '+gini.toFixed(4)+' — 4.5x below world average (0.38) — near-perfect equality!'
        : 'Aequitas target: Gini < 0.30  ·  World average: 0.38  •  World average: 38%';
      ctx.fillText(noteText, px(0.5), noteY);
    }

  } catch(e){ console.error('Lorenz error:',e); }
}


function drawWcapSlideChart() {
  const canvas = document.getElementById('wcap-slide-chart');
  if (!canvas || !canvas.offsetParent) return;
  canvas.width = canvas.offsetWidth;
  const ctx = canvas.getContext('2d');
  const W = canvas.width, H = canvas.height;
  ctx.clearRect(0,0,W,H);
  const pad = {l:44,r:20,t:36,b:32};
  const cW = W-pad.l-pad.r, cH = H-pad.t-pad.b;
  const maxN = 28;
  const bw = cW/maxN;
  // horizontal reference lines
  [5,10,15,20,25].forEach(function(v){
    var y=H-pad.b-cH*(v/25);
    ctx.strokeStyle=v===25?'rgba(0,255,209,0.2)':'rgba(139,92,246,0.08)'; ctx.lineWidth=1;
    ctx.beginPath(); ctx.moveTo(pad.l,y); ctx.lineTo(W-pad.r,y); ctx.stroke();
    ctx.fillStyle='rgba(200,168,76,0.7)'; ctx.font='10px JetBrains Mono,monospace'; ctx.textAlign='right';
    ctx.fillText(v+'x', pad.l-5, y+4);
  });
  // bars
  for (var n=1;n<=maxN;n++) {
    var mult=Math.max(5,Math.min(n,25));
    var bh=cH*(mult/25), bx=pad.l+(n-1)*bw+1, bw2=bw-2;
    var y=H-pad.b-bh, r=Math.min(3,bw2/2);
    var barGrad;
    if (n>25) { barGrad='rgba(255,255,255,0.06)'; }
    else if (n===25) { var g=ctx.createLinearGradient(0,y,0,H-pad.b); g.addColorStop(0,'rgba(0,255,209,0.8)'); g.addColorStop(1,'rgba(0,255,209,0.25)'); barGrad=g; }
    else if (n>=20) { var g2=ctx.createLinearGradient(0,y,0,H-pad.b); g2.addColorStop(0,'rgba(200,168,76,0.85)'); g2.addColorStop(1,'rgba(200,168,76,0.28)'); barGrad=g2; }
    else { var g3=ctx.createLinearGradient(0,y,0,H-pad.b); g3.addColorStop(0,'rgba(200,168,76,0.6)'); g3.addColorStop(1,'rgba(200,168,76,0.18)'); barGrad=g3; }
    // rounded top bar
    ctx.beginPath();
    ctx.moveTo(bx+r,y); ctx.lineTo(bx+bw2-r,y);
    ctx.arcTo(bx+bw2,y,bx+bw2,y+r,r);
    ctx.lineTo(bx+bw2,H-pad.b); ctx.lineTo(bx,H-pad.b); ctx.lineTo(bx,y+r);
    ctx.arcTo(bx,y,bx+r,y,r); ctx.closePath();
    if (n===25){ctx.save();ctx.shadowColor='rgba(0,255,209,0.55)';ctx.shadowBlur=8;}
    ctx.fillStyle=barGrad; ctx.fill();
    if (n===25) ctx.restore();
    // labels at key N values
    if (n===1||n===5||n===10||n===15||n===20||n===25) {
      ctx.fillStyle=n===25?'rgba(0,255,209,0.9)':'rgba(200,168,76,0.85)';
      ctx.font='bold 9px JetBrains Mono,monospace'; ctx.textAlign='center';
      ctx.fillText(mult+'x', bx+bw2/2, y-4);
      ctx.fillStyle='rgba(255,255,255,0.4)'; ctx.font='8px JetBrains Mono,monospace';
      ctx.fillText('N='+n, bx+bw2/2, H-pad.b+13);
    }
  }
  // lock line at N=25
  var lockY=H-pad.b-cH;
  ctx.save(); ctx.shadowColor='rgba(0,255,209,0.5)'; ctx.shadowBlur=5;
  ctx.strokeStyle='rgba(0,255,209,0.55)'; ctx.lineWidth=1.5; ctx.setLineDash([5,4]);
  ctx.beginPath(); ctx.moveTo(pad.l+(25-1)*bw,lockY); ctx.lineTo(W-pad.r,lockY); ctx.stroke();
  ctx.setLineDash([]); ctx.restore();
  ctx.fillStyle='rgba(0,255,209,0.8)'; ctx.font='bold 9px JetBrains Mono,monospace'; ctx.textAlign='left';
  ctx.fillText('LOCKED AT 25x', pad.l+25*bw+4, lockY-4);
  // title
  ctx.fillStyle='rgba(200,168,76,0.35)'; ctx.font='10px Inter,sans-serif'; ctx.textAlign='left';
  ctx.fillText('WEALTH CAP  —  BOOTSTRAP MULTIPLIER  ·  max(5, min(N, 25))×', pad.l, 20);
}

function drawPriceChart() {
  const canvas = document.getElementById('price-chart');
  if (!canvas || !priceHistory.length || !canvas.offsetParent) return;
  canvas.width = canvas.offsetWidth;
  const ctx = canvas.getContext('2d');
  const W = canvas.width, H = canvas.height;
  ctx.clearRect(0, 0, W, H);
  const pad = {l:58, r:24, t:36, b:36};
  var now = Date.now();
  var ciMs = (typeof chartIntervalMs !== 'undefined') ? chartIntervalMs : 0;
  var pts = ciMs > 0
    ? priceHistory.filter(function(p){ return now - p.t <= ciMs; })
    : priceHistory;
  if (!pts.length) {
    ctx.fillStyle='rgba(139,92,246,0.45)'; ctx.font='11px Inter,sans-serif'; ctx.textAlign='center';
    ctx.fillText('No price data in this interval yet — wait a few minutes or select a wider range', W/2, H/2);
    return;
  }
  pts = pts.filter(function(p){return p.p>0;}); const prices = pts.map(function(p){return p.p;});
  const minP = Math.min.apply(null,prices), maxP = Math.max.apply(null,prices);
  const range = maxP - minP || minP * 0.01 || 0.0001;
  const padR = range * 0.1;
  const lo = minP - padR, hi = maxP + padR;
  const cW = W-pad.l-pad.r, cH = H-pad.t-pad.b;
  const toX = function(i){return pad.l + cW * i / Math.max(pts.length - 1, 1);};
  const toY = function(p){return pad.t + cH * (1 - (p - lo) / (hi - lo));};
  // grid
  for (var gi=0;gi<=4;gi++) {
    var gy = pad.t + cH*gi/4;
    ctx.strokeStyle = gi===4?'rgba(139,92,246,0.2)':'rgba(139,92,246,0.08)'; ctx.lineWidth=1;
    ctx.beginPath(); ctx.moveTo(pad.l,gy); ctx.lineTo(W-pad.r,gy); ctx.stroke();
    var gv = hi - (hi-lo)*gi/4;
    ctx.fillStyle='rgba(139,92,246,0.75)'; ctx.font='10px JetBrains Mono,monospace'; ctx.textAlign='right';
    ctx.fillText(gv.toFixed(4), pad.l-5, gy+4);
  }
  // bezier fill
  ctx.beginPath(); ctx.moveTo(toX(0),toY(pts[0].p));
  for (var bi=1;bi<pts.length-1;bi++) {
    var mx=(toX(bi)+toX(bi+1))/2, my=(toY(pts[bi].p)+toY(pts[bi+1].p))/2;
    ctx.quadraticCurveTo(toX(bi),toY(pts[bi].p),mx,my);
  }
  if (pts.length>1) ctx.lineTo(toX(pts.length-1),toY(pts[pts.length-1].p));
  ctx.lineTo(toX(pts.length-1),H-pad.b); ctx.lineTo(toX(0),H-pad.b); ctx.closePath();
  var grad=ctx.createLinearGradient(0,pad.t,0,H-pad.b);
  grad.addColorStop(0,'rgba(139,92,246,0.38)'); grad.addColorStop(0.65,'rgba(139,92,246,0.1)'); grad.addColorStop(1,'rgba(139,92,246,0.01)');
  ctx.fillStyle=grad; ctx.fill();
  // glowing bezier line
  ctx.save(); ctx.shadowColor='rgba(139,92,246,0.7)'; ctx.shadowBlur=12;
  ctx.strokeStyle='#8B5CF6'; ctx.lineWidth=2.5;
  ctx.beginPath(); ctx.moveTo(toX(0),toY(pts[0].p));
  for (var li=1;li<pts.length-1;li++) {
    var mx2=(toX(li)+toX(li+1))/2, my2=(toY(pts[li].p)+toY(pts[li+1].p))/2;
    ctx.quadraticCurveTo(toX(li),toY(pts[li].p),mx2,my2);
  }
  if (pts.length>1) ctx.lineTo(toX(pts.length-1),toY(pts[pts.length-1].p));
  ctx.stroke(); ctx.restore();
  // last price dot
  var lx=toX(pts.length-1), ly=toY(prices[prices.length-1]);
  ctx.save(); ctx.shadowColor='rgba(139,92,246,0.9)'; ctx.shadowBlur=16;
  ctx.beginPath(); ctx.arc(lx,ly,5,0,2*Math.PI); ctx.fillStyle='#8B5CF6'; ctx.fill(); ctx.restore();
  ctx.beginPath(); ctx.arc(lx,ly,2.5,0,2*Math.PI); ctx.fillStyle='#fff'; ctx.fill();
  var pLabel=prices[prices.length-1].toFixed(6)+' tUSD';
  ctx.fillStyle='rgba(139,92,246,0.95)'; ctx.font='bold 11px JetBrains Mono,monospace';
  ctx.textAlign = lx>W*0.75?'right':'left';
  ctx.fillText(pLabel, lx+(lx>W*0.75?-8:8), ly-9);
  // x-axis time labels
  ctx.fillStyle='rgba(139,92,246,0.5)'; ctx.font='9px JetBrains Mono,monospace'; ctx.textAlign='center';
  [0, Math.floor(pts.length/2), pts.length-1].forEach(function(i){
    if (i<0||i>=pts.length) return;
    var dd=new Date(pts[i].t);
    var ts=dd.getHours().toString().padStart(2,'0')+':'+dd.getMinutes().toString().padStart(2,'0')+':'+dd.getSeconds().toString().padStart(2,'0');
    ctx.fillText(ts, toX(i), H-pad.b+16);
  });
  // title
  ctx.fillStyle='rgba(139,92,246,0.38)'; ctx.font='10px Inter,sans-serif'; ctx.textAlign='left';
  ctx.fillText('AEQ / tUSD  —  LIVE PRICE  (x·y = k  AMM)', pad.l, 20);
}

let allBlocks = [];

async function loadBlocks() {
  try {
    const blocks = await (await fetch('/api/blocks')).json();
    const list = document.getElementById('blocks-list');
    if (!blocks || !blocks.length) { list.innerHTML = '<div class="empty">No blocks yet</div>'; return; }
    allBlocks = blocks;
    const blockMap = {};
    blocks.forEach(function(bk){ blockMap[bk.hash] = bk; });
    // Deduplicate by height: keep only the block with the most parents per height.
    // In BlockDAG, multiple nodes produce parallel blocks at the same height.
    // The merge block (most parents) is the canonical one to display.
    const byHeight = {};
    blocks.forEach(function(b) {
      const h = b.height;
      const pc = (b.parent_hashes || []).length;
      if (!byHeight[h] || pc > (byHeight[h].parent_hashes || []).length) byHeight[h] = b;
    });
    const dedupedBlocks = Object.values(byHeight).sort(function(a,b){ return b.height - a.height; });
    document.getElementById('block-count').textContent = dedupedBlocks.length + ' blocks';
    list.innerHTML = dedupedBlocks.map(b => {
      const merge = b.parent_hashes && b.parent_hashes.length > 1;
      const hasTx = b.transactions && b.transactions.length > 0;
      const validator = b.proposer ? short(b.proposer, 6, 4) : '—';
      // For MERGE blocks: show unique parent proposers so all contributing nodes are visible
      let parentNodes = '';
      if (merge && b.parent_hashes) {
        const parentProps = [...new Set(b.parent_hashes
          .map(function(h){ const pb = blockMap[h]; return pb && pb.proposer ? short(pb.proposer, 6, 4) : null; })
          .filter(Boolean))];
        if (parentProps.length > 0) {
          parentNodes = ' ← ' + parentProps.map(function(p){ return '<span style="color:var(--purple)">' + sanitize(p) + '</span>'; }).join(', ');
        }
      }
      return '<div class="block-item" onclick="openBlock(\'' + sanitize(b.hash) + '\')">' +
        '<div class="block-num">#' + b.height + '</div>' +
        '<div><div class="block-hash">' + short(b.hash) +
          (merge ? ' <span class="bm">MERGE</span>' : '') +
          (hasTx ? ' <span class="bt">TX</span>' : '') +
          '</div>' +
          '<div class="block-parents">' + (b.parent_hashes ? b.parent_hashes.length + ' parent(s)' : '') +
          ' · <span style="color:var(--teal)">' + validator + '</span>' + parentNodes + '</div>' +
        '</div>' +
        '<div class="block-right"><div class="block-humans">' + (b.humans || 0) + ' humans</div>' +
        '<div class="block-time">' + timeAgo(b.timestamp) + '</div></div>' +
        '</div>';
    }).join('');
  } catch (e) {}
}

function openBlock(hash) {
  const b = allBlocks.find(x => x.hash === hash);
  if (!b) return;
  document.getElementById('bdc-title').textContent = 'Block #' + b.height;
  const ts = new Date(b.timestamp * 1000);
  // All peer-supplied block fields go through sanitize() before innerHTML
  // to prevent XSS — an authorized validator can sign arbitrary content
  // in parent_hashes, state_root, and proposer.
  const parentList = (b.parent_hashes || []).map(function(h){
    const pb = allBlocks.find(function(x){ return x.hash === h; });
    const pProp = pb && pb.proposer ? ' <span style="color:var(--purple);font-size:0.5rem">(' + short(pb.proposer,6,4) + ')</span>' : '';
    return '<div style="margin-bottom:2px">' + sanitize(h) + pProp + '</div>';
  }).join('') || '—';
  const isMerge = b.parent_hashes && b.parent_hashes.length > 1;
  let html = '';
  html += '<div class="bdc-row"><div class="bdc-k">Height</div><div class="bdc-v">#' + sanitize(String(b.height)) + (b.is_genesis ? ' <span class="bm">GENESIS</span>' : '') + '</div></div>';
  html += '<div class="bdc-row"><div class="bdc-k">Full Hash</div><div class="bdc-v" style="font-size:0.55rem">' + sanitize(b.hash || '') + '</div></div>';
  html += '<div class="bdc-row"><div class="bdc-k">Timestamp</div><div class="bdc-v">' + sanitize(ts.toUTCString()) + '</div></div>';
  html += '<div class="bdc-row"><div class="bdc-k">Node P2P ID</div><div class="bdc-v" style="color:var(--teal);word-break:break-all;font-size:0.54rem">' + sanitize(b.proposer || '—') + '</div></div>';
  html += '<div class="bdc-row"><div class="bdc-k" style="color:var(--muted);font-size:0.54rem">i</div><div class="bdc-v" style="color:var(--muted);font-size:0.52rem">Signing address of the block producer (Ethereum address derived from RELAYER_PRIVATE_KEY)</div></div>';
  html += '<div class="bdc-row"><div class="bdc-k">Humans</div><div class="bdc-v">' + sanitize(String(b.humans || 0)) + '</div></div>';
  html += '<div class="bdc-row"><div class="bdc-k">Type</div><div class="bdc-v">' + (isMerge ? '<span class="bm">MERGE BLOCK</span> &mdash; ' + sanitize(String(b.parent_hashes.length)) + ' parents merged' : 'Standard block &mdash; 1 parent') + '</div></div>';
  html += '<div class="bdc-row"><div class="bdc-k">Parent(s)</div><div class="bdc-v" style="font-size:0.55rem">' + parentList + '</div></div>';
  if (b.state_root) html += '<div class="bdc-row"><div class="bdc-k">State Root</div><div class="bdc-v" style="font-size:0.55rem">' + sanitize(b.state_root) + '</div></div>';
  const txs = b.transactions || [];
  if (txs.length > 0) {
    html += '<div class="bdc-tx-hdr">Transactions (' + txs.length + ')</div>';
    txs.forEach(tx => {
      html += '<div class="bdc-tx">TYPE: ' + sanitize(tx.type || '?') +
        '<br>WALLET: ' + sanitize(tx.wallet || '—') +
        (tx.amount ? '<br>AMOUNT: ' + tx.amount + ' AEQ' : '') +
        (tx.tx_hash ? '<br>TX HASH: ' + sanitize(tx.tx_hash) : '') +
        '</div>';
    });
  } else {
    html += '<div class="bdc-row"><div class="bdc-k">Transactions</div><div class="bdc-v" style="color:var(--muted)">None</div></div>';
  }
  document.getElementById('bdc-content').innerHTML = html;
  document.getElementById('block-detail-overlay').classList.add('open');
  document.body.style.overflow = 'hidden';
}

function closeBlock() {
  document.getElementById('block-detail-overlay').classList.remove('open');
  document.body.style.overflow = '';
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
var priceHistory = [];
var chartIntervalMs = 60000;
var priceHistoryLoaded = false;

// Preload price history from DB so interval buttons show real historical data.
// Fetches the last 4 hours of price snapshots saved after each swap/liquidity.
async function preloadPriceHistory() {
  if (priceHistoryLoaded) return;
  try {
    var d = await (await fetch('/api/price-history?minutes=14400&limit=5000')).json();
    var hist = d.history || [];
    if (hist.length > 0) {
      // Merge DB history with any in-memory points, de-duplicate by timestamp
      var existing = new Set(priceHistory.map(function(p){ return p.t; }));
      hist.forEach(function(pt) {
        if (!existing.has(pt.t)) priceHistory.push({t: pt.t, p: pt.p});
      });
      priceHistory.sort(function(a,b){ return a.t - b.t; });
      priceHistoryLoaded = true;
      drawPriceChart();
    }
  } catch(_) {}
}

function setChartInterval(ms) {
  chartIntervalMs = ms;
  var btnIds = ['ci-1m','ci-5m','ci-30m','ci-1h','ci-4h','ci-all'];
  var btnVals = [60000,300000,1800000,3600000,14400000,0];
  for (var bi = 0; bi < btnIds.length; bi++) {
    var btnEl = document.getElementById(btnIds[bi]);
    if (btnEl) btnEl.className = 'ci-btn' + (btnVals[bi] === ms ? ' ci-active' : '');
  }
  drawPriceChart();
}

function swapLog(msg, type) {
  const el = document.getElementById('swap-log');
  el.innerHTML += '<div><span class="' + (type || 'info') + '">' + msg + '</span></div>';
  el.scrollTop = el.scrollHeight;
}

function sanitize(s) {
  const d = document.createElement('div');
  d.textContent = String(s);
  return d.innerHTML;
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
    const total = (d.reserve_aeq || 0) + (d.reserve_tusd || 0);
    if (total > 0) {
      const aeqPct = (d.reserve_aeq / total * 100).toFixed(1);
      const depthFill = document.getElementById('depth-aeq-fill');
      const aeqPctEl = document.getElementById('depth-aeq-pct');
      const tusdPctEl = document.getElementById('depth-tusd-pct');
      if (depthFill) depthFill.style.width = aeqPct + '%';
      if (aeqPctEl) aeqPctEl.textContent = aeqPct + '%';
      if (tusdPctEl) tusdPctEl.textContent = (100 - parseFloat(aeqPct)).toFixed(1) + '%';
    }
    const desc = document.getElementById('swap-addliq-desc');
    if (desc) {
      desc.textContent = d.reserve_aeq > 0
        ? ('Pool ratio: 1 AEQ ≈ ' + d.price_aeq_in_tusd.toFixed(4) + ' tUSD — match this ratio when depositing')
        : 'Be the first to deposit — your ratio sets the starting price.';
    }
    if (d.reserve_aeq > 0 && d.price_aeq_in_tusd > 0) {
      priceHistory.push({ t: Date.now(), p: d.price_aeq_in_tusd });
      if (priceHistory.length > 1000) priceHistory.shift();
      drawPriceChart();
    }
    updateFeeEstimate();
  } catch (e) {}
}

function setSwapDirection(dir) {
  swapDirection = dir;
  const fromIcon = document.getElementById('swap-from-icon');
  const fromSym = document.getElementById('swap-from-sym');
  const toIcon = document.getElementById('swap-to-icon');
  const toSym = document.getElementById('swap-to-sym');
  if (dir === 'aeq_to_tusd') {
    if (fromIcon) fromIcon.textContent = '🔶'; if (fromSym) fromSym.textContent = 'AEQ';
    if (toIcon) toIcon.textContent = '💵'; if (toSym) toSym.textContent = 'tUSD';
  } else {
    if (fromIcon) fromIcon.textContent = '💵'; if (fromSym) fromSym.textContent = 'tUSD';
    if (toIcon) toIcon.textContent = '🔶'; if (toSym) toSym.textContent = 'AEQ';
  }
  // Sync balance labels in the from/to panels
  const fromBal = document.getElementById('swap-from-bal');
  const toBal = document.getElementById('swap-to-bal');
  if (fromBal) fromBal.textContent = dir === 'aeq_to_tusd' ? (fmt(myAEQBalance) + ' AEQ') : (fmt(myTUSDBalance) + ' tUSD');
  if (toBal) toBal.textContent = dir === 'aeq_to_tusd' ? (fmt(myTUSDBalance) + ' tUSD') : (fmt(myAEQBalance) + ' AEQ');
  updateFeeEstimate();
}

function reverseSwapDir() {
  setSwapDirection(swapDirection === 'aeq_to_tusd' ? 'tusd_to_aeq' : 'aeq_to_tusd');
  document.getElementById('swap-amount').value = '';
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
  const aeqToTusd = swapDirection === 'aeq_to_tusd';
  const unit = aeqToTusd ? 'AEQ' : 'tUSD';
  const outUnit = aeqToTusd ? 'tUSD' : 'AEQ';
  const fee = amt * 0.001;
  const feeEl = document.getElementById('swap-fee-est');
  if (feeEl) feeEl.textContent = fee > 0 ? (fee.toFixed(6) + ' ' + unit) : '—';

  const panel = document.getElementById('swap-details-panel');
  const goBtn = document.getElementById('swap-btn-go');
  const warnEl = document.getElementById('swap-warn');

  if (amt <= 0) {
    if (panel) panel.style.display = 'none';
    warnEl.style.display = 'none';
    const od = document.getElementById('swap-out-est-dex'); if (od) od.textContent = '—';
    if (swapWaddr) goBtn.disabled = false;
    return;
  }
  if (currentPoolAEQ <= 0 || currentPoolTUSD <= 0) {
    if (panel) panel.style.display = 'none';
    warnEl.textContent = '⚠ Pool has no liquidity yet — deposit some below before swapping.';
    warnEl.style.display = 'block';
    if (swapWaddr) goBtn.disabled = true;
    return;
  }
  const est = estimateSwapOutput(amt, aeqToTusd);
  if (est && est.tooLarge) {
    if (panel) panel.style.display = 'none';
    // Binary-search the largest input that stays safely under the
    // reserve, so the warning can suggest a concrete number instead of
    // just saying "too much" — 99% of the output reserve as a safety
    // margin, since the pool could shift slightly before this submits.
    let lo = 0, hi = amt;
    for (let i = 0; i < 30; i++) {
      const mid = (lo + hi) / 2;
      const midEst = estimateSwapOutput(mid, aeqToTusd);
      if (midEst && midEst.amountOut < (aeqToTusd ? currentPoolTUSD : currentPoolAEQ) * 0.99) lo = mid;
      else hi = mid;
    }
    warnEl.innerHTML = '⚠ Too large for current pool liquidity. Try up to ~' + lo.toFixed(4) + ' ' + unit + '.';
    warnEl.style.display = 'block';
    if (swapWaddr) goBtn.disabled = true;
  } else if (est) {
    // Show swap details panel with price impact calculation
    if (panel) {
      panel.style.display = 'block';
      const outEl = document.getElementById('swap-out-est');
      const outDex = document.getElementById('swap-out-est-dex');
      const outStr = est.amountOut.toFixed(6) + ' ' + outUnit;
      if (outEl) outEl.textContent = outStr;
      if (outDex) outDex.textContent = outStr;
      // Price impact = how far execution price deviates from spot price
      const spotPrice = aeqToTusd ? (currentPoolTUSD / currentPoolAEQ) : (currentPoolAEQ / currentPoolTUSD);
      const amtAfterFee = amt - est.fee;
      const execPrice = amtAfterFee > 0 ? est.amountOut / amtAfterFee : 0;
      const impact = spotPrice > 0 ? Math.max(0, (1 - execPrice / spotPrice) * 100) : 0;
      const impEl = document.getElementById('swap-price-impact');
      if (impEl) {
        impEl.textContent = impact.toFixed(2) + '%';
        impEl.style.color = impact < 1 ? 'var(--neon)' : impact < 3 ? 'var(--gold)' : 'var(--red)';
      }
      const rateEl = document.getElementById('swap-rate-display');
      if (rateEl) rateEl.textContent = aeqToTusd
        ? ('1 AEQ = ' + (est.amountOut / amt).toFixed(4) + ' tUSD')
        : ('1 tUSD = ' + (est.amountOut / amt).toFixed(4) + ' AEQ');
      if (impact >= 5) {
        warnEl.innerHTML = '⚠ High price impact (' + impact.toFixed(2) + '%). Consider a smaller amount.';
        warnEl.style.display = 'block';
      } else {
        warnEl.style.display = 'none';
      }
    } else {
      warnEl.innerHTML = 'You will receive ≈ ' + est.amountOut.toFixed(6) + ' ' + outUnit;
      warnEl.style.display = 'block';
    }
    if (swapWaddr) goBtn.disabled = false;
  }
}

async function connectSwapWallet() {
  if (!window.ethereum) {
    const _isMobS = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);
    if (_isMobS) { const _dl = 'https://metamask.app.link/dapp/' + window.location.host; swapLog('🦊 MetaMask nicht gefunden. Mobile: <a href="' + _dl + '" style="color:var(--gold)">In MetaMask App öffnen</a>', 'warn'); } else { swapLog('🦊 MetaMask not found — <a href="https://metamask.io/download/" target="_blank" style="color:var(--gold)">install MetaMask</a>', 'warn'); }
    return;
  }
  try {
    await addToMetaMask();
    const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
    swapWaddr = accounts[0];
    waddr = swapWaddr;
    localStorage.setItem('aeq_wallet', swapWaddr);
    document.getElementById('swap-wbox').style.display = 'block';
    document.getElementById('swap-wadr').textContent = swapWaddr;
    const btn = document.getElementById('swap-btn-conn');
    btn.textContent = swapWaddr.slice(0, 10) + '...' + swapWaddr.slice(-4);
    btn.style.background = 'var(--green)';
    btn.style.color = '#050A14';
    const swapDBtn = document.getElementById('swap-btn-disconnect');
    if (swapDBtn) swapDBtn.style.display = 'block';
    // Sync register tab wallet display
    const regBox = document.getElementById('wbox');
    const regAdr = document.getElementById('wadr');
    const regBtn = document.getElementById('btn-conn');
    const regDBtn = document.getElementById('btn-disconnect');
    if (regBox) regBox.style.display = 'block';
    if (regAdr) regAdr.textContent = swapWaddr;
    if (regBtn) { regBtn.textContent = swapWaddr.slice(0, 10) + '...' + swapWaddr.slice(-4); regBtn.style.background = 'var(--green)'; regBtn.style.color = '#050A14'; }
    if (regDBtn) regDBtn.style.display = 'block';
    await refreshSwapBalances();
    await loadLPPosition();
    document.getElementById('swap-btn-go').disabled = false;
    document.getElementById('swap-btn-faucet').disabled = false;
    document.getElementById('swap-btn-addliq').disabled = false;
    setSwapDirection('aeq_to_tusd');
  } catch (e) {
    swapLog('Connection failed: ' + sanitize(e.message), 'err');
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
    // Update DEX from/to panel balance labels
    const fromBal = document.getElementById('swap-from-bal');
    const toBal = document.getElementById('swap-to-bal');
    if (fromBal) fromBal.textContent = swapDirection === 'aeq_to_tusd' ? (fmt(myAEQBalance) + ' AEQ') : (fmt(myTUSDBalance) + ' tUSD');
    if (toBal) toBal.textContent = swapDirection === 'aeq_to_tusd' ? (fmt(myTUSDBalance) + ' tUSD') : (fmt(myAEQBalance) + ' AEQ');
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
    const floored = Math.floor(myAEQBalance * pct * 1e6) / 1e6;
    document.getElementById('addliq-aeq').value = floored > 0 ? floored : '';
    updateLiquidityRatio('aeq');
  } else {
    const floored = Math.floor(myTUSDBalance * pct * 1e6) / 1e6;
    document.getElementById('addliq-tusd').value = floored > 0 ? floored : '';
    updateLiquidityRatio('tusd');
  }
}

function setSwapPct(pct) {
  const bal = swapDirection === 'aeq_to_tusd' ? myAEQBalance : myTUSDBalance;
  const amt = bal * pct;
  document.getElementById('swap-amount').value = amt > 0 ? amt.toFixed(6) : '';
  updateFeeEstimate();
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
    const nonceResp = await fetch('/api/nonce?wallet=' + swapWaddr);
    const { nonce } = await nonceResp.json();
    const timestamp = Math.floor(Date.now() / 1000);
    const message = 'Aequitas Swap: ' + swapDirection + ' ' + amount.toFixed(8) + ' nonce:' + nonce + ' ts:' + timestamp;
    swapLog('Sign the message in MetaMask to confirm this swap...', 'info');
    const signature = await signMessage(message);

    const resp = await fetch('/api/swap', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ wallet: swapWaddr, direction: swapDirection, amount, nonce, timestamp, signature })
    });
    const data = await resp.json();
    if (data.success) {
      swapLog('✓ Swapped! Received ' + data.amount_out.toFixed(6) + ' ' + (swapDirection === 'aeq_to_tusd' ? 'tUSD' : 'AEQ'), 'ok');
      document.getElementById('swap-bal-aeq').textContent = fmt(data.new_aeq_balance) + ' AEQ';
      document.getElementById('swap-bal-tusd').textContent = fmt(data.new_tusd_balance) + ' tUSD';
      loadPoolStatus();
    } else {
      swapLog('✗ Swap failed: ' + sanitize(data.message), 'err');
    }
  } catch (e) {
    swapLog('✗ Error: ' + sanitize(e.message), 'err');
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
      swapLog('✗ Faucet claim failed: ' + sanitize(data.message), 'err');
      document.getElementById('swap-btn-faucet').disabled = false;
    }
  } catch (e) {
    swapLog('✗ Error: ' + sanitize(e.message), 'err');
    document.getElementById('swap-btn-faucet').disabled = false;
  }
}

// When the pool already has liquidity, typing one amount auto-fills the
// other at the pool's current ratio — matches what AddLiquidity itself
// requires (within 1% tolerance), so users don't have to calculate it
// by hand and then get rejected for a slightly-off ratio.
function updateLiquidityRatio(changed) {
  if (currentPoolAEQ <= 0 || currentPoolTUSD <= 0) return;
  const aeqInput = document.getElementById('addliq-aeq');
  const tusdInput = document.getElementById('addliq-tusd');
  if (changed === 'aeq') {
    const aeq = parseFloat(aeqInput.value || '0');
    if (aeq > 0) tusdInput.value = Math.floor(aeq * (currentPoolTUSD / currentPoolAEQ) * 1e6) / 1e6;
  } else {
    const tusd = parseFloat(tusdInput.value || '0');
    if (tusd > 0) aeqInput.value = Math.floor(tusd * (currentPoolAEQ / currentPoolTUSD) * 1e6) / 1e6;
  }
}

async function doAddLiquidity() {
  if (!swapWaddr) return;
  const amountAEQ = parseFloat(document.getElementById('addliq-aeq').value || '0');
  const amountTUSD = parseFloat(document.getElementById('addliq-tusd').value || '0');
  if (amountAEQ <= 0 || amountTUSD <= 0) { swapLog('Enter both AEQ and tUSD amounts', 'err'); return; }

  document.getElementById('swap-btn-addliq').disabled = true;
  try {
    const nonceRespL = await fetch('/api/nonce?wallet=' + swapWaddr);
    const { nonce: nonce_l } = await nonceRespL.json();
    const timestamp = Math.floor(Date.now() / 1000);
    const message = 'Aequitas Add Liquidity: ' + amountAEQ.toFixed(8) + ' AEQ + ' + amountTUSD.toFixed(8) + ' tUSD nonce:' + nonce_l + ' ts:' + timestamp;
    swapLog('Sign the message in MetaMask to confirm this deposit...', 'info');
    const signature = await signMessage(message);

    const resp = await fetch('/api/add-liquidity', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ wallet: swapWaddr, amount_aeq: amountAEQ, amount_tusd: amountTUSD, nonce: nonce_l, timestamp, signature })
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
      swapLog('✗ Add liquidity failed: ' + sanitize(data.message), 'err');
    }
  } catch (e) {
    swapLog('✗ Error: ' + sanitize(e.message), 'err');
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
  var aeq = myFullWithdrawableAEQ * removePct;
  var tusd = myFullWithdrawableTUSD * removePct;
  var preview = aeq.toFixed(4) + ' AEQ + ' + tusd.toFixed(4) + ' tUSD';
  document.getElementById('lp-remove-preview').textContent = preview;
  // Also update the prominent inline preview
  var inline = document.getElementById('lp-remove-inline');
  if (inline) {
    inline.style.display = removePct > 0 ? 'block' : 'none';
    var aeqEl = document.getElementById('lp-inline-aeq');
    var tusdEl = document.getElementById('lp-inline-tusd');
    if (aeqEl) aeqEl.textContent = fmt(aeq);
    if (tusdEl) tusdEl.textContent = fmt(tusd);
  }
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
    const nonceRespR = await fetch('/api/nonce?wallet=' + swapWaddr);
    const { nonce: nonce_r } = await nonceRespR.json();
    const timestamp = Math.floor(Date.now() / 1000);
    const message = 'Aequitas Remove Liquidity: ' + sharesToBurn.toFixed(8) + ' shares nonce:' + nonce_r + ' ts:' + timestamp;
    swapLog('Sign the message in MetaMask to confirm this withdrawal...', 'info');
    const signature = await signMessage(message);

    const resp = await fetch('/api/remove-liquidity', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ wallet: swapWaddr, shares: sharesToBurn, nonce: nonce_r, timestamp, signature })
    });
    const data = await resp.json();
    if (data.success) {
      swapLog('✓ Removed liquidity: received ' + data.amount_aeq.toFixed(4) + ' AEQ + ' + data.amount_tusd.toFixed(4) + ' tUSD', 'ok');
      await refreshSwapBalances();
      await loadPoolStatus();
      await loadLPPosition();
    } else {
      swapLog('✗ Remove liquidity failed: ' + sanitize(data.message), 'err');
    }
  } catch (e) {
    swapLog('✗ Error: ' + sanitize(e.message), 'err');
  }
  document.getElementById('swap-btn-removeliq').disabled = false;
}

function activateTabFromPath(path) {
  const tabNames = ['register','explorer','index','network','exchange'];
  const parts = (path || '').replace(/^\//, '').split('/');
  let name = parts[0];
  const stabSlug = parts[1] || '';
  // Backwards-compat: /swap -> /exchange
  if (name === 'swap') name = 'exchange';
  if (!name || !tabNames.includes(name)) return;
  let tabEl = null;
  document.querySelectorAll('.tab').forEach(t => {
    if ((t.getAttribute('onclick') || '').includes("'" + name + "'")) tabEl = t;
  });
  if (!tabEl) return;
  document.documentElement.removeAttribute('data-active');
  document.querySelectorAll('.tab-content').forEach(t => {
    t.classList.remove('active');
    t.style.display = '';
  });
  document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
  const tabContent = document.getElementById('tab-' + name);
  if (!tabContent) return;
  tabContent.classList.add('active');
  tabEl.classList.add('active');
  // Activate stab-panel: use URL slug if present, otherwise first panel
  const stabMap = {
    explorer:  {blocks:'sep-blocks', humans:'sep-humans'},
    index:     {score:'eqi-score', distribution:'eqi-lorenz', economy:'eqi-economy', story:'eqi-story'},
    network:   {overview:'net-overview', node:'net-runnode', protocol:'net-protocol'},
    exchange:  {swap:'exch-swap', liquidity:'exch-liquidity'}
  };
  const panels = tabContent.querySelectorAll('.stab-panel');
  const stabs  = tabContent.querySelectorAll('.stab');
  if (panels.length) {
    panels.forEach(p => p.classList.remove('active'));
    stabs.forEach(s => s.classList.remove('active'));
    const targetId = stabSlug && stabMap[name] && stabMap[name][stabSlug];
    const targetEl = targetId ? document.getElementById(targetId) : panels[0];
    if (targetEl) targetEl.classList.add('active');
    // Activate matching stab button
    const stabBtn = targetId
      ? tabContent.querySelector('.stab[onclick*=\"' + targetId + '\"]')
      : stabs[0];
    if (stabBtn) stabBtn.classList.add('active');
    else if (stabs[0]) stabs[0].classList.add('active');
  }
  if (name === 'exchange') { loadPoolStatus(); preloadPriceHistory(); }
  if (name === 'index') {
    setTimeout(function() {
      const active = tabContent.querySelector('.stab-panel.active');
      if (!active) return;
      if (active.id === 'eqi-score') { drawGiniHistoryChart(); drawLorenzCurve(); }
      else if (active.id === 'eqi-economy') drawWcapSlideChart();
    }, 50);
  }
}

document.addEventListener('DOMContentLoaded', () => {
  const amtInput = document.getElementById('swap-amount');
  if (amtInput) amtInput.addEventListener('input', updateFeeEstimate);

  // Activate correct tab on initial load based on URL path
  activateTabFromPath(window.location.pathname);
  // Belt-and-suspenders: re-run after first paint in case any async init
  // (MetaMask restore, pool fetch) overwrites the tab state
  requestAnimationFrame(() => activateTabFromPath(window.location.pathname));
});

// Back/forward navigation: restore the tab that matches the URL
window.addEventListener('popstate', () => activateTabFromPath(window.location.pathname));

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
    fetch('/api/prove/get/' + proofId).then(r => r.json()).then(pd => {
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
    const _isMobC = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);
    if (_isMobC) { const _dl = 'https://metamask.app.link/dapp/' + window.location.host; addLog('🦊 Mobile: <a href="' + _dl + '" style="color:var(--gold)">In MetaMask App öffnen</a>', 'warn'); } else { addLog('🦊 MetaMask not found — <a href="https://metamask.io/download/" target="_blank" style="color:var(--gold)">install MetaMask</a>', 'warn'); }
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
    swapWaddr = waddr;
    localStorage.setItem('aeq_wallet', waddr);
    document.getElementById('wbox').style.display = 'block';
    document.getElementById('wadr').textContent = waddr;
    const btn = document.getElementById('btn-conn');
    btn.textContent = waddr.slice(0, 10) + '...' + waddr.slice(-4);
    btn.style.background = 'var(--green)';
    btn.style.color = '#050A14';
    const dBtn = document.getElementById('btn-disconnect');
    if (dBtn) dBtn.style.display = 'block';

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

    const proveResp = await fetch('/api/prove', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ bio: pendingBioHash, salt: salt, wallet: waddr })
    });
    if (!proveResp.ok) {
      let err = {};
      try { err = await proveResp.json(); } catch(e) { err = { error: 'HTTP ' + proveResp.status }; }
      if (err.registered) {
        addLog('This identity is already registered.', 'ok');
        document.getElementById('btn-reg').disabled = true;
        document.getElementById('btn-reg').textContent = 'ALREADY REGISTERED';
        return;
      }
      addLog('Proof generation failed: ' + sanitize(err.error || 'unknown error'), 'err');
      return;
    }
    proofData = await proveResp.json();
    document.getElementById('pbox').style.display = 'block';
    document.getElementById('pval').textContent = 'Proof ready for ' + waddr.slice(0, 10) + '...';
    document.getElementById('btn-reg').disabled = false;
    document.getElementById('btn-reg').textContent = 'PROOF READY — CLICK TO REGISTER';
    addLog('Proof generated for your wallet. Click REGISTER to continue.', 'ok');
  } catch (e) {
    addLog('Connection failed: ' + sanitize(e.message), 'err');
  }
}

async function connectWallet() {
  if (!window.ethereum) {
    const _isMobW = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);
    if (_isMobW) { const _dl = 'https://metamask.app.link/dapp/' + window.location.host; addLog('🦊 Mobile: <a href="' + _dl + '" style="color:var(--gold)">In MetaMask App öffnen</a>', 'warn'); } else { addLog('🦊 MetaMask not found — <a href="https://metamask.io/download/" target="_blank" style="color:var(--gold)">install MetaMask</a>', 'warn'); }
    return;
  }
  try {
    await addToMetaMask();
    const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
    waddr = accounts[0];
    swapWaddr = waddr;
    localStorage.setItem('aeq_wallet', waddr);
    document.getElementById('wbox').style.display = 'block';
    document.getElementById('wadr').textContent = waddr;
    const btn = document.getElementById('btn-conn');
    btn.textContent = waddr.slice(0, 10) + '...' + waddr.slice(-4);
    btn.style.background = 'var(--green)';
    btn.style.color = '#050A14';
    const dBtn = document.getElementById('btn-disconnect');
    if (dBtn) dBtn.style.display = 'block';
    // Sync swap tab wallet display
    const swapBox = document.getElementById('swap-wbox');
    const swapAdr = document.getElementById('swap-wadr');
    const swapBtn = document.getElementById('swap-btn-conn');
    const swapDBtn = document.getElementById('swap-btn-disconnect');
    if (swapBox) swapBox.style.display = 'block';
    if (swapAdr) swapAdr.textContent = waddr;
    if (swapBtn) { swapBtn.textContent = waddr.slice(0, 10) + '...' + waddr.slice(-4); swapBtn.style.background = 'var(--green)'; swapBtn.style.color = '#050A14'; }
    if (swapDBtn) swapDBtn.style.display = 'block';
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
    addLog('Connection failed: ' + sanitize(e.message), 'err');
  }
}

function copyAddr(id, btn) {
  const addr = document.getElementById(id).textContent;
  if (!addr || addr === '—') return;
  navigator.clipboard.writeText(addr).then(() => {
    const orig = btn.textContent;
    btn.textContent = '✓ Copied';
    setTimeout(() => { btn.textContent = orig; }, 1500);
  });
}

function addLog(msg, type) {
  const el = document.getElementById('rlog');
  if (!el) return;
  el.innerHTML += '<div><span class="' + (type||'info') + '">' + msg + '</span></div>';
  el.scrollTop = el.scrollHeight;
}

async function registerViaBrowser() {
  if (!navigator.credentials || !window.PublicKeyCredential) {
    addLog('WebAuthn not supported in this browser.', 'err');
    return;
  }
  document.getElementById('web-reg-warn').style.display = 'block';
  addLog('Creating device credential (biometric or PIN prompt)...', 'info');
  try {
    const challenge = crypto.getRandomValues(new Uint8Array(32));
    const userId = crypto.getRandomValues(new Uint8Array(16));
    const credential = await navigator.credentials.create({
      publicKey: {
        challenge,
        rp: { name: 'Aequitas', id: window.location.hostname },
        user: { id: userId, name: 'aequitas-user', displayName: 'Aequitas User' },
        pubKeyCredParams: [{ alg: -7, type: 'public-key' }, { alg: -257, type: 'public-key' }],
        timeout: 60000,
        attestation: 'none',
        authenticatorSelection: { userVerification: 'preferred' }
      }
    });
    // Hash credential.rawId bytes into a BigInt, then reduce mod BN254 field size
    const credBytes = new Uint8Array(credential.rawId);
    let bioNum = BigInt(0);
    for (const b of credBytes) bioNum = (bioNum << 8n) | BigInt(b);
    const FIELD_SIZE = BigInt('21888242871839275222246405745257275088548364400416034343698204186575808495617');
    pendingBioHash = (bioNum % FIELD_SIZE).toString();
    addLog('Device identity hashed. Connecting wallet...', 'ok');
    await connectWalletAndProve();
  } catch (e) {
    addLog('WebAuthn error: ' + sanitize(e.message), 'err');
  }
}

async function doRegister() {
  if (!waddr || !proofData) return;
  try {
    addLog('Preparing signature...', 'info');
    document.getElementById('btn-reg').disabled = true;

    // commitment is pubSignals[0] — must match exactly what the contract reads
    const commitment = proofData.pubSignals[0];

    // Nullifier: prefer ZK-circuit-derived (pubSignals[1] from v2 circuit) since
    // it is cryptographically attested by the proof. Fall back to SHA256 for v1.
    let nullifier = '';
    if (proofData.zkNullifier) {
      // v2 circuit: nullifier is pubSignals[1], proven inside the ZK statement.
      const zkN = BigInt(proofData.zkNullifier);
      nullifier = zkN.toString(16).padStart(64, '0');
      addLog('Using ZK-bound nullifier (circuit v2)', 'info');
    } else {
      // v1 fallback: derive from bioHash off-chain (server validates derivation)
      const bioHashForNullifier = pendingBioHash || proofData.bioHashKey || '';
      if (bioHashForNullifier) {
        const enc = new TextEncoder();
        const buf = await crypto.subtle.digest('SHA-256', enc.encode(bioHashForNullifier + ':aequitas-ubi-v1'));
        nullifier = Array.from(new Uint8Array(buf)).map(b => b.toString(16).padStart(2, '0')).join('');
      }
    }
    if (!nullifier) {
      addLog('Error: could not compute nullifier — biometric hash unavailable', 'err');
      document.getElementById('btn-reg').disabled = false;
      return;
    }

    // Build the EXACT same hash the contract computes:
    // keccak256(abi.encodePacked(block.chainid, address(this), "register", commitment, nullifier))
    const messageHash = ethers.solidityPackedKeccak256(
      ['uint256', 'address', 'string', 'uint256', 'bytes32'],
      [1926, V7_CONTRACT, 'register', commitment, '0x' + nullifier]
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
        bioHash: pendingBioHash || '',
        bioHashKey: proofData.bioHashKey || '',
        nullifier: nullifier,
        circuitVersion: proofData.circuitVersion || 1,
        zkNullifier: proofData.zkNullifier || null
      })
    });
    const d = await r.json();
    if (!d.success) { addLog('Error: ' + d.message, 'err'); document.getElementById('btn-reg').disabled = false; return; }
    addLog('Registered! ' + d.message, 'ok');
    setTimeout(() => { window.location.href = '/registered?wallet=' + waddr; }, 1500);
  } catch (e) { addLog('Error: ' + sanitize(e.message), 'err'); document.getElementById('btn-reg').disabled = false; }
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

function disconnectWallet() {
  waddr = '';
  swapWaddr = '';
  localStorage.removeItem('aeq_wallet');
  // Reset register tab
  const wbox = document.getElementById('wbox');
  const wadr = document.getElementById('wadr');
  const bConn = document.getElementById('btn-conn');
  const bDisc = document.getElementById('btn-disconnect');
  const bReg = document.getElementById('btn-reg');
  if (wbox) wbox.style.display = 'none';
  if (wadr) wadr.textContent = '—';
  if (bConn) { bConn.textContent = '🦊 CONNECT METAMASK'; bConn.style.background = ''; bConn.style.color = ''; }
  if (bDisc) bDisc.style.display = 'none';
  if (bReg) { bReg.disabled = true; bReg.textContent = 'REGISTER ON-CHAIN'; }
  // Reset swap tab
  const swapBox = document.getElementById('swap-wbox');
  const swapAdr = document.getElementById('swap-wadr');
  const swapConn = document.getElementById('swap-btn-conn');
  const swapDisc = document.getElementById('swap-btn-disconnect');
  const swapGo = document.getElementById('swap-btn-go');
  if (swapBox) swapBox.style.display = 'none';
  if (swapAdr) swapAdr.textContent = '—';
  if (swapConn) { swapConn.textContent = '🦊 CONNECT METAMASK'; swapConn.style.background = ''; swapConn.style.color = ''; }
  if (swapDisc) swapDisc.style.display = 'none';
  if (swapGo) swapGo.disabled = true;
  addLog('✓ Wallet disconnected locally. To fully revoke, open MetaMask → Connected Sites.', 'info');
}

async function restoreWalletFromStorage() {
  const saved = localStorage.getItem('aeq_wallet');
  if (!saved || !window.ethereum) return;
  try {
    const accounts = await window.ethereum.request({ method: 'eth_accounts' });
    if (accounts && accounts[0] && accounts[0].toLowerCase() === saved.toLowerCase()) {
      waddr = accounts[0];
      swapWaddr = accounts[0];
      // Restore register tab UI
      const wbox = document.getElementById('wbox');
      const wadr = document.getElementById('wadr');
      const bConn = document.getElementById('btn-conn');
      const bDisc = document.getElementById('btn-disconnect');
      if (wbox) wbox.style.display = 'block';
      if (wadr) { wadr.textContent = accounts[0]; wadr.title = accounts[0]; }
      if (bConn) { bConn.textContent = accounts[0].slice(0,10)+'...'+accounts[0].slice(-4); bConn.style.background='var(--green)'; bConn.style.color='#050A14'; }
      if (bDisc) bDisc.style.display = 'block';
      // Restore swap tab UI
      const swapBox = document.getElementById('swap-wbox');
      const swapAdr = document.getElementById('swap-wadr');
      const swapConn = document.getElementById('swap-btn-conn');
      const swapDBtn = document.getElementById('swap-btn-disconnect');
      if (swapBox) swapBox.style.display = 'block';
      if (swapAdr) { swapAdr.textContent = accounts[0]; swapAdr.title = accounts[0]; }
      if (swapConn) { swapConn.textContent = accounts[0].slice(0,10)+'...'+accounts[0].slice(-4); swapConn.style.background='var(--green)'; swapConn.style.color='#050A14'; }
      if (swapDBtn) swapDBtn.style.display = 'block';
      const goBtn = document.getElementById('swap-btn-go');
      const faucetBtn = document.getElementById('swap-btn-faucet');
      const addliqBtn = document.getElementById('swap-btn-addliq');
      if (goBtn) goBtn.disabled = false;
      if (faucetBtn) faucetBtn.disabled = false;
      if (addliqBtn) addliqBtn.disabled = false;
      setSwapDirection('aeq_to_tusd');
      refreshSwapBalances();
      loadLPPosition();
      // Check registration status silently — no popup
      try {
        const br = await fetch('/api/balance?wallet=' + accounts[0]);
        const bd = await br.json();
        if (bd.is_human) {
          const bReg = document.getElementById('btn-reg');
          if (bReg) { bReg.disabled = true; bReg.textContent = 'ALREADY REGISTERED ✓'; }
          addLog('✓ Wallet restored. Balance: ' + (bd.balance || 0).toFixed(4) + ' AEQ · Already registered.', 'ok');
        }
      } catch(_) {}
    } else {
      localStorage.removeItem('aeq_wallet');
    }
  } catch(e) {}
}

checkProofParams();
restoreWalletFromStorage();
loadStatus();
loadBlocks();
loadHumans();
setInterval(loadStatus, 6000);
setInterval(loadBlocks, 6000);
setInterval(loadHumans, 10000);
setInterval(loadPoolStatus, 8000);
// Observe each canvas individually so charts redraw when they become visible.
// We observe the canvas containers, not document.body (which fires on every
// DOM change and would cause constant redraws killing performance).
(function() {
  if (typeof ResizeObserver === 'undefined') return;
  function observeCanvas(canvasId, drawFn) {
    var canvas = document.getElementById(canvasId);
    if (!canvas) return;
    var ro = new ResizeObserver(function(entries) {
      for (var e of entries) {
        if (e.contentRect.width > 0) drawFn();
      }
    });
    ro.observe(canvas);
  }
  observeCanvas('gini-history-chart', drawGiniHistoryChart);
  observeCanvas('lorenz-chart', drawLorenzCurve);
  observeCanvas('price-chart', drawPriceChart);
})();

async function registerValidatorKey() {
  var statusEl = document.getElementById('vk-status');
  var signingAddr = document.getElementById('vk-signing-addr').value.trim().toLowerCase();
  var signingKeySig = document.getElementById('vk-signing-sig').value.trim();
  if (!signingAddr.startsWith('0x') || signingAddr.length !== 42) {
    statusEl.textContent = 'Enter a valid signing address (0x... 42 chars)'; return;
  }
  if (!signingKeySig) {
    statusEl.textContent = 'Enter the signing key signature from the curl command'; return;
  }
  if (!window.ethereum) { statusEl.textContent = 'MetaMask not found'; return; }
  try {
    var accs = await window.ethereum.request({ method: 'eth_requestAccounts' });
    var humanWallet = accs[0].toLowerCase();
    statusEl.textContent = 'Sign with human wallet in MetaMask...';
    statusEl.style.color = 'var(--gold)';
    var humanMsg = 'Aequitas: authorize validator key ' + signingAddr;
    var humanSig = await window.ethereum.request({ method: 'personal_sign', params: [humanMsg, humanWallet] });
    statusEl.textContent = 'Submitting...';
    var resp = await fetch('/api/register-validator-key', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        signing_address: signingAddr,
        human_wallet: humanWallet,
        human_signature: humanSig,
        signing_key_signature: signingKeySig
      })
    });
    var data = await resp.json();
    if (data.success) {
      statusEl.textContent = 'Validator key registered! Node blocks are now accepted.';
      statusEl.style.color = 'var(--teal)';
    } else {
      statusEl.textContent = sanitize(data.error || 'Registration failed');
      statusEl.style.color = '#f87171';
    }
  } catch(e) {
    statusEl.textContent = sanitize(e.message);
    statusEl.style.color = '#f87171';
  }
}

function generateNodeGuidePDF() {
  var lang = curLang || 'en';
  if (window.jspdf) { try { _buildNodeGuidePDF(lang); } catch(e) { alert('PDF-Fehler: ' + e.message); } return; }
  var s = document.createElement('script');
  s.src = 'https://cdnjs.cloudflare.com/ajax/libs/jspdf/2.5.1/jspdf.umd.min.js';
  s.onload = function() { try { _buildNodeGuidePDF(lang); } catch(e) { alert('PDF-Fehler: ' + e.message); } };
  s.onerror = function() { alert('PDF-Bibliothek konnte nicht geladen werden. Bitte Internetverbindung prüfen.'); };
  document.head.appendChild(s);
}

function _buildNodeGuidePDF(lang) {
  var jsPDF = window.jspdf.jsPDF;
  var doc = new jsPDF({orientation:'portrait',unit:'mm',format:'a4'});
  var W=210, MG=18, CW=W-2*MG, y=20;
  function np(){doc.addPage();y=22;hdr();}
  function ck(n){if(y+n>272)np();}
  function hdr(){doc.setFont('helvetica','bold');doc.setFontSize(7);doc.setTextColor(180,150,60);doc.text('AEQUITAS · NODE OPERATOR GUIDE',MG,10);doc.setDrawColor(180,150,60);doc.setLineWidth(0.2);doc.line(MG,12,W-MG,12);}
  function h1(t){ck(18);y+=5;doc.setFont('helvetica','bold');doc.setFontSize(12);doc.setTextColor(139,92,246);var ls=doc.splitTextToSize(t,CW);doc.text(ls,MG,y);doc.setDrawColor(139,92,246);doc.setLineWidth(0.4);doc.line(MG,y+2,MG+CW,y+2);y+=ls.length*7+5;doc.setTextColor(30,30,30);}
  function h2(t){ck(10);doc.setFont('helvetica','bold');doc.setFontSize(9.5);doc.setTextColor(80,80,200);var ls=doc.splitTextToSize(t,CW);doc.text(ls,MG,y);y+=ls.length*6+3;doc.setTextColor(30,30,30);}
  function tx(t){doc.setFont('helvetica','normal');doc.setFontSize(8.5);doc.setTextColor(40,40,40);var ls=doc.splitTextToSize(t,CW);ls.forEach(function(l){ck(6);doc.text(l,MG,y);y+=5.2;});y+=1.5;}
  function cd(t){var ls=t.split('\n'),lh=4.8,bh=ls.length*lh+8;ck(bh+4);doc.setFillColor(8,10,22);doc.setDrawColor(80,50,180);doc.setLineWidth(0.3);doc.roundedRect(MG,y,CW,bh,2,2,'FD');doc.setFont('courier','normal');doc.setFontSize(7);doc.setTextColor(0,220,170);ls.forEach(function(l,i){doc.text(l,MG+4,y+6+i*lh);});y+=bh+4;doc.setFont('helvetica','normal');doc.setTextColor(40,40,40);}
  function tbl(hdrs,rows,cws){var nC=hdrs.length;if(!cws){cws=[];for(var i=0;i<nC;i++)cws.push(CW/nC);}var lh=7,needH=lh+rows.length*lh+4;ck(needH);doc.setFillColor(25,15,70);doc.rect(MG,y,CW,lh,'F');doc.setFont('helvetica','bold');doc.setFontSize(7.5);doc.setTextColor(255,255,255);var x0=MG;hdrs.forEach(function(h,i){doc.text(h,x0+3,y+5);x0+=cws[i];});y+=lh;rows.forEach(function(row,ri){doc.setFillColor(ri%2===0?243:250,ri%2===0?241:249,ri%2===0?255:255);doc.rect(MG,y,CW,lh,'F');doc.setFont('helvetica','normal');doc.setFontSize(7.2);doc.setTextColor(20,20,50);var x=MG;row.forEach(function(cell,ci){var wrapped=doc.splitTextToSize(String(cell||''),cws[ci]-4);doc.text(wrapped[0]||'',x+3,y+5);x+=cws[ci];});y+=lh;});y+=3;doc.setTextColor(40,40,40);}
  function bl(items){doc.setFont('helvetica','normal');doc.setFontSize(8.5);doc.setTextColor(40,40,40);items.forEach(function(item){ck(7);var ls=doc.splitTextToSize('• '+item,CW-3);ls.forEach(function(l,i){doc.text(l,MG+(i>0?4:2),y);y+=5;});});y+=2;}
  var C={
    en:{title:'Aequitas Node Operator Guide',sub:'Complete step-by-step guide · Aequitas Chain (Chain ID 1926)',badge:'v1.0 · Open Source · Permissionless · No stake required',
      s1:'1. Overview',what:'What a node does',wtxt:'An Aequitas node participates fully in the network: produces blocks in the BlockDAG consensus, validates Groth16 zero-knowledge biometric proofs for new human registrations, enforces wealth caps and demurrage at protocol level, syncs state with peers via libp2p + HTTP, and optionally runs daily pool distributions. Every node runs the full chain — there are no light clients.',
      earn:'What you earn',etxt:'Set NODE_OPERATOR_WALLET to a registered human wallet. The Validators Pool accumulates 40% of all protocol fees (swap fees, demurrage, wealth cap overflow). Every 24 h the primary node distributes the pool balance proportionally among all registered node operator wallets. No stake required — block production is fully permissionless.',
      s2:'2. Requirements',rh:['Component','Minimum','Recommended'],rr:[['OS','Linux / Docker-capable host','Ubuntu 22.04 LTS'],['RAM','512 MB','1 GB (EVM needs headroom)'],['CPU','1 vCPU','2 vCPU (Groth16 is CPU-bound)'],['Storage','2 GB','10 GB SSD (chain grows over time)'],['Database','PostgreSQL 14+','Railway or Supabase (managed)'],['Network','Public IP / port forward','TCP 8080 open, stable uptime']],
      s3:'3. Environment Variables',e3:'Set these before starting the node. Variables marked YES are required; "For rewards" is needed to earn validator payouts.',eh:['Variable','Purpose','Required?'],er:[['DATABASE_URL','PostgreSQL connection string: postgres://user:pass@host:5432/aequitas','YES'],['RELAYER_PRIVATE_KEY','Private key (0x...) of the EOA that signs on-chain human registrations','YES'],['NODE_OPERATOR_WALLET','Registered human wallet address that receives daily validator pool rewards','For rewards'],['RELAYER_ADDRESS','EOA address matching RELAYER_PRIVATE_KEY. Has a hardcoded fallback but set explicitly.','Recommended'],['PORT','HTTP port for API + JSON-RPC. Default: 8080','NO'],['PEER_SECRET','Shared secret authorising this node as validator. ALL nodes must use the SAME value. Get it from the network operator.','For multi-node'],['SELF_URL','This node public HTTPS URL (e.g. https://my-node.up.railway.app). Required for peer discovery self-exclusion.','For multi-node'],['PRIMARY_NODE_URL','Primary node URL for automatic peer discovery. Set to https://aequitas.digital','For multi-node'],['BOOTSTRAP_SNAPSHOT_URL','Set to https://aequitas.digital/api/snapshot on a fresh node. Downloads and imports full network state at startup (only when local DB has 0 humans). Fixes StateRoot mismatch immediately.','Multi-node'],['BOOTSTRAP_SIGNER','Ethereum signing address of the primary node. Get from https://aequitas.digital/api/status → signing_address. Required when BOOTSTRAP_SNAPSHOT_URL is set.','Multi-node'],['SNAPSHOT_TOKEN','Bearer token for /api/snapshot authentication. Must match the primary node\'s SNAPSHOT_TOKEN. Get from network operator.','Multi-node'],['PEER_NODES','Static peer URLs (legacy). Use PRIMARY_NODE_URL for auto-discovery.','Optional'],['NODE_KEY','Base64-encoded libp2p private key for stable peer identity. If not set: auto-generated and printed to stderr as "SAVE THIS AS NODE_KEY ENVIRONMENT VAR: <base64>". Copy and set it.','NO'],['IS_PRIMARY_NODE','"true" only on the designated primary. DEPRECATED: distributions now use a DB-level lock (TryLockDistribution) — not required. Safe to leave unset.','NO'],['RESET_STATE','"true" wipes the database on startup. DESTRUCTIVE — never use in production.','NO']],
      s4:'4. Quick Start — Railway (Recommended)',r4:'Railway is the fastest way to get running. The free tier meets all requirements. Estimated setup time: 10–15 minutes.',rs:['Fork the repo: https://github.com/hanoi96international-gif/Aequitas','Create a Railway account at railway.app and start a new project','Click "Deploy from GitHub Repo" and select your fork','In the project: + New → Database → Add PostgreSQL — DATABASE_URL is auto-set by Railway','Go to your service → Settings → Variables and add the env vars from Section 3','Set PRIMARY_NODE_URL=https://aequitas.digital so your node syncs from the primary','Set NODE_OPERATOR_WALLET=<your registered AEQ human wallet> to receive daily validator rewards','Set RELAYER_PRIVATE_KEY=<your EOA private key> for signing on-chain registrations','Click "Deploy" — the Dockerfile in the repo root handles the build (~3 min for Go compilation)','Watch the deploy logs for: "Aequitas Node Running ✓" and "[NODE] Registered node operator wallet"','Open YOUR-RAILWAY-URL/api/status to confirm the node is live and block height is climbing','Add your node\'s RPC to MetaMask: Chain ID 1926, Symbol AEQ, URL https://YOUR-URL/rpc'],rn:'Railway assigns a random subdomain; custom domains can be set in project settings. Only port 8080 needs to be exposed — P2P is managed internally by the node.',
      s5:'5. Quick Start — Docker / VPS',d5:'For VPS (Contabo, Hetzner, DigitalOcean). Docker + PostgreSQL required. You can reuse Railway\'s PostgreSQL as external DB. NODE_OPERATOR_WALLET must be a registered Aequitas human.',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\ndocker build -t aequitas-node .\n\n# Step 1: First start (no NODE_KEY yet — it will be shown in logs)\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgresql://user:pass@host:5432/dbname" \\\n  -e RELAYER_PRIVATE_KEY="YOUR_HEX_PRIVATE_KEY" \\\n  -e RELAYER_ADDRESS="0xYOUR_NODE_SIGNING_ADDRESS" \\\n  -e NODE_OPERATOR_WALLET="0xYOUR_REGISTERED_HUMAN_WALLET" \\\n  -e PEER_SECRET="get-from-network-operator" \\\n  -e SELF_URL="http://YOUR-SERVER-IP:8080" \\\n  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n  -e BOOTSTRAP_SNAPSHOT_URL="https://aequitas.digital/api/snapshot" \\\n  -e BOOTSTRAP_SIGNER="0x92cbedec9d348b4762cb9af99500ee6139c5b671" \\\n  -e SNAPSHOT_TOKEN="get-from-network-operator" \\\n  -p 8080:8080 -p 4001:4001 aequitas-node\n\n# Step 2: Get NODE_KEY from logs (one-time setup):\ndocker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY"\n\n# Step 3: Final command with NODE_KEY (use this permanently):\ndocker stop aequitas-node && docker rm aequitas-node\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgresql://user:pass@host:5432/dbname" \\\n  -e RELAYER_PRIVATE_KEY="YOUR_HEX_PRIVATE_KEY" \\\n  -e RELAYER_ADDRESS="0xYOUR_NODE_SIGNING_ADDRESS" \\\n  -e NODE_OPERATOR_WALLET="0xYOUR_REGISTERED_HUMAN_WALLET" \\\n  -e NODE_KEY="base64-from-step-2" \\\n  -e PEER_SECRET="get-from-network-operator" \\\n  -e SELF_URL="http://YOUR-SERVER-IP:8080" \\\n  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n  -e BOOTSTRAP_SNAPSHOT_URL="https://aequitas.digital/api/snapshot" \\\n  -e BOOTSTRAP_SIGNER="0x92cbedec9d348b4762cb9af99500ee6139c5b671" \\\n  -e SNAPSHOT_TOKEN="get-from-network-operator" \\\n  -p 8080:8080 -p 4001:4001 aequitas-node\n\n# Open firewall:\nufw allow 8080/tcp && ufw allow 4001/tcp',dn:'Port 8080 required (open inbound). Port 4001 optional (P2P). Tip: use --env-file /root/.aequitas.env instead of inline -e flags to keep secrets out of shell history.',
      s6:'6. Verify Your Node',v6:'Once running, check these endpoints to confirm the node is synced and healthy.',vc:'# 1. Node status (height should match the primary node within 1-2 blocks)\ncurl https://YOUR-NODE-URL/api/status\n# Expect: { "height": N, "total_humans": N, "index": N }\n\n# 2. EVM JSON-RPC (EVM compatibility check)\ncurl -X POST https://YOUR-NODE-URL/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'\n\n# 3. In startup logs, confirm:\n#    [NODE] Registered node operator wallet: 0xYOUR_WALLET\n#    Aequitas Node Running V  (blocks every ~6 seconds)\n\n# MetaMask: RPC URL = https://YOUR-NODE-URL/rpc | Chain ID: 1926 | Symbol: AEQ',
      s7:'7. P2P Networking & Sync',p7:'Set PRIMARY_NODE_URL=https://aequitas.digital in your environment. The node auto-registers with the primary on startup, receives the full peer list, and begins syncing. The libp2p multiaddress below is for advanced/manual setups:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'PRIMARY_NODE_URL (HTTPS) is the recommended method. The libp2p multiaddress may change if the primary node is redeployed.',
      s8:'8. Earning Validator Rewards',w8:'Validator rewards come from the Validators Pool (40% of all protocol fees: swap fees, demurrage, wealth cap overflow). Steps to receive rewards:',b8:['First register as a human on Aequitas: install the Android app and complete biometric registration to receive your wallet address and 1,000 AEQ','Set NODE_OPERATOR_WALLET to that registered wallet address in your node\'s environment variables','Start (or restart) your node — it calls RegisterNode() on startup. Confirm in logs: "[NODE] Registered node operator wallet: 0xYOUR_WALLET"','The primary node distributes rewards every 24 h to all registered node operator wallets proportional to blocks produced','Secondary nodes do NOT need to trigger distribution — just keep your node running and synced','No minimum uptime required, but consistently offline nodes contribute less to block production and proportionally less to the pool share'],
      s8b:'8b. Register Validator Key (Required for Rewards)',w8b:'Link your node signing key to your registered human wallet. Required to earn daily validator rewards. Requires SNAPSHOT_TOKEN for authentication.',b8b:['On the VPS server itself, or from any machine with the SNAPSHOT_TOKEN: curl -H "Authorization: Bearer YOUR-SNAPSHOT-TOKEN" "http://YOUR-NODE-IP:8080/api/sign-validator-challenge?wallet=0xYOUR_HUMAN_WALLET"','The command returns a JSON with a "signature" field. Copy that value.','Open the Network tab on the Aequitas website → Node Guide → scroll to Step 5b','Enter your RELAYER_ADDRESS (your node signing address, shown at node startup as "proposer addr: 0x...") and paste the signature','Click "Sign with MetaMask" using your registered human wallet — this proves you control both keys','Watch the node logs: [PEERS] Auto-authorized validator: 0xYOUR-SIGNING-ADDR should appear on the primary node','Also open port 8080 inbound on your VPS: ufw allow 8080/tcp — needed so the primary syncs your blocks'],
      s9:'9. Troubleshooting',th:['Symptom','Likely Cause','Fix'],tr:[['Height stays at 0 after start','PRIMARY_NODE_URL not set or wrong','Set PRIMARY_NODE_URL=https://aequitas.digital and SELF_URL to your public URL/IP'],['"no code at address" in logs','V7 contract not deployed','Verify RELAYER_ADDRESS is set; node auto-deploys V7 on startup if missing'],['"NODE_OPERATOR_WALLET not set" in logs','Missing env var','Set NODE_OPERATOR_WALLET to your registered human wallet address'],['DATABASE_URL error on startup','Wrong connection string or DB unreachable','Check format: postgres://user:pass@host:5432/dbname and that PostgreSQL is running. On Railway: DATABASE_URL is auto-set if PostgreSQL is in the same project.'],['Port 8080 not reachable from outside','Firewall blocking inbound','Run: ufw allow 8080/tcp — required for the primary node to sync your blocks'],['⚠ P2P bootstrap unreachable (HTTP sync still works)','Port 4001 firewalled — normal on Railway/Docker without -p 4001:4001','Not critical. HTTP sync handles blocks automatically. Add -p 4001:4001 and ufw allow 4001/tcp to also enable P2P.'],['Node not appearing in block explorer','Step 5b not done OR port 8080 closed','1. Open port 8080. 2. Complete Step 5b with your SNAPSHOT_TOKEN. 3. Your blocks will then appear in the explorer as MERGE blocks.'],['StateRoot mismatch warnings','Fresh DB — node has no registered humans','Set BOOTSTRAP_SNAPSHOT_URL=https://aequitas.digital/api/snapshot, BOOTSTRAP_SIGNER=0x92cbedec9d348b4762cb9af99500ee6139c5b671, SNAPSHOT_TOKEN=<from operator>. Restart — state imports automatically.'],['NODE_KEY generating new key on every restart','NODE_KEY env var not set','On first start: docker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY". Copy the base64 value, set as NODE_KEY, restart once.'],['MetaMask shows 0 AEQ or wrong balance','Stale cached network config','MetaMask → Settings → Networks → delete all Aequitas Chain entries → re-add via + ADD AEQUITAS NETWORK button on the website.'],['Bootstrap failed: snapshot too old','Node offline for >24 hours before bootstrap','Set SNAPSHOT_MAX_AGE_SECONDS=172800 (48h) on the secondary node env vars to allow older snapshots.']],
      s10:'10. MetaMask Configuration',m10:'To use your own node as the RPC endpoint in MetaMask or any EVM-compatible wallet:',mh:['Field','Value'],mr:[['Network Name','Aequitas Chain'],['RPC URL','https://YOUR-NODE-URL/rpc'],['Chain ID','1926  (hex: 0x786)'],['Currency Symbol','AEQ'],['Decimals','18'],['Block Explorer','https://aequitas.digital']],
      foot:'Open source · Permissionless · No admin keys · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'},
    de:{title:'Aequitas Node-Betreiber-Handbuch',sub:'Vollständige Schritt-für-Schritt-Anleitung · Aequitas Chain (Chain ID 1926)',badge:'v1.0 · Open Source · Erlaubnisfrei · Kein Stake erforderlich',
      s1:'1. Überblick',what:'Was ein Node leistet',wtxt:'Ein Aequitas-Node nimmt vollständig am Netzwerk teil: produziert Blöcke im BlockDAG-Konsens, validiert Groth16-Zero-Knowledge-Biometriebeweise für neue Menschenregistrierungen, setzt Vermögensobergrenzen und Demurrage auf Protokollebene durch, synchronisiert den Zustand mit Peers via libp2p + HTTP und führt optional tägliche Pool-Verteilungen durch. Jeder Node führt die vollständige Chain aus — es gibt keine Light-Clients.',
      earn:'Was du verdienst',etxt:'NODE_OPERATOR_WALLET auf eine als Mensch registrierte Wallet-Adresse setzen. Der Validators-Pool erhält 40% aller Protokollgebühren (Swap-Gebühren, Demurrage, Vermögensobergrenze-Überschuss). Alle 24 Stunden verteilt der primäre Node den Pool-Saldo proportional auf alle registrierten Node-Betreiber-Wallets. Kein Stake erforderlich — Blockproduktion ist vollständig erlaubnisfrei.',
      s2:'2. Voraussetzungen',rh:['Komponente','Minimum','Empfohlen'],rr:[['Betriebssystem','Linux / Docker-fähiger Host','Ubuntu 22.04 LTS'],['RAM','512 MB','1 GB (EVM braucht Spielraum)'],['CPU','1 vCPU','2 vCPU (Groth16 ist CPU-gebunden)'],['Speicher','2 GB','10 GB SSD (Chain wächst kontinuierlich)'],['Datenbank','PostgreSQL 14+','Railway oder Supabase (verwaltet)'],['Netzwerk','Öffentliche IP / Port-Weiterleitung','TCP 8080 offen, stabile Verfügbarkeit']],
      s3:'3. Umgebungsvariablen',e3:'Diese vor dem Start des Nodes setzen. Mit JA markierte Variablen sind Pflicht; "Für Bel." wird benötigt um Validator-Auszahlungen zu erhalten.',eh:['Variable','Zweck','Pflicht?'],er:[['DATABASE_URL','PostgreSQL-Verbindungsstring: postgres://user:pass@host:5432/aequitas','JA'],['RELAYER_PRIVATE_KEY','Privater Schlüssel (0x...) des EOA der On-Chain-Menschenregistrierungen signiert','JA'],['NODE_OPERATOR_WALLET','Registrierte Mensch-Wallet-Adresse die täglich Validator-Pool-Bel. erhält','Für Bel.'],['RELAYER_ADDRESS','EOA-Adresse passend zu RELAYER_PRIVATE_KEY. Hat Fallback, aber explizit setzen.','Empfohlen'],['PORT','HTTP-Port für API + JSON-RPC. Standard: 8080','NEIN'],['PEER_SECRET','Geteiltes Geheimnis das diesen Node als Validator autorisiert. ALLE Nodes müssen denselben Wert nutzen. Vom Netzwerkbetreiber erhalten.','Für Multi-Node'],['SELF_URL','Eigene öffentliche HTTPS-URL dieses Nodes. In Railway: Settings > Networking.','Für Multi-Node'],['PRIMARY_NODE_URL','Primär-Node für automatische Peer-Discovery. Auf https://aequitas.digital setzen.','Für Multi-Node'],['BOOTSTRAP_SNAPSHOT_URL','Auf https://aequitas.digital/api/snapshot setzen bei einem frischen Node. Lädt den vollständigen Netzwerk-State beim Start herunter (nur wenn DB leer, d.h. 0 Menschen). Behebt StateRoot-Mismatch sofort.','Multi-Node'],['BOOTSTRAP_SIGNER','Ethereum-Signing-Adresse des Primär-Nodes. Aus https://aequitas.digital/api/status → Feld signing_address. Pflicht wenn BOOTSTRAP_SNAPSHOT_URL gesetzt.','Multi-Node'],['SNAPSHOT_TOKEN','Bearer-Token für /api/snapshot-Authentifizierung. Muss mit dem SNAPSHOT_TOKEN des Primär-Nodes übereinstimmen. Vom Netzwerkbetreiber erfragen.','Multi-Node'],['PEER_NODES','Statische Peer-URLs (veraltet). PRIMARY_NODE_URL für Auto-Discovery verwenden.','Optional'],['NODE_KEY','Base64-kodierter libp2p-Private-Key. Wenn nicht gesetzt: wird generiert und in stderr ausgegeben als "SAVE THIS AS NODE_KEY: <base64>". Kopieren und setzen.','NEIN'],['IS_PRIMARY_NODE','"true" nur auf dem designierten Primär-Node. VERALTET: Verteilungen nutzen jetzt DB-Lock — nicht mehr erforderlich.','NEIN'],['RESET_STATE','"true" löscht die DB beim Start. DESTRUKTIV — niemals in Produktion.','NEIN']],
      s4:'4. Schnellstart — Railway (Empfohlen)',r4:'Railway ist der schnellste Einstieg. Der kostenlose Tarif erfüllt die Mindestanforderungen derzeit. Geschätzte Einrichtungszeit: 10–15 Minuten.',rs:['Repo forken: https://github.com/hanoi96international-gif/Aequitas','Railway-Konto auf railway.app erstellen und neues Projekt starten','"Deploy from GitHub Repo" anklicken und den Fork auswählen','Im Projekt: + New → Database → Add PostgreSQL — DATABASE_URL wird automatisch gesetzt','Service → Settings → Variables aufrufen und Umgebungsvariablen aus Abschnitt 3 hinzufügen','PRIMARY_NODE_URL=https://aequitas.digital setzen','NODE_OPERATOR_WALLET=<deine registrierte AEQ-Mensch-Wallet> für tägliche Validator-Bel. setzen','RELAYER_PRIVATE_KEY=<EOA-Privatschlüssel für On-Chain-Registrierungssignaturen> setzen','"Deploy" klicken — das Dockerfile im Root-Verzeichnis steuert den Build (~3 Min. für Go-Kompilierung)','Deploy-Logs auf "Aequitas Node Running ✓" und "[NODE] Registered node operator wallet" prüfen','DEINE-RAILWAY-URL/api/status aufrufen um Synchronisation des Nodes zu bestätigen','Node-RPC zu MetaMask hinzufügen: Chain-ID 1926, Symbol AEQ, URL https://DEINE-URL/rpc'],rn:'Railway vergibt eine zufällige Subdomain; benutzerdefinierte Domains in den Projekteinstellungen konfigurierbar. Nur Port 8080 muss exponiert werden — P2P wird intern verwaltet.',
      s5:'5. Schnellstart — Docker / VPS',d5:'Für VPS (Contabo, Hetzner, DigitalOcean). Docker + PostgreSQL erforderlich. Railway-PostgreSQL kann als externe DB genutzt werden. NODE_OPERATOR_WALLET muss ein registrierter Aequitas-Human sein.',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\ndocker build -t aequitas-node .\n\n# Schritt 1: Erster Start (noch kein NODE_KEY — erscheint in den Logs)\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgresql://user:pass@host:5432/dbname" \\\n  -e RELAYER_PRIVATE_KEY="DEIN_HEX_PRIVATER_SCHLUESSEL" \\\n  -e RELAYER_ADDRESS="0xDEINE_NODE_SIGNING_ADRESSE" \\\n  -e NODE_OPERATOR_WALLET="0xDEINE_REGISTRIERTE_HUMAN_WALLET" \\\n  -e PEER_SECRET="vom-netzwerkbetreiber" \\\n  -e SELF_URL="http://DEINE-SERVER-IP:8080" \\\n  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n  -e BOOTSTRAP_SNAPSHOT_URL="https://aequitas.digital/api/snapshot" \\\n  -e BOOTSTRAP_SIGNER="0x92cbedec9d348b4762cb9af99500ee6139c5b671" \\\n  -e SNAPSHOT_TOKEN="vom-netzwerkbetreiber" \\\n  -p 8080:8080 -p 4001:4001 aequitas-node\n\n# Schritt 2: NODE_KEY aus Logs kopieren (einmalig):\ndocker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY"\n\n# Schritt 3: Endgültiger Befehl mit NODE_KEY (dauerhaft verwenden):\ndocker stop aequitas-node && docker rm aequitas-node\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgresql://user:pass@host:5432/dbname" \\\n  -e RELAYER_PRIVATE_KEY="DEIN_HEX_PRIVATER_SCHLUESSEL" \\\n  -e RELAYER_ADDRESS="0xDEINE_NODE_SIGNING_ADRESSE" \\\n  -e NODE_OPERATOR_WALLET="0xDEINE_REGISTRIERTE_HUMAN_WALLET" \\\n  -e NODE_KEY="base64-aus-schritt-2" \\\n  -e PEER_SECRET="vom-netzwerkbetreiber" \\\n  -e SELF_URL="http://DEINE-SERVER-IP:8080" \\\n  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n  -e BOOTSTRAP_SNAPSHOT_URL="https://aequitas.digital/api/snapshot" \\\n  -e BOOTSTRAP_SIGNER="0x92cbedec9d348b4762cb9af99500ee6139c5b671" \\\n  -e SNAPSHOT_TOKEN="vom-netzwerkbetreiber" \\\n  -p 8080:8080 -p 4001:4001 aequitas-node\n\n# Firewall:\nufw allow 8080/tcp && ufw allow 4001/tcp',dn:'Port 8080 muss von außen erreichbar sein (Pflicht). Port 4001 optional (P2P). Tipp: --env-file /root/.aequitas.env statt einzelner -e Flags verwenden — Secrets bleiben aus der Shell-History.',
      s6:'6. Node verifizieren',v6:'Sobald der Node läuft, diese Endpunkte prüfen um Synchronisation und Gesundheit zu bestätigen.',vc:'# 1. Node-Status (Höhe sollte mit Primär-Node übereinstimmen)\ncurl https://DEINE-NODE-URL/api/status\n# Erwartet: { "height": N, "total_humans": N, "index": N }\n\n# 2. EVM JSON-RPC prüfen\ncurl -X POST https://DEINE-NODE-URL/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'\n\n# 3. In Start-Logs bestätigen:\n#    [NODE] Registered node operator wallet: 0xDEINE_WALLET\n#    Aequitas Node Running V  (Blöcke alle ~6 Sekunden)\n\n# MetaMask: RPC-URL https://DEINE-NODE-URL/rpc | Chain-ID: 1926 | Symbol: AEQ',
      s7:'7. P2P-Netzwerk & Synchronisation',p7:'PRIMARY_NODE_URL=https://aequitas.digital setzen. Der Node verbindet sich automatisch und synchronisiert den vollständigen Chain-Zustand via libp2p-Gossip (Echtzeit) und HTTP-Pulls von Peers (Fallback). Libp2p-Multiaddresse des Primär-Nodes:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'PRIMARY_NODE_URL (HTTPS) ist die empfohlene Methode. Die libp2p-Multiaddresse kann sich bei einem Neudeployment ändern.',
      s8:'8. Validator-Belöhnungen erhalten',w8:'Validator-Belöhnungen kommen aus dem Validators-Pool (40% aller Protokollgebühren). Schritte um Belöhnungen zu erhalten:',b8:['Zuerst als Mensch auf Aequitas registrieren: Android-App installieren und biometrische Registrierung abschließen um Wallet-Adresse und 1.000 AEQ zu erhalten','NODE_OPERATOR_WALLET auf diese registrierte Wallet-Adresse in den Umgebungsvariablen des Nodes setzen','Node starten (oder neu starten) — er ruft RegisterNode() beim Start auf. In Logs bestätigen: "[NODE] Registered node operator wallet: 0xDEINE_WALLET"','Der Primär-Node verteilt Belöhnungen alle 24 Stunden an alle registrierten Node-Betreiber-Wallets proportional zur Blockproduktion','Sekundäre Nodes müssen die Verteilung NICHT auslösen — einfach den Node laufen lassen und synchronisiert halten','Keine Mindest-Verfügbarkeit erforderlich, aber dauerhaft offline Nodes tragen weniger zur Blockproduktion und zum Pool-Anteil bei'],
      s8b:'8b. Validator-Key registrieren (für Rewards)',w8b:'Signing-Key mit registrierter Human-Wallet verknüpfen. Pflicht für tägliche Validator-Belohnungen. Benötigt SNAPSHOT_TOKEN zur Authentifizierung.',b8b:['Auf dem VPS-Server selbst oder von jedem Rechner mit SNAPSHOT_TOKEN: curl -H "Authorization: Bearer DEIN-SNAPSHOT-TOKEN" "http://DEINE-NODE-IP:8080/api/sign-validator-challenge?wallet=0xDEINE_HUMAN_WALLET"','Befehl gibt JSON zurück — das "signature"-Feld kopieren','Website öffnen: Network-Tab → Node Guide → zu Schritt 5b scrollen','RELAYER_ADDRESS (Signing-Adresse des Nodes, bei Start als "proposer addr: 0x..." angezeigt) und Signatur eintragen','"Sign with MetaMask" mit der registrierten Human-Wallet klicken','In Node-Logs: [PEERS] Auto-authorized validator: 0xDEINE-ADRESSE sollte auf dem Primär-Node erscheinen','Port 8080 öffnen: ufw allow 8080/tcp — damit der Primär-Node deine Blöcke syncen kann'],
      s9:'9. Fehlerbehebung',th:['Symptom','Wahrscheinliche Ursache','Lösung'],tr:[['Höhe bleibt bei 0 nach Start','PRIMARY_NODE_URL nicht gesetzt oder falsch','PRIMARY_NODE_URL=https://aequitas.digital und SELF_URL auf eigene IP/URL setzen'],['"no code at address" in Logs','V7-Contract nicht im EVM deployed','RELAYER_ADDRESS prüfen; Node deployed V7 automatisch beim Start wenn fehlend'],['"NODE_OPERATOR_WALLET not set" in Logs','Fehlende Umgebungsvariable','NODE_OPERATOR_WALLET auf registrierte Mensch-Wallet-Adresse setzen'],['DATABASE_URL-Fehler beim Start','Falscher Verbindungsstring oder DB nicht erreichbar','Format prüfen: postgres://user:pass@host:5432/dbname. Bei Railway: DATABASE_URL wird automatisch gesetzt wenn PostgreSQL im selben Projekt.'],['Port 8080 nicht erreichbar','Firewall blockiert eingehend','ufw allow 8080/tcp — Pflicht damit Primär-Node deine Blöcke syncen kann'],['⚠ P2P Bootstrap nicht erreichbar (HTTP-Sync funktioniert trotzdem)','Port 4001 geblockt — normal bei Railway/Docker ohne -p 4001:4001','Nicht kritisch. HTTP-Sync läuft automatisch. Für P2P: -p 4001:4001 und ufw allow 4001/tcp hinzufügen.'],['Node nicht im Block Explorer / keine MERGE-Blöcke','Schritt 5b nicht abgeschlossen ODER Port 8080 geschlossen','1. Port 8080 öffnen. 2. Schritt 5b mit SNAPSHOT_TOKEN durchführen. 3. Blöcke erscheinen dann als MERGE im Explorer.'],['StateRoot-Mismatch-Warnungen in Logs','Frischer Node ohne registrierte Menschen in DB','BOOTSTRAP_SNAPSHOT_URL=https://aequitas.digital/api/snapshot, BOOTSTRAP_SIGNER=0x92cbedec9d348b4762cb9af99500ee6139c5b671, SNAPSHOT_TOKEN=<vom Betreiber> setzen. Neustart — State wird automatisch importiert.'],['NODE_KEY generiert bei jedem Neustart neuen Key','NODE_KEY Umgebungsvariable nicht gesetzt','Beim ersten Start: docker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY". Base64-Wert kopieren, als NODE_KEY setzen, einmal neu starten.'],['MetaMask zeigt 0 AEQ oder falschen Saldo','Veraltete Netzwerkkonfiguration in MetaMask','MetaMask → Einstellungen → Netzwerke → alle "Aequitas Chain"-Einträge löschen → über "+ AEQUITAS-NETZWERK HINZUFÜGEN" auf der Website neu hinzufügen.']],
      s10:'10. MetaMask-Konfiguration',m10:'Um deinen eigenen Node als RPC-Endpunkt in MetaMask oder einer anderen EVM-kompatiblen Wallet zu verwenden:',mh:['Feld','Wert'],mr:[['Netzwerkname','Aequitas Chain'],['RPC-URL','https://DEINE-NODE-URL/rpc'],['Chain-ID','1926  (hex: 0x786)'],['Währungssymbol','AEQ'],['Dezimalstellen','18'],['Block-Explorer','https://aequitas.digital']],
      foot:'Open Source · Erlaubnisfrei · Keine Admin-Schlüssel · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'},
    es:{title:'Guia del Operador de Nodos Aequitas',sub:'Guia completa paso a paso · Aequitas Chain (Chain ID 1926)',badge:'v1.0 · Codigo Abierto · Sin permisos · Sin stake requerido',
      s1:'1. Vision General',what:'Que hace un nodo',wtxt:'Un nodo Aequitas participa plenamente en la red: produce bloques en el consenso BlockDAG, valida pruebas biometricas Groth16 de conocimiento cero para nuevos registros humanos, aplica limites de riqueza y demurrage a nivel de protocolo, sincroniza el estado con pares via libp2p + HTTP y ejecuta distribuciones diarias de pools. Cada nodo ejecuta la cadena completa: no hay clientes ligeros.',
      earn:'Que ganas',etxt:'Establece NODE_OPERATOR_WALLET en una billetera humana registrada. El Pool de Validadores acumula el 40% de todas las tarifas del protocolo. Cada 24 h el nodo primario distribuye el saldo proporcionalmente entre todos los operadores registrados. No se requiere stake.',
      s2:'2. Requisitos',rh:['Componente','Minimo','Recomendado'],rr:[['SO','Linux / host con Docker','Ubuntu 22.04 LTS'],['RAM','512 MB','1 GB'],['CPU','1 vCPU','2 vCPU'],['Almacenamiento','2 GB','10 GB SSD'],['Base de datos','PostgreSQL 14+','Railway o Supabase'],['Red','IP publica / reenvio de puerto','TCP 8080 abierto']],
      s3:'3. Variables de Entorno',e3:'Configura estas variables antes de iniciar el nodo. Las marcadas SI son obligatorias.',eh:['Variable','Proposito','Requerida?'],er:[['DATABASE_URL','Cadena de conexion PostgreSQL: postgres://user:pass@host:5432/aequitas','SI'],['RELAYER_PRIVATE_KEY','Clave privada (0x...) del EOA que firma registros on-chain','SI'],['NODE_OPERATOR_WALLET','Billetera humana registrada que recibe recompensas diarias del pool','Para recomp.'],['RELAYER_ADDRESS','Direccion EOA. Tiene fallback pero configurar explicitamente.','Recomendado'],['PORT','Puerto HTTP. Por defecto: 8080','NO'],['PEER_SECRET','Secreto compartido que autoriza este nodo como validador. TODOS los nodos deben usar el MISMO valor.','Multi-nodo'],['SELF_URL','URL HTTPS publica de este nodo. Necesaria para excluirse en el descubrimiento de pares.','Multi-nodo'],['PRIMARY_NODE_URL','URL del nodo primario para descubrimiento automatico. Establecer en https://aequitas.digital.','Multi-nodo'],['PEER_NODES','URLs de pares estaticos (legado). Usar PRIMARY_NODE_URL.','Opcional'],['NODE_KEY','Clave libp2p base64 para identidad P2P estable. Si no se establece: auto-generada en stderr como "SAVE THIS AS NODE_KEY: <base64>".','NO'],['IS_PRIMARY_NODE','"true" solo en el nodo primario. OBSOLETO: las distribuciones usan DB-lock — no requerido.','NO'],['RESET_STATE','"true" borra la BD al iniciar. DESTRUCTIVO.','NO']],
      s4:'4. Inicio Rapido — Railway (Recomendado)',r4:'Railway es la forma mas rapida de comenzar. El nivel gratuito cubre los requisitos minimos para el lanzamiento. Tiempo estimado: 10-15 minutos.',rs:['Haz un fork del repo: https://github.com/hanoi96international-gif/Aequitas','Crea una cuenta en railway.app e inicia un nuevo proyecto','Haz clic en "Deploy from GitHub Repo" y selecciona tu fork','En el proyecto: + New → Database → Add PostgreSQL','Ve a tu servicio → Variables y agrega las variables de la Seccion 3','Establece PRIMARY_NODE_URL=https://aequitas.digital','Establece NODE_OPERATOR_WALLET=<tu billetera humana AEQ>','Establece RELAYER_PRIVATE_KEY=<tu clave privada EOA>','Haz clic en "Deploy" — el Dockerfile gestiona la compilacion (~3 min)','En los logs busca: "Aequitas Node Running" y "[NODE] Registered node operator wallet"','Abre TU-URL/api/status para confirmar que el nodo esta activo','Agrega el RPC a MetaMask: Chain ID 1926, Simbolo AEQ, URL https://TU-URL/rpc'],rn:'Railway asigna un subdominio aleatorio; los dominios personalizados se configuran en ajustes del proyecto.',
      s5:'5. Inicio Rapido — Docker',d5:'Para VPS, VM en la nube o servidor local. Requiere Docker y PostgreSQL disponibles. NODE_OPERATOR_WALLET debe ser una billetera humana registrada en Aequitas.',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\n\n# Construir imagen (~3 min)\ndocker build -t aequitas-node .\n\n# Primer inicio (NODE_KEY sera impreso en logs)\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n  -e RELAYER_PRIVATE_KEY="0xTU_CLAVE_PRIVADA" \\\n  -e RELAYER_ADDRESS="0xTU_DIRECCION" \\\n  -e NODE_OPERATOR_WALLET="0xTU_BILLETERA_HUMANA" \\\n  -e PEER_SECRET="del-operador-de-red" \\\n  -e SELF_URL="http://TU-IP-SERVIDOR:8080" \\\n  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n  -e BOOTSTRAP_SNAPSHOT_URL="https://aequitas.digital/api/snapshot" \\\n  -e BOOTSTRAP_SIGNER="0x92cbedec9d348b4762cb9af99500ee6139c5b671" \\\n  -e SNAPSHOT_TOKEN="del-operador-de-red" \\\n  -p 8080:8080 -p 4001:4001 aequitas-node\n\n# Copiar NODE_KEY de logs del primer inicio:\ndocker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY"\n\n# Reiniciar con NODE_KEY para ID P2P estable:\ndocker stop aequitas-node && docker rm aequitas-node\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e NODE_KEY="<base64 de arriba>" \\\n  ... (todas las otras vars) ...',dn:'Puerto 8080 debe ser accesible desde fuera (ufw allow 8080/tcp). Puerto 4001 es opcional para P2P (ufw allow 4001/tcp). HTTP sync funciona sin P2P.',
      s6:'6. Verificar el Nodo',v6:'Una vez en ejecucion, comprueba estos endpoints:',vc:'curl https://TU-NODO-URL/api/status\n# Esperado: {"height": N, "total_humans": N}\n\ncurl -X POST https://TU-NODO-URL/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'',
      s7:'7. Red P2P y Sincronizacion',p7:'Establece PRIMARY_NODE_URL=https://aequitas.digital en tu entorno. El nodo sincroniza la cadena automaticamente. Multidireccion libp2p del nodo primario:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'PRIMARY_NODE_URL (HTTPS) es el metodo recomendado. La multidireccion libp2p puede cambiar.',
      s8:'8. Ganar Recompensas de Validador',w8:'Las recompensas provienen del Pool de Validadores (40% de todas las tarifas del protocolo). Pasos:',b8:['Registrate como humano en Aequitas: instala la app Android y completa el registro biometrico','Establece NODE_OPERATOR_WALLET en esa direccion de billetera registrada','Reinicia tu nodo y confirma en logs: "[NODE] Registered node operator wallet: 0x..."','El nodo primario distribuye recompensas cada 24 h a todos los operadores registrados','Los nodos secundarios NO necesitan activar la distribucion — solo manten tu nodo en ejecucion'],
      s9:'9. Solucion de Problemas',th:['Sintoma','Causa probable','Solucion'],tr:[['Altura permanece en 0','PRIMARY_NODE_URL no configurado o incorrecto','Establece PRIMARY_NODE_URL=https://aequitas.digital y SELF_URL a tu URL/IP publica y redespliega'],['"no code at address" en logs','Contrato V7 no desplegado aun','Normal en el primer inicio — el nodo despliega V7 automaticamente'],['Error DATABASE_URL','Cadena de conexion incorrecta','Verifica el formato: postgres://usuario:clave@host:5432/dbname'],['Puerto 8080 no accesible','Firewall o configuracion del proveedor','ufw allow 8080/tcp — necesario para que el nodo primario sincronice tus bloques'],['Bootstrap snapshot fallido / StateRoot mismatch','SNAPSHOT_TOKEN no configurado o BOOTSTRAP_SIGNER incorrecto','Establece BOOTSTRAP_SNAPSHOT_URL=https://aequitas.digital/api/snapshot, BOOTSTRAP_SIGNER=0x92cbedec9d348b4762cb9af99500ee6139c5b671, SNAPSHOT_TOKEN=<del operador>. Reinicia — el estado se importa automaticamente.'],['NODE_KEY generando nueva clave en cada reinicio','Variable NODE_KEY no configurada','En el primer inicio: docker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY". Copia el valor base64, configuralo como NODE_KEY y reinicia una vez.'],['MetaMask muestra 0 AEQ o saldo incorrecto','Configuracion de red obsoleta en MetaMask','MetaMask → Configuracion → Redes → eliminar todas las entradas de "Aequitas Chain" → volver a agregar con el boton en el sitio web.'],['Nodo rechazado como validador','NODE_OPERATOR_WALLET no es una billetera humana registrada','Registrate primero en Aequitas (app Android), luego establece NODE_OPERATOR_WALLET a esa billetera registrada.']],
      s10:'10. Configuracion de MetaMask',m10:'Para usar tu nodo como endpoint RPC en MetaMask:',mh:['Campo','Valor'],mr:[['Nombre de red','Aequitas Chain'],['URL RPC','https://TU-NODO-URL/rpc'],['Chain ID','1926  (hex: 0x786)'],['Simbolo','AEQ'],['Decimales','18'],['Explorador','https://aequitas.digital']],
      foot:'Codigo abierto · Sin permisos · Sin claves de administrador · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'},
    it:{title:'Guida per Operatori di Nodi Aequitas',sub:'Guida completa passo dopo passo · Aequitas Chain (Chain ID 1926)',badge:'v1.0 · Open Source · Senza permessi · Nessuno stake richiesto',
      s1:'1. Panoramica',what:'Cosa fa un nodo',wtxt:'Un nodo Aequitas partecipa pienamente alla rete: produce blocchi nel consenso BlockDAG, valida prove biometriche Groth16 a conoscenza zero per le nuove registrazioni umane, applica limiti di ricchezza e demurrage a livello di protocollo, sincronizza lo stato con i peer via libp2p + HTTP ed esegue distribuzioni giornaliere dei pool. Ogni nodo esegue la catena completa: non esistono client leggeri.',
      earn:'Cosa guadagni',etxt:'Imposta NODE_OPERATOR_WALLET su un indirizzo wallet registrato come umano. Il Pool Validatori accumula il 40% di tutte le commissioni di protocollo. Ogni 24 h il nodo primario distribuisce il saldo del pool proporzionalmente tra tutti i wallet degli operatori registrati. Nessuno stake richiesto.',
      s2:'2. Requisiti',rh:['Componente','Minimo','Consigliato'],rr:[['SO','Linux / host con Docker','Ubuntu 22.04 LTS'],['RAM','512 MB','1 GB'],['CPU','1 vCPU','2 vCPU'],['Archiviazione','2 GB','10 GB SSD'],['Database','PostgreSQL 14+','Railway o Supabase'],['Rete','IP pubblica / port forward','TCP 8080 aperto']],
      s3:'3. Variabili di Ambiente',e3:'Configura queste variabili prima di avviare il nodo. Quelle contrassegnate con SI sono obbligatorie.',eh:['Variabile','Scopo','Richiesta?'],er:[['DATABASE_URL','Stringa di connessione PostgreSQL: postgres://user:pass@host:5432/aequitas','SI'],['RELAYER_PRIVATE_KEY','Chiave privata (0x...) dell\'EOA che firma le registrazioni on-chain','SI'],['NODE_OPERATOR_WALLET','Wallet umano registrato che riceve le ricompense giornaliere del pool','Per ricomp.'],['RELAYER_ADDRESS','Indirizzo EOA corrispondente a RELAYER_PRIVATE_KEY. Ha un fallback ma impostalo esplicitamente.','Consigliato'],['PORT','Porta HTTP per API + JSON-RPC. Default: 8080','NO'],['PEER_NODES','URL dei peer bootstrap (legacy). Usare PRIMARY_NODE_URL.','Facoltativo'],['PEER_SECRET','Segreto condiviso: TUTTI i nodi devono usare lo STESSO valore.','Multi-nodo'],['SELF_URL','URL HTTPS pubblica per self-exclusion.','Multi-nodo'],['PRIMARY_NODE_URL','Nodo primario per peer discovery (https://aequitas.digital).','Multi-nodo'],['NODE_KEY','Chiave libp2p base64. Auto-generata in stderr come "SAVE THIS AS NODE_KEY: <base64>" se omessa.','NO'],['IS_PRIMARY_NODE','"true" solo sul nodo primario. OBSOLETO: le distribuzioni usano DB-lock — non necessario.','NO'],['RESET_STATE','"true" cancella il DB all\'avvio. DISTRUTTIVO.','NO']],
      s4:'4. Avvio Rapido — Railway (Consigliato)',r4:'Railway e il modo piu veloce per iniziare. Il livello gratuito soddisfa i requisiti minimi per il lancio. Tempo stimato: 10-15 minuti.',rs:['Fai un fork del repo: https://github.com/hanoi96international-gif/Aequitas','Crea un account su railway.app e avvia un nuovo progetto','Clicca "Deploy from GitHub Repo" e seleziona il tuo fork','Nel progetto: + New → Database → Add PostgreSQL','Vai al tuo servizio → Variables e aggiungi le variabili della Sezione 3','Imposta PRIMARY_NODE_URL=https://aequitas.digital','Imposta NODE_OPERATOR_WALLET=<il tuo wallet umano AEQ>','Imposta RELAYER_PRIVATE_KEY=<la tua chiave privata EOA>','Clicca "Deploy" — il Dockerfile gestisce la compilazione (~3 min)','Nei log cerca: "Aequitas Node Running" e "[NODE] Registered node operator wallet"','Apri TUO-URL/api/status per confermare che il nodo e attivo','Aggiungi il tuo RPC a MetaMask: Chain ID 1926, Simbolo AEQ, URL https://TUO-URL/rpc'],rn:'Railway assegna un sottodominio casuale; i domini personalizzati si configurano nelle impostazioni del progetto.',
      s5:'5. Avvio Rapido — Docker',d5:'Per VPS, VM cloud o server locale. Prerequisiti: Docker installato e PostgreSQL disponibile. NODE_OPERATOR_WALLET deve essere un wallet umano registrato su Aequitas.',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\n\n# Crea immagine (~3 min)\ndocker build -t aequitas-node .\n\n# Primo avvio (NODE_KEY verra stampato nei log)\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n  -e RELAYER_PRIVATE_KEY="0xLA_TUA_CHIAVE_PRIVATA" \\\n  -e RELAYER_ADDRESS="0xIL_TUO_INDIRIZZO" \\\n  -e NODE_OPERATOR_WALLET="0xIL_TUO_WALLET_UMANO" \\\n  -e PEER_SECRET="dall-operatore-di-rete" \\\n  -e SELF_URL="http://IL-TUO-IP-SERVER:8080" \\\n  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n  -e BOOTSTRAP_SNAPSHOT_URL="https://aequitas.digital/api/snapshot" \\\n  -e BOOTSTRAP_SIGNER="0x92cbedec9d348b4762cb9af99500ee6139c5b671" \\\n  -e SNAPSHOT_TOKEN="dall-operatore-di-rete" \\\n  -p 8080:8080 -p 4001:4001 aequitas-node\n\n# Copia NODE_KEY dai log del primo avvio:\ndocker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY"\n\n# Riavvia con NODE_KEY per ID P2P stabile:\ndocker stop aequitas-node && docker rm aequitas-node\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e NODE_KEY="<base64 dal passo precedente>" \\\n  ... (tutte le altre variabili) ...',dn:'La porta 8080 deve essere raggiungibile dall\'esterno (ufw allow 8080/tcp). La porta 4001 e opzionale per P2P (ufw allow 4001/tcp). La sincronizzazione HTTP funziona senza P2P.',
      s6:'6. Verifica il Nodo',v6:'Una volta avviato, controlla questi endpoint per confermare che il nodo e sincronizzato.',vc:'curl https://TUO-NODO-URL/api/status\n# Atteso: {"height": N, "total_humans": N}\n\ncurl -X POST https://TUO-NODO-URL/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'',
      s7:'7. Rete P2P e Sincronizzazione',p7:'Imposta PRIMARY_NODE_URL=https://aequitas.digital nel tuo ambiente. Il nodo si connette e sincronizza la catena automaticamente. Multiindirizzo libp2p del nodo primario:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'PRIMARY_NODE_URL (HTTPS) e il metodo raccomandato. Il multiindirizzo libp2p puo cambiare.',
      s8:'8. Guadagnare Ricompense da Validatore',w8:'Le ricompense provengono dal Pool Validatori (40% di tutte le commissioni di protocollo). Passaggi:',b8:['Prima registrati come umano su Aequitas: installa l\'app Android e completa la registrazione biometrica','Imposta NODE_OPERATOR_WALLET su quell\'indirizzo wallet registrato','Riavvia il nodo e conferma nei log: "[NODE] Registered node operator wallet: 0x..."','Il nodo primario distribuisce le ricompense ogni 24 h a tutti i wallet degli operatori registrati','I nodi secondari NON devono attivare la distribuzione — tieni solo il nodo in esecuzione'],
      s9:'9. Risoluzione dei Problemi',th:['Sintomo','Causa probabile','Soluzione'],tr:[['L\'altezza rimane a 0','PRIMARY_NODE_URL non impostato o errato','Imposta PRIMARY_NODE_URL=https://aequitas.digital e SELF_URL al tuo URL/IP pubblico e ridistribuisci'],['"no code at address" nei log','Contratto V7 non ancora distribuito','Normale al primo avvio — il nodo distribuisce V7 automaticamente'],['Errore DATABASE_URL','Stringa di connessione errata','Verifica il formato: postgres://utente:password@host:5432/dbname'],['Porta 8080 non raggiungibile','Firewall o configurazione del provider','ufw allow 8080/tcp — necessario perche il nodo primario sincronizzi i tuoi blocchi'],['Bootstrap snapshot fallito / StateRoot mismatch','SNAPSHOT_TOKEN mancante o BOOTSTRAP_SIGNER errato','Imposta BOOTSTRAP_SNAPSHOT_URL=https://aequitas.digital/api/snapshot, BOOTSTRAP_SIGNER=0x92cbedec9d348b4762cb9af99500ee6139c5b671, SNAPSHOT_TOKEN=<dall operatore>. Riavvia — lo stato viene importato automaticamente.'],['NODE_KEY genera nuova chiave ad ogni riavvio','Variabile NODE_KEY non impostata','Al primo avvio: docker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY". Copia il valore base64 e impostalo come NODE_KEY, poi riavvia una volta.'],['MetaMask mostra 0 AEQ o saldo errato','Configurazione di rete obsoleta in MetaMask','MetaMask → Impostazioni → Reti → elimina tutte le voci "Aequitas Chain" → aggiungi di nuovo tramite il pulsante sul sito web.'],['Nodo rifiutato come validatore','NODE_OPERATOR_WALLET non e un wallet umano registrato','Registrati prima su Aequitas (app Android), poi imposta NODE_OPERATOR_WALLET a quel wallet registrato.']],
      s10:'10. Configurazione MetaMask',m10:'Per usare il tuo nodo come endpoint RPC in MetaMask:',mh:['Campo','Valore'],mr:[['Nome rete','Aequitas Chain'],['URL RPC','https://TUO-NODO-URL/rpc'],['Chain ID','1926  (hex: 0x786)'],['Simbolo','AEQ'],['Decimali','18'],['Block Explorer','https://aequitas.digital']],
      foot:'Open source · Senza permessi · Senza chiavi admin · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'},
    id:{title:'Panduan Operator Node Aequitas',sub:'Panduan lengkap langkah demi langkah · Aequitas Chain (Chain ID 1926)',badge:'v1.0 · Open Source · Tanpa Izin · Tidak perlu stake',
      s1:'1. Gambaran Umum',what:'Apa yang dilakukan node',wtxt:'Node Aequitas berpartisipasi penuh dalam jaringan: memproduksi blok dalam konsensus BlockDAG, memvalidasi bukti biometrik Groth16 zero-knowledge untuk pendaftaran manusia baru, menerapkan batas kekayaan dan demurrage di tingkat protokol, menyinkronkan status dengan peer via libp2p + HTTP, dan menjalankan distribusi pool harian. Setiap node menjalankan rantai penuh — tidak ada klien ringan.',
      earn:'Apa yang kamu dapatkan',etxt:'Atur NODE_OPERATOR_WALLET ke alamat wallet manusia terdaftar. Pool Validator mengumpulkan 40% dari semua biaya protokol. Setiap 24 jam, node primer mendistribusikan saldo pool secara proporsional ke semua wallet operator node terdaftar. Tidak perlu stake.',
      s2:'2. Persyaratan',rh:['Komponen','Minimum','Direkomendasikan'],rr:[['OS','Linux / host berkemampuan Docker','Ubuntu 22.04 LTS'],['RAM','512 MB','1 GB'],['CPU','1 vCPU','2 vCPU'],['Penyimpanan','2 GB','10 GB SSD'],['Database','PostgreSQL 14+','Railway atau Supabase'],['Jaringan','IP publik / port forward','TCP 8080 terbuka']],
      s3:'3. Variabel Lingkungan',e3:'Atur variabel ini sebelum memulai node. Yang ditandai YA wajib diisi.',eh:['Variabel','Tujuan','Wajib?'],er:[['DATABASE_URL','String koneksi PostgreSQL: postgres://user:pass@host:5432/aequitas','YA'],['RELAYER_PRIVATE_KEY','Kunci privat (0x...) EOA yang menandatangani pendaftaran on-chain','YA'],['NODE_OPERATOR_WALLET','Wallet manusia terdaftar yang menerima hadiah validator harian','Untuk hadiah'],['RELAYER_ADDRESS','Alamat EOA yang cocok dengan RELAYER_PRIVATE_KEY. Ada fallback tapi atur secara eksplisit.','Direkomendasikan'],['PORT','Port HTTP untuk API + JSON-RPC. Default: 8080','TIDAK'],['PEER_SECRET','Rahasia bersama untuk mengotorisasi node sebagai validator. SEMUA node harus menggunakan nilai yang SAMA.','Multi-node'],['SELF_URL','URL HTTPS publik node ini. Diperlukan untuk self-exclusion di peer discovery.','Multi-node'],['PRIMARY_NODE_URL','Node primer untuk peer discovery otomatis. Atur ke https://aequitas.digital.','Multi-node'],['PEER_NODES','URL peer statis (lama). Gunakan PRIMARY_NODE_URL.','Opsional'],['NODE_KEY','Kunci libp2p base64 untuk identitas P2P stabil. Jika tidak diatur: auto-dibuat, cetak di stderr sebagai "SAVE THIS AS NODE_KEY: <base64>".','TIDAK'],['IS_PRIMARY_NODE','"true" hanya pada node primer. USANG: distribusi kini menggunakan DB-lock — tidak wajib.','TIDAK'],['RESET_STATE','"true" menghapus database saat startup. DESTRUKTIF.','TIDAK']],
      s4:'4. Mulai Cepat — Railway (Direkomendasikan)',r4:'Railway adalah cara tercepat untuk memulai. Tingkat gratis memenuhi persyaratan minimum untuk node. Perkiraan waktu: 10-15 menit.',rs:['Fork repo: https://github.com/hanoi96international-gif/Aequitas','Buat akun di railway.app dan mulai proyek baru','Klik "Deploy from GitHub Repo" dan pilih fork kamu','Di proyek: + New → Database → Add PostgreSQL','Buka layanan kamu → Variables dan tambahkan variabel dari Bagian 3','Atur PRIMARY_NODE_URL=https://aequitas.digital','Atur NODE_OPERATOR_WALLET=<wallet manusia AEQ kamu>','Atur RELAYER_PRIVATE_KEY=<kunci privat EOA kamu>','Klik "Deploy" — Dockerfile mengelola kompilasi (~3 menit)','Di log cari: "Aequitas Node Running" dan "[NODE] Registered node operator wallet"','Buka URL-KAMU/api/status untuk memastikan node aktif','Tambahkan RPC ke MetaMask: Chain ID 1926, Simbol AEQ, URL https://URL-KAMU/rpc'],rn:'Railway menetapkan subdomain acak; domain kustom dapat diatur di pengaturan proyek.',
      s5:'5. Mulai Cepat — Docker',d5:'Untuk VPS, VM cloud, atau server lokal. Prasyarat: Docker terinstal, PostgreSQL tersedia. NODE_OPERATOR_WALLET harus berupa wallet manusia yang terdaftar di Aequitas.',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\n\n# Buat image (~3 menit)\ndocker build -t aequitas-node .\n\n# Start pertama (NODE_KEY akan dicetak di log)\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n  -e RELAYER_PRIVATE_KEY="0xKUNCI_PRIVAT_KAMU" \\\n  -e RELAYER_ADDRESS="0xALAMAT_KAMU" \\\n  -e NODE_OPERATOR_WALLET="0xWALLET_MANUSIA_KAMU" \\\n  -e PEER_SECRET="dari-operator-jaringan" \\\n  -e SELF_URL="http://IP-SERVER-KAMU:8080" \\\n  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n  -e BOOTSTRAP_SNAPSHOT_URL="https://aequitas.digital/api/snapshot" \\\n  -e BOOTSTRAP_SIGNER="0x92cbedec9d348b4762cb9af99500ee6139c5b671" \\\n  -e SNAPSHOT_TOKEN="dari-operator-jaringan" \\\n  -p 8080:8080 -p 4001:4001 aequitas-node\n\n# Salin NODE_KEY dari log start pertama:\ndocker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY"\n\n# Restart dengan NODE_KEY untuk Peer ID stabil:\ndocker stop aequitas-node && docker rm aequitas-node\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e NODE_KEY="<base64 dari atas>" \\\n  ... (semua variabel lainnya) ...',dn:'Port 8080 harus dapat diakses dari luar (ufw allow 8080/tcp). Port 4001 opsional untuk P2P (ufw allow 4001/tcp). Sinkronisasi HTTP berfungsi tanpa P2P.',
      s6:'6. Verifikasi Node',v6:'Setelah berjalan, periksa endpoint ini untuk memastikan node tersinkronisasi.',vc:'curl https://URL-NODE-KAMU/api/status\n# Diharapkan: {"height": N, "total_humans": N}\n\ncurl -X POST https://URL-NODE-KAMU/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'',
      s7:'7. Jaringan P2P dan Sinkronisasi',p7:'Atur PRIMARY_NODE_URL=https://aequitas.digital di environment. Node terhubung dan menyinkronkan rantai penuh secara otomatis. Multialamat libp2p node primer:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'PRIMARY_NODE_URL (HTTPS) adalah metode yang direkomendasikan. Multialamat libp2p dapat berubah.',
      s8:'8. Mendapatkan Hadiah Validator',w8:'Hadiah berasal dari Pool Validator (40% dari semua biaya protokol). Langkah-langkah:',b8:['Pertama daftar sebagai manusia di Aequitas: instal aplikasi Android dan selesaikan pendaftaran biometrik','Atur NODE_OPERATOR_WALLET ke alamat wallet terdaftar tersebut','Mulai ulang node kamu dan konfirmasi di log: "[NODE] Registered node operator wallet: 0x..."','Node primer mendistribusikan hadiah setiap 24 jam ke semua wallet operator terdaftar','Node sekunder TIDAK perlu memicu distribusi — cukup jalankan node kamu'],
      s9:'9. Pemecahan Masalah',th:['Gejala','Kemungkinan Penyebab','Solusi'],tr:[['Tinggi tetap di 0','PRIMARY_NODE_URL tidak diatur atau salah','Atur PRIMARY_NODE_URL=https://aequitas.digital dan SELF_URL ke URL/IP publik kamu, lalu deploy ulang'],['"no code at address" di log','Kontrak V7 belum di-deploy','Normal saat pertama kali — node men-deploy V7 secara otomatis'],['Error DATABASE_URL','String koneksi salah','Periksa format: postgres://user:pass@host:5432/dbname'],['Port 8080 tidak dapat diakses','Firewall atau konfigurasi provider','ufw allow 8080/tcp — diperlukan agar node primer dapat menyinkronkan blok kamu'],['Bootstrap snapshot gagal / StateRoot mismatch','SNAPSHOT_TOKEN tidak diatur atau BOOTSTRAP_SIGNER salah','Atur BOOTSTRAP_SNAPSHOT_URL=https://aequitas.digital/api/snapshot, BOOTSTRAP_SIGNER=0x92cbedec9d348b4762cb9af99500ee6139c5b671, SNAPSHOT_TOKEN=<dari operator>. Restart — state diimpor otomatis.'],['NODE_KEY menghasilkan kunci baru setiap restart','Variabel NODE_KEY tidak diatur','Saat start pertama: docker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY". Salin nilai base64, atur sebagai NODE_KEY, restart sekali.'],['MetaMask menampilkan 0 AEQ atau saldo salah','Konfigurasi jaringan yang kedaluwarsa di MetaMask','MetaMask → Pengaturan → Jaringan → hapus semua entri "Aequitas Chain" → tambahkan kembali via tombol di website.'],['Node ditolak sebagai validator','NODE_OPERATOR_WALLET bukan wallet manusia terdaftar','Daftar di Aequitas terlebih dahulu (app Android), lalu atur NODE_OPERATOR_WALLET ke wallet terdaftar tersebut.']],
      s10:'10. Konfigurasi MetaMask',m10:'Untuk menggunakan node kamu sebagai endpoint RPC di MetaMask:',mh:['Kolom','Nilai'],mr:[['Nama Jaringan','Aequitas Chain'],['URL RPC','https://URL-NODE-KAMU/rpc'],['Chain ID','1926  (hex: 0x786)'],['Simbol','AEQ'],['Desimal','18'],['Block Explorer','https://aequitas.digital']],
      foot:'Open source · Tanpa izin · Tanpa kunci admin · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'},
    fr:{title:'Guide de l\'Operateur de Noeud Aequitas',sub:'Guide complet etape par etape · Aequitas Chain (Chain ID 1926)',badge:'v1.0 · Open Source · Sans permission · Aucun stake requis',
      s1:'1. Presentation',what:'Role d\'un noeud',wtxt:'Un noeud Aequitas participe pleinement au reseau : produit des blocs dans le consensus BlockDAG, valide les preuves biometriques Groth16 ZK pour les nouvelles inscriptions humaines, applique les plafonds de richesse et le demurrage, synchronise l\'etat avec les pairs via libp2p + HTTP et execute les distributions quotidiennes des pools. Chaque noeud execute la chaine complete — pas de clients legers.',
      earn:'Ce que vous gagnez',etxt:'Definissez NODE_OPERATOR_WALLET sur une adresse de portefeuille humain enregistre. Le Pool Validateurs accumule 40% de tous les frais de protocole. Toutes les 24h, le noeud principal distribue proportionnellement le solde du pool a tous les operateurs enregistres. Aucun stake requis.',
      s2:'2. Prerequis',rh:['Composant','Minimum','Recommande'],rr:[['OS','Linux / hote Docker','Ubuntu 22.04 LTS'],['RAM','512 Mo','1 Go'],['CPU','1 vCPU','2 vCPU'],['Stockage','2 Go','10 Go SSD'],['Base de donnees','PostgreSQL 14+','Railway ou Supabase'],['Reseau','IP publique / redirection de port','TCP 8080 ouvert']],
      s3:'3. Variables d\'Environnement',e3:'Definir ces variables avant de demarrer le noeud. Variables marquees OUI sont obligatoires.',eh:['Variable','Fonction','Requise?'],er:[['DATABASE_URL','Chaine de connexion PostgreSQL : postgres://user:pass@host:5432/aequitas','OUI'],['RELAYER_PRIVATE_KEY','Cle privee (0x...) de l\'EOA qui signe les inscriptions on-chain','OUI'],['NODE_OPERATOR_WALLET','Portefeuille humain enregistre qui recoit les recompenses de validateur quotidiennes','Pour recomp.'],['RELAYER_ADDRESS','Adresse EOA correspondant a RELAYER_PRIVATE_KEY. Fallback disponible mais a definir.','Recommande'],['PORT','Port HTTP pour API + JSON-RPC. Defaut : 8080','NON'],['PEER_NODES','URLs de pairs statiques (legacy). Utiliser PRIMARY_NODE_URL.','Optionnel'],['PEER_SECRET','Secret partage: TOUS les noeuds doivent utiliser la MEME valeur.','Multi-noeud'],['SELF_URL','URL HTTPS publique pour self-exclusion.','Multi-noeud'],['PRIMARY_NODE_URL','Noeud principal pour decouverte auto (https://aequitas.digital).','Multi-noeud'],['NODE_KEY','Cle libp2p base64 pour identite P2P stable. Auto-generee en stderr comme "SAVE THIS AS NODE_KEY: <base64>" si omise.','NON'],['IS_PRIMARY_NODE','"true" uniquement sur le noeud principal. OBSOLETE: les distributions utilisent DB-lock — non requis.','NON'],['RESET_STATE','"true" efface la BD au demarrage. DESTRUCTIF.','NON']],
      s4:'4. Demarrage Rapide — Railway (Recommande)',r4:'Railway est le moyen le plus rapide de commencer. Le niveau gratuit couvre les exigences minimales pour le lancement. Duree estimee : 10 a 15 minutes.',rs:['Forker le depot : https://github.com/hanoi96international-gif/Aequitas','Creer un compte sur railway.app et demarrer un nouveau projet','Cliquer sur "Deploy from GitHub Repo" et selectionner votre fork','Dans le projet : + New → Database → Add PostgreSQL','Aller dans votre service → Variables et ajouter les variables de la section 3','Definir PRIMARY_NODE_URL=https://aequitas.digital','Definir NODE_OPERATOR_WALLET=<votre portefeuille humain AEQ>','Definir RELAYER_PRIVATE_KEY=<votre cle privee EOA>','Cliquer sur "Deploy" — le Dockerfile gere la compilation (~3 min)','Verifier dans les logs : "Aequitas Node Running" et "[NODE] Registered node operator wallet"','Ouvrir VOTRE-URL/api/status pour confirmer que le noeud est actif','Ajouter le RPC a MetaMask : Chain ID 1926, Symbole AEQ, URL https://VOTRE-URL/rpc'],rn:'Railway attribue un sous-domaine aleatoire ; domaines personnalises configurables dans les parametres du projet.',
      s5:'5. Demarrage Rapide — Docker',d5:'Pour VPS, VM cloud ou serveur local. Prerequis : Docker installe, PostgreSQL disponible. NODE_OPERATOR_WALLET doit etre un portefeuille humain enregistre sur Aequitas.',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\n\n# Construire l\'image (~3 min)\ndocker build -t aequitas-node .\n\n# Premier demarrage (NODE_KEY sera affiche dans les logs)\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n  -e RELAYER_PRIVATE_KEY="0xVOTRE_CLE_PRIVEE" \\\n  -e RELAYER_ADDRESS="0xVOTRE_ADRESSE" \\\n  -e NODE_OPERATOR_WALLET="0xVOTRE_PORTEFEUILLE_HUMAIN" \\\n  -e PEER_SECRET="aupres-de-l-operateur-reseau" \\\n  -e SELF_URL="http://VOTRE-IP-SERVEUR:8080" \\\n  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n  -e BOOTSTRAP_SNAPSHOT_URL="https://aequitas.digital/api/snapshot" \\\n  -e BOOTSTRAP_SIGNER="0x92cbedec9d348b4762cb9af99500ee6139c5b671" \\\n  -e SNAPSHOT_TOKEN="aupres-de-l-operateur-reseau" \\\n  -p 8080:8080 -p 4001:4001 aequitas-node\n\n# Copier NODE_KEY depuis les logs du premier demarrage :\ndocker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY"\n\n# Redemarrer avec NODE_KEY pour un Peer ID stable :\ndocker stop aequitas-node && docker rm aequitas-node\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e NODE_KEY="<base64 ci-dessus>" \\\n  ... (toutes les autres variables) ...',dn:'Le port 8080 doit etre accessible de l\'exterieur (ufw allow 8080/tcp). Le port 4001 est optionnel pour P2P (ufw allow 4001/tcp). La synchronisation HTTP fonctionne sans P2P.',
      s6:'6. Verifier le Noeud',v6:'Une fois en cours d\'execution, verifier ces endpoints :',vc:'curl https://VOTRE-NOEUD-URL/api/status\n# Attendu : {"height": N, "total_humans": N}\n\ncurl -X POST https://VOTRE-NOEUD-URL/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'',
      s7:'7. Reseau P2P et Synchronisation',p7:'Definir PRIMARY_NODE_URL=https://aequitas.digital dans l\'environnement. Le noeud se connecte et synchronise automatiquement. Multiadresse libp2p du noeud principal :',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'PRIMARY_NODE_URL (HTTPS) est la methode recommandee. La multiadresse libp2p peut changer.',
      s8:'8. Gagner des Recompenses de Validateur',w8:'Les recompenses viennent du Pool Validateurs (40% de tous les frais de protocole). Etapes :',b8:['D\'abord s\'inscrire comme humain sur Aequitas : installer l\'app Android et completer l\'inscription biometrique','Definir NODE_OPERATOR_WALLET sur cette adresse de portefeuille enregistree','Redemarrer le noeud et confirmer dans les logs : "[NODE] Registered node operator wallet: 0x..."','Le noeud principal distribue les recompenses toutes les 24h a tous les operateurs enregistres','Les noeuds secondaires n\'ont PAS besoin de declencher la distribution'],
      s9:'9. Depannage',th:['Symptome','Cause probable','Solution'],tr:[['Hauteur reste a 0','PRIMARY_NODE_URL non definie ou incorrecte','Definir PRIMARY_NODE_URL=https://aequitas.digital et SELF_URL sur votre URL/IP publique, puis redeployer'],['"no code at address" dans les logs','Contrat V7 non encore deploye','Normal au premier demarrage — le noeud deploie V7 automatiquement'],['Erreur DATABASE_URL','Chaine de connexion incorrecte','Verifier le format : postgres://user:pass@host:5432/dbname'],['Port 8080 inaccessible','Pare-feu ou configuration du fournisseur','ufw allow 8080/tcp — necessite pour que le noeud primaire synchronise vos blocs'],['Bootstrap snapshot echoue / StateRoot mismatch','SNAPSHOT_TOKEN absent ou BOOTSTRAP_SIGNER incorrect','Definir BOOTSTRAP_SNAPSHOT_URL=https://aequitas.digital/api/snapshot, BOOTSTRAP_SIGNER=0x92cbedec9d348b4762cb9af99500ee6139c5b671, SNAPSHOT_TOKEN=<aupres operateur>. Redemarrer — l\'etat est importe automatiquement.'],['NODE_KEY genere une nouvelle cle a chaque redemarrage','Variable NODE_KEY non definie','Au premier demarrage : docker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY". Copier la valeur base64, la definir comme NODE_KEY, redemarrer une fois.'],['MetaMask affiche 0 AEQ ou mauvais solde','Configuration reseau obsolete dans MetaMask','MetaMask → Parametres → Reseaux → supprimer toutes les entrees "Aequitas Chain" → les rajouter via le bouton sur le site web.'],['Noeud refuse comme validateur','NODE_OPERATOR_WALLET n\'est pas un portefeuille humain enregistre','S\'inscrire d\'abord sur Aequitas (app Android), puis definir NODE_OPERATOR_WALLET sur ce portefeuille enregistre.']],
      s10:'10. Configuration MetaMask',m10:'Pour utiliser votre noeud comme endpoint RPC dans MetaMask :',mh:['Champ','Valeur'],mr:[['Nom du reseau','Aequitas Chain'],['URL RPC','https://VOTRE-NOEUD-URL/rpc'],['Chain ID','1926  (hex: 0x786)'],['Symbole','AEQ'],['Decimales','18'],['Explorateur','https://aequitas.digital']],
      foot:'Open source · Sans permission · Sans cle admin · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'},
    pt:{title:'Guia do Operador de Node Aequitas',sub:'Guia completo passo a passo · Aequitas Chain (Chain ID 1926)',badge:'v1.0 · Open Source · Sem permissao · Sem stake necessario',
      s1:'1. Visao Geral',what:'O que um node faz',wtxt:'Um node Aequitas participa totalmente da rede: produz blocos no consenso BlockDAG, valida provas biometricas Groth16 ZK para novos registros humanos, aplica tetos de riqueza e demurrage, sincroniza estado com peers via libp2p + HTTP e executa distribuicoes diarias dos pools. Cada node executa a cadeia completa — sem clientes leves.',
      earn:'O que voce ganha',etxt:'Defina NODE_OPERATOR_WALLET para um endereco de carteira humano registrado. O Pool de Validadores acumula 40% de todas as taxas do protocolo. A cada 24h o node principal distribui proporcionalmente o saldo do pool entre todos os operadores registrados. Sem stake necessario.',
      s2:'2. Requisitos',rh:['Componente','Minimo','Recomendado'],rr:[['OS','Linux / host Docker','Ubuntu 22.04 LTS'],['RAM','512 MB','1 GB'],['CPU','1 vCPU','2 vCPU'],['Armazenamento','2 GB','10 GB SSD'],['Banco de dados','PostgreSQL 14+','Railway ou Supabase'],['Rede','IP publico / redirecionamento de porta','TCP 8080 aberto']],
      s3:'3. Variaveis de Ambiente',e3:'Defina estas variaveis antes de iniciar o node. Variaveis marcadas SIM sao obrigatorias.',eh:['Variavel','Funcao','Necessaria?'],er:[['DATABASE_URL','String de conexao PostgreSQL: postgres://user:pass@host:5432/aequitas','SIM'],['RELAYER_PRIVATE_KEY','Chave privada (0x...) do EOA que assina registros on-chain','SIM'],['NODE_OPERATOR_WALLET','Carteira humana registrada que recebe recompensas de validador diarias','Para recomp.'],['RELAYER_ADDRESS','Endereco EOA correspondente a RELAYER_PRIVATE_KEY. Tem fallback mas defina explicitamente.','Recomendado'],['PORT','Porta HTTP para API + JSON-RPC. Padrao: 8080','NAO'],['PEER_SECRET','Segredo compartilhado para autorizar este node como validador. TODOS os nodes devem usar o MESMO valor.','Multi-node'],['SELF_URL','URL HTTPS publica deste node. Necessaria para self-exclusion no peer discovery.','Multi-node'],['PRIMARY_NODE_URL','Node principal para descoberta automatica de pares. Definir como https://aequitas.digital.','Multi-node'],['PEER_NODES','URLs de pares estaticos (legado). Usar PRIMARY_NODE_URL.','Opcional'],['NODE_KEY','Chave libp2p base64 para identidade P2P. Se omitida: auto-gerada em stderr como "SAVE THIS AS NODE_KEY: <base64>".','NAO'],['IS_PRIMARY_NODE','"true" apenas no node principal. OBSOLETO: distribuicoes usam DB-lock — nao necessario.','NAO'],['RESET_STATE','"true" apaga o BD na inicializacao. DESTRUTIVO.','NAO']],
      s4:'4. Inicio Rapido — Railway (Recomendado)',r4:'Railway e a forma mais rapida de comecar. O nivel gratuito atende os requisitos minimos para el lanzamiento. Tempo estimado: 10-15 minutos.',rs:['Fazer fork do repositorio: https://github.com/hanoi96international-gif/Aequitas','Criar conta em railway.app e iniciar novo projeto','Clicar em "Deploy from GitHub Repo" e selecionar seu fork','No projeto: + New → Database → Add PostgreSQL','Ir para seu servico → Variables e adicionar variaveis da Secao 3','Definir PRIMARY_NODE_URL=https://aequitas.digital','Definir NODE_OPERATOR_WALLET=<sua carteira humana AEQ>','Definir RELAYER_PRIVATE_KEY=<sua chave privada EOA>','Clicar em "Deploy" — o Dockerfile gerencia a compilacao (~3 min)','Verificar nos logs: "Aequitas Node Running" e "[NODE] Registered node operator wallet"','Abrir SUA-URL/api/status para confirmar que o node esta ativo','Adicionar RPC ao MetaMask: Chain ID 1926, Simbolo AEQ, URL https://SUA-URL/rpc'],rn:'Railway atribui subdominio aleatorio; dominios personalizados nas configuracoes do projeto.',
      s5:'5. Inicio Rapido — Docker',d5:'Para VPS, VM na nuvem ou servidor local. Prerequisitos: Docker instalado, PostgreSQL disponivel. NODE_OPERATOR_WALLET deve ser uma carteira humana registrada no Aequitas.',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\n\n# Criar imagem (~3 min)\ndocker build -t aequitas-node .\n\n# Primeiro inicio (NODE_KEY sera impresso nos logs)\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n  -e RELAYER_PRIVATE_KEY="0xSUA_CHAVE_PRIVADA" \\\n  -e RELAYER_ADDRESS="0xSEU_ENDERECO" \\\n  -e NODE_OPERATOR_WALLET="0xSUA_CARTEIRA_HUMANA" \\\n  -e PEER_SECRET="do-operador-da-rede" \\\n  -e SELF_URL="http://SEU-IP-SERVIDOR:8080" \\\n  -e PRIMARY_NODE_URL="https://aequitas.digital" \\\n  -e BOOTSTRAP_SNAPSHOT_URL="https://aequitas.digital/api/snapshot" \\\n  -e BOOTSTRAP_SIGNER="0x92cbedec9d348b4762cb9af99500ee6139c5b671" \\\n  -e SNAPSHOT_TOKEN="do-operador-da-rede" \\\n  -p 8080:8080 -p 4001:4001 aequitas-node\n\n# Copiar NODE_KEY dos logs do primeiro inicio:\ndocker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY"\n\n# Reiniciar com NODE_KEY para Peer ID estavel:\ndocker stop aequitas-node && docker rm aequitas-node\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e NODE_KEY="<base64 acima>" \\\n  ... (todas as outras variaveis) ...',dn:'A porta 8080 deve ser acessivel de fora (ufw allow 8080/tcp). A porta 4001 e opcional para P2P (ufw allow 4001/tcp). A sincronizacao HTTP funciona sem P2P.',
      s6:'6. Verificar o Node',v6:'Apos iniciar, verificar estes endpoints:',vc:'curl https://SEU-NODE-URL/api/status\n# Esperado: {"height": N, "total_humans": N}\n\ncurl -X POST https://SEU-NODE-URL/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'',
      s7:'7. Rede P2P e Sincronizacao',p7:'Definir PRIMARY_NODE_URL=https://aequitas.digital no ambiente. O node sincroniza automaticamente. Multiendereco libp2p do node principal:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'PRIMARY_NODE_URL (HTTPS) e o metodo recomendado. O multiendereco libp2p pode mudar.',
      s8:'8. Ganhar Recompensas de Validador',w8:'Recompensas vem do Pool de Validadores (40% de todas as taxas). Passos:',b8:['Primeiro registrar como humano no Aequitas: instalar app Android e completar registro biometrico','Definir NODE_OPERATOR_WALLET para esse endereco registrado','Reiniciar node e confirmar nos logs: "[NODE] Registered node operator wallet: 0x..."','Node principal distribui recompensas a cada 24h para todos os operadores registrados','Nodes secundarios NAO precisam acionar a distribuicao — apenas manter o node ativo'],
      s9:'9. Resolucao de Problemas',th:['Sintoma','Causa provavel','Solucao'],tr:[['Altura fica em 0','PRIMARY_NODE_URL nao definida ou incorreta','Definir PRIMARY_NODE_URL=https://aequitas.digital e SELF_URL para sua URL/IP publico, depois reimplantar'],['"no code at address" nos logs','Contrato V7 ainda nao implantado','Normal na primeira vez — node implanta V7 automaticamente'],['Erro DATABASE_URL','String de conexao incorreta','Verificar formato: postgres://user:pass@host:5432/dbname'],['Porta 8080 inacessivel','Firewall ou configuracao do provedor','ufw allow 8080/tcp — necessario para o node principal sincronizar seus blocos'],['Bootstrap snapshot falhou / StateRoot mismatch','SNAPSHOT_TOKEN ausente ou BOOTSTRAP_SIGNER incorreto','Definir BOOTSTRAP_SNAPSHOT_URL=https://aequitas.digital/api/snapshot, BOOTSTRAP_SIGNER=0x92cbedec9d348b4762cb9af99500ee6139c5b671, SNAPSHOT_TOKEN=<do operador>. Reiniciar — estado importado automaticamente.'],['NODE_KEY gera nova chave a cada reinicio','Variavel NODE_KEY nao definida','No primeiro inicio: docker logs aequitas-node 2>&1 | grep "SAVE THIS AS NODE_KEY". Copiar o valor base64, definir como NODE_KEY, reiniciar uma vez.'],['MetaMask mostra 0 AEQ ou saldo incorreto','Configuracao de rede desatualizada no MetaMask','MetaMask → Configuracoes → Redes → excluir todas as entradas "Aequitas Chain" → adicionar novamente via botao no site.'],['Node recusado como validador','NODE_OPERATOR_WALLET nao e carteira humana registrada','Registrar-se primeiro no Aequitas (app Android), depois definir NODE_OPERATOR_WALLET para essa carteira registrada.']],
      s10:'10. Configuracao MetaMask',m10:'Para usar seu node como endpoint RPC no MetaMask:',mh:['Campo','Valor'],mr:[['Nome da rede','Aequitas Chain'],['URL RPC','https://SEU-NODE-URL/rpc'],['Chain ID','1926  (hex: 0x786)'],['Simbolo','AEQ'],['Decimais','18'],['Explorador','https://aequitas.digital']],
      foot:'Open source · Sem permissao · Sem chaves admin · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'},
    tr:{title:'Aequitas Dugum Operatoru Rehberi',sub:'Adim adim tam rehber · Aequitas Chain (Chain ID 1926)',badge:'v1.0 · Acik Kaynak · Izinsiz · Stake gerekmiyor',
      s1:'1. Genel Bakis',what:'Bir dugum ne yapar',wtxt:'Bir Aequitas dugumu agda tam olarak yer alir: BlockDAG uzlasmasinda blok uretir, yeni insan kayitlari icin Groth16 ZK biyometrik kanitlari dogrular, servet tavanlarini ve demurrage\'i protokol seviyesinde uygular, libp2p + HTTP araciligiyla eslerle durum senkronize eder ve gunluk havuz dagitimlarini calistirir.',
      earn:'Ne kazanirsiniz',etxt:'NODE_OPERATOR_WALLET\'i kayitli bir insan cuzdani adresine ayarlayin. Dogrulayicilar Havuzu tum protokol ucretlerinin %40\'ini biriktirir. Her 24 saatte bir ana dugum havuz bakiyesini tum kayitli operatörlere orantili olarak dagitir. Stake gerekmiyor.',
      s2:'2. Gereksinimler',rh:['Bilesen','Minimum','Onerilir'],rr:[['OS','Linux / Docker destekli sunucu','Ubuntu 22.04 LTS'],['RAM','512 MB','1 GB'],['CPU','1 vCPU','2 vCPU'],['Depolama','2 GB','10 GB SSD'],['Veritabani','PostgreSQL 14+','Railway veya Supabase'],['Ag','Genel IP / port yonlendirme','TCP 8080 acik']],
      s3:'3. Ortam Degiskenleri',e3:'Dugumu baslatmadan once bu degiskenleri ayarlayin. EVET olarak isaretlenenler zorunludur.',eh:['Degisken','Amac','Gerekli?'],er:[['DATABASE_URL','PostgreSQL baglanti dizesi: postgres://user:pass@host:5432/aequitas','EVET'],['RELAYER_PRIVATE_KEY','On-chain kayitlari imzalayan EOA\'nin ozel anahtari (0x...)','EVET'],['NODE_OPERATOR_WALLET','Gunluk dogrulayici odullerini alan kayitli insan cuzdani adresi','Oduller icin'],['RELAYER_ADDRESS','RELAYER_PRIVATE_KEY ile eslesen EOA adresi. Yedegi var ama acikca ayarlayin.','Onerilir'],['PORT','API + JSON-RPC icin HTTP portu. Varsayilan: 8080','HAYIR'],['PEER_NODES','Statik peer adresleri (eski). PRIMARY_NODE_URL kullanin.','Opsiyonel'],['PEER_SECRET','TUM dugumler AYNI PEER_SECRET degerini kullanmalidir.','Cok dugumlu'],['SELF_URL','Dugumun HTTPS adresi (self-exclusion icin gerekli).','Cok dugumlu'],['PRIMARY_NODE_URL','Birincil dugum (https://aequitas.digital).','Cok dugumlu'],['NODE_KEY','Kararli P2P kimligi icin libp2p base64 anahtari. Atlanirsa stderr\'de "SAVE THIS AS NODE_KEY: <base64>" olarak olusturulur.','HAYIR'],['IS_PRIMARY_NODE','Yalnizca belirlenmis birincil dugumde "true". Gunluk dagitimi etkinlestirir.','HAYIR'],['RESET_STATE','"true" baslatmada veritabanini siler. YIKICI.','HAYIR']],
      s4:'4. Hizli Baslangic — Railway (Onerilir)',r4:'Railway en hizli baslangi yoludur. Ucretsiz plan BETA gereksinimlerini karsilar. Tahmini kurulum suresi: 10-15 dakika.',rs:['Depoyu fork\'layin: https://github.com/hanoi96international-gif/Aequitas','railway.app adresinde hesap olusturun ve yeni proje baslatın','"Deploy from GitHub Repo" butonuna tiklayin ve fork\'unuzu secin','Projede: + New → Database → Add PostgreSQL','Servisinize gidin → Variables ve Bolum 3\'teki degiskenleri ekleyin','PRIMARY_NODE_URL=https://aequitas.digital ayarlayin','NODE_OPERATOR_WALLET=<AEQ insan cuzdaniniz> ayarlayin (gunluk oduller icin)','RELAYER_PRIVATE_KEY=<EOA ozel anahtariniz> ayarlayin','"Deploy" butonuna tiklayin — Dockerfile derlemeyi yonetir (~3 dk)','Loglarda kontrol edin: "Aequitas Node Running" ve "[NODE] Registered node operator wallet"','DUGUM-URL/api/status acarak dugumun aktif oldugunu dogrulayin','MetaMask\'a RPC ekleyin: Chain ID 1926, Sembol AEQ, URL https://URL\'NIZI/rpc'],rn:'Railway rastgele bir alt alan adi atar; ozel alan adlari proje ayarlarindan yapilandirilabilir.',
      s5:'5. Hizli Baslangic — Docker',d5:'VPS, bulut VM veya yerel sunucu icin. On kosullar: Docker kurulu, PostgreSQL mevcut.',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\n\n# Imaji olustur (~3 dk)\ndocker build -t aequitas-node .\n\n# Dugumu calistir\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n  -e RELAYER_PRIVATE_KEY="0xOZEL_ANAHTARINIZ" \\\n  -e RELAYER_ADDRESS="0xADRESINIZ" \\\n  -e NODE_OPERATOR_WALLET="0xINSAN_CUZDAN" \\\n  -e PEER_NODES="https://aequitas.digital" \\\n  -p 8080:8080 aequitas-node\n\ndocker logs -f aequitas-node',dn:'Konteyner 8080 portunu acktirir. Guvenlik duvarinda TCP 8080 girisini acin.',
      s6:'6. Dugumu Dogrulama',v6:'Calistiktan sonra bu endpoint\'leri kontrol edin:',vc:'curl https://DUGUM-URL/api/status\n# Beklenen: {"height": N, "total_humans": N}\n\ncurl -X POST https://DUGUM-URL/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'',
      s7:'7. P2P Ag ve Senkronizasyon',p7:'PRIMARY_NODE_URL=https://aequitas.digital ortama ayarlayin. Dugum otomatik baglanilar ve zinciri senkronize eder. Ana dugum libp2p multiadresi:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'PRIMARY_NODE_URL (HTTPS) onerilen yontemdir. Libp2p multiadresi degisebilir.',
      s8:'8. Dogrulayici Odulleri Kazanma',w8:'Oduller Dogrulayicilar Havuzu\'ndan gelir (%40 protokol ucreti). Adimlar:',b8:['Once Aequitas\'ta insan olarak kayit olun: Android uygulamasini indirin ve biyometrik kaydı tamamlayin','NODE_OPERATOR_WALLET\'i o kayitli cüzdan adresine ayarlayin','Dugumu yeniden baslatin ve loglarda dogrulayin: "[NODE] Registered node operator wallet: 0x..."','Ana dugum her 24 saatte odulleri tum kayitli operatörlere dagitir','Ikincil dugumler dagitimi baslatmak zorunda DEGILDIR — sadece dugumunuzu calisir tutun'],
      s9:'9. Sorun Giderme',th:['Belirti','Olasilik Nedeni','Cozum'],tr:[['Yukseklik 0\'da kaliyor','PEER_NODES ayarlanmadi','PRIMARY_NODE_URL=https://aequitas.digital ayarlayin ve yeniden dagitın'],['"no code at address" loglarda','V7 sozlesmesi henuz dagitilmamis','Ilk baslatmada normal — dugum V7\'yi otomatik dagitir'],['DATABASE_URL hatasi','Yanlis baglanti dizesi','Format: postgres://user:pass@host:5432/dbname'],['8080 portu erisilebilir degil','Guvenlik duvari veya saglayici ayarlari','Guvenlik duvarinda TCP 8080 girisini acin']],
      s10:'10. MetaMask Yapilandirmasi',m10:'Kendi dugumunuzu MetaMask\'ta RPC endpoint olarak kullanmak icin:',mh:['Alan','Deger'],mr:[['Ag Adi','Aequitas Chain'],['RPC URL','https://DUGUM-URL/rpc'],['Chain ID','1926  (hex: 0x786)'],['Para Birimi Sembolu','AEQ'],['Ondalik','18'],['Blok Gezgini','https://aequitas.digital']],
      foot:'Acik kaynak · Izinsiz · Yonetici anahtari yok · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'}
  };
  var c=C[lang]||C['en'];
  var fn='aequitas-node-operator-guide-'+lang+'.pdf';

  // Cover page
  doc.setFillColor(6,9,26);doc.rect(0,0,210,297,'F');
  doc.setFillColor(245,166,35);doc.rect(0,0,210,3,'F');
  y=55;doc.setFont('helvetica','bold');doc.setFontSize(30);doc.setTextColor(245,166,35);
  doc.text('AEQUITAS',105,y,{align:'center'});y+=10;
  doc.setFontSize(8.5);doc.setTextColor(90,110,160);
  doc.text('PROOF OF HUMANITY · DECENTRALIZED HUMAN CURRENCY',105,y,{align:'center'});y+=28;
  doc.setFontSize(20);doc.setTextColor(230,235,255);
  var tl=doc.splitTextToSize(c.title,160);doc.text(tl,105,y,{align:'center'});y+=tl.length*11+8;
  doc.setFontSize(9.5);doc.setTextColor(100,80,200);
  var sl=doc.splitTextToSize(c.sub,150);doc.text(sl,105,y,{align:'center'});y+=sl.length*7+16;
  doc.setFillColor(22,14,65);doc.setDrawColor(100,70,220);doc.setLineWidth(0.5);
  doc.roundedRect(45,y-5,120,13,4,4,'FD');
  doc.setFont('helvetica','bold');doc.setFontSize(7.5);doc.setTextColor(139,92,246);
  doc.text(c.badge,105,y+3.5,{align:'center'});
  doc.setFont('helvetica','normal');doc.setFontSize(7.5);doc.setTextColor(55,65,95);
  doc.text(c.link,105,282,{align:'center'});
  var dStr=new Date().toLocaleDateString(lang==='de'?'de-DE':'en-US',{year:'numeric',month:'long',day:'numeric'});
  doc.text(dStr,105,288,{align:'center'});
  doc.setFillColor(245,166,35);doc.rect(0,294,210,3,'F');

  // Content pages
  doc.addPage();hdr();y=22;
  h1(c.s1);h2(c.what);tx(c.wtxt);h2(c.earn);tx(c.etxt);
  h1(c.s2);tbl(c.rh,c.rr,[45,55,74]);
  h1(c.s3);tx(c.e3);tbl(c.eh,c.er,[52,100,22]);
  h1(c.s4);tx(c.r4);bl(c.rs);tx(c.rn);
  h1(c.s5);tx(c.d5);cd(c.dc);tx(c.dn);
  h1(c.s6);tx(c.v6);cd(c.vc);
  h1(c.s7);tx(c.p7);cd(c.pa);tx(c.pn);
  h1(c.s8);tx(c.w8);bl(c.b8);
  if(c.s8b){h1(c.s8b);tx(c.w8b);bl(c.b8b);}
  h1(c.s9);tbl(c.th,c.tr,[52,60,62]);
  h1(c.s10);tx(c.m10);tbl(c.mh,c.mr,[45,129]);
  ck(20);y+=6;
  doc.setDrawColor(200,160,40);doc.setLineWidth(0.3);doc.line(MG,y,W-MG,y);y+=8;
  doc.setFont('helvetica','italic');doc.setFontSize(7.5);doc.setTextColor(140,110,40);
  doc.text(c.foot,105,y,{align:'center'});y+=6;
  doc.setFont('helvetica','normal');doc.setTextColor(100,75,185);doc.text(c.link,105,y,{align:'center'});
  var pc=doc.getNumberOfPages();
  for(var pi=2;pi<=pc;pi++){doc.setPage(pi);doc.setFont('helvetica','normal');doc.setFontSize(7);doc.setTextColor(160,160,160);doc.text((pi-1)+' / '+(pc-1),W-MG,290,{align:'right'});}
  doc.save(fn);
}
window.addEventListener('resize', () => {
  const gd = document.getElementById('gini-history-chart');
  if (gd && gd._data) drawGiniHistoryChart(gd._data);
  const n = parseInt(document.getElementById('idx-humans2')?.textContent || '0');
  if (n > 0) drawWcapSlideChart(n);
});

</script>
</body>
</html>`
