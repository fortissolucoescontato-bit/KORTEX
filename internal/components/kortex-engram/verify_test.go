package kortexengram

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
)

func TestVerifyInstalled(t *testing.T) {
	original := lookPath
	t.Cleanup(func() { lookPath = original })

	lookPath = func(string) (string, error) { return "/opt/homebrew/bin/KortexEngram", nil }
	if err := VerifyInstalled(); err != nil {
		t.Fatalf("VerifyInstalled() error = %v", err)
	}

	lookPath = func(string) (string, error) { return "", errors.New("missing") }
	if err := VerifyInstalled(); err == nil {
		t.Fatalf("VerifyInstalled() expected missing binary error")
	}
}

func TestVerifyHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	if err := VerifyHealth(context.Background(), server.URL); err != nil {
		t.Fatalf("VerifyHealth() error = %v", err)
	}

	badServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer badServer.Close()

	if err := VerifyHealth(context.Background(), badServer.URL); err == nil {
		t.Fatalf("VerifyHealth() expected non-200 error")
	}
}

func TestVerifyVersion(t *testing.T) {
	originalLook := lookPath
	originalExec := execCommand
	t.Cleanup(func() {
		lookPath = originalLook
		execCommand = originalExec
	})

	lookPath = func(string) (string, error) { return "/usr/bin/kortex-engram", nil }

	// RED -> GREEN: Mockando execCommand para retornar versão fictícia
	execCommand = func(name string, arg ...string) *exec.Cmd {
		// Padrão Go para mock de exec.Command: retornar um comando que executa o próprio teste
		// com uma flag específica para agir como o binário mockado.
		// Mas aqui simplificaremos apenas injetando o comportamento via override se possível.
		// Na verdade, a implementação de verify.go usa execCommand(cmdName, "version").
		// Vou simplificar o verify.go para permitir injeção direta de Output se necessário, 
		// ou usar o padrão oficial de TestProcessResponse.
		return exec.Command("echo", "v1.2.3")
	}

	version, err := VerifyVersion()
	if err != nil {
		t.Fatalf("VerifyVersion() error = %v", err)
	}
	if version != "v1.2.3" {
		t.Fatalf("VerifyVersion() = %q, want %q", version, "v1.2.3")
	}
}
