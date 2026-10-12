"""Assemble the horizontal and stacked Ports lockups from the symbol + wordmark.

The wordmark is outlined Geist SemiBold (see outline_wordmark.py); this
composes the two parts on a fixed grid so relative sizes and spacing are
locked. Height 256 per the skill's lockup convention.
"""
from fontTools.ttLib import TTFont
from fontTools.pens.svgPathPen import SVGPathPen

FONT = "Geist-600.ttf"
TEXT = "Ports"

font = TTFont(FONT)
gs = font.getGlyphSet()
cmap = font.getBestCmap()
upm = font["head"].unitsPerEm
cap = font["OS/2"].sCapHeight

# Word layout in font units.
x = 0
glyphs = []
for ch in TEXT:
    g = gs[cmap[ord(ch)]]
    pen = SVGPathPen(gs)
    g.draw(pen)
    d = pen.getCommands()
    if d:
        glyphs.append((x, d))
    x += g.width
word_w_font = x  # 2656

# --- Geometry (SVG units) ---
CAP = 64.0            # wordmark cap height
SYM = 128.0           # symbol outer diameter = 2x cap
cap_scale = CAP / cap
word_w = word_w_font * cap_scale
gap = 44.0            # symbol -> wordmark
pad = 40.0

def word_group(tx, baseline_y):
    """Glyph paths scaled to CAP, baseline at baseline_y, left edge at tx."""
    out = []
    for x0, d in glyphs:
        out.append(
            f'<path transform="translate({tx:.2f},{baseline_y:.2f}) '
            f'scale({cap_scale:.6f}) translate({x0},0) scale(1,-1)" d="{d}"/>'
        )
    return out

def symbol_group(tx, ty, diameter):
    """The 256-canvas symbol scaled to `diameter`, top-left at (tx,ty)."""
    s = diameter / 256.0
    return (
        f'<g transform="translate({tx:.2f},{ty:.2f}) scale({s:.6f})">'
        '<path fill-rule="evenodd" d="M210.72 113.41 A84 84 0 1 1 142.59 45.28 '
        'L135.64 84.67 A44 44 0 1 0 171.33 120.36 Z"/>'
        '<circle cx="173.25" cy="82.75" r="20"/></g>'
    )

INK = "#0A0A0A"

# ---- Horizontal lockup ----
h_h = 256.0
sym_top = (h_h - SYM) / 2
baseline = h_h / 2 + CAP / 2
word_x = pad + SYM + gap
h_w = word_x + word_w + pad
hg = [symbol_group(pad, sym_top, SYM)]
hg += word_group(word_x, baseline)

with open("ports-horizontal.svg", "w") as f:
    f.write(
        f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {h_w:.0f} {h_h:.0f}" '
        f'width="{h_w:.0f}" height="{h_h:.0f}" role="img" aria-labelledby="t">\n'
        f'  <title id="t">Ports</title>\n'
        f'  <g fill="{INK}">\n    ' + "\n    ".join(hg) + "\n  </g>\n</svg>\n"
    )

# ---- Stacked lockup ----
s_h = pad + SYM + 40 + CAP + pad
s_w = max(SYM, word_w) + pad * 2
sym_x = (s_w - SYM) / 2
stack_base = pad + SYM + 40 + CAP
sg = [symbol_group(sym_x, pad, SYM)]
sg += word_group((s_w - word_w) / 2, stack_base)

with open("ports-stacked.svg", "w") as f:
    f.write(
        f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {s_w:.0f} {s_h:.0f}" '
        f'width="{s_w:.0f}" height="{s_h:.0f}" role="img" aria-labelledby="t">\n'
        f'  <title id="t">Ports</title>\n'
        f'  <g fill="{INK}">\n    ' + "\n    ".join(sg) + "\n  </g>\n</svg>\n"
    )

print(f"horizontal: {h_w:.0f}x{h_h:.0f}")
print(f"stacked:    {s_w:.0f}x{s_h:.0f}")
