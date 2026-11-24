# Avito Kitchen

MVP сервис для заказа еды из ресторанов

Проект состоит из двух Go сервисов

* **Avito Kitchen** Основной сервис. В этом сервисе люди регистрируются, создают заказы и узнают о ресторанах и их продукции
* **Restaurant Service** Отдельный сервис. Пример ресторанного API. Хранит меню, продукты, заказы. Отправляет изменения статуса в Avito Kitchen

### Архитектура

Сценарий:

1. Пользователь авторизуется на сервисе Avito Kitchen
2. Пользователь получает список ресторанов
3. Пользователь получает меню ресторана
4. Пользователь создаёт заказ
5. Avito Kitchen пытается создать заказ и проверяет их наличие в меню у Restaurant API
6. Заказ сохраняется в базе Avito Kitchen
7. Restaurant Service создаёт внутренний заказ и возвращает внутренний id
8. Restaurant Service изменяет статус заказа и отправляет изменения в Avito Kitchen

Для связи заказов используются два идентификатора:

orders.id - идентификатор заказа в Avito Kitchen
restaurant_order_id - идентификатор заказа в Restaurant Service

## Схема базы данных

### Avito Kitchen DB

**users**

* `id`
* `name`
* `email`
* `password`
* `created_at`
* `updated_at`

**restaurants**

* `id`
* `name`
* `base_url`
* `api_key`

**orders**

* `id`
* `restaurant_id`
* `user_id`
* `restaurant_order_id`
* `status`
* `total_price`
* `created_at`

Связи между таблицами

```text
users -> orders one-to-many
restaurants -> orders one-to-many
orders -> order_items one-to-many
```

**order_items**

* `order_id`
* `restaurant_product_id`
* `name`
* `price`
* `quantity`

Название и цена товара сохраняются в заказе как snapshot, поэтому изменение меню ресторана не изменяет уже созданный заказ

### Restaurant DB

**products**

* `id`
* `name`
* `price`
* `is_available`
* `created_at`
* `updated_at`

**orders**

* `id`
* `kitchen_order_id`
* `status`
* `created_at`
* `updated_at`

`kitchen_order_id` связывает заказ ресторана с заказом в Avito Kitchen.

**order_items**

* `order_id`
* `product_id`
* `quantity`
* `price`

Связи:

```text
products -> order_items one-to-many
orders   -> order_items one-to-many
```

## Запуск

```bash
docker compose up --build
```

После чего поднимаются два сервиса, доступные по:<br>
http://loclhost:3001 - Avito Kitchen <br>
http://localhost:3000 - Restaurant Service

## Документация

```text
docs/
├── architecture
│   ├── c4-level-2.puml
│   └── c4-level-3.puml
├── cjm
│   ├── restaurant-order.puml
│   └── user-order.puml
└── openapi
    ├── avito_kitchen.yaml
    └── restaurant.yaml
```