

**Исполнитель:** Садыков Озоджон Бахтиёрович  
**Почта:** ozodbaht8@gmail.com

**balanceapp - API для работы с балансом пользователей**

Оглавление:
-----------
 
**[ЧАСТЬ 1: Задание](#tz)**

**[ЧАСТЬ 2: API](#api)**  
&emsp;**[2.1 Метод начисления средств на баланс](#m1)**  
&emsp;**[2.2 Метод списания средств с баланса](#m2)**  
&emsp;**[2.3 Метод перевода средств от пользователя к пользователю](#m3)**  
&emsp;**[2.4 Метод получения текущего баланса пользователя](#m4)**  
&emsp;**[2.5 Метод получения списка транзакций](#m5)**  


<a name="tz">ЧАСТЬ 1: Задание</a>
---------------------------------
**Задача:**

Необходимо реализовать микросервис для работы с балансом пользователей (зачисление средств, списание средств, перевод средств от пользователя к пользователю, а также метод получения баланса пользователя). Сервис должен предоставлять HTTP API и принимать/отдавать запросы/ответы в формате JSON. 

**Сценарии использования:**

Далее описаны несколько упрощенных кейсов приближенных к реальности.
1. Сервис биллинга с помощью внешних мерчантов (аля через visa/mastercard) обработал зачисление денег на наш счет. Теперь биллингу нужно добавить эти деньги на баланс пользователя. 
2. Пользователь хочет купить у нас какую-то услугу. Для этого у нас есть специальный сервис управления услугами, который перед применением услуги проверяет баланс и потом списывает необходимую сумму. 
3. В ближайшем будущем планируется дать пользователям возможность перечислять деньги друг-другу внутри нашей платформы. Мы решили заранее предусмотреть такую возможность и заложить ее в архитектуру нашего сервиса. 

**Требования к коду:**

1. Язык разработки: Go/PHP. Мы готовы рассматривать решения на любом языке, но приоритетом для нас являются именно эти.
2. Фреймворки и библиотеки можно использовать любые
3. Реляционная СУБД: MySQL или PostgreSQL
4. Весь код должен быть выложен на Github с Readme файлом с инструкцией по запуску и примерами запросов/ответов (можно просто описать в Readme методы, можно через Postman, можно в Readme curl запросы скопировать, вы поняли идею...)
5. Если есть потребность, можно подключить кеши(Redis) и/или очереди(RabbitMQ, Kafka)
6. При возникновении вопросов по ТЗ оставляем принятие решения за кандидатом (в таком случае Readme файле к проекту должен быть указан список вопросов с которыми кандидат столкнулся и каким образом он их решил)
7. Разработка интерфейса в браузере НЕ ТРЕБУЕТСЯ. Взаимодействие с АПИ предполагается посредством запросов из кода другого сервиса. Для тестирования можно использовать любой удобный инструмент. Например: в терминале через curl или Postman.

**Будет плюсом:**

1. Использование docker и docker-compose для поднятия и развертывания dev-среды.
2. Методы АПИ возвращают человеко-читабельные описания ошибок и соответвующие статус коды при их возникновении.
3. Все реализовано на GO, все-же мы собеседуем на GO разработчика. HINT: На собеседовании так или иначе будут вопросы по Go. Кто прочитал, тот молодец :)
4. Написаны unit/интеграционные тесты.

**Основное задание (минимум):**

Метод начисления средств на баланс. Принимает id пользователя и сколько средств зачислить.

Метод списания средств с баланса. Принимает id пользователя и сколько средств списать. 

Метод перевода средств от пользователя к пользователю. Принимает id пользователя с которого нужно списать средства, id пользователя которому должны зачислить средства, а также сумму.

Метод получения текущего баланса пользователя. Принимает id пользователя. Баланс всегда в рублях.

**Детали по заданию:**

1. Методы начисления и списания можно объединить в один, если это позволяет общая архитектура.
2. По умолчанию сервис не содержит в себе никаких данных о балансах (пустая табличка в БД). Данные о балансе появляются при первом зачислении денег. 
3. Валидацию данных и обработку ошибок оставляем на усмотрение кандидата. 
4. Список полей к методам не фиксированный. Перечислен лишь необходимый минимум. В рамках выполнения доп. заданий возможны дополнительные поля.
5. Механизм миграции не нужен. Достаточно предоставить конечный SQL файл с созданием всех необходимых таблиц в БД. 
6. Баланс пользователя - очень важные данные в которых недопустимы ошибки (фактически мы работаем тут с реальными деньгами). Необходимо всегда держать баланс в актуальном состоянии и не допускать ситуаций когда баланс может уйти в минус. 
7. Валюта баланса по умолчанию всегда рубли.

**Дополнительные задания**

Далее перечислены доп. задания. Они не являются обязательными, но их выполнение даст существенный плюс перед другими кандидатами. 

*Доп. задание 1:*

Эффективные менеджеры захотели добавить в наши приложения товары и услуги в различных от рубля валютах. Необходима возможность вывода баланса пользователя в отличной от рубля валюте.

Задача: добавить к методу получения баланса доп. параметр. Пример: ?currency=USD. 
Если этот параметр присутствует, то мы должны конвертировать баланс пользователя с рубля на указанную валюту. Данные по текущему курсу валют можно взять отсюда https://exchangeratesapi.io/ или из любого другого открытого источника. 

Примечание: напоминаем, что базовая валюта которая хранится на балансе у нас всегда рубль. В рамках этой задачи конвертация всегда происходит с базовой валюты.

*Доп. задание 2:*

Пользователи жалуются, что не понимают за что были списаны (или зачислены) средства. 

Задача: необходимо предоставить метод получения списка транзакций с комментариями откуда и зачем были начислены/списаны средства с баланса. Необходимо предусмотреть пагинацию и сортировку по сумме и дате. 


<a name="api">ЧАСТЬ 2: API</a>
------------

**ФОРМАТЫ ВЫХОДНЫХ ДАННЫХ:** `JSON`  


### <a name="m1">2.1 Метод начисления средств на баланс</a>

**URL:http://localhost:8080/balances/income**  


**METHOD: POST**

**Request body:**
```javascript
{
  "toId": 1, // int, идентификатор баланса
  "amount": "300", // decimal, сумма начисления 
  "reason": "For test" // string, причина начисления 
}
```
**Статус-коды:**  
`200` - успешно  
`400` - неверные URL параметры


### <a name="m2">2.2 Метод списания средств с баланса</a>

**URL:http://localhost:8080/balances/outcome**  


**METHOD: POST**

**Request body:**
```javascript
{
  "fromId": 1, // int, идентификатор баланса
  "amount": "300", // decimal, сумма списания 
  "reason": "For test" // string, причина списания 
}
```
**Статус-коды:**  
`200` - успешно  
`400` - неверные URL параметры


### <a name="m3">2.3 Метод перевода средств от пользователя к пользователю</a>

**URL:http://localhost:8080/balances/transfer**  


**METHOD: POST**

**Request body:**
```javascript
{
  "fromId": 1, // int, идентификатор баланса, с которого переводим
  "toId" : 2, // int, идентификатор баланса, на который переводим
  "amount": "300", // decimal, сумма списания 
  "reason": "For test" // string, причина перевода 
}
```
**Статус-коды:**  
`200` - успешно  
`400` - неверные URL параметры

 
### <a name="m4">2.4 Метод получения текущего баланса пользователя</a>

**URL:http://localhost:8080/balances/{id}?currency=USD**  

**METHOD: GET**

```javascript
  id, // int, идентификатор баланса, обязательный параметр в URL
  currency, // string, значение валюты, в которой хотим получить баланс, необязательный параметр в URL
```


**Статус-коды:**  
`200` - успешно  
`400` - неверные URL параметры


### <a name="m5">2.5 Метод получения списка транзакций</a>

**URL:http://localhost:8080/balances/history/{id}?order_by=amount&limit=5&offset=2**  

**METHOD: GET**

```javascript
  id, // int, идентификатор баланса, обязательный параметр в URL
  order_by, // string, параметр сортировки из значений: amount - по сумме, date - по времени, необязательный параметр в URL
  limit, // int, количество транзакций, которое хотим получить, необязательный параметр в URL
  offset, // int, количество транзакций, которое хотим пропустить, необязательный параметр в URL
```

**Response body:**
```javascript
{
  "id": 3,// int, идентификатор транзакции
  "balance_id": 1, // int, идентификатор баланса, для которого произведена транзакция 
  "from_id": 2, // int, идентификатор баланса, с которым связана транзакция, если был сделан перевод от одного к другому пользователю, либо равен нулю при использовании методов начисления и списания  
  "amount": "100", // decimal, сумма списания или начисления 
  "reason": "For something", // string, причина транзакции 
  "type": "outcome", // string, тип транзакции, outcome - списание, income - начисление
  "date": "2020-09-28 17:01:55" //time, время совершения транзакции
}

**Статус-коды:**  
`200` - успешно  
`400` - неверные URL параметры

