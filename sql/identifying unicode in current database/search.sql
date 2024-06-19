-- Find names with unicode.
SELECT name
FROM Players
WHERE name ~ '[^\x00-\x7F]'
  -- AND name NOT SIMILAR TO '%\u00A0%'     -- optional filter out the currently found unicode names.
