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

INSERT INTO models (id, explanation)
VALUES (1, 'the most recently updated version'),
       (2, '');

INSERT INTO views(id, name, model_id)
VALUES (1, 'latest', 1),
       (2, '最新版本', 1),
       (3, '梗', 2),
       (4, 'meme', 2),
       (5, '典故', 2),
       (6, '模因', 2),
       (7, '成语', 2),
       (8, 'neta', 2);

INSERT INTO existences(id, view_id, source, quote, reason, tag)
VALUES (1, 1, 'GoVersionChronology/缘起', '所以你可以猜到我会选择 latest。', 'go get example.com/m@latest', 2),
       (2, 3, 'GoVersionChronology/meme/用词的梗', '', '短', 0);

-- This is the SQL that could output the original manual table format.
-- To be honest, I would rather not implement the combine logic in SQL.
-- Just a proof of concept that the DDL could store all the data.
SELECT name, alternatives, reason, tag, source
FROM (SELECT model_id,
             group_concat(name, '，') || iif(length(m.explanation) > 0, ',' || m.explanation, '') AS alternatives
      FROM views
               LEFT JOIN existences e ON views.id = e.view_id
               LEFT JOIN models m on views.model_id = m.id
      WHERE e.id IS NULL
      GROUP BY model_id, e.id) AS sub
         LEFT JOIN views v ON sub.model_id = v.model_id
         LEFT JOIN existences e ON v.id = e.view_id
WHERE e.id IS NOT NULL;

SELECT *
FROM views
         LEFT JOIN models ON views.model_id = models.id;