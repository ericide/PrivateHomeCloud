CREATE TABLE `chat`.`user_login_record`  (
                                    `id` char(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                                    `user_id` char(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                                    `device` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                                    `device_name` char(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                                    `push_token` char(192) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                                    `invalid` int NOT NULL DEFAULT 0,
                                    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;