# loglint - Линтер для лог-сообщений в го

## Встроенные правила
- Лог-сообщения должны начинаться со строчной буквы
- Лог-сообщения должны быть только на английском языке
- Лог-сообщения не должны содержать спецсимволы или эмодзи
- Лог-сообщения не должны содержать потенциально чувствительные данные

## Поддерживаемые логгеры
- `log/slog` (Info, Error, Warn, Debug, Trace, Fatal, Panic)
- `log` (Print, Println, Printf, Fatal, Panic)

## Установка
1. Установите golangci-lint:
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```
2. Уcтановите плагин линтера
   ```bash
   git clone https://github.com/yxatabl/go-linter.git
   cd go-linter
   make plugin
   ```
3. Настройте плагин для golangci-lint
   Создайте `.golangci.yml` в вашем проекте
   ```bash
   linters-settings:
    custom:
      loglint:
        path: /path/to/your/go/bin/loglint.so
        description: Checks log messages for style violations
        original-url: github.com/yxatabl/go-linter/m/pkg/loglint
  
   linters:
     enable:
       - loglint
   ```

## Использование
```bash
golangci-lint run --enable=loglint
```
