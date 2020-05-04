# Desafio Golang
- Essa é uma API que cadastra super heróis ou vilões utilizando como fonte de dados a superHeroAPI (https://superheroapi.com/)

## Pré-requisitos
Go 1.14 e postgres

## Como rodar?
- Crie uma cópia do arquivo .env-example com o nome de .env, preenchendo as configurações de banco de dados, além de seu token de acesso que pode ser obtido no site da superHeroAPI
- Na raiz do projeto, execute o seguinte comando: 
```
go run main.go
```

## Funcionalidades
- Cadastrar um super: Mande uma requisição POST para ```localhost:8080/supers``` com o seguinte body ```{"name":"nome_super"}```

- Listar todos os supers cadastrados: Requisição GET para ```localhost:8080/supers```

- Listar apenas os heróis ou vilões: Requisição GET para ```localhost:8080/supers/alignment/{good/bad}``` utilizando a chave *good* para heróis ou *bad* para vilões

- Buscar por nome: Requisição GET para ```localhost:8080/supers/search/{nome_super}```

- Buscar por uuid: Requisição GET para ```localhost:8080/supers/{uuid_super}```

- Remover o super: Requisição DELETE para ```localhost:8080/supers/{uuid_super}```

## Testes
- Para executar os testes dos modelos, execute o comando ```go test -v``` a partir da pasta tests > model
- Para executar os testes dos controllers, execute o comando ```go test -v``` a partir da pasta tests > controller

## Dependências
- Mux (http://github.com/gorilla/mux)
- Gorm (http://github.com/jinzhu/gorm)
- GoDotEnv (http://github.com/joho/godotenv)
- Assert (http://gopkg.in/go-playground/assert.v1)
