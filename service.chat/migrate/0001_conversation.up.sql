CREATE TABLE `chat`.`conversation`  (
    `id` char(128)  NOT NULL,
    `type`   char(32)  NOT NULL,
    `chat_id` char(192)  NOT NULL,
    `owner_id` char(128)  NOT NULL,
    `name` char(255) NOT NULL DEFAULT '',
    `last_read_time` datetime  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `create_time` datetime  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE
);