Olá dev, tudo bem?

Neste desafio vamos aplicar o que aprendemos sobre webserver http, contextos,
banco de dados e manipulação de arquivos com Go.

Você precisará nos entregar dois sistemas em Go:
- client.go
- server.go

Os requisitos para cumprir este desafio são:


## Server (server.go)
- Utilizar API com câmbio de Dólar e Real: https://economia.awesomeapi.com.br/json/last/USD-BRL timeout máximo de 200ms
- Retornar no formato JSON o resultado para o cliente.
- Registrar no banco de dados SQLite cada cotação recebida com timeout de 10ms
- Endpoint a ser exposto: 8080/cotacao

Utilizar context

## Client  
- Realizar requisição HTTP para o server solicitando a cotação do dólar com timeout de 300ms
- Precisará receber apenas o valor atual do câmbio "bid"
- Salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}

Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.
