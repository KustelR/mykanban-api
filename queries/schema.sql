create table Projects (
    id varchar(30) not null,
    name varchar(20) not null,
    primary key (id)
);
create table Columns (
    id varchar(30) not null,
    project_id varchar(30) not null,
    name varchar(20) not null,
    draw_order int unsigned not null,
    primary key (id),
    foreign key (project_id) references projects (id) ON DELETE CASCADE
);
create table Cards (
    id varchar(30) not null,
    column_id varchar(30) not null,
    name varchar(20) not null,
    description text not null,
    draw_order int unsigned not null,
    primary key (id),
    foreign key (column_id) references Columns (id) ON DELETE CASCADE
);
CREATE TABLE Tags (
    id varchar(30) not null,
    project_id varchar(30) not null,
    name varchar(20) not null,
    color varchar(7) not null,
    primary key (id),
    foreign key (project_id) references Projects (id) ON DELETE CASCADE
);
create table CardsTags (
    id int unsigned not null auto_increment,
    card_id varchar(30) not null,
    tag_id varchar(30) not null,
    primary key (id),
    foreign key (card_id) references cards (id) ON DELETE CASCADE,
    foreign key (tag_id) references tags (id) ON DELETE CASCADE
);
