// Fig 3 — the restart interlock. A keyed lock bolt passes through the bored eye of
// a trip lever. The key must be turned before the bolt clears the eye, and only
// then can the lever travel and trip the process. This is the app's destructive
// action rule as a machine: the action is keyed, it cannot fire by accident, and
// the eye is what holds it.
//
// The lock body also carries the read head that takes the command line before the
// lever moves — the app reads the command line before stopping a process, because
// the command line is gone with the process.
import * as k from "./kit/iso-kit.mjs";
import * as G from "./kit/lathe.mjs";
import * as A from "./kit/audit.mjs";
import { at as atOf, byDepth as byDepthOf, line, plateStyle, dots } from "./figure-theme.mjs";

const WIDTH = 560;
const HEIGHT = 430;

const PLATE = { x: -142, y: -74, w: 288, d: 148, r: 10 };
const PLATE_H = 12;
const FOOT_H = 5;

// The lock body, on the right: the key on top, the bolt leaving its near face.
const BODY = { x: 34, y: -26, w: 76, d: 52, r: 4 };
const BODY_H = 34;
const BODY_TOP = BODY_H;
const KEY = [88, 0];
const KEY_R = 13;
const KEY_Z = BODY_TOP + 4;
const KEY_H = 9;
const KEY_TOP = KEY_Z + KEY_H;
const KEY_GUARD = 5;

// The bolt runs out of the body's near face along -x, through the lever's eye.
const BOLT_Z = 26;
const BOLT_R = 3.2;
const BOLT_FROM = BODY.x;
// The bolt's tip is the leftmost thing on its axis, so the eye it passes through
// sits just inside the tip and the rod reads as running right through it.
const BOLT_TO = -30;
const BOLT_SPAN = BOLT_FROM - BOLT_TO;
const BOLT_TRAVEL = 22;
// Two cradles the bolt rests in, clear of the eye. A bore in the middle of the rod
// would wrap it, and nothing on screen could then order the two; a cradle meets it
// tangentially, so the rod stays in front of it everywhere.
const CRADLES = [14, -4];

// The trip lever rides two rails; its eye is the only thing the bolt holds.
const RAIL_X = [-134, -14];
const RAIL_D = 128;
const RAIL_H = 4;
const PAD_Z = RAIL_H;
const PAD_H = 4;
const CARRIAGE = { x: -100, y: -64, w: 56, d: 128, r: 4 };
const CARRIAGE_Z = PAD_Z + PAD_H;
const CARRIAGE_H = 5;
const CARRIAGE_TOP = CARRIAGE_Z + CARRIAGE_H;
// The fork the bolt drops into. It is cut open at the top, because a closed eye
// around the rod is a genuine interlock: no drawing order can put one in front of
// the other. A fork holds the lever in the one direction that matters and can be
// ordered — its near jaw in front of the rod, its far jaw behind it.
const FORK = { x: -46, w: 18 };
const JAW_D = 5;
const JAW_IN = BOLT_R + 1.4;
const JAW_OUT = JAW_IN + JAW_D;
const POST = { x: FORK.x, y: -JAW_OUT, w: FORK.w, d: JAW_OUT * 2, r: 2 };
const POST_TOP = BOLT_Z - BOLT_R - 0.8;
const FORK_TOP = BOLT_Z + BOLT_R + 1.2;
const TRIP = 30;
const HANDLE = [CARRIAGE.x - 6, 0];
const GRIP_Z = CARRIAGE_TOP + 8;

const P = k.fitProjection(
  [
    ...k.boxCorners(PLATE, -PLATE_H - FOOT_H, 0),
    [BODY.x + BODY.w, BODY.y, BODY_TOP],
    [KEY[0] + KEY_R + KEY_GUARD, KEY[1], KEY_TOP],
    [BODY.x, BODY.y + BODY.d, BODY_TOP],
    [RAIL_X[0], -RAIL_D / 2 - 6, 10],
    [RAIL_X[1], RAIL_D / 2 + 6, 10],
    [CARRIAGE.x - TRIP, CARRIAGE.y, GRIP_Z + 6],
    [FORK.x, -JAW_OUT, FORK_TOP],
    [FORK.x + FORK.w, JAW_OUT, FORK_TOP],
    [-60, -96, 0],
    [30, 84, 0],
  ],
  WIDTH,
  HEIGHT,
  { pad: 48, azimuth: 42 },
);

const at = atOf(P);
const byDepth = byDepthOf(P);
const S = plateStyle;
const D = dots;
const up = (x, y, z) => G.frameAlong([x, y, z], [0, 0, 1]);
const alongX = (x, y, z) => G.frameAlong([x, y, z], [-1, 0, 0]);

// ── Base ─────────────────────────────────────────────────────────────────
function base() {
  const out = [`<path class="iso-halo" d="${k.haloOf(PLATE, -PLATE_H - FOOT_H, PLATE_H, P)}"/>`];
  const shapes = [];
  for (const [x, y] of byDepth(k.corners(PLATE, 24))) {
    out.push(S(k.cylinder(x, y, 7, -PLATE_H - FOOT_H, 1.6, P, 28), { tone: "lo" }));
    out.push(S(k.cylinder(x, y, 5.4, -PLATE_H - FOOT_H + 1.6, FOOT_H - 1.6, P, 28), { tone: "mid" }));
    shapes.push(A.cylinder(x, y, 7, -PLATE_H - FOOT_H, 1.6, { name: `base.foot${shapes.length}` }), A.cylinder(x, y, 5.4, -PLATE_H - FOOT_H + 1.6, FOOT_H - 1.6, { name: `base.foot${shapes.length}b` }));
  }
  out.push(S(k.slabOf(PLATE, -PLATE_H, PLATE_H, P, 10, 1.6), { tone: "mid" }));
  shapes.push(A.slab(PLATE, -PLATE_H, PLATE_H, { name: "base.plate", ground: true }, 10));
  out.push(line(k.planOutline(k.insetPlan(PLATE, 6), 0, P), { tone: "lo" }));
  const screws = k.corners(PLATE, 11);
  out.push(line(screws.map(([x, y]) => k.ring(x, y, 2.2, 0, P, 16)).join(""), { tone: "lo" }));
  out.push(D(screws.map(([x, y]) => at([x, y, 0])), { size: 0.5 }));
  return { out: out.join(""), shapes };
}

function scaleStrip() {
  const strip = { x: -132, y: PLATE.y + PLATE.d - 4, w: 132, d: 4, r: 1 };
  const out = [S(k.slabOf(strip, 0, 0.6, P, 3), { tone: "lo", crease: "none" })];
  const ticks = k.topTicks(strip.x + 4, strip.x + strip.w - 4, 6, 4, strip.y + 0.8, 0.6, [2, 3.6], P);
  out.push(line(ticks.minor, { tone: "lo" }), line(ticks.major, { tone: "mid" }));
  out.push(D([at([strip.x + 2, strip.y + 3.5, 0.6]), at([strip.x + strip.w - 2, strip.y + 3.5, 0.6])], { size: 0.45 }));
  return { out: out.join(""), shapes: [A.slab(strip, 0, 0.6, { name: "base.scale" }, 3)] };
}

// The warning plate beside the key. Its rule lines stand for printed text, so they
// are the one place a line may stop on a bare face.
function warningPlate() {
  const plan = { x: -128, y: PLATE.y + PLATE.d - 40, w: 74, d: 22, r: 2 };
  const out = [S(k.slabOf(plan, 0, 1, P, 4), { tone: "lo", crease: "none" }), line(k.planOutline(k.insetPlan(plan, 2), 1, P), { tone: "faint" })];
  for (const [row, length] of [[0, 52], [1, 38], [2, 24]]) {
    const lineY = plan.y + 6 + row * 5;
    out.push(line(k.lineOnTop([plan.x + 5, lineY], [plan.x + 5 + length, lineY], 1, P), { tone: "lo", free: true }));
  }
  out.push(D([at([plan.x + 4, plan.y + 4, 1]), at([plan.x + plan.w - 4, plan.y + plan.d - 4, 1])], { size: 0.45 }));
  return { out: out.join(""), shapes: [A.slab(plan, 0, 1, { name: "base.warning" }, 4)] };
}

// ── Rails and end stops ──────────────────────────────────────────────────
function rails() {
  const out = [];
  const shapes = [];
  for (const y of [-RAIL_D / 2, RAIL_D / 2 - 7]) {
    const rail = { x: RAIL_X[0], y, w: RAIL_X[1] - RAIL_X[0], d: 7, r: 1 };
    out.push(S(k.slabOf(rail, 0, RAIL_H, P, 3, 0.4), { tone: "mid" }));
    shapes.push(A.slab(rail, 0, RAIL_H, { name: `rail${y}` }, 3));
    out.push(D([rail.x + 10, rail.x + rail.w - 10].map((x) => at([x, y + 3.5, RAIL_H])), { size: 0.4 }));
  }
  // End stops are a block per rail, not a bar across both: a bar would run through
  // the bolt cradles, and nothing could then order the two on screen.
  for (const x of [RAIL_X[0] - 7, RAIL_X[1]]) {
    for (const y of [-RAIL_D / 2 - 4, RAIL_D / 2 - 3]) {
      const stop = { x, y, w: 7, d: 7, r: 1.4 };
      out.push(S(k.slabOf(stop, 0, 11, P, 3, 0.5), { tone: "mid" }));
      shapes.push(A.slab(stop, 0, 11, { name: `rail.stop${x}.${y}` }, 3));
      out.push(D([at([x + 3.5, y + 3.5, 11])], { size: 0.45 }));
    }
  }
  return { out: out.join(""), shapes };
}

// ── The lock body: housing, read-head cover, key switch ──────────────────
function body() {
  const out = [];
  const shapes = [];
  out.push(S(k.slabOf(BODY, 0, BODY_H, P, 6, 1.1), { tone: "mid" }));
  shapes.push(A.slab(BODY, 0, BODY_H, { name: "body.housing" }, 6));
  out.push(line(k.sideSeam(BODY, BODY_H - 7, P, 6), { tone: "faint" }));
  out.push(line(k.planOutline(k.insetPlan(BODY, 5), BODY_TOP, P), { tone: "faint" }));
  const screws = k.corners(BODY, 9);
  out.push(line(screws.map(([x, y]) => k.ring(x, y, 2, BODY_TOP, P, 14)).join(""), { tone: "lo" }));
  out.push(D(screws.map(([x, y]) => at([x, y, BODY_TOP])), { size: 0.45 }));

  // The read head, under a cover plate on the body's top: the command line is
  // taken here, before the lever is allowed to move.
  const cover = { x: BODY.x + 6, y: BODY.y + 7, w: 12, d: BODY.d - 14, r: 3 };
  out.push(S(k.slabOf(cover, BODY_TOP, 7, P, 4, 0.5), { tone: "mid" }));
  shapes.push(A.slab(cover, BODY_TOP, 7, { name: "body.cover" }, 4));
  out.push(line(k.sideSeam(cover, BODY_TOP + 2.8, P, 4), { tone: "faint" }));
  out.push(D([at([cover.x + 4, cover.y + 6, BODY_TOP + 7]), at([cover.x + 4, cover.y + cover.d - 6, BODY_TOP + 7])], { size: 0.42 }));

  // The key switch: a guarded barrel with a knurled rim and an index ring.
  const F = up(KEY[0], KEY[1], BODY_TOP);
  out.push(S(G.ringBand(0, 4, KEY_R, KEY_R + KEY_GUARD, F, P), { tone: "mid", crease: "lo" }));
  shapes.push(A.ring(F, 0, 4, KEY_R, KEY_R + KEY_GUARD, { name: "key.guard" }));
  out.push(line(G.arcOf(F, 4, KEY_R + KEY_GUARD, P, {}), { tone: "lo" }));
  out.push(D(G.dotsOf(F, 2, KEY_R + KEY_GUARD + 2.4, 6, P, { all: true }), { size: 0.45, tone: "lo" }));
  out.push(S(k.cylinder(KEY[0], KEY[1], KEY_R, KEY_Z, KEY_H, P, 48, 0.6), { tone: "hi", crease: "lo" }));
  shapes.push(A.cylinder(KEY[0], KEY[1], KEY_R, KEY_Z, KEY_H, { name: "key.barrel" }));
  out.push(line(k.knurl(KEY[0], KEY[1], KEY_R, KEY_Z + 1, KEY_TOP - 1, 48, P, 0.35), { tone: "lo" }));
  out.push(line(k.sideArc(KEY[0], KEY[1], KEY_R, KEY_Z + 1, P, 34), { tone: "faint" }));
  out.push(line(k.ring(KEY[0], KEY[1], KEY_R - 2.6, KEY_TOP, P, 48), { tone: "lo" }));
  const ticks = k.radialTicks(KEY[0], KEY[1], KEY_R - 3.4, 8, 2, KEY_TOP, [1.4, 3], P);
  out.push(line(ticks.minor, { tone: "lo" }), line(ticks.major, { tone: "mid" }));
  out.push(line(k.ring(KEY[0], KEY[1], 1.5, KEY_TOP, P, 14), { tone: "mid" }));

  // The key's pointer. It turns, so it is redrawn each frame, but it belongs to the
  // body group — a line has no volume for the audit to check.
  out.push(`<path data-part="pointer" class="iso-line" data-tone="lit" data-free="start" d="${pointerAt(0)}" style="stroke-width:1.3"/>`);
  return { out: out.join(""), shapes, at: [KEY[0], KEY[1], BODY_H / 2] };
}

// The pointer that shows how far the key has turned. It is the one part the reader
// watches while the machine unlocks, so it carries the lit accent.
const pointerAt = (turn) => {
  const [x, y] = KEY;
  return k.lineOnTop([x + Math.sin(turn) * 2, y - Math.cos(turn) * 2], [x + Math.sin(turn) * (KEY_R - 3), y - Math.cos(turn) * (KEY_R - 3)], KEY_TOP + 0.3, P);
};

// ── The bolt: one rod out of a gland, through the lever's eye ────────────
function bolt() {
  const out = [];
  const F = alongX(BOLT_FROM, 0, BOLT_Z);
  const rod = G.rodOut(F, 0, BOLT_SPAN, BOLT_R, P, { stroke: BOLT_TRAVEL });
  out.push(`<clipPath id="iso-bolt-clip"><path data-hold d="${rod.clip}"/></clipPath>`);
  out.push(`<g clip-path="url(#iso-bolt-clip)">${S(rod.paths, { tone: "hi", lit: true, crease: "lo" })}</g>`);
  // The gland it leaves, and the seal ring on its face.
  out.push(S(G.ringBand(0, 3.4, BOLT_R + 0.6, BOLT_R + 4, F, P), { tone: "mid", crease: "lo" }));
  out.push(line(G.arcOf(F, 3.4, BOLT_R + 4, P, {}), { tone: "lo" }));
  out.push(line(k.sideArc(BOLT_FROM - 9, 0, BOLT_R, BOLT_Z, P, 22), { tone: "faint" }));
  return {
    out: out.join(""),
    shapes: [A.slide(F, 0, BOLT_SPAN, BOLT_R, { name: "bolt.rod" }), A.ring(F, 0, 3.4, BOLT_R + 0.6, BOLT_R + 4, { name: "bolt.gland" })],
    at: [BOLT_FROM - BOLT_SPAN / 2, 0, BOLT_Z],
  };
}

// Two cradles the bolt rests in, with a shallow seat cut across each top.
function cradles() {
  const out = [];
  const shapes = [];
  for (const x of CRADLES) {
    const block = { x: x - 8, y: -11, w: 16, d: 22, r: 2 };
    const top = BOLT_Z - BOLT_R;
    out.push(S(k.slabOf(block, 0, top, P, 4, 0.5), { tone: "mid" }));
    shapes.push(A.slab(block, 0, top, { name: `cradle${x}` }, 4));
    const seat = { x: x - 5, y: -5, w: 10, d: 10, r: 2 };
    out.push(line(k.planOutline(seat, top, P, 4), { tone: "faint" }));
    out.push(D([at([block.x + 4, block.y + block.d - 4, top]), at([block.x + block.w - 4, block.y + block.d - 4, top])], { size: 0.42 }));
  }
  return { out: out.join(""), shapes };
}

// ── The trip lever: a carriage on the rails, with the eye the bolt holds ─
function lever() {
  const out = [`<path class="iso-halo" d="${k.haloOf(k.insetPlan(CARRIAGE, 4), CARRIAGE_Z, 0.1, P)}" style="filter:blur(5px)"/>`];
  const shapes = [];
  // Pads resting on the rails' top faces. Nothing wraps a rail, so the carriage is
  // in front of every rail wherever the two overlap.
  for (const y of [-RAIL_D / 2 + 3.5, RAIL_D / 2 - 3.5]) {
    const pad = { x: CARRIAGE.x + 6, y: y - 3.5, w: CARRIAGE.w - 12, d: 7, r: 1 };
    out.push(S(k.slabOf(pad, PAD_Z, PAD_H, P, 3, 0.5), { tone: "mid" }));
    shapes.push(A.slab(pad, PAD_Z, PAD_H, { name: `lever.pad${y}`, owner: "lever" }, 3));
    out.push(D([pad.x + 6, pad.x + pad.w - 6].map((x) => at([x, y, PAD_Z + PAD_H])), { size: 0.42 }));
  }
  out.push(S(k.slabOf(CARRIAGE, CARRIAGE_Z, CARRIAGE_H, P, 6, 0.7), { tone: "mid" }));
  shapes.push(A.slab(CARRIAGE, CARRIAGE_Z, CARRIAGE_H, { name: "lever.plate", owner: "lever" }, 6));
  out.push(line(k.sideSeam(CARRIAGE, CARRIAGE_Z + CARRIAGE_H - 1.8, P, 6), { tone: "faint" }));
  out.push(D([CARRIAGE.y + 14, CARRIAGE.y + CARRIAGE.d - 14].map((y) => at([CARRIAGE.x + CARRIAGE.w - 9, y, CARRIAGE_TOP])), { size: 0.4, tone: "lo" }));

  // The post the fork stands on, rising from the carriage to just under the rod.
  out.push(S(k.slabOf(POST, CARRIAGE_TOP, POST_TOP - CARRIAGE_TOP, P, 4, 0.5), { tone: "mid" }));
  shapes.push(A.slab(POST, CARRIAGE_TOP, POST_TOP - CARRIAGE_TOP, { name: "lever.post", owner: "lever" }, 4));
  out.push(D([at([POST.x + 4, POST.y + POST.d - 4, POST_TOP]), at([POST.x + POST.w - 4, POST.y + 4, POST_TOP])], { size: 0.42 }));
  // A rib along the post's near face, so the load path reads.
  out.push(line(k.segment([POST.x + 5, POST.y + POST.d, CARRIAGE_TOP + 3], [POST.x + 5, POST.y + POST.d, POST_TOP - 3], P), { tone: "faint" }));

  // The trip handle: a knurled grip the operator pulls.
  out.push(S(k.slabOf({ x: HANDLE[0] - 5, y: -7, w: 11, d: 14, r: 2 }, CARRIAGE_TOP, 8, P, 3, 0.5), { tone: "mid" }));
  shapes.push(A.slab({ x: HANDLE[0] - 5, y: -7, w: 11, d: 14, r: 2 }, CARRIAGE_TOP, 8, { name: "lever.neck", owner: "lever" }, 3));
  const grip = [HANDLE[0] + 0.5, 0];
  out.push(S(k.cylinder(grip[0], grip[1], 6.5, GRIP_Z, 6, P, 32, 0.5), { tone: "hi", crease: "lo" }));
  shapes.push(A.cylinder(grip[0], grip[1], 6.5, GRIP_Z, 6, { name: "lever.grip", owner: "lever" }));
  out.push(line(k.knurl(grip[0], grip[1], 6.5, GRIP_Z + 0.6, GRIP_Z + 5.4, 40, P, 0.3), { tone: "lo" }));
  out.push(line(k.sideArc(grip[0], grip[1], 6.5, GRIP_Z + 0.4, P, 26), { tone: "faint" }));
  out.push(line(k.ring(grip[0], grip[1], 6.5, GRIP_Z + 6, P, 32), { tone: "mid" }));
  return { out: out.join(""), shapes, at: [CARRIAGE.x + CARRIAGE.w / 2, 0, CARRIAGE_Z + CARRIAGE_H] };
}

// ── The fork: two jaws that straddle the rod ─────────────────────────────
// The far jaw is drawn before the bolt and the near jaw after it, which is what
// makes the two orderable. Each jaw is a block with a screw, and the pair is the
// mouth the rod sits in.
function fork() {
  const jaw = (side) => {
    const y = side === "far" ? -JAW_OUT : JAW_IN;
    const plan = { x: FORK.x, y, w: FORK.w, d: JAW_D, r: 1.5 };
    const out = [S(k.slabOf(plan, POST_TOP, FORK_TOP - POST_TOP, P, 3, 0.5), { tone: "hi", crease: "lo" })];
    out.push(D([at([plan.x + 4, y + JAW_D / 2, FORK_TOP]), at([plan.x + plan.w - 4, y + JAW_D / 2, FORK_TOP])], { size: 0.4, tone: "lo" }));
    // The jaw's mouth is the face toward the rod, so the line reads as the grip.
    out.push(line(k.segment([plan.x, y + (side === "far" ? JAW_D : 0), POST_TOP + 0.6], [plan.x + plan.w, y + (side === "far" ? JAW_D : 0), POST_TOP + 0.6], P), { tone: "faint" }));
    return { out: out.join(""), shapes: [A.slab(plan, POST_TOP, FORK_TOP - POST_TOP, { name: `fork.${side}`, owner: "lever" }, 3)], at: [plan.x + plan.w / 2, y + JAW_D / 2, POST_TOP] };
  };
  return { far: jaw("far"), near: jaw("near") };
}

// ── Assemble ─────────────────────────────────────────────────────────────
export function buildInterlock() {
  const R = A.recorder(P);
  const items = [];
  const put = (name, point, svg, shapes, bias = 0) =>
    items.push({ name, at: point, key: k.depthOf(point, P) + bias, svg, shapes: shapes.map((shape) => R.solid(shape)) });

  const parts = { base: base(), scale: scaleStrip(), warning: warningPlate(), rails: rails(), cradles: cradles(), body: body(), bolt: bolt(), lever: lever(), fork: fork() };
  put("base.plate", [0, 0, -PLATE_H / 2], parts.base.out, parts.base.shapes, -8);
  put("base.scale", [-60, PLATE.y + PLATE.d - 10, 0.3], parts.scale.out, parts.scale.shapes, 0.2);
  put("base.warning", [-92, PLATE.y + PLATE.d - 28, 0.5], parts.warning.out, parts.warning.shapes, 0.2);
  put("rails", [0, 0, 2], parts.rails.out, parts.rails.shapes, -1);
  put("cradles", [0, 0, BOLT_Z / 2], parts.cradles.out, parts.cradles.shapes, 0.6);
  put("body", parts.body.at, parts.body.out, parts.body.shapes, 1.4);

  // The bolt withdraws into the body; the lever travels only once it has cleared.
  // The two moves are away from each other, so they can never collide in flight.
  // The far jaw is painted before the bolt and the near jaw after it: that split is
  // what lets a rod pass through the fork with no drawing order left undecided.
  items.push({ name: "fork.far", at: parts.fork.far.at, key: k.depthOf(parts.bolt.at, P) - 0.5, svg: `<g data-part="fork" transform="${k.translateAlong(P, [0, 0, 0])}">${parts.fork.far.out}</g>`, shapes: parts.fork.far.shapes.map((shape) => R.solid(shape)), move: [-TRIP, 0, 0], ride: "trip" });
  items.push({
    name: "bolt",
    at: parts.bolt.at,
    key: k.depthOf(parts.bolt.at, P) + 1.0,
    svg: `<g data-part="bolt" transform="${k.translateAlong(P, [0, 0, 0])}">${parts.bolt.out}</g>`,
    shapes: parts.bolt.shapes.map((shape) => R.solid(shape)),
    move: [BOLT_TRAVEL, 0, 0],
    ride: "unlock",
  });
  items.push({
    name: "lever",
    at: parts.lever.at,
    key: k.depthOf(parts.lever.at, P) + 1.2,
    svg: `<g data-part="lever" transform="${k.translateAlong(P, [0, 0, 0])}">${parts.lever.out}</g>`,
    shapes: parts.lever.shapes.map((shape) => R.solid(shape)),
    move: [-TRIP, 0, 0],
    ride: "trip",
  });
  items.push({ name: "fork.near", at: parts.fork.near.at, key: k.depthOf(parts.bolt.at, P) + 2, svg: `<g data-part="fork" transform="${k.translateAlong(P, [0, 0, 0])}">${parts.fork.near.out}</g>`, shapes: parts.fork.near.shapes.map((shape) => R.solid(shape)), move: [-TRIP, 0, 0], ride: "trip" });

  R.items = items;
  R.order = ["unlock", "trip"];
  R.spring = { stagger: 0.3, stiffness: 52, damping: 13 };
  A.settle(R, P, { states: A.explodeStates(R.order, R.spring).depth });
  items.sort((a, b) => a.key - b.key);

  const svg = k.figureSvg({
    width: WIDTH,
    height: HEIGHT,
    label: `A keyed restart interlock. A lock bolt runs out of the lock body and through the bored eye of a trip lever. Turning the key withdraws the bolt from the eye; only then can the lever travel along its rails and trip the process. The lock body carries the read head that takes the command line before the lever moves.`,
    body: items.map((item) => item.svg),
  });

  const plate = k.plateHtml({
    fig: "Fig 3",
    title: "Restart interlock",
    hint: "Hover to key and trip",
    readout: "keyed · bolt in the eye · holding",
    keys: [
      { mark: "lit", label: "The bolt: while it is in the eye, the lever is held" },
      { mark: "raised", label: "The key: turned before anything moves" },
      { mark: "flat", label: "Trip lever: the process it stops" },
      { mark: "dim", label: "Read head: takes the command line first" },
    ],
    caption: `Every destructive action in the app asks first and names the process. Here that is a keyed interlock: the bolt sits in the lever's eye until the key is turned, and the lever cannot travel while it is held. The lock body's read head takes the command line before the lever moves, because the command line is gone with the process — which is what makes restart possible.`,
    body: `<div class="iso iso-stage" data-figure="interlock" tabindex="0" role="group" aria-label="Restart interlock. Move the pointer or use the arrow keys to turn the key and trip the lever.">${svg}<span class="iso-sr" role="status" aria-live="polite" data-status></span></div>`,
  });

  const live = `
const P = ${JSON.stringify(P)};
const KEY = ${JSON.stringify(KEY)};
const KEY_R = ${KEY_R}, KEY_TOP = ${KEY_TOP}, TURN = ${Math.PI / 2};
const BOLT_TRAVEL = ${BOLT_TRAVEL}, TRIP = ${TRIP};
const pointerAt = ${pointerAt.toString().replace("k.lineOnTop", "lineOnTop")};
const clamp = (value, low, high) => Math.min(high, Math.max(low, value));

function mount(root) {
  const stage = root.querySelector(".iso-stage") || root;
  const svg = stage.querySelector("svg");
  const bolt = svg.querySelector('[data-part="bolt"]');
  const lever = svg.querySelector('[data-part="lever"]');
  const pointer = svg.querySelector('[data-part="pointer"]');
  // This figure's own readout: three figures share one page, and the first one in
  // the document is not necessarily this one.
  const readout = (stage.closest(".iso-figure") || document).querySelector("[data-readout]");
  const status = stage.querySelector("[data-status]");
  const still = matchMedia("(prefers-reduced-motion: reduce)");
  let phase = 0, turn = 0, trip = 0, turnRate = 0, tripRate = 0;
  let frame = 0, last = 0, clock = 0, touring = true, visible = false, idleTimer = 0;

  // Three phases: holding, keyed with the bolt clear, tripped. The bolt leads the
  // lever by a wide stagger, so the lever can never move while it is held.
  const shares = (want) => ({ unlock: clamp(want, 0, 1), trip: clamp(want - 1, 0, 1) });

  function draw() {
    const want = shares(phase);
    bolt.setAttribute("transform", translateAlong(P, [BOLT_TRAVEL * want.unlock, 0, 0]));
    holdStill(bolt, axisVector(P, [BOLT_TRAVEL * want.unlock, 0, 0]));
    lever.setAttribute("transform", translateAlong(P, [-TRIP * want.trip, 0, 0]));
    pointer.setAttribute("d", pointerAt(turn));
    const text = want.trip > 0.5 ? "tripped \\u00b7 lever released" : want.unlock > 0.5 ? "keyed \\u00b7 bolt clear" : "keyed \\u00b7 bolt in the eye \\u00b7 holding";
    if (readout.textContent !== text) readout.textContent = text;
    const said = want.trip > 0.5 ? "Tripped: the lever has released the process" : want.unlock > 0.5 ? "Keyed: the bolt has cleared the eye" : "Holding: the bolt is in the eye";
    if (status.textContent !== said) status.textContent = said;
  }

  function tick(now) {
    const dt = last ? Math.min((now - last) / 1000, 1 / 30) : 1 / 60;
    last = now;
    const calm = still.matches;
    if (touring) {
      clock += dt;
      const at = clock % 7.2;
      phase = at < 2.4 ? 0 : at < 5 ? 1 : 2;
    }
    const want = shares(phase);
    const step = (value, goal, rate) => {
      if (calm) return [goal, 0];
      const stiff = 180, drag = 2 * Math.sqrt(180);
      for (let i = 0; i < 4; i++) { rate += ((goal - value) * stiff - rate * drag) * (dt / 4); value += rate * (dt / 4); }
      if (Math.abs(goal - value) < 0.003 && Math.abs(rate) < 0.003) return [goal, 0];
      return [value, rate];
    };
    [turn, turnRate] = step(turn, want.unlock * TURN, turnRate);
    [trip, tripRate] = step(trip, want.trip * -TRIP, tripRate);
    draw();
    const settled = !touring && Math.abs(turn - want.unlock * TURN) < 1e-6 && Math.abs(trip - want.trip * -TRIP) < 1e-6;
    frame = visible && !settled ? requestAnimationFrame(tick) : 0;
  }

  const run = () => { if (!frame && visible) { last = 0; frame = requestAnimationFrame(tick); } };
  const steer = (next) => {
    touring = false;
    phase = clamp(next, 0, 2);
    clearTimeout(idleTimer);
    idleTimer = setTimeout(() => { touring = !still.matches; clock = 0; run(); }, 3600);
    run();
  };

  new IntersectionObserver((entries) => {
    visible = entries[entries.length - 1].isIntersecting;
    if (visible) run();
  }, { rootMargin: "120px 0px" }).observe(stage);

  if (still.matches) touring = false;
  stage.addEventListener("pointermove", (event) => {
    const box = svg.getBoundingClientRect();
    steer(clamp((event.clientX - box.left) / box.width, 0, 1) * 2);
  });
  stage.addEventListener("pointerleave", () => steer(phase));
  stage.addEventListener("keydown", (event) => {
    let to = null;
    if (event.key === "ArrowRight" || event.key === "ArrowDown") to = Math.min(2, Math.round(phase) + 1);
    if (event.key === "ArrowLeft" || event.key === "ArrowUp") to = Math.max(0, Math.round(phase) - 1);
    if (event.key === "Home") to = 0;
    if (event.key === "End") to = 2;
    if (event.key === " " || event.key === "Enter") to = Math.round(phase) >= 2 ? 0 : Math.round(phase) + 1;
    if (to === null) return;
    event.preventDefault();
    steer(to);
  });
  draw();
}
`;
  return { plate, live: `(function(){\n${live}\nmount(document.querySelector('[data-figure="interlock"]'));\n})();`, R, P, solids: R.solids.length, items: items.length };
}










