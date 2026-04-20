-- IAM 系统数据库初始化脚本
-- MySQL 8.0

CREATE DATABASE IF NOT EXISTS `iam` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `iam`;

-- ============================================
-- 1. tenants — 租户表
-- ============================================
CREATE TABLE `tenants` (
    `id`         BIGINT       NOT NULL AUTO_INCREMENT,
    `name`       VARCHAR(100) NOT NULL,
    `status`     TINYINT      NOT NULL DEFAULT 1 COMMENT '1=启用 2=禁用 3=过期',
    `max_users`  INT          NOT NULL DEFAULT 100 COMMENT '最大用户配额',
    `max_apps`   INT          NOT NULL DEFAULT 10  COMMENT '最大应用配额',
    `expire_at`  DATETIME              DEFAULT NULL COMMENT '租户过期时间',
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

-- ============================================
-- 2. users — 用户表
-- ============================================
CREATE TABLE `users` (
    `id`                  BIGINT        NOT NULL AUTO_INCREMENT,
    `tenant_id`           BIGINT        NOT NULL COMMENT '租户 ID',
    `email`               VARCHAR(100)  NOT NULL COMMENT '邮箱（登录账号）',
    `phone`               VARCHAR(20)            DEFAULT NULL COMMENT '手机号',
    `password_hash`       VARCHAR(255)  NOT NULL COMMENT 'bcrypt 哈希',
    `status`              TINYINT       NOT NULL DEFAULT 1 COMMENT '1=启用 2=禁用 3=锁定',
    `mfa_enabled`         TINYINT       NOT NULL DEFAULT 0 COMMENT '是否开启 MFA',
    `mfa_secret`          VARCHAR(100)           DEFAULT NULL COMMENT 'TOTP 密钥',
    `last_login_at`       DATETIME               DEFAULT NULL COMMENT '最后登录时间',
    `password_changed_at` DATETIME               DEFAULT NULL COMMENT '密码最后修改时间',
    `created_at`          DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_email` (`tenant_id`, `email`),
    INDEX `idx_tenant_status` (`tenant_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 2.1 user_invitations — 用户邀请表
-- ============================================
CREATE TABLE `user_invitations` (
    `id`              BIGINT       NOT NULL AUTO_INCREMENT,
    `tenant_id`       BIGINT       NOT NULL COMMENT '目标租户 ID',
    `email`           VARCHAR(100) NOT NULL COMMENT '受邀邮箱',
    `inviter_user_id` BIGINT                DEFAULT NULL COMMENT '邀请人用户 ID',
    `token_hash`      VARCHAR(255) NOT NULL COMMENT '邀请令牌哈希',
    `status`          TINYINT      NOT NULL DEFAULT 1 COMMENT '1=待接受 2=已接受 3=已过期 4=已撤销',
    `expires_at`      DATETIME     NOT NULL COMMENT '过期时间',
    `accepted_at`     DATETIME              DEFAULT NULL COMMENT '接受时间',
    `created_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_token_hash` (`token_hash`),
    INDEX `idx_tenant_email_status` (`tenant_id`, `email`, `status`),
    INDEX `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户邀请表';

-- ============================================
-- 3. user_groups — 用户组表
-- ============================================
CREATE TABLE `user_groups` (
    `id`          BIGINT       NOT NULL AUTO_INCREMENT,
    `tenant_id`   BIGINT       NOT NULL COMMENT '租户 ID',
    `name`        VARCHAR(100) NOT NULL COMMENT '用户组名称',
    `description` TEXT                  DEFAULT NULL,
    `parent_id`   BIGINT                DEFAULT NULL COMMENT '父组 ID',
    `level`       INT          NOT NULL DEFAULT 1 COMMENT '层级深度',
    `path`        VARCHAR(500)          DEFAULT NULL COMMENT '完整路径',
    `is_system`   TINYINT      NOT NULL DEFAULT 0 COMMENT '是否系统组',
    `group_type`  VARCHAR(20)  NOT NULL DEFAULT 'NORMAL' COMMENT '组类型',
    `sort_order`  INT          NOT NULL DEFAULT 0 COMMENT '排序号',
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_tenant` (`tenant_id`),
    INDEX `idx_parent` (`parent_id`),
    INDEX `idx_path` (`path`(100))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户组表';

-- ============================================
-- 4. user_group_members — 用户组成员表
-- ============================================
CREATE TABLE `user_group_members` (
    `id`         BIGINT   NOT NULL AUTO_INCREMENT,
    `tenant_id`  BIGINT   NOT NULL COMMENT '租户 ID',
    `group_id`   BIGINT   NOT NULL COMMENT '用户组 ID',
    `user_id`    BIGINT   NOT NULL COMMENT '用户 ID',
    `created_by` BIGINT            DEFAULT NULL COMMENT '操作人',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_group_user` (`group_id`, `user_id`),
    INDEX `idx_tenant` (`tenant_id`),
    INDEX `idx_group` (`group_id`),
    INDEX `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户组成员表';

-- ============================================
-- 5. permissions — 权限定义表
-- ============================================
CREATE TABLE `permissions` (
    `id`          BIGINT       NOT NULL AUTO_INCREMENT,
    `code`        VARCHAR(100) NOT NULL COMMENT '权限编码，如 user:read',
    `name`        VARCHAR(100) NOT NULL COMMENT '权限名称',
    `resource`    VARCHAR(100) NOT NULL COMMENT '资源类型',
    `action`      VARCHAR(50)  NOT NULL COMMENT '操作：read/write/delete',
    `app_code`    VARCHAR(50)           DEFAULT NULL COMMENT '归属应用（NULL=平台级）',
    `description` TEXT                  DEFAULT NULL,
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    INDEX `idx_app_code` (`app_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限定义表';

-- ============================================
-- 6. roles — 角色表
-- ============================================
CREATE TABLE `roles` (
    `id`          BIGINT       NOT NULL AUTO_INCREMENT,
    `tenant_id`   BIGINT       NOT NULL COMMENT '租户 ID',
    `name`        VARCHAR(100) NOT NULL COMMENT '角色名称',
    `code`        VARCHAR(100) NOT NULL COMMENT '角色编码',
    `type`        TINYINT      NOT NULL DEFAULT 2 COMMENT '1=系统内置 2=自定义',
    `status`      TINYINT      NOT NULL DEFAULT 1 COMMENT '1=启用 2=禁用',
    `description` TEXT                  DEFAULT NULL,
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`),
    INDEX `idx_tenant` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ============================================
-- 7. role_permissions — 角色权限关联表
-- ============================================
CREATE TABLE `role_permissions` (
    `id`            BIGINT      NOT NULL AUTO_INCREMENT,
    `role_id`       BIGINT      NOT NULL COMMENT '角色 ID',
    `permission_id` BIGINT      NOT NULL COMMENT '权限 ID',
    `data_scope`    VARCHAR(50) NOT NULL DEFAULT 'all' COMMENT 'all/dept/dept_and_sub/personal/custom',
    `created_at`    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_perm` (`role_id`, `permission_id`),
    INDEX `idx_role` (`role_id`),
    INDEX `idx_permission` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- ============================================
-- 8. user_roles — 用户角色关联表
-- ============================================
CREATE TABLE `user_roles` (
    `id`         BIGINT       NOT NULL AUTO_INCREMENT,
    `tenant_id`  BIGINT       NOT NULL COMMENT '租户 ID',
    `user_id`    BIGINT       NOT NULL COMMENT '用户 ID',
    `role_id`    BIGINT       NOT NULL COMMENT '角色 ID',
    `app_code`   VARCHAR(50)           DEFAULT NULL COMMENT '应用编码（角色生效范围）',
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_user_role_app` (`tenant_id`, `user_id`, `role_id`, `app_code`),
    INDEX `idx_tenant` (`tenant_id`),
    INDEX `idx_user` (`user_id`),
    INDEX `idx_role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ============================================
-- 9. role_constraints — 角色约束表（SoD 约束）
-- ============================================
CREATE TABLE `role_constraints` (
    `id`         BIGINT   NOT NULL AUTO_INCREMENT,
    `tenant_id`  BIGINT   NOT NULL COMMENT '租户 ID',
    `type`       TINYINT  NOT NULL COMMENT '1=静态SoD 2=动态SoD',
    `role_a`     BIGINT   NOT NULL COMMENT '冲突角色 A',
    `role_b`     BIGINT   NOT NULL COMMENT '冲突角色 B',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_tenant` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色约束表（SoD）';

-- ============================================
-- 10. applications — 应用表
-- ============================================
CREATE TABLE `applications` (
    `id`          BIGINT       NOT NULL AUTO_INCREMENT,
    `tenant_id`   BIGINT       NOT NULL COMMENT '租户 ID',
    `code`        VARCHAR(50)  NOT NULL COMMENT '应用编码（租户内唯一）',
    `name`        VARCHAR(100) NOT NULL COMMENT '应用名称',
    `description` TEXT                  DEFAULT NULL,
    `status`      TINYINT      NOT NULL DEFAULT 1 COMMENT '1=启用 2=禁用',
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`),
    INDEX `idx_tenant` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用表';

-- ============================================
-- 11. user_app_authorizations — 用户应用授权表
-- ============================================
CREATE TABLE `user_app_authorizations` (
    `id`         BIGINT   NOT NULL AUTO_INCREMENT,
    `tenant_id`  BIGINT   NOT NULL COMMENT '租户 ID',
    `user_id`    BIGINT   NOT NULL COMMENT '用户 ID',
    `app_id`     BIGINT   NOT NULL COMMENT '应用 ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_user_app` (`tenant_id`, `user_id`, `app_id`),
    INDEX `idx_tenant` (`tenant_id`),
    INDEX `idx_user` (`user_id`),
    INDEX `idx_app` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户应用授权表';

-- ============================================
-- 12. clients — 内部客户端表
-- ============================================
CREATE TABLE `clients` (
    `id`                   BIGINT       NOT NULL AUTO_INCREMENT,
    `client_id`            VARCHAR(64)  NOT NULL COMMENT '客户端标识',
    `name`                 VARCHAR(100) NOT NULL COMMENT '客户端名称',
    `allowed_scopes`       JSON         NOT NULL COMMENT '允许的 scopes',
    `access_token_ttl_sec` INT          NOT NULL DEFAULT 600 COMMENT 'Access Token TTL（秒）',
    `status`               TINYINT      NOT NULL DEFAULT 1 COMMENT '1=启用 2=禁用',
    `created_at`           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_client_id` (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内部客户端表';

-- ============================================
-- 12.1 client_credentials — 客户端凭证表
-- ============================================
CREATE TABLE `client_credentials` (
    `id`            BIGINT       NOT NULL AUTO_INCREMENT,
    `client_id`     BIGINT       NOT NULL COMMENT '客户端主键 ID',
    `access_key_id` VARCHAR(64)  NOT NULL COMMENT 'AK 标识',
    `secret_hash`   VARCHAR(255) NOT NULL COMMENT 'SK 哈希值',
    `secret_hint`   VARCHAR(16)           DEFAULT NULL COMMENT 'SK 提示信息',
    `status`        TINYINT      NOT NULL DEFAULT 1 COMMENT '1=启用 2=禁用 3=过期',
    `expires_at`    DATETIME              DEFAULT NULL COMMENT '过期时间',
    `last_used_at`  DATETIME              DEFAULT NULL COMMENT '最近使用时间',
    `rotated_at`    DATETIME              DEFAULT NULL COMMENT '轮换时间',
    `created_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_access_key_id` (`access_key_id`),
    INDEX `idx_client_status` (`client_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户端凭证表';

-- ============================================
-- 13. audit_logs — 操作审计日志表
-- ============================================
CREATE TABLE `audit_logs` (
    `id`            BIGINT       NOT NULL AUTO_INCREMENT,
    `tenant_id`     BIGINT       NOT NULL COMMENT '租户 ID',
    `app_id`        BIGINT                DEFAULT NULL COMMENT '应用 ID',
    `user_id`       BIGINT       NOT NULL COMMENT '操作人',
    `action`        VARCHAR(100) NOT NULL COMMENT '操作类型',
    `resource_type` VARCHAR(50)  NOT NULL COMMENT '资源类型',
    `resource_id`   BIGINT                DEFAULT NULL COMMENT '资源 ID',
    `detail`        JSON                  DEFAULT NULL COMMENT '操作详情',
    `ip`            VARCHAR(45)           DEFAULT NULL COMMENT '操作 IP',
    `created_at`    DATETIME     NOT NULL COMMENT '操作时间',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant` (`tenant_id`),
    INDEX `idx_app` (`app_id`),
    INDEX `idx_user` (`user_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作审计日志表';

-- ============================================
-- 14. login_logs — 登录日志表
-- ============================================
CREATE TABLE `login_logs` (
    `id`          BIGINT        NOT NULL AUTO_INCREMENT,
    `tenant_id`   BIGINT        NOT NULL COMMENT '租户 ID',
    `user_id`     BIGINT                 DEFAULT NULL COMMENT 'NULL=登录失败',
    `app_id`      BIGINT                 DEFAULT NULL COMMENT '应用 ID',
    `email`       VARCHAR(100)  NOT NULL COMMENT '登录账号',
    `status`      TINYINT       NOT NULL COMMENT '1=成功 2=失败 3=MFA待验证',
    `fail_reason` VARCHAR(200)           DEFAULT NULL COMMENT '失败原因',
    `login_type`  VARCHAR(30)   NOT NULL COMMENT 'password/code/oauth/mfa',
    `ip`          VARCHAR(45)            DEFAULT NULL COMMENT '登录 IP',
    `user_agent`  VARCHAR(500)           DEFAULT NULL COMMENT '用户代理',
    `created_at`  DATETIME      NOT NULL COMMENT '登录时间',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant` (`tenant_id`),
    INDEX `idx_app` (`app_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';

-- ============================================
-- 15. password_policies — 密码策略表
-- ============================================
CREATE TABLE `password_policies` (
    `id`                  BIGINT   NOT NULL AUTO_INCREMENT,
    `tenant_id`           BIGINT   NOT NULL COMMENT '每租户一条',
    `min_length`          INT      NOT NULL DEFAULT 8 COMMENT '最小密码长度',
    `require_uppercase`   TINYINT  NOT NULL DEFAULT 1 COMMENT '需要大写字母',
    `require_lowercase`   TINYINT  NOT NULL DEFAULT 1 COMMENT '需要小写字母',
    `require_digit`       TINYINT  NOT NULL DEFAULT 1 COMMENT '需要数字',
    `require_special`     TINYINT  NOT NULL DEFAULT 1 COMMENT '需要特殊字符',
    `history_count`       INT      NOT NULL DEFAULT 3 COMMENT '历史密码检查次数',
    `expire_days`         INT      NOT NULL DEFAULT 0 COMMENT '密码过期天数（0=永不过期）',
    `max_login_attempts`  INT      NOT NULL DEFAULT 5 COMMENT '最大登录失败次数',
    `lockout_minutes`     INT      NOT NULL DEFAULT 30 COMMENT '锁定时长（分钟）',
    `updated_at`          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='密码策略表';

-- ============================================
-- 16. password_history — 密码历史表
-- ============================================
CREATE TABLE `password_history` (
    `id`            BIGINT       NOT NULL AUTO_INCREMENT,
    `user_id`       BIGINT       NOT NULL COMMENT '用户 ID',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '历史密码哈希',
    `created_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_user_created` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='密码历史表';
