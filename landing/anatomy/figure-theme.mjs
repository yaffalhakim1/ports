// Shared chrome for the landing page's figures: the brand-toned anatomy palette
// and the plate styles. The page ships one theme (dark graphite), so these
// tokens are the whole set — there is no light build to guard against.
import * as k from "./kit/iso-kit.mjs";

export const FIGURE_CSS = `
/* The figure takes the page's plate, ink and brand green, at one theme. */
.iso{
  --anatomy-card:oklch(15% 0.006 275);
  --anatomy-paper:oklch(19.5% 0.008 275);
  --anatomy-face:oklch(24% 0.010 275);
  --anatomy-top:oklch(25% 0.010 275);
  --anatomy-shade-0:oklch(16% 0.006 275);
  --anatomy-shade-1:oklch(17.5% 0.007 275);
  --anatomy-shade-2:oklch(19% 0.007 275);
  --anatomy-shade-3:oklch(21% 0.008 275);
  --anatomy-lit-top:oklch(40% 0.100 155);
  --anatomy-lit-shade:oklch(30% 0.070 155);
  --anatomy-hi:oklch(93% 0.006 275);
  --anatomy-mid:oklch(62% 0.010 275);
  --anatomy-lo:oklch(44% 0.010 275);
  --anatomy-faint:oklch(33% 0.010 275);
  --anatomy-lit:oklch(80% 0.190 155);
  --anatomy-halo:oklch(0% 0 0);
  --anatomy-halo-opacity:.55;
  --anatomy-dot:oklch(88% 0.008 275);
  --anatomy-red:oklch(66% 0.180 25);
  --anatomy-green:oklch(80% 0.190 155);
  --anatomy-blue:oklch(70% 0.130 250);
  --anatomy-weight:.7px;
  --anatomy-ease:cubic-bezier(.23,1,.32,1);
}
.iso-plate{
  background:var(--color-plate);
  box-shadow:inset 0 0 0 1px var(--color-rule-2);
  border-radius:var(--radius);
  padding:56px var(--space-md) 64px;
}
.iso-plate-corners{inset:var(--space-md) var(--space-md) 0}
.iso-plate-corner{font-family:var(--font-mono);font-size:var(--text-2xs);letter-spacing:.06em;color:var(--color-dim)}
.iso-plate-corner[data-corner="fig"]{color:var(--color-accent)}
.iso-plate-corner[data-corner="title"]{color:var(--color-ink-2)}
.iso-plate-corner[data-corner="readout"]{color:var(--color-ink)}
.iso-plate-corner[data-corner="hint"]{color:var(--color-dim)}
.iso-figure{margin:0}
.iso-legend{margin-top:var(--space-sm);font-size:var(--text-sm);color:var(--color-muted);line-height:1.6}
.iso-keys{gap:6px 20px}
.iso-key{gap:8px}
.iso-caption{max-width:var(--measure);display:block}
.iso-stage{outline:none;touch-action:pan-y}
.iso-stage:focus-visible{outline:2px solid var(--color-focus);outline-offset:6px;border-radius:var(--radius)}
.iso-sr{position:absolute;width:1px;height:1px;overflow:hidden;clip-path:inset(50%);white-space:nowrap}
/* On a phone the card is narrower than the drawing, so a wide row would flatten
   into a line of dots. Stretch the figure past the card by the kit's own
   measured range and clip it: the trim lands inside the drawing's outer margin,
   never on a part. */
@media (max-width: 40rem){
  .iso-plate{overflow:hidden;padding:48px 12px 68px}
  .iso-plate-corners{inset:var(--space-sm) var(--space-sm) 0}
  .iso-plate-corner[data-corner="hint"]{bottom:16px}
  .iso-stage[data-wide] svg{width:116%;margin:-8% -8%}
}
`;

export const plateStyle = (paths, style) => k.solidSvg(paths, style);
export const at = (P) => (point) => k.iso(point, P);
export const line = (d, style) => k.lineSvg(d, style);
export const dots = (points, style) => k.dotsSvg(points, style);

// Feet first so the near ones paint over the far ones.
export const byDepth = (P) => (points) =>
  [...points].sort((a, b) => k.depthOf([a[0], a[1], 0], P) - k.depthOf([b[0], b[1], 0], P));

export const wrap = (css) => `<style>${k.ISO_CSS}${css}</style>`;
