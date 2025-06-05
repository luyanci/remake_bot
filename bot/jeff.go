package jeff

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "time"
)

var filename := "jeff.txt"

func getLineCount() (int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineCount := 0
    for scanner.Scan() {
        lineCount++
    }

    if err := scanner.Err(); err != nil {
        return 0, err
    }

    return lineCount, nil
}

func getRandomLine() (string, error) {
    lineCount, err := getLineCount(filename)
    if err != nil {
        return "", err
    }

    rand.Seed(time.Now().UnixNano())
    randomLineNumber := rand.Intn(lineCount)

    file, err := os.Open(filename)
    if err != nil {
        return "", err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    currentLine := 0
    for scanner.Scan() {
        if currentLine == randomLineNumber {
            return scanner.Text(), nil
        }
        currentLine++
    }

    if err := scanner.Err(); err != nil {
        return "", err
    }

    return "", nil
}

func replacePlaceholder(text string, replacement string) string {
    return fmt.Sprintf(text, replacement)
}
