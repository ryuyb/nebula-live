package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordConfig 密码配置参数
type PasswordConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultPasswordConfig 默认密码配置
var DefaultPasswordConfig = &PasswordConfig{
	Memory:      64 * 1024, // 64MB
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// HashPassword 使用 Argon2id 哈希密码
func HashPassword(password string) (string, error) {
	return HashPasswordWithConfig(password, DefaultPasswordConfig)
}

// HashPasswordWithConfig 使用指定配置哈希密码
func HashPasswordWithConfig(password string, config *PasswordConfig) (string, error) {
	// 生成随机盐
	salt, err := generateRandomBytes(config.SaltLength)
	if err != nil {
		return "", err
	}

	// 生成哈希
	hash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	// 编码为 base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 格式: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, config.Memory, config.Iterations, config.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, encodedHash string) (bool, error) {
	// 解析编码的哈希
	config, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// 使用相同参数哈希输入的密码
	otherHash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	// 使用常数时间比较防止时序攻击
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

// generateRandomBytes 生成指定长度的随机字节
func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// decodeHash 解析编码的哈希值
func decodeHash(encodedHash string) (config *PasswordConfig, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, fmt.Errorf("invalid hash format")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, fmt.Errorf("incompatible version of argon2")
	}

	config = &PasswordConfig{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &config.Memory, &config.Iterations, &config.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	config.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	config.KeyLength = uint32(len(hash))

	return config, salt, hash, nil
}
