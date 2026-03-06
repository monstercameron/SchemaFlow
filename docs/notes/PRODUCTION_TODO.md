# Production TODO

Last updated: 2026-03-06

## Verification snapshot

Completed in this review:
- `go test ./...` in the main module: pass
- `go vet ./...` in the main module: pass
- `go run ./examples/tools all`: pass
- `go build ./...` in `examples/smarttodo`: pass
- Numbered example smoke run with `SCHEMAFLOW_PROVIDER=local`: 26 passed, 19 failed
- Live OpenAI smoke run after provider fixes: `Extract`, `Validate`, and `Rank` passed
- Full live OpenAI smoke run via `scripts/live_llm_ops_smoke.go`: 65 / 65 exported LLM ops passed
- Three-model live OpenAI matrix after remapping intelligence:
  - `Smart` / `gpt-5.4`: 65 / 65 passed
  - `Fast` / `gpt-5-mini`: 65 / 65 passed
  - `Quick` / `gpt-5-nano`: 65 / 65 passed

Numbered examples that passed under the local provider:
- `02-transform`, `07-summarize`, `08-classify`, `09-score`, `10-compare`, `11-similar`, `12-validate`, `13-merge`, `14-decide`, `15-guard`, `16-infer`, `17-diff`, `18-explain`, `19-parse`, `20-complete`, `21-redact`, `22-suggest`, `36-negotiate`, `37-resolve`, `38-derive`, `39-conform`, `41-arbitrate`, `42-project`, `43-audit`, `44-compose`, `45-pivot`

Numbered examples that failed under the local provider:
- `01-extract`, `03-generate`, `04-choose`, `05-filter`, `06-sort`, `23-annotate`, `24-cluster`, `25-rank`, `26-compress`, `27-decompose`, `28-enrich`, `29-normalize`, `30-match`, `31-critique`, `32-synthesize`, `33-predict`, `34-verify`, `35-question`, `40-interpolate`

## Fixed in this pass

- [x] OpenAI requests no longer send `temperature` to `gpt-5-*` models, which were failing with HTTP 400 in live usage: `internal/llm/provider.go:109`
- [x] OpenAI response parsing now scans all `output` items instead of assuming the first item contains text; GPT-5 responses can emit a `reasoning` item before the final `message`: `internal/llm/provider.go:161`
- [x] `NewEnrichOptions()` now validates successfully and can run without explicit derive fields by falling back to generic derivation instructions: `internal/ops/enrich.go:60`, `internal/ops/enrich.go:221`, `internal/ops/enrich.go:342`
- [x] `NewRankOptions()` now creates valid defaults: `internal/ops/rank.go:49`
- [x] `internal/ops/explain_test.go` was updated to match the current `types.OpOptions` field count: `internal/ops/explain_test.go:84`
- [x] Shared LLM execution now propagates JSON contracts, semantic grounding, and `Steering` into the actual provider request: `internal/ops/llm_helper.go:47`
- [x] GPT-5 Responses requests now inject a JSON hint into `input`, use minimal reasoning, and report incomplete responses clearly instead of surfacing as empty content: `internal/llm/provider.go:93`, `internal/llm/provider.go:121`, `internal/llm/provider.go:176`
- [x] Structured parsers were hardened for live-model variability in `Compress`, `DecomposeToSlice`, `Suggest`, `Synthesize`, and `Negotiate`: `internal/ops/compress.go:290`, `internal/ops/decompose.go:410`, `internal/ops/suggest.go:301`, `internal/ops/synthesize.go:334`, `internal/ops/negotiate.go:187`
- [x] Default OpenAI model mapping now targets current public names rather than stale snapshots: `internal/config/config.go:143`
- [x] `Sort` now has a scoring fallback when a model returns an incomplete ordering, which closed the last cross-model live smoke failure: `internal/ops/collection.go:253`
- [x] Tracing env handling now accepts both `SCHEMAFLOW_TRACE` and `SCHEMAFLOW_ENABLE_TRACING`, and tests cover both paths: `internal/config/config.go:68`, `telemetry/tracing.go:43`, `internal/config/config_test.go:5`, `telemetry/tracing_test.go:84`
- [x] `examples/smarttodo` was migrated off legacy `schemaflow.OpOptions` calls and now builds again: `examples/smarttodo/internal/processor/processor.go:57`, `examples/smarttodo/internal/processor/ai_quotes.go:68`, `examples/smarttodo/internal/tui/views_apikey.go:98`, `examples/smarttodo/cmd/smarttodo/main.go:48`

## Release blockers

- [ ] The local/mock provider is not a credible integration substitute for the typed API surface. Its default branch returns plain text such as `Mock response for: ...`, which is incompatible with structured operations like `Rank`, `Enrich`, `Predict`, `Verify`, and `Question`: `internal/llm/provider.go:685`
- [ ] Example programs are not a reliable release gate today. 19 of 45 numbered examples failed under the local provider because the mock layer does not satisfy their JSON contracts. These failures block a claim that the documented feature set is continuously runnable.
- [ ] Many examples hard-fail if a `.env` file is not found even when the required environment variables are already present in the process. This is not production-grade example behavior and makes CI and Windows usage brittle: representative loader pattern in `examples/24-cluster/main.go:18`

## Not production-ready features

The following tool families are explicitly stubbed and should not be marketed as production-ready until they have real integrations, tests, configuration docs, and failure semantics:
- [ ] Search/browser tooling: `internal/tools/http.go:193`, `internal/tools/http.go:213`, `internal/tools/http.go:235`
- [ ] Geo/weather tooling: `internal/tools/time.go:500`, `internal/tools/time.go:519`
- [ ] Vector database tooling: `internal/tools/database.go:402`
- [ ] File watching: `internal/tools/file.go:515`
- [ ] JWT/encryption helpers: `internal/tools/cache_security.go:405`, `internal/tools/cache_security.go:453`, `internal/tools/cache_security.go:472`
- [ ] Additional stub families documented in the tools example README: archive, image processing, audio, messaging, and several AI helpers in `examples/tools/README.md`

## Next work to reach a real release bar

- [ ] Add a first-class way for each operation to declare whether it requires JSON output, and enforce that through `CompletionRequest`
- [ ] Upgrade the local provider so every documented operation has a schema-compatible mock response; then convert the numbered example smoke run into CI
- [ ] Remove duplicated `.env` loaders from examples and let them honor normal environment variables first
- [ ] Split tool capabilities into `implemented` vs `stubbed` documentation so users do not confuse placeholders with working features
- [ ] Add Linux-based `-race` CI coverage, because the current machine cannot run `go test -race ./...` on Windows/arm64
