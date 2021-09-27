<?php

namespace SwagTestEnvironment\Migration;

use Shopware\Core\Framework\Migration\MigrationStep;
use Doctrine\DBAL\Connection;

class Migration1575569953createTable extends MigrationStep
{
    public function getCreationTimestamp(): int
    {
        return 1575569953;
    }

    public function update(Connection $connection): void
    {
        $connection->executeStatement('CREATE TABLE `sw_mail_catcher` (
    `id` BINARY(16) NOT NULL,
    `sender` JSON NOT NULL,
    `receiver` JSON NOT NULL,
    `subject` VARCHAR(255) NOT NULL,
    `plainText` LONGTEXT NULL,
    `htmlText` LONGTEXT NULL,
    `eml` LONGTEXT NULL,
    `created_at` DATETIME(3) NOT NULL,
    `updated_at` DATETIME(3) NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `json.sw_mail_catcher.receiver` CHECK (JSON_VALID(`receiver`))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
');
    }

    public function updateDestructive(Connection $connection): void
    {
    }
}
