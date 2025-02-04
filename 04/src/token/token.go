package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//식별자 + 리터럴
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"
	//연산자
	ASSIGN = "="
	PLUS   = "+"

	//구분자
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	//예약어
	FUNCTION = "FUNCTION"
	LET      = "LET"

	BANG     = "!"
	MINUS    = "-"
	SLASH    = "/"
	ASTERISK = "*"
	LT       = "<"
	GT       = ">"
	IF       = "if"
	RETURN   = "return"
	TRUE     = "true"
	ELSE     = "else"
	FALSE    = "false"
	EQ       = "=="
	NOT_EQ   = "!="
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"return": RETURN,
	"true":   TRUE,
	"else":   ELSE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
