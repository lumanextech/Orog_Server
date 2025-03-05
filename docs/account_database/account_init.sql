CREATE TABLE "account" (
	-- 用户地址
	"address" VARCHAR(255) UNIQUE,
	-- 链
	"chain" VARCHAR(255),
	-- 用户名
	"username" VARCHAR(255),
	-- 创建时间
	"created_at" TIMESTAMP,
	"updated_at" TIMESTAMP,
	-- 初始资金
	"initial_funding" float8,
	-- 用户余额
	"bakance" float8,
    -- 角色id
    "role_id" INTEGER,
    -- 邀请人地址
    "invited_address" VARCHAR(255),
    -- 邀请人地址
    "invited_code" VARCHAR(255),
	PRIMARY KEY("address")
);
COMMENT ON COLUMN account.address IS '用户地址';
COMMENT ON COLUMN account.chain IS '链';
COMMENT ON COLUMN account.username IS '用户名';
COMMENT ON COLUMN account.created_at IS '创建时间';
COMMENT ON COLUMN account.initial_funding IS '初始资金';
COMMENT ON COLUMN account.bakance IS '用户余额';


CREATE TABLE "adily_money_change" (
	-- 用户地址
	"address" VARCHAR(255) NOT NULL UNIQUE,
	-- 链
	"chain" VARCHAR(255),
	-- 变动后金额
	"amount_change" DECIMAL,
	-- 变动理由
	"reason" DECIMAL,
	"created_at" TIMESTAMP,
	"id" INTEGER NOT NULL,
	PRIMARY KEY("id")
);
COMMENT ON COLUMN adily_money_change.address IS '用户地址';
COMMENT ON COLUMN adily_money_change.chain IS '链';
COMMENT ON COLUMN adily_money_change.amount_change IS '变动后金额';
COMMENT ON COLUMN adily_money_change.reason IS '变动理由';


CREATE TABLE "user_token_follow" (
	"address" VARCHAR(255) NOT NULL UNIQUE,
	"chain" VARCHAR(255),
	"token_address" VARCHAR(255),
	"status" SMALLINT,  
	"followed_at" TIMESTAMP,
	"unfollowed_at" TIMESTAMP,
	"created_at" TIMESTAMP,
	"updated_at" TIMESTAMP,
	"id" INTEGER NOT NULL,
	PRIMARY KEY("id")
);


ALTER TABLE "adily_money_change"
ADD FOREIGN KEY("address") REFERENCES "account"("address")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "user_token_follow"
ADD FOREIGN KEY("address") REFERENCES "account"("address")
ON UPDATE NO ACTION ON DELETE NO ACTION;