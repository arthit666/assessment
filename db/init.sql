-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS expenses_id_seq;

-- Table Definition
CREATE TABLE "expenses" (
    "id" int4 NOT NULL DEFAULT nextval('expenses_id_seq'::regclass),
    "title" text,
    "amount" float,
    "note" text,
    "tags" text ARRAY,
    PRIMARY KEY ("id")
);

INSERT INTO "expenses" ("title", "amount", "note", "tags") VALUES ('test-title', 99.0, 'test-note',ARRAY['test-tags1', 'test-tags2']);