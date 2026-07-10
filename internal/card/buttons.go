package card

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type btnTheme struct{ bg, border, text string }

var btnThemes = map[string]btnTheme{
	"dark":  {bg: "#21262d", border: "#3d444d", text: "#c9d1d9"},
	"light": {bg: "#f6f8fa", border: "#d1d9e0", text: "#24292f"},
}

const btnFont = "-apple-system,BlinkMacSystemFont,'Segoe UI',Helvetica,Arial,sans-serif"

const (
	btnH    = 32.0
	btnFS   = 14.0
	btnPadL = 12.0
	btnIC   = 16.0
	btnGap  = 8.0
	btnPadR = 12.0
	btnRX   = 6.0
)

// glyphAdv is per-glyph advance width in 1000-em units (Helvetica AFM, a close
// proxy for the system sans stack) — used to size and center each button by its
// label's real proportional width. Unmapped runes use glyphAdvDefault.
const glyphAdvDefault = 556.0

var glyphAdv = map[rune]float64{
	' ': 278, '!': 278, '"': 355, '#': 556, '$': 556, '%': 889, '&': 667, '\'': 191,
	'(': 333, ')': 333, '*': 389, '+': 584, ',': 278, '-': 333, '.': 278, '/': 278,
	'0': 556, '1': 556, '2': 556, '3': 556, '4': 556, '5': 556, '6': 556, '7': 556,
	'8': 556, '9': 556, ':': 278, ';': 278, '<': 584, '=': 584, '>': 584, '?': 556, '@': 1015,
	'A': 667, 'B': 667, 'C': 722, 'D': 722, 'E': 667, 'F': 611, 'G': 778, 'H': 722,
	'I': 278, 'J': 500, 'K': 667, 'L': 556, 'M': 833, 'N': 722, 'O': 778, 'P': 667,
	'Q': 778, 'R': 722, 'S': 667, 'T': 611, 'U': 722, 'V': 667, 'W': 944, 'X': 667,
	'Y': 667, 'Z': 611, '[': 278, '\\': 278, ']': 278, '^': 469, '_': 556, '`': 333,
	'a': 556, 'b': 556, 'c': 500, 'd': 556, 'e': 556, 'f': 278, 'g': 556, 'h': 556,
	'i': 222, 'j': 222, 'k': 500, 'l': 222, 'm': 833, 'n': 556, 'o': 556, 'p': 556,
	'q': 556, 'r': 333, 's': 500, 't': 278, 'u': 556, 'v': 500, 'w': 722, 'x': 500,
	'y': 500, 'z': 500, '{': 334, '|': 260, '}': 334, '~': 584,
}

func labelWidth(label string) float64 {
	var units float64
	for _, r := range label {
		if a, ok := glyphAdv[r]; ok {
			units += a
		} else {
			units += glyphAdvDefault
		}
	}
	return units / 1000 * btnFS
}

type Button struct{ Label, Icon, File string }

func siteMark(x, y, size float64, color string) string {
	s := size / 1000
	return fmt.Sprintf(`<g transform="translate(%.2f %.2f) scale(%.4f)" fill="%s">`+
		`<g transform="matrix(0.110282,0,0,-0.112389,-24.2299,923.254)">`+
		`<path d="%s"/></g></g>`, x, y, s, color, siteMarkPath)
}

func lucideLinkedIn(x, y, size float64, color string) string {
	s := size / 24
	return fmt.Sprintf(`<g transform="translate(%.2f %.2f) scale(%.4f)" fill="none" stroke="%s" `+
		`stroke-width="2" stroke-linecap="round" stroke-linejoin="round">`+
		`<path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"/>`+
		`<rect width="4" height="12" x="2" y="9"/><circle cx="4" cy="4" r="2"/></g>`, x, y, s, color)
}

func lucideMail(x, y, size float64, color string) string {
	s := size / 24
	return fmt.Sprintf(`<g transform="translate(%.2f %.2f) scale(%.4f)" fill="none" stroke="%s" `+
		`stroke-width="2" stroke-linecap="round" stroke-linejoin="round">`+
		`<path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"/>`+
		`<rect x="2" y="4" width="20" height="16" rx="2"/></g>`, x, y, s, color)
}

func btnNaturalWidth(label string) float64 {
	return math.Ceil(btnPadL + btnIC + btnGap + labelWidth(label) + btnPadR)
}

func renderButton(themeName, label, icon string, width float64) string {
	t := btnThemes[themeName]
	textW := labelWidth(label)
	block := btnIC + btnGap + textW
	iconX := (width - block) / 2
	iconY := (btnH - btnIC) / 2
	textX := iconX + btnIC + btnGap
	textY := btnH/2 + btnFS*0.34

	var mark string
	switch icon {
	case "website":
		mark = siteMark(iconX, iconY, btnIC, t.text)
	case "linkedin":
		mark = lucideLinkedIn(iconX, iconY, btnIC, t.text)
	case "email":
		mark = lucideMail(iconX, iconY, btnIC, t.text)
	}

	var o strings.Builder
	o.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" `+
		`viewBox="0 0 %.0f %.0f" font-family="%s" role="img" aria-label="%s">`,
		width, btnH, width, btnH, btnFont, esc(label)))
	o.WriteString(fmt.Sprintf(`<rect x="0.5" y="0.5" width="%.0f" height="%.0f" rx="%g" fill="%s" stroke="%s"/>`,
		width-1, btnH-1, btnRX, t.bg, t.border))
	o.WriteString(mark)
	o.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%g" font-weight="500" fill="%s">%s</text>`,
		textX, textY, btnFS, t.text, esc(label)))
	o.WriteString(`</svg>`)
	return o.String()
}

// WriteButtons renders the connect buttons for both themes into
// dir/btn_<file>_<theme>.svg.
func WriteButtons(dir, websiteLabel string) error {
	btns := []Button{
		{Label: websiteLabel, Icon: "website", File: "website"},
		{Label: "LinkedIn", Icon: "linkedin", File: "linkedin"},
		{Label: "Email", Icon: "email", File: "email"},
	}
	var width float64
	for _, b := range btns {
		if w := btnNaturalWidth(b.Label); w > width {
			width = w
		}
	}
	for _, b := range btns {
		for _, themeName := range []string{"dark", "light"} {
			path := fmt.Sprintf("%s/btn_%s_%s.svg", dir, b.File, themeName)
			if err := os.WriteFile(path, []byte(renderButton(themeName, b.Label, b.Icon, width)+"\n"), 0o644); err != nil {
				return fmt.Errorf("write %s: %w", path, err)
			}
		}
	}
	return nil
}

// siteMarkPath is the xqsit.dev monogram (viewBox 0 0 1000 1000), drawn
// monochrome in the button's text color via the matrix in siteMark.
const siteMarkPath = `M5665,7028C5382,7014 5100,6948 4791,6823C4637,6761 4552,6717 4395,6617C4020,6378 3747,6106 3473,5698C3404,5595 3343,5513 3339,5516C3335,5518 3278,5596 3213,5688C3148,5780 3028,5950 2945,6066C2809,6258 2715,6392 2570,6598L2519,6670L1510,6670C954,6670 500,6667 500,6664C500,6661 521,6629 547,6592C573,6555 639,6460 695,6380C872,6125 892,6097 1080,5829C1311,5501 1595,5096 1654,5010C1679,4974 1726,4907 1759,4860C1813,4784 1925,4624 2032,4470C2053,4440 2099,4375 2134,4325C2218,4208 2350,4013 2350,4007C2350,4004 2262,3872 2153,3713C2045,3555 1931,3387 1900,3340C1869,3293 1792,3179 1729,3085C1665,2992 1467,2697 1288,2430C969,1954 870,1806 623,1443C556,1343 500,1259 500,1256C500,1253 950,1250 1499,1250L2498,1250L2533,1298C2565,1341 3265,2391 3374,2558C3399,2597 3424,2630 3428,2630C3432,2629 3451,2606 3470,2578C3537,2481 3624,2384 3795,2215C4063,1949 4255,1809 4550,1662C5075,1399 5676,1311 6250,1412C6313,1423 6383,1434 6405,1437L6445,1442L6506,1343C6540,1289 6607,1181 6656,1102C6841,804 6979,672 7195,585C7398,504 7442,500 8289,500C8905,500 9011,502 9007,514C9004,522 8963,588 8917,662C8871,735 8751,928 8650,1090C8550,1252 8419,1462 8361,1555C8302,1649 8224,1772 8188,1830C8152,1888 8087,1992 8042,2061C7998,2131 7950,2211 7934,2238L7905,2289L7925,2314C7936,2327 7981,2383 8026,2437C8516,3031 8740,3808 8644,4574C8578,5100 8375,5573 8033,6000C7864,6212 7652,6401 7395,6572C7063,6790 6727,6923 6305,7002C6177,7026 5868,7039 5665,7028ZM6235,6510C7031,6372 7695,5836 7990,5092C8112,4783 8168,4461 8157,4125C8142,3668 7989,3236 7700,2833C7632,2738 7610,2715 7610,2739C7610,2754 7486,2950 7314,3209L7231,3332L7254,3369C7429,3645 7509,4033 7470,4410C7447,4636 7348,4927 7236,5100C6989,5482 6614,5748 6190,5843C6109,5861 6062,5864 5865,5864C5659,5865 5624,5863 5530,5842C5231,5775 4982,5645 4777,5451C4690,5368 4583,5247 4525,5164C4398,4985 3958,4333 3590,3780C3510,3660 3348,3419 3230,3243C3112,3068 2844,2666 2633,2350C2423,2034 2245,1767 2238,1757C2226,1741 2199,1740 1832,1742L1439,1745L1528,1880C1577,1954 1656,2074 1703,2145C1750,2217 1985,2568 2225,2925C2465,3283 2763,3726 2886,3910C3009,4094 3122,4263 3137,4285C3152,4307 3226,4418 3302,4532C3377,4646 3453,4759 3470,4782C3486,4806 3593,4964 3706,5134C3998,5573 4020,5603 4158,5763C4483,6139 4927,6406 5376,6495C5555,6531 5661,6539 5890,6535C6059,6532 6145,6526 6235,6510ZM2318,6088C2379,6005 2510,5820 2768,5455C2815,5389 2895,5276 2947,5204C2998,5132 3040,5067 3040,5061C3040,5054 2968,4942 2879,4812C2791,4682 2702,4551 2683,4522C2663,4493 2642,4471 2636,4473C2630,4475 2571,4557 2505,4656C2439,4755 2336,4905 2276,4990C2216,5075 2094,5249 2006,5375C1918,5502 1817,5647 1781,5698C1610,5939 1460,6160 1460,6171C1460,6176 1618,6180 1855,6180L2250,6180L2318,6088ZM6050,5368C6196,5340 6367,5264 6513,5159C6745,4994 6911,4721 6956,4431C6975,4309 6974,4108 6955,3999C6919,3790 6843,3630 6693,3448C6658,3405 6630,3368 6630,3364C6630,3356 6745,3171 6825,3050C6859,2998 6916,2908 6950,2850C6985,2793 7056,2678 7108,2595C7161,2513 7264,2348 7338,2230C7412,2112 7482,1999 7495,1980C7548,1898 7608,1803 7785,1520C7828,1451 7919,1306 7987,1198C8054,1090 8110,999 8110,996C8110,986 7603,990 7525,1000C7422,1015 7314,1065 7251,1129C7210,1170 7181,1213 7053,1420C6956,1577 6885,1690 6673,2025C6649,2064 6566,2194 6490,2315C6294,2625 6294,2624 6201,2717C6062,2854 5942,2923 5743,2981C5684,2998 5588,3025 5530,3042C5340,3096 5161,3196 4991,3343C4845,3470 4668,3682 4575,3843C4560,3869 4545,3890 4541,3890C4537,3890 4512,3853 4485,3808C4458,3762 4399,3669 4354,3600C4238,3421 4244,3440 4287,3384C4506,3092 4809,2827 5069,2700C5199,2636 5251,2617 5420,2570C5500,2547 5592,2519 5624,2507C5721,2473 5819,2402 5886,2319C5948,2242 6060,2085 6060,2075C6060,2072 6068,2057 6079,2042C6134,1963 6172,1898 6166,1893C6153,1879 6055,1870 5865,1864C5632,1856 5481,1871 5281,1920C4917,2009 4663,2133 4355,2372C4251,2453 3953,2752 3870,2858C3788,2963 3730,3050 3730,3069C3730,3083 3753,3118 3946,3409C4185,3768 4326,3980 4567,4341C4998,4990 5052,5059 5225,5182C5325,5253 5504,5333 5625,5361C5749,5389 5925,5392 6050,5368Z`
