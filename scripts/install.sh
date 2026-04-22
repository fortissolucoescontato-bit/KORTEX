#!/usr/bin/env bash
set -euo pipefail

# ============================================================================
# kortex — Script de Instalação
# Um único comando para configurar qualquer agente de codificação IA em qualquer SO.
#
# Uso:
#   curl -sL https://raw.githubusercontent.com/fortissolucoescontato-bit/kortex/main/scripts/install.sh | bash
#
# Ou baixe e execute:
#   curl -sLO https://raw.githubusercontent.com/fortissolucoescontato-bit/kortex/main/scripts/install.sh
#   chmod +x install.sh
#   ./install.sh
# ============================================================================

GITHUB_OWNER="fortissolucoescontato-bit"
GITHUB_REPO="kortex"
BINARY_NAME="kortex"
BREW_TAP="fortissolucoescontato-bit/homebrew-tap"

# ============================================================================
# Suporte a cores
# ============================================================================

setup_colors() {
    if [ -t 1 ] && [ "${TERM:-}" != "dumb" ]; then
        RED='\033[0;31m'
        GREEN='\033[0;32m'
        YELLOW='\033[1;33m'
        BLUE='\033[0;34m'
        CYAN='\033[0;36m'
        BOLD='\033[1m'
        DIM='\033[2m'
        NC='\033[0m'
    else
        RED='' GREEN='' YELLOW='' BLUE='' CYAN='' BOLD='' DIM='' NC=''
    fi
}

# ============================================================================
# Helpers de logging
# ============================================================================

info()    { echo -e "${BLUE}[info]${NC}    $*"; }
success() { echo -e "${GREEN}[ok]${NC}      $*"; }
warn()    { echo -e "${YELLOW}[aviso]${NC}   $*"; }
error()   { echo -e "${RED}[erro]${NC}    $*" >&2; }
fatal()   { error "$@"; exit 1; }
step()    { echo -e "\n${CYAN}${BOLD}==>${NC} ${BOLD}$*${NC}"; }

# ============================================================================
# Ajuda
# ============================================================================

show_help() {
    cat <<EOF
${BOLD}instalador kortex${NC}

Uso: install.sh [OPÇÕES]

Opções:
  --method MÉTODO   Força o método de instalação: brew, go, binary (padrão: autodetectar)
  --dir DIR         Diretório de instalação personalizado para o método binary
  -h, --help        Mostra esta ajuda

Métodos de instalação (autodetectados por ordem de prioridade):
  1. brew    — Homebrew tap (recomendado)
  2. go      — go install a partir do código-fonte
  3. binary  — Binário pré-compilado do GitHub Releases

Exemplos:
  curl -sL https://raw.githubusercontent.com/${GITHUB_OWNER}/${GITHUB_REPO}/main/scripts/install.sh | bash
  ./install.sh --method binary
  ./install.sh --method binary --dir \$HOME/.local/bin

EOF
}

# ============================================================================
# Detecção de plataforma
# ============================================================================

detect_platform() {
    local uname_os uname_arch

    uname_os="$(uname -s)"
    uname_arch="$(uname -m)"

    case "$uname_os" in
        Darwin) OS="darwin"; OS_LABEL="macOS"; GORELEASER_OS="darwin" ;;
        Linux)  OS="linux";  OS_LABEL="Linux"; GORELEASER_OS="linux" ;;
        *)      fatal "SO não suportado: $uname_os. Apenas macOS e Linux são suportados." ;;
    esac

    case "$uname_arch" in
        x86_64|amd64)   ARCH="amd64" ;;
        arm64|aarch64)  ARCH="arm64" ;;
        *)              fatal "Arquitetura não suportada: $uname_arch. Apenas amd64 e arm64 são suportadas." ;;
    esac

    success "Plataforma: ${OS_LABEL} (${OS}/${ARCH})"
}

# ============================================================================
# Nomeação de arquivo GoReleaser
# ============================================================================

get_archive_name() {
    local version="$1"
    echo "${BINARY_NAME}_${version}_${GORELEASER_OS}_${ARCH}.tar.gz"
}

# ============================================================================
# Pré-requisitos
# ============================================================================

check_prerequisites() {
    step "Verificando pré-requisitos"

    local missing=()

    if ! command -v curl &>/dev/null; then
        missing+=("curl")
    fi

    if ! command -v git &>/dev/null; then
        missing+=("git")
    fi

    if [ ${#missing[@]} -gt 0 ]; then
        fatal "Ferramentas necessárias ausentes: ${missing[*]}. Instale-as e tente novamente."
    fi

    success "curl e git estão disponíveis"
}

# ============================================================================
# Detecção do método de instalação
# ============================================================================

detect_install_method() {
    if [ -n "${FORCE_METHOD:-}" ]; then
        case "$FORCE_METHOD" in
            brew|go|binary) INSTALL_METHOD="$FORCE_METHOD" ;;
            *) fatal "Método de instalação desconhecido: $FORCE_METHOD. Use: brew, go ou binary" ;;
        esac
        info "Usando método forçado: $INSTALL_METHOD"
        return
    fi

    step "Detectando o melhor método de instalação"

    if command -v brew &>/dev/null; then
        INSTALL_METHOD="brew"
        success "Homebrew encontrado — instalando via brew tap"
    else
        INSTALL_METHOD="binary"
        info "Baixando binário pré-compilado do GitHub Releases"
    fi
}

# ============================================================================
# Instalação via Homebrew
# ============================================================================

install_brew() {
    step "Instalando via Homebrew"

    info "Atualizando ${BREW_TAP}..."
    brew untap "$BREW_TAP" 2>/dev/null || true
    if ! brew tap "$BREW_TAP"; then
        fatal "Falha ao adicionar o tap $BREW_TAP"
    fi

    if brew list "$BINARY_NAME" &>/dev/null; then
        info "Já instalado, atualizando ${BINARY_NAME}..."
        if brew upgrade "$BINARY_NAME" 2>/dev/null; then
            success "${BINARY_NAME} atualizado via Homebrew"
        else
            success "${BINARY_NAME} já está na versão mais recente"
        fi
    else
        info "Instalando ${BINARY_NAME}..."
        if brew install "$BINARY_NAME"; then
            success "${BINARY_NAME} instalado via Homebrew"
        else
            fatal "Falha ao instalar ${BINARY_NAME} via Homebrew"
        fi
    fi
}

# ============================================================================
# Instalação via go install
# ============================================================================

install_go() {
    step "Instalando via go install"

    local go_package="github.com/${GITHUB_OWNER,,}/${GITHUB_REPO}/cmd/${BINARY_NAME}@latest"

    info "Executando: go install ${go_package}"
    if ! go install "$go_package"; then
        fatal "Falha ao instalar via go install. Verifique se o Go está configurado corretamente."
    fi

    local gobin
    gobin="$(go env GOBIN)"
    if [ -z "$gobin" ]; then
        gobin="$(go env GOPATH)/bin"
    fi

    if [[ ":$PATH:" != *":$gobin:"* ]]; then
        warn "${gobin} não está no seu PATH"
        warn "Adicione isto ao seu perfil do shell: export PATH=\"\$PATH:${gobin}\""
    fi

    success "${BINARY_NAME} instalado via go install"
}

# ============================================================================
# Instalação via download de binário
# ============================================================================

get_latest_version() {
    local url="https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPO}/releases/latest"

    info "Buscando o último lançamento no GitHub..."

    local response
    response="$(curl -sL -w "\n%{http_code}" "$url")" || fatal "Falha ao buscar o último lançamento"

    local http_code body
    http_code="$(echo "$response" | tail -n1)"
    body="$(echo "$response" | sed '$d')"

    if [ "$http_code" != "200" ]; then
        fatal "A API do GitHub retornou HTTP $http_code. Limite de requisições excedido? Tente novamente mais tarde ou use --method brew/go"
    fi

    LATEST_VERSION="$(echo "$body" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -1)"

    if [ -z "$LATEST_VERSION" ]; then
        fatal "Não foi possível determinar a última versão a partir da resposta da API do GitHub"
    fi

    VERSION_NUMBER="${LATEST_VERSION#v}"

    success "Última versão: ${LATEST_VERSION}"
}

install_binary() {
    step "Instalando binário pré-compilado"

    get_latest_version

    local archive_name
    archive_name="$(get_archive_name "$VERSION_NUMBER")"
    local download_url="https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}/releases/download/${LATEST_VERSION}/${archive_name}"
    local checksums_url="https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}/releases/download/${LATEST_VERSION}/checksums.txt"

    local tmpdir
    tmpdir="$(mktemp -d)"
    trap '[ -n "${tmpdir:-}" ] && rm -rf "$tmpdir"' EXIT

    info "Baixando ${archive_name}..."
    if ! curl -sfL -o "${tmpdir}/${archive_name}" "$download_url"; then
        fatal "Falha ao baixar ${download_url}"
    fi

    local file_size
    file_size="$(wc -c < "${tmpdir}/${archive_name}" | tr -d '[:space:]')"
    if [ "$file_size" -lt 1000 ]; then
        fatal "O arquivo baixado é suspeitosamente pequeno (${file_size} bytes). O pacote pode não existir para esta plataforma."
    fi

    success "${archive_name} baixado (${file_size} bytes)"

    info "Verificando checksum..."
    if curl -sL -o "${tmpdir}/checksums.txt" "$checksums_url"; then
        local expected_checksum
        expected_checksum="$(grep "${archive_name}" "${tmpdir}/checksums.txt" 2>/dev/null | awk '{print $1}' || true)"

        if [ -n "$expected_checksum" ]; then
            local actual_checksum
            if command -v sha256sum &>/dev/null; then
                actual_checksum="$(sha256sum "${tmpdir}/${archive_name}" | awk '{print $1}')"
            elif command -v shasum &>/dev/null; then
                actual_checksum="$(shasum -a 256 "${tmpdir}/${archive_name}" | awk '{print $1}')"
            else
                warn "sha256sum ou shasum não encontrados — pulando verificação de checksum"
                actual_checksum="$expected_checksum"
            fi

            if [ "$actual_checksum" != "$expected_checksum" ]; then
                fatal "Erro no checksum!\n  Esperado: ${expected_checksum}\n  Obtido:   ${actual_checksum}"
            fi
            success "Checksum verificado"
        else
            warn "Arquivo não encontrado no checksums.txt — pulando verificação"
        fi
    else
        warn "Não foi possível baixar o checksums.txt — pulando verificação"
    fi

    info "Extraindo ${BINARY_NAME}..."
    if ! tar -xzf "${tmpdir}/${archive_name}" -C "$tmpdir"; then
        fatal "Falha ao extrair o arquivo"
    fi

    if [ ! -f "${tmpdir}/${BINARY_NAME}" ]; then
        fatal "Binário '${BINARY_NAME}' não encontrado no pacote"
    fi

    local install_dir="${INSTALL_DIR:-}"

    if [ -z "$install_dir" ]; then
        if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
            install_dir="/usr/local/bin"
        elif [ "$(id -u)" = "0" ]; then
            install_dir="/usr/local/bin"
        else
            install_dir="${HOME}/.local/bin"
        fi
    fi

    mkdir -p "$install_dir"

    info "Instalando em ${install_dir}/${BINARY_NAME}..."
    if cp "${tmpdir}/${BINARY_NAME}" "${install_dir}/${BINARY_NAME}" 2>/dev/null; then
        chmod +x "${install_dir}/${BINARY_NAME}"
    elif command -v sudo &>/dev/null; then
        warn "Permissão negada. Tentando com sudo..."
        sudo cp "${tmpdir}/${BINARY_NAME}" "${install_dir}/${BINARY_NAME}"
        sudo chmod +x "${install_dir}/${BINARY_NAME}"
    else
        fatal "Não é possível gravar em ${install_dir}. Execute com sudo ou use --dir para especificar um diretório com permissão de escrita."
    fi

    success "${BINARY_NAME} instalado em ${install_dir}/${BINARY_NAME}"

    if [[ ":$PATH:" != *":${install_dir}:"* ]]; then
        warn "${install_dir} não está no seu PATH"
        echo ""
        warn "Adicione isto ao seu perfil do shell (~/.bashrc, ~/.zshrc, etc.):"
        echo -e "  ${DIM}export PATH=\"\$PATH:${install_dir}\"${NC}"
        echo ""
    fi
}

# ============================================================================
# Verificar instalação
# ============================================================================

verify_installation() {
    step "Verificando instalação"

    hash -r 2>/dev/null || true

    if command -v "$BINARY_NAME" &>/dev/null; then
        local version_output
        version_output="$("$BINARY_NAME" version 2>&1 || true)"
        success "${BINARY_NAME} está instalado: ${version_output}"
        return 0
    fi

    local locations=(
        "/usr/local/bin/${BINARY_NAME}"
        "${HOME}/.local/bin/${BINARY_NAME}"
        "$(go env GOPATH 2>/dev/null || echo "")/bin/${BINARY_NAME}"
    )

    for loc in "${locations[@]}"; do
        if [ -n "$loc" ] && [ -x "$loc" ]; then
            local version_output
            version_output="$("$loc" version 2>&1 || true)"
            success "Encontrado ${BINARY_NAME} em ${loc}: ${version_output}"
            warn "A localização do binário não está no seu PATH. Adicione-a para usar o '${BINARY_NAME}' diretamente."
            return 0
        fi
    done

    warn "Não foi possível verificar a instalação. Você pode precisar reiniciar o seu shell."
    return 0
}

# ============================================================================
# Imprimir próximos passos
# ============================================================================

print_banner() {
    echo ""
    echo -e "${CYAN}${BOLD}"
    echo "   ____            _   _              _    ___ "
    echo "  / ___| ___ _ __ | |_| | ___        / \  |_ _|"
    echo " | |  _ / _ \ '_ \| __| |/ _ \_____ / _ \  | | "
    echo " | |_| |  __/ | | | |_| |  __/_____/ ___ \ | | "
    echo "  \____|\___|_| |_|\__|_|\___|    /_/   \_\___|"
    echo -e "${NC}"
    echo -e "  ${DIM}Um único comando para configurar qualquer agente de codificação IA em qualquer SO${NC}"
    echo ""
}

print_next_steps() {
    echo ""
    echo -e "${GREEN}${BOLD}Instalação concluída!${NC}"
    echo ""
    echo -e "${BOLD}Próximos passos:${NC}"
    echo -e "  ${CYAN}1.${NC} Execute ${BOLD}${BINARY_NAME}${NC} para iniciar o instalador TUI"
    echo -e "  ${CYAN}2.${NC} Selecione seu(s) agente(s) de IA e ferramentas para configurar"
    echo -e "  ${CYAN}3.${NC} Siga as instruções interativas"
    echo ""
    echo -e "${DIM}Para ajuda: ${BINARY_NAME} --help${NC}"
    echo -e "${DIM}Docs:       https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}${NC}"
    echo ""
}

# ============================================================================
# Main
# ============================================================================

main() {
    setup_colors

    FORCE_METHOD=""
    INSTALL_DIR=""

    while [ $# -gt 0 ]; do
        case "$1" in
            --method)
                [ $# -lt 2 ] && fatal "--method requer um argumento"
                FORCE_METHOD="$2"; shift 2
                ;;
            --dir)
                [ $# -lt 2 ] && fatal "--dir requer um argumento"
                INSTALL_DIR="$2"; shift 2
                ;;
            -h|--help)
                setup_colors
                show_help
                exit 0
                ;;
            *)
                fatal "Opção desconhecida: $1. Use --help para ver o uso."
                ;;
        esac
    done

    print_banner

    step "Detectando plataforma"
    detect_platform

    check_prerequisites
    detect_install_method

    case "$INSTALL_METHOD" in
        brew)   install_brew ;;
        go)     install_go ;;
        binary) install_binary ;;
    esac

    verify_installation
    print_next_steps
}

main "$@"
