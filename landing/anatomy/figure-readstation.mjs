// Fig 2 — the parameter read station. Two processes of the same program ride a
// shuttle; a fixed spring feeler descends from the top arm and reads the working
// directory engraved on the block presented under it. This is the app's answer to
// "two node.exe — which project owns each?": the port says what, the folder says
// which.
//
// The reading is read-only by construction: the feeler has no cutter, it is a
// spring probe. The blocks move; the head does not.
import * as k from "./kit/iso-kit.mjs";
import * as G from "./kit/lathe.mjs";
import * as A from "./kit/audit.mjs";
import { at as atOf, byDepth as byDepthOf, line, plateStyle, dots } from "./figure-theme.mjs";

const WIDTH = 560;
const HEIGHT = 440;

// The two rows the app shows side by side: one program, two projects.
export const BLOCKS = [
  { process: "node.exe", folder: "portfolio", port: 3000 },
  { process: "node.exe", folder: "waku", port: 5173 },
];

const PLATE = { x: -140, y: -76, w: 292, d: 152, r: 10 };
const PLATE_H = 12;
const FOOT_H = 5;
const RAIL_X = [-86, 66];
const RAIL_D = 140;
const SHUTTLE = { x: -96, y: -58, w: 186, d: 116, r: 5 };
const SHUTTLE_Z = 4;
const SHUTTLE_H = 6;
const BLOCK = { x: -32, w: 64, d: 36, r: 3 };
const BLOCK_Z = SHUTTLE_Z + SHUTTLE_H;
const BLOCK_H = 34;
const BLOCK_TOP = BLOCK_Z + BLOCK_H;
const BLOCK_Y = [-48, 4];
const STRIP_H = 2.5;
const COLUMN = { x: -140, y: -44, w: 36, d: 76, r: 3 };
const COLUMN_TOP = 150;
const ARM = { x: -140, y: -46, w: 172, d: 84, r: 3 };
const ARM_Z = COLUMN_TOP;
const ARM_H = 24;
const HEAD = [0, -30];
const COLLAR_Z = ARM_Z;
const COLLAR_H = 5;
const BORE_R = 2.4;
const COLLAR_R = 6.4;
const FEELER_R = 2.2;
const TIP_H = 3.2;
const ROD_END = COLLAR_Z - (BLOCK_TOP + STRIP_H);
const LIFT = 16;

// The shuttle presents one block at a time under the fixed head, so its travel is
// the distance between the two block centres.
const CENTRES = BLOCK_Y.map((y) => y + BLOCK.d / 2);
const HEAD_Y = HEAD[1];
const SHUTTLE_AT = CENTRES.map((centre) => HEAD_Y - centre);
const HOME = 0;

const P = k.fitProjection(
  [
    ...k.boxCorners(PLATE, -PLATE_H - FOOT_H, 0),
    [COLUMN.x, COLUMN.y, COLUMN_TOP],
    [COLUMN.x + COLUMN.w, COLUMN.y + COLUMN.d, COLUMN_TOP],
    [ARM.x + ARM.w, ARM.y, ARM_Z + ARM_H],
    [HEAD[0] + COLLAR_R, HEAD[1], COLLAR_Z],
    [SHUTTLE.x, SHUTTLE.y, BLOCK_TOP + STRIP_H],
    [SHUTTLE.x + SHUTTLE.w, SHUTTLE.y + SHUTTLE.d, BLOCK_TOP + STRIP_H],
    [BLOCK.x + BLOCK.w, BLOCK_Y[0], BLOCK_TOP],
    [BLOCK.x, BLOCK_Y[1] + BLOCK.d, BLOCK_TOP],
    [-20, -110, 0],
    [20, 80, 0],
  ],
  WIDTH,
  HEIGHT,
  { pad: 46, azimuth: 42 },
);

const at = atOf(P);
const byDepth = byDepthOf(P);
const S = plateStyle;
const D = dots;
const up = (x, y, z) => G.frameAlong([x, y, z], [0, 0, 1]);
const down = (x, y, z) => G.frameAlong([x, y, z], [0, 0, -1]);
const blockAt = (index) => ({ x: BLOCK.x, y: BLOCK_Y[index], w: BLOCK.w, d: BLOCK.d, r: BLOCK.r });

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

// The graduated strip along the front edge: the shuttle's position, in the same
// units as the rail's rack.
function scaleStrip() {
  const strip = { x: -70, y: PLATE.y + PLATE.d - 15, w: 130, d: 7, r: 1 };
  const out = [S(k.slabOf(strip, 0, 0.6, P, 3), { tone: "lo", crease: "none" })];
  const ticks = k.topTicks(strip.x + 4, strip.x + strip.w - 4, 6, 4, strip.y + 0.8, 0.6, [2, 3.6], P);
  out.push(line(ticks.minor, { tone: "lo" }), line(ticks.major, { tone: "mid" }));
  out.push(D([at([strip.x + 2, strip.y + 3.5, 0.6]), at([strip.x + strip.w - 2, strip.y + 3.5, 0.6])], { size: 0.45 }));
  return { out: out.join(""), shapes: [A.slab(strip, 0, 0.6, { name: "base.scale" }, 3)] };
}

function labelPlate() {
  const plan = { x: 80, y: PLATE.y + PLATE.d - 23, w: 58, d: 15, r: 2 };
  const out = [S(k.slabOf(plan, 0, 0.8, P, 4), { tone: "lo", crease: "none" }), line(k.planOutline(k.insetPlan(plan, 1.8), 0.8, P), { tone: "faint" })];
  for (const length of [40, 27, 17]) out.push(line(k.lineOnTop([plan.x + 5, plan.y + plan.d / 2], [plan.x + 5 + length, plan.y + plan.d / 2], 0.8, P), { tone: "lo", free: true }));
  out.push(D([at([plan.x + 4, plan.y + 4, 0.8]), at([plan.x + plan.w - 4, plan.y + plan.d - 4, 0.8])], { size: 0.45 }));
  return { out: out.join(""), shapes: [A.slab(plan, 0, 0.8, { name: "base.plate-label" }, 4)] };
}

// ── Rails, rack and the shuttle that presents each block ─────────────────
function rails() {
  const out = [];
  const shapes = [];
  for (const x of RAIL_X) {
    const rail = { x, y: -RAIL_D / 2, w: 7, d: RAIL_D, r: 1 };
    out.push(S(k.slabOf(rail, 0, 4, P, 3, 0.4), { tone: "mid" }));
    shapes.push(A.slab(rail, 0, 4, { name: `rail${x}` }, 3));
    out.push(D([-RAIL_D / 2 + 8, RAIL_D / 2 - 8].map((dy) => at([x + 3.5, dy, 4])), { size: 0.4 }));
  }
  // End stops at both ends of the travel, with a screw each.
  for (const y of [-RAIL_D / 2 - 6, RAIL_D / 2]) {
    const stop = { x: RAIL_X[0] - 4, y, w: RAIL_X[1] - RAIL_X[0] + 15, d: 6, r: 1.4 };
    out.push(S(k.slabOf(stop, 0, 9, P, 3, 0.5), { tone: "mid" }));
    shapes.push(A.slab(stop, 0, 9, { name: `rail.stop${y}` }, 3));
    out.push(D([at([stop.x + 5, y + 3, 9]), at([stop.x + stop.w - 5, y + 3, 9])], { size: 0.45 }));
  }
  // The rack the shuttle's pointer reads, hanging from the front rail's face.
  const rack = k.sideTicks(RAIL_X[0] + 10, RAIL_X[1] - 10, 8, 5, RAIL_D / 2, 4, [2.2, 4], P);
  out.push(line(rack.minor, { tone: "lo" }), line(rack.major, { tone: "mid" }));
  return { out: out.join(""), shapes };
}

function shuttle() {
  const out = [`<path class="iso-halo" d="${k.haloOf(k.insetPlan(SHUTTLE, 4), SHUTTLE_Z, 0.1, P)}" style="filter:blur(5px)"/>`];
  const shapes = [];
  out.push(S(k.slabOf(SHUTTLE, SHUTTLE_Z, SHUTTLE_H, P, 6, 0.7), { tone: "mid" }));
  shapes.push(A.slab(SHUTTLE, SHUTTLE_Z, SHUTTLE_H, { name: "shuttle.plate", owner: "shuttle" }, 6));
  out.push(line(k.sideSeam(SHUTTLE, SHUTTLE_Z + SHUTTLE_H - 2.4, P, 6), { tone: "faint" }));
  // Bearing blocks where the shuttle meets each rail, and the rivet row.
  for (const x of RAIL_X) {
    const pad = { x: x - 3, y: SHUTTLE.y + 6, w: 13, d: SHUTTLE.d - 12, r: 2 };
    out.push(S(k.slabOf(pad, SHUTTLE_Z + SHUTTLE_H, 4, P, 3, 0.5), { tone: "mid" }));
    shapes.push(A.slab(pad, SHUTTLE_Z + SHUTTLE_H, 4, { name: `shuttle.bearing${x}`, owner: "shuttle" }, 3));
    out.push(D([SHUTTLE.y + 12, SHUTTLE.y + SHUTTLE.d - 12].map((y) => at([x + 3.5, y, SHUTTLE_Z + SHUTTLE_H + 4])), { size: 0.45 }));
  }
  out.push(D([-70, -40, 40, 70].map((dx) => at([dx, SHUTTLE.y + SHUTTLE.d - 3, SHUTTLE_Z + SHUTTLE_H])), { size: 0.4, tone: "lo" }));
  // The pointer that reads the rack.
  out.push(S(k.slabOf({ x: -12, y: SHUTTLE.y - 7, w: 24, d: 7, r: 1.5 }, SHUTTLE_Z + SHUTTLE_H, 3, P, 3, 0.4), { tone: "mid" }));
  shapes.push(A.slab({ x: -12, y: SHUTTLE.y - 7, w: 24, d: 7, r: 1.5 }, SHUTTLE_Z + SHUTTLE_H, 3, { name: "shuttle.pointer", owner: "shuttle" }, 3));
  out.push(line(k.lineOnTop([0, SHUTTLE.y - 7], [0, SHUTTLE.y - 13], SHUTTLE_Z + SHUTTLE_H + 3, P), { tone: "hi", free: "end" }));
  return { out: out.join(""), shapes, at: [0, 0, SHUTTLE_Z + SHUTTLE_H / 2] };
}

// One process parameter block: the housing, the folder field on top, the command
// line engraved on the front face, and the port it listens on.
function parameterBlock(index) {
  const plan = blockAt(index);
  const out = [];
  const shapes = [];
  out.push(S(k.slabOf(plan, BLOCK_Z, BLOCK_H, P, 5, 0.9), { tone: "hi", crease: "lo" }));
  shapes.push(A.slab(plan, BLOCK_Z, BLOCK_H, { name: `block${index}.body`, owner: `block${index}` }, 5));
  out.push(line(k.sideSeam(plan, BLOCK_Z + BLOCK_H - 8, P, 5), { tone: "faint" }));
  out.push(D([plan.x + 6, plan.x + plan.w - 6].map((x) => at([x, plan.y + plan.d - 6, BLOCK_Z + BLOCK_H])), { size: 0.45 }));

  // The folder field: a raised strip whose rule lines stand for the engraved
  // directory, the one place a line may stop on a bare face.
  const strip = k.insetPlan(plan, 5);
  out.push(S(k.slabOf(strip, BLOCK_TOP, STRIP_H, P, 4, 0.4), { tone: "mid", crease: "lo" }));
  shapes.push(A.slab(strip, BLOCK_TOP, STRIP_H, { name: `block${index}.folder`, owner: `block${index}` }, 4));
  const face = BLOCK_TOP + STRIP_H;
  for (const [row, length] of [[0, 40], [1, 30], [2, 22]]) {
    const lineY = strip.y + 7 + row * 6;
    out.push(line(k.lineOnTop([strip.x + 6, lineY], [strip.x + 6 + length, lineY], face, P), { tone: "lo", free: true }));
  }
  out.push(D([at([strip.x + 4, strip.y + strip.d - 4, face]), at([strip.x + strip.w - 4, strip.y + 4, face])], { size: 0.42 }));

  // The command line, engraved on the front face and running edge to edge.
  const front = plan.y + plan.d;
  for (const z of [BLOCK_Z + 12, BLOCK_Z + 18, BLOCK_Z + 24]) {
    out.push(line(k.segment([plan.x + 4, front, z], [plan.x + plan.w - 4, front, z], P), { tone: "faint" }));
  }
  return { out: out.join(""), shapes, at: [plan.x + plan.w / 2, plan.y + plan.d / 2, BLOCK_Z + BLOCK_H / 2] };
}

// ── The fixed head: column, top arm, gland, feeler ───────────────────────
function frame_() {
  const out = [];
  const shapes = [];
  const F = up(COLUMN.x + COLUMN.w / 2, COLUMN.y + COLUMN.d / 2, 0);
  out.push(S(G.ringBand(0, 6, 8, 11, F, P), { tone: "mid", crease: "lo" }));
  shapes.push(A.ring(F, 0, 6, 8, 11, { name: "frame.column-foot" }));
  out.push(S(k.slabOf(COLUMN, 6, COLUMN_TOP - 6, P, 6, 1.2), { tone: "mid" }));
  shapes.push(A.slab(COLUMN, 6, COLUMN_TOP - 6, { name: "frame.column" }, 6));
  out.push(line(k.sideSeam(COLUMN, 40, P, 6), { tone: "faint" }));
  // Stiffening ribs on the column's near face. They sit against the face rather
  // than in it, and clear of the handwheel's boss.
  for (const z of [22, 100]) {
    const rib = { x: COLUMN.x + COLUMN.w, y: COLUMN.y + 8, w: 5, d: COLUMN.d - 16, r: 1 };
    out.push(S(k.slabOf(rib, z, 42, P, 3, 0.4), { tone: "mid" }));
    shapes.push(A.slab(rib, z, 42, { name: `frame.rib${z}` }, 3));
  }
  out.push(D([at([COLUMN.x + 18, COLUMN.y + COLUMN.d, 22]), at([COLUMN.x + 18, COLUMN.y + COLUMN.d, COLUMN_TOP - 10])], { size: 0.45 }));
  const screws = k.corners(COLUMN, 9);
  out.push(D(screws.map(([x, y]) => at([x, y, COLUMN_TOP])), { size: 0.45 }));
  out.push(line(screws.map(([x, y]) => k.ring(x, y, 2, COLUMN_TOP, P, 14)).join(""), { tone: "lo" }));

  out.push(S(k.slabOf(ARM, ARM_Z, ARM_H, P, 6, 1), { tone: "mid" }));
  shapes.push(A.slab(ARM, ARM_Z, ARM_H, { name: "frame.arm" }, 6));
  out.push(line(k.sideSeam(ARM, ARM_Z + ARM_H - 5, P, 6), { tone: "faint" }));
  out.push(line(k.planOutline(k.insetPlan(ARM, 7), ARM_Z + ARM_H, P), { tone: "faint" }));
  const armScrews = k.corners(ARM, 12);
  out.push(line(armScrews.map(([x, y]) => k.ring(x, y, 2, ARM_Z + ARM_H, P, 14)).join(""), { tone: "lo" }));
  out.push(D(armScrews.map(([x, y]) => at([x, y, ARM_Z + ARM_H])), { size: 0.45 }));
  const rack = k.sideTicks(ARM.x + 10, ARM.x + ARM.w - 10, 8, 5, ARM.y + ARM.d, ARM_Z + ARM_H, [2.2, 4], P);
  out.push(line(rack.minor, { tone: "lo" }), line(rack.major, { tone: "mid" }));
  return { out: out.join(""), shapes, at: [COLUMN.x, COLUMN.y, COLUMN_TOP / 2] };
}

function head() {
  const out = [];
  const shapes = [];
  const [x, y] = HEAD;
  // The gland the feeler slides through, hung directly under the arm: a collar with
  // a bore, so the rod passes through it instead of into it.
  const F = down(x, y, COLLAR_Z);
  out.push(S(G.ringBand(0, COLLAR_H, BORE_R, COLLAR_R, F, P), { tone: "hi", crease: "lo" }));
  shapes.push(A.ring(F, 0, COLLAR_H, BORE_R, COLLAR_R, { name: "head.collar", owner: "head.collar" }));
  out.push(line(G.arcOf(F, COLLAR_H, COLLAR_R, P, {}), { tone: "lo" }));
  out.push(line(k.ring(x, y, COLLAR_R, COLLAR_Z, P, 32), { tone: "lo" }));
  out.push(line(k.ring(x, y, BORE_R, COLLAR_Z - COLLAR_H, P, 32), { tone: "mid" }));
  out.push(D(G.dotsOf(F, 2.5, COLLAR_R + 2.4, 6, P, { all: true }), { size: 0.45, tone: "lo" }));
  // A small lamp housing beside the bore: the read indicator.
  const lamp = { x: x + 12, y: y - 4, w: 8, d: 8, r: 2 };
  out.push(S(k.slabOf(lamp, COLLAR_Z - 5, 5, P, 3, 0.4), { tone: "mid" }));
  shapes.push(A.slab(lamp, COLLAR_Z - 5, 5, { name: "head.lamp", owner: "head" }, 3));
  return { out: out.join(""), shapes, at: [x, y, COLLAR_Z] };
}

// The feeler: one rod built at full length, its hidden part clipped away by a clip
// that stays on the gland face while the shuttle indexes. It is the figure's one
// lit part — the reading itself.
function feeler() {
  const F = down(HEAD[0], HEAD[1], COLLAR_Z);
  const rod = G.rodOut(F, COLLAR_H, ROD_END, FEELER_R, P, { stroke: LIFT });
  const tipZ = COLLAR_Z - ROD_END;
  const tipF = down(HEAD[0], HEAD[1], tipZ + TIP_H);
  const tip = G.solidOf([[0, FEELER_R], [0.6, FEELER_R], [TIP_H, 0.7]], tipF, P);
  const out = [
    `<clipPath id="iso-feeler-clip"><path data-hold d="${rod.clip}"/></clipPath>`,
    `<g clip-path="url(#iso-feeler-clip)">${S(rod.paths, { tone: "hi", lit: true, crease: "lo" })}</g>`,
    S(tip, { tone: "hi", lit: true, crease: "lo" }),
    line(k.sideArc(HEAD[0], HEAD[1], FEELER_R, tipZ + TIP_H - 0.8, P, 20), { tone: "faint" }),
  ];
  return {
    out: out.join(""),
    shapes: [A.slide(F, COLLAR_H, ROD_END, FEELER_R, { name: "head.feeler", owner: "head.feeler" }), A.solid(tipF, [[0, FEELER_R], [0.6, FEELER_R], [TIP_H, 0.7]], { name: "head.tip", owner: "head.feeler" })],
    at: [HEAD[0], HEAD[1], tipZ + TIP_H],
  };
}

// The handwheel that indexes the shuttle: the control the reader turns. It is on a
// boss off the column's near face, on the axis the reader faces.
// The index handwheel, on its own pedestal beside the rails. It lies flat, so it
// rests on its pedestal's top face the way the dial rests on Fig 1's plate — a
// disc mounted on a vertical face can never be ordered against that face.
function handwheel() {
  const [x, y] = [118, 34];
  const out = [];
  const shapes = [];
  const PED = { x: x - 15, y: y - 13, w: 30, d: 26, r: 3 };
  out.push(S(k.slabOf(PED, 0, 8, P, 5, 0.6), { tone: "mid" }));
  shapes.push(A.slab(PED, 0, 8, { name: "wheel.pedestal" }, 5));
  out.push(line(k.sideSeam(PED, 5.6, P, 5), { tone: "faint" }));
  out.push(D([at([PED.x + 5, PED.y + PED.d, 8]), at([PED.x + PED.w - 5, PED.y + PED.d, 8])], { size: 0.42 }));
  const F = up(x, y, 8);
  out.push(S(G.ringBand(0, 4, 5, 7.4, F, P), { tone: "mid", crease: "lo" }));
  shapes.push(A.ring(F, 0, 4, 5, 7.4, { name: "wheel.boss" }));
  out.push(S(k.cylinder(x, y, 5, 12, 10, P, 28), { tone: "mid", crease: "none" }));
  shapes.push(A.cylinder(x, y, 5, 12, 10, { name: "wheel.spindle" }));
  out.push(line(k.sideArc(x, y, 5, 16, P, 20), { tone: "faint" }));
  out.push(S(k.cylinder(x, y, 17, 22, 7, P, 64, 0.8), { tone: "hi", crease: "lo" }));
  shapes.push(A.cylinder(x, y, 17, 22, 7, { name: "wheel.rim" }));
  out.push(line(k.knurl(x, y, 17, 22.8, 28.2, 72, P, 0.35), { tone: "lo" }));
  out.push(line(k.sideArc(x, y, 17, 22.4, P, 40), { tone: "faint" }));
  out.push(line(k.ring(x, y, 14.6, 29, P, 64), { tone: "lo" }));
  const ticks = k.radialTicks(x, y, 14.6, 24, 3, 29, [1.6, 3.4], P);
  out.push(line(ticks.minor, { tone: "lo" }), line(ticks.major, { tone: "mid" }));
  out.push(line(k.ring(x, y, 4.4, 29, P, 24), { tone: "mid" }));
  // The knob you grip, on the wheel's top face at the rim.
  out.push(S(k.cylinder(x + 11, y, 3.4, 29, 8, P, 24, 0.5), { tone: "mid", crease: "lo" }));
  shapes.push(A.cylinder(x + 11, y, 3.4, 29, 8, { name: "wheel.knob" }));
  out.push(line(k.knurl(x + 11, y, 3.4, 29.6, 36.4, 22, P, 0.3), { tone: "lo" }));
  out.push(line(k.ring(x + 11, y, 3.4, 37, P, 24), { tone: "mid" }));
  out.push(line(k.lineOnTop([x, y], [x + 12, y], 29.2, P), { tone: "hi" }));
  return { out: out.join(""), shapes, at: [x, y, 20] };
}

// ── Assemble ─────────────────────────────────────────────────────────────
export function buildReadStation() {
  const R = A.recorder(P);
  const items = [];
  const put = (name, point, svg, shapes, bias = 0) =>
    items.push({ name, at: point, key: k.depthOf(point, P) + bias, svg, shapes: shapes.map((shape) => R.solid(shape)) });

  const parts = { base: base(), scale: scaleStrip(), label: labelPlate(), rails: rails(), shuttle: shuttle(), frame: frame_(), head: head(), feeler: feeler(), wheel: handwheel() };
  put("base.plate", [0, 0, -PLATE_H / 2], parts.base.out, parts.base.shapes, -8);
  put("base.scale", [-5, PLATE.y + PLATE.d - 10, 0.3], parts.scale.out, parts.scale.shapes, 0.2);
  put("base.label", [109, PLATE.y + PLATE.d - 15, 0.4], parts.label.out, parts.label.shapes, 0.2);
  put("rails", [0, 0, 2], parts.rails.out, parts.rails.shapes, -1);
  put("frame", parts.frame.at, parts.frame.out, parts.frame.shapes, 0.5);
  put("head", parts.head.at, parts.head.out, parts.head.shapes, 1.6);
  put("wheel", parts.wheel.at, parts.wheel.out, parts.wheel.shapes, 1.5);

  // The shuttle rides the rails and carries both blocks, so they move together:
  // the head reads whichever block the shuttle has presented.
  const shuttleParts = [{ name: "shuttle", out: parts.shuttle.out, shapes: parts.shuttle.shapes, at: parts.shuttle.at }];
  BLOCKS.forEach((_, index) => {
    const block = parameterBlock(index);
    shuttleParts.push({ name: `block${index}`, out: block.out, shapes: block.shapes, at: block.at });
  });
  items.push({
    name: "carriage",
    at: [0, 0, BLOCK_TOP / 2],
    key: k.depthOf([0, 0, BLOCK_TOP / 2], P) + 0.8,
    svg: `<g data-part="carriage" transform="${k.translateAlong(P, [0, SHUTTLE_AT[HOME], 0])}">${shuttleParts.map((part) => part.out).join("")}</g>`,
    shapes: shuttleParts.flatMap((part) => part.shapes.map((shape) => R.solid(shape))),
    move: [0, SHUTTLE_AT[SHUTTLE_AT.length - 1] - SHUTTLE_AT[0], 0],
    ride: "index",
  });
  items.push({
    name: "feeler",
    at: parts.feeler.at,
    key: k.depthOf(parts.feeler.at, P) + 1.8,
    svg: `<g data-part="feeler" transform="${k.translateAlong(P, [0, 0, 0])}">${parts.feeler.out}</g>`,
    shapes: parts.feeler.shapes.map((shape) => R.solid(shape)),
    move: [0, 0, LIFT],
    ride: "lift",
  });

  R.items = items;
  R.order = ["index", "lift"];
  R.spring = { stagger: 0.03, stiffness: 58, damping: 13.6 };
  A.settle(R, P, { states: A.explodeStates(R.order, R.spring).depth });
  items.sort((a, b) => a.key - b.key);

  const svg = k.figureSvg({
    width: WIDTH,
    height: HEIGHT,
    label: `A parameter read station. Two process parameter blocks ride a shuttle along two rails; a fixed spring feeler descends from the top arm and reads the working directory engraved on whichever block is presented under it. Both blocks are the same program, node.exe, and only the folder tells them apart.`,
    body: items.map((item) => item.svg),
  });

  const plate = k.plateHtml({
    fig: "Fig 2",
    title: "Parameter read station",
    hint: "Hover to index the shuttle",
    readout: `${BLOCKS[HOME].process} · ${BLOCKS[HOME].folder} · port ${BLOCKS[HOME].port}`,
    keys: [
      { mark: "lit", label: "The feeler: reading the working directory" },
      { mark: "raised", label: "Parameter block: one process, its folder and command line" },
      { mark: "flat", label: "Shuttle: presents one block at a time" },
      { mark: "dim", label: "Rack: where the shuttle stands" },
    ],
    caption: `Two rows the app shows side by side: the same program, told apart by the folder each one serves. The shuttle indexes a block under the fixed feeler; the feeler lifts while the work moves, then settles on the folder field. The reading is read-only — a spring probe, not a cutter.`,
    body: `<div class="iso iso-stage" data-figure="readstation" tabindex="0" role="group" aria-label="Parameter read station. Move the pointer or use the arrow keys to index each process block and read its folder.">${svg}<span class="iso-sr" role="status" aria-live="polite" data-status></span></div>`,
  });

  const live = `
const BLOCKS = ${JSON.stringify(BLOCKS.map((block, index) => ({ ...block, at: SHUTTLE_AT[index] })))};
const P = ${JSON.stringify(P)};
const HOME = ${HOME}, LIFT = ${LIFT}, BASE = ${SHUTTLE_AT[HOME]};
const TRAVEL = ${JSON.stringify(SHUTTLE_AT)};
const STOPS = TRAVEL;
const clamp = (value, low, high) => Math.min(high, Math.max(low, value));

function mount(root) {
  const stage = root.querySelector(".iso-stage") || root;
  const svg = stage.querySelector("svg");
  const carriage = svg.querySelector('[data-part="carriage"]');
  const feeler = svg.querySelector('[data-part="feeler"]');
  // This figure's own readout: three figures share one page, and the first one in
  // the document is not necessarily this one.
  const readout = (stage.closest(".iso-figure") || document).querySelector("[data-readout]");
  const status = stage.querySelector("[data-status]");
  const still = matchMedia("(prefers-reduced-motion: reduce)");
  let y = STOPS[HOME], target = y, lift = 0, rate = 0;
  let frame = 0, last = 0, clock = 0, touring = true, visible = false, idleTimer = 0;

  const nearest = (value) => STOPS.reduce((best, stop, index) => (Math.abs(stop - value) < Math.abs(STOPS[best] - value) ? index : best), 0);

  function draw() {
    carriage.setAttribute("transform", translateAlong(P, [0, y, 0]));
    feeler.setAttribute("transform", translateAlong(P, [0, 0, lift]));
    holdStill(feeler, axisVector(P, [0, 0, lift]));
    const index = nearest(y);
    const block = BLOCKS[index];
    const text = block.process + " \\u00b7 " + block.folder + " \\u00b7 port " + block.port;
    if (readout.textContent !== text) readout.textContent = text;
    const said = "Block " + (index + 1) + " of " + BLOCKS.length + ": " + block.process + ", serving " + block.folder + ", on port " + block.port;
    if (status.textContent !== said) status.textContent = said;
  }

  function tick(now) {
    const dt = last ? Math.min((now - last) / 1000, 1 / 30) : 1 / 60;
    last = now;
    const calm = still.matches;
    if (touring) {
      clock += dt;
      target = STOPS[Math.floor(clock / 2.6) % STOPS.length];
    }
    y = calm ? target : y + (target - y) * (1 - Math.exp(-dt / 0.24));
    if (Math.abs(target - y) < 0.02) y = target;
    // The feeler lifts while the work is still moving, so it never drags across a
    // block face; it settles back onto the field once the shuttle has arrived.
    const want = LIFT * Math.min(1, Math.abs(target - y) / 8);
    const stiff = 420, drag = calm ? 2 * Math.sqrt(420) : 24;
    for (let i = 0; i < 4; i++) { rate += ((want - lift) * stiff - rate * drag) * (dt / 4); lift += rate * (dt / 4); }
    if (Math.abs(want - lift) < 0.004 && Math.abs(rate) < 0.004) { lift = want; rate = 0; }
    draw();
    const settled = !touring && y === target && rate === 0 && lift === want;
    frame = visible && !settled ? requestAnimationFrame(tick) : 0;
  }

  const run = () => { if (!frame && visible) { last = 0; frame = requestAnimationFrame(tick); } };
  const steer = (index) => {
    touring = false;
    target = clamp(STOPS[clamp(index, 0, STOPS.length - 1)], Math.min(...STOPS), Math.max(...STOPS));
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
    const share = clamp((event.clientY - box.top) / box.height, 0, 1);
    steer(share < 0.5 ? 0 : 1);
  });
  stage.addEventListener("pointerleave", () => steer(nearest(target)));
  stage.addEventListener("keydown", (event) => {
    const here = nearest(target);
    let to = null;
    if (event.key === "ArrowRight" || event.key === "ArrowDown") to = Math.min(STOPS.length - 1, here + 1);
    if (event.key === "ArrowLeft" || event.key === "ArrowUp") to = Math.max(0, here - 1);
    if (event.key === "Home") to = 0;
    if (event.key === "End") to = STOPS.length - 1;
    if (to === null) return;
    event.preventDefault();
    steer(to);
  });
  draw();
}
`;
  return { plate, live: `(function(){\n${live}\nmount(document.querySelector('[data-figure="readstation"]'));\n})();`, R, P, solids: R.solids.length, items: items.length };
}
















