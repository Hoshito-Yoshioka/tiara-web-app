package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// allowedImageExtensions は許可される画像ファイルの拡張子。
var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// allowedImageMIMETypes は許可される画像ファイルの MIME タイプ。
var allowedImageMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// maxImageSize はアップロード可能な画像の最大サイズ (10MB)。
const maxImageSize = 10 * 1024 * 1024

// validateImageFile はアップロードされた画像ファイルのバリデーションを行う。
// 拡張子のホワイトリストチェックとマジックバイト（MIME タイプ）チェックを実施する。
func validateImageFile(file *multipart.FileHeader) error {
	// ファイルサイズチェック
	if file.Size > maxImageSize {
		return fmt.Errorf("file size must be less than 10MB")
	}

	// 拡張子チェック
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExtensions[ext] {
		return fmt.Errorf("allowed file types: jpg, jpeg, png, gif, webp")
	}

	// マジックバイト（MIME タイプ）チェック
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to read file")
	}
	defer src.Close() //nolint:errcheck // best-effort close

	// 先頭 512 バイトを読み取って MIME タイプを判定
	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file")
	}

	mimeType := http.DetectContentType(buf[:n])
	if !allowedImageMIMETypes[mimeType] {
		return fmt.Errorf("file content does not match an allowed image type")
	}

	// ファイルポインタを先頭に戻す
	if seeker, ok := src.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	return nil
}
