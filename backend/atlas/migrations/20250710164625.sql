-- Create "instances" table
CREATE TABLE "public"."instances" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_instances_name" to table: "instances"
CREATE UNIQUE INDEX "idx_instances_name" ON "public"."instances" ("name");
