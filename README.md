# Objetivo
Objetivo: Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema deverá implementar OTEL(Open Telemetry) e Zipkin.

## Rodando a aplicação:

Há um arquivo docker-compose.yml no repositório, só é necessário rodar docker compose up.

Endpoints expostos:

Service A: http://localhost:8082/weather/{postalCode} (Vai fazer a requisição pro service B que vai chamar a viacep api)

Service B: http://localhost:8081/weather/{postalCode} (Chama a viacep api)

Zipkin: http://localhost:9411/zipkin/ (UI pra analisar os spans e traces das requisições.)

