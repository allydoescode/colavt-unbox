CREATE TABLE IF NOT EXISTS `items` (
    `id` INTEGER PRIMARY KEY,
    `name` TEXT NOT NULL,
    `desc` TEXT NOT NULL,
    `image` TEXT NOT NULL,
    `rarity` TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS `inventory` (
    `id` INTEGER PRIMARY KEY,
    `item_id` INTEGER NOT NULL,
    `owner` TEXT NOT NULL,
    `wear` TEXT NOT NULL,
    `float` REAL NOT NULL,
    FOREIGN KEY(`item_id`) REFERENCES `items`(`id`)
);