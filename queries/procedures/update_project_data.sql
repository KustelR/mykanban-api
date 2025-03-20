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