
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
    foreign key (project_id) references Projects (id) ON DELETE CASCADE
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
    foreign key (card_id) references Cards (id) ON DELETE CASCADE,
    foreign key (tag_id) references Tags (id) ON DELETE CASCADE
);



DELIMITER $$

DROP PROCEDURE add_card$$

CREATE PROCEDURE add_card(p_column_id CHAR(30), p_id CHAR(30), p_name CHAR(20), p_description TEXT, p_draw_order INT)
BEGIN
    INSERT Cards (
        id,
        column_id,
        name,
        description,
        draw_order
    ) VALUES (
        p_id,
        p_column_id,
        p_name,
        p_description,
        p_draw_order
    ) ON DUPLICATE KEY UPDATE
        column_id = p_column_id,
        name = p_name,
        description = p_description,
        draw_order = p_draw_order;
 END$$

 DELIMITER ;


DELIMITER $$

DROP PROCEDURE add_column$$

CREATE PROCEDURE add_column(p_project_id char(30), p_id char(30), p_name char(20), p_draw_order INT)
BEGIN
    INSERT Columns (
        id,
        project_id,
        name,
        draw_order
    ) VALUES (
        p_id, 
        p_project_id, 
        p_name, 
        p_draw_order
    ) ON DUPLICATE KEY UPDATE 
        id=p_id, 
        project_id=p_project_id, 
        name=p_name, 
        draw_order=p_draw_order;
END$$

DELIMITER ;


DELIMITER $$

DROP PROCEDURE pop_card_reorder$$


CREATE PROCEDURE pop_card_reorder(p_column_id CHAR(30), p_popped_order INT)
BEGIN
    UPDATE Cards 
    SET 
    draw_order = draw_order - 1
    WHERE draw_order > p_popped_order AND column_id = p_column_id;
END$$

DELIMITER ;


DELIMITER $$

DROP PROCEDURE pop_column_reorder$$


CREATE PROCEDURE pop_column_reorder(p_project_id CHAR(30), p_popped_order INT)
BEGIN
    UPDATE Columns 
    SET 
    draw_order = draw_order - 1
    WHERE draw_order > p_popped_order AND project_id = p_project_id;
END$$

DELIMITER ;


DELIMITER $$

DROP PROCEDURE update_card$$


CREATE PROCEDURE update_card(p_id CHAR(30), p_column_id CHAR(30), p_name CHAR(20), p_description TEXT, p_draw_order INT)
BEGIN
    UPDATE Cards 
    SET 
    name = p_name,
    description = p_description,
    draw_order = p_draw_order,
    column_id = p_column_id
    WHERE id = p_id;
END$$

DELIMITER ;



DELIMITER $$

DROP PROCEDURE update_column_data$$


CREATE PROCEDURE update_column_data(p_id CHAR(30), p_name CHAR(20), p_draw_order INT)
BEGIN
    UPDATE Columns 
    SET 
    name = p_name,
    draw_order = p_draw_order
    WHERE id = p_id;
END$$

DELIMITER ;



DELIMITER $$

DROP PROCEDURE update_project_data$$


CREATE PROCEDURE update_project_data(p_id CHAR(30), p_name CHAR(20))
BEGIN
    UPDATE Projects 
    SET 
    name = p_name
    WHERE id = p_id;
END$$

DELIMITER ;