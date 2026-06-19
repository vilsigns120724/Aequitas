package keeper

const explorerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0">
<title>Aequitas — Proof of Humanity Chain</title>
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
@media(max-width:480px){.stats-grid{grid-template-columns:repeat(2,1fr)}.stat-val{font-size:1.4rem}header{height:52px}.logo-text{font-size:0.85rem;letter-spacing:2px}.badge-dag{display:none}.main-grid{padding:0 12px 12px}.hero{padding:14px 12px 0}}
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
    <div class="wbox" id="wbox"><div class="wlbl" data-i18n="conn-wallet">CONNECTED WALLET</div><div class="wadr" id="wadr">—</div></div>
    <div id="demurrage-notice" style="display:none"></div>
    <div class="pbox" id="pbox"><div class="plbl" data-i18n="proof-recv">⚡ ZK PROOF RECEIVED</div><div class="pval" id="pval" data-i18n="proof-hint">Connect wallet to register</div></div>
    <button class="rbtn bc" id="btn-conn" onclick="connectWallet()" data-i18n="btn-conn">🦊 CONNECT METAMASK</button>
    <button class="rbtn br" id="btn-reg" onclick="doRegister()" disabled data-i18n="btn-reg">🔐 REGISTER ON-CHAIN</button>
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
    <div class="wbox" id="swap-wbox"><div class="wlbl" data-i18n="conn-wallet">CONNECTED WALLET</div><div class="wadr" id="swap-wadr">—</div></div>
    <div id="demurrage-notice" style="display:none;font-size:13px;padding:10px 12px;border-radius:8px;background:rgba(255,179,0,0.1);border:1px solid rgba(255,179,0,0.3);color:var(--gold);margin:10px 0"></div>
    <div class="ic-row" style="margin:8px 0"><span class="ic-key" data-i18n="swap-your-aeq">Your AEQ</span><span class="ic-val go" id="swap-bal-aeq">—</span></div>
    <div class="ic-row" style="margin-bottom:16px"><span class="ic-key" data-i18n="swap-your-tusd">Your tUSD</span><span class="ic-val go" id="swap-bal-tusd">—</span></div>

    <div class="swap-dir" id="swap-dir-box" style="display:flex;gap:8px;margin-bottom:12px">
      <button class="rbtn" id="swap-dir-a2t" onclick="setSwapDirection('aeq_to_tusd')" data-i18n="swap-aeq-to-tusd" style="flex:1">AEQ → tUSD</button>
      <button class="rbtn" id="swap-dir-t2a" onclick="setSwapDirection('tusd_to_aeq')" data-i18n="swap-tusd-to-aeq" style="flex:1">tUSD → AEQ</button>
    </div>
    <input type="number" id="swap-amount" placeholder="Amount" oninput="updateFeeEstimate()" style="width:100%;padding:14px;border-radius:8px;border:1px solid var(--border);background:#0A1220;color:#E8EDF5;font-size:16px;margin-bottom:8px;box-sizing:border-box">
    <div id="swap-details-panel" class="sd-panel" style="display:none">
      <div class="sd-header" data-i18n="swap-details-hdr">Swap Details</div>
      <div class="sd-row"><span class="sd-key">You receive (est.)</span><span class="sd-val" id="swap-out-est" style="color:var(--neon)">—</span></div>
      <div class="sd-row"><span class="sd-key">Price impact</span><span class="sd-val" id="swap-price-impact">—</span></div>
      <div class="sd-row"><span class="sd-key" data-i18n="swap-fee-est">Protocol fee (0.1%)</span><span class="sd-val" id="swap-fee-est" style="color:var(--muted)">—</span></div>
      <div class="sd-row"><span class="sd-key">Exchange rate</span><span class="sd-val" id="swap-rate-display" style="color:var(--purple)">—</span></div>
    </div>
    <div id="swap-warn" style="display:none;font-size:13px;padding:10px 12px;border-radius:8px;background:rgba(255,179,0,0.1);border:1px solid rgba(255,179,0,0.3);color:var(--gold);margin-bottom:10px"></div>

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
    <div class="ic-title" data-i18n="swap-pool-title">AEQ / tUSD — Pool Status</div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-pool-price">Spot Price</span><span class="ic-val go" id="pool-price">—</span></div>
    <div class="ic-row"><span class="ic-key" data-i18n="swap-pool-aeq">AEQ Reserve</span><span class="ic-val p" id="pool-reserve-aeq">—</span></div>
    <div class="ic-row" style="margin-bottom:4px"><span class="ic-key" data-i18n="swap-pool-tusd">tUSD Reserve</span><span class="ic-val b" id="pool-reserve-tusd">—</span></div>
    <div style="margin:12px 0 4px">
      <div style="font-size:0.54rem;color:var(--muted);margin-bottom:6px;font-weight:600;letter-spacing:1.5px;text-transform:uppercase">Pool Composition</div>
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
      <div style="font-size:0.54rem;color:var(--purple);font-weight:700;letter-spacing:1.2px;text-transform:uppercase;margin-bottom:6px">x × y = k — Constant Product AMM</div>
      <div class="amm-formula">AEQ_reserve × tUSD_reserve = k (constant)</div>
      <div class="amm-text">When you swap AEQ for tUSD, AEQ reserve grows and tUSD reserve shrinks — their product always stays equal to k. Every swap moves the price. Larger swaps relative to pool size cause greater price impact. The 0.1% fee is taken from the input before the formula is applied, ensuring the pool earns on every trade.</div>
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
    <div class="idx-desc" data-i18n="idx-desc">The Aequitas Index is derived from the <strong style="color:var(--teal)">Gini coefficient</strong> — the international standard for measuring wealth inequality, adopted by the World Bank, OECD, and UN. Unlike a simple richest-vs-poorest ratio, the Gini coefficient captures the <em style="color:var(--text)">entire distribution</em> across every verified human simultaneously, in a single number. <strong style="color:var(--neon)">0 = perfect equality</strong> (every wallet holds exactly the same AEQ). <strong style="color:var(--red)">100 = total concentration</strong> (one wallet holds all AEQ in existence). For context: Bitcoin Gini ≈ 0.85 (Index 85) · most unequal country on Earth (South Africa) ≈ 0.63 · Scandinavia ≈ 0.25. Aequitas is mathematically engineered to stay below 20 — enforced automatically, no governance vote, no admin key required.</div>
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
        <div style="font-size:0.64rem;color:var(--muted);line-height:1.9" data-i18n="gini-what-text">Developed by Italian statistician Corrado Gini (1912). It measures wealth distribution by comparing the actual balance distribution against a hypothetical perfectly equal baseline — visualized as the Lorenz curve. The coefficient equals the ratio of the area between the Lorenz curve and the diagonal of equality to the total area below that diagonal. Scale: 0 means every person holds identical wealth. 1 means one person holds all wealth in existence. Used by the World Bank, OECD, and UN to compare countries. Reference values: Bitcoin ≈ 0.85 · South Africa (world record) ≈ 0.63 · Brazil ≈ 0.53 · USA ≈ 0.41 · Germany ≈ 0.31 · Sweden ≈ 0.27 · Aequitas target: below 0.20.</div>
      </div>
      <div style="background:var(--card2);border:1px solid rgba(139,92,246,0.2);border-radius:var(--radius-sm);padding:16px">
        <div style="font-size:0.6rem;color:var(--purple);letter-spacing:1.5px;text-transform:uppercase;margin-bottom:10px;font-weight:600" data-i18n="gini-calc-title">How is the Aequitas Index calculated?</div>
        <div style="font-size:0.62rem;color:var(--purple);font-family:var(--font-mono);text-align:center;margin:8px 0;padding:10px;background:rgba(139,92,246,0.08);border-radius:6px;border:1px solid rgba(139,92,246,0.15)">G = Σ|xi − xj| / (2 × n² × x̄)<br><span style="color:var(--muted);font-size:0.58rem">Aequitas Index = G × 100</span></div>
        <div style="font-size:0.64rem;color:var(--muted);line-height:1.9;margin-top:8px" data-i18n="gini-calc-text">All AEQ balances of verified humans are collected (x₁ through xₙ). The formula computes the mean absolute difference between every possible pair of balances, normalized by the number of people squared (n²) and the mean balance (x̄). The result ranges 0–1 and is multiplied by 100 to produce the Aequitas Index. Updated on-chain after every registration, every monthly demurrage run, every pool payout, and every wealth cap enforcement event — via the keeper calling updateGini().</div>
      </div>
    </div>
    <div style="margin-top:10px;display:grid;grid-template-columns:repeat(4,1fr);gap:8px">
      <div style="background:rgba(0,255,209,0.06);border:1px solid rgba(0,255,209,0.25);border-radius:var(--radius-sm);padding:14px;text-align:center">
        <div style="font-size:1.05rem;font-weight:700;color:var(--neon);font-family:var(--font-display)">0 – 20</div>
        <div style="font-size:0.6rem;color:var(--neon);margin-top:5px;font-weight:700;letter-spacing:0.5px">IDEAL</div>
        <div style="font-size:0.56rem;color:var(--muted);margin-top:5px;line-height:1.7">Near-perfect equality. Better than any country on Earth. Wealth cap and demurrage passively maintaining balance. No additional protocol action.</div>
      </div>
      <div style="background:rgba(96,165,250,0.06);border:1px solid rgba(96,165,250,0.25);border-radius:var(--radius-sm);padding:14px;text-align:center">
        <div style="font-size:1.05rem;font-weight:700;color:var(--blue);font-family:var(--font-display)">20 – 40</div>
        <div style="font-size:0.6rem;color:var(--blue);margin-top:5px;font-weight:700;letter-spacing:0.5px">GOOD</div>
        <div style="font-size:0.56rem;color:var(--muted);margin-top:5px;line-height:1.7">Mild inequality — comparable to Scandinavia. Redistribution mechanisms actively flattening the distribution. Demurrage and wealth cap intensifying.</div>
      </div>
      <div style="background:rgba(245,166,35,0.06);border:1px solid rgba(245,166,35,0.25);border-radius:var(--radius-sm);padding:14px;text-align:center">
        <div style="font-size:1.05rem;font-weight:700;color:var(--gold);font-family:var(--font-display)">40 – 65</div>
        <div style="font-size:0.6rem;color:var(--gold);margin-top:5px;font-weight:700;letter-spacing:0.5px">WARNING</div>
        <div style="font-size:0.56rem;color:var(--muted);margin-top:5px;line-height:1.7">Noticeable concentration — comparable to developing countries. Protocol phase advancing. Redistribution pressure at maximum for current phase.</div>
      </div>
      <div style="background:rgba(248,113,113,0.06);border:1px solid rgba(248,113,113,0.25);border-radius:var(--radius-sm);padding:14px;text-align:center">
        <div style="font-size:1.05rem;font-weight:700;color:var(--red);font-family:var(--font-display)">65 – 100</div>
        <div style="font-size:0.6rem;color:var(--red);margin-top:5px;font-weight:700;letter-spacing:0.5px">CRITICAL</div>
        <div style="font-size:0.56rem;color:var(--muted);margin-top:5px;line-height:1.7">Worse than Bitcoin (85) or any nation on Earth (max 63). Protocol at maximum intervention. Phase 3 forced. Wealth cap at 3× mean.</div>
      </div>
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
      <div style="font-size:0.6rem;color:var(--muted);margin-bottom:6px">current pool balance</div>
      <div id="pool-u" class="ubi-pool-amount">0.0000 AEQ</div>
      <div class="ubi-fill-track"><div id="ubi-fill-bar" class="ubi-fill-bar"></div></div>
      <div style="font-size:0.61rem;color:var(--muted);line-height:1.85;margin-top:6px" data-i18n="ubi-hero-desc">Split equally among all verified humans · paid every 24 h · pool resets to zero after each payout · no minimum balance required to receive</div>
    </div>

    <!-- UBI SOURCE BREAKDOWN -->
    <div style="font-size:0.54rem;color:var(--muted);letter-spacing:2.5px;text-transform:uppercase;font-weight:600;margin:16px 0 8px">How the UBI Pool fills up</div>
    <div class="ubi-src-grid">
      <div class="ubi-src-card" style="border-color:rgba(6,182,212,0.2)">
        <div class="ubi-src-pct" style="color:var(--teal)">20%</div>
        <div class="ubi-src-name" style="color:var(--teal)">Swap Fees</div>
        <div class="ubi-src-desc">Every AEQ↔tUSD swap contributes 20% of its 0.1% fee here. More trading activity = faster pool fill.</div>
      </div>
      <div class="ubi-src-card" style="border-color:rgba(245,166,35,0.2)">
        <div class="ubi-src-pct" style="color:var(--gold)">variable</div>
        <div class="ubi-src-name" style="color:var(--gold)">Demurrage</div>
        <div class="ubi-src-desc">Idle AEQ (3+ months inactive) decays at 0.5%/month. The decayed amount enters the 40/30/20/10 split — 20% goes to UBI.</div>
      </div>
      <div class="ubi-src-card" style="border-color:rgba(139,92,246,0.2)">
        <div class="ubi-src-pct" style="color:var(--purple)">variable</div>
        <div class="ubi-src-name" style="color:var(--purple)">Wealth Cap Overflow</div>
        <div class="ubi-src-desc">Wallets exceeding 25× average balance have the excess confiscated instantly. 20% flows to UBI immediately.</div>
      </div>
    </div>

    <!-- ALL FOUR POOLS GRID -->
    <div style="font-size:0.54rem;color:var(--muted);letter-spacing:2.5px;text-transform:uppercase;font-weight:600;margin:16px 0 10px">All four redistribution pools</div>
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
        <div class="pool4-amount" style="color:var(--gold)">see countdown above</div>
        <div class="pool4-timer" style="color:var(--gold)">⏰ countdown displayed above</div>
        <div class="pool4-desc" data-i18n="ubi-pool-desc">20% of swap fees + demurrage + wealth cap overflow → divided equally among all verified humans every 24 hours. Even with zero trading, demurrage and wealth cap ensure the pool always fills.</div>
      </div>
      <div class="pool4-card" style="border-color:rgba(96,165,250,0.2)" onmouseover="this.style.borderColor='rgba(96,165,250,0.4)'" onmouseout="this.style.borderColor='rgba(96,165,250,0.2)'">
        <div class="pool4-head">
          <span class="pool4-name" style="color:var(--blue)" data-i18n="treasury">TREASURY</span>
          <span class="pool4-badge">10% of fees</span>
        </div>
        <div id="pool-t" class="pool4-amount" style="color:var(--blue)">0.0000 AEQ</div>
        <div class="pool4-timer" style="color:var(--blue)">Accumulates — no timer</div>
        <div class="pool4-desc" data-i18n="treasury-desc">Protocol development, infrastructure, security audits, and future upgrades. Governed by the Aequitas team with full on-chain transparency.</div>
      </div>
    </div>
  </div>
  <div class="idx">
    <div class="idx-title" data-i18n="phases-title">Protocol Phases</div>
    <div class="idx-desc" data-i18n="phases-desc">Aequitas evolves automatically as the number of verified humans grows. Each phase adjusts the wealth cap multiplier to maintain fairness at scale. Phase transitions are triggered by human count — no voting, no governance, no admin keys.</div>
    <table class="spect">
      <tr><td><strong style="color:var(--neon)">Phase 0</strong></td><td style="color:var(--neon)" data-i18n="p0">Bootstrap · &lt;100 humans · Wealth Cap: 50× average balance · Currently active</td></tr>
      <tr><td><strong style="color:var(--blue)">Phase 1</strong></td><td style="color:var(--blue)" data-i18n="p1">Growth · 100–10,000 humans · Wealth Cap: 20× average balance</td></tr>
      <tr><td><strong style="color:var(--gold)">Phase 2</strong></td><td style="color:var(--gold)" data-i18n="p2">Stability · 10,000–1M humans · Wealth Cap: 10× average balance</td></tr>
      <tr><td><strong style="color:var(--purple)">Phase 3</strong></td><td style="color:var(--purple)" data-i18n="p3">Maturity · 1M+ humans · Wealth Cap: 3× average balance · Maximum redistribution</td></tr>
    </table>
    <div class="hlbox" data-i18n="wealth-cap-explain">The <strong>Wealth Cap</strong> is set as a multiple of the current average AEQ balance across all verified humans — not a fixed number. This means the cap automatically adjusts as the network grows and average wealth changes, always maintaining relative fairness regardless of the total supply.</div>
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
    <div class="idx-title" data-i18n="story-title">The Story of Aequitas — Why This Exists</div>
    <div class="story" data-i18n="story-text"><p>The year is 2009. Satoshi Nakamoto releases Bitcoin. For the first time, value can transfer between any two people without a bank. A genuine revolution. But something goes wrong almost immediately.</p><p>Early miners accumulate millions of coins at almost zero cost. By 2021, the top 1% of Bitcoin addresses control over 90% of all Bitcoin. Bitcoin's estimated Gini coefficient exceeds 0.85 — higher than any country on Earth. The cryptocurrency that was supposed to democratize finance created the most extreme wealth concentration in human history.</p><p><span style="color:var(--gold)">Aequitas</span> — Latin for "fairness" and "equality" — was created to answer a single question: <em style="color:var(--gold)">"What would a cryptocurrency look like if designed from first principles to be fair to every human being?"</em></p><p>The answer is simple: <strong style="color:var(--text)">Money exists because people exist. Therefore, every person should have an equal share of money simply by virtue of being human.</strong></p><p>Aequitas implements this principle mathematically. Every verified human receives 1,000 AEQ. No mining, no staking, no early-adopter advantage. The wealth cap, demurrage, and redistribution pools ensure that inequality cannot accumulate indefinitely. The Gini coefficient and Aequitas Index are calculated on-chain in real time, and the protocol adjusts automatically.</p><p>The Aequitas network launched in June 2026. Currently in Phase 0 (Bootstrap). The goal: demonstrate that money can be distributed fairly, equality maintained through mathematical governance, and financial inclusion achieved at global scale — without any central authority.</p><p><em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em></p></div>
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
      <button onclick="document.getElementById('node-guide').style.display=document.getElementById('node-guide').style.display==='none'?'block':'none'" style="display:inline-flex;align-items:center;gap:8px;background:var(--gold);color:#06091A;padding:12px 20px;border-radius:8px;font-size:0.7rem;font-weight:700;cursor:pointer;border:none;font-family:var(--font-body);transition:opacity 0.2s" onmouseover="this.style.opacity=0.87" onmouseout="this.style.opacity=1">
        📄 Node Operator Guide
      </button>
      <a href="https://github.com/hanoi96international-gif/Aequitas" target="_blank" style="display:inline-flex;align-items:center;gap:8px;background:rgba(139,92,246,0.15);color:var(--purple);border:1px solid rgba(139,92,246,0.3);padding:12px 20px;border-radius:8px;font-size:0.7rem;font-weight:700;text-decoration:none;transition:all 0.2s" onmouseover="this.style.opacity=0.87" onmouseout="this.style.opacity=1">
        🐙 View Source on GitHub
      </a>
    </div>
    <!-- INLINE NODE GUIDE -->
    <div id="node-guide" style="display:none;background:var(--card);border:1px solid rgba(245,166,35,0.2);border-radius:var(--radius);padding:24px;margin-top:4px">
      <div style="font-size:0.58rem;color:var(--gold);letter-spacing:2.5px;text-transform:uppercase;font-weight:700;margin-bottom:18px;display:flex;align-items:center;gap:8px">
        📄 Aequitas Node Operator Guide — BETA
        <span style="font-size:0.52rem;background:rgba(245,166,35,0.12);border:1px solid rgba(245,166,35,0.3);color:var(--gold);padding:2px 8px;border-radius:10px">v0.1</span>
      </div>

      <div style="display:grid;grid-template-columns:1fr 1fr;gap:10px;margin-bottom:18px">
        <div style="background:var(--card2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:14px">
          <div style="font-size:0.58rem;color:var(--neon);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:8px">What a node does</div>
          <div style="font-size:0.62rem;color:var(--muted);line-height:1.85">Produces blocks in the Aequitas BlockDAG, validates Zero-Knowledge biometric proofs, enforces wealth caps and demurrage, distributes daily pool payouts, and syncs state with peer nodes via libp2p + HTTP. Every node runs the full chain — there are no light clients.</div>
        </div>
        <div style="background:var(--card2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:14px">
          <div style="font-size:0.58rem;color:var(--neon);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:8px">What you earn</div>
          <div style="font-size:0.62rem;color:var(--muted);line-height:1.85">40% of all protocol fees (swap fees, demurrage, wealth cap overflow) are distributed to the Validators Pool and paid out daily. The more blocks you produce proportional to the network, the larger your share. There is no staking requirement — block production is permissionless.</div>
        </div>
      </div>

      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">1. Requirements</div>
      <table style="width:100%;border-collapse:collapse;margin-bottom:16px">
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)"><td style="font-size:0.62rem;color:var(--muted);padding:6px 0;width:40%">OS</td><td style="font-size:0.62rem;color:var(--text);padding:6px 0">Linux (recommended: Ubuntu 22.04) or any Docker-capable host</td></tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)"><td style="font-size:0.62rem;color:var(--muted);padding:6px 0">RAM</td><td style="font-size:0.62rem;color:var(--text);padding:6px 0">Minimum 512 MB · Recommended 1 GB (EVM engine needs headroom)</td></tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)"><td style="font-size:0.62rem;color:var(--muted);padding:6px 0">CPU</td><td style="font-size:0.62rem;color:var(--text);padding:6px 0">1 vCPU minimum (Groth16 proof verification is CPU-bound)</td></tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)"><td style="font-size:0.62rem;color:var(--muted);padding:6px 0">Storage</td><td style="font-size:0.62rem;color:var(--text);padding:6px 0">2 GB+ (chain grows with every block; PostgreSQL recommended)</td></tr>
        <tr style="border-bottom:1px solid rgba(139,92,246,0.08)"><td style="font-size:0.62rem;color:var(--muted);padding:6px 0">Database</td><td style="font-size:0.62rem;color:var(--text);padding:6px 0">PostgreSQL 14+ (Railway, Supabase, or self-hosted) — set DATABASE_URL</td></tr>
        <tr><td style="font-size:0.62rem;color:var(--muted);padding:6px 0">Network</td><td style="font-size:0.62rem;color:var(--text);padding:6px 0">Public IP or port forwarding · TCP 8080 (API + RPC) · P2P port (auto)</td></tr>
      </table>

      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">2. Environment Variables</div>
      <div style="font-size:0.62rem;font-family:var(--font-mono);background:rgba(0,0,0,0.3);border:1px solid rgba(139,92,246,0.15);border-radius:var(--radius-sm);padding:14px;margin-bottom:16px;line-height:2.2;overflow-x:auto">
        <span style="color:var(--muted)"># Required</span><br>
        <span style="color:var(--neon)">DATABASE_URL</span>=<span style="color:var(--gold)">postgres://user:pass@host:5432/aequitas</span><br>
        <span style="color:var(--neon)">RELAYER_PRIVATE_KEY</span>=<span style="color:var(--gold)">0xYOUR_PRIVATE_KEY</span>  <span style="color:var(--muted)"># EOA that signs registrations</span><br>
        <span style="color:var(--neon)">RELAYER_ADDRESS</span>=<span style="color:var(--gold)">0xYOUR_ADDRESS</span><br>
        <span style="color:var(--muted)"># Optional</span><br>
        <span style="color:var(--teal)">PORT</span>=<span style="color:var(--gold)">8080</span><br>
        <span style="color:var(--teal)">PEER_NODES</span>=<span style="color:var(--gold)">https://aequitas-production-9fba.up.railway.app</span>  <span style="color:var(--muted)"># Comma-separated bootstrap peers</span><br>
        <span style="color:var(--teal)">NODE_KEY</span>=<span style="color:var(--gold)">hex32bytes</span>  <span style="color:var(--muted)"># Stable P2P identity (generated if omitted)</span><br>
        <span style="color:var(--teal)">IS_PRIMARY_NODE</span>=<span style="color:var(--gold)">false</span>  <span style="color:var(--muted)"># true only on the network's designated primary</span><br>
        <span style="color:var(--teal)">RESET_STATE</span>=<span style="color:var(--gold)">false</span>  <span style="color:var(--muted)"># true wipes the DB on startup — destructive!</span>
      </div>

      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">3. Quick Start — Railway (recommended)</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9;margin-bottom:10px">Railway is the fastest way to get a node running. Free tier supports the minimum requirements during BETA.</div>
      <div style="font-size:0.62rem;font-family:var(--font-mono);background:rgba(0,0,0,0.3);border:1px solid rgba(139,92,246,0.15);border-radius:var(--radius-sm);padding:14px;margin-bottom:16px;line-height:2.2">
        <span style="color:var(--muted)"># 1. Fork or clone the repo</span><br>
        <span style="color:var(--neon)">git clone</span> https://github.com/hanoi96international-gif/Aequitas<br>
        <span style="color:var(--muted)"># 2. Create a new Railway project, connect your repo</span><br>
        <span style="color:var(--muted)"># 3. Add a PostgreSQL plugin to the project</span><br>
        <span style="color:var(--muted)"># 4. Set the environment variables above in Railway Settings → Variables</span><br>
        <span style="color:var(--muted)"># 5. Set PEER_NODES to the primary node URL so your node syncs on startup</span><br>
        <span style="color:var(--neon)">PEER_NODES</span>=https://aequitas-production-9fba.up.railway.app<br>
        <span style="color:var(--muted)"># 6. Deploy — the node starts producing blocks automatically</span>
      </div>

      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">4. Quick Start — Docker</div>
      <div style="font-size:0.62rem;font-family:var(--font-mono);background:rgba(0,0,0,0.3);border:1px solid rgba(139,92,246,0.15);border-radius:var(--radius-sm);padding:14px;margin-bottom:16px;line-height:2.2">
        <span style="color:var(--neon)">docker build</span> -t aequitas-node .<br>
        <span style="color:var(--neon)">docker run</span> -d --name aequitas \<br>
        &nbsp;&nbsp;-e DATABASE_URL=<span style="color:var(--gold)">postgres://...</span> \<br>
        &nbsp;&nbsp;-e RELAYER_PRIVATE_KEY=<span style="color:var(--gold)">0x...</span> \<br>
        &nbsp;&nbsp;-e RELAYER_ADDRESS=<span style="color:var(--gold)">0x...</span> \<br>
        &nbsp;&nbsp;-e PEER_NODES=https://aequitas-production-9fba.up.railway.app \<br>
        &nbsp;&nbsp;-p 8080:8080 \<br>
        &nbsp;&nbsp;aequitas-node
      </div>

      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">5. Verify Your Node</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9;margin-bottom:10px">Once running, check the health endpoint. You should see your node's block height climbing and human count matching the primary node within a few seconds of sync.</div>
      <div style="font-size:0.62rem;font-family:var(--font-mono);background:rgba(0,0,0,0.3);border:1px solid rgba(139,92,246,0.15);border-radius:var(--radius-sm);padding:14px;margin-bottom:16px;line-height:2.2">
        <span style="color:var(--neon)">curl</span> https://YOUR-NODE-URL/api/status | jq<br>
        <span style="color:var(--muted)"># Expect: { "height": ..., "total_humans": ..., "index": ... }</span><br><br>
        <span style="color:var(--muted)"># Add MetaMask network pointing to your node:</span><br>
        <span style="color:var(--muted)"># RPC URL: https://YOUR-NODE-URL/rpc</span><br>
        <span style="color:var(--muted)"># Chain ID: 1926 · Symbol: AEQ · Name: Aequitas</span>
      </div>

      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">6. P2P Networking</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9;margin-bottom:16px">Nodes discover each other via libp2p. Set <span style="color:var(--neon);font-family:var(--font-mono)">PEER_NODES</span> to at least one known bootstrap peer (the primary node URL works). Your node will then auto-discover additional peers over time. Block sync happens via both P2P gossip and a periodic HTTP pull from peers.</div>

      <div style="font-size:0.58rem;color:var(--purple);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid var(--border);padding-bottom:6px">7. Earning Rewards</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9;margin-bottom:16px">Your node must have a registered human wallet (RELAYER_ADDRESS) to receive validator rewards. Rewards from the Validators Pool (40% of all fees) are distributed daily to all active node operator wallets proportional to blocks produced. Make sure your relayer wallet is registered as a human — run the Aequitas Android app on a device and register with the same wallet address as RELAYER_ADDRESS.</div>

      <div style="font-size:0.58rem;color:var(--gold);font-weight:700;letter-spacing:1px;text-transform:uppercase;margin-bottom:10px;border-bottom:1px solid rgba(245,166,35,0.2);padding-bottom:6px">Questions / Feedback</div>
      <div style="font-size:0.62rem;color:var(--muted);line-height:1.9">Open an issue on <a href="https://github.com/hanoi96international-gif/Aequitas" target="_blank" style="color:var(--purple)">GitHub</a> or reach the Aequitas team via the repository. BETA feedback on node setup friction, performance, and documentation gaps is especially welcome.</div>
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
    <div style="font-size:0.6rem;color:var(--muted);margin-top:10px;line-height:1.7">Set in your environment: <span style="color:var(--purple);font-family:var(--font-mono)">PEER_NODES=https://aequitas-production-9fba.up.railway.app</span></div>
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
      <tr><td>RPC URL</td><td style="color:var(--blue);font-size:0.52rem">https://aequitas-production-9fba.up.railway.app/rpc</td></tr>
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
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="ca-title">Contract &amp; Network Addresses</div>
    <div style="font-size:0.65rem;color:var(--muted);line-height:1.8;margin-bottom:10px">The Aequitas V7 contract is deployed on Aequitas Chain (Chain ID 1926). It handles human registration, balance tracking, UBI distribution, and all governance parameters. The BioVerifier contract validates Groth16 proofs on-chain before any registration is accepted.</div>
    <div class="hlbox" data-i18n="ca-text">Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier (Groth16 on-chain verifier): 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Main contract): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78</div>
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
    <div class="hlbox" data-i18n="cap-box">Cap = 25× current average AEQ balance of all verified humans<br>Automatically adjusts as the network grows and balances change<br>Applies to ALL addresses except the 4 protocol pool addresses<br>Excess AEQ is instantly redistributed to the 4 redistribution pools<br>No manual intervention required — enforced at the protocol level on every incoming transfer</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="ubi-title">5. UNIVERSAL BASIC INCOME — Daily Redistribution</div>
    <div class="hlbox" data-i18n="ubi-box">Sources of UBI Pool income:<br>· 20% of all swap fees from the AEQ↔tUSD AMM pool<br>· Overflow from wealth cap enforcement<br>· Demurrage charges from inactive accounts<br>· Inactive escrow released after 4 years<br><br>Distribution: Every 24 hours, the entire UBI pool balance is divided equally among all registered verified humans. The pool resets to zero and begins filling again immediately from ongoing protocol activity.</div>
  </div>
  <div class="idx" style="margin-bottom:12px">
    <div class="idx-title" data-i18n="inf-title">6. NO ALGORITHMIC INFLATION — Fixed Supply Formula</div>
    <div class="hlbox" data-i18n="inf-box">The ONLY event that creates new AEQ: a new verified human registers.<br><br>Total Supply = Verified Humans × 1,000 AEQ<br><br>This is not a policy — it is enforced by the protocol. No admin can mint additional AEQ, no governance vote can change the issuance, no founder allocation was pre-mined. AEQ is the only cryptocurrency where the total supply is determined solely by the number of verified living humans.</div>
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
  'idx-desc':'The Aequitas Index is derived from the <strong style="color:var(--teal)">Gini coefficient</strong> — the international standard for measuring wealth inequality, adopted by the World Bank, OECD, and UN. It captures the complete balance distribution across every verified human simultaneously. <strong style="color:var(--neon)">0 = perfect equality</strong> (every wallet holds the same AEQ). <strong style="color:var(--red)">100 = total concentration</strong> (one wallet holds all AEQ). Bitcoin Gini ≈ 0.85 (Index 85) · South Africa (world record) ≈ 0.63 · Scandinavia ≈ 0.27 · Aequitas target: below 0.20 — enforced automatically, no governance required.',
  'gini-what-title':'What is the Gini Coefficient?',
  'gini-what-text':'Developed by Italian statistician Corrado Gini (1912). Measures wealth distribution by comparing actual balances against a hypothetical perfectly equal baseline — visualized as the Lorenz curve. Scale: 0 (everyone holds the same) to 1 (one person holds everything). Used by World Bank, OECD, UN to compare countries. Reference values: Bitcoin ≈ 0.85 · South Africa (world record) ≈ 0.63 · USA ≈ 0.41 · Germany ≈ 0.31 · Sweden ≈ 0.27 · Aequitas target: below 0.20.',
  'gini-calc-title':'How is the Aequitas Index calculated?',
  'gini-calc-text':'All AEQ balances of verified humans are collected. The formula computes the mean absolute difference between every possible pair of balances, normalized by population squared (n²) and the mean balance (x̄). Result 0–1 multiplied by 100 = Aequitas Index. Updated on-chain after every registration, monthly demurrage run, pool payout, and wealth cap event — via keeper calling updateGini().',
  'gini-why-title':'Why Gini — and not a simpler metric?',
  'gini-why-text':'A simple richest-vs-poorest ratio is easy to game: 10,000 wallets could show a low spread but 90% of AEQ concentrated in 100 hands — Gini detects this, a ratio does not. The coefficient captures the complete distribution across all verified humans in one auditable number. Aequitas publishes this on-chain — transparent, tamper-evident, globally verifiable. It is the primary signal for automatic phase transitions, wealth cap calibration, and redistribution intensity. No human can override the index reading or the mechanisms it triggers.',
  'curr-idx':'Current Index','bar-0':'0 — Perfect Equality','bar-100':'100 — Max Inequality',
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
  'phases-desc':'Phase transitions are triggered automatically by human count — no voting, no governance, no admin keys required.',
  'p0':'Bootstrap · &lt;100 humans · Wealth Cap: 50× average balance · Currently active',
  'p1':'Growth · 100–10,000 humans · Wealth Cap: 20× average balance',
  'p2':'Stability · 10,000–1M humans · Wealth Cap: 10× average balance',
  'p3':'Maturity · 1M+ humans · Wealth Cap: 3× average balance · Maximum redistribution',
  'wealth-cap-explain':'The Wealth Cap is set as a multiple of the current average AEQ balance across all verified humans — not a fixed number. This means it automatically adjusts as the network grows, always maintaining relative fairness regardless of total supply.',
  'demurrage-title':'Demurrage — Incentive to Circulate',
  'demurrage-desc':'Aequitas implements a demurrage mechanism inspired by historical complementary currencies. Idle AEQ balances slowly lose value to discourage hoarding and incentivize economic participation.',
  'dem-rate-k':'Decay Rate','dem-rate-v':'0.5% per month (continuous, not stepped)',
  'dem-grace-k':'Grace Period','dem-grace-v':'3 months of inactivity before decay begins',
  'dem-reset-k':'Clock Reset','dem-reset-v':'Any transfer, swap, or liquidity action resets the timer to zero',
  'dem-dest-k':'Decayed AEQ goes to','dem-dest-v':'Redistribution pools (40/30/20/10 split)',
  'dem-warn-k':'Warning System','dem-warn-v':'14-day notice (shown once) + 7-day repeated reminder at each login',
  'story-title':'The Story of Aequitas — Why This Exists',
  'story-text':'<p>The year is 2009. Satoshi Nakamoto releases Bitcoin. For the first time, value can transfer between any two people without a bank. A genuine revolution. But something goes wrong almost immediately.</p><p>Early miners accumulate millions of coins at almost zero cost. By 2021, the top 1% of Bitcoin addresses control over 90% of all Bitcoin. Bitcoin\'s estimated Gini coefficient exceeds 0.85 — higher than any country on Earth. The cryptocurrency that was supposed to democratize finance created the most extreme wealth concentration in human history.</p><p><span style="color:var(--gold)">Aequitas</span> — Latin for "fairness" and "equality" — was created to answer a single question: <em style="color:var(--gold)">"What would a cryptocurrency look like if designed from first principles to be fair to every human being?"</em></p><p>The answer is simple: <strong style="color:var(--text)">Money exists because people exist. Therefore, every person should have an equal share of money simply by virtue of being human.</strong></p><p>Aequitas implements this mathematically. Every verified human receives 1,000 AEQ. No mining, no staking, no early-adopter advantage. The wealth cap, demurrage, and redistribution pools ensure inequality cannot accumulate indefinitely. The protocol adjusts automatically as the network grows.</p><p>The Aequitas network launched in June 2026. Currently in Phase 0. The goal: demonstrate that money can be distributed fairly, equality maintained through mathematical governance, and financial inclusion achieved at global scale — without any central authority.</p><p><em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em></p>',
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
  'ca-text':'Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Main): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'poa-title':'1. PROOF OF ALIVE','poa-text':'<p>What happens to AEQ when people die or disappear? In Bitcoin, millions of BTC are permanently lost. In Aequitas, if someone is inactive for an extended period, their AEQ eventually returns to the community through the UBI pool.</p>',
  'poa-box':'Year 0-2: Normal usage<br>Year 2: Warning 1 — Guardian can respond<br>Year 2+60d: Warning 2<br>Year 2+120d: Warning 3<br>Year 2+180d: AEQ goes to PERSONAL ESCROW<br>Year 4: If still inactive — returns to UBI Pool',
  'guard-title':'2. GUARDIAN SYSTEM','guard-text':'<p>What if someone cannot access their device for months? A trusted Guardian — another verified human — can confirm they are still alive, without any transaction rights.</p>',
  'guard-box':'1 Guardian per human (must be another verified human)<br>Guardian can ONLY call confirmAlive() — zero transaction rights<br>Guardian CANNOT move funds or transfer AEQ<br>Max 3 wards · 7-day timelock · No circular relationships allowed',
  'dem-title':'3. DEMURRAGE — Anti-Hoarding Mechanism',
  'dem-box':'Rate: 0.5%/month after 3 months grace period<br>Clock resets on any transfer, swap, or liquidity action<br>Decayed AEQ redistributed to pools (not burned)',
  'dem-text':'<p>Historical precedent: The Wörgl experiment (Austria, 1932) used a demurrage currency and reduced unemployment by 25% in one year. The Chiemgauer (Germany, 2003) has operated successfully for over 20 years using a similar mechanism.</p>',
  'cap-title':'4. WEALTH CAP — Mathematical Fairness','cap-box':'Cap = 25× current average balance of all verified humans<br>Automatically adjusts as the network grows<br>Excess AEQ instantly redistributed to redistribution pools',
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
  'idx-desc':'Der Aequitas-Index wird aus dem <strong style="color:var(--teal)">Gini-Koeffizienten</strong> abgeleitet — dem internationalen Standard zur Messung wirtschaftlicher Ungleichheit, genutzt von Weltbank, OECD und UN. Er erfasst die vollständige Bilanzverteilung aller verifizierten Menschen gleichzeitig. <strong style="color:var(--neon)">0 = perfekte Gleichheit</strong> (jede Wallet hält gleich viel AEQ). <strong style="color:var(--red)">100 = totale Konzentration</strong> (eine Wallet hält alles). Bitcoin-Gini ≈ 0,85 (Index 85) · Südafrika (Weltrekord) ≈ 0,63 · Skandinavien ≈ 0,27 · Aequitas-Ziel: unter 0,20 — automatisch durchgesetzt, keine Governance nötig.',
  'gini-what-title':'Was ist der Gini-Koeffizient?',
  'gini-what-text':'Entwickelt vom italienischen Statistiker Corrado Gini (1912). Misst die Vermögensverteilung durch Vergleich mit einer perfekt gleichen Verteilung — visualisiert als Lorenz-Kurve. Skala: 0 (alle halten gleich viel) bis 1 (eine Person hält alles). Genutzt von Weltbank, OECD, UN. Referenzwerte: Bitcoin ≈ 0,85 · Südafrika (Weltrekord) ≈ 0,63 · USA ≈ 0,41 · Deutschland ≈ 0,31 · Schweden ≈ 0,27 · Aequitas-Ziel: unter 0,20.',
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
  'phases-desc':'Phasenübergänge werden automatisch durch die Menschenanzahl ausgelöst — keine Abstimmung, keine Governance, keine Admin-Schlüssel erforderlich.',
  'p0':'Bootstrap · &lt;100 Menschen · Vermögensobergrenze: 50× Durchschnittsguthaben · Derzeit aktiv',
  'p1':'Wachstum · 100–10.000 Menschen · Vermögensobergrenze: 20× Durchschnittsguthaben',
  'p2':'Stabilität · 10.000–1M Menschen · Vermögensobergrenze: 10× Durchschnittsguthaben',
  'p3':'Reife · 1M+ Menschen · Vermögensobergrenze: 3× Durchschnittsguthaben · Maximale Umverteilung',
  'wealth-cap-explain':'Die Vermögensobergrenze wird als Vielfaches des aktuellen Durchschnittsguthabens aller verifizierten Menschen berechnet — keine feste Zahl. Sie passt sich automatisch an wenn das Netzwerk wächst.',
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
  'ca-title':'Contract- & Netzwerk-Adressen','ca-text':'Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier (Groth16 On-Chain-Verifier): 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Haupt-Contract): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
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
  'cap-title':'4. VERMÖGENSOBERGRENZE — Mathematische Fairness-Durchsetzung','cap-box':'Obergrenze = 25× aktuelles Durchschnittsguthaben aller verifizierten Menschen<br>Passt sich automatisch an während das Netzwerk wächst und sich Guthaben ändern<br>Gilt für ALLE Adressen außer den 4 Protokoll-Pool-Adressen<br>Überschuss-AEQ wird sofort an die 4 Umverteilungspools weitergeleitet<br>Keine manuelle Eingriffe erforderlich — auf Protokollebene bei jeder eingehenden Überweisung erzwungen',
  'ubi-title':'5. UNIVERSELLES GRUNDEINKOMMEN — Tägliche Umverteilung','ubi-box':'Quellen des UBI-Pool-Einkommens:<br>· 20% aller Swap-Gebühren aus dem AEQ↔tUSD AMM-Pool<br>· Überschuss aus der Vermögensobergrenze-Durchsetzung<br>· Demurrage-Gebühren von inaktiven Konten<br>· Inaktive Treuhand nach 4 Jahren freigegeben<br><br>Ausschüttung: Alle 24 Stunden wird der gesamte UBI-Pool-Saldo gleichmäßig unter allen registrierten verifizierten Menschen aufgeteilt. Der Pool setzt sich auf null zurück und beginnt sofort wieder aus der laufenden Protokollaktivität aufzufüllen.',
  'inf-title':'6. KEINE ALGORITHMISCHE INFLATION — Feste Mengenformel','inf-box':'Das EINZIGE Ereignis das neues AEQ schafft: ein neuer verifizierter Mensch registriert sich.<br><br>Gesamtmenge = Verifizierte Menschen × 1.000 AEQ<br><br>Dies ist keine Richtlinie — es wird durch das Protokoll erzwungen. Kein Admin kann zusätzliches AEQ prägen, kein Governance-Votum kann die Ausgabe ändern, keine Gründer-Zuteilung wurde vorab gemint. AEQ ist die einzige Kryptowährung bei der die Gesamtmenge ausschließlich durch die Anzahl verifizierter lebender Menschen bestimmt wird.'
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
  'phases-desc':'Las transiciones de fase se activan automáticamente por el número de humanos — sin votación, sin gobernanza, sin claves de administrador.',
  'p0':'Bootstrap · &lt;100 humanos · Límite de riqueza: 50× saldo promedio · Actualmente activo',
  'p1':'Crecimiento · 100–10,000 humanos · Límite de riqueza: 20× saldo promedio',
  'p2':'Estabilidad · 10,000–1M humanos · Límite de riqueza: 10× saldo promedio',
  'p3':'Madurez · 1M+ humanos · Límite de riqueza: 3× saldo promedio · Redistribución máxima',
  'wealth-cap-explain':'El Límite de Riqueza se establece como múltiplo del saldo promedio actual de todos los humanos verificados. Se ajusta automáticamente a medida que crece la red.',
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
  'ca-title':'Contratos y Direcciones de Red','ca-text':'Cadena: Aequitas Chain (ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier (verificador Groth16 on-chain): 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (contrato principal): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
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
  'cap-title':'4. LÍMITE DE RIQUEZA — Aplicación de Justicia Matemática','cap-box':'Límite = 25× saldo promedio actual de todos los humanos verificados<br>Se ajusta automáticamente mientras la red crece y los saldos cambian<br>Se aplica a TODAS las direcciones excepto las 4 direcciones del pool de protocolo<br>El exceso de AEQ se redistribuye instantáneamente a los 4 pools de redistribución<br>Sin intervención manual — aplicado a nivel de protocolo en cada transferencia entrante',
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
  'phases-desc':'Переходы фаз запускаются автоматически по количеству людей — без голосования, без управления, без административных ключей.',
  'p0':'Начальный этап · &lt;100 людей · Лимит богатства: 50× средний баланс · Сейчас активен',
  'p1':'Рост · 100–10 000 людей · Лимит богатства: 20× средний баланс',
  'p2':'Стабильность · 10 000–1М людей · Лимит богатства: 10× средний баланс',
  'p3':'Зрелость · 1М+ людей · Лимит богатства: 3× средний баланс · Максимальное перераспределение',
  'wealth-cap-explain':'Лимит Богатства устанавливается как кратное текущего среднего баланса всех верифицированных людей. Автоматически корректируется по мере роста сети.',
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
  'ca-title':'Адреса Контрактов','ca-text':'Цепь: Aequitas Chain (ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7: 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'poa-title':'1. ДОКАЗАТЕЛЬСТВО ЖИЗНИ — Восстановление Неактивных Балансов','poa-text':'<p>Что происходит с AEQ когда люди умирают или становятся недееспособными? В Bitcoin потерянные кошельки означают навсегда потерянный объём. Aequitas решает это через многоуровневую систему: если кошелёк не проявляет активности в течение длительного периода, его баланс постепенно возвращается сообществу через пул UBI.</p>',
  'poa-box':'Год 0–2: Обычное использование — без ограничений<br>Год 2: Предупреждение 1 — Guardian может ответить от имени<br>Год 2+60д: Предупреждение 2 — нарастающая срочность<br>Год 2+120д: Предупреждение 3 — последнее уведомление<br>Год 2+180д: AEQ перемещён в личный ЭСКРОУ (ещё восстановимо)<br>Год 4: При сохранении бездействия — ЭСКРОУ в Пул UBI',
  'guard-title':'2. СИСТЕМА GUARDIAN — Человеческая Защита','guard-text':'<p>Что если кто-то госпитализирован или иначе не может получить доступ к устройству месяцами? Система Guardian позволяет доверенному лицу — другому верифицированному человеку — подтвердить что владелец кошелька жив. Guardian имеет строго нулевой финансовый доступ: он может только сбросить таймер бездействия.</p>',
  'guard-box':'1 Guardian на человека · должен быть верифицированным человеком в Aequitas<br>Guardian может ТОЛЬКО вызывать confirmAlive() — ноль прав транзакций<br>Guardian НЕ МОЖЕТ перемещать средства, переводить AEQ или получать доступ к кошельку<br>Максимум 3 подопечных · Блокировка 7 дней при назначении · Без круговых отношений',
  'dem-title':'3. ДЕМЕРЕДЖ — Механизм Против Накопления',
  'dem-box':'Ставка: 0,5%/месяц после 3 месяцев бездействия (непрерывно, не ступенчато)<br>Таймер сбрасывается при любом переводе, свопе или операции с ликвидностью<br>Decayed AEQ перераспределяется в пулы — никогда не сжигается',
  'dem-text':'<p>Демередж — стоимость хранения денег. Эксперимент Вёрглена (Австрия, 1932) сократил местную безработицу на 25% за год. Chiemgauer (Германия, 2003) работает по тому же принципу уже более 20 лет.</p>',
  'cap-title':'4. ЛИМИТ БОГАТСТВА — Математическое Обеспечение Справедливости','cap-box':'Лимит = 25× текущий средний баланс всех верифицированных людей<br>Автоматически корректируется · Применяется ко всем адресам кроме 4 протокольных пулов<br>Избыточный AEQ мгновенно перераспределяется в 4 пула · Без ручного вмешательства',
  'ubi-title':'5. УНИВЕРСАЛЬНЫЙ БАЗОВЫЙ ДОХОД — Ежедневное Перераспределение','ubi-box':'Источники: Комиссии свопов (20%) · Превышение лимита богатства · Демередж · Эскроу после 4 лет<br><br>Ежедневно: весь пул UBI делится поровну между всеми зарегистрированными людьми. Пул сбрасывается и сразу наполняется снова.',
  'inf-title':'6. НИКАКОЙ АЛГОРИТМИЧЕСКОЙ ИНФЛЯЦИИ — Фиксированная Формула','inf-box':'ЕДИНСТВЕННОЕ событие создающее новый AEQ: регистрируется новый верифицированный человек.<br><br>Общий Объём = Верифицированные Люди × 1 000 AEQ<br><br>Это не политика — обеспечивается протоколом. AEQ — единственная криптовалюта где объём определяется исключительно числом верифицированных живых людей.'
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
  'ca-title':'合约地址','ca-text':'链：Aequitas Chain（链ID：1926 · 0x786）<br>RPC：https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier：0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7：0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'poa-title':'1. 生存证明 — 非活跃余额恢复','poa-text':'<p>当人们死亡或永久失去行为能力时AEQ会怎样？在比特币中，丢失的钱包意味着永久丢失的供应量。Aequitas通过多阶段非活跃恢复系统解决这个问题：如果一个钱包长时间没有活动，其余额会逐渐通过UBI池返回社区。</p>',
  'poa-box':'第0–2年：正常使用 — 无限制<br>第2年：警告1 — 监护人可以代表回应<br>第2年+60天：警告2 — 紧迫性增加<br>第2年+120天：警告3 — 最终通知<br>第2年+180天：AEQ移至个人托管（仍可恢复）<br>第4年：如果仍不活跃 — 托管释放至UBI池',
  'guard-title':'2. 监护人系统 — 人类安全保障','guard-text':'<p>如果有人住院或因其他原因数月无法访问其设备怎么办？监护人系统允许可信任的人——另一个经过验证的人类——确认钱包所有者仍然活着。监护人拥有严格为零的财务访问权限：只能调用重置非活跃计时器的单一函数。在任何情况下都不能移动、花费或访问资金。</p>',
  'guard-box':'每人1个监护人 · 必须是Aequitas上的经过验证的人类<br>监护人只能调用confirmAlive() — 零交易权限<br>监护人不能移动资金、转移AEQ或访问钱包<br>每个监护人最多3名受监护人 · 分配7天时间锁 · 不允许循环关系',
  'dem-title':'3. 滞期费 — 防囤积机制',
  'dem-box':'费率：3个月非活跃后每月0.5%（连续，非分步）<br>任何转账、互换或流动性操作会自动重置计时器<br>衰减的AEQ重新分配到四个池中 — 从不销毁<br>14天通知显示一次 · 每次活跃会话重复7天提醒',
  'dem-text':'<p>滞期费是货币的持有成本——一种使囤积变得昂贵、流通变得有吸引力的负利率。沃尔格实验（奥地利，1932年）使用滞期费货币在一年内将当地失业率降低了25%。奥地利中央银行正因为它运作得太好而关闭了它。Chiemgauer（德国，2003年）按照相同原则成功运营了20多年。</p>',
  'cap-title':'4. 财富上限 — 数学公平执行','cap-box':'上限 = 所有经过验证的人类当前平均AEQ余额的25倍<br>随网络增长和余额变化自动调整<br>适用于除4个协议池地址外的所有地址<br>超额AEQ立即重新分配到4个再分配池<br>无需手动干预 — 在每次入账转账时在协议级别执行',
  'ubi-title':'5. 普遍基本收入 — 每日再分配','ubi-box':'UBI池收入来源：<br>· AEQ↔tUSD AMM池所有互换费用的20%<br>· 财富上限执行的溢出<br>· 非活跃账户的滞期费<br>· 4年后释放的非活跃托管<br><br>分配：每24小时，整个UBI池余额在所有注册的经过验证的人类中平均分配。池重置为零并立即开始从持续的协议活动中重新填充。',
  'inf-title':'6. 无算法通胀 — 固定供应公式','inf-box':'创建新AEQ的唯一事件：新的经过验证的人类注册。<br><br>总供应量 = 经过验证的人类 × 1,000 AEQ<br><br>这不是政策——它由协议执行。没有管理员可以铸造额外的AEQ，没有治理投票可以改变发行，没有预挖矿的创始人分配。AEQ是唯一一种总供应量完全由经过验证的活人数量决定的加密货币。'
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
  'phases-desc':'Transisi fase dipicu secara otomatis oleh jumlah manusia — tanpa pemungutan suara, tata kelola, atau kunci admin.',
  'p0':'Bootstrap · &lt;100 manusia · Batas kekayaan: 50× saldo rata-rata · Aktif saat ini',
  'p1':'Pertumbuhan · 100–10.000 manusia · Batas kekayaan: 20× saldo rata-rata',
  'p2':'Stabilitas · 10.000–1J manusia · Batas kekayaan: 10× saldo rata-rata',
  'p3':'Kematangan · 1J+ manusia · Batas kekayaan: 3× saldo rata-rata · Redistribusi maksimum',
  'wealth-cap-explain':'Batas Kekayaan ditetapkan sebagai kelipatan saldo rata-rata semua manusia terverifikasi saat ini — bukan angka tetap. Secara otomatis menyesuaikan seiring pertumbuhan jaringan.',
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
  'ca-title':'Alamat Kontrak','ca-text':'Rantai: Aequitas Chain (ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7: 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'poa-title':'1. BUKTI KEHIDUPAN — Pemulihan Saldo Tidak Aktif','poa-text':'<p>Apa yang terjadi dengan AEQ ketika orang meninggal atau menjadi tidak mampu secara permanen? Di Bitcoin, dompet yang hilang berarti pasokan yang hilang selamanya. Aequitas menyelesaikan ini melalui sistem pemulihan ketidakaktifan multi-tahap: jika dompet tidak menunjukkan aktivitas untuk jangka waktu yang lama, saldonya secara bertahap dikembalikan ke komunitas melalui pool UBI.</p>',
  'poa-box':'Tahun 0–2: Penggunaan normal — tanpa batasan<br>Tahun 2: Peringatan 1 — Guardian dapat merespons atas nama<br>Tahun 2+60h: Peringatan 2 — urgensi meningkat<br>Tahun 2+120h: Peringatan 3 — pemberitahuan terakhir<br>Tahun 2+180h: AEQ dipindahkan ke ESCROW pribadi (masih dapat dipulihkan)<br>Tahun 4: Jika masih tidak aktif — ESCROW dirilis ke Pool UBI',
  'guard-title':'2. SISTEM GUARDIAN — Perlindungan Manusia','guard-text':'<p>Bagaimana jika seseorang dirawat di rumah sakit atau tidak dapat mengakses perangkatnya selama berbulan-bulan? Sistem Guardian memungkinkan orang terpercaya — manusia terverifikasi lainnya — mengonfirmasi bahwa pemilik dompet masih hidup. Guardian memiliki nol akses keuangan: hanya dapat memanggil satu fungsi yang mereset timer ketidakaktifan. Tidak dapat memindahkan, membelanjakan, atau mengakses dana dalam keadaan apapun.</p>',
  'guard-box':'1 Guardian per manusia · harus manusia terverifikasi di Aequitas<br>Guardian HANYA dapat memanggil confirmAlive() — nol hak transaksi<br>Guardian TIDAK DAPAT memindahkan dana, mentransfer AEQ, atau mengakses dompet<br>Maksimal 3 wali per Guardian · Kunci waktu 7 hari · Tanpa hubungan melingkar',
  'dem-title':'3. DEMURRAGE — Mekanisme Anti-Penimbunan',
  'dem-box':'Tingkat: 0,5%/bulan setelah 3 bulan ketidakaktifan (berkelanjutan, tidak bertahap)<br>Timer direset secara otomatis dengan transfer, swap, atau tindakan likuiditas apapun<br>AEQ yang meluruh didistribusikan ulang ke empat pool — tidak pernah dibakar<br>Pemberitahuan 14 hari ditampilkan sekali · 7 hari diulang di setiap sesi aktif',
  'dem-text':'<p>Demurrage adalah biaya kepemilikan uang — suku bunga negatif yang membuat penimbunan mahal dan sirkulasi menarik. Eksperimen Wörgl (Austria, 1932) mengurangi pengangguran lokal 25% dalam satu tahun. Bank Sentral Austria menutupnya justru karena bekerja terlalu baik. Chiemgauer (Jerman, 2003) beroperasi dengan prinsip yang sama dengan sukses selama lebih dari 20 tahun.</p>',
  'cap-title':'4. BATAS KEKAYAAN — Penerapan Keadilan Matematis','cap-box':'Batas = 25× saldo AEQ rata-rata semua manusia terverifikasi saat ini<br>Otomatis menyesuaikan seiring pertumbuhan jaringan dan perubahan saldo<br>Berlaku untuk SEMUA alamat kecuali 4 alamat pool protokol<br>Kelebihan AEQ langsung didistribusikan ulang ke 4 pool redistribusi<br>Tanpa intervensi manual — diterapkan di tingkat protokol pada setiap transfer masuk',
  'ubi-title':'5. PENDAPATAN DASAR UNIVERSAL — Redistribusi Harian','ubi-box':'Sumber pendapatan Pool UBI:<br>· 20% semua biaya swap dari pool AMM AEQ↔tUSD<br>· Overflow dari penerapan batas kekayaan<br>· Biaya demurrage dari akun tidak aktif<br>· Escrow tidak aktif dirilis setelah 4 tahun<br><br>Distribusi: Setiap 24 jam, seluruh saldo pool UBI dibagi rata di antara semua manusia terverifikasi yang terdaftar. Pool direset ke nol dan segera mulai diisi ulang dari aktivitas protokol yang berkelanjutan.',
  'inf-title':'6. TANPA INFLASI ALGORITMIK — Formula Pasokan Tetap','inf-box':'SATU-SATUNYA peristiwa yang menciptakan AEQ baru: manusia terverifikasi baru mendaftar.<br><br>Total Pasokan = Manusia Terverifikasi × 1.000 AEQ<br><br>Ini bukan kebijakan — ini diterapkan oleh protokol. Tidak ada admin yang dapat mencetak AEQ tambahan, tidak ada suara tata kelola yang dapat mengubah penerbitan. AEQ adalah satu-satunya cryptocurrency di mana total pasokan ditentukan semata-mata oleh jumlah manusia hidup yang terverifikasi.'
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
  'phases-desc':'Le transizioni di fase vengono attivate automaticamente dal numero di umani — nessun voto, nessuna governance, nessuna chiave admin necessaria.',
  'p0':'Bootstrap · &lt;100 umani · Limite ricchezza: 50× saldo medio · Attualmente attivo',
  'p1':'Crescita · 100–10.000 umani · Limite ricchezza: 20× saldo medio',
  'p2':'Stabilità · 10.000–1M umani · Limite ricchezza: 10× saldo medio',
  'p3':'Maturità · 1M+ umani · Limite ricchezza: 3× saldo medio · Massima redistribuzione',
  'wealth-cap-explain':'Il Limite di Ricchezza è impostato come multiplo del saldo medio attuale di tutti gli umani verificati — non un numero fisso. Si adatta automaticamente man mano che la rete cresce, mantenendo sempre l\'equità relativa.',
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
  'ca-title':'Indirizzi Contratto','ca-text':'Chain: Aequitas Chain (ID: 1926 · 0x786)<br>RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>BioVerifier: 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>AequitasV7 (Principale): 0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78',
  'poa-title':'1. PROVA DI VITA — Recupero Saldi Inattivi','poa-text':'<p>Cosa succede all\'AEQ quando le persone muoiono o diventano permanentemente incapaci? In Bitcoin, i portafogli persi significano fornitura persa permanentemente. Aequitas risolve questo con un sistema di recupero dell\'inattività a più fasi: se un portafoglio non mostra attività per un periodo prolungato, il suo saldo viene gradualmente restituito alla comunità attraverso il pool UBI.</p>',
  'poa-box':'Anno 0–2: Uso normale — nessuna restrizione<br>Anno 2: Avviso 1 — il Guardian può rispondere a nome<br>Anno 2+60g: Avviso 2 — urgenza crescente<br>Anno 2+120g: Avviso 3 — avviso finale<br>Anno 2+180g: AEQ spostato in ESCROW personale (ancora recuperabile)<br>Anno 4: Se ancora inattivo — ESCROW rilasciato al Pool UBI',
  'guard-title':'2. SISTEMA GUARDIAN — Protezione Umana','guard-text':'<p>E se qualcuno è ricoverato in ospedale o non riesce ad accedere al proprio dispositivo per mesi? Il sistema Guardian permette a una persona di fiducia — un altro umano verificato — di confermare che il proprietario del portafoglio è ancora vivo. Il Guardian ha accesso finanziario strettamente nullo: può solo chiamare una singola funzione che reimposta il timer di inattività. Non può spostare, spendere o accedere ai fondi in nessuna circostanza.</p>',
  'guard-box':'1 Guardian per umano · deve essere un umano verificato su Aequitas<br>Il Guardian può SOLO chiamare confirmAlive() — zero diritti di transazione<br>Il Guardian NON PUÒ spostare fondi, trasferire AEQ o accedere al portafoglio<br>Massimo 3 tutelati per Guardian · Blocco di 7 giorni all\'assegnazione · Nessuna relazione circolare',
  'dem-title':'3. DEMURRAGE — Meccanismo Anti-Accumulo',
  'dem-box':'Tasso: 0,5%/mese dopo 3 mesi di inattività (continuo, non a gradini)<br>Il timer si azzera automaticamente con qualsiasi trasferimento, swap o azione di liquidità<br>AEQ decaduto ridistribuito ai quattro pool — mai bruciato<br>Avviso di 14 giorni mostrato una volta · 7 giorni ripetuto in ogni sessione attiva',
  'dem-text':'<p>Il demurrage è un costo di detenzione sul denaro — un tasso di interesse negativo che rende costoso accumulare e attraente la circolazione. L\'esperimento di Wörgl (Austria, 1932) usò una valuta con demurrage e ridusse la disoccupazione locale del 25% in un anno. La Banca Centrale austriaca lo chiuse proprio perché funzionava troppo bene. Il Chiemgauer (Germania, 2003) opera con lo stesso principio con successo da oltre 20 anni.</p>',
  'cap-title':'4. LIMITE DI RICCHEZZA — Applicazione dell\'Equità Matematica','cap-box':'Limite = 25× saldo AEQ medio attuale di tutti gli umani verificati<br>Si adatta automaticamente man mano che la rete cresce e i saldi cambiano<br>Si applica a TUTTI gli indirizzi tranne i 4 indirizzi del pool del protocollo<br>L\'eccesso di AEQ viene immediatamente ridistribuito ai 4 pool di redistribuzione<br>Nessun intervento manuale richiesto — applicato a livello di protocollo ad ogni trasferimento in entrata',
  'ubi-title':'5. REDDITO UNIVERSALE DI BASE — Ridistribuzione Giornaliera','ubi-box':'Fonti di reddito del Pool UBI:<br>· 20% di tutte le commissioni di swap del pool AMM AEQ↔tUSD<br>· Overflow dall\'applicazione del limite di ricchezza<br>· Addebiti di demurrage da account inattivi<br>· Escrow inattivo rilasciato dopo 4 anni<br><br>Distribuzione: Ogni 24 ore, l\'intero saldo del pool UBI viene diviso equamente tra tutti gli umani verificati registrati. Il pool si azzera e inizia immediatamente a riempirsi di nuovo dall\'attività continua del protocollo.',
  'inf-title':'6. NESSUNA INFLAZIONE ALGORITMICA — Formula di Fornitura Fissa','inf-box':'L\'UNICO evento che crea nuovo AEQ: un nuovo umano verificato si registra.<br><br>Offerta Totale = Umani Verificati × 1.000 AEQ<br><br>Questo non è una politica — è applicato dal protocollo. Nessun amministratore può coniare AEQ aggiuntivo, nessun voto di governance può modificare l\'emissione. AEQ è l\'unica criptovaluta in cui l\'offerta totale è determinata esclusivamente dal numero di esseri umani vivi verificati.'
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
        bioHash: pendingBioHash || '',
        bioHashKey: proofData.bioHashKey || ''
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
