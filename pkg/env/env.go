// Package env fornece helpers pra ler variaveis de ambiente com fallback.
package env

import (
	"os"
	"strconv"
)

// GetOrDefault retorna o valor da env var ou o fallback se vazia.
func GetOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// GetIntOrDefault retorna o valor inteiro da env var ou o fallback se
// estiver vazia ou nao for um numero valido.
func GetIntOrDefault(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
