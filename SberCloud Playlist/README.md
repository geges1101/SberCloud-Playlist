# Тестовое задание для поступления в GoCloudCamp


## 1. Вопросы для разогрева

1. Опишите самую интересную задачу в программировании, которую вам приходилось решать?

Улучшал существующие RESTful API, добавлял дополнительные функции и возможности обработки ошибок для заказов платформы

2. Расскажите о своем самом большом факапе? Что вы предприняли для решения проблемы?

Переписал алгоритм фильтрации заказов, при тестировании ошибок не было выявлено. В итоге код поломал на небольшое время сервис, пришлось откатывать.

3. Каковы ваши ожидания от участия в буткемпе?

Хочу получить хороший опыт разработки на Go, развиться как разработчик

---

## 2. Разработка музыкального плейлиста

### Часть 1. Разработка основного модуля работы с плейлистом

Требуется разработать модуль для обеспечения работы с плейлистом. Модуль должен обладать следующими возможностями:
- Play - начинает воспроизведение
- Pause - приостанавливает воспроизведение
- AddSong - добавляет в конец плейлиста песню
- Next воспроизвести след песню
- Prev воспроизвести предыдущую песню


### Часть 2: Построение API для музыкального плейлиста

Реализовать сервис, который позволит управлять музыкальным плейлистом. Доступ к сервису должен осуществляться с помощью API, который имеет возможность выполнять CRUD операции с песнями в плейлисте, а также воспроизводить, приостанавливать, переходить к следующему и предыдущему трекам. Конфигурация может храниться в любом источнике данных, будь то файл на диске, либо база данных. Для удобства интеграции с сервисом может быть реализована клиентская библиотека.
