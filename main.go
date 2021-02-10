package main

import "fmt"

func main() {
    data, err := loadGameContent()
    if err != nil {
        fmt.Printf("failed to load game data, error: %v", err)
        return
    }
    
    // Print
    for _, wonder := range data.Wonders {
        wonder.print()
        fmt.Printf("\n")
    }
}
