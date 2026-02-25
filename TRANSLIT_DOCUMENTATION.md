# Dhivehi Transliteration Engine — Architecture Documentation

This document describes the four transliteration implementations (`translit1`–`translit4`), their design, mapping strategy, engine logic, context handling, options, and determinism. It does not modify any transliteration logic.

---

## 1. Overview

| Version | Path | Primary focus |
|--------|------|----------------|
| V1 | `internal/translit1` | Array lookups, optional gemination/glottal rules |
| V2 | `internal/translit2` | Map-based, letter names, context rules |
| V3 | `internal/translit3` | Array lookups + Options (Gemination, NormalizeArabic, Nishaan) |
| V4 | `internal/translit4` | Byte-level UTF-8, bitmasks, no options |

---

## 2. Design Approach

| Aspect | translit1 | translit2 | translit3 | translit4 |
|--------|-----------|-----------|-----------|-----------|
| **Input unit** | `[]rune` (Unicode code points) | `[]rune` | `[]rune` | Raw bytes (UTF-8) |
| **Lookup style** | Arrays indexed by `(r - 0x0780)` | `map[rune]string` (Akuru, Fili, etc.) | Arrays + optional normalized consonant table | Byte pairs: `0xDE` + second byte; bitmasks for sets |
| **State** | `lastRune`, `lastLatin`, `pending`, `posInWord`, `geminateNext`, `lastVowel` | Index-based loop, no explicit “pending” | Same as V1 minus `lastVowel` | `prevIdx`, byte index `i` |
| **Output** | `strings.Builder` | `strings.Builder` | `strings.Builder` | Pre-allocated `[]byte`, `unsafe.String` at end |
| **Punctuation** | Not handled | `Nishaan` map (Arabic punctuation → Latin) | `nishaanLat` / `nishaanOk` arrays | Inline UTF-8 branches for `،;?` |

---

## 3. Mapping Strategy

| Aspect | translit1 | translit2 | translit3 | translit4 |
|--------|-----------|-----------|-----------|-----------|
| **Consonants** | `ConsonantMap` → init to `cLat`/`cOk` arrays | `Akuru` map (rune → Latin) | `consonantData` → `cLat`; `consonantNormData` → `cLatNorm` (when `NormalizeArabic`) | `akuruValues` array; index = second byte − 0x80 |
| **Vowels** | `VowelMap` → `vLat`/`vOk`; sukun excluded in init | `Fili` map | `vowelData` → `vLat`/`vOk` | `filiValues` array |
| **Arabic-derived** | Normalized in map (e.g. ޝ→"sh", ޢ→"'") | Distinct forms in Akuru (e.g. apostrophe forms) | V2-style in `cLat`; normalized in `cLatNorm` when option set | Single set in `akuruValues` (e.g. kh, sh', etc.) |
| **Standalone names** | No | `AkuruNames` (e.g. “haa”, “alifu”) for bare consonants | `akuruNameData` → `akNames`/`akNamesOk` (present but not used in main path in same way as V2) | `akuruNameValues`; used when “bare akuru” |
| **Sukun overrides** | Inline (ށ→h, ނ+ބ/ޕ→m, އ+cons→geminate) | `SukunOverrides` map (e.g. thaalu→"iy", ainu→"u") | `sukunOverrideData` → `skOver`/`skOverOk` | `sukunOvrdValues` + `sukunOvrdMask` |

---

## 4. Engine Logic Differences

- **translit1**: Single pass; word boundaries on space/newline/tab; sukun handled with special cases for ށ, ނ+ބ/ޕ, އ; Alifu can insert glottal stop (with diphthong/position checks); final Alifu+sukun → `h`.
- **translit2**: Look-ahead: akuru+fili (two runes), akuru+sukun (two runes), then bare akuru; Raa between fili/akuru → `r`; Noonu between fili and next akuru → `n'`; no Options struct.
- **translit3**: Same flow as V1 but with nishaan first, sukun overrides table, and Alifu never outputs glottal before vowel (V2-style). Gemination and NormalizeArabic via Options.
- **translit4**: Byte scanner; detects Thaana by `0xDE` + next byte; uses bitmasks (`akuruMask`, `filiMask`, etc.) to classify; same semantic rules as V2 (ainu+fili, sukun, noonu, raa, bare names) but no rune allocation in hot path.

---

## 5. Context Handling (Ainu, Sukun, Tashdid)

| Context | translit1 | translit2 | translit3 | translit4 |
|---------|-----------|-----------|-----------|-----------|
| **Ainu (ޢ) + fili** | Not special (carrier empty) | First char of fili + `'` + rest | Same: first char of vowel + `'` + rest | Same (byte copy) |
| **Sukun** | ށ→h; ނ+ބ/ޕ→first of next; އ+cons→geminate, else h; else `lastLatin` | Overrides map; Alifu/Shaviyani→h or next’s first; Noonu+meemu/baa/paviyani→first | Overrides first; then Shaviyani/Alifu/Noonu rules; else `lastLatin` | Same as V2 with bitmask checks |
| **Tashdid (gemination)** | Only if `Options.Gemination`: cons+sukun+same cons → doubled | Not implemented | Only if `Options.Gemination` | Not implemented |
| **Alifu + sukun** | Next consonant → set `geminateNext`; else `h` | Next consonant → output next’s first; else `h` | Same (geminateNext or h) | Same |
| **Word boundary** | Whitespace resets state | Implicit (no “word” state) | Whitespace resets state | `prevIdx` used for raa/noonu context |

---

## 6. Options Struct Usage

| Version | Options | Notes |
|---------|---------|-------|
| **translit1** | `Options{Gemination, SuppressGlottalStop}` | Gemination: cons+sukun+same → double. SuppressGlottalStop: no `'` between vowels (diphthong/position still apply). |
| **translit2** | None | No options; single behavior. |
| **translit3** | `Options{Gemination, SuppressGlottalStop, NormalizeArabic}` | Same as V1 for first two; NormalizeArabic uses `cLatNorm` for Arabic-derived letters. |
| **translit4** | None | No options; fixed behavior (no gemination, no normalize toggle). |

---

## 7. Determinism

- **translit1**: Deterministic for given input and options; no randomness or map iteration order in output.
- **translit2**: Deterministic; map iteration is not used for output order (explicit sequence of checks).
- **translit3**: Deterministic; same as V1/V2.
- **translit4**: Deterministic; fixed arrays and bitmasks.

---

## 8. Structured Comparison Table

| Criterion | translit1 | translit2 | translit3 | translit4 |
|-----------|-----------|-----------|-----------|-----------|
| **Design** | Rune loop, array lookups | Rune loop, map lookups | Rune loop, array lookups, nishaan | Byte loop, bitmasks, arrays |
| **Mapping** | Init from maps to arrays | Maps only | Init from maps; dual consonant tables | Static arrays + bitmasks |
| **Ainu** | Treated as empty carrier | First vowel char + `'` + rest | Same as V2 | Same as V2 |
| **Sukun** | Inline cases | Overrides + Alifu/Shaviyani/Noonu | Override table + same rules | Same as V2 (bitmask) |
| **Tashdid** | Via option | No | Via option | No |
| **Glottal (Alifu)** | Optional suppression, diphthong rule | N/A (no glottal) | No glottal (V2 style) | N/A |
| **Nishaan** | No | Yes (map) | Yes (array) | Yes (inline bytes) |
| **Normalize Arabic** | Built-in (map) | No (distinct letters) | Option | No |
| **Standalone names** | No | Yes | Tables present | Yes |
| **API** | `Transliterate`, `TransliterateWithOptions` | `Transliterate` | Same as V1 | `Transliterate` |

---

## 9. Benchmark & accuracy artifacts

- **Performance**: Run `go test ./benchmark/ -run TestWriteBenchmarkCSV -v` to produce `benchmark_results.csv` (version, ns_per_op, allocs_per_op, bytes_per_op) using a shared dataset of 10,000+ mixed Dhivehi words.
- **Accuracy**: Run `go test ./benchmark/ -run TestWriteAccuracyReport -v` to produce `accuracy_report.json` (exact match %, character-level edit distance) against `testdata/golden_cases.txt` (Qawaaidu-aligned expected output).
- **Graphs**: Run `go run ./cmd/benchgraph/` (from project root, after CSV and JSON exist) to generate `benchmarks_speed.png` and `benchmarks_accuracy.png`.
---

## 10. Final Summary & Recommendation


- **Fastest version**: **translit4** (lowest ns/op; byte-level UTF-8 and bitmasks, no options).
- **Most accurate version**: **translit3** (100% exact match on Qawaaidu golden set; Options for Gemination, NormalizeArabic, Nishaan).
- **Best tradeoff version**: **translit3** when Qawaaidu alignment is required; **translit4** when throughput is the priority and near-Qawaaidu accuracy is acceptable (translit2 and translit4 both ~98.7% on the same golden set).

### Recommendation

- **Production default for accuracy-critical use**: Use **translit3** when compliance with the Dhivehi Bas Latin Akurun Liyumuge Qawaaidu is required. It is the only version that scores 100% on the golden set and supports Gemination and NormalizeArabic via options.
- **Production default for throughput-critical use**: Use **translit4** when processing large volumes and speed matters more than optional features. It is roughly 3–4× faster than translit3 on the same 10k-word dataset with comparable accuracy to translit2 (~98.7%).
- **Current `cmd` default**: The CLI currently defaults to **v4**; keep this for speed. For strict Qawaaidu output, callers should select **v3** (e.g. `-v3` flag).
