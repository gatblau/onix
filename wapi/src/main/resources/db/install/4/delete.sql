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
DO
  $$
    BEGIN

    CREATE OR REPLACE FUNCTION ox_get_delete_result(
      row_count int
    )
    RETURNS TABLE(result char(1))
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
    AS
    $BODY$
    DECLARE
      result char(1);
    BEGIN
      IF row_count > 0 THEN
        result := 'D';
      ELSE
        result := 'N';
      END IF;
      RETURN QUERY SELECT result;
    END
    $BODY$;

    ALTER FUNCTION ox_get_delete_result(int)
      OWNER TO onix;

     /*
      ox_delete_partition
     */
      CREATE OR REPLACE FUNCTION ox_delete_partition(
        key_param character varying,
        role_key_param character varying[]
      )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      DECLARE
        rows_affected INTEGER;
      BEGIN
        -- checks the role can modify this role
        PERFORM ox_can_manage_partition(role_key_param);

        DELETE
        FROM partition p
        WHERE p.key = key_param;

        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
      END
      $BODY$;

      ALTER FUNCTION ox_delete_partition(character varying, character varying[])
        OWNER TO onix;

     /*
       ox_delete_role
      */
     CREATE OR REPLACE FUNCTION ox_delete_role(
       key_param character varying,
       role_key_param character varying[]
     )
       RETURNS TABLE(result char(1))
       LANGUAGE 'plpgsql'
       COST 100
       VOLATILE
     AS
     $BODY$
     DECLARE
       rows_affected INTEGER;
     BEGIN
       -- checks the role can modify this role
       PERFORM ox_can_manage_partition(role_key_param);

       DELETE
       FROM role r
       WHERE r.key = key_param;

       GET DIAGNOSTICS rows_affected := ROW_COUNT;
       RETURN QUERY SELECT ox_get_delete_result(rows_affected);
     END
     $BODY$;

     ALTER FUNCTION ox_delete_role(character varying, character varying[])
       OWNER TO onix;

      /*
      ox_delete_model
     */
      CREATE OR REPLACE FUNCTION ox_delete_model(
          key_param character varying,
          role_key_param character varying[]
        )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      DECLARE
        rows_affected INTEGER;
      BEGIN
        DELETE
        FROM model m
          USING partition p, privilege pr, role r
          WHERE m.key = key_param
          AND m.partition_id = p.id
          AND pr.partition_id = p.id
          AND pr.can_delete = TRUE
          AND pr.role_id = r.id
          AND r.key = ANY(role_key_param);

        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
      END
      $BODY$;

      ALTER FUNCTION ox_delete_model(character varying, character varying[])
        OWNER TO onix;

    /*
      ox_delete_user
     */
    CREATE OR REPLACE FUNCTION ox_delete_user(
        key_param character varying,
        role_key_param character varying[]
    )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
    AS
    $BODY$
    DECLARE
        rows_affected INTEGER;
    BEGIN
        -- only users in level 2 roles can delete other users
        -- if not super admin then raise exception
        PERFORM ox_is_super_admin(role_key_param, TRUE);

        DELETE FROM "user" u
        WHERE u.key = key_param;

        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
    END
    $BODY$;

    ALTER FUNCTION ox_delete_model(character varying, character varying[])
        OWNER TO onix;

    /*
      ox_delete_membership
     */
    CREATE OR REPLACE FUNCTION ox_delete_membership(
        key_param character varying,
        role_key_param character varying[]
    )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
    AS
    $BODY$
    DECLARE
        rows_affected INTEGER;
    BEGIN
        -- only users in level 2 roles can delete other users
        -- if not super admin then raise exception
        PERFORM ox_is_super_admin(role_key_param, TRUE);

        DELETE FROM membership m
        WHERE m.key = key_param;

        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
    END
    $BODY$;

    ALTER FUNCTION ox_delete_membership(key_param character varying, role_key_param character varying[])
        OWNER TO onix;

    /*
        ox_delete_item
       */
      CREATE OR REPLACE FUNCTION ox_delete_item(
        key_param character varying,
        role_key_param character varying[]
      )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      DECLARE
        rows_affected      integer;
      BEGIN
        DELETE
        FROM item i
        USING partition p, privilege pr, role r
        WHERE i.key = key_param
          AND p.id = i.partition_id
          AND pr.can_delete = TRUE
          AND r.key = ANY(role_key_param);

        GET DIAGNOSTICS rows_affected := ROW_COUNT;

        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
      END
      $BODY$;

      ALTER FUNCTION ox_delete_item(character varying, character varying[])
        OWNER TO onix;

      /*
        ox_delete_all_items
       */
      CREATE OR REPLACE FUNCTION ox_delete_all_items(
        role_key_param character varying[]
      )
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN
        DELETE
        FROM item i
        USING partition p, privilege pr, role r
        WHERE i.partition_id = p.id
        AND p.id = pr.partition_id
        AND pr.can_delete = TRUE
        AND pr.role_id = r.id
        AND r.key = ANY(role_key_param);
      END
      $BODY$;

      ALTER FUNCTION ox_delete_all_items(character varying[])
        OWNER TO onix;

      /*
        ox_delete_item_type
       */
      CREATE OR REPLACE FUNCTION ox_delete_item_type(
        key_param character varying,
        role_key_param character varying[]
      )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      DECLARE
        rows_affected INTEGER;
      BEGIN
        DELETE
        FROM item_type it
        USING model m, partition p, privilege pr, role r
        WHERE it.key = key_param
          AND it.model_id = m.id
          AND m.partition_id = p.id
          AND p.id = pr.partition_id
          AND pr.role_id = r.id
          AND pr.can_delete = TRUE
          AND r.key = ANY(role_key_param);

        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
      END;
      $BODY$;

      ALTER FUNCTION ox_delete_item_type(character varying, character varying[])
        OWNER TO onix;

    /*
        ox_delete_item_type_attribute
    */
    CREATE OR REPLACE FUNCTION ox_delete_item_type_attribute(
        item_type_key_param character varying,
        type_attr_key_param character varying,
        role_key_param character varying[]
    )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
    AS
    $BODY$
    DECLARE
        rows_affected INTEGER;
    BEGIN
        DELETE
        FROM type_attribute ta
            USING model m, partition p, privilege pr, role r, item_type it
        WHERE ta.key = type_attr_key_param
          AND it.key = item_type_key_param
          AND ta.item_type_id = it.id
          AND it.model_id = m.id
          AND m.partition_id = p.id
          AND p.id = pr.partition_id
          AND pr.role_id = r.id
          AND pr.can_delete = TRUE
          AND r.key = ANY(role_key_param);
        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
    END;
    $BODY$;

    ALTER FUNCTION ox_delete_item_type_attribute(character varying, character varying, character varying[])
        OWNER TO onix;

    /*
        ox_delete_link_type_attribute
    */
    CREATE OR REPLACE FUNCTION ox_delete_link_type_attribute(
        link_type_key_param character varying,
        type_attr_key_param character varying,
        role_key_param character varying[]
    )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
    AS
    $BODY$
    DECLARE
        rows_affected INTEGER;
    BEGIN
        DELETE
        FROM type_attribute ta
            USING model m, partition p, privilege pr, role r, link_type lt
        WHERE ta.key = type_attr_key_param
          AND lt.key = link_type_key_param
          AND ta.link_type_id = lt.id
          AND lt.model_id = m.id
          AND m.partition_id = p.id
          AND p.id = pr.partition_id
          AND pr.role_id = r.id
          AND pr.can_delete = TRUE
          AND r.key = ANY(role_key_param);
        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
    END;
    $BODY$;

    ALTER FUNCTION ox_delete_link_type_attribute(character varying, character varying, character varying[])
        OWNER TO onix;

      /*
        ox_delete_link
       */
      CREATE OR REPLACE FUNCTION ox_delete_link(
        key_param character varying,
        role_key_param character varying[]
      )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      DECLARE
        rows_affected INTEGER;
      BEGIN
        DELETE
        FROM link l
        USING link_type lt, model m, partition p, privilege pr, role r
        WHERE l.key = key_param
          AND lt.id = l.link_type_id
          AND lt.model_id = m.id
          AND m.partition_id = p.id
          AND pr.partition_id = p.id
          AND pr.role_id = r.id
          AND r.key = ANY(role_key_param)
          AND pr.can_delete = TRUE;

        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
      END
      $BODY$;

      ALTER FUNCTION ox_delete_link(
        character varying,
        character varying[] -- role_key_param
      )
        OWNER TO onix;

      /*
        ox_delete_link_type
       */
      CREATE OR REPLACE FUNCTION ox_delete_link_type(
        key_param character varying,
        role_key_param character varying[]
      )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      DECLARE
        rows_affected INTEGER;
      BEGIN
        DELETE
        FROM link_type lt
        USING model m, partition p, privilege pr, role r
        WHERE lt.key = key_param
          AND r.key = ANY(role_key_param)
          AND m.partition_id = p.id
          AND pr.partition_id = p.id
          AND pr.role_id = r.id
          AND pr.can_delete = TRUE;

        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
      END
      $BODY$;

      ALTER FUNCTION ox_delete_link_type(
        character varying,
        character varying[]
      )
        OWNER TO onix;

      /*
        ox_clear_all: deletes all instance data
       */
      CREATE OR REPLACE FUNCTION ox_clear_all(
        role_key_param character varying[]
      )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      DECLARE
        rows_affected INTEGER;
      BEGIN
        DELETE FROM tag;
        PERFORM ox_delete_models(role_key_param);
        PERFORM ox_delete_link_types(role_key_param);
        PERFORM ox_delete_item_types(role_key_param);
        PERFORM ox_delete_link_rules(role_key_param);
        PERFORM ox_delete_all_items(role_key_param);
        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
      END
      $BODY$;

      ALTER FUNCTION ox_clear_all(character varying[])
        OWNER TO onix;

    /*
        ox_delete_models: deletes all models
       */
    CREATE OR REPLACE FUNCTION ox_delete_models(
        role_key_param character varying[]
    )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
    AS
    $BODY$
    DECLARE
        rows_affected INTEGER;
    BEGIN
        DELETE FROM model m
            USING partition p, privilege pr, role r
        WHERE m.partition_id = p.id
          AND pr.partition_id = p.id
          AND pr.can_delete = TRUE
          AND pr.role_id = r.id
          AND r.key = ANY(role_key_param);

        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
    END
    $BODY$;

    ALTER FUNCTION ox_delete_models(character varying[])
        OWNER TO onix;

    /*
        ox_delete_item_types: deletes all item types
       */
      CREATE OR REPLACE FUNCTION ox_delete_item_types(
        role_key_param character varying[]
      )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      DECLARE
        rows_affected INTEGER;
      BEGIN
        DELETE FROM item_type it
          USING partition p, privilege pr, role r, model m
          WHERE it.model_id = m.id
          AND m.partition_id = p.id
          AND pr.partition_id = p.id
          AND pr.can_delete = TRUE
          AND pr.role_id = r.id
          AND r.key = ANY(role_key_param);

        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
      END
      $BODY$;

      ALTER FUNCTION ox_delete_item_types(character varying[])
        OWNER TO onix;

      /*
        ox_delete_link_types: deletes all link types
       */
      CREATE OR REPLACE FUNCTION ox_delete_link_types(
        role_key_param character varying[]
      )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      DECLARE
        rows_affected INTEGER;
      BEGIN
        DELETE FROM link_type lt
        USING partition p, privilege pr, role r, model m
        WHERE lt.model_id = m.id
          AND m.partition_id = p.id
          AND pr.partition_id = p.id
          AND pr.can_delete = TRUE
          AND pr.role_id = r.id
          AND r.key = ANY(role_key_param);

        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
      END
      $BODY$;

      ALTER FUNCTION ox_delete_link_types(character varying[])
        OWNER TO onix;

    /*
        ox_delete_link_rule: deletes the specified link rule
       */
    CREATE OR REPLACE FUNCTION ox_delete_link_rule(
        key_param character varying,
        role_key_param character varying[]
    )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
    AS
    $BODY$
    DECLARE
        rows_affected INTEGER;
    BEGIN
        DELETE
        FROM link_rule lr
            USING link_type lt, model m, partition p, privilege pr, role r
        WHERE lr.link_type_id = lt.id
          AND lt.model_id = m.id
          AND m.partition_id = p.id
          AND pr.partition_id = p.id
          AND r.id = pr.role_id
          AND pr.can_delete = TRUE
          AND lr.key = key_param
          AND r.key = ANY(role_key_param);

        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
    END
    $BODY$;

    ALTER FUNCTION ox_delete_link_rule(
        character varying, -- key_param
        character varying[] -- role_key_param
        )
        OWNER TO onix;

    /*
        ox_delete_link_rules: deletes all link rules
       */
      CREATE OR REPLACE FUNCTION ox_delete_link_rules(
        role_key_param character varying[]
      )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      DECLARE
        rows_affected INTEGER;
      BEGIN
        DELETE
        FROM link_rule lr
        USING link_type lt, model m, partition p, privilege pr, role r
          WHERE lr.link_type_id = lt.id
          AND lt.model_id = m.id
          AND m.partition_id = p.id
          AND pr.partition_id = p.id
          AND r.id = pr.role_id
          AND pr.can_delete = TRUE
          AND r.key = ANY(role_key_param);

        GET DIAGNOSTICS rows_affected := ROW_COUNT;
        RETURN QUERY SELECT ox_get_delete_result(rows_affected);
      END
      $BODY$;

      ALTER FUNCTION ox_delete_link_rules(
        character varying[] -- role_key_param
      )
      OWNER TO onix;

    /*
     ox_delete_privilege()
    */
    CREATE OR REPLACE FUNCTION ox_delete_privilege(
        key_param character varying,
        logged_role_key_param character varying[]
    )
        RETURNS TABLE(result char(1))
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
    AS
    $BODY$
    DECLARE
        role_id_value      bigint;
        partition_id_value bigint;
        role_owner         character varying;
        partition_owner    character varying;
        logged_role_level  integer;
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
        FROM privilege p
            INNER JOIN role r ON r.id = p.role_id
        WHERE p.key = key_param
            INTO role_owner, role_id_value;

        -- fins the owner of the partition to add the privilege to
        SELECT p.owner, p.id
        FROM privilege pr
            INNER JOIN partition p ON p.id = pr.partition_id
        WHERE pr.key = key_param
            INTO partition_owner, partition_id_value;

        IF (logged_role_level = 0) THEN
            -- logged role cannot mess with privileges
            RAISE EXCEPTION 'Role level %: "%" is not authorised to remove privilege.', logged_role_level, logged_role_key_param;
        ELSEIF (logged_role_level = 1) THEN
            IF NOT(role_owner = ANY(logged_role_key_param) AND partition_owner = ANY(logged_role_key_param)) THEN
                -- logged role can only remove privileges if it owns both the role and partition, so cannot do it in this case
                RAISE EXCEPTION 'Role level %: "%" is not authorised to remove privilege because it does not own privilege or role to add the privilege to. Role owner is "%" and Partition owner is "%".', logged_role_level, logged_role_key_param, role_owner, partition_owner;
            END IF;
        END IF;

        -- logged role is either level 1 owning role and partition or level 2
        DELETE FROM privilege p
        WHERE p.partition_id = partition_id_value
          AND p.role_id = role_id_value
          AND p.key = key_param;

        RETURN QUERY SELECT 'D'::char(1);
    END;
    $BODY$;

    ALTER FUNCTION ox_delete_privilege(
        character varying,
        character varying[]
        )
        OWNER TO onix;
    END
    $$;