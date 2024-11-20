# Development Guide: Building with Podman

[English](#english) | [Français](#français)

---
# English

## Why Podman?

Podman was chosen for this project for several key reasons:
1. **Daemonless Architecture**: Unlike Docker, Podman doesn't require a daemon process
2. **Enhanced Security**: Rootless containers by default
3. **OCI Compliance**: Full compatibility with container standards
4. **Cross-Platform**: Works on Linux, macOS, and Windows
5. **No Dependencies**: No need to install Go locally

## Development Environment Setup

### Required Tools
- Podman (latest version)
- Git for version control
- Text editor of your choice

### Project Structure
```
.
├── Dockerfile
├── Makefile
├── main.go
├── frenchdiceware.txt
├── diceware-fr-alt.txt
├── go.mod
└── README.md
```

### Building Process

The Dockerfile is configured in two stages:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
COPY go.mod .
COPY frenchdiceware.txt .
COPY diceware-fr-alt.txt .
RUN go mod tidy && \
    go build -ldflags="-s -w" -o diceware

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/diceware .
COPY frenchdiceware.txt .
COPY diceware-fr-alt.txt .
ENTRYPOINT ["./diceware"]
```

### Makefile Commands
```makefile
.PHONY: build run clean

build:
    podman build -t diceware .

run:
    podman run -it diceware

clean:
    podman rmi diceware
```

## Development Workflow

1. **Initial Setup**
   ```bash
   git clone [repository]
   cd diceware-generator
   make build
   ```

2. **Testing Changes**
   ```bash
   make clean
   make build
   make run
   ```

## Troubleshooting Common Issues

1. **Build Errors**
   - Verify go.mod is present and correct
   - Ensure all required files are in place
   - Check file permissions

2. **Runtime Errors**
   - Verify wordlist files are properly formatted
   - Check terminal compatibility for colors
   - Ensure interactive input is working

---
# Français

## Pourquoi Podman ?

Podman a été choisi pour ce projet pour plusieurs raisons clés :
1. **Architecture sans démon** : Contrairement à Docker, Podman ne nécessite pas de processus démon
2. **Sécurité renforcée** : Conteneurs rootless par défaut
3. **Conformité OCI** : Compatibilité totale avec les standards des conteneurs
4. **Multi-plateforme** : Fonctionne sur Linux, macOS et Windows
5. **Sans dépendances** : Pas besoin d'installer Go localement

## Configuration de l'environnement de développement

### Outils nécessaires
- Podman (dernière version)
- Git pour le contrôle de version
- Éditeur de texte au choix

### Structure du projet
```
.
├── Dockerfile
├── Makefile
├── main.go
├── frenchdiceware.txt
├── diceware-fr-alt.txt
├── go.mod
└── README.md
```

### Processus de build

Le Dockerfile est configuré en deux étapes :
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
COPY go.mod .
COPY frenchdiceware.txt .
COPY diceware-fr-alt.txt .
RUN go mod tidy && \
    go build -ldflags="-s -w" -o diceware

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/diceware .
COPY frenchdiceware.txt .
COPY diceware-fr-alt.txt .
ENTRYPOINT ["./diceware"]
```

### Commandes Makefile
```makefile
.PHONY: build run clean

build:
    podman build -t diceware .

run:
    podman run -it diceware

clean:
    podman rmi diceware
```

## Workflow de développement

1. **Configuration initiale**
   ```bash
   git clone [dépôt]
   cd diceware-generator
   make build
   ```

2. **Test des modifications**
   ```bash
   make clean
   make build
   make run
   ```

## Résolution des problèmes courants

1. **Erreurs de build**
   - Vérifier que go.mod est présent et correct
   - S'assurer que tous les fichiers requis sont en place
   - Vérifier les permissions des fichiers

2. **Erreurs d'exécution**
   - Vérifier le format des fichiers wordlist
   - Vérifier la compatibilité du terminal pour les couleurs
   - S'assurer que l'entrée interactive fonctionne

## Notes importantes

- Le build multi-stage permet d'obtenir une image finale légère
- L'utilisation de Alpine Linux minimise la taille de l'image
- Le flag `-it` est nécessaire pour l'interactivité
- Les wordlists sont intégrées à l'image finale

## Sécurité du processus de développement

- Pas de secrets dans l'image
- Pas de dépendances externes sauf Go
- Build reproductible
- Validation des entrées utilisateur
