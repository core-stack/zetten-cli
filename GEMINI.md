# Projeto: Zetten CLI
Zetten é um CLI em Go para gerenciamento e compartilhamento de arquivos entre múltiplos projetos de forma simples, com suporte a controle de versão via tags/branches.

## 📁 Estrutura do Projeto
```
.
├── config/
│   ├── config.go                # Leitura geral de arquivos de config
│   ├── package.go               # Parser do zetten-package.yaml
│   └── project.go               # Parser do zetten.yaml
├── internal/
│   ├── auth/                    # Módulo de autenticação
│   ├── commands/
│   │   ├── initialize/          # Implementação do comando `init`
│   │   ├── install/             # Implementação do comando `install`
│   │   ├── mirror/              # Implementação do comando `mirror`
│   │   ├── promote/             # Implementação do comando `promote`
│   │   ├── uninstall/           # Implementação do comando `remove`
│   │   └── update/              # Implementação do comando `update`
│   ├── git_util/                # Abstrações para git (clone, switch, etc)
│   ├── prompt/                  # Entrada interativa do usuário
│   └── util/                    # Utilitários (file system, helpers, etc)
├── main.go                      # Entrypoint da aplicação CLI
├── main_test.go                 # Testes de integração
├── GEMINI.md                    # Documentação técnica do projeto
├── go.mod
├── go.sum
└── Makefile                     # Tasks de build/lint/test
```

# 🔧 Comandos Disponíveis
| Comando |	Descrição                                           |
| ------- | --------------------------------------------------  |
| init	  | Adiciona o Zetten no projeto                        |
| install |	Instala dependência de outro repositório/tag/branch |
| remove  |	Remove uma dependência e listeners relacionados     |
| update  |	Atualiza branch/tag de dependência instalada        |
| promote |	Cria nova tag com base em personalização            |
| mirror  |	Espelha arquivos entre projetos em tempo real       |

## ⚙️ Arquivos de Configuração
| zetten.yaml
```yaml
project-name: my-project
version: 1.0.0
dependencies:
  shared-ui: "main"
```

| zetten-package.yaml
```yaml
tag: "v1.2.0"
branch: "main"
repository: "https://github.com/org/shared-ui.git"
```

| zetten-dependencies.yaml
```yaml
shared-ui: "v1.2.0"
utils-lib: "main"
```


## 🔄 Hooks de Eventos
Eventos podem ser registrados em todas as operações:

**install / remove / update / promote**
- before
- after
- error

Com parâmetros: **package, branch or tag, current branch or tag, error**

**mirror**<br/>
**file-change**: detecta alteração e espelha (type: add/edit/remove)<br/>
**add**: adiciona projeto ao espelho<br/>
**remove**: remove projeto do espelho

## ✅ Próximos passos possíveis
- Criar testes para todos os comandos em internal/commands/
- Adicionar serialização YAML em config/* com fallback de validação
- Implementar watchers em mirror.go com debounce para performance
- Gerar documentação interativa via CLI (ex: zetten info)

