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