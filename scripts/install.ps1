#Requires -Version 5.1
<#
.SYNOPSIS
    kortex — Script de Instalação para Windows
    Um único comando para configurar qualquer agente de codificação IA em qualquer SO.

.DESCRIPTION
    Baixa e instala o binário kortex para Windows.
    Suporta instalação via Go ou binário pré-compilado do GitHub Releases.

.EXAMPLE
    # Executar diretamente:
    irm https://raw.githubusercontent.com/fortissolucoescontato-bit/kortex/main/scripts/install.ps1 | iex

    # Ou baixar e executar:
    Invoke-WebRequest -Uri https://raw.githubusercontent.com/fortissolucoescontato-bit/kortex/main/scripts/install.ps1 -OutFile install.ps1
    .\install.ps1

    # Forçar um método específico:
    .\install.ps1 -Method binary
    .\install.ps1 -Method go
#>

[CmdletBinding()]
param(
    [ValidateSet("auto", "go", "binary")]
    [string]$Method = "auto",

    [string]$InstallDir = ""
)

$ErrorActionPreference = "Stop"

$GITHUB_OWNER = "fortissolucoescontato-bit"
$GITHUB_REPO = "kortex"
$BINARY_NAME = "kortex"

# ============================================================================
# Helpers de logging
# ============================================================================

function Write-Info    { param([string]$Message) Write-Host "[info]    $Message" -ForegroundColor Blue }
function Write-Success { param([string]$Message) Write-Host "[ok]      $Message" -ForegroundColor Green }
function Write-Warn    { param([string]$Message) Write-Host "[aviso]   $Message" -ForegroundColor Yellow }
function Write-Err     { param([string]$Message) Write-Host "[erro]    $Message" -ForegroundColor Red }
function Write-Step    { param([string]$Message) Write-Host "`n==> $Message" -ForegroundColor Cyan }

function Stop-WithError {
    param([string]$Message)
    Write-Err $Message
    exit 1
}

# ============================================================================
# Banner
# ============================================================================

function Show-Banner {
    Write-Host ""
    Write-Host "   ____            _   _              _    ___ " -ForegroundColor Cyan
    Write-Host "  / ___| ___ _ __ | |_| | ___        / \  |_ _|" -ForegroundColor Cyan
    Write-Host " | |  _ / _ \ '_ \| __| |/ _ \_____ / _ \  | | " -ForegroundColor Cyan
    Write-Host " | |_| |  __/ | | | |_| |  __/_____/ ___ \ | | " -ForegroundColor Cyan
    Write-Host "  \____|\___|_| |_|\__|_|\___|    /_/   \_\___|" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "  Um único comando para configurar qualquer agente de codificação IA em qualquer SO" -ForegroundColor DarkGray
    Write-Host ""
}

# ============================================================================
# Detecção de plataforma
# ============================================================================

function Get-Platform {
    Write-Step "Detectando plataforma"

    $arch = if ([Environment]::Is64BitOperatingSystem) {
        if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
    } else {
        Stop-WithError "Windows 32-bit não é suportado."
    }

    Write-Success "Plataforma: Windows ($arch)"
    return $arch
}

# ============================================================================
# Pré-requisitos
# ============================================================================

function Test-Prerequisites {
    Write-Step "Verificando pré-requisitos"

    $missing = @()
    if (-not (Get-Command "curl" -ErrorAction SilentlyContinue)) { $missing += "curl" }
    if (-not (Get-Command "git" -ErrorAction SilentlyContinue))  { $missing += "git" }

    if ($missing.Count -gt 0) {
        Stop-WithError "Ferramentas necessárias ausentes: $($missing -join ', '). Instale-as e tente novamente."
    }

    Write-Success "curl e git estão disponíveis"
}

# ============================================================================
# Detecção do método de instalação
# ============================================================================

function Get-InstallMethod {
    param([string]$Forced)

    if ($Forced -ne "auto") {
        Write-Info "Usando método forçado: $Forced"
        return $Forced
    }

    Write-Step "Detectando o melhor método de instalação"

    Write-Info "Baixando binário pré-compilado do GitHub Releases"
    return "binary"
}

# ============================================================================
# Instalação via go install
# ============================================================================

function Install-ViaGo {
    Write-Step "Instalando via go install"

    $goPackage = "github.com/$($GITHUB_OWNER.ToLower())/$GITHUB_REPO/cmd/$BINARY_NAME@latest"
    Write-Info "Executando: go install $goPackage"

    & go install $goPackage
    if ($LASTEXITCODE -ne 0) {
        Stop-WithError "Falha ao instalar via go install. Verifique se o Go está configurado corretamente."
    }

    $gobin = & go env GOBIN 2>$null
    if (-not $gobin) {
        $gopath = & go env GOPATH 2>$null
        $gobin = Join-Path $gopath "bin"
    }

    if ($env:PATH -notlike "*$gobin*") {
        Write-Warn "$gobin não está no seu PATH"
        Write-Warn "Adicione-o à sua variável de ambiente PATH."
    }

    Write-Success "$BINARY_NAME instalado via go install"
}

# ============================================================================
# Instalação via download de binário
# ============================================================================

function Get-LatestVersion {
    Write-Info "Buscando o último lançamento no GitHub..."

    $url = "https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/releases/latest"

    try {
        $response = Invoke-RestMethod -Uri $url -Headers @{ "User-Agent" = "kortex-installer" }
    } catch {
        Stop-WithError "Falha ao buscar o último lançamento. Limite de requisições excedido? Tente novamente mais tarde ou use -Method go"
    }

    $version = $response.tag_name
    if (-not $version) {
        Stop-WithError "Não foi possível determinar a última versão a partir da resposta da API do GitHub"
    }

    Write-Success "Última versão: $version"
    return $version
}

function Install-ViaBinary {
    param([string]$Arch)

    Write-Step "Instalando binário pré-compilado"

    $version = Get-LatestVersion
    $versionNumber = $version.TrimStart("v")

    $archiveName = "${BINARY_NAME}_${versionNumber}_windows_${Arch}.zip"
    $downloadUrl = "https://github.com/$GITHUB_OWNER/$GITHUB_REPO/releases/download/$version/$archiveName"
    $checksumsUrl = "https://github.com/$GITHUB_OWNER/$GITHUB_REPO/releases/download/$version/checksums.txt"

    $tmpDir = Join-Path $env:TEMP "kortex-install-$(Get-Random)"
    New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

    try {
        # Download archive
        Write-Info "Baixando $archiveName..."
        $archivePath = Join-Path $tmpDir $archiveName
        Invoke-WebRequest -Uri $downloadUrl -OutFile $archivePath -UseBasicParsing

        $fileSize = (Get-Item $archivePath).Length
        if ($fileSize -lt 1000) {
            Stop-WithError "O arquivo baixado é suspeitosamente pequeno ($fileSize bytes). O pacote pode não existir para esta plataforma."
        }
        Write-Success "$archiveName baixado ($fileSize bytes)"

        # Verify checksum
        Write-Info "Verificando checksum..."
        try {
            $checksumsPath = Join-Path $tmpDir "checksums.txt"
            Invoke-WebRequest -Uri $checksumsUrl -OutFile $checksumsPath -UseBasicParsing

            $checksums = Get-Content $checksumsPath
            $expectedLine = $checksums | Where-Object { $_ -match $archiveName }
            if ($expectedLine) {
                $expectedChecksum = ($expectedLine -split "\s+")[0]
                $actualChecksum = (Get-FileHash -Path $archivePath -Algorithm SHA256).Hash.ToLower()

                if ($actualChecksum -ne $expectedChecksum) {
                    Stop-WithError "Erro no checksum!`n  Esperado: $expectedChecksum`n  Obtido:   $actualChecksum"
                }
                Write-Success "Checksum verificado"
            } else {
                Write-Warn "Arquivo não encontrado no checksums.txt — pulando verificação"
            }
        } catch {
            Write-Warn "Não foi possível baixar o checksums.txt — pulando verificação"
        }

        # Extract binary
        Write-Info "Extraindo $BINARY_NAME..."
        Expand-Archive -Path $archivePath -DestinationPath $tmpDir -Force

        $binaryPath = Join-Path $tmpDir "$BINARY_NAME.exe"
        if (-not (Test-Path $binaryPath)) {
            Stop-WithError "Binário '$BINARY_NAME.exe' não encontrado no pacote"
        }

        # Determine install directory
        $installDir = $InstallDir
        if (-not $installDir) {
            $installDir = Join-Path $env:LOCALAPPDATA "kortex\bin"
        }

        if (-not (Test-Path $installDir)) {
            New-Item -ItemType Directory -Path $installDir -Force | Out-Null
        }

        # Install binary
        $destPath = Join-Path $installDir "$BINARY_NAME.exe"
        Write-Info "Instalando em $destPath..."
        Copy-Item -Path $binaryPath -Destination $destPath -Force

        Write-Success "$BINARY_NAME instalado em $destPath"

        # Check if install dir is in PATH
        if ($env:PATH -notlike "*$installDir*") {
            Write-Warn "$installDir não está no seu PATH"
            Write-Host ""
            Write-Warn "Execute isto para adicioná-lo permanentemente:"
            Write-Host "  [Environment]::SetEnvironmentVariable('PATH', `$env:PATH + ';$installDir', 'User')" -ForegroundColor DarkGray
            Write-Host ""
        }
    } finally {
        Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

# ============================================================================
# Verificar instalação
# ============================================================================

function Test-Installation {
    Write-Step "Verificando instalação"

    # Refresh PATH for current session
    $env:PATH = [Environment]::GetEnvironmentVariable("PATH", "Machine") + ";" + [Environment]::GetEnvironmentVariable("PATH", "User")

    $cmd = Get-Command $BINARY_NAME -ErrorAction SilentlyContinue
    if ($cmd) {
        $versionOutput = & $BINARY_NAME version 2>&1
        Write-Success "$BINARY_NAME está instalado: $versionOutput"
        return
    }

    # Check common locations
    $gopath = $null
    if (Get-Command "go" -ErrorAction SilentlyContinue) {
        $gopath = & go env GOPATH 2>$null
    }
    $locations = @(
        (Join-Path $env:LOCALAPPDATA "kortex\bin\$BINARY_NAME.exe")
    )
    if ($gopath) {
        $locations += (Join-Path $gopath "bin\$BINARY_NAME.exe")
    }

    foreach ($loc in $locations) {
        if ($loc -and (Test-Path $loc)) {
            $versionOutput = & $loc version 2>&1
            Write-Success "Encontrado $BINARY_NAME em $loc`: $versionOutput"
            Write-Warn "A localização do binário não está no seu PATH. Adicione-a para usar o '$BINARY_NAME' diretamente."
            return
        }
    }

    Write-Warn "Não foi possível verificar a instalação. Você pode precisar reiniciar o seu terminal."
}

# ============================================================================
# Próximos passos
# ============================================================================

function Show-NextSteps {
    Write-Host ""
    Write-Host "Instalação concluída!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Próximos passos:" -ForegroundColor White
    Write-Host "  1. Execute '$BINARY_NAME' para iniciar o instalador TUI" -ForegroundColor Cyan
    Write-Host "  2. Selecione seu(s) agente(s) de IA e ferramentas para configurar" -ForegroundColor Cyan
    Write-Host "  3. Siga as instruções interativas" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Para ajuda: $BINARY_NAME --help" -ForegroundColor DarkGray
    Write-Host "Docs:       https://github.com/$GITHUB_OWNER/$GITHUB_REPO" -ForegroundColor DarkGray
    Write-Host ""
}

# ============================================================================
# Main
# ============================================================================

function Main {
    Show-Banner

    $arch = Get-Platform
    Test-Prerequisites

    $installMethod = Get-InstallMethod -Forced $Method

    switch ($installMethod) {
        "go"     { Install-ViaGo }
        "binary" { Install-ViaBinary -Arch $arch }
    }

    Test-Installation
    Show-NextSteps
}

Main
