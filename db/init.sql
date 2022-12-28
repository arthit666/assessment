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

INSERT INTO "expenses" ("id", "title", "amount", "note", "tags") VALUES (1, 'test-title', 50.0, 'test-note',ARRAY['test-tags1', 'test-tags2']);