"""Outline "Ports" from Geist 600 into SVG paths, for the wordmark.

The logo master must not depend on an installed font, so the wordmark is
converted to outlines here. Run once; commit the result.
"""
import sys
from fontTools.ttLib import TTFont
from fontTools.pens.svgPathPen import SVGPathPen
from fontTools.pens.transformPen import TransformPen
from fontTools.misc.transform import Transform

FONT = "Geist-600.ttf"
TEXT = "Ports"
TRACKING = 0.0  # em; tighten slightly for a display wordmark
SIZE = 1000  # em units per 1000 -> we scale the whole word to a target cap

font = TTFont(FONT)
glyph_set = font.getGlyphSet()
cmap = font.getBestCmap()
upm = font["head"].unitsPerEm
cap = font["OS/2"].sCapHeight if hasattr(font["OS/2"], "sCapHeight") else int(upm * 0.7)

# Lay out the word, glyph by glyph, in font units.
x = 0
advance_tracking = int(TRACKING * upm)
paths = []
for ch in TEXT:
    gname = cmap[ord(ch)]
    glyph = glyph_set[gname]
    pen = SVGPathPen(glyph_set)
    glyph.draw(pen)
    d = pen.getCommands()
    if d:
        paths.append((x, d))
    x += glyph.width + advance_tracking

total_width = x - advance_tracking

# Flip the y axis (font space is y-up, SVG is y-down) and drop the baseline.
# Emit at font scale; the caller positions it.
print(f"<!-- {TEXT} in Geist SemiBold, outlined. upm={upm} cap={cap} width={total_width} -->")
for x0, d in paths:
    print(f'<path transform="translate({x0},0) scale(1,-1)" d="{d}"/>')
print(f"<!-- total advance width in font units: {total_width} -->")
