# 1. 렉싱 (Lexing)

## 정의

- "구문 → 토큰 → 추상구문트리"에서 **구문 → 토큰**으로 변환하는 과정
- 소스 코드를 읽고 의미 있는 단위(토큰)로 분리하는 역할

## 렉서 구조체

- **현재 읽고 있는 문자의 포지션**
- **다음 읽을 문자의 포지션**
- **현재 가리키고 있는 문자**
- **입력 문자열 (input)**

## 토큰 구조체

- **타입(Token Type)과 리터럴(Literal) 필드로 구성**
- 리터럴은 `string` 타입을 사용하여 모든 타입을 수용하고 복잡도를 낮춤

## 렉서 동작 과정

1. **렉서 초기화**
    - `readChar` 메서드를 호출하여 첫 번째 문자를 읽고, 구조체의 멤버 변수를 초기화함
2. **토큰 분리 (nextToken 메서드 실행)**
    - **식별자 or 예약어**: 공백이나 특수문자가 나올 때까지 문자를 읽고, 예약어 맵과 비교하여 판별 (`ReadIdentifier`)
    - **정수인지 판별**: 숫자가 끝날 때까지 읽고 저장 (`ReadNumber`)
    - **연산자 판별**: `==`, `!=` 등의 연산자 처리
    - **공백 제거**
3. **REPL 구현 후 테스트 수행**

---

# 2. 파싱 (Parsing)

## 정의

- **토큰 → 추상구문트리(AST)** 로 변환하는 과정
- 프로그램의 문법을 분석하고, 구조적인 표현(AST)으로 변환하는 역할

## 파서 제너레이터

- 문법을 입력하면 자동으로 파서를 생성해주는 소프트웨어

## Let문 파싱 (`let <identifier> = <expression>`)

- 필드: **식별자(Identifier)**, **표현식(Expression)**, **let 토큰**

### `LetStatement` 구조체

```
    type LetStatement struct {
        Token token.Token // 'let' 토큰
        Name  *Identifier // 식별자 (변수명)
        Value Expression  // 값 (표현식)
    }
```

✅ **각 필드의 값 예시:**

| 필드명 | 값 (예시) | 설명 |
| --- | --- | --- |
| `Token` | `{Type: token.LET, Literal: "let"}` | `let` 토큰을 저장 |
| `Name` | `&Identifier{Token: {Type: token.IDENT, Literal: "x"}, Value: "x"}` | 변수명 `x`를 저장 |
| `Value` | `&IntegerLiteral{Token: {Type: token.INT, Literal: "5"}, Value: 5}` | 값 `5`를 저장 |

### `Identifier` 구조체

```
    type Identifier struct {
        Token token.Token // 식별자 토큰 (변수명)
        Value string      // 실제 변수명
    }
```

✅ **각 필드의 값 예시:**

| 필드명 | 값 (예시) | 설명 |
| --- | --- | --- |
| `Token` | `{Type: token.IDENT, Literal: "x"}` | 변수명 `x`의 토큰 |
| `Value` | `"x"` | 실제 변수명 |

## 파서 동작 과정

### `p.nextToken()`을 두 번 호출하는 이유

- **한 번만 호출하면 `curToken`이 빈 값이 되고, `peekToken`이 첫 번째 토큰이 됨**
- 두 번 호출해야 `curToken`과 `peekToken`이 올바르게 초기화됨

### 파서 프로그램의 진입점: `ParseProgram()`

1. **루트 노드(`Program`)와 `Statement` 슬라이스 초기화**
2. **현재 토큰이 `let`인지, `return`인지 판별**
    - `let`이면 `parseLetStatement()` 호출
    - `return`이면 `parseReturnStatement()` 호출

### `parseLetStatement()`

1. **현재 위치의 토큰을 이용해 `LetStatement` 노드 생성**
2. **다음 토큰이 식별자인지 판별** (아니면 오류 추가)
3. **Identifier 노드 생성**
4. **다음 토큰이 `=`인지 판별** (아니면 오류 추가)

## 전체적인 파싱 로직 순서

1. **`ParseProgram()`** → 루트 노드 및 `Statement` 슬라이스 초기화
2. **`parseStatement()`** → 현재 토큰이 `let`, `return` 등인지 판별하여 분기 처리
3. **`parseLetStatement()`** → `Token`, `Identifier`, `Expression`을 저장하고 문법 오류 처리
4. **`ParseProgram()`** → `Statement`를 슬라이스에 추가 후, 다음 `Statement` 파싱 수행

---

✅ **이렇게 구현하면 `let x = 5; let y = 10;` 같은 코드가 올바르게 AST로 변환됨!**
