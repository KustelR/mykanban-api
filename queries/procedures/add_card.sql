DELIMITER $$

DROP PROCEDURE add_card$$

CREATE PROCEDURE add_card(p_column_id CHAR(30), p_id CHAR(30), p_name CHAR(20), p_description TEXT, p_draw_order INT)
BEGIN
    INSERT cards (
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