USE ginger_finance;

DROP FUNCTION IF EXISTS record_trigger_action;
DELIMITER $$
CREATE FUNCTION record_trigger_action(vname VARCHAR(1024), vrowid INT, vstatus VARCHAR(128), vmsg VARCHAR(1024)) RETURNS TINYINT
BEGIN
    /*触发器日志*/
    INSERT INTO system_trigger_log(`name`, `row_id`, `status`, `msg`, `created_at`) values(vname, vrowid, vstatus, vmsg, NOW());
    RETURN 1;
END
$$
DELIMITER ;

DROP FUNCTION IF EXISTS should_ignore_balance;
DELIMITER $$
CREATE FUNCTION should_ignore_balance(code VARCHAR(32)) RETURNS TINYINT
BEGIN
    /*账户余额不需要更新*/
    return LOCATE("_0", code) || LOCATE("user_rmb_", code);
END
$$
DELIMITER ;

DROP FUNCTION IF EXISTS can_overdraw_account;
DELIMITER $$
CREATE FUNCTION can_overdraw_account(code VARCHAR(32), is_debtable TINYINT) RETURNS TINYINT
BEGIN
    /*账户可透支*/
    return LOCATE("_0", code) || is_debtable;
END
$$
DELIMITER ;

DROP TRIGGER IF EXISTS transfer_insert_trigger;
DELIMITER $$
CREATE TRIGGER transfer_insert_trigger AFTER INSERT ON ginger_finance.account_transfer FOR EACH ROW 
BEGIN
    DECLARE trigger_name VARCHAR(1024);
    DECLARE fn_result TINYINT;
    DECLARE source_account_code VARCHAR(128);
    DECLARE dest_account_code VARCHAR(128);
    DECLARE source_account_is_debtable TINYINT;
    
    SET trigger_name = "transfer_insert_trigger";
    SET fn_result = record_trigger_action(trigger_name, new.id, "triggered", "");
    
    IF new.source_amount < 0 || new.dest_amount < 0 THEN 
        SET fn_result = record_trigger_action(trigger_name, new.id, "failed", "交易金额不能为负");
        SIGNAL SQLSTATE 'HY000' SET MESSAGE_TEXT = "[mysql trigger]: 交易金额不能为负";
    END IF;
        
    /*先扣除源账户金额*/
    SELECT code, is_debtable into source_account_code, source_account_is_debtable
    FROM account_account WHERE id=new.source_account_id;

    IF !should_ignore_balance(source_account_code) THEN
        IF !can_overdraw_account(source_account_code, source_account_is_debtable) THEN
            UPDATE account_account SET balance=balance-new.source_amount, updated_at=NOW()
            WHERE id=new.source_account_id AND balance-frozen_amount>=new.source_amount;
        ELSE
            UPDATE account_account SET balance=balance-new.source_amount, updated_at=NOW()
            WHERE id=new.source_account_id;
        END IF;

        IF ROW_COUNT() <> 1 THEN
            SET fn_result = record_trigger_action(trigger_name, new.id, "failed", "账户余额不足");
            SIGNAL SQLSTATE 'HY000' SET MESSAGE_TEXT = "[mysql trigger]: 账户余额不足";
        END IF;
        SET fn_result = record_trigger_action(trigger_name, new.id, "done", "succeed");
    END IF;

    /*再增加目标账户金额*/
    SELECT code INTO dest_account_code FROM account_account WHERE id=new.dest_account_id;
    IF !should_ignore_balance(dest_account_code) THEN
        UPDATE account_account SET balance=balance+new.dest_amount, updated_at=NOW()
        WHERE id=new.dest_account_id;
        SET fn_result = record_trigger_action(trigger_name, new.id, "done", "succeed");
    END IF;
END
$$
DELIMITER ;