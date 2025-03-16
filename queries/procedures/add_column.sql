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