# image-preview
Загрузка и обработка изображений.

## Требования:
- возможность загружать изображения по заданному URL
- возможность загружать несколько изображений в обном запросе
- создание квадратного превью изображения размером 100px на 100px

## Как запустить
Собраные исходники находятся на [https://github.com/dbond762/image-preview/releases/tag/v1.0](https://github.com/dbond762/image-preview/releases/tag/v1.0)

Для того, что бы собрать, выполните
```
go build
```

## Пример использования
[http://localhost:8000/preview?url=https://golang.org/lib/godoc/images/footer-gopher.jpg&url=https://golang.org/lib/godoc/images/home-gopher.png&url=https://i.pinimg.com/originals/bc/75/22/bc75225ef044d29d1f2d1c051d9b8063.gif](http://localhost:8000/preview?url=https://golang.org/lib/godoc/images/footer-gopher.jpg&url=https://golang.org/lib/godoc/images/home-gopher.png&url=https://i.pinimg.com/originals/bc/75/22/bc75225ef044d29d1f2d1c051d9b8063.gif)
