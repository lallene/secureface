# Architecture Technique — SecureFace Access API

## 1. Présentation

SecureFace Access est une plateforme de contrôle d’accès biométrique basée sur la reconnaissance faciale.

Le backend Go assure la logique métier, la sécurité, les permissions, la journalisation et la communication avec le service IA.

## 2. Architecture globale

```text
Utilisateur / Caméra / Frontend
        |
        v
Backend Go API (Gin)
        |
        |---- PostgreSQL
        |
        |---- Service IA Python FastAPI / DeepFace