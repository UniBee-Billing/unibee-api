

-- alter table merchant_discount_code modify column `code` varchar(200) CHARACTER SET utf8 COLLATE utf8_bin

/* SELECT LOWER(code) AS lower_name, COUNT(*)
FROM merchant_discount_code
GROUP BY LOWER(code),merchant_id COLLATE utf8_general_ci
HAVING COUNT(*) > 1; */
