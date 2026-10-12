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
