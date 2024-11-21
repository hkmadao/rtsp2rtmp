-- public.camera definition

-- Drop table

-- DROP TABLE public.camera;

CREATE TABLE public.camera (
	id varchar NOT NULL, -- id
	code varchar NULL, -- 摄像头编号
	rtsp_url varchar NULL, -- rtsp地址
	rtmp_url varchar NULL, -- rtmp地址
	play_auth_code varchar NULL, -- 播放识别码
	online_status int2 NULL, -- 是否在线：1.在线；0.不在线；
	enabled int2 NULL, -- 是否启用：1.启用；0.禁用；
	created timestamp(0) NULL, -- 创建时间
	save_video int2 NULL, -- 是否保留录像：1.保留；0.不保留；
	live int2 NULL, -- 开启直播状态：1.开启；0.关闭；
	rtmp_push_status int2 NULL, -- 开启rtmp推送状态：1.开启；0.关闭；
	CONSTRAINT camera_pk PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.camera.id IS 'id';
COMMENT ON COLUMN public.camera.code IS '摄像头编号';
COMMENT ON COLUMN public.camera.rtsp_url IS 'rtsp地址';
COMMENT ON COLUMN public.camera.rtmp_url IS 'rtmp地址';
COMMENT ON COLUMN public.camera.play_auth_code IS '播放识别码';
COMMENT ON COLUMN public.camera.online_status IS '是否在线：1.在线；0.不在线；';
COMMENT ON COLUMN public.camera.enabled IS '是否启用：1.启用；0.禁用；';
COMMENT ON COLUMN public.camera.created IS '创建时间';
COMMENT ON COLUMN public.camera.save_video IS '是否保留录像：1.保留；0.不保留；';
COMMENT ON COLUMN public.camera.live IS '开启直播状态：1.开启；0.关闭；';
COMMENT ON COLUMN public.camera.rtmp_push_status IS '开启rtmp推送状态：1.开启；0.关闭；';

-- public.camera_share definition

-- Drop table

-- DROP TABLE public.camera_share;

CREATE TABLE public.camera_share (
	id varchar NOT NULL,
	camera_id varchar NULL, -- 摄像头标识
	auth_code varchar NULL, -- 播放权限码
	enabled varchar NULL, -- 启用状态：1.启用；0.禁用；
	created timestamp(0) NULL, -- 创建时间
	deadline timestamp(0) NULL, -- 截止日期
	"name" varchar NULL, -- 分享说明
	start_time timestamp(0) NULL, -- 开始生效时间
	CONSTRAINT camera_share_pk PRIMARY KEY (id)
);
COMMENT ON TABLE public.camera_share IS '摄像头分享表';

-- Column comments

COMMENT ON COLUMN public.camera_share.camera_id IS '摄像头标识';
COMMENT ON COLUMN public.camera_share.auth_code IS '播放权限码';
COMMENT ON COLUMN public.camera_share.enabled IS '启用状态：1.启用；0.禁用；';
COMMENT ON COLUMN public.camera_share.created IS '创建时间';
COMMENT ON COLUMN public.camera_share.deadline IS '截止日期';
COMMENT ON COLUMN public.camera_share."name" IS '分享说明';
COMMENT ON COLUMN public.camera_share.start_time IS '开始生效时间';