# Monkey 인터프리터

## 🟨 0. 인터프리터 모식도

```go
┌─────────────────────────┐
│        REPL Loop        │
│  (Start 함수 실행)      │
└────────────┬────────────┘
             │
             ▼
    ┌─────────────────┐
    │    입력 (line)  │
    └─────────────────┘
             │
             ▼
    ┌─────────────────┐
    │     Lexer       │
    └─────────────────┘
             │
             ▼
    ┌─────────────────┐
    │     Parser      │
    └─────────────────┘
             │
             ▼
    ┌─────────────────┐
    │      AST        │
    │ (Program, 등)   │
    └─────────────────┘
             │
             ▼
    ┌─────────────────┐
    │   Evaluator     │
    │  (Eval 함수)    │
    └─────────────────┘
             │
             ▼
    ┌─────────────────┐
    │  결과 출력      │
    └─────────────────┘
```

## 🟨 1. 렉싱 (Lexing / 어휘 분석)

### 🔹 정의

소스 코드를 읽어 의미 있는 단위인 **토큰(Token)** 으로 분리하는 과정이다.

📌 즉, **"구문 → 토큰"** 으로 변환하는 단계이다.

### 🔹 주요 구성 요소

✅ **토큰 정의**

각 토큰은 **타입(Token Type)** 과 **리터럴(Literal)** 정보를 가진다.

📍 예) `LET`, `IDENT`, `INT`, `+`, `-`, `;` 등

📍 리터럴은 값을 문자열로 표현하여 단순화를 꾀한다.

✅ **렉서(Lexer) 구조체**

렉서는 다음과 같은 필드를 관리한다.

🔸 **현재 읽고 있는 문자 위치** (`current position`)

🔸 **다음 읽을 문자의 위치** (`next position`)

🔸 **현재 가리키는 문자** (`current character`)

🔸 **입력 문자열** (`input`)

### 🔹 동작 과정

1️⃣ **초기화**

- `readChar` 메서드를 호출하여 첫 문자를 읽고, 내부 포지션을 설정한다.

2️⃣ **토큰 분리 (nextToken 메서드)**

- **식별자/예약어**:문자가 공백이나 특수문자가 나오기 전까지 읽고, 예약어 맵과 비교하여 `IDENT` 또는 예약어 토큰으로 결정한다.→ `readIdentifier` 함수 사용
- **숫자 (정수) 판별**:숫자를 연속해서 읽어 정수 리터럴 토큰을 생성한다.→ `readNumber` 함수 사용
- **연산자 및 특수문자 판별**:단일 또는 복합 연산자(`==`, `!=` 등)를 인식한다.
- **공백 제거**:의미 없는 공백을 건너뛴다.

3️⃣ **REPL과 테스트**

- 렉서가 올바르게 토큰을 생성하는지 **REPL 환경에서 검증**한다.

---

## 🟨 2. 파싱 (Parsing)

### 🔹 정의

렉서가 만든 **토큰 스트림을 받아** 추상 구문 트리(**AST, Abstract Syntax Tree**) 로 변환하는 과정이다.

📌 즉, **"토큰 → AST"** 단계이다.

### 🔹 동작 순서 (예시)

1. **`ParseProgram()`** → 파싱 프로그램 생성
2. **`ParseStatement()`** → 토큰에 따라 노드 생성 처리 분기
3. **`parseExpressionStatement()`** → `ExpressionStatement` 생성
4. **`parseExpression()`  → 5번 메소드 호출**
5. **`parsePrefixExpression()`** → 토큰 타입에 따라 파싱 함수 호출 후, `Right`에 대한 토큰 타입에 따라 재귀적으로 파싱 함수 호출
6. **`parseExpression()`** → 정수 리터럴 또는 식별자에 대한 `Expression` 생성

✅ **이렇게 구현하면 `-5` 또는 `!true` 같은 코드가 올바르게 AST로 변환됨!**

### 🔹 AST의 기본 구조

```go

Node (interface)
├── Statement (interface)
│   ├── LetStatement (struct)
│   ├── ReturnStatement (struct)
│   └── ExpressionStatement (struct)
└── Expression (interface)
    ├── Identifier (struct)
    ├── IntegerLiteral (struct)
    ├── PrefixExpression (struct)
    └── InfixExpression (struct)

```

---

### 📌 2.1. Let 문 파싱

### 🔹 문법

```go

let <identifier> = <expression
```

### 🔹 구조체

```go
type LetStatement struct {
    Token token.Token // 'let' 토큰
    Name  *Identifier // 변수명을 나타내는 Identifier 노드
    Value Expression  // 할당될 표현식
}
```

### 🔹 동작

1️⃣ 현재 토큰(`let`)을 확인한다.

2️⃣ 다음 토큰이 **식별자(IDENT)** 인지 검사하고, `Identifier` 노드를 생성한다.

3️⃣ `=` 토큰과 그 뒤의 표현식을 파싱하여 `Value` 필드에 할당한다.

### 🔹 **각 필드의 값 예시:**

| 필드명 | 값 (예시) | 설명 |
| --- | --- | --- |
| `Token` | `{Type: token.LET, Literal: "let"}` | `let` 토큰을 저장 |
| `Name` | `&Identifier{Token: {Type: token.IDENT, Literal: "x"}, Value: "x"}` | 변수명 `x`를 저장 |
| `Value` | `&IntegerLiteral{Token: {Type: token.INT, Literal: "5"}, Value: 5}` | 값 `5`를 저장 |

---

### 📌 2.2. Return 문 파싱

### 🔹 문법

```

return <expression>;
```

### 🔹 동작

1️⃣ `return` 토큰을 만나면,

2️⃣ 뒤따르는 표현식을 파싱하여 `ReturnStatement` AST 노드를 생성한다.

---

### 📌 2.3. 표현식 파싱 & Pratt 파싱 기법

### 🔹 표현식(Expression)

산술 연산, 변수 참조, 함수 호출 등 **다양한 식(expression)을 처리**한다.

### 🔹 **Pratt 파싱 기법**

✅ **개념**

- 연산자의 **우선순위**와 **결합 규칙**을 자연스럽게 처리하기 위한 파싱 기법이다.

✅ **동작**

1️⃣ 각 토큰에 대해 **"좌측 표현식"** 과 **"우측 표현식"** 을 재귀적으로 파싱한다.

2️⃣ 연산자 우선순위에 따라 올바른 **AST** 를 구성한다.

---

### 📌 2.4. 식별자(Identifier)와 정수 리터럴(Integer Literal)

✅ **식별자 (Identifier)**

```go
go
복사편집
type Identifier struct {
    Token token.Token // 변수명을 나타내는 토큰
    Value string      // 실제 변수명 (예: "x")
}

```

✅ **정수 리터럴 (IntegerLiteral)**

- 정수 값을 담은 노드이다.
- 렉서에서 읽은 숫자 문자열을 **정수 값으로 변환**하여 저장한다.

---

### 📌 2.5. 전위 연산자 파싱 (Prefix Expression)

### 🔹 문법

```

<operator> <expression>  // 예: -5, !true
```

### 🔹 구조체

```go

type PrefixExpression struct {
    Token    token.Token // 전위 연산자 토큰
    Operator string      // 연산자 기호 (예: "-", "!")
    Right    Expression  // 오른쪽 표현식
}
```

### 🔹 동작

1️⃣ 전위 연산자 토큰을 만나면,

2️⃣ 해당 토큰과 함께 **다음 표현식**을 재귀적으로 파싱하여 `Right` 필드를 채운다.

---

### 📌 2.6. 중위 연산자 파싱 (Infix Expression)

### 🔹 문법

```

<left expression> <operator> <right expression>  // 예: 5 + 5
```

### 🔹 구조체

```go

type InfixExpression struct {
    Token    token.Token // 연산자 토큰
    Left     Expression  // 왼쪽 표현식
    Operator string      // 연산자 기호 (예: "+")
    Right    Expression  // 오른쪽 표현식
}
```

### 🔹 동작

1️⃣ 먼저 **왼쪽 표현식**을 파싱한다.

2️⃣ 현재 토큰이 **중위 연산자**이면,

3️⃣ 연산자 **우선순위를 고려**해 **오른쪽 표현식**을 재귀적으로 파싱한다.

---

### 📌 2.7. 불리언 리터럴 (Boolean Literal)

### 🔹 문법

```

true 또는 false
```

### 🔹 동작

- `true`, `false` 토큰을 만나면, **불리언 값을 담은 AST 노드**를 생성한다.

---

### 📌 2.8. 그룹 표현식 (Grouped Expression)

### 🔹 문법

```

( <expression> )
```

### 🔹 동작

1️⃣ **괄호**로 묶인 표현식을 만나면,

2️⃣ 내부 표현식을 별도로 **파싱하여 그룹화된 AST 노드**로 생성한다.

📌 이를 통해 **우선순위를 명확하게 조절**한다.

---

### 📌 2.9. 조건 표현식 (If Expression)

### 🔹 문법

```

if (<condition>) <consequence> else <alternative>

```

### 🔹 구조체

```go

type IfExpression struct {
    Token       token.Token     // 'if' 토큰
    Condition   Expression      // 조건식
    Consequence *BlockStatement // 참일 때 실행할 블록
    Alternative *BlockStatement // (옵션) 거짓일 때 실행할 블록
}
```

### 🔹 동작

1️⃣ `if` 토큰을 만나면, **조건식과 실행 블록**을 파싱한다.

2️⃣ `else` 블록이 있으면 함께 파싱하여 하나의 **조건 표현식 AST 노드**로 구성한다.

---

### 📌 2.10. 함수 리터럴 & 호출 표현식 (Function Literal & Call Expression)

✅ **함수 리터럴(Function Literal)**

```go

type FunctionLiteral struct {
    Token      token.Token   // 'fn' 토큰
    Parameters []*Identifier // 매개변수 목록
    Body       *BlockStatement // 함수 본문
}
```

✅ **호출 표현식(Call Expression)**

```go

type CallExpression struct {
    Token     token.Token   // '(' 토큰
    Function  Expression    // 호출 대상 (예: 식별자 또는 함수 리터럴)
    Arguments []Expression  // 인자 목록
}
```

📌 **동작:**

- **함수 리터럴**은 매개변수와 본문을 파싱하여 `FunctionLiteral` 노드를 생성한다.
- **함수 호출**은 호출 대상과 인자들을 파싱하여 `CallExpression` 노드를 구성한다. 🚀

## 🟨 3. 평가 (Evaluation)

**객체 시스템을 정의하고 표현하는 과정**으로, 프로그래밍 언어에서 데이터 타입을 다룰 수 있도록 한다.

---

### 📌 3.1. 객체 표현하기 (Object Representation)

객체 시스템은 **언어 내부에서 다루는 값(value)을 표현**하는 방법이다.

✅ **기본 객체 타입**

- **정수 (Integer)**
- **불리언 (Boolean)**
- **널 (Null)**

### 🔹 객체 시스템의 기초

객체는 **Object 인터페이스**를 기반으로 한다.

```go

type ObjectType string

type Object interface {
    Type() ObjectType
    Inspect() string
}
```

📌 **`Type()`**: 객체의 타입을 반환한다.

📌 **`Inspect()`**: 객체의 값을 문자열로 변환한다.

---

### 🔹 정수(Integer) 객체

```go

type Integer struct {
    Value int64
}
```

📌 **정수를 표현하는 객체**이며, `Value` 필드에 값을 저장한다.

---

### 🔹 불리언(Boolean) 객체

```go

type Boolean struct {
    Value bool
}
```

📌 **`true` 또는 `false` 값을 가지는 객체**이다.

---

### 🔹 Null 객체

```go

type Null struct{}
```

📌 **"값이 없음"을 나타내는 객체**이다.

---

### 📌 3.2. 표현식 평가 (Expression Evaluation)

📌 **표현식을 해석하고 적절한 값을 반환하는 과정**이다.

📌 **내부 평가(Eval) 함수의 재귀적 동작:**

- **AST 노드별 평가**:
    
    `Eval` 함수는 AST의 각 노드 타입(예: ExpressionStatement, BlockStatement, InfixExpression 등)에 대해 자신을 재귀적으로 호출하여 하위 노드를 평가한다.
    
- **함수 호출 및 블록 평가**:
    
    함수 리터럴이나 블록 구문을 평가할 때, 새 환경(Env)을 생성하고, 해당 환경에서 내부 구문을 재귀적으로 평가한다.
    

### ✅ 정수 리터럴 평가

```go

func evalIntegerLiteral(node *ast.IntegerLiteral) Object {
    return &Integer{Value: node.Value}
}
```

📌 **AST 노드에서 값을 읽어 Integer 객체를 생성한다.**

---

### ✅ REPL

📌 **REPL(Read-Eval-Print Loop)은 사용자가 입력한 코드를 평가하여 즉시 결과를 출력하는 환경이다.**

1️⃣ **입력을 받는다.**

2️⃣ **파싱하여 AST를 만든다.**

3️⃣ **AST를 평가하여 결과 객체를 반환한다.**

4️⃣ **객체의 `Inspect()` 메서드를 호출해 출력한다.**

```go
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParserProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
```

---

### ✅ 불리터럴 평가

```go

func evalBooleanLiteral(node *ast.Boolean) Object {
    return &Boolean{Value: node.Value}
}
```

📌 **`true` 또는 `false` 값을 가진 객체를 생성한다.**

---

### ✅ Null 평가

```go

var NULL = &Null{}
```

📌 **항상 동일한 NULL 객체를 반환한다.**

---

### ✅ 전위 표현식 평가

```go

func evalPrefixExpression(operator string, right Object) Object {
    switch operator {
    case "!":
        return evalBangOperatorExpression(right)
    case "-":
        return evalMinusPrefixOperatorExpression(right)
    }
    return NULL
}
```

📌 **예제**

```

-5  → -5
!true  → false
```

---

### ✅ 중위 표현식 평가

```go

func evalInfixExpression(operator string, left, right Object) Object {
    switch {
    case operator == "+":
        return evalAdditionExpression(left, right)
    case operator == "-":
        return evalSubtractionExpression(left, right)
    }
    return NULL
}
```

---

### 📌 3.3. 조건식 평가 (If Expression)

📌 **조건식의 결과에 따라 실행할 코드 블록을 결정한다.**

```go

func evalIfExpression(ie *ast.IfExpression) Object {
    condition := eval(ie.Condition)
    if isTruthy(condition) {
        return eval(ie.Consequence)
    } else if ie.Alternative != nil {
        return eval(ie.Alternative)
    }
    return NULL
}
```

📌 `isTruthy()` 함수로 조건이 참인지 검사하여 적절한 블록을 실행한다.

---

### 📌 3.4. Return 문 평가

📌 **함수 실행을 중단하고 값을 반환하는 역할을 한다.**

```go

type ReturnValue struct {
    Value Object
}
```

- 평가 시, `return` 문을 만나면 해당 표현식을 평가한 뒤, 그 값을 **ReturnValue** 객체에 감싸서 반환한다.
- 예를 들어, `return 5;` 는 `&object.ReturnValue{Value: 5}`와 같이 평가된다.

📌 **`ReturnValue` 객체를 만나면 즉시 평가를 종료할 수 있는 이유**

- **evalProgram** 및 **evalBlockStatement** 함수에서는 각 문장을 평가한 후 결과의 타입을 확인한다.
- 만약 평가 결과가 `ReturnValue` (또는 에러) 타입이면, 바로 상위 호출로 해당 값을 반환하고 더 이상 후속 문장을 평가하지 않는다.
- 이렇게 함으로써 함수나 블록 내에서 Return 문이 등장하면 즉시 평가를 종료할 수 있다.
- 예시코드(evalProgram 일부):

```go
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
    var result object.Object
    for _, statement := range program.Statements {
        result = Eval(statement, env)
        switch result := result.(type) {
        case *object.ReturnValue:
            return result.Value  // ReturnValue를 만나면 즉시 반환하여 종료
        case *object.Error:
            return result
        }
    }
    return result
}

```

### 📌 3.5. 함수와 환경 (Function and Environment)

- 함수 리터럴 생성 시점의 환경 캡처
    
    함수 리터럴(예: `fn(x) { ... }`)이 평가될 때, 현재의 변수 바인딩(환경)을 함께 캡처한다.
    
    → 이 캡처된 환경은 함수 객체의 한 부분(예: `Env` 필드)으로 저장된다.
    
    예시 : 
    
    ```go
    
    // 함수 리터럴이 생성될 때의 환경(env)을 함께 캡처하여 함수 객체 생성
    &object.Function{Parameters: params, Body: body, Env: env}
    ```
    
- 환경(Environment) 구조체
    
    인터프리터는 변수와 그 값을 저장하기 위해 환경 객체를 사용한다.
    
    함수 리터럴이 생성될 때의 환경을 함수 객체에 저장함으로써,
    
    함수가 호출될 때 원래 정의된 환경(즉, 클로저)을 그대로 유지할 수 있게 된다.
    
    예시 : 
    
    ```go
    
    func NewEnvironment() *Environment {
    	s := make(map[string]Object)
    	return &Environment{store: s, outer: nil}
    }
    ```
    

---

### 📌 3.6. 클로저

- 클로저(Closure)란?
    
    클로저는 자신이 생성될 때의 환경(스코프)을 기억하는 함수이다.
    
    → 함수가 반환된 이후에도, 그 내부에서 외부 변수에 접근할 수 있게 된다.
    
    😎 즉, 클로저 덕분에 함수는 자신이 정의된 시점의 변수 바인딩을 계속 유지한다.
    
- 구현 방식
    1. 함수 리터럴 평가 시
        
        파서가 함수 리터럴을 만나면, 함수의 매개변수 목록과 본문, 그리고 현재의 환경을 포함하는 함수 객체를 생성한다.
        
        ```go
        
        	case *ast.FunctionLiteral:
        		params := node.Parameters
        		body := node.Body
        		return &object.Function{Parameters: params, Body: body, Env: env}
        ```
        
    2. 함수 호출 시 환경 확장
        
        함수를 호출할 때(CallExpression), `extendFunctionEnv` 함수를 사용하여 함수가 캡처한 환경을 기반으로 새로운 환경(Enclosed Environment) 을 생성한다.
        
        이 새로운 환경은 함수의 원래 환경(클로저) 위에, 호출 시 전달된 인자들을 매개변수 이름에 바인딩한다.
        
        ```go
        	func Eval(node ast.Node, env *object.Environment) object.Object{
        		//[...]
        		case *ast.CallExpression:
        			function := Eval(node.Function, env)
        			if isError(function) {
        				return function
        			}
        			args := evalExpressions(node.Arguments, env)
        			if len(args) == 1 && isError(args[0]) {
        				return args[0]
        			}
        			return applyFunction(function, args)
        		//[...]
        	}
        			
        		func applyFunction(fn object.Object, args []object.Object) object.Object {
        			switch fn := fn.(type) {
        			case *object.Function:
        				extendedEnv := extendFunctionEnv(fn, args)
        				evaluated := Eval(fn.Body, extendedEnv)
        				return unwrapReturnValue(evaluated)
        			
        			case *object.Builtin:
        				return fn.Fn(args...)
        			
        			default:
        				return newError("not a function: %s", fn.Type())
        			}
        		}
        		
        		// 함수가 가진 환경이 새로 만들어진 환경을 감싸고, 새로 만들어진 환경안에서 인수들을 바인딩
        		func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
        			env := object.NewEnclosedEnvironment(fn.Env)
        		
        			for paramIdx, param := range fn.Parameters {
        				env.Set(param.Value, args[paramIdx])
        			}
        			return env
        		}
        ```
        
        👉 이렇게 하면 함수 내부에서는 함수가 정의될 때의 외부 변수들과, 호출 시 전달된 인자 모두에 접근할 수 있다.
        

---

### 📌 3.7. 고차 함수 (Higher-order Functions)

- 고차 함수(Higher-order Functions)란?
    
    고차 함수는 함수를 인자로 받거나 함수를 반환하는 함수를 의미한다.
    
- 클로저를 통한 고차 함수 지원
    
    클로저는 함수가 자신의 환경을 기억하기 때문에,
    
    함수를 값으로 전달하거나 반환할 때도 그 환경 정보가 함께 전달된다.
    
    → 반환된 함수가 다른 함수의 인자로 전달되더라도, 원래 함수가 정의될 때의 환경(즉, 변수 바인딩)을 그대로 유지하여 정상 동작한다.
    
- 예시
    
    ```go
    
    func newAdder(x int) func(int) int {
        // 여기서 x의 값(예: 2)이 함수와 함께 캡처된다.
        return func(y int) int {
            return x + y // x는 newAdder가 생성될 때 캡처된 값이다.
        }
    }
    
    func main() {
        addTwo := newAdder(2)  // 클로저 생성: x=2를 캡처 😎\n      
        fmt.Println(addTwo(3)) // 2 + 3 = 5, addTwo는 고차 함수로서 사용됨 🚀
    }
    
    ```
    
    위 예제에서 `newAdder` 함수는 인자로 받은 `x`를 캡처하여 내부 익명 함수를 반환한다.
    
    반환된 함수는 클로저가 되어, 이후 호출 시에도 원래의 `x` 값에 접근할 수 있다.
    
    이 메커니즘을 통해 인터프리터는 고차 함수를 지원할 수 있다.
    

---

## 🟨 4. 확장 (Extensions)

**문자열, 내장 함수, 배열, 해시를 추가하여 기능을 확장한다.**

---

### 📌 4.1. 문자열 (String)

✅ **문자열 파싱**

```go

type String struct {
    Value string
```

📌 **문자열 리터럴을 `String` 객체로 변환한다.**

✅ **문자열 평가**

```go

func evalStringLiteral(node *ast.StringLiteral) Object {
    return &String{Value: node.Value}
}
```

📌 **AST에서 문자열 값을 읽어 `String` 객체를 생성한다.**

✅ **문자열 결합 (Concatenation)**

```go

func evalStringConcatenation(left, right Object) Object {
    return &String{Value: left.(*String).Value + right.(*String).Value}
}
```

📌 **예제**

```

"Hello" + " World!"  → "Hello World!"
```

---

### 📌 4.2. 내장 함수 (Built-in Functions)

✅ **`len` 함수 구현**

```go

var builtins = map[string]*Builtin{
    "len": &Builtin{
        Fn: func(args ...Object) Object {
            return &Integer{Value: int64(len(args[0].(*String).Value))}
        },
    },
}
```

📌 **예제**

```

len("Hello")  → 5
```

---

### 📌 4.3. 배열 (Array)

✅ **배열 표현**

```go

type Array struct {
    Elements []Object
}
```

📌 **여러 개의 객체를 저장하는 데이터 구조이다.**

✅ **배열 인덱싱**

```go

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}
```

📌 **예제**

```

arr = [1, 2, 3]
arr[0]  → 1
```

---

### 📌 4.4. 해시 (Hash)

✅ **해시 테이블 표현**

```go
type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
    Pairs map[HashKey]HashPair
}
```

✅ **해시 키**

```go

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
```

📌 **HashKey를 도입함으로써 값이 같은 문자열이 서로 다른 메모리에 위치해 올바른 해시값을 가져올 수 없는 상황을 방지하였다.**

📌 **해시는 키-값 저장소로 활용된다.**

✅ **해시 조회**

```go
func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}
```

📌 **예제**

```

h = {"name": "Alice", "age": 25}
h["name"]  → "Alice"
```

---

## 🟨 5. 얻은점

- Go 언어의 특징
    - 컴파일 시 바이트코드를 생성 후 실행 파일을 실행하고 삭제한다.
    - `:=`를 사용해 동적 타입처럼 보이지만, 사실은 정적 타입 언어이다.
    - 내장 테스트 기능 제공 → `"go test 파일이름"` 명령어로 `_test.go` 파일을 찾아 실행한다.
- TDD(테스트 주도 개발) 경험
    - 테스트 기반으로 점진적 개발을 진행하며, 소프트웨어 공학 이론을 실제로 경험하였다.
    - 에러 처리 로직을 철저히 해야 디버깅이 쉬워진다는 것을 몸소 느꼈다.
- 느낀점 및 다짐
    - 처음에는 Go 문법조차 익숙하지 않고 내용자체도 난해하여서 속도가 매우 느렸으나 끝까지 완주했다는 것에 대해서 뿌듯함을 느꼈다.
    - 프로그램이 복잡해질수록 코드의 가독성이 떨어졌고 로직순서도 헷갈려서 계속 디버깅을 돌려가면서 구조를 그림으로 그리고 파악하였다. 이 과정을 통해 프로그램이 복잡해질수록 다이어그램과 리팩토링이 필수적이라는 것을 느꼈다.
    - 더욱 정진하여서 훗날 지금보다 실력이 성장한 상태에서 이 책을 다시 보면서 지금을 회상할 것이다.
