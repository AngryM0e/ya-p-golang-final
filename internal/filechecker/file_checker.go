package filechecker

import (
    "os"
)

// CheckFileExists проверяет, существует ли файл по указанному пути
func CheckFileExists(filePath string) bool {
    _, err := os.Stat(filePath)
    if os.IsNotExist(err) {
        return false
    }
    return true
}