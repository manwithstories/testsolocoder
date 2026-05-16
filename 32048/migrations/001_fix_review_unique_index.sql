-- 修复评价表唯一索引，支持买卖双方对同一交易进行评价
-- 执行此脚本前请先备份数据库

BEGIN;

-- 删除旧的唯一索引（如果存在）
DROP INDEX IF EXISTS idx_reviews_transaction_id;
DROP INDEX IF EXISTS uix_reviews_transaction_id;

-- 创建新的联合唯一索引：交易ID + 评价人ID
CREATE UNIQUE INDEX IF NOT EXISTS idx_transaction_reviewer 
ON reviews (transaction_id, reviewer_id);

COMMIT;
