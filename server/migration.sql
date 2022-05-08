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


DELIMITER $$
$$
CREATE TRIGGER tgr_b_i_services
BEFORE INSERT
ON services FOR EACH ROW
BEGIN 
	set new.created_at = UNIX_TIMESTAMP();
	set new.updated_at = UNIX_TIMESTAMP();
END
$$

$$
CREATE TRIGGER trg_b_u_services
BEFORE UPDATE
ON services FOR EACH ROW
BEGIN 
	SET new.updated_at = UNIX_TIMESTAMP();
END
$$

DELIMITER ;


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



DELIMITER $$
$$
CREATE TRIGGER tgr_b_i_policies
BEFORE INSERT
ON policies FOR EACH ROW
BEGIN 
	set new.created_at = UNIX_TIMESTAMP();
	set new.updated_at = UNIX_TIMESTAMP();
END
$$

$$
CREATE TRIGGER trg_b_u_policies
BEFORE UPDATE
ON policies FOR EACH ROW
BEGIN 
	SET new.updated_at = UNIX_TIMESTAMP();
END
$$

DELIMITER ;



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


DELIMITER $$
$$
CREATE TRIGGER tgr_b_i_service_policy_mesh
BEFORE INSERT
ON service_policy_mesh FOR EACH ROW
BEGIN 
	set new.created_at = UNIX_TIMESTAMP();
	set new.updated_at = UNIX_TIMESTAMP();
END
$$

$$
CREATE TRIGGER trg_b_u_service_policy_mesh
BEFORE UPDATE
ON service_policy_mesh FOR EACH ROW
BEGIN 
	SET new.updated_at = UNIX_TIMESTAMP();
END
$$

DELIMITER ;





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
  UNIQUE KEY `objects_UN` (`global_id`),
  UNIQUE KEY `objects_golbal_UN` (`external_id`,`service_id`),
  KEY `objects_global_id_IDX` (`global_id`) USING BTREE,
  KEY `objects_external_id_IDX` (`external_id`) USING BTREE,
  KEY `objects_service_id_IDX` (`service_id`) USING BTREE,
  KEY `objects_status_IDX` (`status`) USING BTREE,
  KEY `objects_token_IDX` (`token`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



DELIMITER $$
$$
CREATE TRIGGER tgr_b_i_objects
BEFORE INSERT
ON objects FOR EACH ROW
BEGIN 
	set new.created_at = UNIX_TIMESTAMP();
	set new.updated_at = UNIX_TIMESTAMP();
END
$$

$$
CREATE TRIGGER trg_b_u_objects
BEFORE UPDATE
ON objects FOR EACH ROW
BEGIN 
	SET new.updated_at = UNIX_TIMESTAMP();
END
$$

DELIMITER ;






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


DELIMITER $$
$$
CREATE TRIGGER tgr_b_i_object_policy_mesh
BEFORE INSERT
ON object_policy_mesh FOR EACH ROW
BEGIN 
	set new.created_at = UNIX_TIMESTAMP();
	set new.updated_at = UNIX_TIMESTAMP();
END
$$

$$
CREATE TRIGGER trg_b_u_object_policy_mesh
BEFORE UPDATE
ON object_policy_mesh FOR EACH ROW
BEGIN 
	SET new.updated_at = UNIX_TIMESTAMP();
END
$$

DELIMITER ;
