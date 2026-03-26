-- Запрос для получения списка активных устройств с количеством связанных configs.
SELECT
    d.id,
    d.name,
    COUNT(c.id) AS config_count
FROM
    devices AS d
    LEFT JOIN configs AS c ON c.device_id = d.id
WHERE
    d.is_active = TRUE
GROUP BY
    d.id,
    d.name
ORDER BY
    d.name;

-- Запрос для получения последних N записей из logs для конкретного устройства.
-- N = 10
SELECT
    *
FROM
    logs
WHERE
    device_id = 123
ORDER BY
    created_at DESC
LIMIT
    10;


-- Пример "тяжёлого" запроса
-- Найти устройства без конфигураций, но с логами за последнюю неделю (более 50 логов).
SELECT 
    d.id,
    d.hostname,
    d.ip,
    d.location,
    d.created_at,
    COUNT(l.id) as logs_count
FROM devices AS d
LEFT JOIN configs AS c ON c.device_id = d.id
INNER JOIN logs AS l ON l.device_id = d.id
WHERE d.is_deleted = FALSE
    AND c.id IS NULL
    AND l.created_at >= NOW() - INTERVAL '7 days'
GROUP BY d.id, d.hostname, d.ip, d.location, d.created_at
HAVING COUNT(l.id) > 50
ORDER BY logs_count DESC;

-- Индекс для поиска отсутствующих конфигураций
CREATE INDEX idx_configs_device ON configs(device_id);
CREATE INDEX idx_logs_device_created ON logs(device_id, created_at);
