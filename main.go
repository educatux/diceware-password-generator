package main

import (
    "bufio"
    "crypto/rand"
    "embed"
    "fmt"
    "log"
    "math"
    "math/big"
    "strings"
)

const (
    colorReset  = "\033[0m"
    colorRed    = "\033[31m"
    colorGreen  = "\033[32m"
    colorYellow = "\033[33m"
    colorBlue   = "\033[34m"
    colorPurple = "\033[35m"
    colorCyan   = "\033[36m"
)

//go:embed frenchdiceware.txt diceware-fr-alt.txt
var wordlistFS embed.FS

type DicewareGenerator struct {
    primaryWordlist   map[string]string
    secondaryWordlist map[string]string
    specialChars      string
    numberChars       string
}

func NewDicewareGenerator() (*DicewareGenerator, error) {
    dg := &DicewareGenerator{
        primaryWordlist:   make(map[string]string),
        secondaryWordlist: make(map[string]string),
        specialChars:      "!@#$%^&*",
        numberChars:       "0123456789",
    }

    if err := dg.loadWordlist("frenchdiceware.txt", dg.primaryWordlist); err != nil {
        return nil, err
    }

    if err := dg.loadWordlist("diceware-fr-alt.txt", dg.secondaryWordlist); err != nil {
        return nil, err
    }

    return dg, nil
}

func (dg *DicewareGenerator) loadWordlist(filename string, wordlist map[string]string) error {
    file, err := wordlistFS.Open(filename)
    if err != nil {
        return fmt.Errorf("erreur d'ouverture de %s: %w", filename, err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) > 6 {
            number := line[:5]
            word := strings.TrimSpace(line[5:])
            wordlist[number] = word
        }
    }
    return nil
}

func (dg *DicewareGenerator) rollDice() (string, error) {
    var roll strings.Builder
    for i := 0; i < 5; i++ {
        n, err := rand.Int(rand.Reader, big.NewInt(6))
        if err != nil {
            return "", err
        }
        roll.WriteString(fmt.Sprintf("%d", n.Int64()+1))
    }
    return roll.String(), nil
}

func (dg *DicewareGenerator) getRandomChar(chars string) (string, error) {
    n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
    if err != nil {
        return "", err
    }
    return string(chars[n.Int64()]), nil
}

func (dg *DicewareGenerator) getRandomSpecialChar() (string, error) {
    return dg.getRandomChar(dg.specialChars)
}

func (dg *DicewareGenerator) getRandomNumber() (string, error) {
    return dg.getRandomChar(dg.numberChars)
}

func (dg *DicewareGenerator) replaceRandomChar(word string) (string, error) {
    if len(word) == 0 {
        return word, nil
    }
    
    useSpecial, err := rand.Int(rand.Reader, big.NewInt(2))
    if err != nil {
        return "", err
    }

    var newChar string
    if useSpecial.Int64() == 0 {
        newChar, err = dg.getRandomSpecialChar()
    } else {
        newChar, err = dg.getRandomNumber()
    }
    if err != nil {
        return "", err
    }

    pos, err := rand.Int(rand.Reader, big.NewInt(int64(len(word))))
    if err != nil {
        return "", err
    }

    return word[:pos.Int64()] + newChar + word[pos.Int64()+1:], nil
}

func (dg *DicewareGenerator) GenerateChoices() ([]string, error) {
    choices := make([]string, 5)
    for i := 0; i < 5; i++ {
        useSecondary, err := rand.Int(rand.Reader, big.NewInt(2))
        if err != nil {
            return nil, err
        }

        roll, err := dg.rollDice()
        if err != nil {
            return nil, err
        }

        var word string
        if useSecondary.Int64() == 1 {
            word = dg.secondaryWordlist[roll]
        } else {
            word = dg.primaryWordlist[roll]
        }
        choices[i] = word
    }
    return choices, nil
}

func (dg *DicewareGenerator) calculateEntropy(wordCount int, hasSpecialChars bool) float64 {
    baseEntropy := float64(wordCount) * 25.8 // log2(7776*7776)

    if hasSpecialChars {
        specialCharEntropy := math.Log2(18)
        positionEntropy := math.Log2(float64(wordCount))
        charPositionEntropy := math.Log2(6)
        baseEntropy += specialCharEntropy + positionEntropy + charPositionEntropy
    }

    return baseEntropy
}

func (dg *DicewareGenerator) transformWords(words []string, randomCapitalize, firstLastCapitalize, addSpecial bool) error {
    if firstLastCapitalize {
        words[0] = strings.Title(words[0])
        words[len(words)-1] = strings.Title(words[len(words)-1])
    }

    if randomCapitalize {
        for i := range words {
            n, err := rand.Int(rand.Reader, big.NewInt(2))
            if err != nil {
                return err
            }
            if n.Int64() == 1 {
                words[i] = strings.Title(words[i])
            }
        }
    }

    if addSpecial {
        var nonCapitalizedIndices []int
        for i, word := range words {
            if word[0] >= 'a' && word[0] <= 'z' {
                nonCapitalizedIndices = append(nonCapitalizedIndices, i)
            }
        }

        specialCount := 2
        if len(words) < 5 {
            specialCount = 1
        }

        if len(nonCapitalizedIndices) >= specialCount {
            for i := 0; i < specialCount; i++ {
                n, err := rand.Int(rand.Reader, big.NewInt(int64(len(nonCapitalizedIndices))))
                if err != nil {
                    return err
                }
                idx := nonCapitalizedIndices[n.Int64()]
                words[idx], err = dg.replaceRandomChar(words[idx])
                if err != nil {
                    return err
                }
                nonCapitalizedIndices = append(nonCapitalizedIndices[:n.Int64()], nonCapitalizedIndices[n.Int64()+1:]...)
            }
        }
    }
    return nil
}

func (dg *DicewareGenerator) formatResult(words []string, hasSpecialChars bool) string {
    passphrase := strings.Join(words, " ")
    entropy := dg.calculateEntropy(len(words), hasSpecialChars)

    result := fmt.Sprintf("\n%sPassphrase : %s%s%s\n",
        colorGreen,
        colorCyan,
        passphrase,
        colorReset)
    result += fmt.Sprintf("%sEntropie : %.1f bits%s",
        colorYellow,
        entropy,
        colorReset)

    return result
}

func (dg *DicewareGenerator) GeneratePassphrase(wordCount int, randomCapitalize, firstLastCapitalize, addSpecial bool) (string, error) {
    if wordCount < 2 || wordCount > 8 {
        return "", fmt.Errorf("le nombre de mots doit être entre 2 et 8")
    }

    var words []string
    for i := 0; i < wordCount; i++ {
        useSecondary, err := rand.Int(rand.Reader, big.NewInt(2))
        if err != nil {
            return "", err
        }

        roll, err := dg.rollDice()
        if err != nil {
            return "", err
        }

        var word string
        if useSecondary.Int64() == 1 {
            word = dg.secondaryWordlist[roll]
        } else {
            word = dg.primaryWordlist[roll]
        }
        words = append(words, word)
    }

    err := dg.transformWords(words, randomCapitalize, firstLastCapitalize, addSpecial)
    if err != nil {
        return "", err
    }

    return dg.formatResult(words, addSpecial), nil
}

func (dg *DicewareGenerator) InteractiveGeneration(wordCount int, randomCapitalize, firstLastCapitalize, addSpecial bool) (string, error) {
    var selectedWords []string
    
    for len(selectedWords) < wordCount {
        choices, err := dg.GenerateChoices()
        if err != nil {
            return "", err
        }
        
        fmt.Print(colorBlue + "\nChoisissez un mot (1-5) parmi:\n" + colorReset)
        for i, word := range choices {
            fmt.Printf("%s%d: %s%s\n", colorCyan, i+1, word, colorReset)
        }
        
        var choice int
        for {
            fmt.Print(colorBlue + "Votre choix : " + colorReset)
            fmt.Scan(&choice)
            if choice >= 1 && choice <= 5 {
                break
            }
            fmt.Println(colorRed + "Erreur : Choisissez un nombre entre 1 et 5" + colorReset)
        }
        
        selectedWords = append(selectedWords, choices[choice-1])
        remainingWords := wordCount - len(selectedWords)
        if remainingWords > 0 {
            suffix := ""
            if remainingWords > 1 {
                suffix = "s"
            }
            fmt.Printf("%s\nEncore %d mot%s à choisir%s\n",
                colorYellow, remainingWords, suffix, colorReset)
        }
    }

    err := dg.transformWords(selectedWords, randomCapitalize, firstLastCapitalize, addSpecial)
    if err != nil {
        return "", err
    }

    return dg.formatResult(selectedWords, addSpecial), nil
}

func main() {
    fmt.Print(colorPurple + "\n=== Générateur de Passphrase Diceware ===" + colorReset + "\n\n")

    generator, err := NewDicewareGenerator()
    if err != nil {
        log.Fatal(colorRed + "Erreur : " + err.Error() + colorReset)
    }

    var mode string
    for {
        fmt.Print(colorBlue + "Mode (1: Automatique, 2: Interactif) : " + colorReset)
        fmt.Scan(&mode)
        if mode == "1" || mode == "2" {
            break
        }
        fmt.Println(colorRed + "Erreur : Choisissez 1 ou 2" + colorReset)
    }

    var wordCount int
    var capitalize, special string

    for {
        fmt.Print(colorBlue + "Nombre de mots (2-8) : " + colorReset)
        fmt.Scan(&wordCount)
        if wordCount >= 2 && wordCount <= 8 {
            break
        }
        fmt.Println(colorRed + "Erreur : Le nombre doit être entre 2 et 8" + colorReset)
    }

    for {
        fmt.Print(colorBlue + "Choisissez la capitalisation:\n" +
            "  1: Premier et dernier mot\n" +
            "  2: Aléatoire\n" +
            "  N: Aucune\n" +
            "Votre choix : " + colorReset)
        fmt.Scan(&capitalize)
        capitalize = strings.ToUpper(capitalize)
        if capitalize == "1" || capitalize == "2" || capitalize == "N" {
            break
        }
        fmt.Println(colorRed + "Erreur : Choix invalide" + colorReset)
    }

    useSpecial := wordCount >= 5
    specialCount := 2
    if !useSpecial {
        specialCount = 1
    }

    pluralS := ""
    if specialCount > 1 {
        pluralS = "s"
    }

    for {
        fmt.Printf(colorBlue+"Caractères spéciaux sur %d mot%s non capitalisé%s (O/N) ? "+colorReset,
            specialCount, pluralS, pluralS)
        fmt.Scan(&special)
        special = strings.ToUpper(special)
        if special == "O" || special == "N" {
            break
        }
        fmt.Println(colorRed + "Erreur : Répondez par O ou N" + colorReset)
    }

    var randomCap, firstLastCap bool
    switch capitalize {
    case "1":
        firstLastCap = true
    case "2":
        randomCap = true
    }

    var result string
    if mode == "2" {
        result, err = generator.InteractiveGeneration(
            wordCount,
            randomCap,
            firstLastCap,
            special == "O",
        )
    } else {
        result, err = generator.GeneratePassphrase(
            wordCount,
            randomCap,
            firstLastCap,
            special == "O",
        )
    }

    if err != nil {
        log.Fatal(colorRed + "Erreur : " + err.Error() + colorReset)
    }

    fmt.Println(result)
}
