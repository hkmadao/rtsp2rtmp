-- camera definition

-- Drop table

-- DROP TABLE camera;

CREATE TABLE camera (
	id varchar NOT NULL, -- id
	code varchar NULL, -- 摄像头编号
	rtsp_url varchar NULL, -- rtsp地址
	rtmp_url varchar NULL, -- rtmp地址
	play_auth_code varchar NULL, -- 播放识别码
	online_status INTEGER NULL, -- 是否在线：1.在线；0.不在线；
	enabled INTEGER NULL, -- 是否启用：1.启用；0.禁用；
	created timestamp NULL, -- 创建时间
	save_video INTEGER NULL, -- 是否保留录像：1.保留；0.不保留；
	live INTEGER NULL, -- 开启直播状态：1.开启；0.关闭；
	rtmp_push_status INTEGER NULL, -- 开启rtmp推送状态：1.开启；0.关闭；
	CONSTRAINT camera_pk PRIMARY KEY (id)
);

-- camera_share definition

-- Drop table

-- DROP TABLE camera_share;

CREATE TABLE camera_share (
	id varchar NOT NULL,
	camera_id varchar NULL, -- 摄像头标识
	auth_code varchar NULL, -- 播放权限码
	enabled varchar NULL, -- 启用状态：1.启用；0.禁用；
	created timestamp NULL, -- 创建时间
	deadline timestamp NULL, -- 截止日期
	"name" varchar NULL, -- 分享说明
	start_time timestamp NULL, -- 开始生效时间
	CONSTRAINT camera_share_pk PRIMARY KEY (id)
);