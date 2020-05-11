/*
    Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software distributed under
    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
    either express or implied.
    See the License for the specific language governing permissions and limitations under the License.

    Contributors to this project, hereby assign copyright in this code to the project,
    to be licensed under the same terms as the rest of the code.

*/
DO $$
  BEGIN
    /*
      ox_set_version(...)
      Inserts a new record in the version control table.
     */
    CREATE OR REPLACE FUNCTION ox_set_version(
      application_version_param CHARACTER VARYING(25),
      database_version_param    CHARACTER VARYING(25),
      description_param         TEXT,
      scripts_source_param      CHARACTER VARYING(250)
    )
    RETURNS VOID
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
    AS
    $BODY$
    BEGIN
      INSERT INTO version (
         application_version,
         database_version,
         description,
         time,
         scripts_source
      )
      VALUES (
         application_version_param,
         database_version_param,
         description_param,
         current_timestamp,
         scripts_source_param
      );
    END;
    $BODY$;

    /*
      ox_set_partition(...)
      Inserts a new or updates an existing partition.
     */
    CREATE OR REPLACE FUNCTION ox_set_partition(
      key_param character varying,
      name_param character varying,
      description_param text,
      local_version_param bigint,
      changed_by_param character varying,
      role_key_param character varying[]
    )
      RETURNS TABLE(result char(1))
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
    AS
    $BODY$
    DECLARE
      result          char(1); -- the result status for the upsert
      current_version bigint; -- the version of the row before the update or null if no row
      rows_affected   integer;
    BEGIN
      -- gets the current item type version
      SELECT version FROM partition WHERE key = key_param INTO current_version;

      -- checks the role can modify this role
      PERFORM ox_can_manage_partition(role_key_param);

      IF (current_version IS NULL) THEN
        INSERT INTO partition (
          id,
          key,
          name,
          description,
          version,
          created,
          updated,
          changed_by,
          owner
        )
        VALUES (
           nextval('partition_id_seq'),
           key_param,
           name_param,
           description_param,
           1,
           current_timestamp,
           null,
           changed_by_param,
           role_key_param
        );
        result := 'I';
      ELSE
        UPDATE partition
        SET name         = name_param,
            description  = description_param,
            version      = version + 1,
            updated      = current_timestamp,
            changed_by   = changed_by_param
        WHERE key = key_param
          -- concurrency management - optimistic locking
          AND (local_version_param = current_version OR local_version_param IS NULL OR local_version_param = 0)
          AND (
              name != name_param OR
              description != description_param
          );
        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        SELECT ox_get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
      END IF;
      RETURN QUERY SELECT result;
    END;
    $BODY$;

    ALTER FUNCTION ox_set_partition(
        character varying, -- key
        character varying, -- name
        text, -- description
        bigint, -- client version
        character varying, -- changed by
        role_key_param character varying[]
      )
      OWNER TO onix;

    /*
      ox_set_role(...)
      Inserts a new or updates an existing role.
    */
    CREATE OR REPLACE FUNCTION ox_set_role(
      key_param character varying,
      name_param character varying,
      description_param text,
      role_level_param integer,
      local_version_param bigint,
      changed_by_param character varying,
      role_key_param character varying[]
    )
      RETURNS TABLE(result char(1))
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
    AS
    $BODY$
    DECLARE
      result          char(1); -- the result status for the upsert
      current_version bigint; -- the version of the row before the update or null if no row
      rows_affected   integer;
    BEGIN
      -- gets the current item type version
      SELECT version FROM role WHERE key = key_param INTO current_version;

      -- checks the role can modify this role
      PERFORM ox_can_manage_partition(role_key_param);

      IF (current_version IS NULL) THEN
        INSERT INTO role (
          id,
          key,
          name,
          description,
          version,
          created,
          updated,
          changed_by,
          owner,
          level
        )
        VALUES (
           nextval('role_id_seq'),
           key_param,
           name_param,
           description_param,
           1,
           current_timestamp,
           null,
           changed_by_param,
           role_key_param,
           role_level_param
        );
        result := 'I';
      ELSE
        UPDATE role
        SET name         = name_param,
            description  = description_param,
            version      = version + 1,
            updated      = current_timestamp,
            changed_by   = changed_by_param
        WHERE key = key_param
          -- concurrency management - optimistic locking
          AND (local_version_param = current_version OR local_version_param IS NULL OR local_version_param = 0)
          AND (
              name != name_param OR
              description != description_param OR
              level != role_level_param
          );
        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        SELECT ox_get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
      END IF;
      RETURN QUERY SELECT result;
    END;
    $BODY$;

    ALTER FUNCTION ox_set_role(
        character varying, -- key
        character varying, -- name
        text, -- description
        integer, -- role level
        bigint, -- client version
        character varying, -- changed by
        character varying[] -- role keys
      )
      OWNER TO onix;

    /*
      ox_set_user(...)
      Inserts a new or updates an existing user.
    */
    CREATE OR REPLACE FUNCTION ox_set_user(
        key_param character varying,
        name_param character varying,
        email_param character varying,
        pwd_param character varying,
        salt_param character varying,
        expires_param timestamp(6) with time zone,
        local_version_param bigint,
        changed_by_param character varying,
        role_key_param character varying[])
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
    AS
    $BODY$
    DECLARE
        result          char(1); -- the result status for the upsert
        current_version bigint; -- the version of the row before the update or null if no row
        rows_affected   integer;
        new_salt        character varying;
    BEGIN
        -- only users in level 2 roles can create or update other users
        -- if not super admin then raise exception
        PERFORM ox_is_super_admin(role_key_param, TRUE);

        -- gets the current user version
        SELECT version FROM "user" WHERE key = key_param INTO current_version;

        IF (current_version IS NULL) THEN
            INSERT INTO "user" (
                id,
                key,
                name,
                email,
                pwd,
                salt,
                expires,
                version,
                created,
                updated,
                changed_by
            )
            VALUES (
               nextval('user_id_seq'),
               key_param,
               name_param,
               email_param,
               pwd_param,
               salt_param,
               expires_param,
               1,
               current_timestamp,
               null,
               changed_by_param
            );
            result := 'I';
        ELSE
            -- NOTE: if a password has been provided, even if it is the same as the originally
            -- stored in the database, it would have got here encrypted with a new randomly generated
            -- salt and therefore, it would look different to the database server
            -- so it would get updated and would look different both pwd and salt in the database
            IF (pwd_param IS NOT NULL) THEN
                -- has to update the salt otherwise the app will not be able to authenticate
                -- the new pwd in the future
                new_salt = salt_param;
            END IF;

            UPDATE "user"
            SET name         = name_param,
                email        = email_param,
                pwd          = COALESCE(pwd_param, pwd),  -- if the passed-in password is NULL, then do not change it
                salt         = COALESCE(new_salt, salt),  -- if new_salt is NOT NULL, the update the salt
                expires     = expires_param,
                version      = version + 1,
                updated      = current_timestamp,
                changed_by   = changed_by_param
            WHERE key = key_param
              -- concurrency management - optimistic locking
              AND (local_version_param = current_version OR local_version_param IS NULL OR local_version_param = 0)
              AND (
                    name != name_param OR
                    email != email_param OR
                    expires != expires_param OR
                    pwd != pwd_param AND pwd_param IS NOT NULL
                );
            GET DIAGNOSTICS rows_affected := ROW_COUNT;
            SELECT ox_get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
        END IF;
        RETURN QUERY SELECT result;
    END;
    $BODY$;

    ALTER FUNCTION ox_set_user(
            character varying, -- key
            character varying, -- name
            character varying, -- email
            character varying, -- pwd
            character varying, -- salt
            timestamp(6) with time zone, -- expires
            bigint, -- client version
            character varying, -- changed by
            character varying[] -- role keys
        )
        OWNER TO onix;

    /*
      ox_add_membership(...)
      creates a new membership.
    */
    CREATE OR REPLACE FUNCTION ox_add_membership(
        key_param character varying,
        user_key_param character varying,
        role_key_param character varying,
        changed_by_param character varying,
        logged_role_key_param character varying[])
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
    AS
    $BODY$
    DECLARE
        current_version bigint; -- the version of the row before the update or null if no row
        role_id_value   bigint;
        user_id_value   bigint;
    BEGIN
        -- only users in level 2 roles can create or update other users
        -- if not super admin then raise exception
        PERFORM ox_is_super_admin(logged_role_key_param, TRUE);

        -- gets the current membership version
        SELECT version FROM membership WHERE key = key_param INTO current_version;

        -- gets the role id from the passed-in key
        SELECT id FROM role WHERE key = role_key_param INTO role_id_value;

        -- if role Id not found raise exception
        IF (role_id_value IS NULL) THEN
            RAISE EXCEPTION 'Role with key ''%'' has not been found.', role_key_param
                USING HINT = 'Check the role you specified exists.';
        END IF;

        -- gets the user id from the passed-in key
        SELECT id FROM "user" WHERE key = user_key_param INTO user_id_value;

        -- if user Id not found raise exception
        IF (user_id_value IS NULL) THEN
            RAISE EXCEPTION 'User with key ''%'' has not been found.', user_key_param
                USING HINT = 'Check the user you specified exists.';
        END IF;

        IF (current_version IS NULL) THEN
            INSERT INTO membership (
                id,
                key,
                role_id,
                user_id,
                version,
                created,
                updated,
                changed_by
            )
            VALUES (
               nextval('membership_id_seq'),
               key_param,
               role_id_value,
               user_id_value,
               1,
               current_timestamp,
               null,
               changed_by_param
            );
        END IF;
        RETURN QUERY SELECT 'I'::char(1);
    END;
    $BODY$;

    ALTER FUNCTION ox_add_membership(
        character varying, -- key
        character varying, -- user_key_param
        character varying, -- role_key_param
        character varying, -- changed_by_param
        character varying[] -- role_key_param
        )
        OWNER TO onix;

    /*
      ox_set_model(...)
      Inserts a new or updates an existing meta model.
     */
    CREATE OR REPLACE FUNCTION ox_set_model(
         key_param character varying,
         name_param character varying,
         description_param text,
         managed_param boolean,
         local_version_param bigint,
         changed_by_param character varying,
         partition_key_param character varying,
         role_key_param character varying[]
      )
      RETURNS TABLE(result char(1))
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
    AS
    $BODY$
    DECLARE
      result             char(1); -- the result status for the upsert
      current_version    bigint; -- the version of the row before the update or null if no row
      rows_affected      integer;
      partition_id_value bigint;
    BEGIN
      -- if there is no partition key, then use the REF partition
      IF (partition_key_param IS NULL OR partition_key_param = '') THEN
        partition_key_param = 'REF';
      END IF;

      -- gets the current model version
      SELECT version FROM model WHERE key = key_param INTO current_version;

      IF (current_version IS NULL) THEN
        -- as no model exists yet, it finds the partition associated with the role
        SELECT p.id
        FROM partition p
           INNER JOIN privilege pr on p.id = pr.partition_id
           INNER JOIN role r on pr.role_id = r.id
        AND pr.can_create = TRUE -- has create permission
        AND r.key = ANY(role_key_param) -- the user role
        AND p.key = partition_key_param -- the requested partition
        LIMIT 1
           INTO partition_id_value;

        IF (partition_id_value IS NULL) THEN
          RAISE EXCEPTION 'Role % is not authorised to create a Model on partition %.', role_key_param, partition_key_param
            USING HINT = 'The role needs to be granted CREATE privilege on the specified partition, or a different role or partition should be used instead.';
        END IF;

        INSERT INTO model (
           id,
           key,
           name,
           description,
           managed,
           version,
           created,
           updated,
           changed_by,
           partition_id
        )
        VALUES (
          nextval('role_id_seq'),
          key_param,
          name_param,
          description_param,
          managed_param,
          1,
          current_timestamp,
          null,
          changed_by_param,
          partition_id_value
        );
        result := 'I';
      ELSE
        -- a model exists therefore, it finds the partition associated with the model
        -- and check the role has create / update rights
        -- NOTE: the existing partition is used to determine rights (instead of the passed in partition)
        SELECT p.id
        FROM partition p
          INNER JOIN privilege pr on p.id = pr.partition_id
          INNER JOIN role r on pr.role_id = r.id
          INNER JOIN model m on p.id = m.partition_id
          AND pr.can_create = TRUE -- whether the role has update permission on the model
          AND r.key = ANY(role_key_param) -- the user role requesting the update
          AND m.key = key_param -- the model to be updated
          LIMIT 1
             INTO partition_id_value;

        IF (partition_id_value IS NULL) THEN
          RAISE EXCEPTION 'Role % is not authorised to update the Model on partition %.', role_key_param, partition_key_param
            USING HINT = 'The role needs to be granted CREATE privilege on the specified partition, or a different role or partition should be used instead.';
        END IF;

        UPDATE model
        SET name         = name_param,
            description  = description_param,
            managed      = managed_param,
            version      = version + 1,
            updated      = current_timestamp,
            changed_by   = changed_by_param
        WHERE key = key_param
          -- concurrency management - optimistic locking
          AND (local_version_param = current_version OR local_version_param IS NULL OR local_version_param = 0)
          AND (
            name != name_param OR
            description != description_param
          );
        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        SELECT ox_get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
      END IF;
      RETURN QUERY SELECT result;
    END;
    $BODY$;

    ALTER FUNCTION ox_set_model(
        character varying, -- key
        character varying, -- name
        text, -- description
        boolean, -- managed
        bigint, -- client version
        character varying, -- changed by
        partition_key_param character varying,
        role_key_param character varying[]
      )
      OWNER TO onix;

    /*
      ox_set_item(...)
      Inserts a new or updates an existing item.
      Concurrency Management:
       - If the item is found in the database, the function attempts an update of the existing record.
          In this case, if a null regex is passed as local_version_param, no optimistic locking is performed.
          If a regex is specified for local_version_param, the update is only performed if and only if the version in the database matches the passed in version.
       - If the item is not found in the database, then the local_version_param is ignored and a record with version 1 is inserted.
     */
    CREATE OR REPLACE FUNCTION ox_set_item(
        key_param character varying,
        name_param character varying,
        description_param text,
        meta_param jsonb,
        meta_enc_param boolean,
        txt_param text,
        txt_enc_param boolean,
        enc_key_ix_param smallint,
        tag_param text[],
        attribute_param hstore,
        status_param smallint,
        item_type_key_param character varying,
        local_version_param bigint,
        changed_by_param character varying,
        partition_key_param character varying,
        role_key_param character varying[]
      )
      RETURNS TABLE(result char(1))
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
    AS
    $BODY$
    DECLARE
      result                     char(1); -- the result status for the upsert
      current_version            bigint; -- the version of the row before the update or null if no row
      rows_affected              integer;
      item_type_id_value         integer;
      meta_schema_value          jsonb;
      is_meta_valid              boolean;
      partition_id_value         bigint;
      current_partition_id_value bigint;
      current_partition_key      character varying; -- the key of the item partition if the item exists (update case)
    BEGIN
      -- if a partition is not specified
      IF (partition_key_param IS NULL OR partition_key_param = '') THEN
        -- it defaults to the INS (instance) default partition
        partition_key_param = 'INS';
      END IF;

      -- find the item type surrogate key from the provided natural key
      SELECT id FROM item_type WHERE key = item_type_key_param INTO item_type_id_value;
      IF (item_type_id_value IS NULL) THEN
        -- the provided natural key is not in the item type table, cannot proceed
        RAISE EXCEPTION 'Item Type Key --> % not found.', item_type_key_param
          USING hint = 'Check an Item Type with the key exist in the database.';
      END IF;

      -- checks that the attributes passed in comply with the validation in the item_type
      PERFORM ox_check_item_attr(item_type_key_param, attribute_param);

      -- checks that the meta field complies with the json schema defined by the item type
      IF (meta_param IS NOT NULL) THEN
        SELECT meta_schema FROM item_type it WHERE it.key = item_type_key_param INTO meta_schema_value;
        IF (meta_schema_value IS NOT NULL) THEN
          SELECT ox_validate_json_schema(meta_schema_value, meta_param) INTO is_meta_valid;
          IF (NOT is_meta_valid) THEN
            RAISE EXCEPTION 'Meta field % for Item % is not valid as defined by the schema % in its type %.', meta_param, key_param, meta_schema_value, item_type_key_param
              USING hint = 'Check the JSON value meets the requirement of the schema defined by the item type.';
          END IF;
        END IF;
      END IF;

      -- get the item current version
      SELECT version FROM item WHERE key = key_param INTO current_version;

      -- if no version is found then go for an insert
      IF (current_version IS NULL) THEN
        -- finds the partition id for the specified key and role / privilege
        SELECT p.id
        FROM partition p
        INNER JOIN privilege pr on p.id = pr.partition_id
        INNER JOIN role r on pr.role_id = r.id
          AND p.key = partition_key_param -- the requested partition key
          AND pr.can_create = TRUE -- has create permission
          AND r.key = ANY(role_key_param) -- the user role
        LIMIT 1
             INTO partition_id_value;

        IF (partition_id_value IS NULL) THEN
          RAISE EXCEPTION 'Role % is not authorised to create Item % in partition %.', role_key_param, key_param, partition_key_param
            USING hint = 'The role needs to be granted CREATE privilege or a new role should be used instead.';
        END IF;

        INSERT INTO item (
            id,
            key,
            name,
            description,
            meta,
            meta_enc,
            txt,
            txt_enc,
            enc_key_ix,
            tag,
            attribute,
            status,
            item_type_id,
            version,
            created,
            updated,
            changed_by,
            partition_id
        )
        VALUES (
            nextval('item_id_seq'),
            key_param,
            name_param,
            description_param,
            meta_param,
            meta_enc_param,
            txt_param,
            txt_enc_param,
            enc_key_ix_param,
            tag_param,
            attribute_param,
            status_param,
            item_type_id_value,
            1,
            current_timestamp,
            null,
            changed_by_param,
            partition_id_value
        );
        result := 'I';
      ELSE
        -- check the role specified partition is the same as the item partition
        SELECT p.key
        FROM item i
        INNER JOIN partition p ON i.partition_id = p.id
        WHERE i.key = key_param
        INTO current_partition_key;

        IF (current_partition_key != partition_key_param) THEN
            RAISE EXCEPTION 'Item % exist within Partition % but a different partition % has been specified.', key_param, current_partition_key, partition_key_param
                USING hint = 'The item partition cannot be changed, consider posting changes under the current partition.';
        END IF;

        -- checks the role has privilege on the current partition
        SELECT i.partition_id
        FROM partition p
        INNER JOIN privilege pr on p.id = pr.partition_id
        INNER JOIN role r on pr.role_id = r.id
        INNER JOIN item i on p.id = i.partition_id
          AND i.partition_id = p.id -- the partition associated to the existing item
          AND pr.can_create = TRUE -- has create permission
          AND r.key = ANY(role_key_param) -- the user role
        LIMIT 1
             INTO current_partition_id_value;

        IF (current_partition_id_value IS NULL) THEN
          RAISE EXCEPTION 'Role % is not authorised to update Item % on the existing partition.', role_key_param, key_param
            USING hint = 'The role needs to be granted CREATE privilege or a new role should be used instead.';
        END IF;

        -- if a partition has been specified
        IF (partition_key_param IS NOT NULL) THEN
          -- checks the role has privilege on the passed-in partition
          SELECT p.id
          FROM partition p
               INNER JOIN privilege pr on p.id = pr.partition_id
               INNER JOIN role r on pr.role_id = r.id
               INNER JOIN item i on p.id = i.partition_id
            AND p.key = partition_key_param -- the passed in partition
            AND pr.can_create = TRUE -- has create permission
            AND r.key = ANY(role_key_param) -- the user role
               INTO partition_id_value;

          -- if the specified partition does not have privilege then raises error
          IF (partition_id_value IS NULL) THEN
            RAISE EXCEPTION 'Role % is not authorised to update Item % on partition %.', role_key_param, key_param, partition_key_param
              USING hint = 'The role needs to be granted CREATE privilege or a new role should be used instead.';
          END IF;
        END IF;

        -- if a version is found, go for an update
        UPDATE item
        SET name         = name_param,
            description  = description_param,
            meta         = meta_param,
            meta_enc     = meta_enc_param,
            txt          = txt_param,
            txt_enc      = txt_enc_param,
            enc_key_ix   = enc_key_ix_param,
            tag          = tag_param,
            attribute    = attribute_param,
            status       = status_param,
            item_type_id = item_type_id_value,
            version      = version + 1,
            updated      = current_timestamp,
            changed_by   = changed_by_param
        WHERE key = key_param
          -- the database record has not been modified by someone else
          -- if a null regex is passed as local version then it does not perform optimistic locking
          AND (local_version_param = current_version OR local_version_param IS NULL OR local_version_param = 0)
          AND (
          -- the fields to be updated have not changed
            name != name_param OR
            description != description_param OR
            status != status_param OR
            item_type_id != item_type_id_value OR
            meta != meta_param OR
            meta_enc != meta_enc_param OR
            txt != txt_param OR
            txt_enc != txt_enc_param OR
            enc_key_ix != enc_key_ix_param OR
            tag != tag_param OR
            avals(attribute) != avals(attribute_param)
          );
        -- determines if the update has gone ahead
        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        -- works out the update status
        SELECT ox_get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
      END IF;
      RETURN QUERY SELECT result;
    END;
    $BODY$;

    ALTER FUNCTION ox_set_item(
        character varying,
        character varying,
        text,
        jsonb, -- meta
        boolean, -- meta_enc
        text, -- txt
        boolean, -- txt_enc
        smallint, -- enc_key_ix
        text[],
        hstore,
        smallint,
        character varying,
        bigint,
        character varying,
        character varying, -- partition_key_param
        character varying[] -- role_key_param
      )
      OWNER TO onix;

    /*
      ox_set_item_type(...)
      Inserts a new or updates an existing item item.
      Concurrency Management:
       - If the item type is found in the database, the function attempts an update of the existing record.
          In this case, if a null regex is passed as local_version_param, no optimistic locking is performed.
          If a regex is specified for local_version_param, the update is only performed if and only if the version in the database matches the passed in version.
       - If the item type is not found in the database, then the local_version_param is ignored and a record with version 1 is inserted.
     */
    CREATE OR REPLACE FUNCTION ox_set_item_type(
        key_param character varying,
        name_param character varying,
        description_param text,
        filter_param jsonb,
        meta_schema_param jsonb,
        local_version_param bigint,
        model_key_param character varying,
        changed_by_param character varying,
        notify_change_param char,
        tag_param text[],
        encrypt_meta_param boolean,
        encrypt_txt_param boolean,
        style_param jsonb,
        role_key_param character varying[]
      )
      RETURNS TABLE(result char(1))
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
    AS
    $BODY$
    DECLARE
      result             char(1); -- the result status for the upsert
      current_version    bigint; -- the version of the row before the update or null if no row
      rows_affected      integer;
      model_id_value     integer;
      partition_id_value bigint;
    BEGIN
      -- checks a model has been specified
      IF (model_key_param IS NULL) THEN
        RAISE EXCEPTION 'Model not specified when trying to set Item Type with key %', key_param;
      END IF;

      -- gets the model id associated with the model key
      SELECT m.id FROM model m WHERE m.key = model_key_param INTO model_id_value;

      IF (model_id_value IS NULL) THEN
        RAISE EXCEPTION 'Model % not found.', model_key_param
          USING hint = 'Check a Model with the specified key exist in the database.';
      END IF;

      -- finds the partition associated with the model
      -- for the item type that has create rights for the specified role
      SELECT p.id
      FROM partition p
      INNER JOIN model m on p.id = m.partition_id
      INNER JOIN privilege pr on p.id = pr.partition_id
      INNER JOIN role r on pr.role_id = r.id
        AND pr.can_create = TRUE -- has create permission
        AND r.key = ANY(role_key_param) -- the user role
        AND m.key = model_key_param -- the model the item type is in
      LIMIT 1
      INTO partition_id_value;

      IF (partition_id_value IS NULL) THEN
        RAISE EXCEPTION 'Role % is not authorised to create Item Type %.', role_key_param, key_param
          USING hint = 'The role needs to be granted CREATE privilege or a new role should be used instead.';
      END IF;

      -- gets the current item type version
      SELECT version FROM item_type WHERE key = key_param INTO current_version;

      IF (current_version IS NULL) THEN
        INSERT INTO item_type (
          id,
          key,
          name,
          description,
          filter,
          meta_schema,
          version,
          created,
          updated,
          changed_by,
          model_id,
          notify_change,
          tag,
          encrypt_meta,
          encrypt_txt,
          style
        )
        VALUES (
          nextval('item_type_id_seq'),
          key_param,
          name_param,
          description_param,
          filter_param,
          meta_schema_param,
          1,
          current_timestamp,
          null,
          changed_by_param,
          model_id_value,
          notify_change_param,
          tag_param,
          encrypt_meta_param,
          encrypt_txt_param,
          style_param
        );
        result := 'I';
      ELSE
        UPDATE item_type
        SET name          = name_param,
            description   = description_param,
            filter        = filter_param,
            meta_schema   = meta_schema_param,
            version       = version + 1,
            updated       = current_timestamp,
            changed_by    = changed_by_param,
            model_id      = model_id_value,
            notify_change = notify_change_param,
            tag           = tag_param,
            encrypt_meta  = encrypt_meta_param,
            encrypt_txt   = encrypt_txt_param,
            style         = style_param
        WHERE key = key_param
          -- concurrency management - optimistic locking
          AND (local_version_param = current_version OR local_version_param IS NULL OR local_version_param = 0)
          AND (
                name != name_param OR
                description != description_param OR
                filter != filter_param OR
                meta_schema != meta_schema_param OR
                model_id != model_id_value OR
                notify_change != notify_change_param OR
                tag != tag_param OR
                encrypt_meta != encrypt_meta_param OR
                encrypt_txt != encrypt_txt_param OR
                style != style_param
          );
        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        SELECT ox_get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
      END IF;
      RETURN QUERY SELECT result;
    END;
    $BODY$;

    ALTER FUNCTION ox_set_item_type(
        character varying, -- key
        character varying, -- name
        text, -- description
        jsonb, -- meta query filter
        jsonb, -- meta json schema
        bigint, -- client version
        character varying, -- meta model key
        character varying, -- changed by
        char, -- notify change
        text[], -- tag
        boolean, -- encrypt meta
        boolean, -- encrypt txt
        jsonb, -- style
        character varying[] -- role_key_param
      )
      OWNER TO onix;

    /*
      ox_set_link_type(...)
      Inserts a new or updates an existing link type.
      Concurrency Management:
       - If the link type is found in the database, the function attempts an update of the existing record.
          In this case, if a null regex is passed as local_version_param, no optimistic locking is performed.
          If a regex is specified for local_version_param, the update is only performed if and only if the version in the database matches the passed in version.
       - If the link type is not found in the database, then the local_version_param is ignored and a record with version 1 is inserted.
     */
    CREATE OR REPLACE FUNCTION ox_set_link_type(
        key_param character varying,
        name_param character varying,
        description_param text,
        meta_schema_param jsonb,
        tag_param text[],
        encrypt_meta_param boolean,
        encrypt_txt_param boolean,
        style_param jsonb,
        local_version_param bigint,
        model_key_param character varying,
        changed_by_param character varying,
        role_key_param character varying[]
      )
      RETURNS TABLE(result char(1))
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
    AS
    $BODY$
    DECLARE
      result          char(1); -- the result status for the upsert
      current_version bigint; -- the version of the row before the update or null if no row
      rows_affected   integer;
      model_id_value  integer;
    BEGIN
      -- checks a model has been specified
      IF (model_key_param IS NULL) THEN
        RAISE EXCEPTION 'Meta Model not specified when trying to set Link Type with key %', key_param;
      END IF;

      -- gets the model id associated with the model key
      SELECT m.id FROM model m WHERE m.key = model_key_param INTO model_id_value;

      IF (model_id_value IS NULL) THEN
        RAISE EXCEPTION 'Meta Model % not found.', model_key_param
          USING hint = 'Check a meta model with the specified key exist in the database.';
      END IF;

      -- finds the partition associated with the model
      -- for the link type that has create rights for the specified role
      SELECT COUNT(p.id)
      FROM partition p
      INNER JOIN model m on p.id = m.partition_id
      INNER JOIN privilege pr on p.id = pr.partition_id
      INNER JOIN role r on pr.role_id = r.id
        AND pr.can_create = TRUE -- has create permission
        AND r.key = ANY(role_key_param) -- the user role
        AND m.key = model_key_param -- the model the item type is in
      INTO rows_affected;

      IF (rows_affected = 0) THEN
        RAISE EXCEPTION 'Role % is not authorised to create Link Type %.', role_key_param, key_param
          USING hint = 'The role needs to be granted CREATE privilege or a new role should be used instead.';
      END IF;

      -- gets the link type current version
      SELECT version FROM link_type WHERE key = key_param INTO current_version;

      IF (current_version IS NULL) THEN
        INSERT INTO link_type (
           id,
           key,
           name,
           description,
           meta_schema,
           tag,
           encrypt_meta,
           encrypt_txt,
           style,
           version,
           created,
           updated,
           changed_by,
           model_id)
        VALUES (nextval('link_type_id_seq'),
                key_param,
                name_param,
                description_param,
                meta_schema_param,
                tag_param,
                encrypt_meta_param,
                encrypt_txt_param,
                style_param,
                1,
                current_timestamp,
                null,
                changed_by_param,
                model_id_value);
        result := 'I';
      ELSE
        UPDATE link_type
        SET name        = name_param,
            description = description_param,
            meta_schema = meta_schema_param,
            tag         = tag_param,
            encrypt_meta= encrypt_meta_param,
            encrypt_txt = encrypt_txt_param,
            style       = style_param,
            version     = version + 1,
            updated     = current_timestamp,
            changed_by  = changed_by_param,
            model_id    = model_id_value
        WHERE key = key_param
          -- concurrency management - optimistic locking
          AND (local_version_param = current_version OR local_version_param IS NULL OR local_version_param = 0)
          AND (
            name != name_param OR
            description != description_param OR
            meta_schema != meta_schema_param OR
            tag != tag_param OR
            encrypt_meta != encrypt_meta_param OR
            encrypt_txt != encrypt_txt_param OR
            style != style_param OR
            model_id != model_id_value
          );
        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        SELECT ox_get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
      END IF;
      RETURN QUERY SELECT result;
    END;
    $BODY$;

    ALTER FUNCTION ox_set_link_type(
      character varying, -- key
      character varying, -- name
      text, -- description
      jsonb, -- meta json schema validation
      text[], -- tag
      boolean, -- encrypt meta
      boolean, -- encrypt txt
      jsonb, -- style
      bigint, -- client version
      character varying, -- meta model key
      character varying, -- changed by
      character varying[] -- role_key_param
      )
      OWNER TO onix;

    /*
      ox_set_link(...)
      Inserts a new or updates an existing link.
      Concurrency Management:
       - If the link is found in the database, the function attempts an update of the existing record.
          In this case, if a null regex is passed as local_version_param, no optimistic locking is performed.
          If a regex is specified for local_version_param, the update is only performed if and only if the version in the database matches the passed in version.
       - If the link is not found in the database, then the local_version_param is ignored and a record with version 1 is inserted.
     */
    CREATE OR REPLACE FUNCTION ox_set_link(
      key_param character varying,
      link_type_key_param character varying,
      start_item_key_param character varying,
      end_item_key_param character varying,
      description_param text,
      meta_param jsonb,
      meta_enc_param boolean,
      txt_param text,
      txt_enc_param boolean,
      enc_key_ix_param smallint,
      tag_param text[],
      attribute_param hstore,
      local_version_param bigint,
      changed_by_param character varying,
      role_key_param character varying[]
    )
      RETURNS TABLE(result char(1))
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
    AS
    $BODY$
    DECLARE
      result                    char(1); -- the result status for the upsert
      current_version           bigint; -- the version of the row before the update or null if no row
      rows_affected             integer;
      start_item_id_value       bigint;
      end_item_id_value         bigint;
      link_type_id_value        integer;
      start_item_type_key_value character varying;
      end_item_type_key_value   character varying;
      meta_schema_value         jsonb;
      is_meta_valid             boolean;
    BEGIN
      -- find the link type surrogate key from the provided natural key
      SELECT id FROM link_type WHERE key = link_type_key_param INTO link_type_id_value;
      IF (link_type_id_value IS NULL) THEN
        -- the provided natural key is not in the link type table, cannot proceed
        RAISE EXCEPTION 'Link Type Key --> % not found.', link_type_key_param
          USING hint = 'Check a Link Type with the key exist in the database.';
      END IF;

      SELECT i.id, t.key
      FROM item i
      INNER JOIN item_type t
        ON i.item_type_id = t.id
      WHERE i.key = start_item_key_param INTO start_item_id_value, start_item_type_key_value;

      IF (start_item_id_value IS NULL) THEN
        -- the start item does not exist
        RAISE EXCEPTION 'Start item with key --> % does not exist.', start_item_key_param
          USING hint = 'Check an item with the specified key exist in the database.';
      END IF;

      SELECT i.id, t.key
      FROM item i
      INNER JOIN item_type t ON i.item_type_id = t.id
      WHERE i.key = end_item_key_param
        INTO end_item_id_value, end_item_type_key_value;

      IF (end_item_id_value IS NULL) THEN
        -- the end item does not exist
        RAISE EXCEPTION 'End item with key --> % does not exist.', end_item_key_param
          USING hint = 'Check an item with the specified key exist in the database.';
      END IF;

      -- checks that the link is allowed
      PERFORM ox_check_link(link_type_key_param, start_item_type_key_value, end_item_type_key_value);

      -- checks that the attributes passed in comply with the validation in the link_type
      PERFORM ox_check_link_attr(link_type_key_param, attribute_param);

      -- checks that the meta field complies with the json schema defined by the item type
      IF (meta_param IS NOT NULL) THEN
        SELECT meta_schema FROM link_type it WHERE it.key = link_type_key_param INTO meta_schema_value;
        IF (meta_schema_value IS NOT NULL) THEN
          SELECT ox_validate_json_schema(meta_schema_value, meta_param) INTO is_meta_valid;
          IF (NOT is_meta_valid) THEN
            RAISE EXCEPTION 'Meta field for Link % is not valid as defined by the schema in its type %.', key_param, link_type_key_param
              USING hint = 'Check the JSON value meets the requirement of the schema defined by the link type.';
          END IF;
        END IF;
      END IF;

      SELECT version FROM link WHERE key = key_param INTO current_version;
      IF (current_version IS NULL) THEN
        -- finds if the link can be created by the role
        SELECT COUNT(*)
        FROM link_type lt
          INNER JOIN model m ON lt.model_id = m.id
          INNER JOIN partition p ON m.partition_id = p.id
          INNER JOIN privilege pr ON p.id = pr.partition_id
          INNER JOIN role r ON pr.role_id = r.id
          WHERE r.key = ANY(role_key_param)
          AND pr.can_create = TRUE
          AND lt.key = link_type_key_param
          INTO rows_affected;

        IF (rows_affected = 0) THEN
          RAISE EXCEPTION 'The Role % is not authorised to create the Link % or specified type %.', role_key_param, key_param, link_type_key_param
            USING hint = 'The role needs to be granted CREATE privilege or a new role should be used instead.';
        END IF;

        INSERT INTO link (
          id,
          key,
          link_type_id,
          start_item_id,
          end_item_id,
          description,
          meta,
          meta_enc,
          txt,
          txt_enc,
          enc_key_ix,
          tag,
          attribute,
          version,
          created,
          updated,
          changed_by
        )
        VALUES (
          nextval('link_id_seq'),
          key_param,
          link_type_id_value,
          start_item_id_value,
          end_item_id_value,
          description_param,
          meta_param,
          meta_enc_param,
          txt_param,
          txt_enc_param,
          enc_key_ix_param,
          tag_param,
          attribute_param,
          1,
          current_timestamp,
          null,
          changed_by_param
        );
        result := 'I';
      ELSE
        -- finds if the link can be updated by the role - using the passed in link type
        SELECT COUNT(*)
        FROM link_type lt
           INNER JOIN model m ON lt.model_id = m.id
           INNER JOIN partition p ON m.partition_id = p.id
           INNER JOIN privilege pr ON p.id = pr.partition_id
           INNER JOIN role r ON pr.role_id = r.id
        WHERE r.key = ANY(role_key_param)
          AND pr.can_create = TRUE
          AND lt.key = link_type_key_param
          INTO rows_affected;

        IF (rows_affected = 0) THEN
          RAISE EXCEPTION 'The Role % is not authorised to create the Link % of specified type %.', role_key_param, key_param, link_type_key_param
            USING hint = 'The role needs to be granted CREATE privilege or a new role should be used instead.';
        END IF;

        -- finds if the link can be updated by the role - using the current link type this time
        SELECT COUNT(*)
        FROM link_type lt
           INNER JOIN model m ON lt.model_id = m.id
           INNER JOIN partition p ON m.partition_id = p.id
           INNER JOIN privilege pr ON p.id = pr.partition_id
           INNER JOIN role r ON pr.role_id = r.id
           INNER JOIN link l ON lt.id = l.link_type_id
        WHERE r.key = ANY(role_key_param)
          AND pr.can_create = TRUE
          AND lt.id = l.link_type_id
          AND l.key = key_param
          INTO rows_affected;

        IF (rows_affected = 0) THEN
          RAISE EXCEPTION 'The Role % is not authorised to create the Link % of current type %.', role_key_param, key_param, link_type_key_param
            USING hint = 'The role needs to be granted CREATE privilege or a new role should be used instead.';
        END IF;

        UPDATE link
        SET description   = description_param,
            meta          = meta_param,
            meta_enc      = meta_enc_param,
            txt           = txt_param,
            txt_enc       = txt_enc_param,
            enc_key_ix    = enc_key_ix_param,
            tag           = tag_param,
            attribute     = attribute_param,
            link_type_id  = link_type_id_value,
            start_item_id = start_item_id_value,
            end_item_id   = end_item_id_value,
            version       = version + 1,
            updated       = current_timestamp,
            changed_by    = changed_by_param
        WHERE key = key_param
          -- concurrency management - optimistic locking
          AND (local_version_param = current_version OR local_version_param IS NULL OR local_version_param = 0)
          AND (
            meta != meta_param OR
            meta_enc != meta_enc_param OR
            txt != txt_param OR
            txt_enc != txt_enc_param OR
            enc_key_ix != enc_key_ix_param OR
            description != description_param OR
            tag != tag_param OR
            attribute != attribute_param OR
            link_type_id != link_type_id_value OR
            start_item_id != start_item_id_value OR
            end_item_id != end_item_id_value
          );
        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        SELECT ox_get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
      END IF;
      RETURN QUERY SELECT result;
    END;
    $BODY$;

    ALTER FUNCTION ox_set_link(
        character varying,
        character varying,
        character varying,
        character varying,
        text,
        jsonb,
        boolean, -- meta_enc
        text,
        boolean, -- txt_enc
        smallint, -- enc_key
        text[],
        hstore,
        bigint,
        character varying,
        character varying[] -- role_key_param
      )
      OWNER TO onix;

    /*
      ox_set_link_rule(...)
      Inserts a new or updates an existing link rule.
      Concurrency Management:
       - If the link rule is found in the database, the function attempts an update of the existing record.
          In this case, if a null regex is passed as local_version_param, no optimistic locking is performed.
          If a regex is specified for local_version_param, the update is only performed if and only if the version in the database matches the passed in version.
       - If the link rule is not found in the database, then the local_version_param is ignored and a record with version 1 is inserted.
     */
    CREATE OR REPLACE FUNCTION ox_set_link_rule(
      key_param character varying,
      name_param character varying,
      description_param text,
      link_type_key_param character varying,
      start_item_type_key_param character varying,
      end_item_type_key_param character varying,
      local_version_param bigint,
      changed_by_param character varying,
      role_key_param character varying[]
    )
      RETURNS TABLE(result char(1))
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
    AS
    $BODY$
    DECLARE
      result                   char(1); -- the result status for the upsert
      current_version          bigint; -- the version of the row before the update or null if no row
      rows_affected            integer;
      link_type_id_value       integer;
      start_item_type_id_value integer;
      end_item_type_id_value   integer;
    BEGIN
      -- gets the current version
      SELECT version FROM link_rule WHERE key = key_param INTO current_version;

      -- gets the required id's
      SELECT id FROM link_type WHERE key = link_type_key_param INTO link_type_id_value;

      -- the link type must exist
      IF link_type_id_value IS NULL THEN
        RAISE EXCEPTION 'The specified link type "%" could not be found or does not exist.', link_type_key_param
          USING hint = 'A link rule needs a link type, make sure it exists before creating the rule.';
      END IF;

      SELECT id FROM item_type WHERE key = start_item_type_key_param INTO start_item_type_id_value;

      -- the start item type must exist
      IF start_item_type_id_value IS NULL THEN
        RAISE EXCEPTION 'The specified item type "%" for the start item could not be found or does not exist.', start_item_type_key_param
          USING hint = 'A link rule needs item types for start and end items the link is connecting, make sure they exists before creating the rule.';
      END IF;

      SELECT id FROM item_type WHERE key = end_item_type_key_param INTO end_item_type_id_value;

      -- the end item type must exist
      IF end_item_type_id_value IS NULL THEN
        RAISE EXCEPTION 'The specified item type "%" for the end item could not be found or does not exist.', end_item_type_key_param
          USING hint = 'A link rule needs item types for start and end items the link is connecting, make sure they exists before creating the rule.';
      END IF;

      IF (current_version IS NULL) THEN
        -- check if the role can create the link rule
        SELECT COUNT(*)
        FROM link_type lt
               INNER JOIN model m ON lt.model_id = m.id
               INNER JOIN partition p ON m.partition_id = p.id
               INNER JOIN privilege pr ON p.id = pr.partition_id
               INNER JOIN role r ON pr.role_id = r.id
        WHERE r.key = ANY(role_key_param)
          AND pr.can_create = TRUE
          AND lt.key = link_type_key_param
          INTO rows_affected;

        IF (rows_affected = 0) THEN
          RAISE EXCEPTION 'The Role % is not authorised to create the Link Rule %.', role_key_param, key_param
            USING hint = 'The role needs to be granted CREATE privilege or a new role should be used instead.';
        END IF;

        INSERT INTO link_rule (
           id,
           key,
           name,
           description,
           link_type_id,
           start_item_type_id,
           end_item_type_id,
           version,
           created,
           updated,
           changed_by
        )
        VALUES (
          nextval('link_rule_id_seq'),
          key_param,
          name_param,
          description_param,
          link_type_id_value,
          start_item_type_id_value,
          end_item_type_id_value,
          1,
          current_timestamp,
          null,
          changed_by_param
        );
        result := 'I';
      ELSE
        -- checks if the existing link rule can be updated by the specified role
        SELECT COUNT(*)
        FROM link_rule lr
               INNER JOIN link_type lt ON lr.link_type_id = lt.id
               INNER JOIN model m ON lt.model_id = m.id
               INNER JOIN partition p on m.partition_id = p.id
               INNER JOIN privilege pr on p.id = pr.partition_id
               INNER JOIN role r on pr.role_id = r.id
        WHERE r.key = ANY(role_key_param)
          AND pr.can_create = TRUE
          AND lr.key = key_param
          AND lt.key = link_type_key_param
          INTO rows_affected;

        IF (rows_affected = 0) THEN
          RAISE EXCEPTION 'The Role % is not authorised to update the Link Rule %.', role_key_param, key_param
            USING hint = 'The role needs to be granted CREATE privilege or a new role should be used instead.';
        END IF;

        UPDATE link_rule
        SET name               = name_param,
            description        = description_param,
            link_type_id       = link_type_id_value,
            start_item_type_id = start_item_type_id_value,
            end_item_type_id   = end_item_type_id_value,
            version            = version + 1,
            updated            = current_timestamp,
            changed_by         = changed_by_param
        WHERE key = key_param
          -- concurrency management - optimistic locking (disabled if local_version_param is null)
          AND (local_version_param = current_version OR local_version_param IS NULL OR local_version_param = 0)
          AND (
            name != name_param OR
            description != description_param OR
            link_type_id != link_type_id_value OR
            start_item_type_id != start_item_type_id_value OR
            end_item_type_id != end_item_type_id_value
          );
        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        SELECT ox_get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
      END IF;
      RETURN QUERY SELECT result;
    END;
    $BODY$;

    ALTER FUNCTION ox_set_link_rule(
      character varying,
      character varying,
      text,
      character varying,
      character varying,
      character varying,
      bigint,
      character varying,
      character varying[] -- role_key_param
    )
      OWNER TO onix;

    /*
      ox_set_privilege()
    */
    CREATE OR REPLACE FUNCTION ox_set_privilege(
      key_param character varying,
      partition_key_param character varying,
      role_key_param character varying,
      can_create_param boolean,
      can_read_param boolean,
      can_delete_param boolean,
      local_version_param bigint,
      changed_by_param character varying,
      logged_role_key_param character varying[]
    )
      RETURNS TABLE(result char(1))
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
    AS
    $BODY$
    DECLARE
      result             char(1); -- the result status for the upsert
      role_id_value      bigint;
      partition_id_value bigint;
      current_version    bigint;
      role_owner         character varying;
      partition_owner    character varying;
      logged_role_level  integer;
      rows_affected      integer;
    BEGIN
      -- finds the level of the logged role
      SELECT r.level
      FROM role r
      WHERE r.key = ANY(logged_role_key_param)
      ORDER BY r.level DESC
      LIMIT 1
        INTO logged_role_level;

      -- finds the owner of the role to add the privilege to
      SELECT r.owner, r.id
      FROM role r
      WHERE r.key = role_key_param
        INTO role_owner, role_id_value;

      -- finds the owner of the partition to add the privilege to
      SELECT p.owner, p.id
      FROM partition p
      WHERE p.key = partition_key_param
        INTO partition_owner, partition_id_value;

      IF (logged_role_level = 0) THEN
        -- logged role cannot mess with privileges
        RAISE EXCEPTION 'Role level %: "%" is not authorised to set privilege.', logged_role_level, logged_role_key_param;
      ELSEIF (logged_role_level = 1) THEN
        IF NOT(role_owner = ANY(logged_role_key_param) AND partition_owner = ANY(logged_role_key_param)) THEN
          -- logged role can only add privileges if it owns both the role and partition, so cannot do it in this case
          RAISE EXCEPTION 'Role level %: "%" is not authorised to set privilege because it does not own privilege or role to add the privilege to. Role owner is "%" and Partition owner is "%".', logged_role_level, logged_role_key_param, role_owner, partition_owner;
        END IF;
      END IF;

      -- get the privilege current version
      SELECT version FROM privilege WHERE key = key_param INTO current_version;

      -- if no version is found then go for an insert
      IF (current_version IS NULL) THEN
          -- logged role is either level 1 owning role and partition or level 2
          INSERT INTO privilege(
            key,
            partition_id,
            role_id,
            can_create,
            can_read,
            can_delete,
            version,
            created,
            updated,
            changed_by
          )
          VALUES(
            key_param,
            partition_id_value,
            role_id_value,
            can_create_param,
            can_read_param,
            can_delete_param,
            1,
            current_timestamp,
            null,
            changed_by_param
          );
          result := 'I';
      ELSE
          UPDATE privilege
          SET can_create   = can_create_param,
              can_read     = can_read_param,
              can_delete   = can_delete_param,
              version      = version + 1,
              updated      = current_timestamp,
              changed_by   = changed_by_param
          WHERE key = key_param
            -- concurrency management - optimistic locking
            AND (local_version_param = current_version OR local_version_param IS NULL OR local_version_param = 0)
            AND (
                  can_create != can_create_param OR
                  can_read != can_read_param OR
                  can_delete != can_delete_param OR
                  role_id != role_id_value OR
                  partition_id != partition_id_value
              );
          GET DIAGNOSTICS rows_affected := ROW_COUNT;
          SELECT ox_get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
      END IF;
      RETURN QUERY SELECT result;
    END;
    $BODY$;

    ALTER FUNCTION ox_set_privilege(
        character varying,
        character varying,
        character varying,
        boolean,
        boolean,
        boolean,
        bigint,
        character varying,
        character varying[]
    )
    OWNER TO onix;

    /*
      ox_set_type_attribute(...)
      Inserts a new or updates an existing type attribute.
      Concurrency Management:
       - If the item type is found in the database, the function attempts an update of the existing record.
          In this case, if a null regex is passed as local_version_param, no optimistic locking is performed.
          If a regex is specified for local_version_param, the update is only performed if and only if the version in the database matches the passed in version.
       - If the item type is not found in the database, then the local_version_param is ignored and a record with version 1 is inserted.
     */
    CREATE OR REPLACE FUNCTION ox_set_type_attribute(
        key_param character varying,
        name_param character varying,
        description_param text,
        type_param character varying,
        def_value_param character varying,
        required_param boolean,
        regex_param character varying,
        item_type_key_param character varying,
        link_type_key_param character varying,
        local_version_param bigint,
        changed_by_param character varying,
        role_key_param character varying[]
    )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
    AS
    $BODY$
    DECLARE
        result             char(1); -- the result status for the upsert
        current_version    bigint; -- the version of the row before the update or null if no row
        rows_affected      integer;
        model_id_value     integer;
        partition_id_value bigint;
        link_type_id_value bigint;
        item_type_id_value bigint;
    BEGIN
        -- gets the model id linked to the item_type
        IF (item_type_key_param IS NOT NULL) THEN
            SELECT m.id, it2.id INTO model_id_value, item_type_id_value
            FROM model m
                INNER JOIN item_type it2 on m.id = it2.model_id
            WHERE it2.key = item_type_key_param;
        ELSIF (link_type_key_param IS NOT NULL) THEN
            -- the type param could be associated with a link type instead of an item type
            -- gets the model id linked to the link_type
            SELECT m.id, lt2.id INTO model_id_value, link_type_id_value
            FROM model m
                 INNER JOIN link_type lt2 on m.id = lt2.model_id
            WHERE lt2.key = link_type_key_param;
        END IF;

        IF (model_id_value IS NULL) THEN
            RAISE EXCEPTION 'Missing item and link type in the definition of the attribute.'
                USING hint = 'Check either an item type or a link type is associated with the type attribute.';
        END IF;

        -- finds the partition associated with the model
        -- for the item type that has create rights for the specified role
        SELECT p.id
        FROM partition p
                 INNER JOIN model m on p.id = m.partition_id
                 INNER JOIN privilege pr on p.id = pr.partition_id
                 INNER JOIN role r on pr.role_id = r.id
            AND pr.can_create = TRUE -- has create permission
            AND r.key = ANY(role_key_param) -- the user role
            AND m.id = model_id_value -- the model the item type is in
        LIMIT 1
        INTO partition_id_value;

        IF (partition_id_value IS NULL) THEN
            RAISE EXCEPTION 'Role % is not authorised to create Type Attribute %.', role_key_param, key_param
                USING hint = 'The role needs to be granted CREATE privilege or a new role should be used instead.';
        END IF;

        -- gets the current type attribute version
        SELECT version FROM type_attribute WHERE key = key_param INTO current_version;

        IF (current_version IS NULL) THEN
            INSERT INTO type_attribute (
                id,
                key,
                name,
                description,
                type,
                def_value,
                required,
                regex,
                item_type_id,
                link_type_id,
                version,
                created,
                updated,
                changed_by
            )
            VALUES (
                nextval('type_attribute_id_seq'),
                key_param,
                name_param,
                description_param,
                type_param,
                def_value_param,
                required_param,
                regex_param,
                item_type_id_value,
                link_type_id_value,
                1,
                current_timestamp,
                null,
                changed_by_param
           );
            result := 'I';
        ELSE
            UPDATE type_attribute
            SET name        = name_param,
                description = description_param,
                type = type_param,
                def_value = def_value_param,
                required = required_param,
                regex = regex_param,
                item_type_id = item_type_id_value,
                link_type_id = link_type_id_value,
                version     = version + 1,
                updated     = current_timestamp,
                changed_by  = changed_by_param
            WHERE key = key_param
              -- concurrency management - optimistic locking
              AND (local_version_param = current_version OR local_version_param IS NULL OR local_version_param = 0)
              AND (
                    name != name_param OR
                    description != description_param OR
                    type != type_param OR
                    def_value != def_value_param OR
                    required != required_param OR
                    regex != regex_param OR
                    item_type_id != item_type_id_value OR
                    link_type_id != link_type_id_value
                );
            GET DIAGNOSTICS rows_affected := ROW_COUNT;
            SELECT ox_get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
        END IF;
        RETURN QUERY SELECT result;
    END;
    $BODY$;

    ALTER FUNCTION ox_set_type_attribute(
        character varying, -- key_param
        character varying, -- name_param
        text, -- description_param
        character varying, -- type_param
        character varying, -- def_value_param
        boolean, -- required_param
        character varying, -- regex_param
        character varying, -- item_type_key_param
        character varying, -- link_type_key_param
        bigint, -- local_version_param
        character varying, -- changed_by_param
        character varying[] -- role_key_param
        )
        OWNER TO onix;

  END
  $$;