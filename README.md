# ai-models

Mirror of [sherpa-onnx](https://github.com/k2-fsa/sherpa-onnx) streaming Zipformer speech models used by **Lumiquo**.

Binaries are served via **GitHub Releases** (not committed to the tree) so the app can fetch them over stable download URLs without Git LFS bandwidth limits.

## Releases

| Tag | Language | Source |
|-----|----------|--------|
| `zipformer-en-kroko-2025-08-06` | English | csukuangfj/sherpa-onnx-streaming-zipformer-en-kroko-2025-08-06 |
| `zipformer-fr-kroko-2025-08-06` | French | csukuangfj/sherpa-onnx-streaming-zipformer-fr-kroko-2025-08-06 |
| `zipformer-es-kroko-2025-08-06` | Spanish | csukuangfj/sherpa-onnx-streaming-zipformer-es-kroko-2025-08-06 |
| `zipformer-pt-kroko-64l-int8` | Portuguese (Brazil) | hudaiapa88/sherpa-stt-onnx (pt/kroko_64l) |

Each release contains the loose `encoder`, `decoder`, `joiner`, and `tokens.txt` artifacts.

Download URL pattern:

```
https://github.com/apoema-ai/ai-models/releases/download/<tag>/<file>
```
