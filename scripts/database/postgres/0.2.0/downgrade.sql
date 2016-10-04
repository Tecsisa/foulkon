DO $$
    DECLARE COL_EXIST NUMERIC(10);
    BEGIN
        ------------------------
        -- ALTER TABLE groups --
        ------------------------
        -- Delete the updateAt column
        SELECT COUNT(column_name) INTO COL_EXIST
            FROM information_schema.columns
            WHERE table_name LIKE 'groups' AND column_name LIKE 'update_at';

        IF COL_EXIST = 1 THEN
            EXECUTE 'ALTER TABLE groups DROP COLUMN update_at';
            RAISE NOTICE '[INFO] Alter table groups to remove next column: update_at';
        ELSE
            RAISE NOTICE '[WARN] The groups column could not be removed. Check if this column exists';
        END IF;

        ------------------------
        -- ALTER TABLE users --
        ------------------------
        -- Delete the updateAt column
        SELECT COUNT(column_name) INTO COL_EXIST
            FROM information_schema.columns
            WHERE table_name LIKE 'users' AND column_name LIKE 'update_at';

        IF COL_EXIST = 1 THEN
            EXECUTE 'ALTER TABLE users DROP COLUMN update_at';
            RAISE NOTICE '[INFO] Alter table users to remove next column: update_at';
        ELSE
            RAISE NOTICE '[WARN] The users column could not be removed. Check if this column exists';
        END IF;

        --------------------------
        -- ALTER TABLE policies --
        --------------------------
        -- Delete the updateAt column
        SELECT COUNT(column_name) INTO COL_EXIST
            FROM information_schema.columns
            WHERE table_name LIKE 'policies' AND column_name LIKE 'update_at';

        IF COL_EXIST = 1 THEN
            EXECUTE 'ALTER TABLE policies DROP COLUMN update_at';
            RAISE NOTICE '[INFO] Alter table policies to remove next column: update_at';
        ELSE
            RAISE NOTICE '[WARN] The policies column could not be removed. Check if this column exists';
        END IF;

        ----------------------------------------
        -- ALTER TABLE group_policy_relations --
        ----------------------------------------
        -- Delete the createAt column
        SELECT COUNT(column_name) INTO COL_EXIST
            FROM information_schema.columns
            WHERE table_name LIKE 'group_policy_relations' AND column_name LIKE 'create_at';

        IF COL_EXIST = 1 THEN
            EXECUTE 'ALTER TABLE group_policy_relations DROP COLUMN create_at';
            RAISE NOTICE '[INFO] Alter table group_policy_relations to remove next column: create_at';
        ELSE
            RAISE NOTICE '[WARN] The group_policy_relations column could not be removed. Check if this column exists';
        END IF;

        --------------------------------------
        -- ALTER TABLE group_user_relations --
        --------------------------------------
        -- Delete the createAt column
        SELECT COUNT(column_name) INTO COL_EXIST
            FROM information_schema.columns
            WHERE table_name LIKE 'group_user_relations' AND column_name LIKE 'create_at';

        IF COL_EXIST = 1 THEN
            EXECUTE 'ALTER TABLE group_user_relations DROP COLUMN create_at';
            RAISE NOTICE '[INFO] Alter table group_user_relations to remove next column: create_at';
        ELSE
            RAISE NOTICE '[WARN] The group_user_relations column could not be removed. Check if this column exists';
        END IF;

    END $$
;