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

     /*
      ox_delete_partition
     */
      CREATE OR REPLACE FUNCTION ox_delete_partition(
        key_param character varying,
        role_key_param character varying[]
      )
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
      BEGIN
        -- checks the role can modify this role
        PERFORM ox_can_manage_partition(role_key_param);

        DELETE
        FROM partition p
        WHERE p.key = key_param;
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
       RETURNS VOID
       LANGUAGE 'plpgsql'
       COST 100
       VOLATILE
     AS
     $BODY$
     BEGIN
       -- checks the role can modify this role
       PERFORM ox_can_manage_partition(role_key_param);

       DELETE
       FROM role r
       WHERE r.key = key_param;
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
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
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
      END
      $BODY$;

      ALTER FUNCTION ox_delete_model(character varying, character varying[])
        OWNER TO onix;

      /*
        ox_delete_item
       */
      CREATE OR REPLACE FUNCTION ox_delete_item(
        key_param character varying,
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
        WHERE i.key = key_param
          AND p.id = i.partition_id
          AND pr.can_delete = TRUE
          AND r.key = ANY(role_key_param);
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
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
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
      END;
      $BODY$;

      ALTER FUNCTION ox_delete_item_type(character varying, character varying[])
        OWNER TO onix;

      /*
        ox_delete_link
       */
      CREATE OR REPLACE FUNCTION ox_delete_link(
        key_param character varying,
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
        RETURNS VOID
        LANGUAGE 'plpgsql'
        COST 100
        VOLATILE
      AS
      $BODY$
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
        PERFORM ox_delete_link_types(role_key_param);
        PERFORM ox_delete_item_types(role_key_param);
        PERFORM ox_delete_link_rules(role_key_param);
      END
      $BODY$;

      ALTER FUNCTION ox_clear_all(character varying[])
        OWNER TO onix;

      /*
        ox_delete_item_types: deletes all item types
       */
      CREATE OR REPLACE FUNCTION ox_delete_item_types(
        role_key_param character varying[]
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
          AND r.key = ANY(role_key_param);
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
          AND r.key = ANY(role_key_param);
      END
      $BODY$;

      ALTER FUNCTION ox_delete_link_types(character varying[])
        OWNER TO onix;

      /*
        ox_delete_link_rules: deletes all link rules
       */
      CREATE OR REPLACE FUNCTION ox_delete_link_rules(
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
        FROM link_rule lr
        USING link_type lt, model m, partition p, privilege pr, role r
          WHERE lr.link_type_id = lt.id
          AND lt.model_id = m.id
          AND m.partition_id = p.id
          AND pr.partition_id = p.id
          AND r.id = pr.role_id
          AND pr.can_delete = TRUE
          AND r.key = ANY(role_key_param);
      END
      $BODY$;

      ALTER FUNCTION ox_delete_link_rules(
        character varying[] -- role_key_param
      )
      OWNER TO onix;

    END
    $$;