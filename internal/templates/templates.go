package templates

import (
	"fmt"
	"html/template"
	"math"
	"strings"
	"time"
)

func Parse() *template.Template {
	return template.Must(template.New("").Funcs(funcMap).Parse(rawTemplates))
}

var funcMap = template.FuncMap{
	"timeAgo":        timeAgo,
	"statusColor":    statusColor,
	"statusIcon":     statusIcon,
	"priorityLabel":  priorityLabel,
	"priorityColor":  priorityColor,
	"runStatusColor": runStatusColor,
	"formatCost":     formatCost,
	"formatTokens":   formatTokens,
	"truncate":       truncate,
	"deref":          deref,
	"diffLines":      diffLines,
	"nl2br":          nl2br,
	"upper":          strings.ToUpper,
	"shortID":        func(s string) string { if len(s) > 8 { return s[:8] }; return s },
	"seq":            seq,
	"add":            func(a, b int) int { return a + b },
	"sub":            func(a, b int) int { return a - b },
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1m ago"
		}
		return fmt.Sprintf("%dm ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1h ago"
		}
		return fmt.Sprintf("%dh ago", h)
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1d ago"
		}
		return fmt.Sprintf("%dd ago", days)
	}
}

func statusColor(s string) string {
	switch s {
	case "todo":
		return "bg-zinc-600/80 text-zinc-200"
	case "in_progress":
		return "bg-blue-500/15 text-blue-400 ring-1 ring-inset ring-blue-500/25"
	case "in_review":
		return "bg-amber-500/15 text-amber-400 ring-1 ring-inset ring-amber-500/25"
	case "done":
		return "bg-emerald-500/15 text-emerald-400 ring-1 ring-inset ring-emerald-500/25"
	case "blocked":
		return "bg-red-500/15 text-red-400 ring-1 ring-inset ring-red-500/25"
	case "cancelled":
		return "bg-zinc-500/10 text-zinc-500 ring-1 ring-inset ring-zinc-500/20"
	default:
		return "bg-zinc-700 text-zinc-300"
	}
}

func statusIcon(s string) template.HTML {
	switch s {
	case "todo":
		return `<svg class="w-3 h-3 shrink-0" fill="none" viewBox="0 0 16 16"><circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.5"/></svg>`
	case "in_progress":
		return `<svg class="w-3 h-3 shrink-0" fill="none" viewBox="0 0 16 16"><circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.5"/><path d="M8 2a6 6 0 0 0 0 12" fill="currentColor" opacity=".35"/></svg>`
	case "in_review":
		return `<svg class="w-3 h-3 shrink-0" fill="none" viewBox="0 0 16 16"><circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.5"/><circle cx="8" cy="8" r="2.5" fill="currentColor" opacity=".4"/></svg>`
	case "done":
		return `<svg class="w-3 h-3 shrink-0" fill="none" viewBox="0 0 16 16"><circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.5"/><path d="M5.5 8l2 2 3-3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>`
	case "blocked":
		return `<svg class="w-3 h-3 shrink-0" fill="none" viewBox="0 0 16 16"><circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.5"/><path d="M5.5 5.5l5 5M10.5 5.5l-5 5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>`
	case "cancelled":
		return `<svg class="w-3 h-3 shrink-0" fill="none" viewBox="0 0 16 16"><circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.5" stroke-dasharray="3 2"/></svg>`
	default:
		return ""
	}
}

func priorityLabel(p int) string {
	switch p {
	case 1:
		return "Low"
	case 2:
		return "Medium"
	case 3:
		return "High"
	case 4:
		return "Urgent"
	default:
		return ""
	}
}

func priorityColor(p int) string {
	switch p {
	case 1:
		return "text-emerald-400"
	case 2:
		return "text-yellow-400"
	case 3:
		return "text-orange-400"
	case 4:
		return "text-red-400"
	default:
		return "text-zinc-500"
	}
}

func runStatusColor(s string) string {
	switch s {
	case "running":
		return "bg-blue-500/15 text-blue-400 ring-1 ring-inset ring-blue-500/25"
	case "completed":
		return "bg-emerald-500/15 text-emerald-400 ring-1 ring-inset ring-emerald-500/25"
	case "failed":
		return "bg-red-500/15 text-red-400 ring-1 ring-inset ring-red-500/25"
	case "cancelled":
		return "bg-zinc-500/10 text-zinc-500 ring-1 ring-inset ring-zinc-500/20"
	default:
		return "bg-zinc-700 text-zinc-300"
	}
}

func formatCost(f float64) string {
	if f == 0 {
		return "$0.00"
	}
	if f < 0.01 {
		return fmt.Sprintf("$%.4f", f)
	}
	return fmt.Sprintf("$%.2f", f)
}

func formatTokens(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(n)/1000000)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type DiffLine struct {
	Class   string
	Content string
}

func diffLines(diff string) []DiffLine {
	if diff == "" {
		return nil
	}
	var lines []DiffLine
	for _, line := range strings.Split(diff, "\n") {
		dl := DiffLine{Content: line}
		switch {
		case strings.HasPrefix(line, "@@"):
			dl.Class = "text-indigo-400 bg-indigo-950/20"
		case strings.HasPrefix(line, "+++"), strings.HasPrefix(line, "---"), strings.HasPrefix(line, "diff "):
			dl.Class = "text-zinc-500 font-medium"
		case strings.HasPrefix(line, "+"):
			dl.Class = "text-emerald-400 bg-emerald-950/30"
		case strings.HasPrefix(line, "-"):
			dl.Class = "text-red-400 bg-red-950/30"
		default:
			dl.Class = "text-zinc-500"
		}
		lines = append(lines, dl)
	}
	return lines
}

func nl2br(s string) template.HTML {
	return template.HTML(strings.ReplaceAll(template.HTMLEscapeString(s), "\n", "<br>"))
}

func seq(n int) []int {
	s := make([]int, int(math.Max(0, float64(n))))
	for i := range s {
		s[i] = i
	}
	return s
}

var rawTemplates = layoutTpl + dashboardTpl + issuesTpl + issueDetailTpl + agentsTpl + agentDetailTpl + runDetailTpl

var layoutTpl = `
{{define "layout"}}<!DOCTYPE html>
<html lang="en" class="h-full">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>TLO - {{block "title" .}}Dashboard{{end}}</title>
<script src="https://cdn.tailwindcss.com"></script>
<script src="https://unpkg.com/htmx.org@2.0.4"></script>
<style>
  @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');
  body{font-family:'Inter',system-ui,-apple-system,sans-serif}
  *{scrollbar-width:thin;scrollbar-color:#3f3f46 transparent}
  ::-webkit-scrollbar{width:6px;height:6px}
  ::-webkit-scrollbar-track{background:transparent}
  ::-webkit-scrollbar-thumb{background:#3f3f46;border-radius:3px}
  .fade-in{animation:fadeIn .15s ease-out}
  @keyframes fadeIn{from{opacity:0;transform:translateY(-4px)}to{opacity:1;transform:translateY(0)}}
  .toast-in{animation:toastIn .3s ease-out}
  @keyframes toastIn{from{opacity:0;transform:translateX(100%)}to{opacity:1;transform:translateX(0)}}
  .toast-out{animation:toastOut .2s ease-in forwards}
  @keyframes toastOut{to{opacity:0;transform:translateX(100%)}}
  .kbd{font-size:.6rem;padding:1px 5px;border-radius:4px;background:#27272a;border:1px solid #3f3f46;color:#a1a1aa;font-family:ui-monospace,monospace}
  details summary::-webkit-details-marker{display:none}
  details summary{list-style:none}
  details[open] .chev{transform:rotate(90deg)}
  .line-clamp-2{display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden}
</style>
</head>
<body class="h-full bg-zinc-900 text-zinc-100 antialiased">
<div class="flex h-full">

  <!-- Sidebar -->
  <nav class="w-52 shrink-0 border-r border-zinc-800 bg-zinc-950 flex flex-col fixed h-full z-30">
    <div class="px-4 py-4 border-b border-zinc-800">
      <a href="/dashboard" class="flex items-center gap-2 text-zinc-100 font-semibold text-sm tracking-tight">
        <svg class="w-5 h-5 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.455 2.456L21.75 6l-1.036.259a3.375 3.375 0 00-2.455 2.456z"/></svg>
        TLO
      </a>
    </div>
    <div class="flex-1 py-3 px-2 space-y-0.5">
      <a href="/dashboard" class="flex items-center gap-2.5 px-3 py-1.5 rounded-md text-[13px] text-zinc-400 hover:text-zinc-100 hover:bg-zinc-800/80 transition-colors">
        <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-4 0h4"/></svg>
        Dashboard
      </a>
      <a href="/issues" class="flex items-center gap-2.5 px-3 py-1.5 rounded-md text-[13px] text-zinc-400 hover:text-zinc-100 hover:bg-zinc-800/80 transition-colors">
        <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/></svg>
        Issues
      </a>
      <a href="/agents" class="flex items-center gap-2.5 px-3 py-1.5 rounded-md text-[13px] text-zinc-400 hover:text-zinc-100 hover:bg-zinc-800/80 transition-colors">
        <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
        Agents
      </a>
    </div>
    <div class="px-3 py-3 border-t border-zinc-800">
      <button onclick="togglePalette()" class="w-full flex items-center justify-between px-3 py-1.5 rounded-md text-xs text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800/80 transition-colors">
        <span>Search</span>
        <span class="flex gap-0.5"><span class="kbd">&#8984;</span><span class="kbd">K</span></span>
      </button>
    </div>
  </nav>

  <!-- Main -->
  <main class="flex-1 ml-52 min-h-screen">
    <div class="max-w-6xl mx-auto px-8 py-8">
      {{block "content" .}}{{end}}
    </div>
  </main>
</div>

<!-- Command Palette -->
<div id="palette" class="fixed inset-0 z-50 hidden" onclick="if(event.target===this)togglePalette()">
  <div class="absolute inset-0 bg-black/60 backdrop-blur-sm"></div>
  <div class="relative mx-auto mt-[18vh] w-full max-w-lg">
    <div class="bg-zinc-900 border border-zinc-700/80 rounded-xl shadow-2xl overflow-hidden fade-in">
      <div class="flex items-center gap-3 px-4 py-3 border-b border-zinc-800">
        <svg class="w-4 h-4 text-zinc-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"/></svg>
        <input id="pal-input" type="text" placeholder="Search issues, agents..." autocomplete="off" class="flex-1 bg-transparent text-sm text-zinc-100 placeholder-zinc-500 outline-none" oninput="searchPalette(this.value)" onkeydown="paletteNav(event)">
        <span class="kbd text-[10px]">ESC</span>
      </div>
      <div id="pal-results" class="max-h-80 overflow-y-auto py-1">
        <a href="/issues" class="pal-item flex items-center gap-2.5 px-4 py-2 text-sm text-zinc-300 hover:bg-zinc-800/70"><span class="text-indigo-400 text-base leading-none">+</span> Create issue</a>
      </div>
      <div class="border-t border-zinc-800 px-4 py-2 flex items-center gap-4 text-[10px] text-zinc-600">
        <span><span class="kbd" style="font-size:9px">&#8593;</span><span class="kbd" style="font-size:9px">&#8595;</span> Navigate</span>
        <span><span class="kbd" style="font-size:9px">&#8629;</span> Open</span>
      </div>
    </div>
  </div>
</div>

<!-- Toast Container -->
<div id="toasts" class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 pointer-events-none" style="max-width:360px"></div>

<script>
(function(){
  var rt;
  function connect(){
    var es=new EventSource('/events');
    es.addEventListener('comment',function(e){
      try{var d=JSON.parse(e.data);var p=d.body.length>80?d.body.substring(0,80)+'...':d.body;toast(d.author+' commented',p,'/issues/'+d.issue_key)}catch(x){}
    });
    es.addEventListener('run_complete',function(e){
      try{var d=JSON.parse(e.data);toast('Run '+d.status,(d.agent_name||'Agent')+' finished',d.run_id?'/runs/'+d.run_id:null)}catch(x){}
    });
    es.onerror=function(){es.close();if(rt)clearTimeout(rt);rt=setTimeout(connect,3000)};
  }
  connect();
})();

function toast(title,body,href){
  var c=document.getElementById('toasts'),el=document.createElement('div');
  el.className='pointer-events-auto bg-zinc-800 border border-zinc-700/80 rounded-lg px-4 py-3 shadow-xl toast-in cursor-default';
  var h='<div class="text-[13px] font-medium text-zinc-100">'+esc(title)+'</div>';
  if(body)h+='<div class="text-xs text-zinc-400 mt-0.5 line-clamp-2">'+esc(body)+'</div>';
  if(href){el.innerHTML='<a href="'+href+'" class="block">'+h+'</a>';el.style.cursor='pointer'}else{el.innerHTML=h}
  c.appendChild(el);
  setTimeout(function(){el.classList.remove('toast-in');el.classList.add('toast-out');setTimeout(function(){el.remove()},250)},5000);
}
function esc(s){var d=document.createElement('div');d.textContent=s;return d.innerHTML}

var palOpen=false,palIdx=-1,palItems=[];
var debounce;
function togglePalette(){
  palOpen=!palOpen;
  var el=document.getElementById('palette');
  if(palOpen){
    el.classList.remove('hidden');
    var inp=document.getElementById('pal-input');inp.value='';inp.focus();
    palIdx=-1;palItems=[];
    document.getElementById('pal-results').innerHTML='<a href="/issues" class="pal-item flex items-center gap-2.5 px-4 py-2 text-sm text-zinc-300 hover:bg-zinc-800/70"><span class="text-indigo-400 text-base leading-none">+</span> Create issue</a>';
  }else{el.classList.add('hidden')}
}
document.addEventListener('keydown',function(e){
  if((e.metaKey||e.ctrlKey)&&e.key==='k'){e.preventDefault();togglePalette()}
  if(e.key==='Escape'&&palOpen)togglePalette();
});
function searchPalette(q){
  palIdx=-1;
  if(debounce)clearTimeout(debounce);
  if(!q.trim()){
    document.getElementById('pal-results').innerHTML='<a href="/issues" class="pal-item flex items-center gap-2.5 px-4 py-2 text-sm text-zinc-300 hover:bg-zinc-800/70"><span class="text-indigo-400 text-base leading-none">+</span> Create issue</a>';
    palItems=[];return;
  }
  debounce=setTimeout(function(){
    fetch('/search?q='+encodeURIComponent(q)).then(function(r){return r.json()}).then(function(items){
      palItems=items||[];renderPal(q);
    });
  },120);
}
function renderPal(q){
  var el=document.getElementById('pal-results');
  if(!palItems.length){
    el.innerHTML='<div class="px-4 py-6 text-center text-sm text-zinc-500">No results</div>'+
      '<a href="#" onclick="event.preventDefault();createFromPal(\''+esc(q||'')+'\')" class="pal-item flex items-center gap-2.5 px-4 py-2 text-sm text-zinc-400 hover:bg-zinc-800/70 border-t border-zinc-800"><span class="text-indigo-400">+</span> Create &quot;'+esc(q)+'&quot;</a>';
    return;
  }
  var h='';
  palItems.forEach(function(r,i){
    var a=i===palIdx;
    var ico=r.type==='issue'
      ?'<svg class="w-3.5 h-3.5 text-indigo-400 shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"/></svg>'
      :'<svg class="w-3.5 h-3.5 text-emerald-400 shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>';
    h+='<a href="'+r.url+'" class="pal-item flex items-center gap-2.5 px-4 py-2 text-sm '+(a?'bg-zinc-800 text-zinc-100':'text-zinc-300 hover:bg-zinc-800/60')+'">'+ico+'<span class="text-zinc-500 text-xs font-mono w-14 shrink-0">'+esc(r.key)+'</span><span class="truncate">'+esc(r.title)+'</span></a>';
  });
  h+='<div class="border-t border-zinc-800 mt-1"><a href="#" onclick="event.preventDefault();createFromPal(\''+esc(q||'')+'\')" class="pal-item flex items-center gap-2.5 px-4 py-2 text-sm text-zinc-400 hover:bg-zinc-800/70"><span class="text-indigo-400">+</span> Create &quot;'+esc(q)+'&quot;</a></div>';
  el.innerHTML=h;
}
function createFromPal(title){
  var f=document.createElement('form');f.method='POST';f.action='/issues';
  var i=document.createElement('input');i.name='title';i.value=title;
  f.appendChild(i);document.body.appendChild(f);f.submit();
}
function paletteNav(e){
  var items=document.querySelectorAll('.pal-item');
  if(!items.length)return;
  if(e.key==='ArrowDown'){e.preventDefault();palIdx=Math.min(palIdx+1,items.length-1)}
  else if(e.key==='ArrowUp'){e.preventDefault();palIdx=Math.max(palIdx-1,0)}
  else if(e.key==='Enter'&&palIdx>=0){e.preventDefault();items[palIdx].click();return}
  else return;
  items.forEach(function(el,i){el.classList.toggle('bg-zinc-800',i===palIdx);el.classList.toggle('text-zinc-100',i===palIdx)});
}
</script>
</body>
</html>{{end}}
`

var dashboardTpl = `
{{define "dashboard"}}{{template "layout" .}}{{end}}
{{define "title"}}Dashboard{{end}}
{{define "content"}}
<div>
  <div class="mb-8">
    <h1 class="text-lg font-semibold text-zinc-100 tracking-tight">Dashboard</h1>
    <p class="text-[13px] text-zinc-500 mt-0.5">Organization overview</p>
  </div>

  <!-- Stats -->
  <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 mb-8">
    <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-4">
      <div class="text-[11px] font-medium text-zinc-500 uppercase tracking-wider">Agents</div>
      <div class="mt-2 flex items-baseline gap-1.5">
        <span class="text-2xl font-semibold text-zinc-100 tabular-nums">{{.Stats.ActiveAgents}}</span>
        <span class="text-xs text-zinc-500">/ {{.Stats.TotalAgents}}</span>
      </div>
    </div>
    <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-4">
      <div class="text-[11px] font-medium text-zinc-500 uppercase tracking-wider">Open Issues</div>
      <div class="mt-2 flex items-baseline gap-1.5">
        <span class="text-2xl font-semibold text-zinc-100 tabular-nums">{{.Stats.OpenIssues}}</span>
        <span class="text-xs text-zinc-500">/ {{.Stats.TotalIssues}}</span>
      </div>
    </div>
    <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-4">
      <div class="text-[11px] font-medium text-zinc-500 uppercase tracking-wider">Running</div>
      <div class="mt-2">
        <span class="text-2xl font-semibold text-blue-400 tabular-nums">{{.Stats.RunningRuns}}</span>
      </div>
    </div>
    <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-4">
      <div class="text-[11px] font-medium text-zinc-500 uppercase tracking-wider">Cost Today</div>
      <div class="mt-2">
        <span class="text-2xl font-semibold text-zinc-100 tabular-nums">{{formatCost .Stats.TotalCostToday}}</span>
      </div>
    </div>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
    <!-- Recent Issues -->
    <div class="lg:col-span-2">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-[13px] font-medium text-zinc-300">Recent Issues</h2>
        <a href="/issues" class="text-xs text-zinc-500 hover:text-zinc-300 transition-colors">View all</a>
      </div>
      <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg divide-y divide-zinc-800/60 overflow-hidden">
        {{range .Issues}}
        <a href="/issues/{{.Key}}" class="flex items-center gap-3 px-4 py-2.5 hover:bg-zinc-800/40 transition-colors group">
          <span class="{{statusColor .Status}} inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] font-medium shrink-0">{{statusIcon .Status}}</span>
          <span class="font-mono text-xs text-zinc-500 w-14 shrink-0">{{.Key}}</span>
          <span class="text-[13px] text-zinc-200 flex-1 truncate group-hover:text-zinc-100">{{.Title}}</span>
          {{if gt .Priority 0}}<span class="text-[11px] {{priorityColor .Priority}} shrink-0">{{priorityLabel .Priority}}</span>{{end}}
          {{if .AssigneeName}}<span class="text-[11px] text-zinc-500 shrink-0">{{.AssigneeName}}</span>{{end}}
          <span class="text-[11px] text-zinc-600 shrink-0 tabular-nums">{{timeAgo .CreatedAt}}</span>
        </a>
        {{else}}
        <div class="px-4 py-10 text-center text-sm text-zinc-500">No issues yet. Press <span class="kbd">&#8984;</span><span class="kbd">K</span> to create one.</div>
        {{end}}
      </div>
    </div>

    <!-- Agents -->
    <div>
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-[13px] font-medium text-zinc-300">Agents</h2>
        <a href="/agents" class="text-xs text-zinc-500 hover:text-zinc-300 transition-colors">Manage</a>
      </div>
      <div class="space-y-1.5">
        {{range .Agents}}
        <a href="/agents/{{.Slug}}" class="flex items-center gap-3 bg-zinc-800/50 border border-zinc-800 rounded-lg px-4 py-3 hover:border-zinc-700 transition-colors group">
          <div class="w-2 h-2 rounded-full shrink-0 {{if .Active}}bg-emerald-400{{else}}bg-zinc-600{{end}}"></div>
          <div class="flex-1 min-w-0">
            <div class="text-[13px] font-medium text-zinc-200 truncate group-hover:text-zinc-100">{{.Name}}</div>
            <div class="text-[11px] text-zinc-500 flex items-center gap-1.5 mt-0.5">
              <span>{{.ArchetypeSlug}}</span>
              <span class="text-zinc-700">&middot;</span>
              <span>{{.Model}}</span>
              {{if .HeartbeatEnabled}}<span class="text-zinc-700">&middot;</span><span class="text-emerald-500">heartbeat</span>{{end}}
            </div>
          </div>
        </a>
        {{else}}
        <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg px-4 py-8 text-center text-sm text-zinc-500">No agents configured.</div>
        {{end}}
      </div>
    </div>
  </div>
</div>
{{end}}
`

var issuesTpl = `
{{define "issues"}}{{template "layout" .}}{{end}}
{{define "title"}}Issues{{end}}
{{define "content"}}
<div>
  <div class="flex items-center justify-between mb-6">
    <h1 class="text-lg font-semibold text-zinc-100 tracking-tight">Issues</h1>
    <button onclick="document.getElementById('new-issue').classList.toggle('hidden')" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-[13px] font-medium bg-indigo-600 text-white hover:bg-indigo-500 transition-colors">
      <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15"/></svg>
      New Issue
    </button>
  </div>

  <!-- Create Form -->
  <div id="new-issue" class="hidden mb-5 fade-in">
    <form method="POST" action="/issues" class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-5 space-y-4">
      <div>
        <label class="block text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Title</label>
        <input name="title" type="text" required placeholder="What needs to be done?" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 outline-none focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500 transition-shadow">
      </div>
      <div>
        <label class="block text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Description</label>
        <textarea name="description" rows="3" placeholder="Add details..." class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 outline-none focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500 resize-none transition-shadow"></textarea>
      </div>
      <div class="flex items-end gap-3">
        <div class="flex-1">
          <label class="block text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Assignee</label>
          <select name="assignee_slug" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-200 outline-none focus:ring-1 focus:ring-indigo-500">
            <option value="">Auto (CEO)</option>
            {{range .Agents}}<option value="{{.Slug}}">{{.Name}}</option>{{end}}
          </select>
        </div>
        <div class="w-36">
          <label class="block text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Priority</label>
          <select name="priority" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-200 outline-none focus:ring-1 focus:ring-indigo-500">
            <option value="0">None</option>
            <option value="1">Low</option>
            <option value="2">Medium</option>
            <option value="3">High</option>
            <option value="4">Urgent</option>
          </select>
        </div>
        <button type="submit" class="px-5 py-2 rounded-md text-sm font-medium bg-indigo-600 text-white hover:bg-indigo-500 transition-colors">Create</button>
      </div>
    </form>
  </div>

  <!-- Status Tabs -->
  <div class="flex gap-1 mb-4 border-b border-zinc-800 pb-3">
    {{$cs := .CurrentStatus}}
    <a href="/issues" class="px-3 py-1.5 rounded-md text-xs font-medium transition-colors {{if eq $cs ""}}bg-zinc-800 text-zinc-200{{else}}text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800/50{{end}}">All</a>
    <a href="/issues?status=todo" class="px-3 py-1.5 rounded-md text-xs font-medium transition-colors {{if eq $cs "todo"}}bg-zinc-800 text-zinc-200{{else}}text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800/50{{end}}">Todo</a>
    <a href="/issues?status=in_progress" class="px-3 py-1.5 rounded-md text-xs font-medium transition-colors {{if eq $cs "in_progress"}}bg-zinc-800 text-zinc-200{{else}}text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800/50{{end}}">In Progress</a>
    <a href="/issues?status=in_review" class="px-3 py-1.5 rounded-md text-xs font-medium transition-colors {{if eq $cs "in_review"}}bg-zinc-800 text-zinc-200{{else}}text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800/50{{end}}">Review</a>
    <a href="/issues?status=done" class="px-3 py-1.5 rounded-md text-xs font-medium transition-colors {{if eq $cs "done"}}bg-zinc-800 text-zinc-200{{else}}text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800/50{{end}}">Done</a>
    <a href="/issues?status=blocked" class="px-3 py-1.5 rounded-md text-xs font-medium transition-colors {{if eq $cs "blocked"}}bg-zinc-800 text-zinc-200{{else}}text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800/50{{end}}">Blocked</a>
  </div>

  <!-- Issue List -->
  <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg divide-y divide-zinc-800/60 overflow-hidden">
    {{range .Issues}}
    <a href="/issues/{{.Key}}" class="flex items-center gap-3 px-4 py-2.5 hover:bg-zinc-800/40 transition-colors group">
      <span class="{{statusColor .Status}} inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] font-medium shrink-0">{{statusIcon .Status}}</span>
      <span class="font-mono text-xs text-zinc-500 w-14 shrink-0 group-hover:text-indigo-400 transition-colors">{{.Key}}</span>
      {{if gt .Priority 0}}<span class="text-[11px] {{priorityColor .Priority}} shrink-0">{{priorityLabel .Priority}}</span>{{end}}
      <span class="text-[13px] text-zinc-200 flex-1 truncate">{{.Title}}</span>
      {{if .ParentIssueKey}}<span class="text-[11px] text-zinc-600 font-mono shrink-0">^{{deref .ParentIssueKey}}</span>{{end}}
      {{if .AssigneeName}}<span class="text-[11px] text-zinc-500 bg-zinc-800/80 px-2 py-0.5 rounded shrink-0">{{.AssigneeName}}</span>{{end}}
      <span class="px-2 py-0.5 rounded text-[11px] font-medium {{statusColor .Status}} shrink-0">{{.Status}}</span>
      <span class="text-[11px] text-zinc-600 shrink-0 tabular-nums">{{timeAgo .CreatedAt}}</span>
    </a>
    {{else}}
    <div class="px-4 py-16 text-center">
      <svg class="w-8 h-8 mx-auto text-zinc-700 mb-2" fill="none" viewBox="0 0 24 24" stroke-width="1" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/></svg>
      <p class="text-sm text-zinc-500">No issues found</p>
    </div>
    {{end}}
  </div>
</div>
{{end}}
`

var issueDetailTpl = `
{{define "issue_detail"}}{{template "layout" .}}{{end}}
{{define "title"}}{{.Issue.Key}} - {{.Issue.Title}}{{end}}
{{define "content"}}
<div>
  <div class="flex items-center gap-1.5 mb-5 text-xs">
    <a href="/issues" class="text-zinc-500 hover:text-zinc-300 transition-colors">Issues</a>
    <span class="text-zinc-700">/</span>
    <span class="text-zinc-400 font-mono">{{.Issue.Key}}</span>
  </div>

  <div class="flex gap-6">
    <!-- Main -->
    <div class="flex-1 min-w-0 space-y-6">

      <!-- Title & Description -->
      <div>
        <div class="flex items-start gap-3 mb-3">
          <span class="{{statusColor .Issue.Status}} inline-flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-medium mt-0.5 shrink-0">{{statusIcon .Issue.Status}} {{.Issue.Status}}</span>
          <h1 class="text-xl font-semibold text-zinc-100 leading-snug">{{.Issue.Title}}</h1>
        </div>
        {{if .Issue.Description}}
        <div class="text-[13px] text-zinc-400 whitespace-pre-wrap leading-relaxed bg-zinc-800/40 border border-zinc-800 rounded-lg p-4">{{.Issue.Description}}</div>
        {{end}}
      </div>

      <!-- Children -->
      {{if .Children}}
      <div>
        <h3 class="text-[13px] font-medium text-zinc-300 mb-2">Sub-issues</h3>
        <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg divide-y divide-zinc-800/60 overflow-hidden">
          {{range .Children}}
          <a href="/issues/{{.Key}}" class="flex items-center gap-3 px-4 py-2 hover:bg-zinc-800/40 transition-colors">
            <span class="{{statusColor .Status}} inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] font-medium shrink-0">{{statusIcon .Status}}</span>
            <span class="font-mono text-xs text-zinc-500 w-14 shrink-0">{{.Key}}</span>
            <span class="text-[13px] text-zinc-200 flex-1 truncate">{{.Title}}</span>
            <span class="text-[11px] text-zinc-600 tabular-nums">{{timeAgo .UpdatedAt}}</span>
          </a>
          {{end}}
        </div>
      </div>
      {{end}}

      <!-- Runs -->
      {{if .Runs}}
      <div>
        <h3 class="text-[13px] font-medium text-zinc-300 mb-2">Runs</h3>
        <div class="space-y-1.5">
          {{range .Runs}}
          <details class="bg-zinc-800/50 border border-zinc-800 rounded-lg group">
            <summary class="flex items-center gap-3 px-4 py-2.5 cursor-pointer select-none hover:bg-zinc-800/40 transition-colors">
              <svg class="w-3.5 h-3.5 text-zinc-600 transition-transform chev shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5"/></svg>
              <span class="{{runStatusColor .Status}} inline-flex items-center px-2 py-0.5 rounded text-[11px] font-medium shrink-0">{{.Status}}</span>
              <a href="/runs/{{.ID}}" class="font-mono text-xs text-zinc-400 hover:text-indigo-400 transition-colors" onclick="event.stopPropagation()">{{.ID | truncate 8}}</a>
              <span class="text-xs text-zinc-500">{{.Mode}}</span>
              <span class="flex-1"></span>
              <span class="text-xs text-zinc-500 tabular-nums">{{formatCost .TotalCostUSD}}</span>
              <span class="text-[11px] text-zinc-600 tabular-nums">{{timeAgo .StartedAt}}</span>
            </summary>
            {{if .Diff}}
            <div class="border-t border-zinc-800 p-3 overflow-x-auto">
              <pre class="text-xs font-mono leading-relaxed">{{range diffLines .Diff}}<div class="{{.Class}} px-2 rounded-sm">{{.Content}}</div>{{end}}</pre>
            </div>
            {{end}}
          </details>
          {{end}}
        </div>
      </div>
      {{end}}

      <!-- Comments -->
      <div>
        <h3 class="text-[13px] font-medium text-zinc-300 mb-2">Comments</h3>
        <div class="space-y-2 mb-3">
          {{range .Comments}}
          <div class="bg-zinc-800/40 border border-zinc-800 rounded-lg px-4 py-3">
            <div class="flex items-center gap-2 mb-1.5">
              <span class="text-xs font-medium text-zinc-300">{{.Author}}</span>
              <span class="text-[11px] text-zinc-600">{{timeAgo .CreatedAt}}</span>
            </div>
            <div class="text-[13px] text-zinc-400 whitespace-pre-wrap leading-relaxed">{{.Body}}</div>
          </div>
          {{end}}
        </div>

        <form method="POST" action="/issues/{{.Issue.Key}}" class="bg-zinc-800/40 border border-zinc-800 rounded-lg p-3">
          <input type="hidden" name="action" value="comment">
          <textarea name="body" rows="2" placeholder="Add a comment..." class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-200 placeholder-zinc-600 outline-none focus:ring-1 focus:ring-indigo-500 resize-none transition-shadow mb-2"></textarea>
          <div class="flex justify-end">
            <button type="submit" class="px-3 py-1.5 rounded-md text-xs font-medium bg-zinc-700 text-zinc-200 hover:bg-zinc-600 transition-colors">Comment</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Sidebar -->
    <div class="w-56 shrink-0 space-y-3">
      <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
        <div class="text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Status</div>
        <form method="POST" action="/issues/{{.Issue.Key}}">
          <select name="status" onchange="this.form.submit()" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-2.5 py-1.5 text-xs text-zinc-200 outline-none focus:ring-1 focus:ring-indigo-500 cursor-pointer">
            <option value="todo" {{if eq .Issue.Status "todo"}}selected{{end}}>Todo</option>
            <option value="in_progress" {{if eq .Issue.Status "in_progress"}}selected{{end}}>In Progress</option>
            <option value="in_review" {{if eq .Issue.Status "in_review"}}selected{{end}}>In Review</option>
            <option value="done" {{if eq .Issue.Status "done"}}selected{{end}}>Done</option>
            <option value="blocked" {{if eq .Issue.Status "blocked"}}selected{{end}}>Blocked</option>
            <option value="cancelled" {{if eq .Issue.Status "cancelled"}}selected{{end}}>Cancelled</option>
          </select>
        </form>
      </div>

      <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
        <div class="text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Assignee</div>
        {{if .Assignee}}
        <a href="/agents/{{.Assignee.Slug}}" class="text-[13px] text-indigo-400 hover:text-indigo-300 transition-colors">{{.Assignee.Name}}</a>
        <div class="text-[11px] text-zinc-500 mt-0.5">{{.Assignee.ArchetypeSlug}}</div>
        {{else}}
        <span class="text-xs text-zinc-500">Unassigned</span>
        {{end}}
        <form method="POST" action="/issues/{{.Issue.Key}}" class="mt-2">
          <input type="hidden" name="action" value="assign">
          <select name="assignee_slug" onchange="this.form.submit()" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-2.5 py-1.5 text-xs text-zinc-200 outline-none focus:ring-1 focus:ring-indigo-500 cursor-pointer">
            <option value="">Reassign...</option>
            {{range .Agents}}<option value="{{.Slug}}">{{.Name}}</option>{{end}}
          </select>
        </form>
      </div>

      <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
        <div class="text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Priority</div>
        <span class="text-[13px] {{priorityColor .Issue.Priority}}">{{if gt .Issue.Priority 0}}{{priorityLabel .Issue.Priority}}{{else}}None{{end}}</span>
      </div>

      {{if .Issue.ParentIssueKey}}
      <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
        <div class="text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Parent</div>
        <a href="/issues/{{deref .Issue.ParentIssueKey}}" class="text-[13px] text-indigo-400 font-mono hover:text-indigo-300 transition-colors">{{deref .Issue.ParentIssueKey}}</a>
      </div>
      {{end}}

      <!-- Actions -->
      <div class="space-y-1.5">
        <form method="POST" action="/issues/{{.Issue.Key}}">
          <input type="hidden" name="action" value="restart">
          <button type="submit" class="w-full flex items-center justify-center gap-1.5 bg-zinc-800 border border-zinc-700 hover:bg-zinc-700 text-zinc-300 text-xs px-3 py-2 rounded-md transition-colors">
            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182"/></svg>
            Restart
          </button>
        </form>
        <form method="POST" action="/issues/{{.Issue.Key}}">
          <input type="hidden" name="action" value="cancel">
          <button type="submit" class="w-full flex items-center justify-center gap-1.5 bg-zinc-800 border border-zinc-700 hover:bg-red-950/40 hover:border-red-800/50 hover:text-red-400 text-zinc-400 text-xs px-3 py-2 rounded-md transition-colors">
            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
            Cancel
          </button>
        </form>
      </div>

      <div class="text-[11px] text-zinc-600 space-y-0.5 pt-1">
        <div>Created {{timeAgo .Issue.CreatedAt}}</div>
        <div>Updated {{timeAgo .Issue.UpdatedAt}}</div>
      </div>
    </div>
  </div>
</div>
{{end}}
`

var agentsTpl = `
{{define "agents"}}{{template "layout" .}}{{end}}
{{define "title"}}Agents{{end}}
{{define "content"}}
<div>
  <div class="flex items-center justify-between mb-6">
    <div>
      <h1 class="text-lg font-semibold text-zinc-100 tracking-tight">Agents</h1>
      <p class="text-[13px] text-zinc-500 mt-0.5">{{len .Agents}} configured</p>
    </div>
    <button onclick="document.getElementById('new-agent').classList.toggle('hidden')" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-[13px] font-medium bg-indigo-600 text-white hover:bg-indigo-500 transition-colors">
      <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15"/></svg>
      New Agent
    </button>
  </div>

  <!-- Create Form -->
  <div id="new-agent" class="hidden mb-5 fade-in">
    <form method="POST" action="/agents" class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-5 space-y-4">
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Name</label>
          <input name="name" required placeholder="Agent name" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 outline-none focus:ring-1 focus:ring-indigo-500 transition-shadow">
        </div>
        <div>
          <label class="block text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Archetype</label>
          <select name="archetype_slug" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-200 outline-none focus:ring-1 focus:ring-indigo-500">
            <option value="ceo">CEO</option>
            <option value="fullstack" selected>Fullstack</option>
            <option value="backend">Backend</option>
            <option value="frontend">Frontend</option>
            <option value="architect">Architect</option>
            <option value="devops">DevOps</option>
            <option value="qa">QA</option>
            <option value="designer">Designer</option>
            <option value="product">Product</option>
            <option value="reviewer">Reviewer</option>
            <option value="other">Other</option>
          </select>
        </div>
      </div>
      <div class="grid grid-cols-3 gap-4">
        <div>
          <label class="block text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Model</label>
          <select name="model" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-200 outline-none focus:ring-1 focus:ring-indigo-500">
            <option value="sonnet">Sonnet</option>
            <option value="opus">Opus</option>
            <option value="haiku">Haiku</option>
          </select>
        </div>
        <div>
          <label class="block text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Max Turns</label>
          <input name="max_turns" type="number" value="50" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-100 outline-none focus:ring-1 focus:ring-indigo-500 transition-shadow">
        </div>
        <div>
          <label class="block text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Timeout (s)</label>
          <input name="timeout_sec" type="number" value="600" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-100 outline-none focus:ring-1 focus:ring-indigo-500 transition-shadow">
        </div>
      </div>
      <div>
        <label class="block text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1.5">Working Directory</label>
        <input name="working_dir" placeholder="." class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 outline-none focus:ring-1 focus:ring-indigo-500 transition-shadow">
      </div>
      <div class="flex items-center gap-6">
        <label class="flex items-center gap-2 text-[13px] text-zinc-300 cursor-pointer">
          <input type="checkbox" name="heartbeat_enabled" class="rounded border-zinc-600 bg-zinc-900 text-indigo-600 focus:ring-indigo-500">
          Heartbeat
        </label>
        <label class="flex items-center gap-2 text-[13px] text-zinc-300 cursor-pointer">
          <input type="checkbox" name="chrome_enabled" class="rounded border-zinc-600 bg-zinc-900 text-indigo-600 focus:ring-indigo-500">
          Chrome
        </label>
      </div>
      <div class="flex justify-end gap-2">
        <button type="button" onclick="document.getElementById('new-agent').classList.add('hidden')" class="px-3 py-1.5 rounded-md text-sm text-zinc-400 hover:text-zinc-200 hover:bg-zinc-800 transition-colors">Cancel</button>
        <button type="submit" class="px-5 py-1.5 rounded-md text-sm font-medium bg-indigo-600 text-white hover:bg-indigo-500 transition-colors">Create</button>
      </div>
    </form>
  </div>

  <!-- Grid -->
  {{if .Agents}}
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
    {{range .Agents}}
    <a href="/agents/{{.Slug}}" class="block bg-zinc-800/50 border border-zinc-800 rounded-lg p-4 hover:border-zinc-700 hover:bg-zinc-800/70 transition-all group">
      <div class="flex items-start justify-between mb-2.5">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-md bg-zinc-900 border border-zinc-700 flex items-center justify-center text-[11px] font-semibold text-zinc-400 group-hover:text-indigo-400 group-hover:border-indigo-500/30 transition-colors uppercase">{{printf "%.2s" .Name}}</div>
          <div>
            <div class="text-[13px] font-medium text-zinc-200 group-hover:text-zinc-100 transition-colors">{{.Name}}</div>
            <div class="text-[11px] text-zinc-500 font-mono">{{.Slug}}</div>
          </div>
        </div>
        <div class="w-2 h-2 rounded-full mt-1.5 {{if .Active}}bg-emerald-400{{else}}bg-zinc-600{{end}}"></div>
      </div>
      <div class="flex flex-wrap items-center gap-1.5 text-[11px]">
        <span class="px-1.5 py-0.5 bg-zinc-900/80 text-zinc-500 rounded">{{.ArchetypeSlug}}</span>
        <span class="px-1.5 py-0.5 bg-zinc-900/80 text-zinc-500 rounded">{{.Model}}</span>
        {{if .HeartbeatEnabled}}<span class="px-1.5 py-0.5 bg-emerald-950/40 text-emerald-500 rounded">heartbeat</span>{{end}}
        {{if .ChromeEnabled}}<span class="px-1.5 py-0.5 bg-blue-950/40 text-blue-500 rounded">chrome</span>{{end}}
      </div>
    </a>
    {{end}}
  </div>
  {{else}}
  <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg px-4 py-16 text-center">
    <svg class="w-8 h-8 mx-auto text-zinc-700 mb-2" fill="none" viewBox="0 0 24 24" stroke-width="1" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
    <p class="text-sm text-zinc-500">No agents configured yet</p>
  </div>
  {{end}}
</div>
{{end}}
`

var agentDetailTpl = `
{{define "agent_detail"}}{{template "layout" .}}{{end}}
{{define "title"}}{{.Agent.Name}}{{end}}
{{define "content"}}
<div>
  <div class="flex items-center gap-1.5 mb-5 text-xs">
    <a href="/agents" class="text-zinc-500 hover:text-zinc-300 transition-colors">Agents</a>
    <span class="text-zinc-700">/</span>
    <span class="text-zinc-400 font-mono">{{.Agent.Slug}}</span>
  </div>

  <div class="flex gap-6">
    <div class="flex-1 min-w-0 space-y-6">

      <!-- Header -->
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-lg bg-zinc-800 border border-zinc-700 flex items-center justify-center text-sm font-semibold text-indigo-400 uppercase">{{printf "%.2s" .Agent.Name}}</div>
        <div>
          <div class="flex items-center gap-2.5">
            <h1 class="text-xl font-semibold text-zinc-100">{{.Agent.Name}}</h1>
            <span class="px-2 py-0.5 rounded text-[11px] bg-zinc-800 text-zinc-400">{{.Agent.ArchetypeSlug}}</span>
            <div class="w-2 h-2 rounded-full {{if .Agent.Active}}bg-emerald-400{{else}}bg-zinc-600{{end}}"></div>
          </div>
          <div class="text-xs text-zinc-500 mt-0.5 font-mono">{{.Agent.Slug}}</div>
        </div>
      </div>

      <!-- Usage -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
          <div class="text-[11px] text-zinc-500 uppercase tracking-wider">Today Tokens</div>
          <div class="text-lg font-semibold text-zinc-100 mt-1 tabular-nums">{{formatTokens .TodayTokens}}</div>
        </div>
        <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
          <div class="text-[11px] text-zinc-500 uppercase tracking-wider">Today Cost</div>
          <div class="text-lg font-semibold text-zinc-100 mt-1 tabular-nums">{{formatCost .TodayCost}}</div>
        </div>
        <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
          <div class="text-[11px] text-zinc-500 uppercase tracking-wider">Total Tokens</div>
          <div class="text-lg font-semibold text-zinc-100 mt-1 tabular-nums">{{formatTokens .TotalTokens}}</div>
        </div>
        <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
          <div class="text-[11px] text-zinc-500 uppercase tracking-wider">Total Cost</div>
          <div class="text-lg font-semibold text-zinc-100 mt-1 tabular-nums">{{formatCost .TotalCost}}</div>
        </div>
      </div>

      <!-- Inbox -->
      {{if .Issues}}
      <div>
        <h3 class="text-[13px] font-medium text-zinc-300 mb-2">Inbox</h3>
        <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg divide-y divide-zinc-800/60 overflow-hidden">
          {{range .Issues}}
          <a href="/issues/{{.Key}}" class="flex items-center gap-3 px-4 py-2.5 hover:bg-zinc-800/40 transition-colors">
            <span class="{{statusColor .Status}} inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] font-medium shrink-0">{{statusIcon .Status}}</span>
            <span class="font-mono text-xs text-zinc-500 w-14 shrink-0">{{.Key}}</span>
            <span class="text-[13px] text-zinc-200 flex-1 truncate">{{.Title}}</span>
            <span class="text-[11px] {{priorityColor .Priority}} shrink-0">{{priorityLabel .Priority}}</span>
          </a>
          {{end}}
        </div>
      </div>
      {{end}}

      <!-- Assign -->
      <form method="POST" action="/agents/{{.Agent.Slug}}/assign" class="flex items-center gap-2">
        <select name="issue_key" required class="flex-1 bg-zinc-900 border border-zinc-700 rounded-md px-3 py-1.5 text-xs text-zinc-200 outline-none focus:ring-1 focus:ring-indigo-500">
          <option value="">Assign an issue...</option>
          {{range .Issues}}<option value="{{.Key}}">{{.Key}}: {{.Title}}</option>{{end}}
        </select>
        <button type="submit" class="px-3 py-1.5 rounded-md text-xs font-medium bg-indigo-600 text-white hover:bg-indigo-500 transition-colors">Assign</button>
      </form>

      <!-- Runs -->
      {{if .Runs}}
      <div>
        <h3 class="text-[13px] font-medium text-zinc-300 mb-2">Recent Runs</h3>
        <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg divide-y divide-zinc-800/60 overflow-hidden">
          {{range .Runs}}
          <a href="/runs/{{.ID}}" class="flex items-center gap-3 px-4 py-2.5 hover:bg-zinc-800/40 transition-colors">
            <span class="{{runStatusColor .Status}} inline-flex items-center px-2 py-0.5 rounded text-[11px] font-medium shrink-0">{{.Status}}</span>
            <span class="font-mono text-xs text-zinc-400 shrink-0">{{.ID | truncate 8}}</span>
            <span class="text-xs text-zinc-500">{{.Mode}}</span>
            {{if .IssueKey}}<span class="font-mono text-xs text-zinc-500">{{deref .IssueKey}}</span>{{end}}
            <span class="flex-1"></span>
            <span class="text-xs text-zinc-500 tabular-nums">{{formatCost .TotalCostUSD}}</span>
            <span class="text-[11px] text-zinc-600 tabular-nums">{{timeAgo .StartedAt}}</span>
          </a>
          {{end}}
        </div>
      </div>
      {{end}}
    </div>

    <!-- Sidebar -->
    <div class="w-56 shrink-0 space-y-3">
      <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3 space-y-2.5 text-xs">
        <div class="text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1">Config</div>
        <div class="flex justify-between"><span class="text-zinc-500">Model</span><span class="text-zinc-300 font-mono">{{.Agent.Model}}</span></div>
        <div class="flex justify-between"><span class="text-zinc-500">Max turns</span><span class="text-zinc-300">{{.Agent.MaxTurns}}</span></div>
        <div class="flex justify-between"><span class="text-zinc-500">Timeout</span><span class="text-zinc-300">{{.Agent.TimeoutSec}}s</span></div>
        <div class="flex justify-between gap-2"><span class="text-zinc-500 shrink-0">Working dir</span><span class="text-zinc-300 truncate font-mono" title="{{.Agent.WorkingDir}}">{{.Agent.WorkingDir}}</span></div>
        <div class="flex justify-between"><span class="text-zinc-500">Heartbeat</span>
          {{if .Agent.HeartbeatEnabled}}<span class="text-emerald-400">on</span>{{else}}<span class="text-zinc-600">off</span>{{end}}
        </div>
        <div class="flex justify-between"><span class="text-zinc-500">Chrome</span>
          {{if .Agent.ChromeEnabled}}<span class="text-blue-400">on</span>{{else}}<span class="text-zinc-600">off</span>{{end}}
        </div>
      </div>

      {{if .Agent.HeartbeatEnabled}}
      <form method="POST" action="/agents/{{.Agent.Slug}}/heartbeat">
        <button type="submit" class="w-full flex items-center justify-center gap-1.5 bg-zinc-800 border border-zinc-700 hover:bg-zinc-700 text-zinc-300 text-xs px-3 py-2 rounded-md transition-colors">
          <svg class="w-3.5 h-3.5 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z"/></svg>
          Trigger Heartbeat
        </button>
      </form>
      {{end}}

      <details class="bg-zinc-800/50 border border-zinc-800 rounded-lg">
        <summary class="flex items-center gap-2 px-3 py-2.5 text-xs text-zinc-400 cursor-pointer hover:text-zinc-300 transition-colors">
          <svg class="w-3.5 h-3.5 text-zinc-600 transition-transform chev shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5"/></svg>
          Edit Config
        </summary>
        <form method="POST" action="/agents/{{.Agent.Slug}}" class="p-3 pt-0 space-y-2 border-t border-zinc-800 mt-1">
          <div>
            <label class="text-[11px] text-zinc-500">Name</label>
            <input name="name" value="{{.Agent.Name}}" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-2.5 py-1.5 text-xs text-zinc-100 outline-none focus:ring-1 focus:ring-indigo-500">
          </div>
          <div>
            <label class="text-[11px] text-zinc-500">Model</label>
            <select name="model" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-2.5 py-1.5 text-xs text-zinc-200 outline-none focus:ring-1 focus:ring-indigo-500">
              <option value="sonnet" {{if eq .Agent.Model "sonnet"}}selected{{end}}>Sonnet</option>
              <option value="opus" {{if eq .Agent.Model "opus"}}selected{{end}}>Opus</option>
              <option value="haiku" {{if eq .Agent.Model "haiku"}}selected{{end}}>Haiku</option>
            </select>
          </div>
          <div>
            <label class="text-[11px] text-zinc-500">Working Dir</label>
            <input name="working_dir" value="{{.Agent.WorkingDir}}" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-2.5 py-1.5 text-xs text-zinc-100 outline-none focus:ring-1 focus:ring-indigo-500">
          </div>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="text-[11px] text-zinc-500">Max Turns</label>
              <input name="max_turns" type="number" value="{{.Agent.MaxTurns}}" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-2.5 py-1.5 text-xs text-zinc-100 outline-none focus:ring-1 focus:ring-indigo-500">
            </div>
            <div>
              <label class="text-[11px] text-zinc-500">Timeout</label>
              <input name="timeout_sec" type="number" value="{{.Agent.TimeoutSec}}" class="w-full bg-zinc-900 border border-zinc-700 rounded-md px-2.5 py-1.5 text-xs text-zinc-100 outline-none focus:ring-1 focus:ring-indigo-500">
            </div>
          </div>
          <div class="flex items-center gap-4">
            <label class="flex items-center gap-1.5 text-[11px] text-zinc-400 cursor-pointer"><input type="checkbox" name="heartbeat_enabled" {{if .Agent.HeartbeatEnabled}}checked{{end}} class="rounded border-zinc-600 bg-zinc-900"> Heartbeat</label>
            <label class="flex items-center gap-1.5 text-[11px] text-zinc-400 cursor-pointer"><input type="checkbox" name="chrome_enabled" {{if .Agent.ChromeEnabled}}checked{{end}} class="rounded border-zinc-600 bg-zinc-900"> Chrome</label>
          </div>
          <button type="submit" class="w-full bg-indigo-600 hover:bg-indigo-500 text-white text-xs px-3 py-1.5 rounded-md transition-colors">Save</button>
        </form>
      </details>

      <div class="text-[11px] text-zinc-600 space-y-0.5 pt-1">
        <div>Created {{timeAgo .Agent.CreatedAt}}</div>
        <div>Updated {{timeAgo .Agent.UpdatedAt}}</div>
      </div>
    </div>
  </div>
</div>
{{end}}
`

var runDetailTpl = `
{{define "run_detail"}}{{template "layout" .}}{{end}}
{{define "title"}}Run {{.Run.ID | truncate 8}}{{end}}
{{define "content"}}
<div>
  <div class="flex items-center gap-1.5 mb-5 text-xs">
    {{if .Agent}}<a href="/agents/{{.Agent.Slug}}" class="text-zinc-500 hover:text-zinc-300 transition-colors">{{.Agent.Name}}</a>
    <span class="text-zinc-700">/</span>{{end}}
    <span class="text-zinc-400">Run</span>
    <span class="text-zinc-700">/</span>
    <span class="text-zinc-400 font-mono">{{.Run.ID | truncate 8}}</span>
  </div>

  <div class="flex items-center gap-3 mb-6">
    <h1 class="text-xl font-semibold text-zinc-100">
      {{if .Agent}}{{.Agent.Name}}{{end}} &mdash; {{.Run.Mode}}
      {{if .Run.IssueKey}}<a href="/issues/{{deref .Run.IssueKey}}" class="text-indigo-400 font-mono text-base hover:text-indigo-300 transition-colors ml-1">{{deref .Run.IssueKey}}</a>{{end}}
    </h1>
    <span class="{{runStatusColor .Run.Status}} inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium">
      {{if eq .Run.Status "running"}}<span class="w-1.5 h-1.5 rounded-full bg-blue-400 animate-pulse"></span>{{end}}
      {{.Run.Status}}
    </span>
  </div>

  <!-- Stats -->
  <div class="grid grid-cols-2 sm:grid-cols-5 gap-3 mb-6">
    <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
      <div class="text-[11px] text-zinc-500 uppercase tracking-wider">Cost</div>
      <div class="text-sm font-semibold text-zinc-100 mt-1 tabular-nums">{{formatCost .Run.TotalCostUSD}}</div>
    </div>
    <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
      <div class="text-[11px] text-zinc-500 uppercase tracking-wider">Input</div>
      <div class="text-sm font-semibold text-zinc-100 mt-1 tabular-nums">{{formatTokens .Run.InputTokens}}</div>
    </div>
    <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
      <div class="text-[11px] text-zinc-500 uppercase tracking-wider">Output</div>
      <div class="text-sm font-semibold text-zinc-100 mt-1 tabular-nums">{{formatTokens .Run.OutputTokens}}</div>
    </div>
    <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
      <div class="text-[11px] text-zinc-500 uppercase tracking-wider">Cache Read</div>
      <div class="text-sm font-semibold text-zinc-100 mt-1 tabular-nums">{{formatTokens .Run.CacheReadTokens}}</div>
    </div>
    <div class="bg-zinc-800/50 border border-zinc-800 rounded-lg p-3">
      <div class="text-[11px] text-zinc-500 uppercase tracking-wider">Cache Write</div>
      <div class="text-sm font-semibold text-zinc-100 mt-1 tabular-nums">{{formatTokens .Run.CacheCreateTokens}}</div>
    </div>
  </div>

  <!-- Live Stdout -->
  <div class="mb-6">
    <div class="flex items-center justify-between mb-2">
      <h3 class="text-[13px] font-medium text-zinc-300">Output</h3>
      {{if eq .Run.Status "running"}}
      <span class="flex items-center gap-1.5 text-xs text-blue-400"><span class="w-1.5 h-1.5 rounded-full bg-blue-400 animate-pulse"></span> Live</span>
      {{end}}
    </div>
    <div class="bg-zinc-950 border border-zinc-800 rounded-lg p-4 min-h-[200px] overflow-auto max-h-[600px]"
         id="stdout-container"
         {{if eq .Run.Status "running"}}hx-get="/runs/{{.Run.ID}}/stdout" hx-trigger="every 2s" hx-target="#stdout-container" hx-swap="innerHTML"{{end}}>
      {{if .Run.Stdout}}
      <pre class="text-xs font-mono text-zinc-300 whitespace-pre-wrap leading-relaxed">{{.Run.Stdout}}</pre>
      {{else}}
      <div class="text-sm text-zinc-600">{{if eq .Run.Status "running"}}Waiting for output...{{else}}No output captured.{{end}}</div>
      {{end}}
    </div>
  </div>

  <!-- Diff -->
  {{if .Run.Diff}}
  <div>
    <h3 class="text-[13px] font-medium text-zinc-300 mb-2">Git Diff</h3>
    <div class="bg-zinc-950 border border-zinc-800 rounded-lg p-4 overflow-x-auto">
      <pre class="text-xs font-mono leading-relaxed">{{range diffLines .Run.Diff}}<div class="{{.Class}} px-2 rounded-sm">{{.Content}}</div>{{end}}</pre>
    </div>
  </div>
  {{end}}

  <div class="text-[11px] text-zinc-600 mt-6 space-y-0.5">
    <div>Started {{timeAgo .Run.StartedAt}}</div>
    {{if .Run.CompletedAt}}<div>Completed {{timeAgo .Run.CompletedAt}}</div>{{end}}
    <div class="font-mono">{{.Run.ID}}</div>
  </div>
</div>
{{end}}
`
