// Command upload_r2 mirrors the staged model directories to a Cloudflare R2
// bucket over the S3 API (SigV4, unsigned payload over TLS — no SDK needed).
//
// Usage:
//
//	R2_ENDPOINT=https://<account>.r2.cloudflarestorage.com \
//	R2_BUCKET=ai-models \
//	R2_ACCESS_KEY_ID=... \
//	R2_SECRET_ACCESS_KEY=... \
//	go run tool/upload_r2.go <dir>
//
// <dir> contains one subdirectory per release tag (e.g. zipformer-en-…/)
// and is uploaded preserving paths, so object keys look like
// "<tag>/<file>" — the same shape the GitHub release URLs had.
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	endpoint := strings.TrimRight(os.Getenv("R2_ENDPOINT"), "/")
	bucket := os.Getenv("R2_BUCKET")
	keyID := os.Getenv("R2_ACCESS_KEY_ID")
	secret := os.Getenv("R2_SECRET_ACCESS_KEY")
	if endpoint == "" || bucket == "" || keyID == "" || secret == "" {
		fatal("R2_ENDPOINT, R2_BUCKET, R2_ACCESS_KEY_ID and R2_SECRET_ACCESS_KEY are required")
	}
	if len(os.Args) != 2 {
		fatal("usage: upload_r2.go <dir>")
	}
	root := os.Args[1]

	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fatal(err.Error())
	}
	sort.Strings(files)
	if len(files) == 0 {
		fatal("no files under " + root)
	}

	client := &http.Client{Timeout: 30 * time.Minute}
	for _, path := range files {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			fatal(err.Error())
		}
		key := filepath.ToSlash(rel)
		if err := put(client, endpoint, bucket, key, path, keyID, secret); err != nil {
			fatal(fmt.Sprintf("%s: %v", key, err))
		}
		fmt.Printf("uploaded %s\n", key)
	}
	fmt.Printf("done: %d objects\n", len(files))
}

func put(client *http.Client, endpoint, bucket, key, path, keyID, secret string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}

	u := endpoint + "/" + bucket + "/" + key
	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	scope := dateStamp + "/auto/s3/aws4_request"

	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonHeaders := "host:" + hostOf(endpoint) + "\n" +
		"x-amz-content-sha256:UNSIGNED-PAYLOAD\n" +
		"x-amz-date:" + amzDate + "\n"
	canonReq := "PUT\n/" + bucket + "/" + key + "\n\n" +
		canonHeaders + "\n" + signedHeaders + "\nUNSIGNED-PAYLOAD"
	sum := sha256.Sum256([]byte(canonReq))
	toSign := "AWS4-HMAC-SHA256\n" + amzDate + "\n" + scope + "\n" + hex.EncodeToString(sum[:])

	k := []byte("AWS4" + secret)
	for _, part := range []string{dateStamp, "auto", "s3", "aws4_request"} {
		k = hmacSHA256(k, part)
	}
	signature := hex.EncodeToString(hmacSHA256(k, toSign))

	req, err := http.NewRequest(http.MethodPut, u, f)
	if err != nil {
		return err
	}
	req.ContentLength = info.Size()
	req.Header.Set("x-amz-content-sha256", "UNSIGNED-PAYLOAD")
	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("Authorization",
		"AWS4-HMAC-SHA256 Credential="+keyID+"/"+scope+
			", SignedHeaders="+signedHeaders+", Signature="+signature)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
	if res.StatusCode/100 != 2 {
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func hostOf(endpoint string) string {
	return strings.TrimPrefix(strings.TrimPrefix(endpoint, "https://"), "http://")
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, "upload_r2:", msg)
	os.Exit(1)
}
