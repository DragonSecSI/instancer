-- Modify "challenges" table
ALTER TABLE "public"."challenges" ADD COLUMN "type" bigint;
UPDATE "public"."challenges" SET "type" = 0 WHERE "type" IS NULL;
ALTER TABLE "public"."challenges" ALTER COLUMN "type" SET NOT NULL;
