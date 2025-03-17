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