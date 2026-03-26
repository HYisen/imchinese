-- This file contains SQL that could be used after schema.sql.
-- For dummy data initialization and SQL example playground.

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

INSERT INTO existences(id, view_id, source, quote, reason, tag, why_not)
VALUES (1, 1, 'GoVersionChronology/缘起', '所以你可以猜到我会选择 latest。', 'go get example.com/m@latest', 2, ''),
       (2, 3, 'GoVersionChronology/meme/用词的梗', '', '短', 0, '');

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

SELECT views.id, model_id, name, COUNT(existences.id)
FROM existences
         JOIN views ON existences.view_id = views.id
GROUP BY view_id;

SELECT *
FROM existences
         JOIN views ON existences.view_id = views.id
         JOIN models ON views.model_id = models.id;
