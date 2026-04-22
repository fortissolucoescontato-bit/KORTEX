package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ArchiveEntry describes a single file to include in a tar.gz archive (for
// creation) or a file that was extracted from one (for extraction).
type ArchiveEntry struct {
	// RelPath is the relative path inside the archive (e.g. "files/.config/foo.json").
	RelPath string
	// SourcePath is the absolute path on disk (source when creating; destination when extracting).
	SourcePath string
	// Mode is the file permission bits.
	Mode os.FileMode
}

// CreateArchive writes a tar.gz archive at archivePath containing all entries.
// Each entry's content is read from SourcePath and stored under RelPath inside
// the archive. Intermediate directories in archivePath are created automatically.
//
// The archive uses filepath.ToSlash for cross-platform tar path compatibility.
func CreateArchive(archivePath string, entries []ArchiveEntry) error {
	if err := os.MkdirAll(filepath.Dir(archivePath), 0o755); err != nil {
		return fmt.Errorf("falha ao criar diretório do arquivo %q: %w", filepath.Dir(archivePath), err)
	}

	f, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo compactado %q: %w", archivePath, err)
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)

	for _, entry := range entries {
		if err := addFileToTar(tw, entry); err != nil {
			return err
		}
	}

	// Flush gzip writer before the deferred closes fire; any error here is surfaced.
	if err := tw.Close(); err != nil {
		return fmt.Errorf("falha ao fechar gravador tar: %w", err)
	}
	if err := gw.Close(); err != nil {
		return fmt.Errorf("falha ao fechar gravador gzip: %w", err)
	}

	return nil
}

// addFileToTar appends a single file to the tar writer.
func addFileToTar(tw *tar.Writer, entry ArchiveEntry) error {
	data, err := os.ReadFile(entry.SourcePath)
	if err != nil {
		return fmt.Errorf("falha ao ler arquivo de origem %q: %w", entry.SourcePath, err)
	}

	hdr := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     filepath.ToSlash(entry.RelPath),
		Mode:     int64(entry.Mode),
		Size:     int64(len(data)),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("falha ao gravar cabeçalho tar para %q: %w", entry.RelPath, err)
	}

	if _, err := tw.Write(data); err != nil {
		return fmt.Errorf("falha ao gravar conteúdo tar para %q: %w", entry.RelPath, err)
	}

	return nil
}

// ExtractArchive extracts a tar.gz archive at archivePath into destDir.
// It returns the list of extracted entries (RelPath = path inside archive,
// SourcePath = absolute path of the extracted file inside destDir).
//
// Path traversal protection: any entry whose RelPath contains ".." is rejected
// with an error and extraction is aborted.
func ExtractArchive(archivePath string, destDir string) ([]ArchiveEntry, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir arquivo compactado %q: %w", archivePath, err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar leitor gzip para %q: %w", archivePath, err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	var extracted []ArchiveEntry

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("falha ao ler entrada tar: %w", err)
		}

		// FIX-1: Only process regular files; skip symlinks, hardlinks, directories, etc.
		if hdr.Typeflag != tar.TypeReg && hdr.Typeflag != tar.TypeRegA {
			continue
		}

		destPath := filepath.Join(destDir, filepath.FromSlash(hdr.Name))

		// Path traversal protection: verify the cleaned destination is inside destDir.
		cleanDest := filepath.Clean(destPath)
		cleanBase := filepath.Clean(destDir) + string(filepath.Separator)
		if !strings.HasPrefix(cleanDest, cleanBase) {
			return nil, fmt.Errorf("entrada do arquivo %q escapa do diretório de destino", hdr.Name)
		}

		// FIX-3: Reject "." entry — resolves to destDir itself, not a file inside it.
		if cleanDest == filepath.Clean(destDir) {
			return nil, fmt.Errorf("entrada do arquivo %q resolve para o próprio diretório de destino", hdr.Name)
		}

		mode := os.FileMode(hdr.Mode)

		// Create intermediate directories.
		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return nil, fmt.Errorf("falha ao criar diretórios para %q: %w", destPath, err)
		}

		outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode.Perm())
		if err != nil {
			return nil, fmt.Errorf("falha ao criar arquivo extraído %q: %w", destPath, err)
		}

		if _, err := io.Copy(outFile, tr); err != nil {
			_ = outFile.Close()
			return nil, fmt.Errorf("falha ao gravar arquivo extraído %q: %w", destPath, err)
		}

		if err := outFile.Close(); err != nil {
			return nil, fmt.Errorf("falha ao fechar arquivo extraído %q: %w", destPath, err)
		}

		extracted = append(extracted, ArchiveEntry{
			RelPath:    hdr.Name,
			SourcePath: destPath,
			Mode:       mode,
		})
	}

	return extracted, nil
}
