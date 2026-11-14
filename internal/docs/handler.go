package docs

import (
	_ "embed"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

//go:embed openapi.yaml
var openAPISpec []byte

//go:embed swagger.html
var swaggerUI []byte

func RegisterRoutes(e *echo.Echo) {
	e.GET("/", serveHome)
	e.GET("/health", serveHealthPage)
	e.GET("/healthz", serveHealthJSON)
	e.GET("/api/v1/docs", serveSwaggerUI)
	e.GET("/api/v1/docs/", serveSwaggerUI)
	e.GET("/api/v1/docs/openapi.yaml", serveOpenAPISpec)
}

func serveHealthPage(c *echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	return c.String(http.StatusOK, `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="theme-color" content="#07111f">
    <title>System status | GoTix</title>
    <style>
      :root { color-scheme: dark; --ink:#f8fafc; --muted:#93a4ba; --line:rgba(255,255,255,.12); --green:#34d399; }
      * { box-sizing:border-box; }
      body { margin:0; min-height:100vh; color:var(--ink); font:16px/1.6 Inter,ui-sans-serif,system-ui,-apple-system,sans-serif; background:radial-gradient(circle at 50% 0,rgba(16,185,129,.15),transparent 30rem),#07111f; }
      .shell { width:min(920px,calc(100% - 40px)); margin:auto; }
      nav { display:flex; justify-content:space-between; align-items:center; padding:28px 0; }
      a { color:var(--muted); text-decoration:none; } a:hover { color:var(--ink); }
      .brand { color:var(--ink); font-weight:800; letter-spacing:-.03em; font-size:1.15rem; }
      .brand b { display:inline-grid; place-items:center; width:34px; height:34px; margin-right:9px; border-radius:11px; color:#07111f; background:linear-gradient(135deg,#67e8f9,#a78bfa); }
      .panel { margin-top:70px; padding:clamp(28px,6vw,58px); border:1px solid var(--line); border-radius:25px; background:rgba(15,28,47,.72); box-shadow:0 30px 90px rgba(0,0,0,.25); backdrop-filter:blur(18px); }
      .status { display:flex; align-items:center; gap:11px; color:var(--green); font-size:.78rem; font-weight:800; letter-spacing:.13em; text-transform:uppercase; }
      .pulse { width:11px; height:11px; border-radius:50%; background:var(--green); box-shadow:0 0 0 0 rgba(52,211,153,.5); animation:pulse 2s infinite; }
      @keyframes pulse { 70% { box-shadow:0 0 0 10px rgba(52,211,153,0); } 100% { box-shadow:0 0 0 0 rgba(52,211,153,0); } }
      h1 { margin:18px 0 10px; font-size:clamp(2.4rem,7vw,4.8rem); line-height:1; letter-spacing:-.07em; }
      .intro { margin:0; color:var(--muted); font-size:1.05rem; }
      .grid { display:grid; grid-template-columns:repeat(3,1fr); gap:12px; margin-top:42px; }
      .metric { padding:18px; border:1px solid var(--line); border-radius:15px; background:rgba(255,255,255,.035); }
      .metric label { display:block; color:var(--muted); font-size:.78rem; text-transform:uppercase; letter-spacing:.08em; }
      .metric strong { display:block; margin-top:5px; font-size:1.05rem; }
      footer { margin-top:22px; color:#66768b; font-size:.8rem; }
      @media (max-width:620px) { .panel { margin-top:35px; } .grid { grid-template-columns:1fr; } }
    </style>
  </head>
  <body>
    <div class="shell">
      <nav><a class="brand" href="/"><b>G</b>GoTix</a><a href="/api/v1/docs">API documentation &#8594;</a></nav>
      <main class="panel">
        <div class="status"><span class="pulse"></span> All systems operational</div>
        <h1>Everything is running smoothly.</h1>
        <p class="intro">The GoTix API is online and ready to power your next event.</p>
        <div class="grid">
          <div class="metric"><label>API service</label><strong>Operational</strong></div>
          <div class="metric"><label>Response</label><strong>&lt; 100 ms</strong></div>
          <div class="metric"><label>Checked at</label><strong id="checked">Loading...</strong></div>
        </div>
      </main>
      <footer>Automated status check &middot; <a href="/healthz">Health check endpoint</a></footer>
    </div>
    <script>
      fetch("/healthz").then(r => r.json()).then(data => {
        document.getElementById("checked").textContent = new Date(data.timestamp).toLocaleTimeString();
      });
    </script>
  </body>
</html>`)
}

func serveHealthJSON(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status":    "ok",
		"service":   "gotix-api",
		"timestamp": time.Now().UTC(),
	})
}

func serveHome(c *echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	return c.String(http.StatusOK, `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="theme-color" content="#111827">
    <title>GoTix | Event infrastructure</title>
    <style>
      :root {
        color-scheme: dark;
        --ink: #f8fafc;
        --muted: #a7b0c2;
        --line: rgba(255,255,255,.12);
        --panel: rgba(17,24,39,.72);
        --purple: #a78bfa;
        --cyan: #67e8f9;
      }
      * { box-sizing: border-box; }
      body {
        margin: 0;
        min-height: 100vh;
        color: var(--ink);
        font: 16px/1.6 Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
        background:
          radial-gradient(circle at 15% 15%, rgba(124,58,237,.32), transparent 34rem),
          radial-gradient(circle at 85% 20%, rgba(8,145,178,.22), transparent 30rem),
          #070b16;
      }
      .shell { width: min(1120px, calc(100% - 40px)); margin: auto; }
      nav { display:flex; align-items:center; justify-content:space-between; padding: 28px 0; }
      .brand { display:flex; gap:11px; align-items:center; color:var(--ink); text-decoration:none; font-weight:800; letter-spacing:-.03em; font-size:1.2rem; }
      .mark { display:grid; place-items:center; width:35px; height:35px; border-radius:11px; color:#111827; background:linear-gradient(135deg,var(--cyan),var(--purple)); box-shadow:0 0 28px rgba(103,232,249,.25); }
      .nav-link { color:var(--muted); text-decoration:none; font-size:.92rem; }
      .nav-link:hover { color:var(--ink); }
      .hero { display:grid; grid-template-columns: 1.15fr .85fr; gap:70px; align-items:center; padding: 76px 0 92px; }
      .eyebrow { display:inline-flex; align-items:center; gap:9px; color:var(--cyan); font-size:.75rem; font-weight:800; letter-spacing:.14em; text-transform:uppercase; }
      .eyebrow i { width:7px; height:7px; border-radius:50%; background:#34d399; box-shadow:0 0 14px #34d399; }
      h1 { max-width: 700px; margin:18px 0 22px; font-size:clamp(3.2rem, 8vw, 6.5rem); line-height:.94; letter-spacing:-.075em; }
      h1 span { color:transparent; background:linear-gradient(110deg,var(--cyan),var(--purple) 72%); background-clip:text; -webkit-background-clip:text; }
      .lead { max-width:570px; color:var(--muted); font-size:1.14rem; }
      .actions { display:flex; flex-wrap:wrap; gap:13px; margin-top:32px; }
      .button { display:inline-flex; align-items:center; gap:9px; padding:13px 18px; border:1px solid var(--line); border-radius:12px; color:var(--ink); text-decoration:none; font-weight:750; transition:transform .2s, border-color .2s, background .2s; }
      .button:hover { transform:translateY(-2px); border-color:rgba(167,139,250,.7); }
      .primary { color:#101526; border-color:transparent; background:linear-gradient(110deg,var(--cyan),var(--purple)); box-shadow:0 12px 30px rgba(103,232,249,.14); }
      .visual { position:relative; min-height:390px; display:grid; place-items:center; }
      .orb { position:absolute; width:280px; height:280px; border:1px solid rgba(167,139,250,.4); border-radius:50%; box-shadow:inset 0 0 60px rgba(124,58,237,.18), 0 0 80px rgba(8,145,178,.14); }
      .orb:before, .orb:after { content:""; position:absolute; inset:28px; border:1px solid rgba(103,232,249,.22); border-radius:50%; }
      .orb:after { inset:68px; border-color:rgba(255,255,255,.16); }
      .console { position:relative; width:min(100%,360px); padding:22px; border:1px solid var(--line); border-radius:22px; background:var(--panel); backdrop-filter:blur(18px); box-shadow:0 24px 70px rgba(0,0,0,.3); }
      .console-top { display:flex; gap:7px; margin-bottom:22px; }
      .dot { width:8px; height:8px; border-radius:50%; background:#fb7185; } .dot:nth-child(2){background:#fbbf24}.dot:nth-child(3){background:#34d399}
      .code { color:#c4b5fd; font: .88rem/2 ui-monospace, SFMono-Regular, Menlo, monospace; }
      .code b { color:#67e8f9; font-weight:500; } .code em { color:#94a3b8; font-style:normal; }
      .features { display:grid; grid-template-columns:repeat(3,1fr); gap:14px; padding-bottom:70px; }
      .card { padding:22px; border:1px solid var(--line); border-radius:17px; background:rgba(255,255,255,.035); }
      .icon { display:grid; place-items:center; width:36px; height:36px; margin-bottom:16px; border-radius:10px; color:#111827; background:linear-gradient(135deg,var(--cyan),var(--purple)); font-weight:900; }
      .card h2 { margin:0 0 6px; font-size:1rem; } .card p { margin:0; color:var(--muted); font-size:.9rem; }
      footer { padding:22px 0 30px; border-top:1px solid var(--line); color:#718096; font-size:.82rem; }
      @media (max-width: 760px) { .hero { grid-template-columns:1fr; gap:30px; padding-top:50px; } .visual { min-height:300px; } .features { grid-template-columns:1fr; } h1 { font-size:clamp(3.2rem, 18vw, 5rem); } }
    </style>
  </head>
  <body>
    <div class="shell">
      <nav>
        <a class="brand" href="/"><span class="mark">G</span> GoTix</a>
        <a class="nav-link" href="/health">System status <span style="color:#34d399">●</span></a>
      </nav>
      <main>
        <section class="hero">
          <div>
            <div class="eyebrow"><i></i> Event infrastructure, reimagined</div>
            <h1>Make every <span>moment</span> count.</h1>
            <p class="lead">A fast, reliable API for creating unforgettable events, managing ticket inventory, and giving every guest a seamless booking experience.</p>
            <div class="actions">
              <a class="button primary" href="/api/v1/docs">Explore the API <span>→</span></a>
              <a class="button" href="/health">View system status</a>
            </div>
          </div>
          <div class="visual" aria-hidden="true">
            <div class="orb"></div>
            <div class="console">
              <div class="console-top"><i class="dot"></i><i class="dot"></i><i class="dot"></i></div>
              <div class="code"><em>GET</em> <b>/api/v1/events</b><br><em>POST</em> <b>/api/v1/bookings</b><br><em>AUTH</em> <b>Bearer JWT</b><br><br><span style="color:#34d399">200</span> <em>ready for your next event</em></div>
            </div>
          </div>
        </section>
        <section class="features">
          <article class="card"><div class="icon">01</div><h2>Built for momentum</h2><p>Simple REST endpoints that help teams ship event experiences quickly.</p></article>
          <article class="card"><div class="icon">02</div><h2>Secure by default</h2><p>JWT authentication and role-based access for every important action.</p></article>
          <article class="card"><div class="icon">03</div><h2>Inventory aware</h2><p>Transactional ticket reservations keep availability accurate under load.</p></article>
        </section>
      </main>
      <footer>GoTix API &middot; Event management and ticket booking platform</footer>
    </div>
  </body>
</html>`)
}

func serveSwaggerUI(c *echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	return c.String(http.StatusOK, string(swaggerUI))
}

func serveOpenAPISpec(c *echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/yaml; charset=utf-8")
	return c.String(http.StatusOK, string(openAPISpec))
}
