DELIMITER $$

DROP PROCEDURE pop_column_reorder$$


CREATE PROCEDURE pop_column_reorder(p_project_id CHAR(30), p_popped_order INT)
BEGIN
    UPDATE Cards 
    SET 
    draw_order = draw_order - 1
    WHERE draw_order > p_popped_order AND project_id = p_project_id;
END$$

DELIMITER ;