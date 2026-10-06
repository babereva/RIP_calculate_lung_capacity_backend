# Лабораторная 3. REST API

В рамках лабораторной работы реализован REST API для работы с категориями пациентов, лайками и пользователями.


## Методы API

| Метод | URL | Назначение |
| --- | --- | --- |
| GET | `/api/patient-categories?ageMin=20&ageMax=60` | Список опубликованных категорий с фильтрацией по возрасту |
| GET | `/api/patient-categories/feed` | Лента, первая опубликованная категория |
| GET | `/api/patient-categories/feed/:id?next=true` | Лента по идентификатору, с `next=true` следующая запись |
| GET | `/api/patient-categories/draft` | Черновик текущего пользователя |
| POST | `/api/patient-categories` | Создание черновика и загрузка фото и видео |
| PUT | `/api/patient-categories/:id/publish` | Заполнение полей и публикация черновика |
| DELETE | `/api/patient-categories/:id` | Логическое удаление категории создателем |
| POST | `/api/patient-categories/:id/like` | Установка или снятие лайка, тело `{ "value": 1 }` или `{ "value": 0 }` |
| POST | `/api/users/register` | Регистрация пользователя |
| POST | `/api/users/login` | Заглушка авторизации для лабораторной работы 4 |
| POST | `/api/users/logout` | Заглушка деавторизации для лабораторной работы 4 |

## Таблицы базы данных

### patient_category_users

`patient_category_user_id` — первичный ключ, `patient_category_username`, `patient_category_password`.

### patient_categories

`patient_category_id` — первичный ключ, `patient_category_title`, `patient_category_description`, `patient_category_status`, `patient_category_image_url`, `patient_category_video_url`, `patient_category_age`, `patient_category_height`, `patient_category_created_at`, `patient_category_published_at`, `patient_category_creator_id` — внешний ключ на `patient_category_users`.

### patient_category_likes

`patient_category_like_id` — первичный ключ, `patient_category_user_id` и `patient_category_id` — внешние ключи на пользователя и категорию, пара уникальна.
