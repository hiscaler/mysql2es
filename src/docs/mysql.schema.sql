CREATE TABLE `m2e_log` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`table_name` VARCHAR(100) NOT NULL,
	`pk_name` VARCHAR(100) NOT NULL,
	`pk_string_value` VARCHAR(100) NOT NULL,
	`pk_int_value` INT(11) NOT NULL DEFAULT '0',
	`last_modify` INT(11) NOT NULL,
	`version` VARCHAR(50) NOT NULL,
	PRIMARY KEY (`id`)
)
COLLATE='latin1_swedish_ci'
ENGINE=InnoDB
AUTO_INCREMENT=1
;
