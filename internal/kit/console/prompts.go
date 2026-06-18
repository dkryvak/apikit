package console

import (
    "bufio"
    "fmt"
    "strings"
)

func Prompt(reader *bufio.Reader, message string) (string, error) {
    for {
        fmt.Printf("%s: ", message)
        input, err := reader.ReadString('\n')
        if err != nil {
            return "", err
        }
        input = strings.TrimSpace(input)
        if input != "" {
            return input, nil
        }
        fmt.Println("  This field is required. Please enter a value.")
    }
}

// PromptOptional reads a single line and allows an empty value (returns "").
func PromptOptional(reader *bufio.Reader, message string) (string, error) {
    fmt.Printf("%s: ", message)
    input, err := reader.ReadString('\n')
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(input), nil
}

func PromptYesNo(reader *bufio.Reader, message string) (bool, error) {
    for {
        fmt.Printf("%s (y/n): ", message)
        input, err := reader.ReadString('\n')
        if err != nil {
            return false, err
        }
        input = strings.ToLower(strings.TrimSpace(input))
        if input == "y" || input == "yes" {
            return true, nil
        }
        if input == "n" || input == "no" {
            return false, nil
        }
        fmt.Println("  Please enter 'y' or 'n'")
    }
}

func PromptMultiline(reader *bufio.Reader, message string) (string, error) {
    fmt.Println(message)
    var lines []string
    emptyLineCount := 0

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            return "", err
        }

        line = strings.TrimRight(line, "\n\r")

        if strings.TrimSpace(line) == "" {
            emptyLineCount++
            if emptyLineCount >= 2 {
                break
            }
            continue
        }

        emptyLineCount = 0
        lines = append(lines, line)
    }

    return strings.Join(lines, ""), nil
}
