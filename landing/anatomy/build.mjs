// The landing page's build. Every figure is generated from the app's own model and
// audited before the page is written, so the drawing and the page cannot drift.
//
//   node build.mjs            writes index.html, figure.css and figure.js
//   node build.mjs --audit    proves the geometry of every figure, exits non-zero on any fault
//   node build.mjs --fig 2    writes figure-2.html, for reviewing one drawing alone
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import * as k from "./kit/iso-kit.mjs";
import * as A from "./kit/audit.mjs";
import { kitScript } from "./scripts/inline-kit.mjs";
import { FIGURE_CSS } from "./figure-theme.mjs";
import { buildIndexer } from "./figure-indexer.mjs";
import { buildReadStation } from "./figure-readstation.mjs";
import { buildInterlock } from "./figure-interlock.mjs";

const HERE = dirname(fileURLToPath(import.meta.url));
const argv = process.argv.slice(2);
const auditOnly = argv.includes("--audit");
const one = argv.includes("--fig") ? Number(argv[argv.indexOf("--fig") + 1]) : 0;

const FIGURES = [
  { name: "Port indexer", build: buildIndexer },
  { name: "Parameter read station", build: buildReadStation },
  { name: "Restart interlock", build: buildInterlock },
];

const built = FIGURES.map((figure) => ({ ...figure, ...figure.build() }));

// The audit is the proof, not a formality: it runs on every figure in every build.
for (const figure of built) A.auditOrExit(figure.R, figure.P);
if (auditOnly) process.exit(0);

const script = `${kitScript()}\n${built.map((figure) => figure.live).join("\n")}`;
// The kit's own source contains a literal closing script tag (inside its page
// builder), so anything inlined into a script element has to escape it.
const inline = (source) => source.replace(/<\/(script)/gi, "<\\/$1");

writeFileSync(join(HERE, "figure.css"), `${k.ISO_CSS}\n${FIGURE_CSS}`);
writeFileSync(join(HERE, "figure.js"), script);

if (one) {
  const figure = built[one - 1];
  mkdirSync(join(HERE, "review"), { recursive: true });
  writeFileSync(
    join(HERE, "review", `figure-${one}.html`),
    k
      .pageHtml({ title: figure.name, theme: "dark", body: figure.plate, script: "", width: 780 })
      .replace("</head>", () => `<style>${FIGURE_CSS}</style></head>`)
      .replace("</body>", () => `<script type="module">${inline(`${kitScript()}\n${figure.live}`)}</script></body>`),
  );
  console.log(`wrote review/figure-${one}.html`);
  process.exit(0);
}

const template = readFileSync(join(HERE, "page.template.html"), "utf8");
const page = template
  .replace("{{figure1}}", () => built[0].plate)
  .replace("{{figure2}}", () => built[1].plate)
  .replace("{{figure3}}", () => built[2].plate)
  .replace("{{script}}", () => inline(script));
writeFileSync(join(HERE, "..", "index.html"), page);

console.log(`wrote index.html · figure.css · figure.js · ${built.length} figures · ${built.reduce((total, figure) => total + figure.solids, 0)} solids`);
