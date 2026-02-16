# GameKeep - Game Save Manager

GameKeep é um gerenciador profissional de checkpoints de saves de jogos com interface gráfica moderna e CLI opcional.

![GameKeep](https://img.shields.io/badge/version-1.0.0-blue)
![Go](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## 🎯 Características

- ✅ **Interface Gráfica Moderna** - UI intuitiva com Fyne
- ✅ **Gerenciamento de Múltiplos Jogos** - Organize saves de vários jogos
- ✅ **Checkpoints Compactados** - Saves salvos em ZIP
- ✅ **Restauração Rápida** - Volte a qualquer ponto anterior
- ✅ **Verificação de Integridade** - SHA256 para garantir dados íntegros
- ✅ **Metadados em JSON** - Configuração simples e portável
- ✅ **CLI Opcional** - Para usuários avançados e automação
- ✅ **Cross-platform** - Windows, macOS e Linux

## 📦 Instalação

### Pré-requisitos

- Go 1.21 ou superior
- Dependências do Fyne (para GUI):
  - **Linux**: `sudo apt-get install libgl1-mesa-dev xorg-dev`
  - **macOS**: Xcode command line tools
  - **Windows**: Nenhuma dependência adicional

### Build

```bash
# Clone o repositório
git clone https://github.com/gamekeep/gamekeep.git
cd gamekeep

# Baixar dependências
go mod download

# Build GUI e CLI
make build
```

## 🚀 Uso

### Interface Gráfica (Recomendado)

```bash
./gamekeep-gui
```

### CLI

```bash
# Adicionar jogo
gamekeep add-game --name "The Witcher 3" --path "/path/to/saves"

# Criar checkpoint
gamekeep checkpoint --game witcher3 --name "Before Boss" --note "Level 25"

# Listar checkpoints
gamekeep list --game witcher3

# Restaurar
gamekeep restore --checkpoint <id>
```

## 📁 Estrutura

```
~/.gamekeep/
├── config/
│   ├── games.json
│   └── checkpoints.json
└── vault/
    └── {game_id}/
        └── {checkpoint_id}.zip
```

Para mais detalhes, veja a documentação completa.
