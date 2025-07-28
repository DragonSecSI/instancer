-- Modify "instances" table
ALTER TABLE "public"."instances" ADD COLUMN "challenge_type" bigint;
UPDATE "public"."instances" SET "challenge_type" = 1 WHERE "challenge_type" IS NULL;
ALTER TABLE "public"."instances" ALTER COLUMN "challenge_type" SET NOT NULL;
