# SecureFace

SecureFace est une plateforme de contrôle d'accès intelligente basée sur la reconnaissance faciale. L'application permet de gérer les utilisateurs, les portes sécurisées, les permissions d'accès et les journaux d'audit tout en intégrant un moteur de reconnaissance faciale.

## Fonctionnalités

### Authentification et sécurité

* Authentification JWT
* Gestion des rôles utilisateurs
* Contrôle des accès
* Journalisation des actions
* Chiffrement des données sensibles

### Gestion des portes

* Création et gestion des portes
* Verrouillage et déverrouillage à distance
* Historique des accès

### Reconnaissance faciale

* Enregistrement des visages
* Vérification d'identité
* Vérification à partir d'images
* Intégration avec un moteur IA de reconnaissance faciale

### Audit et conformité

* Historique des accès
* Traçabilité des opérations
* Nettoyage automatique des données expirées

## Architecture

```text
Frontend
    │
    ▼
SecureFace API (Go + Gin)
    │
    ├── PostgreSQL
    │
    └── Face Recognition AI Service
```

## Stack Technique

### Backend

* Go
* Gin Framework
* GORM
* JWT

### Base de données

* PostgreSQL

### DevOps

* Docker
* Docker Compose
* GitHub Actions
* GitHub Container Registry (GHCR)
* Trivy Security Scanner

## Installation Locale

### Prérequis

* Docker
* Docker Compose

### Cloner le projet

```bash
git clone https://github.com/lallene/secureface.git
cd secureface
```

### Variables d'environnement

Créer un fichier `.env` :

```env
APP_PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=smartface_db

JWT_SECRET=YOUR_SECRET

FACE_ENCRYPTION_KEY=YOUR_ENCRYPTION_KEY

FACE_AI_URL=http://host.docker.internal:8000/recognize
```

### Lancer l'application

```bash
docker compose up --build
```

API disponible sur :

```text
http://localhost:8080
```

## CI/CD

Le projet utilise GitHub Actions pour :

* Compilation automatique
* Tests automatiques
* Construction des images Docker
* Analyse de sécurité avec Trivy
* Publication automatique des images Docker dans GitHub Container Registry

## Docker Image

Image publiée automatiquement :

```text
ghcr.io/lallene/secureface:latest
```

## Endpoints principaux

### Authentification

```http
POST /api/auth/register
POST /api/auth/login
GET  /api/users/profile
```

### Portes

```http
POST /api/doors/
GET  /api/doors/
POST /api/doors/:id/unlock
POST /api/doors/:id/lock
```

### Reconnaissance Faciale

```http
POST /api/faces/register
GET  /api/faces/
POST /api/faces/verify
POST /api/faces/verify-image
```

### Audit

```http
POST /api/access/open
GET  /api/access/logs
```

## Sécurité

* JWT Authentication
* Chiffrement des données biométriques
* Analyse des vulnérabilités avec Trivy
* Pipeline DevSecOps
* Gestion des permissions d'accès

## Auteur

Lallene Cedric

Master Expert Systèmes d'Information & Sécurité

GitHub : https://github.com/lallene
