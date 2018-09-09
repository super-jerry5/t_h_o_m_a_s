DROP DATABASE IF EXISTS d_thomas;
CREATE DATABASE d_thomas;

DROP TABLE IF EXISTS `d_thomas`.`job_info`;
CREATE TABLE `d_thomas`.`job_info` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `input_url` varchar(2048) COLLATE utf8mb4_bin NOT NULL,
  `input_url_backup` varchar(2048) COLLATE utf8mb4_bin NOT NULL,
  `output_url` varchar(2048) COLLATE utf8mb4_bin NOT NULL,
  `scale`  varchar(256) COLLATE utf8mb4_bin NOT NULL,
  `vcodec_type` int(11) NOT NULL COMMENT '0 h264; 1 h265',
  `vcodec_bitrate` int(11) NOT NULL,
  `watermark_text` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `watermark_param` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'fontsize:width:height',
  `audio_param`  varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `gpu_number` int(11) NOT NULL,
  `gpu_deint` int(11) NOT NULL,
  `worker` varchar(256) COLLATE utf8mb4_bin NOT NULL,
  `op` int(11) NOT NULL COMMENT '0 none; 1 start; 2 restart; 3 stop',
  `state` int(11) NOT NULL COMMENT '0 none; 1 todo; 2 doing; 3 done',
  `score` int(11) NOT NULL,
  `create_time` timestamp NOT NULL,
  `update_time` timestamp NOT NULL,
  `modify_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `state` (`state`) USING BTREE,
  INDEX `op` (`op`) USING BTREE,
  INDEX `worker` (`worker`) USING BTREE,
  INDEX `update_time` (`update_time`) USING BTREE,
  INDEX `create_time` (`create_time`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


