CREATE TABLE `migrations` (
  `id`        INT(10) UNSIGNED                        NOT NULL AUTO_INCREMENT,
  `file_name` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci