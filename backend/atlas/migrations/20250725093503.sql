-- Modify "challenges" table
ALTER TABLE "public"."challenges" ADD COLUMN "values" text NOT NULL;
-- Modify "instances" table
ALTER TABLE "public"."instances" ADD COLUMN "active" boolean NOT NULL DEFAULT true;
