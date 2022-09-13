CREATE TABLE `chat`.`user`  (
        `id` char(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
        `phone` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
        `password` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
        `name` char(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
        `avatar_url` char(255) NOT NULL DEFAULT '',
        PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;