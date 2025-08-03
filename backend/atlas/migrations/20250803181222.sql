-- Create "challenges" table
CREATE TABLE "public"."challenges" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "description" text NOT NULL,
  "category" text NOT NULL,
  "type" bigint NOT NULL,
  "remote_id" text NOT NULL,
  "flag" text NOT NULL,
  "flag_type" bigint NOT NULL,
  "duration" bigint NOT NULL,
  "repository" text NOT NULL,
  "chart" text NOT NULL,
  "chart_version" text NOT NULL,
  "values" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_challenges_name" to table: "challenges"
CREATE UNIQUE INDEX "idx_challenges_name" ON "public"."challenges" ("name");
-- Create index "idx_challenges_remote_id" to table: "challenges"
CREATE UNIQUE INDEX "idx_challenges_remote_id" ON "public"."challenges" ("remote_id");
-- Create "teams" table
CREATE TABLE "public"."teams" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "remote_id" text NOT NULL,
  "token" text NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_teams_name" to table: "teams"
CREATE UNIQUE INDEX "idx_teams_name" ON "public"."teams" ("name");
-- Create index "idx_teams_remote_id" to table: "teams"
CREATE UNIQUE INDEX "idx_teams_remote_id" ON "public"."teams" ("remote_id");
-- Create index "idx_teams_token" to table: "teams"
CREATE UNIQUE INDEX "idx_teams_token" ON "public"."teams" ("token");
-- Create "instances" table
CREATE TABLE "public"."instances" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "flag" text NOT NULL,
  "team_id" bigint NOT NULL,
  "challenge_id" bigint NOT NULL,
  "challenge_type" bigint NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "active" boolean NOT NULL DEFAULT true,
  "duration" bigint NOT NULL DEFAULT 1800,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_instances_challenge" FOREIGN KEY ("challenge_id") REFERENCES "public"."challenges" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_instances_team" FOREIGN KEY ("team_id") REFERENCES "public"."teams" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_instances_active" to table: "instances"
CREATE INDEX "idx_instances_active" ON "public"."instances" ("active");
-- Create index "idx_instances_challenge_id" to table: "instances"
CREATE INDEX "idx_instances_challenge_id" ON "public"."instances" ("challenge_id");
-- Create index "idx_instances_flag" to table: "instances"
CREATE UNIQUE INDEX "idx_instances_flag" ON "public"."instances" ("flag");
-- Create index "idx_instances_name" to table: "instances"
CREATE UNIQUE INDEX "idx_instances_name" ON "public"."instances" ("name");
-- Create index "idx_instances_team_id" to table: "instances"
CREATE INDEX "idx_instances_team_id" ON "public"."instances" ("team_id");
