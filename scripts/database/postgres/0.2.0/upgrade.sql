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

        IF COL_EXIST = 0 THEN
            EXECUTE 'ALTER TABLE groups ADD update_at BIGINT';
            EXECUTE 'UPDATE groups SET update_at = create_at';
            EXECUTE 'ALTER TABLE groups ALTER COLUMN update_at SET NOT NULL';
            RAISE NOTICE '[INFO] Alter table groups to add next column: update_at';
        ELSE
            RAISE NOTICE '[WARN] The groups column could not be created. Check if this column already exists';
        END IF;

        ------------------------
        -- ALTER TABLE users --
        ------------------------
        -- Add the updateAt column
        SELECT COUNT(column_name) INTO COL_EXIST
            FROM information_schema.columns
            WHERE table_name LIKE 'users' AND column_name LIKE 'update_at';

        IF COL_EXIST = 0 THEN
            EXECUTE 'ALTER TABLE users ADD update_at BIGINT';
            EXECUTE 'UPDATE users SET update_at = create_at';
            EXECUTE 'ALTER TABLE users ALTER COLUMN update_at SET NOT NULL';
            RAISE NOTICE '[INFO] Alter table users to add next column: update_at';
        ELSE
            RAISE NOTICE '[WARN] The users column could not be created. Check if this column already exists';
        END IF;

        ------------------------
        -- ALTER TABLE policies --
        ------------------------
        -- Add the updateAt column
        SELECT COUNT(column_name) INTO COL_EXIST
            FROM information_schema.columns
            WHERE table_name LIKE 'policies' AND column_name LIKE 'update_at';

        IF COL_EXIST = 0 THEN
            EXECUTE 'ALTER TABLE policies ADD update_at BIGINT';
            EXECUTE 'UPDATE policies SET update_at = create_at';
            EXECUTE 'ALTER TABLE policies ALTER COLUMN update_at SET NOT NULL';
            RAISE NOTICE '[INFO] Alter table policies to add next column: update_at';
        ELSE
            RAISE NOTICE '[WARN] The policies column could not be created. Check if this column already exists';
        END IF;

    END $$
;