# ai-models

Mirror of [sherpa-onnx](https://github.com/k2-fsa/sherpa-onnx) streaming Zipformer speech models used by **Lumiquo**, **Teramigo** and **Exatiq**.

Binaries are served via **Cloudflare R2** (object keys `<tag>/<file>`, public custom domain) and mirrored on **GitHub Releases** as a fallback source.

## Releases

| Tag | Language | Source |
|-----|----------|--------|
| `zipformer-en-kroko-2025-08-06` | English | csukuangfj/sherpa-onnx-streaming-zipformer-en-kroko-2025-08-06 |
| `zipformer-fr-kroko-2025-08-06` | French | csukuangfj/sherpa-onnx-streaming-zipformer-fr-kroko-2025-08-06 |
| `zipformer-es-kroko-2025-08-06` | Spanish | csukuangfj/sherpa-onnx-streaming-zipformer-es-kroko-2025-08-06 |
| `zipformer-pt-kroko-64l-int8` | Portuguese (Brazil) | hudaiapa88/sherpa-stt-onnx (pt/kroko_64l) |

Each release contains the loose `encoder`, `decoder`, `joiner`, and `tokens.txt` artifacts.

Download URL patterns:

```
# Primary: Cloudflare R2 public bucket
https://<PUBLIC_BASE_URL>/<tag>/<file>

# Fallback: GitHub Releases
https://github.com/apoema-ai/ai-models/releases/download/<tag>/<file>
```

## Uploading to R2

Stage the release assets in a directory named after their tag, then run the
SigV4 uploader (stdlib-only Go, no SDK):

```bash
R2_ENDPOINT=https://<account>.r2.cloudflarestorage.com \
R2_BUCKET=ai-models \
R2_ACCESS_KEY_ID=... \
R2_SECRET_ACCESS_KEY=... \
go run tool/upload_r2.go /path/to/staged/models
```

Object keys preserve the `<tag>/<file>` layout, so the public URL
`https://<PUBLIC_BASE_URL>/<tag>/<file>` is a drop-in replacement for the
GitHub Releases URL. Apps pin each file's SHA-256, so re-uploads must be
byte-identical.
