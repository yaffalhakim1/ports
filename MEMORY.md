# Project memory

Facts below are recalled into every session.

## Project identity: `ports` is a Go desktop app built with `github.com/egoist/mygo
Project identity: `ports` is a Go desktop app built with `github.com/egoist/mygo` (see mygo-maintenance skill).

## Vet command: run `go tool mygo vet .` (not plain `go vet`) to lint this project.
Vet command: run `go tool mygo vet .` (not plain `go vet`) to lint this project.

## Repo: public GitHub repo at https://github.com/yaffalhakim1/ports.
Repo: public GitHub repo at https://github.com/yaffalhakim1/ports.

## Appearance control convention: theme selector is a one-click menu (`System / Lig
Appearance control convention: theme selector is a one-click menu (`System / Light / Dark`) at `view_header.go:169` — no cycle/double-click behavior.

## Build artifacts to clean before commit: stray `ports.exe`, `out.txt`, `err.txt` 
Build artifacts to clean before commit: stray `ports.exe`, `out.txt`, `err.txt` are junk and must not be committed.

## CI/release workflow: intentionally omitted; add a GitHub Actions workflow only w
CI/release workflow: intentionally omitted; add a GitHub Actions workflow only when auto-published tagged installers are wanted.

## close-to-tray: implemented in `tray.go` (`OnClose` → `PreventDefault` + `Hide`);
close-to-tray: implemented in `tray.go` (`OnClose` → `PreventDefault` + `Hide`); tray menu has Show/Quit

## no-terminal build: MyGo builds with `-H=windowsgui` so PE Subsystem = 2 (GUI); n
no-terminal build: MyGo builds with `-H=windowsgui` so PE Subsystem = 2 (GUI); no console window opens

## install layout: per-user install at `%LOCALAPPDATA%\Programs\ports\ports.exe` wi
install layout: per-user install at `%LOCALAPPDATA%\Programs\ports\ports.exe` with Start Menu `ports.lnk` and registered uninstaller (no admin prompt); window title is "Ports"

## app identity: Go desktop app built with github.com/egoist/mygo (see mygo-mainten
app identity: Go desktop app built with github.com/egoist/mygo (see mygo-maintenance skill)

## ports UI stack: project uses mygo (pure-Go, no-cgo, no-webview, GPU-drawn deskto
ports UI stack: project uses mygo (pure-Go, no-cgo, no-webview, GPU-drawn desktop UI toolkit), not a webview stack.

## mygo toolchain lock-in: driven by `mygo.json` config + the `go tool mygo` comman
mygo toolchain lock-in: driven by `mygo.json` config + the `go tool mygo` command, not just a package import — migrating off is a project.

## mygo version: currently v0.3.5 (0.x) — minor bumps require migration guides, so 
mygo version: currently v0.3.5 (0.x) — minor bumps require migration guides, so the API is still churning/unstable.

## mygo testing: `ui.NewTester` provides real headless UI tests (click by label, as
mygo testing: `ui.NewTester` provides real headless UI tests (click by label, assert focus, snapshot PNGs).

## mygo API model: immediate-mode build with retained identity via `Handle`/`Bind`/
mygo API model: immediate-mode build with retained identity via `Handle`/`Bind`/`Key`, with stale-element panics in dev.

## ports-app-skeleton: The `ports` app is a MyGo tray-resident Windows utility buil
ports-app-skeleton: The `ports` app is a MyGo tray-resident Windows utility built on a reusable skeleton — list + per-item actions + tray icon + 2s polling loop — which is intended to be repurposed for other small utilities (services manager, docker tray, etc.) by swapping the `scan()`.

## ports-app identity: `ports` is a Windows-native single-binary dev tool with a li
ports-app identity: `ports` is a Windows-native single-binary dev tool with a list/actions/tray UI whose job is showing "what's running on my machine."

## ports-roadmap decision: preferred next step is to deepen `ports` by adding servi
ports-roadmap decision: preferred next step is to deepen `ports` by adding services and containers as row types (reusing existing list/actions/tray skeleton), not start a second app.

## windows-incumbent heuristic: on Windows, "boring + useful" usually means Microso
windows-incumbent heuristic: on Windows, "boring + useful" usually means Microsoft or a strong free tool already owns it (EarTrumpet, Win+V, powercfg, services.msc) — so a good incumbent = dead idea.

## ports-new-app exception: only defended new app is an SSH tunnel / port-forward m
ports-new-app exception: only defended new app is an SSH tunnel / port-forward manager (no good native Windows incumbent, reuses existing port code).

## agent-self-correction gotcha: agent tends to emit "opportunity" listicles for Wi
agent-self-correction gotcha: agent tends to emit "opportunity" listicles for Windows utilities without first checking incumbents; retracted earlier sloppy list.

## ports stack: `ports` is a Go, tray-resident Windows desktop app.
ports stack: `ports` is a Go, tray-resident Windows desktop app.

## ports UI skeleton: header/table/footer layout is the reusable shell for new sibl
ports UI skeleton: header/table/footer layout is the reusable shell for new sibling tools.

## ports scan discipline: filesystem/enumeration scans use a cached index + backgro
ports scan discipline: filesystem/enumeration scans use a cached index + background refresh pattern (see `scan.go`) rather than blocking scans.

## ports feature set: list model, per-row actions, tray, background poll every N se
ports feature set: list model, per-row actions, tray, background poll every N seconds, search/filter, settings — reusable across sibling apps.

## Windows gotcha: no native file-lock detective; LockHunter is old nagware and `ha
Windows gotcha: no native file-lock detective; LockHunter is old nagware and `handle.exe` is CLI-only (handle-enumeration work needed).

## Idea selection lens: an idea only qualifies if its incumbent is cloud, subscript
Idea selection lens: an idea only qualifies if its incumbent is cloud, subscription, or Electron-heavy — that's the only context where a ~10 MB offline native binary is a feature, not a nicety; the product is "offline, private, instant," not "a better X."

## Candidate directions under evaluation: (1) local dictation/speech-to-text — top 
Candidate directions under evaluation: (1) local dictation/speech-to-text — top pick, hard part is audio capture + bundling whisper.cpp; (2) per-process network/bandwidth monitor — easiest, shippable this week; (3) native local LLM runner — visible win over Electron incumbents but churns fast.

## MyGo fit criteria for candidates: system-wide hotkey, tray-resident, instant sta
MyGo fit criteria for candidates: system-wide hotkey, tray-resident, instant startup, no webview, always-running single-purpose tool.

## ports brand palette: reusable brand colors sourced from the `ports` project; reu
ports brand palette: reusable brand colors sourced from the `ports` project; reuse for consistent theming across apps.

## Kata project: Japanese → Indonesian vocabulary learn app built with MyGo, hosted
Kata project: Japanese → Indonesian vocabulary learn app built with MyGo, hosted at https://github.com/yaffalhakim1/kata.

## MyGo apps: scaffold a native MyGo project, then study the widget/UI API surface 
MyGo apps: scaffold a native MyGo project, then study the widget/UI API surface before writing views.

## MyGo gotcha: select widgets render blank when no default is set — default the To
MyGo gotcha: select widgets render blank when no default is set — default the Topic select to "All".

## MyGo conventions: split app by responsibility (state, views, model, icons) with 
MyGo conventions: split app by responsibility (state, views, model, icons) with pure quiz logic kept testable.

## MyGo verification: render screenshots to verify look and CJK/Japanese font rende
MyGo verification: render screenshots to verify look and CJK/Japanese font rendering; run tests plus launch the app before finishing.

## Go GUI layout gotcha: `Grow` in a column grows vertically, not horizontally — us
Go GUI layout gotcha: `Grow` in a column grows vertically, not horizontally — use `FillWidth().MaxWidth(contentMax)` to fill width while capping it.

## View files are split per screen: `view_common.go`, `view_study.go`, `view_quiz.g
View files are split per screen: `view_common.go`, `view_study.go`, `view_quiz.go`.

## Convention: review with `better-layout` + `better-ui` by inspecting rendered sta
Convention: review with `better-layout` + `better-ui` by inspecting rendered states (wide, dark, 420px minimum, empty), not source alone.

## Convention: convert throwaway narrow-width renders into real regression tests as
Convention: convert throwaway narrow-width renders into real regression tests asserting controls stay inside the window.

## repo: kata — Japanese→Indonesian vocabulary app at https://github.com/yaffalhaki
repo: kata — Japanese→Indonesian vocabulary app at https://github.com/yaffalhakim1/kata

## screenshots: generated deterministically by tests, not manual captures — `go tes
screenshots: generated deterministically by tests, not manual captures — `go test -run TestRender` writes them to `testdata/`

## study screen: search + topic filter; each word shows 日本語 + romaji, Indonesian me
study screen: search + topic filter; each word shows 日本語 + romaji, Indonesian meaning in brand green, topic badge, example sentence

## quiz setup: user picks direction (日本語→Indonesia or reverse) and question count (
quiz setup: user picks direction (日本語→Indonesia or reverse) and question count (5/10/20/All)

## quiz UX: answered state marks correct choice green with check icon and shows wor
quiz UX: answered state marks correct choice green with check icon and shows word breakdown + example so wrong answers teach

## theming: app uses a brand light and brand dark palette (both rendered for screen
theming: app uses a brand light and brand dark palette (both rendered for screenshots)

## Teaching model decided from competitor research: teach kana first, row by row, w
Teaching model decided from competitor research: teach kana first, row by row, with mnemonics (Tofugu style); teach by concept — scripts, politeness, particles, sentence order (Tae Kim style) — not word lists.

## Test gotchas: tests assumed the old default tab and must explicitly set their ta
Test gotchas: tests assumed the old default tab and must explicitly set their tab; word tests require the `studyApp()` helper; layout tests reference the new "Words" tab and `TabQuiz` with the Words source.

## Verification: app builds, tests pass, and four pages are rendered to screenshots
Verification: app builds, tests pass, and four pages are rendered to screenshots; narrow-width (420px) is checked for clipping.

## screenshot-rendering: Screenshots are generated reproducibly via `go test -run T
screenshot-rendering: Screenshots are generated reproducibly via `go test -run TestRender` and committed, then referenced in the README.

## japanese-learning-app: App teaches Japanese with sections: Basics, Kana (46 hira
japanese-learning-app: App teaches Japanese with sections: Basics, Kana (46 hiragana with mnemonics), Words, and Quiz.

## theme-support: UI supports light and dark modes (word cards rendered in both).
theme-support: UI supports light and dark modes (word cards rendered in both).

## brand-color-usage: Indonesian meanings on word cards are styled in the brand gre
brand-color-usage: Indonesian meanings on word cards are styled in the brand green accent.

## quiz-reveal-ux: Quiz marks the correct choice green with a check and reveals the
quiz-reveal-ux: Quiz marks the correct choice green with a check and reveals the answer (e.g. `いい（ii）— bagus`) so wrong answers still teach.

## Screenshots: the `kata` repo (yaffalhakim1/kata) stores UI screenshots in the `t
Screenshots: the `kata` repo (yaffalhakim1/kata) stores UI screenshots in the `testdata/` folder, which are referenced in the README.

## Repo: `kata` lives at https://github.com/yaffalhakim1/kata (Japanese/kana learni
Repo: `kata` lives at https://github.com/yaffalhakim1/kata (Japanese/kana learning app).

## Learning flow: Basics → Kana → Words → Quiz; beginners start at Basics.
Learning flow: Basics → Kana → Words → Quiz; beginners start at Basics.

## App tab structure: The app has four pages/tabs — Basics (default first tab), Kan
App tab structure: The app has four pages/tabs — Basics (default first tab), Kana, Words, Quiz — with in-app buttons that jump between them.

## Onboarding guide location: The beginner "Your path" step-by-step guide lives as 
Onboarding guide location: The beginner "Your path" step-by-step guide lives as a card at the top of the Basics page, and its content is intentionally mirrored in the README so learners see it in-app and in the repo.

## Minimum window width: The UI must render cleanly down to a 420px minimum width (
Minimum window width: The UI must render cleanly down to a 420px minimum width (also checked at 820px); the guide's button row is the tight case to watch for clipping.

## Verification workflow: Changes are verified with Go tests, `go vet`, and manual 
Verification workflow: Changes are verified with Go tests, `go vet`, and manual render checks at 820px/420px; regression tests are added for interactive behavior (e.g. buttons actually switching tabs).

## Brand palette locked in palette.go (graphite + brand green): preserve it for any
Brand palette locked in palette.go (graphite + brand green): preserve it for any UI/landing work in this repo.

## Landing typography decision: Zen Maru Gothic (display), Geist (body), Geist Mono
Landing typography decision: Zen Maru Gothic (display), Geist (body), Geist Mono (mono).

## Hallmark skill: first run for this project (no prior diversification stamp to ho
Hallmark skill: first run for this project (no prior diversification stamp to honor).

## ports: Personal Go utility project built with github.com/egoist/mygo — has tray 
ports: Personal Go utility project built with github.com/egoist/mygo — has tray integration, process control, and tests; uses the mygo `ui` package and `mygo build`.

## kata: Personal Go project built with MyGo — contains kana, concepts, a quiz engi
kata: Personal Go project built with MyGo — contains kana, concepts, a quiz engine, and a landing page (landing page modeled as a Workbench page).

## Pulse reference: egoist/pulse is a cross-platform energy/network monitor built o
Pulse reference: egoist/pulse is a cross-platform energy/network monitor built on the same MyGo stack; pages are structured as `view_*.go` files, and it demonstrates menu-bar integration and per-process icon retrieval.

## Reference workflow: User's convention when learning from a peer app is to clone 
Reference workflow: User's convention when learning from a peer app is to clone it and steal structure (not pixels), then port one idea back into their own project (e.g. per-process icons or tray readout into ports).

## ports repo layout: restructured into internal/port (pure model), internal/host
ports repo layout: restructured into internal/port (pure model), internal/host (system reads + process control, all //go:build windows), internal/app (state, views, theme, tray, settings); main.go stays at root so mygo.json resolves. Deps flow main -> app -> host -> port.

## ports release flow: `git tag vX.Y.Z && git push origin vX.Y.Z` triggers
ports release flow: `git tag vX.Y.Z && git push origin vX.Y.Z` triggers .github/workflows/release.yml, which tests, builds both Windows arches, checks the tag against mygo.json's version, and publishes a GitHub release with installers + signed update files.

## ports updater: app uses the native updater plugin (github.com/egoist/mygo/plugins/updater/native) with updates.publicKey + updates.github in mygo.json; tray has a "Check for Updates" item. Private key lives at %APPDATA%\mygo\update-keys\mygo-update.key and as the MYGO_UPDATER_PRIVATE_KEY repo secret — losing it strands installed apps.

## mygo release gotcha: both Windows architectures build an installer with the same
mygo release gotcha: both Windows architectures build an installer with the same basename ("ports Setup X.Y.Z.exe"), and a GitHub release cannot hold two same-named assets — so build signed WITHOUT -upload and upload per-arch-renamed installers; the .tar.gz and update-*.json names must stay as built because the manifests reference them.

## ports landing page: Hallmark-built static page in landing/ (index.html, tokens.css, assets/, favicon.png), published by .github/workflows/pages.yml to https://yaffalhakim1.github.io/ports/ on changes under landing/. Apple-design pass applied (translucent sticky nav, press feedback, reduced-motion/transparency fallbacks).

## omp skill registry root: omp indexes `~/.agents/skills`, NOT `~/.omp/skills`; sk
omp skill registry root: omp indexes `~/.agents/skills`, NOT `~/.omp/skills`; skills present only under `~/.omp/skills/` are invisible to `skill://` and absent from the available-skills list.

## anatomy skill location: full skill on disk at `C:/Users/yafit/.omp/skills/anatom
anatomy skill location: full skill on disk at `C:/Users/yafit/.omp/skills/anatomy/` (13 MB; SKILL.md, kit/, references/, scripts/, examples/, evals/, agents/openai.yaml, tests/), byte-identical duplicate at `~/.claude/skills/anatomy`; frontmatter valid (`name: anatomy`, description 971 chars < 1024 cap) but unloadable due to registry root mismatch.

## anatomy skill purpose: builds crafted isometric SVG/WebGL figures explaining ide
anatomy skill purpose: builds crafted isometric SVG/WebGL figures explaining ideas as physical machines — four-tone shaded solids, painter's order, framed card with live readout; kit/ drawing library (`iso-kit`, `lathe.mjs`, `tube.mjs`, `gl.mjs`, `turn.mjs`); verification via `node build.mjs --audit` geometry proof, `scripts/lines.mjs` line-end check, `scripts/capture.mjs`/`drive.mjs` headless-Chrome screenshots.

## anatomy skill requirements: Node >= 20.10, zero deps, Chrome required for checks
anatomy skill requirements: Node >= 20.10, zero deps, Chrome required for checks. Local env satisfies this — Node v22.23.2, Chrome at `C:/Program Files/Google/Chrome/Application/chrome.exe`.

## anatomy fix option: to make loadable, symlink/junction or copy `anatomy` into `~
anatomy fix option: to make loadable, symlink/junction or copy `anatomy` into `~/.agents/skills/anatomy` (13 MB duplicate), or add `~/.omp/skills` to omp's skill roots if supported.

## skill list anomaly: available list vs `~/.agents/skills` contents match exactly 
skill list anomaly: available list vs `~/.agents/skills` contents match exactly except `mygo-maintenance` (project-local) and a stray `scripts/` dir.

## omp skill registry ignores junctions: a skill junctioned/symlinked into the skil
omp skill registry ignores junctions: a skill junctioned/symlinked into the skills dir is NOT registered; read its instruction + kit files directly from disk.

## ports app surfaces: rescans listening ports every 2 s; `0.0.0.0` binds denote th
ports app surfaces: rescans listening ports every 2 s; `0.0.0.0` binds denote the "exposed" state; resolved-folder list is the other main view.

## landing-page concept chosen: "bench-top socket analyser" 3D figure — index head 
landing-page concept chosen: "bench-top socket analyser" 3D figure — index head traverses a row of port sockets, spring probe drops into occupied sockets, odometer counts, marks paper tape.

## 3D kit gotcha: a rod inside a bore trips `solidClearances` as a sunk pair; make 
3D kit gotcha: a rod inside a bore trips `solidClearances` as a sunk pair; make the probe a `ringBand` bore so the rod isn't flagged.

## 3D kit recipe: sockets are pucks the probe rests on (dial-indicator contact reci
3D kit recipe: sockets are pucks the probe rests on (dial-indicator contact recipe, passes audit); probe is a `rodOut` plunger retracting into the head.

## 3D kit gotcha: gantry posts can intersect nearby parts (e.g. dial puck at ~4.8 d
3D kit gotcha: gantry posts can intersect nearby parts (e.g. dial puck at ~4.8 depth), triggering depth-cycle audit failures; move posts inboard, dials outboard.

## Kit verification: build must reach audit (`solidClearances`) zero and line-check
Kit verification: build must reach audit (`solidClearances`) zero and line-check zero before render is accepted.

## Deployment target: GitHub Pages at https://yaffalhakim1.github.io/ports, publish
Deployment target: GitHub Pages at https://yaffalhakim1.github.io/ports, publishing the `landing/` directory via `.github/workflows/pages.yml` (upload-pages-artifact path: landing → deploy-pages).

## Deploy trigger gotcha: the Pages workflow only fires on push to `main` that touc
Deploy trigger gotcha: the Pages workflow only fires on push to `main` that touches `landing/**` (or manual `workflow_dispatch`); changes to `README.md` or `MEMORY.md` alone will not deploy.

## Commit policy: do not `git add`/`commit`/`push` without the user explicitly aski
Commit policy: do not `git add`/`commit`/`push` without the user explicitly asking; when asked, propose the commit contents first.

## Repo layout: static landing site lives in `landing/` (with `landing/index.html`,
Repo layout: static landing site lives in `landing/` (with `landing/index.html`, `landing/anatomy/`), alongside `README.md` and `MEMORY.md` at the project root.

## Repo/deploy: `yaffalhakim1/ports` (Go) auto-deploys GitHub Pages from `main`; li
Repo/deploy: `yaffalhakim1/ports` (Go) auto-deploys GitHub Pages from `main`; live at https://yaffalhakim1.github.io/ports/.

## Vet command: use `go tool mygo vet .` — plain `go vet` is not the project's base
Vet command: use `go tool mygo vet .` — plain `go vet` is not the project's baseline.

## Windows quirk: `curl -o /dev/null` prints write warnings on this machine; they a
Windows quirk: `curl -o /dev/null` prints write warnings on this machine; they are not request failures.

## Verification convention: after pushing to `main`, verify against the *deployed* 
Verification convention: after pushing to `main`, verify against the *deployed* URL, not just green CI — `cmp` live assets against local (byte-identical), drive the live page in Chromium, assert rendered counts and console-clean.

## Lazy-loading gotcha: `loading="lazy"` images report as not-loaded in an initial 
Lazy-loading gotcha: `loading="lazy"` images report as not-loaded in an initial page check (e.g. `imgsOk:false`); they load on scroll.

## Landing-page figure: implemented in `anatomy/figure.js` + `anatomy/figure.css`, 
Landing-page figure: implemented in `anatomy/figure.js` + `anatomy/figure.css`, embedded from `index.html`; renders 57 solids / 183 strokes with pointer + keyboard indexing and a live readout.

## landing-page redesign: full industrial/anatomy rebuild in one theme (dark graphite
landing-page redesign: full industrial/anatomy rebuild in one theme (dark graphite paper `oklch(15% 0.006 275)` + brand green `#26E085` as a signal only). Macrostructure is Map / Diagram — the previous build stamped Workbench, and Hallmark's diversification rule forbids repeating it.

## landing-page figures: three machines, one per idea — Fig 1 port indexer (socket
landing-page figures: three machines, one per idea — Fig 1 port indexer (socket per listener, post height = port number, wide collar = exposed); Fig 2 parameter read station (two `node.exe` blocks on a shuttle, fixed feeler reads the engraved folder); Fig 3 keyed restart interlock (bolt in the trip lever's fork, key must turn before the lever travels). Built in `landing/anatomy/figure-{indexer,readstation,interlock}.mjs`, assembled by `build.mjs`.

## anatomy build gotcha: the kit's own source contains a literal `</script>` (insid
anatomy build gotcha: the kit's own source contains a literal `</script>` (inside `pageHtml`), so inlining `kitScript()` into a `<script>` element closes it early and breaks the page. Escape it: `source.replace(/<\/(script)/gi, "<\\/$1")`.

## anatomy multi-figure gotcha: on a page with several figures, each live layer must
anatomy multi-figure gotcha: on a page with several figures, each live layer must scope its own DOM — `mount(document)` grabs the first stage and the first `[data-readout]`, so the second and third figures crash on null refs and write to the wrong readout. Give each stage a `data-figure` and mount into it.

## anatomy geometry gotchas learned on Fig 3: a bore in the middle of a rod wraps it
anatomy geometry gotchas learned on Fig 3: a bore in the middle of a rod wraps it and no drawing order can exist (use a cradle, a tangent contact); a closed eye around a rod is a real interlock (cut it open into a fork and split the jaws across the rod in paint order); a disc mounted on a vertical face cannot be ordered against it (put the wheel on a pedestal so it rests on a top face, as the dial does).

## landing-page verification (must stay green): `node build.mjs --audit` prints `audit
landing-page verification (must stay green): `node build.mjs --audit` prints `audit passed` for all three figures; `node scripts/lines.mjs ../index.html` prints `line ends: 0`; responsive at 320/375/414/768/1024/1440 with no sideways scroll, no overflow, no wrapped buttons, and all three figures at one scale.

## landing-page review output lives in `landing/anatomy/review/` (gitignored); `nod
landing-page review output lives in `landing/anatomy/review/` (gitignored); `node build.mjs --fig n` writes one figure on its own page there. The old committed `figure.html` / `figure-light.html` review pages are gone, and there is no light theme any more.

## landing page keeps real app captures: `landing/assets/screenshot-dark.png` and `
landing page keeps real app captures: `landing/assets/screenshot-dark.png` and `-light.png` (copies of `docs/`) plus the empty-search and window-disabled states. Hallmark forbids re-drawn UI chrome, so these are real screenshots in a hairline frame, never fake window bars.
