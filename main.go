/*
MIT License

Copyright (c) 2024 [Votre nom]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
    "bufio"
    "crypto/rand"
    "embed"
    "fmt"
    "log"
    "math/big"
    "os"
    "strings"
    "sync"
    "unicode"

    "github.com/skip2/go-qrcode"
)

// Number of choices presented in interactive mode
const CHOICES_COUNT = 6

//go:embed *.txt
var wordlistFS embed.FS

// Terminal color codes
const (
    colorReset  = "\033[0m"
    colorRed    = "\033[31m"
    colorGreen  = "\033[32m"
    colorYellow = "\033[33m"
    colorBlue   = "\033[34m"
    colorPurple = "\033[35m"
    colorCyan   = "\033[36m"
    colorWhite  = "\033[97m"
)

// Global QR code display flag
var showQR bool

// DicewareGenerator handles the passphrase generation logic
type DicewareGenerator struct {
    primaryWordlist    map[string]string
    secondaryWordlist  map[string]string
    specialChars       string
    numberChars        string
    entropyPool        []byte
    entropyMutex       sync.Mutex
    sessionSalt        []byte
}

// PassphraseResult contains the generated passphrase and its metadata
type PassphraseResult struct {
    Passphrase   string
    SelectedList []string
    WordIndices  []string
    Entropy      float64
    QRCode       string
}

// NewDicewareGenerator initializes a new generator with secure random seed
func NewDicewareGenerator() (*DicewareGenerator, error) {
    sessionSalt := make([]byte, 32)
    if _, err := rand.Read(sessionSalt); err != nil {
        return nil, fmt.Errorf("erreur génération salt: %w", err)
    }

    dg := &DicewareGenerator{
        primaryWordlist:   make(map[string]string),
        secondaryWordlist: make(map[string]string),
        specialChars:      "!@#$%^&*",
        numberChars:       "0123456789",
        entropyPool:       make([]byte, 1024),
        sessionSalt:       sessionSalt,
    }

    if err := dg.loadWordlists(); err != nil {
        return nil, err
    }

    return dg, nil
}

// loadWordlists reads wordlists from embedded files
func (dg *DicewareGenerator) loadWordlists() error {
    lists := []struct {
        filename string
        wordmap  *map[string]string
    }{
        {"frenchdiceware.txt", &dg.primaryWordlist},
        {"diceware-fr-alt.txt", &dg.secondaryWordlist},
    }

    for _, list := range lists {
        file, err := wordlistFS.Open(list.filename)
        if err != nil {
            return fmt.Errorf("erreur ouverture %s: %w", list.filename, err)
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            line := strings.TrimSpace(scanner.Text())
            if len(line) > 5 {
                number := line[:5]
                word := strings.TrimSpace(line[5:])
                (*list.wordmap)[number] = word
            }
        }
    }
    return nil
}

// getSecureRandom generates a cryptographically secure random number
func (dg *DicewareGenerator) getSecureRandom(max int64) (int64, error) {
    n, err := rand.Int(rand.Reader, big.NewInt(max))
    if err != nil {
        return 0, err
    }
    return n.Int64(), nil
}

// rollDice simulates rolling 5 six-sided dice
func (dg *DicewareGenerator) rollDice() (string, error) {
    var result strings.Builder
    for i := 0; i < 5; i++ {
        n, err := dg.getSecureRandom(6)
        if err != nil {
            return "", err
        }
        result.WriteString(fmt.Sprintf("%d", n+1))
    }
    return result.String(), nil
}

// generateQRCode creates a terminal-displayable QR code
func (dg *DicewareGenerator) generateQRCode(text string) (string, error) {
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return "", err
    }

    qr.DisableBorder = false
    
    var result strings.Builder
    result.WriteString("\n")
    
    bitmap := qr.Bitmap()
    for _, row := range bitmap {
        for _, cell := range row {
            if cell {
                result.WriteString("  ")
            } else {
                result.WriteString("██")
            }
        }
        result.WriteString("\n")
    }

    return result.String(), nil
}

// generateWordChoices generates CHOICES_COUNT random words for selection
func (dg *DicewareGenerator) generateWordChoices() ([]string, []bool, []string, error) {
    words := make([]string, CHOICES_COUNT)
    isSecondary := make([]bool, CHOICES_COUNT)
    indices := make([]string, CHOICES_COUNT)

    for i := 0; i < CHOICES_COUNT; i++ {
        useSecondary, err := dg.getSecureRandom(2)
        if err != nil {
            return nil, nil, nil, err
        }

        index, err := dg.rollDice()
        if err != nil {
            return nil, nil, nil, err
        }

        if useSecondary == 1 {
            words[i] = dg.secondaryWordlist[index]
            isSecondary[i] = true
        } else {
            words[i] = dg.primaryWordlist[index]
            isSecondary[i] = false
        }
        indices[i] = index
    }

    return words, isSecondary, indices, nil
}

// transformWord applies capitalization and special character rules
func (dg *DicewareGenerator) transformWord(word string, capitalize bool, addSpecial bool) (string, error) {
    if capitalize {
        word = strings.Title(word)
    }

    if addSpecial {
        pos, err := dg.getSecureRandom(int64(len(word)))
        if err != nil {
            return "", err
        }

        useSpecial, err := dg.getSecureRandom(2)
        if err != nil {
            return "", err
        }

        var char string
        if useSpecial == 0 {
            index, err := dg.getSecureRandom(int64(len(dg.specialChars)))
            if err != nil {
                return "", err
            }
            char = string(dg.specialChars[index])
        } else {
            index, err := dg.getSecureRandom(int64(len(dg.numberChars)))
            if err != nil {
                return "", err
            }
            char = string(dg.numberChars[index])
        }

        word = word[:pos] + char + word[pos+1:]
    }

    return word, nil
}

// InteractiveGeneration handles interactive passphrase generation
func (dg *DicewareGenerator) InteractiveGeneration(wordCount int, randomCapitalize, firstLastCapitalize, addSpecial bool) (*PassphraseResult, error) {
    var finalWords []string
    var selectedLists []string
    var wordIndices []string

    // Generate and select words one by one
    for len(finalWords) < wordCount {
        words, isSecondary, indices, err := dg.generateWordChoices()
        if err != nil {
            return nil, err
        }

        fmt.Print(colorBlue + "\nChoisissez un mot (1-" + fmt.Sprintf("%d", CHOICES_COUNT) + ") parmi les propositions :\n" + colorReset)
        for i, word := range words {
            listType := "Liste 1"
            if isSecondary[i] {
                listType = "Liste 2"
            }
            fmt.Printf("%s%d: %s (%s)%s\n", colorCyan, i+1, word, listType, colorReset)
        }


	var choice int
	for {
	    fmt.Print(colorBlue + "Votre choix : " + colorReset)
	    
	    // Lecture de l'entrée utilisateur
	    if _, err := fmt.Scan(&choice); err != nil {
		fmt.Println(colorRed + "Erreur : Entrée invalide. Veuillez saisir un nombre." + colorReset)
		continue // Recommencer la boucle en cas d'erreur
	    }

	    // Vérification de la plage du choix
	    if choice >= 1 && choice <= CHOICES_COUNT {
		break // Sortie de la boucle si le choix est valide
	    }

	    fmt.Printf(colorRed+"Erreur : Choisissez un nombre entre 1 et %d\n"+colorReset, CHOICES_COUNT)
	}

        idx := choice - 1
        finalWords = append(finalWords, words[idx])
        if isSecondary[idx] {
            selectedLists = append(selectedLists, "secondary")
        } else {
            selectedLists = append(selectedLists, "primary")
        }
        wordIndices = append(wordIndices, indices[idx])

        remaining := wordCount - len(finalWords)
        if remaining > 0 {
            suffix := ""
            if remaining > 1 {
                suffix = "s"
            }
            fmt.Printf("%s\nEncore %d mot%s à choisir%s\n",
                colorYellow, remaining, suffix, colorReset)
        }
    }

    // Apply transformations to selected words
    for i := range finalWords {
        if firstLastCapitalize && (i == 0 || i == len(finalWords)-1) {
            finalWords[i] = strings.Title(finalWords[i])
        } else if randomCapitalize {
            shouldCap, err := dg.getSecureRandom(2)
            if err != nil {
                return nil, err
            }
            if shouldCap == 1 {
                finalWords[i] = strings.Title(finalWords[i])
            }
        }

        if addSpecial && !unicode.IsUpper(rune(finalWords[i][0])) {
            shouldAddSpecial, err := dg.getSecureRandom(2)
            if err != nil {
                return nil, err
            }
            if shouldAddSpecial == 1 {
                finalWords[i], err = dg.transformWord(finalWords[i], false, true)
                if err != nil {
                    return nil, err
                }
            }
        }
    }

    passphrase := strings.Join(finalWords, " ")
    
    // Calculate entropy (log2(2*7776) ≈ 13.92 bits per word)
    baseEntropy := float64(wordCount) * 13.92
    if addSpecial {
        baseEntropy += 4.17 * float64(wordCount/2)
    }

    var qrCode string
    if showQR {
        qrCode, _ = dg.generateQRCode(passphrase)
    }

    return &PassphraseResult{
        Passphrase:   passphrase,
        SelectedList: selectedLists,
        WordIndices:  wordIndices,
        Entropy:      baseEntropy,
        QRCode:       qrCode,
    }, nil
}

// GeneratePassphrase handles automatic passphrase generation
func (dg *DicewareGenerator) GeneratePassphrase(wordCount int, randomCapitalize, firstLastCapitalize, addSpecial bool) (*PassphraseResult, error) {
    if wordCount < 2 || wordCount > 8 {
        return nil, fmt.Errorf("le nombre de mots doit être entre 2 et 8")
    }

    var words []string
    var selectedLists []string
    var wordIndices []string

    // Generate all words automatically
    for i := 0; i < wordCount; i++ {
        useSecondary, err := dg.getSecureRandom(2)
        if err != nil {
            return nil, err
        }

        index, err := dg.rollDice()
        if err != nil {
            return nil, err
        }

        var word string
        if useSecondary == 1 {
            word = dg.secondaryWordlist[index]
            selectedLists = append(selectedLists, "secondary")
        } else {
            word = dg.primaryWordlist[index]
            selectedLists = append(selectedLists, "primary")
        }
        wordIndices = append(wordIndices, index)

        if firstLastCapitalize && (i == 0 || i == wordCount-1) {
            word = strings.Title(word)
        } else if randomCapitalize {
            shouldCap, err := dg.getSecureRandom(2)
            if err != nil {
                return nil, err
            }
            if shouldCap == 1 {
                word = strings.Title(word)
            }
        }

        if addSpecial && !unicode.IsUpper(rune(word[0])) {
            shouldAddSpecial, err := dg.getSecureRandom(2)
            if err != nil {
                return nil, err
            }
            if shouldAddSpecial == 1 {
                word, err = dg.transformWord(word, false, true)
                if err != nil {
                    return nil, err
                }
            }
        }

        words = append(words, word)
    }

    passphrase := strings.Join(words, " ")
    baseEntropy := float64(wordCount) * 13.92
    if addSpecial {
        baseEntropy += 4.17 * float64(wordCount/2)
    }

    var qrCode string
    if showQR {
        qrCode, _ = dg.generateQRCode(passphrase)
    }

    return &PassphraseResult{
        Passphrase:   passphrase,
        SelectedList: selectedLists,
        WordIndices:  wordIndices,
        Entropy:      baseEntropy,
        QRCode:       qrCode,
    }, nil
}

// readUserInput safely reads user input with proper error handling
func readUserInput() (string, error) {
    reader := bufio.NewReader(os.Stdin)
    input, err := reader.ReadString('\n')
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(input), nil
}

func main() {
    // Display program header
    fmt.Print(colorPurple + "\n=== Générateur de Passphrase Diceware ===\n" + colorReset)
    fmt.Println(colorYellow + "Utilisation de deux listes de mots pour une meilleure entropie" + colorReset)

    // Initialize generator
    generator, err := NewDicewareGenerator()
    if err != nil {
        log.Fatal(colorRed + "Erreur d'initialisation : " + err.Error() + colorReset)
    }

    // Get generation mode
    var mode string
    for {
        fmt.Print(colorBlue + "\nMode (1: Automatique, 2: Interactif) : " + colorReset)
        mode, err = readUserInput()
        if err != nil {
            fmt.Println(colorRed + "Erreur de saisie. Veuillez réessayer." + colorReset)
            continue
        }
        if mode == "1" || mode == "2" {
            break
        }
        fmt.Println(colorRed + "Erreur : Choisissez 1 ou 2" + colorReset)
    }

	var wordCount int
	for {
	    fmt.Print(colorBlue + "\nNombre de mots (2-8) : " + colorReset)
	    
	    // Lecture de l'entrée utilisateur
	    input, err := readUserInput()
	    if err != nil {
		fmt.Println(colorRed + "Erreur de saisie. Veuillez réessayer." + colorReset)
		continue // Redemande à l'utilisateur
	    }
	    
	    // Tentative de conversion de l'entrée en entier
	    if _, err := fmt.Sscan(input, &wordCount); err != nil {
		fmt.Println(colorRed + "Erreur : entrée invalide. Veuillez saisir un nombre entier." + colorReset)
		continue // Redemande à l'utilisateur
	    }

	    // Vérification si le nombre est dans la plage valide
	    if wordCount >= 2 && wordCount <= 8 {
		break // Sort de la boucle
	    }

	    fmt.Println(colorRed + "Erreur : Le nombre doit être entre 2 et 8." + colorReset)
	}


    // Get capitalization preference
    var capitalize string
    for {
        fmt.Print(colorBlue + "\nChoisissez la capitalisation :\n" +
            "  1: Premier et dernier mot\n" +
            "  2: Aléatoire\n" +
            "  N: Aucune\n" +
            "Votre choix : " + colorReset)
        capitalize, err = readUserInput()
        if err != nil {
            fmt.Println(colorRed + "Erreur de saisie. Veuillez réessayer." + colorReset)
            continue
        }
        capitalize = strings.ToUpper(capitalize)
        if capitalize == "1" || capitalize == "2" || capitalize == "N" {
            break
        }
        fmt.Println(colorRed + "Erreur : Choix invalide" + colorReset)
    }

    // Get special character preference
    var special string
    specialCount := 2
    if wordCount < 5 {
        specialCount = 1
    }
    pluralS := ""
    if specialCount > 1 {
        pluralS = "s"
    }

    for {
        fmt.Printf(colorBlue+"Caractères spéciaux sur %d mot%s non capitalisé%s (O/N) ? "+colorReset,
            specialCount, pluralS, pluralS)
        special, err = readUserInput()
        if err != nil {
            fmt.Println(colorRed + "Erreur de saisie. Veuillez réessayer." + colorReset)
            continue
        }
        special = strings.ToUpper(special)
        if special == "O" || special == "N" {
            break
        }
        fmt.Println(colorRed + "Erreur : Répondez par O ou N" + colorReset)
    }

    // Get QR code preference
    var qrChoice string
    for {
        fmt.Print(colorBlue + "Afficher un QR Code scannable (O/N) ? " + colorReset)
        qrChoice, err = readUserInput()
        if err != nil {
            fmt.Println(colorRed + "Erreur de saisie. Veuillez réessayer." + colorReset)
            continue
        }
        qrChoice = strings.ToUpper(qrChoice)
        if qrChoice == "O" || qrChoice == "N" {
            break
        }
        fmt.Println(colorRed + "Erreur : Répondez par O ou N" + colorReset)
    }
    showQR = qrChoice == "O"

    // Generate passphrase
    var result *PassphraseResult
    if mode == "2" {
        result, err = generator.InteractiveGeneration(
            wordCount,
            capitalize == "2",
            capitalize == "1",
            special == "O",
        )
    } else {
        result, err = generator.GeneratePassphrase(
            wordCount,
            capitalize == "2",
            capitalize == "1",
            special == "O",
        )
    }

    if err != nil {
        log.Fatal(colorRed + "Erreur : " + err.Error() + colorReset)
    }

    // Display results
    fmt.Printf("\n%sPassphrase : %s%s%s\n",
        colorGreen, colorCyan, result.Passphrase, colorReset)
    fmt.Printf("%sEntropie   : %.1f bits%s\n",
        colorYellow, result.Entropy, colorReset)
    
    if showQR && result.QRCode != "" {
        fmt.Printf("%sQR Code :\n%s%s%s\n",
            colorPurple, colorWhite, result.QRCode, colorReset)
    }
}
