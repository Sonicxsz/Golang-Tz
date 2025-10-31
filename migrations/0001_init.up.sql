-- Создание таблицы подписок
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

                             -- Проверка что end_date >= start_date, если указан
                             CONSTRAINT valid_date_range CHECK (end_date IS NULL OR end_date >= start_date)
    );

-- Индексы для оптимизации запросов
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_service_name ON subscriptions(service_name);
CREATE INDEX idx_subscriptions_start_date ON subscriptions(start_date);
CREATE INDEX idx_subscriptions_end_date ON subscriptions(end_date);

-- Составной индекс для фильтрации по пользователю и сервису
CREATE INDEX idx_subscriptions_user_service ON subscriptions(user_id, service_name);

-- Индекс для запросов по периоду
CREATE INDEX idx_subscriptions_date_range ON subscriptions(start_date, end_date);

-- Комментарии к таблице и колонкам
COMMENT ON TABLE subscriptions IS 'Таблица для хранения информации о подписках пользователей';
COMMENT ON COLUMN subscriptions.id IS 'Уникальный идентификатор подписки';
COMMENT ON COLUMN subscriptions.service_name IS 'Название сервиса подписки';
COMMENT ON COLUMN subscriptions.price IS 'Стоимость месячной подписки в рублях (целое число)';
COMMENT ON COLUMN subscriptions.user_id IS 'UUID пользователя';
COMMENT ON COLUMN subscriptions.start_date IS 'Дата начала подписки (месяц и год)';
COMMENT ON COLUMN subscriptions.end_date IS 'Дата окончания подписки (опционально)';
COMMENT ON COLUMN subscriptions.created_at IS 'Дата создания записи';
COMMENT ON COLUMN subscriptions.updated_at IS 'Дата последнего обновления записи';