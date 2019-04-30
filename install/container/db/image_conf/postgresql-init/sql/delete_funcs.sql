/*
    Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

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

      /*
      delete_model
     */
      CREATE OR REPLACE FUNCTION delete_model(
          key_param character varying,
          force boolean,
          role_key_param character varying
        )
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      DECLARE
        partition_id_value bigint;
      BEGIN
        -- finds the partition associated with the model
        -- for the item type that has create rights for the specified role
        SELECT p.id
        FROM partition p
        INNER JOIN model m on p.id = m.partition_id
        INNER JOIN privilege pr on p.id = pr.partition_id
        INNER JOIN role r on pr.role_id = r.id
          AND pr.can_delete = TRUE -- has create permission
          AND r.key = role_key_param -- the user role
          AND m.key = key_param -- the model
             INTO partition_id_value;

        IF (partition_id_value IS NULL) THEN
          RAISE EXCEPTION 'Role % is not authorised to delete Model.', role_key_param
            USING hint = 'The role needs to be granted DELETE privilege or a new role should be used instead.';
        END IF;

        IF (force = TRUE) THEN
          DELETE
          FROM item_type it USING model m
          WHERE m.id = it.model_id;

          DELETE
          FROM link_type lt USING model m
          WHERE m.id = lt.model_id;
        END IF;

        DELETE
        FROM model
        WHERE key = key_param;
      END
      $BODY$;

      ALTER FUNCTION delete_model(character varying, boolean, character varying)
        OWNER TO onix;

      /*
        delete_item
       */
      CREATE OR REPLACE FUNCTION delete_item(key_param character varying)
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN

        DELETE
        FROM item
        WHERE key = key_param;
      END
      $BODY$;

      ALTER FUNCTION delete_item(character varying)
        OWNER TO onix;

      /*
        delete_all_items
       */
      CREATE OR REPLACE FUNCTION delete_all_items()
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN
        DELETE FROM link_rule;
        DELETE FROM tag;
        DELETE FROM link;
        DELETE FROM item;
      END
      $BODY$;

      ALTER FUNCTION delete_all_items()
        OWNER TO onix;

      /*
        delete_item_type
       */
      CREATE OR REPLACE FUNCTION delete_item_type(
        key_param character varying,
        force boolean,
        role_key_param character varying
      )
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      DECLARE
        partition_id_value bigint;
        item_type_id_value bigint;
      BEGIN
        -- return if the item_type does not exist
        SELECT id FROM item_type WHERE key = key_param INTO item_type_id_value;
        IF (item_type_id_value IS NULL) THEN
          RETURN;
        END IF;

        -- finds the partition associated with the model
        -- for the item type that has create rights for the specified role
        SELECT p.id
        FROM partition p
        INNER JOIN model m on p.id = m.partition_id
        INNER JOIN privilege pr on p.id = pr.partition_id
        INNER JOIN role r on pr.role_id = r.id
        INNER JOIN item_type it on m.id = it.model_id
          AND pr.can_delete = TRUE -- has create permission
          AND r.key = role_key_param -- the user role
          AND it.key = key_param -- the item type
             INTO partition_id_value;

        IF (partition_id_value IS NULL) THEN
          RAISE EXCEPTION 'Role % is not authorised to delete Item Type %.', role_key_param, key_param
            USING hint = 'The role needs to be granted DELETE privilege or a new role should be used instead.';
        END IF;

        IF (force = TRUE) THEN
          -- if forcing then it removes all items of this item type
          DELETE
          FROM item
          WHERE item_type_id IN (
            SELECT id
            FROM item_type
            WHERE key = key_param
          );

          DELETE
          FROM link_rule r USING item_type it
          WHERE r.start_item_type_id = it.id;

          DELETE
          FROM link_rule r USING item_type it
          WHERE r.end_item_type_id = it.id;
        END IF;

        DELETE
        FROM item_type
        WHERE key = key_param;
      END;
      $BODY$;

      ALTER FUNCTION delete_item_type(character varying, boolean, character varying)
        OWNER TO onix;

      /*
        delete_link
       */
      CREATE OR REPLACE FUNCTION delete_link(key_param character varying)
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN
        DELETE
        FROM link
        WHERE key = key_param;
      END
      $BODY$;

      ALTER FUNCTION delete_link(character varying)
        OWNER TO onix;

      /*
        delete_link_type
       */
      CREATE OR REPLACE FUNCTION delete_link_type(
        key_param character varying,
        force boolean,
        role_key_param character varying
      )
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN
        IF (force = TRUE) THEN
          -- if forcing then it removes all links of this link type
          DELETE
          FROM link l
          USING link_type lt, model m, partition p, privilege pr, role r
          WHERE lt.id = l.link_type_id
            AND lt.key = key_param
            AND r.key = role_key_param
            AND m.partition_id = p.id
            AND pr.partition_id = p.id
            AND pr.role_id = r.id
            AND pr.can_delete = TRUE;
        END IF;
        DELETE
        FROM link_type lt
        USING model m, partition p, privilege pr, role r
        WHERE lt.key = key_param
          AND r.key = role_key_param
          AND m.partition_id = p.id
          AND pr.partition_id = p.id
          AND pr.role_id = r.id
          AND pr.can_delete = TRUE;
      END
      $BODY$;

      ALTER FUNCTION delete_link_type(character varying, boolean, character varying)
        OWNER TO onix;

      /*
        clear_all: deletes all instance data
       */
      CREATE OR REPLACE FUNCTION clear_all(
        role_key_param character varying
      )
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN
        DELETE FROM tag;
        DELETE FROM link_rule;
        DELETE FROM link;
        DELETE FROM item;
        PERFORM delete_link_types(role_key_param);
        PERFORM delete_item_types(role_key_param);
        PERFORM delete_link_rules();
      END
      $BODY$;

      ALTER FUNCTION clear_all(character varying)
        OWNER TO onix;

      /*
        delete_item_types: deletes all item types
       */
      CREATE OR REPLACE FUNCTION delete_item_types(
        role_key_param character varying
      )
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN
        DELETE FROM item_type it
          USING partition p, privilege pr, role r, model m
          WHERE it.model_id = m.id
          AND m.partition_id = p.id
          AND pr.partition_id = p.id
          AND pr.can_delete = TRUE
          AND pr.role_id = r.id
          AND r.key = role_key_param;
      END
      $BODY$;

      ALTER FUNCTION delete_item_types(character varying)
        OWNER TO onix;

      /*
        delete_link_types: deletes all link types
       */
      CREATE OR REPLACE FUNCTION delete_link_types(
        role_key_param character varying
      )
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN
        DELETE FROM link_type lt
        USING partition p, privilege pr, role r, model m
        WHERE lt.model_id = m.id
          AND m.partition_id = p.id
          AND pr.partition_id = p.id
          AND pr.can_delete = TRUE
          AND pr.role_id = r.id
          AND r.key = role_key_param;
      END
      $BODY$;

      ALTER FUNCTION delete_link_types(character varying)
        OWNER TO onix;

      /*
        delete_link_rules: deletes all link rules
       */
      CREATE OR REPLACE FUNCTION delete_link_rules()
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN
        DELETE FROM link_rule;
      END
      $BODY$;

      ALTER FUNCTION delete_link_rules()
        OWNER TO onix;

    END
    $$;