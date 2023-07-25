# Estágio de construção - usando a imagem base golang para compilar o código
FROM golang:1.20 AS build

# Definindo o diretório de trabalho
WORKDIR /app

# Copiando os arquivos do projeto
COPY . .

# Baixando as dependências do projeto
RUN go mod tidy
RUN go mod download

# Compilando o arquivo Go para um binário executável
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o namespace-cleaner main.go

# Estágio de execução - usando uma imagem base alpine para a imagem final (para manter o tamanho da imagem pequeno)
FROM alpine:3.13

# Instalando ca-certificates
RUN apk --no-cache add ca-certificates

# Copiando o binário compilado do estágio de construção para a imagem final
COPY --from=build /app/namespace-cleaner /namespace-cleaner

# Configurando o comando padrão para executar quando a imagem for iniciada
CMD ["/namespace-cleaner"]
