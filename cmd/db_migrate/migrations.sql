-- +migrate Up

CREATE DATABASE IF NOT EXISTS `bsc` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `bsc`;


DROP TABLE IF EXISTS `blocks`;
CREATE TABLE `blocks` (
                          `number` int(10) UNSIGNED NOT NULL,
                          `hash` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
                          `time` int(11) NOT NULL,
                          `parent_hash` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
                          `transactions` mediumtext COLLATE utf8mb4_unicode_ci NOT NULL,
                          `confirmed` tinyint(4) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


ALTER TABLE `blocks`
    ADD UNIQUE KEY `number` (`number`);

DROP TABLE IF EXISTS `transactions`;
CREATE TABLE `transactions` (
                                `hash` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
                                `block_hash` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
                                `from_addr` varchar(42) COLLATE utf8mb4_unicode_ci NOT NULL,
                                `to_addr` varchar(42) COLLATE utf8mb4_unicode_ci NOT NULL,
                                `nonce` mediumint(8) UNSIGNED NOT NULL,
                                `data` text COLLATE utf8mb4_unicode_ci NOT NULL,
                                `value` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
                                `logs` text COLLATE utf8mb4_unicode_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `transactions`
    ADD UNIQUE KEY `hash` (`hash`);
COMMIT;

-- +migrate Down
DROP TABLE `blocks`;
DROP TABLE `transactions`;