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
        partition_key_value character varying;
      BEGIN
        -- gets the partition the model is in
        SELECT p.key
        FROM model m
        INNER JOIN partition p
           ON p.id = m.partition_id
        WHERE m.key = key_param INTO partition_key_value;

        -- check the role have delete rights on the partition, if not the raise exception
        PERFORM check_partition_privilege(role_key_param, partition_key_value, false, false, false, true);

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
      CREATE OR REPLACE FUNCTION delete_item_type(key_param character varying, force boolean default false)
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN
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
      END
      $BODY$;

      ALTER FUNCTION delete_item_type(character varying, boolean)
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
      CREATE OR REPLACE FUNCTION delete_link_type(key_param character varying, force boolean default false)
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
          FROM link
          WHERE link_type_id IN (
            SELECT id
            FROM link_type
            WHERE key = key_param
          );
        END IF;
        DELETE
        FROM link_type
        WHERE key = key_param;
      END
      $BODY$;

      ALTER FUNCTION delete_link_type(character varying, boolean)
        OWNER TO onix;

      /*
        clear_all: deletes all instance data
       */
      CREATE OR REPLACE FUNCTION clear_all()
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
        PERFORM delete_link_types();
        PERFORM delete_item_types();
        PERFORM delete_link_rules();
      END
      $BODY$;

      ALTER FUNCTION clear_all()
        OWNER TO onix;

      /*
        delete_item_types: deletes all item types
       */
      CREATE OR REPLACE FUNCTION delete_item_types()
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN
        DELETE FROM item_type;
      END
      $BODY$;

      ALTER FUNCTION delete_item_types()
        OWNER TO onix;

      /*
        delete_link_types: deletes all link types
       */
      CREATE OR REPLACE FUNCTION delete_link_types()
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN
        DELETE FROM link_type;
      END
      $BODY$;

      ALTER FUNCTION delete_link_types()
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