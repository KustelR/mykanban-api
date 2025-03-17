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