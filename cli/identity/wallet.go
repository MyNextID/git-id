package identity

// Package identity provides simple utilities to create, load, and fetch Ed25519 identities.

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Identity holds an Ed25519 private and public key.
type Identity struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

// LoadOrCreate loads an identity from the given path, or generates a new one if it doesn't exist.
func LoadOrCreate(idPath string) (*Identity, error) {
	_, err := os.Stat(idPath)
	if err == nil {
		return ReadIdentity(idPath)
	}
	if os.IsNotExist(err) {
		fmt.Printf("Generating new identity at %s\n", idPath)
		return GenerateIdentity(idPath, false)
	}
	return nil, err
}

// ReadIdentity reads an existing identity from disk.
func ReadIdentity(path string) (*Identity, error) {
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read identity file: %w", err)
	}
	privKey, err := pemToPrivateKey(pemBytes)
	if err != nil {
		return nil, err
	}
	pubKey := privKey.Public().(ed25519.PublicKey)
	return &Identity{
		PrivateKey: privKey,
		PublicKey:  pubKey,
	}, nil
}

// GenerateIdentity creates a new Ed25519 identity and saves it locally.
// If force is true, it will overwrite existing files.
func GenerateIdentity(path string, force bool) (*Identity, error) {
	// Check private key file
	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			return nil, fmt.Errorf("path %q is a directory, not a file", path)
		}
		if !force {
			return nil, fmt.Errorf("file %q already exists", path)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat path %q: %w", path, err)
	}

	publicKeyPath := filepath.Join(filepath.Dir(path), "gid.pem")

	// Check public key file
	if info, err := os.Stat(publicKeyPath); err == nil && !info.IsDir() {
		if !force {
			return nil, fmt.Errorf("public key file %q already exists", publicKeyPath)
		}
	} else if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat public key path %q: %w", publicKeyPath, err)
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, fmt.Errorf("failed to create parent directory for %q: %w", path, err)
	}

	// Generate key pair
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Marshal and write private key
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}
	privPem := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	if err := os.WriteFile(path, privPem, 0600); err != nil {
		return nil, fmt.Errorf("failed to write private key to %q: %w", path, err)
	}

	// Marshal and write public key
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	pubPem := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})

	if err := os.WriteFile(publicKeyPath, pubPem, 0600); err != nil {
		return nil, fmt.Errorf("failed to write public key to %q: %w", publicKeyPath, err)
	}

	return &Identity{
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}

// FetchPublicKeyFromGitHub fetches a PEM-encoded public key from GitHub.
// handler should be like "user/repo", ref like "main", and keyPath like "gid/gid.pem".
func FetchPublicKeyFromGitHub(handler, ref, keyPath string) (ed25519.PublicKey, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", handler, ref, keyPath)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch public key: HTTP %d", resp.StatusCode)
	}

	pemData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key data: %w", err)
	}

	return pemToPublicKey(pemData)
}

// pemToPrivateKey parses a PEM-encoded Ed25519 private key.
func pemToPrivateKey(pemBytes []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM format for private key")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	privKey, ok := priv.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parsed key is not Ed25519 private key")
	}
	return privKey, nil
}

// pemToPublicKey parses a PEM-encoded Ed25519 public key.
func pemToPublicKey(pemBytes []byte) (ed25519.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM format for public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	pubKey, ok := pub.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("parsed key is not Ed25519 public key")
	}
	return pubKey, nil
}
