// Package card renders the neofetch-style profile SVG (light + dark themes).
package card

import (
	"fmt"
	"html"
	"math"
	"os"
	"strings"

	"github.com/xqsit94/xqsit94/internal/profile"
)

// Data carries every value the card renders.
type Data struct {
	Name, Host                string
	Title                     string
	Bio                       []string
	Portrait                  []string
	Company, Role, Location   string
	Uptime                    string
	Institution, Degree       string
	Langs, Frameworks, AI     string
	Website, LinkedIn, GitHub string
	Commits, PullRequests     int
	Stars, Issues             int
	Contributed               int
}

// Stats is the subset of GitHub stats the card shows.
type Stats struct {
	Commits, PullRequests, Stars, Issues, Contributed int
}

const (
	padX     = 30.0
	contentY = 30.0
	pFS, pLH = 12.0, 13.4 // portrait font-size / line-height
	pCW      = 7.22       // portrait char advance
	iFS, iLH = 15.0, 23.0 // info font-size / line-height
	iCW      = 9.0        // info char advance
	bFS, bLH = 13.5, 19.0 // bio font-size / line-height
	bCW      = 8.12
	font     = "'SFMono-Regular','SF Mono',Consolas,'Liberation Mono',Menlo,monospace"
)

type theme struct {
	win, bar, border, title    string
	name, host, topkey, subkey string
	arrow, branch, value, dash string
	pg                         [3]string
	palette                    [8]string
}

var themes = map[string]theme{
	"dark": {
		win: "#0d1117", bar: "#161b22", border: "#30363d", title: "#7d8590",
		name: "#f47067", host: "#6cb6ff", topkey: "#d2a8ff", subkey: "#e3b341",
		arrow: "#56d4dd", branch: "#39929e", value: "#adbac7", dash: "#444c56",
		pg:      [3]string{"#f47067", "#d2a8ff", "#6cb6ff"},
		palette: [8]string{"#f47067", "#57ab5a", "#e3b341", "#6cb6ff", "#d2a8ff", "#56d4dd", "#adbac7", "#444c56"},
	},
	"light": {
		win: "#ffffff", bar: "#f6f8fa", border: "#d0d7de", title: "#57606a",
		name: "#cf222e", host: "#0969da", topkey: "#8250df", subkey: "#9a6700",
		arrow: "#1b7c83", branch: "#1b7c83", value: "#24292f", dash: "#d8dee4",
		pg:      [3]string{"#cf222e", "#8250df", "#0969da"},
		palette: [8]string{"#cf222e", "#1a7f37", "#9a6700", "#0969da", "#8250df", "#1b7c83", "#24292f", "#afb8c1"},
	},
}

type seg struct{ text, class string } // a colored run within a rich row
type row struct {
	kind, key, val string // kind: top | mid | end | blank | stats
	segs           []seg  // for kind == "stats"
}

func (d Data) rows() []row {
	return []row{
		{kind: "top", key: "Host", val: d.Company},
		{kind: "top", key: "OS", val: d.Title},
		{kind: "mid", key: "Role", val: d.Role},
		{kind: "mid", key: "Uptime", val: d.Uptime},
		{kind: "end", key: "Location", val: d.Location},
		{kind: "blank"},
		{kind: "top", key: "DEV", val: "Stack"},
		{kind: "mid", key: "Languages", val: d.Langs},
		{kind: "mid", key: "Frameworks", val: d.Frameworks},
		{kind: "end", key: "AI", val: d.AI},
		{kind: "blank"},
		{kind: "top", key: "STATS", val: "Past Year"},
		{kind: "stats", segs: []seg{
			{fmtNum(d.Commits), "vl"}, {" commits", "sk"},
			{"  ·  ", "br"},
			{fmtNum(d.PullRequests), "vl"}, {" PRs", "sk"},
			{"  ·  ", "br"},
			{fmtNum(d.Stars), "vl"}, {" stars", "sk"},
		}},
		{kind: "blank"},
		{kind: "top", key: "EDU", val: d.Institution},
		{kind: "end", key: "Degree", val: d.Degree},
		{kind: "blank"},
		{kind: "top", key: "WWW", val: "Links"},
		{kind: "mid", key: "Website", val: d.Website},
		{kind: "mid", key: "LinkedIn", val: d.LinkedIn},
		{kind: "end", key: "GitHub", val: d.GitHub},
	}
}

func esc(s string) string { return html.EscapeString(s) }

// fmtNum renders an int with thousands separators, e.g. 1429 -> "1,429".
func fmtNum(n int) string {
	s := fmt.Sprintf("%d", n)
	neg := ""
	if strings.HasPrefix(s, "-") {
		neg, s = "-", s[1:]
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteByte(s[i])
	}
	return neg + b.String()
}

func rowChars(r row) int {
	if r.kind == "blank" {
		return 0
	}
	if r.kind == "stats" {
		n := 3 // " └ "
		for _, s := range r.segs {
			n += len([]rune(s.text))
		}
		return n
	}
	pre := 0
	if r.kind != "top" {
		pre = 3 // " ├ "
	}
	return pre + len(r.key) + 3 + len(r.val) // key + " → " + val
}

// Render returns the SVG document for the given theme ("dark" / "light").
func Render(themeName string, d Data) string {
	t := themes[themeName]
	pl := d.Portrait

	nCols := 0
	for _, l := range pl {
		if len(l) > nCols {
			nCols = len(l)
		}
	}
	pW := float64(nCols) * pCW
	infoX := padX + pW + 46

	// panel width = widest of rows / header / bio / dash
	maxPx := float64(28) * iCW // dash line
	if h := float64(len(d.Name)+1+len(d.Host)) * iCW; h > maxPx {
		maxPx = h
	}
	for _, b := range d.Bio {
		if w := float64(len([]rune(b))) * bCW; w > maxPx {
			maxPx = w
		}
	}
	for _, r := range d.rows() {
		if w := float64(rowChars(r)) * iCW; w > maxPx {
			maxPx = w
		}
	}
	infoW := math.Ceil(maxPx) + 12
	W := infoX + infoW + padX

	// promptBase is the prompt-line baseline; bodyTop drops the logo+panel below
	// it to leave room for the prompt line above them.
	promptBase := contentY + iFS
	bodyTop := contentY + iLH + 14

	// Pre-measure the panel's full height so the portrait can be centered against
	// it; each step below mirrors a drawn element in the render pass.
	y := bodyTop + iFS // header baseline
	y += iLH           // dash
	y += bLH * 0.4     // gap before bio
	for range d.Bio {
		y += bLH
	}
	y += iLH * 0.5 // gap after bio
	for _, r := range d.rows() {
		y += iLH
		_ = r
	}
	panelBottom := y
	paletteY := panelBottom + 14
	panelH := paletteY + 17 - bodyTop

	portraitH := float64(len(pl)) * pLH
	pTop := bodyTop
	if panelH > portraitH {
		pTop = bodyTop + (panelH-portraitH)/2
	}
	portraitBottom := pTop + portraitH

	H := math.Max(portraitBottom, paletteY+17) + 22

	var o strings.Builder
	p := func(s string) { o.WriteString(s); o.WriteByte('\n') }

	p(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="%s" role="img" aria-label="%s@%s — neofetch profile card">`,
		int(W), int(H), int(W), int(H), font, esc(d.Name), esc(d.Host)))
	p(`<style>` +
		`@keyframes blink{50%{opacity:0}}.cur{animation:blink 1.1s step-end infinite}` +
		`@keyframes fin{from{opacity:0}to{opacity:1}}.fade{animation:fin .8s ease both}` +
		fmt.Sprintf(`.tk{fill:%s;font-weight:700}.sk{fill:%s}`, t.topkey, t.subkey) +
		fmt.Sprintf(`.ar{fill:%s}.br{fill:%s}.vl{fill:%s}`, t.arrow, t.branch, t.value) +
		`text,tspan{white-space:pre}` +
		`</style>`)
	p(fmt.Sprintf(`<defs><linearGradient id="pg" x1="0" y1="0" x2="0" y2="1">`+
		`<stop offset="0" stop-color="%s"/><stop offset="0.5" stop-color="%s"/>`+
		`<stop offset="1" stop-color="%s"/></linearGradient></defs>`, t.pg[0], t.pg[1], t.pg[2]))
	p(fmt.Sprintf(`<rect x="0.5" y="0.5" width="%.0f" height="%.0f" rx="12" fill="%s" stroke="%s"/>`,
		W-1, H-1, t.win, t.border))

	// starship-style prompt line: ~ ❯ neofetch
	p(fmt.Sprintf(`<text x="%.0f" y="%.1f" font-size="%g">`+
		`<tspan fill="%s" font-weight="700">~ </tspan>`+
		`<tspan fill="%s" font-weight="700">❯ </tspan>`+
		`<tspan fill="%s">neofetch</tspan></text>`,
		padX, promptBase, iFS, t.arrow, t.palette[1], t.value))

	p(fmt.Sprintf(`<g class="fade" fill="url(#pg)" font-size="%g">`, pFS))
	py := pTop + pFS
	for _, l := range pl {
		if strings.TrimSpace(l) != "" {
			p(fmt.Sprintf(`<text x="%.0f" y="%.1f">%s</text>`, padX, py, esc(l)))
		}
		py += pLH
	}
	p(`</g>`)

	p(fmt.Sprintf(`<g font-size="%g">`, iFS))
	yy := bodyTop + iFS
	p(fmt.Sprintf(`<text x="%.0f" y="%.1f" font-weight="700">`+
		`<tspan fill="%s">%s</tspan><tspan fill="%s">@</tspan><tspan fill="%s">%s</tspan></text>`,
		infoX, yy, t.name, esc(d.Name), t.title, t.host, esc(d.Host)))
	curX := infoX + float64(len(d.Name)+1+len(d.Host))*iCW + 5
	p(fmt.Sprintf(`<rect class="cur" x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>`,
		curX, yy-iFS+3, iFS*0.55, iFS-1, t.name))
	yy += iLH
	p(fmt.Sprintf(`<text x="%.0f" y="%.1f" fill="%s">%s</text>`, infoX, yy, t.dash, strings.Repeat("─", 28)))
	yy += bLH * 0.4
	for _, b := range d.Bio {
		yy += bLH
		p(fmt.Sprintf(`<text x="%.0f" y="%.1f" font-size="%g" fill="%s">%s</text>`, infoX, yy, bFS, t.value, esc(b)))
	}
	yy += iLH * 0.5
	for _, r := range d.rows() {
		yy += iLH
		if r.kind == "blank" {
			continue
		}
		if r.kind == "stats" {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf(`<text x="%.0f" y="%.1f"><tspan class="br"> └ </tspan>`, infoX, yy))
			for _, s := range r.segs {
				sb.WriteString(fmt.Sprintf(`<tspan class="%s">%s</tspan>`, s.class, esc(s.text)))
			}
			sb.WriteString(`</text>`)
			p(sb.String())
			continue
		}
		if r.kind == "top" {
			p(fmt.Sprintf(`<text x="%.0f" y="%.1f"><tspan class="tk">%s</tspan>`+
				`<tspan class="ar"> → </tspan><tspan class="vl">%s</tspan></text>`,
				infoX, yy, esc(r.key), esc(r.val)))
		} else {
			br := " ├ "
			if r.kind == "end" {
				br = " └ "
			}
			p(fmt.Sprintf(`<text x="%.0f" y="%.1f"><tspan class="br">%s</tspan>`+
				`<tspan class="sk">%s</tspan><tspan class="ar"> → </tspan>`+
				`<tspan class="vl">%s</tspan></text>`, infoX, yy, br, esc(r.key), esc(r.val)))
		}
	}
	p(`</g>`)

	const sq, gap = 17.0, 7.0
	for i, c := range t.palette {
		p(fmt.Sprintf(`<rect x="%.0f" y="%.1f" width="%.0f" height="%.0f" rx="3.5" fill="%s"/>`,
			infoX+float64(i)*(sq+gap), paletteY, sq, sq, c))
	}

	p(`</svg>`)
	return o.String()
}

// Assemble builds card Data from the profile, deriving the current role and
// uptime from experience and taking the top entry from education.
func Assemble(p profile.Profile, exps []profile.Experience, edu []profile.Education, totalExp string, portrait []string, s Stats) Data {
	d := Data{
		Name: p.Name, Host: p.Host,
		Title:    p.Title,
		Bio:      p.Bio,
		Portrait: portrait,
		Langs:    p.Stack.Languages, Frameworks: p.Stack.Frameworks, AI: p.Stack.AI,
		Website: p.Links.Website.Label, LinkedIn: p.Links.LinkedIn.Label, GitHub: p.Links.GitHub.Label,
		Commits: s.Commits, PullRequests: s.PullRequests,
		Stars: s.Stars, Issues: s.Issues, Contributed: s.Contributed,
	}
	if len(exps) > 0 {
		cur := exps[0]
		for _, e := range exps {
			if e.IsCurrent {
				cur = e
				break
			}
		}
		d.Company, d.Role, d.Location = cur.Company, cur.Role, cur.Location

		earliest := 0
		for _, e := range exps {
			if y := e.StartDate.Year(); y > 1 && (earliest == 0 || y < earliest) {
				earliest = y
			}
		}
		if earliest > 0 {
			d.Uptime = fmt.Sprintf("~%s · since %d", totalExp, earliest)
		} else {
			d.Uptime = "~" + totalExp
		}
	}
	if len(edu) > 0 {
		d.Institution = edu[0].Institution
		d.Degree = edu[0].Degree
		if edu[0].ShortDegree != "" {
			d.Degree = edu[0].ShortDegree
		}
	}
	return d
}

// WriteAll renders both themes into dir/{dark,light}_mode.svg.
func WriteAll(dir string, d Data) error {
	for _, name := range []string{"dark", "light"} {
		path := fmt.Sprintf("%s/%s_mode.svg", dir, name)
		if err := os.WriteFile(path, []byte(Render(name, d)+"\n"), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}
	return nil
}
