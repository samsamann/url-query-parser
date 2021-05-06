package querystring

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isInteger(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isEqualSign(ch rune) bool {
	return ch == '='
}

func isAmpersand(ch rune) bool {
	return ch == '&'
}

func isOpenBracket(ch rune) bool {
	return ch == '['
}

func isCloseBracket(ch rune) bool {
	return ch == ']'
}

func isOperator(tok Token) bool {
	if _, ok := operators[tok]; ok {
		return true
	}
	return false
}

func isPageKeyword(tok Token) bool {
	if _, ok := pageKeywords[tok]; ok {
		return true
	}
	return false
}
