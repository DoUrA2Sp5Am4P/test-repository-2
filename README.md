# Для проверяющих
Для более полного понимания ваших оценок, если вам будет несложно, прошу оставить комментарий в [Issues](https://github.com/DoUrA2Sp5Am4P/test-repository-2/issues)
# Распределенный вычислитель арифметических выражений
Представим, что такие простые операции, как сложение, вычитание, умножение, деление занимают много времени.
Поэтому я создал приложение на Go, которое производит вычисление выражений параллельно.
# Как установить?
```
mkdir Distributed_arithmetic_expression_evaluator
cd Distributed_arithmetic_expression_evaluator
git init
git pull https://github.com/DoUrA2Sp5Am4P/test-repository-2.git
cd server
go build .
cd ..
cd client
go build .
Запустите ./server/server.exe
Затем запустите ./client/client.exe
```
**Примечание: чтобы корректно завершить программу используйте специальный пункт в меню. Это важно, так если в базу данных ещё не сохранились данные (а они сохраняются через промежутки времени localhost:8080/page7.html), то после перезагрузки они не восстановятся!**
![CloseProgram](https://github.com/DoUrA2Sp5Am4P/test-repository-png/blob/main/closeProgram.png?raw=true)
# Как протестировать?
```
cd Distributed_arithmetic_expression_evaluator
cd server
go test -v -run All
cd ..
cd client
go test -v -run All
```
# Как использовать?
## Не используйте 127.0.0.1:8080, используйте localhost:8080
```
0. По желанию отредактируйте файл client/CONCURRENCY.go, чтобы изменить количество воркеров (По умолчанию 5)
1. Перейдите на localhost:8080
Убедитесь, что на localhost:8080 ничего не запущено!
2. Наслаждайтесь приложением. Документация по использованию каждой страницы размещена на wiki
```

# [Wiki](https://github.com/DoUrA2Sp5Am4P/test-repository-2/wiki) <-- активная ссылка на вики
### На wiki размещены примеры запросов, **подробно описанные тестовые сценарии**
