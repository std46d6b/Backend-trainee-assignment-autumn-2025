CREATE TYPE "pull_request_status" AS ENUM (
  'OPEN',
  'MERGED'
);

CREATE TABLE "teams" (
  "team_name" text PRIMARY KEY
);

CREATE TABLE "users" (
  "user_id" text PRIMARY KEY,
  "username" text NOT NULL,
  "team_name" text NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true
);

CREATE TABLE "pull_requests" (
  "pull_request_id" text PRIMARY KEY,
  "pull_request_name" text NOT NULL,
  "author_id" text NOT NULL,
  "status" pull_request_status NOT NULL DEFAULT 'OPEN',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "merged_at" timestamptz
);

CREATE TABLE "assigned_reviewers" (
  "pull_request_id" text NOT NULL,
  "user_id" text NOT NULL,
  PRIMARY KEY ("pull_request_id", "user_id")
);

CREATE INDEX "idx_users_team_name" ON "users" ("team_name");

CREATE INDEX "idx_assigned_reviewers_user_id" ON "assigned_reviewers" ("user_id");

ALTER TABLE "users" ADD FOREIGN KEY ("team_name") REFERENCES "teams" ("team_name");

ALTER TABLE "pull_requests" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("user_id");

ALTER TABLE "assigned_reviewers" ADD FOREIGN KEY ("pull_request_id") REFERENCES "pull_requests" ("pull_request_id");

ALTER TABLE "assigned_reviewers" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");
