# Ports — Logo Guidelines

## 1. The logo
- **Idea**: A **port opening** — a bold ring with a gap, and the live indicator dot sitting in it. The dot is the listener the app is about.
- **Mark type**: abstract symbol (socket), used with a wordmark lockup.
- **Versions**: horizontal (primary) · stacked · symbol-only · wordmark-only.
- **Files**:
  - `ports-symbol.svg` — the mark alone (favicon, tray, app icon).
  - `ports-horizontal.svg` — primary lockup.
  - `ports-stacked.svg` — stacked lockup (tight spaces, square avatars).
  - `kit/ports-{black,white,mono-26e085}.svg` — one-colour versions.
  - `kit/ports-app-icon.svg` + `kit/*.png`, `kit/favicon.ico`, `kit/site.webmanifest` — the web/app icon set.

## 2. Clear space
Keep a clear zone of **half the symbol's height** around the lockup on all sides. The zone scales with the logo — never use a fixed distance.

## 3. Minimum size
| Version | Screen | Notes |
|---|---|---|
| Horizontal | 96 px wide | below this, use the symbol alone |
| Symbol | 16 px | use the simplified `-favicon` drawing below 24 px |

## 4. Colour
| Name | HEX | RGB | Use |
|---|---|---|---|
| Brand green | `#26E085` | 38 224 133 | the mark, on dark or white |
| Graphite | `#17181C` | 23 24 28 | the app-icon tile; type |
| Ink | `#0A0A0A` | 10 10 10 | one-colour mark on light |
| Paper | `#F6F5F6` | 246 245 246 | light canvas |

**Approved pairs**: brand green on graphite · brand green on white · ink on white · white on graphite.
On the dark tray, the mark is drawn as a **template icon** (black on transparent) and the system recolours it.

## 5. Typography
- Wordmark: **Geist SemiBold**, outlined (paths, not live text).
- UI: Geist; web fallback `ui-sans-serif, system-ui, sans-serif`.

## 6. Don'ts
Don't stretch, squash or rotate · don't recolour outside the palette · don't add shadows, outlines or gradients · don't reflow the lockup (symbol left, wordmark right) · don't recreate the wordmark by typing it · don't close the ring's gap or move the dot out of it.

## 7. Files
Masters live in `brand/logo/`. The app ships `resources/icon.png` (installer/app icon); the landing page ships `landing/favicon.png`.
