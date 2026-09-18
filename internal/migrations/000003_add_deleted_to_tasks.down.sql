ALTER TABLE tasks DROP COLUMN deleted;
/* 
ALTER TABLE tasks — изменить таблицу tasks.

ADD COLUMN deleted — добавить колонку deleted.

BOOLEAN — тип (аналог bool в Go).

NOT NULL — колонка не может быть NULL. Это важно! Иначе будут три состояния (NULL, true, false), и в Go придётся использовать *bool.

DEFAULT FALSE — значение по умолчанию для всех существующих записей — false.
 */