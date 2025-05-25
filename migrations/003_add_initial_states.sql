-- Добавление начальных состояний книг
INSERT INTO states (id, name, created_at, updated_at) VALUES
(1, 'available', NOW(), NOW()),  -- Доступна для обмена
(2, 'trading', NOW(), NOW()),    -- В процессе обмена
(3, 'traded', NOW(), NOW());     -- Обменена 