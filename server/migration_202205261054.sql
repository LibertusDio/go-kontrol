-- -------------------------------------------------------------
-- TablePlus 4.6.6(422)
--
-- https://tableplus.com/
--
-- Database: auth_db
-- Generation Time: 2022-05-26 10:54:04.8490
-- -------------------------------------------------------------


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


CREATE TABLE `object_policy_mesh` (
  `id` varchar(36) NOT NULL,
  `object_id` varchar(36) NOT NULL,
  `created_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  `policy_id` varchar(36) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `object_policy_mesh_UN` (`object_id`,`policy_id`),
  KEY `object_policy_mesh_object_id_IDX` (`object_id`,`policy_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `objects` (
  `id` varchar(36) NOT NULL,
  `created_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  `global_id` varchar(36) NOT NULL DEFAULT '',
  `external_id` varchar(36) NOT NULL,
  `service_id` varchar(36) NOT NULL,
  `status` varchar(10) NOT NULL DEFAULT '',
  `token` varchar(100) NOT NULL DEFAULT '',
  `expiry_date` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `objects_golbal_UN` (`external_id`,`service_id`),
  UNIQUE KEY `objects_UN` (`global_id`),
  KEY `objects_global_id_IDX` (`global_id`) USING BTREE,
  KEY `objects_external_id_IDX` (`external_id`) USING BTREE,
  KEY `objects_service_id_IDX` (`service_id`) USING BTREE,
  KEY `objects_status_IDX` (`status`) USING BTREE,
  KEY `objects_token_IDX` (`token`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `policies` (
  `id` varchar(36) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  `created_at` bigint(20) NOT NULL,
  `name` varchar(100) NOT NULL DEFAULT '',
  `service_id` varchar(36) NOT NULL,
  `status` varchar(10) NOT NULL DEFAULT '',
  `apply_from` bigint(20) NOT NULL DEFAULT '0',
  `apply_to` bigint(20) NOT NULL DEFAULT '2147483647',
  `permission` longtext NOT NULL,
  PRIMARY KEY (`id`),
  KEY `policies_service_id_IDX` (`service_id`) USING BTREE,
  KEY `policies_status_IDX` (`status`) USING BTREE,
  KEY `policies_apply_from_IDX` (`apply_from`) USING BTREE,
  KEY `policies_apply_to_IDX` (`apply_to`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `service_policy_mesh` (
  `id` varchar(36) NOT NULL,
  `service_id` varchar(36) NOT NULL,
  `created_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  `policy_id` varchar(36) NOT NULL,
  `type` varchar(10) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `service_policy_mesh_UN` (`service_id`,`policy_id`,`type`),
  KEY `service_policy_mesh_service_id_IDX` (`service_id`,`policy_id`,`type`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `services` (
  `id` varchar(36) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  `created_at` bigint(20) NOT NULL,
  `service_id` varchar(36) NOT NULL,
  `key` varchar(150) NOT NULL,
  `status` varchar(10) NOT NULL DEFAULT '',
  `expiry_date` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `services_service_id_IDX` (`service_id`) USING BTREE,
  KEY `services_status_IDX` (`status`) USING BTREE,
  KEY `services_expiry_date_IDX` (`expiry_date`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `objects` (`id`, `created_at`, `updated_at`, `global_id`, `external_id`, `service_id`, `status`, `token`, `expiry_date`) VALUES
('2e041cf1-9a73-4eb3-a940-76c6007b0eaf', 1653475139, 1653533999, 'cc9c471c-dc9f-11ec-9d64-0242ac120002', 'hrd_id_1', '88503398-db0f-11ec-9d64-0242ac120004', 'enable', '0PcdoH-jm0cQr5gRMnXW-RSNznvQ4I6Z1VjkrmSh2hI=', 1653476939),
('96d4c8f5-473b-474c-85b7-2f87aa2f391a', 1653475088, 1653534002, 'e59745be-dc9f-11ec-9d64-0242ac120002', 'idt_id_1', '88503398-db0f-11ec-9d64-0242ac120002', 'enable', 'hByzOcYFMx8zXEw7rBgL_WnVQR6w_zwKCSVsSJ9_vh4=', 1653476888),
('d1e49882-b07a-4fce-b4cb-d317733d266f', 1653362262, 1653534016, 'f31eadf8-dc9f-11ec-9d64-0242ac120002', 'external_id_1', '88503398-db0f-11ec-9d64-0242ac120002', 'enable', 'hVhYAtji7tzvlk3qXzkB4O_W2LxlZUcHxL8p2R7NHIE=', 1653456377),
('e1021392-38c1-4003-adf6-90582a7e528d', 1653474967, 1653534034, 'fc4831ec-dc9f-11ec-9d64-0242ac120002', 'adt_id_1', '88503398-db0f-11ec-9d64-0242ac120003', 'enable', 'KmkYs7MtAOq_TYSPD-ZgpLs_jlLFdT7G6Gc3YoyN7DQ=', 1653476767);

INSERT INTO `policies` (`id`, `updated_at`, `created_at`, `name`, `service_id`, `status`, `apply_from`, `apply_to`, `permission`) VALUES
('9941806f-e4a6-4360-aa65-03744aefc785', 1653363881, 1653363881, 'administrator', '88503398-db0f-11ec-9d64-0242ac120002', 'enable', 1653016960, 1684552960, '{\"delete_profile\":0,\"edit_profile\":1,\"view_profile\":1}');

INSERT INTO `service_policy_mesh` (`id`, `service_id`, `created_at`, `updated_at`, `policy_id`, `type`) VALUES
('bb456508-db13-11ec-9d64-0242ac120002', '88503398-db0f-11ec-9d64-0242ac120002', 1653363987, 1653363987, '9941806f-e4a6-4360-aa65-03744aefc785', 'default');

INSERT INTO `services` (`id`, `updated_at`, `created_at`, `service_id`, `key`, `status`, `expiry_date`) VALUES
('88503398-db0f-11ec-9d64-0242ac120002', 1653474132, 1653362097, 'idt', 'HWFxDt1CvZrbIQG8CeiRX3tNIuxitU31I4DmkYOA_KA=', 'enable', 1684552960),
('88503398-db0f-11ec-9d64-0242ac120003', 1653474135, 1653468025, 'adt', 'HWFxDt1CvZrbIQG8CeiRX3tNIuxitU31I4DmkYOA_KA=', 'enable', 1684552960),
('88503398-db0f-11ec-9d64-0242ac120004', 1653474138, 1653468077, 'hrd', 'HWFxDt1CvZrbIQG8CeiRX3tNIuxitU31I4DmkYOA_KA=', 'enable', 1684552960);



/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;