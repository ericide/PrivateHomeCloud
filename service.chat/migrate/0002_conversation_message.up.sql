CREATE TABLE `chat`.`conversation_message`  (
        `id` char(128)  NOT NULL,
        `chat_id` char(192)  NOT NULL,
        `type`   char(32)  NOT NULL,
        `sender_id` char(128)  NOT NULL,
        `content` text NOT NULL,
        `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
        PRIMARY KEY (`id`) USING BTREE
);