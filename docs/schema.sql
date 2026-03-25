DROP TABLE models;
DROP TABLE views;
DROP TABLE existences;

CREATE TABLE models
(
    id          INTEGER PRIMARY KEY ASC,
    explanation TEXT NOT NULL
) STRICT;

INSERT INTO models (id, explanation)
VALUES (0, 'undefined');

CREATE TABLE views
(
    id       INTEGER PRIMARY KEY ASC,
    name     TEXT    NOT NULL,
    model_id INTEGER NOT NULL,
    FOREIGN KEY (model_id) REFERENCES models (id)
) STRICT;

CREATE INDEX idx_views_model_id ON views (model_id);

CREATE INDEX idx_views_name ON views (name);

CREATE TABLE existences
(
    id      INTEGER PRIMARY KEY ASC,
    view_id INTEGER                                 NOT NULL,
    source  TEXT                                    NOT NULL,
    quote   TEXT                                    NOT NULL,
    reason  TEXT                                    NOT NULL, -- why that view of model is chosen in this existence
    tag     INTEGER CHECK ( tag >= -1 AND tag <= 3) NOT NULL, -- -1 for undefined.
    why_not TEXT                                    NOT NULL, -- exists but invalid, shall be ignored reason, empty if valid
    FOREIGN KEY (view_id) REFERENCES views (id)
) STRICT;

CREATE UNIQUE INDEX idx_existences_view_id_source_quote ON existences (view_id, source, quote);


-- tag is defined as a common group of definitions that the row matches.
-- The enum has been defined in the manual form long ago, we just keep it.
-- | 级别 | 描述                                                      |
-- | ---- | --------------------------------------------------------- |
-- | 0    | 我使用了某个汉字词，即使知道其英文词，也不认为适合使用。  |
-- | 1    | 我确信对应中文，但因为要装逼等原因我 TMD 就要用 English。 |
-- | 2    | 我知道对应中文，但权衡后认为使用英文更合适。              |
-- | 3    | 我一时间还想不起来对应中文，在这里机翻或现造是不对的。    |
