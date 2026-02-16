# SaveSync - Developer Documentation

## 📋 Visão Geral da Arquitetura

SaveSync segue princípios de Clean Architecture com separação clara de responsabilidades.

### Camadas

```
┌─────────────────────────────────────┐
│           CLI Layer                 │
│         (cmd/main.go)               │
│  - Parsing de comandos              │
│  - Formatação de output             │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│       Business Logic Layer          │
│      (internal/core/service.go)     │
│  - Validações                       │
│  - Orquestração                     │
│  - Regras de negócio                │
└──────────────┬──────────────────────┘
               │
        ┌──────┴──────┐
        │             │
┌───────▼──────┐ ┌───▼──────────────┐
│   Storage    │ │   Vault Manager  │
│ (metadata)   │ │  (file ops)      │
│              │ │                  │
│ - JSON I/O   │ │ - ZIP/Unzip      │
│ - Atomic     │ │ - SHA256         │
│   writes     │ │ - File copy      │
└──────────────┘ └──────────────────┘
```

## 🏗️ Estrutura de Pacotes

### `/cmd`
**Responsabilidade**: Entry point e interface CLI

**Arquivos**:
- `main.go` - Inicialização e comando dispatcher

**Funções principais**:
- Parse de argumentos
- Formatação de output para usuário
- Coordenação entre comandos e service layer

### `/internal/models`
**Responsabilidade**: Definição de tipos e erros

**Arquivos**:
- `models.go` - Estruturas Game e Checkpoint
- `errors.go` - Erros customizados

**Princípios**:
- Estruturas imutáveis (sem setters)
- Validação em métodos próprios
- JSON tags para serialização

### `/internal/storage`
**Responsabilidade**: Persistência de metadados

**Arquivos**:
- `json_store.go` - Implementação com JSON files

**Interface**:
```go
type MetadataStore interface {
    SaveGames(games []Game) error
    LoadGames() ([]Game, error)
    SaveCheckpoints(checkpoints []Checkpoint) error
    LoadCheckpoints() ([]Checkpoint, error)
}
```

**Características**:
- **Atomic writes**: Temp file + rename para evitar corrupção
- **Thread-safe**: RWMutex para concurrent access
- **Lazy loading**: Arquivos só criados quando necessário

### `/internal/vault`
**Responsabilidade**: Operações com arquivos de checkpoint

**Arquivos**:
- `manager.go` - Compressão, descompressão e hashing

**Operações principais**:
- `CreateCheckpoint()` - Comprime diretório → ZIP
- `RestoreCheckpoint()` - Extrai ZIP → diretório
- `VerifyCheckpoint()` - Valida hash SHA256
- `DeleteCheckpoint()` - Remove arquivo do vault

**Segurança**:
- Prevenção contra zip slip
- Verificação de integridade via hash
- Cleanup em caso de erro

### `/internal/core`
**Responsabilidade**: Lógica de negócio

**Arquivos**:
- `service.go` - Orquestração de operações

**Funções principais**:
```go
AddGame(name, path string) (*Game, error)
GetGame(identifier string) (*Game, error)
CreateCheckpoint(gameID, name, note string) (*Checkpoint, error)
RestoreCheckpoint(checkpointID string) error
ListCheckpoints(gameID string) ([]Checkpoint, error)
DeleteCheckpoint(checkpointID string) error
```

## 🔄 Fluxos de Dados

### Criação de Checkpoint

```
1. User → CLI: savesync checkpoint --game witcher3 --name "Boss"
2. CLI → Service: CreateCheckpoint("witcher3", "Boss", "")
3. Service → Storage: LoadGames()
4. Service: Valida que game existe
5. Service: Gera UUID para checkpoint
6. Service → VaultMgr: CreateCheckpoint(gameID, uuid, savePath)
7. VaultMgr: Cria ZIP do diretório
8. VaultMgr: Calcula SHA256
9. VaultMgr → Service: (vaultFile, hash)
10. Service: Cria objeto Checkpoint
11. Service → Storage: SaveCheckpoints(checkpoints)
12. Storage: Atomic write do JSON
13. Service → CLI: checkpoint object
14. CLI → User: Mensagem de sucesso
```

### Restauração de Checkpoint

```
1. User → CLI: savesync restore --checkpoint abc123
2. CLI → Service: RestoreCheckpoint("abc123")
3. Service → Storage: LoadCheckpoints()
4. Service: Encontra checkpoint por ID
5. Service → Storage: LoadGames()
6. Service: Encontra game associado
7. Service → VaultMgr: VerifyCheckpoint(vaultFile, hash)
8. VaultMgr: Calcula hash atual e compara
9. Service → VaultMgr: RestoreCheckpoint(vaultFile, savePath)
10. VaultMgr: Remove diretório existente
11. VaultMgr: Extrai ZIP para savePath
12. Service → CLI: success
13. CLI → User: Mensagem de sucesso
```

## 🔐 Garantias de Segurança

### Escrita Atômica
```go
// Padrão usado:
1. Marshal data → JSON
2. Write JSON → tempfile
3. Rename tempfile → target (atomic)
```

### Verificação de Integridade
```go
// Sempre antes de restore:
actualHash := CalculateHash(zipFile)
if actualHash != checkpoint.Hash {
    return ErrHashMismatch
}
```

### Prevenção Zip Slip
```go
// Verificação em cada arquivo:
if !strings.HasPrefix(filepath.Clean(path), baseDir) {
    return error // Path traversal detected
}
```

## 🧪 Testabilidade

### Interfaces
Todas as dependências principais são interfaces:
- `MetadataStore` - Permite mock de storage
- Pode adicionar `VaultManager` interface futuramente

### Dependency Injection
```go
// Service não cria dependências
service := NewService(store, vaultMgr)
```

### Testes Sugeridos

```go
// internal/core/service_test.go
func TestAddGame(t *testing.T) {
    mockStore := &MockStore{}
    mockVault := &MockVault{}
    service := NewService(mockStore, mockVault)
    
    game, err := service.AddGame("Test", "/path")
    // assertions...
}
```

## 📊 Formato de Dados

### Game JSON
```json
{
  "id": "the_witcher_3",          // Sanitized from name
  "name": "The Witcher 3",         // Display name
  "save_path": "/path/to/saves"    // Absolute path
}
```

### Checkpoint JSON
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",  // UUID v4
  "game_id": "the_witcher_3",                     // Reference
  "name": "Before Boss",                          // Display name
  "note": "Level 25, fire build",                 // Optional
  "vault_file": "the_witcher_3/{uuid}.zip",       // Relative path
  "hash": "a1b2c3d4e5f6...",                      // SHA256
  "created_at": "2024-01-15T10:30:00Z"           // RFC3339 UTC
}
```

## 🎯 Decisões de Design

### Por que JSON e não SQLite?
- **Simplicidade**: Sem dependências externas
- **Portabilidade**: Arquivos legíveis e editáveis
- **Backup**: Fácil fazer backup manual
- **Volume**: Quantidade de dados é pequena

### Por que UUID em vez de auto-increment?
- **Distributed-safe**: Funciona sem coordenação
- **Collision-free**: Probabilidade desprezível
- **URL-friendly**: Pode ser usado em futuras APIs

### Por que ZIP em vez de tar.gz?
- **Cross-platform**: Windows tem suporte nativo
- **Random access**: Pode extrair arquivos individuais
- **Standard library**: Go tem excelente suporte

### Por que SHA256?
- **Segurança**: Resistente a colisões
- **Performance**: Rápido o suficiente
- **Detecção**: Identifica corrupção de dados

## 🔧 Extensões Futuras

### Possíveis Melhorias

1. **Compression levels**
```go
func (m *Manager) CreateCheckpoint(..., compressionLevel int)
```

2. **Incremental backups**
```go
func (m *Manager) CreateIncrementalCheckpoint(baseCheckpoint string)
```

3. **Encryption**
```go
func (m *Manager) CreateEncryptedCheckpoint(..., password string)
```

4. **Cloud sync**
```go
type CloudBackend interface {
    Upload(checkpoint) error
    Download(checkpointID) error
}
```

5. **Checkpoint tags**
```go
type Checkpoint struct {
    // ...
    Tags []string `json:"tags"`
}
```

6. **Auto-checkpoint**
```go
func (s *Service) StartAutoCheckpoint(gameID string, interval time.Duration)
```

## 🐛 Debugging

### Verificar Metadados
```bash
# Ver jogos registrados
cat ~/.savesync/config/games.json | jq

# Ver checkpoints
cat ~/.savesync/config/checkpoints.json | jq
```

### Verificar Vault
```bash
# Listar arquivos no vault
find ~/.savesync/vault -type f -name "*.zip"

# Verificar conteúdo de um checkpoint
unzip -l ~/.savesync/vault/{game_id}/{checkpoint_id}.zip
```

### Verificar Hash
```bash
# Calcular hash de um checkpoint
sha256sum ~/.savesync/vault/{game_id}/{checkpoint_id}.zip
```

## 📚 Recursos Adicionais

### Go Best Practices
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Arquitetura
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

### Segurança
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Zip Slip Vulnerability](https://snyk.io/research/zip-slip-vulnerability)
