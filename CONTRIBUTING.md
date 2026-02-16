# Contributing to SaveSync

Obrigado pelo interesse em contribuir com o SaveSync! 🎮

## 🤝 Como Contribuir

### Reportar Bugs

Se você encontrou um bug, por favor abra uma issue com:

- **Descrição clara** do problema
- **Passos para reproduzir**
- **Comportamento esperado** vs **comportamento atual**
- **Ambiente**: OS, versão do Go, versão do SaveSync
- **Logs/Erros** se aplicável

### Sugerir Features

Para sugestões de novas funcionalidades:

1. Verifique se já não existe uma issue similar
2. Descreva o caso de uso
3. Explique por que seria útil
4. Se possível, sugira uma implementação

### Pull Requests

#### Antes de começar

1. Fork o repositório
2. Clone seu fork: `git clone https://github.com/seu-usuario/savesync.git`
3. Crie uma branch: `git checkout -b feature/minha-feature`

#### Durante o desenvolvimento

**Padrões de Código**:
```bash
# Formatar código
go fmt ./...

# Verificar com go vet
go vet ./...

# Rodar testes
go test ./...
```

**Commits**:
- Use mensagens descritivas em português ou inglês
- Prefixos sugeridos:
  - `feat:` - Nova funcionalidade
  - `fix:` - Correção de bug
  - `docs:` - Documentação
  - `refactor:` - Refatoração
  - `test:` - Testes
  - `chore:` - Tarefas de manutenção

Exemplos:
```
feat: adicionar suporte para compressão incremental
fix: corrigir race condition em atomic writes
docs: melhorar documentação de instalação
```

#### Checklist antes do PR

- [ ] Código formatado com `go fmt`
- [ ] Sem warnings do `go vet`
- [ ] Testes passando
- [ ] Documentação atualizada (README.md, DEVELOPER.md)
- [ ] Comentários em código complexo
- [ ] Commit messages descritivos

#### Submetendo o PR

1. Push para seu fork: `git push origin feature/minha-feature`
2. Abra um Pull Request no GitHub
3. Descreva as mudanças claramente
4. Referencie issues relacionadas (ex: "Closes #123")

## 📋 Guia de Estilo

### Nomenclatura

**Variáveis**:
```go
// ✅ Bom
userID := "123"
checkpointName := "Boss Fight"

// ❌ Evitar
uid := "123"
n := "Boss Fight"
```

**Funções**:
```go
// ✅ Bom - Verbo + Substantivo
CreateCheckpoint()
ValidateGame()
LoadCheckpoints()

// ❌ Evitar
DoCheckpoint()
Check()
```

**Interfaces**:
```go
// ✅ Bom - Substantivo + "er"
type MetadataStore interface {}
type CheckpointManager interface {}

// ❌ Evitar
type IMetadataStore interface {}
```

### Estrutura de Funções

```go
func FunctionName(param1 Type1, param2 Type2) (ReturnType, error) {
    // 1. Validação de input
    if param1 == "" {
        return nil, ErrInvalidInput
    }
    
    // 2. Lógica principal
    result := processData(param1, param2)
    
    // 3. Tratamento de erros
    if err != nil {
        return nil, fmt.Errorf("failed to process: %w", err)
    }
    
    // 4. Return
    return result, nil
}
```

### Comentários

**Funções públicas** (exportadas):
```go
// CreateCheckpoint creates a new checkpoint for the specified game.
// It compresses the save directory, calculates the hash, and stores metadata.
//
// Parameters:
//   - gameID: The unique identifier of the game
//   - name: Display name for the checkpoint
//   - note: Optional note (can be empty)
//
// Returns the created checkpoint and any error encountered.
func CreateCheckpoint(gameID, name, note string) (*Checkpoint, error) {
    // ...
}
```

**Código complexo**:
```go
// Atomic write pattern: write to temp file then rename
// This prevents corruption if process is killed mid-write
tmpFile := filepath + ".tmp"
if err := os.WriteFile(tmpFile, data, 0644); err != nil {
    return err
}
return os.Rename(tmpFile, filepath)
```

### Error Handling

```go
// ✅ Bom - Wrap errors com contexto
if err := saveFile(path); err != nil {
    return fmt.Errorf("failed to save config file: %w", err)
}

// ❌ Evitar - Perder contexto
if err := saveFile(path); err != nil {
    return err
}

// ✅ Bom - Errors customizados
var ErrGameNotFound = errors.New("game not found")

// ❌ Evitar - Strings genéricas
return errors.New("not found")
```

## 🧪 Testes

### Estrutura de Teste

```go
func TestFunctionName(t *testing.T) {
    // Arrange
    service := setupTestService(t)
    input := "test-input"
    
    // Act
    result, err := service.Function(input)
    
    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

### Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid input", "test", false},
        {"empty input", "", true},
        {"invalid chars", "test@#$", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 🏗️ Arquitetura

### Princípios

1. **Separation of Concerns**: Cada pacote tem uma responsabilidade clara
2. **Dependency Injection**: Dependências passadas como parâmetros
3. **Interface-based Design**: Programar para interfaces, não implementações
4. **Error Handling**: Sempre propagar erros com contexto

### Adicionando Nova Feature

**Exemplo**: Adicionar suporte para tags em checkpoints

1. **Model** (`internal/models/models.go`):
```go
type Checkpoint struct {
    // ... campos existentes
    Tags []string `json:"tags,omitempty"`
}
```

2. **Service** (`internal/core/service.go`):
```go
func (s *Service) CreateCheckpoint(..., tags []string) (*Checkpoint, error) {
    checkpoint := &Checkpoint{
        // ... campos existentes
        Tags: tags,
    }
    // ...
}
```

3. **CLI** (`cmd/main.go`):
```go
func (c *CLI) createCheckpoint(args []string) error {
    fs := flag.NewFlagSet("checkpoint", flag.ExitOnError)
    // ... flags existentes
    tags := fs.String("tags", "", "Comma-separated tags")
    // ...
}
```

4. **Tests**:
```go
func TestCreateCheckpointWithTags(t *testing.T) {
    // ...
}
```

5. **Documentation**:
- Atualizar README.md
- Atualizar DEVELOPER.md

## 📝 Documentação

Ao adicionar features, sempre atualizar:

- [ ] README.md - Exemplos de uso
- [ ] DEVELOPER.md - Detalhes técnicos
- [ ] Comentários no código
- [ ] Exemplos em `test-workflow.sh` se aplicável

## 🔍 Code Review

Revisores vão verificar:

✅ **Qualidade**:
- Código limpo e legível
- Funções pequenas e focadas
- Nomes descritivos

✅ **Segurança**:
- Validação de input
- Tratamento de erros
- Prevenção de vulnerabilidades

✅ **Performance**:
- Sem loops desnecessários
- Uso eficiente de memória
- Operações I/O otimizadas

✅ **Testes**:
- Cobertura adequada
- Casos edge testados
- Testes passando

## 💬 Comunicação

- **GitHub Issues**: Para discussões técnicas e bugs
- **Pull Requests**: Para revisão de código
- **Discussions**: Para perguntas gerais

## 🌟 Reconhecimento

Contribuidores serão listados no README.md!

---

Obrigado por contribuir com o SaveSync! 🚀
