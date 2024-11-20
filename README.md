# Diceware Passphrase Generator
[English](#english) | [Français](#français) | [한국어](#한국어)

---
# English

## About
A secure passphrase generator using the Diceware method with dual wordlists for enhanced entropy, developed in collaboration with Claude (Anthropic AI).

## Key Features
- **Dual Wordlist System**: Enhanced entropy using two distinct word lists
- **Two Operating Modes**:
  - Automatic: Quick passphrase generation
  - Interactive: Choose from 6 words at each step
- **Customization Options**:
  - Word count: 2-8 words
  - Capitalization: First/Last, Random, or None
  - Special characters injection
  - QR Code display
- **Security Features**:
  - Cryptographic random generation
  - Enhanced entropy (13.92 bits/word)
  - Secure character injection
  - Input validation

## Installation
```bash
git clone [repository_url]
cd diceware-generator
go build
```

## Usage
```bash
./diceware
```

Follow the prompts to:
1. Select mode (Automatic/Interactive)
2. Specify word count (2-8)
3. Choose capitalization method
4. Enable/disable special characters
5. Enable/disable QR code display

## Security Metrics
- Base entropy: 13.92 bits/word
- 4 words: 55.68 bits
- 6 words: 83.52 bits
- 8 words: 111.36 bits
- Special characters: +4.17 bits each

## Technical Details
- Written in Go
- Uses crypto/rand for secure random generation
- MIT License
- Input validation and error handling
- Terminal-friendly colored output

---
# Français

## À propos
Générateur de phrases de passe sécurisé utilisant la méthode Diceware avec double liste de mots pour une entropie améliorée, développé en collaboration avec Claude (IA Anthropic).

## Fonctionnalités principales
- **Système à double liste** : Entropie améliorée avec deux listes distinctes
- **Deux modes de fonctionnement** :
  - Automatique : Génération rapide
  - Interactif : Choix parmi 6 mots à chaque étape
- **Options de personnalisation** :
  - Nombre de mots : 2-8 mots
  - Capitalisation : Premier/Dernier, Aléatoire, Aucune
  - Injection de caractères spéciaux
  - Affichage QR Code
- **Caractéristiques de sécurité** :
  - Génération aléatoire cryptographique
  - Entropie améliorée (13,92 bits/mot)
  - Injection sécurisée de caractères
  - Validation des entrées

## Installation
```bash
git clone [url_du_dépôt]
cd diceware-generator
go build
```

## Utilisation
```bash
./diceware
```

Suivez les instructions pour :
1. Choisir le mode (Automatique/Interactif)
2. Spécifier le nombre de mots (2-8)
3. Choisir la méthode de capitalisation
4. Activer/désactiver les caractères spéciaux
5. Activer/désactiver l'affichage du QR Code

## Métriques de sécurité
- Entropie de base : 13,92 bits/mot
- 4 mots : 55,68 bits
- 6 mots : 83,52 bits
- 8 mots : 111,36 bits
- Caractères spéciaux : +4,17 bits chacun

## Détails techniques
- Écrit en Go
- Utilise crypto/rand pour la génération aléatoire sécurisée
- Licence MIT
- Validation des entrées et gestion des erreurs
- Sortie colorée adaptée au terminal

---
# 한국어

## 소개
Diceware 방식을 사용하는 보안 암호구 생성기로, 이중 단어 목록을 통해 향상된 엔트로피를 제공합니다. Claude(Anthropic AI)와의 협업으로 개발되었습니다.

## 주요 기능
- **이중 단어 목록 시스템**: 두 개의 개별 단어 목록을 통한 향상된 엔트로피
- **두 가지 작동 모드**:
  - 자동: 빠른 암호구 생성
  - 대화형: 각 단계마다 6개 단어 중 선택
- **사용자 정의 옵션**:
  - 단어 수: 2-8 단어
  - 대문자 설정: 첫/마지막, 무작위, 없음
  - 특수 문자 삽입
  - QR 코드 표시
- **보안 기능**:
  - 암호화 난수 생성
  - 향상된 엔트로피 (단어당 13.92 비트)
  - 안전한 문자 삽입
  - 입력 유효성 검사

## 설치
```bash
git clone [저장소_URL]
cd diceware-generator
go build
```

## 사용법
```bash
./diceware
```

다음 단계를 따르세요:
1. 모드 선택 (자동/대화형)
2. 단어 수 지정 (2-8)
3. 대문자 설정 방법 선택
4. 특수 문자 활성화/비활성화
5. QR 코드 표시 활성화/비활성화

## 보안 메트릭
- 기본 엔트로피: 단어당 13.92 비트
- 4단어: 55.68 비트
- 6단어: 83.52 비트
- 8단어: 111.36 비트
- 특수 문자: 각각 +4.17 비트

## 기술 세부사항
- Go 언어로 작성
- 보안 난수 생성을 위한 crypto/rand 사용
- MIT 라이선스
- 입력 유효성 검사 및 오류 처리
- 터미널 친화적인 컬러 출력

---
## License / Licence / 라이선스

MIT License - See LICENSE file for details.  
Licence MIT - Voir le fichier LICENSE pour plus de détails.  
MIT 라이선스 - 자세한 내용은 LICENSE 파일을 참조하세요.

