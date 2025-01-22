# Monkey_Interpreter

## 1. 렉싱

- 정의 :  “구문 → 토큰 → 추상구문트리” 에서 구문→트리 로 변환하는 과정
- 렉서 구조체 : 현재 읽고있는 문자의 포지션, 다음 읽을 문자 포지션, 현재 가르키고 있는 문자, input
- 토큰 구조체 : 타입과 리터럴이 존재하고 리터럴은 string으로 설정(모든 타입을 받을 수 있고 복잡도를 낮추기 위함)
- 렉서를 먼저 생성한 뒤 생성됨과 동시에 readch 메서드를 호출해 구문의 첫번째 문자를 읽고 구조체의 멤버변수를 초기화한다.
- nextToken 메서드를 호출해서 식별자인지 예약어인지 정수인지 switch안에 if문으로 각각을 구별하고 토큰으로 분리한다.
- 식별자인지 예약어인지 → 특수문자(공백)이 나올때까지 문자를 읽으면서 문자인지 판별한뒤(IsLetter) 해당 자리까지 글자를 변수에 담은 뒤(ReadIdentifier) 해당 글자가 예약어인지 이미 존재하는 맵에 비교한 후 판별한다. (ReadIdentifier)
- 정수인지 → 정수가 끝날때까지 문자를 읽으면서 정수인지 판별한 뒤(isDigit) 해당 자리까지 정수를 변수에 담는다. (ReadNumber)
- 공백처리함수, 스위치문에서 EQ(==)와 NOT_EQ(≠) 처리, {if, return, else, false, function} 등등 다양한 예약어 지원, REPL 구현 후 테스트가 즉각적으로 이루어지는지

## 2. 파싱

- 정의 : 토큰 → 추상구문트리로 변환시켜주는 프로그램
- 파서 제너레이터 → 문법을 입력하면 자동으로 파서를 생성해주는 소프트웨어
- Let문 파싱 → let <identifier> = <expression> 으로 구성되어 있음.
    - 그러므로 필드는 identifier, expression이 필요하고 렉서에서 개발한 token.Let도 필드로 가져야 한
    
    ```go
    type LetStatement struct{
    	Token token.Token
    	Name *Identifier
    	value Expression
    }
    ```
    
- 파서 → p.nextToken() 를 두번 호출하는 이유는 한번만 호출하면 curToken이 빈값이 들어가고 peekToken이 첫번째 토큰이 들어가기 때문에 한번 더 호출해서 제대로 값을 설정해준다.
- 파서 프로그램의 진입점 : ParseProgram()
    - ParseProgram은 루트노드와 statment의 슬라이스를 초기화 시킨 뒤 토큰을 순차적으로 가져와서 let문이라면 let파싱, return 이라면 return파싱을 수행한다.
- parseLetStatement() : let문 파싱
    - 현재 위치에 있는 토큰으로 LetStatement 노드를 만듦
    - 다음토큰이 식별자인지 판별 (아니면 error추가)
    - Identifier 노드를 생성
    - 다음토큰이 등호인지 판별 (아니면 error 추가)
