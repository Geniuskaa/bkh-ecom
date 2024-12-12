-- +goose Up
DO $$
    BEGIN
        -- Генерируем 50 случайных баннеров
        FOR i IN 1..50 LOOP
                INSERT INTO banners (name)
                VALUES (concat('Banner_', i));  -- Название баннера будет 'Banner_1', 'Banner_2', ..., 'Banner_50'
            END LOOP;
    END $$;
