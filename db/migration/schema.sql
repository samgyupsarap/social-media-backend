CREATE TABLE `comments` (
  `comment_uuid` varchar(36) NOT NULL,
  `comment_content` varchar(1000) DEFAULT NULL,
  `user_uuid` varchar(36) DEFAULT NULL,
  `post_uuid` varchar(36) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `commentscol` varchar(45) DEFAULT 'CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP',
  PRIMARY KEY (`comment_uuid`),
  UNIQUE KEY `comment_uuid_UNIQUE` (`comment_uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



CREATE TABLE `posts` (
  `post_uuid` varchar(36) NOT NULL,
  `post_content` longtext,
  `post_tags` varchar(200) DEFAULT NULL,
  `user_uuid` varchar(45) DEFAULT NULL,
  `likes` int DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` varchar(45) DEFAULT 'CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP',
  PRIMARY KEY (`post_uuid`),
  UNIQUE KEY `post_uuid_UNIQUE` (`post_uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



CREATE TABLE `users` (
  `user_uuid` varchar(36) NOT NULL,
  `full_name` varchar(200) DEFAULT NULL,
  `email` varchar(200) DEFAULT NULL,
  `user_name` varchar(100) DEFAULT NULL,
  `password` varchar(45) DEFAULT NULL,
  `profile_picture` blob,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_uuid`),
  UNIQUE KEY `user_uuid_UNIQUE` (`user_uuid`),
  UNIQUE KEY `email_UNIQUE` (`email`),
  UNIQUE KEY `user_name_UNIQUE` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

