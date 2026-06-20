package keeper

const explorerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0">
<title>Aequitas — Proof of Humanity Chain</title>
<meta name="description" content="Aequitas Chain — a Proof of Humanity blockchain with built-in Universal Basic Income, demurrage, and wealth cap. Chain ID 1926.">
<meta name="theme-color" content="#080010">
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
  --bg:#080010;--card:#0D0820;--card2:#110B28;--border:rgba(139,92,246,0.18);
  --purple:#8B5CF6;--teal:#06B6D4;--neon:#00FFD1;--gold:#F5A623;
  --green:#10B981;--red:#F87171;--blue:#60A5FA;
  --text:#F0EAFF;--muted:#7C6FA0;
  --glow-purple:0 0 20px rgba(139,92,246,0.25);
  --glow-teal:0 0 20px rgba(6,182,212,0.25);
  --glow-strong:0 0 40px rgba(139,92,246,0.4);
  --grad:linear-gradient(135deg,var(--purple),var(--teal));
  --radius:14px;--radius-sm:8px;
}
*{box-sizing:border-box;margin:0;padding:0}
body{background:var(--bg);color:var(--text);font-family:var(--font-body);min-height:100vh;overflow-x:hidden;line-height:1.5;background-image:radial-gradient(ellipse 80% 50% at 20% 0%,rgba(139,92,246,0.12) 0%,transparent 60%),radial-gradient(ellipse 60% 40% at 80% 100%,rgba(6,182,212,0.08) 0%,transparent 60%)}
body::before{content:'';position:fixed;top:0;left:0;width:100%;height:100%;pointer-events:none;z-index:0;background-image:radial-gradient(1px 1px at 10% 20%,rgba(139,92,246,0.6) 0%,transparent 100%),radial-gradient(1px 1px at 30% 60%,rgba(6,182,212,0.4) 0%,transparent 100%),radial-gradient(1px 1px at 50% 10%,rgba(0,255,209,0.5) 0%,transparent 100%),radial-gradient(1px 1px at 70% 40%,rgba(139,92,246,0.4) 0%,transparent 100%),radial-gradient(1px 1px at 90% 70%,rgba(6,182,212,0.6) 0%,transparent 100%),radial-gradient(1px 1px at 55% 90%,rgba(139,92,246,0.5) 0%,transparent 100%),radial-gradient(1px 1px at 85% 15%,rgba(6,182,212,0.3) 0%,transparent 100%);animation:starFloat 20s ease-in-out infinite alternate}
@keyframes starFloat{0%{transform:translateY(0)}100%{transform:translateY(-8px)}}
header{background:rgba(8,0,16,0.85);backdrop-filter:blur(20px);border-bottom:1px solid rgba(139,92,246,0.2);padding:0 24px;position:sticky;top:0;z-index:100;display:flex;align-items:center;justify-content:space-between;height:60px;gap:10px;box-shadow:0 1px 30px rgba(139,92,246,0.1)}
header::before{content:'';position:absolute;top:0;left:0;right:0;height:2px;background:var(--grad);opacity:0.8}
.logo-wrap{display:flex;align-items:center;gap:12px;flex-shrink:0;position:relative;z-index:1}
.logo-icon{width:34px;height:34px;border-radius:9px;background:var(--grad);display:flex;align-items:center;justify-content:center;font-size:17px;box-shadow:var(--glow-purple)}
.logo-text{font-size:1rem;font-weight:900;letter-spacing:3px;background:var(--grad);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text}
.logo-sub{font-size:0.48rem;color:var(--muted);letter-spacing:2.5px;text-transform:uppercase}
.header-right{display:flex;gap:8px;align-items:center;position:relative;z-index:1}
.badge{display:flex;align-items:center;gap:5px;padding:5px 11px;border-radius:20px;font-size:0.58rem;letter-spacing:0.5px;font-weight:600}
.badge-live{background:rgba(0,255,209,0.08);border:1px solid rgba(0,255,209,0.2);color:var(--neon)}
.badge-dag{background:rgba(139,92,246,0.08);border:1px solid rgba(139,92,246,0.2);color:var(--purple)}
.pulse{width:5px;height:5px;border-radius:50%;background:var(--neon);box-shadow:0 0 6px var(--neon);animation:pulse 2s infinite}
@keyframes pulse{0%,100%{opacity:1;transform:scale(1)}50%{opacity:0.4;transform:scale(0.7)}}
.lang-sel{background:rgba(139,92,246,0.08);color:var(--muted);border:1px solid var(--border);border-radius:6px;padding:5px 10px;font-family:var(--font-body);font-size:0.62rem;outline:none;cursor:pointer}
.tabs{background:rgba(8,0,16,0.7);backdrop-filter:blur(10px);border-bottom:1px solid rgba(139,92,246,0.12);padding:0 24px;display:flex;overflow-x:auto;-webkit-overflow-scrolling:touch;scrollbar-width:none;gap:2px;position:relative;z-index:1}
.tabs::-webkit-scrollbar{display:none}
.tab{padding:16px 16px;font-size:0.65rem;color:var(--muted);cursor:pointer;border-bottom:2px solid transparent;letter-spacing:0.5px;font-weight:600;white-space:nowrap;transition:all 0.2s;flex-shrink:0}
.tab:hover{color:var(--text)}.tab.active{color:var(--purple);border-bottom-color:var(--purple);text-shadow:0 0 10px rgba(139,92,246,0.5)}
.tab-content{display:none;position:relative;z-index:1}.tab-content.active{display:block}
.hero{padding:20px 20px 0;position:relative;z-index:1}
.section-label{font-size:0.6rem;color:var(--muted);letter-spacing:3px;text-transform:uppercase;margin-bottom:14px;font-weight:600}
.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:1px;background:rgba(139,92,246,0.1);border:1px solid var(--border);border-radius:var(--radius);overflow:hidden;margin-bottom:20px;box-shadow:var(--glow-purple)}
.stat{background:var(--card);padding:20px 16px;position:relative;transition:all 0.2s;cursor:default}
.stat:hover{background:var(--card2);box-shadow:inset 0 0 30px rgba(139,92,246,0.08)}
.stat-accent{position:absolute;top:0;left:0;right:0;height:2px}
.stat-icon{font-size:1rem;margin-bottom:8px}
.stat-lbl{font-size:0.58rem;color:var(--muted);letter-spacing:1.5px;text-transform:uppercase;margin-bottom:6px;font-weight:500}
.stat-val{font-size:1.7rem;font-weight:900;line-height:1;margin-bottom:4px;font-family:var(--font-display)}
.stat-sub{font-size:0.57rem;color:var(--muted);line-height:1.5}
.c-green .stat-val{color:var(--neon);text-shadow:0 0 15px rgba(0,255,209,0.4)}.c-green .stat-accent{background:linear-gradient(90deg,var(--neon),transparent)}
.c-blue .stat-val{color:var(--teal);text-shadow:0 0 15px rgba(6,182,212,0.4)}.c-blue .stat-accent{background:linear-gradient(90deg,var(--teal),transparent)}
.c-gold .stat-val{color:var(--gold);text-shadow:0 0 15px rgba(245,166,35,0.4)}.c-gold .stat-accent{background:linear-gradient(90deg,var(--gold),transparent)}
.c-purple .stat-val{color:var(--purple);text-shadow:0 0 15px rgba(139,92,246,0.4)}.c-purple .stat-accent{background:linear-gradient(90deg,var(--purple),transparent)}
.c-teal .stat-val{color:var(--teal)}.c-teal .stat-accent{background:linear-gradient(90deg,var(--teal),transparent)}
.info-banner{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:20px;margin-bottom:20px;display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:20px;box-shadow:var(--glow-purple)}
.ib-icon{font-size:1.4rem;margin-bottom:8px}
.ib-title{font-size:0.7rem;color:var(--gold);font-weight:700;margin-bottom:8px;letter-spacing:0.5px}
.ib-text{font-size:0.65rem;color:var(--muted);line-height:1.8}
.main-grid{display:grid;grid-template-columns:1fr 310px;gap:16px;padding:0 20px 20px;position:relative;z-index:1}
@media(max-width:800px){.main-grid{grid-template-columns:1fr}.right-col{display:none}}
.section{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);overflow:hidden;box-shadow:var(--glow-purple)}
.sec-head{padding:14px 18px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;background:rgba(139,92,246,0.05)}
.sec-title{font-size:0.65rem;color:var(--muted);letter-spacing:1px;text-transform:uppercase;display:flex;align-items:center;gap:8px;font-weight:600}
.sec-dot{width:6px;height:6px;border-radius:50%;background:var(--neon);box-shadow:0 0 8px var(--neon)}
.sec-count{font-size:0.6rem;color:var(--muted);background:var(--card2);padding:3px 8px;border-radius:10px;border:1px solid var(--border)}
.sec-desc{padding:10px 18px;font-size:0.65rem;color:var(--muted);background:rgba(139,92,246,0.03);border-bottom:1px solid var(--border);line-height:1.7}
.block-item{padding:12px 18px;border-bottom:1px solid rgba(139,92,246,0.08);display:grid;grid-template-columns:60px 1fr auto;gap:10px;align-items:center;transition:all 0.15s}
.block-item:hover{background:rgba(139,92,246,0.05)}.block-item:last-child{border-bottom:none}
.block-num{font-size:0.8rem;font-weight:700;color:var(--purple);font-family:var(--font-mono);text-shadow:0 0 8px rgba(139,92,246,0.4)}
.block-hash{font-size:0.63rem;color:var(--muted);margin-bottom:2px;display:flex;align-items:center;gap:4px;flex-wrap:wrap;font-family:var(--font-mono)}
.block-parents{font-size:0.57rem;color:rgba(139,92,246,0.3)}
.block-right{text-align:right}
.block-humans{font-size:0.65rem;color:var(--gold);margin-bottom:2px;font-weight:600}
.block-time{font-size:0.57rem;color:var(--neon)}
.bm{background:rgba(139,92,246,0.1);color:var(--purple);font-size:0.53rem;padding:2px 6px;border-radius:4px;border:1px solid rgba(139,92,246,0.2)}
.bt{background:rgba(0,255,209,0.08);color:var(--neon);font-size:0.53rem;padding:2px 6px;border-radius:4px;border:1px solid rgba(0,255,209,0.15)}
.empty{padding:40px;text-align:center;color:var(--muted);font-size:0.7rem;line-height:2.5}
.right-col{display:flex;flex-direction:column;gap:12px;position:relative;z-index:1}
.ic{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:18px;box-shadow:var(--glow-purple)}
.ic-title{font-size:0.6rem;color:var(--purple);letter-spacing:1.5px;text-transform:uppercase;margin-bottom:14px;font-weight:600}
.ic-row{display:flex;justify-content:space-between;align-items:center;padding:8px 0;border-bottom:1px solid rgba(139,92,246,0.08)}
.ic-row:last-child{border-bottom:none}
.ic-key{font-size:0.63rem;color:var(--muted)}
.ic-val{font-size:0.63rem;color:var(--text);text-align:right;max-width:58%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;font-weight:500}
.ic-val.g{color:var(--neon)}.ic-val.b{color:var(--teal)}.ic-val.go{color:var(--gold)}.ic-val.p{color:var(--purple)}
.mm-card{background:rgba(6,182,212,0.05);border:1px solid rgba(6,182,212,0.15);border-radius:var(--radius);padding:16px;box-shadow:var(--glow-teal)}
.mm-title{font-size:0.6rem;color:var(--teal);letter-spacing:1.5px;margin-bottom:12px;font-weight:700;text-transform:uppercase}
.mm-row{display:flex;justify-content:space-between;padding:6px 0;border-bottom:1px solid rgba(6,182,212,0.08)}
.mm-row:last-child{border-bottom:none}
.mm-key{font-size:0.6rem;color:var(--muted)}.mm-val{font-size:0.6rem;color:var(--purple);font-family:var(--font-mono)}
.mm-btn{width:100%;margin-top:12px;padding:11px;background:var(--grad);color:#fff;border:none;border-radius:var(--radius-sm);cursor:pointer;font-family:var(--font-body);font-size:0.68rem;font-weight:700;letter-spacing:0.5px;transition:all 0.2s;box-shadow:var(--glow-purple)}
.mm-btn:hover{opacity:0.87;transform:translateY(-1px);box-shadow:var(--glow-strong)}
.phil-card{background:linear-gradient(135deg,rgba(139,92,246,0.1),rgba(6,182,212,0.05));border:1px solid rgba(139,92,246,0.2);border-radius:var(--radius);padding:22px;text-align:center;box-shadow:var(--glow-purple)}
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
.idx{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:24px;box-shadow:var(--glow-purple)}
.idx-title{font-size:0.6rem;color:var(--purple);letter-spacing:2px;text-transform:uppercase;margin-bottom:10px;font-weight:600}
.idx-desc{font-size:0.67rem;color:var(--muted);line-height:1.8;margin-bottom:16px}
.idx-big{font-size:2.8rem;font-weight:900;line-height:1;font-family:var(--font-display);background:var(--grad);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text}
.idx-lbl{font-size:0.6rem;color:var(--muted);margin-top:4px}
.bar-bg{height:8px;background:rgba(139,92,246,0.1);border-radius:6px;overflow:hidden;margin:14px 0 6px;border:1px solid rgba(139,92,246,0.1)}
.bar-fill{height:100%;border-radius:6px;background:linear-gradient(90deg,var(--neon),var(--teal),var(--purple),var(--gold),var(--red));transition:width 1.5s ease;box-shadow:0 0 8px rgba(139,92,246,0.4)}
.bar-lbl{display:flex;justify-content:space-between;font-size:0.55rem;color:var(--muted)}
.mrow{display:grid;grid-template-columns:repeat(2,1fr);gap:8px;margin-top:14px}
.mbox{background:var(--card2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:12px;text-align:center;transition:all 0.2s}
.mbox:hover{border-color:rgba(139,92,246,0.4);box-shadow:var(--glow-purple)}
.mval{font-size:1.15rem;font-weight:700;color:var(--teal);font-family:var(--font-display);text-shadow:0 0 10px rgba(6,182,212,0.3)}
.mlbl{font-size:0.57rem;color:var(--muted);margin-top:3px;font-weight:500}
.story{font-size:0.7rem;line-height:2;color:var(--muted)}
.story p{margin-bottom:14px}
.hlbox{background:rgba(139,92,246,0.05);border-left:3px solid var(--purple);border-radius:0 var(--radius-sm) var(--radius-sm) 0;padding:14px 18px;margin:16px 0;font-size:0.67rem;color:var(--text);line-height:1.9}
.ns{padding:20px;display:grid;grid-template-columns:1fr 1fr;gap:16px;position:relative;z-index:1}
@media(max-width:700px){.ns{grid-template-columns:1fr}}
.nc{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:20px;box-shadow:var(--glow-purple)}
.nc-title{font-size:0.6rem;color:var(--purple);letter-spacing:1.5px;text-transform:uppercase;margin-bottom:14px;font-weight:600}
.nbox{background:var(--card2);border-radius:var(--radius-sm);padding:14px;border:1px solid var(--border);margin-bottom:10px}
.nstat{display:flex;align-items:center;gap:6px;font-size:0.67rem;color:var(--neon);margin-bottom:5px;font-weight:600}
.ndot{width:7px;height:7px;border-radius:50%;background:var(--neon);box-shadow:0 0 8px var(--neon)}
.nurl{font-size:0.58rem;color:var(--muted);word-break:break-all;margin-bottom:3px;font-family:var(--font-mono)}
.ndesc{font-size:0.58rem;color:rgba(139,92,246,0.4)}
.spect{width:100%;border-collapse:collapse}
.spect td{padding:8px 0;border-bottom:1px solid rgba(139,92,246,0.08);font-size:0.63rem}
.spect tr:last-child td{border-bottom:none}
.spect td:first-child{color:var(--muted);width:45%}
.spect td:last-child{text-align:right;font-family:var(--font-mono);font-size:0.6rem;color:var(--purple)}
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
.rcard{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:20px;margin-bottom:14px;box-shadow:var(--glow-purple)}
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
    <div class="rhero-sub" data-i18n="reg-sub">Join the Aequitas network and receive your 1,000 AEQ Universal Basic Income grant. Registration is one-time, permanent, and completely gasless. No personal data is ever stored — only a cryptographic proof that you are a unique human being.</div>
    <a href="/download/app.apk" style="display:inline-flex;align-items:center;gap:10px;margin-top:18px;background:var(--grad);color:#fff;padding:13px 28px;border-radius:10px;font-size:0.75rem;font-weight:700;text-decoration:none;letter-spacing:0.5px;box-shadow:var(--glow-purple);transition:all 0.2s" onmouseover="this.style.opacity='0.87';this.style.transform='translateY(-2px)'" onmouseout="this.style.opacity='1';this.style.transform='translateY(0)'">
      <span style="font-size:1.1rem">📱</span>
      <span data-i18n="btn-download-app">DOWNLOAD AEQUITASBIO APP</span>
    </a>
    <div style="font-size:0.55rem;color:rgba(255,255,255,0.35);margin-top:8px">Android APK · direct download · BETA</div>
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
        <div style="font-size:0.67rem;font-weight:700;color:var(--purple);margin-bottom:4px" data-i18n="usp-c3-title">Just a smartphone</div>
        <div style="font-size:0.59rem;color:var(--muted);line-height:1.75" data-i18n="usp-c3-desc">No computer, no bank account, no ID document. An Android phone with a fingerprint sensor is all you need to join.</div>
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
    <div class="aonly-text" data-i18n="app-text">Proof of Humanity requires biometric verification on your personal device. Your fingerprint or face scan is processed exclusively by the Hardware Secure Element (HSE) inside your phone — raw biometric data never leaves your device, never touches any server, and is never stored anywhere. The app generates a Zero-Knowledge Proof that mathematically proves your uniqueness without revealing any personal information. Download the AequitasBio app, scan your biometrics, connect your MetaMask wallet, and your <strong style="color:var(--gold)">1,000 AEQ will be credited automatically</strong> within seconds.</div>
  </div>
  <div class="rsteps">
    <div class="rstep"><div class="snum">1</div><div class="stitle" data-i18n="s1t">Biometric Scan</div><div class="sdesc" data-i18n="s1d">Open the AequitasBio app and scan your fingerprint or use face recognition. Your biometric data is processed by your phone's Hardware Secure Element and never leaves your device.</div></div>
    <div class="rstep"><div class="snum">2</div><div class="stitle" data-i18n="s2t">ZK Proof Generation</div><div class="sdesc" data-i18n="s2d">A Groth16 Zero-Knowledge Proof is generated on our proof server. This cryptographically proves your uniqueness without revealing your identity — your hash is never exposed.</div></div>
    <div class="rstep"><div class="snum">3</div><div class="stitle" data-i18n="s3t">Connect Wallet</div><div class="sdesc" data-i18n="s3d">The app opens MetaMask on this page. Connect your Ethereum wallet — this is the address that will receive your 1,000 AEQ grant. The proof is cryptographically bound to your wallet.</div></div>
    <div class="rstep"><div class="snum">4</div><div class="stitle" data-i18n="s4t">1,000 AEQ Granted</div><div class="sdesc" data-i18n="s4d">Your registration is confirmed on the Aequitas BlockDAG within 6 seconds. 1,000 AEQ is credited to your wallet instantly, gasless. Your identity is permanently recorded as verified human.</div></div>
  </div>
  <div class="pbar" data-i18n="priv-bar">🔒 Hardware Secure Element · Groth16 Zero-Knowledge Proof · Biometric data never leaves your device · No gas fees · One registration per human · Permanent &amp; immutable</div>
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
</div>
</div>

<!-- EXPLORER -->
<div id="tab-explorer" class="tab-content">
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

<!-- HUMANS -->
<div id="tab-humans" class="tab-content">
<div class="hero">
  <div class="section-label" data-i18n="humans-title">Verified Humans on Aequitas Chain</div>
  <div class="info-banner">
    <div>
      <div class="ib-icon">🔒</div>
      <div class="ib-title" data-i18n="h-what">What is a Verified Human?</div>
      <div class="ib-text" data-i18n="h-what-t">A Verified Human is a wallet address cryptographically proven to belong to a unique living human being. The verification happens through biometric authentication on your personal device — your fingerprint or face scan unlocks a key pair stored in your phone's Hardware Secure Element. Only a mathematical proof of uniqueness is ever transmitted. Your biometric data never leaves your device, never touches a server, and is never stored anywhere.</div>
    </div>
    <div>
      <div class="ib-icon">🧮</div>
      <div class="ib-title" data-i18n="h-zkp">Zero-Knowledge Proof System</div>
      <div class="ib-text" data-i18n="h-zkp-t">Aequitas uses the Groth16 proving system on the BN128 (alt-bn128) elliptic curve — the same curve used by Ethereum and Zcash. A ZK proof allows one party to prove they know a secret without revealing the secret itself. In Aequitas, this means proving "I am a unique human" without revealing who you are or what your biometrics look like. Proof size: ~200 bytes. Verification time: ~10ms. The proof is generated client-side on the proof server after your device authenticates.</div>
    </div>
    <div>
      <div class="ib-icon">🛡</div>
      <div class="ib-title" data-i18n="h-sybil">Sybil Attack Prevention</div>
      <div class="ib-text" data-i18n="h-sybil-t">A Sybil attack is when one person creates multiple identities to gain an unfair advantage. Aequitas prevents this by storing a permanent keccak256 hash of each biometric identity. Attempting to register a second wallet with the same fingerprint is immediately rejected — the hash is already in the database. One human, one wallet, forever. <strong style="color:var(--gold)">⚠ Current limitation:</strong> The biometric hash is derived from a device key — switching phones creates a new hash. A physiological sensor (MAX30102 PPG heart-rate sensor) is planned to provide truly device-independent identity verification.</div>
    </div>
    <div>
      <div class="ib-icon">🌍</div>
      <div class="ib-title" data-i18n="h-global">Global Financial Inclusion</div>
      <div class="ib-text" data-i18n="h-global-t">1.4 billion adults worldwide have no bank account. Aequitas requires nothing more than an Android smartphone with a fingerprint or face sensor — a device over 3 billion people already own. No bank account, no credit card, no prior cryptocurrency, no government ID. Just being human is enough to participate in the Aequitas economy.</div>
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

<!-- SWAP -->
<div id="tab-swap" class="tab-content">
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

    <div style="margin-bottom:8px">
      <div style="font-size:0.6rem;color:var(--muted);text-transform:uppercase;letter-spacing:1px;margin-bottom:4px" data-i18n="swap-sell-label">Sell</div>
      <div style="display:flex;gap:6px">
        <button class="rbtn" id="swap-dir-a2t" onclick="setSwapDirection('aeq_to_tusd')" style="flex:1">AEQ</button>
        <button class="rbtn" id="swap-dir-t2a" onclick="setSwapDirection('tusd_to_aeq')" style="flex:1">tUSD</button>
      </div>
    </div>
    <input type="number" id="swap-amount" placeholder="Amount" oninput="updateFeeEstimate()" style="width:100%;padding:14px;border-radius:8px;border:1px solid var(--border);background:#0A1220;color:#E8EDF5;font-size:16px;margin-bottom:4px;box-sizing:border-box">
    <div class="pct-row" style="display:flex;gap:6px;margin-bottom:8px">
      <button class="rbtn pctbtn" onclick="setSwapPct(0.25)" style="flex:1;padding:8px;font-size:12px">25%</button>
      <button class="rbtn pctbtn" onclick="setSwapPct(0.5)" style="flex:1;padding:8px;font-size:12px">50%</button>
      <button class="rbtn pctbtn" onclick="setSwapPct(0.75)" style="flex:1;padding:8px;font-size:12px">75%</button>
      <button class="rbtn pctbtn" onclick="setSwapPct(1)" style="flex:1;padding:8px;font-size:12px">MAX</button>
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
</div>
</div>

<!-- INDEX -->
<div id="tab-index" class="tab-content">
<div class="is">
  <div class="idx" style="grid-column:1/-1">
    <div class="idx-title" data-i18n="idx-title">Aequitas Index — Real-Time Economic Equality Score</div>
    <div class="idx-desc" data-i18n="idx-desc">The Aequitas Index is derived from the <strong style="color:var(--teal)">Gini coefficient</strong> — the international standard for measuring wealth inequality, adopted by the World Bank, OECD, and UN. Unlike a simple richest-vs-poorest ratio, the Gini coefficient captures the <em style="color:var(--text)">entire distribution</em> across every verified human simultaneously, in a single number. <strong style="color:var(--neon)">0 = perfect equality</strong> (every wallet holds exactly the same AEQ). <strong style="color:var(--red)">100 = total concentration</strong> (one wallet holds all AEQ in existence). For context: Bitcoin Gini ≈ 0.85 (Index 85) · most unequal country on Earth (South Africa) ≈ 0.63 · Scandinavia ≈ 0.25. Aequitas targets Gini below 0.35 (Index below 35) at scale — comparable to the most equal developed economies — enforced automatically by the wealth cap and redistribution pools, no governance vote required.</div>
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
        <div style="font-size:0.64rem;color:var(--muted);line-height:1.9" data-i18n="gini-what-text">Developed by Italian statistician Corrado Gini (1912). It measures wealth distribution by comparing the actual balance distribution against a hypothetical perfectly equal baseline — visualized as the Lorenz curve. The coefficient equals the ratio of the area between the Lorenz curve and the diagonal of equality to the total area below that diagonal. Scale: 0 means every person holds identical wealth. 1 means one person holds all wealth in existence. Used by the World Bank, OECD, and UN to compare countries. Reference values: Bitcoin ≈ 0.85 · South Africa (world record) ≈ 0.63 · Brazil ≈ 0.53 · USA ≈ 0.41 · Germany ≈ 0.31 · Sweden ≈ 0.27 · Aequitas long-term target: Gini below 0.35 — comparable to Scandinavia and Germany, enforced by the wealth cap at scale (bootstrap: sliding cap 5×→25× per human).</div>
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
    <div class="idx-title">Gini Index History</div>
    <div style="font-size:0.63rem;color:var(--muted);margin-bottom:12px">Recorded after each UBI distribution. Shows how equality evolves as the network grows.</div>
    <canvas id="gini-history-chart" height="120" style="width:100%;border-radius:6px;background:var(--card2)"></canvas>
    <div id="gini-history-empty" style="display:none;text-align:center;padding:24px;color:var(--muted);font-size:0.63rem">No snapshots yet — first one saved after the next UBI distribution.</div>
  </div>
  <div class="idx" style="grid-column:1/-1">
    <div class="idx-title">Lorenz Curve — Wealth Distribution Across Humans</div>
    <div style="font-size:0.63rem;color:var(--muted);margin-bottom:12px">Each point = cumulative % of AEQ held by the poorest X% of humans. The diagonal = perfect equality. The further the curve bows below the diagonal, the higher the Gini.</div>
    <canvas id="lorenz-chart" height="220" style="width:100%;border-radius:6px;background:var(--card2)"></canvas>
  </div>
  <div class="idx" style="grid-column:1/-1">
    <div class="idx-title">Wealth Cap Multiplier — Bootstrap Slider</div>
    <div style="font-size:0.63rem;color:var(--muted);margin-bottom:12px">Formula: <code style="color:var(--teal)">max(5, min(N, 25))×</code> average AEQ balance. Each new human slides the cap up by 1×, until the 25th human locks it at 25× permanently.</div>
    <canvas id="wcap-slide-chart" height="90" style="width:100%;border-radius:6px;background:var(--card2)"></canvas>
  </div>
  <div class="idx" style="grid-column:1/-1">
    <div class="idx-title" data-i18n="story-title">The Story of Aequitas — Why This Exists</div>
    <div class="story" data-i18n="story-text"><p>The year is 2009. Satoshi Nakamoto releases Bitcoin. For the first time, value can transfer between any two people without a bank. A genuine revolution. But something goes wrong almost immediately.</p><p>Early miners accumulate millions of coins at almost zero cost. By 2021, the top 1% of Bitcoin addresses control over 90% of all Bitcoin. Bitcoin's estimated Gini coefficient exceeds 0.85 — higher than any country on Earth. The cryptocurrency that was supposed to democratize finance created the most extreme wealth concentration in human history.</p><p><span style="color:var(--gold)">Aequitas</span> — Latin for "fairness" and "equality" — was created to answer a single question: <em style="color:var(--gold)">"What would a cryptocurrency look like if designed from first principles to be fair to every human being?"</em></p><p>The answer is simple: <strong style="color:var(--text)">Money exists because people exist. Therefore, every person should have an equal share of money simply by virtue of being human.</strong></p><p>Aequitas implements this principle mathematically. Every verified human receives 1,000 AEQ. No mining, no staking, no early-adopter advantage. The wealth cap, demurrage, and redistribution pools ensure that inequality cannot accumulate indefinitely. The Gini coefficient and Aequitas Index are calculated on-chain in real time, and the protocol adjusts automatically.</p><p>The Aequitas network launched in June 2026. Currently in Phase 0 (Bootstrap). The goal: demonstrate that money can be distributed fairly, Gini coefficient held below 0.35 (comparable to the most equal developed nations), and financial inclusion achieved at global scale — without any central authority.</p><p><em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em></p></div>
  </div>
</div>
</div>

<!-- NETWORK -->
<div id="tab-network" class="tab-content">
<div class="ns">
  <div class="nc" style="grid-column:1/-1;background:linear-gradient(135deg,rgba(245,166,35,0.06),rgba(13,8,32,0.9));border-color:rgba(245,166,35,0.2)">
    <div class="nc-title" style="color:var(--gold)" data-i18n="run-node-title">Run Your Own Node — Help Secure the Network</div>
    <div style="font-size:0.67rem;color:var(--muted);line-height:1.9;margin-bottom:16px" data-i18n="run-node-desc">Anyone can run an Aequitas node — no permission, no stake, no application required. Nodes participate in block production, validate the human registry, and synchronize the BlockDAG. Node operators earn a share of protocol fees via the Validators Pool (40% of all swap fees, distributed daily). The more nodes that run, the more decentralized and resilient the network becomes.</div>
    <div style="display:flex;gap:12px;flex-wrap:wrap;margin-bottom:16px">
      <button onclick="generateNodeGuidePDF()" style="display:inline-flex;align-items:center;gap:8px;background:var(--gold);color:#06091A;padding:12px 20px;border-radius:8px;font-size:0.7rem;font-weight:700;cursor:pointer;border:none;font-family:var(--font-body);transition:opacity 0.2s" onmouseover="this.style.opacity=0.87" onmouseout="this.style.opacity=1">
        📄 Node Operator Guide (PDF)
      </button>
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
            <span style="font-size:0.52rem;background:rgba(245,166,35,0.12);border:1px solid rgba(245,166,35,0.3);color:var(--gold);padding:2px 8px;border-radius:10px">BETA v0.1</span>
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
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="color:var(--gold);font-weight:700;min-width:16px">4.</span><span><strong style="color:var(--text)">A dedicated node wallet:</strong> Your node needs its own Ethereum wallet to sign transactions. This is NOT your personal AEQ wallet. Install MetaMask (metamask.io), create a new account specifically for your node, then export its private key: MetaMask &rarr; click account icon &rarr; Account Details &rarr; Show Private Key &rarr; enter your MetaMask password &rarr; copy. Keep this key strictly private &mdash; treat it like a password.</span></div>
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
      <div style="background:rgba(220,50,50,0.06);border:1px solid rgba(220,50,50,0.2);border-radius:6px;padding:10px 14px;margin-bottom:12px;font-size:0.6rem"><span style="color:#f87171;font-weight:700">Security Warning: </span><span style="color:var(--muted)">Your RELAYER_PRIVATE_KEY is like a master password. Anyone who has it controls your node wallet. Never share it publicly, never paste it in chat or email. Use a dedicated wallet just for the node &mdash; not your personal AEQ wallet.</span></div>
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
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--muted);padding:8px">PEER_NODES</td>
          <td style="font-size:0.6rem;color:var(--teal);padding:8px">Recommended</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Set to: <span style="font-family:var(--font-mono);color:var(--neon)">https://aequitas.digital</span> &mdash; this is the primary node your node syncs from on startup. Without this your node starts at block 0 instead of the current chain height.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--muted);padding:8px">PORT</td>
          <td style="font-size:0.6rem;color:var(--muted);padding:8px">No</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Leave unset on Railway &mdash; Railway sets this automatically. Default is 8080.</td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--muted);padding:8px">NODE_KEY</td>
          <td style="font-size:0.6rem;color:var(--muted);padding:8px">No</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">32-byte hex string that gives your node a stable P2P identity. Auto-generated if omitted, but changes on every restart (peers temporarily lose your node). To generate one: <span style="font-family:var(--font-mono)">node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"</span></td>
        </tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08);background:rgba(0,0,0,0.1)">
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:var(--muted);padding:8px">IS_PRIMARY_NODE</td>
          <td style="font-size:0.6rem;color:var(--muted);padding:8px">No</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Leave unset or false. Only the official Aequitas primary node should set this to true. Setting it to true on a secondary node causes double pool distributions &mdash; do not do this.</td>
        </tr>
        <tr>
          <td style="font-size:0.61rem;font-family:var(--font-mono);color:#f87171;padding:8px">RESET_STATE</td>
          <td style="font-size:0.6rem;color:var(--muted);padding:8px">No</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">DANGEROUS: Setting this to true wipes your entire database on every restart. Development use only. Never in production.</td>
        </tr>
      </table>

      <!-- Step 4 Railway -->
      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">Step 4 &mdash; Deploy on Railway (Recommended)</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9;margin-bottom:12px">Railway is the easiest way to run your node &mdash; no server setup, no command line required. The free tier covers all BETA requirements. Total time: about 10&ndash;15 minutes.</div>
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
          <span style="color:var(--neon)">PEER_NODES</span> = <span style="color:var(--gold)">https://aequitas.digital</span>
        </div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="display:inline-flex;align-items:center;justify-content:center;background:rgba(139,92,246,0.2);color:var(--purple);font-weight:700;font-size:0.58rem;min-width:22px;height:22px;border-radius:50%">5</span><span>Click <strong style="color:var(--text)">Deploy</strong> (or save variables to trigger auto-redeploy). Build takes ~3 minutes while Go compiles the node binary.</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="display:inline-flex;align-items:center;justify-content:center;background:rgba(139,92,246,0.2);color:var(--purple);font-weight:700;font-size:0.58rem;min-width:22px;height:22px;border-radius:50%">6</span><span>Watch <strong style="color:var(--text)">Deploy Logs</strong>. Success looks like: <span style="font-family:var(--font-mono);color:var(--teal)">Aequitas Node Running V</span> and <span style="font-family:var(--font-mono);color:var(--teal)">[NODE] Registered node operator wallet: 0x...</span></span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:6px"><span style="display:inline-flex;align-items:center;justify-content:center;background:rgba(139,92,246,0.2);color:var(--purple);font-weight:700;font-size:0.58rem;min-width:22px;height:22px;border-radius:50%">7</span><span>Go to <strong style="color:var(--text)">Settings</strong> &rarr; <strong style="color:var(--text)">Networking</strong> &rarr; <strong style="color:var(--text)">Generate Domain</strong> to get your public URL</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px"><span style="display:inline-flex;align-items:center;justify-content:center;background:rgba(139,92,246,0.2);color:var(--purple);font-weight:700;font-size:0.58rem;min-width:22px;height:22px;border-radius:50%">8</span><span>Open <span style="font-family:var(--font-mono);color:var(--neon)">https://YOUR-URL/api/status</span> in your browser &mdash; you should see JSON with <strong style="color:var(--text)">height</strong> climbing every ~6 seconds</span></div>
      </div>

      <!-- Step 4b Docker -->
      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">Step 4b &mdash; Alternative: Deploy with Docker (Advanced)</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9;margin-bottom:8px">Use this if you have your own server (VPS, home server, cloud VM). Requires Docker installed and a PostgreSQL database (Railway or Supabase both work as external databases). Skip to Step 5 if you used Railway above.</div>
      <div style="font-family:var(--font-mono);background:rgba(0,0,0,0.3);border:1px solid rgba(139,92,246,0.15);border-radius:6px;padding:14px;margin-bottom:18px;font-size:0.61rem;line-height:2.2;overflow-x:auto">
        <span style="color:var(--muted)"># 1. Download the code</span><br>
        git clone https://github.com/hanoi96international-gif/Aequitas &amp;&amp; cd Aequitas<br><br>
        <span style="color:var(--muted)"># 2. Build the node image (takes ~3 min for Go compilation)</span><br>
        docker build -t aequitas-node .<br><br>
        <span style="color:var(--muted)"># 3. Start the node &mdash; replace all placeholder values</span><br>
        docker run -d --name aequitas-node --restart unless-stopped \<br>
        &nbsp;&nbsp;-e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \<br>
        &nbsp;&nbsp;-e RELAYER_PRIVATE_KEY="0xYOUR_PRIVATE_KEY" \<br>
        &nbsp;&nbsp;-e RELAYER_ADDRESS="0xYOUR_NODE_WALLET_ADDRESS" \<br>
        &nbsp;&nbsp;-e NODE_OPERATOR_WALLET="0xYOUR_HUMAN_WALLET" \<br>
        &nbsp;&nbsp;-e PEER_NODES="https://aequitas.digital" \<br>
        &nbsp;&nbsp;-p 8080:8080 \<br>
        &nbsp;&nbsp;aequitas-node<br><br>
        <span style="color:var(--muted)"># 4. Watch the live logs</span><br>
        docker logs -f aequitas-node
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
      <div style="background:rgba(0,220,170,0.05);border:1px solid rgba(0,220,170,0.15);border-radius:6px;padding:10px 14px;margin-bottom:18px;font-size:0.62rem;color:var(--muted)">The block height should match the primary node within 1&ndash;2 blocks within seconds of startup. If it stays at 0, check that PEER_NODES is set correctly and the primary node URL is reachable.</div>

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
      <div style="background:rgba(245,166,35,0.05);border:1px solid rgba(245,166,35,0.2);border-radius:6px;padding:12px 14px;margin-bottom:12px;font-size:0.62rem;color:var(--muted);line-height:1.9">The Validators Pool collects 40% of all protocol fees (swap fees, demurrage, wealth cap overflow). Every 24 hours the primary node distributes the pool balance to all registered node operator wallets proportionally. The more consistently your node runs, the larger your share.</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:2.2;margin-bottom:18px">
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:4px"><span style="color:var(--gold);font-weight:700;min-width:16px">1.</span><span>Make sure you are registered as a human on Aequitas. If not: install the Android app and complete biometric registration first. You will receive a wallet address and 1,000 AEQ.</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:4px"><span style="color:var(--gold);font-weight:700;min-width:16px">2.</span><span>Set <span style="font-family:var(--font-mono);color:var(--neon)">NODE_OPERATOR_WALLET</span> = your Aequitas human wallet address in your Railway Variables</span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:4px"><span style="color:var(--gold);font-weight:700;min-width:16px">3.</span><span>Save &mdash; Railway redeploys automatically. On Docker: <span style="font-family:var(--font-mono);color:var(--teal)">docker restart aequitas-node</span></span></div>
        <div style="display:flex;align-items:flex-start;gap:10px;margin-bottom:4px"><span style="color:var(--gold);font-weight:700;min-width:16px">4.</span><span>In your node logs, confirm: <span style="font-family:var(--font-mono);color:var(--teal)">[NODE] Registered node operator wallet: 0x...</span></span></div>
        <div style="display:flex;align-items:flex-start;gap:10px"><span style="color:var(--gold);font-weight:700;min-width:16px">5.</span><span>Rewards are distributed automatically every 24 hours. Just keep your node running &mdash; no further action needed.</span></div>
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
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">PEER_NODES not set or wrong</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Set PEER_NODES=https://aequitas.digital and redeploy</td>
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
        <tr>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Docker build fails with module error</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">No internet access during build</td>
          <td style="font-size:0.61rem;color:var(--muted);padding:8px">Docker build needs outbound internet to download Go modules. Railway handles this automatically.</td>
        </tr>
      </table>

      <!-- Footer -->
      <div style="font-size:0.58rem;color:var(--gold);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid rgba(245,166,35,0.2);padding-bottom:6px">Questions / Feedback</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9">Open an issue on <a href="https://github.com/hanoi96international-gif/Aequitas" target="_blank" style="color:var(--purple)">GitHub</a> or reach the Aequitas team via the repository. BETA feedback on node setup, performance, and documentation gaps is especially welcome. Download this guide as a PDF in your selected language using the button above.</div>
    </div>
  </div>
  <div class="nc" style="grid-column:1/-1">
    <div class="nc-title" data-i18n="nodes-title">Active Nodes — Current Network Topology</div>
    <div style="font-size:0.65rem;color:var(--muted);line-height:1.8;margin-bottom:12px" data-i18n="nodes-desc">The Aequitas network currently operates on two geographically distributed nodes. Both participate in block production, state synchronization, and API serving. They communicate peer-to-peer via libp2p and synchronize block state via HTTP. Both share access to the same PostgreSQL database for persistent state. The network is designed to support additional nodes — any third-party operator can join by setting the bootstrap peer address.</div>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:8px">
      <div class="nbox">
        <div class="nstat"><span class="ndot"></span><span data-i18n="node1">Node 1 — Railway (Primary)</span></div>
        <div class="nurl">aequitas-production-9fba.up.railway.app</div>
        <div class="ndesc" data-i18n="node1-desc">Primary API · Block producer · UBI distribution · P2P bootstrap · PostgreSQL · RPC for MetaMask</div>
        <div style="margin-top:6px;font-size:0.57rem;color:rgba(0,255,209,0.5)">IS_PRIMARY_NODE=true · Daily pool distributions</div>
      </div>
      <div class="nbox">
        <div class="nstat"><span class="ndot"></span><span data-i18n="node2">Node 2 — Render (Secondary)</span></div>
        <div class="nurl">aequitas-node-2.onrender.com</div>
        <div class="ndesc" data-i18n="node2-desc">Secondary API · Block producer · P2P peer · HTTP sync · Shared PostgreSQL state</div>
        <div style="margin-top:6px;font-size:0.57rem;color:rgba(139,92,246,0.5)">Redundancy · Geographic distribution</div>
      </div>
    </div>
  </div>
  <div class="nc">
    <div class="nc-title" data-i18n="bootstrap-title">Connect a New Node</div>
    <div style="font-size:0.63rem;color:var(--muted);line-height:1.8;margin-bottom:10px" data-i18n="bootstrap-desc">To run your own Aequitas node, set the PEER_NODES environment variable to the bootstrap node address below. Your node will automatically sync the full chain state and begin participating in block production.</div>
    <div style="font-size:0.6rem;color:var(--muted);margin-bottom:6px;letter-spacing:1px">LIBP2P MULTIADDRESS</div>
    <div class="bsbox">/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R</div>
    <div style="font-size:0.6rem;color:var(--muted);margin-top:10px;line-height:1.7">Set in your environment: <span style="color:var(--purple);font-family:var(--font-mono)">PEER_NODES=https://aequitas.digital</span></div>
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

<!-- PROTOCOL V7 -->
<div id="tab-protocol" class="tab-content">
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
    <div style="font-size:0.63rem;color:var(--muted);line-height:1.9">Currently the network runs on 1–2 nodes (Railway + experimental). Decentralization is a staged process:<br><br>
    <span style="color:var(--neon)">Phase 0 (now):</span> Single-operator bootstrapping. Trust is established through code transparency, not node count.<br>
    <span style="color:var(--blue)">Phase 1 (100+ humans):</span> Open node join — anyone can run a full node and earn validator rewards from the 40% pool.<br>
    <span style="color:var(--gold)">Phase 2 (1,000+ humans):</span> Minimum 10 independent node operators required. Node diversity enforced by smart contract.<br>
    <span style="color:var(--purple)">Phase 3 (10,000+ humans):</span> Fully decentralized BlockDAG. No single operator can censor or halt the chain.<br><br>
    The node operator guide (PDF) is available on the Network tab. Each new node operator earns from the 40% validator pool — the more nodes, the more resilient the network.</div>
  </div>
</div>
</div>

<script>
const PS = 'https://aequitas-proof-server-production.up.railway.app';
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
  'app-text':'Proof of Humanity requires biometric verification on your personal device. Your fingerprint or face scan is processed exclusively by the Hardware Secure Element inside your phone — raw biometric data never leaves your device, never touches any server. The app generates a Zero-Knowledge Proof that proves your uniqueness without revealing any personal information. Download AequitasBio, scan your biometrics, connect MetaMask, and your <strong style="color:var(--gold)">1,000 AEQ will be credited automatically</strong>.',
  's1t':'Biometric Scan','s1d':'Open AequitasBio app · scan fingerprint or face · Hardware Secure Element processes locally · biometric data never leaves your device',
  's2t':'ZK Proof Generation','s2d':'Groth16 Zero-Knowledge Proof is generated on the proof server · your uniqueness is verified cryptographically · your identity is never revealed',
  's3t':'Connect Wallet','s3d':'The app opens MetaMask on this page · connect your Ethereum wallet · the proof is cryptographically bound to your wallet address',
  's4t':'1,000 AEQ Granted','s4d':'Registration confirmed on Aequitas BlockDAG within 6 seconds · 1,000 AEQ credited instantly · your identity is permanently recorded as a verified human',
  'priv-bar':'🔒 Hardware Secure Element · Groth16 Zero-Knowledge Proof · Biometric data never leaves your device · No gas fees · One registration per human · Permanent &amp; immutable',
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
  'h-what':'What is a Verified Human?','h-what-t':'A Verified Human is a wallet address cryptographically proven to belong to a unique living human through biometric Zero-Knowledge Proof. Biometric data is never transmitted or stored — only the mathematical proof of uniqueness.',
  'h-zkp':'Zero-Knowledge Proof System','h-zkp-t':'Aequitas uses the Groth16 proving system over the BN128 elliptic curve. Proof size: ~200 bytes. Verification time: ~10ms. The proof mathematically demonstrates uniqueness without revealing any identifying information.',
  'h-sybil':'Sybil Attack Prevention','h-sybil-t':'Each biometric hash is stored permanently using keccak256. Attempting to register twice is immediately rejected. One human, one wallet, forever. ⚠ Test phase: current verification is device-bound. A physiological sensor (MAX30102 PPG) is planned for fully device-independent identity verification in a future update.',
  'h-global':'Global Financial Inclusion','h-global-t':'No bank account, no credit card, no prior cryptocurrency required. Just an Android smartphone with a fingerprint or face sensor. Aequitas is designed to be accessible to every human on Earth.',
  'reg-humans':'Registered Humans','h-desc':'Every address below has been verified as a unique human through biometric ZKP. Each received exactly 1,000 AEQ. The registry is permanent, immutable, and on-chain.',
  'no-humans':'No humans registered yet.\n\nDownload the Aequitas Android App and be the first human on the chain!',
  'reg-stats':'Registry Stats','total-humans':'Total Humans',
  'idx-title':'Aequitas Index — Real-Time Economic Equality Score',
  'idx-desc':'The Aequitas Index is derived from the <strong style="color:var(--teal)">Gini coefficient</strong> — the international standard for measuring wealth inequality, adopted by the World Bank, OECD, and UN. It captures the complete balance distribution across every verified human simultaneously. <strong style="color:var(--neon)">0 = perfect equality</strong> (every wallet holds the same AEQ). <strong style="color:var(--red)">100 = total concentration</strong> (one wallet holds all AEQ). Bitcoin Gini ≈ 0.85 (Index 85) · South Africa (world record) ≈ 0.63 · Scandinavia ≈ 0.27 · Aequitas long-term target: Gini below 0.35 (Index below 35) — comparable to the most equal developed economies, enforced by the wealth cap and redistribution pools.',
  'gini-what-title':'What is the Gini Coefficient?',
  'gini-what-text':'Developed by Italian statistician Corrado Gini (1912). Measures wealth distribution by comparing actual balances against a hypothetical perfectly equal baseline — visualized as the Lorenz curve. Scale: 0 (everyone holds the same) to 1 (one person holds everything). Used by World Bank, OECD, UN to compare countries. Reference values: Bitcoin ≈ 0.85 · South Africa (world record) ≈ 0.63 · USA ≈ 0.41 · Germany ≈ 0.31 · Sweden ≈ 0.27 · Aequitas long-term target: Gini below 0.35 — comparable to Scandinavia and Germany, enforced by the wealth cap at scale (bootstrap: sliding cap 5×→25× per human).',
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
  'story-text':'<p>The year is 2009. Satoshi Nakamoto releases Bitcoin. For the first time, value can transfer between any two people without a bank. A genuine revolution. But something goes wrong almost immediately.</p><p>Early miners accumulate millions of coins at almost zero cost. By 2021, the top 1% of Bitcoin addresses control over 90% of all Bitcoin. Bitcoin\'s estimated Gini coefficient exceeds 0.85 — higher than any country on Earth. The cryptocurrency that was supposed to democratize finance created the most extreme wealth concentration in human history.</p><p><span style="color:var(--gold)">Aequitas</span> — Latin for "fairness" and "equality" — was created to answer a single question: <em style="color:var(--gold)">"What would a cryptocurrency look like if designed from first principles to be fair to every human being?"</em></p><p>The answer is simple: <strong style="color:var(--text)">Money exists because people exist. Therefore, every person should have an equal share of money simply by virtue of being human.</strong></p><p>Aequitas implements this mathematically. Every verified human receives 1,000 AEQ. No mining, no staking, no early-adopter advantage. The wealth cap, demurrage, and redistribution pools ensure inequality cannot accumulate indefinitely. The protocol adjusts automatically as the network grows.</p><p>The Aequitas network launched in June 2026. Currently in Phase 0. The goal: demonstrate that money can be distributed fairly, Gini coefficient held below 0.35 (comparable to the most equal developed nations), and financial inclusion achieved at global scale — without any central authority.</p><p><em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em></p>',
  'nodes-title':'Active Nodes — Current Network Topology',
  'nodes-desc':'The Aequitas network currently operates on two geographically distributed nodes. Both participate in block production, state synchronization, and API serving. They communicate peer-to-peer via libp2p and synchronize block state via HTTP. Both share access to the same PostgreSQL database for persistent state. The network is designed to support additional nodes — any operator can join.',
  'node1':'Node 1 — Railway (Primary)','node1-desc':'Primary API · Block producer · UBI distribution · P2P bootstrap · PostgreSQL · RPC for MetaMask',
  'node2':'Node 2 — Render (Secondary)','node2-desc':'Secondary API · Block producer · P2P peer · HTTP sync · Shared PostgreSQL state',
  'run-node-title':'Run Your Own Node — Help Secure the Network',
  'run-node-desc':'Anyone can run an Aequitas node — no permission, no stake, no application required. Nodes participate in block production, validate the human registry, and synchronize the BlockDAG. Node operators earn a share of protocol fees via the Validators Pool (40% of all swap fees, distributed daily).',
  'bootstrap-title':'Connect a New Node','bootstrap-desc':'To run your own Aequitas node, set the PEER_NODES environment variable to the bootstrap node address below. Your node will automatically sync the full chain state and begin participating in block production.',
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
  'inf-title':'6. NO ALGORITHMIC INFLATION','inf-box':'The ONLY event that creates new AEQ: a new verified human registers<br><br>Total Supply = Verified Humans × 1,000 AEQ — always, exactly.'
},
de:{
  'logo-sub':'MENSCHLICHKEITSNACHWEIS','live':'LIVE',
  'tab-register':'🔐 Registrieren','tab-explorer':'🔍 Explorer','tab-humans':'👥 Menschen','tab-index':'📊 Index','tab-network':'🌐 Netzwerk','tab-protocol':'📜 Protokoll V7','tab-swap':'🔄 Tauschen',
  'reg-title':'🔐 Als verifizierter Mensch registrieren',
  'reg-sub':'Tritt dem Aequitas-Netzwerk bei und erhalte dein Universelles Grundeinkommen von 1.000 AEQ. Einmalig, permanent und vollständig gebührenfrei. Keine persönlichen Daten werden jemals gespeichert.',
  'app-title':'REGISTRIERUNG NUR ÜBER ANDROID-APP',
  'app-text':'Der Menschlichkeitsnachweis erfordert biometrische Verifizierung auf deinem Gerät. Fingerabdruck oder Gesichtserkennung werden ausschliesslich durch das Hardware Secure Element verarbeitet — rohe biometrische Daten verlassen niemals dein Gerät, berühren keinen Server. Die App erstellt einen Zero-Knowledge-Beweis der deine Einzigartigkeit mathematisch beweist ohne persönliche Informationen preiszugeben. Lade AequitasBio herunter, scanne deine Biometrie, verbinde MetaMask, und deine <strong style="color:var(--gold)">1.000 AEQ werden automatisch gutgeschrieben</strong>.',
  's1t':'Biometrischer Scan','s1d':'AequitasBio-App öffnen · Fingerabdruck oder Gesicht scannen · Hardware Secure Element verarbeitet lokal · biometrische Daten verlassen nie dein Gerät',
  's2t':'ZK-Beweis-Erzeugung','s2d':'Groth16 Zero-Knowledge-Beweis wird auf dem Proof-Server erzeugt · Einzigartigkeit wird kryptografisch verifiziert · deine Identität wird nie preisgegeben',
  's3t':'Wallet verbinden','s3d':'Die App öffnet MetaMask auf dieser Seite · verbinde deine Ethereum-Wallet · der Beweis ist kryptografisch an deine Wallet-Adresse gebunden',
  's4t':'1.000 AEQ gutgeschrieben','s4d':'Registrierung auf Aequitas BlockDAG innerhalb von 6 Sekunden bestätigt · 1.000 AEQ sofort gutgeschrieben · deine Identität ist dauerhaft als verifizierter Mensch gespeichert',
  'priv-bar':'🔒 Hardware Secure Element · Groth16 Zero-Knowledge-Beweis · Biometrische Daten verlassen nie dein Gerät · Keine Gasgebühren · Eine Registrierung pro Mensch · Permanent und unveränderlich',
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
  'h-what':'Was ist ein verifizierter Mensch?','h-what-t':'Ein verifizierter Mensch ist eine Wallet-Adresse die kryptografisch bewiesen hat einem einzigartigen lebenden Menschen zu gehören. Biometrische Daten werden nie übertragen oder gespeichert.',
  'h-zkp':'Zero-Knowledge-Beweissystem','h-zkp-t':'Aequitas verwendet das Groth16-System über der BN128-Kurve. Beweissgröße: ~200 Bytes. Verifizierungszeit: ~10ms. Der Beweis demonstriert mathematisch die Einzigartigkeit ohne Identität preiszugeben.',
  'h-sybil':'Sybil-Angriff-Prävention','h-sybil-t':'Jeder biometrische Hash wird dauerhaft mit keccak256 gespeichert. Doppelte Registrierungsversuche werden sofort abgelehnt. Ein Mensch, eine Wallet, für immer. ⚠ Testphase: Aktuelle Verifizierung ist gerätegebunden. Ein physiologischer Sensor (MAX30102 PPG) ist für vollständig geräteunabhängige Identifizierung geplant.',
  'h-global':'Globale finanzielle Inklusion','h-global-t':'Kein Bankkonto, keine Kreditkarte, keine Kryptowährung erforderlich. Nur ein Android-Smartphone mit Fingerabdruck- oder Gesichtssensor.',
  'reg-humans':'Registrierte Menschen','h-desc':'Jede Adresse wurde als einzigartiger Mensch durch biometrisches ZKP verifiziert. Jeder erhielt genau 1.000 AEQ. Dauerhaft, unveränderlich, on-chain.',
  'no-humans':'Noch keine Menschen registriert.\n\nLade die Aequitas Android App herunter und sei der erste Mensch auf der Chain!',
  'reg-stats':'Registrierungsstatistiken','total-humans':'Gesamtmenschen',
  'idx-title':'Aequitas-Index — Echtzeit-Wirtschaftsgleichheits-Score',
  'idx-desc':'Der Aequitas-Index wird aus dem <strong style="color:var(--teal)">Gini-Koeffizienten</strong> abgeleitet — dem internationalen Standard zur Messung wirtschaftlicher Ungleichheit, genutzt von Weltbank, OECD und UN. Er erfasst die vollständige Bilanzverteilung aller verifizierten Menschen gleichzeitig. <strong style="color:var(--neon)">0 = perfekte Gleichheit</strong> (jede Wallet hält gleich viel AEQ). <strong style="color:var(--red)">100 = totale Konzentration</strong> (eine Wallet hält alles). Bitcoin-Gini ≈ 0,85 (Index 85) · Südafrika (Weltrekord) ≈ 0,63 · Skandinavien ≈ 0,27 · Aequitas-Langzeitziel: Gini unter 0,35 (Index unter 35) — vergleichbar mit den gleichheitsstärksten Industrieländern, automatisch durchgesetzt durch den Vermögensobergrenze-Mechanismus.',
  'gini-what-title':'Was ist der Gini-Koeffizient?',
  'gini-what-text':'Entwickelt vom italienischen Statistiker Corrado Gini (1912). Misst die Vermögensverteilung durch Vergleich mit einer perfekt gleichen Verteilung — visualisiert als Lorenz-Kurve. Skala: 0 (alle halten gleich viel) bis 1 (eine Person hält alles). Genutzt von Weltbank, OECD, UN. Referenzwerte: Bitcoin ≈ 0,85 · Südafrika (Weltrekord) ≈ 0,63 · USA ≈ 0,41 · Deutschland ≈ 0,31 · Schweden ≈ 0,27 · Aequitas-Langzeitziel: Gini unter 0,35 — vergleichbar mit Skandinavien und Deutschland, durchgesetzt durch den Vermögensdeckel (Bootstrap: gleitender Deckel 5×→25× pro Mensch).',
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
  'bootstrap-title':'Neuen Node verbinden','bootstrap-desc':'Um einen eigenen Aequitas-Node zu betreiben, setze die PEER_NODES-Umgebungsvariable auf die unten stehende Bootstrap-Adresse. Dein Node synchronisiert automatisch den vollständigen Chain-Zustand und beginnt mit der Blockproduktion.',
  'tech-title':'Technische Spezifikationen','mm-config':'MetaMask-Konfiguration',
  'k-lang':'Sprache','k-src':'Quellcode','evm-yes':'Ja — JSON-RPC /rpc · MetaMask-kompatibel',
  'proto-label':'Aequitas V7 Protokoll — Technische Dokumentation',
  'ca-title':'Contract- & Netzwerk-Adressen','ca-text':'Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier (Groth16 On-Chain-Verifier): 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Haupt-Contract): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7 ist die einzige Wahrheitsquelle der gesamten Aequitas-Wirtschaft. Jedes AEQ-Guthaben, jede Menschenregistrierung, jede UBI-Auszahlung und jede Durchsetzung der Vermögensobergrenze wird durch diesen einen unveränderlichen Contract geregelt — deployed auf der Aequitas Chain, einer maßgeschneiderten EVM-kompatiblen Blockchain mit BlockDAG-Konsens. Es gibt keinen Admin-Schlüssel, keinen Upgrade-Proxy, keine Governance-Abstimmung die eine einzige Zeile seiner Logik ändern könnte. Der Code der heute läuft ist der Code der in zehn Jahren läuft.<br><br>Der BioVerifier-Contract empfängt Groth16-Zero-Knowledge-Beweise die vollständig auf dem Android-Gerät des Nutzers erzeugt werden. Er verifiziert mathematisch on-chain in ~10 ms dass ein neuer Registrierungskandidat ein einzigartiger lebender Mensch ist — ohne jemals seinen Namen, seine Identität oder seine biometrischen Daten zu erfahren. Das ist es was die gasfreie, investitionsfreie Registrierung möglich macht: Der Beweis ist das Einzige was das Gerät je verlässt.<br><br>Zusammen machen diese zwei Contracts etwas möglich das in keinem Währungssystem der Geschichte je existiert hat: eine Geldmenge deren Regeln — wie viel existiert, wer es bekommt, wie es umverteilt wird — von keiner Person, keinem Unternehmen und keiner Regierung je geändert werden können. Niemals.',
  'ib-poh':'Menschlichkeitsnachweis','ib-poh-t':'Jeder AEQ-Inhaber muss kryptographisch beweisen dass er ein einzigartiger lebender Mensch ist. Keine Bots, keine Unternehmen, keine KI, keine Duplikate. Biometrische Daten verlassen niemals dein Gerät — nur ein mathematischer Einzigartigkeitsbeweis wird übertragen. Das bedeutet: AEQ ist die erste Währung die ausschließlich menschlich ist.',
  'ib-fair':'Radikal faire Verteilung','ib-fair-t':'Jeder verifizierte Mensch erhält bei der Registrierung genau 1.000 AEQ — nicht mehr, nicht weniger. Kein Pre-Mining, keine Gründer-Zuteilung, keine Investorenrunden. Die Gesamtmenge ist immer und exakt gleich der Anzahl verifizierter Menschen multipliziert mit 1.000. Dies wird mathematisch erzwungen, nicht durch Richtlinien.',
  'ib-dag':'BlockDAG-Architektur','ib-dag-t':'Im Gegensatz zu traditionellen Blockchains wo nur ein Block pro Höhe existieren kann, verwendet Aequitas eine DAG-Struktur. Mehrere Blöcke können gleichzeitig von verschiedenen Nodes produziert und später in den DAG zusammengeführt werden. Dies ermöglicht höheren Durchsatz, niedrigere Latenz und eliminiert Einzelknoten-Engpässe. Merge-Ereignisse werden im Explorer mit einem speziellen Badge markiert.',
  'ib-gas':'Wirklich gebührenfrei','ib-gas-t':'Alle Registrierungen und AEQ-Übertragungen kosten absolut nichts. Kein ETH, BNB oder MATIC erforderlich. Keine Kreditkarte, kein Bankkonto, keine vorherige Kryptowährung nötig. Der Relayer übernimmt alle Transaktionskosten. Wenn du ein Mensch mit einem Smartphone bist, kannst du teilnehmen — unabhängig von deiner wirtschaftlichen Situation.',
  'h-what':'Was ist ein verifizierter Mensch?','h-what-t':'Ein verifizierter Mensch ist eine Wallet-Adresse, die kryptographisch bewiesen gehört zu einem einzigartigen lebenden Menschen. Die Verifikation erfolgt durch biometrische Authentifizierung auf deinem Gerät — dein Fingerabdruck oder Gesicht entsperrt ein Schlüsselpaar das im Hardware-Sicherheitselement deines Telefons gespeichert ist. Nur ein mathematischer Einzigartigkeitsbeweis wird übertragen. Deine biometrischen Daten verlassen niemals dein Gerät, berühren keinen Server und werden nirgends gespeichert.',
  'h-zkp':'Zero-Knowledge-Proof-System','h-zkp-t':'Aequitas verwendet das Groth16-Beweissystem auf der BN128-Kurve — dieselbe Kurve wie Ethereum und Zcash. Ein ZK-Beweis ermöglicht es einer Partei zu beweisen dass sie ein Geheimnis kennt, ohne das Geheimnis selbst preiszugeben. Bei Aequitas bedeutet dies: "Ich bin ein einzigartiger Mensch" zu beweisen ohne zu enthüllen wer du bist. Beweisgrße: ~200 Byte. Verifikationszeit: ~10ms.',
  'h-sybil':'Sybil-Angriff-Prävention','h-sybil-t':'Ein Sybil-Angriff ist wenn eine Person mehrere Identitäten erstellt um einen unfairen Vorteil zu erlangen. Aequitas verhindert dies durch permanentes Speichern eines keccak256-Hashes jeder biometrischen Identität. Der Versuch eine zweite Wallet mit demselben Fingerabdruck zu registrieren wird sofort abgelehnt. Ein Mensch, eine Wallet, für immer. <strong style="color:var(--gold)">⚠ Aktuelle Einschränkung:</strong> Der biometrische Hash ist gerätegebunden. Ein Sensor (MAX30102 PPG) ist für wirklich geräteunabhängige Verifikation geplant.',
  'h-global':'Globale finanzielle Inklusion','h-global-t':'1,4 Milliarden Erwachsene weltweit haben kein Bankkonto. Aequitas benötigt nur ein Android-Smartphone mit einem Fingerabdruck- oder Gesichtssensor — ein Gerät das über 3 Milliarden Menschen bereits besitzen. Kein Bankkonto, keine Kreditkarte, keine vorherige Kryptowährung, kein Personalausweis. Einfach Mensch zu sein reicht aus.',
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
  'usp-c3-title':'Nur ein Smartphone nötig','usp-c3-desc':'Kein Computer, kein Bankkonto, kein Personalausweis. Ein Android-Gerät mit Fingerabdrucksensor reicht aus um dem Netzwerk beizutreten.',
  'usp-c4-title':'Täglich UBI empfangen','usp-c4-desc':'Nach der Registrierung erhältst du automatisch täglich einen Anteil der UBI-Ausschüttung — jeden Tag, ohne Aktion, solange du AEQ hältst.',
  'v7-intro-title':'Was ist AequitasV7?',
  'v7-intro-text':'AequitasV7 ist der zentrale Smart Contract des Aequitas-Protokolls. "V7" steht für die 7. Hauptversion des Fairness-Contracts — das Ergebnis iterativer Designverbesserung. Er ist unveränderlich auf der Aequitas Chain (Chain ID 1926) deployed und regelt jeden Aspekt des Protokolls: Menschenregistrierung, ZK-Beweisverifizierung, Guthabenverwaltung, Vermögensobergrenze, UBI-Ausschüttung, Swap-Gebühren und alle Governance-Parameter. Kein Admin kann den Contract upgraden oder ersetzen — er ist das unveränderliche Gesetz der Aequitas-Wirtschaft.'
},
es:{
  'logo-sub':'PRUEBA DE HUMANIDAD','live':'EN VIVO',
  'tab-register':'🔐 Registrar','tab-explorer':'🔍 Explorador','tab-humans':'👥 Humanos','tab-index':'📊 Índice','tab-network':'🌐 Red','tab-protocol':'📜 Protocolo V7','tab-swap':'🔄 Intercambiar',
  'reg-title':'🔐 Regístrate como Humano Verificado',
  'reg-sub':'Únete a la red Aequitas y recibe tu subsidio de Renta Básica Universal de 1,000 AEQ. Único, permanente y completamente gratuito. Ningún dato personal es almacenado.',
  'app-title':'REGISTRO SOLO VÍA APP ANDROID',
  'app-text':'La Prueba de Humanidad requiere verificación biométrica en tu dispositivo. Tu huella o reconocimiento facial se procesa exclusivamente por el Elemento Seguro de Hardware — los datos biométricos nunca salen de tu dispositivo. La app genera una Prueba de Conocimiento Cero que demuestra tu unicidad sin revelar información personal. Descarga AequitasBio, escanea tu biometría, conecta MetaMask, y tus <strong style="color:var(--gold)">1,000 AEQ serán acreditados automáticamente</strong>.',
  's1t':'Escaneo Biométrico','s1d':'Abrir AequitasBio · escanear huella o cara · el Elemento Seguro de Hardware procesa localmente · los datos biométricos nunca salen del dispositivo',
  's2t':'Generación de Prueba ZK','s2d':'La Prueba Groth16 de Conocimiento Cero se genera en el servidor · tu unicidad se verifica criptográficamente · tu identidad nunca se revela',
  's3t':'Conectar Wallet','s3d':'La app abre MetaMask en esta página · conecta tu wallet Ethereum · la prueba está criptográficamente vinculada a tu dirección',
  's4t':'1,000 AEQ Acreditados','s4d':'Registro confirmado en el BlockDAG de Aequitas en 6 segundos · 1,000 AEQ acreditados instantáneamente · tu identidad queda permanentemente registrada',
  'priv-bar':'🔒 Elemento Seguro de Hardware · Prueba ZK Groth16 · Datos biométricos nunca salen del dispositivo · Sin tarifas de gas · Un registro por humano · Permanente e inmutable',
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
  'h-what':'¿Qué es un Humano Verificado?','h-what-t':'Un Humano Verificado es una dirección wallet demostrada criptográficamente que pertenece a un humano único vivo mediante Prueba de Conocimiento Cero. Los datos biométricos nunca se transmiten ni almacenan.',
  'h-zkp':'Sistema de Prueba ZK','h-zkp-t':'Aequitas usa Groth16 sobre BN128. Tamaño de prueba: ~200 bytes. Verificación: ~10ms. La prueba demuestra matemáticamente la unicidad sin revelar información identificable.',
  'h-sybil':'Prevención de Ataques Sybil','h-sybil-t':'Cada hash biométrico se almacena permanentemente con keccak256. Intentar registrarse dos veces se rechaza inmediatamente. ⚠ Fase de prueba: verificación vinculada al dispositivo. Se planea sensor MAX30102 PPG para identificación independiente.',
  'h-global':'Inclusión Financiera Global','h-global-t':'Sin cuenta bancaria, tarjeta de crédito ni criptomoneda previa. Solo un smartphone Android con sensor biométrico.',
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
  'usp-c3-title':'Solo un smartphone','usp-c3-desc':'Sin ordenador, sin cuenta bancaria, sin documento de identidad. Un Android con lector de huella es suficiente.',
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
  'bootstrap-title':'Conectar un Nuevo Nodo','bootstrap-desc':'Para ejecutar tu propio nodo, establece la variable de entorno PEER_NODES a la dirección de bootstrap. Tu nodo sincronizará automáticamente el estado completo de la cadena.',
  'tech-title':'Especificaciones Técnicas','mm-config':'Configuración MetaMask',
  'k-lang':'Idioma','k-src':'Código Fuente','evm-yes':'Sí — JSON-RPC /rpc · Compatible con MetaMask',
  'proto-label':'Protocolo Aequitas V7 — Documentación Técnica',
  'ca-title':'Contratos y Direcciones de Red','ca-text':'Cadena: Aequitas Chain (ID: 1926 · 0x786)<br>RPC: https://aequitas.digital/rpc<br><br>BioVerifier (verificador Groth16 on-chain): 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (contrato principal): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'ca-desc':'AequitasV7 es la única fuente de verdad para toda la economía Aequitas. Cada saldo AEQ, cada registro humano, cada pago UBI y cada aplicación del límite de riqueza está gobernado por este único contrato inmutable — desplegado en Aequitas Chain, una blockchain personalizada compatible con EVM que ejecuta un motor de consenso BlockDAG. No hay clave de administrador, no hay proxy de actualización, no hay votación de gobernanza que pueda cambiar una sola línea de su lógica. El código que funciona hoy es el código que funcionará en diez años.<br><br>El contrato BioVerifier recibe pruebas de conocimiento cero Groth16 generadas completamente en el dispositivo Android del usuario. Verifica matemáticamente on-chain en ~10 ms que un nuevo registrante es un ser humano único y vivo — sin conocer jamás su nombre, identidad o datos biométricos. Esto es lo que hace posible el registro sin gas y sin inversión: la prueba es lo único que sale del dispositivo.<br><br>Juntos, estos dos contratos hacen posible algo que nunca ha existido en ningún sistema monetario de la historia: una oferta monetaria cuyas reglas — quién la recibe, cuánto existe, cómo se redistribuye — no puede ser alterada por ninguna persona, empresa o gobierno. Jamás.',
  'ib-poh':'Prueba de Humanidad','ib-poh-t':'Cada titular de AEQ debe probar criptográficamente que es un ser humano único y vivo. Sin bots, sin corporaciones, sin IA, sin duplicados. Los datos biométricos nunca salen de tu dispositivo — solo se transmite una prueba matemática de unicidad. AEQ es la primera moneda que es exclusivamente humana.',
  'ib-fair':'Distribución Radicalmente Justa','ib-fair-t':'Cada humano verificado recibe exactamente 1.000 AEQ al registrarse. Sin pre-minado, sin asignación a fundadores, sin rondas de inversores. El suministro total es siempre y exactamente igual al número de humanos verificados multiplicado por 1.000. Esto se aplica matemáticamente, no por política.',
  'ib-dag':'Arquitectura BlockDAG','ib-dag-t':'A diferencia de las blockchains tradicionales donde solo puede existir un bloque por altura, Aequitas usa una estructura DAG. Múltiples bloques pueden producirse simultáneamente por diferentes nodos y luego fusionarse en el DAG. Esto permite mayor rendimiento, menor latencia y elimina cuellos de botella. Los eventos de fusión se marcan con una insignia especial en el explorador.',
  'ib-gas':'Verdaderamente Sin Gas','ib-gas-t':'Todos los registros y transferencias de AEQ no cuestan absolutamente nada. No se necesita ETH, BNB ni MATIC. Sin tarjeta de crédito, sin cuenta bancaria, sin criptomoneda previa. El relayer cubre todos los costos de transacción. Si eres humano con un smartphone, puedes participar independientemente de tu situación económica.',
  'h-what':'¿Qué es un Humano Verificado?','h-what-t':'Un Humano Verificado es una dirección wallet demostrada criptográficamente que pertenece a un ser humano único y vivo. La verificación ocurre a través de autenticación biométrica en tu dispositivo personal — tu huella o cara desbloquea un par de claves almacenado en el Elemento Seguro de Hardware. Solo se transmite una prueba matemática de unicidad. Tus datos biométricos nunca salen de tu dispositivo, nunca tocan un servidor y nunca se almacenan.',
  'h-zkp':'Sistema de Prueba ZK','h-zkp-t':'Aequitas usa el sistema de prueba Groth16 en la curva BN128 — la misma curva que Ethereum y Zcash. Una prueba ZK permite probar que conoces un secreto sin revelar el secreto mismo. En Aequitas significa probar "soy un humano único" sin revelar quién eres. Tamaño de prueba: ~200 bytes. Tiempo de verificación: ~10ms.',
  'h-sybil':'Prevención de Ataques Sybil','h-sybil-t':'Un ataque Sybil es cuando una persona crea múltiples identidades para obtener ventaja injusta. Aequitas lo previene almacenando permanentemente un hash keccak256 de cada identidad biométrica. Intentar registrar una segunda wallet con la misma huella se rechaza inmediatamente. Un humano, una wallet, para siempre. <strong style="color:var(--gold)">⚠ Limitación actual:</strong> El hash biométrico está vinculado al dispositivo. Un sensor MAX30102 PPG está planificado para verificación verdaderamente independiente del dispositivo.',
  'h-global':'Inclusión Financiera Global','h-global-t':'1.400 millones de adultos en todo el mundo no tienen cuenta bancaria. Aequitas solo requiere un smartphone Android con sensor biométrico — un dispositivo que más de 3.000 millones de personas ya poseen. Sin cuenta bancaria, sin tarjeta de crédito, sin criptomoneda previa, sin documento de identidad. Simplemente ser humano es suficiente.',
  'poa-title':'1. PRUEBA DE VIDA — Recuperación de Saldos Inactivos','poa-text':'<p>¿Qué pasa con AEQ cuando las personas mueren o quedan permanentemente incapacitadas? En Bitcoin, las wallets perdidas significan suministro perdido permanentemente. Aequitas soluciona esto mediante un sistema de recuperación por inactividad de múltiples etapas: si una wallet no muestra actividad durante un período prolongado, su saldo se devuelve gradualmente a la comunidad a través del pool UBI.</p>',
  'poa-box':'Año 0–2: Uso normal — sin restricciones<br>Año 2: Aviso 1 — el Guardian puede responder en nombre<br>Año 2+60d: Aviso 2 — urgencia creciente<br>Año 2+120d: Aviso 3 — aviso final<br>Año 2+180d: AEQ movido a CUSTODIA personal (aún recuperable)<br>Año 4: Si aún inactivo — CUSTODIA liberada al Pool UBI',
  'guard-title':'2. SISTEMA GUARDIAN — Salvaguarda Humana','guard-text':'<p>¿Y si alguien está hospitalizado, encarcelado o de algún modo incapaz de acceder a su dispositivo por meses? El sistema Guardian permite a una persona de confianza — otro humano verificado — confirmar que el propietario de la wallet sigue vivo. El Guardian tiene estrictamente cero acceso financiero: solo puede llamar una función que reinicia el temporizador de inactividad.</p>',
  'guard-box':'1 Guardian por humano · debe ser un humano verificado en Aequitas<br>Guardian SOLO puede llamar confirmAlive() — cero derechos de transacción<br>Guardian NO PUEDE mover fondos, transferir AEQ ni acceder a la wallet<br>Máximo 3 tutelados por Guardian (evita centralización de confianza)<br>Bloqueo de 7 días en asignación de Guardian (evita asignación forzada)<br>No se permiten relaciones Guardian circulares',
  'dem-title':'3. DEMURRAGE — Mecanismo Anti-Acaparamiento',
  'dem-box':'Tasa: 0,5% por mes después de 3 meses de inactividad (continuo, no escalonado)<br>El reloj se reinicia automáticamente con cualquier transferencia, swap o acción de liquidez<br>AEQ decaído redistribuido a los cuatro pools — nunca destruido<br>Aviso de 14 días mostrado una vez · aviso de 7 días repetido en cada sesión activa',
  'dem-text':'<p>El demurrage es un costo de tenencia sobre el dinero — una tasa de interés negativa que hace costoso acumular y atractivo circular. El experimento de Wörgl (Austria, 1932) usó una moneda con demurrage y redujo el desempleo local un 25% en un año. El Banco Central de Austria lo cerró precisamente porque funcionó demasiado bien. El Chiemgauer (Alemania, 2003) opera según el mismo principio con éxito desde hace más de 20 años.</p>',
  'cap-title':'4. LÍMITE DE RIQUEZA — Aplicación de Justicia Matemática','cap-box':'Límite bootstrap: max(5,min(N,25))× saldo promedio actual<br>1–4 humanos: 5× · +1× por humano · 25+: 25× permanente<br>Se aplica a TODAS las direcciones excepto las 4 pools del protocolo<br>Exceso AEQ redistribuido instantáneamente · Sin intervención manual',
  'ubi-title':'5. RENTA BÁSICA UNIVERSAL — Redistribución Diaria','ubi-box':'Fuentes de ingresos del Pool UBI:<br>· 20% de todas las comisiones de swap del pool AMM AEQ↔tUSD<br>· Desbordamiento de la aplicación del límite de riqueza<br>· Cargos de demurrage de cuentas inactivas<br>· Custodia inactiva liberada después de 4 años<br><br>Distribución: Cada 24 horas, todo el saldo del pool UBI se divide igualmente entre todos los humanos verificados registrados. El pool se reinicia a cero y comienza a llenarse inmediatamente de la actividad continua del protocolo.',
  'inf-title':'6. SIN INFLACIÓN ALGORÍTMICA — Fórmula de Suministro Fijo','inf-box':'El ÚNICO evento que crea nuevo AEQ: un nuevo humano verificado se registra.<br><br>Suministro Total = Humanos Verificados × 1.000 AEQ<br><br>Esto no es una política — es aplicado por el protocolo. Ningún administrador puede acuñar AEQ adicional, ningún voto de gobernanza puede cambiar la emisión. AEQ es la única criptomoneda donde el suministro total está determinado únicamente por el número de humanos vivos verificados.'
},
ru:{
  'logo-sub':'ДОКАЗАТЕЛЬСТВО ЧЕЛОВЕЧНОСТИ','live':'ОНЛАЙН',
  'tab-register':'🔐 Регистрация','tab-explorer':'🔍 Проводник','tab-humans':'👥 Люди','tab-index':'📊 Индекс','tab-network':'🌐 Сеть','tab-protocol':'📜 Протокол V7','tab-swap':'🔄 Обмен',
  'reg-title':'🔐 Зарегистрируйтесь как Верифицированный Человек',
  'reg-sub':'Присоединитесь к сети Aequitas и получите 1 000 AEQ в качестве Универсального Базового Дохода. Однократно, постоянно и полностью бесплатно. Никакие личные данные никогда не сохраняются.',
  'app-title':'РЕГИСТРАЦИЯ ТОЛЬКО ЧЕРЕЗ ANDROID-ПРИЛОЖЕНИЕ',
  'app-text':'Доказательство Человечности требует биометрической верификации на вашем устройстве. Ваш отпечаток пальца или распознавание лица обрабатывается исключительно аппаратным защищённым элементом — биометрические данные никогда не покидают ваше устройство. Приложение создаёт Доказательство с Нулевым Разглашением, которое математически подтверждает вашу уникальность без раскрытия личной информации. Скачайте AequitasBio, отсканируйте биометрию, подключите MetaMask, и ваши <strong style="color:var(--gold)">1 000 AEQ будут зачислены автоматически</strong>.',
  's1t':'Биометрическое Сканирование','s1d':'Открыть AequitasBio · сканировать отпечаток или лицо · аппаратный элемент обрабатывает локально · биометрия никогда не покидает устройство',
  's2t':'Создание ZK-Доказательства','s2d':'Доказательство Groth16 создаётся на сервере · уникальность верифицируется криптографически · личность никогда не раскрывается',
  's3t':'Подключение Кошелька','s3d':'Приложение открывает MetaMask на этой странице · подключите кошелёк Ethereum · доказательство криптографически привязано к вашему адресу',
  's4t':'1 000 AEQ Зачислены','s4d':'Регистрация подтверждена на BlockDAG Aequitas за 6 секунд · 1 000 AEQ зачислены мгновенно · личность навсегда записана как верифицированный человек',
  'priv-bar':'🔒 Аппаратный Защищённый Элемент · Доказательство Groth16 с Нулевым Разглашением · Биометрия никогда не покидает устройство · Без комиссий · Одна регистрация на человека · Постоянно и неизменно',
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
  'h-what':'Что такое Верифицированный Человек?','h-what-t':'Верифицированный Человек — адрес кошелька, криптографически доказанный принадлежащим уникальному живому человеку через биометрическое Доказательство с Нулевым Разглашением.',
  'h-zkp':'Система ZK-Доказательств','h-zkp-t':'Aequitas использует Groth16 на BN128. Размер доказательства: ~200 байт. Время верификации: ~10мс.',
  'h-sybil':'Защита от Атак Сибиллы','h-sybil-t':'Каждый биометрический хэш хранится постоянно с keccak256. Двойная регистрация немедленно отклоняется. ⚠ Тестовая фаза: верификация привязана к устройству. Планируется сенсор MAX30102 PPG.',
  'h-global':'Глобальная Финансовая Инклюзия','h-global-t':'Банковский счёт, кредитная карта или криптовалюта не требуются. Только Android-смартфон с биометрическим датчиком.',
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
  'bootstrap-title':'Подключить Новую Ноду','bootstrap-desc':'Установите PEER_NODES на адрес bootstrap-ноды ниже. Нода автоматически синхронизируется и начнёт производство блоков.',
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
  'phases-desc':'Границы фаз определяют вехи роста сети. Мультипликатор лимита богатства в настоящее время зафиксирован на 25× (константа кода Go: wealthCapMultiplier = 25.0) — автоматическая корректировка по фазам запланирована как будущее обновление протокола.',
  'p0':'Bootstrap · &lt;100 людей · Лимит богатства: 25× средний баланс · Активен сейчас',
  'p1':'Рост · 100–10 000 людей · Лимит богатства: 25× (планируемое снижение: 20×)',
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
  'usp-c3-title':'Только смартфон','usp-c3-desc':'Без компьютера, без банковского счёта, без документа. Android-телефон со сканером отпечатка — всё что нужно.',
  'usp-c4-title':'Ежедневный UBI навсегда','usp-c4-desc':'После регистрации вы автоматически получаете ежедневную долю выплат UBI — каждый день, без каких-либо действий.',
  'v7-intro-title':'Что такое AequitasV7?',
  'v7-intro-text':'AequitasV7 — центральный смарт-контракт протокола Aequitas. "V7" — 7-я основная версия контракта справедливости. Развёрнут неизменяемым образом в Aequitas Chain (ID 1926) и управляет всем: регистрация людей, верификация ZK, управление балансами, лимит богатства, распределение UBI, комиссии свопов. Ни один администратор не может обновить его. Шесть механизмов образуют самоусиливающуюся систему.',
},
zh:{
  'logo-sub':'人类证明','live':'实时',
  'tab-register':'🔐 注册','tab-explorer':'🔍 浏览器','tab-humans':'👥 人类','tab-index':'📊 指数','tab-network':'🌐 网络','tab-protocol':'📜 协议 V7','tab-swap':'🔄 兑换',
  'reg-title':'🔐 注册成为经过验证的人类',
  'reg-sub':'加入Aequitas网络并获得1,000 AEQ的普遍基本收入补贴。一次性、永久性且完全免费。永远不会存储任何个人数据。',
  'app-title':'仅通过安卓应用注册',
  'app-text':'人类证明需要在您的设备上进行生物特征验证。指纹或面部识别由手机内的硬件安全元件独立处理——原始生物特征数据永远不会离开您的设备。应用程序生成零知识证明，在不透露个人信息的情况下数学证明您的唯一性。下载AequitasBio，扫描您的生物特征，连接MetaMask，您的<strong style="color:var(--gold)">1,000 AEQ将自动记入</strong>。',
  's1t':'生物特征扫描','s1d':'打开AequitasBio · 扫描指纹或面部 · 硬件安全元件本地处理 · 生物特征数据永不离开设备',
  's2t':'ZK证明生成','s2d':'在证明服务器上生成Groth16零知识证明 · 唯一性得到密码学验证 · 身份永不泄露',
  's3t':'连接钱包','s3d':'应用在此页面打开MetaMask · 连接您的以太坊钱包 · 证明与您的地址密码绑定',
  's4t':'获得1,000 AEQ','s4d':'在6秒内在Aequitas BlockDAG上确认注册 · 立即记入1,000 AEQ · 身份永久记录为经过验证的人类',
  'priv-bar':'🔒 硬件安全元件 · Groth16零知识证明 · 生物特征数据永不离开设备 · 无Gas费 · 每人一次注册 · 永久不可更改',
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
  'h-what':'什么是已验证人类？','h-what-t':'已验证人类是通过生物特征零知识证明，密码学证明属于独特活人的钱包地址。生物特征数据从不传输或存储。',
  'h-zkp':'零知识证明系统','h-zkp-t':'Aequitas在BN128上使用Groth16。证明大小：约200字节。验证时间：约10毫秒。',
  'h-sybil':'女巫攻击防护','h-sybil-t':'每个生物特征哈希使用keccak256永久存储。尝试两次注册立即被拒绝。⚠ 测试阶段：当前验证与设备绑定。计划使用MAX30102 PPG传感器实现独立于设备的识别。',
  'h-global':'全球金融包容','h-global-t':'无需银行账户、信用卡或加密货币。只需一台带生物特征传感器的安卓手机。',
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
  'usp-c3-title':'只需一部智能手机','usp-c3-desc':'无需电脑，无需银行账户，无需身份证件。一部带指纹传感器的Android手机就足够了。',
  'usp-c4-title':'永久每日UBI','usp-c4-desc':'注册后，您每天自动获得UBI支付份额——每天，无需任何操作。',
  'v7-intro-title':'什么是 AequitasV7？',
  'v7-intro-text':'AequitasV7是Aequitas协议的核心智能合约。"V7"指公平合约的第7个主要版本。它不可更改地部署在Aequitas Chain（链ID 1926）上，处理所有方面：人类注册、ZK证明验证、余额管理、财富上限、UBI分配、兑换手续费。没有管理员可以升级它。六个机制形成自我强化系统。',
},
id:{
  'logo-sub':'BUKTI KEMANUSIAAN','live':'LANGSUNG',
  'tab-register':'🔐 Daftar','tab-explorer':'🔍 Penjelajah','tab-humans':'👥 Manusia','tab-index':'📊 Indeks','tab-network':'🌐 Jaringan','tab-protocol':'📜 Protokol V7','tab-swap':'🔄 Tukar',
  'reg-title':'🔐 Daftar sebagai Manusia Terverifikasi',
  'reg-sub':'Bergabunglah dengan jaringan Aequitas dan terima hibah Pendapatan Dasar Universal sebesar 1.000 AEQ. Satu kali, permanen, dan sepenuhnya gratis. Tidak ada data pribadi yang pernah disimpan.',
  'app-title':'PENDAFTARAN HANYA MELALUI APLIKASI ANDROID',
  'app-text':'Bukti Kemanusiaan memerlukan verifikasi biometrik pada perangkat Anda. Sidik jari atau pengenalan wajah diproses secara eksklusif oleh Elemen Aman Perangkat Keras — data biometrik mentah tidak pernah meninggalkan perangkat Anda. Aplikasi menghasilkan Bukti Pengetahuan Nol yang membuktikan keunikan Anda secara matematis tanpa mengungkapkan informasi pribadi. Unduh AequitasBio, pindai biometrik Anda, hubungkan MetaMask, dan <strong style="color:var(--gold)">1.000 AEQ Anda akan dikreditkan secara otomatis</strong>.',
  's1t':'Pemindaian Biometrik','s1d':'Buka AequitasBio · pindai sidik jari atau wajah · Elemen Aman Perangkat Keras memproses secara lokal · data biometrik tidak pernah meninggalkan perangkat',
  's2t':'Pembuatan Bukti ZK','s2d':'Bukti Groth16 Pengetahuan Nol dibuat di server · keunikan diverifikasi secara kriptografis · identitas tidak pernah terungkap',
  's3t':'Hubungkan Dompet','s3d':'Aplikasi membuka MetaMask di halaman ini · hubungkan dompet Ethereum Anda · bukti terikat secara kriptografis ke alamat Anda',
  's4t':'1.000 AEQ Dikreditkan','s4d':'Pendaftaran dikonfirmasi di BlockDAG Aequitas dalam 6 detik · 1.000 AEQ dikreditkan seketika · identitas Anda dicatat permanen sebagai manusia terverifikasi',
  'priv-bar':'🔒 Elemen Aman Perangkat Keras · Bukti ZK Groth16 · Data biometrik tidak pernah meninggalkan perangkat · Tanpa biaya gas · Satu pendaftaran per manusia · Permanen &amp; tidak dapat diubah',
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
  'h-what':'Apa itu Manusia Terverifikasi?','h-what-t':'Manusia Terverifikasi adalah alamat dompet yang terbukti secara kriptografis milik manusia hidup yang unik melalui Bukti Pengetahuan Nol biometrik. Data biometrik tidak pernah ditransmisikan atau disimpan.',
  'h-zkp':'Sistem Bukti ZK','h-zkp-t':'Aequitas menggunakan Groth16 pada BN128. Ukuran bukti: ~200 byte. Waktu verifikasi: ~10ms.',
  'h-sybil':'Pencegahan Serangan Sybil','h-sybil-t':'Setiap hash biometrik disimpan permanen dengan keccak256. Mencoba mendaftar dua kali langsung ditolak. ⚠ Fase uji coba: verifikasi terikat perangkat. Sensor fisiologis MAX30102 PPG direncanakan.',
  'h-global':'Inklusi Keuangan Global','h-global-t':'Tidak perlu rekening bank, kartu kredit, atau cryptocurrency sebelumnya. Hanya smartphone Android dengan sensor biometrik.',
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
  'usp-c3-title':'Hanya butuh smartphone','usp-c3-desc':'Tanpa komputer, tanpa rekening bank, tanpa dokumen ID. Ponsel Android dengan sensor sidik jari sudah cukup.',
  'usp-c4-title':'UBI harian selamanya','usp-c4-desc':'Setelah terdaftar, Anda secara otomatis menerima bagian harian dari pembayaran UBI — setiap hari, tanpa tindakan apa pun.',
  'v7-intro-title':'Apa itu AequitasV7?',
  'v7-intro-text':'AequitasV7 adalah kontrak pintar inti dari protokol Aequitas. "V7" mengacu pada versi utama ke-7 dari kontrak keadilan. Dikerahkan secara tidak dapat diubah di Aequitas Chain (ID 1926) dan menangani setiap aspek: pendaftaran manusia, verifikasi ZK, manajemen saldo, batas kekayaan, distribusi UBI, biaya swap. Tidak ada admin yang dapat memperbaruinya. Keenam mekanisme membentuk sistem yang saling memperkuat.',
},
it:{
  'logo-sub':'PROVA DI UMANITÀ','live':'LIVE',
  'tab-register':'🔐 Registrati','tab-explorer':'🔍 Explorer','tab-humans':'👥 Umani','tab-index':'📊 Indice','tab-network':'🌐 Rete','tab-protocol':'📜 Protocollo V7','tab-swap':'🔄 Scambia',
  'reg-title':'🔐 Registrati come Umano Verificato',
  'reg-sub':'Unisciti alla rete Aequitas e ricevi il tuo sussidio di Reddito Universale di Base di 1.000 AEQ. Una tantum, permanente e completamente gratuito. Nessun dato personale viene mai memorizzato.',
  'app-title':'REGISTRAZIONE SOLO VIA APP ANDROID',
  'app-text':'La Prova di Umanità richiede la verifica biometrica sul tuo dispositivo personale. La tua impronta digitale o il riconoscimento facciale viene elaborato esclusivamente dall\'Elemento Sicuro Hardware del tuo telefono — i dati biometrici grezzi non lasciano mai il tuo dispositivo. L\'app genera una Prova a Conoscenza Zero che dimostra matematicamente la tua unicità senza rivelare informazioni personali. Scarica AequitasBio, scansiona la tua biometria, connetti MetaMask, e i tuoi <strong style="color:var(--gold)">1.000 AEQ saranno accreditati automaticamente</strong>.',
  's1t':'Scansione Biometrica','s1d':'Apri AequitasBio · scansiona impronta o volto · l\'Elemento Sicuro Hardware elabora localmente · i dati biometrici non lasciano mai il dispositivo',
  's2t':'Generazione Prova ZK','s2d':'La Prova Groth16 a Conoscenza Zero viene generata sul server · l\'unicità viene verificata crittograficamente · la tua identità non viene mai rivelata',
  's3t':'Connetti Wallet','s3d':'L\'app apre MetaMask su questa pagina · connetti il tuo wallet Ethereum · la prova è crittograficamente legata al tuo indirizzo',
  's4t':'1.000 AEQ Accreditati','s4d':'Registrazione confermata su Aequitas BlockDAG entro 6 secondi · 1.000 AEQ accreditati istantaneamente · la tua identità è registrata permanentemente come umano verificato',
  'priv-bar':'🔒 Elemento Sicuro Hardware · Prova Groth16 a Conoscenza Zero · Dati biometrici non lasciano il dispositivo · Nessuna commissione gas · Una registrazione per umano · Permanente e immutabile',
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
  'h-what':'Cos\'è un Umano Verificato?','h-what-t':'Un Umano Verificato è un indirizzo wallet dimostrato crittograficamente appartenere a un essere umano unico e vivente tramite Prova a Conoscenza Zero biometrica. I dati biometrici non vengono mai trasmessi o memorizzati.',
  'h-zkp':'Sistema di Prova a Conoscenza Zero','h-zkp-t':'Aequitas usa il sistema di prova Groth16 sulla curva ellittica BN128. Dimensione prova: ~200 byte. Tempo di verifica: ~10ms. La prova dimostra matematicamente l\'unicità senza rivelare alcuna informazione identificativa.',
  'h-sybil':'Prevenzione Attacchi Sybil','h-sybil-t':'Ogni hash biometrico viene memorizzato permanentemente con keccak256. Tentare di registrarsi due volte viene immediatamente rifiutato. Un umano, un wallet, per sempre. ⚠ Fase di test: la verifica è attualmente legata al dispositivo. È previsto un sensore fisiologico (MAX30102 PPG) per la verifica indipendente dal dispositivo.',
  'h-global':'Inclusione Finanziaria Globale','h-global-t':'Nessun conto bancario, nessuna carta di credito, nessuna criptovaluta precedente necessaria. Solo uno smartphone Android con sensore biometrico. Aequitas è progettato per essere accessibile a ogni essere umano sulla Terra.',
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
  'bootstrap-title':'Connettere un Nuovo Node','bootstrap-desc':'Per eseguire il tuo node, imposta la variabile d\'ambiente PEER_NODES sull\'indirizzo bootstrap qui sotto. Il tuo node si sincronizzerà automaticamente con lo stato completo della chain.',
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
  'usp-c3-title':'Solo uno smartphone','usp-c3-desc':'Senza computer, senza conto bancario, senza documento d\'identità. Un telefono Android con sensore di impronte digitali è tutto ciò che serve.',
  'usp-c4-title':'UBI quotidiano per sempre','usp-c4-desc':'Una volta registrato, ricevi automaticamente una quota giornaliera dei pagamenti UBI — ogni giorno, senza alcuna azione richiesta.',
  'v7-intro-title':'Cos\'è AequitasV7?',
  'v7-intro-text':'AequitasV7 è il contratto intelligente centrale del protocollo Aequitas. "V7" si riferisce alla 7ª versione principale del contratto di equità. È distribuito immutabilmente su Aequitas Chain (ID 1926) e gestisce ogni aspetto: registrazione umana, verifica ZK, gestione saldi, limite di ricchezza, distribuzione UBI, commissioni swap. Nessun amministratore può aggiornarlo. I sei meccanismi formano un sistema auto-rinforzante.'
}
};

function showTab(name, el) {
  document.querySelectorAll('.tab-content').forEach(t => t.classList.remove('active'));
  document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
  document.getElementById('tab-' + name).classList.add('active');
  el.classList.add('active');
  if (name === 'swap') loadPoolStatus();
  if (name === 'index') { drawGiniHistoryChart(); drawLorenzCurve(); drawWcapSlideChart(); }
  history.pushState(null, '', '/' + name);
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
    document.getElementById('idx-gini').textContent = typeof d.gini === 'number' ? d.gini.toFixed(3) : '—';
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

    // UBI countdown timer + fill bar (shows time elapsed since last payout)
    if (d.ubi_next_payout_secs !== undefined) {
      startUBITimer(d.ubi_next_payout_secs);
      const fillPct = Math.min(100, Math.max(0, (86400 - d.ubi_next_payout_secs) / 86400 * 100));
      const fillBar = document.getElementById('ubi-fill-bar');
      if (fillBar) fillBar.style.width = fillPct.toFixed(1) + '%';
    }

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
    const pad = {l:36,r:16,t:14,b:8};
    ctx.strokeStyle = 'rgba(255,255,255,0.05)'; ctx.lineWidth = 1;
    for (let i = 0; i <= 4; i++) {
      const y = pad.t + (H-pad.t-pad.b)*(1-i/4);
      ctx.beginPath(); ctx.moveTo(pad.l,y); ctx.lineTo(W-pad.r,y); ctx.stroke();
      ctx.fillStyle = 'rgba(255,255,255,0.3)'; ctx.font = '9px monospace'; ctx.textAlign = 'right';
      ctx.fillText((i*25)+'', pad.l-3, y+3);
    }
    const targetY = pad.t + (H-pad.t-pad.b)*(1-35/100);
    ctx.strokeStyle = 'rgba(0,255,209,0.3)'; ctx.setLineDash([4,4]);
    ctx.beginPath(); ctx.moveTo(pad.l,targetY); ctx.lineTo(W-pad.r,targetY); ctx.stroke();
    ctx.setLineDash([]); ctx.fillStyle = 'rgba(0,255,209,0.6)'; ctx.font = '8px monospace'; ctx.textAlign = 'left';
    ctx.fillText('Target 35', W-pad.r-48, targetY-3);
    const xs = (W-pad.l-pad.r)/Math.max(history.length-1,1);
    ctx.strokeStyle = '#C9A84C'; ctx.lineWidth = 2;
    ctx.beginPath();
    history.forEach((pt,i)=>{ const x=pad.l+i*xs, y=pad.t+(H-pad.t-pad.b)*(1-pt.idx/100); if(i===0)ctx.moveTo(x,y); else ctx.lineTo(x,y); });
    ctx.stroke();
    ctx.fillStyle = '#C9A84C';
    history.forEach((pt,i)=>{ const x=pad.l+i*xs, y=pad.t+(H-pad.t-pad.b)*(1-pt.idx/100); ctx.beginPath(); ctx.arc(x,y,3,0,2*Math.PI); ctx.fill(); });
  } catch(e) {}
}

async function drawLorenzCurve() {
  const canvas = document.getElementById('lorenz-chart');
  if (!canvas || !canvas.offsetParent) return;
  canvas.width = canvas.offsetWidth;
  const ctx = canvas.getContext('2d');
  const W = canvas.width, H = canvas.height;
  ctx.clearRect(0,0,W,H);
  const pad = {l:36,r:16,t:14,b:28};
  try {
    const d = await (await fetch('/api/humans')).json();
    const humans = d.humans || [];
    if (humans.length < 2) {
      ctx.fillStyle = 'rgba(255,255,255,0.3)'; ctx.font = '11px monospace'; ctx.textAlign = 'center';
      ctx.fillText('Not enough humans registered yet', W/2, H/2); return;
    }
    const bals = humans.map(h=>h.balance||0).sort((a,b)=>a-b);
    const total = bals.reduce((s,b)=>s+b,0);
    const n = bals.length;
    ctx.strokeStyle = 'rgba(255,255,255,0.12)'; ctx.setLineDash([4,4]); ctx.lineWidth = 1;
    ctx.beginPath(); ctx.moveTo(pad.l,H-pad.b); ctx.lineTo(W-pad.r,pad.t); ctx.stroke();
    ctx.setLineDash([]);
    for (let i=0;i<=4;i++) {
      const x=pad.l+(W-pad.l-pad.r)*i/4, y=pad.t+(H-pad.t-pad.b)*(1-i/4);
      ctx.strokeStyle='rgba(255,255,255,0.05)'; ctx.beginPath(); ctx.moveTo(x,pad.t); ctx.lineTo(x,H-pad.b); ctx.stroke();
      ctx.fillStyle='rgba(255,255,255,0.25)'; ctx.font='9px monospace'; ctx.textAlign='center'; ctx.fillText((i*25)+'%',x,H-4);
    }
    ctx.strokeStyle='#C9A84C'; ctx.lineWidth=2;
    ctx.beginPath(); ctx.moveTo(pad.l,H-pad.b);
    let cum=0;
    bals.forEach((b,i)=>{ cum+=b; const x=pad.l+(W-pad.l-pad.r)*(i+1)/n, y=(H-pad.b)-(H-pad.t-pad.b)*(cum/total); ctx.lineTo(x,y); });
    ctx.stroke();
    ctx.fillStyle='rgba(200,168,76,0.12)'; ctx.beginPath(); ctx.moveTo(pad.l,H-pad.b);
    cum=0; bals.forEach((b,i)=>{ cum+=b; const x=pad.l+(W-pad.l-pad.r)*(i+1)/n, y=(H-pad.b)-(H-pad.t-pad.b)*(cum/total); ctx.lineTo(x,y); });
    ctx.lineTo(W-pad.r,H-pad.b); ctx.closePath(); ctx.fill();
    ctx.fillStyle='rgba(255,255,255,0.3)'; ctx.font='9px monospace'; ctx.textAlign='center';
    ctx.fillText('% of AEQ (cumulative)', pad.l-30, H/2);
  } catch(e) {}
}

function drawWcapSlideChart() {
  const canvas = document.getElementById('wcap-slide-chart');
  if (!canvas || !canvas.offsetParent) return;
  canvas.width = canvas.offsetWidth;
  const ctx = canvas.getContext('2d');
  const W = canvas.width, H = canvas.height;
  ctx.clearRect(0,0,W,H);
  const pad = {l:36,r:16,t:8,b:24};
  const maxN = 28;
  const bw = (W-pad.l-pad.r)/maxN;
  for (let n=1; n<=maxN; n++) {
    const mult = Math.max(5,Math.min(n,25));
    const bh = (H-pad.t-pad.b)*(mult/25);
    const x = pad.l+(n-1)*bw, y = H-pad.b-bh;
    ctx.fillStyle = n<=25 ? (n===25?'rgba(0,255,209,0.55)':'rgba(200,168,76,0.55)') : 'rgba(255,255,255,0.1)';
    ctx.fillRect(x+1,y,bw-2,bh);
    if (n===1||n===5||n===10||n===15||n===20||n===25) {
      ctx.fillStyle='rgba(255,255,255,0.5)'; ctx.font='8px monospace'; ctx.textAlign='center';
      ctx.fillText(mult+'×',x+bw/2,y-1);
      ctx.fillStyle='rgba(255,255,255,0.3)';
      ctx.fillText('N='+n,x+bw/2,H-2);
    }
  }
  ctx.strokeStyle='rgba(0,255,209,0.3)'; ctx.lineWidth=1; ctx.setLineDash([3,3]);
  const lockY = H-pad.b-(H-pad.t-pad.b);
  ctx.beginPath(); ctx.moveTo(pad.l+(25-1)*bw,lockY); ctx.lineTo(W-pad.r,lockY); ctx.stroke();
  ctx.setLineDash([]);
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
      if (outEl) outEl.textContent = est.amountOut.toFixed(6) + ' ' + outUnit;
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
      swapLog('✗ Remove liquidity failed: ' + sanitize(data.message), 'err');
    }
  } catch (e) {
    swapLog('✗ Error: ' + sanitize(e.message), 'err');
  }
  document.getElementById('swap-btn-removeliq').disabled = false;
}

function activateTabFromPath(path) {
  const tabNames = ['register','explorer','humans','index','network','protocol','swap'];
  const name = (path || '').replace(/^\//, '').split('/')[0];
  if (!name || !tabNames.includes(name)) return;
  const tabEl = document.querySelector('.tab[onclick*="showTab(\'' + name + '\'"]');
  if (!tabEl) return;
  document.querySelectorAll('.tab-content').forEach(t => t.classList.remove('active'));
  document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
  document.getElementById('tab-' + name).classList.add('active');
  tabEl.classList.add('active');
  if (name === 'swap') loadPoolStatus();
}

document.addEventListener('DOMContentLoaded', () => {
  const amtInput = document.getElementById('swap-amount');
  if (amtInput) amtInput.addEventListener('input', updateFeeEstimate);

  // Activate correct tab on initial load based on URL path
  activateTabFromPath(window.location.pathname);
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
  el.innerHTML += '<div><span class="' + type + '">' + msg + '</span></div>';
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

    // Compute nullifier = SHA256(bioHash + ":aequitas-ubi-v1")
    // Stored on-chain to prevent the same biometric registering a second wallet.
    let nullifier = '';
    const bioHashForNullifier = pendingBioHash || proofData.bioHashKey || '';
    if (bioHashForNullifier) {
      const enc = new TextEncoder();
      const buf = await crypto.subtle.digest('SHA-256', enc.encode(bioHashForNullifier + ':aequitas-ubi-v1'));
      nullifier = Array.from(new Uint8Array(buf)).map(b => b.toString(16).padStart(2, '0')).join('');
    }

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
        nullifier: nullifier
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

function generateNodeGuidePDF() {
  var lang = curLang || 'en';
  if (window.jspdf) { _buildNodeGuidePDF(lang); return; }
  var s = document.createElement('script');
  s.src = 'https://cdnjs.cloudflare.com/ajax/libs/jspdf/2.5.1/jspdf.umd.min.js';
  s.onload = function() { _buildNodeGuidePDF(lang); };
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
    en:{title:'Aequitas Node Operator Guide',sub:'Complete step-by-step guide · Aequitas Chain (Chain ID 1926)',badge:'BETA v0.1 · Open Source · Permissionless · No stake required',
      s1:'1. Overview',what:'What a node does',wtxt:'An Aequitas node participates fully in the network: produces blocks in the BlockDAG consensus, validates Groth16 zero-knowledge biometric proofs for new human registrations, enforces wealth caps and demurrage at protocol level, syncs state with peers via libp2p + HTTP, and optionally runs daily pool distributions. Every node runs the full chain — there are no light clients.',
      earn:'What you earn',etxt:'Set NODE_OPERATOR_WALLET to a registered human wallet. The Validators Pool accumulates 40% of all protocol fees (swap fees, demurrage, wealth cap overflow). Every 24 h the primary node distributes the pool balance proportionally among all registered node operator wallets. No stake required — block production is fully permissionless.',
      s2:'2. Requirements',rh:['Component','Minimum','Recommended'],rr:[['OS','Linux / Docker-capable host','Ubuntu 22.04 LTS'],['RAM','512 MB','1 GB (EVM needs headroom)'],['CPU','1 vCPU','2 vCPU (Groth16 is CPU-bound)'],['Storage','2 GB','10 GB SSD (chain grows over time)'],['Database','PostgreSQL 14+','Railway or Supabase (managed)'],['Network','Public IP / port forward','TCP 8080 open, stable uptime']],
      s3:'3. Environment Variables',e3:'Set these before starting the node. Variables marked YES are required; "For rewards" is needed to earn validator payouts.',eh:['Variable','Purpose','Required?'],er:[['DATABASE_URL','PostgreSQL connection string: postgres://user:pass@host:5432/aequitas','YES'],['RELAYER_PRIVATE_KEY','Private key (0x...) of the EOA that signs on-chain human registrations','YES'],['NODE_OPERATOR_WALLET','Registered human wallet address that receives daily validator pool rewards','For rewards'],['RELAYER_ADDRESS','EOA address matching RELAYER_PRIVATE_KEY. Has a hardcoded fallback but set explicitly.','Recommended'],['PORT','HTTP port for API + JSON-RPC. Default: 8080','NO'],['PEER_NODES','Bootstrap peer URLs (comma-separated). Use the primary node URL.','Recommended'],['NODE_KEY','32-byte hex for stable libp2p peer identity. Auto-generated if omitted (changes on restart).','NO'],['IS_PRIMARY_NODE','"true" only on the designated primary. Enables daily UBI + Validator + LP distributions.','NO'],['RESET_STATE','"true" wipes the database on startup. DESTRUCTIVE — never use in production.','NO']],
      s4:'4. Quick Start — Railway (Recommended for BETA)',r4:'Railway is the fastest way to get running. The free tier meets minimum requirements for BETA. Estimated setup time: 10–15 minutes.',rs:['Fork the repo: https://github.com/hanoi96international-gif/Aequitas','Create a Railway account at railway.app and start a new project','Click "Deploy from GitHub Repo" and select your fork','In the project: + New → Database → Add PostgreSQL — DATABASE_URL is auto-set by Railway','Go to your service → Settings → Variables and add the env vars from Section 3','Set PEER_NODES=https://aequitas.digital so your node syncs from the primary','Set NODE_OPERATOR_WALLET=<your registered AEQ human wallet> to receive daily validator rewards','Set RELAYER_PRIVATE_KEY=<your EOA private key> for signing on-chain registrations','Click "Deploy" — the Dockerfile in the repo root handles the build (~3 min for Go compilation)','Watch the deploy logs for: "Aequitas Node Running ✓" and "[NODE] Registered node operator wallet"','Open YOUR-RAILWAY-URL/api/status to confirm the node is live and block height is climbing','Add your node\'s RPC to MetaMask: Chain ID 1926, Symbol AEQ, URL https://YOUR-URL/rpc'],rn:'Railway assigns a random subdomain; custom domains can be set in project settings. Only port 8080 needs to be exposed — P2P is managed internally by the node.',
      s5:'5. Quick Start — Docker',d5:'For VPS, cloud VM, or local server. Prerequisites: Docker installed, PostgreSQL available (managed service or a second container).',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\n\n# Build the image (~3 min for Go compilation)\ndocker build -t aequitas-node .\n\n# Run the node\ndocker run -d --name aequitas-node \\\n  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n  -e RELAYER_PRIVATE_KEY="0xYOUR_PRIVATE_KEY" \\\n  -e RELAYER_ADDRESS="0xYOUR_ADDRESS" \\\n  -e NODE_OPERATOR_WALLET="0xYOUR_HUMAN_WALLET" \\\n  -e PEER_NODES="https://aequitas.digital" \\\n  -p 8080:8080 aequitas-node\n\n# Follow logs\ndocker logs -f aequitas-node',dn:'The container exposes port 8080 (API + RPC). Ensure TCP 8080 is open inbound in your firewall or cloud security group.',
      s6:'6. Verify Your Node',v6:'Once running, check these endpoints to confirm the node is synced and healthy.',vc:'# 1. Node status (height should match the primary node within 1-2 blocks)\ncurl https://YOUR-NODE-URL/api/status\n# Expect: { "height": N, "total_humans": N, "index": N }\n\n# 2. EVM JSON-RPC (EVM compatibility check)\ncurl -X POST https://YOUR-NODE-URL/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'\n\n# 3. In startup logs, confirm:\n#    [NODE] Registered node operator wallet: 0xYOUR_WALLET\n#    Aequitas Node Running V  (blocks every ~6 seconds)\n\n# MetaMask: RPC URL = https://YOUR-NODE-URL/rpc | Chain ID: 1926 | Symbol: AEQ',
      s7:'7. P2P Networking & Sync',p7:'Set PEER_NODES to at least one known bootstrap URL. The node connects and syncs the full chain automatically using libp2p gossip (real-time) plus periodic HTTP pulls from peers (fallback). The primary node libp2p multiaddress is:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'The HTTP URL in PEER_NODES is more stable for bootstrap. The libp2p multiaddress above may change if the primary node is redeployed on Railway. When in doubt, use the HTTPS URL.',
      s8:'8. Earning Validator Rewards',w8:'Validator rewards come from the Validators Pool (40% of all protocol fees: swap fees, demurrage, wealth cap overflow). Steps to receive rewards:',b8:['First register as a human on Aequitas: install the Android app and complete biometric registration to receive your wallet address and 1,000 AEQ','Set NODE_OPERATOR_WALLET to that registered wallet address in your node\'s environment variables','Start (or restart) your node — it calls RegisterNode() on startup. Confirm in logs: "[NODE] Registered node operator wallet: 0xYOUR_WALLET"','The primary node distributes rewards every 24 h to all registered node operator wallets proportional to blocks produced','Secondary nodes do NOT need to trigger distribution — just keep your node running and synced','No minimum uptime required for BETA, but consistently offline nodes contribute less to block production and proportionally less to the pool share'],
      s9:'9. Troubleshooting',th:['Symptom','Likely Cause','Fix'],tr:[['Height stays at 0 after start','PEER_NODES not set or unreachable','Set PEER_NODES=https://aequitas.digital and redeploy'],['"no code at address" in logs','V7 contract not deployed in EVM','Verify RELAYER_ADDRESS is set; node auto-deploys V7 on startup if missing'],['"NODE_OPERATOR_WALLET not set" in logs','Missing env var','Set NODE_OPERATOR_WALLET to your registered human wallet address'],['DATABASE_URL error on startup','Wrong connection string or DB unreachable','Check format: postgres://user:pass@host:5432/dbname and that PostgreSQL is running'],['Port 8080 not reachable','Firewall or cloud provider config','Open TCP 8080 inbound; check Railway/Render/VPS port settings'],['Docker build fails with go module error','No internet access during build','Ensure Docker build environment has outbound internet; Railway handles this automatically']],
      s10:'10. MetaMask Configuration',m10:'To use your own node as the RPC endpoint in MetaMask or any EVM-compatible wallet:',mh:['Field','Value'],mr:[['Network Name','Aequitas Chain'],['RPC URL','https://YOUR-NODE-URL/rpc'],['Chain ID','1926  (hex: 0x786)'],['Currency Symbol','AEQ'],['Decimals','18'],['Block Explorer','https://aequitas.digital']],
      foot:'Open source · Permissionless · No admin keys · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'},
    de:{title:'Aequitas Node-Betreiber-Handbuch',sub:'Vollständige Schritt-für-Schritt-Anleitung · Aequitas Chain (Chain ID 1926)',badge:'BETA v0.1 · Open Source · Erlaubnisfrei · Kein Stake erforderlich',
      s1:'1. Überblick',what:'Was ein Node leistet',wtxt:'Ein Aequitas-Node nimmt vollständig am Netzwerk teil: produziert Blöcke im BlockDAG-Konsens, validiert Groth16-Zero-Knowledge-Biometriebeweise für neue Menschenregistrierungen, setzt Vermögensobergrenzen und Demurrage auf Protokollebene durch, synchronisiert den Zustand mit Peers via libp2p + HTTP und führt optional tägliche Pool-Verteilungen durch. Jeder Node führt die vollständige Chain aus — es gibt keine Light-Clients.',
      earn:'Was du verdienst',etxt:'NODE_OPERATOR_WALLET auf eine als Mensch registrierte Wallet-Adresse setzen. Der Validators-Pool erhält 40% aller Protokollgebühren (Swap-Gebühren, Demurrage, Vermögensobergrenze-Überschuss). Alle 24 Stunden verteilt der primäre Node den Pool-Saldo proportional auf alle registrierten Node-Betreiber-Wallets. Kein Stake erforderlich — Blockproduktion ist vollständig erlaubnisfrei.',
      s2:'2. Voraussetzungen',rh:['Komponente','Minimum','Empfohlen'],rr:[['Betriebssystem','Linux / Docker-fähiger Host','Ubuntu 22.04 LTS'],['RAM','512 MB','1 GB (EVM braucht Spielraum)'],['CPU','1 vCPU','2 vCPU (Groth16 ist CPU-gebunden)'],['Speicher','2 GB','10 GB SSD (Chain wächst kontinuierlich)'],['Datenbank','PostgreSQL 14+','Railway oder Supabase (verwaltet)'],['Netzwerk','Öffentliche IP / Port-Weiterleitung','TCP 8080 offen, stabile Verfügbarkeit']],
      s3:'3. Umgebungsvariablen',e3:'Diese vor dem Start des Nodes setzen. Mit JA markierte Variablen sind Pflicht; "Für Bel." wird benötigt um Validator-Auszahlungen zu erhalten.',eh:['Variable','Zweck','Pflicht?'],er:[['DATABASE_URL','PostgreSQL-Verbindungsstring: postgres://user:pass@host:5432/aequitas','JA'],['RELAYER_PRIVATE_KEY','Privater Schlüssel (0x...) des EOA der On-Chain-Menschenregistrierungen signiert','JA'],['NODE_OPERATOR_WALLET','Registrierte Mensch-Wallet-Adresse die täglich Validator-Pool-Bel. erhält','Für Bel.'],['RELAYER_ADDRESS','EOA-Adresse passend zu RELAYER_PRIVATE_KEY. Hat Fallback, aber explizit setzen.','Empfohlen'],['PORT','HTTP-Port für API + JSON-RPC. Standard: 8080','NEIN'],['PEER_NODES','Bootstrap-Peer-URLs (kommagetrennt). Primär-Node-URL verwenden.','Empfohlen'],['NODE_KEY','32-Byte-Hex für stabile libp2p-Identität. Auto-generiert wenn nicht gesetzt.','NEIN'],['IS_PRIMARY_NODE','"true" nur auf dem designierten Primär-Node. Aktiviert tägliche Verteilungen.','NEIN'],['RESET_STATE','"true" löscht die DB beim Start. DESTRUKTIV — niemals in Produktion.','NEIN']],
      s4:'4. Schnellstart — Railway (Empfohlen für BETA)',r4:'Railway ist der schnellste Einstieg. Der kostenlose Tarif erfüllt die Mindestanforderungen während der BETA. Geschätzte Einrichtungszeit: 10–15 Minuten.',rs:['Repo forken: https://github.com/hanoi96international-gif/Aequitas','Railway-Konto auf railway.app erstellen und neues Projekt starten','"Deploy from GitHub Repo" anklicken und den Fork auswählen','Im Projekt: + New → Database → Add PostgreSQL — DATABASE_URL wird automatisch gesetzt','Service → Settings → Variables aufrufen und Umgebungsvariablen aus Abschnitt 3 hinzufügen','PEER_NODES=https://aequitas.digital setzen','NODE_OPERATOR_WALLET=<deine registrierte AEQ-Mensch-Wallet> für tägliche Validator-Bel. setzen','RELAYER_PRIVATE_KEY=<EOA-Privatschlüssel für On-Chain-Registrierungssignaturen> setzen','"Deploy" klicken — das Dockerfile im Root-Verzeichnis steuert den Build (~3 Min. für Go-Kompilierung)','Deploy-Logs auf "Aequitas Node Running ✓" und "[NODE] Registered node operator wallet" prüfen','DEINE-RAILWAY-URL/api/status aufrufen um Synchronisation des Nodes zu bestätigen','Node-RPC zu MetaMask hinzufügen: Chain-ID 1926, Symbol AEQ, URL https://DEINE-URL/rpc'],rn:'Railway vergibt eine zufällige Subdomain; benutzerdefinierte Domains in den Projekteinstellungen konfigurierbar. Nur Port 8080 muss exponiert werden — P2P wird intern verwaltet.',
      s5:'5. Schnellstart — Docker',d5:'Für VPS, Cloud-VM oder lokalen Server. Voraussetzungen: Docker installiert, PostgreSQL verfügbar (verwalteter Dienst oder zweiter Container).',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\n\n# Image erstellen (~3 Minuten für Go-Kompilierung)\ndocker build -t aequitas-node .\n\n# Node starten\ndocker run -d --name aequitas-node \\\n  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n  -e RELAYER_PRIVATE_KEY="0xDEIN_PRIVATER_SCHLUESSEL" \\\n  -e RELAYER_ADDRESS="0xDEINE_ADRESSE" \\\n  -e NODE_OPERATOR_WALLET="0xDEINE_MENSCH_WALLET" \\\n  -e PEER_NODES="https://aequitas.digital" \\\n  -p 8080:8080 aequitas-node\n\ndocker logs -f aequitas-node',dn:'Container exponiert Port 8080 (API + RPC). TCP 8080 eingehend in Firewall oder Cloud-Security-Group öffnen.',
      s6:'6. Node verifizieren',v6:'Sobald der Node läuft, diese Endpunkte prüfen um Synchronisation und Gesundheit zu bestätigen.',vc:'# 1. Node-Status (Höhe sollte mit Primär-Node übereinstimmen)\ncurl https://DEINE-NODE-URL/api/status\n# Erwartet: { "height": N, "total_humans": N, "index": N }\n\n# 2. EVM JSON-RPC prüfen\ncurl -X POST https://DEINE-NODE-URL/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'\n\n# 3. In Start-Logs bestätigen:\n#    [NODE] Registered node operator wallet: 0xDEINE_WALLET\n#    Aequitas Node Running V  (Blöcke alle ~6 Sekunden)\n\n# MetaMask: RPC-URL https://DEINE-NODE-URL/rpc | Chain-ID: 1926 | Symbol: AEQ',
      s7:'7. P2P-Netzwerk & Synchronisation',p7:'PEER_NODES auf mindestens eine bekannte Bootstrap-URL setzen. Der Node verbindet sich automatisch und synchronisiert den vollständigen Chain-Zustand via libp2p-Gossip (Echtzeit) und HTTP-Pulls von Peers (Fallback). Libp2p-Multiaddresse des Primär-Nodes:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'Die HTTP-URL in PEER_NODES ist stabiler für Bootstrap. Die Multiaddresse kann sich ändern wenn der Primär-Node auf Railway neu bereitgestellt wird. Im Zweifelsfall die HTTPS-URL verwenden.',
      s8:'8. Validator-Belöhnungen erhalten',w8:'Validator-Belöhnungen kommen aus dem Validators-Pool (40% aller Protokollgebühren). Schritte um Belöhnungen zu erhalten:',b8:['Zuerst als Mensch auf Aequitas registrieren: Android-App installieren und biometrische Registrierung abschließen um Wallet-Adresse und 1.000 AEQ zu erhalten','NODE_OPERATOR_WALLET auf diese registrierte Wallet-Adresse in den Umgebungsvariablen des Nodes setzen','Node starten (oder neu starten) — er ruft RegisterNode() beim Start auf. In Logs bestätigen: "[NODE] Registered node operator wallet: 0xDEINE_WALLET"','Der Primär-Node verteilt Belöhnungen alle 24 Stunden an alle registrierten Node-Betreiber-Wallets proportional zur Blockproduktion','Sekundäre Nodes müssen die Verteilung NICHT auslösen — einfach den Node laufen lassen und synchronisiert halten','In der BETA keine Mindest-Verfügbarkeit erforderlich, aber dauerhaft offline Nodes tragen weniger zur Blockproduktion und zum Pool-Anteil bei'],
      s9:'9. Fehlerbehebung',th:['Symptom','Wahrscheinliche Ursache','Lösung'],tr:[['Höhe bleibt bei 0 nach Start','PEER_NODES nicht gesetzt oder nicht erreichbar','PEER_NODES=https://aequitas.digital setzen und neu deployen'],['"no code at address" in Logs','V7-Contract nicht im EVM deployed','RELAYER_ADDRESS prüfen; Node deployed V7 automatisch beim Start wenn fehlend'],['"NODE_OPERATOR_WALLET not set" in Logs','Fehlende Umgebungsvariable','NODE_OPERATOR_WALLET auf registrierte Mensch-Wallet-Adresse setzen'],['DATABASE_URL-Fehler beim Start','Falscher Verbindungsstring oder DB nicht erreichbar','Format prüfen: postgres://user:pass@host:5432/dbname und PostgreSQL-Erreichbarkeit'],['Port 8080 nicht erreichbar','Firewall oder Cloud-Provider-Konfiguration','TCP 8080 eingehend öffnen; Railway/Render/VPS-Port-Einstellungen prüfen'],['Docker-Build mit Go-Modul-Fehler','Kein Internetzugang beim Build','Build-Umgebung muss ausgehenden Internetzugang haben; Railway handhabt dies automatisch']],
      s10:'10. MetaMask-Konfiguration',m10:'Um deinen eigenen Node als RPC-Endpunkt in MetaMask oder einer anderen EVM-kompatiblen Wallet zu verwenden:',mh:['Feld','Wert'],mr:[['Netzwerkname','Aequitas Chain'],['RPC-URL','https://DEINE-NODE-URL/rpc'],['Chain-ID','1926  (hex: 0x786)'],['Währungssymbol','AEQ'],['Dezimalstellen','18'],['Block-Explorer','https://aequitas.digital']],
      foot:'Open Source · Erlaubnisfrei · Keine Admin-Schlüssel · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'},
    es:{title:'Guia del Operador de Nodos Aequitas',sub:'Guia completa paso a paso · Aequitas Chain (Chain ID 1926)',badge:'BETA v0.1 · Codigo Abierto · Sin permisos · Sin stake requerido',
      s1:'1. Vision General',what:'Que hace un nodo',wtxt:'Un nodo Aequitas participa plenamente en la red: produce bloques en el consenso BlockDAG, valida pruebas biometricas Groth16 de conocimiento cero para nuevos registros humanos, aplica limites de riqueza y demurrage a nivel de protocolo, sincroniza el estado con pares via libp2p + HTTP y ejecuta distribuciones diarias de pools. Cada nodo ejecuta la cadena completa: no hay clientes ligeros.',
      earn:'Que ganas',etxt:'Establece NODE_OPERATOR_WALLET en una billetera humana registrada. El Pool de Validadores acumula el 40% de todas las tarifas del protocolo. Cada 24 h el nodo primario distribuye el saldo proporcionalmente entre todos los operadores registrados. No se requiere stake.',
      s2:'2. Requisitos',rh:['Componente','Minimo','Recomendado'],rr:[['SO','Linux / host con Docker','Ubuntu 22.04 LTS'],['RAM','512 MB','1 GB'],['CPU','1 vCPU','2 vCPU'],['Almacenamiento','2 GB','10 GB SSD'],['Base de datos','PostgreSQL 14+','Railway o Supabase'],['Red','IP publica / reenvio de puerto','TCP 8080 abierto']],
      s3:'3. Variables de Entorno',e3:'Configura estas variables antes de iniciar el nodo. Las marcadas SI son obligatorias.',eh:['Variable','Proposito','Requerida?'],er:[['DATABASE_URL','Cadena de conexion PostgreSQL: postgres://user:pass@host:5432/aequitas','SI'],['RELAYER_PRIVATE_KEY','Clave privada (0x...) del EOA que firma registros on-chain','SI'],['NODE_OPERATOR_WALLET','Billetera humana registrada que recibe recompensas diarias del pool','Para recomp.'],['RELAYER_ADDRESS','Direccion EOA. Tiene fallback pero configurar explicitamente.','Recomendado'],['PORT','Puerto HTTP. Por defecto: 8080','NO'],['PEER_NODES','URLs de bootstrap (separadas por coma). Usar URL del nodo primario.','Recomendado'],['NODE_KEY','Hex 32 bytes para identidad P2P estable. Se genera automaticamente si no se establece.','NO'],['IS_PRIMARY_NODE','"true" solo en el nodo primario designado. Activa distribuciones diarias.','NO'],['RESET_STATE','"true" borra la BD al iniciar. DESTRUCTIVO.','NO']],
      s4:'4. Inicio Rapido — Railway (Recomendado)',r4:'Railway es la forma mas rapida de comenzar. El nivel gratuito cubre los requisitos minimos para BETA. Tiempo estimado: 10-15 minutos.',rs:['Haz un fork del repo: https://github.com/hanoi96international-gif/Aequitas','Crea una cuenta en railway.app e inicia un nuevo proyecto','Haz clic en "Deploy from GitHub Repo" y selecciona tu fork','En el proyecto: + New → Database → Add PostgreSQL','Ve a tu servicio → Variables y agrega las variables de la Seccion 3','Establece PEER_NODES=https://aequitas.digital','Establece NODE_OPERATOR_WALLET=<tu billetera humana AEQ>','Establece RELAYER_PRIVATE_KEY=<tu clave privada EOA>','Haz clic en "Deploy" — el Dockerfile gestiona la compilacion (~3 min)','En los logs busca: "Aequitas Node Running" y "[NODE] Registered node operator wallet"','Abre TU-URL/api/status para confirmar que el nodo esta activo','Agrega el RPC a MetaMask: Chain ID 1926, Simbolo AEQ, URL https://TU-URL/rpc'],rn:'Railway asigna un subdominio aleatorio; los dominios personalizados se configuran en ajustes del proyecto.',
      s5:'5. Inicio Rapido — Docker',d5:'Para VPS, VM en la nube o servidor local. Requiere Docker y PostgreSQL disponibles.',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\n\n# Construir imagen (~3 min)\ndocker build -t aequitas-node .\n\n# Ejecutar nodo\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n  -e RELAYER_PRIVATE_KEY="0xTU_CLAVE_PRIVADA" \\\n  -e RELAYER_ADDRESS="0xTU_DIRECCION" \\\n  -e NODE_OPERATOR_WALLET="0xTU_BILLETERA_HUMANA" \\\n  -e PEER_NODES="https://aequitas.digital" \\\n  -p 8080:8080 aequitas-node\n\ndocker logs -f aequitas-node',dn:'El contenedor expone el puerto 8080. Abre TCP 8080 entrante en tu firewall.',
      s6:'6. Verificar el Nodo',v6:'Una vez en ejecucion, comprueba estos endpoints:',vc:'curl https://TU-NODO-URL/api/status\n# Esperado: {"height": N, "total_humans": N}\n\ncurl -X POST https://TU-NODO-URL/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'',
      s7:'7. Red P2P y Sincronizacion',p7:'Establece PEER_NODES en al menos una URL de bootstrap. El nodo sincroniza la cadena automaticamente. Multidireccion libp2p del nodo primario:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'La URL HTTP en PEER_NODES es mas estable para bootstrap. La multidireccion puede cambiar si el nodo primario se reimplementa.',
      s8:'8. Ganar Recompensas de Validador',w8:'Las recompensas provienen del Pool de Validadores (40% de todas las tarifas del protocolo). Pasos:',b8:['Registrate como humano en Aequitas: instala la app Android y completa el registro biometrico','Establece NODE_OPERATOR_WALLET en esa direccion de billetera registrada','Reinicia tu nodo y confirma en logs: "[NODE] Registered node operator wallet: 0x..."','El nodo primario distribuye recompensas cada 24 h a todos los operadores registrados','Los nodos secundarios NO necesitan activar la distribucion — solo manten tu nodo en ejecucion'],
      s9:'9. Solucion de Problemas',th:['Sintoma','Causa probable','Solucion'],tr:[['Altura permanece en 0','PEER_NODES no configurado','Establece PEER_NODES=https://aequitas.digital y redespliega'],['"no code at address" en logs','Contrato V7 no desplegado aun','Normal en el primer inicio — el nodo despliega V7 automaticamente'],['Error DATABASE_URL','Cadena de conexion incorrecta','Verifica el formato: postgres://usuario:clave@host:5432/dbname'],['Puerto 8080 no accesible','Firewall o configuracion del proveedor','Abre TCP 8080 entrante en tu firewall']],
      s10:'10. Configuracion de MetaMask',m10:'Para usar tu nodo como endpoint RPC en MetaMask:',mh:['Campo','Valor'],mr:[['Nombre de red','Aequitas Chain'],['URL RPC','https://TU-NODO-URL/rpc'],['Chain ID','1926  (hex: 0x786)'],['Simbolo','AEQ'],['Decimales','18'],['Explorador','https://aequitas.digital']],
      foot:'Codigo abierto · Sin permisos · Sin claves de administrador · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'},
    it:{title:'Guida per Operatori di Nodi Aequitas',sub:'Guida completa passo dopo passo · Aequitas Chain (Chain ID 1926)',badge:'BETA v0.1 · Open Source · Senza permessi · Nessuno stake richiesto',
      s1:'1. Panoramica',what:'Cosa fa un nodo',wtxt:'Un nodo Aequitas partecipa pienamente alla rete: produce blocchi nel consenso BlockDAG, valida prove biometriche Groth16 a conoscenza zero per le nuove registrazioni umane, applica limiti di ricchezza e demurrage a livello di protocollo, sincronizza lo stato con i peer via libp2p + HTTP ed esegue distribuzioni giornaliere dei pool. Ogni nodo esegue la catena completa: non esistono client leggeri.',
      earn:'Cosa guadagni',etxt:'Imposta NODE_OPERATOR_WALLET su un indirizzo wallet registrato come umano. Il Pool Validatori accumula il 40% di tutte le commissioni di protocollo. Ogni 24 h il nodo primario distribuisce il saldo del pool proporzionalmente tra tutti i wallet degli operatori registrati. Nessuno stake richiesto.',
      s2:'2. Requisiti',rh:['Componente','Minimo','Consigliato'],rr:[['SO','Linux / host con Docker','Ubuntu 22.04 LTS'],['RAM','512 MB','1 GB'],['CPU','1 vCPU','2 vCPU'],['Archiviazione','2 GB','10 GB SSD'],['Database','PostgreSQL 14+','Railway o Supabase'],['Rete','IP pubblica / port forward','TCP 8080 aperto']],
      s3:'3. Variabili di Ambiente',e3:'Configura queste variabili prima di avviare il nodo. Quelle contrassegnate con SI sono obbligatorie.',eh:['Variabile','Scopo','Richiesta?'],er:[['DATABASE_URL','Stringa di connessione PostgreSQL: postgres://user:pass@host:5432/aequitas','SI'],['RELAYER_PRIVATE_KEY','Chiave privata (0x...) dell\'EOA che firma le registrazioni on-chain','SI'],['NODE_OPERATOR_WALLET','Wallet umano registrato che riceve le ricompense giornaliere del pool','Per ricomp.'],['RELAYER_ADDRESS','Indirizzo EOA corrispondente a RELAYER_PRIVATE_KEY. Ha un fallback ma impostalo esplicitamente.','Consigliato'],['PORT','Porta HTTP per API + JSON-RPC. Default: 8080','NO'],['PEER_NODES','URL dei peer bootstrap (separati da virgola). Usa l\'URL del nodo primario.','Consigliato'],['NODE_KEY','Hex 32 byte per identita P2P stabile. Auto-generato se omesso.','NO'],['IS_PRIMARY_NODE','"true" solo sul nodo primario designato. Abilita distribuzioni giornaliere.','NO'],['RESET_STATE','"true" cancella il DB all\'avvio. DISTRUTTIVO.','NO']],
      s4:'4. Avvio Rapido — Railway (Consigliato)',r4:'Railway e il modo piu veloce per iniziare. Il livello gratuito soddisfa i requisiti minimi per la BETA. Tempo stimato: 10-15 minuti.',rs:['Fai un fork del repo: https://github.com/hanoi96international-gif/Aequitas','Crea un account su railway.app e avvia un nuovo progetto','Clicca "Deploy from GitHub Repo" e seleziona il tuo fork','Nel progetto: + New → Database → Add PostgreSQL','Vai al tuo servizio → Variables e aggiungi le variabili della Sezione 3','Imposta PEER_NODES=https://aequitas.digital','Imposta NODE_OPERATOR_WALLET=<il tuo wallet umano AEQ>','Imposta RELAYER_PRIVATE_KEY=<la tua chiave privata EOA>','Clicca "Deploy" — il Dockerfile gestisce la compilazione (~3 min)','Nei log cerca: "Aequitas Node Running" e "[NODE] Registered node operator wallet"','Apri TUO-URL/api/status per confermare che il nodo e attivo','Aggiungi il tuo RPC a MetaMask: Chain ID 1926, Simbolo AEQ, URL https://TUO-URL/rpc'],rn:'Railway assegna un sottodominio casuale; i domini personalizzati si configurano nelle impostazioni del progetto.',
      s5:'5. Avvio Rapido — Docker',d5:'Per VPS, VM cloud o server locale. Prerequisiti: Docker installato e PostgreSQL disponibile.',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\n\n# Crea immagine (~3 min)\ndocker build -t aequitas-node .\n\n# Avvia nodo\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n  -e RELAYER_PRIVATE_KEY="0xLA_TUA_CHIAVE_PRIVATA" \\\n  -e RELAYER_ADDRESS="0xIL_TUO_INDIRIZZO" \\\n  -e NODE_OPERATOR_WALLET="0xIL_TUO_WALLET_UMANO" \\\n  -e PEER_NODES="https://aequitas.digital" \\\n  -p 8080:8080 aequitas-node\n\ndocker logs -f aequitas-node',dn:'Il container espone la porta 8080. Apri TCP 8080 in entrata nel firewall o nel gruppo di sicurezza cloud.',
      s6:'6. Verifica il Nodo',v6:'Una volta avviato, controlla questi endpoint per confermare che il nodo e sincronizzato.',vc:'curl https://TUO-NODO-URL/api/status\n# Atteso: {"height": N, "total_humans": N}\n\ncurl -X POST https://TUO-NODO-URL/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'',
      s7:'7. Rete P2P e Sincronizzazione',p7:'Imposta PEER_NODES su almeno un URL di bootstrap noto. Il nodo si connette e sincronizza la catena automaticamente. Multiindirizzo libp2p del nodo primario:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'L\'URL HTTP in PEER_NODES e piu stabile per il bootstrap. Il multiindirizzo puo cambiare se il nodo primario viene ridistribuito su Railway.',
      s8:'8. Guadagnare Ricompense da Validatore',w8:'Le ricompense provengono dal Pool Validatori (40% di tutte le commissioni di protocollo). Passaggi:',b8:['Prima registrati come umano su Aequitas: installa l\'app Android e completa la registrazione biometrica','Imposta NODE_OPERATOR_WALLET su quell\'indirizzo wallet registrato','Riavvia il nodo e conferma nei log: "[NODE] Registered node operator wallet: 0x..."','Il nodo primario distribuisce le ricompense ogni 24 h a tutti i wallet degli operatori registrati','I nodi secondari NON devono attivare la distribuzione — tieni solo il nodo in esecuzione'],
      s9:'9. Risoluzione dei Problemi',th:['Sintomo','Causa probabile','Soluzione'],tr:[['L\'altezza rimane a 0','PEER_NODES non impostato','Imposta PEER_NODES=https://aequitas.digital e ridistribuisci'],['"no code at address" nei log','Contratto V7 non ancora distribuito','Normale al primo avvio — il nodo distribuisce V7 automaticamente'],['Errore DATABASE_URL','Stringa di connessione errata','Verifica il formato: postgres://utente:password@host:5432/dbname'],['Porta 8080 non raggiungibile','Firewall o configurazione del provider','Apri TCP 8080 in entrata nel tuo firewall']],
      s10:'10. Configurazione MetaMask',m10:'Per usare il tuo nodo come endpoint RPC in MetaMask:',mh:['Campo','Valore'],mr:[['Nome rete','Aequitas Chain'],['URL RPC','https://TUO-NODO-URL/rpc'],['Chain ID','1926  (hex: 0x786)'],['Simbolo','AEQ'],['Decimali','18'],['Block Explorer','https://aequitas.digital']],
      foot:'Open source · Senza permessi · Senza chiavi admin · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'},
    id:{title:'Panduan Operator Node Aequitas',sub:'Panduan lengkap langkah demi langkah · Aequitas Chain (Chain ID 1926)',badge:'BETA v0.1 · Open Source · Tanpa Izin · Tidak perlu stake',
      s1:'1. Gambaran Umum',what:'Apa yang dilakukan node',wtxt:'Node Aequitas berpartisipasi penuh dalam jaringan: memproduksi blok dalam konsensus BlockDAG, memvalidasi bukti biometrik Groth16 zero-knowledge untuk pendaftaran manusia baru, menerapkan batas kekayaan dan demurrage di tingkat protokol, menyinkronkan status dengan peer via libp2p + HTTP, dan menjalankan distribusi pool harian. Setiap node menjalankan rantai penuh — tidak ada klien ringan.',
      earn:'Apa yang kamu dapatkan',etxt:'Atur NODE_OPERATOR_WALLET ke alamat wallet manusia terdaftar. Pool Validator mengumpulkan 40% dari semua biaya protokol. Setiap 24 jam, node primer mendistribusikan saldo pool secara proporsional ke semua wallet operator node terdaftar. Tidak perlu stake.',
      s2:'2. Persyaratan',rh:['Komponen','Minimum','Direkomendasikan'],rr:[['OS','Linux / host berkemampuan Docker','Ubuntu 22.04 LTS'],['RAM','512 MB','1 GB'],['CPU','1 vCPU','2 vCPU'],['Penyimpanan','2 GB','10 GB SSD'],['Database','PostgreSQL 14+','Railway atau Supabase'],['Jaringan','IP publik / port forward','TCP 8080 terbuka']],
      s3:'3. Variabel Lingkungan',e3:'Atur variabel ini sebelum memulai node. Yang ditandai YA wajib diisi.',eh:['Variabel','Tujuan','Wajib?'],er:[['DATABASE_URL','String koneksi PostgreSQL: postgres://user:pass@host:5432/aequitas','YA'],['RELAYER_PRIVATE_KEY','Kunci privat (0x...) EOA yang menandatangani pendaftaran on-chain','YA'],['NODE_OPERATOR_WALLET','Wallet manusia terdaftar yang menerima hadiah validator harian','Untuk hadiah'],['RELAYER_ADDRESS','Alamat EOA yang cocok dengan RELAYER_PRIVATE_KEY. Ada fallback tapi atur secara eksplisit.','Direkomendasikan'],['PORT','Port HTTP untuk API + JSON-RPC. Default: 8080','TIDAK'],['PEER_NODES','URL peer bootstrap (dipisah koma). Gunakan URL node primer.','Direkomendasikan'],['NODE_KEY','Hex 32 byte untuk identitas P2P stabil. Dibuat otomatis jika tidak diatur.','TIDAK'],['IS_PRIMARY_NODE','"true" hanya pada node primer yang ditunjuk. Mengaktifkan distribusi harian.','TIDAK'],['RESET_STATE','"true" menghapus database saat startup. DESTRUKTIF.','TIDAK']],
      s4:'4. Mulai Cepat — Railway (Direkomendasikan)',r4:'Railway adalah cara tercepat untuk memulai. Tingkat gratis memenuhi persyaratan minimum untuk BETA. Perkiraan waktu: 10-15 menit.',rs:['Fork repo: https://github.com/hanoi96international-gif/Aequitas','Buat akun di railway.app dan mulai proyek baru','Klik "Deploy from GitHub Repo" dan pilih fork kamu','Di proyek: + New → Database → Add PostgreSQL','Buka layanan kamu → Variables dan tambahkan variabel dari Bagian 3','Atur PEER_NODES=https://aequitas.digital','Atur NODE_OPERATOR_WALLET=<wallet manusia AEQ kamu>','Atur RELAYER_PRIVATE_KEY=<kunci privat EOA kamu>','Klik "Deploy" — Dockerfile mengelola kompilasi (~3 menit)','Di log cari: "Aequitas Node Running" dan "[NODE] Registered node operator wallet"','Buka URL-KAMU/api/status untuk memastikan node aktif','Tambahkan RPC ke MetaMask: Chain ID 1926, Simbol AEQ, URL https://URL-KAMU/rpc'],rn:'Railway menetapkan subdomain acak; domain kustom dapat diatur di pengaturan proyek.',
      s5:'5. Mulai Cepat — Docker',d5:'Untuk VPS, VM cloud, atau server lokal. Prasyarat: Docker terinstal, PostgreSQL tersedia.',dc:'git clone https://github.com/hanoi96international-gif/Aequitas\ncd Aequitas\n\n# Buat image (~3 menit)\ndocker build -t aequitas-node .\n\n# Jalankan node\ndocker run -d --name aequitas-node --restart unless-stopped \\\n  -e DATABASE_URL="postgres://user:pass@host:5432/aequitas" \\\n  -e RELAYER_PRIVATE_KEY="0xKUNCI_PRIVAT_KAMU" \\\n  -e RELAYER_ADDRESS="0xALAMAT_KAMU" \\\n  -e NODE_OPERATOR_WALLET="0xWALLET_MANUSIA_KAMU" \\\n  -e PEER_NODES="https://aequitas.digital" \\\n  -p 8080:8080 aequitas-node\n\ndocker logs -f aequitas-node',dn:'Container mengekspos port 8080. Buka TCP 8080 inbound di firewall atau security group cloud kamu.',
      s6:'6. Verifikasi Node',v6:'Setelah berjalan, periksa endpoint ini untuk memastikan node tersinkronisasi.',vc:'curl https://URL-NODE-KAMU/api/status\n# Diharapkan: {"height": N, "total_humans": N}\n\ncurl -X POST https://URL-NODE-KAMU/rpc \\\n  -H "Content-Type: application/json" \\\n  -d \'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}\'',
      s7:'7. Jaringan P2P dan Sinkronisasi',p7:'Atur PEER_NODES ke setidaknya satu URL bootstrap yang diketahui. Node terhubung dan menyinkronkan rantai penuh secara otomatis. Multialamat libp2p node primer:',pa:'/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R',pn:'URL HTTP di PEER_NODES lebih stabil untuk bootstrap. Multialamat libp2p dapat berubah jika node primer di-redeploy di Railway.',
      s8:'8. Mendapatkan Hadiah Validator',w8:'Hadiah berasal dari Pool Validator (40% dari semua biaya protokol). Langkah-langkah:',b8:['Pertama daftar sebagai manusia di Aequitas: instal aplikasi Android dan selesaikan pendaftaran biometrik','Atur NODE_OPERATOR_WALLET ke alamat wallet terdaftar tersebut','Mulai ulang node kamu dan konfirmasi di log: "[NODE] Registered node operator wallet: 0x..."','Node primer mendistribusikan hadiah setiap 24 jam ke semua wallet operator terdaftar','Node sekunder TIDAK perlu memicu distribusi — cukup jalankan node kamu'],
      s9:'9. Pemecahan Masalah',th:['Gejala','Kemungkinan Penyebab','Solusi'],tr:[['Tinggi tetap di 0','PEER_NODES tidak diatur','Atur PEER_NODES=https://aequitas.digital dan deploy ulang'],['"no code at address" di log','Kontrak V7 belum di-deploy','Normal saat pertama kali — node men-deploy V7 secara otomatis'],['Error DATABASE_URL','String koneksi salah','Periksa format: postgres://user:pass@host:5432/dbname'],['Port 8080 tidak dapat diakses','Firewall atau konfigurasi provider','Buka TCP 8080 inbound di firewall kamu']],
      s10:'10. Konfigurasi MetaMask',m10:'Untuk menggunakan node kamu sebagai endpoint RPC di MetaMask:',mh:['Kolom','Nilai'],mr:[['Nama Jaringan','Aequitas Chain'],['URL RPC','https://URL-NODE-KAMU/rpc'],['Chain ID','1926  (hex: 0x786)'],['Simbol','AEQ'],['Desimal','18'],['Block Explorer','https://aequitas.digital']],
      foot:'Open source · Tanpa izin · Tanpa kunci admin · Aequitas Chain V7 · Chain ID 1926',link:'github.com/hanoi96international-gif/Aequitas'}
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
</script>
</body>
</html>`
