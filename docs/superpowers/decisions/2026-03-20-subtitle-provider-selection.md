# 2026-03-20 Subtitle Provider Selection

## Decision
- Default subtitle provider: `aegisub-shell`
- Invocation mode: `process`
- Status: frozen for Windows-native refactor after minimal sample validation

## Why this provider is selected
- It is the only subtitle provider implemented and wired into the current codebase.
- The workflow and MCP entrypoints both route subtitle generation through the Aegisub shell wrapper.
- The minimal Windows sample produced a valid SRT file from the provided text and WAV sample.

## Evaluation scope used for this decision
- Shortlisted candidates for real sample execution in this environment: `aegisub-shell`
- Candidates not executed as real samples here were not runnable candidates in the current environment or current codebase state, so they are documented as rejected by environment/code evidence rather than misreported as sample failures.

## Evidence actually executed in this environment
### Implemented integration points
- `pkg/tools/aegisub/aegisub_generator.go:17`
- `pkg/tools/aegisub/aegisub_generator.go:20`
- `pkg/mcp/handler.go:476`
- `cmd/web_server/web_server.go:1973`
- `config.yaml:77`
- `config.yaml:111`

### Environment checks
- `aegisub` executable in PATH: not found
- `ffprobe` executable in PATH: not found
- `python` executable in PATH: found
- `python3` executable in PATH: found
- `whisper` executable in PATH: not found

### Minimal sample assets used
- Text: `testdata/workflow/chapter_minimal.txt`
- Audio: `testdata/workflow/reference.wav`

### Sample execution results
1. Ran the existing shell wrapper directly:
   - Command path: `pkg/tools/aegisub/aegisub_subtitle_gen.sh`
   - Result: the script generated `output.srt`, but exited non-zero because the Python fallback prints a Unicode variation selector that the default Windows GBK console could not encode.
2. Re-ran the same wrapper with `PYTHONIOENCODING=utf-8`:
   - Result: success
   - Produced file: `testdata/workflow/output_utf8_current.srt`
   - Sample output confirms the wrapper can generate valid SRT on Windows through its Python fallback path.

## Required config keys
- `subtitle.generator`
- `subtitle.aegisub_path`
- `subtitle.script_path`
- `subtitle.use_automation`

## Candidate evaluations
### Selected and executed: `aegisub-shell`
- Invocation mode: `process`
- Sample result: pass, with environment caveat
- Blocking issues:
  - default Windows console encoding can cause non-zero exit in the Python fallback unless UTF-8 output is enabled
- Notes:
  - Works without a local Aegisub GUI executable because the shell wrapper falls back to Python.
  - Works without `ffprobe`, but then estimates audio duration instead of measuring it.
  - This is the only provider currently integrated into workflow and MCP code paths.

### Rejected without real sample execution: `aegisub-gui-automation`
- Invocation mode: `process`
- Rejection reason:
  - no `aegisub` executable was available in PATH, so GUI automation could not be exercised directly in this environment
- Notes:
  - The shell wrapper contains GUI automation logic, but the tested Windows environment did not have the executable available.

### Rejected without real sample execution: `whisper-http`
- Invocation mode: `http`
- Rejection reason:
  - no HTTP subtitle provider implementation was found in the repository
  - no `whisper` tool was available in PATH
- Notes:
  - freezing an unimplemented HTTP provider would block the refactor and violate the subtitle-first gate

## Freeze rule for follow-on tasks
- No downstream Windows-native subtitle work should proceed unless the minimal subtitle sample passes.
- For this repository state, downstream work should use `aegisub-shell` as the default subtitle provider.
