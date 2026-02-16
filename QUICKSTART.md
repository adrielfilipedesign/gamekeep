# SaveSync - Quick Start Guide

## 🚀 Instalação Rápida

### Pré-requisitos
- Go 1.21+ instalado
- ~10MB de espaço em disco

### Passo 1: Clone o repositório
```bash
git clone https://github.com/savesync/savesync.git
cd savesync
```

### Passo 2: Build
```bash
go build -o savesync cmd/main.go
```

### Passo 3: (Opcional) Instalar globalmente
```bash
# Linux/macOS
sudo mv savesync /usr/local/bin/

# Ou usar o script de instalação
chmod +x install.sh
./install.sh
```

## 📖 Primeiros Passos

### 1. Adicione seu primeiro jogo

```bash
savesync add-game --name "The Witcher 3" --path "C:/Users/YourName/Documents/The Witcher 3"
```

**Dicas**:
- Use o nome completo do jogo
- O path deve ser o diretório onde os saves estão salvos
- No Windows, use barras normais (`/`) ou duplas (`\\`)

### 2. Verifique os jogos registrados

```bash
savesync list-games
```

Você verá algo como:
```
Registered Games (1):

ID             NAME             SAVE PATH
──             ────             ─────────
the_witcher_3  The Witcher 3   C:/Users/YourName/Documents/The Witcher 3
```

### 3. Crie seu primeiro checkpoint

```bash
savesync checkpoint --game "the_witcher_3" --name "Antes de Kaer Morhen" --note "Level 18, build sinais"
```

**Dicas**:
- Use o `ID` ou `name` do jogo
- Nome do checkpoint deve ser descritivo
- A nota é opcional mas recomendada

### 4. Liste seus checkpoints

```bash
savesync list --game "the_witcher_3"
```

Saída:
```
Checkpoints for the_witcher_3 (1):

ID        NAME                  CREATED           NOTE
──        ────                  ───────           ────
a1b2c3d4  Antes de Kaer Morhen  2024-01-15 14:30  Level 18, build sinais
```

### 5. Restaure um checkpoint

```bash
savesync restore --checkpoint a1b2c3d4
```

⚠️ **ATENÇÃO**: Isso vai **substituir** seus saves atuais!

## 🎮 Casos de Uso Comuns

### Antes de lutas difíceis

```bash
# Criar checkpoint
savesync checkpoint --game darksouls --name "Pre-Boss" --note "Todos items prontos"

# Lute contra o boss...
# Se morrer, restaure:
savesync restore --checkpoint <id>
```

### Experimentar builds diferentes

```bash
# Save atual
savesync checkpoint --game skyrim --name "Build Mago Inicial"

# Mude sua build, teste...

# Se não gostar, volte:
savesync restore --checkpoint <id>

# Se gostar, crie novo checkpoint:
savesync checkpoint --game skyrim --name "Build Mago Melhorado"
```

### Antes de decisões importantes

```bash
# Antes de escolha que afeta história
savesync checkpoint --game witcher3 --name "Antes de escolher Triss/Yennefer"

# Faça a escolha, veja o resultado...
# Se quiser ver a outra opção:
savesync restore --checkpoint <id>
```

### Speedruns e practice

```bash
# Checkpoint no início de cada seção
savesync checkpoint --game celeste --name "Chapter 3 Start"
savesync checkpoint --game celeste --name "Chapter 3 Mid"
savesync checkpoint --game celeste --name "Chapter 3 End"

# Practice section específica
savesync restore --checkpoint <chapter-3-mid>
```

## 🔍 Comandos Úteis

### Ver todos os jogos
```bash
savesync list-games
```

### Ver checkpoints de um jogo
```bash
savesync list --game <nome-ou-id>
```

### Deletar checkpoint antigo
```bash
savesync delete --checkpoint <id>
```

### Ver ajuda
```bash
savesync help
```

## 💡 Dicas Pro

### 1. Use IDs curtos
Você não precisa digitar o UUID completo:
```bash
# Se o ID é: 550e8400-e29b-41d4-a716-446655440000
# Pode usar apenas:
savesync restore --checkpoint 550e8400
```

### 2. Busca flexível de jogos
Pode usar nome parcial:
```bash
savesync checkpoint --game witcher --name "..."
# Funciona se só houver um jogo com "witcher" no nome
```

### 3. Organize com notas
Use notas para lembrar contexto:
```bash
--note "Level 50, equipamento lendário completo, 1000 horas jogadas"
```

### 4. Checkpoint antes de mods
```bash
savesync checkpoint --game skyrim --name "Vanilla" --note "Antes de instalar mods"
```

## 📁 Onde os arquivos ficam?

```
~/.savesync/
├── config/
│   ├── games.json          # Lista de jogos
│   └── checkpoints.json    # Lista de checkpoints
└── vault/
    └── {game_id}/
        └── {checkpoint}.zip # Arquivos compactados
```

**Dicas**:
- Faça backup da pasta `~/.savesync` periodicamente
- Cada checkpoint é um arquivo ZIP independente
- Checkpoints podem ser grandes dependendo do tamanho dos saves

## ⚠️ Avisos Importantes

### ⚠️ SEMPRE teste primeiro!
```bash
# 1. Crie checkpoint de teste
savesync checkpoint --game test --name "Teste"

# 2. Modifique algo no save

# 3. Restaure e verifique
savesync restore --checkpoint <id>
```

### ⚠️ Saves em cloud (Steam Cloud, etc)
Se o jogo usa cloud saves:
1. Desabilite sync antes de restaurar
2. Ou o cloud pode sobrescrever seu restore

### ⚠️ Espaço em disco
Checkpoints ocupam espaço! Limpe antigos:
```bash
savesync delete --checkpoint <id-antigo>
```

## 🆘 Problemas Comuns

### "Game not found"
- Verifique com: `savesync list-games`
- Use o ID exato ou nome completo

### "Save path does not exist"
- Verifique se o path está correto
- Alguns jogos mudam o local dos saves

### "Permission denied"
- Execute com permissões adequadas
- No Windows, pode precisar executar como administrador

### "Hash mismatch"
- Arquivo do vault foi corrompido
- Delete o checkpoint e crie um novo

## 🎓 Exemplos Práticos

### Elden Ring - Practice boss

```bash
# Registrar jogo
savesync add-game --name "Elden Ring" --path "C:/Users/You/AppData/Roaming/EldenRing"

# Checkpoint antes de cada tentativa
savesync checkpoint --game elden --name "Malenia Start"

# Practice... practice... practice...

# Restaurar para tentar de novo
savesync restore --checkpoint <id>
```

### Múltiplos finais
```bash
# Checkpoint antes da decisão final
savesync checkpoint --game game --name "Pre-Final Choice"

# Ver final A
# Restaurar
savesync restore --checkpoint <id>

# Ver final B
# Restaurar novamente
savesync restore --checkpoint <id>

# Ver final C
```

## 📚 Próximos Passos

- Leia o [README.md](README.md) completo
- Veja [DEVELOPER.md](DEVELOPER.md) para contribuir
- Reporte bugs no GitHub Issues

---

**Happy Gaming! 🎮**
