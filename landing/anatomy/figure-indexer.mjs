// Fig 1 — the port indexer. One socket per listening port, each post as tall as
// its port number on a 0-65535 scale; a wide collar marks a port another machine
// can reach. A head on a gantry sweeps the row and settles a spring probe on each
// socket while a dial reads the number it has found.
//
// Every number comes from the app's own model: Port.Number sets a post's height,
// Port.Exposed (anything but a loopback bind) sets the collar, and the head sweeps
// the row the way the app rescans every two seconds.
import * as k from "./kit/iso-kit.mjs";
import * as G from "./kit/lathe.mjs";
import * as A from "./kit/audit.mjs";
import { at as atOf, byDepth as byDepthOf, line, plateStyle, dots } from "./figure-theme.mjs";

const WIDTH = 560;
const HEIGHT = 470;
const PORT_MAX = 65535;

// Sorted by port, as the app lists them. The first five are the README's rows;
// the rest are ports a Node project actually opens.
export const ROWS = [
  { port: 3000, name: "node.exe", folder: "portfolio", addr: "127.0.0.1", pid: 14892 },
  { port: 5173, name: "node.exe", folder: "waku", addr: "127.0.0.1", pid: 17152 },
  { port: 5432, name: "postgres.exe", folder: "portfolio", addr: "127.0.0.1", pid: 14892 },
  { port: 8644, name: "python.exe", folder: "hermes", addr: "::", pid: 21740 },
  { port: 8888, name: "python.exe", folder: "hermes", addr: "127.0.0.1", pid: 21740 },
  { port: 34123, name: "waku-daemon", folder: "waku", addr: "0.0.0.0", pid: 17152 },
  { port: 51467, name: "opencode.exe", folder: "portfolio", addr: "127.0.0.1", pid: 32920 },
  { port: 63502, name: "ABDownload.exe", folder: "system32", addr: "127.0.0.1", pid: 2392 },
];

// A loopback bind is this machine only; anything else another machine can reach.
const isLoopback = (addr) => addr === "::1" || addr === "localhost" || addr.startsWith("127.");
export const SOCKETS = ROWS.map((row, index) => ({ ...row, index, exposed: !isLoopback(row.addr) }));
export const EXPOSED = SOCKETS.filter((socket) => socket.exposed).length;

const PLATE = { x: -150, y: -76, w: 300, d: 152, r: 10 };
const PLATE_H = 12;
const FOOT_H = 5;
const ROW_Y = 26;
const SOCKET_STEP = 32;
const socketX = (index) => -110 + SOCKET_STEP * index;
const PAD_R = 9;
const POST_R = 4.4;
const CAP_R = 5.2;
const COLLAR_OUT = 11.6;
const POST_X = [-130, 120];
const POST_Y = -40;
const POST_R_ = 3.6;
const BEAM_Z = 66;
const BEAM_H = 10;
const BEAM_D = 14;
const BEAM_W = 268;
const CARRIAGE_Z = BEAM_Z + BEAM_H;
const CARRIAGE_H = 12;
const ARM_Z = CARRIAGE_Z;
const ARM_H = 8;
const ARM_D = 54;
const KNOB_R = 5;
const COLLAR_Z = ARM_Z;
const COLLAR_H = 5;
const BORE_R = 2.6;
const COLLAR_R = 5.2;
const PROBE_R = 2.2;
const TIP_H = 3;
const SAFE_LIFT = 20;
const DIAL = [150, -66];
const DIAL_R = 18;
const DIAL_Z = 70;
const DIAL_H = 10;
const DIAL_TOP = DIAL_Z + DIAL_H;
const REST = 0;

// The row's height band. A linear 0-65535 scale would leave the five ports a dev
// machine really uses 2-4 units tall, which reads as a flat, bland row; the band
// spreads the low end while keeping the port number the only thing that sets a
// socket's height.
const POST_Z = 6;
const POST_SPAN = 46;
const postHeight = (port) => POST_Z + POST_SPAN * (port / PORT_MAX);
const tops = SOCKETS.map((socket) => postHeight(socket.port));
SOCKETS.forEach((socket, index) => (socket.top = tops[index]));
const REST_X = socketX(REST);
const REST_TOP = tops[REST];
const COLLAR_BOTTOM = COLLAR_Z - COLLAR_H;
const TALLEST = Math.max(...tops);
const ROD_END = COLLAR_Z - REST_TOP;
const RETRACT = TALLEST - REST_TOP;
const SWEEP = socketX(SOCKETS.length - 1) - REST_X;

const P = k.fitProjection(
  [
    ...k.boxCorners(PLATE, -PLATE_H - FOOT_H, 0),
    [POST_X[0], POST_Y, BEAM_Z + BEAM_H],
    [POST_X[1], POST_Y, BEAM_Z + BEAM_H],
    [REST_X, POST_Y, CARRIAGE_Z + CARRIAGE_H + KNOB_R * 2],
    [REST_X, ROW_Y, ARM_Z + ARM_H],
    [DIAL[0] + DIAL_R, DIAL[1], DIAL_TOP],
    [DIAL[0] - DIAL_R, DIAL[1], DIAL_TOP],
    [socketX(SOCKETS.length - 1) + COLLAR_OUT, ROW_Y, 0],
    [socketX(0) - COLLAR_OUT, ROW_Y, 0],
    [socketX(SOCKETS.length - 1), ROW_Y, TALLEST],
    [socketX(0), ROW_Y, TALLEST],
  ],
  WIDTH,
  HEIGHT,
  { pad: 62, azimuth: 45 },
);

const at = atOf(P);
const byDepth = byDepthOf(P);
const S = plateStyle;
const D = dots;
const up = (x, y, z) => G.frameAlong([x, y, z], [0, 0, 1]);
const down = (x, y, z) => G.frameAlong([x, y, z], [0, 0, -1]);

// The recorder and the item list are built fresh per call: the build script calls
// this once, but a module-scope list would silently accumulate if it ever didn't.

// ── Base: plate on feet, graduated scale, maker's plate ──────────────────
function base() {
  const out = [`<path class="iso-halo" d="${k.haloOf(PLATE, -PLATE_H - FOOT_H, PLATE_H, P)}"/>`];
  const shapes = [];
  for (const [x, y] of byDepth(k.corners(PLATE, 26))) {
    out.push(S(k.cylinder(x, y, 7, -PLATE_H - FOOT_H, 1.6, P, 28), { tone: "lo" }));
    out.push(S(k.cylinder(x, y, 5.4, -PLATE_H - FOOT_H + 1.6, FOOT_H - 1.6, P, 28), { tone: "mid" }));
    shapes.push(A.cylinder(x, y, 7, -PLATE_H - FOOT_H, 1.6, { name: `base.foot${shapes.length}` }), A.cylinder(x, y, 5.4, -PLATE_H - FOOT_H + 1.6, FOOT_H - 1.6, { name: `base.foot${shapes.length}b` }));
  }
  out.push(S(k.slabOf(PLATE, -PLATE_H, PLATE_H, P, 10, 1.6), { tone: "mid" }));
  shapes.push(A.slab(PLATE, -PLATE_H, PLATE_H, { name: "base.plate", ground: true }, 10));
  out.push(line(k.planOutline(k.insetPlan(PLATE, 6), 0, P), { tone: "lo" }));
  const screws = k.corners(PLATE, 12);
  out.push(line(screws.map(([x, y]) => k.ring(x, y, 2.2, 0, P, 16)).join(""), { tone: "lo" }));
  out.push(D(screws.map(([x, y]) => at([x, y, 0])), { size: 0.5 }));
  return { out: out.join(""), shapes };
}

function scaleStrip() {
  const strip = { x: -100, y: PLATE.y + PLATE.d - 14, w: 140, d: 7, r: 1 };
  const out = [S(k.slabOf(strip, 0, 0.6, P, 3), { tone: "lo", crease: "none" })];
  const ticks = k.topTicks(strip.x + 4, strip.x + strip.w - 4, 5, 4, strip.y + 0.8, 0.6, [2, 3.6], P);
  out.push(line(ticks.minor, { tone: "lo" }), line(ticks.major, { tone: "mid" }));
  out.push(D([at([strip.x + 2, strip.y + 3.5, 0.6]), at([strip.x + strip.w - 2, strip.y + 3.5, 0.6])], { size: 0.45 }));
  return { out: out.join(""), shapes: [A.slab(strip, 0, 0.6, { name: "base.scale" }, 3)] };
}

// The maker's plate. Its rule lines stand for printed text, so they are the one
// place a line may stop on a bare face.
function labelPlate() {
  const plan = { x: 60, y: PLATE.y + PLATE.d - 22, w: 70, d: 16, r: 2 };
  const out = [S(k.slabOf(plan, 0, 0.8, P, 4), { tone: "lo", crease: "none" }), line(k.planOutline(k.insetPlan(plan, 1.8), 0.8, P), { tone: "faint" })];
  for (const length of [44, 30, 20]) {
    out.push(line(k.lineOnTop([plan.x + 5, plan.y + plan.d / 2], [plan.x + 5 + length, plan.y + plan.d / 2], 0.8, P), { tone: "lo", free: true }));
  }
  out.push(D([at([plan.x + 4, plan.y + 4, 0.8]), at([plan.x + plan.w - 4, plan.y + plan.d - 4, 0.8])], { size: 0.45 }));
  return { out: out.join(""), shapes: [A.slab(plan, 0, 0.8, { name: "base.plate-label" }, 4)] };
}

// ── The socket row ───────────────────────────────────────────────────────
function socket(socket_) {
  const x = socketX(socket_.index);
  const out = [];
  const shapes = [];
  out.push(S(k.cylinder(x, ROW_Y, PAD_R, 0, POST_Z, P, 40), { tone: "lo" }));
  shapes.push(A.cylinder(x, ROW_Y, PAD_R, 0, POST_Z, { name: `row.${socket_.port}.pad` }));
  out.push(line(k.sideArc(x, ROW_Y, PAD_R, POST_Z, P, 40), { tone: "faint" }));

  if (socket_.exposed) {
    // A wide collar: another machine can reach this port.
    const F = up(x, ROW_Y, 0);
    out.push(S(G.ringBand(0, 1.8, PAD_R, COLLAR_OUT, F, P), { tone: "mid", crease: "lo" }));
    shapes.push(A.ring(F, 0, 1.8, PAD_R, COLLAR_OUT, { name: `row.${socket_.port}.collar` }));
    out.push(line(G.arcOf(F, 1.8, COLLAR_OUT, P, {}), { tone: "lo" }));
    out.push(D(G.dotsOf(F, 1.8, (PAD_R + COLLAR_OUT) / 2, 8, P, { all: true }), { size: 0.45, tone: "lo" }));
  }

  const height = postHeight(socket_.port) - POST_Z;
  out.push(S(k.cylinder(x, ROW_Y, POST_R, POST_Z, height, P, 36), { tone: "hi", crease: "lo" }));
  shapes.push(A.cylinder(x, ROW_Y, POST_R, POST_Z, height, { name: `row.${socket_.port}.post` }));
  out.push(line(k.sideArc(x, ROW_Y, POST_R, POST_Z + 0.8, P, 24), { tone: "faint" }));

  const F = up(x, ROW_Y, socket_.top - 1.6);
  out.push(S(G.ringBand(0, 1.6, POST_R, CAP_R, F, P), { tone: "mid", crease: "lo" }));
  shapes.push(A.ring(F, 0, 1.6, POST_R, CAP_R, { name: `row.${socket_.port}.cap` }));
  out.push(line(k.ring(x, ROW_Y, CAP_R, socket_.top, P, 36), { tone: "mid" }));
  out.push(line(k.ring(x, ROW_Y, POST_R, socket_.top, P, 36), { tone: "lo" }));
  return { out: out.join(""), shapes, at: [x, ROW_Y, socket_.top / 2] };
}

// ── The gantry: two posts, a beam with a rack, end stops ─────────────────
function gantry() {
  const out = [];
  const shapes = [];
  for (const x of POST_X) {
    const F = up(x, POST_Y, 0);
    out.push(S(G.ringBand(0, 5, POST_R_, 5.5, F, P), { tone: "mid", crease: "lo" }));
    shapes.push(A.ring(F, 0, 5, POST_R_, 5.5, { name: `gantry.post-base${x}` }));
    out.push(S(k.cylinder(x, POST_Y, POST_R_, 5, BEAM_Z - 5, P, 32), { tone: "mid", crease: "none" }));
    shapes.push(A.cylinder(x, POST_Y, POST_R_, 5, BEAM_Z - 5, { name: `gantry.post${x}` }));
    out.push(line([0.25, 0.5, 0.75].map((share) => k.sideArc(x, POST_Y, POST_R_, 5 + (BEAM_Z - 10) * share, P, 16)).join(""), { tone: "faint" }));
    const collar = up(x, POST_Y, BEAM_Z - 5);
    out.push(S(G.ringBand(0, 5, POST_R_, 4.6, collar, P), { tone: "mid", crease: "lo" }));
    shapes.push(A.ring(collar, 0, 5, POST_R_, 4.6, { name: `gantry.post-collar${x}` }));
  }

  const beam = { x: POST_X[0] - 6, y: POST_Y - BEAM_D / 2, w: BEAM_W, d: BEAM_D, r: 2 };
  out.push(S(k.slabOf(beam, BEAM_Z, BEAM_H, P, 6, 0.8), { tone: "mid" }));
  shapes.push(A.slab(beam, BEAM_Z, BEAM_H, { name: "gantry.beam" }, 6));
  out.push(line(k.sideSeam(beam, BEAM_Z + BEAM_H - 3, P, 6), { tone: "faint" }));
  const rack = k.sideTicks(beam.x + 8, beam.x + beam.w - 8, 7, 5, beam.y + beam.d, BEAM_Z + BEAM_H, [2.4, 4.4], P);
  out.push(line(rack.minor, { tone: "lo" }), line(rack.major, { tone: "mid" }));
  for (const x of [beam.x - 4, beam.x + beam.w]) {
    out.push(S(k.slabOf({ x, y: beam.y - 2, w: 4, d: beam.d + 4, r: 1 }, BEAM_Z - 1, BEAM_H + 3, P, 3, 0.5), { tone: "mid" }));
    shapes.push(A.slab({ x, y: beam.y - 2, w: 4, d: beam.d + 4, r: 1 }, BEAM_Z - 1, BEAM_H + 3, { name: `gantry.stop${x}` }, 3));
  }
  out.push(D([at([beam.x + 6, beam.y + 3, BEAM_Z + BEAM_H]), at([beam.x + beam.w - 6, beam.y + 3, BEAM_Z + BEAM_H])], { size: 0.45 }));
  return { out: out.join(""), shapes, at: [0, POST_Y, BEAM_Z + BEAM_H / 2] };
}

// ── The head: block, index knob, arm, gland collar ───────────────────────
function head() {
  const out = [];
  const shapes = [];
  const block = { x: REST_X - 14, y: POST_Y + BEAM_D / 2, w: 28, d: 24, r: 2 };
  out.push(S(k.slabOf(block, CARRIAGE_Z, CARRIAGE_H, P, 5, 0.7), { tone: "mid" }));
  shapes.push(A.slab(block, CARRIAGE_Z, CARRIAGE_H, { name: "head.block", owner: "head" }, 5));
  out.push(line(k.sideSeam(block, CARRIAGE_Z + CARRIAGE_H - 3, P, 5), { tone: "faint" }));
  out.push(D([-10, -3, 4, 11].map((dx) => at([REST_X + dx, block.y + block.d, CARRIAGE_Z + CARRIAGE_H])), { size: 0.4, tone: "lo" }));
  out.push(line([k.ring(REST_X - 8, block.y + 6, 1.6, CARRIAGE_Z + CARRIAGE_H, P, 14), k.ring(REST_X + 8, block.y + 6, 1.6, CARRIAGE_Z + CARRIAGE_H, P, 14)].join(""), { tone: "lo" }));

  // The index knob, and the pointer that reads the rack behind it.
  const knob = [REST_X, block.y + 6];
  out.push(S(k.cylinder(knob[0], knob[1], KNOB_R, CARRIAGE_Z + CARRIAGE_H, 6, P, 32), { tone: "mid" }));
  shapes.push(A.cylinder(knob[0], knob[1], KNOB_R, CARRIAGE_Z + CARRIAGE_H, 6, { name: "head.knob", owner: "head" }));
  out.push(line(k.knurl(knob[0], knob[1], KNOB_R, CARRIAGE_Z + CARRIAGE_H, CARRIAGE_Z + CARRIAGE_H + 6, 36, P), { tone: "lo" }));
  out.push(line(k.lineOnTop([knob[0], knob[1] + 1], [knob[0] + 4.6, knob[1] + 1], CARRIAGE_Z + CARRIAGE_H + 6, P), { tone: "hi" }));
  out.push(line(k.lineOnTop([REST_X, block.y - 1], [REST_X, block.y - 9], CARRIAGE_Z + CARRIAGE_H - 4, P), { tone: "hi", free: "end" }));

  // The arm reaches forward over the row, with a rib under its root.
  const arm = { x: REST_X - 8, y: POST_Y + 12, w: 16, d: ARM_D, r: 2 };
  out.push(S(k.slabOf(arm, ARM_Z, ARM_H, P, 5, 0.7), { tone: "mid" }));
  shapes.push(A.slab(arm, ARM_Z, ARM_H, { name: "head.arm", owner: "head" }, 5));
  out.push(line(k.sideSeam(arm, ARM_Z + ARM_H - 2.6, P, 5), { tone: "faint" }));
  out.push(D([ROW_Y - 18, ROW_Y - 4, ROW_Y + 10].map((y) => at([REST_X, y, ARM_Z + ARM_H])), { size: 0.4, tone: "lo" }));
  const rib = { x: REST_X - 4, y: POST_Y + 12, w: 8, d: 20, r: 1.5 };
  out.push(S(k.slabOf(rib, ARM_Z - 8, 8, P, 4, 0.5), { tone: "mid" }));
  shapes.push(A.slab(rib, ARM_Z - 8, 8, { name: "head.rib", owner: "head" }, 4));

  // The gland the probe slides through: a real collar with a bore, so the rod
  // passes through it instead of into it.
  const F = down(REST_X, ROW_Y, COLLAR_Z);
  out.push(S(G.ringBand(0, COLLAR_H, BORE_R, COLLAR_R, F, P), { tone: "hi", crease: "lo" }));
  shapes.push(A.ring(F, 0, COLLAR_H, BORE_R, COLLAR_R, { name: "head.collar", owner: "head" }));
  out.push(line(G.arcOf(F, COLLAR_H, COLLAR_R, P, {}), { tone: "lo" }));
  out.push(line(k.ring(REST_X, ROW_Y, COLLAR_R, COLLAR_Z, P, 32), { tone: "lo" }));
  out.push(line(k.ring(REST_X, ROW_Y, BORE_R, COLLAR_Z - COLLAR_H, P, 32), { tone: "mid" }));
  return { out: out.join(""), shapes, at: [REST_X, ROW_Y, ARM_Z] };
}

// The probe is one rod built at full length; the part inside the gland is clipped
// away by a clip that stays put while the head moves (craft.md, "Rod out of a
// gland"). It is the figure's one lit part: the port being read.
function probe() {
  const F = down(REST_X, ROW_Y, COLLAR_Z);
  const rod = G.rodOut(F, COLLAR_H, ROD_END, PROBE_R, P, { stroke: RETRACT });
  const tipZ = COLLAR_Z - ROD_END;
  const tipF = down(REST_X, ROW_Y, tipZ + TIP_H);
  const tip = G.solidOf([[0, PROBE_R], [0.6, PROBE_R], [TIP_H, 0.9]], tipF, P);
  const out = [
    `<clipPath id="iso-probe-clip"><path data-hold d="${rod.clip}"/></clipPath>`,
    `<g clip-path="url(#iso-probe-clip)">${S(rod.paths, { tone: "hi", lit: true, crease: "lo" })}</g>`,
    S(tip, { tone: "hi", lit: true, crease: "lo" }),
    line(k.sideArc(REST_X, ROW_Y, PROBE_R, tipZ + TIP_H - 0.8, P, 20), { tone: "faint" }),
  ];
  return {
    out: out.join(""),
    shapes: [A.slide(F, COLLAR_H, ROD_END, PROBE_R, { name: "head.rod", owner: "head" }), A.solid(tipF, [[0, PROBE_R], [0.6, PROBE_R], [TIP_H, 0.9]], { name: "head.tip", owner: "head" })],
    at: [REST_X, ROW_Y, tipZ + TIP_H],
  };
}

// ── The dial: the same 0-65535 scale, read as an angle ───────────────────
const needleAt = (port) => {
  const angle = (port / PORT_MAX) * Math.PI * 2;
  const [x, y] = DIAL;
  return k.lineOnTop([x - Math.sin(angle) * 3.4, y + Math.cos(angle) * 3.4], [x + Math.sin(angle) * (DIAL_R - 4), y - Math.cos(angle) * (DIAL_R - 4)], DIAL_TOP + 0.2, P);
};

function dial() {
  const [x, y] = DIAL;
  const out = [];
  const shapes = [];
  const F = up(x, y, 0);
  out.push(S(G.ringBand(0, 6, 4.6, 7.4, F, P), { tone: "mid", crease: "lo" }));
  shapes.push(A.ring(F, 0, 6, 4.6, 7.4, { name: "dial.foot" }));
  out.push(S(k.cylinder(x, y, 4.6, 6, DIAL_Z - 6, P, 28), { tone: "mid", crease: "none" }));
  shapes.push(A.cylinder(x, y, 4.6, 6, DIAL_Z - 6, { name: "dial.column" }));
  out.push(line([24, 44, 62].map((z) => k.sideArc(x, y, 4.6, z, P, 16)).join(""), { tone: "faint" }));
  out.push(S(k.cylinder(x, y, DIAL_R, DIAL_Z, DIAL_H, P, 64, 0.8), { tone: "hi", crease: "lo" }));
  shapes.push(A.cylinder(x, y, DIAL_R, DIAL_Z, DIAL_H, { name: "dial.puck" }));
  out.push(line(k.knurl(x, y, DIAL_R, DIAL_TOP - 3.4, DIAL_TOP, 84, P, 0.3), { tone: "lo" }));
  out.push(line(k.sideArc(x, y, DIAL_R, DIAL_TOP - 3.4, P, 40), { tone: "faint" }));
  out.push(line(k.ring(x, y, DIAL_R - 2.2, DIAL_TOP, P, 64), { tone: "lo" }));
  const ticks = k.radialTicks(x, y, DIAL_R - 3, 40, 5, DIAL_TOP, [1.3, 2.8], P);
  out.push(line(ticks.minor, { tone: "lo" }), line(ticks.major, { tone: "mid" }));
  out.push(line(k.ring(x, y, 1.6, DIAL_TOP, P, 16), { tone: "mid" }));
  const crown = { x: x + DIAL_R, y: y - 4, w: 6, d: 8, r: 1.5 };
  out.push(S(k.slabOf(crown, DIAL_Z + 2, DIAL_H - 4, P, 3, 0.4), { tone: "mid" }));
  shapes.push(A.slab(crown, DIAL_Z + 2, DIAL_H - 4, { name: "dial.crown" }, 3));
  out.push(`<path data-part="needle" class="iso-line" data-tone="lit" data-free="start" d="${needleAt(ROWS[REST].port)}" style="stroke-width:1.2"/>`);
  return { out: out.join(""), shapes, at: [x, y, DIAL_TOP] };
}

// ── Assemble ─────────────────────────────────────────────────────────────
export function buildIndexer() {
  const R = A.recorder(P);
  const items = [];
  const put = (name, point, svg, shapes, bias = 0) =>
    items.push({ name, at: point, key: k.depthOf(point, P) + bias, svg, shapes: shapes.map((shape) => R.solid(shape)) });

  const parts = { base: base(), scale: scaleStrip(), label: labelPlate(), gantry: gantry(), head: head(), probe: probe(), dial: dial() };
  put("base.plate", [0, 0, -PLATE_H / 2], parts.base.out, parts.base.shapes, -6);
  put("base.scale", [-30, PLATE.y + PLATE.d - 10, 0.3], parts.scale.out, parts.scale.shapes, 0.2);
  put("base.label", [95, PLATE.y + PLATE.d - 14, 0.4], parts.label.out, parts.label.shapes, 0.2);
  put("gantry", parts.gantry.at, parts.gantry.out, parts.gantry.shapes, 0.4);

  const socketParts = SOCKETS.map((socket_) => ({ socket_, ...socket(socket_) }));
  socketParts.forEach(({ socket_, out, shapes, at: point }) => put(`row.${socket_.port}`, point, out, shapes, socket_.exposed ? 0.6 : 0.3));

  // The head carries the probe; they move independently — the head sweeps the row,
  // the probe retracts into its gland as the head meets a taller socket.
  items.push({ name: "head", at: parts.head.at, key: k.depthOf(parts.head.at, P) + 0.9, svg: `<g data-part="carriage" transform="${k.translateAlong(P, [0, 0, 0])}">${parts.head.out}</g>`, shapes: parts.head.shapes.map((shape) => R.solid(shape)), move: [SWEEP, 0, 0], ride: "sweep" });
  items.push({ name: "probe", at: [REST_X, ROW_Y, (COLLAR_BOTTOM + REST_TOP) / 2], key: k.depthOf([REST_X, ROW_Y, (COLLAR_BOTTOM + REST_TOP) / 2], P) + 1.1, svg: `<g data-part="probe" transform="${k.translateAlong(P, [0, 0, 0])}">${parts.probe.out}</g>`, shapes: parts.probe.shapes.map((shape) => R.solid(shape)), move: [SWEEP, 0, RETRACT], ride: "probe" });
  put("dial", parts.dial.at, parts.dial.out, parts.dial.shapes, 1.2);

  R.items = items;
  R.order = ["sweep", "probe"];
  R.spring = { stagger: 0.02, stiffness: 64, damping: 14.2 };
  A.settle(R, P, { states: A.explodeStates(R.order, R.spring).depth });
  items.sort((a, b) => a.key - b.key);

  const svg = k.figureSvg({
    width: WIDTH,
    height: HEIGHT,
    label: `A bench indexer for the listening ports on a Windows machine. Eight sockets stand in a row, one per port: each post is as tall as its port number on a 0 to 65535 scale, and the ${EXPOSED} ports another machine can reach sit in wide collars. A head on a gantry sweeps the row, lifting between sockets and settling a spring probe on each, while a dial reads the port number it has found.`,
    body: items.map((item) => item.svg),
  });

  const plate = k.plateHtml({
    fig: "Fig 1",
    title: "Port indexer",
    hint: "Hover to index the row",
    readout: `${ROWS[REST].port} · ${ROWS[REST].name} · ${ROWS[REST].folder} · ${SOCKETS[REST].exposed ? "exposed" : "this machine only"}`,
    keys: [
      { mark: "lit", label: "The probe: the port being read" },
      { mark: "raised", label: "Wide collar: another machine can reach it" },
      { mark: "flat", label: "Socket: one listening port" },
      { mark: "dim", label: "Dial: the port number, 0 to 65535" },
    ],
    caption: `An example scan: ${ROWS.length} listeners, ${EXPOSED} of them reachable from another machine. Socket height and the needle both read the port number on one 0-65535 scale. The head sweeps the row the way the app rescans every two seconds, lifting between sockets and settling on each.`,
    body: `<div class="iso iso-stage" data-wide data-figure="indexer" tabindex="0" role="group" aria-label="Port indexer. Move the pointer across the row to read each listening port, or use the arrow keys.">${svg}<span class="iso-sr" role="status" aria-live="polite" data-status></span></div>`,
  });

  const live = `
const ROWS = ${JSON.stringify(SOCKETS.map((socket_) => ({ port: socket_.port, name: socket_.name, folder: socket_.folder, addr: socket_.addr, exposed: socket_.exposed, x: socketX(socket_.index), top: socket_.top, fx: at([socketX(socket_.index), ROW_Y, 0])[0] / WIDTH })))};
const P = ${JSON.stringify(P)};
const DIAL = ${JSON.stringify(DIAL)};
const DIAL_R = ${DIAL_R}, DIAL_TOP = ${DIAL_TOP}, PORT_MAX = ${PORT_MAX};
const REST_X = ${REST_X}, REST_TOP = ${REST_TOP}, SAFE_LIFT = ${SAFE_LIFT};
const REACH = [ROWS[0].x, ROWS[ROWS.length - 1].x];
const needleAt = ${needleAt.toString().replace("k.lineOnTop", "lineOnTop")};

function mount(root) {
  const stage = root.querySelector(".iso-stage") || root;
  const svg = stage.querySelector("svg");
  const carriage = svg.querySelector('[data-part="carriage"]');
  const probe = svg.querySelector('[data-part="probe"]');
  const needle = svg.querySelector('[data-part="needle"]');
  // This figure's own readout: three figures share one page, and the first one in
  // the document is not necessarily this one.
  const readout = (stage.closest(".iso-figure") || document).querySelector("[data-readout]");
  const status = stage.querySelector("[data-status]");
  const still = matchMedia("(prefers-reduced-motion: reduce)");
  const clamp = (value, low, high) => Math.min(high, Math.max(low, value));
  const rowAt = (x) => clamp((x - ROWS[0].x) / (ROWS[1].x - ROWS[0].x), 0, ROWS.length - 1);
  const nearest = (x) => ROWS[Math.round(rowAt(x))];
  let x = REST_X, target = REST_X, probeZ = REST_TOP, rate = 0;
  let frame = 0, last = 0, clock = 0, touring = true, visible = false, idleTimer = 0;

  const rowAtScreen = (fx) => {
    let best = 0;
    for (let index = 1; index < ROWS.length; index++) if (Math.abs(ROWS[index].fx - fx) < Math.abs(ROWS[best].fx - fx)) best = index;
    return best;
  };

  function draw() {
    const dx = x - REST_X, dz = probeZ - REST_TOP;
    carriage.setAttribute("transform", translateAlong(P, [dx, 0, 0]));
    probe.setAttribute("transform", translateAlong(P, [dx, 0, dz]));
    // The clip must stay on the gland, which moves with the head only, so it is
    // held against the probe's own vertical travel and nothing else.
    holdStill(probe, axisVector(P, [0, 0, dz]));
    const row = nearest(x);
    needle.setAttribute("d", needleAt(row.port));
    const text = row.port + " \\u00b7 " + row.name + " \\u00b7 " + row.folder + " \\u00b7 " + (row.exposed ? "exposed" : "this machine only");
    if (readout.textContent !== text) readout.textContent = text;
    const said = "Port " + row.port + ", " + row.name + ", in " + row.folder + (row.exposed ? ", reachable from another machine" : ", this machine only");
    if (status.textContent !== said) status.textContent = said;
  }

  function tick(now) {
    const dt = last ? Math.min((now - last) / 1000, 1 / 30) : 1 / 60;
    last = now;
    const calm = still.matches;
    if (touring) {
      clock += dt;
      target = REACH[0] + ((REACH[1] - REACH[0]) * ((clock / 1.5) % ROWS.length)) / (ROWS.length - 1);
    }
    x = calm ? target : x + (target - x) * (1 - Math.exp(-dt / 0.26));
    if (Math.abs(target - x) < 0.02) x = target;
    // The probe rides high while the head still has ground to cover, so it lifts
    // off one socket and drops onto the next.
    const under = nearest(x);
    const want = under.top + SAFE_LIFT * Math.min(1, (Math.abs(target - x) / (REACH[1] - REACH[0])) * ((ROWS.length - 1) / 0.3));
    const stiff = 460, drag = calm ? 2 * Math.sqrt(460) : 27;
    for (let i = 0; i < 4; i++) { rate += ((want - probeZ) * stiff - rate * drag) * (dt / 4); probeZ += rate * (dt / 4); }
    if (Math.abs(want - probeZ) < 0.004 && Math.abs(rate) < 0.004) { probeZ = want; rate = 0; }
    draw();
    const settled = !touring && x === target && rate === 0 && probeZ === want;
    frame = visible && !settled ? requestAnimationFrame(tick) : 0;
  }

  const run = () => { if (!frame && visible) { last = 0; frame = requestAnimationFrame(tick); } };
  const steer = (next) => {
    touring = false;
    target = clamp(next, REACH[0], REACH[1]);
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
    const fx = (event.clientX - box.left) / box.width;
    // Snap to the socket the pointer is nearest, so the head settles on a
    // listener rather than stopping between two.
    steer(ROWS[rowAtScreen(clamp(fx, ROWS[0].fx, ROWS[ROWS.length - 1].fx))].x);
  });
  stage.addEventListener("pointerleave", () => steer(target));
  stage.addEventListener("keydown", (event) => {
    const here = Math.round(rowAt(target));
    let to = null;
    if (event.key === "ArrowRight" || event.key === "ArrowDown") to = Math.min(ROWS.length - 1, here + 1);
    if (event.key === "ArrowLeft" || event.key === "ArrowUp") to = Math.max(0, here - 1);
    if (event.key === "Home") to = 0;
    if (event.key === "End") to = ROWS.length - 1;
    if (to === null) return;
    event.preventDefault();
    steer(ROWS[to].x);
  });
  draw();
}
`;
  return { plate, live: `(function(){\n${live}\nmount(document.querySelector('[data-figure="indexer"]'));\n})();`, solids: R.solids.length, items: items.length, R, P };
}





