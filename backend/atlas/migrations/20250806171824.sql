-- Modify "challenges" table
ALTER TABLE "public"."challenges" ADD COLUMN "cooldown" bigint NOT NULL DEFAULT 0;
