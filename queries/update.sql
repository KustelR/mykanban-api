# rename project
UPDATE Projects
SET Name = "test proj 2"
WHERE Id = "12adf2823a";
# rename column
UPDATE Columns
SET Name = "test col name 2"
WHERE Id = "13dsfa3g243";
#update card
UPDATE Cards
SET
name = ? # card name
description = ? # card description
column_id = ? # column id
order = ? # card order
WHERE id = ?; # card id