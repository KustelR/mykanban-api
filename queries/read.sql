# get project
select * from projects where id="12adf2823a"; # project id

# get all tags in project
select * from tags where project_id="12adf2823a" # project id

# get all columns from project
select columns.* from projects join
columns on columns.project_id = projects.id
where projects.id="12adf2823a"; # project id

# get all cards in column
select cards.* from columns join
cards on cards.column_id = columns.id
where columns.id="13dsfa3g43"; # column id

# get all tags in card
select Tags.* from Tags join
CardsTags on Tags.id = CardsTags.tag_id join
Cards on Cards.id = CardsTags.card_id
where cards.id="13def342a"; # card id

# get last order of thing
SELECT MAX(draw_order) FROM ?;