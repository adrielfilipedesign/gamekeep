# GameKeep - Guia Rápido de Build

## 🚀 Build Rápido

### 1. Extrair o projeto
```bash
tar -xzf gamekeep.tar.gz
cd gamekeep
```

### 2. Instalar dependências do Fyne (apenas para GUI)

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get install libgl1-mesa-dev xorg-dev
```

**macOS:**
```bash
xcode-select --install
```

**Windows:**
Nenhuma dependência adicional necessária.

### 3. Build

```bash
# Baixar dependências Go
go mod download

# Build GUI + CLI
make build

# Ou build separadamente:
make build-gui  # Cria: gamekeep-gui
make build-cli  # Cria: gamekeep
```

### 4. Executar

```bash
# GUI (recomendado)
./gamekeep-gui

# CLI
./gamekeep help
```

## 📦 Estrutura do Projeto

```
gamekeep/
├── cmd/
│   ├── gamekeep-gui/main.go    # GUI application
│   └── main.go                 # CLI application
├── internal/
│   ├── core/                   # Business logic
│   ├── models/                 # Data structures
│   ├── storage/                # JSON persistence
│   └── vault/                  # File operations
├── ui/                         # GUI components (Fyne)
│   ├── main_ui.go
│   ├── games_view.go
│   ├── checkpoints_view.go
│   ├── dialogs.go
│   └── theme.go
├── go.mod                      # Dependencies
├── Makefile                    # Build commands
└── README.md                   # Documentation
```

## 🎯 Comandos Makefile

```bash
make build          # Build GUI + CLI
make build-gui      # Build apenas GUI
make build-cli      # Build apenas CLI
make build-all      # Build multi-plataforma
make run-gui        # Run GUI em modo dev
make run-cli        # Run CLI em modo dev
make test           # Run testes
make clean          # Limpar builds
make install        # Instalar em $GOPATH/bin
make help           # Ver todos comandos
```

## 🖥️ Dependências

### Go Modules (gerenciadas automaticamente)
- `fyne.io/fyne/v2` - GUI framework
- `github.com/google/uuid` - UUID generation

### Sistema (apenas para GUI)
- **Linux**: OpenGL, X11
- **macOS**: Xcode command line tools
- **Windows**: Nenhuma

## 🔧 Troubleshooting

### Erro: "fyne not found"
```bash
go mod download
```

### Erro: "GL/gl.h: No such file" (Linux)
```bash
sudo apt-get install libgl1-mesa-dev xorg-dev
```

### Erro: "xcrun: error" (macOS)
```bash
xcode-select --install
```

### Build apenas CLI (sem dependências GUI)
```bash
go build -o gamekeep cmd/main.go
```

## 📝 Primeiro Uso

### Via GUI
1. Execute `./gamekeep-gui`
2. Clique em "➕ Add Game"
3. Preencha nome e path
4. Selecione o jogo
5. Clique em "➕ Create Checkpoint"

### Via CLI
```bash
# 1. Adicionar jogo
./gamekeep add-game --name "Meu Jogo" --path "/path/to/saves"

# 2. Criar checkpoint
./gamekeep checkpoint --game meu_jogo --name "Save Inicial"

# 3. Listar checkpoints
./gamekeep list --game meu_jogo

# 4. Restaurar (se necessário)
./gamekeep restore --checkpoint <id>
```

## 🎮 Exemplo Completo

```bash
# Build
make build

# Adicionar um jogo
./gamekeep add-game --name "Dark Souls 3" --path "/home/user/.steam/steam/userdata/123/374320/remote"

# Criar checkpoint antes de boss
./gamekeep checkpoint --game dark_souls_3 --name "Pre Nameless King" --note "SL 80"

# Listar checkpoints
./gamekeep list --game dark_souls_3

# Se morrer muito, restaurar :)
./gamekeep restore --checkpoint abc12345

# Ou usar a GUI para tudo isso!
./gamekeep-gui
```

## 📚 Documentação Completa

- `README.md` - Overview e features
- `QUICKSTART.md` - Guia passo-a-passo
- `DEVELOPER.md` - Arquitetura técnica
- `CONTRIBUTING.md` - Como contribuir

## 🆘 Suporte

- GitHub Issues: Para bugs e features
- Documentação: Veja os arquivos .md no projeto

---

**Boa sorte e não perca mais saves! 🎮**