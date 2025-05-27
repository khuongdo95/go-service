-- Create "ip_access_controls" table
CREATE TABLE `ip_access_controls` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL,
  `created_by` varchar(255) NOT NULL,
  `updated_by` varchar(255) NULL,
  `ip_address` varchar(255) NOT NULL,
  `tcp_stage` bool NOT NULL DEFAULT false,
  `tcp_live` bool NOT NULL DEFAULT false,
  `rabbitmq_stage` bool NOT NULL DEFAULT false,
  `rabbitmq_live` bool NOT NULL DEFAULT false,
  `translations_rmq` bool NOT NULL DEFAULT false,
  `translations_tcp` bool NOT NULL DEFAULT false,
  `active` bool NOT NULL DEFAULT true,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "users" table
CREATE TABLE `users` (
  `id` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL,
  `deleted_at` timestamp NULL,
  `deleted_by` varchar(255) NULL,
  `created_by` varchar(255) NOT NULL,
  `updated_by` varchar(255) NULL,
  `name` varchar(255) NULL,
  `email` varchar(255) NULL,
  `tenant_id` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "identities" table
CREATE TABLE `identities` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL,
  `username` varchar(255) NULL,
  `password` varchar(255) NULL,
  `user_id` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `identity_username` (`username`),
  CONSTRAINT `identities_users_identity` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "user_ip_white_lists" table
CREATE TABLE `user_ip_white_lists` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL,
  `ip_address` varchar(255) NOT NULL,
  `user_id` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `user_ip_white_lists_users_ip_white_list` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
