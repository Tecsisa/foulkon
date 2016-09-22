DO $$
    DECLARE COL_EXIST NUMERIC(10);
    BEGIN
        ------------------------
        -- ALTER TABLE groups --
        ------------------------
        -- Add the updateAt column
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
        -- Add the updateAt column
        SELECT COUNT(column_name) INTO COL_EXIST
            FROM information_schema.columns
            WHERE table_name LIKE 'users' AND column_name LIKE 'update_at';

        IF COL_EXIST = 1 THEN
            EXECUTE 'ALTER TABLE users DROP COLUMN update_at';
            RAISE NOTICE '[INFO] Alter table users to remove next column: update_at';
        ELSE
            RAISE NOTICE '[WARN] The users column could not be removed. Check if this column exists';
        END IF;

        ------------------------
        -- ALTER TABLE policies --
        ------------------------
        -- Add the updateAt column
        SELECT COUNT(column_name) INTO COL_EXIST
            FROM information_schema.columns
            WHERE table_name LIKE 'policies' AND column_name LIKE 'update_at';

        IF COL_EXIST = 1 THEN
            EXECUTE 'ALTER TABLE policies DROP COLUMN update_at';
            RAISE NOTICE '[INFO] Alter table policies to remove next column: update_at';
        ELSE
            RAISE NOTICE '[WARN] The policies column could not be removed. Check if this column exists';
        END IF;

    END $$
;