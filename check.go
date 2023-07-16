package lambda

func CheckNotEmpty(strs ...string) bool {
    for _, str := range strs {
        if str == "" {
            return false
        }
    }
    return true
}

func CheckLengthMIN(str string, min int) bool {
    return len(str) >= min
}

func CheckLengthMAX(str string, max int) bool {
    return len(str) <= max
}

func CheckLength(str string, min, max int) bool {
    l := len(str)
    return min <= l && l <= max
}
