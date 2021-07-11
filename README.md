# urlshortener

Программа для создания короткой ссылки URL через POST запросы

Приложены 2 изображения работы - кодирование с записью в sqlite и декодирование.

Реализованы следующие методы:

1. На вход поступает длинная ссылка, возвращается сокращённая ссылка
Request:
POST /short {"url": "long-url-here"}
Response:
{"url": "short-url-here"}
2. На вход поступает сокращённая ссылка, возвращается полная ссылка
Request:
POST /long {"url": "short-url-here"}
Response:
{"url": "long-url-here"}


