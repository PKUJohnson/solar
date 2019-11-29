package fileutil

import (
	"encoding/json"

	std "github.com/PKUJohnson/solar/std"
	tucrypto "github.com/PKUJohnson/solar/std/toolkit/crypto"
)

const (
	fileEncryptSecret = "file_4uma8-9dj0e"
)

type FileInfo struct {
	Bucket   string
	Key      string
	ViewerId int64
	Deadline int64
	Preview  bool
}

func NewFileInfo(bucket, key string, viewerId, deadline int64) *FileInfo {
	return &FileInfo{
		Bucket:   bucket,
		Key:      key,
		ViewerId: viewerId,
		Deadline: deadline,
	}
}

func NewPreviewFileInfo(bucket, key string, viewerId, deadline int64) *FileInfo {
	return &FileInfo{
		Bucket:   bucket,
		Key:      key,
		ViewerId: viewerId,
		Deadline: deadline,
		Preview:  true,
	}
}

func FileInfoFromToken(content string) (*FileInfo, error) {
	data, err := tucrypto.Base64Decode(content)
	if err != nil {
		std.LogInfoLn("FileInfo ParseToken Base64Decode Error:", err)
		return nil, std.ErrBadFileToken
	}
	data, err = tucrypto.Decrypt(data, []byte(fileEncryptSecret))
	if err != nil {
		std.LogInfoLn("FileInfo ParseToken Decrypt Error:", err)
		return nil, std.ErrBadFileToken
	}
	fileInfo := new(FileInfo)
	if err = json.Unmarshal(data, &fileInfo); err != nil {
		std.LogInfoLn("FileInfo ParseToken Unmarshal Error:", err)
		return nil, std.ErrBadFileToken
	}
	return fileInfo, nil
}

func (fi *FileInfo) FileTokenString() string {
	data, err := json.Marshal(fi)
	if err != nil {
		std.LogInfoLn("FileTokenString Error:", err)
		return ""
	}
	data, err = tucrypto.Encrypt(data, []byte(fileEncryptSecret))
	if err != nil {
		std.LogInfoLn("FileTokenString Encrypt Error:", err)
		return ""
	}
	return tucrypto.Base64Encode(data)
}
