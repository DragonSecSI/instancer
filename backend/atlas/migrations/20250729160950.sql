-- Modify "challenges" table
ALTER TABLE "public"."challenges" ALTER COLUMN "category" DROP DEFAULT;
-- Modify "instances" table
ALTER TABLE "public"."instances" ADD COLUMN "duration" bigint NOT NULL DEFAULT 1800;
-- Create index "idx_instances_active" to table: "instances"
CREATE INDEX "idx_instances_active" ON "public"."instances" ("active");
