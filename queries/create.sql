# add project
insert projects (
    name,
    id    
) values (
 'test 1', # name
 '12adf2823a' # id
);
# add column
insert columns (
    id,
    project_id,
    name,
    draw_order
) values (
    '13dsfa3g243', # id
    '12adf2823a', # associated project id
    'test col 2', # name
    1 # order
);
# add card
insert cards (
    id,
    column_id,
    name,
    description,
    draw_order
) values (
    "13def344a", # id
    "13dsfa3g243", # associated column id
    "test card 3", # name
    "test description 3", # card description
    2 # order
);
insert cards (
    id,
    column_id,
    name,
    description
) values (
    "234dfa3223", # id
    "13dsfa3g243", # associated column id
    "test card 2", # name
    "test description 2" # card description
);
# add tag
insert tags (
    id,
    project_id,
    name,
    color
) values (
    "1dcv82cg23", # id
    "12adf2823a", # project id
    "test tag 1", # name
    "#ff0000" # color
);
insert CardsTags (card_id, tag_id) values (

    "13def342a", # card_id
    "1dcv82cg23" # tag_id
);