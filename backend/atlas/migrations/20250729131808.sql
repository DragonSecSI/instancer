-- Modify "challenges" table
ALTER TABLE "public"."challenges" ADD COLUMN "category" text NOT NULL DEFAULT '';
