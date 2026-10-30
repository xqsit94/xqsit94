package card

import (
	_ "embed"
	"fmt"
	"strings"
)

// The code texture is this package's own renderer source, so the portrait is
// drawn in the code that draws it.
//
//go:embed card.go
var codeSrc string

const codeBuckets = 8 // lit cells are split into this many scramble groups

func flatCode() []rune {
	var b strings.Builder
	space := false
	for _, r := range codeSrc {
		switch {
		case r == ' ' || r == '\t' || r == '\n':
			if !space {
				b.WriteByte(' ')
			}
			space = true
		case r > ' ' && r < 0x7f:
			b.WriteRune(r)
			space = false
		}
	}
	return []rune(b.String())
}

// Each bucket loops every codeLoop seconds: scramble for codeFrames*codeStep,
// resolve into the real code, hold, scramble out, then show only dim code
// until the loop restarts. Bucket delays are staggered across that dark gap,
// so the face as a whole never vanishes. The scramble glyphs flip on one
// global codeStep clock; a bucket only reveals them during its windows.
const (
	codeLoop   = 10.0
	codeFrames = 10
	codeStep   = 0.055
	codeIn     = 0.1            // scramble-in starts
	codeOut    = codeLoop - 1.8 // scramble-out starts
	codeGlyphs = `01<>{}[]=+*#%$/\|`
	codeScroll = 60  // extra dim rows the background steps up through
	codeRowDur = 0.9 // seconds per background row
)

// keyframes renders a step-end opacity track that is on only inside windows
// (seconds within codeLoop).
func keyframes(name string, windows [][2]float64) string {
	pct := func(t float64) float64 { return t / codeLoop * 100 }
	s := "@keyframes " + name + "{0%{opacity:0}"
	for _, w := range windows {
		s += fmt.Sprintf(`%.2f%%{opacity:1}%.2f%%{opacity:0}`, pct(w[0]), pct(w[1]))
	}
	return s + "100%{opacity:0}}"
}

func codeStyle() string {
	spin := codeFrames * codeStep
	s := keyframes("cr", [][2]float64{{codeIn + spin, codeOut}}) +
		keyframes("cw", [][2]float64{{codeIn, codeIn + spin}, {codeOut, codeOut + spin}}) +
		fmt.Sprintf(`.cr{animation:cr %gs step-end infinite backwards}.cw{animation:cw %gs step-end infinite backwards}`, codeLoop, codeLoop)
	s += fmt.Sprintf(`@keyframes cg{0%%{opacity:1}%.2f%%,100%%{opacity:0}}`, 100.0/codeFrames)
	// negative delays put frame j in phase from t=0 without a backwards-fill flash
	for j := range codeFrames {
		s += fmt.Sprintf(`.cg%d{animation:cg %.3fs step-end %.3fs infinite}`, j, spin, float64(j)*codeStep-spin)
	}
	// spread bucket starts across the dark gap (scramble-out end to the next
	// scramble-in) so some bucket is always lit
	gap := codeLoop - (codeOut + spin) + codeIn
	for k := range codeBuckets {
		s += fmt.Sprintf(`.ck%d>g{animation-delay:%.2fs}`, k, 0.2+float64(k)*gap/codeBuckets)
	}
	s += fmt.Sprintf(`@keyframes csc{to{transform:translateY(-%.1fpx)}}.csc{animation:csc %gs steps(%d) infinite}`,
		codeScroll*pLH, codeScroll*codeRowDur, codeScroll)
	return s + `@media (prefers-reduced-motion:reduce){.cr,.csc{animation:none}.csg{display:none}}`
}

// codeDefs holds the holo gradient and glow, plus a mask that clips the
// scrolling dim code to the portrait box and fades its top and bottom rows.
func codeDefs(t theme, cols, rows int, top float64) string {
	x0, y0, w, h := padX, top, float64(cols)*pCW, float64(rows)*pLH
	bx := box(cols, rows, top)
	return fmt.Sprintf(`<linearGradient id="holo" gradientUnits="userSpaceOnUse" x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f">`+
		`<stop offset="0" stop-color="%s"/><stop offset="0.5" stop-color="%s"/><stop offset="1" stop-color="%s"/></linearGradient>`+
		`<filter id="glow" x="-10%%" y="-10%%" width="120%%" height="120%%"><feGaussianBlur stdDeviation="%g" result="b"/>`+
		`<feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge></filter>`+
		`<linearGradient id="cfade" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#000"/>`+
		`<stop offset="0.07" stop-color="#fff"/><stop offset="0.93" stop-color="#fff"/><stop offset="1" stop-color="#000"/></linearGradient>`+
		`<mask id="cmask" maskUnits="userSpaceOnUse" %s><rect %s fill="url(#cfade)"/></mask>`,
		x0, y0, x0+w, y0+h*0.45, t.holo[0], t.holo[1], t.holo[2], t.glow, bx, bx)
}

const inkTiers = 3

// inkTier maps the portrait's density ramp to a brightness level 1..inkTiers,
// so the lit code keeps the shading of the ASCII art; 0 means blank.
func inkTier(r rune) int {
	switch {
	case r == ' ':
		return 0
	case strings.ContainsRune(".`'", r):
		return 1
	case strings.ContainsRune(":;!", r):
		return 2
	}
	return 3
}

type grid [][]rune

func newGrid(rows, cols int) grid {
	g := make(grid, rows)
	for y := range g {
		g[y] = []rune(strings.Repeat(" ", cols))
	}
	return g
}

// codePortrait draws the portrait box as scrolling dim code. Inked cells show
// a bold copy of the same scrolling code, so the face stays put while its
// characters flow through it: a fixed holo rect is masked by the bold text,
// which is masked by the cell shapes (gray = ink tier, animated per bucket).
// The scramble is codeFrames grids of random glyphs, masked by the same cell
// shapes during each bucket's scramble windows.
func codePortrait(t theme, pl []string, nCols int, top float64) string {
	src := flatCode()
	// rows repeat every codeScroll, so the scroll's jump back to row 0 lands on
	// identical text
	dim := newGrid(len(pl)+codeScroll, nCols)
	for y := range dim {
		for x := range nCols {
			dim[y][x] = src[((y%codeScroll)*nCols+x)%len(src)]
		}
	}
	var cells [codeBuckets][inkTiers]strings.Builder // cell paths, one per bucket and tier
	var scr [codeFrames]grid
	for j := range scr {
		scr[j] = newGrid(len(pl), nCols)
	}
	for y, l := range pl {
		lr := []rune(l)
		for x := range min(nCols, len(lr)) {
			tier := inkTier(lr[x])
			if tier == 0 {
				continue
			}
			k := int(cellHash(x, y, 0) % codeBuckets)
			// cell rect spans from ~cap height above the row's baseline to just below it
			fmt.Fprintf(&cells[k][tier-1], "M%.1f %.1fh%gv%gh-%gz",
				padX+float64(x)*pCW, top+pFS+float64(y)*pLH-pLH*0.78, pCW, pLH, pCW)
			for j := range scr {
				scr[j][y][x] = rune(codeGlyphs[cellHash(x, y, j+1)%uint32(len(codeGlyphs))])
			}
		}
	}

	var o strings.Builder
	// one <text> per non-empty row, trimmed and shifted so blank margins cost nothing
	rows := func(g grid) {
		for y, r := range g {
			s := string(r)
			lead := len(s) - len(strings.TrimLeft(s, " "))
			if s = strings.TrimSpace(s); s != "" {
				fmt.Fprintf(&o, `<text x="%.1f" y="%.1f">%s</text>`+"\n", padX+float64(lead)*pCW, top+pFS+float64(y)*pLH, esc(s))
			}
		}
	}
	bx := box(nCols, len(pl), top)
	// cellMask emits one group per bucket, animated by class cls, of that
	// bucket's cell paths; tiered sets each tier's ink level as its gray.
	cellMask := func(id, cls string, tiered bool) {
		fmt.Fprintf(&o, `<mask id="%s" maskUnits="userSpaceOnUse" %s>`+"\n", id, bx)
		for k := range codeBuckets {
			fmt.Fprintf(&o, `<g class="ck%d"><g class="%s" fill="#fff">`, k, cls)
			for i := range inkTiers {
				if cells[k][i].Len() == 0 {
					continue
				}
				op := 1.0
				if tiered {
					op = 0.05 + 0.95*float64(i+1)/inkTiers
				}
				fmt.Fprintf(&o, `<use href="#cp%d_%d" fill-opacity="%.2f"/>`, k, i, op)
			}
			o.WriteString("</g></g>\n")
		}
		o.WriteString("</mask>\n")
	}

	o.WriteString("<defs>")
	for k := range codeBuckets {
		for i := range inkTiers {
			if d := cells[k][i].String(); d != "" {
				fmt.Fprintf(&o, `<path id="cp%d_%d" d="%s"/>`, k, i, d)
			}
		}
	}
	o.WriteString("</defs>\n")
	cellMask("ccell", "cr", true)
	cellMask("cscr", "cw", false)
	fmt.Fprintf(&o, `<g font-size="%g">`+"\n", pFS)
	fmt.Fprintf(&o, `<mask id="clit" maskUnits="userSpaceOnUse" %s><g mask="url(#ccell)"><g class="csc" fill="#fff" font-weight="700">`+"\n", bx)
	rows(dim)
	o.WriteString("</g></g></mask>\n")

	// the dim code is a faint copy of the holo rect, seen through the scrolling text
	fmt.Fprintf(&o, `<mask id="cdim" maskUnits="userSpaceOnUse" %s><g mask="url(#cmask)"><g class="csc" fill="#fff">`+"\n", bx)
	rows(dim)
	o.WriteString("</g></g></mask>\n")
	fmt.Fprintf(&o, `<rect %s fill="url(#holo)" fill-opacity="%g" mask="url(#cdim)"/>`+"\n", bx, t.dimA)
	filter := ""
	if t.glow > 0 {
		filter = ` filter="url(#glow)"`
	}
	fmt.Fprintf(&o, `<g%s>`+"\n"+`<rect %s fill="url(#holo)" mask="url(#clit)"/>`+"\n", filter, bx)
	// the window-colored stroke under each scramble glyph hides the dim glyph beneath it
	fmt.Fprintf(&o, `<g class="csg" mask="url(#cscr)" fill="url(#holo)" fill-opacity=".8" font-weight="700" stroke="%s" stroke-width="2.5" paint-order="stroke">`+"\n", t.win)
	for j, g := range scr {
		fmt.Fprintf(&o, `<g class="cg%d">`+"\n", j)
		rows(g)
		o.WriteString("</g>\n")
	}
	o.WriteString("</g>\n</g>\n</g>\n")
	return o.String()
}

// cellHash mixes all bits before any modulo, so buckets and glyphs scatter
// instead of tiling in a lattice (a plain x*a^y*b mod 8 keeps only low bits).
func cellHash(x, y, salt int) uint32 {
	h := uint32(x)*374761393 + uint32(y)*668265263 + uint32(salt)*2246822519
	h = (h ^ h>>13) * 1274126177
	return h ^ h>>16
}

// box is the portrait area as SVG x/y/width/height attributes; the +4 keeps
// the last row's descenders inside it.
func box(cols, rows int, top float64) string {
	return fmt.Sprintf(`x="%.0f" y="%.0f" width="%.0f" height="%.0f"`, padX, top, float64(cols)*pCW, float64(rows)*pLH+4)
}
