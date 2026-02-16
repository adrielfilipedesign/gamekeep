# SaveSync - Estrutura do Projeto

## 📂 Estrutura Completa

```
savesync/
│
├── cmd/
│   └── main.go                    # Entry point da aplicação
│
├── internal/
│   ├── core/
│   │   └── service.go             # Lógica de negócio principal
│   │
│   ├── models/
│   │   ├── models.go              # Definições de Game e Checkpoint
│   │   └── errors.go              # Erros customizados
│   │
│   ├── storage/
│   │   └── json_store.go          # Persistência em JSON com escrita atômica
│   │
│   └── vault/
│       └── manager.go             # Gerenciamento de arquivos ZIP e hashing
│
├── config/                         # (Gerado em runtime)
│   ├── games.json                 # Lista de jogos registrados
│   └── checkpoints.json           # Lista de checkpoints
│
├── vault/                          # (Gerado em runtime)
│   └── {game_id}/
│       └── {checkpoint_id}.zip    # Arquivos de save compactados
│
├── go.mod                          # Dependências do Go
├── Makefile                        # Comandos de build e teste
├── .gitignore                      # Arquivos ignorados pelo Git
├── LICENSE                         # Licença MIT
│
├── README.md                       # Documentação principal
├── QUICKSTART.md                   # Guia de início rápido
├── DEVELOPER.md                    # Documentação técnica detalhada
├── CONTRIBUTING.md                 # Guia de contribuição
│
├── install.sh                      # Script de instalação
└── test-workflow.sh                # Script de teste do workflow

```

## 📝 Descrição dos Arquivos

### Core Application Files

**`cmd/main.go`** (422 linhas)
- Entry point da aplicação
- Parser de comandos CLI
- Formatação de output para usuário
- Inicialização de dependências

**`internal/core/service.go`** (281 linhas)
- Camada de lógica de negócio
- Orquestração entre storage e vault
- Validações e regras de negócio
- Funções: AddGame, CreateCheckpoint, RestoreCheckpoint, etc.

**`internal/models/models.go`** (42 linhas)
- Estruturas de dados: Game e Checkpoint
- Métodos de validação
- Tags JSON para serialização

**`internal/models/errors.go`** (18 linhas)
- Definições de erros customizados
- Permite tratamento específico de erros

**`internal/storage/json_store.go`** (111 linhas)
- Interface MetadataStore
- Implementação JSONStore
- Escrita atômica (temp file + rename)
- Thread-safe com RWMutex

**`internal/vault/manager.go`** (249 linhas)
- Compressão de diretórios em ZIP
- Extração de ZIPs
- Cálculo de hash SHA256
- Prevenção contra zip slip vulnerability

### Documentation Files

**`README.md`**
- Overview do projeto
- Instalação
- Exemplos de uso
- Estrutura de dados

**`QUICKSTART.md`**
- Guia passo-a-passo para iniciantes
- Casos de uso comuns
- Dicas práticas
- Troubleshooting

**`DEVELOPER.md`**
- Arquitetura detalhada
- Decisões de design
- Guia de extensão
- Recursos para desenvolvimento

**`CONTRIBUTING.md`**
- Como contribuir
- Padrões de código
- Processo de PR
- Guia de estilo

### Build & Configuration

**`go.mod`**
- Definição do módulo Go
- Dependência: github.com/google/uuid

**`Makefile`**
- Comandos de build
- Cross-compilation
- Testes
- Instalação

**`.gitignore`**
- Binários
- Diretórios runtime (config/, vault/)
- Arquivos de IDE

**`LICENSE`**
- Licença MIT

### Scripts

**`install.sh`**
- Script de instalação automatizada
- Verifica dependências
- Build e instalação

**`test-workflow.sh`**
- Script de demonstração
- Testa workflow completo
- Útil para validar instalação

## 🎯 Comandos Implementados

### Gerenciamento de Jogos
```bash
savesync add-game --name "Nome" --path "/path/to/saves"
savesync list-games
```

### Gerenciamento de Checkpoints
```bash
savesync checkpoint --game <id> --name "Nome" --note "Nota"
savesync list --game <id>
savesync restore --checkpoint <id>
savesync delete --checkpoint <id>
```

### Utilitários
```bash
savesync version
savesync help
```

## 🔧 Características Técnicas

### Segurança
✅ Escrita atômica de JSON (temp file + rename)
✅ Verificação de integridade SHA256
✅ Prevenção contra zip slip
✅ Thread-safe operations (RWMutex)

### Arquitetura
✅ Clean Architecture
✅ Separation of Concerns
✅ Dependency Injection
✅ Interface-based design

### Qualidade de Código
✅ Error wrapping com contexto
✅ Validações completas
✅ Comentários detalhados
✅ Código modular e testável

## 📊 Estatísticas

- **Total de arquivos Go**: 6
- **Total de linhas de código**: ~1400
- **Dependências externas**: 1 (UUID)
- **Pacotes internos**: 4 (core, models, storage, vault)
- **Comandos CLI**: 7

## 🚀 Como Usar

### Build
```bash
cd savesync
go mod download
go build -o savesync cmd/main.go
```

### Run
```bash
./savesync help
```

### Install
```bash
chmod +x install.sh
./install.sh
```

## 📦 Dependências

### Externas
- `github.com/google/uuid` - Geração de UUIDs

### Standard Library
- `archive/zip` - Compressão/descompressão
- `crypto/sha256` - Hashing
- `encoding/json` - Serialização
- `flag` - Parsing de argumentos CLI
- `os` - File I/O
- `path/filepath` - Manipulação de paths
- `time` - Timestamps

## 🎓 Conceitos Demonstrados

1. **Clean Architecture**: Separação clara de camadas
2. **SOLID Principles**: Especialmente SRP e DIP
3. **Error Handling**: Wrapping e propagação de erros
4. **Concurrency**: Safe concurrent access com mutexes
5. **Data Integrity**: Atomic writes e checksums
6. **Security**: Input validation e path traversal prevention
7. **CLI Design**: User-friendly interface design

## 🔮 Possíveis Extensões Futuras

- [ ] Testes unitários completos
- [ ] Testes de integração
- [ ] Compressão incremental
- [ ] Encriptação de checkpoints
- [ ] Sync com cloud storage
- [ ] Interface gráfica (GUI)
- [ ] Auto-checkpoint em intervalos
- [ ] Comparação de checkpoints (diff)
- [ ] Tags e categorias
- [ ] Busca full-text em notas
- [ ] Export/import de configurações
- [ ] Suporte para múltiplos vaults

## 📄 Licença

MIT License - Veja LICENSE para detalhes

---

**Projeto criado como MVP profissional demonstrando boas práticas em Go**
