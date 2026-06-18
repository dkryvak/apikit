package console

import (
    "fmt"
    "strings"
)

type Schema struct {
    Title   string
    Intro   string
    Fields  []SchemaField
    Example string
    Notes   []string
}

type SchemaField struct {
    Name        string
    Required    bool
    Type        string
    Description string
    Example     string
}

func PrintSchema(schema Schema) {
    sep := strings.Repeat("-", 72)

    fmt.Println(sep)
    fmt.Println(schema.Title)
    fmt.Println(sep)

    if strings.TrimSpace(schema.Intro) != "" {
        fmt.Println(schema.Intro)
        fmt.Println()
    }

    fmt.Println("Fields:")
    for _, f := range schema.Fields {
        req := "optional"
        if f.Required {
            req = "required"
        }

        fmt.Printf("  %-12s  %s\n", f.Name, req)

        if f.Type != "" {
            fmt.Printf("    Type: %s\n", f.Type)
        }
        if f.Description != "" {
            fmt.Printf("    Description: %s\n", f.Description)
        }
        if f.Example != "" {
            fmt.Printf("    Example: %s\n", f.Example)
        }
        fmt.Println()
    }

    if len(schema.Notes) > 0 {
        fmt.Println("Notes:")
        for _, n := range schema.Notes {
            if strings.TrimSpace(n) == "" {
                continue
            }
            fmt.Printf("  - %s\n", n)
        }
        fmt.Println()
    }

    if strings.TrimSpace(schema.Example) != "" {
        fmt.Println("Example:")
        fmt.Println(schema.Example)
    }

    fmt.Println(sep)
}
