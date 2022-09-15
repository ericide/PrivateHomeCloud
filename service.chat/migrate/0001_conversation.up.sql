CREATE TABLE `chat`.`conversation`  (
    `id` char(128)  NOT NULL,
    `type`   char(32)  NOT NULL,
    `chat_id` char(192)  NOT NULL,
    `owner_id` char(128)  NOT NULL,
    `oppo_id` char(36) NOT NULL DEFAULT '',
    `name` char(255) NOT NULL DEFAULT '',
    `last_read_time` bigint(20) NOT NULL DEFAULT 0,
    `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    PRIMARY KEY (`id`) USING BTREE
);